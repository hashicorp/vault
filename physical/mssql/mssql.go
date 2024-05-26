// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mssql

import (
	"context"
	"database/sql"
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

	logLevel := getConfValue("loglevel", "0")

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

	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mssql: %w", err)
	}

	db.SetMaxOpenConns(maxParInt)

	if _, err := db.Exec("IF NOT EXISTS(SELECT * FROM sys.databases WHERE name = ?) CREATE DATABASE "+database, database); err != nil {
		return nil, fmt.Errorf("failed to create mssql database: %w", err)
	}

	dbTable := database + "." + schema + "." + table
	createQuery := "IF NOT EXISTS(SELECT 1 FROM " + database + ".INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_NAME=? AND TABLE_SCHEMA=?) CREATE TABLE " + dbTable + " (Path VARCHAR(512) PRIMARY KEY, Value VARBINARY(MAX))"

	if schema != "dbo" {

		var num int
		err = db.QueryRow("SELECT 1 FROM "+database+".sys.schemas WHERE name = ?", schema).Scan(&num)

		switch {
		case err == sql.ErrNoRows:
			if _, err := db.Exec("USE " + database + "; EXEC ('CREATE SCHEMA " + schema + "')"); err != nil {
				return nil, fmt.Errorf("failed to create mssql schema: %w", err)
			}

		case err != nil:
			return nil, fmt.Errorf("failed to check if mssql schema exists: %w", err)
		}
	}

	if _, err := db.Exec(createQuery, table, schema); err != nil {
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
		"put": "IF EXISTS(SELECT 1 FROM " + dbTable + " WHERE Path = ?) UPDATE " + dbTable + " SET Value = ? WHERE Path = ?" +
			" ELSE INSERT INTO " + dbTable + " VALUES(?, ?)",
		"get":    "SELECT Value FROM " + dbTable + " WHERE Path = ?",
		"delete": "DELETE FROM " + dbTable + " WHERE Path = ?",
		"list":   "SELECT Path FROM " + dbTable + " WHERE Path LIKE ?",
	}

	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
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

	_, err := m.statements["put"].Exec(entry.Key, entry.Value, entry.Key, entry.Key, entry.Value)
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
