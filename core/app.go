// Package core provides the core application functionality and business logic.
package core

import (
	"github.com/dominikbraun/todo/model"
)

// App represents the core application. At this time, it merely consists of an
// arbitrary storage.Storage implementation for accessing ToDo items.
type App struct {
	storage Storage
}

// NewApp creates a new App instance that persists data to the given storage.
func NewApp(storage Storage) *App {
	return &App{
		storage: storage,
	}
}

// CreateToDo creates a new ToDo item. The provided item should not have an ID.
func (a *App) CreateToDo(toDo model.ToDo) (model.ToDo, error) {
	return a.storage.CreateToDo(toDo)
}

// GetToDos returns a list of all stored ToDo items.
func (a *App) GetToDos() ([]model.ToDo, error) {
	return a.storage.FindToDos()
}

// GetToDo returns the ToDo with the given ID or an error if it doesn't exist.
func (a *App) GetToDo(id int64) (model.ToDo, error) {
	return a.storage.FindToDoByID(id)
}

// UpdateToDo updates a ToDo item by replacing the stored item with the given ID
// with the provided item. Depending on the underlying storage, the IDs of the
// ToDo's sub-tasks may change.
func (a *App) UpdateToDo(id int64, toDo model.ToDo) error {
	return a.storage.UpdateToDo(id, toDo)
}

// DeleteToDo deletes the ToDo item with the given ID along with its sub-tasks.
func (a *App) DeleteToDo(id int64) error {
	return a.storage.DeleteToDo(id)
}
