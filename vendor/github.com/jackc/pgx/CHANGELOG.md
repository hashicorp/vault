# 3.3.0 (December 1, 2018)

## Features

* Add CopyFromReader and CopyToWriter (Murat Kabilov)
* Add MacaddrArray (Anthony Regeda)
* Add float types to FieldDescription.Type (David Yamnitsky)
* Add CheckedOutConnections helper method (MOZGIII)
* Add host query parameter to support Unix sockets (Jörg Thalheim)
* Custom cancelation hook for use with PostgreSQL-like databases (James Hartig)
* Added LastStmtSent for safe retry logic (James Hartig)

## Fixes

* Do not silently ignore assign NULL to \*string
* Fix issue with JSON and driver.Valuer conversion
* Fix race with stdlib Driver.configs Open (Greg Curtis)

## Changes

* Connection pool uses connections in queue order instead of stack. This
  minimized the time any connection is idle vs. any other connection.
  (Anthony Regeda)
* FieldDescription.Modifier is int32 instead of uint32
* tls: stop sending ssl_renegotiation_limit in startup message (Tejas Manohar)

# 3.2.0 (August 7, 2018)

## Features

* Support sslkey, sslcert, and sslrootcert URI params (Sean Chittenden)
* Allow any scheme in ParseURI (for convenience with cockroachdb) (Sean Chittenden)
* Add support for domain types
* Add zerolog logging adaptor (Justin Reagor)
* Add driver.Connector support / Go 1.10 support (James Lawrence)
* Allow nested database/sql/driver.Drivers (Jackson Owens)
* Support int64 and uint64 numeric array (Anthony Regeda)
* Add nul support to pgtype.Bool (Tarik Demirci)
* Add types to decode error messages (Damir Vandic)


## Fixes

* Fix Rows.Values returning same value for multiple columns of same complex type
* Fix StartReplication() syntax (steampunkcoder)
* Fix precision loss for test format geometric types
* Allows scanning jsonb column into `*json.RawMessage`
* Allow recovery to savepoint in failed transaction
* Fix deadlock when CopyFromSource panics
* Include PreferSimpleProtocol in config Merge (Murat Kabilov)

## Changes

* pgtype.JSON(B).Value now returns []byte instead of string. This allows
  database/sql to scan json(b) into \*json.RawMessage. This is a tiny behavior
  change, but database/sql Scan should automatically convert []byte to string, so
  there shouldn't be any incompatibility.

# 3.1.0 (January 15, 2018)

## Features

* Add QueryEx, QueryRowEx, ExecEx, and RollbackEx to Tx
* Add more ColumnType support (Timothée Peignier)
* Add UUIDArray type (Kelsey Francis)
* Add zap log adapter (Kelsey Francis)
* Add CreateReplicationSlotEx that consistent_point and snapshot_name (Mark Fletcher)
* Add BeginBatch to Tx (Gaspard Douady)
* Support CrateDB (Felix Geisendörfer)
* Allow use of logrus logger with fields configured (André Bierlein)
* Add array of enum support
* Add support for bit type
* Handle timeout parameters (Timothée Peignier)
* Allow overriding connection info (James Lawrence)
* Add support for bpchar type (Iurii Krasnoshchok)
* Add ConnConfig.PreferSimpleProtocol

## Fixes

* Fix numeric EncodeBinary bug (Wei Congrui)
* Fix logrus updated package name (Damir Vandic)
* Fix some invalid one round trip execs failing to return non-nil error. (Kelsey Francis)
* Return ErrClosedPool when Acquire() with closed pool (Mike Graf)
* Fix decoding row with same type values
* Always return non-nil \*Rows from Query to fix QueryRow (Kelsey Francis)
* Fix pgtype types that can Set database/sql/driver.driver.Valuer
* Prefix types in namespaces other than pg_catalog or public (Kelsey Francis)
* Fix incomplete selects during batch (Gaspard Douady and Jack Christensen)
* Support nil pointers to value implementing driver.Valuer
* Fix time logging for QueryEx
* Fix ranges with text format where end is unbounded
* Detect erroneous JSON(B) encoding
* Fix missing interval mapping
* ConnPool begin should not retry if ctx is done (Gaspard Douady)
* Fix reading interrupted messages could break connection
* Return error on unknown oid while decoding record instead of panic (Iurii Krasnoshchok)

## Changes

* Align sslmode "require" more closely to libpq (Johan Brandhorst)

# 3.0.1 (August 12, 2017)

## Fixes

* Fix compilation on 32-bit platform
* Fix invalid MarshalJSON of types with status Undefined
* Fix pid logging

# 3.0.0 (July 24, 2017)

## Changes

