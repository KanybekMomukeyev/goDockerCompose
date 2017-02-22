package main

import (
	"fmt"
	"time"
)

func main() {

	//bridge := make(chan int) //creating the channel, channel reference
	//go my_func(bridge) //creating a new go routine and sharing the reference of the channel
	//
	//for i := 0; i < 10; i++ {
	//	bridge <- i //sending integer through the bridge(channel), this is blocked till its recieved in the goroutine
	//	fmt.Println("in the main function ")
	//
	//}
	//time.Sleep(1 * time.Second)

	/////
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()

	x, y := 1, 1
	for
	{
		select {
		case c <- x:
			x, y = y, x + y
		case <-quit:
			fmt.Println("quit")
			return
		//default:
		//	fmt.Println("default")
		}
	}

	////
	//go forever()
	//select {} // block forever
}

func forever() {
	for {
		fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}

func my_func(referenceToMainBridge chan int) {

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		<-referenceToMainBridge // the recieved value is not used, its just disposed
		//main is blocked till its recieved here
		fmt.Println("In the Go routine ")
	}
}

func fibonacci(c, quit chan int) {
	x, y := 1, 1
	for
	{
		select {
		case c <- x:
			x, y = y, x + y
		case <-quit:
			fmt.Println("quit")
			return
		//default:
		//	fmt.Println("default")
		}
	}
}
