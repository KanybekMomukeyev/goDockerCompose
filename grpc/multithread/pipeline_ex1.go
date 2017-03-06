package main

import (
	"fmt"
	"time"
)

func main() {
	// Set up the pipeline.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.

	//fmt.Println(<-out) // 4
	//fmt.Println(<-out) // 9

	//for n := range sq(sq(gen(2, 3))) {
	//	fmt.Println(n) // 16 then 81
	//}

	go printChannels(out)


	time.Sleep(2 * time.Second)

	//c := make(k<- chan bool) - you can only write to k (channel of bool)
	//func readFromChannel(input <-chan string) {}

	//r := make(<-chan bool) - can only read from k (channel of bool)
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func printChannels(quit<-chan int) {
	for
	{
		select {
		case msg := <-quit:
			fmt.Println(msg)
			return
		default:
			fmt.Println("default")
		}
	}
}