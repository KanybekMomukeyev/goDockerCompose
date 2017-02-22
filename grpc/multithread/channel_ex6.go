package main

import (
	"sync"
	"fmt"
	"time"
)

func main() {
	testQueue()
}

func NewQueue() chan func() {
	//create channel of type function
	var queueChannel = make(chan func())

	//spawn go routine to read and run functions in the channel
	go func() {
		for {
			select {
			case nextFunction := <-queueChannel:
				nextFunction()
			//default:
			//	fmt.Println("default")
			}
		}
	}()

	return queueChannel
}

var channelOfFunc chan func() = NewQueue()

func testQueue() {

	//Create new serial queue
	//serailQueue := New()

	//Number of times to loop
	var loops = 10

	//Queue output will be added here
	var queueOutput string

	//WaitGroup for determining when queue output is finished
	var wg sync.WaitGroup

	//Create function to place in queue
	var printTest = func(index int) {
		fmt.Println("index is =", index)
		queueOutput = fmt.Sprintf("%v%v",queueOutput, index)
		time.Sleep(1 * time.Second)
		wg.Done()
	}

	//Add functions to queue
	var index int;
	for index = 0; index < loops; index ++ {
		wg.Add(1)

		localIndex := index
		funcToQueue := func() {
			printTest(localIndex)
		}

		channelOfFunc <- funcToQueue
	}

	//Wait until all functions in queue are done
	wg.Wait()

	//Generate correct output
	var correctOutput string
	for index = 0; index < loops; index ++ {
		correctOutput = fmt.Sprintf("%v%v", correctOutput, index)
	}

	//Compare queue output with correct output
	if queueOutput != correctOutput {
		fmt.Println("Not equal")
		fmt.Println("Serial Queue produced %v, want %v", queueOutput, correctOutput)
	}

	fmt.Println("Serial Queue produced %v, want %v", queueOutput, correctOutput)
}