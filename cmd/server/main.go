package main

import (
	"fmt"
	"github.com/DA-Melentev/go_final_project/cmd/internal/db"
	"github.com/DA-Melentev/go_final_project/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := initDb(); err != nil {
		log.Fatal("Unable to load database: ", err)
		return
	}

	if err := startListening(); err != nil {
		log.Fatal("Unable to load server: ", err)
		return
	}
}

func initDb() error {
	dbPath := os.Getenv("TODO_DBFILE")
	if len(dbPath) == 0 {
		dbPath = config.DbFile
	}
	dbIsInstall, err := db.IsDbFileExists(dbPath)
	if err != nil {
		return err
	}

	if err := db.Connect(dbPath); err != nil {
		return err
	}

	if !dbIsInstall {
		log.Println("Database is not initialized. Run migration...")
		if err := db.RunMigration(); err != nil {
			return err
		}
	}

	return nil
}

func startListening() error {
	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		port = config.Port
	}

	r := chi.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	log.Printf("Start listening on %s", "localhost:"+config.Port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Println("Unable to load server: ", err)
		return err
	}
	return nil
}
