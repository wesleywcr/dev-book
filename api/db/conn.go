package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wesleywcr/dev-book/api/config"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.StringConectionDB)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
