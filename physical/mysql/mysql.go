// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mysql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	metrics "github.com/armon/go-metrics"
	mysql "github.com/go-sql-driver/mysql"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/physical"
)

// Verify MySQLBackend satisfies the correct interfaces
var (
	_ physical.Backend   = (*MySQLBackend)(nil)
	_ physical.HABackend = (*MySQLBackend)(nil)
	_ physical.Lock      = (*MySQLHALock)(nil)
)

// Unreserved tls key
// Reserved values are "true", "false", "skip-verify"
const mysqlTLSKey = "default"

// MySQLBackend is a physical backend that stores data
// within MySQL database.
type MySQLBackend struct {
	dbTable      string
	dbLockTable  string
	client       *sql.DB
	statements   map[string]*sql.Stmt
	logger       log.Logger
	permitPool   *permitpool.Pool
	conf         map[string]string
	redirectHost string
	redirectPort int64
	haEnabled    bool
}

// NewMySQLBackend constructs a MySQL backend using the given API client and
// server address and credential for accessing mysql database.
func NewMySQLBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	var err error

	db, err := NewMySQLClient(conf, logger)
	if err != nil {
		return nil, err
	}

	database := conf["database"]
	if database == "" {
		database = "vault"
	}
	table := conf["table"]
	if table == "" {
		table = "vault"
	}

	err = validateDBTable(database, table)
	if err != nil {
		return nil, err
	}

	dbTable := fmt.Sprintf("`%s`.`%s`", database, table)

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	} else {
		maxParInt = physical.DefaultParallelOperations
	}

	// Check schema exists
	var schemaExist bool
	schemaRows, err := db.Query("SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", database)
	if err != nil {
		return nil, fmt.Errorf("failed to check mysql schema exist: %w", err)
	}
	defer schemaRows.Close()
	schemaExist = schemaRows.Next()

	// Check table exists
	var tableExist bool
	tableRows, err := db.Query("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?", table, database)
	if err != nil {
		return nil, fmt.Errorf("failed to check mysql table exist: %w", err)
	}
	defer tableRows.Close()
	tableExist = tableRows.Next()

	// Create the required database if it doesn't exists.
	if !schemaExist {
		if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS `" + database + "`"); err != nil {
			return nil, fmt.Errorf("failed to create mysql database: %w", err)
		}
	}

	// Create the required table if it doesn't exists.
	if !tableExist {
		create_query := "CREATE TABLE IF NOT EXISTS " + dbTable +
			" (vault_key varbinary(3072), vault_value mediumblob, PRIMARY KEY (vault_key))"
		if _, err := db.Exec(create_query); err != nil {
			return nil, fmt.Errorf("failed to create mysql table: %w", err)
		}
	}

	// Default value for ha_enabled
	haEnabledStr, ok := conf["ha_enabled"]
	if !ok {
		haEnabledStr = "false"
	}
	haEnabled, err := strconv.ParseBool(haEnabledStr)
	if err != nil {
		return nil, fmt.Errorf("value [%v] of 'ha_enabled' could not be understood", haEnabledStr)
	}

	locktable, ok := conf["lock_table"]
	if !ok {
		locktable = table + "_lock"
	}

	dbLockTable := "`" + database + "`.`" + locktable + "`"

	// Only create lock table if ha_enabled is true
	if haEnabled {
		// Check table exists
		var lockTableExist bool
		lockTableRows, err := db.Query("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?", locktable, database)
		if err != nil {
			return nil, fmt.Errorf("failed to check mysql table exist: %w", err)
		}
		defer lockTableRows.Close()
		lockTableExist = lockTableRows.Next()

		// Create the required table if it doesn't exists.
		if !lockTableExist {
			create_query := "CREATE TABLE IF NOT EXISTS " + dbLockTable +
				" (node_job varbinary(512), current_leader varbinary(512), PRIMARY KEY (node_job))"
			if _, err := db.Exec(create_query); err != nil {
				return nil, fmt.Errorf("failed to create mysql table: %w", err)
			}
		}
	}

	// Setup the backend.
	m := &MySQLBackend{
		dbTable:     dbTable,
		dbLockTable: dbLockTable,
		client:      db,
		statements:  make(map[string]*sql.Stmt),
		logger:      logger,
		permitPool:  permitpool.New(maxParInt),
		conf:        conf,
		haEnabled:   haEnabled,
	}

	// Prepare all the statements required
	statements := map[string]string{
		"put": "INSERT INTO " + dbTable +
			" VALUES( ?, ? ) ON DUPLICATE KEY UPDATE vault_value=VALUES(vault_value)",
		"get":    "SELECT vault_value FROM " + dbTable + " WHERE vault_key = ?",
		"delete": "DELETE FROM " + dbTable + " WHERE vault_key = ?",
		"list":   "SELECT vault_key FROM " + dbTable + " WHERE vault_key LIKE ?",
	}

	// Only prepare ha-related statements if we need them
	if haEnabled {
		statements["get_lock"] = "SELECT current_leader FROM " + dbLockTable + " WHERE node_job = ?"
		statements["used_lock"] = "SELECT IS_USED_LOCK(?)"
	}

	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// validateDBTable to prevent SQL injection attacks. This ensures that the database and table names only have valid
