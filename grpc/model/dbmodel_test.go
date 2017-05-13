package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var db *sqlx.DB

func init() {
	db, _ = NewTestDB("datasource")
	DeleteAccountIfNotExsists(db)
	CreateAccountIfNotExsists(db)
}

func TestSomething(t *testing.T) {

	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(t, 123, 456, "they should not be equal")

	// assert for nil (good for errors)
	//assert.Nil(t, dbMng)

	// assert for not nil (good when you expect something)
	if assert.NotNil(t, db) {

	}
}

func TestAccountStore(t *testing.T) {

	accountReq := new(pb.AccountRequest)
	accountReq.AccountId = 100
	accountReq.CustomerId = 1100
	accountReq.SupplierId = 1200
	accountReq.Balance = 2000

	tx := db.MustBegin()

	countId, err := StoreAccount(tx, accountReq)
	assert.Nil(t, err)

	assert.Equal(t, countId, uint64(1), "")

	err = tx.Commit()
	assert.Nil(t, err)

	savedAccReq, err := AccountForCustomer(db, 1100)
	assert.Nil(t, err)

	assert.Equal(t, savedAccReq.CustomerId, uint64(1100), "")
	assert.Equal(t, savedAccReq.SupplierId, uint64(1200), "")
	assert.Equal(t, savedAccReq.Balance, float64(2000), "")
	assert.Equal(t, savedAccReq.AccountId, uint64(1), "")
}