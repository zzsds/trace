package main

import (
	"fmt"
	"log"

	"github.com/zzsds/trace"
)

func main() {
	trace := trace.NewTrace(func(o *trace.Options) {
		o.Name = "New Product"
	})

	log.Println(trace.Name())
	go trace.Run()

	select {}
	fmt.Println(123)
}
