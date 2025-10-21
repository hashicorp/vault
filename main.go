// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package main // import "github.com/hashicorp/vault"

import (
	"os"

	"github.com/hashicorp/vault/command"
)

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
