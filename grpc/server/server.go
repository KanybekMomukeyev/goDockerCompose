package main

import (
	"log"
	"net"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"

	"github.com/jmoiron/sqlx"
	"github.com/KanybekMomukeyev/goDockerCompose/grpc/model"

	"fmt"
)

const (
	port = ":50051"
)

// server is used to implement customer.CustomerServer.
type server struct {
	savedCustomers []*pb.CustomerRequest
}

// CreateCustomer creates a new Customer
func (s *server) CreateCustomer(ctx context.Context, customerReq *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	s.savedCustomers = append(s.savedCustomers, customerReq)
	unique_key, storeError := model.StoreCustomer2(db, customerReq)
	if storeError != nil {
		return nil, storeError
	}
	fmt.Printf("unique_key ==> %#v\n", unique_key)
	return &pb.CustomerResponse{Id: unique_key, Success: true}, nil
}

// GetCustomers returns all customers by given filter
func (s *server) GetCustomers(filter *pb.CustomerFilter, stream pb.CustomerService_GetCustomersServer) error {

	//customers, _ := model.AllCustomers(db)
	customers, _ := model.AllCustomersAuto(db)

	for _, customer := range customers {
		fmt.Printf("%#v\n", customer)
	}

	for _, customer := range s.savedCustomers {
		if filter.Keyword != "" {
			if !strings.Contains(customer.Name, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil
}

var db *sqlx.DB

func main() {

	var databaseError error
	db, databaseError = model.NewDB("datasource")
	if databaseError != nil {
		log.Fatalf("failed to listen: %v", databaseError)
	}
	model.CreateTableIfNotExsists(db)


	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterCustomerServiceServer(s, &server{})
	s.Serve(lis)
}