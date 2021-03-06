// Package controller provides application controllers that convert incoming
// requests to domain models, run business logic on them and return the results.
package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dominikbraun/todo/core"
	"github.com/dominikbraun/todo/model"
	"github.com/dominikbraun/todo/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type errorResponse struct {
	Error string `json:"error"`
}

// RESTController represents a controller capable of handling HTTP requests and
// yielding a corresponding JSON result.
type RESTController struct {
	app *core.App
}

// NewRESTController returns a new REST controller instance that will use the
// provided app instance for running business logic.
func NewRESTController(app *core.App) *RESTController {
	return &RESTController{
		app: app,
	}
}

// CreateToDo processes a POST request for creating a ToDo item. It expects a
// ToDo item without ID and returns an item containing the ID.
func (r *RESTController) CreateToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var toDo model.ToDo

		if err := json.NewDecoder(request.Body).Decode(&toDo); err != nil {
			respond(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		createdToDo, err := r.app.CreateToDo(toDo)
		if err != nil {
			respond(writer, request, statusCodeForError(err), err)
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
			respond(writer, request, statusCodeForError(err), err)
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
			respond(writer, request, statusCodeForError(err), err)
			return
		}

		toDo, err := r.app.GetToDo(int64(id))
		if err != nil {
			respond(writer, request, statusCodeForError(err), err)
			return
		}

		respond(writer, request, http.StatusOK, toDo)
	}
}

// UpdateToDo processes a PUT request for updating a ToDo item. The item with
// the given ID will be overridden by the item in the request body.
//
// Expects the `id` URL parameter.
func (r *RESTController) UpdateToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respond(writer, request, statusCodeForError(err), err)
			return
		}

		var toDo model.ToDo

		if err := json.NewDecoder(request.Body).Decode(&toDo); err != nil {
			respond(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		err = r.app.UpdateToDo(int64(id), toDo)
		if err != nil {
			respond(writer, request, statusCodeForError(err), err)
			return
		}

		respond(writer, request, http.StatusOK, nil)
	}
}

// DeleteToDo processes a DELETE request for deleting a single ToDo item. This
// will also delete all of the item's sub-tasks.
//
// Expects the `id` URL parameter.
func (r *RESTController) DeleteToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respond(writer, request, statusCodeForError(err), err)
			return
		}

		err = r.app.DeleteToDo(int64(id))
		if err != nil {
			respond(writer, request, statusCodeForError(err), err)
			return
		}

		respond(writer, request, http.StatusOK, nil)
	}
}

// respond writes the status code as well as the JSON body to an HTTP response.
//
// If v is nil, the response body will be empty. In addition, an errorResponse
// instance will be rendered automatically if v is an error value.
func respond(writer http.ResponseWriter, request *http.Request, status int, v interface{}) {
	render.Status(request, status)

	response := v

	if err, isError := v.(error); isError {
		response = errorResponse{Error: err.Error()}
	}

	if response != nil {
		render.JSON(writer, request, response)
	}
}

// statusCodeForError returns an appropriate HTTP status code for a given error.
func statusCodeForError(err error) int {
	statusCodes := map[error]int{
		storage.ErrToDoNotFound:    http.StatusNotFound,
		core.ErrNameMustNotBeEmpty: http.StatusUnprocessableEntity,
		nil:                        http.StatusOK,
	}

	statusCode, isRegistered := statusCodes[err]

	// Return status 500 for all errors that are not nil and not registered.
	if err != nil && !isRegistered {
		return http.StatusInternalServerError
	}

	return statusCode
}
