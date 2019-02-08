package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var Instance *sql.DB

func Init(connection string) error {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return err
	}
	//defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	Instance = db

	return nil
}

