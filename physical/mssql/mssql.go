// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/physical"
	mssql "github.com/microsoft/go-mssqldb"
)

// Verify MSSQLBackend satisfies the correct interfaces
var (
	_ physical.Backend = (*MSSQLBackend)(nil)
)

type MSSQLBackend struct {
	dbTable    string
	client     *sql.DB
	statements map[string]*sql.Stmt
	logger     log.Logger
	permitPool *physical.PermitPool
}

func NewMSSQLBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	var err error
	validIdentifierRE := regexp.MustCompile(`^[\p{L}_][\p{L}\p{Nd}@#$_]*$`)

	// <-- CLOSURE FUNCTION: get config value with defaults
	getConfValue := func(confKey string, defaultValue string) string {
		confVal, ok := conf[confKey]
		if ok || confVal != "" {
			return confVal
		}
		return defaultValue
	} // --> END

	username := getConfValue("username", "")
	password := getConfValue("password", "")

	server := getConfValue("server", "")
	if server == "" {
		return nil, fmt.Errorf("missing server")
	}

	port := getConfValue("port", "")

	maxParInt := physical.DefaultParallelOperations
	maxParStr := getConfValue("max_parallel", "")
	if maxParStr != "" {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	database := getConfValue("database", "Vault")
	if !validIdentifierRE.MatchString(database) {
		return nil, fmt.Errorf("invalid database name")
	}

	databaseCollation := getConfValue("databaseCollation", "")

	table := getConfValue("table", "Vault")
	if !validIdentifierRE.MatchString(table) {
		return nil, fmt.Errorf("invalid table name")
	}

	appname := getConfValue("appname", "Vault")

	connectionTimeout := getConfValue("connectiontimeout", "30")

	logLevel := getConfValue("logLevel", "0")
	// SetLogger only if needed
	if logLevel != "0" {
		mssql.SetLogger(logger.StandardLogger(nil))
	}

	schema := getConfValue("schema", "dbo")
	if !validIdentifierRE.MatchString(schema) {
		return nil, fmt.Errorf("invalid schema name")
	}

	connectionString := fmt.Sprintf("server=%s;app name=%s;connection timeout=%s;log=%s", server, appname, connectionTimeout, logLevel)
	if username != "" {
		connectionString += ";user id=" + username
	}

	if password != "" {
		connectionString += ";password=" + password
	}

	if port != "" {
		connectionString += ";port=" + port
	}

	connectionDatabase := ";database=" + database

	// <-- CLOSURE FUNCTION: openConnection
	openConnection := func(connectionStringEx string) (*sql.DB, error) {
		db, err := sql.Open("mssql", connectionString+connectionStringEx)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to mssql: %w", err)
		}
		return db, nil
	} // --> END

	// <-- CLOSURE FUNCTION: createDatabase
	errNoDefaultDb := errors.New("mssql: Cannot open requested database")
	createDatabase := func(db *sql.DB, database string, collation string) error {
		exQuery := database
		if collation != "" {
			exQuery += " COLLATE " + collation
		}
		_, err := db.Exec("IF NOT EXISTS(SELECT * FROM sys.databases WHERE name = ?) CREATE DATABASE "+exQuery, database)
		if err != nil {
			var sqlErr mssql.Error
			if errors.As(err, &sqlErr) && sqlErr.SQLErrorNumber() == 4063 {
				return errNoDefaultDb
			}
			err = fmt.Errorf("failed to create mssql database: %w", err)
		}
		return err
	} // --> END

	// Open connection with database parameter
	var db *sql.DB
	db, err = openConnection(connectionDatabase)
	if err != nil {
		return nil, err
	}

	// Create database if exist and empty
	err = createDatabase(db, database, databaseCollation)
	if err != nil {
		if err == errNoDefaultDb {
			// Database not exist
			// 4063: Cannot open database that was requested by the login. Using the user default database instead.
			err = db.Close()
			if err == nil {
				// if ok, Reopen connection without database parameter
				db, err = openConnection("")
				if err == nil {
					// if ok, Create database
					err = createDatabase(db, database, databaseCollation)
					if err == nil {
						err = db.Close()
						if err == nil {
							// if ok, Reopen connection with database parameter
							db, err = openConnection(connectionDatabase)
						}
					}
				}
			}
		}
		// Fail if there are errors
		if err != nil {
			return nil, err
		}
	}

	// Mismatched data types on table and parameter may cause long running queries
	// README.md: https://github.com/microsoft/go-mssqldb
	if databaseCollation != "" {
		dbCollation := ""
		err = db.QueryRow("SELECT DATABASEPROPERTYEX('" + database + "', 'Collation')").Scan(&dbCollation)
		if err != nil || dbCollation == "" {
			logger.Warn("Cannot get database collation", err)
		} else if dbCollation != databaseCollation {
			logger.Warn("Database and vault config collation mismatch. This may cause long running queries!", dbCollation, databaseCollation)
		}
	}

	db.SetMaxOpenConns(maxParInt)

	dbTable := schema + "." + table
	createQuery := "IF NOT EXISTS(SELECT 1 FROM " + database + ".INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_NAME=? AND TABLE_SCHEMA=?) CREATE TABLE " + dbTable + " (Path VARCHAR(512) PRIMARY KEY, Value VARBINARY(MAX))"

	if schema != "dbo" {
		var errFmt string = "failed to check if mssql schema exists: %w"
		var num int
		err = db.QueryRow("SELECT 1 FROM "+database+".sys.schemas WHERE name = ?", schema).Scan(&num)
		if err != nil && err == sql.ErrNoRows {
			// CREATE SCHEMA
			errFmt = "failed to create mssql schema: %w"
			_, err = db.Exec("USE " + database + "; EXEC ('CREATE SCHEMA " + schema + "')")
		}
		// Fail if there are errors
		if err != nil {
			return nil, fmt.Errorf(errFmt, err)
		}
	}

	_, err = db.Exec(createQuery, table, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create mssql table: %w", err)
	}

	m := &MSSQLBackend{
		dbTable:    dbTable,
		client:     db,
		statements: make(map[string]*sql.Stmt),
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}

	statements := map[string]string{
		"put": "DECLARE @qP VARCHAR(512) = CAST(:1 AS VARCHAR(512));" +
			" IF EXISTS(SELECT 1 FROM " + dbTable + " WHERE Path = @qP) UPDATE " + dbTable + " SET Value = :2 WHERE Path = @qP" +
			" ELSE INSERT INTO " + dbTable + " VALUES(@qP, :2)",
		"get":    "SELECT Value FROM " + dbTable + " WHERE Path = CAST(? AS VARCHAR(512))",
		"delete": "DELETE FROM " + dbTable + " WHERE Path = CAST(? AS VARCHAR(512))",
		"list":   "SELECT Path FROM " + dbTable + " WHERE Path LIKE CAST(? AS VARCHAR(512))",
	}

	for name, query := range statements {
		err = m.prepare(name, query)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *MSSQLBackend) prepare(name, query string) error {
	stmt, err := m.client.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare %q: %w", name, err)
	}

	m.statements[name] = stmt

	return nil
}

func (m *MSSQLBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"mssql", "put"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	_, err := m.statements["put"].Exec(entry.Key, entry.Value)
	if err != nil {
		return err
	}

	return nil
}

func (m *MSSQLBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"mssql", "get"}, time.Now())

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

func (m *MSSQLBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"mssql", "delete"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	_, err := m.statements["delete"].Exec(key)
	if err != nil {
		return err
	}

	return nil
}

func (m *MSSQLBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"mssql", "list"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	likePrefix := prefix + "%"
	rows, err := m.statements["list"].Query(likePrefix)
	if err != nil {
		return nil, err
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
			keys = append(keys, key)
		} else if i != -1 {
			keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
		}
	}

	sort.Strings(keys)

	return keys, nil
}
