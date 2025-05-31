package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var myDB *sql.DB

func CreateConnection() (*sql.DB, error) {
	dsn := "root:mypass@tcp(localhost:3306)/kontak"
	myDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return myDB, nil
}

// singleton pattern
func InitDB() (*sql.DB, error) {
	if myDB == nil {
		return CreateConnection()
	}
	return myDB, nil
}

// Database Pooling (5)
//
