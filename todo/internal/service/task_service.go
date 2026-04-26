package service

import (
	"crypto/rand"
	"encoding/hex"

	"hello/todo/internal/domain"
)

// TaskService contains all business logic for managing tasks.
// It depends on the domain.TaskRepository interface, not on a concrete type
// (Dependency-Inversion Principle).
type TaskService struct {
	repo domain.TaskRepository
}

// NewTaskService wires the service with the given repository.
func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// GetAll returns every task in display order.
func (s *TaskService) GetAll() ([]*domain.Task, error) {
	return s.repo.FindAll()
}

// Create validates text, generates an ID, persists the task, and returns it.
func (s *TaskService) Create(text string) (*domain.Task, error) {
	if text == "" {
		return nil, domain.ErrEmptyTaskText
	}
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	task := &domain.Task{ID: id, Text: text}
	if err := s.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

// Update applies partial changes (text and/or completed flag) to an existing task.
// Passing nil for either pointer leaves that field unchanged.
func (s *TaskService) Update(id string, text *string, completed *bool) (*domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if text != nil {
		if *text == "" {
			return nil, domain.ErrEmptyTaskText
		}
		task.Text = *text
	}
	if completed != nil {
		task.Completed = *completed
	}
	if err := s.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

// Delete removes a task by ID.
func (s *TaskService) Delete(id string) error {
	return s.repo.Delete(id)
}

// Reorder sets a new display order from an ordered list of IDs.
func (s *TaskService) Reorder(ids []string) error {
	return s.repo.Reorder(ids)
}

// generateID produces a cryptographically random 16-char hex string.
func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
