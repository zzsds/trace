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

func TestBuffer(t *testing.T) {
	go func() {
		for {
			select {
			case result := <-matchup.Buffer():
				t.Logf("%s, %s, %#v, %#v", result.Bid.Name(), result.Trigger.String(), result.Trigger.Unit, result.Trades)
			}
		}

	}()
}

func TestRun(t *testing.T) {
	b := bid.NewBid(bid.Name("product"))
	matchup.Register(b)
	go matchup.Run()

	for i := 0; i < 1000; i++ {
		traceType := b.Buy()
		if i%2 != 0 {
			traceType = b.Sell()
		}
		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
		data, _ := b.Add(traceType, &bid.Unit{
			Name:   "xlj",
			Amount: i + 1,
			Price:  price,
			UID:    int(i),
			ID:     int(i),
		})
		t.Log(traceType.Name(), data.Content)
	}

	t.Log("截断")
	<-time.After(1 * time.Second)

	// for _, v := range b.Buy().List() {
	// 	t.Logf("Buy %#v", v.Content)
	// }
	// for _, v := range b.Sell().List() {
	// 	t.Logf("Sell %#v", v.Content)
	// }

	t.Log("End")
	t.Run("TestBuffer", TestBuffer)
	<-time.After(1 * time.Second)
}
