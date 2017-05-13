package testpackage

import (
	"testing"
	"fmt"
	"sync"
)

func New() chan func() {
	//create channel of type function
	var queue = make(chan func())

	//spawn go routine to read and run functions in the channel
	go func() {
		for nextFunction := range queue {
			nextFunction()
		}
	}()

	return queue
}

func TestQueue(t *testing.T) {
	//Create new serial queue
	queue := New()

	//Number of times to loop
	var loops = 100

	//Queue output will be added here
	var queueOutput string

	//WaitGroup for determining when queue output is finished
	var waitGroup sync.WaitGroup

	//Create function to place in queue
	var printTest = func(i int) {
		queueOutput = fmt.Sprintf("%v%v",queueOutput, i)
		waitGroup.Done()
	}

	//Add functions to queue
	var i int;
	for i=0; i<loops; i++ {
		waitGroup.Add(1)

		t:=i
		f := func() {
			printTest(t)
		}

		queue <- f
	}

	//Generate correct output
	var correctOutput string
	for i=0; i<loops; i++ {
		correctOutput = fmt.Sprintf("%v%v", correctOutput, i)
	}

	//Wait until all functions in queue are done
	waitGroup.Wait()

	//Compare queue output with correct output
	if queueOutput != correctOutput {
		t.Errorf("Serial Queue produced %v, want %v", queueOutput, correctOutput);
	}
}