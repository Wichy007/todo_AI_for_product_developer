package repository_test

import (
	"testing"

	"todo/internal/domain"
	"todo/internal/repository"
)

func newRepo() *repository.MemoryTaskRepository {
	return repository.NewMemoryTaskRepository()
}

// ─── Save / FindAll ───────────────────────────────────────────────────────────

func TestMemory_SaveAndFindAll(t *testing.T) {
	repo := newRepo()
	_ = repo.Save(&domain.Task{ID: "1", Text: "Buy milk"})

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll error: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Text != "Buy milk" {
		t.Errorf("unexpected tasks: %+v", tasks)
	}
}

func TestMemory_Save_Update(t *testing.T) {
	repo := newRepo()
	task := &domain.Task{ID: "1", Text: "Original"}
	_ = repo.Save(task)

	task.Text = "Updated"
	_ = repo.Save(task)

	tasks, _ := repo.FindAll()
	if tasks[0].Text != "Updated" {
		t.Errorf("got %q, want %q", tasks[0].Text, "Updated")
	}
}

// ─── FindByID ─────────────────────────────────────────────────────────────────

func TestMemory_FindByID_NotFound(t *testing.T) {
	_, err := newRepo().FindByID("missing")
	if err != domain.ErrTaskNotFound {
		t.Errorf("error = %v, want ErrTaskNotFound", err)
	}
}

func TestMemory_FindByID_ReturnsCopy(t *testing.T) {
	repo := newRepo()
	_ = repo.Save(&domain.Task{ID: "1", Text: "Original"})

	got, _ := repo.FindByID("1")
	got.Text = "Mutated"

	got2, _ := repo.FindByID("1")
	if got2.Text != "Original" {
		t.Error("FindByID should return a defensive copy")
	}
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestMemory_Delete(t *testing.T) {
	repo := newRepo()
	_ = repo.Save(&domain.Task{ID: "1", Text: "Delete me"})
	if err := repo.Delete("1"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	tasks, _ := repo.FindAll()
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestMemory_Delete_NotFound(t *testing.T) {
	err := newRepo().Delete("ghost")
	if err != domain.ErrTaskNotFound {
		t.Errorf("error = %v, want ErrTaskNotFound", err)
	}
}

// ─── Reorder ──────────────────────────────────────────────────────────────────

func TestMemory_Reorder(t *testing.T) {
	repo := newRepo()
	for _, id := range []string{"a", "b", "c"} {
		_ = repo.Save(&domain.Task{ID: id, Text: id})
	}
	if err := repo.Reorder([]string{"c", "a", "b"}); err != nil {
		t.Fatalf("Reorder error: %v", err)
	}
	tasks, _ := repo.FindAll()
	if tasks[0].ID != "c" || tasks[1].ID != "a" || tasks[2].ID != "b" {
		t.Errorf("unexpected order: %v %v %v", tasks[0].ID, tasks[1].ID, tasks[2].ID)
	}
}

func TestMemory_Reorder_WrongLength(t *testing.T) {
	repo := newRepo()
	_ = repo.Save(&domain.Task{ID: "1", Text: "x"})
	err := repo.Reorder([]string{"1", "2"})
	if err != domain.ErrInvalidReorder {
		t.Errorf("error = %v, want ErrInvalidReorder", err)
	}
}

func TestMemory_Reorder_UnknownID(t *testing.T) {
	repo := newRepo()
	_ = repo.Save(&domain.Task{ID: "1", Text: "x"})
	err := repo.Reorder([]string{"ghost"})
	if err != domain.ErrTaskNotFound {
		t.Errorf("error = %v, want ErrTaskNotFound", err)
	}
}

func TestMemory_Reorder_DuplicateID(t *testing.T) {
	repo := newRepo()
	for _, id := range []string{"1", "2"} {
		_ = repo.Save(&domain.Task{ID: id, Text: id})
	}
	err := repo.Reorder([]string{"1", "1"})
	if err != domain.ErrInvalidReorder {
		t.Errorf("error = %v, want ErrInvalidReorder", err)
	}
}
