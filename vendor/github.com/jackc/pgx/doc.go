// Package pgx is a PostgreSQL database driver.
/*
pgx provides lower level access to PostgreSQL than the standard database/sql.
It remains as similar to the database/sql interface as possible while
providing better speed and access to PostgreSQL specific features. Import
github.com/jackc/pgx/stdlib to use pgx as a database/sql compatible driver.

Query Interface

pgx implements Query and Scan in the familiar database/sql style.

    var sum int32

    // Send the query to the server. The returned rows MUST be closed
    // before conn can be used again.
    rows, err := conn.Query("select generate_series(1,$1)", 10)
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
        return err
    }

    // No errors found - do something with sum

pgx also implements QueryRow in the same style as database/sql.

    var name string
    var weight int64
    err := conn.QueryRow("select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
    if err != nil {
        return err
    }

Use Exec to execute a query that does not return a result set.

    commandTag, err := conn.Exec("delete from widgets where id=$1", 42)
    if err != nil {
        return err
    }
    if commandTag.RowsAffected() != 1 {
        return errors.New("No row found to delete")
    }

Connection Pool

Connection pool usage is explicit and configurable. In pgx, a connection can be
created and managed directly, or a connection pool with a configurable maximum
connections can be used. The connection pool offers an after connect hook that
allows every connection to be automatically setup before being made available in
the connection pool.

It delegates methods such as QueryRow to an automatically checked out and
released connection so you can avoid manually acquiring and releasing
connections when you do not need that level of control.

    var name string
    var weight int64
    err := pool.QueryRow("select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
    if err != nil {
        return err
    }

Base Type Mapping

pgx maps between all common base types directly between Go and PostgreSQL. In
particular:

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

pgx can map nulls in two ways. The first is package pgtype provides types that
have a data field and a status field. They work in a similar fashion to
database/sql. The second is to use a pointer to a pointer.

    var foo pgtype.Varchar
    var bar *string
    err := conn.QueryRow("select foo, bar from widgets where id=$1", 42).Scan(&foo, &bar)
    if err != nil {
        return err
    }

Array Mapping

pgx maps between int16, int32, int64, float32, float64, and string Go slices
and the equivalent PostgreSQL array type. Go slices of native types do not
support nulls, so if a PostgreSQL array that contains a null is read into a
native Go slice an error will occur. The pgtype package includes many more
array types for PostgreSQL types that do not directly map to native Go types.

JSON and JSONB Mapping

pgx includes built-in support to marshal and unmarshal between Go types and
the PostgreSQL JSON and JSONB.

Inet and CIDR Mapping

pgx encodes from net.IPNet to and from inet and cidr PostgreSQL types. In
addition, as a convenience pgx will encode from a net.IP; it will assume a /32
netmask for IPv4 and a /128 for IPv6.

Custom Type Support

pgx includes support for the common data types like integers, floats, strings,
dates, and times that have direct mappings between Go and SQL. In addition,
pgx uses the github.com/jackc/pgx/pgtype library to support more types. See
documention for that library for instructions on how to implement custom
types.

See example_custom_type_test.go for an example of a custom type for the
PostgreSQL point type.

pgx also includes support for custom types implementing the database/sql.Scanner
and database/sql/driver.Valuer interfaces.

If pgx does cannot natively encode a type and that type is a renamed type (e.g.
type MyTime time.Time) pgx will attempt to encode the underlying type. While
this is usually desired behavior it can produce suprising behavior if one the
underlying type and the renamed type each implement database/sql interfaces and
the other implements pgx interfaces. It is recommended that this situation be
avoided by implementing pgx interfaces on the renamed type.

Raw Bytes Mapping

[]byte passed as arguments to Query, QueryRow, and Exec are passed unmodified
to PostgreSQL.

Transactions

Transactions are started by calling Begin or BeginEx. The BeginEx variant
can create a transaction with a specified isolation level.

    tx, err := conn.Begin()
    if err != nil {
        return err
    }
    // Rollback is safe to call even if the tx is already closed, so if
    // the tx commits successfully, this is a no-op
    defer tx.Rollback()

    _, err = tx.Exec("insert into foo(id) values (1)")
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

Copy Protocol

Use CopyFrom to efficiently insert multiple rows at a time using the PostgreSQL
copy protocol. CopyFrom accepts a CopyFromSource interface. If the data is already
in a [][]interface{} use CopyFromRows to wrap it in a CopyFromSource interface. Or
implement CopyFromSource to avoid buffering the entire data set in memory.

    rows := [][]interface{}{
        {"John", "Smith", int32(36)},
        {"Jane", "Doe", int32(29)},
    }

    copyCount, err := conn.CopyFrom(
        pgx.Identifier{"people"},
        []string{"first_name", "last_name", "age"},
        pgx.CopyFromRows(rows),
    )

CopyFrom can be faster than an insert with as few as 5 rows.

Listen and Notify

pgx can listen to the PostgreSQL notification system with the
WaitForNotification function. It takes a maximum time to wait for a
notification.

    err := conn.Listen("channelname")
    if err != nil {
        return nil
    }

    if notification, err := conn.WaitForNotification(time.Second); err != nil {
        // do something with notification
    }

TLS

The pgx ConnConfig struct has a TLSConfig field. If this field is
nil, then TLS will be disabled. If it is present, then it will be used to
configure the TLS connection. This allows total configuration of the TLS
connection.

pgx has never explicitly supported Postgres < 9.6's `ssl_renegotiation` option.
As of v3.3.0, it doesn't send `ssl_renegotiation: 0` either to support Redshift
(https://github.com/jackc/pgx/pull/476). If you need TLS Renegotiation,
consider supplying `ConnConfig.TLSConfig` with a non-zero `Renegotiation`
value and if it's not the default on your server, set `ssl_renegotiation`
via `ConnConfig.RuntimeParams`.

Logging

pgx defines a simple logger interface. Connections optionally accept a logger
that satisfies this interface. Set LogLevel to control logging verbosity.
Adapters for github.com/inconshreveable/log15, github.com/sirupsen/logrus, and
the testing log are provided in the log directory.
*/
package pgx
