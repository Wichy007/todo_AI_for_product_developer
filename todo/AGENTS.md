# Todo Web App — Agent Instructions

> อ่านไฟล์นี้ก่อนทำงานใดๆ ใน project นี้เสมอ

## Project Overview

Web-based todo list application ที่เขียนด้วย **Go** (backend REST API) + **Vanilla JS/HTML/CSS** (frontend)  
ทำงานเป็น single binary ที่ embed static files ไว้ใน binary ผ่าน `//go:embed`

**Entry point:** `todo/main.go`  
**Module:** `hello`

---

## Repository Structure

```
todo/
├── main.go                          # HTTP server, wiring, 12-factor config
├── internal/
│   ├── domain/task.go               # Entity + Repository interface (contracts)
│   ├── repository/
│   │   ├── memory.go                # In-memory TaskRepository implementation
│   │   └── memory_test.go
│   ├── service/
│   │   ├── task_service.go          # Business logic
│   │   └── task_service_test.go
│   └── handler/
│       ├── task_handler.go          # HTTP handlers (REST)
│       └── task_handler_test.go
└── web/static/
    ├── index.html
    ├── style.css                    # Mobile-first responsive
    └── app.js                       # Drag-drop + REST API client
```

**ADR docs:** `.github/adr/`  
**Instruction files:** `.github/instructions/`

---

## Mandatory Engineering Practices

ดูรายละเอียดใน `.github/instructions/` — สรุปสั้นๆ:

| Practice | กฎเหล็ก |
|---|---|
| **TDD** | เขียน test ก่อน production code เสมอ — ห้าม push code ที่ไม่มี test |
| **Clean Code** | function ≤ 20 บรรทัด, ชื่อสื่อความหมาย, ไม่มี comment อธิบายสิ่งที่ code พูดแทนได้ |
| **SOLID** | interface อยู่ใน `domain/`, concrete อยู่ใน `repository/` — ห้าม import ย้อนกลับ |
| **12-Factor** | config จาก env var เท่านั้น (ไม่ hardcode), bind `0.0.0.0` |
| **DRY** | ห้าม copy-paste logic — extract เป็น function/helper |
| **YAGNI** | ห้าม implement feature ที่ยังไม่มี requirement |

---

## Build & Test Commands

```bash
# Run all tests
go test ./todo/... -v

# Build binary
go build -o /tmp/todo_web ./todo/

# Run server (default port 8080)
go run ./todo/

# Custom port
PORT=3000 go run ./todo/
```

---

## API Contract

| Method | Path | Description |
|---|---|---|
| GET | `/api/tasks` | List all tasks |
| POST | `/api/tasks` | Create task `{"text":"..."}` |
| PUT | `/api/tasks/{id}` | Update task `{"completed":true}` or `{"text":"..."}` |
| DELETE | `/api/tasks/{id}` | Delete task |
| PUT | `/api/tasks/reorder` | Reorder `{"ids":["a","b","c"]}` |

Static files served at `/` via embedded `web/static/`.

---

## Layer Dependency Rule (STRICT)

```
handler → service → domain ← repository
```

- **handler** ห้าม import `repository` โดยตรง
- **repository** ห้าม import `service` หรือ `handler`
- **domain** ไม่ import layer อื่น (pure Go, no side effects)

---

## Current Storage

In-memory เท่านั้น — ข้อมูลหายเมื่อ restart  
ดู [ADR-003](.github/adr/003-in-memory-storage.md) สำหรับ rationale และ migration path

---

## Key Decisions (ADR Index)

| # | Decision | Status |
|---|---|---|
| [001](.github/adr/001-go-stdlib-only-router.md) | ใช้ Go 1.22 stdlib router แทน third-party | Accepted |
| [002](.github/adr/002-embed-static-files.md) | Embed static files ใน binary | Accepted |
| [003](.github/adr/003-in-memory-storage.md) | In-memory storage สำหรับ MVP | Accepted |
| [004](.github/adr/004-vanilla-js-frontend.md) | Vanilla JS แทน framework | Accepted |
| [005](.github/adr/005-clean-architecture-layers.md) | Clean Architecture 4 layers | Accepted |
