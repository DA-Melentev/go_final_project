package repositories

import (
	"database/sql"
	"fmt"
	"github.com/DA-Melentev/go_final_project/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(database *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: database,
	}
}

func (r *TaskRepository) AddTask(task models.Task) (int, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *TaskRepository) PutTask(task models.Task) error {
	query := "UPDATE scheduler SET date = ?,  title = ?, comment = ?, repeat = ? WHERE id = ?"

	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return fmt.Errorf("task with id=%d not found", task.Id)
	}

	return nil
}

func (r *TaskRepository) GetTaskById(id int) (models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"

	row := r.db.QueryRow(query, id)

	var task models.Task
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return models.Task{}, err
	}

	if row.Err() != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetAllTasks() ([]models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := rowSetToTaskList(rows)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetTasksBySearch(search string) ([]models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE " +
		"title LIKE :search OR comment LIKE :search ORDER BY date"

	rows, err := r.db.Query(query, sql.Named("search", "%"+search+"%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := rowSetToTaskList(rows)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetTasksByDate(date string) ([]models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE " +
		"date LIKE ?"

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := rowSetToTaskList(rows)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(id int) error {
	query := "DELETE FROM scheduler WHERE id = ?"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if c == 0 {
		return fmt.Errorf("task with id=%d not found", id)
	}

	return nil
}

func rowSetToTaskList(rows *sql.Rows) ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
