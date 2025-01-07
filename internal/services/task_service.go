package services

import (
	"github.com/DA-Melentev/go_final_project/internal/models"
	"github.com/DA-Melentev/go_final_project/internal/repositories"
	"time"
)

type TaskService struct {
	repo *repositories.TaskRepository
}

func NewTaskService(repository *repositories.TaskRepository) *TaskService {
	return &TaskService{
		repo: repository,
	}
}

func (s *TaskService) AddTask(task models.Task) (int, error) {
	return s.repo.AddTask(task)
}

func (s *TaskService) PutTask(task models.Task) error {
	return s.repo.PutTask(task)
}

func (s *TaskService) GetTaskById(id int) (models.Task, error) {
	return s.repo.GetTaskById(id)
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetTasksBySearch(search string) ([]models.Task, error) {
	return s.repo.GetTasksBySearch(search)
}

func (s *TaskService) GetTasksByDate(date time.Time) ([]models.Task, error) {
	dateString := date.Format("20060102")
	return s.repo.GetTasksByDate(dateString)
}
