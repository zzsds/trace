package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zzsds/trade"
	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/match"
)

var t trade.Server

//go:generate go version
func main() {
	t = trade.Newtrade(func(o *trade.Options) {
		o.Name = "New Product"
	})
	m := match.NewMatch(match.Name("goods")).Register(bid.NewBid(bid.Name("USDT-BTC")))
	t.Register(m)
	go t.Run()

	go func() {
		m, err := t.Load(m.Name())
		if err != nil {
			os.Exit(0)
		}
		for {
			select {
			case msg := <-m.Buffer():
				_ = msg.Amount
				// log.Println(msg)
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// 输出到STDOUT展示处理已经开始
		fmt.Fprint(os.Stdout, "processing request\n")
		// 通过select监听多个channel
		select {
		case <-time.After(2 * time.Second):
			// 如果两秒后接受到了一个消息后，意味请求已经处理完成
			// 我们写入"request processed"作为响应
			w.Write([]byte("request processed"))
		case <-ctx.Done():

			// 如果处理完成前取消了，在STDERR中记录请求被取消的消息
			fmt.Fprint(os.Stderr, "request cancelled\n")
		}
	})
	http.HandleFunc("/start", func(rw http.ResponseWriter, r *http.Request) {
		m, err := t.Load(m.Name())
		if err != nil {
			os.Exit(0)
		}

		fmt.Println(m.Start(), m.State(), m.Name())
		rw.Write([]byte("start request"))
	})
	http.HandleFunc("/stop", func(rw http.ResponseWriter, r *http.Request) {
		m, err := t.Load(m.Name())
		if err != nil {
			os.Exit(0)
		}

		fmt.Println(m.Stop(), m.State(), m.Name())
		rw.Write([]byte("stop request"))
	})
	http.HandleFunc("/print", func(rw http.ResponseWriter, r *http.Request) {
		m, err := t.Load(m.Name())
		if err != nil {
			os.Exit(0)
		}

		fmt.Println(m.State(), m.Name())
		var buf bytes.Buffer
		buf.WriteString("\n\t")
		for n := m.Bid().Buy().Front(); n != nil; n = n.Next() {
			b, _ := json.Marshal(n.Value)

			buf.Write(b)
			buf.WriteString("\n\t")
		}
		for n := m.Bid().Sell().Front(); n != nil; n = n.Next() {
			b, _ := json.Marshal(n.Value)

			buf.Write(b)
			buf.WriteString("\n\t")
		}
		rw.Write(buf.Bytes())
	})
	http.HandleFunc("/add", func(rw http.ResponseWriter, r *http.Request) {
		m, err := t.Load(m.Name())
		if err != nil {
			os.Exit(0)
		}
		go func() {
			for i := 1; i < 10000000; i++ {
				price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
				traceType := bid.Type_Buy
				if i%2 != 0 {
					traceType = bid.Type_Sell
				}
				m.Bid().Add(bid.NewUnit(func(u *bid.Unit) {
					u.Type = traceType
					u.Name = "xlj"
					u.Amount = rand.Intn(1000)
					u.Price = price
					u.UID = rand.Intn(1000)
					u.ID = int(i)
				}))
			}
		}()
		rw.Write([]byte("add request"))
	})
	// 创建一个监听8000端口的服务器
	http.ListenAndServe(":8000", nil)
}
