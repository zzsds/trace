package match

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/zzsds/trade/bid"
)

var matchup Server

func TestMain(t *testing.M) {
	matchup = NewMatch(Name("test"))
	t.Run()
}

func TestRun(t *testing.T) {
	b := bid.NewBid(bid.Name("product"))

	for i := 0; i < 500; i++ {
		traceType := b.Buy()
		if i%2 != 0 {
			traceType = b.Sell()
		}
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		b.Add(traceType, &bid.Unit{
			Name:    "xlj",
			Amount:  int(rand.Intn(1000)),
			Price:   price,
			UID:     int(i),
			TradeID: int(i),
		})
	}

	go matchup.Bid(b).Run()
	select {
	case <-time.After(3 * time.Second):
	}
	t.Log("End")
}
