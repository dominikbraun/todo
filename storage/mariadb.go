package storage

import (
	"github.com/dominikbraun/todo/model"
)

type mariaDB struct{}

func NewMariaDB() *mariaDB {
	return &mariaDB{}
}

func (m *mariaDB) CreateTodo(toDo model.ToDo) error {
	panic("implement me")
}

func (m *mariaDB) FindToDos(filter ToDoFilter) ([]model.ToDo, error) {
	panic("implement me")
}

func (m *mariaDB) FindToDoById(id int) (model.ToDo, error) {
	panic("implement me")
}

func (m *mariaDB) UpdateToDo(toDo model.ToDo) error {
	panic("implement me")
}

func (m *mariaDB) DeleteToDo(toDo model.ToDo) error {
	panic("implement me")
}
