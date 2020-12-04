package bid

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/zzsds/trade/queue"
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
			case buf := <-bid.Buffer():
				message := buf.(*BufferMessage)
				t.Log(message.Queue.Name(), message.Price)
			}
		}
	}()
}

func TestAdd(t *testing.T) {
	// t.Run("TestBuffer", TestBuffer)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		bid.Buy().Push(queue.NewData(&Unit{
			Name:    "xlj-" + strconv.Itoa(i),
			Price:   price / 3,
			Number:  1,
			UID:     i,
			TradeID: i,
		}))
		// bid.Add(bid.Buy(), &Unit{
		// 	Name:    "xlj-" + strconv.Itoa(i),
		// 	Price:   price / 3.5,
		// 	Number:  1,
		// 	UID:     i,
		// 	TradeID: i,
		// })
	}

	for i := 0; i < 10; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		bid.Sell().Unshift(queue.NewData(&Unit{
			Name:    "wj-" + strconv.Itoa(i),
			Price:   price / 2,
			Number:  2,
			UID:     i,
			TradeID: i,
		}))
		// bid.Add(bid.Sell(), &Unit{
		// 	Name:    "wj-" + strconv.Itoa(i),
		// 	Price:   price / 2.5,
		// 	Number:  2,
		// 	UID:     i,
		// 	TradeID: i,
		// })
	}

	for n := bid.Buy().Front(); n != nil; n = n.Next() {
		t.Log(n.Data.Content)
	}

	for n := bid.Sell().Front(); n != nil; n = n.Next() {
		t.Log(n.Data.Content)
	}

	// select {
	// case timeout := <-time.After(1 * time.Second):
	// 	t.Fatal(timeout)
	// }
}
