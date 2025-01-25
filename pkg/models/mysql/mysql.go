package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ManagerDB struct {
	Database *sql.DB
}

func OpenDB() (db *sql.DB, err error) {
	// Настроить конфигурацию
	db, err = sql.Open("mysql", `root@tcp(localhost:3306)/couch_film_critic_db?&charset=utf8&interpolateParams=true&parseTime=true`)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while open database in OpenDB(). Error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while ping database in OpenDB(). Error: %w", err)
	}

	return db, nil
}

func (db *ManagerDB) CloseDB() {
	db.Database.Close()
}
