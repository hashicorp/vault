package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/plugins/database/postgresql"
)

func main() {
	err := postgresql.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
