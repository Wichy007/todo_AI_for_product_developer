---
description: "Use when adding new layers, packages, files, or changing how layers communicate. Covers Clean Architecture rules, import boundaries, and layer responsibilities."
applyTo: "todo/internal/**/*.go"
---

# Architecture Guidelines

## Layer Map

```
┌─────────────────────────────────────────────┐
│  main.go  (composition root)                │
│  ─ wires all dependencies                   │
│  ─ starts HTTP server                       │
└──────────┬──────────────────────────────────┘
           │ imports
┌──────────▼──────────────────────────────────┐
│  handler/  (Delivery layer)                 │
│  ─ HTTP request/response                    │
│  ─ JSON encode/decode                       │
│  ─ maps errors → HTTP status codes          │
└──────────┬──────────────────────────────────┘
           │ imports
┌──────────▼──────────────────────────────────┐
│  service/  (Use-case layer)                 │
│  ─ business rules                           │
│  ─ validation                               │
│  ─ orchestrates repository calls            │
└──────────┬──────────────────────────────────┘
           │ imports
┌──────────▼──────────────────────────────────┐
│  domain/   (Enterprise layer)               │
│  ─ Task entity                              │
│  ─ TaskRepository interface                 │
│  ─ Sentinel errors                          │
│  ─ ZERO external imports                    │
└──────────▲──────────────────────────────────┘
           │ implements
┌──────────┴──────────────────────────────────┐
│  repository/  (Infrastructure layer)        │
│  ─ MemoryTaskRepository                     │
│  ─ (future: PostgresTaskRepository etc.)    │
└─────────────────────────────────────────────┘
```

## Import Rules (Enforced)

| Package | สามารถ import | ห้าม import |
|---|---|---|
| `domain` | stdlib เท่านั้น | ทุก internal package |
| `repository` | `domain`, stdlib | `service`, `handler` |
| `service` | `domain`, stdlib | `repository` (concrete), `handler` |
| `handler` | `service`, `domain`, stdlib | `repository` (concrete) |
| `main` | ทุก layer | — |

> **ทำไม:** ถ้า `service` import `repository` concrete type, เราไม่สามารถ swap storage ได้โดยไม่แตะ service

## Adding a New Feature — Checklist

1. **domain**: เพิ่ม field ใน `Task` หรือ method ใน `TaskRepository` interface (ถ้าจำเป็น)
2. **repository**: implement method ใหม่ใน `MemoryTaskRepository`
3. **service**: เพิ่ม method ใน `TaskService` พร้อม validation
4. **handler**: เพิ่ม route ใน `Register()` + handler method
5. **test**: เขียน test ทุก layer ก่อน implement (TDD)

## Sentinel Errors Pattern

error ทุกตัวต้องนิยามใน `domain/task.go` เป็น `var Err... = errors.New("...")`

```go
// ✅ ใน domain/task.go
var ErrTaskNotFound = errors.New("task not found")

// ✅ ใน handler — compare ด้วย errors.Is()
case errors.Is(err, domain.ErrTaskNotFound):
    http.Error(w, err.Error(), http.StatusNotFound)
```

ห้าม return `fmt.Errorf("task not found")` เพราะ caller ใช้ `errors.Is()` ไม่ได้

## Expanding Storage

เมื่อต้องการเปลี่ยน/เพิ่ม storage backend:

1. สร้าง `internal/repository/postgres.go` ที่ implement `domain.TaskRepository`
2. แก้ `main.go` เพื่อ wire implementation ใหม่
3. **ไม่ต้องแตะ** `service/` หรือ `handler/` เลย

## Static Files

- อยู่ใน `web/static/` — embedded ผ่าน `//go:embed web/static`
- ทุก path ที่ไม่ตรงกับ `/api/` จะถูก serve เป็น static file
- ห้ามใส่ secret หรือ config ใน static files
