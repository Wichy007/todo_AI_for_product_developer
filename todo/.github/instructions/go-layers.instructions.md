---
description: "Use when writing or reviewing Go source files in handler, service, repository, or domain layers. Covers Go idioms, error handling patterns, and concurrency rules."
applyTo: "todo/internal/**/*.go"
---

# Go Layer Guidelines

## Package Declarations

- `internal/domain` → `package domain`
- `internal/repository` → `package repository` / test: `package repository_test`
- `internal/service` → `package service` / test: `package service_test`
- `internal/handler` → `package handler` / test: `package handler_test`

ใช้ `_test` suffix สำหรับ test package เสมอ (black-box testing)

## Error Handling

### Sentinel Errors
```go
// นิยามใน domain/task.go
var ErrTaskNotFound = errors.New("task not found")

// wrap เมื่อเพิ่ม context
return fmt.Errorf("FindByID %s: %w", id, domain.ErrTaskNotFound)

// unwrap ด้วย errors.Is()
if errors.Is(err, domain.ErrTaskNotFound) { ... }
```

### ห้าม panic ใน business logic
- `panic` ใช้ได้เฉพาะ `init()` หรือ `main()` เมื่อ setup ล้มเหลว
- handler ต้อง return error response ไม่ crash server

## Struct & Constructor

```go
// ✅ constructor ทุกตัวต้องชื่อ New{Type}
func NewTaskService(repo domain.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

// ✅ unexported fields — บังคับให้ใช้ constructor
type TaskService struct {
    repo domain.TaskRepository
}
```

## Concurrency (Repository layer)

- `MemoryTaskRepository` ใช้ `sync.RWMutex`
- read operations: `RLock/RUnlock`
- write operations: `Lock/Unlock`
- ห้าม hold lock ขณะ call function ภายนอก

```go
func (r *MemoryTaskRepository) FindAll() ([]*domain.Task, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // ...
}
```

## HTTP Handler Pattern

```go
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
    // 1. decode request
    var body struct { Text string `json:"text"` }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    // 2. call service
    task, err := h.svc.Create(body.Text)
    if err != nil {
        writeError(w, err)   // centralized error mapping
        return
    }
    // 3. write response
    writeJSON(w, http.StatusCreated, task)
}
```

ทุก handler: decode → call service → writeJSON/writeError — ห้ามมี logic อื่น

## ID Generation

ใช้ `crypto/rand` เท่านั้น — ห้าม `math/rand` (insecure):

```go
func generateID() (string, error) {
    b := make([]byte, 8)
    if _, err := rand.Read(b); err != nil { return "", err }
    return hex.EncodeToString(b), nil
}
```

## Go Version

ใช้ Go ≥ 1.22 — ใช้ method+pattern routing ของ stdlib ได้:

```go
mux.HandleFunc("GET /api/tasks", h.listTasks)
mux.HandleFunc("PUT /api/tasks/{id}", h.updateTask)
// PathValue ใช้ได้ใน Go 1.22+
id := r.PathValue("id")
```
