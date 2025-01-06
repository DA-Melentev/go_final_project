package handlers

import (
	"fmt"
	"github.com/DA-Melentev/go_final_project/utils"
	"log"
	"net/http"
	"time"

	"github.com/DA-Melentev/go_final_project/internal/services"
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

	log.Printf("New nextdate request now:%s date:%s repeat:%s", nowStr, date, repeat)

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid now param", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("20060102", date)
	if err != nil {
		http.Error(w, "Invalid date param", http.StatusBadRequest)
		return
	}

	result, err := utils.NextDate(now, date, repeat)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
	log.Printf("response sent")
}
