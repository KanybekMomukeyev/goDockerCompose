package model

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"fmt"
	"log"
)

var schemaCustomerDelete = `
DROP TABLE IF EXISTS customer;
`

var schemaCustomer = `
CREATE TABLE IF NOT EXISTS customer (
    cid serial PRIMARY KEY NOT NULL,
    first_name text,
    email text,
    phone text
);
`
type Customer2 struct {
	CustomerId uint64 `db:"cid"`
	FirstName    string `db:"first_name"`
	Email  string
	Phone string
}

func CreateTableIfNotExsists(db *sqlx.DB) {
	// for some migrations
	//db.MustExec(schemaCustomerDelete)
	db.MustExec(schemaCustomer)
}

func StoreCustomer(db *sqlx.DB, customer *pb.ExampleRequest) (uint64, error)  {

	tx := db.MustBegin()

	var lastInsertId uint64
	err := tx.QueryRow("INSERT INTO customer (first_name, phone, email) VALUES($1, $2, $3) returning cid;", customer.Name, customer.Phone, customer.Email).Scan(&lastInsertId)
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

func AllCustomers(db *sqlx.DB) ([]*Customer2, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}


	rows, err := db.Queryx("SELECT cid, first_name, email, phone FROM customer")
	if err != nil {
		print("error")
	}

	customers := make([]*Customer2, 0)
	for rows.Next() {
		bk := new(Customer2)
		err := rows.Scan(&bk.CustomerId, &bk.FirstName, &bk.Email, &bk.Phone)
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

func AllCustomersAuto(db *sqlx.DB) ([]*Customer2, error) {

	customers := []*Customer2{}
	db.Select(&customers, "SELECT cid, first_name, email, phone FROM customer ORDER BY first_name ASC")

	return customers, nil
}
