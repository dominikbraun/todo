// Package storage provides a generic, database-agnostic storage interface as
// well as a default implementation for MariaDB.
package storage

import "github.com/dominikbraun/todo/model"

// Storage represents a storage backend.
type Storage interface {

	// Install initializes the storage so that calls to the other functions like
	// CreateToDo are safe to perform.
	//
	// For example, a MongoDB backend should create the required buckets to
	// store ToDo items, while a MariaDB backend should create the corresponding
	// tables.
	Install() error

	// CreateToDo stores a new ToDo item and returns the inserted entity.
	CreateToDo(toDo model.ToDo) (model.ToDo, error)

	// FindToDos returns a list of all stored ToDo items.
	FindToDos() ([]model.ToDo, error)

	// FindToDoById returns the ToDo item with the given ID. In case the item
	// cannot be found, an error will be returned.
	FindToDoByID(id int64) (model.ToDo, error)

	// UpdateToDo overwrites the ToDo item with the given ID. In case the item
	// cannot be found, an error will be returned.
	UpdateToDo(id int64, toDo model.ToDo) (model.ToDo, error)

	// DeleteToDo deletes the ToDo item with the given ID. In case the item
	// cannot be found, an error will be returned.
	DeleteToDo(id int64) error
}
