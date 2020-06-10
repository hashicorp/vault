package cockroachdb

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/ory/dockertest"

	_ "github.com/lib/pq"
)

func prepareCockroachDBTestContainer(t *testing.T) (cleanup func(), retURL, tableName string) {
	tableName = os.Getenv("CR_TABLE")
	if tableName == "" {
		tableName = defaultTableName
	}
	t.Logf("Table name: %s", tableName)
	retURL = os.Getenv("CR_URL")
	if retURL != "" {
		return func() {}, retURL, tableName
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "cockroachdb/cockroach",
		Tag:        "release-1.0",
		Cmd:        []string{"start", "--insecure"},
	}
	resource, err := pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local CockroachDB docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL = fmt.Sprintf("postgresql://root@localhost:%s/?sslmode=disable", resource.GetPort("26257/tcp"))
	database := "vault"
	tableName = fmt.Sprintf("%s.%s", database, tableName)

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", database))
		return err
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return cleanup, retURL, tableName
}

func TestCockroachDBBackend(t *testing.T) {
	cleanup, connURL, table := prepareCockroachDBTestContainer(t)
	defer cleanup()

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewCockroachDBBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
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
