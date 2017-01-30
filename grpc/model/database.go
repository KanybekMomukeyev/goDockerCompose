package model

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"log"
	"fmt"
	"os"
)


func NewDB(dataSourceName string) (*sqlx.DB, error) {

	connInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_ENV_POSTGRES_USER"),
		os.Getenv("DB_ENV_POSTGRES_DATABASENAME"),
		os.Getenv("DB_ENV_POSTGRES_PASSWORD"),
		os.Getenv("GODOCKERCOMPOSE_POSTGRES_1_PORT_5432_TCP_ADDR"),
		os.Getenv("GODOCKERCOMPOSE_POSTGRES_1_PORT_5432_TCP_PORT"),
	)

	fmt.Println(connInfo)
	//db, err := sqlx.Connect("postgres", connInfo)
	//db, err := sqlx.Connect("postgres", "user=kanybek dbname=databasename password=nazgulum host=172.17.0.4 port=5432 sslmode=disable")
	db, err := sqlx.Connect("postgres", "user=kanybek dbname=databasename password=nazgulum host=localhost port=5432 sslmode=disable")
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

