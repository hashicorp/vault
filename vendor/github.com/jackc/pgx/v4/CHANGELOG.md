# 4.18.3 (March 9, 2024)

Use spaces instead of parentheses for SQL sanitization.

This still solves the problem of negative numbers creating a line comment, but this avoids breaking edge cases such as
`set foo to $1` where the substitution is taking place in a location where an arbitrary expression is not allowed.

# 4.18.2 (March 4, 2024)

Fix CVE-2024-27289

SQL injection can occur when all of the following conditions are met:

1. The non-default simple protocol is used.
2. A placeholder for a numeric value must be immediately preceded by a minus.
3. There must be a second placeholder for a string value after the first placeholder; both must be on the same line.
4. Both parameter values must be user-controlled.

Thanks to Paul Gerste for reporting this issue.

Fix CVE-2024-27304

SQL injection can occur if an attacker can cause a single query or bind message to exceed 4 GB in size. An integer
overflow in the calculated message size can cause the one large message to be sent as multiple messages under the
attacker's control.

Thanks to Paul Gerste for reporting this issue.

* Fix *dbTx.Exec not checking if it is already closed

# 4.18.1 (February 27, 2023)

* Fix: Support pgx v4 and v5 stdlib in same program (Tomáš Procházka)

# 4.18.0 (February 11, 2023)

* Upgrade pgconn to v1.14.0
* Upgrade pgproto3 to v2.3.2
* Upgrade pgtype to v1.14.0
* Fix query sanitizer when query text contains Unicode replacement character
* Fix context with value in BeforeConnect (David Harju)
* Support pgx v4 and v5 stdlib in same program (Vitalii Solodilov)

# 4.17.2 (September 3, 2022)

* Fix panic when logging batch error (Tom Möller)

# 4.17.1 (August 27, 2022)

* Upgrade puddle to v1.3.0 - fixes context failing to cancel Acquire when acquire is creating resource which was introduced in v4.17.0 (James Hartig)
* Fix atomic alignment on 32-bit platforms

# 4.17.0 (August 6, 2022)

* Upgrade pgconn to v1.13.0
* Upgrade pgproto3 to v2.3.1
* Upgrade pgtype to v1.12.0
* Allow background pool connections to continue even if cause is canceled (James Hartig)
* Add LoggerFunc (Gabor Szabad)
* pgxpool: health check should avoid going below minConns (James Hartig)
* Add pgxpool.Conn.Hijack()
* Logging improvements (Stepan Rabotkin)

# 4.16.1 (May 7, 2022)

* Upgrade pgconn to v1.12.1
* Fix explicitly prepared statements with describe statement cache mode

# 4.16.0 (April 21, 2022)

* Upgrade pgconn to v1.12.0
* Upgrade pgproto3 to v2.3.0
* Upgrade pgtype to v1.11.0
* Fix: Do not panic when context cancelled while getting statement from cache.
* Fix: Less memory pinning from old Rows.
* Fix: Support '\r' line ending when sanitizing SQL comment.
* Add pluggable GSSAPI support (Oliver Tan)

# 4.15.0 (February 7, 2022)

* Upgrade to pgconn v1.11.0
* Upgrade to pgtype v1.10.0
* Upgrade puddle to v1.2.1
* Make BatchResults.Close safe to be called multiple times

# 4.14.1 (November 28, 2021)

* Upgrade pgtype to v1.9.1 (fixes unintentional change to timestamp binary decoding)
* Start pgxpool background health check after initial connections

# 4.14.0 (November 20, 2021)

* Upgrade pgconn to v1.10.1
* Upgrade pgproto3 to v2.2.0
* Upgrade pgtype to v1.9.0
* Upgrade puddle to v1.2.0
* Add QueryFunc to BatchResults
* Add context options to zerologadapter (Thomas Frössman)
* Add zerologadapter.NewContextLogger (urso)
* Eager initialize minpoolsize on connect (Daniel)
* Unpin memory used by large queries immediately after use

# 4.13.0 (July 24, 2021)

* Trimmed pseudo-dependencies in Go modules from other packages tests
* Upgrade pgconn -- context cancellation no longer will return a net.Error
* Support time durations for simple protocol (Michael Darr)

# 4.12.0 (July 10, 2021)

