package repository

import (
	"sync"

	"todo/internal/domain"
)

// MemoryTaskRepository is a thread-safe, in-memory implementation of
// domain.TaskRepository. Swap this for a DB-backed version without changing
// any other layer (Open/Closed, Dependency-Inversion).
type MemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
	order []string // preserves insertion / drag order
}

// NewMemoryTaskRepository returns an empty repository.
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{tasks: make(map[string]*domain.Task)}
}

// FindAll returns tasks in their current order. Each returned value is a
// defensive copy so callers cannot mutate repository state.
func (r *MemoryTaskRepository) FindAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Task, 0, len(r.order))
	for _, id := range r.order {
		if t, ok := r.tasks[id]; ok {
			cp := *t
			result = append(result, &cp)
		}
	}
	return result, nil
}

// FindByID returns a defensive copy of the task with the given ID.
func (r *MemoryTaskRepository) FindByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	cp := *t
	return &cp, nil
}

// Save inserts or updates a task. New tasks are appended to the end.
func (r *MemoryTaskRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		task.Order = len(r.order)
		r.order = append(r.order, task.ID)
	}
	cp := *task
	r.tasks[task.ID] = &cp
	return nil
}

// Delete removes a task; returns ErrTaskNotFound when the ID is unknown.
func (r *MemoryTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(r.tasks, id)
	for i, oid := range r.order {
		if oid == id {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	return nil
}

// Reorder replaces the task ordering with the supplied ID slice.
// The slice must contain every existing task ID exactly once.
func (r *MemoryTaskRepository) Reorder(ids []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(ids) != len(r.order) {
		return domain.ErrInvalidReorder
	}
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := r.tasks[id]; !ok {
			return domain.ErrTaskNotFound
		}
		if _, dup := seen[id]; dup {
			return domain.ErrInvalidReorder
		}
		seen[id] = struct{}{}
	}
	r.order = make([]string, len(ids))
	copy(r.order, ids)
	for i, id := range ids {
		r.tasks[id].Order = i
	}
	return nil
}
