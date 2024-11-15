// Package stdlib is the compatibility layer from pgx to database/sql.
//
// A database/sql connection can be established through sql.Open.
//
//	db, err := sql.Open("pgx", "postgres://pgx_md5:secret@localhost:5432/pgx_test?sslmode=disable")
//	if err != nil {
//		return err
//	}
//
// Or from a DSN string.
//
//	db, err := sql.Open("pgx", "user=postgres password=secret host=localhost port=5432 database=pgx_test sslmode=disable")
//	if err != nil {
//		return err
//	}
//
// Or a pgx.ConnConfig can be used to set configuration not accessible via connection string. In this case the
// pgx.ConnConfig must first be registered with the driver. This registration returns a connection string which is used
// with sql.Open.
//
//	connConfig, _ := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
//	connConfig.Logger = myLogger
//	connStr := stdlib.RegisterConnConfig(connConfig)
//	db, _ := sql.Open("pgx", connStr)
//
// pgx uses standard PostgreSQL positional parameters in queries. e.g. $1, $2.
// It does not support named parameters.
//
//	db.QueryRow("select * from users where id=$1", userID)
//
// In Go 1.13 and above (*sql.Conn) Raw() can be used to get a *pgx.Conn from the standard
// database/sql.DB connection pool. This allows operations that use pgx specific functionality.
//
//	// Given db is a *sql.DB
//	conn, err := db.Conn(context.Background())
//	if err != nil {
//		// handle error from acquiring connection from DB pool
//	}
//
//	err = conn.Raw(func(driverConn interface{}) error {
//		conn := driverConn.(*stdlib.Conn).Conn() // conn is a *pgx.Conn
//		// Do pgx specific stuff with conn
//		conn.CopyFrom(...)
//		return nil
//	})
//	if err != nil {
//		// handle error that occurred while using *pgx.Conn
//	}
package stdlib

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Only intrinsic types should be binary format with database/sql.
var databaseSQLResultFormats pgx.QueryResultFormatsByOID

var pgxDriver *Driver

type ctxKey int

var ctxKeyFakeTx ctxKey = 0

var ErrNotPgx = errors.New("not pgx *sql.DB")

func init() {
	pgxDriver = &Driver{
		configs: make(map[string]*pgx.ConnConfig),
	}
	fakeTxConns = make(map[*pgx.Conn]*sql.Tx)

	// if pgx driver was already registered by different pgx major version then we
	// skip registration under the default name.
	if !contains(sql.Drivers(), "pgx") {
		sql.Register("pgx", pgxDriver)
	}
	sql.Register("pgx/v4", pgxDriver)

	databaseSQLResultFormats = pgx.QueryResultFormatsByOID{
		pgtype.BoolOID:        1,
		pgtype.ByteaOID:       1,
		pgtype.CIDOID:         1,
		pgtype.DateOID:        1,
		pgtype.Float4OID:      1,
		pgtype.Float8OID:      1,
		pgtype.Int2OID:        1,
		pgtype.Int4OID:        1,
		pgtype.Int8OID:        1,
		pgtype.OIDOID:         1,
		pgtype.TimestampOID:   1,
		pgtype.TimestamptzOID: 1,
		pgtype.XIDOID:         1,
	}
}

// TODO replace by slices.Contains when experimental package will be merged to stdlib
// https://pkg.go.dev/golang.org/x/exp/slices#Contains
func contains(list []string, y string) bool {
	for _, x := range list {
		if x == y {
			return true
		}
	}
	return false
}

var (
	fakeTxMutex sync.Mutex
	fakeTxConns map[*pgx.Conn]*sql.Tx
)

// OptionOpenDB options for configuring the driver when opening a new db pool.
type OptionOpenDB func(*connector)

