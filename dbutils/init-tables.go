package dbutils

import (
	"database/sql"
	"log"
)

func Initialize(dbDriver *sql.DB) {
	statement, driverError := dbDriver.Prepare(video)
	if driverError != nil {
		log.Println(driverError)
	}
	// Create video table
	_, statementError := statement.Exec()
	if statementError != nil {
		log.Println("Table already exists!")
	}
	statement, _ = dbDriver.Prepare(assets)
	statement.Exec()
	statement, _ = dbDriver.Prepare(schedule)
	statement.Exec()
	log.Println("All tables created/initialized successfully!")
}
