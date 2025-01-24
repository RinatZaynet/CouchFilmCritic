package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ClientDB struct {
	DB *sql.DB
}

func OpenDB() (db *sql.DB, err error) {
	// Настроить конфигурацию
	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/coursera?&charset=utf8&interpolateParams=true")
	if err != nil {
		return nil, fmt.Errorf("an error occurred while open database in NewClientDB. Error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while ping database in NewClientDB. Error: %w", err)
	}
	return db, nil
}