// OptionBeforeConnect provides a callback for before connect. It is passed a shallow copy of the ConnConfig that will
// be used to connect, so only its immediate members should be modified.
func OptionBeforeConnect(bc func(context.Context, *pgx.ConnConfig) error) OptionOpenDB {
	return func(dc *connector) {
		dc.BeforeConnect = bc
	}
}

// OptionAfterConnect provides a callback for after connect.
func OptionAfterConnect(ac func(context.Context, *pgx.Conn) error) OptionOpenDB {
	return func(dc *connector) {
		dc.AfterConnect = ac
	}
}

// OptionResetSession provides a callback that can be used to add custom logic prior to executing a query on the
// connection if the connection has been used before.
// If ResetSessionFunc returns ErrBadConn error the connection will be discarded.
func OptionResetSession(rs func(context.Context, *pgx.Conn) error) OptionOpenDB {
	return func(dc *connector) {
		dc.ResetSession = rs
	}
}

// RandomizeHostOrderFunc is a BeforeConnect hook that randomizes the host order in the provided connConfig, so that a
// new host becomes primary each time. This is useful to distribute connections for multi-master databases like
// CockroachDB. If you use this you likely should set https://golang.org/pkg/database/sql/#DB.SetConnMaxLifetime as well
// to ensure that connections are periodically rebalanced across your nodes.
func RandomizeHostOrderFunc(ctx context.Context, connConfig *pgx.ConnConfig) error {
	if len(connConfig.Fallbacks) == 0 {
		return nil
	}

	newFallbacks := append([]*pgconn.FallbackConfig{&pgconn.FallbackConfig{
		Host:      connConfig.Host,
		Port:      connConfig.Port,
		TLSConfig: connConfig.TLSConfig,
	}}, connConfig.Fallbacks...)

	rand.Shuffle(len(newFallbacks), func(i, j int) {
		newFallbacks[i], newFallbacks[j] = newFallbacks[j], newFallbacks[i]
	})

	// Use the one that sorted last as the primary and keep the rest as the fallbacks
	newPrimary := newFallbacks[len(newFallbacks)-1]
	connConfig.Host = newPrimary.Host
	connConfig.Port = newPrimary.Port
	connConfig.TLSConfig = newPrimary.TLSConfig
	connConfig.Fallbacks = newFallbacks[:len(newFallbacks)-1]
	return nil
}

func GetConnector(config pgx.ConnConfig, opts ...OptionOpenDB) driver.Connector {
	c := connector{
		ConnConfig:    config,
		BeforeConnect: func(context.Context, *pgx.ConnConfig) error { return nil }, // noop before connect by default
		AfterConnect:  func(context.Context, *pgx.Conn) error { return nil },       // noop after connect by default
		ResetSession:  func(context.Context, *pgx.Conn) error { return nil },       // noop reset session by default
		driver:        pgxDriver,
	}

	for _, opt := range opts {
		opt(&c)
	}
	return c
}

func OpenDB(config pgx.ConnConfig, opts ...OptionOpenDB) *sql.DB {
	c := GetConnector(config, opts...)
	return sql.OpenDB(c)
}

type connector struct {
	pgx.ConnConfig
	BeforeConnect func(context.Context, *pgx.ConnConfig) error // function to call before creation of every new connection
	AfterConnect  func(context.Context, *pgx.Conn) error       // function to call after creation of every new connection
	ResetSession  func(context.Context, *pgx.Conn) error       // function is called before a connection is reused
	driver        *Driver
}

// Connect implement driver.Connector interface
func (c connector) Connect(ctx context.Context) (driver.Conn, error) {
	var (
		err  error
		conn *pgx.Conn
	)

	// Create a shallow copy of the config, so that BeforeConnect can safely modify it
	connConfig := c.ConnConfig
	if err = c.BeforeConnect(ctx, &connConfig); err != nil {
		return nil, err
	}

	if conn, err = pgx.ConnectConfig(ctx, &connConfig); err != nil {
		return nil, err
	}

	if err = c.AfterConnect(ctx, conn); err != nil {
		return nil, err
	}

	return &Conn{conn: conn, driver: c.driver, connConfig: connConfig, resetSessionFunc: c.ResetSession}, nil
}

