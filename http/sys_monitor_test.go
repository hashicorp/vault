// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/vault"
)

func TestSysMonitorUnknownLogLevel(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: Handler,
		NumCores:    1,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	request := client.NewRequest("GET", "/v1/sys/monitor")
	request.Params.Add("log_level", "haha")
	_, err := client.RawRequest(request)

	if err == nil {
		t.Fatal("expected to get an error, but didn't")
	} else {
		if !strings.Contains(err.Error(), "Code: 400") {
			t.Fatalf("expected to receive a 400 error, but got %s instead", err)
		}

		if !strings.Contains(err.Error(), "unknown log level") {
			t.Fatalf("expected to receive a message indicating an unknown log level, but got %s instead", err)
		}
	}
}

func TestSysMonitorUnknownLogFormat(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: Handler,
		NumCores:    1,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	request := client.NewRequest("GET", "/v1/sys/monitor")
	request.Params.Add("log_format", "haha")
	_, err := client.RawRequest(request)

	if err == nil {
		t.Fatal("expected to get an error, but didn't")
	} else {
		if !strings.Contains(err.Error(), "Code: 400") {
			t.Fatalf("expected to receive a 400 error, but got %s instead", err)
		}

		if !strings.Contains(err.Error(), "unknown log format") {
			t.Fatalf("expected to receive a message indicating an unknown log format, but got %s instead", err)
		}
	}
}

func TestSysMonitorStreamingLogs(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: Handler,
		NumCores:    1,
	})
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	stopCh := testhelpers.GenerateDebugLogs(t, client)
	defer close(stopCh)

	for _, lf := range []string{"standard", "json"} {
		t.Run(lf, func(t *testing.T) {
			debugCount := 0

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logCh, err := client.Sys().Monitor(ctx, "DEBUG", lf)
			if err != nil {
				t.Fatal(err)
			}

			type jsonlog struct {
				Level     string `json:"@level"`
				Message   string `json:"@message"`
				TimeStamp string `json:"@timestamp"`
			}
			jsonLog := &jsonlog{}

			timeCh := time.After(5 * time.Second)

			for {
				select {
				case log := <-logCh:
					if lf == "json" {
						err := json.Unmarshal([]byte(log), jsonLog)
						if err != nil {
							t.Fatal("Expected JSON log from channel")
						}
						if strings.Contains(jsonLog.Level, "debug") {
							debugCount++
						}
					} else if strings.Contains(log, "[DEBUG]") {
						debugCount++
					}
					if debugCount > 3 {
						// If we've seen multiple lines that match what we want,
						// it's probably safe to assume streaming is working
						return
					}
				case <-timeCh:
					t.Fatal("Failed to get a DEBUG message after 5 seconds")
				}
			}
		})
	}
}
