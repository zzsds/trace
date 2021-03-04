package match

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/zzsds/trade/bid"
)

var srv Server

func TestMain(t *testing.M) {
	srv = NewMatch(Name("test"))
	t.Run()
}

func TestQueue(t *testing.T) {
	go func() {
		for {
			select {
			case res := <-srv.Queue():
				result := res.(Result)
				t.Logf("%s, %#v, %#v", result.Bid.Name(), result.Trigger, result.Trades)
			}
		}

	}()
}

func TestRun(t *testing.T) {
	t.Run("TestQueue", TestQueue)
	b := srv.Bid()
	srv.Start()
	defer srv.Close()
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
	t.Log(srv.State())

	t.Log("End")
	<-time.After(5 * time.Second)
}
