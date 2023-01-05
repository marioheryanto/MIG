package database

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Init() *sql.DB {

	db, err := sql.Open("mysql", os.Getenv("DSN_CONNECTION_MYSQL"))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
