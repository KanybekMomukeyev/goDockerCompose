package main

import (
	log "github.com/Sirupsen/logrus"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"google.golang.org/grpc/credentials"

	"github.com/jmoiron/sqlx"
	"github.com/KanybekMomukeyev/goDockerCompose/grpc/model"

	"flag"
	"io"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"runtime"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var (
	certFile   = flag.String("cert_file", "certfiles/ssl.crt", "The TLS cert file")
	keyFile    = flag.String("key_file", "certfiles/ssl.key", "The TLS key file")
)

const (
	port = ":50051"
)

// server is used to implement customer.CustomerServer.
type server struct {
	savedCustomers []*pb.ExampleRequest
}

func isAuthorized(ctx context.Context) error{

	md, ok := metadata.FromIncomingContext(ctx)

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
func (s *server) CreateExample(ctx context.Context, customerReq *pb.ExampleRequest) (*pb.ExampleResponse, error) {

	//authorizeError := isAuthorized(ctx)
	//if authorizeError != nil {
	//	return nil, authorizeError
	//}

	s.savedCustomers = append(s.savedCustomers, customerReq)

	//tx := db.MustBegin()
	//unique_key, err := model.StoreCustomer(tx, customerReq)
	//if err != nil {
	//	tx.Rollback()
	//	log.WithFields(log.Fields{"err": err}).Warn("")
	//	return nil, err
	//}
	//
	//err = tx.Commit()
	//if err != nil {
	//	log.WithFields(log.Fields{"err": err}).Warn("")
	//	return nil, err
	//}

	return &pb.ExampleResponse{Id: 101, Success: true}, nil
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

// -------------------------- CATEGORY ---------------------------------- //
func (s *server) CreateCategory(ctx context.Context, categoryReq *pb.CategoryRequest) (*pb.CategoryRequest, error) {

	tx := db.MustBegin()
	unique_key, err := model.StoreCategory(tx, categoryReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	categoryReq.CategoryId = unique_key

	return categoryReq, nil
}

func (s *server) UpdateCategory(ctx context.Context, categoryReq *pb.CategoryRequest) (*pb.CategoryRequest, error) {

	tx := db.MustBegin()
	_, err := model.UpdateCategory(tx, categoryReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return categoryReq, nil
}

func (s *server) AllCategoriesForInitial(ctx context.Context, catFilter *pb.CategoryFilter) (*pb.AllCategoryResponse, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	categories, error := model.AllCategory(db)
	if error != nil {
		return nil, error
	}
	allCategoryResponse := new(pb.AllCategoryResponse)
	allCategoryResponse.CategoryRequest = categories

	return allCategoryResponse, nil
}

func (s *server) CheckCategoriesForUpdate(ctx context.Context, catFilter *pb.CategoryFilter) (*pb.AllCategoryResponse, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	categories, error := model.AllUpdatedCategories(db, catFilter)
	if error != nil {
		return nil, error
	}
	allCategoryResponse := new(pb.AllCategoryResponse)
	allCategoryResponse.CategoryRequest = categories

	return allCategoryResponse, nil
}

// ---------------------------- ADDITTION -------------------------------- //
func (s *server) CreateProductWith(ctx context.Context, createPrReq *pb.CreateProductRequest) (*pb.CreateProductRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	productSerialKey, err := model.StoreProduct(tx, createPrReq.Product)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createPrReq.Product.ProductId = productSerialKey
	createPrReq.OrderDetail.ProductId = productSerialKey

	orderDetailSerialKey, err := model.StoreOrderDetails(tx, createPrReq.OrderDetail)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createPrReq.OrderDetail.OrderDetailId = orderDetailSerialKey

	return createPrReq, nil
}

func (s *server) UpdateProductWith(ctx context.Context, createPrReq *pb.CreateProductRequest) (*pb.CreateProductRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	_, err := model.UpdateProduct(tx, createPrReq.Product)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	orderDetailSerialKey, err := model.StoreOrderDetails(tx, createPrReq.OrderDetail)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createPrReq.OrderDetail.OrderDetailId = orderDetailSerialKey
	return createPrReq, nil
}

func (s *server) AllProductsForInitial(ctx context.Context, prFilter *pb.ProductFilter) (*pb.AllProductsResponse, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createProdRequests := make([]*pb.CreateProductRequest, 0)
	products, _ := model.AllProducts(db)

	for _, productReq := range products {

		orderDetReq, err := model.RecentOrderDetailForProduct(db, productReq)
		if err != nil {
			break
			return nil, err
		}

		createProductReq := new(pb.CreateProductRequest)
		createProductReq.OrderDetail = orderDetReq
		createProductReq.Product = productReq

		createProdRequests = append(createProdRequests, createProductReq)
	}

	allProdResponse := new(pb.AllProductsResponse)
	allProdResponse.ProductRequest = createProdRequests

	return allProdResponse,nil
}

func (s *server) CheckProductsForUpdate(ctx context.Context, prFilter *pb.ProductFilter) (*pb.AllProductsResponse, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createProdRequests := make([]*pb.CreateProductRequest, 0)
	products, _ := model.AllProductsForUpdate(db, prFilter)

	for _, productReq := range products {

		orderDetReq, err := model.RecentOrderDetailForProduct(db, productReq)
		if err != nil {
			break
			return nil, err
		}

		createProductReq := new(pb.CreateProductRequest)
		createProductReq.OrderDetail = orderDetReq
		createProductReq.Product = productReq

		createProdRequests = append(createProdRequests, createProductReq)
	}

	allProdResponse := new(pb.AllProductsResponse)
	allProdResponse.ProductRequest = createProdRequests

	return allProdResponse,nil
}

func (s *server) CheckOrderDetailsForUpdate(ctx context.Context, oDetFilter *pb.OrderDetailFilter) (*pb.AllOrderDetailResponse, error) {
	return nil, nil
}

func (s *server) CheckTransactionsForUpdate(ctx context.Context, transFilter *pb.TransactionFilter) (*pb.AllTransactionResponse, error) {

	log.WithFields(log.Fields{"rpc CheckTransactionsForUpdate": transFilter, }).Info("")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	transactions, err := model.AllUpdatedTransactions(db, transFilter)
	if err != nil {
		return nil, err
	}

	allTransReq := new(pb.AllTransactionResponse)
	allTransReq.TransactionRequest = transactions

	log.WithFields(log.Fields{"transactions count":  len(transactions), }).Info("")

	return allTransReq, nil
}

func (s *server) UpdateTransactionWith(ctx context.Context, transReq *pb.TransactionRequest) (*pb.TransactionRequest, error) {

	log.WithFields(log.Fields{"transReq": transReq, }).Info("UpdateTransaction RPC")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	_, err := model.UpdateTransaction(tx, transReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return transReq, nil
}

func (s *server) AllOrderDetails(ctx context.Context, oDetFilter *pb.OrderDetailFilter) (*pb.AllOrderDetailResponse, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	orderDetails, error := model.AllOrderDetailsForFilter(db, oDetFilter)
	if error != nil {
		return nil, error
	}

	allOrderDetailResponse := new(pb.AllOrderDetailResponse)
	allOrderDetailResponse.OrderDetailRequest = orderDetails

	return allOrderDetailResponse,nil
}

// ----------------------------  -------------------------------- //
func (s *server) CreateCustomerWith(ctx context.Context, createCustReq *pb.CreateCustomerRequest) (*pb.CreateCustomerRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	customerSerial, err := model.StoreRealCustomer(tx, createCustReq.Customer)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createCustReq.Customer.CustomerId = customerSerial
	createCustReq.Transaction.CustomerId = customerSerial
	createCustReq.Account.CustomerId = customerSerial

	transactionSerial, err := model.StoreTransaction(tx, createCustReq.Transaction)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createCustReq.Transaction.TransactionId = transactionSerial

	accountSerial, err := model.StoreAccount(tx, createCustReq.Account)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createCustReq.Account.AccountId = accountSerial

	return createCustReq, nil
}

func (s *server) UpdateCustomerWith(ctx context.Context, createCustReq *pb.CustomerRequest) (*pb.CustomerRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	_, err := model.UpdateRealCustomer(tx, createCustReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return createCustReq, nil
}

func (s *server) UpdateCustomerBalanceWith(ctx context.Context, updateCustBalanceReq *pb.CreateCustomerRequest) (*pb.CreateCustomerRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	transactionSerial, err := model.StoreTransaction(tx, updateCustBalanceReq.Transaction)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	updateCustBalanceReq.Transaction.TransactionId = transactionSerial

	rowsAffected, err := model.UpdateCustomerBalance(tx, updateCustBalanceReq.Account.CustomerId, updateCustBalanceReq.Account.Balance)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err, "rowsAffected":rowsAffected}).Warn("")
		return nil, err
	}

	return updateCustBalanceReq, nil
}

func (s *server) AllCustomersForInitial(ctx context.Context, custFilter *pb.CustomerFilter) (*pb.AllCustomersResponse, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createCustomerRequests := make([]*pb.CreateCustomerRequest, 0)
	customers, _ := model.AllRealCustomers(db)

	for _, customerReq := range customers {

		transactionReq, err := model.RecentTransactionForCustomer(db, customerReq)
		if err != nil {
			break
			return nil, err
		}

		accountReq, err := model.AccountForCustomer(db, customerReq.CustomerId)
		if err != nil {
			break
			return nil, err
		}

		createCustomerRequest := new(pb.CreateCustomerRequest)
		createCustomerRequest.Customer = customerReq
		createCustomerRequest.Transaction = transactionReq
		createCustomerRequest.Account = accountReq

		createCustomerRequests = append(createCustomerRequests, createCustomerRequest)
	}

	allCustomersResponse := new(pb.AllCustomersResponse)
	allCustomersResponse.CustomerRequest = createCustomerRequests

	return allCustomersResponse, nil
}

func (s *server) CheckCustomersForUpdate(ctx context.Context, custFilter *pb.CustomerFilter) (*pb.AllCustomersResponse, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createCustomerRequests := make([]*pb.CreateCustomerRequest, 0)
	customers, _ := model.AllUpdatedCustomers(db, custFilter)

	for _, customerReq := range customers {

		transactionReq, err := model.RecentTransactionForCustomer(db, customerReq)
		if err != nil {
			break
			return nil, err
		}

		accountReq, err := model.AccountForCustomer(db, customerReq.CustomerId)
		if err != nil {
			break
			return nil, err
		}

		createCustomerRequest := new(pb.CreateCustomerRequest)
		createCustomerRequest.Customer = customerReq
		createCustomerRequest.Transaction = transactionReq
		createCustomerRequest.Account = accountReq

		createCustomerRequests = append(createCustomerRequests, createCustomerRequest)
	}

	allCustomersResponse := new(pb.AllCustomersResponse)
	allCustomersResponse.CustomerRequest = createCustomerRequests

	return allCustomersResponse, nil
}

// ----------------------------  -------------------------------- //
func (s *server) CreateSupplierWith(ctx context.Context, createSuppReq *pb.CreateSupplierRequest) (*pb.CreateSupplierRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}
	tx := db.MustBegin()
	supplierSerial, err := model.StoreSupplier(tx, createSuppReq.Supplier)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createSuppReq.Supplier.SupplierId = supplierSerial
	createSuppReq.Transaction.SupplierId = supplierSerial
	createSuppReq.Account.SupplierId = supplierSerial

	transactionSerial, err := model.StoreTransaction(tx, createSuppReq.Transaction)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createSuppReq.Transaction.TransactionId = transactionSerial

	accountSerial, err := model.StoreAccount(tx, createSuppReq.Account)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createSuppReq.Account.AccountId= accountSerial
	return createSuppReq, nil
}

func (s *server) UpdateSupplierWith(ctx context.Context, createSuppReq *pb.SupplierRequest) (*pb.SupplierRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	_, err := model.UpdateSupplier(tx, createSuppReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return createSuppReq, nil
}

func (s *server) UpdateSupplierBalanceWith(ctx context.Context, createSuppReq *pb.CreateSupplierRequest) (*pb.CreateSupplierRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	transactionSerial, err := model.StoreTransaction(tx, createSuppReq.Transaction)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	createSuppReq.Transaction.TransactionId = transactionSerial

	rowsAffected, err := model.UpdateSupplierBalance(tx, createSuppReq.Account.SupplierId, createSuppReq.Account.Balance)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err, "rowsAffected":rowsAffected}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return createSuppReq, nil
}

func (s *server) AllSuppliersForInitial(ctx context.Context, suppFilter *pb.SupplierFilter) (*pb.AllSuppliersResponse, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createSupplierRequests := make([]*pb.CreateSupplierRequest, 0)
	suppliers, _ := model.AllSuppliers(db)

	for _, supplierReq := range suppliers {

		transactionReq, err := model.RecentTransactionForSupplier(db, supplierReq)
		if err != nil {
			break
			return nil, err
		}

		accountReq, err := model.AccountForSupplier(db, supplierReq.SupplierId)
		if err != nil {
			break
			return nil, err
		}

		createSupplierRequest := new(pb.CreateSupplierRequest)
		createSupplierRequest.Supplier = supplierReq
		createSupplierRequest.Transaction = transactionReq
		createSupplierRequest.Account = accountReq

		createSupplierRequests = append(createSupplierRequests, createSupplierRequest)
	}

	allSuppliersResponse := new(pb.AllSuppliersResponse)
	allSuppliersResponse.SupplierRequest = createSupplierRequests

	return allSuppliersResponse, nil
}

func (s *server) CheckSuppliersForUpdate(ctx context.Context, suppFilter *pb.SupplierFilter) (*pb.AllSuppliersResponse, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createSupplierRequests := make([]*pb.CreateSupplierRequest, 0)
	suppliers, _ := model.AllSuppliersForUpdate(db, suppFilter)

	for _, supplierReq := range suppliers {

		transactionReq, err := model.RecentTransactionForSupplier(db, supplierReq)
		if err != nil {
			break
			return nil, err
		}

		accountReq, err := model.AccountForSupplier(db, supplierReq.SupplierId)
		if err != nil {
			break
			return nil, err
		}

		createSupplierRequest := new(pb.CreateSupplierRequest)
		createSupplierRequest.Supplier = supplierReq
		createSupplierRequest.Transaction = transactionReq
		createSupplierRequest.Account = accountReq

		createSupplierRequests = append(createSupplierRequests, createSupplierRequest)
	}

	allSuppliersResponse := new(pb.AllSuppliersResponse)
	allSuppliersResponse.SupplierRequest = createSupplierRequests

	return allSuppliersResponse, nil
}

func (s *server) AllTransactionsForInitial(ctx context.Context, transFilter *pb.TransactionFilter) (*pb.AllTransactionResponse, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	transactions, error := model.AllTransactionsForFilter(db, transFilter)
	if error != nil {
		return nil, error
	}

	allTransactionResponse := new(pb.AllTransactionResponse)
	allTransactionResponse.TransactionRequest = transactions

	return allTransactionResponse, nil
}

// ----------------------------  -------------------------------- //
func (s *server) CreateStaffWith(ctx context.Context, staffReq *pb.StaffRequest) (*pb.StaffRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	staffSerialKey, err := model.StoreStaff(tx, staffReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	staffReq.StaffId = staffSerialKey

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return staffReq, nil
}

func (s *server) UpdateStaffWith(ctx context.Context, staffReq *pb.StaffRequest) (*pb.StaffRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	_, err := model.UpdateStaff(tx, staffReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return staffReq, nil
}

func (s *server) AllStaffForInitial(ctx context.Context, staffFilter *pb.StaffFilter) (*pb.AllStaffResponse, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	staff, error := model.AllStaff(db)
	if error != nil {
		return nil, error
	}
	allStaffResponse := new(pb.AllStaffResponse)
	allStaffResponse.StaffRequest = staff

	return allStaffResponse, nil
}

func (s *server) CheckStaffForUpdate(ctx context.Context, staffFilter *pb.StaffFilter) (*pb.AllStaffResponse, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	staff, error := model.AllStaffForUpdate(db, staffFilter)
	if error != nil {
		return nil, error
	}
	allStaffResponse := new(pb.AllStaffResponse)
	allStaffResponse.StaffRequest = staff

	return allStaffResponse, nil
}

func (s *server) SignInWith(ctx context.Context, signInReq *pb.SignInRequest) (*pb.StaffRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	stafReq, selectError := model.SignIn(db, signInReq)
	if selectError != nil {
		return nil, selectError
	}
	return  stafReq, nil
}

func (s *server) UpdateStream(stream pb.RentautomationService_UpdateStreamServer) error {
	for {
		stickyNoteReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		resp := &pb.StickyNoteResponse{}
		resp.Message = stickyNoteReq.Message

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) CreateOrderWith(ctx context.Context, creatOrdReq *pb.CreateOrderRequest) (*pb.CreateOrderRequest, error) {

	log.WithFields(log.Fields{
		"creat_order_req": creatOrdReq.Order,
	}).Info("CreateOrderWith rpc method called")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	//log.WithFields(log.Fields{"payment transaction begin": 1, }).Info("")
	//time.Sleep(1 * time.Second)

	// payment
	tx := db.MustBegin()

	paymentSerial, err := model.StorePayment(tx, creatOrdReq.Payment)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	creatOrdReq.Payment.PaymentId = paymentSerial
	creatOrdReq.Order.PaymentId = paymentSerial


	//log.WithFields(log.Fields{"order transaction begin": 1, }).Info("")
	//time.Sleep(6 * time.Second)
	// order
	orderSerial, err := model.StoreOrder(tx, creatOrdReq.Order)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	creatOrdReq.Order.OrderId = orderSerial
	if creatOrdReq.Transaction != nil {
		creatOrdReq.Transaction.OrderId = orderSerial
	}


	//log.WithFields(log.Fields{"transaction transaction begin": 1, }).Info("")
	//time.Sleep(6 * time.Second)
	// transaction
	if creatOrdReq.Transaction != nil {
		transactionSerial, err := model.StoreTransaction(tx, creatOrdReq.Transaction)
		if err != nil {
			tx.Rollback()
			log.WithFields(log.Fields{"err": err}).Warn("")
			return nil, err
		}
		creatOrdReq.Transaction.TransactionId = transactionSerial
	}


	//log.WithFields(log.Fields{"account transaction begin": 1, }).Info("")
	//time.Sleep(6 * time.Second)

	// order document, to update customer/supplier balance
	// also update product amount in stock
	var orderDocument string = ""
	if creatOrdReq.Order.OrderDocument == 0 {
		orderDocument = ".none"

	} else if creatOrdReq.Order.OrderDocument == 1000 {
		orderDocument = ".productOrderSaledToCustomer"

		// customer balance
		acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
		if err != nil {
			return nil, err
		}
		acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
		model.UpdateCustomerBalance(tx, creatOrdReq.Order.CustomerId, acountReq.Balance)

		// product amount
		_, error := model.DecreaseProductsInStock(db, tx, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 2000 {
		orderDocument = ".productOrderSaleEditedToCustomer"

		// customer balance
		acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
		if err != nil {
			return nil, err
		}
		acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
		model.UpdateCustomerBalance(tx, creatOrdReq.Order.CustomerId, acountReq.Balance)

		// product amount
		_, error := model.IncreaseProductsInStock(db, tx, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 3000 {
		orderDocument = ".productOrderReceivedFromSupplier"

		// supplier balance
		acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
		if err != nil {
			return nil, err
		}
		acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
		model.UpdateSupplierBalance(tx, creatOrdReq.Order.SupplierId, acountReq.Balance)

		// product amount
		_, error := model.IncreaseProductsInStock(db, tx, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 4000 {
		orderDocument = ".productOrderReceiveEditedFromSupplier"

		// supplier balance
		acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
		if err != nil {
			return nil, err
		}
		acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
		model.UpdateSupplierBalance(tx, creatOrdReq.Order.SupplierId, acountReq.Balance)

		// product amount
		_, error := model.DecreaseProductsInStock(db, tx, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 5000 {
		orderDocument = ".productReturnedFromCustomer"

		// customer balance
		acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
		if err != nil {
			return nil, err
		}
		acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
		model.UpdateCustomerBalance(tx, creatOrdReq.Order.CustomerId, acountReq.Balance)

		// product amount
		_, error := model.IncreaseProductsInStock(db, tx, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 6000 {
		orderDocument = ".productReturneEditedFromCustomer"

	} else if creatOrdReq.Order.OrderDocument == 5500 {
		orderDocument = ".productReturnedToSupplier"

		// supplier balance
		acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
		if err != nil {
			return nil, err
		}
		acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
		model.UpdateSupplierBalance(tx, creatOrdReq.Order.SupplierId, acountReq.Balance)

		// product amount
		_, error := model.DecreaseProductsInStock(db, tx, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 6600 {
		orderDocument = ".productReturneEditedToSupplier"

	} else if creatOrdReq.Order.OrderDocument == 7000 {
		orderDocument = ".moneyReceived"

		if creatOrdReq.Order.CustomerId > 0 {

			// customer balance
			acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateCustomerBalance(tx, creatOrdReq.Order.CustomerId, acountReq.Balance)

		} else if creatOrdReq.Order.SupplierId > 0 {

			// supplier balance
			acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateSupplierBalance(tx, creatOrdReq.Order.SupplierId, acountReq.Balance)
		}
	} else if creatOrdReq.Order.OrderDocument == 8000 {
		orderDocument = ".moneyReceiveEdited"

		if creatOrdReq.Order.CustomerId > 0 {

			// customer balance
			acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateCustomerBalance(tx, creatOrdReq.Order.CustomerId, acountReq.Balance)

		} else if creatOrdReq.Order.SupplierId > 0 {

			// supplier balance
			acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateSupplierBalance(tx, creatOrdReq.Order.SupplierId, acountReq.Balance)
		}


	} else if creatOrdReq.Order.OrderDocument == 10000 {
		orderDocument = ".moneyGone"

		if creatOrdReq.Order.CustomerId > 0 {
			// customer balance
			acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateCustomerBalance(tx, creatOrdReq.Order.CustomerId, acountReq.Balance)

		} else if creatOrdReq.Order.SupplierId > 0 {
			// supplier balance
			acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateSupplierBalance(tx, creatOrdReq.Order.SupplierId, acountReq.Balance)
		}


	} else if creatOrdReq.Order.OrderDocument == 11000 {
		orderDocument = ".moneyGoneEdited"

		if creatOrdReq.Order.CustomerId > 0 {
			// customer balance
			acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateCustomerBalance(tx, creatOrdReq.Order.CustomerId, acountReq.Balance)

		} else if creatOrdReq.Order.SupplierId > 0 {
			// supplier balance
			acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateSupplierBalance(tx, creatOrdReq.Order.SupplierId, acountReq.Balance)
		}

	} else if creatOrdReq.Order.OrderDocument == 12000 {
		orderDocument = "customerMadePreOrder"
	} else if creatOrdReq.Order.OrderDocument == 13000 {
		orderDocument = "stokTaking"
	}

	// orderDetails
	for _, orderDetailReq := range creatOrdReq.OrderDetails {
		orderDetailReq.OrderId = orderSerial

		orderDetailSerial, err := model.StoreOrderDetails(tx, orderDetailReq)
		if err != nil {
			tx.Rollback()
			log.WithFields(log.Fields{"err": err}).Warn("")
			break
			return nil, err
		}

		orderDetailReq.OrderDetailId = orderDetailSerial
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	log.WithFields(log.Fields{"orderDocument": orderDocument}).Info("")
	return creatOrdReq, nil
}

func (s *server) UpdateOrderWith(ctx context.Context, orderReq *pb.OrderRequest) (*pb.OrderRequest, error) {

	log.WithFields(log.Fields{"orderReq": orderReq, }).Info("UpdateOrderWith RPC")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	tx := db.MustBegin()
	_, err := model.UpdateOrder(tx, orderReq)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	contexErr := ctx.Err()
	if contexErr != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": contexErr}).Warn("")
		return nil, contexErr
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return nil, err
	}

	return orderReq, nil
}

func (s *server) AllOrdersForInitial(ctx context.Context, orderFilter *pb.OrderFilter) (*pb.AllOrderResponse, error) {

	log.WithFields(log.Fields{"OrderDate": orderFilter.OrderDate, }).Info("AllOrdersForInitial")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createOrderRequests := make([]*pb.CreateOrderRequest, 0)
	orders, err := model.AllOrdersForFilter(db, orderFilter)
	log.WithFields(log.Fields{"initial len(orders):": len(orders),}).Info("")

	if err != nil {
		log.WithFields(log.Fields{"error":err,}).Warn("ERROR")
		return nil, err
	}

	for _, order := range orders {
		createOrderRequest := new(pb.CreateOrderRequest)
		createOrderRequest.Order = order

		payment, error := model.PaymentForOrder(db, order)
		if error != nil {
			log.WithFields(log.Fields{"error":err,}).Warn("ERROR")
			break
			return nil, err
		}
		createOrderRequest.Payment = payment

		if order.CustomerId == 0 && order.SupplierId == 0 {
		} else {

			transaction, error := model.TransactionForOrder(db, order)
			if error != nil {
				log.WithFields(log.Fields{"error": error}).Warn("")
				//break
				//return nil, err
			}

			if transaction != nil {
				createOrderRequest.Transaction = transaction
			}

			//account, error := model.AccountForOrder(db, order)
			//if error != nil {
			//	break
			//	return nil, err
			//}
			//createOrderRequest.Account = account
		}

		orderDetails, error := model.AllOrderDetailsForOrder(db, order)
		if error != nil {
			log.WithFields(log.Fields{"error": error}).Warn("")
			break
			return nil, err
		}
		createOrderRequest.OrderDetails = orderDetails
		createOrderRequests = append(createOrderRequests, createOrderRequest)
	}

	allOrderResponse := new(pb.AllOrderResponse)
	allOrderResponse.OrderRequest = createOrderRequests

	log.WithFields(log.Fields{"found count": len(createOrderRequests), }).Info("")

	return allOrderResponse, nil
}

func (s *server) CheckOrdersForUpdate(ctx context.Context, orderFilter *pb.OrderFilter) (*pb.AllOrderResponse, error) {

	log.WithFields(log.Fields{"rpc allOrdersForRecent order filter": orderFilter, }).Info("")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createOrderRequests := make([]*pb.CreateOrderRequest, 0)
	orders, err := model.AllOrdersForRecentFilter(db, orderFilter)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		createOrderRequest := new(pb.CreateOrderRequest)
		createOrderRequest.Order = order

		payment, error := model.PaymentForOrder(db, order)
		if error != nil {
			break
			return nil, err
		}
		createOrderRequest.Payment = payment

		if order.CustomerId == 0 && order.SupplierId == 0 {

		} else {
			transaction, error := model.TransactionForOrder(db, order)
			if error != nil {
				log.WithFields(log.Fields{"error": error}).Warn("")
				//break
				//return nil, err
			}
			if transaction != nil {
				createOrderRequest.Transaction = transaction
			}

			//account, error := model.AccountForOrder(db, order)
			//if error != nil {
			//	break
			//	return nil, err
			//}
			//createOrderRequest.Account = account
		}

		orderDetails, error := model.AllOrderDetailsForOrder(db, order)
		if error != nil {
			break
			return nil, err
		}
		createOrderRequest.OrderDetails = orderDetails
		createOrderRequests = append(createOrderRequests, createOrderRequest)
	}

	allOrderResponse := new(pb.AllOrderResponse)
	allOrderResponse.OrderRequest = createOrderRequests

	log.WithFields(log.Fields{"found count":  len(createOrderRequests), }).Info("")

	return allOrderResponse, nil
}

var db *sqlx.DB

func main() {

	var databaseError error
	db, databaseError = model.NewDB("datasource")
	if databaseError != nil {
		log.WithFields(log.Fields{
			"omg":    databaseError,
			"number": 100,
		}).Fatal("failed to listen:")
	}

	//db.SetMaxIdleConns()
	//db.SetMaxOpenConns()

	log.WithFields(log.Fields{
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other": "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")

	runtime.GOMAXPROCS(4)

	model.CreateStaffIfNotExsists(db)
	model.CreateAccountIfNotExsists(db)
	model.CreateCategoryIfNotExsists(db)
	model.CreateCustomerIfNotExsists(db)
	model.CreateSupplierIfNotExsists(db)

	model.CreateOrderIfNotExsists(db)
	model.CreateOrderDetailsIfNotExsists(db)
	model.CreatePaymentIfNotExsists(db)

	model.CreateProductIfNotExsists(db)
	model.CreateTransactionIfNotExsists(db)

	var err error
	var lis net.Listener
	var grpcServer *grpc.Server
	if false {
		lis, err = net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer = grpc.NewServer()
		println("Listen on " + port + " without TLS")
	} else {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		lis, err = net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
		println("Listen on " + port + " with TLS")
	}

	pb.RegisterRentautomationServiceServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}