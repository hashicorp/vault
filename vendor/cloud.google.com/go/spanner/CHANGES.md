# Changes

## v1.7.0
* Retry:
  - Only retry certain types of internal errors. (#2460)
* Tracing/metrics:
  - Never sample `ping()` trace spans (#2520)
  - Add oc tests for session pool metrics. (#2416)
* Encoding:
  - Allow encoding struct with custom types to mutation (#2529)
* spannertest:
  - Fix evaluation on IN (#2479)
  - Support MIN/MAX aggregation functions (#2411)
* Misc:
  - Fix TestClient_WithGRPCConnectionPoolAndNumChannels_Misconfigured test (#2539)
  - Cleanup backoff files and rename a variable (#2526)
  - Fix TestIntegration_DML test to return err from tx (#2509)
  - Unskip tests for emulator 0.8.0. (#2494)
  - Fix TestIntegration_StartBackupOperation test. (#2418)
  - Fix flakiness in TestIntegration_BatchDML_Error
  - Unskip TestIntegration_BatchDML and TestIntegration_BatchDML_TwoStatements
    for emulator by checking the existence of status.
  - Fix TestStressSessionPool test by taking lock while getting sessions from
    hc.

## v1.6.0
* Sessions:
  - Increase the number of sessions in batches instead of one by one when
    additional sessions are needed. The step size is set to 25, which means
    that whenever the session pool needs at least one more session, it will
    create a batch of 25 sessions.
* Emulator:
  - Run integration tests against the emulator in Kokoro Presubmit.
* RPC retrying:
  - Retry CreateDatabase on retryable codes.
* spannertest:
  - Change internal representation of DATE/TIMESTAMP values.
* spansql:
  - Cleanly parse adjacent comment marker/terminator.
  - Support FROM aliases in SELECT statements.
* Misc:
  - Fix comparing errors in tests.
  - Fix flaky session pool test.
  - Increase timeout in TestIntegration_ReadOnlyTransaction.
  - Fix incorrect instance IDs when deleting instances in tests.
  - Clean up test instances.
  - Clearify docs on Aborted transaction.
  - Fix timeout+staleness bound for test
  - Remove the support for resource-based routing.
  - Fix TestTransaction_SessionNotFound test.

## v1.5.1

* Fix incorrect decreasing metrics, numReads and numWrites.
* Fix an issue that XXX fields/methods are internal to proto and may change
  at any time. XXX_Merge panics in proto v1.4.0. Use proto.Merge instead of
  XXX_Merge.
* spannertest: handle list parameters in RPC interfacea.

## v1.5.0

* Metrics
  - Instrument client library with adding OpenCensus metrics. This allows for
    better monitoring of the session pool.
* Session management
  - Switch the session keepalive method from GetSession to SELECT 1.
* Emulator
  - Use client hooks for admin clients running against an emulator. With
    this change, users can use SPANNER_EMULATOR_HOST for initializing admin
    clients when running against an emulator.
* spansql
  - Add space between constraint name and foreign key def.
* Misc
  - Fix segfault when a non-existent credentials file had been specified.
  - Fix cleaning up instances in integration tests.
  - Fix race condition in batch read-only transaction.
  - Fix the flaky TestLIFOTakeWriteSessionOrder test.
  - Fix ITs to order results in SELECT queries.
  - Fix the documentation of timestamp bounds.
  - Fix the regex issue in managing backups.

## v1.4.0

- Support managed backups. This includes the API methods for CreateBackup,
  GetBackup, UpdateBackup, DeleteBackup and others. Also includes a simple
  wrapper in DatabaseAdminClient to create a backup.
- Update the healthcheck interval. The default interval is updated to 50 mins.
  By default, the first healthcheck is scheduled between 10 and 55 mins and
  the subsequent healthchecks are between 45 and 55 mins. This update avoids
  overloading the backend service with frequent healthchecking.

## v1.3.0

* Query options:
  - Adds the support of providing query options (optimizer version) via
    three ways (precedence follows the order):
    `client-level < environment variables < query-level`. The environment
    variable is set by "SPANNER_OPTIMIZER_VERSION".
* Connection pooling:
  - Use the new connection pooling in gRPC. This change deprecates
    `ClientConfig.numChannels` and users should move to
    `WithGRPCConnectionPool(numChannels)` at their earliest convenience.
    Example:
    ```go
    // numChannels (deprecated):
    err, client := NewClientWithConfig(ctx, database, ClientConfig{NumChannels: 8})

    // gRPC connection pool:
    err, client := NewClientWithConfig(ctx, database, ClientConfig{}, option.WithGRPCConnectionPool(8))
    ```
* Error handling:
  - Do not rollback after failed commit.
  - Return TransactionOutcomeUnknownError if a DEADLINE_EXCEEDED or CANCELED
    error occurs while a COMMIT request is in flight.
* spansql:
  - Added support for IN expressions and OFFSET clauses.
  - Fixed parsing of table constraints.
  - Added support for foreign key constraints in ALTER TABLE and CREATE TABLE.
  - Added support for GROUP BY clauses.
* spannertest:
  - Added support for IN expressions and OFFSET clauses.
  - Added support for GROUP BY clauses.
  - Fixed data race in query execution.
  - No longer rejects reads specifying an index to use.
  - Return last commit timestamp as read timestamp when requested.
  - Evaluate add, subtract, multiply, divide, unary
    negation, unary not, bitwise and/xor/or operations, as well as reporting
    column types for expressions involving any possible arithmetic
    operator.arithmetic expressions.
  - Fixed handling of descending primary keys.
* Misc:
  - Change default healthcheck interval to 30 mins to reduce the GetSession
    calls made to the backend.
  - Add marshal/unmarshal json for nullable types to support NullString,
    NullInt64, NullFloat64, NullBool, NullTime, NullDate.
  - Use ResourceInfo to extract error.
  - Extract retry info from status.

## v1.2.1

- Fix session leakage for ApplyAtLeastOnce. Previously session handles where
  leaked whenever Commit() returned a non-abort, non-session-not-found error,
  due to a missing recycle() call.
- Fix error for WriteStruct with pointers. This fixes a specific check for
  encoding and decoding to pointer types.
- Fix a GRPCStatus issue that returns a Status that has Unknown code if the
  base error is nil. Now, it always returns a Status based on Code field of
  current error.

## v1.2.0

- Support tracking stacktrace of sessionPool.take() that allows the user
  to instruct the session pool to keep track of the stacktrace of each
  goroutine that checks out a session from the pool. This is disabled by
  default, but it can be enabled by setting
  `SessionPoolConfig.TrackSessionHandles: true`.
- Add resource-based routing that includes a step to retrieve the
  instance-specific endpoint before creating the session client when
  creating a new spanner client. This is disabled by default, but it can
  be enabled by setting `GOOGLE_CLOUD_SPANNER_ENABLE_RESOURCE_BASED_ROUTING`.
- Make logger configurable so that the Spanner client can now be configured to
  use a specific logger instead of the standard logger.
- Support encoding custom types that point back to supported basic types.
- Allow decoding Spanner values to custom types that point back to supported
  types.

## v1.1.0

- The String() method of NullString, NullTime and NullDate will now return
  an unquoted string instead of a quoted string. This is a BREAKING CHANGE.
  If you relied on the old behavior, please use fmt.Sprintf("%q", T).
- The Spanner client will now use the new BatchCreateSessions RPC to initialize
  the session pool. This will improve the startup time of clients that are
  initialized with a minimum number of sessions greater than zero
  (i.e. SessionPoolConfig.MinOpened>0).
- Spanner clients that are created with the NewClient method will now default
  to a minimum of 100 opened sessions in the pool
  (i.e. SessionPoolConfig.MinOpened=100). This will improve the performance
  of the first transaction/query that is executed by an application, as a
  session will normally not have to be created as part of the transaction.
  Spanner clients that are created with the NewClientWithConfig method are
  not affected by this change.
- Spanner clients that are created with the NewClient method will now default
  to a write sessions fraction of 0.2 in the pool
  (i.e. SessionPoolConfig.WriteSessions=0.2).
  Spanner clients that are created with the NewClientWithConfig method are
  not affected by this change.
- The session pool maintenance worker has been improved so it keeps better
  track of the actual number of sessions needed. It will now less often delete
  and re-create sessions. This can improve the overall performance of
  applications with a low transaction rate.

## v1.0.0

This is the first tag to carve out spanner as its own module. See:
https://github.com/golang/go/wiki/Modules#is-it-possible-to-add-a-module-to-a-multi-module-repository.
