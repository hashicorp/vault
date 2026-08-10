// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	metrics "github.com/hashicorp/go-metrics/compat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDatabaseErrorSanitizerMiddleware(t *testing.T) {
	type testCase struct {
		inputErr    error
		secretsFunc func() map[string]string

		expectedError error
	}

	tests := map[string]testCase{
		"nil error": {
			inputErr:      nil,
			expectedError: nil,
		},
		"url error": {
			inputErr:      new(url.Error),
			expectedError: errors.New("unable to parse connection url"),
		},
		"nil secrets func": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			expectedError: errors.New("here is my password: iofsd9473tg"),
		},
		"secrets with empty string": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			secretsFunc:   secretFunc(t, "", ""),
			expectedError: errors.New("here is my password: iofsd9473tg"),
		},
		"secrets that do not match": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			secretsFunc:   secretFunc(t, "asdf", "<redacted>"),
			expectedError: errors.New("here is my password: iofsd9473tg"),
		},
		"secrets that do match": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			secretsFunc:   secretFunc(t, "iofsd9473tg", "<redacted>"),
			expectedError: errors.New("here is my password: <redacted>"),
		},
		"multiple secrets": {
			inputErr: errors.New("here is my password: iofsd9473tg"),
			secretsFunc: secretFunc(t,
				"iofsd9473tg", "<redacted>",
				"password", "<this was the word password>",
			),
			expectedError: errors.New("here is my <this was the word password>: <redacted>"),
		},
		"gRPC status error": {
			inputErr:      status.Error(codes.InvalidArgument, "an error with a password iofsd9473tg"),
			secretsFunc:   secretFunc(t, "iofsd9473tg", "<redacted>"),
			expectedError: status.Errorf(codes.InvalidArgument, "an error with a password <redacted>"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := fakeDatabase{}
			mw := NewDatabaseErrorSanitizerMiddleware(db, test.secretsFunc)

			actualErr := mw.sanitize(test.inputErr)
			if !reflect.DeepEqual(actualErr, test.expectedError) {
				t.Fatalf("Actual error: %s\nExpected error: %s", actualErr, test.expectedError)
			}
		})
	}

	t.Run("Initialize", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				initErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.Initialize(context.Background(), InitializeRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 1)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("NewUser", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				newUserErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.NewUser(context.Background(), NewUserRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 1)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				updateUserErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.UpdateUser(context.Background(), UpdateUserRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 1)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				deleteUserErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.DeleteUser(context.Background(), DeleteUserRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 1)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Type", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				typeErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.Type()
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 1)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Close", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				closeErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		err := mw.Close()
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 1)
	})
}

func secretFunc(t *testing.T, vals ...string) func() map[string]string {
	t.Helper()
	if len(vals)%2 != 0 {
		t.Fatalf("Test configuration error: secretFunc must be called with an even number of values")
	}

	m := map[string]string{}

	for i := 0; i < len(vals); i += 2 {
		key := vals[i]
		m[key] = vals[i+1]
	}

	return func() map[string]string {
		return m
	}
}

func TestTracingMiddleware(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.Initialize(context.Background(), InitializeRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 1)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("NewUser", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.NewUser(context.Background(), NewUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 1)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.UpdateUser(context.Background(), UpdateUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 1)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.DeleteUser(context.Background(), DeleteUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 1)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Type", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.Type()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 1)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Close", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		err := mw.Close()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 1)
	})
}

func TestMetricsMiddleware(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.Initialize(context.Background(), InitializeRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 1)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("NewUser", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.NewUser(context.Background(), NewUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 1)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.UpdateUser(context.Background(), UpdateUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 1)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.DeleteUser(context.Background(), DeleteUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 1)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Type", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.Type()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 1)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Close", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		err := mw.Close()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 1)
	})
}

func assertEquals(t *testing.T, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Fatalf("Actual: %d Expected: %d", actual, expected)
	}
}

// newTestMetricsSink installs a fresh in-memory sink as the global metrics sink
// and returns it, so a test can inspect exactly what the middleware emitted.
//
// The global sink is process-wide, so tests using this must not run in
// parallel with each other.
func newTestMetricsSink(t *testing.T) *metrics.InmemSink {
	t.Helper()

	sink := metrics.NewInmemSink(time.Hour, 2*time.Hour)

	// An empty service name keeps the emitted keys unprefixed, so assertions
	// can name the metric exactly as the middleware produces it.
	config := metrics.DefaultConfig("")
	config.EnableHostname = false
	config.EnableTypePrefix = false
	config.EnableRuntimeMetrics = false

	if _, err := metrics.NewGlobal(config, sink); err != nil {
		t.Fatalf("failed to install test metrics sink: %s", err)
	}
	return sink
}

