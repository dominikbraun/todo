package storage

import (
	"testing"

	"github.com/dominikbraun/todo/model"
)

func TestMemory_Initialize(t *testing.T) {
	memory := NewMemory()

	if err := memory.Initialize(); err != nil {
		t.Fatalf("error initializing memory: %s", err.Error())
	}

	if memory.internal == nil {
		t.Error("expected internal memory to be non-nil")
	}
}

func TestMemory_CreateToDo(t *testing.T) {
	memory := NewMemory()
	toDo := model.ToDo{
		ID:   1,
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				Name: "Task 1",
			},
			{
				Name: "Task 2",
			},
		},
	}

	createdToDo, err := memory.CreateToDo(toDo)
	if err != nil {
		t.Fatalf("error creating ToDo: %s", err.Error())
	}

	assertEqual(t, toDo, createdToDo)
}

func TestMemory_FindToDos(t *testing.T) {
	memory := NewMemory()
	toDos := []model.ToDo{
		{
			ID:   1,
			Name: "ToDo 1",
			Tasks: []model.Task{
				{
					Name: "Task 1",
				},
				{
					Name: "Task 2",
				},
			},
		},
		{
			ID:   2,
			Name: "ToDo 2",
			Tasks: []model.Task{
				{
					Name: "Task 3",
				},
				{
					Name: "Task 4",
				},
			},
		},
	}

	for _, toDo := range toDos {
		memory.internal[toDo.ID] = toDo
	}

	foundToDos, err := memory.FindToDos()
	if err != nil {
		t.Fatalf("error finding ToDos: %s", err.Error())
	}

	if len(foundToDos) != len(toDos) {
		t.Fatalf("expected %d ToDos, found %d", len(toDos), len(foundToDos))
	}

	for i, foundToDo := range foundToDos {
		assertEqual(t, foundToDo, toDos[i])
	}
}

func TestMemory_FindToDoByID(t *testing.T) {
	memory := NewMemory()
	toDo := model.ToDo{ID: 1}

	memory.internal[toDo.ID] = toDo

	foundToDo, err := memory.FindToDoByID(toDo.ID)
	if err != nil {
		t.Fatalf("error finding ToDo: %s", err.Error())
	}

	assertEqual(t, toDo, foundToDo)

	_, err = memory.FindToDoByID(42)
	if err == nil {
		t.Fatalf("expected an error for non-existing ToDo")
	}
}

func TestMemory_UpdateToDo(t *testing.T) {
	memory := NewMemory()
	existingToDo := model.ToDo{
		ID:   1,
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
			{
				ID:   2,
				Name: "Task 2",
			},
		},
	}
	newToDo := model.ToDo{
		ID:   1,
		Name: "My ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
			{
				ID:   2,
				Name: "My Task 2",
			},
			{
				Name: "Task 3",
			},
		},
	}

	memory.internal[existingToDo.ID] = existingToDo

	err := memory.UpdateToDo(existingToDo.ID, newToDo)
	if err != nil {
		t.Fatalf("error updating ToDo: %s", err.Error())
	}

	updatedToDo := memory.internal[existingToDo.ID]

	assertEqual(t, newToDo, updatedToDo)
}

func TestMemory_DeleteToDo(t *testing.T) {
	memory := NewMemory()
	toDo := model.ToDo{
		ID:   1,
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
			{
				ID:   2,
				Name: "Task 2",
			},
		},
	}

	memory.internal[toDo.ID] = toDo

	if err := memory.DeleteToDo(toDo.ID); err != nil {
		t.Fatalf("error deleting ToDo: %s", err.Error())
	}

	if _, exists := memory.internal[toDo.ID]; exists {
		t.Errorf("ToDo still exists")
	}
}

// assertEqual tests whether two ToDos, an expected and the actual one, are
// equal. If they aren't, the given testing.T will be used to log an error.
func assertEqual(t *testing.T, expected, actual model.ToDo) {
	if expected.ID != actual.ID {
		format := "expected ID %d, got %d"
		t.Errorf(format, expected.ID, actual.ID)
	}

	if expected.Name != actual.Name {
		format := "expected name %s, got %s"
		t.Errorf(format, expected.Name, actual.Name)
	}

	if expected.Description != actual.Description {
		format := "expected description %s, got %s"
		t.Errorf(format, expected.Description, actual.Description)
	}

	if len(expected.Tasks) != len(actual.Tasks) {
		format := "expected %d tasks, got %d"
		t.Fatalf(format, len(expected.Tasks), len(actual.Tasks))
	}

	for i, task := range expected.Tasks {
		actualTask := actual.Tasks[i]

		if task.ID != actualTask.ID {
			format := "expected task ID %d, got %d"
			t.Errorf(format, task.ID, actualTask.ID)
		}

		if task.Name != actualTask.Name {
			format := "expected task name %s, got %s"
			t.Errorf(format, task.Name, actualTask.Name)
		}

		if task.Description != actualTask.Description {
			format := "expected task description %s, got %s"
			t.Errorf(format, task.Description, actualTask.Description)
		}
	}
}
