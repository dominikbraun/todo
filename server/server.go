// Package server provides an HTTP server implementation serving the ToDo app.
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dominikbraun/todo/controller"
	"github.com/dominikbraun/todo/core"

	"github.com/go-chi/chi"
)

// Server represents an easy-to-use HTTP server that exposes a REST API.
type Server struct {
	router     chi.Router
	internal   *http.Server
	controller *controller.RESTController
}

// New creates a server that uses the given app instance to handle requests.
func New(port uint, app *core.App) *Server {
	address := fmt.Sprintf(":%v", port)

	server := &Server{
		internal:   &http.Server{Addr: address},
		controller: controller.NewRESTController(app),
	}

	server.initializeRouter()
	server.internal.Handler = server.router

	return server
}

// Run starts the server. It will serve requests on the configured address until
// an interrupt signal has been received, e.g. by pressing Ctrl + C.
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
