package bid

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var bid Server

func TestMain(t *testing.M) {
	bid = NewBid(Name("Product"))

	t.Run()
}

func TestBuffer(t *testing.T) {
	go func() {
		for {
			select {
			case message := <-bid.Buffer():
				_ = message
				// t.Log(message.Queue.Name(), message.Node.Value)
			}
		}
	}()
}

func TestAdd(t *testing.T) {
	t.Run("TestBuffer", TestBuffer)
	rand.Seed(time.Now().Unix())
	for i := 1; i < 200; i++ {
		// price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		// price, _ := strconv.ParseFloat(strconv.Itoa(i), 64)
		traceType := Type_Buy
		// if i%2 != 0 {
		// 	traceType = Type_Sell
		// }

		unit := NewUnit(func(u *Unit) {
			u.Type = traceType
			u.Name = "xlj-" + strconv.Itoa(i)
			u.Price = 1.2
			u.Amount = 1
			u.UID = rand.Intn(10) + 1
		})
		bid.Add(unit)
	}

	fmt.Println(bid.BuyData().Map, bid.BuyData().Array, bid.SellData().Array)

	var total float64
	for n := bid.Buy().Front(); n != nil; n = n.Next() {
		unit := n.Value.(*Unit)
		total += float64(unit.Amount)
		t.Logf("%v", n.Value)
	}
	t.Logf("buy length %d %.2f", bid.Buy().Len(), total)

	// for n := bid.Sell().Front(); n != nil; n = n.Next() {
	// 	t.Logf("%#v", n.Value)
	// }
	// t.Logf("sell length %d", bid.Sell().Len())

	<-time.After(1 * time.Millisecond)
}

func BenchmarkAdd(t *testing.B) {
	// go func() {
	// 	for {
	// 		select {
	// 		case msg := <-bid.Buffer():
	// 			_ = msg
	// 		}
	// 	}
	// }()
	for i := 0; i < t.N; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		traceType := Type_Buy
		if i%2 != 0 {
			traceType = Type_Sell
		}
		bid.Add(&Unit{
			Type:     traceType,
			CreateAt: time.Now(),
			Name:     "xlj",
			Amount:   int(rand.Intn(1000)),
			Price:    price / 3.5,
			UID:      int(i),
		})
	}

	t.Log(t.N, bid.Buy().Len(), bid.Sell().Len())
}

var (
	once sync.Once
	mu   *sync.RWMutex
	i    float64
)

func BenchmarkAddParallel(t *testing.B) {
	t.ReportAllocs()
	// once.Do(func() {
	// 	go func() {
	// 		for {
	// 			select {
	// 			case message := <-bid.Buffer():
	// 				_ = message
	// 				// t.Log(message.Queue.Name(), message.Node.Value)
	// 			}
	// 		}
	// 	}()
	// })
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			// rand.Intn(1000)
			// strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
			bid.Add(&Unit{
				Type:     Type_Buy,
				CreateAt: time.Now(),
				Name:     "xlj",
				Amount:   int(rand.Intn(1000)),
				Price:    price / 3.5,
				UID:      rand.Intn(1000),
			})
		}
	})
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
			bid.Add(&Unit{
				Type:     Type_Sell,
				CreateAt: time.Now(),
				Name:     "wj",
				Amount:   int(rand.Intn(1000)),
				Price:    price / 3.5,
				UID:      rand.Intn(1000),
			})
		}
	})
}

func BenchmarkAddBid(t *testing.B) {
	go func() {
		for {
			select {
			case msg := <-bid.Buffer():
				_ = msg
			}
		}
	}()
	t.ReportAllocs()
	t.RunParallel(func(p *testing.PB) {
		// Each goroutine has its own bytes.Buffer.
		for p.Next() {
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
			traceType := Type_Buy
			if rand.Intn(100)%2 == 0 {
				traceType = Type_Sell
			}
			bid.Add(&Unit{
				Type:     traceType,
				CreateAt: time.Now(),
				Name:     "xlj",
				Amount:   int(rand.Intn(1000)),
				Price:    price / 3.5,
				ID:       int(rand.Intn(12321)),
				UID:      rand.Intn(100),
			})
		}
	})

	t.Log(t.N, bid.Buy().Len(), bid.Sell().Len())
}

func BenchmarkTmplExucte(b *testing.B) {
	b.ReportAllocs()
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine has its own bytes.Buffer.
		var buf bytes.Buffer
		for pb.Next() {
			// The loop body is executed b.N times total across all goroutines.
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}
