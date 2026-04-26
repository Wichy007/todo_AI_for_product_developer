---
description: "Use when writing, reviewing, or refactoring any Go or JS code in this project. Covers TDD, Clean Code, SOLID, DRY, YAGNI, and code smell rules that must be followed."
---

# Engineering Practices

## TDD (Test-Driven Development)

**Red → Green → Refactor — ทำตามลำดับนี้เสมอ**

1. เขียน failing test ที่อธิบาย behavior ที่ต้องการก่อน
2. เขียน production code น้อยที่สุดเพื่อให้ test ผ่าน
3. refactor โดยไม่แก้ test

```go
// ✅ เขียน test ก่อนเสมอ
func TestService_Create_EmptyText(t *testing.T) {
    _, err := newService().Create("")
    if err != domain.ErrEmptyTaskText {
        t.Errorf("error = %v, want ErrEmptyTaskText", err)
    }
}
```

**กฎ:**
- test ต้องอยู่ใน `_test.go` ไฟล์เดียวกับ package ที่ทดสอบ
- ใช้ `package xxx_test` (black-box testing) ยกเว้นต้องทดสอบ internal
- ห้ามใช้ test framework ภายนอก — `testing` stdlib เท่านั้น (YAGNI)
- ทุก exported function ต้องมี test อย่างน้อย 1 happy path + 1 error path

---

## Clean Code

### ชื่อ (Naming)
- function/variable: สื่อความหมายโดยไม่ต้องมี comment
- ห้ามใช้ชื่อย่อที่ไม่ชัดเจน: `t`, `x`, `tmp` ยกเว้น loop index
- ชื่อ boolean: ขึ้นต้นด้วย `is`, `has`, `can`, `should`

```go
// ❌
func proc(t *Task, b bool) {}

// ✅
func markCompleted(task *Task, completed bool) {}
```

### Function Size
- Go function: ≤ 20 บรรทัด
- JS function: ≤ 25 บรรทัด
- ถ้าเกิน: แตก function ออก อย่าบีบ code

### Comments
- ห้าม comment อธิบายสิ่งที่ code พูดแทนได้
- comment เฉพาะ **ทำไม** ถึงทำแบบนี้ ไม่ใช่ **ทำอะไร**

```go
// ❌
// increment i by 1
i++

// ✅
// defensive copy — prevent callers from mutating repository state
cp := *task
return &cp, nil
```

---

## SOLID

### S — Single Responsibility
- แต่ละ struct/file ทำงานเดียว
- `TaskHandler` → HTTP routing เท่านั้น
- `TaskService` → business logic เท่านั้น
- `MemoryTaskRepository` → persistence เท่านั้น

### O — Open/Closed
- เพิ่ม behavior ผ่าน implementation ใหม่ ไม่แก้ existing code
- เปลี่ยน storage: สร้าง `PostgresTaskRepository` ที่ implement `domain.TaskRepository` ไม่แตะ service/handler

### L — Liskov Substitution
- ทุก implementation ของ `domain.TaskRepository` ต้องใช้แทนกันได้
- ห้าม method ของ concrete type panic หรือ return error ที่ interface ไม่รู้จัก

### I — Interface Segregation
- interface เล็กและเฉพาะเจาะจง
- `TaskRepository` มีแค่ 5 methods ที่จำเป็น — ห้าม bloat

### D — Dependency Inversion
- layer บนพึ่งพา interface ไม่ใช่ concrete type
- wiring ทำใน `main.go` เท่านั้น

```go
// ✅ service รับ interface
func NewTaskService(repo domain.TaskRepository) *TaskService

// ❌ service รับ concrete type
func NewTaskService(repo *repository.MemoryTaskRepository) *TaskService
```

---

## DRY (Don't Repeat Yourself)

- logic เดียวกัน ≥ 2 ที่ → extract เป็น function
- error mapping ใน `writeError()` — ห้ามทำซ้ำใน handler อื่น
- JS: helper `apiFetch()` ครอบทุก HTTP call — ห้ามเรียก `fetch()` ตรงๆ ใน event listener

---

## YAGNI (You Aren't Gonna Need It)

- ห้าม implement feature ที่ไม่มี requirement ณ ปัจจุบัน
- ห้ามเพิ่ม config option "เผื่อไว้"
- ห้ามสร้าง abstraction ที่มี implementation เดียว (ยกเว้น domain interface ที่วางไว้เพื่อ testability)

---

## Code Smells — สิ่งที่ต้อง Refactor ทันที

| Smell | อาการ | แก้ไข |
|---|---|---|
| God Function | function > 30 บรรทัด | แตก function |
| Magic Number | `if len > 200` | สร้าง constant |
| Deep Nesting | if ซ้อน > 3 ชั้น | early return / extract function |
| Feature Envy | function อ้างอิง field ของ struct อื่นมากกว่าของตัวเอง | ย้าย method |
| Shotgun Surgery | แก้ 1 feature ต้องแตะ > 3 ไฟล์ | consolidate |
| Primitive Obsession | ส่ง `string` แทน type | สร้าง type |

---

## Refactoring Rules

1. ห้าม refactor และเพิ่ม feature พร้อมกัน — แยก commit
2. ทุก refactor ต้องมี test ครอบ ก่อนเริ่ม refactor
3. ใช้ "Strangler Fig" pattern เมื่อ refactor layer ใหญ่ — ทำ parallel จนกว่าจะพร้อม
