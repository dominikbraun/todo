package storage

import "github.com/dominikbraun/todo/model"

type ToDoFilter func(toDo model.ToDo) bool

type Storage interface {
	CreateTodo(toDo model.ToDo) error
	FindToDos(filter ToDoFilter) ([]model.ToDo, error)
	UpdateToDo(toDo model.ToDo) error
	DeleteToDo(toDo model.ToDo) error
}
