package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/internal/handler"
	"todo/internal/repository"
	"todo/internal/service"
)

// newMux wires a fresh in-memory stack and returns a ready-to-use ServeMux.
func newMux() *http.ServeMux {
	svc := service.NewTaskService(repository.NewMemoryTaskRepository())
	h := handler.NewTaskHandler(svc)
	mux := http.NewServeMux()
	h.Register(mux)
	return mux
}

func do(mux *http.ServeMux, method, path, body string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Buffer
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	} else {
		bodyReader = &bytes.Buffer{}
	}
	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func decodeTask(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var task map[string]any
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return task
}

// ─── List ─────────────────────────────────────────────────────────────────────

func TestHandler_ListTasks_ReturnsEmptySlice(t *testing.T) {
	w := do(newMux(), http.MethodGet, "/api/tasks", "")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var tasks []any
	_ = json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 0 {
		t.Errorf("expected empty list, got %d items", len(tasks))
	}
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestHandler_CreateTask_ReturnsCreatedTask(t *testing.T) {
	w := do(newMux(), http.MethodPost, "/api/tasks", `{"text":"Buy milk"}`)
	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	task := decodeTask(t, w)
	if task["text"] != "Buy milk" {
		t.Errorf("text = %v, want %q", task["text"], "Buy milk")
	}
	if task["id"] == "" {
		t.Error("id should not be empty")
	}
}

func TestHandler_CreateTask_EmptyText_ReturnsBadRequest(t *testing.T) {
	w := do(newMux(), http.MethodPost, "/api/tasks", `{"text":""}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_CreateTask_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	w := do(newMux(), http.MethodPost, "/api/tasks", `not json`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestHandler_UpdateTask_ToggleComplete(t *testing.T) {
	mux := newMux()
	created := decodeTask(t, do(mux, http.MethodPost, "/api/tasks", `{"text":"Task"}`))
	id := created["id"].(string)

	w := do(mux, http.MethodPut, "/api/tasks/"+id, `{"completed":true}`)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	updated := decodeTask(t, w)
	if !updated["completed"].(bool) {
		t.Error("expected completed = true")
	}
}

func TestHandler_UpdateTask_NotFound(t *testing.T) {
	w := do(newMux(), http.MethodPut, "/api/tasks/ghost", `{"completed":true}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestHandler_DeleteTask(t *testing.T) {
	mux := newMux()
	created := decodeTask(t, do(mux, http.MethodPost, "/api/tasks", `{"text":"Gone"}`))
	id := created["id"].(string)

	w := do(mux, http.MethodDelete, "/api/tasks/"+id, "")
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestHandler_DeleteTask_NotFound(t *testing.T) {
	w := do(newMux(), http.MethodDelete, "/api/tasks/ghost", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ─── Reorder ──────────────────────────────────────────────────────────────────

func TestHandler_ReorderTasks(t *testing.T) {
	mux := newMux()
	t1 := decodeTask(t, do(mux, http.MethodPost, "/api/tasks", `{"text":"A"}`))
	t2 := decodeTask(t, do(mux, http.MethodPost, "/api/tasks", `{"text":"B"}`))
	t3 := decodeTask(t, do(mux, http.MethodPost, "/api/tasks", `{"text":"C"}`))

	ids := fmt.Sprintf(`{"ids":[%q,%q,%q]}`, t3["id"], t1["id"], t2["id"])
	w := do(mux, http.MethodPut, "/api/tasks/reorder", ids)
	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	var tasks []map[string]any
	_ = json.NewDecoder(do(mux, http.MethodGet, "/api/tasks", "").Body).Decode(&tasks)
	if tasks[0]["id"] != t3["id"] {
		t.Errorf("first task id = %v, want %v", tasks[0]["id"], t3["id"])
	}
}
