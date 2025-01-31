package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ManagerDB struct {
	Database *sql.DB
}

func OpenDB(dsn string) (*sql.DB, error) {
	const fn = "storage.mysql.OpenDB()"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}
	return db, nil
}

func (db *ManagerDB) CloseDB() {
	db.Database.Close()
}
