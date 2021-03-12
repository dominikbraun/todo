// Package storage provides implementations of the core.Storage interface.
package storage

import (
	"github.com/dominikbraun/todo/core"
	"github.com/dominikbraun/todo/model"
)

type memory struct {
	internal map[int64]model.ToDo
	toDoID   int64
	taskID   int64
}

// NewMemory creates an in-memory storage living as long as the server process.
func NewMemory() *memory {
	return &memory{
		internal: make(map[int64]model.ToDo),
		toDoID:   0,
		taskID:   0,
	}
}

// Initialize initializes the in-memory storage by creating a hash map.
func (m *memory) Initialize() error {
	if m.internal == nil {
		m.internal = make(map[int64]model.ToDo)
	}

	return nil
}

// CreateToDo stores the given ToDo item in memory. CreateToDo expects an item
// without ID and will, similarly to the MariaDB storage implementation, assign
// an auto-incremented ID to each task.
func (m *memory) CreateToDo(toDo model.ToDo) (model.ToDo, error) {
	for i, _ := range toDo.Tasks {
		m.taskID++
		toDo.Tasks[i].ID = m.taskID
	}

	m.toDoID++
	toDo.ID = m.toDoID

	m.internal[toDo.ID] = toDo

	return toDo, nil
}

// FindToDos returns all ToDo items stored in memory.
func (m *memory) FindToDos() ([]model.ToDo, error) {
	toDos := make([]model.ToDo, len(m.internal))

	for i, toDo := range m.internal {
		toDos[i] = toDo
	}

	return toDos, nil
}

// FindToDoByID looks for a ToDo item with the provided ID and returns that item
// if it was found. Otherwise, core.ErrToDoNotFound will be returned.
func (m *memory) FindToDoByID(id int64) (model.ToDo, error) {
	if toDo, exists := m.internal[id]; exists {
		return toDo, nil
	}

	return model.ToDo{}, core.ErrToDoNotFound
}

// UpdateToDo overwrites a stored ToDo item with the provided ToDo instance. If
// the requested ToDo cannot be found, core.ErrToDoNotFound will be returned.
//
// Just like mariaDB.UpdateToDo, this function makes sure that IDs of existing
// tasks will not change: If a task has no ID assigned, it is considered to be
// new and will receive an ID. All other tasks, regardless whether they were
// modified or removed, will be overridden with the tasks of the new ToDo item.
func (m *memory) UpdateToDo(id int64, toDo model.ToDo) error {
	if _, exists := m.internal[id]; !exists {
		return core.ErrToDoNotFound
	}

	for i, task := range toDo.Tasks {
		if task.ID == 0 {
			m.taskID++
			toDo.Tasks[i].ID = m.taskID
		}
	}

	m.internal[id] = toDo

	return nil
}

// DeleteToDo deletes the ToDo item with the given ID. If the ToDo item cannot
// be found, core.ErrToDoNotFound will be returned.
func (m *memory) DeleteToDo(id int64) error {
	if _, exists := m.internal[id]; !exists {
		return core.ErrToDoNotFound
	}
	delete(m.internal, id)
	return nil
}
