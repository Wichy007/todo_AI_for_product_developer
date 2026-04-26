/* ─── API layer ──────────────────────────────────────────────────────────── */
const API = "/api/tasks";

async function apiFetch(path, options = {}) {
  const res = await fetch(path, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  if (res.status === 204) return null;
  return res.json();
}

const api = {
  list: () => apiFetch(API),
  create: (text) =>
    apiFetch(API, { method: "POST", body: JSON.stringify({ text }) }),
  update: (id, patch) =>
    apiFetch(`${API}/${id}`, { method: "PUT", body: JSON.stringify(patch) }),
  remove: (id) => apiFetch(`${API}/${id}`, { method: "DELETE" }),
  reorder: (ids) =>
    apiFetch(`${API}/reorder`, {
      method: "PUT",
      body: JSON.stringify({ ids }),
    }),
};

/* ─── DOM refs ───────────────────────────────────────────────────────────── */
const taskList = document.getElementById("task-list");
const emptyMsg = document.getElementById("empty-msg");
const addBtn = document.getElementById("add-btn");
const overlay = document.getElementById("overlay");
const taskInput = document.getElementById("task-input");
const cancelBtn = document.getElementById("cancel-btn");
const confirmBtn = document.getElementById("confirm-btn");

/* ─── Render helpers ─────────────────────────────────────────────────────── */

/** Escape user input to prevent XSS. */
function escapeHtml(str) {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

function toggleEmpty() {
  emptyMsg.hidden = taskList.children.length > 0;
}

/** Build and return a <li> element for the given task. */
function renderTask(task) {
  const li = document.createElement("li");
  li.className = "task-item";
  li.dataset.id = task.id;
  li.innerHTML = `
    <span class="drag-handle" aria-hidden="true">⠿</span>
    <input type="checkbox" class="checkbox" ${task.completed ? "checked" : ""}
           aria-label="Mark complete" />
    <span class="task-text ${task.completed ? "task-text--done" : ""}">${escapeHtml(task.text)}</span>
    <button class="delete-btn" aria-label="Delete task" title="Delete">✕</button>
  `;

  const checkbox = li.querySelector(".checkbox");
  const label = li.querySelector(".task-text");
  const delBtn = li.querySelector(".delete-btn");

  checkbox.addEventListener("change", async () => {
    try {
      await api.update(task.id, { completed: checkbox.checked });
      label.classList.toggle("task-text--done", checkbox.checked);
    } catch (err) {
      console.error("Update failed:", err);
      checkbox.checked = !checkbox.checked; // rollback
    }
  });

  delBtn.addEventListener("click", async () => {
    try {
      await api.remove(task.id);
      li.remove();
      toggleEmpty();
    } catch (err) {
      console.error("Delete failed:", err);
    }
  });

  return li;
}

async function loadTasks() {
  try {
    const tasks = await api.list();
    taskList.innerHTML = "";
    tasks.forEach((t) => taskList.appendChild(renderTask(t)));
    toggleEmpty();
  } catch (err) {
    console.error("Load failed:", err);
  }
}

/* ─── Drag-and-drop (ghost + placeholder) ───────────────────────────────── */

let dragState = null;

/** Clone the source item into a fixed-position ghost that follows the cursor. */
function createGhost(source) {
  const rect = source.getBoundingClientRect();
  const ghost = source.cloneNode(true);
  ghost.classList.add("drag-ghost");
  ghost.style.width = rect.width + "px";
  ghost.style.height = rect.height + "px";
  ghost.style.left = rect.left + "px";
  document.body.appendChild(ghost);
  return ghost;
}

/** Empty list item that shows where the dragged item will land. */
function createPlaceholder(source) {
  const ph = document.createElement("li");
  ph.className = "drag-placeholder";
  ph.style.height = source.getBoundingClientRect().height + "px";
  return ph;
}

/**
 * Returns the task item and whether to insert before it,
 * based on where the pointer currently is.
 */
function dropTarget(clientY) {
  const items = [...taskList.querySelectorAll(".task-item")];
  for (const item of items) {
    const { top, height } = item.getBoundingClientRect();
    if (clientY < top + height / 2) return { el: item, before: true };
  }
  return { el: null, before: false };
}

document.addEventListener("pointerdown", (e) => {
  if (!e.target.closest(".drag-handle")) return;
  const source = e.target.closest(".task-item");
  if (!source) return;
  e.preventDefault();

  const offsetY = e.clientY - source.getBoundingClientRect().top;
  const ghost = createGhost(source);
  const placeholder = createPlaceholder(source);

  ghost.style.top = e.clientY - offsetY + "px";
  source.after(placeholder);
  source.classList.add("drag-source");

  dragState = { source, ghost, placeholder, offsetY };
});

document.addEventListener("pointermove", (e) => {
  if (!dragState) return;
  const { ghost, placeholder, source, offsetY } = dragState;

  ghost.style.top = e.clientY - offsetY + "px";

  const { el, before } = dropTarget(e.clientY);
  if (el) {
    taskList.insertBefore(placeholder, before ? el : el.nextSibling);
  } else {
    taskList.appendChild(placeholder);
  }
});

document.addEventListener("pointerup", async () => {
  if (!dragState) return;
  const { source, ghost, placeholder } = dragState;
  dragState = null;

  ghost.remove();
  source.classList.remove("drag-source");
  placeholder.replaceWith(source);

  const ids = [...taskList.querySelectorAll(".task-item")].map(
    (li) => li.dataset.id,
  );
  try {
    await api.reorder(ids);
  } catch (err) {
    console.error("Reorder failed:", err);
    loadTasks();
  }
});

/* ─── Modal ──────────────────────────────────────────────────────────────── */

function openModal() {
  taskInput.value = "";
  overlay.hidden = false;
  taskInput.focus();
}

function closeModal() {
  overlay.hidden = true;
}

async function addTask() {
  const text = taskInput.value.trim();
  if (!text) return;
  try {
    const task = await api.create(text);
    taskList.appendChild(renderTask(task));
    toggleEmpty();
    closeModal();
  } catch (err) {
    console.error("Create failed:", err);
  }
}

addBtn.addEventListener("click", openModal);
cancelBtn.addEventListener("click", closeModal);
confirmBtn.addEventListener("click", addTask);
taskInput.addEventListener("keydown", (e) => {
  if (e.key === "Enter") addTask();
});
overlay.addEventListener("click", (e) => {
  if (e.target === overlay) closeModal();
});

/* ─── Boot ───────────────────────────────────────────────────────────────── */
loadTasks();
