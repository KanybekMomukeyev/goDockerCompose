package model

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	"fmt"
	"log"
)

var schemaRemoveStaff = `
DROP TABLE IF EXISTS staff;
`

var schemaCreateStaff = `
CREATE TABLE IF NOT EXISTS staff (
    staff_id BIGSERIAL PRIMARY KEY NOT NULL,
    role_id INTEGER,
    staff_image_path VARCHAR (300),
    first_name VARCHAR (300),
    second_name VARCHAR (300),
    email VARCHAR (300) UNIQUE,
    password VARCHAR (300),
    phone_number VARCHAR (300),
    address VARCHAR (300)
);
`
//UNIQUE  set on email!!!!

var schemaCreateIndex = `CREATE INDEX IF NOT EXISTS role_id_idx ON staff (role_id)`

type Staff struct {
	staffId uint64 `db:"staff_id"`
	roleId uint64 `db:"role_id"`
	staffImagePath string `db:"staff_image_path"`
	firstName string `db:"first_name"`
	secondName string `db:"second_name"`
	email string `db:"email"`
	password string `db:"password"`
	phoneNumber string `db:"phone_number"`
	address string `db:"address"`
}

func CreateStaffIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveStaff)
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

	//CheckErr(err)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	commitError := tx.Commit()
	//CheckErr(commitError)
	if commitError != nil {
		fmt.Println(commitError)
		return 0, commitError
	}

	fmt.Println("last inserted staff_id =", lastInsertId)

	return lastInsertId, nil
}

func UpdateStaff(db *sqlx.DB, staff *pb.StaffRequest) (uint64, error)  {

	tx := db.MustBegin()

	stmt, err :=tx.Prepare("UPDATE staff SET role_id=$1, staff_image_path=$2, first_name=$3, second_name=$4, " +
		"email=$5, password=$6, phone_number=$7, address=$8 WHERE staff_id=$9")
	CheckErr(err)

	res, err2 := stmt.Exec(staff.RoleId,
		staff.StaffImagePath,
		staff.FirstName,
		staff.SecondName,
		staff.Email,
		staff.Password,
		staff.PhoneNumber,
		staff.Address,
		staff.StaffId)

	//CheckErr(err2)
	if err2 != nil {
		fmt.Println(err2)
		return 0, err2
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	fmt.Println(affect, "rows changed")

	commitError := tx.Commit()
	//CheckErr(commitError)
	if commitError != nil {
		fmt.Println(commitError)
		return 0, commitError
	}

	return uint64(affect), nil
}

func AllStaff(db *sqlx.DB) ([]*pb.StaffRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT staff_id, role_id, staff_image_path, first_name, second_name, email, password, phone_number, address FROM staff ORDER BY first_name ASC")
	if err != nil {
		print("error")
	}

	staff := make([]*pb.StaffRequest, 0)
	for rows.Next() {
		employee := new(pb.StaffRequest)
		err := rows.Scan(&employee.StaffId, &employee.RoleId, &employee.StaffImagePath, &employee.FirstName, &employee.SecondName, &employee.Email, &employee.Password, &employee.PhoneNumber, &employee.Address)
		if err != nil {
			return nil, err
		}
		staff = append(staff, employee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return staff, nil
}

func AllStaffAuto(db *sqlx.DB) ([]*pb.StaffRequest, error) {

	staff := []*Staff{}
	savedStaff := []*pb.StaffRequest{}

	err := db.Select(&staff, "SELECT staff_id, role_id, staff_image_path, first_name, second_name, email, password, phone_number, address FROM staff ORDER BY first_name ASC")
	if err != nil {
		print("error")
		panic(err)
	}

	for _, employee := range staff {

		staffRequest := &pb.StaffRequest {
			StaffId:    employee.staffId,
			RoleId:  employee.roleId,
			StaffImagePath: employee.staffImagePath,
			FirstName: employee.firstName,
			SecondName: employee.secondName,
			Email: employee.email,
			Password: employee.password,
			PhoneNumber: employee.phoneNumber,
			Address: employee.address,
		}

		savedStaff = append(savedStaff, staffRequest)
	}

	return savedStaff, nil
}
