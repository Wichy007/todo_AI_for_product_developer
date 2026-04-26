# ADR-001: ใช้ Go 1.22 Standard Library Router

**Date:** 2026-04-26  
**Status:** Accepted

---

## Context

Project ต้องการ HTTP routing สำหรับ REST API ที่มี 5 endpoints  
มีตัวเลือกหลัก: Go stdlib `net/http`, [gorilla/mux](https://github.com/gorilla/mux), [chi](https://github.com/go-chi/chi), [gin](https://github.com/gin-gonic/gin)

## Decision

ใช้ **Go 1.22 standard library `net/http`** โดยใช้ method+pattern routing ที่เพิ่มเข้ามาใน Go 1.22

```go
mux.HandleFunc("GET /api/tasks", h.listTasks)
mux.HandleFunc("PUT /api/tasks/reorder", h.reorderTasks)
mux.HandleFunc("PUT /api/tasks/{id}", h.updateTask)
```

## Rationale

| Criteria | stdlib | third-party router |
|---|---|---|
| Dependencies | 0 | เพิ่ม 1+ dependency |
| Method routing | ✅ Go 1.22+ | ✅ ทุกตัว |
| Path parameters | ✅ `r.PathValue("id")` | ✅ ทุกตัว |
| Middleware | manual | built-in chain |
| Learning curve | ต่ำ | ต่ำ-ปานกลาง |

สำหรับ 5 endpoints ไม่จำเป็นต้องใช้ middleware chain หรือ advanced routing — YAGNI

## Consequences

- **ดี:** ไม่มี external dependency, upgrade Go เท่านั้น
- **ดี:** เข้าใจง่าย ไม่ต้องรู้ framework พิเศษ
- **ข้อจำกัด:** ถ้าต้องการ middleware chain หรือ route groups ในอนาคต ควรเปลี่ยนเป็น chi (compatible API)

## Migration Path

ถ้า route > 20 หรือต้องการ middleware → เปลี่ยนเป็น `chi` โดยไม่แตะ handler logic  
(chi ใช้ `http.HandlerFunc` เหมือน stdlib)
