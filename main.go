package main

import (
	"log"

	"github.com/dominikbraun/todo/server"
	"github.com/dominikbraun/todo/storage"
	"github.com/dominikbraun/todo/todo"
)

func main() {
	mariaDB, err := storage.NewMariaDB("", "", "", "")
	if err != nil {
		log.Fatal(err)
	}

	app := todo.New(mariaDB)
	srv := server.New(app)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
