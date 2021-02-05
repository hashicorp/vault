# Changes

## [1.13.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.12.0...v1.13.0) (2021-01-15)


### Features

* **spanner/spannertest:** implement ANY_VALUE aggregation function ([#3428](https://www.github.com/googleapis/google-cloud-go/issues/3428)) ([e16c3e9](https://www.github.com/googleapis/google-cloud-go/commit/e16c3e9b412762b85483f3831ee586a5e6631313))
* **spanner/spannertest:** implement FULL JOIN ([#3218](https://www.github.com/googleapis/google-cloud-go/issues/3218)) ([99f7212](https://www.github.com/googleapis/google-cloud-go/commit/99f7212bd70bb333c1aa1c7a57348b4dfd80d31b))
* **spanner/spannertest:** implement SELECT ... FROM UNNEST(...) ([#3431](https://www.github.com/googleapis/google-cloud-go/issues/3431)) ([deb466f](https://www.github.com/googleapis/google-cloud-go/commit/deb466f497a1e6df78fcad57c3b90b1a4ccd93b4))
* **spanner/spannertest:** support array literals ([#3438](https://www.github.com/googleapis/google-cloud-go/issues/3438)) ([69e0110](https://www.github.com/googleapis/google-cloud-go/commit/69e0110f4977035cd1a705c3034c3ba96cadf36f))
* **spanner/spannertest:** support AVG aggregation function ([#3286](https://www.github.com/googleapis/google-cloud-go/issues/3286)) ([4788415](https://www.github.com/googleapis/google-cloud-go/commit/4788415c908f58c1cc08c951f1a7f17cdaf35aa2))
* **spanner/spannertest:** support Not Null constraint ([#3491](https://www.github.com/googleapis/google-cloud-go/issues/3491)) ([c36aa07](https://www.github.com/googleapis/google-cloud-go/commit/c36aa0785e798b9339d540e691850ca3c474a288))
* **spanner/spannertest:** support UPDATE DML ([#3201](https://www.github.com/googleapis/google-cloud-go/issues/3201)) ([1dec6f6](https://www.github.com/googleapis/google-cloud-go/commit/1dec6f6a31768a3f70bfec7274828301c22ea10b))
* **spanner/spansql:** define structures and parse UPDATE DML statements ([#3192](https://www.github.com/googleapis/google-cloud-go/issues/3192)) ([23b6904](https://www.github.com/googleapis/google-cloud-go/commit/23b69042c58489df512703259f54d075ba0c0722))
* **spanner/spansql:** support DATE and TIMESTAMP literals ([#3557](https://www.github.com/googleapis/google-cloud-go/issues/3557)) ([1961930](https://www.github.com/googleapis/google-cloud-go/commit/196193034a15f84dc3d3c27901990e8be77fca85))
* **spanner/spansql:** support for parsing generated columns ([#3373](https://www.github.com/googleapis/google-cloud-go/issues/3373)) ([9b1d06f](https://www.github.com/googleapis/google-cloud-go/commit/9b1d06fc90a4c07899c641a893dba0b47a1cead9))
* **spanner/spansql:** support NUMERIC data type ([#3411](https://www.github.com/googleapis/google-cloud-go/issues/3411)) ([1bc65d9](https://www.github.com/googleapis/google-cloud-go/commit/1bc65d9124ba22db5bec4c71b6378c27dfc04724))
* **spanner:** Add a DirectPath fallback integration test ([#3487](https://www.github.com/googleapis/google-cloud-go/issues/3487)) ([de821c5](https://www.github.com/googleapis/google-cloud-go/commit/de821c59fb81e9946216d205162b59de8b5ce71c))
* **spanner:** attempt DirectPath by default ([#3516](https://www.github.com/googleapis/google-cloud-go/issues/3516)) ([bbc61ed](https://www.github.com/googleapis/google-cloud-go/commit/bbc61ed368453b28aaf5bed627ca2499a3591f63))
* **spanner:** include User agent ([#3465](https://www.github.com/googleapis/google-cloud-go/issues/3465)) ([4e1ef1b](https://www.github.com/googleapis/google-cloud-go/commit/4e1ef1b3fb536ef950249cdee02cc0b6c2b56e86))
* **spanner:** run E2E test over DirectPath ([#3466](https://www.github.com/googleapis/google-cloud-go/issues/3466)) ([18e3a4f](https://www.github.com/googleapis/google-cloud-go/commit/18e3a4fe2a0c59c6295db2d85c7893ac51688083))
* **spanner:** support NUMERIC in mutations ([#3328](https://www.github.com/googleapis/google-cloud-go/issues/3328)) ([fa90737](https://www.github.com/googleapis/google-cloud-go/commit/fa90737a2adbe0cefbaba4aa1046a6efbba2a0e9))


### Bug Fixes

* **spanner:** fix session leak ([#3461](https://www.github.com/googleapis/google-cloud-go/issues/3461)) ([11fb917](https://www.github.com/googleapis/google-cloud-go/commit/11fb91711db5b941995737980cef7b48b611fefd)), refs [#3460](https://www.github.com/googleapis/google-cloud-go/issues/3460)

## [1.12.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.11.0...v1.12.0) (2020-11-10)


### Features

* **spanner:** add metadata to RowIterator ([#3050](https://www.github.com/googleapis/google-cloud-go/issues/3050)) ([9a2289c](https://www.github.com/googleapis/google-cloud-go/commit/9a2289c3a38492bc2e84e0f4000c68a8718f5c11)), closes [#1805](https://www.github.com/googleapis/google-cloud-go/issues/1805)
* **spanner:** export ToSpannerError ([#3133](https://www.github.com/googleapis/google-cloud-go/issues/3133)) ([b951d8b](https://www.github.com/googleapis/google-cloud-go/commit/b951d8bd194b76da0a8bf2ce7cf85b546d2e051c)), closes [#3122](https://www.github.com/googleapis/google-cloud-go/issues/3122)
* **spanner:** support rw-transaction with options ([#3058](https://www.github.com/googleapis/google-cloud-go/issues/3058)) ([5130694](https://www.github.com/googleapis/google-cloud-go/commit/51306948eef9d26cff70453efc3eb500ddef9117))
* **spanner/spannertest:** make SELECT list aliases visible to ORDER BY ([#3054](https://www.github.com/googleapis/google-cloud-go/issues/3054)) ([7d2d83e](https://www.github.com/googleapis/google-cloud-go/commit/7d2d83ee1cce58d4014d5570bc599bcef1ed9c22)), closes [#3043](https://www.github.com/googleapis/google-cloud-go/issues/3043)

## v1.11.0

* Features:
  - feat(spanner): add KeySetFromKeys function (#2837)
* Misc:
  - test(spanner): check for Aborted error (#3039)
  - test(spanner): fix potential race condition in TestRsdBlockingStates (#3017)
  - test(spanner): compare data instead of struct (#3013)
  - test(spanner): fix flaky oc_test.go (#2838)
  - docs(spanner): document NULL value (#2885)
* spansql/spannertest:
  - Support JOINs (all but FULL JOIN) (#2936, #2924, #2896, #3042, #3037, #2995, #2945, #2931)
  - feat(spanner/spansql): parse CHECK constraints (#3046)
  - fix(spanner/spansql): fix parsing of unary minus and plus (#2997)
  - fix(spanner/spansql): fix parsing of adjacent inline and leading comments (#2851)
  - fix(spanner/spannertest): fix ORDER BY combined with SELECT aliases (#3043)
  - fix(spanner/spannertest): generate query output columns in construction order (#2990)
  - fix(spanner/spannertest): correct handling of NULL AND FALSE (#2991)
  - fix(spanner/spannertest): correct handling of tri-state boolean expression evaluation (#2983)
  - fix(spanner/spannertest): fix handling of NULL with LIKE operator (#2982)
  - test(spanner/spannertest): migrate most test code to integration_test.go (#2977)
  - test(spanner/spansql): add fuzz target for ParseQuery (#2909)
  - doc(spanner/spannertest): document the implementation (#2996)
  - perf(spanner/spannertest): speed up no-wait DDL changes (#2994)
  - perf(spanner/spansql): make fewer allocations during SQL (#2969)
* Backward Incompatible Changes
  - chore(spanner/spansql): use ID type for identifiers throughout (#2889)
  - chore(spanner/spansql): restructure FROM, TABLESAMPLE (#2888)

## v1.10.0

* feat(spanner): add support for NUMERIC data type (#2415)
* feat(spanner): add custom type support to spanner.Key (#2748)
* feat(spanner/spannertest): add support for bool parameter types (#2674)
* fix(spanner): update PDML to take sessions from pool (#2736)
* spanner/spansql: update docs on TableAlteration, ColumnAlteration (#2825)
* spanner/spannertest: support dropping columns (#2823)
* spanner/spannertest: implement GetDatabase (#2802)
* spanner/spannertest: fix aggregation in query evaluation for empty inputs (#2803)

## v1.9.0

* Features:
  - feat(spanner): support custom field type (#2614)
* Bugfixes:
  - fix(spanner): call ctx.cancel after stats have been recorded (#2728)
  - fix(spanner): retry session not found for read (#2724)
  - fix(spanner): specify credentials with SPANNER_EMULATOR_HOST (#2701)
  - fix(spanner): update pdml to retry EOS internal error (#2678)
* Misc:
  - test(spanner): unskip tests for emulator (#2675)
* spansql/spannertest:
  - spanner/spansql: restructure types and parsing for column options (#2656)
  - spanner/spannertest: return error for Read with no keys (#2655)

## v1.8.0

* Features:
  - feat(spanner): support of client-level custom retry settings (#2599)
  - feat(spanner): add a statement-based way to run read-write transaction. (#2545)
* Bugfixes:
  - fix(spanner): set 'gccl' to the request header. (#2609)
  - fix(spanner): add the missing resource prefix (#2605)
  - fix(spanner): fix the upgrade of protobuf. (#2583)
  - fix(spanner): do not copy protobuf messages by value. (#2581)
  - fix(spanner): fix the required resource prefix. (#2580)
  - fix(spanner): add extra field to ignore with cmp (#2577)
  - fix(spanner): remove appengine-specific numChannels. (#2513)
* Misc:
  - test(spanner): log warning instead of fail for stress test (#2559)
  - test(spanner): fix failed TestRsdBlockingStates test (#2597)
  - chore(spanner): cleanup mockserver and mockclient (#2414)

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
