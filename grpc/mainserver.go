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

// CreateCustomer creates a new Customer
// ------------------------------------------------------------ //
func (s *server) CreateExample(ctx context.Context, customerReq *pb.ExampleRequest) (*pb.ExampleResponse, error) {

	s.savedCustomers = append(s.savedCustomers, customerReq)

	unique_key, storeError := model.StoreCustomer(db, customerReq)
	if storeError != nil {
		return nil, storeError
	}
	fmt.Printf("unique_key ==> %#v\n", unique_key)
	return &pb.ExampleResponse{Id: unique_key, Success: true}, nil
}

// GetCustomers returns all customers by given filter
func (s *server) GetExamples(filter *pb.ExampleFilter, stream pb.RentautomationService_GetExamplesServer) error {

	for _, customer := range s.savedCustomers {

		//if filter.Keyword != "" {
		//	if !strings.Contains(customer.Name, filter.Keyword) {
		//		continue
		//	}
		//}
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateAccount(ctx context.Context, accountReq *pb.AccountRequest) (*pb.AccountRequest, error) {

	unique_key, storeError := model.StoreAccount(db, accountReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of staff ==> %v\n", unique_key)
	accountReq.AccountId = unique_key
	fmt.Printf("accountReq.AccountId ==> %v\n", accountReq.AccountId)

	return accountReq, nil
}

func (s *server) GetAccounts(filter *pb.AccountFilter, stream pb.RentautomationService_GetAccountsServer) error {

	accounts, _ := model.AllAccounts(db)
	for _, accountReq := range accounts {
		if err := stream.Send(accountReq); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateCategory(ctx context.Context, categoryReq *pb.CategoryRequest) (*pb.CategoryRequest, error) {

	unique_key, storeError := model.StoreCategory(db, categoryReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of staff ==> %v\n", unique_key)
	categoryReq.CategoryId = unique_key
	fmt.Printf("categoryReq.CategoryId ==> %v\n", categoryReq.CategoryId)

	return categoryReq, nil
}

func (s *server) GetCategories(filter *pb.CategoryFilter, stream pb.RentautomationService_GetCategoriesServer) error {

	categories, _ := model.AllCategory(db)
	for _, categoryReq := range categories {
		if err := stream.Send(categoryReq); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateCustomer(ctx context.Context, customerReq *pb.CustomerRequest) (*pb.CustomerRequest, error) {

	unique_key, storeError := model.StoreRealCustomer(db, customerReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of staff ==> %v\n", unique_key)
	customerReq.CustomerId = unique_key
	fmt.Printf("staffReq.StaffId ==> %v\n", customerReq.CustomerId)

	return customerReq, nil
}

func (s *server) GetCustomers(filter *pb.CustomerFilter, stream pb.RentautomationService_GetCustomersServer) error {

	customers, _ := model.AllRealCustomers(db)
	for _, customerReq := range customers {
		if err := stream.Send(customerReq); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateOrder(ctx context.Context, orderReq *pb.OrderRequest) (*pb.OrderRequest, error) {

	unique_key, storeError := model.StoreOrder(db, orderReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of order ==> %v\n", unique_key)
	orderReq.OrderId = unique_key
	fmt.Printf("orderReq.OrderId ==> %v\n", orderReq.OrderId)

	return orderReq, nil
}

func (s *server) GetOrders(filter *pb.OrderFilter, stream pb.RentautomationService_GetOrdersServer) error {

	orders, _ := model.AllOrders(db)
	for _, orderReq := range orders {
		if err := stream.Send(orderReq); err != nil {
			return err
		}
	}

	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateOrderDetail(ctx context.Context, orderDetailReq *pb.OrderDetailRequest) (*pb.OrderDetailRequest, error) {

	unique_key, storeError := model.StoreOrderDetails(db, orderDetailReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of orderDetail ==> %v\n", unique_key)
	orderDetailReq.OrderDetailId = unique_key
	fmt.Printf("orderDetailReq.OrderDetailId ==> %v\n", orderDetailReq.OrderDetailId)

	return orderDetailReq, nil
}

func (s *server) GetOrderDetails(filter *pb.OrderDetailFilter, stream pb.RentautomationService_GetOrderDetailsServer) error {

	orderDetails, _ := model.AllOrderDetails(db)
	for _, orderDetailReq := range orderDetails {
		if err := stream.Send(orderDetailReq); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreatePayment(ctx context.Context, paymentReq *pb.PaymentRequest) (*pb.PaymentRequest, error) {

	unique_key, storeError := model.StorePayment(db, paymentReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of order ==> %v\n", unique_key)
	paymentReq.PaymentId = unique_key
	fmt.Printf("paymentReq.PaymentId ==> %v\n", paymentReq.PaymentId)

	return paymentReq, nil
}

func (s *server) GetPayments(filter *pb.PaymentFilter, stream pb.RentautomationService_GetPaymentsServer) error {
	payments, _ := model.AllPayments(db)
	for _, paymentReq := range payments {
		if err := stream.Send(paymentReq); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateProduct(ctx context.Context, productReq *pb.ProductRequest) (*pb.ProductRequest, error) {

	unique_key, storeError := model.StoreProduct(db, productReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of order ==> %v\n", unique_key)
	productReq.ProductId = unique_key
	fmt.Printf("productReq.ProductId ==> %v\n", productReq.ProductId)

	return productReq, nil
}

func (s *server) GetProducts(filter *pb.ProductFilter, stream pb.RentautomationService_GetProductsServer) error {
	products, _ := model.AllProducts(db)
	for _, productReq := range products {
		if err := stream.Send(productReq); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateStaff(ctx context.Context, staffReq *pb.StaffRequest) (*pb.StaffRequest, error) {

	unique_key, storeError := model.StoreStaff(db, staffReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of staff ==> %v\n", unique_key)
	staffReq.StaffId = unique_key

	s.savedStaff = append(s.savedStaff, staffReq)

	fmt.Printf("staffReq.StaffId ==> %v\n", staffReq.StaffId)
	return staffReq, nil
}

func (s *server) GetStaff(filter *pb.StaffFilter, stream pb.RentautomationService_GetStaffServer) error {

	//staff, _ := model.AllStaffAuto(db)
	staff, _ := model.AllStaff(db)

	for _, staffRequest := range staff {
		if err := stream.Send(staffRequest); err != nil {
			return err
		}
	}
	return nil
}

// ------------------------------------------------------------ //
func (s *server) CreateTransaction(ctx context.Context, transactionReq *pb.TransactionRequest) (*pb.TransactionRequest, error) {

	unique_key, storeError := model.StoreTransaction(db, transactionReq)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("unique_key of transaction ==> %v\n", unique_key)
	transactionReq.TransactionId = unique_key
	fmt.Printf("transactionReq.TransactionId ==> %v\n", transactionReq.TransactionId)

	return transactionReq, nil
}

func (s *server) GetTransactions(filter *pb.TransactionFilter, stream pb.RentautomationService_GetTransactionsServer) error {
	transactions, _ := model.AllTransactions(db)
	for _, transactionReq := range transactions {
		if err := stream.Send(transactionReq); err != nil {
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

	model.CreateAccountIfNotExsists(db)
	model.CreateCategoryIfNotExsists(db)
	model.CreateCustomerIfNotExsists(db)

	model.CreateOrderIfNotExsists(db)
	model.CreateOrderDetailsIfNotExsists(db)
	model.CreatePaymentIfNotExsists(db)

	model.CreateProductIfNotExsists(db)
	model.CreateStaffIfNotExsists(db)
	model.CreateTransactionIfNotExsists(db)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterRentautomationServiceServer(s, &server{})
	s.Serve(lis)
}