// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main // import "github.com/hashicorp/vault"

import (
	command_server "github.com/hashicorp/vault/command-server"
	"github.com/hashicorp/vault/internal"
	"os"
)

func init() {
	// this is a good place to patch SHA-1 support back into x509
	internal.PatchSha1()
}

func main() {
	os.Exit(command_server.Run(os.Args[1:]))
}
