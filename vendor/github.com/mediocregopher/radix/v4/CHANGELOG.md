Changelog from v4.0.0 and up. v3 changelog can be found in its branch.

# v4.1.4

* Fix bug in `Pool` involving blocking commands whose `Context` times out.
  (#344)

* Fix `Cluster.Clients` not returning correctly if there are no secondaries.
  (#345)

# v4.1.3

* Fixed bug in `Sentinel` where secondaries in the `s_down` state would still
  be included in the active set. (#343)

# v4.1.2

* Fixed `Sentinel` not creating connections to new secondaries properly. (#336)

* Complete refactor of `Conn`, code is simpler and a few cases where `Conn` was
  getting out of sync with its writes/reads are now handled properly.

* Fixed calls to `Unsubscribe` and `PSubscribe` not correctly clearing all
  subscriptions. (#318)

* Fixed the subscriptions made with `PSubscribe` not working correctly. (#333)

# v4.1.1

* Fixed `NewCluster` not returning an error if it can't connect to any of the
  redis instances given. (#319)

* Fix parsing for `CLUSTER SLOTS` command, which changed slightly with redis
  7.0. (#320)

* Fix a bug around discarding of errors in `Conn`. (#323)

* Properly handle the `MEMORY USAGE` command in the context of a cluster. (#325)

# v4.1.0

**New**

* Added `TreatErrorsAsValues` field to `resp.Opts`. (#309)

**Fixes and Improvements**

* Fixed PubSubMessage unmarshaling not correctly handling non-pubsub messages.
  (#306)

# v4.0.0

Below is documented all breaking changes between v3 and v4. There are further
enhancements which don't qualify as breaking changes which may not be documented
here.

**Major Changes**

* Stop using `...opts` pattern for optional parameters across all types, and
  switch instead to a `(Config{}).New` kind of pattern.

* Add `MultiClient` interface which is implemented by `Sentinel` and `Cluster`,
  `Client` has been modified to be implemented only by clients which point at a
  single redis instance (`Conn` and `Pool`). Methods on all affected
  client types have been modified to fit these new interfaces.

  * `Cluster.NewScanner` has been replaced by `ScannerConfig.NewMulti`.

* `Conn` has been completely re-designed. It is now always thread-safe. When
  multiple `Action`s are performed against a single `Conn` concurrently the
  `Conn` will automatically pipeline the `Action`'s read/writes, as appropriate.

  * `Pipeline` has been re-designed as a result as well.

  * `CmdAction` has been removed.

* `Pool` has been completely rewritten to better take advantage of connection
  sharing (previously called "implicit pipelining" in v3) and the new `Conn`
  design.

  * `EvalScript` and `Pipeline` now support connection sharing.

  * Since most `Action`s can be shared on the same `Conn` the `Pool` no longer
    runs the risk of being depleted during too many concurrent `Action`s, and so
    no longer needs to dynamically create/destroy `Conn`s.

  * A Pool size of 0 is no longer supported.

* Brand new `resp/resp3` package which implements the [RESP3][resp3] protocol.
  The new package features more consistent type mappings between go and redis
  and support for streaming types.

* Usage of `context.Context` in many places.

  * Add `context.Context` parameter to `Client.Do`, `PubSub` methods,
    `Scanner.Next`, and `WithConn`.

  * Add `context.Context` parameter to all `Client` and `Conn` creation functions.

  * Add `context.Context` parameter to `Action.Perform` (previously called
    `Action.Run`).

* The `PubSubConn` interface has been redesigned to be simpler to implement and
  use. Naming around pub/sub types has also been made more consistent.


**Minor Changes**

* Remove usage of `xerrors` package.

* Rename `resp.ErrDiscarded` to `resp.ErrConnUsable`, and change some of the
  semantics around using the error. A `resp.ErrConnUnusable` convenience
  function has been added as well.

* `resp.LenReader` now uses `int` instead of `int64` to signify length.

* `resp.Marshaler` and `resp.Unmarshaler` now take an `Opts` argument, to give
  the caller more control over things like byte pools and potentially other
  functionality in the future.

* `resp.Unmarshaler` now takes a `resp.BufferedReader`, rather than
  `*bufio.Reader`. Generally `resp.BufferedReader` will be implemented by a
  `*bufio.Reader`, but this gives more flexibility.

* `Stub` and `PubSubStub` have been renamed to `NewStubConn` and
  `NewPubSubStubConn`, respectively.

* Rename `MaybeNil` to just `Maybe`, and change its semantics a bit.

* The `trace` package has been significantly updated to reflect changes to
  `Pool` and other `Client`s.

* Refactor the `StreamReader` interface to be simpler to use.

[resp3]: https://github.com/antirez/RESP3
