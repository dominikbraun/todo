// Package todo provides the core application functionality and business logic
// triggered by controllers like the REST controller.
package todo

import (
	"database/sql"
	"errors"

	"github.com/dominikbraun/todo/model"
	"github.com/dominikbraun/todo/storage"
)

var (
	// ErrToDoNotFound indicates that a requested ToDo item cannot be found.
	ErrToDoNotFound = errors.New("requested ToDo item not found")
)

// App represents the core application. At this time, it merely consists of an
// arbitrary storage.Storage implementation for accessing ToDo items.
type App struct {
	storage storage.Storage
}

// New creates a new App instance that persists data to the given storage. The
// storage has to be fully initialized and ready-to-use for the app.
func New(storage storage.Storage) *App {
	return &App{
		storage: storage,
	}
}

// CreateToDo creates a new ToDo item. The provided item should not have an ID.
// Instead, CreateToDo will return the freshly created item including its ID.
func (a *App) CreateToDo(toDo model.ToDo) (model.ToDo, error) {
	return a.storage.CreateToDo(toDo)
}

// GetToDos returns a list of all stored ToDo items.
func (a *App) GetToDos() ([]model.ToDo, error) {
	return a.storage.FindToDos()
}

// GetToDo returns the ToDo with the given ID or an error if it doesn't exist.
func (a *App) GetToDo(id int64) (model.ToDo, error) {
	toDo, err := a.storage.FindToDoByID(id)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ToDo{}, ErrToDoNotFound
	} else if err != nil {
		return model.ToDo{}, err
	}

	return toDo, nil
}

// UpdateToDo updates a ToDo item by replacing the stored item with the given ID
// with the provided item. Depending on the underlying storage, the IDs of the
// ToDo's sub-tasks may change.
func (a *App) UpdateToDo(id int64, toDo model.ToDo) (model.ToDo, error) {
	toDo, err := a.storage.UpdateToDo(id, toDo)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ToDo{}, ErrToDoNotFound
	} else if err != nil {
		return model.ToDo{}, err
	}

	return toDo, nil
}

// DeleteToDo deletes the ToDo item with the given ID along with its sub-tasks.
func (a *App) DeleteToDo(id int64) error {
	err := a.storage.DeleteToDo(id)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrToDoNotFound
	} else if err != nil {
		return err
	}

	return nil
}
