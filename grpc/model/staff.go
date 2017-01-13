package model

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"fmt"
)

var schemaRemoveStaff = `
DROP TABLE IF EXISTS staff;
`

var schemaCreateStaff = `
CREATE TABLE IF NOT EXISTS staff (
    staff_id SERIAL PRIMARY KEY NOT NULL,
    role_id INTEGER,
    staff_image_path VARCHAR (300),
    first_name VARCHAR (300),
    second_name VARCHAR (300),
    email VARCHAR (300),
    password VARCHAR (300),
    phone_number VARCHAR (300),
    address VARCHAR (300)
);
`

var schemaCreateIndex = `CREATE INDEX IF NOT EXISTS role_id_idx ON staff (role_id)`

type Staff struct {
	staffId uint64 `db:"staff_id"`
	roleId string `db:"role_id"`
	staffImagePath string `db:"staff_image_path"`
	firstName string `db:"first_name"`
	secondName string `db:"second_name"`
	email string
	password string
	phoneNumber string `db:"phone_number"`
	address string
}

func CreateStaffIfNotExsists(db *sqlx.DB) {
	db.MustExec(schemaCreateStaff)
	db.MustExec(schemaCreateIndex)
}

func StoreStaff(db *sqlx.DB, staff *pb.StaffRequest) (uint64, error)  {

	tx := db.MustBegin()
	var lastInsertId uint64

	err := tx.QueryRow("INSERT INTO staff " +
		"(role_id, staff_image_path, first_name, second_name, email, password, phone_number, address) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8) returning staff_id;",
		staff.RoleId,
		staff.StaffImagePath,
		staff.FirstName,
		staff.SecondName,
		staff.Email,
		staff.Password,
		staff.PhoneNumber,
		staff.Address).Scan(&lastInsertId)

	CheckErr(err)

	commitError := tx.Commit()
	CheckErr(commitError)

	fmt.Println("last inserted id =", lastInsertId)

	return lastInsertId, nil
}

func AllStaffAuto(db *sqlx.DB) ([]*Staff, error) {

	staff := []*Staff{}
	db.Select(&staff, "SELECT staff_id, role_id, staff_image_path, first_name, second_name, email, password, phone_number, address FROM staff ORDER BY first_name ASC")
	return staff, nil
}
