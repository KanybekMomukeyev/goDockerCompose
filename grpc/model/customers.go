package model

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"fmt"
	"database/sql"
	"log"
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

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

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

	lastId, commitError := result.RowsAffected()
	//lastId, commitError:= result.LastInsertId()

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
	//
	//pingError := db.Ping()
	//
	//if pingError != nil {
	//	log.Fatalln(pingError)
	//	print(pingError)
	//}
	//
	//// Query the database, storing results in a []Person (wrapped in []interface{})
	//people := []Person{}
	//db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	//jason, john := people[0], people[1]
	//
	//fmt.Printf("%#v\n%#v", jason, john)

	//// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
	//// Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}
	//
	//// You can also get a single result, a la QueryRow
	//jason = Person{}
	//err := db.Get(&jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
	//fmt.Printf("%#v\n", jason)
	//// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
	//
	//// if you have null fields and use SELECT *, you must use sql.Null* in your struct
	places := []Place{}
	err := db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
	if err != nil {
		fmt.Println(err)
	}
	//usa, singsing, honkers := places[0], places[1], places[2]
	//
	//fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
	//// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	//// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
	//// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
	//
	//// Loop through rows using only one struct
	//place := Place{}
	//rows, err := db.Queryx("SELECT * FROM place")
	//for rows.Next() {
	//	err := rows.StructScan(&place)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	fmt.Printf("%#v\n", place)
	//}
	//// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	//// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
	//// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
	//
	//// Named queries, using `:name` as the bindvar.  Automatic bindvar support
	//// which takes into account the dbtype based on the driverName on sqlx.Open/Connect
	//_, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`,
	//	map[string]interface{}{
	//		"first": "Bin",
	//		"last": "Smuth",
	//		"email": "bensmith@allblacks.nz",
	//	})
	//
	//// Selects Mr. Smith from the database
	//rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "Bin"})
	//
	//// Named queries can also use structs.  Their bind names follow the same rules
	//// as the name -> db mapping, so struct fields are lowercased and the `db` tag
	//// is taken into consideration.
	//rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)


	// ----------

	//rows, err = db.Query("SELECT * FROM customer")
	//tx := db.MustBegin()
	//tx.MustExec("INSERT INTO customer (first_name, email, phone) VALUES ($1, $2, $3) RETURNING cid", "KOO", "123123", "ko@mail.sru")
	//tx.MustExec("INSERT INTO customer (first_name, email, phone) VALUES ($1, $2, $3) RETURNING cid", "KOO", "123123", "ko@mail.sru")
	//tx.MustExec("INSERT INTO customer (first_name, email, phone) VALUES ($1, $2, $3) RETURNING cid", "KOO", "123123", "ko@mail.sru")
	//
	//commitError := tx.Commit()
	//if commitError != nil {
	//	print("error")
	//	print(commitError)
	//}
	//
	//var (
	//	first_name string
	//	email string
	//	phone string
	//)
	//
	//rows_sql, err := db.Query("SELECT first_name, email, phone FROM customer")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows_sql.Close()
	//for rows_sql.Next() {
	//	err := rows_sql.Scan(&first_name, &email, &phone)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Println(first_name, email, phone)
	//}
	//err = rows_sql.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}



	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}


	rows, err := db.Queryx("SELECT first_name, email, phone FROM customer")
	if err != nil {
		print("error")
	}

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
