// Copyright (c) 2015 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2016, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// schema metadata for a keyspace
type KeyspaceMetadata struct {
	Name            string
	DurableWrites   bool
	StrategyClass   string
	StrategyOptions map[string]interface{}
	Tables          map[string]*TableMetadata
	Functions       map[string]*FunctionMetadata
	Aggregates      map[string]*AggregateMetadata
	// Deprecated: use the MaterializedViews field for views and UserTypes field for udts instead.
	Views             map[string]*ViewMetadata
	MaterializedViews map[string]*MaterializedViewMetadata
	UserTypes         map[string]*UserTypeMetadata
}

// schema metadata for a table (a.k.a. column family)
type TableMetadata struct {
	Keyspace          string
	Name              string
	KeyValidator      string
	Comparator        string
	DefaultValidator  string
	KeyAliases        []string
	ColumnAliases     []string
	ValueAlias        string
	PartitionKey      []*ColumnMetadata
	ClusteringColumns []*ColumnMetadata
	Columns           map[string]*ColumnMetadata
	OrderedColumns    []string
}

// schema metadata for a column
type ColumnMetadata struct {
	Keyspace        string
	Table           string
	Name            string
	ComponentIndex  int
	Kind            ColumnKind
	Validator       string
	Type            TypeInfo
	ClusteringOrder string
	Order           ColumnOrder
	Index           ColumnIndexMetadata
}

// FunctionMetadata holds metadata for function constructs
type FunctionMetadata struct {
	Keyspace          string
	Name              string
	ArgumentTypes     []TypeInfo
	ArgumentNames     []string
	Body              string
	CalledOnNullInput bool
	Language          string
	ReturnType        TypeInfo
}

// AggregateMetadata holds metadata for aggregate constructs
type AggregateMetadata struct {
	Keyspace      string
	Name          string
	ArgumentTypes []TypeInfo
	FinalFunc     FunctionMetadata
	InitCond      string
	ReturnType    TypeInfo
	StateFunc     FunctionMetadata
	StateType     TypeInfo

	stateFunc string
	finalFunc string
}

// ViewMetadata holds the metadata for views.
// Deprecated: this is kept for backwards compatibility issues. Use MaterializedViewMetadata.
type ViewMetadata struct {
	Keyspace   string
	Name       string
	FieldNames []string
	FieldTypes []TypeInfo
}

// MaterializedViewMetadata holds the metadata for materialized views.
type MaterializedViewMetadata struct {
	Keyspace                string
	Name                    string
	BaseTableId             UUID
	BaseTable               *TableMetadata
	BloomFilterFpChance     float64
	Caching                 map[string]string
	Comment                 string
	Compaction              map[string]string
	Compression             map[string]string
	CrcCheckChance          float64
	DcLocalReadRepairChance float64
	DefaultTimeToLive       int
	Extensions              map[string]string
	GcGraceSeconds          int
	Id                      UUID
	IncludeAllColumns       bool
	MaxIndexInterval        int
	MemtableFlushPeriodInMs int
	MinIndexInterval        int
	ReadRepairChance        float64
	SpeculativeRetry        string

	baseTableName string
}

type UserTypeMetadata struct {
	Keyspace   string
	Name       string
	FieldNames []string
	FieldTypes []TypeInfo
}

// the ordering of the column with regard to its comparator
type ColumnOrder bool

const (
	ASC  ColumnOrder = false
	DESC ColumnOrder = true
)

type ColumnIndexMetadata struct {
	Name    string
	Type    string
	Options map[string]interface{}
}

type ColumnKind int

const (
	ColumnUnkownKind ColumnKind = iota
	ColumnPartitionKey
	ColumnClusteringKey
	ColumnRegular
	ColumnCompact
	ColumnStatic
)

func (c ColumnKind) String() string {
	switch c {
	case ColumnPartitionKey:
		return "partition_key"
	case ColumnClusteringKey:
		return "clustering_key"
	case ColumnRegular:
		return "regular"
	case ColumnCompact:
		return "compact"
	case ColumnStatic:
		return "static"
	default:
		return fmt.Sprintf("unknown_column_%d", c)
	}
}

func (c *ColumnKind) UnmarshalCQL(typ TypeInfo, p []byte) error {
	if typ.Type() != TypeVarchar {
		return unmarshalErrorf("unable to marshall %s into ColumnKind, expected Varchar", typ)
	}

	kind, err := columnKindFromSchema(string(p))
	if err != nil {
		return err
	}
	*c = kind

	return nil
}

func columnKindFromSchema(kind string) (ColumnKind, error) {
	switch kind {
	case "partition_key":
		return ColumnPartitionKey, nil
	case "clustering_key", "clustering":
		return ColumnClusteringKey, nil
	case "regular":
		return ColumnRegular, nil
	case "compact_value":
		return ColumnCompact, nil
	case "static":
		return ColumnStatic, nil
	default:
		return -1, fmt.Errorf("unknown column kind: %q", kind)
	}
}

