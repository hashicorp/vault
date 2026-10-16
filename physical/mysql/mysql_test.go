// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package mysql

import (
	"bytes"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	mysqlhelper "github.com/hashicorp/vault/helper/testhelpers/mysql"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/stretchr/testify/require"
)

func TestMySQLPlaintextCatch(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MYSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	// Run vault tests
	var buf bytes.Buffer
	log.DefaultOutput = &buf

	logger := logging.NewVaultLogger(log.Debug)

	NewMySQLBackend(map[string]string{
		"address":                      address,
		"database":                     database,
		"table":                        table,
		"username":                     username,
		"password":                     password,
		"plaintext_connection_allowed": "false",
	}, logger)

	str := buf.String()
	dataIdx := strings.IndexByte(str, ' ')
	rest := str[dataIdx+1:]

	if !strings.Contains(rest, "credentials will be sent in plaintext") {
		t.Fatalf("No warning of plaintext credentials occurred")
	}
}

func TestMySQLBackend(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MYSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewMySQLBackend(map[string]string{
		"address":                      address,
		"database":                     database,
		"table":                        table,
		"username":                     username,
		"password":                     password,
		"plaintext_connection_allowed": "true",
		"max_connection_lifetime":      "1",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mysql := b.(*MySQLBackend)
		_, err := mysql.client.Exec("DROP TABLE IF EXISTS " + mysql.dbTable + " ," + mysql.dbLockTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestMySQLHABackend(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MYSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)
	config := map[string]string{
		"address":                      address,
		"database":                     database,
		"table":                        table,
		"username":                     username,
		"password":                     password,
		"ha_enabled":                   "true",
		"plaintext_connection_allowed": "true",
	}

	b, err := NewMySQLBackend(config, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mysql := b.(*MySQLBackend)
		_, err := mysql.client.Exec("DROP TABLE IF EXISTS " + mysql.dbTable + " ," + mysql.dbLockTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	b2, err := NewMySQLBackend(config, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))
}

// TestMySQLHABackend_LockFailPanic is a regression test for the panic shown in
// https://github.com/hashicorp/vault/issues/8203 and patched in
// https://github.com/hashicorp/vault/pull/8229
func TestMySQLHABackend_LockFailPanic(t *testing.T) {
	cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")

	cfg, err := mysql.ParseDSN(connURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := mysqlhelper.TestCredsExist(t, connURL, cfg.User, cfg.Passwd); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	table := "test"
	logger := logging.NewVaultLogger(log.Debug)
	config := map[string]string{
		"address":                      cfg.Addr,
		"database":                     cfg.DBName,
		"table":                        table,
		"username":                     cfg.User,
		"password":                     cfg.Passwd,
		"ha_enabled":                   "true",
		"plaintext_connection_allowed": "true",
	}

	b, err := NewMySQLBackend(config, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	b2, err := NewMySQLBackend(config, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	b1ha := b.(physical.HABackend)
	b2ha := b2.(physical.HABackend)

	// Copied from ExerciseHABackend - ensuring things are normal at this point
	// Get the lock
	lock, err := b1ha.LockWith("foo", "bar")
	if err != nil {
		t.Fatalf("initial lock: %v", err)
	}

	// Attempt to lock
	leaderCh, err := lock.Lock(nil)
	if err != nil {
		t.Fatalf("lock attempt 1: %v", err)
	}
	if leaderCh == nil {
		t.Fatalf("missing leaderCh")
	}

	// Check the value
	held, val, err := lock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Errorf("should be held")
	}
	if val != "bar" {
		t.Errorf("expected value bar: %v", err)
	}

	// Second acquisition should fail
	lock2, err := b2ha.LockWith("foo", "baz")
	if err != nil {
		t.Fatalf("lock 2: %v", err)
	}

	stopCh := make(chan struct{})
	time.AfterFunc(10*time.Second, func() {
		close(stopCh)
	})

	// Attempt to lock - can't lock because lock1 is held - this is normal
	leaderCh2, err := lock2.Lock(stopCh)
	if err != nil {
		t.Fatalf("stop lock 2: %v", err)
	}
	if leaderCh2 != nil {
		t.Errorf("should not have gotten leaderCh: %v", leaderCh2)
	}
	// end normal

	// Clean up the database. When Lock() is called, a new connection is created
	// using the configuration. If that connection cannot be created, there was a
	// panic due to not returning with the connection error. Here we intentionally
	// break the config for b2, so a new connection can't be made, which would
	// trigger the panic shown in https://github.com/hashicorp/vault/issues/8203
	cleanup()

	stopCh2 := make(chan struct{})
	time.AfterFunc(10*time.Second, func() {
		close(stopCh2)
	})
	leaderCh2, err = lock2.Lock(stopCh2)
	if err == nil {
		t.Fatalf("expected error, got none, leaderCh2=%v", leaderCh2)
	}
}

// TestValidateDBTableSQLInjectionChars preserves configuration validation while
// allowing comment-like text that is harmless inside quoted identifiers.
func TestValidateDBTableSQLInjectionChars(t *testing.T) {
	t.Parallel()
	tests := map[string]bool{
		"vault`":        true,
		"vault'":        true,
		"vault`;SELECT": true,
		"vault--":       false,
		"vault_tbl1":    false,
	}
	for name, wantErr := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, names := range [][2]string{{name, "vault"}, {"vault", name}} {
				err := validateDBTable(names[0], names[1])
				if wantErr {
					require.Error(t, err, "invalid configured identifier %q must be rejected", name)
				} else {
					require.NoError(t, err, "existing valid identifier %q must remain accepted", name)
				}
			}
		})
	}
}

// TestQuoteIdentifier verifies escaping independently of validation, including
// custom lock table names that can contain spaces, quotes and backticks.
func TestQuoteIdentifier(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		name string
		want string
	}{
		"plain identifier":   {name: "vault", want: "`vault`"},
		"embedded backtick":  {name: "vault`x", want: "`vault``x`"},
		"multiple backticks": {name: "a`b`c", want: "`a``b``c`"},
		"leading/trailing":   {name: "`vault`", want: "```vault```"},
		"lock table":         {name: "vault lock's", want: "`vault lock's`"},
		"SQL-looking text":   {name: "vault`; SELECT 1--", want: "`vault``; SELECT 1--`"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, quoteIdentifier(tc.name), "identifier must be delimited and escaped")
		})
	}
}

