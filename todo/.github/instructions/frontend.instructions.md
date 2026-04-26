---
description: "Use when editing index.html, style.css, or app.js. Covers JS architecture, CSS token system, accessibility, and mobile-first responsive rules."
applyTo: "todo/web/static/**"
---

# Frontend Guidelines

## JavaScript Architecture

### API Layer — ใช้เสมอ, ห้าม fetch โดยตรง

```js
// ✅ ทุก HTTP call ต้องผ่าน api object
const api = {
  list:    ()          => apiFetch(API),
  create:  (text)      => apiFetch(API, { method: 'POST', body: JSON.stringify({ text }) }),
  update:  (id, patch) => apiFetch(`${API}/${id}`, { method: 'PUT', ... }),
  remove:  (id)        => apiFetch(`${API}/${id}`, { method: 'DELETE' }),
  reorder: (ids)       => apiFetch(`${API}/reorder`, { method: 'PUT', ... }),
};

// ❌
const res = await fetch('/api/tasks', { ... }); // ในตัว event listener
```

### DOM Mutation Pattern

```js
// ✅ renderTask() สร้าง element แล้ว append
taskList.appendChild(renderTask(task));

// ❌ อย่าทำ innerHTML ของ list ทั้งหมดโดยไม่จำเป็น
taskList.innerHTML = tasks.map(renderTask).join(''); // ทำลาย event listeners
```

### XSS Prevention

ทุก user-generated content ต้อง escape ก่อน inject ใน HTML:

```js
// ✅ ใช้ escapeHtml() เสมอ
li.innerHTML = `<span>${escapeHtml(task.text)}</span>`;

// ❌
li.innerHTML = `<span>${task.text}</span>`;
```

### Error Handling ใน Async

```js
// ✅ rollback UI state เมื่อ API ล้มเหลว
checkbox.addEventListener('change', async () => {
  try {
    await api.update(task.id, { completed: checkbox.checked });
  } catch (err) {
    console.error('Update failed:', err);
    checkbox.checked = !checkbox.checked; // rollback
  }
});
```

---

## CSS Design System

### Design Tokens (CSS Variables)

**ห้าม hardcode สี, shadow, radius** — ใช้ token เสมอ:

```css
/* ✅ */
background: var(--surface);
border-radius: var(--radius);
box-shadow: var(--shadow);

/* ❌ */
background: #ffffff;
border-radius: 12px;
```

Token ทั้งหมดนิยามใน `:root` ที่ต้น `style.css`:

| Token | ใช้สำหรับ |
|---|---|
| `--bg` | page background |
| `--surface` | card/modal background |
| `--primary` | interactive elements, FAB |
| `--danger` | delete actions |
| `--text` | primary text |
| `--muted` | placeholder, secondary text |
| `--done` | completed task text |
| `--border` | dividers, input borders |
| `--radius` | border radius ทุกที่ |
| `--shadow` / `--shadow-lg` | card / elevated shadows |

### Mobile-First Responsive

เขียน base styles สำหรับ mobile ก่อน แล้วค่อย override ด้วย `@media (min-width: ...)`:

```css
/* ✅ mobile base */
.overlay { align-items: flex-end; }

/* ✅ desktop override */
@media (min-width: 480px) {
  .overlay { align-items: center; }
}
```

### `hidden` Attribute Override

ถ้า element มี `display: flex/grid` อยู่ ต้องเพิ่ม rule สำหรับ `[hidden]`:

```css
.overlay { display: flex; }
.overlay[hidden] { display: none; }  /* ต้องมีคู่กันเสมอ */
```

---

## Accessibility Minimums

- ทุก interactive element ต้องมี `aria-label` ถ้าไม่มี visible text
- ทุก `<input>` ต้องมี label หรือ `aria-label`
- modal ต้องมี `role="dialog"` และ `aria-modal="true"`
- FAB และปุ่มที่มีแค่ icon ต้องมี `title` attribute

---

## Drag-and-Drop Pattern

ใช้ Pointer Events API (`pointerdown/pointermove/pointerup`) เท่านั้น:
- รองรับทั้ง mouse และ touch บนมือถือโดยไม่ต้องเพิ่ม library
- ต้อง `e.preventDefault()` ใน `pointerdown` เพื่อป้องกัน text selection
- ghost element ต้องเป็น `position: fixed; pointer-events: none`
- cleanup ghost และ placeholder ใน `pointerup` เสมอ
