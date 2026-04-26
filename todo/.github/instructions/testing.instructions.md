---
description: "Use when writing, running, or reviewing tests. Covers test structure, naming, table-driven tests, and what must be tested in each layer."
applyTo: "todo/internal/**/*_test.go"
---

# Testing Guidelines

## Test Structure per Layer

### Repository Tests (`repository/memory_test.go`)
- ทดสอบ CRUD operations ตามสัญญาของ `domain.TaskRepository`
- ทดสอบ thread safety ถ้าเกี่ยวข้อง
- ทดสอบ defensive copy (mutate ผลลัพธ์แล้ว repository state ต้องไม่เปลี่ยน)

### Service Tests (`service/task_service_test.go`)
- สร้าง service ผ่าน `MemoryTaskRepository` จริง (ไม่ใช่ mock — YAGNI)
- ทดสอบ business rule: validation, ID generation uniqueness, error propagation

### Handler Tests (`handler/task_handler_test.go`)
- ใช้ `httptest.NewRecorder()` และ `httptest.NewRequest()` เท่านั้น — ห้าม start real server
- ทดสอบ HTTP status code, response body, และ error cases
- ทดสอบ invalid JSON และ missing fields

## Naming Convention

```
Test{Type}_{Method}_{Scenario}
```

```go
// ✅
func TestMemory_FindByID_NotFound(t *testing.T) { ... }
func TestService_Create_EmptyText(t *testing.T) { ... }
func TestHandler_CreateTask_InvalidJSON_ReturnsBadRequest(t *testing.T) { ... }

// ❌
func TestCreate(t *testing.T) { ... }
func Test1(t *testing.T) { ... }
```

## Table-Driven Tests

ใช้สำหรับ logic ที่มีหลาย input/output combinations:

```go
func TestService_Create(t *testing.T) {
    tests := []struct {
        name    string
        text    string
        wantErr error
    }{
        {"valid text", "Buy milk", nil},
        {"empty text", "", domain.ErrEmptyTaskText},
        {"whitespace only", "   ", domain.ErrEmptyTaskText},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            _, err := newService().Create(tc.text)
            if err != tc.wantErr {
                t.Errorf("got %v, want %v", err, tc.wantErr)
            }
        })
    }
}
```

## Coverage Requirements

| Layer | Minimum Coverage |
|---|---|
| domain | N/A (pure data/interfaces) |
| repository | 90% |
| service | 90% |
| handler | 85% |

ตรวจ coverage:
```bash
go test ./todo/... -cover
```

## Error Assertion Pattern

ใช้ `errors.Is()` เสมอ — ห้าม compare string:

```go
// ✅
if err != domain.ErrTaskNotFound { ... }
if !errors.Is(err, domain.ErrTaskNotFound) { ... }

// ❌
if err.Error() == "task not found" { ... }
```

## Test Helper Pattern

```go
// helper ใน test file สร้าง wired instance ที่พร้อมใช้
func newService() *service.TaskService {
    return service.NewTaskService(repository.NewMemoryTaskRepository())
}
```

## What NOT to Test

- Go stdlib behavior (e.g., `json.Marshal` correctness)
- Third-party library internals
- `main()` function โดยตรง — ทดสอบผ่าน integration หรือ smoke test แทน