// Driver implement driver.Connector interface
func (c connector) Driver() driver.Driver {
	return c.driver
}

// GetDefaultDriver returns the driver initialized in the init function
// and used when the pgx driver is registered.
func GetDefaultDriver() driver.Driver {
	return pgxDriver
}

type Driver struct {
	configMutex sync.Mutex
	configs     map[string]*pgx.ConnConfig
	sequence    int
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Ensure eventual timeout
	defer cancel()

	connector, err := d.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return connector.Connect(ctx)
}

func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	return &driverConnector{driver: d, name: name}, nil
}

func (d *Driver) registerConnConfig(c *pgx.ConnConfig) string {
	d.configMutex.Lock()
	connStr := fmt.Sprintf("registeredConnConfig%d", d.sequence)
	d.sequence++
	d.configs[connStr] = c
	d.configMutex.Unlock()
	return connStr
}

func (d *Driver) unregisterConnConfig(connStr string) {
	d.configMutex.Lock()
	delete(d.configs, connStr)
	d.configMutex.Unlock()
}

type driverConnector struct {
	driver *Driver
	name   string
}

func (dc *driverConnector) Connect(ctx context.Context) (driver.Conn, error) {
	var connConfig *pgx.ConnConfig

	dc.driver.configMutex.Lock()
	connConfig = dc.driver.configs[dc.name]
	dc.driver.configMutex.Unlock()

	if connConfig == nil {
		var err error
		connConfig, err = pgx.ParseConfig(dc.name)
		if err != nil {
			return nil, err
		}
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		conn:             conn,
		driver:           dc.driver,
		connConfig:       *connConfig,
		resetSessionFunc: func(context.Context, *pgx.Conn) error { return nil },
	}

	return c, nil
}

func (dc *driverConnector) Driver() driver.Driver {
	return dc.driver
}

// RegisterConnConfig registers a ConnConfig and returns the connection string to use with Open.
func RegisterConnConfig(c *pgx.ConnConfig) string {
	return pgxDriver.registerConnConfig(c)
}

// UnregisterConnConfig removes the ConnConfig registration for connStr.
func UnregisterConnConfig(connStr string) {
	pgxDriver.unregisterConnConfig(connStr)
}

type Conn struct {
	conn             *pgx.Conn
	psCount          int64 // Counter used for creating unique prepared statement names
	driver           *Driver
	connConfig       pgx.ConnConfig
	resetSessionFunc func(context.Context, *pgx.Conn) error // Function is called before a connection is reused
}

// Conn returns the underlying *pgx.Conn
func (c *Conn) Conn() *pgx.Conn {
	return c.conn
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if c.conn.IsClosed() {
		return nil, driver.ErrBadConn
	}

	name := fmt.Sprintf("pgx_%d", c.psCount)
	c.psCount++

	sd, err := c.conn.Prepare(ctx, name, query)
	if err != nil {
		return nil, err
	}

	return &Stmt{sd: sd, conn: c}, nil
}

func (c *Conn) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return c.conn.Close(ctx)
}

func (c *Conn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if c.conn.IsClosed() {
		return nil, driver.ErrBadConn
	}

	if pconn, ok := ctx.Value(ctxKeyFakeTx).(**pgx.Conn); ok {
		*pconn = c.conn
		return fakeTx{}, nil
	}

	var pgxOpts pgx.TxOptions
	switch sql.IsolationLevel(opts.Isolation) {
	case sql.LevelDefault:
	case sql.LevelReadUncommitted:
		pgxOpts.IsoLevel = pgx.ReadUncommitted
	case sql.LevelReadCommitted:
		pgxOpts.IsoLevel = pgx.ReadCommitted
	case sql.LevelRepeatableRead, sql.LevelSnapshot:
		pgxOpts.IsoLevel = pgx.RepeatableRead
	case sql.LevelSerializable:
		pgxOpts.IsoLevel = pgx.Serializable
	default:
		return nil, fmt.Errorf("unsupported isolation: %v", opts.Isolation)
	}

	if opts.ReadOnly {
		pgxOpts.AccessMode = pgx.ReadOnly
	}

	tx, err := c.conn.BeginTx(ctx, pgxOpts)
	if err != nil {
		return nil, err
	}

	return wrapTx{ctx: ctx, tx: tx}, nil
}

