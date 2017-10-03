package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func openMySQLConnection(user string, pass string, protocol string,
	address string, dbname string) (*sql.DB, error) {

	var source_name string

	source_name = user + ":" + pass + "@" + protocol + "(" + address + ")/" + dbname

	db, err := sql.Open("mysql", source_name)

	return db, err
}
