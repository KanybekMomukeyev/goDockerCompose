package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	model "github.com/KanybekMomukeyev/goDockerCompose/api/models"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"time"
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
	app.Handle("GET", "/", func(ctx context.Context) {
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

	app.Get("/encode", func(ctx context.Context) {
		orderFilter := new(model.OrderFilter)
		orderFilter.OrderDate = uint64(time.Now().UnixNano() / 1000000)
		orderFilter.Limit = 100
		orders, err := model.AllOrdersForFilter(db, orderFilter)
		log.WithFields(log.Fields{"initial len(orders):": len(orders),}).Info("")
		if err != nil {
			log.WithFields(log.Fields{"error":err,}).Warn("ERROR")
		}
		ctx.JSON(orders)
	})
	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"))
}