func (c *Conn) ExecContext(ctx context.Context, query string, argsV []driver.NamedValue) (driver.Result, error) {
	if c.conn.IsClosed() {
		return nil, driver.ErrBadConn
	}

	args := namedValueToInterface(argsV)

	commandTag, err := c.conn.Exec(ctx, query, args...)
	// if we got a network error before we had a chance to send the query, retry
	if err != nil {
		if pgconn.SafeToRetry(err) {
			return nil, driver.ErrBadConn
		}
	}
	return driver.RowsAffected(commandTag.RowsAffected()), err
}

func (c *Conn) QueryContext(ctx context.Context, query string, argsV []driver.NamedValue) (driver.Rows, error) {
	if c.conn.IsClosed() {
		return nil, driver.ErrBadConn
	}

	args := []interface{}{databaseSQLResultFormats}
	args = append(args, namedValueToInterface(argsV)...)

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		if pgconn.SafeToRetry(err) {
			return nil, driver.ErrBadConn
		}
		return nil, err
	}

	// Preload first row because otherwise we won't know what columns are available when database/sql asks.
	more := rows.Next()
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	return &Rows{conn: c, rows: rows, skipNext: true, skipNextMore: more}, nil
}

func (c *Conn) Ping(ctx context.Context) error {
	if c.conn.IsClosed() {
		return driver.ErrBadConn
	}

	err := c.conn.Ping(ctx)
	if err != nil {
		// A Ping failure implies some sort of fatal state. The connection is almost certainly already closed by the
		// failure, but manually close it just to be sure.
		c.Close()
		return driver.ErrBadConn
	}

	return nil
}

func (c *Conn) CheckNamedValue(*driver.NamedValue) error {
	// Underlying pgx supports sql.Scanner and driver.Valuer interfaces natively. So everything can be passed through directly.
	return nil
}

func (c *Conn) ResetSession(ctx context.Context) error {
	if c.conn.IsClosed() {
		return driver.ErrBadConn
	}

	return c.resetSessionFunc(ctx, c.conn)
}

type Stmt struct {
	sd   *pgconn.StatementDescription
	conn *Conn
}

func (s *Stmt) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.conn.conn.Deallocate(ctx, s.sd.Name)
}

func (s *Stmt) NumInput() int {
	return len(s.sd.ParamOIDs)
}

func (s *Stmt) Exec(argsV []driver.Value) (driver.Result, error) {
	return nil, errors.New("Stmt.Exec deprecated and not implemented")
}

func (s *Stmt) ExecContext(ctx context.Context, argsV []driver.NamedValue) (driver.Result, error) {
	return s.conn.ExecContext(ctx, s.sd.Name, argsV)
}

func (s *Stmt) Query(argsV []driver.Value) (driver.Rows, error) {
	return nil, errors.New("Stmt.Query deprecated and not implemented")
}

func (s *Stmt) QueryContext(ctx context.Context, argsV []driver.NamedValue) (driver.Rows, error) {
	return s.conn.QueryContext(ctx, s.sd.Name, argsV)
}

type rowValueFunc func(src []byte) (driver.Value, error)

type Rows struct {
	conn         *Conn
	rows         pgx.Rows
	valueFuncs   []rowValueFunc
	skipNext     bool
	skipNextMore bool

	columnNames []string
}

