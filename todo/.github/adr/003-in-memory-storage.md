# ADR-003: In-Memory Storage สำหรับ MVP

**Date:** 2026-04-26  
**Status:** Accepted

---

## Context

Application ต้องการ persistence layer  
ตัวเลือก: SQLite, PostgreSQL, In-memory, File-based (JSON)

## Decision

ใช้ **in-memory storage** (`MemoryTaskRepository`) สำหรับ MVP  
ข้อมูลหายเมื่อ restart server

## Rationale

| Criteria | In-Memory | SQLite | PostgreSQL |
|---|---|---|---|
| Setup | ทันที | ติดตั้ง driver | ติดตั้ง server |
| Dependencies | 0 | 1 (CGO) | 1 |
| Data persistence | ❌ | ✅ | ✅ |
| Complexity | ต่ำ | ปานกลาง | สูง |
| Testability | ✅ | ต้อง mock | ต้อง mock |
| 12-Factor (stateless) | ✅ | ❌ | ✅ (external) |

MVP เป้าหมายคือ validate UX และ architecture ไม่ใช่ production persistence

## Architecture Protection

`domain.TaskRepository` interface แยก persistence contract ออกจาก business logic  
การเปลี่ยน storage ต้องแตะเฉพาะ:
1. สร้าง `internal/repository/postgres.go` (implement interface)
2. แก้ `main.go` 1 บรรทัด (wire ใหม่)

**ไม่แตะ** `service/`, `handler/`, `domain/` เลย

## Consequences

- **ดี:** zero dependency, ง่ายต่อการ test
- **ดี:** architecture ถูก enforce ให้ swap-able ตั้งแต่ต้น
- **ข้อจำกัด:** ข้อมูลหายเมื่อ restart
- **ข้อจำกัด:** ไม่รองรับ multiple instances

## Migration Path เมื่อต้องการ Persistence

```
Priority 1: SQLite + modernc.org/sqlite (pure Go, no CGO)
Priority 2: PostgreSQL + pgx
Priority 3: External KV store (Redis)
```

ดู task: "Implement SQLite repository" ใน backlog
