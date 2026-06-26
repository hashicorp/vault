// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/observability"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <json-file> <workflow-name> <matrix-id>\n", os.Args[0])
		os.Exit(1)
	}

	jsonFile := os.Args[1]
	workflowName := os.Args[2]
	matrixID := os.Args[3]

	if os.Getenv("INSTANA_AGENT_KEY") == "" {
		fmt.Fprintln(os.Stderr, "INSTANA_AGENT_KEY not set, skipping Instana upload")
		return
	}

	fmt.Printf("Processing test results from: %s\n", jsonFile)
	if err := observability.CreateSpansFromFile(jsonFile, workflowName, matrixID); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating spans: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := observability.FlushSpans(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing spans: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully uploaded test results to Instana")
}
