package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	pb "github.com/KanybekMomukeyev/goDockerCompose/grpc/proto"
)

var schemaRemoveSupplier = `
DROP TABLE IF EXISTS suppliers;
`

var schemaCreateSupplier = `
CREATE TABLE IF NOT EXISTS suppliers (
    supplier_id BIGSERIAL PRIMARY KEY NOT NULL,
    supplier_image_path varchar (400),
    company_name varchar (400),
    contact_fname varchar (400),
    phone_number varchar (400),
    address varchar (400)
);
`

type Supplier struct {
	supplierId uint64 `db:"supplier_id"`
	supplierImagePath string `db:"supplier_image_path"`
	companyName string `db:"company_name"`
	contactFname string `db:"contact_fname"`
	phoneNumber string `db:"phone_number"`
	address float32 `db:"address"`
}

func CreateSupplierIfNotExsists(db *sqlx.DB) {
	//db.MustExec(schemaRemoveSupplier)
	db.MustExec(schemaCreateSupplier)
}

func StoreSupplier(tx *sqlx.Tx, supplierRequest *pb.SupplierRequest) (uint64, error)  {

	var lastInsertId uint64
	err := tx.QueryRow("INSERT INTO suppliers(supplier_image_path, company_name, contact_fname, phone_number, address) VALUES($1, $2, $3, $4, $5) returning supplier_id;",
		supplierRequest.SupplierImagePath,
		supplierRequest.CompanyName,
		supplierRequest.ContactFname,
		supplierRequest.PhoneNumber,
		supplierRequest.Address).Scan(&lastInsertId)

	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"last inserted supplier_id": lastInsertId,
	}).Info("")
	return lastInsertId, nil
}

func UpdateSupplier(tx *sqlx.Tx, supplierRequest *pb.SupplierRequest) (uint64, error)  {

	stmt, err := tx.Prepare("UPDATE suppliers SET supplier_image_path=$1, company_name=$2, contact_fname=$3, " +
		"phone_number=$4, address=$5 WHERE supplier_id=$6")
	if err != nil {
		return ErrorFunc(err)
	}

	res, err := stmt.Exec(supplierRequest.SupplierImagePath,
		supplierRequest.CompanyName,
		supplierRequest.ContactFname,
		supplierRequest.PhoneNumber,
		supplierRequest.Address,
		supplierRequest.SupplierId)
	if err != nil {
		return ErrorFunc(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrorFunc(err)
	}

	log.WithFields(log.Fields{
		"update supplier rows changed":  affect,
	}).Info("")
	return uint64(affect), nil
}

func AllSuppliers(db *sqlx.DB) ([]*pb.SupplierRequest, error) {

	pingError := db.Ping()

	if pingError != nil {
		log.Fatalln(pingError)
		return nil, pingError
	}

	rows, err := db.Queryx("SELECT supplier_id, supplier_image_path, company_name, contact_fname, phone_number, address FROM suppliers ORDER BY supplier_id ASC")
	if err != nil {
		print("error")
	}

	suppliers := make([]*pb.SupplierRequest, 0)
	for rows.Next() {
		supplier := new(pb.SupplierRequest)
		err := rows.Scan(&supplier.SupplierId, &supplier.SupplierImagePath, &supplier.CompanyName, &supplier.ContactFname, &supplier.PhoneNumber, &supplier.Address)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil
}
