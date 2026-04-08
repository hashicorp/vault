package ydb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	log "github.com/hashicorp/go-hclog"
	ydbconsts "github.com/hashicorp/vault/physical/ydb/consts"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type YDBBackend struct {
	db                    *ydb.Driver
	table                 string
	coordinationNode      string
	haEnabled             bool
	transactionMaxEntries int
	transactionMaxSize    int
	logger                log.Logger
}

var (
	_ physical.Backend             = (*YDBBackend)(nil)
	_ physical.HABackend           = (*YDBBackend)(nil)
	_ physical.Transactional       = (*YDBBackend)(nil)
	_ physical.TransactionalLimits = (*YDBBackend)(nil)
)

func NewYDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	var dsn string
	if envDSN := os.Getenv(ydbconsts.EnvDSN); envDSN != "" {
		dsn = strings.TrimSpace(envDSN)
	} else {
		dsn = strings.TrimSpace(conf["dsn"])
		if dsn == "" {
			return &YDBBackend{}, fmt.Errorf("YDB: dsn is not set")
		}
	}

	var table string
	if envTable := os.Getenv(ydbconsts.EnvTable); envTable != "" {
		table = strings.TrimSpace(envTable)
	} else {
		table = strings.TrimSpace(conf["table"])
		if table == "" {
			table = ydbconsts.VAULT_TABLE
		}
	}

	quotedTable, err := quoteYDBIdentifier(table)
	if err != nil {
		return &YDBBackend{}, fmt.Errorf("YDB: invalid table name: %w", err)
	}
	transactionMaxEntries, transactionMaxSize, err := getYDBTransactionLimits(conf)
	if err != nil {
		return &YDBBackend{}, fmt.Errorf("YDB: invalid transaction limits: %w", err)
	}

	opts, err := getYDBOptionsFromConfMap(conf)
	if err != nil {
		return &YDBBackend{}, fmt.Errorf("YDB: invalid options: %w", err)
	}

	ctx := context.TODO()
	db, err := ydb.Open(ctx, dsn, opts...)
	if err != nil {
		errStr := "YDB: failed to open database connection"
		logger.Error(errStr, "error", err)
		return &YDBBackend{}, fmt.Errorf(errStr+": %w", err)
	}

	if err = ensureTableExists(ctx, db, table, logger); err != nil {
		errStr := "YDB: failed to ensure table exists"
		logger.Error(errStr, "table", table, "error", err)
		return &YDBBackend{}, fmt.Errorf(errStr+": %w", err)
	}

	return &YDBBackend{
		db:                    db,
		table:                 quotedTable,
		coordinationNode:      getYDBHACoordinationNodePath(conf, db.Name(), table),
		haEnabled:             getYDBHAEnabled(conf),
		transactionMaxEntries: transactionMaxEntries,
		transactionMaxSize:    transactionMaxSize,
		logger:                logger,
	}, nil
}

func ensureTableExists(ctx context.Context, db *ydb.Driver, tableName string, logger log.Logger) error {
	fullTableName := db.Name() + "/" + tableName
	if strings.HasPrefix(tableName, "/") {
		fullTableName = tableName
	}

	if _, err := db.Scheme().DescribePath(ctx, fullTableName); err == nil {
		logger.Info("YDB: table already exists", "table", tableName)
		return nil
	}

	logger.Info("YDB: creating table", "table", tableName)

	quotedTableName, err := quoteYDBIdentifier(tableName)
	if err != nil {
		return fmt.Errorf("invalid table name %q: %w", tableName, err)
	}

	queryStmt := fmt.Sprintf(`
		CREATE TABLE %s (
			key Text NOT NULL,
			value Bytes,
			updated_at Timestamp,
			PRIMARY KEY (key)
		)`, quotedTableName)

	if err := db.Query().Exec(ctx, queryStmt); err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}

	logger.Info("YDB: table created successfully", "table", tableName)
	return nil
}

func (y *YDBBackend) Put(ctx context.Context, entry *physical.Entry) error {
	stmt := fmt.Sprintf("UPSERT INTO %s (key, value) VALUES ($key, $value)", y.table)
	err := y.db.Query().Exec(ctx,
		stmt,
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$key").Text(entry.Key).
				Param("$value").Bytes(entry.Value).Build()),
	)
	if err != nil {
		return fmt.Errorf("YDB: failed to put entry: "+entry.Key+" %w", err)
	}
	return nil
}

