package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/plugins/database/cassandra"
)

func main() {
	err := cassandra.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
