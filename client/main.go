package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/match"
)

func main() {
	m := match.NewMatch(match.Name("BTC"))
	defer m.Close()
	if err := m.Start(); err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case res := <-m.Queue():
				msg := res.(match.Result)
				// log.Println(msg)
				fmt.Fprintf(os.Stderr, "撮合结果 %v, %s, %d, %.2f \n", msg.Trigger, msg.Trigger.Type.String(), msg.Amount, msg.Price)
			}
		}
	}()
	b := m.Bid()
	http.HandleFunc("/start", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("start request"))
		fmt.Fprintf(os.Stderr, "开始后当前撮合状态：%v, %t, %s  \n", m.Start(), m.State(), m.Name())
	})
	http.HandleFunc("/stop", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("stop request"))
		fmt.Fprintf(os.Stderr, "停止后当前撮合状态：%v, %t, %s  \n", m.Stop(), m.State(), m.Name())
	})
	http.HandleFunc("/add", func(rw http.ResponseWriter, r *http.Request) {
		go func() {
			for i := 1; i < 200; i++ {
				// price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
				price := 1.0
				traceType := bid.Type_Buy
				if i%2 != 0 {
					traceType = bid.Type_Sell
				}
				b.Add(bid.NewUnit(func(u *bid.Unit) {
					u.Type = traceType
					u.Name = "xlj"
					u.Amount = 1
					u.Price = price
					u.UID = rand.Intn(1000) + 1
					u.ID = int(i)
				}))
			}
		}()
		rw.Write([]byte("jayden"))
	})

	http.HandleFunc("/print", func(rw http.ResponseWriter, r *http.Request) {

		fmt.Println(m.State(), m.Name())
		var buf bytes.Buffer
		buf.WriteString("\n\t")

		buf.WriteString(b.Buy().Name())
		buf.WriteString("\n\t")
		for _, n := range b.Buy().NodeList() {
			b, _ := json.Marshal(n.Value)

			buf.Write(b)
			buf.WriteString("\n\t")
		}

		buf.WriteString(b.Sell().Name())
		buf.WriteString("\n\t")
		for _, n := range b.Sell().NodeList() {
			b, _ := json.Marshal(n.Value)

			buf.Write(b)
			buf.WriteString("\n\t")
		}

		rw.Write(buf.Bytes())
	})
	http.ListenAndServe(":8080", nil)
}
