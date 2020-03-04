package http

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/vault"
)

func TestSysMonitorUnknownLogLevel(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{HandlerFunc:Handler})
	cluster.Start()
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

func TestSysMonitorStreamingLogs(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{HandlerFunc:Handler})
	cluster.Start()
	defer cluster.Cleanup()

	stopCh := make(chan struct{})
	success := make(chan struct{})
	client := cluster.Cores[0].Client

	// Make requests that generate logs
	testhelpers.GenerateDebugLogs(t, stopCh, client)

	// Watch for logs that match what we want
	go func() {
		debugCount := 0
		client := cluster.Cores[0].Client
		logCh, err := client.Sys().Monitor("DEBUG", stopCh)

		if err != nil {
			t.Fatal(err)
		}

		for {
			select {
			case log := <-logCh:
				if strings.Contains(log, "[DEBUG]") {
					debugCount++
				}
			case <-stopCh:
				return
			}

			// If we've seen multiple lines that match what we want,
			// it's probably safe to assume streaming is working
			if debugCount > 3 {
				close(success)
				return
			}
		}
	}()

	select {
	case <-success:
		close(stopCh)
		return
	case <-time.After(5 * time.Second):
		close(stopCh)
		t.Fatal("Failed to get a DEBUG message after 5 seconds")
	}
}
