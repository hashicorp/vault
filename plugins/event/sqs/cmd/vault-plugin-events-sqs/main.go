// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"log"
	"os"

	"github.com/hashicorp/vault/plugins/event/sqs"
	"github.com/hashicorp/vault/sdk/event"
)

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Run instantiates a MSSQL object, and runs the RPC server for the plugin
func Run() error {
	event.Serve(sqs.New())
	return nil
}
