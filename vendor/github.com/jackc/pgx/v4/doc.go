// Package pgx is a PostgreSQL database driver.
/*
pgx provides lower level access to PostgreSQL than the standard database/sql. It remains as similar to the database/sql
interface as possible while providing better speed and access to PostgreSQL specific features. Import
github.com/jackc/pgx/v4/stdlib to use pgx as a database/sql compatible driver.

Establishing a Connection

The primary way of establishing a connection is with `pgx.Connect`.

    conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

The database connection string can be in URL or DSN format. Both PostgreSQL settings and pgx settings can be specified
here. In addition, a config struct can be created by `ParseConfig` and modified before establishing the connection with
`ConnectConfig`.

    config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
    if err != nil {
        // ...
    }
    config.Logger = log15adapter.NewLogger(log.New("module", "pgx"))

    conn, err := pgx.ConnectConfig(context.Background(), config)

Connection Pool

`*pgx.Conn` represents a single connection to the database and is not concurrency safe. Use sub-package pgxpool for a
concurrency safe connection pool.

Query Interface

pgx implements Query and Scan in the familiar database/sql style.

    var sum int32

    // Send the query to the server. The returned rows MUST be closed
    // before conn can be used again.
    rows, err := conn.Query(context.Background(), "select generate_series(1,$1)", 10)
    if err != nil {
        return err
    }

    // rows.Close is called by rows.Next when all rows are read
    // or an error occurs in Next or Scan. So it may optionally be
    // omitted if nothing in the rows.Next loop can panic. It is
    // safe to close rows multiple times.
    defer rows.Close()

    // Iterate through the result set
    for rows.Next() {
        var n int32
        err = rows.Scan(&n)
        if err != nil {
            return err
        }
        sum += n
    }

    // Any errors encountered by rows.Next or rows.Scan will be returned here
    if rows.Err() != nil {
        return rows.Err()
    }

    // No errors found - do something with sum

pgx also implements QueryRow in the same style as database/sql.

    var name string
    var weight int64
    err := conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
    if err != nil {
        return err
    }

Use Exec to execute a query that does not return a result set.

    commandTag, err := conn.Exec(context.Background(), "delete from widgets where id=$1", 42)
    if err != nil {
        return err
    }
    if commandTag.RowsAffected() != 1 {
        return errors.New("No row found to delete")
    }

QueryFunc can be used to execute a callback function for every row. This is often easier to use than Query.

    var sum, n int32
	_, err = conn.QueryFunc(
		context.Background(),
		"select generate_series(1,$1)",
		[]interface{}{10},
		[]interface{}{&n},
		func(pgx.QueryFuncRow) error {
            sum += n
			return nil
		},
	)
	if err != nil {
		return err
	}

Base Type Mapping

pgx maps between all common base types directly between Go and PostgreSQL. In particular:

    Go           PostgreSQL
    -----------------------
    string       varchar
                 text

    // Integers are automatically be converted to any other integer type if
    // it can be done without overflow or underflow.
    int8
    int16        smallint
    int32        int
    int64        bigint
    int
    uint8
    uint16
    uint32
    uint64
    uint

    // Floats are strict and do not automatically convert like integers.
    float32      float4
    float64      float8

    time.Time   date
                timestamp
                timestamptz

    []byte      bytea


Null Mapping

pgx can map nulls in two ways. The first is package pgtype provides types that have a data field and a status field.
They work in a similar fashion to database/sql. The second is to use a pointer to a pointer.

    var foo pgtype.Varchar
    var bar *string
    err := conn.QueryRow("select foo, bar from widgets where id=$1", 42).Scan(&foo, &bar)
    if err != nil {
        return err
    }

Array Mapping

pgx maps between int16, int32, int64, float32, float64, and string Go slices and the equivalent PostgreSQL array type.
Go slices of native types do not support nulls, so if a PostgreSQL array that contains a null is read into a native Go
slice an error will occur. The pgtype package includes many more array types for PostgreSQL types that do not directly
map to native Go types.

JSON and JSONB Mapping

pgx includes built-in support to marshal and unmarshal between Go types and the PostgreSQL JSON and JSONB.

Inet and CIDR Mapping

pgx encodes from net.IPNet to and from inet and cidr PostgreSQL types. In addition, as a convenience pgx will encode
from a net.IP; it will assume a /32 netmask for IPv4 and a /128 for IPv6.

Custom Type Support

pgx includes support for the common data types like integers, floats, strings, dates, and times that have direct
mappings between Go and SQL. In addition, pgx uses the github.com/jackc/pgtype library to support more types. See
documention for that library for instructions on how to implement custom types.

See example_custom_type_test.go for an example of a custom type for the PostgreSQL point type.

pgx also includes support for custom types implementing the database/sql.Scanner and database/sql/driver.Valuer
interfaces.

If pgx does cannot natively encode a type and that type is a renamed type (e.g. type MyTime time.Time) pgx will attempt
to encode the underlying type. While this is usually desired behavior it can produce surprising behavior if one the
underlying type and the renamed type each implement database/sql interfaces and the other implements pgx interfaces. It
is recommended that this situation be avoided by implementing pgx interfaces on the renamed type.

Composite types and row values

Row values and composite types are represented as pgtype.Record (https://pkg.go.dev/github.com/jackc/pgtype?tab=doc#Record).
It is possible to get values of your custom type by implementing DecodeBinary interface. Decoding into
pgtype.Record first can simplify process by avoiding dealing with raw protocol directly.

For example:

    type MyType struct {
        a int      // NULL will cause decoding error
        b *string  // there can be NULL in this position in SQL
    }

    func (t *MyType) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
        r := pgtype.Record{
            Fields: []pgtype.Value{&pgtype.Int4{}, &pgtype.Text{}},
        }

        if err := r.DecodeBinary(ci, src); err != nil {
            return err
        }

        if r.Status != pgtype.Present {
            return errors.New("BUG: decoding should not be called on NULL value")
        }

        a := r.Fields[0].(*pgtype.Int4)
        b := r.Fields[1].(*pgtype.Text)

        // type compatibility is checked by AssignTo
        // only lossless assignments will succeed
        if err := a.AssignTo(&t.a); err != nil {
            return err
        }

        // AssignTo also deals with null value handling
        if err := b.AssignTo(&t.b); err != nil {
            return err
        }
        return nil
    }

    result := MyType{}
    err := conn.QueryRow(context.Background(), "select row(1, 'foo'::text)", pgx.QueryResultFormats{pgx.BinaryFormatCode}).Scan(&r)

Raw Bytes Mapping

[]byte passed as arguments to Query, QueryRow, and Exec are passed unmodified to PostgreSQL.

Transactions

Transactions are started by calling Begin.

    tx, err := conn.Begin(context.Background())
    if err != nil {
        return err
    }
    // Rollback is safe to call even if the tx is already closed, so if
    // the tx commits successfully, this is a no-op
    defer tx.Rollback(context.Background())

    _, err = tx.Exec(context.Background(), "insert into foo(id) values (1)")
    if err != nil {
        return err
    }

    err = tx.Commit(context.Background())
    if err != nil {
        return err
    }

The Tx returned from Begin also implements the Begin method. This can be used to implement pseudo nested transactions.
These are internally implemented with savepoints.

Use BeginTx to control the transaction mode.

BeginFunc and BeginTxFunc are variants that begin a transaction, execute a function, and commit or rollback the
transaction depending on the return value of the function. These can be simpler and less error prone to use.

    err = conn.BeginFunc(context.Background(), func(tx pgx.Tx) error {
        _, err := tx.Exec(context.Background(), "insert into foo(id) values (1)")
        return err
    })
    if err != nil {
        return err
    }

Prepared Statements

Prepared statements can be manually created with the Prepare method. However, this is rarely necessary because pgx
includes an automatic statement cache by default. Queries run through the normal Query, QueryRow, and Exec functions are
automatically prepared on first execution and the prepared statement is reused on subsequent executions. See ParseConfig
for information on how to customize or disable the statement cache.

Copy Protocol

Use CopyFrom to efficiently insert multiple rows at a time using the PostgreSQL copy protocol. CopyFrom accepts a
CopyFromSource interface. If the data is already in a [][]interface{} use CopyFromRows to wrap it in a CopyFromSource
interface. Or implement CopyFromSource to avoid buffering the entire data set in memory.

    rows := [][]interface{}{
        {"John", "Smith", int32(36)},
        {"Jane", "Doe", int32(29)},
    }

    copyCount, err := conn.CopyFrom(
        context.Background(),
        pgx.Identifier{"people"},
        []string{"first_name", "last_name", "age"},
        pgx.CopyFromRows(rows),
    )

When you already have a typed array using CopyFromSlice can be more convenient.

    rows := []User{
        {"John", "Smith", 36},
        {"Jane", "Doe", 29},
    }

    copyCount, err := conn.CopyFrom(
        context.Background(),
        pgx.Identifier{"people"},
        []string{"first_name", "last_name", "age"},
        pgx.CopyFromSlice(len(rows), func(i int) ([]interface{}, error) {
            return []interface{}{rows[i].FirstName, rows[i].LastName, rows[i].Age}, nil
        }),
    )

CopyFrom can be faster than an insert with as few as 5 rows.

Listen and Notify

pgx can listen to the PostgreSQL notification system with the `Conn.WaitForNotification` method. It blocks until a
notification is received or the context is canceled.

    _, err := conn.Exec(context.Background(), "listen channelname")
    if err != nil {
        return nil
    }

    if notification, err := conn.WaitForNotification(context.Background()); err != nil {
        // do something with notification
    }


Logging

pgx defines a simple logger interface. Connections optionally accept a logger that satisfies this interface. Set
LogLevel to control logging verbosity. Adapters for github.com/inconshreveable/log15, github.com/sirupsen/logrus,
go.uber.org/zap, github.com/rs/zerolog, and the testing log are provided in the log directory.

Lower Level PostgreSQL Functionality

pgx is implemented on top of github.com/jackc/pgconn a lower level PostgreSQL driver. The Conn.PgConn() method can be
used to access this lower layer.

PgBouncer

pgx is compatible with PgBouncer in two modes. One is when the connection has a statement cache in "describe" mode. The
other is when the connection is using the simple protocol. This can be set with the PreferSimpleProtocol config option.
*/
package pgx
