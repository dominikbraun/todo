// Package main provides the entrypoint and functions for parsing CLI flags.
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

type flags struct {
	mariaDB    storage.MariaDBConfig
	serverPort uint
}

func main() {
	flags := parseCommandLineFlags()

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

// parseCommandLineFlags parses the application configuration from the CLI flags
// as well as environment variables. A configuration value like `port` either
// can be passed to the binary as --port flag or set as TODO_PORT variable.
func parseCommandLineFlags() flags {

	pflag.String("mariadb-user", "admin", "The MariaDB user")
	pflag.String("mariadb-password", "admin", "The MariaDB password")
	pflag.String("mariadb-address", "localhost:3306", "The MariaDB address")
	pflag.Uint("port", 8000, "The port the server should listen on")

	pflag.Parse()

	_ = viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TODO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	flags := flags{
		mariaDB: storage.MariaDBConfig{
			User:     viper.GetString("mariadb-user"),
			Password: viper.GetString("mariadb-password"),
			Address:  viper.GetString("mariadb-address"),
		},
		serverPort: viper.GetUint("port"),
	}

	return flags
}
