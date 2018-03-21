package main

import (
	"os"

	"github.com/hashicorp/vault/apidoc/cmd"
)

func main() {
	os.Exit(cmd.Run())
}
