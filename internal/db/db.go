package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

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
	appPath, err := os.Executable()
	if err != nil {
		return false, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), databaseFile)
	_, err = os.Stat(dbFile)

	exists := false
	if err == nil {
		exists = true
	}
	return exists, nil
}
