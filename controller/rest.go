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

type REST struct {
	app *todo.App
}

func NewRESTController(app *todo.App) *REST {
	return &REST{
		app: app,
	}
}

func (r *REST) CreateToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var toDo model.ToDo

		if err := json.NewDecoder(request.Body).Decode(&toDo); err != nil {
			respondErr(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		createdToDo, err := r.app.CreateToDo(toDo)
		if err != nil {
			respondErr(writer, request, http.StatusInternalServerError, err)
			return
		}

		respond(writer, request, http.StatusOK, createdToDo)
	}
}

func (r *REST) GetToDos() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		toDos, err := r.app.GetToDos()
		if err != nil {
			respondErr(writer, request, http.StatusInternalServerError, err)
			return
		}

		respond(writer, request, http.StatusOK, toDos)
	}
}

func (r *REST) GetToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respondErr(writer, request, http.StatusBadRequest, err)
			return
		}

		toDo, err := r.app.GetToDo(id)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, todo.ErrToDoNotFound) {
				status = http.StatusNotFound
			}
			respondErr(writer, request, status, err)
			return
		}

		respond(writer, request, http.StatusOK, toDo)
	}
}

func (r *REST) UpdateToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respondErr(writer, request, http.StatusBadRequest, err)
			return
		}

		var toDo model.ToDo

		if err := json.NewDecoder(request.Body).Decode(&toDo); err != nil {
			respondErr(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		updatedToDo, err := r.app.UpdateToDo(id, toDo)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, todo.ErrToDoNotFound) {
				status = http.StatusNotFound
			}
			respondErr(writer, request, status, err)
			return
		}

		respond(writer, request, http.StatusOK, updatedToDo)
	}
}

func (r *REST) DeleteToDo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil {
			respondErr(writer, request, http.StatusBadRequest, err)
			return
		}

		if err := r.app.DeleteToDo(id); err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, todo.ErrToDoNotFound) {
				status = http.StatusNotFound
			}
			respondErr(writer, request, status, err)
			return
		}

		respond(writer, request, http.StatusOK, nil)
	}
}

func respondErr(writer http.ResponseWriter, request *http.Request, status int, err error) {
	type body struct {
		Error error `json:"error"`
	}
	respond(writer, request, status, body{err})
}

func respond(writer http.ResponseWriter, request *http.Request, status int, v interface{}) {
	render.Status(request, status)
	render.JSON(writer, request, v)
}
