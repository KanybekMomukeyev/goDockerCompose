package model

import (
	"github.com/jmoiron/sqlx"
	"log"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"fmt"
	"errors"
)

var schemaRemoveAccount = `
DROP TABLE IF EXISTS accounts;
`

var schemaCreateAccount = `
CREATE TABLE IF NOT EXISTS accounts (
    account_id BIGSERIAL PRIMARY KEY NOT NULL,
    customer_id BIGINT,
    supplier_id BIGINT,
    balance REAL
);
`

var schemaCreateIndexForAccount1 = `CREATE INDEX IF NOT EXISTS customer_id_accounts_idx ON accounts (customer_id)`
var schemaCreateIndexForAccount2 = `CREATE INDEX IF NOT EXISTS supplier_id_accounts_idx ON accounts (supplier_id)`

type Account struct {
	accountId uint64 `db:"account_id"`
	supplierId uint64 `db:"supplier_id"`
	categoryId uint64 `db:"category_id"`
	balance float32 `db:"balance"`
}

func CreateAccountIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveAccount)
	db.MustExec(schemaCreateAccount)
	db.MustExec(schemaCreateIndexForAccount1)
	db.MustExec(schemaCreateIndexForAccount2)
}

func StoreAccount(db *sqlx.DB, accountRequest *pb.AccountRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO accounts " +
		"(customer_id, supplier_id, balance) " +
		"VALUES($1, $2, $3) returning account_id;",
		accountRequest.CustomerId,
		accountRequest.SupplierId,
		accountRequest.Balance).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted account_id =", lastInsertId)

	return lastInsertId, nil
}

func UpdateCustomerBalance(db *sqlx.DB, customerId uint64, balance float64) (uint64, error)  {

	tx := db.MustBegin()

	stmt, err :=tx.Prepare("UPDATE accounts SET balance=$1 WHERE customer_id=$2")
	CheckErr(err)

	res, err2 := stmt.Exec(balance,
		customerId)

	CheckErr(err2)

	affect, err := res.RowsAffected()
	CheckErr(err)

	fmt.Println(affect, "rows changed")

	commitError := tx.Commit()
	CheckErr(commitError)

	return uint64(affect), nil
}

func UpdateSupplierBalance(db *sqlx.DB, supplierId uint64, balance float64) (uint64, error)  {

	tx := db.MustBegin()

	stmt, err :=tx.Prepare("UPDATE accounts SET balance=$1 WHERE supplier_id=$2")
	CheckErr(err)

	res, err2 := stmt.Exec(balance,
		supplierId)

	CheckErr(err2)

	affect, err := res.RowsAffected()
	CheckErr(err)

	fmt.Println(affect, "rows changed")

	commitError := tx.Commit()
	CheckErr(commitError)

	return uint64(affect), nil
}

func AllAccounts(db *sqlx.DB) ([]*pb.AccountRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT account_id, customer_id, supplier_id, balance FROM accounts ORDER BY account_id ASC")
	if err != nil {
		print("error")
	}

	accounts := make([]*pb.AccountRequest, 0)
	for rows.Next() {
		account := new(pb.AccountRequest)
		err := rows.Scan(&account.AccountId, &account.CustomerId, &account.SupplierId, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func AccountFor(db *sqlx.DB, order *pb.OrderRequest) (*pb.AccountRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	var rows *sqlx.Rows
	var err error
	if order.CustomerId > 0 {
		rows, err = db.Queryx("SELECT account_id, customer_id, supplier_id, balance " +
			"FROM accounts WHERE customer_id=$1 ORDER BY account_id ASC LIMIT $2", order.CustomerId, 1)
	} else if order.SupplierId > 0 {
		rows, err = db.Queryx("SELECT account_id, customer_id, supplier_id, balance " +
			"FROM accounts WHERE supplier_id=$1 ORDER BY account_id ASC LIMIT $2", order.SupplierId, 1)
	} else {
		rows, err = db.Queryx("SELECT account_id, customer_id, supplier_id, balance " +
		"FROM accounts WHERE customer_id=$1 ORDER BY account_id ASC LIMIT $2", order.CustomerId, 1)
	}

	if err != nil {
		print("error")
	}

	accounts := make([]*pb.AccountRequest, 0)
	for rows.Next() {
		account := new(pb.AccountRequest)
		err := rows.Scan(&account.AccountId, &account.CustomerId, &account.SupplierId, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(accounts) > 0 {
		return accounts[0], nil
	}

	return nil, errors.New("Not found")
}

func AccountForCustomer(db *sqlx.DB, customerId uint64) (*pb.AccountRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	var rows *sqlx.Rows
	var err error
	rows, err = db.Queryx("SELECT account_id, customer_id, supplier_id, balance " +
		"FROM accounts WHERE customer_id=$1 ORDER BY account_id ASC LIMIT $2", customerId, 1)

	if err != nil {
		print("error")
	}

	accounts := make([]*pb.AccountRequest, 0)
	for rows.Next() {
		account := new(pb.AccountRequest)
		err := rows.Scan(&account.AccountId, &account.CustomerId, &account.SupplierId, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(accounts) > 0 {
		return accounts[0], nil
	}

	return nil, errors.New("Not found")
}

func AccountForSupplier(db *sqlx.DB, supplierId uint64) (*pb.AccountRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	var rows *sqlx.Rows
	var err error
	rows, err = db.Queryx("SELECT account_id, customer_id, supplier_id, balance " +
		"FROM accounts WHERE supplier_id=$1 ORDER BY account_id ASC LIMIT $2", supplierId, 1)

	if err != nil {
		print("error")
	}

	accounts := make([]*pb.AccountRequest, 0)
	for rows.Next() {
		account := new(pb.AccountRequest)
		err := rows.Scan(&account.AccountId, &account.CustomerId, &account.SupplierId, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(accounts) > 0 {
		return accounts[0], nil
	}

	return nil, errors.New("Not found")
}