* ResetSession hook is called before a connection is reused from pool for another query (Dmytro Haranzha)
* stdlib: Add RandomizeHostOrderFunc (dkinder)
* stdlib: add OptionBeforeConnect (dkinder)
* stdlib: Do not reuse ConnConfig strings (Andrew Kimball)
* stdlib: implement Conn.ResetSession (Jonathan Amsterdam)
* Upgrade pgconn to v1.9.0
* Upgrade pgtype to v1.8.0

# 4.11.0 (March 25, 2021)

* Add BeforeConnect callback to pgxpool.Config (Robert Froehlich)
* Add Ping method to pgxpool.Conn (davidsbond)
* Added a kitlog level log adapter (Fabrice Aneche)
* Make ScanArgError public to allow identification of offending column (Pau Sanchez)
* Add *pgxpool.AcquireFunc
* Add BeginFunc and BeginTxFunc
* Add prefer_simple_protocol to connection string
* Add logging on CopyFrom (Patrick Hemmer)
* Add comment support when sanitizing SQL queries (Rusakow Andrew)
* Do not panic on double close of pgxpool.Pool (Matt Schultz)
* Avoid panic on SendBatch on closed Tx (Matt Schultz)
* Update pgconn to v1.8.1
* Update pgtype to v1.7.0

# 4.10.1 (December 19, 2020)

* Fix panic on Query error with nil stmtcache.

# 4.10.0 (December 3, 2020)

* Add CopyFromSlice to simplify CopyFrom usage (Egon Elbre)
* Remove broken prepared statements from stmtcache (Ethan Pailes)
* stdlib: consider any Ping error as fatal
* Update puddle to v1.1.3 - this fixes an issue where concurrent Acquires can hang when a connection cannot be established
* Update pgtype to v1.6.2

# 4.9.2 (November 3, 2020)

The underlying library updates fix an issue where appending to a scanned slice could corrupt other data.

* Update pgconn to v1.7.2
* Update pgproto3 to v2.0.6

# 4.9.1 (October 31, 2020)

* Update pgconn to v1.7.1
* Update pgtype to v1.6.1
* Fix SendBatch of all prepared statements with statement cache disabled

# 4.9.0 (September 26, 2020)

* pgxpool now waits for connection cleanup to finish before making room in pool for another connection. This prevents temporarily exceeding max pool size.
* Fix when scanning a column to nil to skip it on the first row but scanning it to a real value on a subsequent row.
* Fix prefer simple protocol with prepared statements. (Jinzhu)
* Fix FieldDescriptions not being available on Rows before calling Next the first time.
* Various minor fixes in updated versions of pgconn, pgtype, and puddle.

# 4.8.1 (July 29, 2020)

* Update pgconn to v1.6.4
    * Fix deadlock on error after CommandComplete but before ReadyForQuery
    * Fix panic on parsing DSN with trailing '='

# 4.8.0 (July 22, 2020)

* All argument types supported by native pgx should now also work through database/sql
* Update pgconn to v1.6.3
* Update pgtype to v1.4.2

# 4.7.2 (July 14, 2020)

* Improve performance of Columns() (zikaeroh)
* Fix fatal Commit() failure not being considered fatal
* Update pgconn to v1.6.2
* Update pgtype to v1.4.1

# 4.7.1 (June 29, 2020)

* Fix stdlib decoding error with certain order and combination of fields

# 4.7.0 (June 27, 2020)

* Update pgtype to v1.4.0
* Update pgconn to v1.6.1
* Update puddle to v1.1.1
* Fix context propagation with Tx commit and Rollback (georgysavva)
* Add lazy connect option to pgxpool (georgysavva)
* Fix connection leak if pgxpool.BeginTx() fail (Jean-Baptiste Bronisz)
* Add native Go slice support for strings and numbers to simple protocol
* stdlib add default timeouts for Conn.Close() and Stmt.Close() (georgysavva)
* Assorted performance improvements especially with large result sets
* Fix close pool on not lazy connect failure (Yegor Myskin)
* Add Config copy (georgysavva)
* Support SendBatch with Simple Protocol (Jordan Lewis)
* Better error logging on rows close (Igor V. Kozinov)
* Expose stdlib.Conn.Conn() to enable database/sql.Conn.Raw()
* Improve unknown type support for database/sql
* Fix transaction commit failure closing connection

# 4.6.0 (March 30, 2020)

