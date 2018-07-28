package spanner

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/physical"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"

	"cloud.google.com/go/spanner"
	"github.com/pkg/errors"
)

// Verify Backend satisfies the correct interfaces
var _ physical.Backend = (*Backend)(nil)
var _ physical.Transactional = (*Backend)(nil)

const (
	// envDatabase is the name of the environment variable to search for the
	// database name.
	envDatabase = "GOOGLE_SPANNER_DATABASE"

	// envHAEnabled is the name of the environment variable to search for the
	// boolean indicating if HA is enabled.
	envHAEnabled = "GOOGLE_SPANNER_HA_ENABLED"

	// envHATable is the name of the environment variable to search for the table
	// name to use for HA.
	envHATable = "GOOGLE_SPANNER_HA_TABLE"

	// envTable is the name of the environment variable to search for the table
	// name.
	envTable = "GOOGLE_SPANNER_TABLE"

	// defaultTable is the default table name if none is specified.
	defaultTable = "Vault"

	// defaultHASuffix is the default suffix to apply to the table name if no
	// HA table is provided.
	defaultHASuffix = "HA"
)

var (
	// metricDelete is the key for the metric for measuring a Delete call.
	metricDelete = []string{"spanner", "delete"}

	// metricGet is the key for the metric for measuring a Get call.
	metricGet = []string{"spanner", "get"}

	// metricList is the key for the metric for measuring a List call.
	metricList = []string{"spanner", "list"}

	// metricPut is the key for the metric for measuring a Put call.
	metricPut = []string{"spanner", "put"}
)

// Backend implements physical.Backend and describes the steps necessary to
// persist data using Google Cloud Spanner.
type Backend struct {
	// database is the name of the database to use for data storage and retrieval.
	// This is supplied as part of user configuration.
	database string

	// table is the name of the table in the database.
	table string

	// haTable is the name of the table to use for HA in the database.
	haTable string

	// haEnabled indicates if high availability is enabled. Default: true.
	haEnabled bool

	// client is the underlying API client for talking to spanner.
	client *spanner.Client

	// logger and permitPool are internal constructs.
	logger     log.Logger
	permitPool *physical.PermitPool
}

// NewBackend creates a new Google Spanner storage backend with the given
// configuration. This uses the official Golang Cloud SDK and therefore supports
// specifying credentials via envvars, credential files, etc.
func NewBackend(c map[string]string, logger log.Logger) (physical.Backend, error) {
	logger.Debug("configuring backend")

	// Database name
	database := os.Getenv(envDatabase)
	if database == "" {
		database = c["database"]
	}
	if database == "" {
		return nil, errors.New("missing database name")
	}

	// Table name
	table := os.Getenv(envTable)
	if table == "" {
		table = c["table"]
	}
	if table == "" {
		table = defaultTable
	}

	// HA table name
	haTable := os.Getenv(envHATable)
	if haTable == "" {
		haTable = c["ha_table"]
	}
	if haTable == "" {
		haTable = table + defaultHASuffix
	}

	// HA configuration
	haEnabled := false
	haEnabledStr := os.Getenv(envHAEnabled)
	if haEnabledStr == "" {
		haEnabledStr = c["ha_enabled"]
	}
	if haEnabledStr != "" {
		var err error
		haEnabled, err = strconv.ParseBool(haEnabledStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed to parse HA enabled: {{err}}", err)
		}
	}

	// Max parallel
	maxParallel, err := extractInt(c["max_parallel"])
	if err != nil {
		return nil, errwrap.Wrapf("failed to parse max_parallel: {{err}}", err)
	}

	logger.Debug("configuration",
		"database", database,
		"table", table,
		"haEnabled", haEnabled,
		"haTable", haTable,
		"maxParallel", maxParallel,
	)
	logger.Debug("creating client")

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, database,
		option.WithUserAgent(useragent.String()),
	)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create spanner client: {{err}}", err)
	}

	return &Backend{
		database:  database,
		table:     table,
		haEnabled: haEnabled,
		haTable:   haTable,

		client:     client,
		permitPool: physical.NewPermitPool(maxParallel),
		logger:     logger,
	}, nil
}

