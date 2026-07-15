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
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	instana "github.com/instana/go-sensor"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

const (
	vaultTestServiceName  = "vault-tests"
	vaultBuildServiceName = "vault-ci-builds"
)

var (
	testSensorOnce sync.Once
	testSensor     atomic.Pointer[instana.Sensor]

	buildSensorOnce sync.Once
	buildSensor     atomic.Pointer[instana.Sensor]
)

// TestEvent represents a single test event from gotestsum JSON output.
type TestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"`
	Output  string    `json:"Output"`
}

// getTestSensor lazily initializes the Instana sensor used for test-result spans.
func getTestSensor() *instana.Sensor {
	testSensorOnce.Do(func() {
		testSensor.Store(newSensor(vaultTestServiceName))
	})
	return testSensor.Load()
}

// getBuildSensor lazily initializes the Instana sensor used for build-metric spans.
// It's kept separate from the test sensor so build stage timings don't pollute the
// test results service in the Instana dashboard.
func getBuildSensor() *instana.Sensor {
	buildSensorOnce.Do(func() {
		buildSensor.Store(newSensor(vaultBuildServiceName))
	})
	return buildSensor.Load()
}

// newSensor creates an Instana sensor for the given service name, or nil if
// INSTANA_AGENT_KEY isn't set.
//
// NOTE: the underlying go-sensor library initializes its agent connection as a
// package-level singleton the first time a sensor is created in a process, so
// only one of getTestSensor/getBuildSensor should actually be exercised per
// process. In practice this holds because test-result and build-metric uploads
// run as separate `go run` invocations (tools/instana-uploader and
// tools/build-metric-uploader respectively).
func newSensor(serviceName string) *instana.Sensor {
	agentKey := os.Getenv("INSTANA_AGENT_KEY")
	if agentKey == "" {
		fmt.Fprintln(os.Stderr, "[Instana] INSTANA_AGENT_KEY not set, skipping initialization")
		return nil
	}

	sensor := instana.NewSensor(serviceName)

	// NewSensor automatically switches to direct backend communication when
	// INSTANA_ENDPOINT_URL is set; otherwise it uses the local agent.
	if endpointURL := os.Getenv("INSTANA_ENDPOINT_URL"); endpointURL != "" {
		fmt.Fprintf(os.Stderr, "[Instana] Using serverless mode with endpoint: %s\n", endpointURL)
	} else {
		agentHost := cmp.Or(os.Getenv("INSTANA_AGENT_HOST"), "localhost")
		agentPort := cmp.Or(os.Getenv("INSTANA_AGENT_PORT"), "42699")
		fmt.Fprintf(os.Stderr, "[Instana] Using local agent mode at %s:%s\n", agentHost, agentPort)
	}

	fmt.Fprintf(os.Stderr, "[Instana] Sensor initialized with service name: %s\n", serviceName)
	return sensor
}

// setGitTags tags a span with git/CI metadata from the GitHub Actions
// environment. GITHUB_HEAD_REF is the PR source branch; falls back to
// GITHUB_REF_NAME for direct pushes. All tags are omitted outside GitHub
// Actions, where these env vars are unset.
func setGitTags(span opentracing.Span) {
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
}

func CreateSpansFromTestResults(reader io.Reader, workflowName, matrixID string) error {
	sensor := getTestSensor()
	if sensor == nil {
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
			span := sensor.Tracer().StartSpan(
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
			setGitTags(span)

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

// SendBuildMetric creates and sends a span to Instana representing the duration of a
// build stage. duration is the elapsed time in seconds.
func SendBuildMetric(stageName, duration, jobName string) error {
	sensor := getBuildSensor()
	if sensor == nil {
		fmt.Fprintln(os.Stderr, "[Instana] sensor is nil, skipping build metric")
		return nil
	}

	durationSeconds, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		return fmt.Errorf("failed to parse duration %q: %w", duration, err)
	}

	fmt.Fprintf(os.Stderr, "[Instana] Sending build metric (stage: %s, job: %s, duration: %.2fs)\n", stageName, jobName, durationSeconds)

	finishTime := time.Now()
	startTime := finishTime.Add(-time.Duration(durationSeconds * float64(time.Second)))

	// Use the stage name as the operation name so the Instana trace list and
	// service endpoint table show actual build stage names instead of a generic label.
	span := sensor.Tracer().StartSpan(
		stageName,
		opentracing.StartTime(startTime),
		opentracing.Tag{Key: "span.kind", Value: "entry"},
	)

	span.SetTag("build.stage", stageName)
	span.SetTag("build.job", jobName)
	span.SetTag("build.duration", durationSeconds)
	setGitTags(span)

	span.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: finishTime,
	})

	return nil
}

// FlushSpans flushes whichever sensor(s) were actually initialized in this
// process via getTestSensor/getBuildSensor. It reads the sensors directly
// (rather than through the getters) so it never force-initializes a sensor
// that was never actually used to create spans.
func FlushSpans(ctx context.Context) error {
	if sensor := testSensor.Load(); sensor != nil {
		if err := sensor.Flush(ctx); err != nil {
			return err
		}
	}

	if sensor := buildSensor.Load(); sensor != nil {
		if err := sensor.Flush(ctx); err != nil {
			return err
		}
	}

	return nil
}
