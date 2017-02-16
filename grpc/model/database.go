package model

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"fmt"
	"os"
	log "github.com/Sirupsen/logrus"
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
	db, err := sqlx.Connect("postgres", connInfo) // for compose
	//db, err := sqlx.Connect("postgres", "user=kanybek dbname=databasename password=nazgulum host=172.17.0.4 port=5432 sslmode=disable") // for single docker app
	//db, err := sqlx.Connect("postgres", "user=kanybek dbname=databasename password=nazgulum host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.WithFields(log.Fields{
			"connection info": connInfo,
			"error": err,
		}).Fatal("Can not connected to database")
		return nil, err
	}

	pingError := db.Ping()

	if pingError != nil {
		log.WithFields(log.Fields{
			"info": connInfo,
			"ping": pingError,
		}).Fatal("Can not connected Ping to database")
		return nil, pingError
	}

	return db, nil
}

