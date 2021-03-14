// Package storage provides a generic storage interface and its implementations.
package storage

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

	// Initialize initializes the storage if it hasn't been set up yet. Methods
	// like CreateToDo must be safe to call after running Initialize.
	//
	// For example, a SQL storage implementation should creates the required
	// database and tables if they don't exist yet.
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

	// Remove removes the storage. It is the inverse operation of Initialize.
	// Must be called before Close when wiping a storage.
	Remove() error

	// Close closes handles and other resources like database connections.
	Close() error
}
