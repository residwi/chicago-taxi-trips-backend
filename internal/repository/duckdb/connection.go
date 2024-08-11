package duckdb

import (
	"database/sql"

	_ "github.com/marcboeker/go-duckdb"
	log "github.com/sirupsen/logrus"
)

func CreateConnection() *sql.DB {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
