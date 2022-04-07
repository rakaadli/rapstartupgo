package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

//not used for a reason

func conn() (*sql.Conn, error) {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil

}
