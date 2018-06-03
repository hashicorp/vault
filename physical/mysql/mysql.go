package mysql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"

	"github.com/armon/go-metrics"
	mysql "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/physical"
)

// Verify MySQLBackend satisfies the correct interfaces
var _ physical.Backend = (*MySQLBackend)(nil)
var _ physical.HABackend = (*MySQLBackend)(nil)
var _ physical.Lock = (*MySQLHALock)(nil)

// Unreserved tls key
// Reserved values are "true", "false", "skip-verify"
const mysqlTLSKey = "default"

// MySQLBackend is a physical backend that stores data
// within MySQL database.
type MySQLBackend struct {
	dbTable    string
	client     *sql.DB
	statements map[string]*sql.Stmt
	logger     log.Logger
	permitPool *physical.PermitPool
}

// NewMySQLBackend constructs a MySQL backend using the given API client and
// server address and credential for accessing mysql database.
func NewMySQLBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	var err error

	// Get the MySQL credentials to perform read/write operations.
	username, ok := conf["username"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing username")
	}
	password, ok := conf["password"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get or set MySQL server address. Defaults to localhost and default port(3306)
	address, ok := conf["address"]
	if !ok {
		address = "127.0.0.1:3306"
	}

	// Get the MySQL database and table details.
	database, ok := conf["database"]
	if !ok {
		database = "vault"
	}
	table, ok := conf["table"]
	if !ok {
		table = "vault"
	}
	dbTable := database + "." + table

	maxIdleConnStr, ok := conf["max_idle_connections"]
	var maxIdleConnInt int
	if ok {
		maxIdleConnInt, err = strconv.Atoi(maxIdleConnStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_idle_connections parameter: {{err}}", err)
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
			return nil, errwrap.Wrapf("failed parsing max_connection_lifetime parameter: {{err}}", err)
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
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	} else {
		maxParInt = physical.DefaultParallelOperations
	}

	dsnParams := url.Values{}
	tlsCaFile, ok := conf["tls_ca_file"]
	if ok {
		if err := setupMySQLTLSConfig(tlsCaFile); err != nil {
			return nil, errwrap.Wrapf("failed register TLS config: {{err}}", err)
		}

		dsnParams.Add("tls", mysqlTLSKey)
	}

	// Create MySQL handle for the database.
	dsn := username + ":" + password + "@tcp(" + address + ")/?" + dsnParams.Encode()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errwrap.Wrapf("failed to connect to mysql: {{err}}", err)
	}
	db.SetMaxOpenConns(maxParInt)
	if maxIdleConnInt != 0 {
		db.SetMaxIdleConns(maxIdleConnInt)
	}
	if maxConnLifeInt != 0 {
		db.SetConnMaxLifetime(time.Duration(maxConnLifeInt) * time.Second)
	}
	// Check schema exists
	var schemaExist bool
	schemaRows, err := db.Query("SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", database)
	if err != nil {
		return nil, errwrap.Wrapf("failed to check mysql schema exist: {{err}}", err)
	}
	defer schemaRows.Close()
	schemaExist = schemaRows.Next()

	// Check table exists
	var tableExist bool
	tableRows, err := db.Query("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?", table, database)

	if err != nil {
		return nil, errwrap.Wrapf("failed to check mysql table exist: {{err}}", err)
	}
	defer tableRows.Close()
	tableExist = tableRows.Next()

	// Create the required database if it doesn't exists.
	if !schemaExist {
		if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
			return nil, errwrap.Wrapf("failed to create mysql database: {{err}}", err)
		}
	}

	// Create the required table if it doesn't exists.
	if !tableExist {
		create_query := "CREATE TABLE IF NOT EXISTS " + dbTable +
			" (vault_key varbinary(512), vault_value mediumblob, PRIMARY KEY (vault_key))"
		if _, err := db.Exec(create_query); err != nil {
			return nil, errwrap.Wrapf("failed to create mysql table: {{err}}", err)
		}
	}

	// Setup the backend.
	m := &MySQLBackend{
		dbTable:    dbTable,
		client:     db,
		statements: make(map[string]*sql.Stmt),
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}

	// Prepare all the statements required
	statements := map[string]string{
		"put": "INSERT INTO " + dbTable +
			" VALUES( ?, ? ) ON DUPLICATE KEY UPDATE vault_value=VALUES(vault_value)",
		"get":    "SELECT vault_value FROM " + dbTable + " WHERE vault_key = ?",
		"delete": "DELETE FROM " + dbTable + " WHERE vault_key = ?",
		"list":   "SELECT vault_key FROM " + dbTable + " WHERE vault_key LIKE ?",
	}
	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// prepare is a helper to prepare a query for future execution
func (m *MySQLBackend) prepare(name, query string) error {
	stmt, err := m.client.Prepare(query)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to prepare %q: {{err}}", name), err)
	}
	m.statements[name] = stmt
	return nil
}

