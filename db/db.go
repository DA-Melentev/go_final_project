package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect(databaseFile string) error {
	var err error
	DB, err = sql.Open("sqlite", databaseFile)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}
	log.Println("Connected to database")
	return nil
}

func IsDbFileExists(databaseFile string) (bool, error) {
	_, err := os.Stat(databaseFile)

	exists := false
	if err == nil {
		exists = true
	}
	return exists, nil
}
