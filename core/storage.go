// Package core provides the core application functionality and business logic
// triggered by controllers like the REST controller.
package core

import (
	"errors"

	"github.com/dominikbraun/todo/model"
)

var (
	// ErrToDoNotFound indicates that a requested ToDo item cannot be found.
	ErrToDoNotFound = errors.New("requested ToDo item not found")
)

// Storage represents a storage backend.
type Storage interface {

	// Initialize sets up the storage if it hasn't been set up yet. Methods like
	// CreateToDo should be safe to call after running Initialize.
	//
	// For example, a SQL storage should create the required database and tables
	// if they don't exist yet. Otherwise, nothing should happen.
	Initialize() error

	// CreateToDo stores a new ToDo item and returns the inserted entity.
	CreateToDo(toDo model.ToDo) (model.ToDo, error)

	// FindToDos returns a list of all stored ToDo items.
	FindToDos() ([]model.ToDo, error)

	// FindToDoById returns the ToDo item with the given ID. In case the item
	// cannot be found, an error will be returned.
	FindToDoByID(id int64) (model.ToDo, error)

	// UpdateToDo overwrites the ToDo item with the given ID. In case the item
	// cannot be found, an error will be returned.
	UpdateToDo(id int64, toDo model.ToDo) error

	// DeleteToDo deletes the ToDo item with the given ID. In case the item
	// cannot be found, an error will be returned.
	DeleteToDo(id int64) error
}
