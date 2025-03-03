// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	metrics "github.com/armon/go-metrics"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/physical"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Verify CockroachDBBackend satisfies the correct interfaces
var (
	_ physical.Backend       = (*CockroachDBBackend)(nil)
	_ physical.Transactional = (*CockroachDBBackend)(nil)
)

const (
	defaultTableName   = "vault_kv_store"
	defaultHATableName = "vault_ha_locks"
)

// CockroachDBBackend Backend is a physical backend that stores data
// within a CockroachDB database.
type CockroachDBBackend struct {
	table           string
	haTable         string
	client          *sql.DB
	rawStatements   map[string]string
	statements      map[string]*sql.Stmt
	rawHAStatements map[string]string
	haStatements    map[string]*sql.Stmt
	logger          log.Logger
	permitPool      *permitpool.Pool
	haEnabled       bool
}

// NewCockroachDBBackend constructs a CockroachDB backend using the given
// API client, server address, credentials, and database.
func NewCockroachDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the CockroachDB credentials to perform read/write operations.
	connURL, ok := conf["connection_url"]
	if !ok || connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	haEnabled := conf["ha_enabled"] == "true"

	dbTable := conf["table"]
	if dbTable == "" {
		dbTable = defaultTableName
	}

	err := validateDBTable(dbTable)
	if err != nil {
		return nil, fmt.Errorf("invalid table: %w", err)
	}

	dbHATable, ok := conf["ha_table"]
	if !ok {
		dbHATable = defaultHATableName
	}

	err = validateDBTable(dbHATable)
	if err != nil {
		return nil, fmt.Errorf("invalid HA table: %w", err)
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
	}

	// Create CockroachDB handle for the database.
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cockroachdb: %w", err)
	}

	// Create the required tables if they don't exist.
	createQuery := "CREATE TABLE IF NOT EXISTS " + dbTable +
		" (path STRING, value BYTES, PRIMARY KEY (path))"
	if _, err := db.Exec(createQuery); err != nil {
		return nil, fmt.Errorf("failed to create CockroachDB table: %w", err)
	}
	if haEnabled {
		createHATableQuery := "CREATE TABLE IF NOT EXISTS " + dbHATable +
			"(ha_key                                      TEXT NOT NULL, " +
			" ha_identity                                 TEXT NOT NULL, " +
			" ha_value                                    TEXT, " +
			" valid_until                                 TIMESTAMP WITH TIME ZONE NOT NULL, " +
			" CONSTRAINT ha_key PRIMARY KEY (ha_key) " +
			");"
		if _, err := db.Exec(createHATableQuery); err != nil {
			return nil, fmt.Errorf("failed to create CockroachDB HA table: %w", err)
		}
	}

	// Setup the backend
	c := &CockroachDBBackend{
		table:   dbTable,
		haTable: dbHATable,
		client:  db,
		rawStatements: map[string]string{
			"put": "INSERT INTO " + dbTable + " VALUES($1, $2)" +
				" ON CONFLICT (path) DO " +
				" UPDATE SET (path, value) = ($1, $2)",
			"get":    "SELECT value FROM " + dbTable + " WHERE path = $1",
			"delete": "DELETE FROM " + dbTable + " WHERE path = $1",
			"list":   "SELECT path FROM " + dbTable + " WHERE path LIKE $1",
		},
		statements: make(map[string]*sql.Stmt),
		rawHAStatements: map[string]string{
			"get": "SELECT ha_value FROM " + dbHATable + " WHERE NOW() <= valid_until AND ha_key = $1",
			"upsert": "INSERT INTO " + dbHATable + " as t (ha_identity, ha_key, ha_value, valid_until)" +
				" VALUES ($1, $2, $3, NOW() + $4) " +
				" ON CONFLICT (ha_key) DO " +
				" UPDATE SET (ha_identity, ha_key, ha_value, valid_until) = ($1, $2, $3, NOW() + $4) " +
				" WHERE (t.valid_until < NOW() AND t.ha_key = $2) OR " +
				" (t.ha_identity = $1 AND t.ha_key = $2)  ",
			"delete": "DELETE FROM " + dbHATable + " WHERE ha_key = $1",
		},
		haStatements: make(map[string]*sql.Stmt),
		logger:       logger,
		permitPool:   permitpool.New(maxParInt),
		haEnabled:    haEnabled,
	}

	// Prepare all the statements required
	for name, query := range c.rawStatements {
		if err := c.prepare(c.statements, name, query); err != nil {
			return nil, err
		}
	}
	if haEnabled {
		for name, query := range c.rawHAStatements {
			if err := c.prepare(c.haStatements, name, query); err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

// prepare is a helper to prepare a query for future execution.
func (c *CockroachDBBackend) prepare(statementMap map[string]*sql.Stmt, name, query string) error {
	stmt, err := c.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare %q: %w", name, err)
	}
	statementMap[name] = stmt
	return nil
}

// Put is used to insert or update an entry.
func (c *CockroachDBBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "put"}, time.Now())

	if err := c.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer c.permitPool.Release()

	_, err := c.statements["put"].Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (c *CockroachDBBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"cockroachdb", "get"}, time.Now())

	if err := c.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer c.permitPool.Release()

	var result []byte
	err := c.statements["get"].QueryRow(key).Scan(&result)
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
func (c *CockroachDBBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "delete"}, time.Now())

	if err := c.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer c.permitPool.Release()

	_, err := c.statements["delete"].Exec(key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (c *CockroachDBBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"cockroachdb", "list"}, time.Now())

	if err := c.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer c.permitPool.Release()

	likePrefix := prefix + "%"
	rows, err := c.statements["list"].Query(likePrefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// Transaction is used to run multiple entries via a transaction
func (c *CockroachDBBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	defer metrics.MeasureSince([]string{"cockroachdb", "transaction"}, time.Now())
	if len(txns) == 0 {
		return nil
	}

	if err := c.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer c.permitPool.Release()

	return crdb.ExecuteTx(context.Background(), c.client, nil, func(tx *sql.Tx) error {
		return c.transaction(tx, txns)
	})
}

func (c *CockroachDBBackend) transaction(tx *sql.Tx, txns []*physical.TxnEntry) error {
	deleteStmt, err := tx.Prepare(c.rawStatements["delete"])
	if err != nil {
		return err
	}
	putStmt, err := tx.Prepare(c.rawStatements["put"])
	if err != nil {
		return err
	}

	for _, op := range txns {
		switch op.Operation {
		case physical.DeleteOperation:
			_, err = deleteStmt.Exec(op.Entry.Key)
		case physical.PutOperation:
			_, err = putStmt.Exec(op.Entry.Key, op.Entry.Value)
		default:
			return fmt.Errorf("%q is not a supported transaction operation", op.Operation)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// validateDBTable against the CockroachDB rules for table names:
// https://www.cockroachlabs.com/docs/stable/keywords-and-identifiers.html#identifiers
//
//   - All values that accept an identifier must:
//   - Begin with a Unicode letter or an underscore (_). Subsequent characters can be letters,
//   - underscores, digits (0-9), or dollar signs ($).
//   - Not equal any SQL keyword unless the keyword is accepted by the element's syntax. For example,
//     name accepts Unreserved or Column Name keywords.
//
// The docs do state that we can bypass these rules with double quotes, however I think it
// is safer to just require these rules across the board.
func validateDBTable(dbTable string) (err error) {
	// Check if this is 'database.table' formatted. If so, split them apart and check the two
	// parts from each other
	split := strings.SplitN(dbTable, ".", 2)
	if len(split) == 2 {
		merr := &multierror.Error{}
		merr = multierror.Append(merr, wrapErr("invalid database: %w", validateDBTable(split[0])))
		merr = multierror.Append(merr, wrapErr("invalid table name: %w", validateDBTable(split[1])))
		return merr.ErrorOrNil()
	}

	// Disallow SQL keywords as the table name
	if sqlKeywords[strings.ToUpper(dbTable)] {
		return fmt.Errorf("name must not be a SQL keyword")
	}

	runes := []rune(dbTable)
	for i, r := range runes {
		if i == 0 && !unicode.IsLetter(r) && r != '_' {
			return fmt.Errorf("must use a letter or an underscore as the first character")
		}

		if !unicode.IsLetter(r) && r != '_' && !unicode.IsDigit(r) && r != '$' {
			return fmt.Errorf("must only contain letters, underscores, digits, and dollar signs")
		}

		if r == '`' || r == '\'' || r == '"' {
			return fmt.Errorf("cannot contain backticks, single quotes, or double quotes")
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
