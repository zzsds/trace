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
				_ = result
				t.Logf("%s, %#v, %#v", result.Bid.Name(), result.Trigger, result.Trades)
			}
		}

	}()
}

func TestRun(t *testing.T) {
	t.Run("TestBuffer", TestBuffer)
	m := matchup.Register(bid.NewBid(bid.Name("product")))
	b := m.Bid()
	m.Start()
	for i := 1; i < 200; i++ {

		price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
		traceType := bid.Type_Buy
		if i%2 != 0 {
			traceType = bid.Type_Sell
		}

		b.Add(bid.NewUnit(func(u *bid.Unit) {
			u.Type = traceType
			u.Name = "xlj"
			u.Amount = 1
			u.Price = price
			u.UID = rand.Intn(1000)
			u.ID = int(i)
		}))
	}

	t.Log("截断")

	// for _, v := range b.Buy().List() {
	// 	t.Logf("Buy %#v", v)
	// }
	// for _, v := range b.Sell().List() {
	// 	t.Logf("Sell %#v", v)
	// }
	t.Log("End")
	<-time.After(5 * time.Second)
}
