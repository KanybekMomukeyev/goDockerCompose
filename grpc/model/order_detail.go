package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"fmt"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveOrderDetail = `
DROP TABLE IF EXISTS orderdetails;
`

var schemaCreateOrderDetail = `
CREATE TABLE IF NOT EXISTS orderdetails (
    	order_detail_id BIGSERIAL PRIMARY KEY NOT NULL,
    	order_id BIGINT,
    	billing_no varchar (400),
    	product_id BIGINT,
    	price REAL,
	order_quantity REAL,
	discount INTEGER
);
`

var schemaCreateIndexForOrderDetail1 = `CREATE INDEX IF NOT EXISTS order_id_order_details_idx ON orderdetails (order_id)`
var schemaCreateIndexForOrderDetail2 = `CREATE INDEX IF NOT EXISTS product_id_order_details_idx ON orderdetails (product_id)`


type OrderDetail struct {
	orderDetailId uint64 `db:"order_detail_id"`
	orderId uint64 `db:"order_id"`
	billingNo string `db:"billing_no"`
	productId uint64 `db:"product_id"`
	price float32 `db:"price"`
	orderQuantity float32 `db:"order_quantity"`
	discount int32 `db:"discount"`
}

func CreateOrderDetailsIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveOrderDetail)
	db.MustExec(schemaCreateOrderDetail)
	db.MustExec(schemaCreateIndexForOrderDetail1)
	db.MustExec(schemaCreateIndexForOrderDetail2)
}

func StoreOrderDetails(db *sqlx.DB, orderDetail *pb.OrderDetailRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO orderdetails " +
		"(order_id, billing_no, product_id, price, order_quantity, discount) " +
		"VALUES($1, $2, $3, $4, $5, $6) returning order_detail_id;",
		orderDetail.OrderId,
		orderDetail.BillingNo,
		orderDetail.ProductId,
		orderDetail.Price,
		orderDetail.OrderQuantity,
		orderDetail.Discount).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted order_detail_id =", lastInsertId)

	return lastInsertId, nil
}

func AllOrderDetails(db *sqlx.DB) ([]*pb.OrderDetailRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT order_detail_id, order_id, billing_no, product_id, " +
		"price, order_quantity, discount FROM orderdetails ORDER BY order_detail_id DESC")

	if err != nil {
		print("error")
	}

	orderDetails := make([]*pb.OrderDetailRequest, 0)
	for rows.Next() {
		orderDetail := new(pb.OrderDetailRequest)
		err := rows.Scan(&orderDetail.OrderDetailId, &orderDetail.OrderId, &orderDetail.BillingNo,
			&orderDetail.ProductId, &orderDetail.Price,
			&orderDetail.OrderQuantity, &orderDetail.Discount)

		if err != nil {
			return nil, err
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orderDetails, nil
}

