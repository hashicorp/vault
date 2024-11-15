package driver

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/SAP/go-hdb/driver/internal/unsafe"
)

const sqlTagKey = "sql"

// parseSQLTag return name, type and options.
func parseSQLTag(tag string) (string, string, sqlTagOptions) {
	name, rest, _ := strings.Cut(tag, ",")
	typ, opts, _ := strings.Cut(rest, ",")
	return name, typ, sqlTagOptions(opts)
}

type sqlTagOptions string

func (o sqlTagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}

// Tagger is an interface used to tag structure fields dynamically.
type Tagger interface {
	Tag(fieldName string) (value string, ok bool)
}

type structColumn struct {
	fieldName  string
	fieldType  reflect.Type
	fieldIndex []int

	sqlName    string
	sqlType    string
	sqlOptions sqlTagOptions
}

func newStructColumn(name string, typ reflect.Type, index []int, tag reflect.StructTag) (*structColumn, bool) {
	c := &structColumn{fieldName: name, fieldType: typ, fieldIndex: index}
	if sqlTag, ok := tag.Lookup(sqlTagKey); ok {
		if sqlTag == "-" { // ignore field
			return nil, false
		}
		c.sqlName, c.sqlType, c.sqlOptions = parseSQLTag(sqlTag)
	}
	return c, true
}

func (c *structColumn) Name() string {
	if c.sqlName != "" {
		return c.sqlName
	}
	return c.fieldName
}

func (c *structColumn) Type() (string, error) {
	if c.sqlType == "" {
		var err error
		if c.sqlType, err = inferSQLDatatype(c.fieldType); err != nil {
			return "", err
		}
	}
	return c.sqlType, nil
}

func (c *structColumn) def() (string, error) {
	typ, err := c.Type()
	if err != nil {
		return "", err
	}
	s := Identifier(c.Name()).String() + " " + typ
	if c.sqlOptions.Contains("not null") {
		s += " not null"
	}
	return s, nil
}

type structColumns []*structColumn

func (c structColumns) defs() (string, error) {
	if len(c) == 0 {
		return "", nil
	}
	buf := []byte{'('}
	definition, err := c[0].def()
	if err != nil {
		return "", err
	}
	buf = append(buf, definition...)
	for i := 1; i < len(c); i++ {
		definition, err := c[i].def()
		if err != nil {
			return "", err
		}
		buf = append(buf, ',')
		buf = append(buf, definition...)
	}
	buf = append(buf, ')')
	return unsafe.ByteSlice2String(buf), nil
}

func (c structColumns) queryPlaceholders() string {
	if len(c) == 0 {
		return ""
	}
	buf := []byte{'(', '?'}
	for i := 1; i < len(c); i++ {
		buf = append(buf, ",?"...)
	}
	buf = append(buf, ')')
	return unsafe.ByteSlice2String(buf)
}

// StructScanner is a database scanner to scan rows into a struct of type S.
// This enables using structs as scan targets for the exported fields of the struct.
// For usage please refer to the example.
type StructScanner[S any] struct {
	columns       structColumns
	nameColumnMap map[string]*structColumn
}

// NewStructScanner returns a new struct scanner.
func NewStructScanner[S any]() (*StructScanner[S], error) {
	var s *S

	rt := reflect.TypeOf(s).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid type %s", rt.Kind())
	}

	tagger, hasTagger := any(s).(Tagger)

	columns := []*structColumn{}
	nameColumnMap := map[string]*structColumn{}

	for _, field := range reflect.VisibleFields(rt) {
		if field.IsExported() {
			fieldTag := field.Tag
			if hasTagger {
				if tag, ok := tagger.Tag(field.Name); ok {
					fieldTag = reflect.StructTag(tag)
				}
			}
			column, ok := newStructColumn(field.Name, field.Type, field.Index, fieldTag)
			if !ok {
				continue
			}
			name := column.Name()
			if _, ok := nameColumnMap[name]; ok {
				return nil, fmt.Errorf("duplicate column name %s", name)
			}
			columns = append(columns, column)
			nameColumnMap[name] = column
		}
	}
	return &StructScanner[S]{columns: columns, nameColumnMap: nameColumnMap}, nil
}

