package domain

import "errors"

// Task represents a single todo item.
type Task struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
}

// Sentinel errors — callers compare with errors.Is().
var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrEmptyTaskText  = errors.New("task text cannot be empty")
	ErrInvalidReorder = errors.New("reorder list must contain all task IDs exactly once")
)

// TaskRepository defines the persistence contract (Dependency-Inversion).
type TaskRepository interface {
	FindAll() ([]*Task, error)
	FindByID(id string) (*Task, error)
	Save(task *Task) error
	Delete(id string) error
	Reorder(ids []string) error
}
