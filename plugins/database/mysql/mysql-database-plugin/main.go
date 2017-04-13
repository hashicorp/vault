package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/plugins/database/mysql"
)

func main() {
	err := mysql.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
