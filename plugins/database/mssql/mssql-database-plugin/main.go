package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/plugins/database/mssql"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)

	err := mssql.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
