package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"

	_ "github.com/lib/pq"
)

type Config struct {
	docker.ServiceURL
	TableName string
}

var _ docker.ServiceConfig = &Config{}

func prepareCockroachDBTestContainer(t *testing.T) (func(), *Config) {
	tableName := os.Getenv("CR_TABLE")
	if tableName == "" {
		tableName = "vault_kv_store"
	}
	if retURL := os.Getenv("CR_URL"); retURL != "" {
		s, err := docker.NewServiceURLParse(retURL)
		if err != nil {
			t.Fatal(err)
		}
		return func() {}, &Config{*s, tableName}
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "cockroachdb/cockroach",
		ImageTag:      "release-1.0",
		ContainerName: "cockroachdb",
		Cmd:           []string{"start", "--insecure"},
		Ports:         []string{"26257/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker CockroachDB: %s", err)
	}
	svc, err := runner.StartService(context.Background(), connectCockroachDB)
	if err != nil {
		t.Fatalf("Could not start docker CockroachDB: %s", err)
	}

	return svc.Cleanup, svc.Config.(*Config)
}

func connectCockroachDB(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword("root", ""),
		Host:     fmt.Sprintf("%s:%d", host, port),
		RawQuery: "sslmode=disable",
	}

	db, err := sql.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	database := "database"
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", database))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServiceURL: *docker.NewServiceURL(u),
		TableName:  database + ".vault_kv",
	}, nil
}

func TestCockroachDBBackend(t *testing.T) {
	cleanup, config := prepareCockroachDBTestContainer(t)
	defer cleanup()

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewCockroachDBBackend(map[string]string{
		"connection_url": config.URL().String(),
		"table":          config.TableName,
	}, logger)

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		truncate(t, b)
	}()

	physical.ExerciseBackend(t, b)
	truncate(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
	truncate(t, b)
	physical.ExerciseTransactionalBackend(t, b)
}

func truncate(t *testing.T, b physical.Backend) {
	crdb := b.(*CockroachDBBackend)
	_, err := crdb.client.Exec("TRUNCATE TABLE " + crdb.table)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
}