// default alias values
const (
	DEFAULT_KEY_ALIAS    = "key"
	DEFAULT_COLUMN_ALIAS = "column"
	DEFAULT_VALUE_ALIAS  = "value"
)

// queries the cluster for schema information for a specific keyspace
type schemaDescriber struct {
	session *Session
	mu      sync.Mutex

	cache map[string]*KeyspaceMetadata
}

// creates a session bound schema describer which will query and cache
// keyspace metadata
func newSchemaDescriber(session *Session) *schemaDescriber {
	return &schemaDescriber{
		session: session,
		cache:   map[string]*KeyspaceMetadata{},
	}
}

// returns the cached KeyspaceMetadata held by the describer for the named
// keyspace.
func (s *schemaDescriber) getSchema(keyspaceName string) (*KeyspaceMetadata, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, found := s.cache[keyspaceName]
	if !found {
		// refresh the cache for this keyspace
		err := s.refreshSchema(keyspaceName)
		if err != nil {
			return nil, err
		}

		metadata = s.cache[keyspaceName]
	}

	return metadata, nil
}

// clears the already cached keyspace metadata
func (s *schemaDescriber) clearSchema(keyspaceName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.cache, keyspaceName)
}

// forcibly updates the current KeyspaceMetadata held by the schema describer
// for a given named keyspace.
func (s *schemaDescriber) refreshSchema(keyspaceName string) error {
	var err error

	// query the system keyspace for schema data
	// TODO retrieve concurrently
	keyspace, err := getKeyspaceMetadata(s.session, keyspaceName)
	if err != nil {
		return err
	}
	tables, err := getTableMetadata(s.session, keyspaceName)
	if err != nil {
		return err
	}
	columns, err := getColumnMetadata(s.session, keyspaceName)
	if err != nil {
		return err
	}
	functions, err := getFunctionsMetadata(s.session, keyspaceName)
	if err != nil {
		return err
	}
	aggregates, err := getAggregatesMetadata(s.session, keyspaceName)
	if err != nil {
		return err
	}
	views, err := getViewsMetadata(s.session, keyspaceName)
	if err != nil {
		return err
	}
	materializedViews, err := getMaterializedViewsMetadata(s.session, keyspaceName)
	if err != nil {
		return err
	}

	// organize the schema data
	compileMetadata(s.session.cfg.ProtoVersion, keyspace, tables, columns, functions, aggregates, views,
		materializedViews, s.session.logger)

	// update the cache
	s.cache[keyspaceName] = keyspace

	return nil
}

// "compiles" derived information about keyspace, table, and column metadata
// for a keyspace from the basic queried metadata objects returned by
// getKeyspaceMetadata, getTableMetadata, and getColumnMetadata respectively;
// Links the metadata objects together and derives the column composition of
// the partition key and clustering key for a table.
func compileMetadata(
	protoVersion int,
	keyspace *KeyspaceMetadata,
	tables []TableMetadata,
	columns []ColumnMetadata,
	functions []FunctionMetadata,
	aggregates []AggregateMetadata,
	views []ViewMetadata,
	materializedViews []MaterializedViewMetadata,
	logger StdLogger,
) {
	keyspace.Tables = make(map[string]*TableMetadata)
	for i := range tables {
		tables[i].Columns = make(map[string]*ColumnMetadata)

		keyspace.Tables[tables[i].Name] = &tables[i]
	}
	keyspace.Functions = make(map[string]*FunctionMetadata, len(functions))
	for i := range functions {
		keyspace.Functions[functions[i].Name] = &functions[i]
	}
	keyspace.Aggregates = make(map[string]*AggregateMetadata, len(aggregates))
	for i, _ := range aggregates {
		aggregates[i].FinalFunc = *keyspace.Functions[aggregates[i].finalFunc]
		aggregates[i].StateFunc = *keyspace.Functions[aggregates[i].stateFunc]
		keyspace.Aggregates[aggregates[i].Name] = &aggregates[i]
	}
	keyspace.Views = make(map[string]*ViewMetadata, len(views))
	for i := range views {
		keyspace.Views[views[i].Name] = &views[i]
	}
	// Views currently holds the types and hasn't been deleted for backward compatibility issues.
	// That's why it's ok to copy Views into Types in this case. For the real Views use MaterializedViews.
	types := make([]UserTypeMetadata, len(views))
	for i := range views {
		types[i].Keyspace = views[i].Keyspace
		types[i].Name = views[i].Name
		types[i].FieldNames = views[i].FieldNames
		types[i].FieldTypes = views[i].FieldTypes
	}
	keyspace.UserTypes = make(map[string]*UserTypeMetadata, len(views))
	for i := range types {
		keyspace.UserTypes[types[i].Name] = &types[i]
	}
	keyspace.MaterializedViews = make(map[string]*MaterializedViewMetadata, len(materializedViews))
	for i, _ := range materializedViews {
		materializedViews[i].BaseTable = keyspace.Tables[materializedViews[i].baseTableName]
		keyspace.MaterializedViews[materializedViews[i].Name] = &materializedViews[i]
	}

	// add columns from the schema data
	for i := range columns {
		col := &columns[i]
		// decode the validator for TypeInfo and order
		if col.ClusteringOrder != "" { // Cassandra 3.x+
			col.Type = getCassandraType(col.Validator, logger)
			col.Order = ASC
			if col.ClusteringOrder == "desc" {
				col.Order = DESC
			}
		} else {
			validatorParsed := parseType(col.Validator, logger)
			col.Type = validatorParsed.types[0]
			col.Order = ASC
			if validatorParsed.reversed[0] {
				col.Order = DESC
			}
		}

		table, ok := keyspace.Tables[col.Table]
		if !ok {
			// if the schema is being updated we will race between seeing
			// the metadata be complete. Potentially we should check for
			// schema versions before and after reading the metadata and
			// if they dont match try again.
			continue
		}

		table.Columns[col.Name] = col
		table.OrderedColumns = append(table.OrderedColumns, col.Name)
	}

	if protoVersion == protoVersion1 {
		compileV1Metadata(tables, logger)
	} else {
		compileV2Metadata(tables, logger)
	}
}

