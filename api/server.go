package main

import (
	"fmt"
	"github.com/astaxie/beego"
	//"github.com/KanybekMomukeyev/goDockerCompose/api/models"
	"strconv"
)

func main() {
	fmt.Print("Hello world!\n")
	beego.Router("/:operation/:num1:int/:num2:int", &mainController{})
	beego.Run()
}

// This is the controller that this application uses
type mainController struct {
	beego.Controller
}

// Get() handles all requests to the route defined above
func (c *mainController) Get() {
	//Obtain the values of the route parameters defined in the route above
	operation := c.Ctx.Input.Param(":operation")
	num1, _ := strconv.Atoi(c.Ctx.Input.Param(":num1"))
	num2, _ := strconv.Atoi(c.Ctx.Input.Param(":num2"))

	//Set the values for use in the template
	c.Data["operation"] = operation
	c.Data["num1"] = num1
	c.Data["num2"] = num2
	c.TplName = "result.html"

	// Perform the calculation depending on the 'operation' route parameter
	switch operation {
	case "sum":
		c.Data["result"] = add(num1, num2)
	case "product":
		c.Data["result"] = multiply(num1, num2)
	default:
		c.TplName = "invalid-route.html"
	}
}

func add(n1, n2 int) int {
	//fmt.Print(SomeMethod1())
	fmt.Print("add function called")
	//models.SomeDatabaseFunction()
	return n1 + n2
}

func multiply(n1, n2 int) int {
	//fmt.Print(SomeMethod2())
	fmt.Print("multiply function called")
	//models.SomeDatabaseFunction()
	return n1 * n2
}

//http://localhost:8080/sum/1/4
//http://192.241.159.66:8080/sum/1/4
//http://138.68.84.55:3000/sum/1/4
//http://138.68.84.55:8080/sum/1/4
