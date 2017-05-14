package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"

)

var schemaRemoveCustomer = `
DROP TABLE IF EXISTS customers;
`

var schemaCreateCustomer = `
CREATE TABLE IF NOT EXISTS customers (
    customer_id BIGSERIAL PRIMARY KEY NOT NULL,
    customer_image_path varchar (400),
    first_name varchar (400),
    second_name varchar (400),
    phone_number varchar (400),
    address varchar (400),
    staff_id BIGINT,
    updated_at BIGINT
);
`

type Customer struct {
	customerId uint64 `db:"customer_id"`
	customerImagePath string `db:"customer_image_path"`
	firstName string `db:"first_name"`
	secondName string `db:"second_name"`
	phoneNumber string `db:"phone_number"`
	address float32 `db:"address"`
}

func DeleteCustomerIfNotExsists(db *sqlx.DB)  {
	db.MustExec(schemaRemoveCustomer)
}

func CreateCustomerIfNotExsists(db *sqlx.DB) {
	db.MustExec(schemaCreateCustomer)
	//db.MustExec("ALTER TABLE customers DROP COLUMN IF EXISTS staff_id")
	db.MustExec("ALTER TABLE customers ADD COLUMN IF NOT EXISTS staff_id BIGINT DEFAULT 0")
	db.MustExec("ALTER TABLE customers ADD COLUMN IF NOT EXISTS updated_at BIGINT DEFAULT 0")
}

func StoreCustomer(tx *sqlx.Tx, customer *pb.ExampleRequest) (uint64, error)  {

	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO customers (first_name, second_name, phone_number) VALUES($1, $2, $3) returning customer_id;", customer.Name, customer.Phone, customer.Email).Scan(&lastInsertId)
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"last inserted customer_id":  lastInsertId,
	}).Info("")
	return lastInsertId, nil
}

func StoreRealCustomer(tx *sqlx.Tx, customerRequest *pb.CustomerRequest) (uint64, error)  {

	var lastInsertId uint64
	err := tx.QueryRow("INSERT INTO customers(customer_image_path, first_name, second_name, phone_number, address, staff_id, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) returning customer_id;",
		customerRequest.CustomerImagePath,
		customerRequest.FirstName,
		customerRequest.SecondName,
		customerRequest.PhoneNumber,
		customerRequest.Address,
		customerRequest.StaffId,
		customerRequest.CustomerUpdatedAt).Scan(&lastInsertId)

	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"last inserted customer_id":  lastInsertId,
	}).Info("")
	return lastInsertId, nil
}

func UpdateRealCustomer(tx *sqlx.Tx, customerReq *pb.CustomerRequest) (uint64, error)  {

	stmt, err :=tx.Prepare("UPDATE customers SET customer_image_path=$1, first_name=$2, second_name=$3, " +
		"phone_number=$4, address=$5, staff_id=$6, updated_at=$7 WHERE customer_id=$8")
	if err != nil {
		return ErrorFunc(err)
	}

	res, err := stmt.Exec(customerReq.CustomerImagePath,
		customerReq.FirstName,
		customerReq.SecondName,
		customerReq.PhoneNumber,
		customerReq.Address,
		customerReq.StaffId,
		customerReq.CustomerUpdatedAt,
		customerReq.CustomerId)
	if err != nil {
		return ErrorFunc(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"update customer rows changed":  affect,
	}).Info("")
	return uint64(affect), nil
}

func AllRealCustomers(db *sqlx.DB) ([]*pb.CustomerRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT customer_id, customer_image_path, first_name, second_name, phone_number, address, staff_id, updated_at FROM customers ORDER BY customer_id ASC")
	if err != nil {
		log.WithFields(log.Fields{"error": err, }).Warn("")
		return nil, err
	}

	realCustomers := make([]*pb.CustomerRequest, 0)
	for rows.Next() {
		customer := new(pb.CustomerRequest)
		err := rows.Scan(&customer.CustomerId, &customer.CustomerImagePath, &customer.FirstName, &customer.SecondName, &customer.PhoneNumber, &customer.Address, &customer.StaffId, &customer.CustomerUpdatedAt)
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

	return realCustomers, nil
}

func AllUpdatedCustomers(db *sqlx.DB, custFilter *pb.CustomerFilter) ([]*pb.CustomerRequest, error)  {

	rows, err := db.Queryx("SELECT customer_id, customer_image_path, first_name, second_name, phone_number, address, staff_id, updated_at FROM customers WHERE updated_at >= $1 LIMIT $2", custFilter.CustomerUpdatedAt, 1000)
	if err != nil {
		log.WithFields(log.Fields{"error": err, }).Warn("")
		return nil, err
	}

	realCustomers := make([]*pb.CustomerRequest, 0)
	for rows.Next() {
		customer := new(pb.CustomerRequest)
		err := rows.Scan(&customer.CustomerId, &customer.CustomerImagePath, &customer.FirstName, &customer.SecondName, &customer.PhoneNumber, &customer.Address, &customer.StaffId, &customer.CustomerUpdatedAt)
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

	return realCustomers, nil
}