// TestMySQLBackendCustomLockTable verifies that custom names remain supported
// with HA enabled and unused lock-table settings are ignored with HA disabled.
func TestMySQLBackendCustomLockTable(t *testing.T) {
	t.Parallel()
	cleanup, connURL := mysqlhelper.PrepareTestContainer(t, false, "secret")
	t.Cleanup(cleanup)
	cfg, err := mysql.ParseDSN(connURL)
	require.NoError(t, err, "parse test connection")
	setupDB, err := sql.Open("mysql", connURL)
	require.NoError(t, err, "open setup connection")
	t.Cleanup(func() { require.NoError(t, setupDB.Close(), "close setup connection") })

	tests := map[string]struct {
		lockTable string
		haEnabled string
	}{
		"spaces and apostrophe": {lockTable: "vault lock's", haEnabled: "true"},
		"backtick":              {lockTable: "vault`lock", haEnabled: "true"},
		"HA disabled":           {lockTable: "", haEnabled: "false"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			id, err := uuid.GenerateUUID()
			require.NoError(t, err, "generate isolated database name")
			database := "vault_quote_" + strings.ReplaceAll(id, "-", "")
			t.Cleanup(func() {
				_, err := setupDB.Exec("DROP DATABASE IF EXISTS " + quoteIdentifier(database))
				require.NoError(t, err, "drop isolated database")
			})
			b, err := NewMySQLBackend(map[string]string{
				"address":                      cfg.Addr,
				"username":                     cfg.User,
				"password":                     cfg.Passwd,
				"database":                     database,
				"table":                        "vault",
				"lock_table":                   tc.lockTable,
				"ha_enabled":                   tc.haEnabled,
				"plaintext_connection_allowed": "true",
			}, logging.NewVaultLogger(log.Debug))
			require.NoError(t, err, "create backend with custom lock-table configuration")
			backend := b.(*MySQLBackend)
			t.Cleanup(func() { require.NoError(t, backend.client.Close(), "close backend") })
			physical.ExerciseBackend(t, b)
			if tc.haEnabled == "true" {
				_, err := backend.client.Exec("INSERT INTO "+backend.dbLockTable+" VALUES (?, ?)", "job", "leader")
				require.NoError(t, err, "write to quoted lock table")
				var leader string
				err = backend.statements["get_lock"].QueryRow("job").Scan(&leader)
				require.NoError(t, err, "read from quoted lock table")
				require.Equal(t, "leader", leader, "lock table must retain stored data")
			}
		})
	}
}

