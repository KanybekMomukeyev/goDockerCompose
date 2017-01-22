package model

import (
	"log"
	"github.com/jmoiron/sqlx"
	"fmt"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveProduct = `
DROP TABLE IF EXISTS products;
`

var schemaCreateProduct = `
CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY NOT NULL,
    product_image_path varchar (400),
    product_name varchar (400),
    supplier_id INTEGER,
    category_id INTEGER,
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


func UpdateProduct(db *sqlx.DB, product *pb.ProductRequest) (uint64, error)  {

	tx := db.MustBegin()

	stmt, err :=tx.Prepare("UPDATE products SET product_image_path=$1, product_name=$2, supplier_id=$3, " +
		"category_id=$4, barcode=$5, quantity_per_unit=$6, sale_unit_price=$7, " +
		"income_unit_price=$8, units_in_stock=$9 WHERE product_id=$10")
	CheckErr(err)

	res, err2 := stmt.Exec(product.ProductImagePath,
		product.ProductName,
		product.SupplierId,
		product.CategoryId,
		product.Barcode,
		product.QuantityPerUnit,
		product.SaleUnitPrice,
		product.IncomeUnitPrice,
		product.UnitsInStock,
		product.ProductId)
	CheckErr(err2)

	affect, err := res.RowsAffected()
	CheckErr(err)

	fmt.Println(affect, "rows changed")

	commitError := tx.Commit()
	CheckErr(commitError)

	return uint64(affect), nil
}

func StoreProduct(db *sqlx.DB, product *pb.ProductRequest) (uint64, error)  {

	tx := db.MustBegin()
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

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
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