// Put creates or updates an entry.
func (b *Backend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince(metricPut, time.Now())

	// Pooling
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	// Insert
	m := spanner.InsertOrUpdateMap(b.table, map[string]interface{}{
		"Key":   entry.Key,
		"Value": entry.Value,
	})
	if _, err := b.client.Apply(ctx, []*spanner.Mutation{m}); err != nil {
		return errwrap.Wrapf("failed to put data: {{err}}", err)
	}
	return nil
}

// Get fetches an entry. If there is no entry, this function returns nil.
func (b *Backend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince(metricList, time.Now())

	// Pooling
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	// Read
	row, err := b.client.Single().ReadRow(ctx, b.table, spanner.Key{key}, []string{"Value"})
	if spanner.ErrCode(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to read value for %q: {{err}}", key), err)
	}

	var value []byte
	if err := row.Column(0, &value); err != nil {
		return nil, errwrap.Wrapf("failed to decode value into bytes: {{err}}", err)
	}

	return &physical.Entry{
		Key:   key,
		Value: value,
	}, nil
}

// Delete deletes an entry with the given key.
func (b *Backend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince(metricDelete, time.Now())

	// Pooling
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	// Delete
	m := spanner.Delete(b.table, spanner.Key{key})
	if _, err := b.client.Apply(ctx, []*spanner.Mutation{m}); err != nil {
		return errwrap.Wrapf("failed to delete key: {{err}}", err)
	}

	return nil
}

// List enumerates all keys with the given prefix.
func (b *Backend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince(metricList, time.Now())

	// Pooling
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	// Sanitize
	safeTable := sanitizeTable(b.table)

	// List
	iter := b.client.Single().Query(ctx, spanner.Statement{
		SQL: "SELECT Key FROM " + safeTable + " WHERE STARTS_WITH(Key, @prefix)",
		Params: map[string]interface{}{
			"prefix": prefix,
		},
	})
	defer iter.Stop()

	var keys []string

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errwrap.Wrapf("failed to read row: {{err}}", err)
		}

		var key string
		if err := row.Column(0, &key); err != nil {
			return nil, errwrap.Wrapf("failed to decode key into string: {{err}}", err)
		}

		// The results will include the full prefix (folder) and any deeply-nested
		// prefixes (subfolders). Vault expects only the top-most things to be
		// included.
		key = strings.TrimPrefix(key, prefix)
		if i := strings.Index(key, "/"); i == -1 {
			// Add objects only from the current 'folder'
			keys = append(keys, key)
		} else {
			// Add truncated 'folder' paths
			keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
		}
	}

	// Sort because the resulting order is not predictable
	sort.Strings(keys)

	return keys, nil
}

// Transaction runs multiple entries via a single transaction.
func (b *Backend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	// Quit early if we can
	if len(txns) == 0 {
		return nil
	}

	// Build all the ops before taking out the pool
	ms := make([]*spanner.Mutation, len(txns))
	for i, tx := range txns {
		op, key, value := tx.Operation, tx.Entry.Key, tx.Entry.Value

		switch op {
		case physical.DeleteOperation:
			ms[i] = spanner.Delete(b.table, spanner.Key{key})
		case physical.PutOperation:
			ms[i] = spanner.InsertOrUpdateMap(b.table, map[string]interface{}{
				"Key":   key,
				"Value": value,
			})
		default:
			return fmt.Errorf("unsupported transaction operation: %q", op)
		}
	}

	// Pooling
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	// Transactivate!
	if _, err := b.client.Apply(ctx, ms); err != nil {
		return errwrap.Wrapf("failed to commit transaction: {{err}}", err)
	}

	return nil
}

// extractInt is a helper function that takes a string and converts that string
// to an int, but accounts for the empty string.
func extractInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// sanitizeTable attempts to sanitize the table name.
func sanitizeTable(s string) string {
	end := strings.IndexRune(s, 0)
	if end > -1 {
		s = s[:end]
	}
	return strings.Replace(s, `"`, `""`, -1)
}
