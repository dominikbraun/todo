package storage

import (
	"github.com/dominikbraun/todo/model"
	"github.com/jmoiron/sqlx"
)

type mariaDB struct {
	db *sqlx.DB
}

func NewMariaDB() *mariaDB {
	return &mariaDB{}
}

func (m *mariaDB) Install() error {
	createToDoTable := `
CREATE TABLE todos (
	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	description VARCHAR(500)
)`

	if _, err := m.db.Exec(createToDoTable); err != nil {
		return err
	}

	createTaskTable := `
CREATE TABLE tasks (
	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	description VARCHAR(500),
	todo_id BIGINT UNSIGNED NOT NULL
)`

	if _, err := m.db.Exec(createTaskTable); err != nil {
		return err
	}

	return nil
}

func (m *mariaDB) CreateTodo(toDo model.ToDo) (model.ToDo, error) {
	insert := `INSERT INTO todos (name, description) VALUES (?, ?)`

	result, err := m.db.Exec(insert, toDo.Name, toDo.Description)
	if err != nil {
		return model.ToDo{}, err
	}

	id, _ := result.LastInsertId()
	toDo.ID = id

	for _, t := range toDo.Tasks {
		task, err := m.createTaskForToDo(toDo.ID, t)
		if err != nil {
			return model.ToDo{}, err
		}
		toDo.Tasks = append(toDo.Tasks, task)
	}

	return toDo, nil
}

func (m *mariaDB) FindToDos(filter ToDoFilter) ([]model.ToDo, error) {
	query := `SELECT id, name, description FROM todos`

	rows, err := m.db.Queryx(query)
	if err != nil {
		return nil, err
	}

	var toDos []model.ToDo

	for rows.Next() {
		var toDo model.ToDo
		if err := rows.StructScan(&toDo); err != nil {
			return nil, err
		}

		tasks, err := m.findTasksByToDo(toDo.ID)
		if err != nil {
			return nil, err
		}

		toDo.Tasks = tasks
		toDos = append(toDos, toDo)
	}

	return toDos, nil
}

func (m *mariaDB) FindToDoById(id int64) (model.ToDo, error) {
	query := `SELECT id, name, description FROM todos WHERE id = ?`

	var toDo model.ToDo

	if err := m.db.QueryRowx(query, id).StructScan(&toDo); err != nil {
		return model.ToDo{}, err
	}

	tasks, err := m.findTasksByToDo(toDo.ID)
	if err != nil {
		return model.ToDo{}, err
	}

	toDo.Tasks = tasks

	return toDo, nil
}

func (m *mariaDB) UpdateToDo(id int64, toDo model.ToDo) (model.ToDo, error) {
	deleteTasks := `DELETE FROM tasks WHERE todo_id = ?`

	_, err := m.db.Exec(deleteTasks, id)
	if err != nil {
		return model.ToDo{}, err
	}

	var newTasks []model.Task

	for _, t := range toDo.Tasks {
		task, err := m.createTaskForToDo(id, t)
		if err != nil {
			return model.ToDo{}, err
		}
		newTasks = append(newTasks, task)
	}

	updateToDo := `UPDATE todos SET name = ?, description = ? WHERE id = ?`

	_, err = m.db.Exec(updateToDo, toDo.Name, toDo.Description, id)
	if err != nil {
		return model.ToDo{}, err
	}

	toDo.ID = id
	toDo.Tasks = newTasks

	return toDo, nil
}

func (m *mariaDB) DeleteToDo(id int64) error {
	deleteTasks := `DELETE FROM tasks WHERE todo_id = ?`

	_, err := m.db.Exec(deleteTasks, id)
	if err != nil {
		return err
	}

	deleteToDo := `DELETE FROM todos WHERE id = ?`

	_, err = m.db.Exec(deleteToDo, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *mariaDB) createTaskForToDo(toDoId int64, task model.Task) (model.Task, error) {
	insert := `INSERT INTO tasks (name, description, todo_id) VALUES (?, ?, ?)`

	result, err := m.db.Exec(insert, task.Name, task.Description, toDoId)
	if err != nil {
		return model.Task{}, err
	}

	id, _ := result.LastInsertId()
	task.ID = id

	return task, nil
}

func (m *mariaDB) findTasksByToDo(toDoID int64) ([]model.Task, error) {
	query := `SELECT id, name, description FROM todos WHERE todo_id = ?`

	rows, err := m.db.Queryx(query, toDoID)
	if err != nil {
		return nil, err
	}

	var tasks []model.Task

	for rows.Next() {
		var task model.Task
		if err := rows.StructScan(&task); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
