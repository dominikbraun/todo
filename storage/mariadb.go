package storage

import (
	"fmt"

	"github.com/dominikbraun/todo/model"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MariaDBConfig struct {
	User     string
	Password string
	Address  string
}

func (m MariaDBConfig) URI() string {
	return fmt.Sprintf("%s:%s@(%s)/", m.User, m.Password, m.Address)
}

type mariaDB struct {
	config MariaDBConfig
	db     *sqlx.DB
}

func NewMariaDB(config MariaDBConfig) (*mariaDB, error) {
	db, err := sqlx.Open("mysql", config.URI())
	if err != nil {
		return nil, err
	}

	return &mariaDB{
		config: config,
		db:     db,
	}, nil
}

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

		tasks, err := m.findTasksByToDo(toDo.ID)
		if err != nil {
			return nil, err
		}

		toDo.Tasks = tasks
		toDos = append(toDos, toDo)
	}

	return toDos, nil
}

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

	tasks, err := m.findTasksByToDo(toDo.ID)
	if err != nil {
		return model.ToDo{}, err
	}

	toDo.Tasks = tasks

	return toDo, nil
}

func (m *mariaDB) UpdateToDo(id int64, toDo model.ToDo) error {
	sql, args, _ := squirrel.
		Delete("tasks").
		Where(squirrel.Eq{"todo_id": id}).
		ToSql()

	_, err := m.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	var newTasks []model.Task

	for _, t := range toDo.Tasks {
		task, err := m.createTaskForToDo(id, t)
		if err != nil {
			return err
		}
		newTasks = append(newTasks, task)
	}

	sql, args, _ = squirrel.
		Update("todos").
		Set("name", toDo.Name).
		Set("description", toDo.Description).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	_, err = m.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	toDo.ID = id
	toDo.Tasks = newTasks

	return nil
}

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

func (m *mariaDB) findTasksByToDo(toDoID int64) ([]model.Task, error) {
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
