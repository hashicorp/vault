package main

import (
	"log"
	"os"
	
	"github.com/hashicorp/vault/api"
	"github.com/fhitchen/vault/plugins/database/couchbase"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	err := couchbase.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
