package main

import (
	"log"
	"os"

	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Run instantiates a MySQL object, and runs the RPC server for the plugin
func Run() error {
	var f func() (interface{}, error)
	f = mysql.New(mysql.DefaultUserNameTemplate)

	dbplugin.ServeMultiplex(f)

	return nil
}
