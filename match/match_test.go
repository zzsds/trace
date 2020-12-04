package match

import (
	"math/rand"
	"testing"

	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/queue"
)

var matchup Server

func TestMain(t *testing.M) {
	matchup = NewMatch(Name("test"))
	t.Run()
}

func TestRun(t *testing.T) {
	b := bid.NewBid(bid.Name("product"))
	buy := b.Buy()

	for i := 0; i < 100; i++ {
		buy.Push(queue.NewData(&bid.Unit{
			Name:    "xlj",
			Number:  int(rand.Intn(1000)),
			Price:   1.0,
			UID:     int(i),
			TradeID: int(i),
		}))
	}

	buy.Loop(func(n *queue.Node) error {
		t.Log(n.Data.Content)
		return nil
	})

	// matchup.Bid(bid).Run()
	t.Log("End")
}
