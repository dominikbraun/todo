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
	ErrToDoNotFound = errors.New("requested ToDo item not found")
)

type App struct {
	storage storage.Storage
}

func New(storage storage.Storage) *App {
	return &App{
		storage: storage,
	}
}

func (a *App) CreateToDo(toDo model.ToDo) (model.ToDo, error) {
	return a.storage.CreateToDo(toDo)
}

func (a *App) GetToDos() ([]model.ToDo, error) {
	return a.storage.FindToDos()
}

func (a *App) GetToDo(id int64) (model.ToDo, error) {
	toDo, err := a.storage.FindToDoByID(id)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ToDo{}, ErrToDoNotFound
	} else if err != nil {
		return model.ToDo{}, err
	}

	return toDo, nil
}

func (a *App) UpdateToDo(id int64, toDo model.ToDo) (model.ToDo, error) {
	toDo, err := a.storage.UpdateToDo(id, toDo)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ToDo{}, ErrToDoNotFound
	} else if err != nil {
		return model.ToDo{}, err
	}

	return toDo, nil
}

func (a *App) DeleteToDo(id int64) error {
	err := a.storage.DeleteToDo(id)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrToDoNotFound
	} else if err != nil {
		return err
	}

	return nil
}
