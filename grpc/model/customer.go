package model

var schemaRemoveCustomer = `
DROP TABLE IF EXISTS customers;
`

var schemaCreateCustomer = `
CREATE TABLE IF NOT EXISTS customers (
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

var schemaCreateIndexForCustomer1 = `CREATE INDEX IF NOT EXISTS supplier_id_products_idx ON products (supplier_id)`
var schemaCreateIndexForCustomer2 = `CREATE INDEX IF NOT EXISTS category_id_products_idx ON products (category_id)`
var schemaCreateIndexForCustomer3 = `CREATE INDEX IF NOT EXISTS barcode_products_idx ON products (barcode)`

type Customer struct {
	customerId uint64 `db:"product_id"`
	customerImagePath string `db:"product_image_path"`
	firstName string `db:"product_name"`
	secondName string `db:"barcode"`
	phoneNumber string `db:"quantity_per_unit"`
	address float32 `db:"sale_unit_price"`
}