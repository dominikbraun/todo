// Package server provides an easy-to-use HTTP server that exposes a REST API
// and serves the ToDo application.
package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dominikbraun/todo/controller"
	"github.com/dominikbraun/todo/todo"

	"github.com/go-chi/chi"
)

// Server represents an HTTP server. It is listening on the configured address
// using a http.Server instance. Incoming HTTP requests will be routed by the
// chi.Router and processed by a controller.RESTController instance.
type Server struct {
	router     chi.Router
	internal   *http.Server
	controller *controller.RESTController
}

// New creates a new HTTP server. The given app instance will be used by the
// REST controller in order to handle requests.
func New(app *todo.App) *Server {
	server := &Server{
		internal:   &http.Server{},
		controller: controller.NewRESTController(app),
	}

	server.initializeRouter()
	server.internal.Handler = server.router

	return server
}

// Run starts the server, listening on the configured address. Can be stopped
// by sending an interrupt signal, e.g. by pressing Ctrl + C.
func (s *Server) Run() error {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt)

	go func() {
		err := s.internal.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Wait until an interrupt signal has been received.
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.internal.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
