package main

import (
	"fmt"
	"log"

	"github.com/zzsds/trace"
	"github.com/zzsds/trace/queue"
)

func main() {
	t := trace.NewTrace(func(o *trace.Options) {
		o.Name = "New Product"
	})

	log.Println(t.Name())
	go t.Run()

	q := queue.NewQueue(func(o *queue.Options) {
		o.Name = "Test"
	})
	log.Println(q.Name(), q.Length())

	select {}
	fmt.Println(123)
}
