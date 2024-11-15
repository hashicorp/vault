# 1.14.3 (March 4, 2024)

* Update golang.org/x/crypto and golang.org/x/text

# 1.14.2 (March 4, 2024)

* Fix CVE-2024-27304. SQL injection can occur if an attacker can cause a single query or bind message to exceed 4 GB in
size. An integer overflow in the calculated message size can cause the one large message to be sent as multiple messages
under the attacker's control.

# 1.14.1 (July 19, 2023)

* Fix: Enable failover efforts when pg_hba.conf disallows non-ssl connections (Brandon Kauffman)
* Fix: connect_timeout is not obeyed for sslmode=allow|prefer (smaher-edb)
* Optimize redundant pgpass parsing in case password is explicitly set (Aleksandr Alekseev)

# 1.14.0 (February 11, 2023)

* Fix: each connection attempt to new node gets own timeout (Nathan Giardina)
* Set SNI for SSL connections (Stas Kelvich)
* Fix: CopyFrom I/O race (Tommy Reilly)
* Minor dependency upgrades

# 1.13.0 (August 6, 2022)

* Add sslpassword support (Eric McCormack and yun.xu)
* Add prefer-standby target_session_attrs support (sergey.bashilov)
* Fix GSS ErrorResponse handling (Oliver Tan)

# 1.12.1 (May 7, 2022)

* Fix: setting krbspn and krbsrvname in connection string (sireax)
* Add support for Unix sockets on Windows (Eno Compton)
* Stop ignoring ErrorResponse during SCRAM auth (Rafi Shamim)

# 1.12.0 (April 21, 2022)

* Add pluggable GSSAPI support (Oliver Tan)
* Fix: Consider any "0A000" error a possible cached plan changed error due to locale
* Better match psql fallback behavior with multiple hosts

# 1.11.0 (February 7, 2022)

* Support port in ip from LookupFunc to override config (James Hartig)
* Fix TLS connection timeout (Blake Embrey)
* Add support for read-only, primary, standby, prefer-standby target_session_attributes (Oscar)
* Fix connect when receiving NoticeResponse

# 1.10.1 (November 20, 2021)

* Close without waiting for response (Kei Kamikawa)
* Save waiting for network round-trip in CopyFrom (Rueian)
* Fix concurrency issue with ContextWatcher
* LRU.Get always checks context for cancellation / expiration (Georges Varouchas)

# 1.10.0 (July 24, 2021)

* net.Timeout errors are no longer returned when a query is canceled via context. A wrapped context error is returned.

# 1.9.0 (July 10, 2021)

* pgconn.Timeout only is true for errors originating in pgconn (Michael Darr)
* Add defaults for sslcert, sslkey, and sslrootcert (Joshua Brindle)
* Solve issue with 'sslmode=verify-full' when there are multiple hosts (mgoddard)
* Fix default host when parsing URL without host but with port
* Allow dbname query parameter in URL conn string
* Update underlying dependencies

# 1.8.1 (March 25, 2021)

* Better connection string sanitization (ip.novikov)
* Use proper pgpass location on Windows (Moshe Katz)
* Use errors instead of golang.org/x/xerrors
* Resume fallback on server error in Connect (Andrey Borodin)

# 1.8.0 (December 3, 2020)

* Add StatementErrored method to stmtcache.Cache. This allows the cache to purge invalidated prepared statements. (Ethan Pailes)

# 1.7.2 (November 3, 2020)

* Fix data value slices into work buffer with capacities larger than length.

# 1.7.1 (October 31, 2020)

* Do not asyncClose after receiving FATAL error from PostgreSQL server

# 1.7.0 (September 26, 2020)

* Exec(Params|Prepared) return ResultReader with FieldDescriptions loaded
* Add ReceiveResults (Sebastiaan Mannem)
* Fix parsing DSN connection with bad backslash
* Add PgConn.CleanupDone so connection pools can determine when async close is complete

# 1.6.4 (July 29, 2020)

* Fix deadlock on error after CommandComplete but before ReadyForQuery
* Fix panic on parsing DSN with trailing '='

# 1.6.3 (July 22, 2020)

* Fix error message after AppendCertsFromPEM failure (vahid-sohrabloo)

# 1.6.2 (July 14, 2020)

* Update pgservicefile library

# 1.6.1 (June 27, 2020)

* Update golang.org/x/crypto to latest
* Update golang.org/x/text to 0.3.3
* Fix error handling for bad PGSERVICE definition
* Redact passwords in ParseConfig errors (Lukas Vogel)

# 1.6.0 (June 6, 2020)

* Fix panic when closing conn during cancellable query
* Fix behavior of sslmode=require with sslrootcert present (Petr Jediný)
* Fix field descriptions available after command concluded (Tobias Salzmann)
* Support connect_timeout (georgysavva)
* Handle IPv6 in connection URLs (Lukas Vogel)
* Fix ValidateConnect with cancelable context
* Improve CopyFrom performance
* Add Config.Copy (georgysavva)

# 1.5.0 (March 30, 2020)

* Update golang.org/x/crypto for security fix
* Implement "verify-ca" SSL mode (Greg Curtis)

# 1.4.0 (March 7, 2020)

* Fix ExecParams and ExecPrepared handling of empty query.
* Support reading config from PostgreSQL service files.

# 1.3.2 (February 14, 2020)

* Update chunkreader to v2.0.1 for optimized default buffer size.

# 1.3.1 (February 5, 2020)

* Fix CopyFrom deadlock when multiple NoticeResponse received during copy

# 1.3.0 (January 23, 2020)

* Add Hijack and Construct.
* Update pgproto3 to v2.0.1.

# 1.2.1 (January 13, 2020)

* Fix data race in context cancellation introduced in v1.2.0.

# 1.2.0 (January 11, 2020)

## Features

* Add Insert(), Update(), Delete(), and Select() statement type query methods to CommandTag.
* Add PgError.SQLState method. This could be used for compatibility with other drivers and databases.

## Performance

* Improve performance when context.Background() is used. (bakape)
* CommandTag.RowsAffected is faster and does not allocate.

## Fixes

* Try to cancel any in-progress query when a conn is closed by ctx cancel.
* Handle NoticeResponse during CopyFrom.
* Ignore errors sending Terminate message while closing connection. This mimics the behavior of libpq PGfinish.

# 1.1.0 (October 12, 2019)

* Add PgConn.IsBusy() method.

# 1.0.1 (September 19, 2019)

* Fix statement cache not properly cleaning discarded statements.
