/*
 * Copyright 2017 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package badger

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/options"
	"github.com/dgraph-io/badger/y"
	"github.com/pkg/errors"
	"golang.org/x/net/trace"
)

// Values have their first byte being byteData or byteDelete. This helps us distinguish between
// a key that has never been seen and a key that has been explicitly deleted.
const (
	bitDelete                 byte = 1 << 0 // Set if the key has been deleted.
	bitValuePointer           byte = 1 << 1 // Set if the value is NOT stored directly next to key.
	bitDiscardEarlierVersions byte = 1 << 2 // Set if earlier versions can be discarded.

	// The MSB 2 bits are for transactions.
	bitTxn    byte = 1 << 6 // Set if the entry is part of a txn.
	bitFinTxn byte = 1 << 7 // Set if the entry is to indicate end of txn in value log.

	mi int64 = 1 << 20
)

type logFile struct {
	path string
	// This is a lock on the log file. It guards the fd’s value, the file’s
	// existence and the file’s memory map.
	//
	// Use shared ownership when reading/writing the file or memory map, use
	// exclusive ownership to open/close the descriptor, unmap or remove the file.
	lock        sync.RWMutex
	fd          *os.File
	fid         uint32
	fmap        []byte
	size        uint32
	loadingMode options.FileLoadingMode
}

// openReadOnly assumes that we have a write lock on logFile.
func (lf *logFile) openReadOnly() error {
	var err error
	lf.fd, err = os.OpenFile(lf.path, os.O_RDONLY, 0666)
	if err != nil {
		return errors.Wrapf(err, "Unable to open %q as RDONLY.", lf.path)
	}

	fi, err := lf.fd.Stat()
	if err != nil {
		return errors.Wrapf(err, "Unable to check stat for %q", lf.path)
	}
	y.AssertTrue(fi.Size() <= math.MaxUint32)
	lf.size = uint32(fi.Size())

	if err = lf.mmap(fi.Size()); err != nil {
		_ = lf.fd.Close()
		return y.Wrapf(err, "Unable to map file")
	}

	return nil
}

func (lf *logFile) mmap(size int64) (err error) {
	if lf.loadingMode != options.MemoryMap {
		// Nothing to do
		return nil
	}
	lf.fmap, err = y.Mmap(lf.fd, false, size)
	if err == nil {
		err = y.Madvise(lf.fmap, false) // Disable readahead
	}
	return err
}

func (lf *logFile) munmap() (err error) {
	if lf.loadingMode != options.MemoryMap {
		// Nothing to do
		return nil
	}
	if err := y.Munmap(lf.fmap); err != nil {
		return errors.Wrapf(err, "Unable to munmap value log: %q", lf.path)
	}
	return nil
}

// Acquire lock on mmap/file if you are calling this
func (lf *logFile) read(p valuePointer, s *y.Slice) (buf []byte, err error) {
	var nbr int64
	offset := p.Offset
	if lf.loadingMode == options.FileIO {
		buf = s.Resize(int(p.Len))
		var n int
		n, err = lf.fd.ReadAt(buf, int64(offset))
		nbr = int64(n)
	} else {
		// Do not convert size to uint32, because the lf.fmap can be of size
		// 4GB, which overflows the uint32 during conversion to make the size 0,
		// causing the read to fail with ErrEOF. See issue #585.
		size := int64(len(lf.fmap))
		valsz := p.Len
		if int64(offset) >= size || int64(offset+valsz) > size {
			err = y.ErrEOF
		} else {
			buf = lf.fmap[offset : offset+valsz]
			nbr = int64(valsz)
		}
	}
	y.NumReads.Add(1)
	y.NumBytesRead.Add(nbr)
	return buf, err
}

func (lf *logFile) doneWriting(offset uint32) error {
	// Sync before acquiring lock.  (We call this from write() and thus know we have shared access
	// to the fd.)
	if err := lf.fd.Sync(); err != nil {
		return errors.Wrapf(err, "Unable to sync value log: %q", lf.path)
	}
	// Close and reopen the file read-only.  Acquire lock because fd will become invalid for a bit.
	// Acquiring the lock is bad because, while we don't hold the lock for a long time, it forces
	// one batch of readers wait for the preceding batch of readers to finish.
	//
	// If there's a benefit to reopening the file read-only, it might be on Windows.  I don't know
	// what the benefit is.  Consider keeping the file read-write, or use fcntl to change
	// permissions.
	lf.lock.Lock()
	defer lf.lock.Unlock()
	if err := lf.munmap(); err != nil {
		return err
	}
	// TODO: Confirm if we need to run a file sync after truncation.
	// Truncation must run after unmapping, otherwise Windows would crap itself.
	if err := lf.fd.Truncate(int64(offset)); err != nil {
		return errors.Wrapf(err, "Unable to truncate file: %q", lf.path)
	}
	if err := lf.fd.Close(); err != nil {
		return errors.Wrapf(err, "Unable to close value log: %q", lf.path)
	}

	return lf.openReadOnly()
}

// You must hold lf.lock to sync()
func (lf *logFile) sync() error {
	return lf.fd.Sync()
}

var errStop = errors.New("Stop iteration")
var errTruncate = errors.New("Do truncate")

type logEntry func(e Entry, vp valuePointer) error

type safeRead struct {
	k []byte
	v []byte

	recordOffset uint32
}

func (r *safeRead) Entry(reader *bufio.Reader) (*Entry, error) {
	var hbuf [headerBufSize]byte
	var err error

	hash := crc32.New(y.CastagnoliCrcTable)
	tee := io.TeeReader(reader, hash)
	if _, err = io.ReadFull(tee, hbuf[:]); err != nil {
		return nil, err
	}

	var h header
	h.Decode(hbuf[:])
	if h.klen > uint32(1<<16) { // Key length must be below uint16.
		return nil, errTruncate
	}
	kl := int(h.klen)
	if cap(r.k) < kl {
		r.k = make([]byte, 2*kl)
	}
	vl := int(h.vlen)
	if cap(r.v) < vl {
		r.v = make([]byte, 2*vl)
	}

	e := &Entry{}
	e.offset = r.recordOffset
	e.Key = r.k[:kl]
	e.Value = r.v[:vl]

	if _, err = io.ReadFull(tee, e.Key); err != nil {
		if err == io.EOF {
			err = errTruncate
		}
		return nil, err
	}
	if _, err = io.ReadFull(tee, e.Value); err != nil {
		if err == io.EOF {
			err = errTruncate
		}
		return nil, err
	}
	var crcBuf [4]byte
	if _, err = io.ReadFull(reader, crcBuf[:]); err != nil {
		if err == io.EOF {
			err = errTruncate
		}
		return nil, err
	}
	crc := binary.BigEndian.Uint32(crcBuf[:])
	if crc != hash.Sum32() {
		return nil, errTruncate
	}
	e.meta = h.meta
	e.UserMeta = h.userMeta
	e.ExpiresAt = h.expiresAt
	return e, nil
}

// iterate iterates over log file. It doesn't not allocate new memory for every kv pair.
// Therefore, the kv pair is only valid for the duration of fn call.
func (vlog *valueLog) iterate(lf *logFile, offset uint32, fn logEntry) (uint32, error) {
	fi, err := lf.fd.Stat()
	if err != nil {
		return 0, err
	}
	if int64(offset) == fi.Size() {
		// We're at the end of the file already. No need to do anything.
		return offset, nil
	}
	if vlog.opt.ReadOnly {
		// We're not at the end of the file. We'd need to replay the entries, or
		// possibly truncate the file.
		return 0, ErrReplayNeeded
	}

	// We're not at the end of the file. Let's Seek to the offset and start reading.
	if _, err := lf.fd.Seek(int64(offset), io.SeekStart); err != nil {
		return 0, errFile(err, lf.path, "Unable to seek")
	}

	reader := bufio.NewReader(lf.fd)
	read := &safeRead{
		k:            make([]byte, 10),
		v:            make([]byte, 10),
		recordOffset: offset,
	}

	var lastCommit uint64
	var validEndOffset uint32
	for {
		e, err := read.Entry(reader)
		if err == io.EOF {
			break
		} else if err == io.ErrUnexpectedEOF || err == errTruncate {
			break
		} else if err != nil {
			return 0, err
		} else if e == nil {
			continue
		}

		var vp valuePointer
		vp.Len = uint32(headerBufSize + len(e.Key) + len(e.Value) + crc32.Size)
		read.recordOffset += vp.Len

		vp.Offset = e.offset
		vp.Fid = lf.fid

		if e.meta&bitTxn > 0 {
			txnTs := y.ParseTs(e.Key)
			if lastCommit == 0 {
				lastCommit = txnTs
			}
			if lastCommit != txnTs {
				break
			}

		} else if e.meta&bitFinTxn > 0 {
			txnTs, err := strconv.ParseUint(string(e.Value), 10, 64)
			if err != nil || lastCommit != txnTs {
				break
			}
			// Got the end of txn. Now we can store them.
			lastCommit = 0
			validEndOffset = read.recordOffset

		} else {
			if lastCommit != 0 {
				// This is most likely an entry which was moved as part of GC.
				// We shouldn't get this entry in the middle of a transaction.
				break
			}
			validEndOffset = read.recordOffset
		}

		if err := fn(*e, vp); err != nil {
			if err == errStop {
				break
			}
			return 0, errFile(err, lf.path, "Iteration function")
		}
	}
	return validEndOffset, nil
}

func (vlog *valueLog) rewrite(f *logFile, tr trace.Trace) error {
	maxFid := atomic.LoadUint32(&vlog.maxFid)
	y.AssertTruef(uint32(f.fid) < maxFid, "fid to move: %d. Current max fid: %d", f.fid, maxFid)
	tr.LazyPrintf("Rewriting fid: %d", f.fid)

	wb := make([]*Entry, 0, 1000)
	var size int64

	y.AssertTrue(vlog.db != nil)
	var count, moved int
	fe := func(e Entry) error {
		count++
		if count%100000 == 0 {
			tr.LazyPrintf("Processing entry %d", count)
		}

		vs, err := vlog.db.get(e.Key)
		if err != nil {
			return err
		}
		if discardEntry(e, vs) {
			return nil
		}

		// Value is still present in value log.
		if len(vs.Value) == 0 {
			return errors.Errorf("Empty value: %+v", vs)
		}
		var vp valuePointer
		vp.Decode(vs.Value)

		if vp.Fid > f.fid {
			return nil
		}
		if vp.Offset > e.offset {
			return nil
		}
		if vp.Fid == f.fid && vp.Offset == e.offset {
			moved++
			// This new entry only contains the key, and a pointer to the value.
			ne := new(Entry)
			ne.meta = 0 // Remove all bits. Different keyspace doesn't need these bits.
			ne.UserMeta = e.UserMeta

			// Create a new key in a separate keyspace, prefixed by moveKey. We are not
			// allowed to rewrite an older version of key in the LSM tree, because then this older
			// version would be at the top of the LSM tree. To work correctly, reads expect the
			// latest versions to be at the top, and the older versions at the bottom.
			if bytes.HasPrefix(e.Key, badgerMove) {
				ne.Key = append([]byte{}, e.Key...)
			} else {
				ne.Key = append([]byte{}, badgerMove...)
				ne.Key = append(ne.Key, e.Key...)
			}

			ne.Value = append([]byte{}, e.Value...)
			wb = append(wb, ne)
			size += int64(e.estimateSize(vlog.opt.ValueThreshold))
			if size >= 64*mi {
				tr.LazyPrintf("request has %d entries, size %d", len(wb), size)
				if err := vlog.db.batchSet(wb); err != nil {
					return err
				}
				size = 0
				wb = wb[:0]
			}
		} else {
			Warningf("This entry should have been caught. %+v\n", e)
		}
		return nil
	}

	_, err := vlog.iterate(f, 0, func(e Entry, vp valuePointer) error {
		return fe(e)
	})
	if err != nil {
		return err
	}

	tr.LazyPrintf("request has %d entries, size %d", len(wb), size)
	batchSize := 1024
	var loops int
	for i := 0; i < len(wb); {
		loops++
		if batchSize == 0 {
			Warningf("We shouldn't reach batch size of zero.")
			return ErrNoRewrite
		}
		end := i + batchSize
		if end > len(wb) {
			end = len(wb)
		}
		if err := vlog.db.batchSet(wb[i:end]); err != nil {
			if err == ErrTxnTooBig {
				// Decrease the batch size to half.
				batchSize = batchSize / 2
				tr.LazyPrintf("Dropped batch size to %d", batchSize)
				continue
			}
			return err
		}
		i += batchSize
	}
	tr.LazyPrintf("Processed %d entries in %d loops", len(wb), loops)
	tr.LazyPrintf("Total entries: %d. Moved: %d", count, moved)
	tr.LazyPrintf("Removing fid: %d", f.fid)
	var deleteFileNow bool
	// Entries written to LSM. Remove the older file now.
	{
		vlog.filesLock.Lock()
		// Just a sanity-check.
		if _, ok := vlog.filesMap[f.fid]; !ok {
			vlog.filesLock.Unlock()
			return errors.Errorf("Unable to find fid: %d", f.fid)
		}
		if vlog.numActiveIterators == 0 {
			delete(vlog.filesMap, f.fid)
			deleteFileNow = true
		} else {
			vlog.filesToBeDeleted = append(vlog.filesToBeDeleted, f.fid)
		}
		vlog.filesLock.Unlock()
	}

	if deleteFileNow {
		vlog.deleteLogFile(f)
	}

	return nil
}

func (vlog *valueLog) deleteMoveKeysFor(fid uint32, tr trace.Trace) error {
	db := vlog.db
	var result []*Entry
	var count, pointers uint64
	tr.LazyPrintf("Iterating over move keys to find invalids for fid: %d", fid)
	err := db.View(func(txn *Txn) error {
		opt := DefaultIteratorOptions
		opt.internalAccess = true
		opt.PrefetchValues = false
		itr := txn.NewIterator(opt)
		defer itr.Close()

		for itr.Seek(badgerMove); itr.ValidForPrefix(badgerMove); itr.Next() {
			count++
			item := itr.Item()
			if item.meta&bitValuePointer == 0 {
				continue
			}
			pointers++
			var vp valuePointer
			vp.Decode(item.vptr)
			if vp.Fid == fid {
				e := &Entry{Key: y.KeyWithTs(item.Key(), item.Version()), meta: bitDelete}
				result = append(result, e)
			}
		}
		return nil
	})
	if err != nil {
		tr.LazyPrintf("Got error while iterating move keys: %v", err)
		tr.SetError()
		return err
	}
	tr.LazyPrintf("Num total move keys: %d. Num pointers: %d", count, pointers)
	tr.LazyPrintf("Number of invalid move keys found: %d", len(result))
	batchSize := 10240
	for i := 0; i < len(result); {
		end := i + batchSize
		if end > len(result) {
			end = len(result)
		}
		if err := db.batchSet(result[i:end]); err != nil {
			if err == ErrTxnTooBig {
				batchSize /= 2
				tr.LazyPrintf("Dropped batch size to %d", batchSize)
				continue
			}
			tr.LazyPrintf("Error while doing batchSet: %v", err)
			tr.SetError()
			return err
		}
		i += batchSize
	}
	tr.LazyPrintf("Move keys deletion done.")
	return nil
}

func (vlog *valueLog) incrIteratorCount() {
	atomic.AddInt32(&vlog.numActiveIterators, 1)
}

func (vlog *valueLog) decrIteratorCount() error {
	num := atomic.AddInt32(&vlog.numActiveIterators, -1)
	if num != 0 {
		return nil
	}

	vlog.filesLock.Lock()
	lfs := make([]*logFile, 0, len(vlog.filesToBeDeleted))
	for _, id := range vlog.filesToBeDeleted {
		lfs = append(lfs, vlog.filesMap[id])
		delete(vlog.filesMap, id)
	}
	vlog.filesToBeDeleted = nil
	vlog.filesLock.Unlock()

	for _, lf := range lfs {
		if err := vlog.deleteLogFile(lf); err != nil {
			return err
		}
	}
	return nil
}

func (vlog *valueLog) deleteLogFile(lf *logFile) error {
	path := vlog.fpath(lf.fid)
	if err := lf.munmap(); err != nil {
		_ = lf.fd.Close()
		return err
	}
	if err := lf.fd.Close(); err != nil {
		return err
	}
	return os.Remove(path)
}

// lfDiscardStats keeps track of the amount of data that could be discarded for
// a given logfile.
type lfDiscardStats struct {
	sync.Mutex
	m map[uint32]int64
}

type valueLog struct {
	buf     bytes.Buffer
	dirPath string
	elog    trace.EventLog

	// guards our view of which files exist, which to be deleted, how many active iterators
	filesLock        sync.RWMutex
	filesMap         map[uint32]*logFile
	filesToBeDeleted []uint32
	// A refcount of iterators -- when this hits zero, we can delete the filesToBeDeleted.
	numActiveIterators int32

	db                *DB
	maxFid            uint32 // accessed via atomics.
	writableLogOffset uint32 // read by read, written by write. Must access via atomics.
	numEntriesWritten uint32
	opt               Options

	garbageCh      chan struct{}
	lfDiscardStats *lfDiscardStats
}

func vlogFilePath(dirPath string, fid uint32) string {
	return fmt.Sprintf("%s%s%06d.vlog", dirPath, string(os.PathSeparator), fid)
}

func (vlog *valueLog) fpath(fid uint32) string {
	return vlogFilePath(vlog.dirPath, fid)
}

func (vlog *valueLog) populateFilesMap() error {
	vlog.filesMap = make(map[uint32]*logFile)

	files, err := ioutil.ReadDir(vlog.dirPath)
	if err != nil {
		return errFile(err, vlog.dirPath, "Unable to open log dir.")
	}

	found := make(map[uint64]struct{})
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".vlog") {
			continue
		}
		fsz := len(file.Name())
		fid, err := strconv.ParseUint(file.Name()[:fsz-5], 10, 32)
		if err != nil {
			return errFile(err, file.Name(), "Unable to parse log id.")
		}
		if _, ok := found[fid]; ok {
			return errFile(err, file.Name(), "Duplicate file found. Please delete one.")
		}
		found[fid] = struct{}{}

		lf := &logFile{
			fid:         uint32(fid),
			path:        vlog.fpath(uint32(fid)),
			loadingMode: vlog.opt.ValueLogLoadingMode,
		}
		vlog.filesMap[uint32(fid)] = lf
		if vlog.maxFid < uint32(fid) {
			vlog.maxFid = uint32(fid)
		}
	}
	return nil
}

func (vlog *valueLog) createVlogFile(fid uint32) (*logFile, error) {
	path := vlog.fpath(fid)
	lf := &logFile{
		fid:         fid,
		path:        path,
		loadingMode: vlog.opt.ValueLogLoadingMode,
	}
	// writableLogOffset is only written by write func, by read by Read func.
	// To avoid a race condition, all reads and updates to this variable must be
	// done via atomics.
	atomic.StoreUint32(&vlog.writableLogOffset, 0)
	vlog.numEntriesWritten = 0

	var err error
	if lf.fd, err = y.CreateSyncedFile(path, vlog.opt.SyncWrites); err != nil {
		return nil, errFile(err, lf.path, "Create value log file")
	}
	if err = syncDir(vlog.dirPath); err != nil {
		return nil, errFile(err, vlog.dirPath, "Sync value log dir")
	}
	if err = lf.mmap(2 * vlog.opt.ValueLogFileSize); err != nil {
		return nil, errFile(err, lf.path, "Mmap value log file")
	}

	vlog.filesLock.Lock()
	vlog.filesMap[fid] = lf
	vlog.filesLock.Unlock()

	return lf, nil
}

func errFile(err error, path string, msg string) error {
	return fmt.Errorf("%s. Path=%s. Error=%v", msg, path, err)
}

func (vlog *valueLog) replayLog(lf *logFile, offset uint32, replayFn logEntry) error {
	// We should open the file in RW mode, so it can be truncated.
	var err error
	lf.fd, err = os.OpenFile(lf.path, os.O_RDWR, 0)
	if err != nil {
		return errFile(err, lf.path, "Open file in RW mode")
	}
	defer lf.fd.Close()

	fi, err := lf.fd.Stat()
	if err != nil {
		return errFile(err, lf.path, "Unable to run file.Stat")
	}

	// Alright, let's iterate now.
	endOffset, err := vlog.iterate(lf, offset, replayFn)
	if err != nil {
		return errFile(err, lf.path, "Unable to replay logfile")
	}
	if int64(endOffset) == fi.Size() {
		return nil
	}

	// End offset is different from file size. So, we should truncate the file
	// to that size.
	y.AssertTrue(int64(endOffset) <= fi.Size())
	if !vlog.opt.Truncate {
		return ErrTruncateNeeded
	}
	if err := lf.fd.Truncate(int64(endOffset)); err != nil {
		return errFile(err, lf.path, fmt.Sprintf(
			"Truncation needed at offset %d. Can be done manually as well.", endOffset))
	}
	return nil
}

func (vlog *valueLog) open(db *DB, ptr valuePointer, replayFn logEntry) error {
	opt := db.opt
	vlog.opt = opt
	vlog.dirPath = opt.ValueDir
	vlog.db = db
	vlog.elog = trace.NewEventLog("Badger", "Valuelog")
	vlog.garbageCh = make(chan struct{}, 1) // Only allow one GC at a time.
	vlog.lfDiscardStats = &lfDiscardStats{m: make(map[uint32]int64)}

	if err := vlog.populateFilesMap(); err != nil {
		return err
	}
	// If no files are found, then create a new file.
	if len(vlog.filesMap) == 0 {
		_, err := vlog.createVlogFile(0)
		return err
	}

	fids := vlog.sortedFids()
	for _, fid := range fids {
		lf, ok := vlog.filesMap[fid]
		y.AssertTrue(ok)

		// This file is before the value head pointer. So, we don't need to
		// replay it, and can just open it in readonly mode.
		if fid < ptr.Fid {
			if err := lf.openReadOnly(); err != nil {
				return err
			}
			continue
		}

		var offset uint32
		if fid == ptr.Fid {
			offset = ptr.Offset + ptr.Len
		}
		Infof("Replaying file id: %d at offset: %d\n", fid, offset)
		now := time.Now()
		// Replay and possible truncation done. Now we can open the file as per
		// user specified options.
		if err := vlog.replayLog(lf, offset, replayFn); err != nil {
			return err
		}
		Infof("Replay took: %s\n", time.Since(now))

		if fid < vlog.maxFid {
			if err := lf.openReadOnly(); err != nil {
				return err
			}
		} else {
			var flags uint32
			switch {
			case vlog.opt.ReadOnly:
				// If we have read only, we don't need SyncWrites.
				flags |= y.ReadOnly
			case vlog.opt.SyncWrites:
				flags |= y.Sync
			}
			var err error
			if lf.fd, err = y.OpenExistingFile(vlog.fpath(fid), flags); err != nil {
				return errFile(err, lf.path, "Open existing file")
			}
		}
	}

	// Seek to the end to start writing.
	last, ok := vlog.filesMap[vlog.maxFid]
	y.AssertTrue(ok)
	lastOffset, err := last.fd.Seek(0, io.SeekEnd)
	if err != nil {
		return errFile(err, last.path, "file.Seek to end")
	}
	vlog.writableLogOffset = uint32(lastOffset)
	// Map the file if needed. When we create a file, it is automatically mapped.
	if err = last.mmap(2 * opt.ValueLogFileSize); err != nil {
		return errFile(err, last.path, "Map log file")
	}
	return nil
}

func (vlog *valueLog) Close() error {
	vlog.elog.Printf("Stopping garbage collection of values.")
	defer vlog.elog.Finish()

	var err error
	for id, f := range vlog.filesMap {
		f.lock.Lock() // We won’t release the lock.
		if munmapErr := f.munmap(); munmapErr != nil && err == nil {
			err = munmapErr
		}

		maxFid := atomic.LoadUint32(&vlog.maxFid)
		if !vlog.opt.ReadOnly && id == maxFid {
			// truncate writable log file to correct offset.
			if truncErr := f.fd.Truncate(
				int64(vlog.woffset())); truncErr != nil && err == nil {
				err = truncErr
			}
		}

		if closeErr := f.fd.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}

// sortedFids returns the file id's not pending deletion, sorted.  Assumes we have shared access to
// filesMap.
func (vlog *valueLog) sortedFids() []uint32 {
	toBeDeleted := make(map[uint32]struct{})
	for _, fid := range vlog.filesToBeDeleted {
		toBeDeleted[fid] = struct{}{}
	}
	ret := make([]uint32, 0, len(vlog.filesMap))
	for fid := range vlog.filesMap {
		if _, ok := toBeDeleted[fid]; !ok {
			ret = append(ret, fid)
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret
}

type request struct {
	// Input values
	Entries []*Entry
	// Output values and wait group stuff below
	Ptrs []valuePointer
	Wg   sync.WaitGroup
	Err  error
}

func (req *request) Wait() error {
	req.Wg.Wait()
	req.Entries = nil
	err := req.Err
	requestPool.Put(req)
	return err
}

// sync is thread-unsafe and should not be called concurrently with write.
func (vlog *valueLog) sync() error {
	if vlog.opt.SyncWrites {
		return nil
	}

	vlog.filesLock.RLock()
	if len(vlog.filesMap) == 0 {
		vlog.filesLock.RUnlock()
		return nil
	}
	maxFid := atomic.LoadUint32(&vlog.maxFid)
	curlf := vlog.filesMap[maxFid]
	curlf.lock.RLock()
	vlog.filesLock.RUnlock()

	dirSyncCh := make(chan error)
	go func() { dirSyncCh <- syncDir(vlog.opt.ValueDir) }()
	err := curlf.sync()
	curlf.lock.RUnlock()
	dirSyncErr := <-dirSyncCh
	if err != nil {
		err = dirSyncErr
	}
	return err
}

func (vlog *valueLog) woffset() uint32 {
	return atomic.LoadUint32(&vlog.writableLogOffset)
}

// write is thread-unsafe by design and should not be called concurrently.
func (vlog *valueLog) write(reqs []*request) error {
	vlog.filesLock.RLock()
	maxFid := atomic.LoadUint32(&vlog.maxFid)
	curlf := vlog.filesMap[maxFid]
	vlog.filesLock.RUnlock()

	toDisk := func() error {
		if vlog.buf.Len() == 0 {
			return nil
		}
		vlog.elog.Printf("Flushing %d blocks of total size: %d", len(reqs), vlog.buf.Len())
		n, err := curlf.fd.Write(vlog.buf.Bytes())
		if err != nil {
			return errors.Wrapf(err, "Unable to write to value log file: %q", curlf.path)
		}
		y.NumWrites.Add(1)
		y.NumBytesWritten.Add(int64(n))
		vlog.elog.Printf("Done")
		atomic.AddUint32(&vlog.writableLogOffset, uint32(n))
		vlog.buf.Reset()

		if vlog.woffset() > uint32(vlog.opt.ValueLogFileSize) ||
			vlog.numEntriesWritten > vlog.opt.ValueLogMaxEntries {
			var err error
			if err = curlf.doneWriting(vlog.writableLogOffset); err != nil {
				return err
			}

			newid := atomic.AddUint32(&vlog.maxFid, 1)
			y.AssertTruef(newid > 0, "newid has overflown uint32: %v", newid)
			newlf, err := vlog.createVlogFile(newid)
			if err != nil {
				return err
			}
			curlf = newlf
		}
		return nil
	}

	for i := range reqs {
		b := reqs[i]
		b.Ptrs = b.Ptrs[:0]
		for j := range b.Entries {
			e := b.Entries[j]
			var p valuePointer

			p.Fid = curlf.fid
			// Use the offset including buffer length so far.
			p.Offset = vlog.woffset() + uint32(vlog.buf.Len())
			plen, err := encodeEntry(e, &vlog.buf) // Now encode the entry into buffer.
			if err != nil {
				return err
			}
			p.Len = uint32(plen)
			b.Ptrs = append(b.Ptrs, p)
		}
		vlog.numEntriesWritten += uint32(len(b.Entries))
		// We write to disk here so that all entries that are part of the same transaction are
		// written to the same vlog file.
		writeNow :=
			vlog.woffset()+uint32(vlog.buf.Len()) > uint32(vlog.opt.ValueLogFileSize) ||
				vlog.numEntriesWritten > uint32(vlog.opt.ValueLogMaxEntries)
		if writeNow {
			if err := toDisk(); err != nil {
				return err
			}
		}
	}
	return toDisk()

	// Acquire mutex locks around this manipulation, so that the reads don't try to use
	// an invalid file descriptor.
}

// Gets the logFile and acquires and RLock() for the mmap. You must call RUnlock on the file
// (if non-nil)
func (vlog *valueLog) getFileRLocked(fid uint32) (*logFile, error) {
	vlog.filesLock.RLock()
	defer vlog.filesLock.RUnlock()
	ret, ok := vlog.filesMap[fid]
	if !ok {
		// log file has gone away, will need to retry the operation.
		return nil, ErrRetry
	}
	ret.lock.RLock()
	return ret, nil
}

// Read reads the value log at a given location.
// TODO: Make this read private.
func (vlog *valueLog) Read(vp valuePointer, s *y.Slice) ([]byte, func(), error) {
	// Check for valid offset if we are reading to writable log.
	maxFid := atomic.LoadUint32(&vlog.maxFid)
	if vp.Fid == maxFid && vp.Offset >= vlog.woffset() {
		return nil, nil, errors.Errorf(
			"Invalid value pointer offset: %d greater than current offset: %d",
			vp.Offset, vlog.woffset())
	}

	buf, cb, err := vlog.readValueBytes(vp, s)
	if err != nil {
		return nil, cb, err
	}
	var h header
	h.Decode(buf)
	n := uint32(headerBufSize) + h.klen
	return buf[n : n+h.vlen], cb, nil
}

func (vlog *valueLog) readValueBytes(vp valuePointer, s *y.Slice) ([]byte, func(), error) {
	lf, err := vlog.getFileRLocked(vp.Fid)
	if err != nil {
		return nil, nil, err
	}

	buf, err := lf.read(vp, s)
	if vlog.opt.ValueLogLoadingMode == options.MemoryMap {
		return buf, lf.lock.RUnlock, err
	}
	// If we are using File I/O we unlock the file immediately
	// and return an empty function as callback.
	lf.lock.RUnlock()
	return buf, nil, err
}

// Test helper
func valueBytesToEntry(buf []byte) (e Entry) {
	var h header
	h.Decode(buf)
	n := uint32(headerBufSize)

	e.Key = buf[n : n+h.klen]
	n += h.klen
	e.meta = h.meta
	e.UserMeta = h.userMeta
	e.Value = buf[n : n+h.vlen]
	return
}

func (vlog *valueLog) pickLog(head valuePointer, tr trace.Trace) (files []*logFile) {
	vlog.filesLock.RLock()
	defer vlog.filesLock.RUnlock()
	fids := vlog.sortedFids()
	if len(fids) <= 1 {
		tr.LazyPrintf("Only one or less value log file.")
		return nil
	} else if head.Fid == 0 {
		tr.LazyPrintf("Head pointer is at zero.")
		return nil
	}

	// Pick a candidate that contains the largest amount of discardable data
	candidate := struct {
		fid     uint32
		discard int64
	}{math.MaxUint32, 0}
	vlog.lfDiscardStats.Lock()
	for _, fid := range fids {
		if fid >= head.Fid {
			break
		}
		if vlog.lfDiscardStats.m[fid] > candidate.discard {
			candidate.fid = fid
			candidate.discard = vlog.lfDiscardStats.m[fid]
		}
	}
	vlog.lfDiscardStats.Unlock()

	if candidate.fid != math.MaxUint32 { // Found a candidate
		tr.LazyPrintf("Found candidate via discard stats: %v", candidate)
		files = append(files, vlog.filesMap[candidate.fid])
	} else {
		tr.LazyPrintf("Could not find candidate via discard stats. Randomly picking one.")
	}

	// Fallback to randomly picking a log file
	var idxHead int
	for i, fid := range fids {
		if fid == head.Fid {
			idxHead = i
			break
		}
	}
	if idxHead == 0 { // Not found or first file
		tr.LazyPrintf("Could not find any file.")
		return nil
	}
	idx := rand.Intn(idxHead) // Don’t include head.Fid. We pick a random file before it.
	if idx > 0 {
		idx = rand.Intn(idx + 1) // Another level of rand to favor smaller fids.
	}
	tr.LazyPrintf("Randomly chose fid: %d", fids[idx])
	files = append(files, vlog.filesMap[fids[idx]])
	return files
}

func discardEntry(e Entry, vs y.ValueStruct) bool {
	if vs.Version != y.ParseTs(e.Key) {
		// Version not found. Discard.
		return true
	}
	if isDeletedOrExpired(vs.Meta, vs.ExpiresAt) {
		return true
	}
	if (vs.Meta & bitValuePointer) == 0 {
		// Key also stores the value in LSM. Discard.
		return true
	}
	if (vs.Meta & bitFinTxn) > 0 {
		// Just a txn finish entry. Discard.
		return true
	}
	return false
}

func (vlog *valueLog) doRunGC(lf *logFile, discardRatio float64, tr trace.Trace) (err error) {
	// Update stats before exiting
	defer func() {
		if err == nil {
			vlog.lfDiscardStats.Lock()
			delete(vlog.lfDiscardStats.m, lf.fid)
			vlog.lfDiscardStats.Unlock()
		}
	}()

	type reason struct {
		total   float64
		discard float64
		count   int
	}

	fi, err := lf.fd.Stat()
	if err != nil {
		tr.LazyPrintf("Error while finding file size: %v", err)
		tr.SetError()
		return err
	}

	// Set up the sampling window sizes.
	sizeWindow := float64(fi.Size()) * 0.1                          // 10% of the file as window.
	countWindow := int(float64(vlog.opt.ValueLogMaxEntries) * 0.01) // 1% of num entries.
	tr.LazyPrintf("Size window: %5.2f. Count window: %d.", sizeWindow, countWindow)

	// Pick a random start point for the log.
	skipFirstM := float64(rand.Int63n(fi.Size())) // Pick a random starting location.
	skipFirstM -= sizeWindow                      // Avoid hitting EOF by moving back by window.
	skipFirstM /= float64(mi)                     // Convert to MBs.
	tr.LazyPrintf("Skip first %5.2f MB of file of size: %d MB", skipFirstM, fi.Size()/mi)
	var skipped float64

	var r reason
	start := time.Now()
	y.AssertTrue(vlog.db != nil)
	s := new(y.Slice)
	var numIterations int
	_, err = vlog.iterate(lf, 0, func(e Entry, vp valuePointer) error {
		numIterations++
		esz := float64(vp.Len) / (1 << 20) // in MBs.
		if skipped < skipFirstM {
			skipped += esz
			return nil
		}

		// Sample until we reach the window sizes or exceed 10 seconds.
		if r.count > countWindow {
			tr.LazyPrintf("Stopping sampling after %d entries.", countWindow)
			return errStop
		}
		if r.total > sizeWindow {
			tr.LazyPrintf("Stopping sampling after reaching window size.")
			return errStop
		}
		if time.Since(start) > 10*time.Second {
			tr.LazyPrintf("Stopping sampling after 10 seconds.")
			return errStop
		}
		r.total += esz
		r.count++

		vs, err := vlog.db.get(e.Key)
		if err != nil {
			return err
		}
		if discardEntry(e, vs) {
			r.discard += esz
			return nil
		}

		// Value is still present in value log.
		y.AssertTrue(len(vs.Value) > 0)
		vp.Decode(vs.Value)

		if vp.Fid > lf.fid {
			// Value is present in a later log. Discard.
			r.discard += esz
			return nil
		}
		if vp.Offset > e.offset {
			// Value is present in a later offset, but in the same log.
			r.discard += esz
			return nil
		}
		if vp.Fid == lf.fid && vp.Offset == e.offset {
			// This is still the active entry. This would need to be rewritten.

		} else {
			vlog.elog.Printf("Reason=%+v\n", r)

			buf, cb, err := vlog.readValueBytes(vp, s)
			if err != nil {
				return errStop
			}
			ne := valueBytesToEntry(buf)
			ne.offset = vp.Offset
			ne.print("Latest Entry Header in LSM")
			e.print("Latest Entry in Log")
			runCallback(cb)
			return errors.Errorf("This shouldn't happen. Latest Pointer:%+v. Meta:%v.",
				vp, vs.Meta)
		}
		return nil
	})

	if err != nil {
		tr.LazyPrintf("Error while iterating for RunGC: %v", err)
		tr.SetError()
		return err
	}
	tr.LazyPrintf("Fid: %d. Skipped: %5.2fMB Num iterations: %d. Data status=%+v\n",
		lf.fid, skipped, numIterations, r)

	// If we couldn't sample at least a 1000 KV pairs or at least 75% of the window size,
	// and what we can discard is below the threshold, we should skip the rewrite.
	if (r.count < countWindow && r.total < sizeWindow*0.75) || r.discard < discardRatio*r.total {
		tr.LazyPrintf("Skipping GC on fid: %d", lf.fid)
		return ErrNoRewrite
	}
	if err = vlog.rewrite(lf, tr); err != nil {
		return err
	}
	tr.LazyPrintf("Done rewriting.")
	return nil
}

func (vlog *valueLog) waitOnGC(lc *y.Closer) {
	defer lc.Done()

	<-lc.HasBeenClosed() // Wait for lc to be closed.

	// Block any GC in progress to finish, and don't allow any more writes to runGC by filling up
	// the channel of size 1.
	vlog.garbageCh <- struct{}{}
}

func (vlog *valueLog) runGC(discardRatio float64, head valuePointer) error {
	select {
	case vlog.garbageCh <- struct{}{}:
		// Pick a log file for GC.
		tr := trace.New("Badger.ValueLog", "GC")
		tr.SetMaxEvents(100)
		defer func() {
			tr.Finish()
			<-vlog.garbageCh
		}()

		var err error
		files := vlog.pickLog(head, tr)
		if len(files) == 0 {
			tr.LazyPrintf("PickLog returned zero results.")
			return ErrNoRewrite
		}
		tried := make(map[uint32]bool)
		for _, lf := range files {
			if _, done := tried[lf.fid]; done {
				continue
			}
			tried[lf.fid] = true
			err = vlog.doRunGC(lf, discardRatio, tr)
			if err == nil {
				return vlog.deleteMoveKeysFor(lf.fid, tr)
			}
		}
		return err
	default:
		return ErrRejected
	}
}

func (vlog *valueLog) updateGCStats(stats map[uint32]int64) {
	vlog.lfDiscardStats.Lock()
	for fid, sz := range stats {
		vlog.lfDiscardStats.m[fid] += sz
	}
	vlog.lfDiscardStats.Unlock()
}
