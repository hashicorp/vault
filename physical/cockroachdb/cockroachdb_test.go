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
)

type Config struct {
	docker.ServiceURL
	TableName   string
	HATableName string
}

var _ docker.ServiceConfig = &Config{}

func prepareCockroachDBTestContainer(t *testing.T) (func(), *Config) {
	if retURL := os.Getenv("CR_URL"); retURL != "" {
		s, err := docker.NewServiceURLParse(retURL)
		if err != nil {
			t.Fatal(err)
		}
		return func() {}, &Config{
			ServiceURL:  *s,
			TableName:   "vault." + defaultTableName,
			HATableName: "vault." + defaultHATableName,
		}
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

	db, err := sql.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	database := "vault"
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", database))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServiceURL:  *docker.NewServiceURL(u),
		TableName:   database + "." + defaultTableName,
		HATableName: database + "." + defaultHATableName,
	}, nil
}

func TestCockroachDBBackend(t *testing.T) {
	cleanup, config := prepareCockroachDBTestContainer(t)
	defer cleanup()

	hae := os.Getenv("CR_HA_ENABLED")
	if hae == "" {
		hae = "true"
	}

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b1, err := NewCockroachDBBackend(map[string]string{
		"connection_url": config.URL().String(),
		"table":          config.TableName,
		"ha_table":       config.HATableName,
		"ha_enabled":     hae,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	b2, err := NewCockroachDBBackend(map[string]string{
		"connection_url": config.URL().String(),
		"table":          config.TableName,
		"ha_table":       config.HATableName,
		"ha_enabled":     hae,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		truncate(t, b1)
		truncate(t, b2)
	}()

	physical.ExerciseBackend(t, b1)
	truncate(t, b1)
	physical.ExerciseBackend_ListPrefix(t, b1)
	truncate(t, b1)
	physical.ExerciseTransactionalBackend(t, b1)
	truncate(t, b1)

	ha1, ok1 := b1.(physical.HABackend)
	ha2, ok2 := b2.(physical.HABackend)
	if !ok1 || !ok2 {
		t.Fatalf("CockroachDB does not implement HABackend")
	}

	if ha1.HAEnabled() && ha2.HAEnabled() {
		logger.Info("Running ha backend tests")
		physical.ExerciseHABackend(t, ha1, ha2)
	}
}

func truncate(t *testing.T, b physical.Backend) {
	crdb := b.(*CockroachDBBackend)
	_, err := crdb.client.Exec("TRUNCATE TABLE " + crdb.table)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	if crdb.haEnabled {
		_, err = crdb.client.Exec("TRUNCATE TABLE " + crdb.haTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}
}

func TestValidateDBTable(t *testing.T) {
	type testCase struct {
		table     string
		expectErr bool
	}

	tests := map[string]testCase{
		"first character is letter":     {"abcdef", false},
		"first character is underscore": {"_bcdef", false},
		"exclamation point":             {"ab!def", true},
		"at symbol":                     {"ab@def", true},
		"hash":                          {"ab#def", true},
		"percent":                       {"ab%def", true},
		"carrot":                        {"ab^def", true},
		"ampersand":                     {"ab&def", true},
		"star":                          {"ab*def", true},
		"left paren":                    {"ab(def", true},
		"right paren":                   {"ab)def", true},
		"dash":                          {"ab-def", true},
		"digit":                         {"a123ef", false},
		"dollar end":                    {"abcde$", false},
		"dollar middle":                 {"ab$def", false},
		"dollar start":                  {"$bcdef", true},
		"backtick prefix":               {"`bcdef", true},
		"backtick middle":               {"ab`def", true},
		"backtick suffix":               {"abcde`", true},
		"single quote prefix":           {"'bcdef", true},
		"single quote middle":           {"ab'def", true},
		"single quote suffix":           {"abcde'", true},
		"double quote prefix":           {`"bcdef`, true},
		"double quote middle":           {`ab"def`, true},
		"double quote suffix":           {`abcde"`, true},
		"underscore with all runes":     {"_bcd123__a__$", false},
		"all runes":                     {"abcd123__a__$", false},
		"default table name":            {defaultTableName, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateDBTable(test.table)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
		t.Run(fmt.Sprintf("database: %s", name), func(t *testing.T) {
			dbTable := fmt.Sprintf("%s.%s", test.table, test.table)
			err := validateDBTable(dbTable)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}