// Compiles derived information from TableMetadata which have had
// ColumnMetadata added already. V1 protocol does not return as much
// column metadata as V2+ (because V1 doesn't support the "type" column in the
// system.schema_columns table) so determining PartitionKey and ClusterColumns
// is more complex.
func compileV1Metadata(tables []TableMetadata, logger StdLogger) {
	for i := range tables {
		table := &tables[i]

		// decode the key validator
		keyValidatorParsed := parseType(table.KeyValidator, logger)
		// decode the comparator
		comparatorParsed := parseType(table.Comparator, logger)

		// the partition key length is the same as the number of types in the
		// key validator
		table.PartitionKey = make([]*ColumnMetadata, len(keyValidatorParsed.types))

		// V1 protocol only returns "regular" columns from
		// system.schema_columns (there is no type field for columns)
		// so the alias information is used to
		// create the partition key and clustering columns

		// construct the partition key from the alias
		for i := range table.PartitionKey {
			var alias string
			if len(table.KeyAliases) > i {
				alias = table.KeyAliases[i]
			} else if i == 0 {
				alias = DEFAULT_KEY_ALIAS
			} else {
				alias = DEFAULT_KEY_ALIAS + strconv.Itoa(i+1)
			}

			column := &ColumnMetadata{
				Keyspace:       table.Keyspace,
				Table:          table.Name,
				Name:           alias,
				Type:           keyValidatorParsed.types[i],
				Kind:           ColumnPartitionKey,
				ComponentIndex: i,
			}

			table.PartitionKey[i] = column
			table.Columns[alias] = column
		}

		// determine the number of clustering columns
		size := len(comparatorParsed.types)
		if comparatorParsed.isComposite {
			if len(comparatorParsed.collections) != 0 ||
				(len(table.ColumnAliases) == size-1 &&
					comparatorParsed.types[size-1].Type() == TypeVarchar) {
				size = size - 1
			}
		} else {
			if !(len(table.ColumnAliases) != 0 || len(table.Columns) == 0) {
				size = 0
			}
		}

		table.ClusteringColumns = make([]*ColumnMetadata, size)

		for i := range table.ClusteringColumns {
			var alias string
			if len(table.ColumnAliases) > i {
				alias = table.ColumnAliases[i]
			} else if i == 0 {
				alias = DEFAULT_COLUMN_ALIAS
			} else {
				alias = DEFAULT_COLUMN_ALIAS + strconv.Itoa(i+1)
			}

			order := ASC
			if comparatorParsed.reversed[i] {
				order = DESC
			}

			column := &ColumnMetadata{
				Keyspace:       table.Keyspace,
				Table:          table.Name,
				Name:           alias,
				Type:           comparatorParsed.types[i],
				Order:          order,
				Kind:           ColumnClusteringKey,
				ComponentIndex: i,
			}

			table.ClusteringColumns[i] = column
			table.Columns[alias] = column
		}

		if size != len(comparatorParsed.types)-1 {
			alias := DEFAULT_VALUE_ALIAS
			if len(table.ValueAlias) > 0 {
				alias = table.ValueAlias
			}
			// decode the default validator
			defaultValidatorParsed := parseType(table.DefaultValidator, logger)
			column := &ColumnMetadata{
				Keyspace: table.Keyspace,
				Table:    table.Name,
				Name:     alias,
				Type:     defaultValidatorParsed.types[0],
				Kind:     ColumnRegular,
			}
			table.Columns[alias] = column
		}
	}
}

