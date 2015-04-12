package main

import (
	"os"

	"github.com/hashicorp/vault/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
