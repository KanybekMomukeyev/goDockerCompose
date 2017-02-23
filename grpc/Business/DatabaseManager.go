package Business

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

type DatabaseManager struct {
	database *sqlx.DB
}

func (db *DatabaseManager) Speak() {
	return "?????"
}

func (db *DatabaseManager) CreateCustomer(createCustReq *pb.CreateCustomerRequest) (*pb.CreateCustomerRequest, error) {

	//customerSerial, storeError := model.StoreRealCustomer(db, createCustReq.Customer)
	//if storeError != nil {
	//	return nil, storeError
	//}
	//createCustReq.Customer.CustomerId = customerSerial
	//createCustReq.Transaction.CustomerId = customerSerial
	//createCustReq.Account.CustomerId = customerSerial
	//
	//transactionSerial, storeError := model.StoreTransaction(db, createCustReq.Transaction)
	//if storeError != nil {
	//	return nil, storeError
	//}
	//createCustReq.Transaction.TransactionId = transactionSerial
	//
	//accountSerial, storeError := db.StoreAccount(db, createCustReq.Account)
	//if storeError != nil {
	//	return nil, storeError
	//}
	//createCustReq.Account.AccountId= accountSerial
	//
	//fmt.Printf("CreateCustomerWith of transaction ==> %v\n", &createCustReq )
	//return createCustReq, nil
	return nil, nil
}


func (db *DatabaseManager) StoreAccount(tx *sqlx.Tx, accountRequest *pb.AccountRequest) (uint64, error)  {

	var lastInsertAccountId uint64

	err := tx.QueryRow("INSERT INTO accounts " +
		"(customer_id, supplier_id, balance) " +
		"VALUES($1, $2, $3) returning account_id;",
		accountRequest.CustomerId,
		accountRequest.SupplierId,
		accountRequest.Balance).Scan(&lastInsertAccountId)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	fmt.Println("last inserted account_id =", lastInsertAccountId)
	return lastInsertAccountId, nil
}