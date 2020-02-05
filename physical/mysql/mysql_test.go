package mysql

import (
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"

	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"

	mysqlhelper "github.com/hashicorp/vault/helper/testhelpers/mysql"
)

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
		"address":  address,
		"database": database,
		"table":    table,
		"username": username,
		"password": password,
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
		"address":    address,
		"database":   database,
		"table":      table,
		"username":   username,
		"password":   password,
		"ha_enabled": "true",
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
	cleanup, connURL := mysqlhelper.PrepareMySQLTestContainer(t, false, "secret")

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
		"address":    cfg.Addr,
		"database":   cfg.DBName,
		"table":      table,
		"username":   cfg.User,
		"password":   cfg.Passwd,
		"ha_enabled": "true",
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

	// Cancel attempt in 50 msec
	stopCh := make(chan struct{})
	time.AfterFunc(3*time.Second, func() {
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

	// Cancel attempt in 50 msec
	stopCh2 := make(chan struct{})
	time.AfterFunc(3*time.Second, func() {
		close(stopCh2)
	})
	leaderCh2, err = lock2.Lock(stopCh2)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}
