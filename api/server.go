package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	model "github.com/KanybekMomukeyev/goDockerCompose/api/models"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"time"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/view"
)

type User struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Age       int    `json:"age"`
}

var db *sqlx.DB

func main() {
	var databaseError error
	db, databaseError = model.NewDBToConnect("dataSourceName")
	if databaseError != nil {
		log.WithFields(log.Fields{"omg":databaseError,}).Fatal("failed to listen:")
	}

	app := iris.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})

	app.AttachView(view.HTML("./views", ".html"))
	app.Get("/template", func(ctx context.Context) {

		//ctx.ViewData("Name", "Iris") // the .Name inside the ./templates/hi.html
		ctx.Gzip(true)               // enable gzip for big files
		ctx.View("index.html")          // render the template with the file name relative to the './templates'
	})

	app.Get("/txt", func(ctx context.Context) {
		ctx.ServeFile("./views/index.html", true)
	})

	//directory := flag.String("d", ".", "build/index.html")
	//app.StaticWeb("/static", "build/index.html")

	v1 := app.Party("/api/v1")
	v1.Use(crs)
	{
		v1.Get("/home", func(ctx context.Context) {
			ctx.WriteString("Hello from /home")
		})

		v1.Get("/encode", func(ctx context.Context) {

			orderFilter := new(model.OrderFilter)
			orderFilter.OrderDate = uint64(time.Now().UnixNano() / 1000000)
			orderFilter.Limit = 100

			createOrderRequests := make([]*model.CreateOrderRequest, 0)

			orders, err := model.AllOrdersForFilter(db, orderFilter)
			for _, order := range orders {
				createOrderRequest := new(model.CreateOrderRequest)

				customer, error := model.CustomerFor(db, order.CustomerId)
				if error != nil {
				}
				if customer != nil {
					order.BillingNo = customer.FirstName + " " + customer.SecondName + " " + customer.PhoneNumber + " " + customer.Address
				}

				createOrderRequest.Order = order

				payment, error := model.PaymentForOrder(db, order)
				if error != nil {
					break
				}
				createOrderRequest.Payment = payment
				orderDetails, error := model.AllOrderDetailsForOrder(db, order)
				if error != nil {
					break
				}
				createOrderRequest.OrderDetails = orderDetails
				createOrderRequests = append(createOrderRequests, createOrderRequest)
			}


			log.WithFields(log.Fields{"initial len(orders):": len(orders),}).Info("")
			if err != nil {
				log.WithFields(log.Fields{"error":err,}).Warn("ERROR")
			}
			ctx.JSON(createOrderRequests)
		})
		v1.Get("/about", func(ctx context.Context) {
			ctx.WriteString("Hello from /about")
		})
		v1.Post("/send", func(ctx context.Context) {
			ctx.WriteString("sent")
		})
	}

	app.Handle("GET", "/sample", func(ctx context.Context) {
		ctx.HTML("<b> Hello world! </b>")
	})

	app.Get("/user", func(ctx context.Context) {
		doe := User{
			Username:  "Johndoe",
			Firstname: "John",
			Lastname:  "Doe",
			City:      "Neither FBI knows!!!",
			Age:       25,
		}

		ctx.JSON(doe)
	})

	log.Info("I'll be logged with common and other field")
	fmt.Println("started server")

	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"))

}
