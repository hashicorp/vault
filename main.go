// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main // import "github.com/hashicorp/vault"

import (
	"os"

	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/internal"
)

func init() {
	// this is a good place to patch SHA-1 support back into x509
	internal.PatchSha1()
}

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
