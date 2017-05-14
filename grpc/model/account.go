package model

import (
	"github.com/jmoiron/sqlx"
	log "github.com/Sirupsen/logrus"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
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

func DeleteAccountIfNotExsists(db *sqlx.DB) {
	db.MustExec(schemaRemoveAccount)
}

func CreateAccountIfNotExsists(db *sqlx.DB) {
	db.MustExec(schemaCreateAccount)
	db.MustExec(schemaCreateIndexForAccount1)
	db.MustExec(schemaCreateIndexForAccount2)
}

func ErrorFunc(err error) (uint64, error) {
	log.WithFields(log.Fields{ "error": err}).Fatal("QueryRow breaks")
	panic(err)
	return 0, err
}

func StoreAccount(tx *sqlx.Tx, accountRequest *pb.AccountRequest) (uint64, error)  {

	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO accounts " +
		"(customer_id, supplier_id, balance) " +
		"VALUES($1, $2, $3) returning account_id;",
		accountRequest.CustomerId,
		accountRequest.SupplierId,
		accountRequest.Balance).Scan(&lastInsertId)

	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{"last inserted account_id":  lastInsertId}).Info("")
	return lastInsertId, nil
}

func UpdateCustomerBalance(tx *sqlx.Tx, customerId uint64, balance float64) (uint64, error)  {

	stmt, err := tx.Prepare("UPDATE accounts SET balance=$1 WHERE customer_id=$2")
	if err != nil {
		return ErrorFunc(err)
	}

	res, err := stmt.Exec(balance, customerId)
	if err != nil {
		return ErrorFunc(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"update customer balance rows changed":  affect,
	}).Info("")
	return uint64(affect), nil
}

func UpdateSupplierBalance(tx *sqlx.Tx, supplierId uint64, balance float64) (uint64, error)  {

	stmt, err := tx.Prepare("UPDATE accounts SET balance=$1 WHERE supplier_id=$2")
	if err != nil {
		return ErrorFunc(err)
	}

	res, err := stmt.Exec(balance, supplierId)
	if err != nil {
		return ErrorFunc(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"update supplier balance rows changed":  affect,
	}).Info("")
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

	log.WithFields(log.Fields{"order.OrderId": order.OrderId}).Warn("")
	return nil, errors.New("Not found AccountFor")
}

func scanAccountRows(rows *sqlx.Rows) ([]*pb.AccountRequest, error) {
	accounts := make([]*pb.AccountRequest, 0)
	for rows.Next() {
		account := new(pb.AccountRequest)
		err := rows.Scan(&account.AccountId,
				&account.CustomerId,
				&account.SupplierId,
				&account.Balance)
		if err != nil {
			log.WithFields(log.Fields{"scanAccountRows":err,}).Warn("ERROR")
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
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

	log.WithFields(log.Fields{"customerId": customerId}).Warn("")
	return nil, errors.New("Not found AccountForCustomer")
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

	log.WithFields(log.Fields{"supplierId": supplierId}).Warn("")
	return nil, errors.New("Not found AccountForSupplier")
}