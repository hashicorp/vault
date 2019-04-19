package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/plugins/database/mongodb"
	"github.com/y0ssar1an/q"
)

func main() {
	q.Q(">>>trying to run")
	fmt.Println(">::: trying to run")
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	err := mongodb.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
