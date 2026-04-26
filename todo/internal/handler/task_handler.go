package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"hello/todo/internal/domain"
	"hello/todo/internal/service"
)

// TaskHandler maps HTTP requests to TaskService methods.
// Each route is registered once via Register(); no logic lives outside methods
// (Single Responsibility Principle).
type TaskHandler struct {
	svc *service.TaskService
}

// NewTaskHandler creates a TaskHandler backed by the given service.
func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// Register mounts all task routes on mux.
// Uses Go 1.22 method+pattern routing — no external router needed (YAGNI).
func (h *TaskHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tasks", h.listTasks)
	mux.HandleFunc("POST /api/tasks", h.createTask)
	mux.HandleFunc("PUT /api/tasks/reorder", h.reorderTasks) // must precede {id} pattern
	mux.HandleFunc("PUT /api/tasks/{id}", h.updateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", h.deleteTask)
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAll()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	task, err := h.svc.Create(body.Text)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Text      *string `json:"text"`
		Completed *bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	task, err := h.svc.Update(id, body.Text, body.Completed)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) reorderTasks(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.svc.Reorder(body.IDs); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrEmptyTaskText):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, domain.ErrInvalidReorder):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
