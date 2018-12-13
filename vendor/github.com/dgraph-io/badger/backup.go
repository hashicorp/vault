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
	"io"
	"sync"

	"github.com/dgraph-io/badger/protos"
	"github.com/dgraph-io/badger/y"
)

func writeTo(entry *protos.KVPair, w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint64(entry.Size())); err != nil {
		return err
	}
	buf, err := entry.Marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// Backup dumps a protobuf-encoded list of all entries in the database into the
// given writer, that are newer than the specified version. It returns a
// timestamp indicating when the entries were dumped which can be passed into a
// later invocation to generate an incremental dump, of entries that have been
// added/modified since the last invocation of DB.Backup()
//
// This can be used to backup the data in a database at a given point in time.
func (db *DB) Backup(w io.Writer, since uint64) (uint64, error) {
	var tsNew uint64
	var skipKey []byte
	err := db.View(func(txn *Txn) error {
		opts := DefaultIteratorOptions
		opts.AllVersions = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.Version() < since || bytes.Equal(skipKey, item.Key()) {
				// Ignore versions less than given timestamp, or skip older
				// versions of the given skipKey.
				continue
			}
			skipKey = skipKey[:0]

			var valCopy []byte
			if !item.IsDeletedOrExpired() {
				// No need to copy value, if item is deleted or expired.
				var err error
				valCopy, err = item.ValueCopy(nil)
				if err != nil {
					Errorf("Key [%x, %d]. Error while fetching value [%v]\n",
						item.Key(), item.Version(), err)
					continue
				}
			}

			// clear txn bits
			meta := item.meta &^ (bitTxn | bitFinTxn)

			entry := &protos.KVPair{
				Key:       y.Copy(item.Key()),
				Value:     valCopy,
				UserMeta:  []byte{item.UserMeta()},
				Version:   item.Version(),
				ExpiresAt: item.ExpiresAt(),
				Meta:      []byte{meta},
			}
			if err := writeTo(entry, w); err != nil {
				return err
			}

			switch {
			case item.DiscardEarlierVersions():
				// If we need to discard earlier versions of this item, add a delete
				// marker just below the current version.
				entry.Version -= 1
				entry.Meta = []byte{bitDelete}
				if err := writeTo(entry, w); err != nil {
					return err
				}
				skipKey = item.KeyCopy(skipKey)

			case item.IsDeletedOrExpired():
				skipKey = item.KeyCopy(skipKey)
			}
		}
		tsNew = txn.readTs
		return nil
	})
	return tsNew, err
}

// Load reads a protobuf-encoded list of all entries from a reader and writes
// them to the database. This can be used to restore the database from a backup
// made by calling DB.Backup().
//
// DB.Load() should be called on a database that is not running any other
// concurrent transactions while it is running.
func (db *DB) Load(r io.Reader) error {
	br := bufio.NewReaderSize(r, 16<<10)
	unmarshalBuf := make([]byte, 1<<10)
	var entries []*Entry
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// func to check for pending error before sending off a batch for writing
	batchSetAsyncIfNoErr := func(entries []*Entry) error {
		select {
		case err := <-errChan:
			return err
		default:
			wg.Add(1)
			return db.batchSetAsync(entries, func(err error) {
				defer wg.Done()
				if err != nil {
					select {
					case errChan <- err:
					default:
					}
				}
			})
		}
	}

	for {
		var sz uint64
		err := binary.Read(br, binary.LittleEndian, &sz)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if cap(unmarshalBuf) < int(sz) {
			unmarshalBuf = make([]byte, sz)
		}

		e := &protos.KVPair{}
		if _, err = io.ReadFull(br, unmarshalBuf[:sz]); err != nil {
			return err
		}
		if err = e.Unmarshal(unmarshalBuf[:sz]); err != nil {
			return err
		}
		entries = append(entries, &Entry{
			Key:       y.KeyWithTs(e.Key, e.Version),
			Value:     e.Value,
			UserMeta:  e.UserMeta[0],
			ExpiresAt: e.ExpiresAt,
			meta:      e.Meta[0],
		})
		// Update nextTxnTs, memtable stores this timestamp in badger head
		// when flushed.
		if e.Version >= db.orc.nextTxnTs {
			db.orc.nextTxnTs = e.Version + 1
		}

		if len(entries) == 1000 {
			if err := batchSetAsyncIfNoErr(entries); err != nil {
				return err
			}
			entries = make([]*Entry, 0, 1000)
		}
	}

	if len(entries) > 0 {
		if err := batchSetAsyncIfNoErr(entries); err != nil {
			return err
		}
	}
	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		// Mark all versions done up until nextTxnTs.
		db.orc.txnMark.Done(db.orc.nextTxnTs - 1)
		return nil
	}
}
