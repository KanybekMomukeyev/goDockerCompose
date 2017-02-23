package Business

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"github.com/KanybekMomukeyev/goDockerCompose/grpc/model"
	log "github.com/Sirupsen/logrus"
)

type DatabaseManager struct {
	database *sqlx.DB
}

func (db *DatabaseManager) Speak() {
	return "?????"
}
func (db *DatabaseManager) CreateExample( customerReq *pb.ExampleRequest) (*pb.ExampleResponse, error) {

	tx := db.database.MustBegin()
	unique_key, err := model.StoreCustomer(tx, customerReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return &pb.ExampleResponse{Id: unique_key, Success: true}, nil
}

// GetCustomers returns all customers by given filter
func (db *DatabaseManager) GetExamples(filter *pb.ExampleFilter, stream pb.RentautomationService_GetExamplesServer) error {
	return nil
}

func (db *DatabaseManager) CreateCategory(categoryReq *pb.CategoryRequest) (*pb.CategoryRequest, error) {

	tx := db.database.MustBegin()
	unique_key, err := model.StoreCategory(tx, categoryReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	categoryReq.CategoryId = unique_key

	return categoryReq, nil
}

func (db *DatabaseManager) UpdateCategory(categoryReq *pb.CategoryRequest) (*pb.CategoryRequest, error) {

	tx := db.database.MustBegin()
	_, err := model.UpdateCategory(tx, categoryReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return categoryReq, nil
}