// characters in them. MySQL allows for more characters that this will allow, but there isn't an easy way of
// representing the full Unicode Basic Multilingual Plane to check against.
// https://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func validateDBTable(db, table string) (err error) {
	merr := &multierror.Error{}
	merr = multierror.Append(merr, wrapErr("invalid database: %w", validate(db)))
	merr = multierror.Append(merr, wrapErr("invalid table: %w", validate(table)))
	return merr.ErrorOrNil()
}

func validate(name string) (err error) {
	if name == "" {
		return fmt.Errorf("missing name")
	}
	// From: https://dev.mysql.com/doc/refman/5.7/en/identifiers.html
	// - Permitted characters in quoted identifiers include the full Unicode Basic Multilingual Plane (BMP), except U+0000:
	//    ASCII: U+0001 .. U+007F
	//    Extended: U+0080 .. U+FFFF
	// - ASCII NUL (U+0000) and supplementary characters (U+10000 and higher) are not permitted in quoted or unquoted identifiers.
	// - Identifiers may begin with a digit but unless quoted may not consist solely of digits.
	// - Database, table, and column names cannot end with space characters.
	//
	// We are explicitly excluding all space characters (it's easier to deal with)
	// The name will be quoted, so the all-digit requirement doesn't apply
	runes := []rune(name)
	validationErr := fmt.Errorf("invalid character found: can only include printable, non-space characters between [0x0001-0xFFFF]")
	for _, r := range runes {
		// U+0000 Explicitly disallowed
		if r == 0x0000 {
			return fmt.Errorf("invalid character: cannot include 0x0000")
		}
		// Cannot be above 0xFFFF
		if r > 0xFFFF {
			return fmt.Errorf("invalid character: cannot include any characters above 0xFFFF")
		}
		if r == '`' {
			return fmt.Errorf("invalid character: cannot include '`' character")
		}
		if r == '\'' || r == '"' {
			return fmt.Errorf("invalid character: cannot include quotes")
		}
		// We are excluding non-printable characters (not mentioned in the docs)
		if !unicode.IsPrint(r) {
			return validationErr
		}
		// We are excluding space characters (not mentioned in the docs)
		if unicode.IsSpace(r) {
			return validationErr
		}
	}
	return nil
}

func wrapErr(message string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(message, err)
}

