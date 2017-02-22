package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func createPizza(number_of_ovens int, out chan<- int) {
	for pizza := 1; pizza <= number_of_ovens; pizza++ {
		time.Sleep(time.Second)
		fmt.Printf("Created Pizza %d \n", pizza)
		out <- pizza
	}
	close(out)
}

func deliverPizza(pizzas int, motoboys int, in <-chan int) {
	for pizzas > 0 {
		for motoboy := 1; motoboy <= motoboys; motoboy++ {
			pizza := <-in
			fmt.Printf("Motoboy %d delivered pizza %d\n", motoboy, pizza)
			pizzas--
			if pizzas == 0 {
				return
			}
		}
		time.Sleep(time.Second)
	}
}

func timeTrack(start time.Time, functionName string) {
	elapsed := time.Since(start)
	fmt.Println(functionName, "took", elapsed)
}

func main() {
	defer timeTrack(time.Now(), "Build Pizzas")
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	number_of_ovens := 5
	motoboys := 3
	pizzaChan := make(chan int)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		createPizza(number_of_ovens, pizzaChan)
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		deliverPizza(number_of_ovens, motoboys, pizzaChan)
	}(&wg)
	wg.Wait()
}