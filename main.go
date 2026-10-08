// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:debug cryptocustomrand=1

package main // import "github.com/hashicorp/vault"

import (
	"os"

	command "github.com/hashicorp/vault/command/server"
)

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