* Pid to PID in accordance with Go naming conventions.
* Conn.Pid changed to accessor method Conn.PID()
* Conn.SecretKey removed
* Remove Conn.TxStatus
* Logger interface reduced to single Log method.
* Replace BeginIso with BeginEx. BeginEx adds support for read/write mode and deferrable mode.
* Transaction isolation level constants are now typed strings instead of bare strings.
* Conn.WaitForNotification now takes context.Context instead of time.Duration for cancellation support.
* Conn.WaitForNotification no longer automatically pings internally every 15 seconds.
* ReplicationConn.WaitForReplicationMessage now takes context.Context instead of time.Duration for cancellation support.
* Reject scanning binary format values into a string (e.g. binary encoded timestamptz to string). See https://github.com/jackc/pgx/issues/219 and https://github.com/jackc/pgx/issues/228
* No longer can read raw bytes of any value into a []byte. Use pgtype.GenericBinary if this functionality is needed.
* Remove CopyTo (functionality is now in CopyFrom)
* OID constants moved from pgx to pgtype package
* Replaced Scanner, Encoder, and PgxScanner interfaces with pgtype system
* Removed ValueReader
* ConnPool.Close no longer waits for all acquired connections to be released. Instead, it immediately closes all available connections, and closes acquired connections when they are released in the same manner as ConnPool.Reset.
* Removed Rows.Fatal(error)
* Removed Rows.AfterClose()
* Removed Rows.Conn()
* Removed Tx.AfterClose()
* Removed Tx.Conn()
* Use Go casing convention for OID, UUID, JSON(B), ACLItem, CID, TID, XID, and CIDR
* Replaced stdlib.OpenFromConnPool with DriverConfig system

## Features

* Entirely revamped pluggable type system that supports approximately 60 PostgreSQL types.
* Types support database/sql interfaces and therefore can be used with other drivers
* Added context methods supporting cancellation where appropriate
* Added simple query protocol support
* Added single round-trip query mode
* Added batch query operations
* Added OnNotice
* github.com/pkg/errors used where possible for errors
* Added stdlib.DriverConfig which directly allows full configuration of underlying pgx connections without needing to use a pgx.ConnPool
* Added AcquireConn and ReleaseConn to stdlib to allow acquiring a connection from a database/sql connection.

# 2.11.0 (June 5, 2017)

## Fixes

* Fix race with concurrent execution of stdlib.OpenFromConnPool (Terin Stock)

## Features

* .pgpass support (j7b)
* Add missing CopyFrom delegators to Tx and ConnPool (Jack Christensen)
* Add ParseConnectionString (James Lawrence)

## Performance

* Optimize HStore encoding (René Kroon)

# 2.10.0 (March 17, 2017)

## Fixes

* database/sql driver created through stdlib.OpenFromConnPool closes connections when requested by database/sql rather than release to underlying connection pool.

# 2.11.0 (June 5, 2017)

## Fixes

* Fix race with concurrent execution of stdlib.OpenFromConnPool (Terin Stock)

## Features

* .pgpass support (j7b)
* Add missing CopyFrom delegators to Tx and ConnPool (Jack Christensen)
* Add ParseConnectionString (James Lawrence)

## Performance

* Optimize HStore encoding (René Kroon)

# 2.10.0 (March 17, 2017)

## Fixes

* Oid underlying type changed to uint32, previously it was incorrectly int32 (Manni Wood)
* Explicitly close checked-in connections on ConnPool.Reset, previously they were closed by GC

## Features

* Add xid type support (Manni Wood)
* Add cid type support (Manni Wood)
* Add tid type support (Manni Wood)
* Add "char" type support (Manni Wood)
* Add NullOid type (Manni Wood)
* Add json/jsonb binary support to allow use with CopyTo
* Add named error ErrAcquireTimeout (Alexander Staubo)
* Add logical replication decoding (Kris Wehner)
* Add PgxScanner interface to allow types to simultaneously support database/sql and pgx (Jack Christensen)
* Add CopyFrom with schema support (Jack Christensen)

## Compatibility

* jsonb now defaults to binary format. This means passing a []byte to a jsonb column will no longer work.
* CopyTo is now deprecated but will continue to work.

# 2.9.0 (August 26, 2016)

## Fixes

* Fix *ConnPool.Deallocate() not deleting prepared statement from map
* Fix stdlib not logging unprepared query SQL (Krzysztof Dryś)
* Fix Rows.Values() with varchar binary format
* Concurrent ConnPool.Acquire calls with Dialer timeouts now timeout in the expected amount of time (Konstantin Dzreev)

## Features

