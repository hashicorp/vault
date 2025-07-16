// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build foundationdb

package foundationdb

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

func connectToFoundationDB(clusterFile string) (*fdb.Database, error) {
	if err := fdb.APIVersion(520); err != nil {
		return nil, fmt.Errorf("failed to set FDB API version: %w", err)
	}

	db, err := fdb.Open(clusterFile, []byte("DB"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &db, nil
}

func cleanupTopDir(clusterFile, topDir string) error {
	db, err := connectToFoundationDB(clusterFile)
	if err != nil {
		return fmt.Errorf("could not connect to FDB for cleanup: %w", err)
	}

	if _, err := directory.Root().Remove(db, []string{topDir}); err != nil {
		return fmt.Errorf("could not remove directory: %w", err)
	}

	return nil
}

func TestFoundationDBPathDecoration(t *testing.T) {
	cases := map[string][]byte{
		"foo":              []byte("/\x01foo"),
		"foo/":             []byte("/\x01foo/"),
		"foo/bar":          []byte("/\x02foo/\x01bar"),
		"foo/bar/":         []byte("/\x02foo/\x01bar/"),
		"foo/bar/baz":      []byte("/\x02foo/\x02bar/\x01baz"),
		"foo/bar/baz/":     []byte("/\x02foo/\x02bar/\x01baz/"),
		"foo/bar/baz/quux": []byte("/\x02foo/\x02bar/\x02baz/\x01quux"),
	}

	for path, expected := range cases {
		decorated, err := decoratePath(path)
		if err != nil {
			t.Fatalf("path %s error: %s", path, err)
		}

		if !bytes.Equal(expected, decorated) {
			t.Fatalf("path %s expected %v got %v", path, expected, decorated)
		}

		undecorated := undecoratePath(decorated)
		if undecorated != path {
			t.Fatalf("expected %s got %s", path, undecorated)
		}
	}
}

func TestFoundationDBBackend(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping in short mode")
	}

	testUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("foundationdb: could not generate UUID to top-level directory: %s", err)
	}

	topDir := fmt.Sprintf("vault-test-%s", testUUID)

	var clusterFile string
	clusterFile = os.Getenv("FOUNDATIONDB_CLUSTER_FILE")
	if clusterFile == "" {
		var cleanup func()
		cleanup, clusterFile = prepareFoundationDBTestDirectory(t, topDir)
		defer cleanup()
	}

	// Remove the test data once done
	defer func() {
		if err := cleanupTopDir(clusterFile, topDir); err != nil {
			t.Fatalf("foundationdb: could not cleanup test data at end of test: %s", err)
		}
	}()

	// Remove any leftover test data before starting
	if err := cleanupTopDir(clusterFile, topDir); err != nil {
		t.Fatalf("foundationdb: could not cleanup test data before starting test: %s", err)
	}

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)
	config := map[string]string{
		"path":         topDir,
		"api_version":  "520",
		"cluster_file": clusterFile,
	}

	b, err := NewFDBBackend(config, logger)
	if err != nil {
		t.Fatalf("foundationdb: failed to create new backend: %s", err)
	}

	b2, err := NewFDBBackend(config, logger)
	if err != nil {
		t.Fatalf("foundationdb: failed to create new backend: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
	physical.ExerciseTransactionalBackend(t, b)
	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))
}

func prepareFoundationDBTestDirectory(t *testing.T, topDir string) (func(), string) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("foundationdb: failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("foundationdb", "5.1.7", nil)
	if err != nil {
		t.Fatalf("foundationdb: could not start container: %s", err)
	}

	tmpFile, err := ioutil.TempFile("", topDir)
	if err != nil {
		t.Fatalf("foundationdb: could not create temporary file for cluster file: %s", err)
	}

	clusterFile := tmpFile.Name()

	cleanup := func() {
		var err error
		for i := 0; i < 10; i++ {
			err = pool.Purge(resource)
			if err == nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
		os.Remove(clusterFile)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	setup := func() error {
		connectString := fmt.Sprintf("foundationdb:foundationdb@127.0.0.1:%s", resource.GetPort("4500/tcp"))

		if err := tmpFile.Truncate(0); err != nil {
			return fmt.Errorf("could not truncate cluster file: %w", err)
		}

		_, err := tmpFile.WriteAt([]byte(connectString), 0)
		if err != nil {
			return fmt.Errorf("could not write cluster file: %w", err)
		}

		if _, err := connectToFoundationDB(clusterFile); err != nil {
			return fmt.Errorf("could not connect to FoundationDB after starting container: %s", err)
		}

		return nil
	}

	if pool.Retry(setup); err != nil {
		cleanup()

		t.Fatalf("foundationdb: could not setup container: %s", err)
	}

	tmpFile.Close()

	return cleanup, clusterFile
}
