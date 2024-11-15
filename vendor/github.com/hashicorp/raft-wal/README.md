# Raft WAL

This library implements a Write-Ahead Log (WAL) suitable for use with
[`hashicorp/raft`](https://github.com/hashicorp/raft).

Specifically the library provides and instance of raft's `LogStore` and
`StableStore` interfaces for storing both raft logs and the other small items
that require stable storage (like which term the node last voted in).

**This library is still considered experimental!** 

It is complete and reasonably well tested so far but we plan to complete more 
rigorous end-to-end testing and performance analysis within our products and 
together with some of our users before we consider this safe for production.

The advantage of this library over `hashicorp/raft-boltdb` that has been used
for many years in HashiCorp products are:
 1. Efficient truncations that don't cause later appends to slow down due to
    free space tracking issues in BoltDB's btree.
 2. More efficient appends due to only one fsync per append vs two in BoltDB.
 3. More efficient and suitable on-disk structure for a log vs a copy-on-write
    BTree.

We aim to provide roughly equivalent resiliency to crashes as respected storage
systems such as SQLite, LevelDB/RocksDB and etcd. BoltDB technically has a
stronger property due to it's page-aligned model (no partial sector overwrites).
We initially [designed a WAL on the same principals](/01-WAL-pages.md),
however felt that the additional complexity it adds wasn't justified given the
weaker assumptions that many other battle-tested systems above use.

Our design goals for crash recovery are:

 - Crashes at any point must not loose committed log entries or result in a
   corrupt file, even if in-flight sector writes are not atomic.
 - We _do_ assume [Powersafe Overwrites](#powersafe-overwrites-psow) where
   partial sectors can be appended to without corrupting existing data even in a
   power failure.
 - Latent errors (i.e. silent corruption in the FS or disk) _may_ be detected
   during a read, but we assume that the file-system and disk are responsible
   for this really. (i.e. we don't validate checksums on every record read).
   This is equivalent to SQLite, LMDB, BoltDB etc.

See the [system assumptions](#system-assumptions) and [crash
safety](#crash-safety) sections for more details.

## Limitations

Here are some notable (but we think acceptable) limitations of this design.

 * Segment files can't be larger than 4GiB. (Current default is 64MiB).
 * Individual records can't be larger than 4GiB without changing the format. 
   (Current limit is 64MiB).
 * Appended log entries must have monotonically increasing `Index` fields with
   no gaps (though may start at any index in an empty log).
 * Only head or tail truncations are supported. `DeleteRange` will error if the
   range is not a prefix of suffix of the log. `hashicorp/raft` never needs
   that.
 * No encryption or compression support.
   * Though we do provide a pluggable entry codec and internally treat each
     entry as opaque bytes so it's possible to supply a custom codec that
     transforms entries in any way desired.
 * If the segment tail file is lost _after_ entries are committed to it due to
   manual intervention or filesystem bug, the WAL can't distinguish that from a
   crash during rotation that left the file missing since we don't update
   metadata on every append for performance reasons. In most other cases,
   missing data would be detected on recovery and fail the recovery to protect
   from silent data loss, but in this particular case that's not possible
   without significantly impacting performance in the steady state by updating
   the last committed entry to meta DB on every append. We assume this is
   reasonable since previous LogStore implementations would also "silently"
   loose data if the database files were removed too.

## Storage Format Overview

The WAL has two types of file: a meta store and one or more log segments.

### Meta Store

We need to provide a `StableStore` interface for small amounts of Raft data. We
also need to store some meta data about log segments to simplify managing them
in an atomic and crash-safe way.

Since this data is _generally_ small we could invent our own storage format with
some sort of double-buffering and limit ourselves to a single page of data etc.
But since performance is not critical for meta-data operations and size is
extremely unlikely to get larger than a few KiB, we choose instead the pragmatic
approach of using BoltDB for our `wal-meta.db`.

The meta database contains two buckets: `stable` containing key/values persisted
by Raft via the `StableStore` interface, and `wal-state` which contains the
source-of-truth meta data about which segment files should be considered part of
the current log.

The `wal-state` bucket contains one record with all the state since it's only
loaded or persisted in one atomic batch and is small. The state is just a JSON
encoded object described by the following structs. JSON encoding is used as this
is not performance sensitive and it's simpler to work with and more human
readable.

```go
type PersistentState struct {
	NextSegmentID uint64
	Segments      []SegmentInfo
}
type SegmentInfo struct {
  ID         uint64
  BaseIndex  uint64
  MinIndex   uint64
  MaxIndex   uint64
  Codec      uint64
  IndexStart uint64
  CreateTime time.Time
  SealTime   time.Time
}
```

The last segment (with highest baseIndex) is the "tail" and must be the only one where
`SealTime = 0` (i.e. it's unsealed). `IndexStart` and `MaxIndex` are also zero until 
the segments is sealed.

Why use BoltDB when the main reason for this library is because the existing
BoltDB `LogStore` has performance issues?

Well, the major performance issue in `raft-boltdb` occurs when a large amount of
log data is written and then truncated, the overhead of tracking all the free
space in the file makes further appends slower.

Our use here is orders of magnitude lighter than storing all log data. As an
example, let's assume we allow 100GiB of logs to be kept around which is at
least an order of magnitude larger than the largest current known Consul user's
worst-case log size, and two orders of magnitude more than the largest Consul
deployments steady-state. Assuming fixed 64MiB segments, that would require
about 1600 segments which encode to about 125 bytes in JSON each. Even at this
extreme, the meta DB only has to hold under 200KiB.

Even if a truncation occurs that reduces that all the way back to a single
segment, 200KiB is only a hundred or so pages (allowing for btree overhead) so
the free list will never be larger than a single 4KB page.

On top of that, we only pay the cost of a write to BoltDB for meta-data
transactions: rotating to a new segment, or truncating. The vast majority of
appends only need to append to a log segment.

### Segment Files

Segment files are pre-allocated (if supported by the filesystem) on creation to 
a fixed size. By default we use 64MiB segment files. This sections defines the 
encoding for those files. All integer types are encoded in little-endian order.

The file starts with a fixed-size header that is written once before the 
first comitted entries.

```
0      1      2      3      4      5      6      7      8
+------+------+------+------+------+------+------+------+
| Magic                     | Reserved           | Vsn  |
+------+------+------+------+------+------+------+------+
| BaseIndex                                             |
+------+------+------+------+------+------+------+------+
| SegmentID                                             |
+------+------+------+------+------+------+------+------+
| Codec                                                 |
+------+------+------+------+------+------+------+------+
```

| Field        | Type      | Description |
| ------------ | --------- | ----------- |
| `Magic`      | `uint32`  | The randomly chosen value `0x58eb6b0d`. |
| `Reserved`   | `[3]byte` | Bytes reserved for future file flags. |
| `Vsn`        | `uint8`   | The version of the file, currently `0x0`. |
| `BaseIndex`  | `uint64`  | The raft Index of the first entry that will be stored in this file. |
| `SegmentID`  | `uint64`  | A unique identifier for this segment file. |
| `Codec`      | `uint64`  | The codec used to write the file. |

Each segment file is named `<BaseIndex>-<SegmentID>.wal`. `BaseIndex` is
formatted in decimal with leading zeros and a fixed width of 20 chars.
`SegmentID` is formatted in lower-case hex with zero padding to 16 chars wide.
This has the nice property of them sorting lexicographically in the directory,
although we don't rely on that.

### Frames

Log entries are stored in consecutive frames after the header. As well as log
entry frames there are a few meta data frame types too. Each frame starts with
an 8-byte header.

```
0      1      2      3      4      5      6      7      8
+------+------+------+------+------+------+------+------+
| Type | Reserved           | Length/CRC                |
+------+------+------+------+------+------+------+------+
```

| Field         | Type        | Description |
| ------------- | ----------- | ----------- |
| `Type`        | `uint8`     | The frame type. See below. |
| `Length/CRC`  | `uint32`    | Depends on Type. See Below |


| Type | Value | Description |
| ---- | ----- | ----------- |
| `Invalid` | `0x0` | The frame is invalid. We make zero value invalid so we can detect unwritten frames cleanly. |
| `Entry`   | `0x1` | The frame contains an entire log entry. |
| `Index`   | `0x2` | The frame contains an index array, not actual log entries. |
| `Commit`  | `0x3` | The frame contains a CRC for all data written in a batch. |

#### Index Frame

An index frame payload is an array of `uint32` file offsets for the 
correspoinding records. The first element of the array contains the file offset 
of the frame containing the first entry in the segment and so on.

`Length` is used to indicate the length in bytes of the array (i.e. number of
entries in the segments is `Length/4`).

Index frames are written only when the segment is sealed and a commit frame
follows to validate the final write.

#### Commit Frame

A Commit frame marks the last write before fsync is called. In order to detect
incomplete or torn writes on recovery the commit frame stores a CRC of all the
bytes appended since the last fsync.

`CRC` is used to specify a CRC32 (Castagnoli) over all bytes written since the
last fsync. That is, since just after the last commit frame, or just after the
file header.

There may also be 4 bytes of padding to keep alignment. Later we could
use these too.

#### Alignment

All frame headers are written with 8-byte alignment to ensure they remain in a
single disk sector. We don't entirely depend on atomic sector writes for
correctness, but it's a simple way to improve our chances or being able to read
through the file on a recovery with some sectors missing.

We add an implicit 0-7 null bytes after each frame to ensure the next frame
header is aligned. This padding is _not_ represented in `Length` but it is
always present and is deterministic by rounding up `Length` to the nearest
multiple of 8. It is always accounted for when reading and CRCs are calculated
over raw bytes written so always include the padding (zero) bytes.

Despite alignment we still don't blindly trust the headers we read are valid. A
CRC mismatch or invalid record format indicate torn writes in the last batch
written and we always safety check the size of lengths read before allocating
memory for them - Entry lengths can't be bigger than the `MaxEntrySize` which 
we default to 64MiB.

### Sealing

Once a segment file has grown larger than the configured soft-limit (64MiB
default), we "seal" it. This process involves:

 1. Write out the in-memory index of record offsets to an index frame.
 2. Write a commit frame to validate all bytes appended in this final append
    (which probably included one or more records that took the segment file over
    the limit).
 3. Return the final `IndexStart` to be stored in `wal-meta.db`

Sealed files can have their indexes read directly on open from the IndexStart in
`wal-meta.db` so records can be looked up in constant time.

## Log Lookup by Index

For an unsealed segment we first lookup the offset in the in-memory index.

For a sealed segment we can discover the index frame location from the metadata
and then perform a read at the right location in the file to lookup the record's
offset. Implementations may choose to cache or memory-map the index array but we
will initially just read the specific entry we need each time and assume the OS
page cache will make that fast for frequently accessed index areas or in-order
traversals. We don't have to read the whole index, just the 4 byte entry we care
about since we can work out it's offset from IndexStart, the BaseIndex of the
segment, and the Index being searched for.

# Crash Safety

Crash safety must be maintained through three type of write operation: appending
a batch of entries, truncating from the head (oldest) entries, and truncating
the newest entries.

## Appending Entries

We want to `fsync` only once for an append batch, however many entries were in
it. We assume [Powersafe Overwrites](#powersafe-overwrites-psow) or PSOW, a
weaker assumption than atomic sector writes in general. Thanks to PSOW, we
assume we can start appending at the tail of the file right after previously
committed entries even if the new entries will be written to the same sector as
the older entries, and that the system will never corrupt the already committed
part of the sector even if it is not atomic and arbitrarily garbles the part of
the sector we actually did write.

At the end of the batch we write a `Commit` frame containing the CRC over the
data written during the current batch.

In a crash one of the following states occurs:
 1. All sectors modified across all frames make it to disk (crash _after_ fsync).
 2. A torn write: one or more sectors, anywhere in the modified tail of the file
    might not be persisted. We don't assume they are zero, they might be
    arbitrarily garbled (crash _before_ fsync).

We can check which one of these is true with the recovery procedure outlined
below. If we find the last batch _was_ torn. It must not have been acknowledged
to Raft yet (since `fsync` can't have returned) and so it is safe to assume that
the previous commit frame is the tail of the log we've actually acknowledged.

### Recovery

We cover recovering the segments generally below since we have to account for
truncations. All segments except the tail were fsynced during seal before the
new tail was added to the meta DB so we can assume they are all made it to disk
if a later tail was added.

On startup we just need to recover the tail log as follows:

 1. If the file doesn't exist, create it from Meta DB information. DONE.
 2. Open file and validate header matches filename. If not delete it and go to 1.
 3. Read all records in the file in sequence, keeping track of the last two
    commit frames observed.
    1. If the file ends with a corrupt frame or non commit frame, discard
       anything after the last commit frame. We're DONE because we wouldn't have
       written extra frames after commit until fsync completed so this commit
       must have been acknowledged.
    1. Else the file ends with a commit frame. Validate its checksum. If it is good DONE.
    2. If CRC is not good then discard everything back to previous commit frame and DONE.
 4. If we read an index frame in that process and the commit frame proceeding it
    is the new tail then mark the segment as sealed and return the seal info
    (crash occured after seal but before updating `wal-meta.db`)

## Head Truncations

The most common form of truncation is a "head" truncation or removing the oldest
prefix of entries after a periodic snapshot has been made to reclaim space.

To be crash safe we can't rely on atomically updating or deleting multiple
segment files. The process looks like this.

 1. In one transaction on Meta DB:
    1. Update the `meta.min_index` to be the new min.
    2. Delete any segments from the `segments` bucket that are sealed and where
       their highest index is less than the new min index.
    3. Commit Txn. This is the commit point for crash recovery.
 2. Update in memory segment state to match (if not done already with a lock
    held).
 3. Delete any segment files we just removed from the meta DB.

### Recovery

The meta data update is crash safe thanks to BoltDB being the source of truth.

 1. Reload meta state from Meta DB.
 2. Walk the files in the dir.
 2. For each one:
    1. Check if that file is present in Meta DB. If not mark it for deletion.
    2. (optionally) validate the file header file size and final block trailer
       to ensure the file appears to be well-formed and contain the expected
       data.
 4. Delete the obsolete segments marked (could be done in a background thread).

 ## Tail Truncations

 Raft occasionally needs to truncate entries from the tail of the log, i.e.
 remove the _most recent_ N entries. This can occur when a follower has
 replicated entries from an old leader that was partitioned with it, but later
 discovers they conflict with entries committed by the new leader in a later
 term. The bounds on how long a partitioned leader can continue to replicate to
 a follower are generally pretty small (30 seconds or so) so it's unlikely that
 the number of records to be truncated will ever be large compared to the size
 of a segment file, but we have to account for needing to delete one or more
 segment files from the tail, as well as truncate older entries out of the new
 tail.

 This follows roughly the same pattern as head-truncation, although there is an
 added complication. A naive implementation that used only the baseIndex as a
 segment file name could in theory get into a tricky state where it's ambiguous
 whether the tail segment is an old one that was logically truncated away but we
 crashed before actually unlinking, or a new replacement with committed data in.

 It's possible to solve this with complex transactional semantics but we take
 the simpler approach of just assigning every segment a unique identifier
 separate from it's baseIndex. So to truncate the tail follows the same
 procedure as the head above: segments we remove from Meta DB can be
 un-ambiguously deleted on recovery because their IDs won't match even if later
 segments end up with the same baseIndex.

 Since these truncations are generally rare and disk space is generally not a
 major bottleneck, we also choose not to try to actually re-use a segment file
 that was previously written and sealed by truncating it etc. Instead we just
 mark it as "sealed" in the Meta DB and with a MaxIndex of the highest index
 left after the truncation (which we check on reads) and start a new segment at
 the next index.

## System Assumptions

There are no straight answers to any question about which guarantees can be
reliably relied on across operating systems, file systems, raid controllers and
hardware devices. We state [our assumptions](#our-assumptions) followed by a
summary of the assumptions made by some other respectable sources for
comparison.

### Our Assumptions

We've tried to make the weakest assumptions we can while still keeping things
relatively simple and performant.

We assume:
 1. That while silent latent errors are possible, they are generally rare and
    there's not a whole lot we can do other than return a `Corrupt` error on
    read. In most cases the hardware or filesystem will detect and return an
    error on read anyway for latent corruption. Not doing so is regarded as a
    bug in the OS/filesystem/hardware. For this reason we don't go out of our
    way to checksum everything to protect against "bitrot". This is roughly
    equivalent to assumptions in BoltDB, LMDB and SQLite.

    While we respect the work in [Protocol Aware Recovery for Consensus-based
    Storage](https://www.usenix.org/system/files/conference/fast18/fast18-alagappan.pdf)
    we choose not to implement a WAL format that allows identifying the index
    and term of "lost" records on read errors so they can be recovered from
    peers. This is mostly for the pragmatic reason that the Raft library this is
    designed to work with would need a major re-write to take advantage of that
    anyway. The proposed format in that paper also seems to make stronger
    assumptions about sector atomicity than we are comfortable with too.
 2. That sector writes are _not_ atomic. (Equivalent to SQLite, weaker than
    almost everything else).
 3. That writing a partial sector does _not_ corrupt any already stored data in
    that sector outside of the range being written (
    [PSOW](#powersafe-overwrites-psow)), (Equivalent to SQLite's defaults,
    RocksDB and Etcd).
 3. That `fsync` as implemented in Go's standard library actually flushes all
    written sectors of the file to persistent media.
 4. That `fsync` on a parent dir is sufficient to ensure newly created files are
    not lost after a crash (assuming the file itself was written and `fsync`ed
    first).
 6. That appending to files may not be atomic since the filesystem metadata
    about the size of the file may not be updated atomically with the data.
    Generally we pre-allocate files where possible without writing all zeros but
    we do potentially extend them if the last batch doesn't fit into the
    allocated space or the filesystem doesn't support pre-allocation. Either way
    we don't rely on the filesystem's reported size and validate the tail is
    coherent on recovery.

### Published Paper on Consensus Disk Recovery

In the paper on [Protocol Aware Recovery for Consensus-based
Storage](https://www.usenix.org/system/files/conference/fast18/fast18-alagappan.pdf)
the authors assume that corruptions of the log can happen due to either torn
writes (for multi-sector appends) or latent corruptions after commit. They
explain the need to detect which it was because torn writes only loose
un-acknowledged records and so are safe to detect and truncate, while corruption
of previously committed records impacts the correctness of the protocol more
generally. Their whole paper seems to indicate that these post-commit
corruptions are a major problem that needs to be correctly handled (which may
well be true). On the flip side, their WAL format design writes a separate index
and log, and explicitly assumes that because the index entries are smaller than
a 512 sector size, that those are safe from corruption during a write.

The core assumptions here are:
  1. Latent, silent corruption of committed data needs to be detected at
     application layer with a checksum per record checked on every read.
  2. Sector writes are atomic.
  3. Sector writes have [powersafe overwrites](#powersafe-overwrites-psow).

### SQLite

The SQLite authors have a [detailed explanation of their system
assumptions](https://www.sqlite.org/atomiccommit.html) which impact correctness
of atomic database commits.

> SQLite assumes that the detection and/or correction of bit errors caused by cosmic rays, thermal noise, quantum fluctuations, device driver bugs, or other mechanisms, is the responsibility of the underlying hardware and operating system. SQLite does not add any redundancy to the database file for the purpose of detecting corruption or I/O errors. SQLite assumes that the data it reads is exactly the same data that it previously wrote.

Is very different from the above paper authors whose main point of their paper
is predicated on how to recover from silent corruptions of the file caused by
hardware, firmware or filesystem errors on read.

Note that this is a pragmatic position rather than a naive one: the authors are
certainly aware that file-systems have bugs, that faulty raid controllers exist
and even that hardware anomalies like high-flying or poorly tracking disk heads
can happen but choose _not_ to protect against that _at all_. See their
[briefing for linux kernel
developers](https://sqlite.org/lpc2019/doc/trunk/briefing.md) for more details
on the uncertainty they understand exists around these areas.

> SQLite has traditionally assumed that a sector write is not atomic.

These statements are on a page with this disclaimer:

> The information in this article applies only when SQLite is operating in "rollback mode", or in other words when SQLite is not using a write-ahead log.

[WAL mode](https://sqlite.org/wal.html) docs are less explicit on assumptions
and how crash recovery is achieved but we can infer some things from the [file
format](https://sqlite.org/fileformat2.html#walformat) and
[code](https://github.com/sqlite/sqlite/blob/master/src/wal.c) though.

> The WAL header is 32 bytes in size...

> Immediately following the wal-header are zero or more frames. Each frame consists of a 24-byte frame-header followed by a page-size bytes of page data.

So each dirty page is appended with a 24 byte header making it _not_ sector
aligned even though pages must be a multiple of sector size.

Commit frames are also appended in the same way (and fsync called if enabled as
an option). If fsync is enabled though (and POWERSAFE_OVERWRITE disabled),
SQLite will "pad" to the next sector boundary (or beyond) by repeating the last
frame until it's passed that boundary. For some reason, they take great care to
write up to the sector boundary, sync then write the rest. I assume this is just
to avoid waiting to flush the redundant padding bytes past the end of the sector
they care about. Padding prevents the next append from potentially overwriting
the committed frame's sector.

But...

> By default, SQLite assumes that an operating system call to write a range of bytes will not damage or alter any bytes outside of that range even if a power loss or OS crash occurs during that write. We call this the "powersafe overwrite" property. Prior to version 3.7.9 (2011-11-01), SQLite did not assume powersafe overwrite. But with the standard sector size increasing from 512 to 4096 bytes on most disk drives, it has become necessary to assume powersafe overwrite in order to maintain historical performance levels and so powersafe overwrite is assumed by default in recent versions of SQLite.

> [assuming no power safe overwrite] In WAL mode, each transaction had to be padded out to the next 4096-byte boundary in the WAL file, rather than the next 512-byte boundary, resulting in thousands of extra bytes being written per transaction.

> SQLite never assumes that database page writes are atomic, regardless of the PSOW setting.(1) And hence SQLite is always able to automatically recover from torn pages induced by a crash. Enabling PSOW does not decrease SQLite's ability to recover from a torn page.

So they basically changed to make SSDs performant and now assume _by default_
that appending to a partial sector won't damage other data. The authors are
explicit that ["powersafe overwrite"](#powersafe-overwrites-psow) is a separate
property from atomicity and they still don't rely on sector atomicity. But they
do now assume powersafe overwrites by default.

To summarize, SQLite authors assume:
 1. Latent, silent corruptions of committed data should be caught by the file
    system or hardware and so should't need to be accounted for in application
    code.
 2. Sector writes are _not_ atomic, but...
 3. Partial sector overwrites can't corrupt committed data in same sector (by
    default).

### Etcd WAL

The authors of etcd's WAL similarly to the authors of the above paper indicate
the need to distinguish between torn writes and silent corruptions.

They maintain a rolling checksum of all records which is used on recovery only
which would imply they only care about torn writes since per-record checksums
are not checked on subsequent reads from the file after recovery. But they have
specific code to distinguish between torn writes and "other" corruption during
recovery.

They are careful to pad every record with 0 to 7 bytes such that the length
prefix for the next record is always 8-byte aligned and so can't span more than
one segment.

But their method of detecting a torn-write (rather than latent corruption)
relies on reading through every 512 byte aligned slice of the set of records
whose checksum has failed to match and seeing if there are any entirely zero
sectors.

This seems problematic in a purely logical way regardless of disk behavior: if a
legitimate record contains more than 1kb of zero bytes and happens to ever be
corrupted after writing, that record will be falsely detected as a torn-write
because at least one sector will be entirely zero bytes. In practice this
doesn't matter much because corruptions caused by anything other than torn
writes are likely very rare but it does make me wonder why bother trying to tell
the difference.

The implied assumptions in their design are:
 1. Latent, silent corruption needs to be detected on recovery, but not on every
    read.
 2. Sector writes are atomic.
 3. Partial sector writes don't corrupt existing data.
 3. Torn writes (caused by multi-segment appends) always leave sectors all-zero.

### LMDB

Symas' Lightning Memory-mapped Database or LMDB is another well-used and
respected DB file format (along with Go-native port BoltDB used by Consul,
etcd and others).

LMDB writes exclusively in whole 4kb pages. LMDB has a copy-on-write design
which reuses free pages and commits transactions using the a double-buffering
technique: writing the root alternately to the first and second pages of the
file. Individual pages do not have checksums and may be larger than the physical
sector size. Dirty pages are written out to new or un-used pages and then
`fsync`ed before the transaction commits so there is no reliance on atomic
sector writes for data pages (a crash might leave pages of a transaction
partially written but they are not linked into the tree root yet so are ignored
on recovery).

The transaction commits only after the double-buffered meta page is written
out. LMDB relies on the fact that the actual content of the meta page is small
enough to fit in a single sector to avoid "torn writes" on the meta page. (See
[the authors
comments](https://ayende.com/blog/162856/reviewing-lightning-memory-mapped-database-library-transactions-commits)
on this blog). Although sector writes are assumed to be atomic, there is no
reliance on partial sector writes due to the paged design.

The implied assumptions in this design are:
 1. Latent, silent corruptions of committed data should be caught by the file
    system or hardware and so should't need to be accounted for in application
    code.
 2. Sector writes _are_ atomic.
 3. No assumptions about Powersafe overwrite since all IO is in whole pages.


### BoltDB

BoltDB is a Go port of LMDB so inherits almost all of the same design
assumptions. One notable different is that the author added a checksum to
metadata page even though it still fits in a single sector. The author noted
in private correspondence that this was probably just a defensive measure
rather than a fix for a specific identified flaw in LMDB's design.

Initially this was _not_ used to revert to the alternate page on failure because
it was still assumed that meta fit in a single sector and that those writes were
atomic. But [a report of Docker causing corruption on a
crash](https://github.com/boltdb/bolt/issues/548) seemed to indicate that the
atomic sector writes assumption _was not safe_ alone and so the checksum was
used to detect non-atomic writes even on the less-than-a-sector meta page.

BoltDB is also an important base case for our WAL since it is used as the
current log store in use for many years within Consul and other HashiCorp
products.

The implied assumptions in this design are:
1. Latent, silent corruptions of committed data should be caught by the file
   system or hardware and so should't need to be accounted for in application
   code.
2. Sector writes are _not_ atomic.
3. No assumptions about Powersafe overwrite since all IO is in whole pages.


### RocksDB WAL

RocksDB is another well-respected storage library based on Google's LevelDB.
RocksDB's [WAL
Format](https://github.com/facebook/rocksdb/wiki/Write-Ahead-Log-File-Format)
uses blocks to allow skipping through files and over corrupt records (which
seems dangerous to me in general but perhaps they assume only torn-write
corruptions are possible?).

Records are packed into 32KiB blocks until they don't fit. Records that are
larger use first/middle/last flags (which inspired this library) to consume
multiple blocks.

RocksDB WAL uses pre-allocated files but also re-uses old files on a circular
buffer pattern since they have tight control of how much WAL is needed. This
means they might be overwriting old records in place.

Each record independently gets a header with a checksum to detect corruption or
incomplete writes, but no attempt is made to avoid sector boundaries or partial
block writes - the current block is just appended to for each write.

Implied assumptions:
 1. No Latent Corruptions? This isn't totally clear from the code or docs, but
    the docs indicate that a record with a mismatching checksum can simply be
    skipped over which would seem to violate basic durability properties for a
    database if they were already committed. That would imply that checksums
    only (correctly) detect torn writes with latent corruption not accounted
    for.
 2. Sector writes _are_ atomic.
 3. Partial sector writes don't corrupt existing data.

### Are Sector Writes Atomic?

Russ Cox asked this on twitter and tweeted a link to an [excellent Stack
Overflow
answer](https://stackoverflow.com/questions/2009063/are-disk-sector-writes-atomic)
about this by one of the authors of the NVME spec.

> TLDR; if you are in tight control of your whole stack from application all the way down the the physical disks (so you can control and qualify the whole lot) you can arrange to have what you need to make use of disk atomicity. If you're not in that situation or you're talking about the general case, you should not depend on sector writes being atomic.

Despite this, _most_ current best-of-breed database libraries (notably except
SQLite and potentially BoltDB), [many linux file
systems](https://lkml.org/lkml/2009/8/24/156), and all academic papers on disk
failure modes I've found so far _do_ assume that sector writes are atomic.

I assume that the authors of these file systems, databases and papers are not
unaware of the complexities described in the above link or the possibility of
non-atomic sector writes, but rather have chosen to put those outside of the
reasonable recoverable behavior of their systems. The actual chances of
encountering a non-atomic sector write in a typical, modern system appear to be
small enough that these authors consider that a reasonable assumption even when
it's not a guarantee that can be 100% relied upon. (Although the Docker bug
linked above for [BoltDB](#boltdb) seems to indicate a real-world case of this
happening in a modern environment.)

### Powersafe Overwrites (PSOW)

A more subtle property that is a weaker assumption that full sector atomicity is
termed by the [SQLite authors as "Powersafe
Overwrites"](https://www.sqlite.org/psow.html) abbreviated PSOW.

> By default, SQLite assumes that an operating system call to write a range of bytes will not damage or alter any bytes outside of that range even if a power loss or OS crash occurs during that write. We call this the "powersafe overwrite" property. Prior to version 3.7.9 (2011-11-01), SQLite did not assume powersafe overwrite. But with the standard sector size increasing from 512 to 4096 bytes on most disk drives, it has become necessary to assume powersafe overwrite in order to maintain historical performance levels and so powersafe overwrite is assumed by default in recent versions of SQLite.

Those who assume atomic sector writes _also_ assume this property but the
reverse need not be true. SQLite's authors in the page above assume nothing
about the atomicity of the actual data written to any sector still even when
POWERSAFE_OVERWRITE is enabled (which is now the default). They simply assume
that no _other_ data is harmed while performing a write that overlaps other
sectors, even if power fails.

It's our view that while there certainly can be cases where this assumption
doesn't hold, it's already weaker than the atomic sector write assumption that
most reliable storage software assumes today and so is safe to assume on for
this case.

### Are fsyncs reliable?

Even when you explicitly `fsync` a file after writing to it, some devices or
even whole operating systems (e.g. macOS) _don't actually flush to disk_ to
improve performance.

In our case, we assume that Go's `os.File.Sync()` method is makes the best
effort it can on all modern OSes. It does now at least behave correctly on macOS
(since Go 1.12). But we can't do anything about a lying hardware device.

# Future Extensions

 * **Auto-tuning segment size.** This format allows for segments to be different
   sizes. We could start with a smaller segment size of say a single 1MiB block
   and then measure how long it takes to fill each segment. If segments fill
   quicker than some target rate we could double the allocated size of the next
   segment. This could mean a burst of writes makes the segments grow and then
   the writes slow down but the log would take a long time to free disk space
   because the segments take so long to fill. Arguably not a terrible problem,
   but we could also have it auto tune segment size down when write rate drops
   too. The only major benefit here would be to allow trivial usages like tests
   not need a whole 64MiB of disk space to just record a handful of log entries.
   But those could also just manually configure a smaller segment size.

# References

In no particular order.

**Files and Crash Recovery**
* [Files are hard](https://danluu.com/file-consistency/)
* [Files are fraught with peril](https://danluu.com/deconstruct-files/)
* [Ensuring data reaches disk](https://lwn.net/Articles/457667/)
* [Write Atomicity and NVME Device Design](https://www.bswd.com/FMS12/FMS12-Rudoff.pdf)
* [Durability: NVME Disks](https://www.evanjones.ca/durability-nvme.html)
* [Intel SSD Durability](https://www.evanjones.ca/intel-ssd-durability.html)
* [Are Disk Sector Writes Atomic?](https://stackoverflow.com/questions/2009063/are-disk-sector-writes-atomic/61832882#61832882)
* [Protocol Aware Recovery for Consensus-based Storage](https://www.usenix.org/system/files/conference/fast18/fast18-alagappan.pdf)
* [Atomic Commit in SQLite](https://www.sqlite.org/atomiccommit.html)
* ["Powersafe Overwrites" in SQLite](https://www.sqlite.org/psow.html)
* [An Analysis of Data Corruption in the Storage Stack](https://www.cs.toronto.edu/~bianca/papers/fast08.pdf)

**DB Design and Storage File layout**
* [BoltDB Implementation](https://github.com/boltdb/bolt)
* LMDB Design: [slides](https://www.snia.org/sites/default/files/SDC15_presentations/database/HowardChu_The_Lighting_Memory_Database.pdf), [talk](https://www.youtube.com/watch?v=tEa5sAh-kVk)
* [SQLite file layout](https://www.sqlite.org/fileformat.html)

**WAL implementations**
* [SQLite WAL Mode](https://sqlite.org/wal.html)
* [RocksDB WAL Format](https://github.com/facebook/rocksdb/wiki/Write-Ahead-Log-File-Format)]
* [etcd implementation](https://github.com/etcd-io/etcd/tree/master/wal)