// labelsForMetric returns the labels attached to the named metric. It reports
// whether the metric was emitted at all.
func labelsForMetric(sink *metrics.InmemSink, name string) ([]metrics.Label, bool) {
	for _, interval := range sink.Data() {
		interval.RLock()
		for _, counter := range interval.Counters {
			if counter.Name == name {
				labels := counter.Labels
				interval.RUnlock()
				return labels, true
			}
		}
		for _, sample := range interval.Samples {
			if sample.Name == name {
				labels := sample.Labels
				interval.RUnlock()
				return labels, true
			}
		}
		interval.RUnlock()
	}
	return nil, false
}

func TestMetricsMiddleware_OperationLabels(t *testing.T) {
	type testCase struct {
		namespace      string
		mountPoint     string
		connectionName string
		expectedLabels []metrics.Label
	}

	tests := map[string]testCase{
		"no metadata emits unlabeled metrics": {
			expectedLabels: nil,
		},
		"namespace, mount point and connection name are all emitted in order": {
			namespace:      "team-a",
			mountPoint:     "database/",
			connectionName: "prod-primary",
			expectedLabels: []metrics.Label{
				{Name: "namespace", Value: "team-a"},
				{Name: "mount_point", Value: "database/"},
				{Name: "connection_name", Value: "prod-primary"},
			},
		},
		"mount point and connection name are both emitted": {
			mountPoint:     "database/",
			connectionName: "prod-primary",
			expectedLabels: []metrics.Label{
				{Name: "mount_point", Value: "database/"},
				{Name: "connection_name", Value: "prod-primary"},
			},
		},
		"mount point alone is emitted": {
			mountPoint: "database/",
			expectedLabels: []metrics.Label{
				{Name: "mount_point", Value: "database/"},
			},
		},
		"connection name alone is emitted": {
			connectionName: "prod-primary",
			expectedLabels: []metrics.Label{
				{Name: "connection_name", Value: "prod-primary"},
			},
		},
		"namespace alone is emitted": {
			namespace: "team-a",
			expectedLabels: []metrics.Label{
				{Name: "namespace", Value: "team-a"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sink := newTestMetricsSink(t)

			mw := databaseMetricsMiddleware{
				next:           &recordingDatabase{},
				typeStr:        "pgx",
				namespace:      test.namespace,
				mountPoint:     test.mountPoint,
				connectionName: test.connectionName,
			}

			if _, err := mw.NewUser(context.Background(), NewUserRequest{}); err != nil {
				t.Fatalf("Expected no error, but got: %s", err)
			}

			// Both the aggregate metric and the per-plugin-type metric must
			// carry the labels.
			for _, metricName := range []string{"database.NewUser", "database.pgx.NewUser"} {
				labels, found := labelsForMetric(sink, metricName)
				if !found {
					t.Fatalf("metric %q was not emitted", metricName)
				}
				if !reflect.DeepEqual(labels, test.expectedLabels) {
					t.Fatalf("metric %q labels: actual: %#v expected: %#v", metricName, labels, test.expectedLabels)
				}
			}
		})
	}
}

// TestMetricsMiddleware_MetricNamesUnchanged guards the compatibility promise
// that opting in adds labels without renaming any metric, so existing queries
// keep working.
func TestMetricsMiddleware_MetricNamesUnchanged(t *testing.T) {
	sink := newTestMetricsSink(t)

	mw := databaseMetricsMiddleware{
		next: &recordingDatabase{
			next: fakeDatabase{
				closeErr: errors.New("close failed"),
			},
		},
		typeStr:        "pgx",
		mountPoint:     "database/",
		connectionName: "prod-primary",
	}

	if err := mw.Close(); err == nil {
		t.Fatal("Expected an error from Close, but got none")
	}

	expectedNames := []string{
		"database.Close",
		"database.pgx.Close",
		"database.Close.error",
		"database.pgx.Close.error",
	}
	for _, metricName := range expectedNames {
		if _, found := labelsForMetric(sink, metricName); !found {
			t.Errorf("metric %q was not emitted", metricName)
		}
	}
}
