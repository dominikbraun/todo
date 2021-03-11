package storage

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/dominikbraun/todo/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	createToDosTable = `
CREATE TABLE todos (
	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	description VARCHAR(500)
)`
	createTasksTable = `
CREATE TABLE tasks (
	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	description VARCHAR(500),
	todo_id BIGINT UNSIGNED NOT NULL
)`
)

type mariaDB struct {
	db *sqlx.DB
}

func NewMariaDB(user, password, address, database string) (*mariaDB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s", user, password, address, database)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &mariaDB{
		db: db,
	}, nil
}

func (m *mariaDB) Install() error {
	if _, err := m.db.Exec(createToDosTable); err != nil {
		return err
	}

	if _, err := m.db.Exec(createTasksTable); err != nil {
		return err
	}

	return nil
}

func (m *mariaDB) CreateToDo(toDo model.ToDo) (model.ToDo, error) {
	sql, args, _ := sq.
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
	sql, _, _ := sq.
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
	sql, args, _ := sq.
		Select("id", "name", "description").
		From("todos").
		Where(sq.Eq{"id": id}).
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
	sql, args, _ := sq.
		Delete("tasks").
		Where(sq.Eq{"todo_id": id}).
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

	sql, args, _ = sq.
		Update("todos").
		Set("name", toDo.Name).
		Set("description", toDo.Description).
		Where(sq.Eq{"id": id}).
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
	sql, args, _ := sq.
		Delete("tasks").
		Where(sq.Eq{"todo_id": id}).
		ToSql()

	_, err := m.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	sql, args, _ = sq.
		Delete("todos").
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err = m.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *mariaDB) createTaskForToDo(toDoId int64, task model.Task) (model.Task, error) {
	sql, args, _ := sq.
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
	sql, args, _ := sq.
		Select("id", "name", "description").
		From("todos").
		Where(sq.Eq{"todo_id": toDoID}).
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
