package db

import (
	"github.com/DA-Melentev/go_final_project/config"
	"log"
	"os"
)

func RunMigration() error {
	content, err := os.ReadFile(config.InitMigrationPath)
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(content))
	if err != nil {
		return err
	}
	log.Println("Migrations completed successfully")
	return nil
}
