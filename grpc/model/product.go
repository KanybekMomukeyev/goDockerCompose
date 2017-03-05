package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveProduct = `
DROP TABLE IF EXISTS products;
`

var schemaCreateProduct = `
CREATE TABLE IF NOT EXISTS products (
    product_id BIGSERIAL PRIMARY KEY NOT NULL,
    product_image_path varchar (400),
    product_name varchar (400),
    supplier_id BIGINT,
    category_id BIGINT,
    barcode VARCHAR (300),
    quantity_per_unit VARCHAR (300),
    sale_unit_price REAL,
    income_unit_price REAL,
    units_in_stock REAL
);
`

var schemaCreateIndexForProduct1 = `CREATE INDEX IF NOT EXISTS supplier_id_products_idx ON products (supplier_id)`
var schemaCreateIndexForProduct2 = `CREATE INDEX IF NOT EXISTS category_id_products_idx ON products (category_id)`
var schemaCreateIndexForProduct3 = `CREATE INDEX IF NOT EXISTS barcode_products_idx ON products (barcode)`

type Product struct {
	productId uint64 `db:"product_id"`
	productImagePath string `db:"product_image_path"`
	productName string `db:"product_name"`
	supplierId uint64 `db:"supplier_id"`
	categoryId uint64 `db:"category_id"`
	barcode string `db:"barcode"`
	quantityPerUnit string `db:"quantity_per_unit"`
	saleUnitPrice float32 `db:"sale_unit_price"`
	incomeUnitPrice float32 `db:"income_unit_price"`
	unitsInStock float32 `db:"units_in_stock"`
}

func CreateProductIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveProduct)
	db.MustExec(schemaCreateProduct)
	db.MustExec(schemaCreateIndexForProduct1)
	db.MustExec(schemaCreateIndexForProduct2)
	db.MustExec(schemaCreateIndexForProduct3)
}


func UpdateProduct(tx *sqlx.Tx, product *pb.ProductRequest) (uint64, error)  {

	stmt, err :=tx.Prepare("UPDATE products SET product_image_path=$1, product_name=$2, supplier_id=$3, " +
		"category_id=$4, barcode=$5, quantity_per_unit=$6, sale_unit_price=$7, " +
		"income_unit_price=$8, units_in_stock=$9 WHERE product_id=$10")
	if err != nil {
		return ErrorFunc(err)
	}

	res, err := stmt.Exec(product.ProductImagePath,
		product.ProductName,
		product.SupplierId,
		product.CategoryId,
		product.Barcode,
		product.QuantityPerUnit,
		product.SaleUnitPrice,
		product.IncomeUnitPrice,
		product.UnitsInStock,
		product.ProductId)
	if err != nil {
		return ErrorFunc(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{"update product rows changed":  affect, }).Info("")
	return uint64(affect), nil
}

func StoreProduct(tx *sqlx.Tx, product *pb.ProductRequest) (uint64, error) {

	var lastInsertId uint64
	err := tx.QueryRow("INSERT INTO products " +
		"(product_image_path, product_name, supplier_id, category_id, barcode," +
		" quantity_per_unit, sale_unit_price, income_unit_price, units_in_stock) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) returning product_id;",
		product.ProductImagePath,
		product.ProductName,
		product.SupplierId,
		product.CategoryId,
		product.Barcode,
		product.QuantityPerUnit,
		product.SaleUnitPrice,
		product.IncomeUnitPrice,
		product.UnitsInStock).Scan(&lastInsertId)

	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"last inserted product_id":  lastInsertId,
	}).Info("")
	return lastInsertId, nil
}

func IncreaseProductsInStock(db *sqlx.DB, orderDetailReqs []*pb.OrderDetailRequest)  (uint64, error)  {

	productIds := make([]uint64, 0)
	updateValues := make(map[uint64]float64)

	for _, orderDetailReq := range orderDetailReqs {
		productIds = append(productIds, orderDetailReq.ProductId)
		updateValues[orderDetailReq.ProductId] = orderDetailReq.OrderQuantity
	}

	products_, err := getProductsForProductIds(db, productIds)
	if err != nil {
		print("error")
		return 0, err
	}

	forUpdatesInDatabase := make(map[uint64]float64)

	for _, product := range products_ {
		quantity := updateValues[product.ProductId]
		product.UnitsInStock = product.UnitsInStock + quantity
		forUpdatesInDatabase[product.ProductId] = product.UnitsInStock
	}

	tx := db.MustBegin()
	rowsAffected, err := updateProductsInStock(tx, forUpdatesInDatabase)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("")
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return ErrorFunc(err)
	}

	return rowsAffected, nil
}

