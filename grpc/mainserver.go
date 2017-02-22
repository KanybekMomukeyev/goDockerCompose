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

	"fmt"
	"flag"
	"io"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

var (
	certFile   = flag.String("cert_file", "certfiles/server.crt", "The TLS cert file")
	keyFile    = flag.String("key_file", "certfiles/server.key", "The TLS key file")
)

const (
	port = ":50051"
)

// server is used to implement customer.CustomerServer.
type server struct {
	savedCustomers []*pb.ExampleRequest
	savedStaff []*pb.StaffRequest
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
func (s *server) CreateExample(ctx context.Context, customerReq *pb.ExampleRequest) (*pb.ExampleResponse, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

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

// -------------------------- CATEGORY ---------------------------------- //
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

func (s *server) UpdateCategory(ctx context.Context, categoryReq *pb.CategoryRequest) (*pb.CategoryRequest, error) {

	rowsAffected, updateError := model.UpdateCategory(db, categoryReq)
	if updateError != nil {
		return nil, updateError
	}
	fmt.Printf("rowsAffected UpdateCategory==> %v\n", rowsAffected)
	return categoryReq, nil
}

// ---------------------------- ADDITTION -------------------------------- //
func (s *server) CreateProductWith(ctx context.Context, createPrReq *pb.CreateProductRequest) (*pb.CreateProductRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	productSerialKey, storeError := model.StoreProduct(db, createPrReq.Product)
	if storeError != nil {
		return nil, storeError
	}

	createPrReq.Product.ProductId = productSerialKey
	createPrReq.OrderDetail.ProductId = productSerialKey

	orderDetailSerialKey, storeError2 := model.StoreOrderDetails(db, createPrReq.OrderDetail)
	if storeError2 != nil {
		return nil, storeError2
	}
	createPrReq.OrderDetail.OrderDetailId = orderDetailSerialKey

	fmt.Printf("CreateProductWith of transaction ==> %v\n", &createPrReq )
	return createPrReq, nil
}

func (s *server) UpdateProductWith(ctx context.Context, createPrReq *pb.CreateProductRequest) (*pb.CreateProductRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	_, updateError := model.UpdateProduct(db, createPrReq.Product)
	if updateError != nil {
		return nil, updateError
	}

	orderDetailSerialKey, storeError2 := model.StoreOrderDetails(db, createPrReq.OrderDetail)
	if storeError2 != nil {
		return nil, storeError2
	}
	createPrReq.OrderDetail.OrderDetailId = orderDetailSerialKey

	fmt.Printf("UpdateProductWith of transaction ==> %v\n", &createPrReq )
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

	customerSerial, storeError := model.StoreRealCustomer(db, createCustReq.Customer)
	if storeError != nil {
		return nil, storeError
	}
	createCustReq.Customer.CustomerId = customerSerial
	createCustReq.Transaction.CustomerId = customerSerial
	createCustReq.Account.CustomerId = customerSerial

	transactionSerial, storeError := model.StoreTransaction(db, createCustReq.Transaction)
	if storeError != nil {
		return nil, storeError
	}
	createCustReq.Transaction.TransactionId = transactionSerial

	accountSerial, storeError := model.StoreAccount(db, createCustReq.Account)
	if storeError != nil {
		return nil, storeError
	}
	createCustReq.Account.AccountId= accountSerial

	fmt.Printf("CreateCustomerWith of transaction ==> %v\n", &createCustReq )
	return createCustReq, nil
}

func (s *server) UpdateCustomerWith(ctx context.Context, createCustReq *pb.CustomerRequest) (*pb.CustomerRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	rowsAffected, updateError := model.UpdateRealCustomer(db, createCustReq)
	if updateError != nil {
		return nil, updateError
	}
	fmt.Printf("rowsAffected ==> %v\n", rowsAffected)
	return createCustReq, nil
}

func (s *server) UpdateCustomerBalanceWith(ctx context.Context, updateCustBalanceReq *pb.CreateCustomerRequest) (*pb.CreateCustomerRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	transactionSerial, storeError := model.StoreTransaction(db, updateCustBalanceReq.Transaction)
	if storeError != nil {
		return nil, storeError
	}
	updateCustBalanceReq.Transaction.TransactionId = transactionSerial

	rowsAffected, storeError := model.UpdateCustomerBalance(db, updateCustBalanceReq.Account.CustomerId, updateCustBalanceReq.Account.Balance)

	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("rowsAffected %v\n", &rowsAffected)
	fmt.Printf("UpdateCustomerBalanceWith of transaction ==> %v\n", &updateCustBalanceReq)
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

// ----------------------------  -------------------------------- //
func (s *server) CreateSupplierWith(ctx context.Context, createSuppReq *pb.CreateSupplierRequest) (*pb.CreateSupplierRequest, error) {
	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	supplierSerial, storeError := model.StoreSupplier(db, createSuppReq.Supplier)
	if storeError != nil {
		return nil, storeError
	}
	createSuppReq.Supplier.SupplierId = supplierSerial
	createSuppReq.Transaction.SupplierId = supplierSerial
	createSuppReq.Account.SupplierId = supplierSerial

	transactionSerial, storeError := model.StoreTransaction(db, createSuppReq.Transaction)
	if storeError != nil {
		return nil, storeError
	}
	createSuppReq.Transaction.TransactionId = transactionSerial

	accountSerial, storeError := model.StoreAccount(db, createSuppReq.Account)
	if storeError != nil {
		return nil, storeError
	}
	createSuppReq.Account.AccountId= accountSerial

	fmt.Printf("unique_key of transaction ==> %v\n", &createSuppReq )
	return createSuppReq, nil
}

func (s *server) UpdateSupplierWith(ctx context.Context, createSuppReq *pb.SupplierRequest) (*pb.SupplierRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	rowsAffected, updateError := model.UpdateSupplier(db, createSuppReq)
	if updateError != nil {
		return nil, updateError
	}
	fmt.Printf("rowsAffected UpdateSupplier==> %v\n", rowsAffected)
	return createSuppReq, nil
}

func (s *server) UpdateSupplierBalanceWith(ctx context.Context, createSuppReq *pb.CreateSupplierRequest) (*pb.CreateSupplierRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	transactionSerial, storeError := model.StoreTransaction(db, createSuppReq.Transaction)
	if storeError != nil {
		return nil, storeError
	}
	createSuppReq.Transaction.TransactionId = transactionSerial

	rowsAffected, storeError := model.UpdateSupplierBalance(db, createSuppReq.Account.SupplierId, createSuppReq.Account.Balance)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("rowsAffected %v\n", &rowsAffected)
	fmt.Printf("UpdateSupplierBalanceWith of transaction ==> %v\n", &createSuppReq)
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

	staffSerialKey, storeError := model.StoreStaff(db, staffReq)
	if storeError != nil {
		return nil, storeError
	}

	staffReq.StaffId = staffSerialKey
	fmt.Printf("CreateStaffWith of transaction ==> %v\n", &staffReq )
	return staffReq, nil
}

func (s *server) UpdateStaffWith(ctx context.Context, staffReq *pb.StaffRequest) (*pb.StaffRequest, error) {

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	rowsAffected, updateError := model.UpdateStaff(db, staffReq)
	if updateError != nil {
		return nil, updateError
	}
	fmt.Printf("rowsAffected UpdateStaffWith==> %v\n", rowsAffected)
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
		"creat_order_req": creatOrdReq,
	}).Info("CreateOrderWith rpc method called")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	// payment
	paymentSerial, storeError := model.StorePayment(db, creatOrdReq.Payment)
	if storeError != nil {
		return nil, storeError
	}
	creatOrdReq.Payment.PaymentId = paymentSerial
	creatOrdReq.Order.PaymentId = paymentSerial

	// order
	orderSerial, storeError := model.StoreOrder(db, creatOrdReq.Order)
	if storeError != nil {
		return nil, storeError
	}
	creatOrdReq.Order.OrderId = orderSerial
	creatOrdReq.Transaction.OrderId = orderSerial

	// transaction
	transactionSerial, storeError := model.StoreTransaction(db, creatOrdReq.Transaction)
	if storeError != nil {
		return nil, storeError
	}
	creatOrdReq.Transaction.TransactionId = transactionSerial

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
		model.UpdateCustomerBalance(db, creatOrdReq.Order.CustomerId, acountReq.Balance)

		// product amount
		_, error := model.DecreaseProductsInStock(db, creatOrdReq.OrderDetails)
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
		model.UpdateCustomerBalance(db, creatOrdReq.Order.CustomerId, acountReq.Balance)

		// product amount
		_, error := model.IncreaseProductsInStock(db, creatOrdReq.OrderDetails)
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
		model.UpdateSupplierBalance(db, creatOrdReq.Order.SupplierId, acountReq.Balance)

		// product amount
		_, error := model.IncreaseProductsInStock(db, creatOrdReq.OrderDetails)
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
		model.UpdateSupplierBalance(db, creatOrdReq.Order.SupplierId, acountReq.Balance)

		// product amount
		_, error := model.DecreaseProductsInStock(db, creatOrdReq.OrderDetails)
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
		model.UpdateCustomerBalance(db, creatOrdReq.Order.CustomerId, acountReq.Balance)

		// product amount
		_, error := model.IncreaseProductsInStock(db, creatOrdReq.OrderDetails)
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
		model.UpdateSupplierBalance(db, creatOrdReq.Order.SupplierId, acountReq.Balance)

		// product amount
		_, error := model.DecreaseProductsInStock(db, creatOrdReq.OrderDetails)
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
			model.UpdateCustomerBalance(db, creatOrdReq.Order.CustomerId, acountReq.Balance)

		} else if creatOrdReq.Order.SupplierId > 0 {

			// supplier balance
			acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance - creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateSupplierBalance(db, creatOrdReq.Order.SupplierId, acountReq.Balance)
		}


	} else if creatOrdReq.Order.OrderDocument == 8000 {
		orderDocument = ".moneyReceiveEdited"


	} else if creatOrdReq.Order.OrderDocument == 10000 {
		orderDocument = ".moneyGone"

		if creatOrdReq.Order.CustomerId > 0 {
			// customer balance
			acountReq, err := model.AccountForCustomer(db, creatOrdReq.Order.CustomerId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateCustomerBalance(db, creatOrdReq.Order.CustomerId, acountReq.Balance)

		} else if creatOrdReq.Order.SupplierId > 0 {
			// supplier balance
			acountReq, err := model.AccountForSupplier(db, creatOrdReq.Order.SupplierId)
			if err != nil {
				return nil, err
			}
			acountReq.Balance = acountReq.Balance + creatOrdReq.Payment.TotalPriceWithDiscount
			model.UpdateSupplierBalance(db, creatOrdReq.Order.SupplierId, acountReq.Balance)
		}


	} else if creatOrdReq.Order.OrderDocument == 11000 {
		orderDocument = ".moneyGoneEdited"
	}

	println(orderDocument)

	// orderDetails
	for _, orderDetailReq := range creatOrdReq.OrderDetails {
		orderDetailReq.OrderId = orderSerial

		orderDetailSerial, storeError := model.StoreOrderDetails(db, orderDetailReq)
		if storeError != nil {
			return nil, storeError
		}
		orderDetailReq.OrderDetailId = orderDetailSerial
	}

	return creatOrdReq, nil
}

