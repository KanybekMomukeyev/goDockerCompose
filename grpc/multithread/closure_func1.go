package main
import "fmt"

func counter(start int) (func() int, func()) { // return int, and void

	// if the value gets mutated, the same is reflected in closure
	ctr := func() int {
		return start
	}

	incr := func() {
		start++
	}

	// both ctr and incr have same reference to start
	// closures are created, but are not called
	return ctr, incr
}

func main() {
	// ctr, incr and ctr1, incr1 are different
	ctr1, incr1 := counter(100)
	ctr2, incr2 := counter(100)

	fmt.Println("counter1 - ", ctr1())
	fmt.Println("counter2 - ", ctr2())

	// incr by 1
	incr1()

	fmt.Println("counter1 - ", ctr1())
	fmt.Println("counter2- ", ctr2())

	// incr1 by 3
	incr2()
	incr2()
	incr2()

	fmt.Println("counter1 - ", ctr1())
	fmt.Println("counter2- ", ctr2())

	testAnotherClosure()
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		fmt.Println("sum=", sum)
		fmt.Println("x=", x)
		sum += x
		fmt.Println("infunction sum=", sum, "x=", x)
		return sum
	}
}

func testAnotherClosure() {

	//neg := adder()
	pos := adder()

	for i := 0; i < 10; i++ {
		fmt.Println("----------------------------------")
		fmt.Println("index=", i, pos(i))
		//fmt.Println("index=", i, pos(i), "||", neg(-2*i))
		fmt.Println("==================================")
	}
}