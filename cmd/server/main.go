package main

import (
	"github.com/DA-Melentev/go_final_project/config"
	"github.com/DA-Melentev/go_final_project/db"
	"github.com/DA-Melentev/go_final_project/internal/handlers"
	"github.com/DA-Melentev/go_final_project/internal/repositories"
	"github.com/DA-Melentev/go_final_project/internal/services"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

var (
	taskHandler *handlers.TaskHandler
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

	taskRepo := repositories.NewTaskRepository(db.DB)
	taskService := services.NewTaskService(taskRepo)
	taskHandler = handlers.NewTaskHandler(taskService)

	r := getRouter()

	log.Printf("Start listening on %s", "localhost:"+config.Port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		return err
	}
	return nil
}

func getRouter() *chi.Mux {
	r := chi.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", taskHandler.NextDate)

	r.Post("/api/task", taskHandler.AddTask)
	r.Get("/api/task", taskHandler.GetTask)
	r.Put("/api/task", taskHandler.PutTask)

	r.Get("/api/tasks", taskHandler.GetTasks)

	return r
}
