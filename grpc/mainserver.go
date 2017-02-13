package main

import (
	"log"
	"net"
	//"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"google.golang.org/grpc/credentials"

	"github.com/jmoiron/sqlx"
	"github.com/KanybekMomukeyev/goDockerCompose/grpc/model"

	"fmt"
	"flag"
	"io"
)

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

// ---------------------------- ACCOUNT -------------------------------- //
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

func (s *server) GetCategories(filter *pb.CategoryFilter, stream pb.RentautomationService_GetCategoriesServer) error {

	categories, _ := model.AllCategory(db)
	for _, categoryReq := range categories {
		if err := stream.Send(categoryReq); err != nil {
			return err
		}
	}
	return nil
}

// ---------------------------- CUSTOMER -------------------------------- //
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

// ---------------------------- ORDER -------------------------------- //
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

// ------------------------- ORDER_DETAIL ----------------------------------- //
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

// -------------------------- PAYMENT ---------------------------------- //
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

// --------------------------- PRODUCT --------------------------------- //
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

// ---------------------------- STAFF -------------------------------- //
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

// ---------------------------- TRANSACTION -------------------------------- //
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



// ---------------------------- ADDITTION -------------------------------- //
func (s *server) CreateProductWith(ctx context.Context, createPrReq *pb.CreateProductRequest) (*pb.CreateProductRequest, error) {

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

	categories, error := model.AllCategory(db)
	if error != nil {
		return nil, error
	}
	allCategoryResponse := new(pb.AllCategoryResponse)
	allCategoryResponse.CategoryRequest = categories

	return allCategoryResponse, nil
}

func (s *server) AllOrderDetails(ctx context.Context, oDetFilter *pb.OrderDetailFilter) (*pb.AllOrderDetailResponse, error) {

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

	rowsAffected, updateError := model.UpdateRealCustomer(db, createCustReq)
	if updateError != nil {
		return nil, updateError
	}
	fmt.Printf("rowsAffected ==> %v\n", rowsAffected)
	return createCustReq, nil
}

func (s *server) UpdateCustomerBalanceWith(ctx context.Context, updateCustBalanceReq *pb.CreateCustomerRequest) (*pb.CreateCustomerRequest, error) {

	transactionSerial, storeError := model.StoreTransaction(db, updateCustBalanceReq.Transaction)
	if storeError != nil {
		return nil, storeError
	}
	updateCustBalanceReq.Transaction.TransactionId = transactionSerial

	rowsAffected, storeError := model.UpdateCustomerBalance(db, updateCustBalanceReq.Account)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("rowsAffected %v\n", &rowsAffected)
	fmt.Printf("UpdateCustomerBalanceWith of transaction ==> %v\n", &updateCustBalanceReq)
	return updateCustBalanceReq, nil
}

