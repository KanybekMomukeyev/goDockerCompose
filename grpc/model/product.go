package model

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
