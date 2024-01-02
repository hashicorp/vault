// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"log"
	"os"

	"github.com/hashicorp/vault/plugins/database/redshift"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

func main() {
	if err := Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Run instantiates a RedShift object, and runs the RPC server for the plugin
func Run() error {
	dbplugin.ServeMultiplex(redshift.New)

	return nil
}
