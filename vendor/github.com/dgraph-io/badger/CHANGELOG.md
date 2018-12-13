# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.5.3] - 2018-07-11
Bug Fixes:
* Fix a panic caused due to item.vptr not copying over vs.Value, when looking
    for a move key.

## [1.5.2] - 2018-06-19
Bug Fixes:
* Fix the way move key gets generated.
* If a transaction has unclosed, or multiple iterators running simultaneously,
    throw a panic. Every iterator must be properly closed. At any point in time,
    only one iterator per transaction can be running. This is to avoid bugs in a
    transaction data structure which is thread unsafe.

* *Warning: This change might cause panics in user code. Fix is to properly
    close your iterators, and only have one running at a time per transaction.*

## [1.5.1] - 2018-06-04
Bug Fixes:
* Fix for infinite yieldItemValue recursion. #503
* Fix recursive addition of `badgerMove` prefix. https://github.com/dgraph-io/badger/commit/2e3a32f0ccac3066fb4206b28deb39c210c5266f
* Use file size based window size for sampling, instead of fixing it to 10MB. #501

Cleanup:
* Clarify comments and documentation.
* Move badger tool one directory level up.

## [1.5.0] - 2018-05-08
* Introduce `NumVersionsToKeep` option. This option is used to discard many
  versions of the same key, which saves space.
* Add a new `SetWithDiscard` method, which would indicate that all the older
  versions of the key are now invalid. Those versions would be discarded during
  compactions.
* Value log GC moves are now bound to another keyspace to ensure latest versions
  of data are always at the top in LSM tree.
* Introduce `ValueLogMaxEntries` to restrict the number of key-value pairs per
  value log file. This helps bound the time it takes to garbage collect one
  file.

## [1.4.0] - 2018-05-04
* Make mmap-ing of value log optional.
* Run GC multiple times, based on recorded discard statistics.
* Add MergeOperator.
* Force compact L0 on clsoe (#439).
* Add truncate option to warn about data loss (#452).
* Discard key versions during compaction (#464).
* Introduce new `LSMOnlyOptions`, to make Badger act like a typical LSM based DB.

Bug fix:
* (Temporary) Check max version across all tables in Get (removed in next
  release).
* Update commit and read ts while loading from backup.
* Ensure all transaction entries are part of the same value log file.
* On commit, run unlock callbacks before doing writes (#413).
* Wait for goroutines to finish before closing iterators (#421).

## [1.3.0] - 2017-12-12
* Add `DB.NextSequence()` method to generate monotonically increasing integer
  sequences.
* Add `DB.Size()` method to return the size of LSM and value log files.
* Tweaked mmap code to make Windows 32-bit builds work.
* Tweaked build tags on some files to make iOS builds work.
* Fix `DB.PurgeOlderVersions()` to not violate some constraints.

## [1.2.0] - 2017-11-30
* Expose a `Txn.SetEntry()` method to allow setting the key-value pair
  and all the metadata at the same time.

## [1.1.1] - 2017-11-28
* Fix bug where txn.Get was returing key deleted in same transaction.
* Fix race condition while decrementing reference in oracle.
* Update doneCommit in the callback for CommitAsync.
* Iterator see writes of current txn.

## [1.1.0] - 2017-11-13
* Create Badger directory if it does not exist when `badger.Open` is called.
* Added `Item.ValueCopy()` to avoid deadlocks in long-running iterations
* Fixed 64-bit alignment issues to make Badger run on Arm v7

## [1.0.1] - 2017-11-06
* Fix an uint16 overflow when resizing key slice

[Unreleased]: https://github.com/dgraph-io/badger/compare/v1.5.3...HEAD
[1.5.3]: https://github.com/dgraph-io/badger/compare/v1.5.2...v1.5.3
[1.5.2]: https://github.com/dgraph-io/badger/compare/v1.5.1...v1.5.2
[1.5.1]: https://github.com/dgraph-io/badger/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/dgraph-io/badger/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/dgraph-io/badger/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/dgraph-io/badger/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/dgraph-io/badger/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/dgraph-io/badger/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/dgraph-io/badger/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/dgraph-io/badger/compare/v1.0.0...v1.0.1
