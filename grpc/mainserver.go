package main

import (
	"log"
	"net"
	//"strings"

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
	savedCustomers []*pb.ExampleRequest
	savedStaff []*pb.StaffRequest
}

func (s *server) CreateStaff(ctx context.Context, staffReq *pb.StaffRequest) (*pb.StaffRequest, error) {
	s.savedStaff = append(s.savedStaff, staffReq)
	fmt.Printf("unique_key ==> %#v\n", staffReq.StaffId)
	return staffReq, nil
}

func (s *server) GetStaff(filter *pb.StaffFilter, stream pb.RentautomationService_GetStaffServer) error {

	for _, staff := range s.savedStaff {
		if err := stream.Send(staff ); err != nil {
			return err
		}
	}
	return nil
}


// CreateCustomer creates a new Customer
func (s *server) CreateExample(ctx context.Context, customerReq *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	s.savedCustomers = append(s.savedCustomers, customerReq)
	unique_key, storeError := model.StoreCustomer2(db, customerReq)
	if storeError != nil {
		return nil, storeError
	}
	fmt.Printf("unique_key ==> %#v\n", unique_key)
	return &pb.ExampleResponse{Id: unique_key, Success: true}, nil
}

// GetCustomers returns all customers by given filter

func (s *server) GetExamples(filter *pb.ExampleFilter, stream pb.RentautomationService_GetExamplesServer) error {

	customers, _ := model.AllCustomers(db)
	print("KANOOOOOO")

	//customers, _ := model.AllCustomersAuto(db)

	for _, customer := range customers {
		fmt.Printf("%#v\n", customer)

		exampleRequest := &pb.ExampleRequest{
			Id:    customer.CustomerId,
			Name:  customer.FirstName,
			Email: customer.Email,
			Phone: customer.Phone,
		}

		if err := stream.Send(exampleRequest); err != nil {
			return err
		}
	}

	//for _, customer := range s.savedCustomers {
	//
	//	if filter.Keyword != "" {
	//		if !strings.Contains(customer.Name, filter.Keyword) {
	//			continue
	//		}
	//	}
	//	if err := stream.Send(customer); err != nil {
	//		return err
	//	}
	//}
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
	pb.RegisterRentautomationServiceServer(s, &server{})
	s.Serve(lis)
}