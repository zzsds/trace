package main

import (
	"fmt"
	"log"

	"github.com/zzsds/trade"
	"github.com/zzsds/trade/queue"
)

//go:generate go run main.go
func main() {
	t := trade.Newtrade(func(o *trade.Options) {
		o.Name = "New Product"
	})

	log.Println(t.Name())
	go t.Run()

	q := queue.NewQueue(queue.Name("Buy"))
	log.Println(q.Name(), q.Length())

	select {}
	fmt.Println(123)
}
