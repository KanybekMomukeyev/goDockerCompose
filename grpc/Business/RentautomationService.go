package Business

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type RentautomationService struct {
	databaseManager DatabaseManager
}

func NewRentautomationService(dbManager DatabaseManager) (*RentautomationService, error) {
	return &RentautomationService{dbManager}, nil
}

func isAuthorized(ctx context.Context) error{
	md, ok := metadata.FromContext(ctx)

	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "Token unsetted")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "Token is uncorrect: appid=%s, appkey=%s", appid, appkey)
	}

	return nil
}

// CreateCustomer creates a new Customer
// ------------------------------------------------------------ //
func (s *RentautomationService) CreateExample(ctx context.Context, customerReq *pb.ExampleRequest) (*pb.ExampleResponse, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	return s.databaseManager.CreateExample(customerReq)
}

// GetCustomers returns all customers by given filter
func (s *RentautomationService) GetExamples(filter *pb.ExampleFilter, stream pb.RentautomationService_GetExamplesServer) error {
	return s.databaseManager.GetExamples(filter, stream)
}
