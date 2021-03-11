// Package controller provides application controllers that convert incoming
// data to domain models, run business logic on them and return the results.
package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dominikbraun/todo/model"
	"github.com/dominikbraun/todo/todo"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// errorResponse will be returned by the REST API in case an error occurred.
type errorResponse struct {
	Error error `json:"error"`
}

// RESTController represents a controller capable of handling incoming HTTP
// requests and yielding a corresponding JSON result.
type RESTController struct {
	app *todo.App
}

// NewRESTController returns a new REST controller instance that will use the
// provided app instance for triggering business logic.
func NewRESTController(app *todo.App) *RESTController {
	return &RESTController{
		app: app,
	}
}

// CreateToDo processes a POST request for creating a ToDo item. It expects
// a ToDo item without ID and returns an item containing the ID.
func (r *RESTController) CreateToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var toDo model.ToDo

		if err := json.NewDecoder(request.Body).Decode(&toDo); err != nil {
			respond(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		createdToDo, err := r.app.CreateToDo(toDo)
		if err != nil {
			respond(writer, request, http.StatusInternalServerError, err)
			return
		}

		respond(writer, request, http.StatusOK, createdToDo)
	}
}

// GetToDos processes a GET request for listing all ToDo items.
func (r *RESTController) GetToDos() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		toDos, err := r.app.GetToDos()
		if err != nil {
			respond(writer, request, http.StatusInternalServerError, err)
			return
		}

		respond(writer, request, http.StatusOK, toDos)
	}
}

// GetToDo processes a GET request for retrieving a single ToDo item by ID.
//
// Expects the `id` URL parameter.
func (r *RESTController) GetToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respond(writer, request, http.StatusBadRequest, err)
			return
		}

		toDo, err := r.app.GetToDo(int64(id))
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, todo.ErrToDoNotFound) {
				status = http.StatusNotFound
			}
			respond(writer, request, status, err)
			return
		}

		respond(writer, request, http.StatusOK, toDo)
	}
}

// UpdateToDo processes a PUT request for updating a ToDo item. The item
// with the given ID will be overridden by the item in the request body.
//
// Expects the `id` URL parameter.
func (r *RESTController) UpdateToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respond(writer, request, http.StatusBadRequest, err)
			return
		}

		var toDo model.ToDo

		if err := json.NewDecoder(request.Body).Decode(&toDo); err != nil {
			respond(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		updatedToDo, err := r.app.UpdateToDo(int64(id), toDo)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, todo.ErrToDoNotFound) {
				status = http.StatusNotFound
			}
			respond(writer, request, status, err)
			return
		}

		respond(writer, request, http.StatusOK, updatedToDo)
	}
}

// DeleteToDo processes a DELETE request for deleting a single ToDo item.
// Deleting an item will also delete all of its sub-tasks.
//
// Expects the `id` URL parameter.
func (r *RESTController) DeleteToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respond(writer, request, http.StatusBadRequest, err)
			return
		}

		if err := r.app.DeleteToDo(int64(id)); err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, todo.ErrToDoNotFound) {
				status = http.StatusNotFound
			}
			respond(writer, request, status, err)
			return
		}

		respond(writer, request, http.StatusOK, nil)
	}
}

// respond renders an HTTP response with the given status code and content.
//
// respond won't dispatch the response. In case the content is an error value,
// it will automatically be wrapped in an errorResponse instance.
func respond(writer http.ResponseWriter, request *http.Request, status int, v interface{}) {
	content := v

	if err, isError := v.(error); isError {
		content = errorResponse{Error: err}
	}

	render.Status(request, status)
	render.JSON(writer, request, content)
}
