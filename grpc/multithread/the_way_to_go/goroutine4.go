// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main
package main

import (
	"fmt"
)

func generate(ch chan int, done chan bool) {

	for j := 2; j <= 100; j++ {
		ch <- j // Send 'i' to channel 'ch'.
	}
	done <- true
}

func filter(in, out chan int) {

	for j := range in {
		if j%2 == 0 {
			out <- j // Send 'i' to channel 'out'.
		}
	}
}

func main() {
	done := make(chan bool, 1)
	ch := make(chan int)

	go generate(ch, done)
	go func() {
		for prime := range ch {

			fmt.Print(prime, " ")
			ch1 := make(chan int)
			go filter(ch, ch1)

			ch = ch1
		}
	}()
	<-done
}
