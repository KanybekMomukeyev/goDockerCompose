package model

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
	log "github.com/Sirupsen/logrus"
	"errors"
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

var schemaCreateIndexForStaff1 = `CREATE INDEX IF NOT EXISTS role_id_staff_idx ON staff (role_id)`
var schemaCreateIndexForStaff2 = `CREATE INDEX IF NOT EXISTS email_staff_idx ON staff (email)`

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
	db.MustExec(schemaCreateIndexForStaff1)
	db.MustExec(schemaCreateIndexForStaff2)
}

func StoreStaff(tx *sqlx.Tx, staff *pb.StaffRequest) (uint64, error)  {

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

	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"last inserted staff_id": lastInsertId,
	}).Info("")
	return lastInsertId, nil
}

func UpdateStaff(tx *sqlx.Tx, staff *pb.StaffRequest) (uint64, error)  {

	stmt, err :=tx.Prepare("UPDATE staff SET role_id=$1, staff_image_path=$2, first_name=$3, second_name=$4, " +
		"email=$5, password=$6, phone_number=$7, address=$8 WHERE staff_id=$9")
	if err != nil {
		return ErrorFunc(err)
	}

	res, err := stmt.Exec(staff.RoleId,
		staff.StaffImagePath,
		staff.FirstName,
		staff.SecondName,
		staff.Email,
		staff.Password,
		staff.PhoneNumber,
		staff.Address,
		staff.StaffId)

	if err != nil {
		return ErrorFunc(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"update staff rows changed":  affect,
	}).Info("")
	return uint64(affect), nil
}

func SignIn(db *sqlx.DB, signInReq *pb.SignInRequest) (*pb.StaffRequest, error)  {

	rows, err := db.Queryx("SELECT staff_id, role_id, staff_image_path, first_name, second_name, email, password, phone_number, address FROM staff WHERE email=$1 AND password=$2 LIMIT 1", signInReq.Email, signInReq.Password)
	if err != nil {
		print("error")
	}
	staffFound := make([]*pb.StaffRequest, 0)
	for rows.Next() {
		employee := new(pb.StaffRequest)
		err := rows.Scan(&employee.StaffId, &employee.RoleId, &employee.StaffImagePath, &employee.FirstName, &employee.SecondName, &employee.Email, &employee.Password, &employee.PhoneNumber, &employee.Address)
		if err != nil {
			return nil, err
		}
		staffFound = append(staffFound, employee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(staffFound) > 0 {
		return staffFound[0], nil
	}


	return nil,errors.New("No such staff")
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
