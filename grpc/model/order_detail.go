package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"fmt"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"errors"
)

var schemaRemoveOrderDetail = `
DROP TABLE IF EXISTS orderdetails;
`

var schemaCreateOrderDetail = `
CREATE TABLE IF NOT EXISTS orderdetails (
    	order_detail_id BIGSERIAL PRIMARY KEY NOT NULL,
    	order_id BIGINT,
    	order_detail_date BIGINT,
    	is_last INTEGER,
    	billing_no varchar (400),
    	product_id BIGINT,
    	price REAL,
	order_quantity REAL,
	discount INTEGER
);
`

var schemaCreateIndexForOrderDetail1 = `CREATE INDEX IF NOT EXISTS order_id_order_details_idx ON orderdetails (order_id)`
var schemaCreateIndexForOrderDetail2 = `CREATE INDEX IF NOT EXISTS product_id_order_details_idx ON orderdetails (product_id)`
var schemaCreateIndexForOrderDetail3 = `CREATE INDEX IF NOT EXISTS order_detail_date_order_details_idx ON orderdetails (order_detail_date)`

type OrderDetail struct {
	orderDetailId uint64 `db:"order_detail_id"`
	orderId uint64 `db:"order_id"`
	orderDetailDate uint64 `db:"order_detail_date"`
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
	db.MustExec(schemaCreateIndexForOrderDetail3)
}

func StoreOrderDetails(db *sqlx.DB, orderDetail *pb.OrderDetailRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO orderdetails " +
		"(order_id, order_detail_date, is_last, billing_no, product_id, price, order_quantity, discount) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8) returning order_detail_id;",
		orderDetail.OrderId,
		orderDetail.OrderDetailDate,
		orderDetail.IsLast,
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

	rows, err := db.Queryx("SELECT order_detail_id, order_id, order_detail_date, is_last, billing_no, product_id, " +
		"price, order_quantity, discount FROM orderdetails ORDER BY order_detail_id DESC")

	if err != nil {
		print("error")
	}

	orderDetails := make([]*pb.OrderDetailRequest, 0)
	for rows.Next() {
		orderDetail := new(pb.OrderDetailRequest)
		err := rows.Scan(&orderDetail.OrderDetailId,
				&orderDetail.OrderId,
				&orderDetail.OrderDetailDate,
				&orderDetail.IsLast,
				&orderDetail.BillingNo,
				&orderDetail.ProductId,
				&orderDetail.Price,
				&orderDetail.OrderQuantity,
				&orderDetail.Discount)

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

func RecentOrderDetailForProduct(db *sqlx.DB, productReq *pb.ProductRequest) (*pb.OrderDetailRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT order_detail_id, order_id, order_detail_date, is_last, billing_no, product_id, " +
		"price, order_quantity, discount FROM orderdetails WHERE product_id=$1 ORDER BY order_detail_date DESC LIMIT $2", productReq.ProductId, 1)

	if err != nil {
		print("error")
	}

	orderDetails := make([]*pb.OrderDetailRequest, 0)
	for rows.Next() {
		orderDetail := new(pb.OrderDetailRequest)
		err := rows.Scan(&orderDetail.OrderDetailId,
			&orderDetail.OrderId,
			&orderDetail.OrderDetailDate,
			&orderDetail.IsLast,
			&orderDetail.BillingNo,
			&orderDetail.ProductId,
			&orderDetail.Price,
			&orderDetail.OrderQuantity,
			&orderDetail.Discount)

		if err != nil {
			return nil, err
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orderDetails) > 0 {
		return orderDetails[0], nil
	}

	return nil, errors.New("Not found")
}
