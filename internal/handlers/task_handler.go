package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DA-Melentev/go_final_project/internal/models"
	"github.com/DA-Melentev/go_final_project/internal/services"
	"github.com/DA-Melentev/go_final_project/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
}

func (h *TaskHandler) NextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	log.Printf("New nextdate request {now:%s date:%s repeat:%s}", nowStr, date, repeat)

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid now param", http.StatusBadRequest)
		log.Printf("Invalid now param: %v", err)
		return
	}

	_, err = time.Parse("20060102", date)
	if err != nil {
		http.Error(w, "Invalid date param", http.StatusBadRequest)
		log.Printf("Invalid date param: %v", err)
		return
	}

	result, err := utils.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("error while calculating nextdate: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		log.Printf("Error while parsing json: %v", err)
		return
	}
	log.Printf("/api/task POST: %v", task)

	status, err := handleTask(&task)
	if err != nil {
		WriteError(w, status, err)
		log.Printf("error: %v", err)
		return
	}

	id, err := h.Service.AddTask(task)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		log.Printf("error: %v", err)
		return
	}

	WriteResponseJSON(w, http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *TaskHandler) PutTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		log.Printf("Error while parsing json: %v", err)
		return
	}
	log.Printf("/api/task PUT: %v", task)

	status, err := handleTask(&task)
	if err != nil {
		WriteError(w, status, err)
		log.Printf("error: %v", err)
		return
	}

	err = h.Service.PutTask(task)
	if err != nil {
		var status int
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		WriteError(w, status, err)
		log.Printf("error: %v", err)
		return
	}

	WriteResponseJSON(w, http.StatusCreated, map[string]interface{}{})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractId(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		log.Printf("error: %v", err)
		return
	}

	log.Printf("/api/task GET ?id=%d", id)

	result, err := h.Service.GetTaskById(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	WriteResponseJSON(w, http.StatusOK, result)
}

func (h *TaskHandler) TaskDone(w http.ResponseWriter, r *http.Request) {
	id, err := extractId(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		log.Printf("error: %v", err)
		return
	}

	log.Printf("/api/task/done POST ?id=%d", id)

	err = h.Service.TaskDone(id)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		log.Printf("error: %v", err)
		return
	}

	WriteResponseJSON(w, http.StatusAccepted, map[string]interface{}{})
}

func (h *TaskHandler) TaskDelete(w http.ResponseWriter, r *http.Request) {
	id, err := extractId(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		log.Printf("error: %v", err)
		return
	}

	log.Printf("/api/task/done DELETE ?id=%d", id)

	err = h.Service.TaskDelete(id)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		log.Printf("error: %v", err)
		return
	}

	WriteResponseJSON(w, http.StatusOK, map[string]interface{}{})
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if len(search) > 0 {
		log.Printf("/api/tasks GET ?search=%s", search)
	} else {
		log.Println("/api/tasks GET")
	}

	var (
		tasks []models.Task
		err   error
	)

	if len(search) == 0 {
		tasks, err = h.Service.GetAllTasks()
	} else {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			tasks, err = h.Service.GetTasksBySearch(search)
		} else {
			tasks, err = h.Service.GetTasksByDate(date)
		}
	}

	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		log.Printf("error: %v", err)
		return
	}

	WriteResponseJSON(w, http.StatusOK, map[string]interface{}{
		"tasks": tasks,
	})
}

func handleTask(task *models.Task) (int, error) {
	if len(task.Title) == 0 {
		err := errors.New("title field is required")
		return http.StatusBadRequest, err
	}

	var (
		date time.Time
		err  error
	)
	if len(task.Date) == 0 {
		date = time.Now().Truncate(24 * time.Hour).UTC()
		task.Date = date.Format("20060102")
	} else {
		date, err = time.Parse("20060102", task.Date)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("wrong date format: %v", err)
		}
	}

	var nextDate string
	if len(task.Repeat) > 0 {
		nextDate, err = utils.NextDate(time.Now().Truncate(24*time.Hour).UTC(), task.Date, task.Repeat)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	now := time.Now().Truncate(24 * time.Hour).UTC()
	date = date.Truncate(24 * time.Hour)

	if date.Before(now) {
		if len(task.Repeat) > 0 {
			task.Date = nextDate
		} else {
			task.Date = time.Now().Truncate(24 * time.Hour).UTC().Format("20060102")
		}
	}

	return 0, nil
}

func extractId(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	if len(idStr) == 0 {
		return 0, errors.New("id param not specified")
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("wrong id param")
	}
	return int(id), nil
}
