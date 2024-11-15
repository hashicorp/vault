# UNRELEASED

# 1.7.0 (June 5th, 2024)

CHANGES

* Raft multi version testing [GH-559](https://github.com/hashicorp/raft/pull/559)

IMPROVEMENTS

* Raft pre-vote extension implementation, activated by default. [GH-530](https://github.com/hashicorp/raft/pull/530)

BUG FIXES

* Fix serialize NetworkTransport data race on ServerAddr(). [GH-591](https://github.com/hashicorp/raft/pull/591)

# 1.6.1 (January 8th, 2024)

CHANGES

* Add reference use of Hashicorp Raft. [GH-584](https://github.com/hashicorp/raft/pull/584)
* [COMPLIANCE] Add Copyright and License Headers. [GH-580](https://github.com/hashicorp/raft/pull/580)

IMPROVEMENTS

* Bump github.com/hashicorp/go-hclog from 1.5.0 to 1.6.2. [GH-583](https://github.com/hashicorp/raft/pull/583)

BUG FIXES

* Fix rare leadership transfer failures when writes happen during transfer. [GH-581](https://github.com/hashicorp/raft/pull/581)

# 1.6.0 (November 15th, 2023)

CHANGES

* Upgrade hashicorp/go-msgpack to v2, with go.mod upgraded from v0.5.5 to v2.1.1. [GH-577](https://github.com/hashicorp/raft/pull/577)

  go-msgpack v2.1.1 is by default binary compatible with v0.5.5 ("non-builtin" encoding of `time.Time`), but can decode messages produced by v1.1.5 as well ("builtin" encoding of `time.Time`).

  However, if users of this libary overrode the version of go-msgpack (especially to v1), this **could break** compatibility if raft nodes are running a mix of versions.

  This compatibility can be configured at runtime in Raft using `NetworkTransportConfig.MsgpackUseNewTimeFormat` -- the default is `false`, which maintains compatibility with `go-msgpack` v0.5.5, but if set to `true`, will be compatible with `go-msgpack` v1.1.5.

IMPROVEMENTS

* Push to notify channel when shutting down. [GH-567](https://github.com/hashicorp/raft/pull/567)
* Add CommitIndex API [GH-560](https://github.com/hashicorp/raft/pull/560)
* Document some Apply error cases better [GH-561](https://github.com/hashicorp/raft/pull/561)

BUG FIXES

* Race with `candidateFromLeadershipTransfer` [GH-570](https://github.com/hashicorp/raft/pull/570)


# 1.5.0 (April 21st, 2023)

IMPROVEMENTS
* Fixed a performance anomaly related to pipelining RPCs that caused large increases in commit latency under high write throughput. Default behavior has changed. For more information see #541.

# 1.4.0 (March 17th, 2023)

FEATURES
* Support log stores with a monotonically increasing index.  Implementing a log store with the `MonotonicLogStore` interface where `IsMonotonic()` returns true will allow Raft to clear all previous logs on user restores of Raft snapshots.

BUG FIXES
* Restoring a snapshot with the raft-wal log store caused a panic due to index gap that is created during snapshot restores.

# 1.3.0 (April 22nd, 2021)

IMPROVEMENTS

* Added metrics for `oldestLogAge` and `lastRestoreDuration` to monitor capacity issues that can cause unrecoverable cluster failure  [[GH-452](https://github.com/hashicorp/raft/pull/452)][[GH-454](https://github.com/hashicorp/raft/pull/454/files)]
* Made `TrailingLogs`, `SnapshotInterval` and `SnapshotThreshold` reloadable at runtime using a new `ReloadConfig` method. This allows recovery from cases where there are not enough logs retained for followers to catchup after a restart. [[GH-444](https://github.com/hashicorp/raft/pull/444)]
* Inclusify the repository by switching to main [[GH-446](https://github.com/hashicorp/raft/pull/446)]
* Add option for a buffered `ApplyCh` if `MaxAppendEntries` is enabled [[GH-445](https://github.com/hashicorp/raft/pull/445)]
* Add string to `LogType` for more human readable debugging [[GH-442](https://github.com/hashicorp/raft/pull/442)]
* Extract fuzzy testing into its own module [[GH-459](https://github.com/hashicorp/raft/pull/459)]

BUG FIXES
* Update LogCache `StoreLogs()` to capture an error that would previously cause a panic [[GH-460](https://github.com/hashicorp/raft/pull/460)]

# 1.2.0 (October 5th, 2020)

IMPROVEMENTS

* Remove `StartAsLeader` configuration option [[GH-364](https://github.com/hashicorp/raft/pull/386)]
* Allow futures to react to `Shutdown()` to prevent a deadlock with `takeSnapshot()` [[GH-390](https://github.com/hashicorp/raft/pull/390)]
* Prevent non-voters from becoming eligible for leadership elections [[GH-398](https://github.com/hashicorp/raft/pull/398)]
* Remove an unneeded `io.Copy` from snapshot writes [[GH-399](https://github.com/hashicorp/raft/pull/399)]
* Log decoded candidate address in `duplicate requestVote` warning [[GH-400](https://github.com/hashicorp/raft/pull/400)]
* Prevent starting a TCP transport when IP address is `nil` [[GH-403](https://github.com/hashicorp/raft/pull/403)]
* Reject leadership transfer requests when in candidate state to prevent indefinite blocking while unable to elect a leader [[GH-413](https://github.com/hashicorp/raft/pull/413)]
* Add labels for metric metadata to reduce cardinality of metric names [[GH-409](https://github.com/hashicorp/raft/pull/409)]
* Add peers metric [[GH-413](https://github.com/hashicorp/raft/pull/431)]

BUG FIXES

* Make `LeaderCh` always deliver the latest leadership transition [[GH-384](https://github.com/hashicorp/raft/pull/384)]
* Handle updating an existing peer in `startStopReplication` [[GH-419](https://github.com/hashicorp/raft/pull/419)]

# 1.1.2 (January 17th, 2020)

FEATURES

* Improve FSM apply performance through batching. Implementing the `BatchingFSM` interface enables this new feature [[GH-364](https://github.com/hashicorp/raft/pull/364)]
* Add ability to obtain Raft configuration before Raft starts with GetConfiguration [[GH-369](https://github.com/hashicorp/raft/pull/369)]

IMPROVEMENTS

* Remove lint violations and add a `make` rule for running the linter.
* Replace logger with hclog [[GH-360](https://github.com/hashicorp/raft/pull/360)]
* Read latest configuration independently from main loop [[GH-379](https://github.com/hashicorp/raft/pull/379)]

BUG FIXES

* Export the leader field in LeaderObservation [[GH-357](https://github.com/hashicorp/raft/pull/357)]
* Fix snapshot to not attempt to truncate a negative range [[GH-358](https://github.com/hashicorp/raft/pull/358)]
* Check for shutdown in inmemPipeline before sending RPCs [[GH-276](https://github.com/hashicorp/raft/pull/276)]

# 1.1.1 (July 23rd, 2019)

FEATURES

* Add support for extensions to be sent on log entries [[GH-353](https://github.com/hashicorp/raft/pull/353)]
* Add config option to skip snapshot restore on startup [[GH-340](https://github.com/hashicorp/raft/pull/340)]
* Add optional configuration store interface [[GH-339](https://github.com/hashicorp/raft/pull/339)]

IMPROVEMENTS

* Break out of group commit early when no logs are present [[GH-341](https://github.com/hashicorp/raft/pull/341)]

BUGFIXES

* Fix 64-bit counters on 32-bit platforms [[GH-344](https://github.com/hashicorp/raft/pull/344)]
* Don't defer closing source in recover/restore operations since it's in a loop [[GH-337](https://github.com/hashicorp/raft/pull/337)]

# 1.1.0 (May 23rd, 2019)

FEATURES

* Add transfer leadership extension [[GH-306](https://github.com/hashicorp/raft/pull/306)]

IMPROVEMENTS

* Move to `go mod` [[GH-323](https://github.com/hashicorp/consul/pull/323)]
* Leveled log [[GH-321](https://github.com/hashicorp/consul/pull/321)]
* Add peer changes to observations [[GH-326](https://github.com/hashicorp/consul/pull/326)]

BUGFIXES

* Copy the contents of an InmemSnapshotStore when opening a snapshot [[GH-270](https://github.com/hashicorp/consul/pull/270)]
* Fix logging panic when converting parameters to strings [[GH-332](https://github.com/hashicorp/consul/pull/332)]

# 1.0.1 (April 12th, 2019)

IMPROVEMENTS

* InMemTransport: Add timeout for sending a message [[GH-313](https://github.com/hashicorp/raft/pull/313)]
* ensure 'make deps' downloads test dependencies like testify [[GH-310](https://github.com/hashicorp/raft/pull/310)]
* Clarifies function of CommitTimeout [[GH-309](https://github.com/hashicorp/raft/pull/309)]
* Add additional metrics regarding log dispatching and committal [[GH-316](https://github.com/hashicorp/raft/pull/316)]

# 1.0.0 (October 3rd, 2017)

v1.0.0 takes the changes that were staged in the library-v2-stage-one branch. This version manages server identities using a UUID, so introduces some breaking API changes. It also versions the Raft protocol, and requires some special steps when interoperating with Raft servers running older versions of the library (see the detailed comment in config.go about version compatibility). You can reference https://github.com/hashicorp/consul/pull/2222 for an idea of what was required to port Consul to these new interfaces.

# 0.1.0 (September 29th, 2017)

v0.1.0 is the original stable version of the library that was in main and has been maintained with no breaking API changes. This was in use by Consul prior to version 0.7.0.
