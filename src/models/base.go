package dbmodel

import (
	"database/sql"
	"log"
)

var sqlDb *sql.DB
var err error

func GetDb() *sql.DB {
	sqlStr := "root:dev@tcp(127.0.0.1:3306)/go-test"
	sqlDb, err := sql.Open("mysql",sqlStr )
	if err!=nil{
		log.Println("Failed to connect to DB.", err)
	}
	return sqlDb
}
