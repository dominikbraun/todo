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

type Flags struct {
	MariaDB    storage.MariaDBConfig
	ServerPort uint
}

func main() {
	flags := parseCommandLineFlags()

	mariaDB, err := storage.NewMariaDB(flags.MariaDB)
	if err != nil {
		log.Fatal(err)
	}

	if err := mariaDB.Initialize(); err != nil {
		log.Fatal(err)
	}

	app := core.NewApp(mariaDB)
	srv := server.New(flags.ServerPort, app)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseCommandLineFlags() Flags {

	pflag.String("mariadb-user", "admin", "The MariaDB user")
	pflag.String("mariadb-password", "admin", "The MariaDB password")
	pflag.String("mariadb-address", "localhost:3306", "The MariaDB address")
	pflag.Uint("port", 8000, "The port the server should listen on")

	pflag.Parse()

	_ = viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TODO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	flags := Flags{
		MariaDB: storage.MariaDBConfig{
			User:     viper.GetString("mariadb-user"),
			Password: viper.GetString("mariadb-password"),
			Address:  viper.GetString("mariadb-address"),
		},
		ServerPort: viper.GetUint("port"),
	}

	return flags
}
