package bid

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	sync "sync"
	"testing"
	"time"
)

var srv Server

func TestMain(t *testing.M) {
	srv = NewBid(Name("Product"))

	t.Run()
}

func TestQueue(t *testing.T) {
	go func() {
		for {
			select {
			case message := <-srv.Queue():
				_ = message
				// t.Log(message.Queue.Name(), message.Node.Value)
			}
		}
	}()
}

func TestAdd(t *testing.T) {
	t.Run("TestQueue", TestQueue)
	rand.Seed(time.Now().Unix())

	for i := 1; i < 200; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		// price, _ := strconv.ParseFloat(strconv.Itoa(i), 64)
		traceType := Type_Buy
		if i%2 != 0 {
			traceType = Type_Sell
		}

		unit := NewUnit(func(u *Unit) {
			u.ID = i
			u.Type = traceType
			u.Name = "xlj-" + strconv.Itoa(i)
			u.Price = price
			u.Amount = 1
			u.UID = rand.Intn(10) + 1
		})
		srv.Add(unit)
	}
	// srv.Buy().Remove()

	var total float64
	for n := srv.Buy().Front(); n != nil; n = n.Next() {
		unit := n.Value.(*Unit)
		total += float64(unit.Amount)
		b, _ := json.MarshalIndent(n.Value, "", "\t")
		s := strings.Replace(string(b), "　", "", -1)
		t.Logf("%s", strings.Replace(s, "\n", "", -1))
	}
	t.Logf("buy length %d %.2f", srv.Buy().Len(), total)

	// srv.Remove(srv.Sell(), 1, 1)
	for n := srv.Sell().Front(); n != nil; n = n.Next() {
		b, _ := json.MarshalIndent(n.Value, "", "\t")
		s := strings.Replace(string(b), "　", "", -1)
		t.Logf("%s", strings.Replace(s, "\n", "", -1))
	}
	t.Logf("sell length %d", srv.Sell().Len())

	<-time.After(1 * time.Millisecond)
}

var (
	once sync.Once
)

func BenchmarkAdd(t *testing.B) {
	once.Do(func() {
		go func() {
			for {
				select {
				case message := <-srv.Queue():
					_ = message
					// t.Log(message.Queue.Name(), message.Node.Value)
				}
			}
		}()
	})
	for i := 0; i < t.N; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		traceType := Type_Buy
		if i%2 != 0 {
			traceType = Type_Sell
		}
		srv.Add(&Unit{
			Type:     traceType,
			CreateAt: time.Now().Unix(),
			Name:     "xlj",
			Amount:   int(rand.Intn(1000)),
			Price:    price / 3.5,
			UID:      int(i),
		})
	}

	t.Log(t.N, srv.Buy().Len(), srv.Sell().Len())
}

func BenchmarkAddParallel(t *testing.B) {
	t.ReportAllocs()
	once.Do(func() {
		go func() {
			for {
				select {
				case message := <-srv.Queue():
					_ = message
					// t.Log(message.Queue.Name(), message.Node.Value)
				}
			}
		}()
	})
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			// rand.Intn(1000)
			// strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
			srv.Add(&Unit{
				Type:     Type_Buy,
				CreateAt: time.Now().Unix(),
				Name:     "xlj",
				Amount:   int(rand.Intn(1000)),
				Price:    price / 3.5,
				UID:      rand.Intn(1000),
			})
		}
	})
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			// strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
			srv.Add(&Unit{
				Type:     Type_Sell,
				CreateAt: time.Now().Unix(),
				Name:     "wj",
				Amount:   int(rand.Intn(1000)),
				Price:    price / 3.5,
				UID:      rand.Intn(1000),
			})
		}
	})

	t.Log(t.N, srv.Buy().Len(), srv.Sell().Len())
}

func BenchmarkAddBid(t *testing.B) {
	t.ReportAllocs()
	once.Do(func() {
		go func() {
			for {
				select {
				case message := <-srv.Queue():
					_ = message
					// t.Log(message.Queue.Name(), message.Node.Value)
				}
			}
		}()
	})
	t.RunParallel(func(p *testing.PB) {
		// Each goroutine has its own bytes.Buffer.
		for p.Next() {
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
			srv.Add(&Unit{
				Type:     Type_Buy,
				CreateAt: time.Now().Unix(),
				Name:     "xlj",
				Amount:   int(rand.Intn(1000)),
				Price:    price / 3.5,
				ID:       int(rand.Intn(12321)),
				UID:      rand.Intn(100),
			})
		}
	})

	t.Log(t.N, srv.Buy().Len(), srv.Sell().Len())
}
