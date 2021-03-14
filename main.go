// Package main provides the application entrypoint as well as functions for
// parsing CLI flags and environment variables.
package main

import (
	"log"
	"strings"

	"github.com/dominikbraun/todo/core"
	"github.com/dominikbraun/todo/server"
	"github.com/dominikbraun/todo/storage"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// config stores all configuration values required to run the ToDo app.
type config struct {
	mariaDB    storage.MariaDBConfig
	serverPort uint
}

func main() {
	flags := parseApplicationConfig()

	mariaDB, err := storage.NewMariaDB(flags.mariaDB)
	if err != nil {
		log.Fatal(err)
	}

	if err := mariaDB.Initialize(); err != nil {
		log.Fatal(err)
	}

	app := core.NewApp(mariaDB)
	srv := server.New(flags.serverPort, app)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

// parseApplicationConfig parses the application configuration from multiple
// sources. Currently, these sources are CLI flags and environment variables.
//
// A configuration value like `port` can either be passed to the binary as a
// --port flag or specified as a TODO_PORT environment variable.
func parseApplicationConfig() config {

	pflag.String("mariadb-user", "admin", "The MariaDB user")
	pflag.String("mariadb-password", "admin", "The MariaDB password")
	pflag.String("mariadb-address", "0.0.0.0:3306", "The MariaDB address")
	pflag.String("mariadb-dbname", "todo_app", "The MariaDB database name")
	pflag.Uint("port", 8000, "The port the server should listen on")

	pflag.Parse()

	_ = viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TODO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	flags := config{
		mariaDB: storage.MariaDBConfig{
			User:     viper.GetString("mariadb-user"),
			Password: viper.GetString("mariadb-password"),
			Address:  viper.GetString("mariadb-address"),
			DBName:   viper.GetString("mariadb-dbname"),
		},
		serverPort: viper.GetUint("port"),
	}

	return flags
}
