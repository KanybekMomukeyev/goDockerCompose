package model

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	"log"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveCategory = `
DROP TABLE IF EXISTS categories;
`

var schemaCreateCategory = `
CREATE TABLE IF NOT EXISTS categories (
    category_id BIGSERIAL PRIMARY KEY NOT NULL,
    category_name varchar (400)
);
`

type Category struct {
	categoryId uint64 `db:"category_id"`
	categoryName string `db:"category_name"`
}

func CreateCategoryIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveCategory)
	db.MustExec(schemaCreateCategory)
}

func StoreCategory(db *sqlx.DB, categoryRequest *pb.CategoryRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO categories " +
		"(category_name) " +
		"VALUES($1) returning category_id;",
		categoryRequest.CategoryName).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
}

func AllCategory(db *sqlx.DB) ([]*pb.CategoryRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT category_id, category_name FROM categories ORDER BY category_name ASC")
	if err != nil {
		print("error")
	}

	categories := make([]*pb.CategoryRequest, 0)
	for rows.Next() {
		category := new(pb.CategoryRequest)
		err := rows.Scan(&category.CategoryId, &category.CategoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}