// The simpler compile case for V2+ protocol
func compileV2Metadata(tables []TableMetadata, logger StdLogger) {
	for i := range tables {
		table := &tables[i]

		clusteringColumnCount := componentColumnCountOfType(table.Columns, ColumnClusteringKey)
		table.ClusteringColumns = make([]*ColumnMetadata, clusteringColumnCount)

		if table.KeyValidator != "" {
			keyValidatorParsed := parseType(table.KeyValidator, logger)
			table.PartitionKey = make([]*ColumnMetadata, len(keyValidatorParsed.types))
		} else { // Cassandra 3.x+
			partitionKeyCount := componentColumnCountOfType(table.Columns, ColumnPartitionKey)
			table.PartitionKey = make([]*ColumnMetadata, partitionKeyCount)
		}

		for _, columnName := range table.OrderedColumns {
			column := table.Columns[columnName]
			if column.Kind == ColumnPartitionKey {
				table.PartitionKey[column.ComponentIndex] = column
			} else if column.Kind == ColumnClusteringKey {
				table.ClusteringColumns[column.ComponentIndex] = column
			}
		}
	}
}

// returns the count of coluns with the given "kind" value.
func componentColumnCountOfType(columns map[string]*ColumnMetadata, kind ColumnKind) int {
	maxComponentIndex := -1
	for _, column := range columns {
		if column.Kind == kind && column.ComponentIndex > maxComponentIndex {
			maxComponentIndex = column.ComponentIndex
		}
	}
	return maxComponentIndex + 1
}

// query only for the keyspace metadata for the specified keyspace from system.schema_keyspace
func getKeyspaceMetadata(session *Session, keyspaceName string) (*KeyspaceMetadata, error) {
	keyspace := &KeyspaceMetadata{Name: keyspaceName}

	if session.useSystemSchema { // Cassandra 3.x+
		const stmt = `
		SELECT durable_writes, replication
		FROM system_schema.keyspaces
		WHERE keyspace_name = ?`

		var replication map[string]string

		iter := session.control.query(stmt, keyspaceName)
		if iter.NumRows() == 0 {
			return nil, ErrKeyspaceDoesNotExist
		}
		iter.Scan(&keyspace.DurableWrites, &replication)
		err := iter.Close()
		if err != nil {
			return nil, fmt.Errorf("error querying keyspace schema: %v", err)
		}

		keyspace.StrategyClass = replication["class"]
		delete(replication, "class")

		keyspace.StrategyOptions = make(map[string]interface{}, len(replication))
		for k, v := range replication {
			keyspace.StrategyOptions[k] = v
		}
	} else {

		const stmt = `
		SELECT durable_writes, strategy_class, strategy_options
		FROM system.schema_keyspaces
		WHERE keyspace_name = ?`

		var strategyOptionsJSON []byte

		iter := session.control.query(stmt, keyspaceName)
		if iter.NumRows() == 0 {
			return nil, ErrKeyspaceDoesNotExist
		}
		iter.Scan(&keyspace.DurableWrites, &keyspace.StrategyClass, &strategyOptionsJSON)
		err := iter.Close()
		if err != nil {
			return nil, fmt.Errorf("error querying keyspace schema: %v", err)
		}

		err = json.Unmarshal(strategyOptionsJSON, &keyspace.StrategyOptions)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid JSON value '%s' as strategy_options for in keyspace '%s': %v",
				strategyOptionsJSON, keyspace.Name, err,
			)
		}
	}

	return keyspace, nil
}