* stdlib: Bail early if preloading rows.Next() results in rows.Err() (Bas van Beek)
* Sanitize time to microsecond accuracy (Andrew Nicoll)
* Update pgtype to v1.3.0
* Update pgconn to v1.5.0
    * Update golang.org/x/crypto for security fix
    * Implement "verify-ca" SSL mode

# 4.5.0 (March 7, 2020)

* Update to pgconn v1.4.0
    * Fixes QueryRow with empty SQL
    * Adds PostgreSQL service file support
* Add Len() to *pgx.Batch (WGH)
* Better logging for individual batch items (Ben Bader)

# 4.4.1 (February 14, 2020)

* Update pgconn to v1.3.2 - better default read buffer size
* Fix race in CopyFrom

# 4.4.0 (February 5, 2020)

* Update puddle to v1.1.0 - fixes possible deadlock when acquire is cancelled
* Update pgconn to v1.3.1 - fixes CopyFrom deadlock when multiple NoticeResponse received during copy
* Update pgtype to v1.2.0
* Add MaxConnIdleTime to pgxpool (Patrick Ellul)
* Add MinConns to pgxpool (Patrick Ellul)
* Fix: stdlib.ReleaseConn closes connections left in invalid state

# 4.3.0 (January 23, 2020)

* Fix Rows.Values panic when unable to decode
* Add Rows.Values support for unknown types
* Add DriverContext support for stdlib (Alex Gaynor)
* Update pgproto3 to v2.0.1 to never return an io.EOF as it would be misinterpreted by database/sql. Instead return io.UnexpectedEOF.

# 4.2.1 (January 13, 2020)

* Update pgconn to v1.2.1 (fixes context cancellation data race introduced in v1.2.0))

# 4.2.0 (January 11, 2020)

* Update pgconn to v1.2.0.
* Update pgtype to v1.1.0.
* Return error instead of panic when wrong number of arguments passed to Exec. (malstoun)
* Fix large objects functionality when PreferSimpleProtocol = true.
* Restore GetDefaultDriver which existed in v3. (Johan Brandhorst)
* Add RegisterConnConfig to stdlib which replaces the removed RegisterDriverConfig from v3.

# 4.1.2 (October 22, 2019)

* Fix dbSavepoint.Begin recursive self call
* Upgrade pgtype to v1.0.2 - fix scan pointer to pointer

# 4.1.1 (October 21, 2019)

* Fix pgxpool Rows.CommandTag() infinite loop / typo

# 4.1.0 (October 12, 2019)

## Potentially Breaking Changes

Technically, two changes are breaking changes, but in practice these are extremely unlikely to break existing code.

* Conn.Begin and Conn.BeginTx return a Tx interface instead of the internal dbTx struct. This is necessary for the Conn.Begin method to signature as other methods that begin a transaction.
* Add Conn() to Tx interface. This is necessary to allow code using a Tx to access the *Conn (and pgconn.PgConn) on which the Tx is executing.

## Fixes

* Releasing a busy connection closes the connection instead of returning an unusable connection to the pool
* Do not mutate config.Config.OnNotification in connect

# 4.0.1 (September 19, 2019)

* Fix statement cache cleanup.
* Corrected daterange OID.
* Fix Tx when committing or rolling back multiple times in certain cases.
* Improve documentation.

# 4.0.0 (September 14, 2019)

v4 is a major release with many significant changes some of which are breaking changes. The most significant are
included below.

* Simplified establishing a connection with a connection string.
* All potentially blocking operations now require a context.Context. The non-context aware functions have been removed.
* OIDs are hard-coded for known types. This saves the query on connection.
* Context cancellations while network activity is in progress is now always fatal. Previously, it was sometimes recoverable. This led to increased complexity in pgx itself and in application code.
* Go modules are required.
* Errors are now implemented in the Go 1.13 style.
* `Rows` and `Tx` are now interfaces.
* The connection pool as been decoupled from pgx and is now a separate, included package (github.com/jackc/pgx/v4/pgxpool).
* pgtype has been spun off to a separate package (github.com/jackc/pgtype).
* pgproto3 has been spun off to a separate package (github.com/jackc/pgproto3/v2).
* Logical replication support has been spun off to a separate package (github.com/jackc/pglogrepl).
* Lower level PostgreSQL functionality is now implemented in a separate package (github.com/jackc/pgconn).
* Tests are now configured with environment variables.
* Conn has an automatic statement cache by default.
* Batch interface has been simplified.
* QueryArgs has been removed.
