// Package storage provides a generic storage interface and its implementations.
package storage

import (
	"testing"

	"github.com/dominikbraun/todo/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// newTestMariaDB creates a MariaDB using a mocked database as sql.DB.
func newTestMariaDB() (*mariaDB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	mariaDB := &mariaDB{
		db: sqlx.NewDb(db, "mysql"),
	}

	return mariaDB, mock
}

func TestMariaDB_CreateToDo(t *testing.T) {
	mariaDB, mock := newTestMariaDB()
	toDo := model.ToDo{
		Tasks: []model.Task{
			{
				Name: "Task 1",
			},
		},
	}

	defer func() {
		_ = mariaDB.db.Close()
	}()

	mock.ExpectExec("INSERT INTO todos").
		WithArgs("", "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs("Task 1", "", int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := mariaDB.CreateToDo(toDo); err != nil {
		t.Fatalf("error creating ToDo: %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestMariaDB_FindToDos(t *testing.T) {
	mariaDB, mock := newTestMariaDB()
	toDo := model.ToDo{
		ID:   1,
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
		},
	}

	defer func() {
		_ = mariaDB.db.Close()
	}()

	toDoRows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(toDo.ID, toDo.Name, toDo.Description)

	taskRows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(toDo.Tasks[0].ID, toDo.Tasks[0].Name, toDo.Tasks[0].Description)

	mock.ExpectQuery("SELECT id, name, description FROM todos").
		WithArgs().
		WillReturnRows(toDoRows)

	mock.ExpectQuery("SELECT id, name, description FROM tasks").
		WithArgs().
		WillReturnRows(taskRows)

	if _, err := mariaDB.FindToDos(); err != nil {
		t.Fatalf("error finding ToDos: %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestMariaDB_FindToDoByID(t *testing.T) {
	mariaDB, mock := newTestMariaDB()
	toDo := model.ToDo{
		ID:   1,
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
		},
	}

	defer func() {
		_ = mariaDB.db.Close()
	}()

	toDoRows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(toDo.ID, toDo.Name, toDo.Description)

	taskRows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(toDo.Tasks[0].ID, toDo.Tasks[0].Name, toDo.Tasks[0].Description)

	mock.ExpectQuery("SELECT id, name, description FROM todos WHERE id = ?").
		WithArgs(toDo.ID).
		WillReturnRows(toDoRows)

	mock.ExpectQuery("SELECT id, name, description FROM tasks WHERE todo_id = ?").
		WithArgs(toDo.ID).
		WillReturnRows(taskRows)

	if _, err := mariaDB.FindToDoByID(1); err != nil {
		t.Fatalf("error finding ToDo: %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestMariaDB_UpdateToDo(t *testing.T) {
	mariaDB, mock := newTestMariaDB()
	toDo := model.ToDo{
		ID:   1,
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
			{
				Name: "Task 2",
			},
		},
	}

	toDoRows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(toDo.ID, toDo.Name, toDo.Description)

	taskRows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(toDo.Tasks[0].ID, toDo.Tasks[0].Name, toDo.Tasks[0].Description)

	mock.ExpectQuery("SELECT id, name, description FROM todos").
		WithArgs().
		WillReturnRows(toDoRows)

	mock.ExpectQuery("SELECT id, name, description FROM tasks").
		WithArgs().
		WillReturnRows(taskRows)

	mock.ExpectExec("UPDATE tasks SET name = ?, description = ? WHERE id = ?").
		WithArgs(toDo.Tasks[0].Name, toDo.Tasks[0].Description, toDo.Tasks[0].ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM tasks").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs("Task 2", "", toDo.ID).
		WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectExec("UPDATE tasks").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := mariaDB.UpdateToDo(1, toDo); err != nil {
		t.Fatalf("error finding ToDo: %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
