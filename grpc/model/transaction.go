package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"fmt"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveTransaction = `
DROP TABLE IF EXISTS transactions;
`

var schemaCreateTransaction = `
CREATE TABLE IF NOT EXISTS transactions (
    	transaction_id BIGSERIAL PRIMARY KEY NOT NULL,
    	transaction_date BIGINT,
    	is_last_transaction INTEGER,
    	transaction_type INTEGER,
    	money_amount REAL,
    	order_id BIGINT,
    	customer_id BIGINT,
    	supplier_id BIGINT,
  	staff_id BIGINT
);
`

var schemaCreateIndexForTransaction1 = `CREATE INDEX IF NOT EXISTS order_id_transactions_idx ON transactions (order_id)`
var schemaCreateIndexForTransaction2 = `CREATE INDEX IF NOT EXISTS customer_id_transactions_idx ON transactions (customer_id)`
var schemaCreateIndexForTransaction3 = `CREATE INDEX IF NOT EXISTS supplier_id_transactions_idx ON transactions (supplier_id)`
var schemaCreateIndexForTransaction4 = `CREATE INDEX IF NOT EXISTS staff_id_transactions_idx ON transactions (staff_id)`

type Transaction struct {
	transactionId uint64 `db:"product_id"`
	transactionDate uint64 `db:"transaction_date"`
	transactionType uint64 `db:"transaction_type"`
	moneyAmount float32 `db:"money_amount"`
	orderId uint64 `db:"order_id"`
	customerId uint64 `db:"customer_id"`
	supplierId uint64 `db:"supplier_id"`
	staffId uint64 `db:"staff_id"`
}

func CreateTransactionIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveTransaction)
	db.MustExec(schemaCreateTransaction)
	db.MustExec(schemaCreateIndexForTransaction1)
	db.MustExec(schemaCreateIndexForTransaction2)
	db.MustExec(schemaCreateIndexForTransaction3)
	db.MustExec(schemaCreateIndexForTransaction4)
}

func StoreTransaction(db *sqlx.DB, transaction *pb.TransactionRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO transactions (transaction_date, is_last_transaction, transaction_type, money_amount, order_id, customer_id, supplier_id, staff_id) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8) returning transaction_id;",
		transaction.TransactionDate,
		transaction.IsLastTransaction,
		transaction.TransactionType,
		transaction.MoneyAmount,
		transaction.OrderId,
		transaction.CustomerId,
		transaction.SupplierId,
		transaction.StaffId).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted transaction_id =", lastInsertId)

	return lastInsertId, nil
}

func AllTransactions(db *sqlx.DB) ([]*pb.TransactionRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT transaction_id, transaction_date, is_last_transaction, transaction_type, money_amount, " +
		"order_id, customer_id, supplier_id, staff_id " +
		"FROM transactions ORDER BY transaction_id DESC")

	if err != nil {
		print("error")
	}

	transactions := make([]*pb.TransactionRequest, 0)
	for rows.Next() {
		transaction := new(pb.TransactionRequest)
		err := rows.Scan(&transaction.TransactionId, &transaction.TransactionDate, &transaction.IsLastTransaction, &transaction.TransactionType,
			&transaction.MoneyAmount, &transaction.OrderId, &transaction.CustomerId, &transaction.SupplierId,
			&transaction.StaffId)

		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}