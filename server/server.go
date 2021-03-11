package server

import (
	"context"
	"errors"
	"github.com/dominikbraun/todo/todo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dominikbraun/todo/controller"

	"github.com/go-chi/chi"
)

type Server struct {
	router     chi.Router
	internal   *http.Server
	controller *controller.REST
}

func New() *Server {
	server := &Server{
		internal:   &http.Server{},
		controller: controller.NewRESTController(todo.New()),
	}

	server.initializeRouter()
	server.internal.Handler = server.router

	return server
}

func (s *Server) Run() error {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt)

	go func() {
		err := s.internal.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.internal.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}