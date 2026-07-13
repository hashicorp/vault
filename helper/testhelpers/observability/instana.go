// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package observability

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	instana "github.com/instana/go-sensor"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

const vaultTestServiceName = "vault-tests"

var globalSensor *instana.Sensor

// TestEvent represents a single test event from gotestsum JSON output.
type TestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"`
	Output  string    `json:"Output"`
}

func init() {
	agentKey := os.Getenv("INSTANA_AGENT_KEY")
	if agentKey == "" {
		fmt.Fprintln(os.Stderr, "[Instana] INSTANA_AGENT_KEY not set, skipping initialization")
		return
	}

	globalSensor = instana.NewSensor(vaultTestServiceName)

	// NewSensor automatically switches to direct backend communication when
	// INSTANA_ENDPOINT_URL is set; otherwise it uses the local agent.
	if endpointURL := os.Getenv("INSTANA_ENDPOINT_URL"); endpointURL != "" {
		fmt.Fprintf(os.Stderr, "[Instana] Using serverless mode with endpoint: %s\n", endpointURL)
	} else {
		agentHost := cmp.Or(os.Getenv("INSTANA_AGENT_HOST"), "localhost")
		agentPort := cmp.Or(os.Getenv("INSTANA_AGENT_PORT"), "42699")
		fmt.Fprintf(os.Stderr, "[Instana] Using local agent mode at %s:%s\n", agentHost, agentPort)
	}

	fmt.Fprintf(os.Stderr, "[Instana] Sensor initialized with service name: %s\n", vaultTestServiceName)
}

func CreateSpansFromTestResults(reader io.Reader, workflowName, matrixID string) error {
	if globalSensor == nil {
		fmt.Fprintln(os.Stderr, "[Instana] sensor is nil, skipping span creation")
		return nil
	}

	fmt.Fprintf(os.Stderr, "[Instana] Processing test results (workflow: %s, matrix: %s)\n", workflowName, matrixID)

	decoder := json.NewDecoder(reader)
	testStarts := make(map[string]TestEvent)
	spanCount := 0

	for {
		var event TestEvent
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to decode test event: %w", err)
		}

		if event.Action == "run" && event.Test != "" {
			testStarts[event.Package+"/"+event.Test] = event
		}

		if (event.Action == "pass" || event.Action == "fail") && event.Test != "" {
			testKey := event.Package + "/" + event.Test
			startEvent, ok := testStarts[testKey]
			if !ok {
				// gotestsum doesn't always emit a "run" event (e.g. cached results),
				// so fall back to reconstructing the start time from elapsed.
				startEvent = event
				startEvent.Time = event.Time.Add(-time.Duration(event.Elapsed * float64(time.Second)))
			}

			// Use the test name as the operation name so the Instana trace list and
			// service endpoint table show actual test names instead of a generic label.
			span := globalSensor.Tracer().StartSpan(
				event.Test,
				opentracing.StartTime(startEvent.Time),
				opentracing.Tag{Key: "span.kind", Value: "entry"},
			)

			span.SetTag("test.name", event.Test)
			span.SetTag("test.package", event.Package)
			span.SetTag("test.status", event.Action)
			span.SetTag("test.duration", event.Elapsed)
			span.SetTag("workflow", workflowName)
			span.SetTag("matrix.id", matrixID)

			// GITHUB_HEAD_REF is the PR source branch; falls back to GITHUB_REF_NAME
			// for direct pushes. Empty outside GitHub Actions.
			if branch := os.Getenv("GITHUB_HEAD_REF"); branch != "" {
				span.SetTag("git.branch", branch)
			} else if branch := os.Getenv("GITHUB_REF_NAME"); branch != "" {
				span.SetTag("git.branch", branch)
			}
			if sha := os.Getenv("GITHUB_SHA"); sha != "" {
				span.SetTag("git.sha", sha)
			}
			if runID := os.Getenv("GITHUB_RUN_ID"); runID != "" {
				span.SetTag("git.run_id", runID)
				serverURL := os.Getenv("GITHUB_SERVER_URL")
				repo := os.Getenv("GITHUB_REPOSITORY")
				if serverURL != "" && repo != "" {
					span.SetTag("git.ci_url", serverURL+"/"+repo+"/actions/runs/"+runID)
				}
			}

			// ext.Error drives "Erroneous calls" in the Instana dashboard;
			// LogFields with otlog.Error separately increments "Error logs".
			if event.Action == "fail" {
				ext.Error.Set(span, true)
				span.LogFields(
					otlog.String("event", "error"),
					otlog.Error(fmt.Errorf("test %s failed", event.Test)),
				)
			}

			span.FinishWithOptions(opentracing.FinishOptions{
				FinishTime: event.Time,
			})

			delete(testStarts, testKey)
			spanCount++
		}
	}

	fmt.Fprintf(os.Stderr, "[Instana] Created %d spans\n", spanCount)
	return nil
}

func CreateSpansFromFile(filePath, workflowName, matrixID string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open test results file: %w", err)
	}
	defer file.Close()

	return CreateSpansFromTestResults(file, workflowName, matrixID)
}

func FlushSpans(ctx context.Context) error {
	if globalSensor == nil {
		return nil
	}

	return globalSensor.Flush(ctx)
}