func (r *Rows) Columns() []string {
	if r.columnNames == nil {
		fields := r.rows.FieldDescriptions()
		r.columnNames = make([]string, len(fields))
		for i, fd := range fields {
			r.columnNames[i] = string(fd.Name)
		}
	}

	return r.columnNames
}

// ColumnTypeDatabaseTypeName returns the database system type name. If the name is unknown the OID is returned.
func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	if dt, ok := r.conn.conn.ConnInfo().DataTypeForOID(r.rows.FieldDescriptions()[index].DataTypeOID); ok {
		return strings.ToUpper(dt.Name)
	}

	return strconv.FormatInt(int64(r.rows.FieldDescriptions()[index].DataTypeOID), 10)
}

const varHeaderSize = 4

// ColumnTypeLength returns the length of the column type if the column is a
// variable length type. If the column is not a variable length type ok
// should return false.
func (r *Rows) ColumnTypeLength(index int) (int64, bool) {
	fd := r.rows.FieldDescriptions()[index]

	switch fd.DataTypeOID {
	case pgtype.TextOID, pgtype.ByteaOID:
		return math.MaxInt64, true
	case pgtype.VarcharOID, pgtype.BPCharArrayOID:
		return int64(fd.TypeModifier - varHeaderSize), true
	default:
		return 0, false
	}
}

// ColumnTypePrecisionScale should return the precision and scale for decimal
// types. If not applicable, ok should be false.
func (r *Rows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	fd := r.rows.FieldDescriptions()[index]

	switch fd.DataTypeOID {
	case pgtype.NumericOID:
		mod := fd.TypeModifier - varHeaderSize
		precision = int64((mod >> 16) & 0xffff)
		scale = int64(mod & 0xffff)
		return precision, scale, true
	default:
		return 0, 0, false
	}
}

// ColumnTypeScanType returns the value type that can be used to scan types into.
func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	fd := r.rows.FieldDescriptions()[index]

	switch fd.DataTypeOID {
	case pgtype.Float8OID:
		return reflect.TypeOf(float64(0))
	case pgtype.Float4OID:
		return reflect.TypeOf(float32(0))
	case pgtype.Int8OID:
		return reflect.TypeOf(int64(0))
	case pgtype.Int4OID:
		return reflect.TypeOf(int32(0))
	case pgtype.Int2OID:
		return reflect.TypeOf(int16(0))
	case pgtype.BoolOID:
		return reflect.TypeOf(false)
	case pgtype.NumericOID:
		return reflect.TypeOf(float64(0))
	case pgtype.DateOID, pgtype.TimestampOID, pgtype.TimestamptzOID:
		return reflect.TypeOf(time.Time{})
	case pgtype.ByteaOID:
		return reflect.TypeOf([]byte(nil))
	default:
		return reflect.TypeOf("")
	}
}

func (r *Rows) Close() error {
	r.rows.Close()
	return r.rows.Err()
}