// query for only the table metadata in the specified keyspace from system.schema_columnfamilies
func getTableMetadata(session *Session, keyspaceName string) ([]TableMetadata, error) {

	var (
		iter *Iter
		scan func(iter *Iter, table *TableMetadata) bool
		stmt string

		keyAliasesJSON    []byte
		columnAliasesJSON []byte
	)

	if session.useSystemSchema { // Cassandra 3.x+
		stmt = `
		SELECT
			table_name
		FROM system_schema.tables
		WHERE keyspace_name = ?`

		switchIter := func() *Iter {
			iter.Close()
			stmt = `
				SELECT
					view_name
				FROM system_schema.views
				WHERE keyspace_name = ?`
			iter = session.control.query(stmt, keyspaceName)
			return iter
		}

		scan = func(iter *Iter, table *TableMetadata) bool {
			r := iter.Scan(
				&table.Name,
			)
			if !r {
				iter = switchIter()
				if iter != nil {
					switchIter = func() *Iter { return nil }
					r = iter.Scan(&table.Name)
				}
			}
			return r
		}
	} else if session.cfg.ProtoVersion == protoVersion1 {
		// we have key aliases
		stmt = `
		SELECT
			columnfamily_name,
			key_validator,
			comparator,
			default_validator,
			key_aliases,
			column_aliases,
			value_alias
		FROM system.schema_columnfamilies
		WHERE keyspace_name = ?`

		scan = func(iter *Iter, table *TableMetadata) bool {
			return iter.Scan(
				&table.Name,
				&table.KeyValidator,
				&table.Comparator,
				&table.DefaultValidator,
				&keyAliasesJSON,
				&columnAliasesJSON,
				&table.ValueAlias,
			)
		}
	} else {
		stmt = `
		SELECT
			columnfamily_name,
			key_validator,
			comparator,
			default_validator
		FROM system.schema_columnfamilies
		WHERE keyspace_name = ?`

		scan = func(iter *Iter, table *TableMetadata) bool {
			return iter.Scan(
				&table.Name,
				&table.KeyValidator,
				&table.Comparator,
				&table.DefaultValidator,
			)
		}
	}

	iter = session.control.query(stmt, keyspaceName)

	tables := []TableMetadata{}
	table := TableMetadata{Keyspace: keyspaceName}

	for scan(iter, &table) {
		var err error

		// decode the key aliases
		if keyAliasesJSON != nil {
			table.KeyAliases = []string{}
			err = json.Unmarshal(keyAliasesJSON, &table.KeyAliases)
			if err != nil {
				iter.Close()
				return nil, fmt.Errorf(
					"invalid JSON value '%s' as key_aliases for in table '%s': %v",
					keyAliasesJSON, table.Name, err,
				)
			}
		}

		// decode the column aliases
		if columnAliasesJSON != nil {
			table.ColumnAliases = []string{}
			err = json.Unmarshal(columnAliasesJSON, &table.ColumnAliases)
			if err != nil {
				iter.Close()
				return nil, fmt.Errorf(
					"invalid JSON value '%s' as column_aliases for in table '%s': %v",
					columnAliasesJSON, table.Name, err,
				)
			}
		}

		tables = append(tables, table)
		table = TableMetadata{Keyspace: keyspaceName}
	}

	err := iter.Close()
	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("error querying table schema: %v", err)
	}

	return tables, nil
}

func (s *Session) scanColumnMetadataV1(keyspace string) ([]ColumnMetadata, error) {
	// V1 does not support the type column, and all returned rows are
	// of kind "regular".
	const stmt = `
		SELECT
				columnfamily_name,
				column_name,
				component_index,
				validator,
				index_name,
				index_type,
				index_options
			FROM system.schema_columns
			WHERE keyspace_name = ?`

	var columns []ColumnMetadata

	rows := s.control.query(stmt, keyspace).Scanner()
	for rows.Next() {
		var (
			column           = ColumnMetadata{Keyspace: keyspace}
			indexOptionsJSON []byte
		)

		// all columns returned by V1 are regular
		column.Kind = ColumnRegular

		err := rows.Scan(&column.Table,
			&column.Name,
			&column.ComponentIndex,
			&column.Validator,
			&column.Index.Name,
			&column.Index.Type,
			&indexOptionsJSON)

		if err != nil {
			return nil, err
		}

		if len(indexOptionsJSON) > 0 {
			err := json.Unmarshal(indexOptionsJSON, &column.Index.Options)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid JSON value '%s' as index_options for column '%s' in table '%s': %v",
					indexOptionsJSON,
					column.Name,
					column.Table,
					err)
			}
		}

		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

func (s *Session) scanColumnMetadataV2(keyspace string) ([]ColumnMetadata, error) {
	// V2+ supports the type column
	const stmt = `
			SELECT
				columnfamily_name,
				column_name,
				component_index,
				validator,
				index_name,
				index_type,
				index_options,
				type
			FROM system.schema_columns
			WHERE keyspace_name = ?`

	var columns []ColumnMetadata

	rows := s.control.query(stmt, keyspace).Scanner()
	for rows.Next() {
		var (
			column           = ColumnMetadata{Keyspace: keyspace}
			indexOptionsJSON []byte
		)

		err := rows.Scan(&column.Table,
			&column.Name,
			&column.ComponentIndex,
			&column.Validator,
			&column.Index.Name,
			&column.Index.Type,
			&indexOptionsJSON,
			&column.Kind,
		)

		if err != nil {
			return nil, err
		}

		if len(indexOptionsJSON) > 0 {
			err := json.Unmarshal(indexOptionsJSON, &column.Index.Options)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid JSON value '%s' as index_options for column '%s' in table '%s': %v",
					indexOptionsJSON,
					column.Name,
					column.Table,
					err)
			}
		}

		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil

}

