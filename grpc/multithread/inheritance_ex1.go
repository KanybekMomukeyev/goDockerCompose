package main

import "fmt"

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human
	school string
}

type Employee struct {
	Human
	company string
}

func (h *Human) sayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

func (e *Employee) sayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone) //Yes you can split into 2 lines here.
}

func (e *Student) sayHi() {
	fmt.Printf("Hi, I am %s, I study at %s. Call me on %s\n", e.name,
		e.school, e.phone) //Yes you can split into 2 lines here.
}

func main() {
	human := Human{"Mark", 25, "222-222-YYYY"}
	sam := Employee{Human{"Sam", 45, "111-888-XXXX"}, "Golang Inc"}
	samSt := Student{Human{"Sam", 45, "111-888-XXXX"}, "BKTL"}

	human.sayHi()
	sam.sayHi()
	samSt.sayHi()
}
