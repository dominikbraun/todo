// Package storage provides a generic storage interface and its implementations.
package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/dominikbraun/todo/model"

	"github.com/google/go-cmp/cmp"
)

const (
	envTestMariaDB         = "TODO_TEST_MARIADB"
	envTestMariaDBUser     = "TODO_TEST_MARIADB_USER"
	envTestMariaDBPassword = "TODO_TEST_MARIADB_PASSWORD"
	envTestMariaDBAddress  = "TODO_TEST_MARIADB_ADDRESS"
	envTestMariaDBDBName   = "TODO_TEST_MARIADB_DBNAME"
)

// loadAndInitializeStorages returns a map containing all initialized storage
// implementations that need to be tested.
//
// Whether the MariaDB implementation should be tested is determined by reading
// the TODO_TEST_MARIADB environment variable. If it is set, MariaDB is tested.
func loadAndInitializeStorages() (map[string]Storage, error) {
	storages := make(map[string]Storage)
	storages["memory"] = NewMemory()

	if os.Getenv(envTestMariaDB) != "" {
		mariaDB, err := NewMariaDB(MariaDBConfig{
			User:     os.Getenv(envTestMariaDBUser),
			Password: os.Getenv(envTestMariaDBPassword),
			Address:  os.Getenv(envTestMariaDBAddress),
			DBName:   os.Getenv(envTestMariaDBDBName),
		})
		if err != nil {
			return storages, err
		}

		storages["mariadb"] = mariaDB
	}

	for _, storage := range storages {
		if err := storage.Initialize(); err != nil {
			return nil, err
		}
	}

	return storages, nil
}

// TestStorage tests all Storage functions for all supported implementations.
//
// TestStorage it not a Unit Test but rather a stateful Integration Test that
// creates a ToDo item and simulates its entire lifecycle.
func TestStorage(t *testing.T) {
	storages, err := loadAndInitializeStorages()
	if err != nil {
		t.Fatalf("failed to initialize storages: %s", err.Error())
	}

	tests := []func(*testing.T, Storage){
		testCreateToDo,
		testFindToDos,
		testFindToDoByID,
		testUpdateToDo,
		testDeleteToDo,
	}

	for _, storage := range storages {
		for _, test := range tests {
			test(t, storage)
		}
	}

	for name, storage := range storages {
		if err := storage.Remove(); err != nil {
			t.Logf("failed to remove %s: %s", name, err.Error())
		}
		_ = storage.Close()
	}
}

func testCreateToDo(t *testing.T, storage Storage) {
	toDo := model.ToDo{
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				Name: "Task 1",
			},
			{
				Name: "Task 2",
			},
		},
	}

	createdToDo, err := storage.CreateToDo(toDo)
	if err != nil {
		t.Fatal(err)
	}

	toDo.ID = createdToDo.ID

	if !cmp.Equal(createdToDo, toDo) {
		t.Fatalf("expected ToDo %v, got %v", toDo, createdToDo)
	}
}

func testFindToDos(t *testing.T, storage Storage) {
	toDos, err := storage.FindToDos()
	if err != nil {
		t.Fatal(err)
	}

	if len(toDos) != 1 {
		t.Fatalf("expected %d ToDos, got %d", 1, len(toDos))
	}

	if len(toDos[0].Tasks) != 2 {
		t.Errorf("expected %d tasks, got %d", 2, len(toDos[0].Tasks))
	}
}

func testFindToDoByID(t *testing.T, storage Storage) {
	toDo, err := storage.FindToDoByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if toDo.ID != 1 {
		t.Errorf("expected ID %d, got %d", 1, toDo.ID)
	}

	if len(toDo.Tasks) != 2 {
		t.Errorf("expected %d tasks, got %d", 2, len(toDo.Tasks))
	}
}

func testUpdateToDo(t *testing.T, storage Storage) {
	toDo := model.ToDo{
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
			{
				Name: "New Task",
			},
		},
	}

	err := storage.UpdateToDo(1, toDo)
	if err != nil {
		t.Fatal(err)
	}

	updatedToDo, err := storage.FindToDoByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(updatedToDo.Tasks) != len(toDo.Tasks) {
		t.Fatalf("expected %d tasks, got %d", len(toDo.Tasks), len(updatedToDo.Tasks))
	}
}

func testDeleteToDo(t *testing.T, storage Storage) {
	if err := storage.DeleteToDo(1); err != nil {
		t.Fatal(err)
	}

	if err := storage.DeleteToDo(1); !errors.Is(err, ErrToDoNotFound) {
		t.Fatalf("expected error %v, got %v", ErrToDoNotFound, err)
	}
}
