package main

import (
	"fmt"
)

func main() {
	//step 1: create the bridge
	bridgeForOne := make(chan int)
	//call the function in a new goroutine
	//and send the reference of channel to the new goroutine
	//using this reference the new go routine will
	//communicate and synchronize with main
	go testRunConcurrent(bridgeForOne)
	//lets onboard a integer onto the bridge
	bridgeForOne <- 1122

	fmt.Println("Integer recieved at the other end of the bridge in the new go routine and hence the main unblocks and control comes here")

	i := 1000
	p := &i // get reference to i
	fmt.Println(*p) // get value from reference
	fmt.Println(p)

	*p = *p / 10
	fmt.Println(i)
	fmt.Println(*p)

	//ic_send_only := make (<-chan int) //a channel that can only send data - arrow going out is sending
	//ic_recv_only := make (chan<- int) //a channel that can only receive a data - arrow going in is receiving

	//pongs chan<- string   ==> pongs <- "hello world" //channel reads from
	//pings <-chan string   ==> msg := <-pings         //channel gives to

	//func ping(pings chan<- string, msg string) { only accepts a channel for sending values
	//pings <- msg
	//}

	//func pong(pings <-chan string, pongs chan<- string) { pong function accepts one channel for receives (pings) and a second for sends (pongs).
	//msg := <-pings
	//pongs <- msg
	//}
	pings := make(chan string)
	pongs := make(chan string)

	go ping(pings, "passed message")
	go pong(pings, pongs)

	fmt.Println(<-pongs)
}

func testRunConcurrent(bridgeReferenceFromMain chan int) {
	//Inside the new Go routine
	fmt.Println("Inside the new goroutine")

	//recieving the integer of the channel bridge

	takingTheIntegerOffTheBridge := <-bridgeReferenceFromMain

	//only after this recieve from the other end of the shared channel,
	//the send through the channel on main function unblocks and continues execution

	fmt.Printf("\nHere is the number sent across the bridge from main: %d\n", takingTheIntegerOffTheBridge)
	//since the function called from your new goroutine ends execution here,
	//the new goroutines also gracefully ends.
}

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}