func TestValidateDBTable(t *testing.T) {
	type testCase struct {
		database  string
		table     string
		expectErr bool
	}

	tests := map[string]testCase{
		"empty database & table":        {"", "", true},
		"empty database":                {"", "a", true},
		"empty table":                   {"a", "", true},
		"ascii database":                {"abcde", "a", false},
		"ascii table":                   {"a", "abcde", false},
		"ascii database & table":        {"abcde", "abcde", false},
		"only whitespace db":            {"     ", "a", true},
		"only whitespace table":         {"a", "     ", true},
		"whitespace prefix db":          {" bcde", "a", true},
		"whitespace middle db":          {"ab de", "a", true},
		"whitespace suffix db":          {"abcd ", "a", true},
		"whitespace prefix table":       {"a", " bcde", true},
		"whitespace middle table":       {"a", "ab de", true},
		"whitespace suffix table":       {"a", "abcd ", true},
		"backtick prefix db":            {"`bcde", "a", true},
		"backtick middle db":            {"ab`de", "a", true},
		"backtick suffix db":            {"abcd`", "a", true},
		"backtick prefix table":         {"a", "`bcde", true},
		"backtick middle table":         {"a", "ab`de", true},
		"backtick suffix table":         {"a", "abcd`", true},
		"single quote prefix db":        {"'bcde", "a", true},
		"single quote middle db":        {"ab'de", "a", true},
		"single quote suffix db":        {"abcd'", "a", true},
		"single quote prefix table":     {"a", "'bcde", true},
		"single quote middle table":     {"a", "ab'de", true},
		"single quote suffix table":     {"a", "abcd'", true},
		"double quote prefix db":        {`"bcde`, "a", true},
		"double quote middle db":        {`ab"de`, "a", true},
		"double quote suffix db":        {`abcd"`, "a", true},
		"double quote prefix table":     {"a", `"bcde`, true},
		"double quote middle table":     {"a", `ab"de`, true},
		"double quote suffix table":     {"a", `abcd"`, true},
		"0x0000 prefix db":              {str(0x0000, 'b', 'c'), "a", true},
		"0x0000 middle db":              {str('a', 0x0000, 'c'), "a", true},
		"0x0000 suffix db":              {str('a', 'b', 0x0000), "a", true},
		"0x0000 prefix table":           {"a", str(0x0000, 'b', 'c'), true},
		"0x0000 middle table":           {"a", str('a', 0x0000, 'c'), true},
		"0x0000 suffix table":           {"a", str('a', 'b', 0x0000), true},
		"unicode > 0xFFFF prefix db":    {str(0x10000, 'b', 'c'), "a", true},
		"unicode > 0xFFFF middle db":    {str('a', 0x10000, 'c'), "a", true},
		"unicode > 0xFFFF suffix db":    {str('a', 'b', 0x10000), "a", true},
		"unicode > 0xFFFF prefix table": {"a", str(0x10000, 'b', 'c'), true},
		"unicode > 0xFFFF middle table": {"a", str('a', 0x10000, 'c'), true},
		"unicode > 0xFFFF suffix table": {"a", str('a', 'b', 0x10000), true},
		"non-printable prefix db":       {str(0x0001, 'b', 'c'), "a", true},
		"non-printable middle db":       {str('a', 0x0001, 'c'), "a", true},
		"non-printable suffix db":       {str('a', 'b', 0x0001), "a", true},
		"non-printable prefix table":    {"a", str(0x0001, 'b', 'c'), true},
		"non-printable middle table":    {"a", str('a', 0x0001, 'c'), true},
		"non-printable suffix table":    {"a", str('a', 'b', 0x0001), true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateDBTable(test.database, test.table)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func str(r ...rune) string {
	return string(r)
}
