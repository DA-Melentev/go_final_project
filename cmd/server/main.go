package main

import (
	"context"
	"errors"
	"github.com/DA-Melentev/go_final_project/config"
	"github.com/DA-Melentev/go_final_project/db"
	"github.com/DA-Melentev/go_final_project/internal/handlers"
	"github.com/DA-Melentev/go_final_project/internal/repositories"
	"github.com/DA-Melentev/go_final_project/internal/services"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	taskHandler *handlers.TaskHandler
)

func main() {
	if err := initDb(); err != nil {
		log.Fatal("Failed to load database: ", err)
		return
	}
	defer closeDatabase()

	mapDependencies()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := startListening()

	<-stop
	log.Println("Shutdown signal received, stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
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
	log.Println("Connected to database")

	if !dbIsInstall {
		log.Println("Database is not initialized. Run migration...")
		if err := db.RunMigration(); err != nil {
			return err
		}
	}

	return nil
}

func startListening() *http.Server {
	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		port = config.Port
	}

	r := getRouter()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server is starting on %s", "localhost:"+port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	return srv
}

func mapDependencies() {
	taskRepo := repositories.NewTaskRepository(db.DB)
	taskService := services.NewTaskService(taskRepo)
	taskHandler = handlers.NewTaskHandler(taskService)
}

func getRouter() *chi.Mux {
	r := chi.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	r.Get("/api/nextdate", taskHandler.NextDate)
	r.Post("/api/signin", handlers.SignIn)

	r.Route("/api/task", func(r chi.Router) {
		r.Use(handlers.Auth)
		r.Get("/", taskHandler.GetTask)
		r.Post("/", taskHandler.AddTask)
		r.Put("/", taskHandler.PutTask)
		r.Delete("/", taskHandler.TaskDelete)

		r.Post("/done", taskHandler.TaskDone)
	})

	r.Route("/api/tasks", func(r chi.Router) {
		r.Use(handlers.Auth)
		r.Get("/", taskHandler.GetTasks)
	})

	return r
}

func closeDatabase() {
	if err := db.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed")
	}
}
