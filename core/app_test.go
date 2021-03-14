// Package core provides the core application functionality and business logic.
package core

import (
	"errors"
	"testing"

	"github.com/dominikbraun/todo/model"
	"github.com/dominikbraun/todo/storage"
)

// newTestApp creates a new app that is backed by an in-memory storage.
func newTestApp() *App {
	return &App{
		storage: storage.NewMemory(),
	}
}

func TestApp_CreateToDo(t *testing.T) {
	app := newTestApp()
	toDo := model.ToDo{
		Tasks: []model.Task{
			{
				Name: "Task 1",
			},
			{
				Name: "Task 2",
			},
		},
	}

	_, err := app.CreateToDo(toDo)
	if !errors.Is(err, ErrNameMustNotBeEmpty) {
		t.Errorf("expected error %v, got %v", ErrNameMustNotBeEmpty, err)
	}

	toDo.Name = "ToDo 1"

	_, err = app.CreateToDo(toDo)
	if err != nil {
		t.Fatalf("error creating ToDo: %s", err.Error())
	}
}

func TestApp_UpdateToDo(t *testing.T) {
	app := newTestApp()
	toDo := model.ToDo{
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

	createdToDo, _ := app.storage.CreateToDo(toDo)
	createdToDo.Name = ""

	err := app.UpdateToDo(createdToDo.ID, createdToDo)
	if !errors.Is(err, ErrNameMustNotBeEmpty) {
		t.Errorf("expected error %v, got %v", ErrNameMustNotBeEmpty, err)
	}

	createdToDo.Name = toDo.Name

	err = app.UpdateToDo(createdToDo.ID, createdToDo)
	if err != nil {
		t.Fatalf("error updating ToDo: %s", err.Error())
	}
}
