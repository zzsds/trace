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

	for i := 0; i < 10; i++ {
		traceType := b.Buy()
		if i%2 != 0 {
			traceType = b.Sell()
		}
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(100)), 64)
		b.Add(traceType, &bid.Unit{
			Name:   "xlj",
			Amount: int(rand.Intn(1111)),
			Price:  price,
			UID:    int(i),
			ID:     int(i),
		})
	}

	match := matchup.Bid(b)
	go match.Run()
	select {
	case result := <-match.Buffer():
		t.Log(result.Bid.Name(), result.Trigger.String(), result.Trigger.Unit, result.Trades)
	case <-time.After(3 * time.Second):
	}

	for _, v := range b.Buy().List() {
		t.Log("Buy", v.Content)
	}
	for _, v := range b.Sell().List() {
		t.Log("Sell", v.Content)
	}

	t.Log("End")
}
