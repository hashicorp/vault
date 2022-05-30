package postgresql

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/postgresql"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestPostgreSQLBackend(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	// Use docker as pg backend if no url is provided via environment variables
	connURL := os.Getenv("PGURL")
	if connURL == "" {
		cleanup, u := postgresql.PrepareTestContainer(t, "11.1")
		defer cleanup()
		connURL = u
	}

	table := os.Getenv("PGTABLE")
	if table == "" {
		table = "vault_kv_store"
	}

	hae := os.Getenv("PGHAENABLED")
	if hae == "" {
		hae = "true"
	}

	// Run vault tests
	logger.Info(fmt.Sprintf("Connection URL: %v", connURL))

	b1, err := NewPostgreSQLBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
		"ha_enabled":     hae,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	b2, err := NewPostgreSQLBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
		"ha_enabled":     hae,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}
	pg := b1.(*PostgreSQLBackend)

	// Read postgres version to test basic connects works
	var pgversion string
	if err = pg.client.QueryRow("SELECT current_setting('server_version_num')").Scan(&pgversion); err != nil {
		t.Fatalf("Failed to check for Postgres version: %v", err)
	}
	logger.Info(fmt.Sprintf("Postgres Version: %v", pgversion))

	setupDatabaseObjects(t, logger, pg)

	defer func() {
		pg := b1.(*PostgreSQLBackend)
		_, err := pg.client.Exec(fmt.Sprintf(" TRUNCATE TABLE %v ", pg.table))
		if err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}()

	logger.Info("Running basic backend tests")
	physical.ExerciseBackend(t, b1)
	logger.Info("Running list prefix backend tests")
	physical.ExerciseBackend_ListPrefix(t, b1)

	ha1, ok := b1.(physical.HABackend)
	if !ok {
		t.Fatalf("PostgreSQLDB does not implement HABackend")
	}

	ha2, ok := b2.(physical.HABackend)
	if !ok {
		t.Fatalf("PostgreSQLDB does not implement HABackend")
	}

	if ha1.HAEnabled() && ha2.HAEnabled() {
		logger.Info("Running ha backend tests")
		physical.ExerciseHABackend(t, ha1, ha2)
		testPostgresSQLLockTTL(t, ha1)
		testPostgresSQLLockRenewal(t, ha1)
	}
}

func TestPostgreSQLBackendMaxIdleConnectionsParameter(t *testing.T) {
	_, err := NewPostgreSQLBackend(map[string]string{
		"connection_url":       "some connection url",
		"max_idle_connections": "bad param",
	}, logging.NewVaultLogger(log.Debug))
	if err == nil {
		t.Error("Expected invalid max_idle_connections param to return error")
	}
	expectedErrStr := "failed parsing max_idle_connections parameter: strconv.Atoi: parsing \"bad param\": invalid syntax"
	if err.Error() != expectedErrStr {
		t.Errorf("Expected: \"%s\" but found \"%s\"", expectedErrStr, err.Error())
	}
}

func TestConnectionURL(t *testing.T) {
	type input struct {
		envar string
		conf  map[string]string
	}

	cases := map[string]struct {
		want  string
		input input
	}{
		"environment_variable_not_set_use_config_value": {
			want: "abc",
			input: input{
				envar: "",
				conf:  map[string]string{"connection_url": "abc"},
			},
		},

		"no_value_connection_url_set_key_exists": {
			want: "",
			input: input{
				envar: "",
				conf:  map[string]string{"connection_url": ""},
			},
		},

		"no_value_connection_url_set_key_doesnt_exist": {
			want: "",
			input: input{
				envar: "",
				conf:  map[string]string{},
			},
		},

		"environment_variable_set": {
			want: "abc",
			input: input{
				envar: "abc",
				conf:  map[string]string{"connection_url": "def"},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			// This is necessary to avoid always testing the branch where the env is set.
			// As long the the env is set --- even if the value is "" --- `ok` returns true.
			if tt.input.envar != "" {
				os.Setenv("VAULT_PG_CONNECTION_URL", tt.input.envar)
				defer os.Unsetenv("VAULT_PG_CONNECTION_URL")
			}

			got := connectionURL(tt.input.conf)

			if got != tt.want {
				t.Errorf("connectionURL(%s): want '%s', got '%s'", tt.input, tt.want, got)
			}
		})
	}
}

// Similar to testHABackend, but using internal implementation details to
// trigger the lock failure scenario by setting the lock renew period for one
// of the locks to a higher value than the lock TTL.
const maxTries = 3

func testPostgresSQLLockTTL(t *testing.T, ha physical.HABackend) {
	t.Log("Skipping testPostgresSQLLockTTL portion of test.")
	return

	for tries := 1; tries <= maxTries; tries++ {
		// Try this several times.  If the test environment is too slow the lock can naturally lapse
		if attemptLockTTLTest(t, ha, tries) {
			break
		}
	}
}

