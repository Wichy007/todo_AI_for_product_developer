package service_test

import (
	"testing"

	"hello/todo/internal/domain"
	"hello/todo/internal/repository"
	"hello/todo/internal/service"
)

func newService() *service.TaskService {
	return service.NewTaskService(repository.NewMemoryTaskRepository())
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestService_Create(t *testing.T) {
	task, err := newService().Create("Buy milk")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if task.ID == "" {
		t.Error("ID should not be empty")
	}
	if task.Text != "Buy milk" {
		t.Errorf("text = %q, want %q", task.Text, "Buy milk")
	}
}

func TestService_Create_EmptyText(t *testing.T) {
	_, err := newService().Create("")
	if err != domain.ErrEmptyTaskText {
		t.Errorf("error = %v, want ErrEmptyTaskText", err)
	}
}

func TestService_Create_IDsAreUnique(t *testing.T) {
	svc := newService()
	t1, _ := svc.Create("A")
	t2, _ := svc.Create("B")
	if t1.ID == t2.ID {
		t.Error("IDs should be unique")
	}
}

// ─── GetAll ───────────────────────────────────────────────────────────────────

func TestService_GetAll_Empty(t *testing.T) {
	tasks, err := newService().GetAll()
	if err != nil {
		t.Fatalf("GetAll error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestService_Update_ToggleComplete(t *testing.T) {
	svc := newService()
	task, _ := svc.Create("Write tests")

	completed := true
	updated, err := svc.Update(task.ID, nil, &completed)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if !updated.Completed {
		t.Error("expected Completed = true")
	}
}

func TestService_Update_RenameTask(t *testing.T) {
	svc := newService()
	task, _ := svc.Create("Old name")

	newText := "New name"
	updated, err := svc.Update(task.ID, &newText, nil)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if updated.Text != "New name" {
		t.Errorf("text = %q, want %q", updated.Text, "New name")
	}
}

func TestService_Update_EmptyText(t *testing.T) {
	svc := newService()
	task, _ := svc.Create("Something")

	empty := ""
	_, err := svc.Update(task.ID, &empty, nil)
	if err != domain.ErrEmptyTaskText {
		t.Errorf("error = %v, want ErrEmptyTaskText", err)
	}
}

func TestService_Update_NotFound(t *testing.T) {
	_, err := newService().Update("ghost", nil, nil)
	if err != domain.ErrTaskNotFound {
		t.Errorf("error = %v, want ErrTaskNotFound", err)
	}
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestService_Delete(t *testing.T) {
	svc := newService()
	task, _ := svc.Create("Delete me")

	if err := svc.Delete(task.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	tasks, _ := svc.GetAll()
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks after delete, got %d", len(tasks))
	}
}

func TestService_Delete_NotFound(t *testing.T) {
	err := newService().Delete("ghost")
	if err != domain.ErrTaskNotFound {
		t.Errorf("error = %v, want ErrTaskNotFound", err)
	}
}

// ─── Reorder ──────────────────────────────────────────────────────────────────

func TestService_Reorder(t *testing.T) {
	svc := newService()
	t1, _ := svc.Create("First")
	t2, _ := svc.Create("Second")
	t3, _ := svc.Create("Third")

	if err := svc.Reorder([]string{t3.ID, t1.ID, t2.ID}); err != nil {
		t.Fatalf("Reorder error: %v", err)
	}
	tasks, _ := svc.GetAll()
	if tasks[0].ID != t3.ID {
		t.Errorf("first task = %s, want %s", tasks[0].ID, t3.ID)
	}
}
