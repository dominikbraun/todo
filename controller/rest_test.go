// Package controller provides application controllers that convert incoming
// requests to domain models, run business logic on them and return the results.
package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dominikbraun/todo/core"
	"github.com/dominikbraun/todo/model"
	"github.com/dominikbraun/todo/storage"

	"github.com/go-chi/chi"
)

// newTestRESTController creates a new REST controller that uses a core.App
// instance backend by an in-memory storage.
func newTestRESTController() *RESTController {
	app := core.NewApp(storage.NewMemory())

	return &RESTController{
		app: app,
	}
}

func TestRESTController_CreateToDo(t *testing.T) {
	restController := newTestRESTController()
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

	toDoBytes, _ := json.Marshal(&toDo)
	body := bytes.NewReader(toDoBytes)

	request := httptest.NewRequest("POST", "/todos", body)
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/todos", restController.CreateToDo())

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestRESTController_GetToDos(t *testing.T) {
	restController := newTestRESTController()
	toDos := []model.ToDo{
		{
			Name: "ToDo 1",
			Tasks: []model.Task{
				{
					Name: "Task 1",
				},
				{
					Name: "Task 2",
				},
			},
		},
		{
			Name: "ToDo 2",
			Tasks: []model.Task{
				{
					Name: "Task 3",
				},
				{
					Name: "Task 4",
				},
			},
		},
	}

	for _, toDo := range toDos {
		_, _ = restController.app.CreateToDo(toDo)
	}

	request := httptest.NewRequest("GET", "/todos", nil)
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/todos", restController.GetToDos())

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response []model.ToDo

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal("could not parse response body")
	}

	if len(response) != len(toDos) {
		t.Fatalf("expected %d ToDos, got %d", len(toDos), len(response))
	}
}

func TestRESTController_GetToDo(t *testing.T) {
	restController := newTestRESTController()
	toDo := model.ToDo{
		Name:        "ToDo 1",
		Description: "My ToDo",
		Tasks: []model.Task{
			{
				Name: "Task 1",
			},
			{
				Name: "Task 2",
			},
		},
	}

	createdToDo, _ := restController.app.CreateToDo(toDo)

	target := fmt.Sprintf("/todos/%d", createdToDo.ID)
	request := httptest.NewRequest("GET", target, nil)
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/todos/{id}", restController.GetToDo())

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response model.ToDo

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal("could not parse response body")
	}
}

func TestRESTController_UpdateToDo(t *testing.T) {
	restController := newTestRESTController()
	toDo := model.ToDo{
		Name: "ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
			{
				ID:   2,
				Name: "Task 2",
			},
		},
	}
	newToDo := model.ToDo{
		Name: "My ToDo 1",
		Tasks: []model.Task{
			{
				ID:   1,
				Name: "Task 1",
			},
			{
				ID:   2,
				Name: "My Task 2",
			},
			{
				Name: "Task 3",
			},
		},
	}

	newToDoBytes, _ := json.Marshal(newToDo)
	body := bytes.NewReader(newToDoBytes)

	createdToDo, _ := restController.app.CreateToDo(toDo)

	target := fmt.Sprintf("/todos/%d", createdToDo.ID)
	request := httptest.NewRequest("PUT", target, body)
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Put("/todos/{id}", restController.UpdateToDo())

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestRESTController_DeleteToDo(t *testing.T) {
	restController := newTestRESTController()
	toDo := model.ToDo{
		Name:        "ToDo 1",
		Description: "My ToDo",
		Tasks: []model.Task{
			{
				Name: "Task 1",
			},
			{
				Name: "Task 2",
			},
		},
	}

	createdToDo, _ := restController.app.CreateToDo(toDo)

	target := fmt.Sprintf("/todos/%d", createdToDo.ID)
	request := httptest.NewRequest("DELETE", target, nil)
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Delete("/todos/{id}", restController.DeleteToDo())

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}