func DecreaseProductsInStock(db *sqlx.DB, orderDetailReqs []*pb.OrderDetailRequest) (uint64, error) {

	productIds := make([]uint64, 0)
	updateValues := make(map[uint64]float64)

	for _, orderDetailReq := range orderDetailReqs {
		productIds = append(productIds, orderDetailReq.ProductId)
		updateValues[orderDetailReq.ProductId] = orderDetailReq.OrderQuantity
	}

	products_, err := getProductsForProductIds(db, productIds)
	if err != nil {
		print("error")
		return 0, err
	}

	forUpdatesInDatabase := make(map[uint64]float64)

	for _, product := range products_ {
		quantity := updateValues[product.ProductId]
		product.UnitsInStock = product.UnitsInStock - quantity
		forUpdatesInDatabase[product.ProductId] = product.UnitsInStock
	}

	tx := db.MustBegin()
	rowsAffected, err := updateProductsInStock(tx, forUpdatesInDatabase)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{"err": err}).Warn("")
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return ErrorFunc(err)
	}

	return rowsAffected, nil
}

func updateProductsInStock(tx *sqlx.Tx, updateValues map[uint64]float64) (uint64, error) {

	var j int64 = 0
	for key, value := range updateValues {
		//fmt.Println("Key:", key, "Value:", value)

		stmt, err := tx.Prepare("UPDATE products SET units_in_stock=$1 WHERE product_id=$2")
		if err != nil {
			break
			return ErrorFunc(err)
		}

		res, err := stmt.Exec(value, key)
		if err != nil {
			break
			return ErrorFunc(err)
		}

		affect, err := res.RowsAffected()
		if err != nil {
			break
			return ErrorFunc(err)
		}

		j = j + affect
	}

	log.WithFields(log.Fields{
		"update products in stock": j,
	}).Info("")
	return uint64(j), nil
}

func getProductsForProductIds(db *sqlx.DB, productIds []uint64) ([]*pb.ProductRequest, error) {

	query, args, err := sqlx.In("SELECT product_id, product_image_path, product_name, supplier_id, " +
				"category_id, barcode, quantity_per_unit, sale_unit_price, " +
				"income_unit_price, units_in_stock FROM products WHERE product_id IN (?)", productIds)

	if err != nil {
		print("error")
		return nil, err
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query) //only if postgres
	rows, err := db.Query(query, args...)

	//var str string
	//for _, value := range productIds {
	//	str += strconv.Itoa(int(value)) + ","
	//}
	//
	//rows, err := db.Queryx("SELECT product_id, product_image_path, product_name, supplier_id, " +
	//	"category_id, barcode, quantity_per_unit, sale_unit_price, " +
	//	"income_unit_price, units_in_stock FROM products WHERE product_id IN (" + str[:len(str)-1] + ")")

	products := make([]*pb.ProductRequest, 0)
	for rows.Next() {
		product := new(pb.ProductRequest)
		err := rows.Scan(&product.ProductId, &product.ProductImagePath, &product.ProductName,
			&product.SupplierId, &product.CategoryId, &product.Barcode, &product.QuantityPerUnit,
			&product.SaleUnitPrice, &product.IncomeUnitPrice, &product.UnitsInStock)

		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func AllProducts(db *sqlx.DB) ([]*pb.ProductRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT product_id, product_image_path, product_name, supplier_id, " +
		"category_id, barcode, quantity_per_unit, sale_unit_price, " +
		"income_unit_price, units_in_stock FROM products ORDER BY product_id DESC")

	if err != nil {
		print("error")
	}

	//defer rows.Close()

	products := make([]*pb.ProductRequest, 0)
	for rows.Next() {
		product := new(pb.ProductRequest)
		err := rows.Scan(&product.ProductId, &product.ProductImagePath, &product.ProductName,
			&product.SupplierId, &product.CategoryId, &product.Barcode, &product.QuantityPerUnit,
			&product.SaleUnitPrice, &product.IncomeUnitPrice, &product.UnitsInStock)

		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}