// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/observability"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <stage-name> <duration-seconds> <job-name>\n", os.Args[0])
		os.Exit(1)
	}

	stageName := os.Args[1]
	duration := os.Args[2]
	jobName := os.Args[3]

	if stageName == "" {
		fmt.Fprintln(os.Stderr, "stage-name must not be empty")
		os.Exit(1)
	}
	if jobName == "" {
		fmt.Fprintln(os.Stderr, "job-name must not be empty")
		os.Exit(1)
	}
	if _, err := strconv.ParseFloat(duration, 64); err != nil {
		fmt.Fprintf(os.Stderr, "duration-seconds must be a number: %v\n", err)
		os.Exit(1)
	}

	if os.Getenv("INSTANA_AGENT_KEY") == "" {
		fmt.Fprintln(os.Stderr, "INSTANA_AGENT_KEY not set, skipping Instana upload")
		return
	}

	fmt.Printf("Sending build metric for stage: %s\n", stageName)
	if err := observability.SendBuildMetric(stageName, duration, jobName); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending build metric: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := observability.FlushSpans(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing spans: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully uploaded build metric to Instana")
}
