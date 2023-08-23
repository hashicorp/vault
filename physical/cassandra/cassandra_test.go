// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"os"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestCassandraBackend(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping in short mode")
	}
	if os.Getenv("VAULT_CI_GO_TEST_RACE") != "" {
		t.Skip("skipping race test in CI pending https://github.com/gocql/gocql/pull/1474")
	}

	host, cleanup := cassandra.PrepareTestContainer(t)
	defer cleanup()

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)
	b, err := NewCassandraBackend(map[string]string{
		"hosts":                       host.ConnectionURL(),
		"protocol_version":            "3",
		"connection_timeout":          "5",
		"initial_connection_timeout":  "5",
		"simple_retry_policy_retries": "3",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestCassandraBackendBuckets(t *testing.T) {
	expectations := map[string][]string{
		"":          {"."},
		"a":         {"."},
		"a/b":       {".", "a"},
		"a/b/c/d/e": {".", "a", "a/b", "a/b/c", "a/b/c/d"},
	}

	b := &CassandraBackend{}
	for input, expected := range expectations {
		actual := b.buckets(input)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("bad: %v expected: %v", actual, expected)
		}
	}
}
