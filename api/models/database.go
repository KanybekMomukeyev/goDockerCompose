package models

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os"
	"errors"
)

type OrderFilter struct {
	OrderKeyword string
	OrderDate    uint64
	Limit        uint32
}

type OrderRequest struct {
	OrderId           uint64
	OrderDocument     uint32 `json:"-"`
	MoneyMovementType uint32 `json:"-"`
	BillingNo         string `json:"Customer"`
	StaffId           uint64 `json:"-"`
	CustomerId        uint64 `json:"-"`
	SupplierId        uint64 `json:"-"`
	OrderDate         uint64
	PaymentId         uint64 `json:"-"`
	ErrorMsg          string `json:"-"`
	Comment           string `json:"-"`
	IsDeleted         uint32 `json:"-"`
	IsPaid            uint32 `json:"-"`
	IsEdited          uint32 `json:"-"`
	OrderUpdatedAt    uint64 `json:"-"`
}

type PaymentRequest struct {
	PaymentId              uint64 `json:"-"`
	TotalOrderPrice        float64
	Discount               float64
	TotalPriceWithDiscount float64
}

type OrderDetailRequest struct {
	OrderDetailId   uint64 `json:"-"`
	OrderId         uint64 `json:"-"`
	OrderDetailDate uint64 `json:"-"`
	IsLast          uint32 `json:"-"`
	ProductId       uint64 `json:"-"`
	BillingNo       string `json:"Product"`
	Price           float64
	OrderQuantity   float64 `json:"Quantity"`
	Discount        int32 `json:"-"`
	QuantityPerUnit string `json:"PerUnit"`
}

type CreateOrderRequest struct {
	Order        *OrderRequest
	Payment      *PaymentRequest
	OrderDetails []*OrderDetailRequest
}

type CustomerRequest struct {
	CustomerId        uint64
	CustomerImagePath string
	FirstName         string
	SecondName        string
	PhoneNumber       string
	Address           string
	StaffId           uint64
	CustomerUpdatedAt uint64
}

type ProductRequest struct {
	ProductId        uint64
	ProductImagePath string
	ProductName      string
	SupplierId       uint64
	CategoryId       uint64
	Barcode          string
	QuantityPerUnit  string
	SaleUnitPrice    float64
	IncomeUnitPrice  float64
	UnitsInStock     float64
	ProductUpdatedAt uint64
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
		" is_editted FROM orders WHERE order_date<=$1 AND order_document=1000 ORDER BY order_date DESC LIMIT $2", orderFilter.OrderDate, orderFilter.Limit)

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

func AllOrderDetailsForOrder(db *sqlx.DB, order *OrderRequest) ([]*OrderDetailRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT order_detail_id, order_id, order_detail_date, is_last, billing_no, product_id, " +
		"price, order_quantity, discount FROM orderdetails WHERE order_id=$1" +
		" ORDER BY order_detail_date DESC", order.OrderId)

	if err != nil {
		print(" AllOrderDetailsForOrder_error ")
	}

	orderDetails := make([]*OrderDetailRequest, 0)
	for rows.Next() {
		orderDetail := new(OrderDetailRequest)
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

	for _, orderDetail := range orderDetails {
		product, err := ProductFor(db, orderDetail.ProductId)
		if err != nil {
			print(" ProductFor_error_not_found ")
		}
		if product != nil {
			orderDetail.BillingNo = product.ProductName
			orderDetail.QuantityPerUnit = product.QuantityPerUnit
		}
	}

	return orderDetails, nil
}

func PaymentForOrder(db *sqlx.DB, order *OrderRequest) (*PaymentRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT payment_id, total_order_price, discount, total_price_with_discount " +
		"FROM payments WHERE payment_id=$1 LIMIT $2", order.PaymentId, 1)

	if err != nil {
		print(" PaymentForOrder_error ")
	}

	payments := make([]*PaymentRequest, 0)
	for rows.Next() {
		payment := new(PaymentRequest)
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

	log.WithFields(log.Fields{"order.OrderId": order.OrderId}).Warn("")
	return nil, errors.New("Not found PaymentForOrder")
}

func ProductFor(db *sqlx.DB, productId uint64) (*ProductRequest, error) {

	rows, err := db.Queryx("SELECT product_id, product_image_path, product_name, supplier_id, " +
		"category_id, barcode, quantity_per_unit, sale_unit_price, " +
		"income_unit_price, units_in_stock FROM products WHERE product_id=$1", productId)

	if err != nil {
		print(" ProductFor_error_Queryx ")
	}

	products := make([]*ProductRequest, 0)
	for rows.Next() {
		product := new(ProductRequest)
		err := rows.Scan(&product.ProductId,
				&product.ProductImagePath,
				&product.ProductName,
				&product.SupplierId,
				&product.CategoryId,
				&product.Barcode,
				&product.QuantityPerUnit,
				&product.SaleUnitPrice,
				&product.IncomeUnitPrice,
				&product.UnitsInStock)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(products) > 0 {
		return products[0], nil
	}

	return nil, errors.New("Not found Product")
}

func CustomerFor(db *sqlx.DB, customerId uint64) (*CustomerRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT customer_id, customer_image_path, first_name, second_name, phone_number, address, staff_id, updated_at FROM customers WHERE customer_id=$1", customerId)
	if err != nil {
		log.WithFields(log.Fields{"error": err, }).Warn("")
		return nil, err
	}

	realCustomers := make([]*CustomerRequest, 0)
	for rows.Next() {
		customer := new(CustomerRequest)
		err := rows.Scan(&customer.CustomerId,
				&customer.CustomerImagePath,
				&customer.FirstName,
				&customer.SecondName,
				&customer.PhoneNumber,
				&customer.Address,
				&customer.StaffId,
				&customer.CustomerUpdatedAt)

		if err != nil {
			log.WithFields(log.Fields{"error": err, }).Warn("")
			return nil, err
		}
		realCustomers = append(realCustomers, customer)
	}

	if err = rows.Err(); err != nil {
		log.WithFields(log.Fields{"error": err, }).Warn("")
		return nil, err
	}

	if len(realCustomers) > 0 {
		return realCustomers[0], nil
	}

	return nil, errors.New("Not found Customer")
}