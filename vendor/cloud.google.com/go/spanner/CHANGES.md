# Changes

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
