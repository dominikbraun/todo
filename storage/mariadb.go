// Package storage provides implementations of the core.Storage interface.
package storage

import (
	"fmt"

	"github.com/dominikbraun/todo/model"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MariaDBConfig stores configuration values for connecting to the MariaDB host.
type MariaDBConfig struct {
	User     string
	Password string
	Address  string
}

// URI creates an URI for connecting to the configured database host. It yields
// a connection string in the form <user>:<password>@<host>:<port>/.
func (m MariaDBConfig) URI() string {
	return fmt.Sprintf("%s:%s@(%s)/", m.User, m.Password, m.Address)
}

type mariaDB struct {
	config MariaDBConfig
	db     *sqlx.DB
}

// NewMariaDB creates a new MariaDB connection using the given configuration.
func NewMariaDB(config MariaDBConfig) (*mariaDB, error) {
	db, err := sqlx.Connect("mysql", config.URI())
	if err != nil {
		return nil, err
	}

	return &mariaDB{
		config: config,
		db:     db,
	}, nil
}

// Initialize creates the database along with the required tables if they don't
// exist yet. After running Initialize without an error, all other operations
// are safe to perform.
func (m *mariaDB) Initialize() error {
	statements := []string{
		`CREATE DATABASE IF NOT EXISTS todo_app`,
		`USE todo_app`,
		`CREATE TABLE IF NOT EXISTS todos (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description VARCHAR(500)
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description VARCHAR(500),
			todo_id BIGINT UNSIGNED NOT NULL
		)`,
	}

	for _, statement := range statements {
		if _, err := m.db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

// CreateToDo inserts the given ToDo item into its table. CreateToDo expects a
// ToDo item without an ID.
func (m *mariaDB) CreateToDo(toDo model.ToDo) (model.ToDo, error) {
	sql, args, _ := squirrel.
		Insert("todos").
		Columns("name", "description").
		Values(toDo.Name, toDo.Description).
		ToSql()

	result, err := m.db.Exec(sql, args...)
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

// FindToDos returns all ToDo items stored in the MariaDB database.
func (m *mariaDB) FindToDos() ([]model.ToDo, error) {
	sql, _, _ := squirrel.
		Select("id", "name", "description").
		From("todos").
		ToSql()

	rows, err := m.db.Queryx(sql)
	if err != nil {
		return nil, err
	}

	var toDos []model.ToDo

	for rows.Next() {
		var toDo model.ToDo
		if err := rows.StructScan(&toDo); err != nil {
			return nil, err
		}

		tasks, err := m.findTasksByToDoID(toDo.ID)
		if err != nil {
			return nil, err
		}

		toDo.Tasks = tasks
		toDos = append(toDos, toDo)
	}

	return toDos, nil
}

// FindToDoByID searches a ToDo item with the provided ID and returns that item
// if it was found. Otherwise, core.ErrToDoNotFound will be returned.
func (m *mariaDB) FindToDoByID(id int64) (model.ToDo, error) {
	sql, args, _ := squirrel.
		Select("id", "name", "description").
		From("todos").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	var toDo model.ToDo

	if err := m.db.QueryRowx(sql, args...).StructScan(&toDo); err != nil {
		return model.ToDo{}, err
	}

	tasks, err := m.findTasksByToDoID(toDo.ID)
	if err != nil {
		return model.ToDo{}, err
	}

	toDo.Tasks = tasks

	return toDo, nil
}

// UpdateToDo overwrites a stored ToDo item with the provided ToDo instance. If
// the requested ToDo cannot be found, core.ErrToDoNotFound will be returned.
//
// The easiest way to update a ToDo along with its sub-tasks would be to delete
// all the tasks and insert the tasks listed in the new ToDo item. However, this
// would change the task IDs, which is probably not expected by an API client.
//
// To solve this problem, UpdateToDo clearly distinguishes between new, modified
// and removed tasks. UpdateToDo adheres to the following rules:
//
//	1. If a task has no ID assigned, it will be inserted.
//	2. If a task has an ID assigned, it will be updated.
//	3. If a task exists in the DB but not in the model, it will be deleted.
//
// For the sake of simplicity, tasks will be updated regardless whether they
// actually changed.
func (m *mariaDB) UpdateToDo(id int64, toDo model.ToDo) error {
	taskIDs := make([]int64, 0)

	for _, task := range toDo.Tasks {
		// If the task has an ID assigned, just update the task.
		if task.ID != 0 {
			sql, args, _ := squirrel.
				Update("tasks").
				Set("name", task.Name).
				Set("description", task.Description).
				Where(squirrel.Eq{"id": task.ID}).
				ToSql()

			if _, err := m.db.Exec(sql, args...); err != nil {
				return err
			}
			taskIDs = append(taskIDs, task.ID)
		}
	}

	// Delete the tasks that are not listed in the ToDo item, i.e. all tasks
	// that haven't been added to the list of valid tasks before.
	sql, args, _ := squirrel.
		Delete("tasks").
		Where(squirrel.And{
			squirrel.Eq{"todo_id": id},
			squirrel.NotEq{"id": taskIDs},
		}).
		ToSql()

	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	// Existing tasks have been updated and removed tasks have been deleted at
	// this point. Finally, insert all new tasks.
	insert := squirrel.
		Insert("tasks").
		Columns("name", "description", "todo_id")

	for _, task := range toDo.Tasks {
		if task.ID == 0 {
			insert.Values(task.Name, task.Description, id)
		}
	}

	sql, args, _ = insert.ToSql()
	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	// All tasks are done - update the ToDo item itself.
	sql, args, _ = squirrel.
		Update("todos").
		Set("name", toDo.Name).
		Set("description", toDo.Description).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	_, err := m.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteToDo deletes the ToDo item with the given ID. If the ToDo item cannot
// be found, core.ErrToDoNotFound will be returned.
func (m *mariaDB) DeleteToDo(id int64) error {
	sql, args, _ := squirrel.
		Delete("tasks").
		Where(squirrel.Eq{"todo_id": id}).
		ToSql()

	_, err := m.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	sql, args, _ = squirrel.
		Delete("todos").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	_, err = m.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

// createTaskForToDo inserts a task that references the given ToDo ID.
func (m *mariaDB) createTaskForToDo(toDoId int64, task model.Task) (model.Task, error) {
	sql, args, _ := squirrel.
		Insert("tasks").
		Columns("name", "description", "todo_id").
		Values(task.Name, task.Description, toDoId).
		ToSql()

	result, err := m.db.Exec(sql, args...)
	if err != nil {
		return model.Task{}, err
	}

	id, _ := result.LastInsertId()
	task.ID = id

	return task, nil
}

// findTasksByToDoID returns all tasks that reference the given ToDo ID.
func (m *mariaDB) findTasksByToDoID(toDoID int64) ([]model.Task, error) {
	sql, args, _ := squirrel.
		Select("id", "name", "description").
		From("todos").
		Where(squirrel.Eq{"todo_id": toDoID}).
		ToSql()

	rows, err := m.db.Queryx(sql, args...)
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