func attemptLockTTLTest(t *testing.T, ha physical.HABackend, tries int) bool {
	// Set much smaller lock times to speed up the test.
	lockTTL := 3
	renewInterval := time.Second * 1
	retryInterval := time.Second * 1
	longRenewInterval := time.Duration(lockTTL*2) * time.Second
	lockkey := "postgresttl"

	var leaderCh <-chan struct{}

	// Get the lock
	origLock, err := ha.LockWith(lockkey, "bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	{
		// set the first lock renew period to double the expected TTL.
		lock := origLock.(*PostgreSQLLock)
		lock.renewInterval = longRenewInterval
		lock.ttlSeconds = lockTTL

		// Attempt to lock
		lockTime := time.Now()
		leaderCh, err = lock.Lock(nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if leaderCh == nil {
			t.Fatalf("failed to get leader ch")
		}

		if tries == 1 {
			time.Sleep(3 * time.Second)
		}
		// Check the value
		held, val, err := lock.Value()
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if !held {
			if tries < maxTries && time.Since(lockTime) > (time.Second*time.Duration(lockTTL)) {
				// Our test environment is slow enough that we failed this, retry
				return false
			}
			t.Fatalf("should be held")
		}
		if val != "bar" {
			t.Fatalf("bad value: %v", val)
		}
	}

	// Second acquisition should succeed because the first lock should
	// not renew within the 3 sec TTL.
	origLock2, err := ha.LockWith(lockkey, "baz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	{
		lock2 := origLock2.(*PostgreSQLLock)
		lock2.renewInterval = renewInterval
		lock2.ttlSeconds = lockTTL
		lock2.retryInterval = retryInterval

		// Cancel attempt in 6 sec so as not to block unit tests forever
		stopCh := make(chan struct{})
		time.AfterFunc(time.Duration(lockTTL*2)*time.Second, func() {
			close(stopCh)
		})

		// Attempt to lock should work
		lockTime := time.Now()
		leaderCh2, err := lock2.Lock(stopCh)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if leaderCh2 == nil {
			t.Fatalf("should get leader ch")
		}
		defer lock2.Unlock()

		// Check the value
		held, val, err := lock2.Value()
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if !held {
			if tries < maxTries && time.Since(lockTime) > (time.Second*time.Duration(lockTTL)) {
				// Our test environment is slow enough that we failed this, retry
				return false
			}
			t.Fatalf("should be held")
		}
		if val != "baz" {
			t.Fatalf("bad value: %v", val)
		}
	}
	// The first lock should have lost the leader channel
	select {
	case <-time.After(longRenewInterval * 2):
		t.Fatalf("original lock did not have its leader channel closed.")
	case <-leaderCh:
	}
	return true
}

// Verify that once Unlock is called, we don't keep trying to renew the original
// lock.
func testPostgresSQLLockRenewal(t *testing.T, ha physical.HABackend) {
	// Get the lock
	origLock, err := ha.LockWith("pgrenewal", "bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// customize the renewal and watch intervals
	lock := origLock.(*PostgreSQLLock)
	// lock.renewInterval = time.Second * 1

	// Attempt to lock
	leaderCh, err := lock.Lock(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh == nil {
		t.Fatalf("failed to get leader ch")
	}

	// Check the value
	held, val, err := lock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "bar" {
		t.Fatalf("bad value: %v", val)
	}

	// Release the lock, which will delete the stored item
	if err := lock.Unlock(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait longer than the renewal time
	time.Sleep(1500 * time.Millisecond)

	// Attempt to lock with new lock
	newLock, err := ha.LockWith("pgrenewal", "baz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	stopCh := make(chan struct{})
	timeout := time.Duration(lock.ttlSeconds)*time.Second + lock.retryInterval + time.Second

	var leaderCh2 <-chan struct{}
	newlockch := make(chan struct{})
	go func() {
		leaderCh2, err = newLock.Lock(stopCh)
		close(newlockch)
	}()

	// Cancel attempt after lock ttl + 1s so as not to block unit tests forever
	select {
	case <-time.After(timeout):
		t.Logf("giving up on lock attempt after %v", timeout)
		close(stopCh)
	case <-newlockch:
		// pass through
	}

	// Attempt to lock should work
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh2 == nil {
		t.Fatalf("should get leader ch")
	}

	// Check the value
	held, val, err = newLock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "baz" {
		t.Fatalf("bad value: %v", val)
	}

	// Cleanup
	newLock.Unlock()
}

func setupDatabaseObjects(t *testing.T, logger log.Logger, pg *PostgreSQLBackend) {
	var err error
	// Setup tables and indexes if not exists.
	createTableSQL := fmt.Sprintf(
		"  CREATE TABLE IF NOT EXISTS %v ( "+
			"  parent_path TEXT COLLATE \"C\" NOT NULL, "+
			"  path        TEXT COLLATE \"C\", "+
			"  key         TEXT COLLATE \"C\", "+
			"  value       BYTEA, "+
			"  CONSTRAINT pkey PRIMARY KEY (path, key) "+
			" ); ", pg.table)

	_, err = pg.client.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	createIndexSQL := fmt.Sprintf(" CREATE INDEX IF NOT EXISTS parent_path_idx ON %v (parent_path); ", pg.table)

	_, err = pg.client.Exec(createIndexSQL)
	if err != nil {
		t.Fatalf("Failed to create index: %v", err)
	}

	createHaTableSQL := " CREATE TABLE IF NOT EXISTS vault_ha_locks ( " +
		" ha_key                                      TEXT COLLATE \"C\" NOT NULL, " +
		" ha_identity                                 TEXT COLLATE \"C\" NOT NULL, " +
		" ha_value                                    TEXT COLLATE \"C\", " +
		" valid_until                                 TIMESTAMP WITH TIME ZONE NOT NULL, " +
		" CONSTRAINT ha_key PRIMARY KEY (ha_key) " +
		" ); "

	_, err = pg.client.Exec(createHaTableSQL)
	if err != nil {
		t.Fatalf("Failed to create hatable: %v", err)
	}
}
