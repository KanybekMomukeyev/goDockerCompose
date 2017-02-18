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