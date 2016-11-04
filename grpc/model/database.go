package model

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"log"
)


func NewDB(dataSourceName string) (*sqlx.DB, error) {

	db, err := sqlx.Connect("postgres", "dbname=blog_test user=kanybek password=nazgulum host=localhost sslmode=disable")

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	return db, nil
}

