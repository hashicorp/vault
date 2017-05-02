package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/plugins/database/postgresql"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)

	err := postgresql.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
