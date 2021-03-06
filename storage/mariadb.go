// Package storage provides a generic storage interface and its implementations.
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
	DBName   string
}

// URI yields a connection string in the form <user>:<password>@<host>:<port>/.
func (m MariaDBConfig) URI() string {
	return fmt.Sprintf("%s:%s@(%s)/", m.User, m.Password, m.Address)
}

type mariaDB struct {
	config        MariaDBConfig
	db            *sqlx.DB
	isInitialized bool
}

// NewMariaDB creates a new MariaDB connection using the given configuration.
func NewMariaDB(config MariaDBConfig) (*mariaDB, error) {
	mariaDB := &mariaDB{
		config: config,
	}

	if err := mariaDB.connect(); err != nil {
		return nil, err
	}

	return mariaDB, nil
}

// connect tries to establish a connection to the configured MariaDB host. If
// MariaDB has already been initialized using Initialize, connect will try to
// directly connect to the database.
func (m *mariaDB) connect() error {
	uri := m.config.URI()
	if m.isInitialized {
		uri = uri + m.config.DBName
	}

	db, err := sqlx.Connect("mysql", uri)
	if err != nil {
		return err
	}

	m.db = db
	return nil
}

// Initialize creates the MariaDB database and tables if they don't exist yet.
func (m *mariaDB) Initialize() error {
	statements := []string{
		`CREATE DATABASE IF NOT EXISTS ` + m.config.DBName,
		`USE ` + m.config.DBName,
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

	m.isInitialized = true

	// Close the database connection and re-connect directly to the database.
	// See https://stackoverflow.com/q/19927879 for more information.
	_ = m.db.Close()
	return m.connect()
}

// CreateToDo inserts the given ToDo item, which is expected to not have an ID.
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

	for i, task := range toDo.Tasks {
		createdTask, err := m.createTaskForToDo(toDo.ID, task)
		if err != nil {
			return model.ToDo{}, err
		}
		toDo.Tasks[i] = createdTask
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

	toDos := make([]model.ToDo, 0)

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

// FindToDoByID looks for a ToDo item with the provided ID and returns that item
// if it was found. Otherwise, ErrToDoNotFound will be returned.
func (m *mariaDB) FindToDoByID(id int64) (model.ToDo, error) {
	sql, args, _ := squirrel.
		Select("id", "name", "description").
		From("todos").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	var toDo model.ToDo

	err := m.db.QueryRowx(sql, args...).StructScan(&toDo)
	if err != nil {
		return model.ToDo{}, ErrToDoNotFound
	}

	tasks, err := m.findTasksByToDoID(toDo.ID)
	if err != nil {
		return model.ToDo{}, err
	}

	toDo.Tasks = tasks

	return toDo, nil
}

// UpdateToDo overwrites a stored ToDo item with the provided ToDo instance. If
// the requested ToDo cannot be found, ErrToDoNotFound will be returned.
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
	if _, err := m.FindToDoByID(id); err != nil {
		return err
	}

	taskIDs := make([]int64, 0)

	for _, task := range toDo.Tasks {
		// If the task has an ID assigned, it is considered to be an existing
		// task that can be updated.
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

	// Delete all tasks that are not listed in the ToDo item, i.e. all tasks
	// that exist in the database but have not just been updated.
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

	insert := squirrel.
		Insert("tasks").
		Columns("name", "description", "todo_id")

	for _, task := range toDo.Tasks {
		if task.ID == 0 {
			insert = insert.Values(task.Name, task.Description, id)
		}
	}

	sql, args, _ = insert.ToSql()

	// Only run the INSERT statement if there are values to insert.
	if sql != "" {
		if _, err := m.db.Exec(sql, args...); err != nil {
			return err
		}
	}

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
// be found, ErrToDoNotFound will be returned.
func (m *mariaDB) DeleteToDo(id int64) error {
	if _, err := m.FindToDoByID(id); err != nil {
		return err
	}

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
		From("tasks").
		Where(squirrel.Eq{"todo_id": toDoID}).
		ToSql()

	rows, err := m.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		if err := rows.StructScan(&task); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Remove drops the configured MariaDB database along with its tables.
func (m *mariaDB) Remove() error {
	sql := `DROP DATABASE ` + m.config.DBName

	if _, err := m.db.Exec(sql); err != nil {
		return err
	}

	return nil
}

// Close attempts to close the database connection.
func (m *mariaDB) Close() error {
	return m.db.Close()
}
