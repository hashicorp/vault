package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/plugins/database/mssql"
)

func main() {
	err := mssql.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
