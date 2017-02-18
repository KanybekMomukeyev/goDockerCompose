package main

import (
	"fmt"
	"time"
)

func main() {

	bridge := make(chan int) //creating the channel
	go my_func(bridge) //creating a new go routine and sharing the reference of the channel
	for i := 0; i < 10; i++ {
		bridge <- i //sending integer through the bridge(channel), this is blocked till its recieved in the goroutine
		fmt.Println("in the main function ")

	}
	//time.Sleep(1 * time.Second)

}

func my_func(referenceToMainBridge chan int) {

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		<-referenceToMainBridge // the recieved value is not used, its just disposed
		//main is blocked till its recieved here
		fmt.Println("In the Go routine ")
	}
}