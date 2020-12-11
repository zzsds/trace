package bid

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var bid Server

func TestMain(t *testing.M) {
	bid = NewBid(Name("Product"))
	go func() {
		for {
			select {
			case <-bid.Buffer():
			}
		}
	}()
	t.Run()
}

func TestBuffer(t *testing.T) {
	go func() {
		for {
			select {
			case message := <-bid.Buffer():
				t.Log(message.Queue.Name(), message.Node.Data().Content)
			}
		}
	}()
}

func TestAdd(t *testing.T) {
	t.Run("TestBuffer", TestBuffer)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)

		traceType := bid.Buy()
		if i%2 != 0 {
			traceType = bid.Sell()
		}
		bid.Add(traceType, &Unit{
			Name:   "xlj-" + strconv.Itoa(i),
			Price:  price / 3.5,
			Amount: i,
			UID:    i,
			ID:     i,
		})
	}

	for n := bid.Buy().Front(); n != nil; n = n.Next() {
		t.Log(n.Data().Content)
	}
	t.Logf("buy length %d", bid.Buy().Len())

	for n := bid.Sell().Front(); n != nil; n = n.Next() {
		t.Log(n.Data().Content)
	}
	t.Logf("sell length %d", bid.Sell().Len())

	select {
	case timeout := <-time.After(3 * time.Second):
		t.Fatal(timeout)
	}
}

func BenchmarkAdd(t *testing.B) {
	for i := 0; i < t.N; i++ {
		traceType := bid.Buy()
		if i%2 != 0 {
			traceType = bid.Sell()
		}
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		bid.Add(traceType, &Unit{
			Name:   "xlj",
			Amount: int(rand.Intn(1000)),
			Price:  price / 3.5,
			UID:    int(i),
			ID:     int(i),
		})
	}
	t.Log(bid.Buy().Len(), bid.Sell().Len())
}