func (s *Session) scanColumnMetadataSystem(keyspace string) ([]ColumnMetadata, error) {
	const stmt = `
			SELECT
				table_name,
				column_name,
				clustering_order,
				type,
				kind,
				position
			FROM system_schema.columns
			WHERE keyspace_name = ?`

	var columns []ColumnMetadata

	rows := s.control.query(stmt, keyspace).Scanner()
	for rows.Next() {
		column := ColumnMetadata{Keyspace: keyspace}

		err := rows.Scan(&column.Table,
			&column.Name,
			&column.ClusteringOrder,
			&column.Validator,
			&column.Kind,
			&column.ComponentIndex,
		)

		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// TODO(zariel): get column index info from system_schema.indexes

	return columns, nil
}

// query for only the column metadata in the specified keyspace from system.schema_columns
func getColumnMetadata(session *Session, keyspaceName string) ([]ColumnMetadata, error) {
	var (
		columns []ColumnMetadata
		err     error
	)

	// Deal with differences in protocol versions
	if session.cfg.ProtoVersion == 1 {
		columns, err = session.scanColumnMetadataV1(keyspaceName)
	} else if session.useSystemSchema { // Cassandra 3.x+
		columns, err = session.scanColumnMetadataSystem(keyspaceName)
	} else {
		columns, err = session.scanColumnMetadataV2(keyspaceName)
	}

	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("error querying column schema: %v", err)
	}

	return columns, nil
}

func getTypeInfo(t string, logger StdLogger) TypeInfo {
	if strings.HasPrefix(t, apacheCassandraTypePrefix) {
		t = apacheToCassandraType(t)
	}
	return getCassandraType(t, logger)
}