func (s *server) AllOrdersForInitial(ctx context.Context, orderFilter *pb.OrderFilter) (*pb.AllOrderResponse, error) {

	log.WithFields(log.Fields{
		"orderFilter": orderFilter,
	}).Info("AllOrdersForInitial rpc method called")

	authorizeError := isAuthorized(ctx)
	if authorizeError != nil {
		return nil, authorizeError
	}

	createOrderRequests := make([]*pb.CreateOrderRequest, 0)
	orders, err := model.AllOrdersForFilter(db, orderFilter)
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

		transaction, error := model.TransactionForOrder(db, order)
		if error != nil {
			break
			return nil, err
		}
		createOrderRequest.Transaction = transaction

		//account, error := model.AccountForOrder(db, order)
		//if error != nil {
		//	break
		//	return nil, err
		//}
		//createOrderRequest.Account = account

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

	log.WithFields(log.Fields{
		"createOrderRequests length":  len(createOrderRequests),
	}).Info("AllOrdersForInitial Found orders")

	return allOrderResponse, nil
}

func (s *server) AllOrdersForRecent(ctx context.Context, orderFilter *pb.OrderFilter) (*pb.AllOrderResponse, error) {

	log.WithFields(log.Fields{
		"orderFilter": orderFilter,
	}).Info("AllOrdersForRecent rpc method called")

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

		transaction, error := model.TransactionForOrder(db, order)
		if error != nil {
			break
			return nil, err
		}
		createOrderRequest.Transaction = transaction

		//account, error := model.AccountForOrder(db, order)
		//if error != nil {
		//	break
		//	return nil, err
		//}
		//createOrderRequest.Account = account

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

	log.WithFields(log.Fields{
		"createOrderRequests length":  len(createOrderRequests),
	}).Info("AllOrdersForRecent Found orders")

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
		"animal": "walrus",
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
	if true {
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