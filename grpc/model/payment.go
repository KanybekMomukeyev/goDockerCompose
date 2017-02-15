package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"errors"
)

var schemaRemovePayment = `
DROP TABLE IF EXISTS payments;
`

var schemaCreatePayment = `
CREATE TABLE IF NOT EXISTS payments (
    payment_id BIGSERIAL PRIMARY KEY NOT NULL,
    total_order_price REAL,
    discount REAL,
    total_price_with_discount REAL
);
`

type Payment struct {
	paymentId uint64 `db:"payment_id"`
	totalOrderPrice float32 `db:"total_order_price"`
	discount float32 `db:"discount"`
	totalPriceWithDiscount float32 `db:"total_price_with_discount"`
}

func CreatePaymentIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemovePayment)
	db.MustExec(schemaCreatePayment)
}

func StorePayment(db *sqlx.DB, payment *pb.PaymentRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO payments " +
		"(total_order_price, discount, total_price_with_discount) " +
		"VALUES($1, $2, $3) returning payment_id;",
		payment.TotalOrderPrice,
		payment.Discount,
		payment.TotalPriceWithDiscount).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	log.WithFields(log.Fields{
		"count payment_id": lastInsertId,
	}).Info("Payment success saved")

	return lastInsertId, nil
}

func AllPayments(db *sqlx.DB) ([]*pb.PaymentRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT payment_id, total_order_price, discount, total_price_with_discount " +
		"FROM payments ORDER BY payment_id DESC")

	if err != nil {
		print("error")
	}

	payments := make([]*pb.PaymentRequest, 0)
	for rows.Next() {
		payment := new(pb.PaymentRequest)
		err := rows.Scan(&payment.PaymentId, &payment.TotalOrderPrice,
			&payment.Discount, &payment.TotalPriceWithDiscount)

		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func PaymentForOrder(db *sqlx.DB, order *pb.OrderRequest) (*pb.PaymentRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT payment_id, total_order_price, discount, total_price_with_discount " +
		"FROM payments WHERE payment_id=$1 LIMIT $2", order.PaymentId, 1)

	if err != nil {
		print("error")
	}

	payments := make([]*pb.PaymentRequest, 0)
	for rows.Next() {
		payment := new(pb.PaymentRequest)
		err := rows.Scan(&payment.PaymentId, &payment.TotalOrderPrice,
			&payment.Discount, &payment.TotalPriceWithDiscount)

		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(payments) > 0 {
		return payments[0], nil
	}

	return nil, errors.New("Not found")
}
