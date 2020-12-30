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
				t.Logf("%s, %#v, %#v", result.Bid.Name(), result.Trigger, result.Trades)
			}
		}

	}()
}

func TestRun(t *testing.T) {
	m := matchup.Register(bid.NewBid(bid.Name("product")))
	go m.Run()
	b := m.Bid()
	go func() {
		for i := 0; i < 2000; i++ {

			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
			traceType := bid.Type_Buy
			if i%2 != 0 {
				traceType = bid.Type_Sell
			}

			data, _ := b.Add(bid.NewUnit(func(u *bid.Unit) {
				u.Name = "xlj"
				u.Amount = i + 1
				u.Price = price
				u.ID = int(i)
			}))
			t.Log(traceType.String(), data)
		}
	}()

	t.Log("截断")

	// for _, v := range b.Buy().List() {
	// 	t.Logf("Buy %#v", v.Value)
	// }
	// for _, v := range b.Sell().List() {
	// 	t.Logf("Sell %#v", v.Value)
	// }

	t.Run("TestBuffer", TestBuffer)
	t.Log("End")
	<-time.After(5 * time.Second)
}
