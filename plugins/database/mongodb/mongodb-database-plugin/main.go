package main

import (
	"log"
	"os"

	"github.com/hashicorp/vault/plugins/database/mongodb"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Run instantiates a MongoDB object, and runs the RPC server for the plugin
func Run() error {
	dbType, err := mongodb.New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database))

	return nil
}