* Add CopyTo
* Add PrepareEx
* Add basic record to []interface{} decoding
* Encode and decode between all Go and PostgreSQL integer types with bounds checking
* Decode inet/cidr to net.IP
* Encode/decode [][]byte to/from bytea[]
* Encode/decode named types whose underlying types are string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64

## Performance

* Substantial reduction in memory allocations

# 2.8.1 (March 24, 2016)

## Features

* Scan accepts nil argument to ignore a column

## Fixes

* Fix compilation on 32-bit architecture
* Fix Tx.status not being set on error on Commit
* Fix Listen/Unlisten with special characters

# 2.8.0 (March 18, 2016)

## Fixes

* Fix unrecognized commit failure
* Fix msgReader.rxMsg bug when msgReader already has error
* Go float64 can no longer be encoded to a PostgreSQL float4
* Fix connection corruption when query with error is closed early

## Features

This release adds multiple extension points helpful when wrapping pgx with
custom application behavior. pgx can now use custom types designed for the
standard database/sql package such as
[github.com/shopspring/decimal](https://github.com/shopspring/decimal).

* Add *Tx.AfterClose() hook
* Add *Tx.Conn()
* Add *Tx.Status()
* Add *Tx.Err()
* Add *Rows.AfterClose() hook
* Add *Rows.Conn()
* Add *Conn.SetLogger() to allow changing logger
* Add *Conn.SetLogLevel() to allow changing log level
* Add ConnPool.Reset method
* Add support for database/sql.Scanner and database/sql/driver.Valuer interfaces
* Rows.Scan errors now include which argument caused error
* Add Encode() to allow custom Encoders to reuse internal encoding functionality
* Add Decode() to allow customer Decoders to reuse internal decoding functionality
* Add ConnPool.Prepare method
* Add ConnPool.Deallocate method
* Add Scan to uint32 and uint64 (utrack)
* Add encode and decode to []uint16, []uint32, and []uint64 (Max Musatov)

## Performance

* []byte skips encoding/decoding

# 2.7.1 (October 26, 2015)

* Disable SSL renegotiation

# 2.7.0 (October 16, 2015)

* Add RuntimeParams to ConnConfig
* ParseURI extracts RuntimeParams
* ParseDSN extracts RuntimeParams
* ParseEnvLibpq extracts PGAPPNAME
* Prepare is now idempotent
* Rows.Values now supports oid type
* ConnPool.Release automatically unlistens connections (Joseph Glanville)
* Add trace log level
* Add more efficient log leveling
* Retry automatically on ConnPool.Begin (Joseph Glanville)
* Encode from net.IP to inet and cidr
* Generalize encoding pointer to string to any PostgreSQL type
* Add UUID encoding from pointer to string (Joseph Glanville)
* Add null mapping to pointer to pointer (Jonathan Rudenberg)
* Add JSON and JSONB type support (Joseph Glanville)

# 2.6.0 (September 3, 2015)

* Add inet and cidr type support
* Add binary decoding to TimestampOid in stdlib driver (Samuel Stauffer)
* Add support for specifying sslmode in connection strings (Rick Snyder)
* Allow ConnPool to have MaxConnections of 1
* Add basic PGSSLMODE to support to ParseEnvLibpq
* Add fallback TLS config
* Expose specific error for TSL refused
* More error details exposed in PgError
* Support custom dialer (Lewis Marshall)

# 2.5.0 (April 15, 2015)

* Fix stdlib nil support (Blaž Hrastnik)
* Support custom Scanner not reading entire value
* Fix empty array scanning (Laurent Debacker)
* Add ParseDSN (deoxxa)
* Add timestamp support to NullTime
* Remove unused text format scanners
* Return error when too many parameters on Prepare
* Add Travis CI integration (Jonathan Rudenberg)
* Large object support (Jonathan Rudenberg)
* Fix reading null byte arrays (Karl Seguin)
* Add timestamptz[] support
* Add timestamp[] support (Karl Seguin)
* Add bool[] support (Karl Seguin)
* Allow writing []byte into text and varchar columns without type conversion (Hari Bhaskaran)
* Fix ConnPool Close panic
* Add Listen / notify example
* Reduce memory allocations (Karl Seguin)

# 2.4.0 (October 3, 2014)

* Add per connection oid to name map
* Add Hstore support (Andy Walker)
* Move introductory docs to godoc from readme
* Fix documentation references to TextEncoder and BinaryEncoder
* Add keep-alive to TCP connections (Andy Walker)
* Add support for EmptyQueryResponse / Allow no-op Exec (Andy Walker)
* Allow reading any type into []byte
* WaitForNotification detects lost connections quicker

# 2.3.0 (September 16, 2014)

* Truncate logged strings and byte slices
* Extract more error information from PostgreSQL
* Fix data race with Rows and ConnPool
