package model

import (
	"github.com/jmoiron/sqlx"
	log "github.com/Sirupsen/logrus"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveCategory = `
DROP TABLE IF EXISTS categories;
`

var schemaCreateCategory = `
CREATE TABLE IF NOT EXISTS categories (
    category_id BIGSERIAL PRIMARY KEY NOT NULL,
    category_name varchar (400),
    category_updated_at BIGINT
);
`

type Category struct {
	categoryId uint64 `db:"category_id"`
	categoryName string `db:"category_name"`
}

func DeleteCategoryIfNotExsists(db *sqlx.DB) {
	db.MustExec(schemaRemoveCategory)
}

func CreateCategoryIfNotExsists(db *sqlx.DB) {
	db.MustExec(schemaCreateCategory)
	db.MustExec("ALTER TABLE categories ADD COLUMN IF NOT EXISTS category_updated_at BIGINT DEFAULT 0")
}

func StoreCategory(tx *sqlx.Tx, categoryRequest *pb.CategoryRequest) (uint64, error)  {

	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO categories " +
		"(category_name, category_updated_at) " +
		"VALUES($1, $2) returning category_id;",
		categoryRequest.CategoryName, categoryRequest.CategoryUpdatedAt).Scan(&lastInsertId)

	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"last inserted category_id":  lastInsertId,
	}).Info("")
	return lastInsertId, nil
}

func UpdateCategory(tx *sqlx.Tx, categoryRequest *pb.CategoryRequest) (uint64, error)  {

	stmt, err := tx.Prepare("UPDATE categories SET category_name=$1, category_updated_at=$2 WHERE category_id=$3")
	if err != nil {
		return ErrorFunc(err)
	}

	res, err := stmt.Exec(categoryRequest.CategoryName, categoryRequest.CategoryUpdatedAt, categoryRequest.CategoryId)
	if err != nil {
		return ErrorFunc(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"rows changed UpdateCategory": affect,
	}).Info("")
	return uint64(affect), nil
}

func AllCategory(db *sqlx.DB) ([]*pb.CategoryRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT category_id, category_name, category_updated_at FROM categories ORDER BY category_name ASC")
	if err != nil {
		print("error")
	}

	categories := make([]*pb.CategoryRequest, 0)
	for rows.Next() {
		category := new(pb.CategoryRequest)
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.CategoryUpdatedAt)
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

func AllUpdatedCategories(db *sqlx.DB, categoryFilter *pb.CategoryFilter) ([]*pb.CategoryRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT category_id, category_name, category_updated_at FROM categories WHERE category_updated_at >= $1 LIMIT $2", categoryFilter.CategoryUpdatedAt, 1000)
	if err != nil {
		print("error")
	}

	categories := make([]*pb.CategoryRequest, 0)
	for rows.Next() {
		category := new(pb.CategoryRequest)
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.CategoryUpdatedAt)
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