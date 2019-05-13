[![](https://godoc.org/github.com/jackc/pgx?status.svg)](https://godoc.org/github.com/jackc/pgx)
[![Build Status](https://travis-ci.org/jackc/pgx.svg)](https://travis-ci.org/jackc/pgx)

# pgx - PostgreSQL Driver and Toolkit

pgx is a pure Go driver and toolkit for PostgreSQL. pgx is different from other drivers such as [pq](http://godoc.org/github.com/lib/pq) because, while it can operate as a database/sql compatible driver, pgx is also usable directly. It offers a native interface similar to database/sql that offers better performance and more features.


```go
var name string
var weight int64
err := conn.QueryRow("select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
if err != nil {
    return err
}
```

## Features

pgx supports many additional features beyond what is available through database/sql.

* Support for approximately 60 different PostgreSQL types
* Batch queries
* Single-round trip query mode
* Full TLS connection control
* Binary format support for custom types (can be much faster)
* Copy protocol support for faster bulk data loads
* Extendable logging support including built-in support for log15 and logrus
* Connection pool with after connect hook to do arbitrary connection setup
* Listen / notify
* PostgreSQL array to Go slice mapping for integers, floats, and strings
* Hstore support
* JSON and JSONB support
* Maps inet and cidr PostgreSQL types to net.IPNet and net.IP
* Large object support
* NULL mapping to Null* struct or pointer to pointer.
* Supports database/sql.Scanner and database/sql/driver.Valuer interfaces for custom types
* Logical replication connections, including receiving WAL and sending standby status updates
* Notice response handling (this is different than listen / notify)

## Performance

pgx performs roughly equivalent to [go-pg](https://github.com/go-pg/pg) and is almost always faster than [pq](http://godoc.org/github.com/lib/pq). When parsing large result sets the percentage difference can be significant (16483 queries/sec for pgx vs. 10106 queries/sec for pq -- 63% faster).

In many use cases a significant cause of latency is network round trips between the application and the server. pgx supports query batching to bundle multiple queries into a single round trip. Even in the case of a connection with the lowest possible latency, a local Unix domain socket, batching as few as three queries together can yield an improvement of 57%. With a typical network connection the results can be even more substantial.

See this [gist](https://gist.github.com/jackc/4996e8648a0c59839bff644f49d6e434) for the underlying benchmark results or checkout [go_db_bench](https://github.com/jackc/go_db_bench) to run tests for yourself.

In addition to the native driver, pgx also includes a number of packages that provide additional functionality.

## github.com/jackc/pgx/stdlib

database/sql compatibility layer for pgx. pgx can be used as a normal database/sql driver, but at any time the native interface may be acquired for more performance or PostgreSQL specific functionality.

## github.com/jackc/pgx/pgtype

Approximately 60 PostgreSQL types are supported including uuid, hstore, json, bytea, numeric, interval, inet, and arrays. These types support database/sql interfaces and are usable even outside of pgx. They are fully tested in pgx and pq. They also support a higher performance interface when used with the pgx driver.

## github.com/jackc/pgx/pgproto3

pgproto3 provides standalone encoding and decoding of the PostgreSQL v3 wire protocol. This is useful for implementing very low level PostgreSQL tooling.

## github.com/jackc/pgx/pgmock

pgmock offers the ability to create a server that mocks the PostgreSQL wire protocol. This is used internally to test pgx by purposely inducing unusual errors. pgproto3 and pgmock together provide most of the foundational tooling required to implement a PostgreSQL proxy or MitM (such as for a custom connection pooler).

## Documentation

pgx includes extensive documentation in the godoc format. It is viewable online at [godoc.org](https://godoc.org/github.com/jackc/pgx).

## Testing

pgx supports multiple connection and authentication types. Setting up a test
environment that can test all of them can be cumbersome. In particular,
Windows cannot test Unix domain socket connections. Because of this pgx will
skip tests for connection types that are not configured.

### Normal Test Environment

To setup the normal test environment, first install these dependencies:

    go get github.com/cockroachdb/apd
    go get github.com/hashicorp/go-version
    go get github.com/jackc/fake
    go get github.com/lib/pq
    go get github.com/pkg/errors
    go get github.com/satori/go.uuid
    go get github.com/shopspring/decimal
    go get github.com/sirupsen/logrus
    go get go.uber.org/zap
    go get gopkg.in/inconshreveable/log15.v2

Then run the following SQL:

    create user pgx_md5 password 'secret';
    create user " tricky, ' } "" \ test user " password 'secret';
    create database pgx_test;
    create user pgx_replication with replication password 'secret';

Connect to database pgx_test and run:

    create extension hstore;
    create domain uint64 as numeric(20,0);

Next open conn_config_test.go.example and make a copy without the
.example. If your PostgreSQL server is accepting connections on 127.0.0.1,
then you are done.

### Connection and Authentication Test Environment

Complete the normal test environment setup and also do the following.

Run the following SQL:

    create user pgx_none;
    create user pgx_pw password 'secret';

Add the following to your pg_hba.conf:

If you are developing on Unix with domain socket connections:

    local  pgx_test  pgx_none  trust
    local  pgx_test  pgx_pw    password
    local  pgx_test  pgx_md5   md5

If you are developing on Windows with TCP connections:

    host  pgx_test  pgx_none  127.0.0.1/32 trust
    host  pgx_test  pgx_pw    127.0.0.1/32 password
    host  pgx_test  pgx_md5   127.0.0.1/32 md5

### Replication Test Environment

Add a replication user:

    create user pgx_replication with replication password 'secret';

Add a replication line to your pg_hba.conf:

    host replication pgx_replication 127.0.0.1/32 md5

Change the following settings in your postgresql.conf:

    wal_level=logical
    max_wal_senders=5
    max_replication_slots=5

Set `replicationConnConfig` appropriately in `conn_config_test.go`.

## Version Policy

pgx follows semantic versioning for the documented public API on stable releases. Branch `v3` is the latest stable release. `master` can contain new features or behavior that will change or be removed before being merged to the stable `v3` branch (in practice, this occurs very rarely). `v2` is the previous stable release.