func (y *YDBBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	stmt := fmt.Sprintf("SELECT key AS Key, value AS Value FROM %s WHERE key = $key", y.table)
	q, err := y.db.Query().QueryRow(ctx,
		stmt,
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$key").Text(key).Build()),
	)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, fmt.Errorf("YDB: failed to get key "+key+" %w", err)
	}

	entry := physical.Entry{}
	if err = q.ScanStruct(&entry, query.WithScanStructAllowMissingColumnsFromSelect()); err != nil {
		return nil, fmt.Errorf("YDB: failed to get key "+key+" %w", err)
	}

	return &entry, nil
}

func (y *YDBBackend) Delete(ctx context.Context, key string) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE key = $key", y.table)
	err := y.db.Query().Exec(ctx,
		stmt,
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$key").Text(key).Build()),
	)
	if err != nil {
		return fmt.Errorf("YDB: failed to drop entry with key "+key+" %w", err)
	}
	return nil
}

func (y *YDBBackend) List(ctx context.Context, prefix string) ([]string, error) {
	errStr := "YDB: failed to list keys by prefix " + prefix
	likePrefix := prefix + "%"
	likeQuery := "WHERE key LIKE $prefix ORDER BY key"

	if prefix == "" {
		likePrefix = "%"
		likeQuery = ""
	}

	stmt := fmt.Sprintf("SELECT key FROM %s "+likeQuery, y.table)
	q, err := y.db.Query().Query(ctx,
		stmt,
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$prefix").Text(likePrefix).Build()),
	)
	if err != nil {
		return nil, fmt.Errorf(errStr+" %w", err)
	}
	defer q.Close(ctx)

	seen := make(map[string]struct{})
	for rs, rerr := range q.ResultSets(ctx) {
		if rerr != nil {
			return nil, fmt.Errorf(errStr+" %w", rerr)
		}
		for row, rerr := range rs.Rows(ctx) {
			if rerr != nil {
				return nil, fmt.Errorf(errStr+" %w", rerr)
			}

			var val string
			if err = row.Scan(&val); err != nil {
				return []string{}, fmt.Errorf("YDB: failed to list keys: %w", err)
			}

			rel := val
			if prefix != "" {
				rel = strings.TrimPrefix(val, prefix)
			}
			if rel == "" {
				continue
			}

			if idx := strings.Index(rel, "/"); idx != -1 {
				seen[rel[:idx+1]] = struct{}{}
			} else {
				seen[rel] = struct{}{}
			}
		}
	}

	lst := make([]string, 0, len(seen))
	for k := range seen {
		lst = append(lst, k)
	}
	return lst, nil
}

func quoteYDBIdentifier(identifier string) (string, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return "", fmt.Errorf("missing identifier")
	}

	hasLeadingSlash := strings.HasPrefix(identifier, "/")
	segments := strings.Split(identifier, "/")
	quotedSegments := make([]string, 0, len(segments))

	for idx, segment := range segments {
		if segment == "" {
			if !(hasLeadingSlash && idx == 0) {
				return "", fmt.Errorf("empty identifier segment")
			}
			continue
		}
		if err := validateYDBIdentifierSegment(segment); err != nil {
			return "", err
		}
		quotedSegments = append(quotedSegments, "`"+strings.ReplaceAll(segment, "`", "``")+"`")
	}

	if len(quotedSegments) == 0 {
		return "", fmt.Errorf("missing identifier")
	}

	quoted := strings.Join(quotedSegments, "/")
	if hasLeadingSlash {
		return "/" + quoted, nil
	}
	return quoted, nil
}

func validateYDBIdentifierSegment(segment string) error {
	if segment == "" {
		return fmt.Errorf("empty identifier segment")
	}
	if segment == "." || segment == ".." {
		return fmt.Errorf("reserved identifier segment %q", segment)
	}

	for _, r := range segment {
		if r == 0 {
			return fmt.Errorf("identifier segment %q contains NUL", segment)
		}
		if !unicode.IsPrint(r) {
			return fmt.Errorf("identifier segment %q contains non-printable characters", segment)
		}
	}

	return nil
}
