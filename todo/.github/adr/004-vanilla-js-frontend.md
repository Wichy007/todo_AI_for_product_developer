# ADR-004: Vanilla JS Frontend (ไม่ใช้ Framework)

**Date:** 2026-04-26  
**Status:** Accepted

---

## Context

Frontend ต้องการ: list tasks, check/uncheck, drag-and-drop reorder, add/delete  
ตัวเลือก: React, Vue, Svelte, HTMX, Vanilla JS

## Decision

ใช้ **Vanilla JavaScript (ES2022)** + HTML + CSS โดยไม่มี build step

## Rationale

| Criteria | Vanilla JS | React/Vue | HTMX |
|---|---|---|---|
| Build step | ❌ ไม่ต้อง | ✅ ต้อง (webpack/vite) | ❌ ไม่ต้อง |
| Bundle size | ~3KB | ~100KB+ | ~14KB |
| Dependencies | 0 | หลาย | 1 |
| Drag-and-drop | Pointer Events | ต้องการ library | ไม่รองรับดี |
| Learning curve | ต่ำ | ปานกลาง-สูง | ต่ำ |
| SSR | N/A | ซับซ้อน | ✅ ง่าย |

Feature ที่ต้องการทั้งหมดทำได้ด้วย Web APIs มาตรฐาน:
- **Pointer Events** สำหรับ drag-and-drop (desktop + mobile)
- **fetch API** สำหรับ REST calls
- **DOM API** สำหรับ rendering

## Consequences

- **ดี:** embed ใน binary ได้ทันที ไม่ต้อง npm/node pipeline
- **ดี:** debug ง่าย ไม่มี virtual DOM abstraction
- **ดี:** ไม่มี dependency security risk จาก npm
- **ข้อจำกัด:** ถ้า UI ซับซ้อนขึ้นมาก (forms, real-time, etc.) ควรย้ายไป framework
- **ข้อจำกัด:** ไม่มี component system — ต้อง discipline ในการ organize code

## Migration Path

ถ้า feature เพิ่มจนถึงจุดที่ต้องการ component system:
1. เพิ่ม `vite` build step ใน `web/`
2. เขียน React/Svelte components
3. `go:embed web/dist` แทน `web/static`
4. ไม่แตะ Go backend เลย

## JS Conventions ที่บังคับใช้

- ทุก HTTP call ผ่าน `api` object (**DRY**)
- ทุก user input ผ่าน `escapeHtml()` (**Security**)
- Async error: log + rollback UI state (**Clean Code**)
- ไม่ใช้ `var` — ใช้ `const`/`let` เท่านั้น