// Put is used to insert or update an entry.
func (m *MySQLBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"mysql", "put"}, time.Now())

	m.permitPool.Acquire()
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

	m.permitPool.Acquire()
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

	m.permitPool.Acquire()
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

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	// Add the % wildcard to the prefix to do the prefix search
	likePrefix := prefix + "%"
	rows, err := m.statements["list"].Query(likePrefix)
	if err != nil {
		return nil, errwrap.Wrapf("failed to execute statement: {{err}}", err)
	}

	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, errwrap.Wrapf("failed to scan rows: {{err}}", err)
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
		key:    key, //+ ":" + value,
		logger: m.logger,
	}
	return l, nil
}

// HAEnabled ...
func (m *MySQLBackend) HAEnabled() bool {
	return true
}

// MySQLHALock is a MySQL Lock implementation for the HABackend
type MySQLHALock struct {
	in     *MySQLBackend
	key    string
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
	lockpath := i.key
	go i.attemptLock(lockpath, didLock, failLock, releaseCh)

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

func (i *MySQLHALock) attemptLock(lockpath string, didLock chan struct{}, failLock chan error, releaseCh chan bool) {
	// Wait to acquire the lock in ZK
	lock := NewMySQLLock(i.in, i.logger, lockpath)
	err := lock.Lock()
	if err != nil {
		failLock <- err
		return
	}
	// Set node value
	i.lock = lock

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
		// monitor this lock to make sure we still hold the lock
		// we have to poll since I know no way to have mysql push to us a notification
		// when a lock is lost. However the only way to lose this lock is if someone is
		// logging into the DB and altering system tables or you lose a connection in
		// which case you will lose the lock anyway.
		err := i.lock.HasLock()
		if err != nil {
			// Somehow we lost the lock.... this should absolutely never happen
			// unless someone is messing around in the DB doing malicious things.
			i.logger.Info("mysql: distributed lock released")
			close(leaderCh)
		}

		time.Sleep(5 * time.Second)
	}
}

func (i *MySQLHALock) unlockInternal() error {
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

func (i *MySQLHALock) Unlock() error {
	var err error

	i.unlockInternal()

	return err
}

func (i *MySQLHALock) hasLockOnKey() error {
	lockRows, err := i.in.client.Query("SELECT IS_USED_LOCK(?)", i.key)
	if err != nil {
		return err
	}

	// Check the row to see if we actually were given the lock or not.
	// Zero means someone else already has the lock where one is returned
	// if we hold the lock.
	defer lockRows.Close()
	lockRows.Next()
	var locknum int
	lockRows.Scan(&locknum)

	if locknum <= 0 {
		return ErrLockHeld
	}

	return nil
}

func (i *MySQLHALock) Value() (bool, string, error) {
	lockpath := i.key
	err := i.hasLockOnKey()
	if err != nil {
		return false, lockpath, err
	}

	return true, lockpath, err
}

// MySQLLock provides an easy way to grab and release mysql
// locks using the built in GET_LOCK function. Note that these
// locks are released when you lose connection to the server.
type MySQLLock struct {
	in     *MySQLBackend
	logger log.Logger
	key    string
}

// Errors specific to trying to grab a lock in MySQL
var (
	// ErrLockHeld is returned when another vault instance already has a lock held for the given key.
	ErrLockHeld = errors.New("mysql: lock already held")
	// ErrUnlockFailed
	ErrUnlockFailed = errors.New("mysql: unable to release lock, already released or not held by this session")
)

// NewMySQLLock helper function
func NewMySQLLock(in *MySQLBackend, l log.Logger, key string) *MySQLLock {
	return &MySQLLock{
		in:     in,
		logger: l,
		key:    key,
	}
}

// HasLock will check if a lock is held and if you are the current owner of it.
// If both conditions are not met an error is returned.
func (i *MySQLLock) HasLock() error {
	lockRows, err := i.in.client.Query("SELECT IS_USED_LOCK(?)", i.key)
	if err != nil {
		return err
	}

	// Check the row to see if we actually were given the lock or not.
	// Zero means someone else already has the lock where one is returned
	// if we hold the lock.
	defer lockRows.Close()
	lockRows.Next()
	var locknum int
	lockRows.Scan(&locknum)

	if locknum <= 0 {
		return ErrLockHeld
	}

	return nil
}

// Lock will try to get a lock for an indefinite amount of time
// based on the given key that has been requested.
func (i *MySQLLock) Lock() error {
	rows, err := i.in.client.Query("SELECT GET_LOCK(?, -1)", i.key)
	if err != nil {
		return err
	}

	// Check the row to see if we actually were given the lock or not.
	// Zero means someone else already has the lock where one is returned
	// if we hold the lock.
	defer rows.Close()
	rows.Next()
	var num int
	rows.Scan(&num)

	if num != 1 {
		return ErrLockHeld
	}

	return nil
}

// Unlock will try to relase the lock that we currently think we are holding.
func (i *MySQLLock) Unlock() error {
	rows, err := i.in.client.Query("SELECT RELEASE_LOCK(?, -1)", i.key)
	if err != nil {
		return err
	}

	// Check the row to see if we actually were able to release the lock.
	// Zero means someone else already has the lock where one is returned
	// if we were able to release the lock.
	defer rows.Close()
	rows.Next()
	var num int
	rows.Scan(&num)

	if num != 1 {
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
