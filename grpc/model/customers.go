package model

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"fmt"
)

var schemaCustomerDelete = `
DROP TABLE IF EXISTS customer;
`

//try self.database.executeUpdate("CREATE TABLE IF NOT EXISTS \(self.staffTableName)
//(staff_id INTEGER PRIMARY KEY, role_id INTEGER, staff_image_path TEXT, first_name TEXT,
//second_name TEXT, email TEXT, password TEXT, phone_number TEXT, address TEXT)", values: nil)

var schemaCustomer = `
CREATE TABLE IF NOT EXISTS customer (
    cid serial PRIMARY KEY NOT NULL,
    first_name text,
    email text,
    phone text
);

CREATE TABLE IF NOT EXISTS person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE IF NOT EXISTS place (
    country text,
    city text NULL,
    telcode integer
);

`

type Customer struct {
	FirstName    string `db:"first_name"`
	Email  string
	Phone string
}

func CreateTableIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaCustomerDelete)
	db.MustExec(schemaCustomer)
}

func StoreCustomer(db *sqlx.DB, customer *pb.ExampleRequest) (uint64, error) {
	tx := db.MustBegin()
	result := tx.MustExec("INSERT INTO customer (first_name, email, phone) VALUES ($1, $2, $3) RETURNING cid", customer.Name, customer.Phone, customer.Email)

	commitError := tx.Commit()
	if commitError != nil {
		return 0, commitError
	}

	lastId, commitError:= result.LastInsertId()

	if commitError != nil {
		return 0, commitError
	}

	return uint64(lastId), nil
}

func StoreCustomer2(db *sqlx.DB, customer *pb.ExampleRequest) (uint64, error)  {

	tx := db.MustBegin()

	var lastInsertId uint64
	err := tx.QueryRow("INSERT INTO customer (first_name, phone, email) VALUES($1, $2, $3) returning cid;", customer.Name, customer.Phone, customer.Email).Scan(&lastInsertId)
	checkErr(err)

	commitError := tx.Commit()
	checkErr(commitError)

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
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

func AllCustomersAuto(db *sqlx.DB) ([]*Customer, error) {

	customers := []*Customer{}
	db.Select(&customers, "SELECT * FROM customer ORDER BY first_name ASC")

	return customers, nil
}
