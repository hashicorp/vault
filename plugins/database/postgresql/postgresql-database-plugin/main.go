package main

import (
	"log"
	"os"

	"github.com/hashicorp/vault/plugins/database/postgresql"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Run instantiates a PostgreSQL object, and runs the RPC server for the plugin
func Run() error {
	dbplugin.ServeMultiplex(postgresql.New)

	return nil
}