func (r *Rows) Next(dest []driver.Value) error {
	ci := r.conn.conn.ConnInfo()
	fieldDescriptions := r.rows.FieldDescriptions()

	if r.valueFuncs == nil {
		r.valueFuncs = make([]rowValueFunc, len(fieldDescriptions))

		for i, fd := range fieldDescriptions {
			dataTypeOID := fd.DataTypeOID
			format := fd.Format

			switch fd.DataTypeOID {
			case pgtype.BoolOID:
				var d bool
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return d, err
				}
			case pgtype.ByteaOID:
				var d []byte
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return d, err
				}
			case pgtype.CIDOID:
				var d pgtype.CID
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			case pgtype.DateOID:
				var d pgtype.Date
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			case pgtype.Float4OID:
				var d float32
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return float64(d), err
				}
			case pgtype.Float8OID:
				var d float64
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return d, err
				}
			case pgtype.Int2OID:
				var d int16
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return int64(d), err
				}
			case pgtype.Int4OID:
				var d int32
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return int64(d), err
				}
			case pgtype.Int8OID:
				var d int64
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return d, err
				}
			case pgtype.JSONOID:
				var d pgtype.JSON
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			case pgtype.JSONBOID:
				var d pgtype.JSONB
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			case pgtype.OIDOID:
				var d pgtype.OIDValue
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			case pgtype.TimestampOID:
				var d pgtype.Timestamp
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			case pgtype.TimestamptzOID:
				var d pgtype.Timestamptz
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			case pgtype.XIDOID:
				var d pgtype.XID
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					if err != nil {
						return nil, err
					}
					return d.Value()
				}
			default:
				var d string
				scanPlan := ci.PlanScan(dataTypeOID, format, &d)
				r.valueFuncs[i] = func(src []byte) (driver.Value, error) {
					err := scanPlan.Scan(ci, dataTypeOID, format, src, &d)
					return d, err
				}
			}
		}
	}

	var more bool
	if r.skipNext {
		more = r.skipNextMore
		r.skipNext = false
	} else {
		more = r.rows.Next()
	}

	if !more {
		if r.rows.Err() == nil {
			return io.EOF
		} else {
			return r.rows.Err()
		}
	}

	for i, rv := range r.rows.RawValues() {
		if rv != nil {
			var err error
			dest[i], err = r.valueFuncs[i](rv)
			if err != nil {
				return fmt.Errorf("convert field %d failed: %v", i, err)
			}
		} else {
			dest[i] = nil
		}
	}

	return nil
}

func valueToInterface(argsV []driver.Value) []interface{} {
	args := make([]interface{}, 0, len(argsV))
	for _, v := range argsV {
		if v != nil {
			args = append(args, v.(interface{}))
		} else {
			args = append(args, nil)
		}
	}
	return args
}

func namedValueToInterface(argsV []driver.NamedValue) []interface{} {
	args := make([]interface{}, 0, len(argsV))
	for _, v := range argsV {
		if v.Value != nil {
			args = append(args, v.Value.(interface{}))
		} else {
			args = append(args, nil)
		}
	}
	return args
}

type wrapTx struct {
	ctx context.Context
	tx  pgx.Tx
}

func (wtx wrapTx) Commit() error { return wtx.tx.Commit(wtx.ctx) }

func (wtx wrapTx) Rollback() error { return wtx.tx.Rollback(wtx.ctx) }

type fakeTx struct{}

func (fakeTx) Commit() error { return nil }

func (fakeTx) Rollback() error { return nil }

// AcquireConn acquires a *pgx.Conn from database/sql connection pool. It must be released with ReleaseConn.
//
// In Go 1.13 this functionality has been incorporated into the standard library in the db.Conn.Raw() method.
func AcquireConn(db *sql.DB) (*pgx.Conn, error) {
	var conn *pgx.Conn
	ctx := context.WithValue(context.Background(), ctxKeyFakeTx, &conn)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		tx.Rollback()
		return nil, ErrNotPgx
	}

	fakeTxMutex.Lock()
	fakeTxConns[conn] = tx
	fakeTxMutex.Unlock()

	return conn, nil
}

// ReleaseConn releases a *pgx.Conn acquired with AcquireConn.
func ReleaseConn(db *sql.DB, conn *pgx.Conn) error {
	var tx *sql.Tx
	var ok bool

	if conn.PgConn().IsBusy() || conn.PgConn().TxStatus() != 'I' {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		conn.Close(ctx)
	}

	fakeTxMutex.Lock()
	tx, ok = fakeTxConns[conn]
	if ok {
		delete(fakeTxConns, conn)
		fakeTxMutex.Unlock()
	} else {
		fakeTxMutex.Unlock()
		return fmt.Errorf("can't release conn that is not acquired")
	}

	return tx.Rollback()
}
