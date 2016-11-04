package model

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaCustomer = `
CREATE TABLE IF NOT EXISTS customer (
    first_name text,
    email text,
    phone text
);
`

type Customer struct {
	FirstName    string
	Email  string
	Phone string
}

func CreateTableIfNotExsists(db *sqlx.DB) {
	db.MustExec(schemaCustomer)
}

func StoreCustomer(db *sqlx.DB, customer *pb.CustomerRequest) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO customer (first_name, email, phone) VALUES ($1, $2, $3)", customer.Name, customer.Phone, customer.Email)
	return tx.Commit()
}

func AllCustomers(db *sqlx.DB) ([]*Customer, error) {

	rows, err := db.Query("SELECT * FROM customer")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	customers := make([]*Customer, 0)
	for rows.Next() {
		bk := new(Customer)
		err := rows.Scan(&bk.FirstName, &bk.Email, &bk.Phone)
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

