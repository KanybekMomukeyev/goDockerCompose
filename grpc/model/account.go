package model

import (
	"github.com/jmoiron/sqlx"
	"log"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"fmt"
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

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
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