func getViewsMetadata(session *Session, keyspaceName string) ([]ViewMetadata, error) {
	if session.cfg.ProtoVersion == protoVersion1 {
		return nil, nil
	}
	var tableName string
	if session.useSystemSchema {
		tableName = "system_schema.types"
	} else {
		tableName = "system.schema_usertypes"
	}
	stmt := fmt.Sprintf(`
		SELECT
			type_name,
			field_names,
			field_types
		FROM %s
		WHERE keyspace_name = ?`, tableName)

	var views []ViewMetadata

	rows := session.control.query(stmt, keyspaceName).Scanner()
	for rows.Next() {
		view := ViewMetadata{Keyspace: keyspaceName}
		var argumentTypes []string
		err := rows.Scan(&view.Name,
			&view.FieldNames,
			&argumentTypes,
		)
		if err != nil {
			return nil, err
		}
		view.FieldTypes = make([]TypeInfo, len(argumentTypes))
		for i, argumentType := range argumentTypes {
			view.FieldTypes[i] = getTypeInfo(argumentType, session.logger)
		}
		views = append(views, view)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return views, nil
}

func getMaterializedViewsMetadata(session *Session, keyspaceName string) ([]MaterializedViewMetadata, error) {
	if !session.useSystemSchema {
		return nil, nil
	}
	var tableName = "system_schema.views"
	stmt := fmt.Sprintf(`
		SELECT
			view_name,
			base_table_id,
			base_table_name,
			bloom_filter_fp_chance,
			caching,
			comment,
			compaction,
			compression,
			crc_check_chance,
			dclocal_read_repair_chance,
			default_time_to_live,
			extensions,
			gc_grace_seconds,
			id,
			include_all_columns,
			max_index_interval,
			memtable_flush_period_in_ms,
			min_index_interval,
			read_repair_chance,
			speculative_retry
		FROM %s
		WHERE keyspace_name = ?`, tableName)

	var materializedViews []MaterializedViewMetadata

	rows := session.control.query(stmt, keyspaceName).Scanner()
	for rows.Next() {
		materializedView := MaterializedViewMetadata{Keyspace: keyspaceName}
		err := rows.Scan(&materializedView.Name,
			&materializedView.BaseTableId,
			&materializedView.baseTableName,
			&materializedView.BloomFilterFpChance,
			&materializedView.Caching,
			&materializedView.Comment,
			&materializedView.Compaction,
			&materializedView.Compression,
			&materializedView.CrcCheckChance,
			&materializedView.DcLocalReadRepairChance,
			&materializedView.DefaultTimeToLive,
			&materializedView.Extensions,
			&materializedView.GcGraceSeconds,
			&materializedView.Id,
			&materializedView.IncludeAllColumns,
			&materializedView.MaxIndexInterval,
			&materializedView.MemtableFlushPeriodInMs,
			&materializedView.MinIndexInterval,
			&materializedView.ReadRepairChance,
			&materializedView.SpeculativeRetry,
		)
		if err != nil {
			return nil, err
		}
		materializedViews = append(materializedViews, materializedView)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materializedViews, nil
}

func getFunctionsMetadata(session *Session, keyspaceName string) ([]FunctionMetadata, error) {
	if session.cfg.ProtoVersion == protoVersion1 || !session.hasAggregatesAndFunctions {
		return nil, nil
	}
	var tableName string
	if session.useSystemSchema {
		tableName = "system_schema.functions"
	} else {
		tableName = "system.schema_functions"
	}
	stmt := fmt.Sprintf(`
		SELECT
			function_name,
			argument_types,
			argument_names,
			body,
			called_on_null_input,
			language,
			return_type
		FROM %s
		WHERE keyspace_name = ?`, tableName)

	var functions []FunctionMetadata

	rows := session.control.query(stmt, keyspaceName).Scanner()
	for rows.Next() {
		function := FunctionMetadata{Keyspace: keyspaceName}
		var argumentTypes []string
		var returnType string
		err := rows.Scan(&function.Name,
			&argumentTypes,
			&function.ArgumentNames,
			&function.Body,
			&function.CalledOnNullInput,
			&function.Language,
			&returnType,
		)
		if err != nil {
			return nil, err
		}
		function.ReturnType = getTypeInfo(returnType, session.logger)
		function.ArgumentTypes = make([]TypeInfo, len(argumentTypes))
		for i, argumentType := range argumentTypes {
			function.ArgumentTypes[i] = getTypeInfo(argumentType, session.logger)
		}
		functions = append(functions, function)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return functions, nil
}

func getAggregatesMetadata(session *Session, keyspaceName string) ([]AggregateMetadata, error) {
	if session.cfg.ProtoVersion == protoVersion1 || !session.hasAggregatesAndFunctions {
		return nil, nil
	}
	var tableName string
	if session.useSystemSchema {
		tableName = "system_schema.aggregates"
	} else {
		tableName = "system.schema_aggregates"
	}

	stmt := fmt.Sprintf(`
		SELECT
			aggregate_name,
			argument_types,
			final_func,
			initcond,
			return_type,
			state_func,
			state_type
		FROM %s
		WHERE keyspace_name = ?`, tableName)

	var aggregates []AggregateMetadata

	rows := session.control.query(stmt, keyspaceName).Scanner()
	for rows.Next() {
		aggregate := AggregateMetadata{Keyspace: keyspaceName}
		var argumentTypes []string
		var returnType string
		var stateType string
		err := rows.Scan(&aggregate.Name,
			&argumentTypes,
			&aggregate.finalFunc,
			&aggregate.InitCond,
			&returnType,
			&aggregate.stateFunc,
			&stateType,
		)
		if err != nil {
			return nil, err
		}
		aggregate.ReturnType = getTypeInfo(returnType, session.logger)
		aggregate.StateType = getTypeInfo(stateType, session.logger)
		aggregate.ArgumentTypes = make([]TypeInfo, len(argumentTypes))
		for i, argumentType := range argumentTypes {
			aggregate.ArgumentTypes[i] = getTypeInfo(argumentType, session.logger)
		}
		aggregates = append(aggregates, aggregate)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return aggregates, nil
}

// type definition parser state
type typeParser struct {
	input  string
	index  int
	logger StdLogger
}

// the type definition parser result
type typeParserResult struct {
	isComposite bool
	types       []TypeInfo
	reversed    []bool
	collections map[string]TypeInfo
}

// Parse the type definition used for validator and comparator schema data
func parseType(def string, logger StdLogger) typeParserResult {
	parser := &typeParser{input: def, logger: logger}
	return parser.parse()
}

const (
	REVERSED_TYPE   = "org.apache.cassandra.db.marshal.ReversedType"
	COMPOSITE_TYPE  = "org.apache.cassandra.db.marshal.CompositeType"
	COLLECTION_TYPE = "org.apache.cassandra.db.marshal.ColumnToCollectionType"
	LIST_TYPE       = "org.apache.cassandra.db.marshal.ListType"
	SET_TYPE        = "org.apache.cassandra.db.marshal.SetType"
	MAP_TYPE        = "org.apache.cassandra.db.marshal.MapType"
)

// represents a class specification in the type def AST
type typeParserClassNode struct {
	name   string
	params []typeParserParamNode
	// this is the segment of the input string that defined this node
	input string
}

// represents a class parameter in the type def AST
type typeParserParamNode struct {
	name  *string
	class typeParserClassNode
}

func (t *typeParser) parse() typeParserResult {
	// parse the AST
	ast, ok := t.parseClassNode()
	if !ok {
		// treat this is a custom type
		return typeParserResult{
			isComposite: false,
			types: []TypeInfo{
				NativeType{
					typ:    TypeCustom,
					custom: t.input,
				},
			},
			reversed:    []bool{false},
			collections: nil,
		}
	}

	// interpret the AST
	if strings.HasPrefix(ast.name, COMPOSITE_TYPE) {
		count := len(ast.params)

		// look for a collections param
		last := ast.params[count-1]
		collections := map[string]TypeInfo{}
		if strings.HasPrefix(last.class.name, COLLECTION_TYPE) {
			count--

			for _, param := range last.class.params {
				// decode the name
				var name string
				decoded, err := hex.DecodeString(*param.name)
				if err != nil {
					t.logger.Printf(
						"Error parsing type '%s', contains collection name '%s' with an invalid format: %v",
						t.input,
						*param.name,
						err,
					)
					// just use the provided name
					name = *param.name
				} else {
					name = string(decoded)
				}
				collections[name] = param.class.asTypeInfo()
			}
		}

		types := make([]TypeInfo, count)
		reversed := make([]bool, count)

		for i, param := range ast.params[:count] {
			class := param.class
			reversed[i] = strings.HasPrefix(class.name, REVERSED_TYPE)
			if reversed[i] {
				class = class.params[0].class
			}
			types[i] = class.asTypeInfo()
		}

		return typeParserResult{
			isComposite: true,
			types:       types,
			reversed:    reversed,
			collections: collections,
		}
	} else {
		// not composite, so one type
		class := *ast
		reversed := strings.HasPrefix(class.name, REVERSED_TYPE)
		if reversed {
			class = class.params[0].class
		}
		typeInfo := class.asTypeInfo()

		return typeParserResult{
			isComposite: false,
			types:       []TypeInfo{typeInfo},
			reversed:    []bool{reversed},
		}
	}
}

func (class *typeParserClassNode) asTypeInfo() TypeInfo {
	if strings.HasPrefix(class.name, LIST_TYPE) {
		elem := class.params[0].class.asTypeInfo()
		return CollectionType{
			NativeType: NativeType{
				typ: TypeList,
			},
			Elem: elem,
		}
	}
	if strings.HasPrefix(class.name, SET_TYPE) {
		elem := class.params[0].class.asTypeInfo()
		return CollectionType{
			NativeType: NativeType{
				typ: TypeSet,
			},
			Elem: elem,
		}
	}
	if strings.HasPrefix(class.name, MAP_TYPE) {
		key := class.params[0].class.asTypeInfo()
		elem := class.params[1].class.asTypeInfo()
		return CollectionType{
			NativeType: NativeType{
				typ: TypeMap,
			},
			Key:  key,
			Elem: elem,
		}
	}

	// must be a simple type or custom type
	info := NativeType{typ: getApacheCassandraType(class.name)}
	if info.typ == TypeCustom {
		// add the entire class definition
		info.custom = class.input
	}
	return info
}

// CLASS := ID [ PARAMS ]
func (t *typeParser) parseClassNode() (node *typeParserClassNode, ok bool) {
	t.skipWhitespace()

	startIndex := t.index

	name, ok := t.nextIdentifier()
	if !ok {
		return nil, false
	}

	params, ok := t.parseParamNodes()
	if !ok {
		return nil, false
	}

	endIndex := t.index

	node = &typeParserClassNode{
		name:   name,
		params: params,
		input:  t.input[startIndex:endIndex],
	}
	return node, true
}

// PARAMS := "(" PARAM { "," PARAM } ")"
// PARAM := [ PARAM_NAME ":" ] CLASS
// PARAM_NAME := ID
func (t *typeParser) parseParamNodes() (params []typeParserParamNode, ok bool) {
	t.skipWhitespace()

	// the params are optional
	if t.index == len(t.input) || t.input[t.index] != '(' {
		return nil, true
	}

	params = []typeParserParamNode{}

	// consume the '('
	t.index++

	t.skipWhitespace()

	for t.input[t.index] != ')' {
		// look for a named param, but if no colon, then we want to backup
		backupIndex := t.index

		// name will be a hex encoded version of a utf-8 string
		name, ok := t.nextIdentifier()
		if !ok {
			return nil, false
		}
		hasName := true

		// TODO handle '=>' used for DynamicCompositeType

		t.skipWhitespace()

		if t.input[t.index] == ':' {
			// there is a name for this parameter

			// consume the ':'
			t.index++

			t.skipWhitespace()
		} else {
			// no name, backup
			hasName = false
			t.index = backupIndex
		}

		// parse the next full parameter
		classNode, ok := t.parseClassNode()
		if !ok {
			return nil, false
		}

		if hasName {
			params = append(
				params,
				typeParserParamNode{name: &name, class: *classNode},
			)
		} else {
			params = append(
				params,
				typeParserParamNode{class: *classNode},
			)
		}

		t.skipWhitespace()

		if t.input[t.index] == ',' {
			// consume the comma
			t.index++

			t.skipWhitespace()
		}
	}

	// consume the ')'
	t.index++

	return params, true
}

func (t *typeParser) skipWhitespace() {
	for t.index < len(t.input) && isWhitespaceChar(t.input[t.index]) {
		t.index++
	}
}

func isWhitespaceChar(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t'
}

// ID := LETTER { LETTER }
// LETTER := "0"..."9" | "a"..."z" | "A"..."Z" | "-" | "+" | "." | "_" | "&"
func (t *typeParser) nextIdentifier() (id string, found bool) {
	startIndex := t.index
	for t.index < len(t.input) && isIdentifierChar(t.input[t.index]) {
		t.index++
	}
	if startIndex == t.index {
		return "", false
	}
	return t.input[startIndex:t.index], true
}

func isIdentifierChar(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '-' ||
		c == '+' ||
		c == '.' ||
		c == '_' ||
		c == '&'
}
