package model

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	"log"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveOrder = `
DROP TABLE IF EXISTS orders;
`

var schemaCreateOrder = `
CREATE TABLE IF NOT EXISTS orders (
    order_id BIGSERIAL PRIMARY KEY NOT NULL,
    order_document INTEGER,
    money_movement INTEGER,
    billing_no varchar (400),

    staff_id BIGINT,
    customer_id BIGINT,
    supplier_id BIGINT,
    order_date BIGINT,
    payment_id BIGINT,

    error_msg varchar (400),
    comment varchar (400),
    is_deleted INTEGER,
    is_paid INTEGER,
    is_editted INTEGER
);
`

var schemaCreateIndexForOrder1 = `CREATE INDEX IF NOT EXISTS customer_id_orders_idx ON orders (customer_id)`
var schemaCreateIndexForOrder2 = `CREATE INDEX IF NOT EXISTS staff_id_orders_idx ON orders (staff_id)`
var schemaCreateIndexForOrder3 = `CREATE INDEX IF NOT EXISTS supplier_id_orders_idx ON orders (supplier_id)`
var schemaCreateIndexForOrder4 = `CREATE INDEX IF NOT EXISTS payment_id_orders_idx ON orders (payment_id)`

type Order struct {
	orderId uint64 `db:"order_id"`
	orderDocument uint32 `db:"order_document"`
	moneyMovementType uint32 `db:"money_movement"`
	billingNo string `db:"billing_no"`

	staffId uint64 `db:"staff_id"`
	customerId uint64 `db:"customer_id"`
	supplierId uint64 `db:"supplier_id"`
	paymentId uint64 `db:"payment_id"`

	orderDate uint64 `db:"order_date"`
	errorMsg string `db:"error_msg"`
	comment string `db:"comment"`

	isDeleted uint32 `db:"is_deleted"`
	isPaid uint32 `db:"is_paid"`
	isEdited uint32 `db:"is_editted"`
}

func CreateOrderIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveOrder)
	db.MustExec(schemaCreateOrder)
	db.MustExec(schemaCreateIndexForOrder1)
	db.MustExec(schemaCreateIndexForOrder2)
	db.MustExec(schemaCreateIndexForOrder3)
	db.MustExec(schemaCreateIndexForOrder4)
}

func StoreOrder(db *sqlx.DB, order *pb.OrderRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO orders " +
		"(order_document, money_movement, billing_no, staff_id, customer_id," +
		" supplier_id, order_date, payment_id, error_msg, comment, is_deleted, is_paid, is_editted) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning order_id;",
		order.OrderDocument,
		order.MoneyMovementType,
		order.BillingNo,
		order.StaffId,
		order.CustomerId,
		order.SupplierId,
		order.OrderDate,
		order.PaymentId,
		order.ErrorMsg,
		order.Comment,
		order.IsDeleted,
		order.IsPaid,
		order.IsEdited).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
}

func AllOrders(db *sqlx.DB) ([]*pb.OrderRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT order_id, order_document, money_movement, billing_no, " +
		"staff_id, customer_id, supplier_id, order_date, " +
		"payment_id, error_msg, comment, is_deleted, is_paid," +
		" is_editted  FROM orders ORDER BY order_date DESC")

	if err != nil {
		print("error")
	}

	orders := make([]*pb.OrderRequest, 0)
	for rows.Next() {
		order := new(pb.OrderRequest)
		err := rows.Scan(&order.OrderId, &order.OrderDocument, &order.MoneyMovementType,
			&order.BillingNo, &order.StaffId, &order.CustomerId,
			&order.SupplierId, &order.OrderDate, &order.PaymentId,
			&order.ErrorMsg, &order.Comment, &order.IsDeleted,
			&order.IsPaid, &order.IsEdited)

		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}