func NewMySQLClient(conf map[string]string, logger log.Logger) (*sql.DB, error) {
	var err error

	// Get the MySQL credentials to perform read/write operations.
	username, ok := conf["username"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing username")
	}
	password, ok := conf["password"]
	if !ok || password == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get or set MySQL server address. Defaults to localhost and default port(3306)
	address, ok := conf["address"]
	if !ok {
		address = "127.0.0.1:3306"
	}

	maxIdleConnStr, ok := conf["max_idle_connections"]
	var maxIdleConnInt int
	if ok {
		maxIdleConnInt, err = strconv.Atoi(maxIdleConnStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_idle_connections parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_idle_connections set", "max_idle_connections", maxIdleConnInt)
		}
	}

	maxConnLifeStr, ok := conf["max_connection_lifetime"]
	var maxConnLifeInt int
	if ok {
		maxConnLifeInt, err = strconv.Atoi(maxConnLifeStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_connection_lifetime parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_connection_lifetime set", "max_connection_lifetime", maxConnLifeInt)
		}
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	} else {
		maxParInt = physical.DefaultParallelOperations
	}

	dsnParams := url.Values{}
	tlsCaFile, tlsOk := conf["tls_ca_file"]
	if tlsOk {
		if err := setupMySQLTLSConfig(tlsCaFile); err != nil {
			return nil, fmt.Errorf("failed register TLS config: %w", err)
		}

		dsnParams.Add("tls", mysqlTLSKey)
	}
	ptAllowed, ptOk := conf["plaintext_connection_allowed"]
	if !(ptOk && strings.ToLower(ptAllowed) == "true") && !tlsOk {
		logger.Warn("No TLS specified, credentials will be sent in plaintext. To mute this warning add 'plaintext_connection_allowed' with a true value to your MySQL configuration in your config file.")
	}

	// Create MySQL handle for the database.
	dsn := username + ":" + password + "@tcp(" + address + ")/?" + dsnParams.Encode()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}
	db.SetMaxOpenConns(maxParInt)
	if maxIdleConnInt != 0 {
		db.SetMaxIdleConns(maxIdleConnInt)
	}
	if maxConnLifeInt != 0 {
		db.SetConnMaxLifetime(time.Duration(maxConnLifeInt) * time.Second)
	}

	return db, err
}

// prepare is a helper to prepare a query for future execution
func (m *MySQLBackend) prepare(name, query string) error {
	stmt, err := m.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare %q: %w", name, err)
	}
	m.statements[name] = stmt
	return nil
}

