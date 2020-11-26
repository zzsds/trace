package main

import (
	"fmt"
	"log"

	"github.com/zzsds/trace"
	"github.com/zzsds/trade/queue"
)

//go:generate go run main.go
func main() {
	t := trace.NewTrace(func(o *trace.Options) {
		o.Name = "New Product"
	})

	log.Println(t.Name())
	go t.Run()

	q := queue.NewQueue(queue.Name("Buy"))
	log.Println(q.Name(), q.Length())

	select {}
	fmt.Println(123)
}
