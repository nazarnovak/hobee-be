package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/nazarnovak/hobee-be/config"
)

var Instance *sql.DB

func Init(cnfDB config.DB, isDev bool) error {
	connectionString := cnfDB.Connection
	if !isDev {
		connectionString = cnfDB.ConnectionProd
	}

	db, err := sql.Open("postgres", connectionString)
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