// Put is used to insert or update an entry.
func (m *MySQLBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"mysql", "put"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer m.permitPool.Release()

	_, err := m.statements["put"].Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch an entry.
func (m *MySQLBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"mysql", "get"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer m.permitPool.Release()

	var result []byte
	err := m.statements["get"].QueryRow(key).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ent := &physical.Entry{
		Key:   key,
		Value: result,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *MySQLBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"mysql", "delete"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer m.permitPool.Release()

	_, err := m.statements["delete"].Exec(key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *MySQLBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"mysql", "list"}, time.Now())

	if err := m.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer m.permitPool.Release()

	// Add the % wildcard to the prefix to do the prefix search
	likePrefix := prefix + "%"
	rows, err := m.statements["list"].Query(likePrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		key = strings.TrimPrefix(key, prefix)
		if i := strings.Index(key, "/"); i == -1 {
			// Add objects only from the current 'folder'
			keys = append(keys, key)
		} else if i != -1 {
			// Add truncated 'folder' paths
			keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
		}
	}

	sort.Strings(keys)
	return keys, nil
}

// LockWith is used for mutual exclusion based on the given key.
func (m *MySQLBackend) LockWith(key, value string) (physical.Lock, error) {
	l := &MySQLHALock{
		in:     m,
		key:    key,
		value:  value,
		logger: m.logger,
	}
	return l, nil
}

func (m *MySQLBackend) HAEnabled() bool {
	return m.haEnabled
}

// MySQLHALock is a MySQL Lock implementation for the HABackend
type MySQLHALock struct {
	in     *MySQLBackend
	key    string
	value  string
	logger log.Logger

	held      bool
	localLock sync.Mutex
	leaderCh  chan struct{}
	stopCh    <-chan struct{}
	lock      *MySQLLock
}

func (i *MySQLHALock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	i.localLock.Lock()
	defer i.localLock.Unlock()
	if i.held {
		return nil, fmt.Errorf("lock already held")
	}

	// Attempt an async acquisition
	didLock := make(chan struct{})
	failLock := make(chan error, 1)
	releaseCh := make(chan bool, 1)
	go i.attemptLock(i.key, i.value, didLock, failLock, releaseCh)

	// Wait for lock acquisition, failure, or shutdown
	select {
	case <-didLock:
		releaseCh <- false
	case err := <-failLock:
		return nil, err
	case <-stopCh:
		releaseCh <- true
		return nil, nil
	}

	// Create the leader channel
	i.held = true
	i.leaderCh = make(chan struct{})

	go i.monitorLock(i.leaderCh)

	i.stopCh = stopCh

	return i.leaderCh, nil
}

func (i *MySQLHALock) attemptLock(key, value string, didLock chan struct{}, failLock chan error, releaseCh chan bool) {
	lock, err := NewMySQLLock(i.in, i.logger, key, value)
	if err != nil {
		failLock <- err
		return
	}

	// Set node value
	i.lock = lock

	err = lock.Lock()
	if err != nil {
		failLock <- err
		return
	}

	// Signal that lock is held
	close(didLock)

	// Handle an early abort
	release := <-releaseCh
	if release {
		lock.Unlock()
	}
}

func (i *MySQLHALock) monitorLock(leaderCh chan struct{}) {
	for {
		// The only way to lose this lock is if someone is
		// logging into the DB and altering system tables or you lose a connection in
		// which case you will lose the lock anyway.
		err := i.hasLock(i.key)
		if err != nil {
			// Somehow we lost the lock.... likely because the connection holding
			// the lock was closed or someone was playing around with the locks in the DB.
			close(leaderCh)
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func (i *MySQLHALock) Unlock() error {
	i.localLock.Lock()
	defer i.localLock.Unlock()
	if !i.held {
		return nil
	}

	err := i.lock.Unlock()

	if err == nil {
		i.held = false
		return nil
	}

	return err
}

// hasLock will check if a lock is held by checking the current lock id against our known ID.
func (i *MySQLHALock) hasLock(key string) error {
	var result sql.NullInt64
	err := i.in.statements["used_lock"].QueryRow(key).Scan(&result)
	if err == sql.ErrNoRows || !result.Valid {
		// This is not an error to us since it just means the lock isn't held
		return nil
	}

	if err != nil {
		return err
	}

	// IS_USED_LOCK will return the ID of the connection that created the lock.
	if result.Int64 != GlobalLockID {
		return ErrLockHeld
	}

	return nil
}

func (i *MySQLHALock) GetLeader() (string, error) {
	defer metrics.MeasureSince([]string{"mysql", "lock_get"}, time.Now())
	var result string
	err := i.in.statements["get_lock"].QueryRow("leader").Scan(&result)
	if err == sql.ErrNoRows {
		return "", err
	}

	return result, nil
}

func (i *MySQLHALock) Value() (bool, string, error) {
	leaderkey, err := i.GetLeader()
	if err != nil {
		return false, "", err
	}

	return true, leaderkey, err
}

// MySQLLock provides an easy way to grab and release mysql
// locks using the built in GET_LOCK function. Note that these
// locks are released when you lose connection to the server.
type MySQLLock struct {
	parentConn *MySQLBackend
	in         *sql.DB
	logger     log.Logger
	statements map[string]*sql.Stmt
	key        string
	value      string
}

// Errors specific to trying to grab a lock in MySQL
var (
	// This is the GlobalLockID for checking if the lock we got is still the current lock
	GlobalLockID int64
	// ErrLockHeld is returned when another vault instance already has a lock held for the given key.
	ErrLockHeld = errors.New("mysql: lock already held")
	// ErrUnlockFailed
	ErrUnlockFailed = errors.New("mysql: unable to release lock, already released or not held by this session")
	// You were unable to update that you are the new leader in the DB
	ErrClaimFailed = errors.New("mysql: unable to update DB with new leader information")
	// Error to throw if between getting the lock and checking the ID of it we lost it.
	ErrSettingGlobalID = errors.New("mysql: getting global lock id failed")
)

// NewMySQLLock helper function
func NewMySQLLock(in *MySQLBackend, l log.Logger, key, value string) (*MySQLLock, error) {
	// Create a new MySQL connection so we can close this and have no effect on
	// the rest of the MySQL backend and any cleanup that might need to be done.
	conn, _ := NewMySQLClient(in.conf, in.logger)

	m := &MySQLLock{
		parentConn: in,
		in:         conn,
		logger:     l,
		statements: make(map[string]*sql.Stmt),
		key:        key,
		value:      value,
	}

	statements := map[string]string{
		"put": "INSERT INTO " + in.dbLockTable +
			" VALUES( ?, ? ) ON DUPLICATE KEY UPDATE current_leader=VALUES(current_leader)",
	}

	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// prepare is a helper to prepare a query for future execution
func (m *MySQLLock) prepare(name, query string) error {
	stmt, err := m.in.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare %q: %w", name, err)
	}
	m.statements[name] = stmt
	return nil
}

// update the current cluster leader in the DB. This is used so
// we can tell the servers in standby who the active leader is.
func (i *MySQLLock) becomeLeader() error {
	_, err := i.statements["put"].Exec("leader", i.value)
	if err != nil {
		return err
	}

	return nil
}

// Lock will try to get a lock for an indefinite amount of time
// based on the given key that has been requested.
func (i *MySQLLock) Lock() error {
	defer metrics.MeasureSince([]string{"mysql", "get_lock"}, time.Now())

	// Lock timeout math.MaxInt32 instead of -1 solves compatibility issues with
	// different MySQL flavours i.e. MariaDB
	rows, err := i.in.Query("SELECT GET_LOCK(?, ?), IS_USED_LOCK(?)", i.key, math.MaxInt32, i.key)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	var lock sql.NullInt64
	var connectionID sql.NullInt64
	rows.Scan(&lock, &connectionID)

	if rows.Err() != nil {
		return rows.Err()
	}

	// 1 is returned from GET_LOCK if it was able to get the lock
	// 0 if it failed and NULL if some strange error happened.
	// https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_get-lock
	if !lock.Valid || lock.Int64 != 1 {
		return ErrLockHeld
	}

	// Since we have the lock alert the rest of the cluster
	// that we are now the active leader.
	err = i.becomeLeader()
	if err != nil {
		return ErrLockHeld
	}

	// This will return the connection ID of NULL if an error happens
	// https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_is-used-lock
	if !connectionID.Valid {
		return ErrSettingGlobalID
	}

	GlobalLockID = connectionID.Int64

	return nil
}

// Unlock just closes the connection. This is because closing the MySQL connection
// is a 100% reliable way to close the lock. If you just release the lock you must
// do it from the same mysql connection_id that you originally created it from. This
// is a huge hastle and I actually couldn't find a clean way to do this although one
// likely does exist. Closing the connection however ensures we don't ever get into a
// state where we try to release the lock and it hangs it is also much less code.
func (i *MySQLLock) Unlock() error {
	err := i.in.Close()
	if err != nil {
		return ErrUnlockFailed
	}

	return nil
}

// Establish a TLS connection with a given CA certificate
// Register a tsl.Config associated with the same key as the dns param from sql.Open
// foo:bar@tcp(127.0.0.1:3306)/dbname?tls=default
func setupMySQLTLSConfig(tlsCaFile string) error {
	rootCertPool := x509.NewCertPool()

	pem, err := ioutil.ReadFile(tlsCaFile)
	if err != nil {
		return err
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return err
	}

	err = mysql.RegisterTLSConfig(mysqlTLSKey, &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		return err
	}

	return nil
}
