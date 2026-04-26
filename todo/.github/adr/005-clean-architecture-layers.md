# ADR-005: Clean Architecture 4-Layer Structure

**Date:** 2026-04-26  
**Status:** Accepted

---

## Context

ต้องการโครงสร้าง code ที่:
- testable โดยไม่ต้องมี HTTP server หรือ database จริง
- swap storage ได้โดยไม่แตะ business logic
- เพิ่ม delivery channel ใหม่ (CLI, gRPC) โดยไม่ duplicate logic

ตัวเลือก: Flat package, MVC, Clean Architecture, Hexagonal Architecture

## Decision

ใช้ **Clean Architecture** แบ่งเป็น 4 layers ภายใน `internal/`:

```
domain → repository (implements domain)
domain ← service (depends on domain interface)
service ← handler (depends on service)
main (wires everything)
```

## Layer Responsibilities

### `domain/` — Enterprise Business Rules
- **Task** entity และ fields
- **TaskRepository** interface (contract)
- **Sentinel errors** (`ErrTaskNotFound`, etc.)
- **Zero imports** จาก internal packages อื่น

### `repository/` — Interface Adapters (Infrastructure)
- Implement `domain.TaskRepository`
- จัดการ concurrency (mutex)
- Defensive copy ป้องกัน mutation

### `service/` — Application Business Rules  
- Validation (empty text, etc.)
- ID generation
- Orchestrate repository calls
- ไม่รู้จัก HTTP หรือ JSON

### `handler/` — Delivery Mechanism
- Parse HTTP request → call service
- Map service errors → HTTP status
- JSON encode response
- ไม่มี business logic

## Rationale

| Criteria | Flat | MVC | Clean Arch |
|---|---|---|---|
| Test without HTTP | ❌ | ❌ | ✅ |
| Swap storage | ❌ | ❌ | ✅ |
| Add delivery channel | ❌ | ❌ | ✅ |
| Code volume (small app) | น้อย | กลาง | มากกว่า |
| Onboarding time | ต่ำ | ต่ำ | ปานกลาง |

App นี้มีแผนขยาย storage → Clean Architecture คุ้มค่า

## Consequences

- **ดี:** test แต่ละ layer แยกกันได้อย่างอิสระ
- **ดี:** เปลี่ยน storage ใน `main.go` 1 บรรทัด
- **ดี:** เพิ่ม CLI หรือ gRPC ได้โดยไม่แตะ service/domain
- **ข้อจำกัด:** boilerplate มากกว่า flat structure สำหรับ app เล็ก
- **ข้อจำกัด:** ต้องเข้าใจ layer rule ก่อนเขียน code

## Import Boundary Enforcement

ตรวจ import violation ด้วย:
```bash
# ตรวจว่า service ไม่ import repository concrete
grep -r "repository\." todo/internal/service/
# ผลลัพธ์ต้องว่างเปล่า
```

หรือใช้ `golang.org/x/tools/go/analysis` สำหรับ CI check
