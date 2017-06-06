package models

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os"
)

type OrderFilter struct {
	OrderKeyword string
	OrderDate    uint64
	Limit        uint32
}

type OrderRequest struct {
	OrderId           uint64
	OrderDocument     uint32
	MoneyMovementType uint32
	BillingNo         string
	StaffId           uint64
	CustomerId        uint64
	SupplierId        uint64
	OrderDate         uint64
	PaymentId         uint64
	ErrorMsg          string
	Comment           string
	IsDeleted         uint32
	IsPaid            uint32
	IsEdited          uint32
	OrderUpdatedAt    uint64
}

func NewDBToConnect(dataSourceName string) (*sqlx.DB, error) {

	connInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_ENV_POSTGRES_USER"),
		os.Getenv("DB_ENV_POSTGRES_DATABASENAME"),
		os.Getenv("DB_ENV_POSTGRES_PASSWORD"),
		os.Getenv("GODOCKERCOMPOSE_POSTGRES_1_PORT_5432_TCP_ADDR"),
		os.Getenv("GODOCKERCOMPOSE_POSTGRES_1_PORT_5432_TCP_PORT"),
	)

	fmt.Println(connInfo)
	db, err := sqlx.Connect("postgres", connInfo) // for compose

	//db, err := sqlx.Connect("postgres", "user=kanybek dbname=databasename password=nazgulum host=172.17.0.4 port=5432 sslmode=disable") // for single docker app
	//db, err := sqlx.Connect("postgres", "user=kanybek dbname=databasename password=nazgulum host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.WithFields(log.Fields{
			"connection info": connInfo,
			"error": err,
		}).Fatal("Can not connected to database")
		return nil, err
	}

	pingError := db.Ping()

	if pingError != nil {
		log.WithFields(log.Fields{
			"info": connInfo,
			"ping": pingError,
		}).Fatal("Can not connected Ping to database")
		return nil, pingError
	}

	return db, nil
}

func AllOrdersForFilter(db *sqlx.DB, orderFilter *OrderFilter) ([]*OrderRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT order_id, order_document, money_movement, billing_no, " +
		"staff_id, customer_id, supplier_id, order_date, " +
		"payment_id, error_msg, comment, is_deleted, is_paid," +
		" is_editted FROM orders WHERE order_date<=$1 ORDER BY order_date DESC LIMIT $2", orderFilter.OrderDate, orderFilter.Limit)

	if err != nil {
		log.WithFields(log.Fields{"error":err,}).Warn("ERROR")
	}

	orders, err := scanOrderRowsWWW(rows)

	if err = rows.Err(); err != nil {
		log.WithFields(log.Fields{"error":err,}).Warn("ERROR")
		return nil, err
	}

	return orders, nil
}

func scanOrderRowsWWW(rows *sqlx.Rows) ([]*OrderRequest, error) {
	orders := make([]*OrderRequest, 0)
	for rows.Next() {
		order := new(OrderRequest)
		err := rows.Scan(&order.OrderId,
			&order.OrderDocument,
			&order.MoneyMovementType,
			&order.BillingNo,
			&order.StaffId,
			&order.CustomerId,
			&order.SupplierId,
			&order.OrderDate,
			&order.PaymentId,
			&order.ErrorMsg,
			&order.Comment,
			&order.IsDeleted,
			&order.IsPaid,
			&order.IsEdited)
		if err != nil {
			log.WithFields(log.Fields{"scanOrderRows":err,}).Warn("ERROR")
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}