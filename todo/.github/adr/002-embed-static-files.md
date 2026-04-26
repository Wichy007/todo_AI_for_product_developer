# ADR-002: Embed Static Files ใน Binary

**Date:** 2026-04-26  
**Status:** Accepted

---

## Context

Frontend files (HTML, CSS, JS) ต้องถูก serve พร้อมกับ Go server  
มีตัวเลือก:
1. Serve จาก disk path (relative/absolute)
2. Embed ใน binary ด้วย `//go:embed`
3. Reverse proxy ไปยัง static server แยก

## Decision

ใช้ `//go:embed` directive เพื่อ embed `web/static/` เข้าไปใน binary

```go
//go:embed web/static
var staticFiles embed.FS

// served via
static, _ := fs.Sub(staticFiles, "web/static")
mux.Handle("/", http.FileServer(http.FS(static)))
```

## Rationale

| Criteria | Disk path | Embed | Separate server |
|---|---|---|---|
| Single binary deploy | ❌ | ✅ | ❌ |
| Hot reload (dev) | ✅ | ❌ (rebuild) | ✅ |
| Complexity | ต่ำ | ต่ำ | สูง |
| Port ที่ต้องจัดการ | 1 | 1 | 2+ |
| 12-Factor friendly | ❌ | ✅ | ✅ |

สำหรับ local development ความต้องการเดียวคือ single binary ที่รันง่าย

## Consequences

- **ดี:** `go build` ได้ binary เดียวที่รันได้ทันที ไม่ต้อง copy static files
- **ดี:** ไม่มี path mismatch ระหว่าง development และ production
- **ข้อจำกัด:** แก้ CSS/JS ต้อง rebuild — ถ้าต้องการ hot reload ให้ใช้ air หรือ serve จาก disk ใน dev mode
- **ข้อจำกัด:** binary size ใหญ่ขึ้นตาม static files

## Dev Hot Reload (Optional)

สร้าง build tag สำหรับ dev mode ที่ serve จาก disk แทน:

```go
//go:build dev
var staticFiles = os.DirFS("web/static")
```