// ScanRow scans the field values of the first row in rows into struct s of type *S and closes rows.
func (sc StructScanner[S]) ScanRow(rows *sql.Rows, s *S) error {
	if rows.Err() != nil {
		return rows.Err()
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
	err := sc.Scan(rows, s)
	if err != nil {
		return err
	}
	return rows.Close()
}

// Scan scans row field values into struct s of type *S.
func (sc StructScanner[S]) Scan(rows *sql.Rows, s *S) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(s).Elem()
	values := make([]any, len(columns))
	for i, name := range columns {
		column, ok := sc.nameColumnMap[name]
		if !ok {
			return fmt.Errorf("field for column name %s not found", name)
		}
		values[i] = rv.FieldByIndex(column.fieldIndex).Addr().Interface()
	}
	return rows.Scan(values...)
}

// columnDefs returns the column definitions for a sql create statement.
// experimental: before 'export' completion of inferSQLType is needed
func (sc StructScanner[S]) columnDefs() (string, error) { return sc.columns.defs() }

func (sc StructScanner[S]) queryPlaceholders() string { return sc.columns.queryPlaceholders() }

var kindSQLDatatypes = map[reflect.Kind]string{
	reflect.Bool:    "boolean",
	reflect.Int:     "bigint",
	reflect.Int8:    "smallint",
	reflect.Int16:   "smallint",
	reflect.Int32:   "integer",
	reflect.Int64:   "bigint",
	reflect.Uint:    "bigint",
	reflect.Uint8:   "tinyint",
	reflect.Uint16:  "smallint",
	reflect.Uint32:  "integer",
	reflect.Uint64:  "bigint",
	reflect.Float32: "real",
	reflect.Float64: "double",
	reflect.String:  "nvarchar(256)",
}

var (
	decimalType     = reflect.TypeFor[Decimal]()
	lobType         = reflect.TypeFor[Lob]()
	timeType        = reflect.TypeFor[time.Time]()
	bytesType       = reflect.TypeFor[[]byte]()
	nullBoolType    = reflect.TypeFor[sql.NullBool]()
	nullByteType    = reflect.TypeFor[sql.NullByte]()
	nullFloat64Type = reflect.TypeFor[sql.NullFloat64]()
	nullInt16Type   = reflect.TypeFor[sql.NullInt16]()
	nullInt32Type   = reflect.TypeFor[sql.NullInt32]()
	nullInt64Type   = reflect.TypeFor[sql.NullInt64]()
	nullStringType  = reflect.TypeFor[sql.NullString]()
	nullTimeType    = reflect.TypeFor[sql.NullTime]()
	nullBytesType   = reflect.TypeFor[NullBytes]()
	nullDecimalType = reflect.TypeFor[NullDecimal]()
	nullLobType     = reflect.TypeFor[NullLob]()
)

var typeSQLDatatypes = map[reflect.Type]string{
	decimalType:     "decimal",
	lobType:         "blob",
	timeType:        "timestamp",
	bytesType:       "varchar(256)",
	nullBoolType:    "boolean",
	nullByteType:    "varchar",
	nullFloat64Type: "double",
	nullInt16Type:   "smallint",
	nullInt32Type:   "integer",
	nullInt64Type:   "bigint",
	nullStringType:  "nvarchar(256)",
	nullTimeType:    "timestamp",
	nullBytesType:   "varchar(256)",
	nullDecimalType: "decimal",
	nullLobType:     "blob",
}

// inferSQLDatatype tries to infer the hdb sql datatype.
func inferSQLDatatype(typ reflect.Type) (string, error) {
	kind := typ.Kind()

	if kind == reflect.Pointer {
		return inferSQLDatatype(typ.Elem())
	}

	// dedicated datatypes.
	for ctyp, sqlType := range typeSQLDatatypes {
		if typ.ConvertibleTo(ctyp) {
			return sqlType, nil
		}
	}

	// generic Null[T].
	if kind == reflect.Struct {
		// see https://github.com/golang/go/issues/54393
		if strings.HasPrefix(typ.String(), "sql.Null[") {
			if f, ok := typ.FieldByName("V"); ok {
				return inferSQLDatatype(f.Type)
			}
		}
	}

	// basic datatypes.
	if sqlType, ok := kindSQLDatatypes[kind]; ok {
		return sqlType, nil
	}
	return "", fmt.Errorf("could not infer sql type kind %s for %s", kind, typ)
}
