package model

var schemaRemoveCategory = `
DROP TABLE IF EXISTS categories;
`

var schemaCreateCategory = `
CREATE TABLE IF NOT EXISTS categories (
    category_id SERIAL PRIMARY KEY NOT NULL,
    category_name varchar (400)
);
`

type Category struct {
	categoryId uint64 `db:"category_id"`
	categoryName string `db:"category_name"`
}