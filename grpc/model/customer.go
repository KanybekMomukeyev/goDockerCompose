package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"fmt"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveCustomer = `
DROP TABLE IF EXISTS customers;
`

var schemaCreateCustomer = `
CREATE TABLE IF NOT EXISTS customers (
    customer_id SERIAL PRIMARY KEY NOT NULL,
    customer_image_path varchar (400),
    first_name varchar (400),
    second_name varchar (400),
    phone_number varchar (400),
    address varchar (400)
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

func CreateCustomerIfNotExsists(db *sqlx.DB) {
	// for some migrations
	//db.MustExec(schemaCustomerDelete)
	db.MustExec(schemaCreateCustomer)
}

func StoreCustomer(db *sqlx.DB, customer *pb.ExampleRequest) (uint64, error)  {

	tx := db.MustBegin()

	var lastInsertId uint64
	err := tx.QueryRow("INSERT INTO customers (first_name, second_name, phone_number) VALUES($1, $2, $3) returning customer_id;", customer.Name, customer.Phone, customer.Email).Scan(&lastInsertId)
	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func AllCustomers(db *sqlx.DB) ([]*Customer, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}


	rows, err := db.Queryx("SELECT cid, first_name, email, phone FROM customer")
	if err != nil {
		print("error")
	}

	customers := make([]*Customer, 0)
	for rows.Next() {
		bk := new(Customer)
		err := rows.Scan(&bk.customerId, &bk.firstName, &bk.secondName, &bk.phoneNumber)
		if err != nil {
			return nil, err
		}
		customers = append(customers, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}

func AllCustomersAuto(db *sqlx.DB) ([]*Customer, error) {

	customers := []*Customer{}
	db.Select(&customers, "SELECT cid, first_name, email, phone FROM customers ORDER BY first_name ASC")

	return customers, nil
}

func StoreRealCustomer(db *sqlx.DB, customerRequest *pb.CustomerRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

		//customer_id SERIAL PRIMARY KEY NOT NULL,
		//customer_image_path varchar (400),
		//first_name varchar (400),
		//second_name varchar (400),
		//phone_number varchar (400),
		//address varchar (400)


	err := tx.QueryRow("INSERT INTO customers(customer_image_path, first_name, second_name, phone_number, address) VALUES($1, $2, $3, $4, $5) returning customer_id;",
		customerRequest.CustomerImagePath,
		customerRequest.FirstName,
		customerRequest.SecondName,
		customerRequest.PhoneNumber,
		customerRequest.Address).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
}

func AllRealCustomers(db *sqlx.DB) ([]*pb.CustomerRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT customer_id, customer_image_path, first_name, second_name, phone_number, address FROM customers ORDER BY customer_id ASC")
	if err != nil {
		print("error")
	}

	realCustomers := make([]*pb.CustomerRequest, 0)
	for rows.Next() {
		customer := new(pb.CustomerRequest)
		err := rows.Scan(&customer.CustomerId, &customer.CustomerImagePath, &customer.FirstName, &customer.SecondName, &customer.PhoneNumber, &customer.Address)
		if err != nil {
			return nil, err
		}
		realCustomers = append(realCustomers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return realCustomers, nil
}







