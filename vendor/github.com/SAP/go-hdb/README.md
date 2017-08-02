go-hdb
======

[![GoDoc](https://godoc.org/github.com/SAP/go-hdb/driver?status.png)](https://godoc.org/github.com/SAP/go-hdb/driver)

Go-hdb is a native Go (golang) HANA database driver for Go's sql package. It implements the SAP HANA SQL command network protocol:  
<http://help.sap.com/hana/SAP_HANA_SQL_Command_Network_Protocol_Reference_en.pdf>

## Installation

```
go get github.com/SAP/go-hdb/driver
```

## Documentation

API documentation and documented examples can be found at <https://godoc.org/github.com/SAP/go-hdb/driver>.

## Tests

For running the driver tests a HANA Database server is required. The test user must have privileges to create a schema.

```
go test -dsn hdb://user:password@host:port
```

## Features

* Native Go implementation (no C libraries, CGO).
* Go <http://golang.org/pkg/database/sql> package compliant.
* Support of databse/sql/driver Execer and Queryer interface for parameter free statements and queries.
* Support of bulk inserts.
* Support of UTF-8 to / from CESU-8 encodings for HANA Unicode types.
* Build-in support of HANA decimals as Go rational numbers <http://golang.org/pkg/math/big>.
* Support of Large Object streaming.
* Support of Stored Procedures with table output parameters. 

## Dependencies

* <http://golang.org/x/text/transform>

## Todo

* Additional Authentication Methods (actually only basic authentication is supported).