func (s *server) AllCustomersForInitial(ctx context.Context, custFilter *pb.CustomerFilter) (*pb.AllCustomersResponse, error) {

	createCustomerRequests := make([]*pb.CreateCustomerRequest, 0)
	customers, _ := model.AllRealCustomers(db)

	for _, customerReq := range customers {

		transactionReq, err := model.RecentTransactionForCustomer(db, customerReq)
		if err != nil {
			break
			return nil, err
		}

		accountReq, err := model.AccountForCustomer(db, customerReq)
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

	rowsAffected, updateError := model.UpdateSupplier(db, createSuppReq)
	if updateError != nil {
		return nil, updateError
	}
	fmt.Printf("rowsAffected UpdateSupplier==> %v\n", rowsAffected)
	return createSuppReq, nil
}

func (s *server) UpdateSupplierBalanceWith(ctx context.Context, createSuppReq *pb.CreateSupplierRequest) (*pb.CreateSupplierRequest, error) {

	transactionSerial, storeError := model.StoreTransaction(db, createSuppReq.Transaction)
	if storeError != nil {
		return nil, storeError
	}
	createSuppReq.Transaction.TransactionId = transactionSerial

	rowsAffected, storeError := model.UpdateSupplierBalance(db, createSuppReq.Account)
	if storeError != nil {
		return nil, storeError
	}

	fmt.Printf("rowsAffected %v\n", &rowsAffected)
	fmt.Printf("UpdateSupplierBalanceWith of transaction ==> %v\n", &createSuppReq)
	return createSuppReq, nil
}

func (s *server) AllSuppliersForInitial(ctx context.Context, suppFilter *pb.SupplierFilter) (*pb.AllSuppliersResponse, error) {

	createSupplierRequests := make([]*pb.CreateSupplierRequest, 0)
	suppliers, _ := model.AllSuppliers(db)

	for _, supplierReq := range suppliers {

		transactionReq, err := model.RecentTransactionForSupplier(db, supplierReq)
		if err != nil {
			break
			return nil, err
		}

		accountReq, err := model.AccountForSupplier(db, supplierReq)
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

	staffSerialKey, storeError := model.StoreStaff(db, staffReq)
	if storeError != nil {
		return nil, storeError
	}

	staffReq.StaffId = staffSerialKey
	fmt.Printf("CreateStaffWith of transaction ==> %v\n", &staffReq )
	return staffReq, nil
}

func (s *server) UpdateStaffWith(ctx context.Context, staffReq *pb.StaffRequest) (*pb.StaffRequest, error) {

	rowsAffected, updateError := model.UpdateStaff(db, staffReq)
	if updateError != nil {
		return nil, updateError
	}
	fmt.Printf("rowsAffected UpdateStaffWith==> %v\n", rowsAffected)
	return staffReq, nil
}

func (s *server) AllStaffForInitial(ctx context.Context, staffFilter *pb.StaffFilter) (*pb.AllStaffResponse, error) {

	staff, error := model.AllStaff(db)
	if error != nil {
		return nil, error
	}
	allStaffResponse := new(pb.AllStaffResponse)
	allStaffResponse.StaffRequest = staff

	return allStaffResponse, nil
}

func (s *server) SignInWith(ctx context.Context, signInReq *pb.SignInRequest) (*pb.StaffRequest, error) {
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
		updateBalanceOfCustomer(creatOrdReq.Account)

		// product amount
		_, error := model.DecreaseProductsInStock(db, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 2000 {
		orderDocument = ".productOrderSaleEditedToCustomer"
		updateBalanceOfCustomer(creatOrdReq.Account)

		// product amount
		_, error := model.IncreaseProductsInStock(db, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 3000 {
		orderDocument = ".productOrderReceivedFromSupplier"
		updateBalanceOfSupplier(creatOrdReq.Account)

		// product amount
		_, error := model.IncreaseProductsInStock(db, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 4000 {
		orderDocument = ".productOrderReceiveEditedFromSupplier"
		updateBalanceOfSupplier(creatOrdReq.Account)

		// product amount
		_, error := model.DecreaseProductsInStock(db, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 5000 {
		orderDocument = ".productReturnedFromCustomer"
		updateBalanceOfCustomer(creatOrdReq.Account)

		// product amount
		_, error := model.IncreaseProductsInStock(db, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 6000 {
		orderDocument = ".productReturneEditedFromCustomer"
		updateBalanceOfCustomer(creatOrdReq.Account)


	} else if creatOrdReq.Order.OrderDocument == 5500 {
		orderDocument = ".productReturnedToSupplier"
		updateBalanceOfSupplier(creatOrdReq.Account)

		// product amount
		_, error := model.DecreaseProductsInStock(db, creatOrdReq.OrderDetails)
		if error != nil {
			return nil, error
		}

	} else if creatOrdReq.Order.OrderDocument == 6600 {
		orderDocument = ".productReturneEditedToSupplier"
		updateBalanceOfSupplier(creatOrdReq.Account)


	} else if creatOrdReq.Order.OrderDocument == 7000 {
		orderDocument = ".moneyReceived"

		if creatOrdReq.Order.CustomerId > 0 {
			updateBalanceOfCustomer(creatOrdReq.Account)
		} else if creatOrdReq.Order.SupplierId > 0 {
			updateBalanceOfSupplier(creatOrdReq.Account)
		}


	} else if creatOrdReq.Order.OrderDocument == 8000 {
		orderDocument = ".moneyReceiveEdited"


	} else if creatOrdReq.Order.OrderDocument == 10000 {
		orderDocument = ".moneyGone"

		if creatOrdReq.Order.CustomerId > 0 {
			updateBalanceOfCustomer(creatOrdReq.Account)
		} else if creatOrdReq.Order.SupplierId > 0 {
			updateBalanceOfSupplier(creatOrdReq.Account)
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

		account, error := model.AccountForOrder(db, order)
		if error != nil {
			break
			return nil, err
		}
		createOrderRequest.Account = account

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

	return allOrderResponse, nil
}

func (s *server) AllOrdersForRecent(ctx context.Context, orderFilter *pb.OrderFilter) (*pb.AllOrderResponse, error) {

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

		account, error := model.AccountForOrder(db, order)
		if error != nil {
			break
			return nil, err
		}
		createOrderRequest.Account = account

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

	return allOrderResponse, nil
}

func updateBalanceOfCustomer(accountReq *pb.AccountRequest) (uint64, error) {
	// customer balance /// ---------------
	rowsAffected, storeError := model.UpdateCustomerBalance(db, accountReq)
	if storeError != nil {
		return 0, storeError
	}
	fmt.Printf("rowsAffected %v\n", &rowsAffected)
	return rowsAffected, nil
}

func updateBalanceOfSupplier(accountReq *pb.AccountRequest) (uint64, error) {
	// customer balance /// ---------------
	rowsAffected, storeError := model.UpdateSupplierBalance(db, accountReq)
	if storeError != nil {
		return 0, storeError
	}
	fmt.Printf("rowsAffected %v\n", &rowsAffected)
	return rowsAffected, nil
}

var db *sqlx.DB

func main() {

	var databaseError error
	db, databaseError = model.NewDB("datasource")
	if databaseError != nil {
		log.Fatalf("failed to listen: %v", databaseError)
	}

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