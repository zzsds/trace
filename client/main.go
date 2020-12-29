package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zzsds/trade"
	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/match"
	"github.com/zzsds/trade/queue"
)

//go:generate go version
func main() {
	t := trade.Newtrade(func(o *trade.Options) {
		o.Name = "New Product"
	})
	m := match.NewMatch(match.Name("goods")).Register(bid.NewBid(bid.Name("test")))

	t.Add(m)
	go t.Run()
	go func() {
		for {
			select {
			case <-m.Buffer():
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
		for n := m.Bid().Buy().Front(); n != nil; n = n.Next() {
			fmt.Println(n.Data().Content)
		}
		rw.Write([]byte("print request"))
	})
	http.HandleFunc("/add", func(rw http.ResponseWriter, r *http.Request) {
		m, err := t.Load(m.Name())
		if err != nil {
			os.Exit(0)
		}
		b := m.Bid()
		for i := 0; i < 1000; i++ {
			traceType := b.Buy()
			if i%2 != 0 {
				traceType = b.Sell()
			}
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(1000)), 64)
			b.Add(traceType, &bid.Unit{
				Name:   "xlj",
				Amount: i + 1,
				Price:  price,
				UID:    int(i),
				ID:     int(i),
			})
		}
		rw.Write([]byte("add request"))
	})
	// 创建一个监听8000端口的服务器
	http.ListenAndServe(":8000", nil)
}

func queueTest() {
	que := queue.NewQueue(queue.Name("Buy"))
	log.Println(que.Name(), que.Len())

	que.Listen(func(n *queue.Node) error {
		if n.Data().ExpireAt == nil {
			que.WriteBuffer(*n.Data())
			que.Remove(n)
		} else if n.Data().ExpireAt.Before(time.Now()) {
			n.Data().ExpireAt = nil
		}
		return nil
	})

	go func() {
		for {
			select {
			case buff := <-que.Buffer():
				log.Println(buff, "出")
			}
		}
	}()

	rand.Seed(time.Now().Unix())
	data := queue.NewExpireData(&bid.Unit{
		Name:   "qwe",
		Amount: int(rand.Intn(1000)),
		Price:  1.0,
		UID:    0,
		ID:     0,
	}, time.Now().Add(3*time.Second))
	que.Push(data)

	time.Sleep(5 * time.Second)
	data = queue.NewData(&bid.Unit{
		Name:   "xlj",
		Amount: int(rand.Intn(1000)),
		Price:  1.0,
		UID:    0,
		ID:     0,
	})
	que.Unshift(data)

	time.Sleep(2 * time.Second)
	fmt.Println("End")
	for k, v := range que.List() {
		log.Println(k, v.Content, "list")
	}
}
