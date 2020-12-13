package db

import (
	"database/sql"
	"github.com/GoogleCloudPlatform/cloudsql-proxy/logging"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"

	"github.com/nazarnovak/hobee-be/config"
)

var Instance *sql.DB

func Init(cnfDB config.DB, isDev bool) error {
	connectionString := cnfDB.Connection
	driver := "postgres"

	if !isDev {
		connectionString = cnfDB.ConnectionProd
		driver = "cloudsqlpostgres"
	}

	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return err
	}
	//defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	Instance = db

	// Supress the "ephemeral certificate for instance hobeechat:europe-west6:myinstance3 will expire soon, refreshing now."
	logging.LogVerboseToNowhere()

	return nil
}

