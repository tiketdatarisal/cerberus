package main

import (
	"fmt"
	"github.com/tiketdatarisal/cerberus"
	"time"
)

func main() {
	tq := cerberus.NewMemQueue()
	go func() { _ = tq.RunEvery(10 * time.Millisecond) }()

	go func() {
		i := 0
		for {
			go tq.Enqueue(func() {
				fmt.Println("Task", i)
				i++
			})
		}
	}()

	fmt.Println("Wait for 2 second")
	time.Sleep(2 * time.Second)
	_ = tq.Close()
}
