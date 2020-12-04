package bid

import (
	"math/rand"
	"strconv"
	"testing"
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

	for i := 0; i < 2; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		bid.Add(bid.Buy(), &Unit{
			Name:    "xlj-" + strconv.Itoa(i),
			Price:   price,
			Number:  int(rand.Intn(1000)),
			UID:     i,
			TradeID: i,
		})
	}

	// for _, v := range bid.Buy().List() {
	// 	t.Log(v.Content)
	// }

	// select {
	// case timeout := <-time.After(1 * time.Second):
	// 	t.Fatal(timeout)
	// }
}
