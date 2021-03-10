// Package storage provides a generic, database-agnostic storage interface
// as well as a default implementation for MariaDB.
package storage

import "github.com/dominikbraun/todo/model"

// ToDoFilter is a function for filtering ToDo items. A storage implementation
// will pass each item to the filter function, and if it returns true for the
// particular item, it will be returned in the result.
//
// For example, to only retrieve ToDo items with 4 sub-tasks:
//
//	filter := func(toDo model.ToDo) bool {
//		return len(toDo.Tasks) == 4
//	}
//	toDos, err := storage.FindToDos(filter)
//
// To retrieve all items, just always return true.
type ToDoFilter func(toDo model.ToDo) bool

// Storage represents a storage backend. This may be an in-memory storage,
// a relational database or a simple key-value store.
type Storage interface {
	Install() error
	CreateTodo(toDo model.ToDo) (model.ToDo, error)
	FindToDos(filter ToDoFilter) ([]model.ToDo, error)
	FindToDoById(id int64) (model.ToDo, error)
	UpdateToDo(id int64, toDo model.ToDo) (model.ToDo, error)
	DeleteToDo(id int64) error
}
