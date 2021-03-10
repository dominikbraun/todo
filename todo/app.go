// Package todo provides the core application functionality and business
// logic triggered by the individual controllers.
package todo

import (
	"errors"

	"github.com/dominikbraun/todo/model"
	"github.com/dominikbraun/todo/storage"
)

var (
	ErrToDoNotFound = errors.New("requested ToDo item not found")
)

type App struct {
	storage storage.Storage
}

func New() *App {
	return &App{
		storage: storage.NewMariaDB(),
	}
}

func (a *App) CreateToDo(toDo model.ToDo) (model.ToDo, error) {
	return model.ToDo{}, nil
}

func (a *App) GetToDos() ([]model.ToDo, error) {
	return nil, nil
}

func (a *App) GetToDo(id int) (model.ToDo, error) {
	return model.ToDo{}, nil
}

func (a *App) UpdateToDo(id int, toDo model.ToDo) (model.ToDo, error) {
	return model.ToDo{}, nil
}

func (a *App) DeleteToDo(id int) error {
	return nil
}
