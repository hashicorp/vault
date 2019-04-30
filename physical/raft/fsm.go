package raft

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	protoio "github.com/gogo/protobuf/io"
	proto "github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	bolt "go.etcd.io/bbolt"
)

const (
	deleteOp uint32 = 1 << iota
	putOp
)

var (
	// bucketName is the value we use for the bucket
	bucketName = []byte("data")
)

// Verify FSM satisfies the correct interfaces
var _ physical.Backend = (*FSM)(nil)
var _ physical.Transactional = (*FSM)(nil)

type FSMApplyResponse struct {
	Success bool
}

type FSM struct {
	l          sync.RWMutex
	path       string
	logger     log.Logger
	permitPool *physical.PermitPool

	db *bolt.DB
}

// NewFSM constructs a FSM using the given directory
func NewFSM(conf map[string]string, logger log.Logger) (*FSM, error) {
	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	dbPath := filepath.Join(path, "vault.db")

	boltDB, err := bolt.Open(dbPath, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// make sure we have a bucket created
	err = boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &FSM{
		path:       path,
		logger:     logger,
		permitPool: physical.NewPermitPool(physical.DefaultParallelOperations),

		db: boltDB,
	}, nil
}

func (b *FSM) Delete(ctx context.Context, path string) error {
	defer metrics.MeasureSince([]string{"raft", "delete"}, time.Now())

	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.l.RLock()
	defer b.l.RUnlock()

	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Delete([]byte(path))
	})
}

func (b *FSM) Get(ctx context.Context, path string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"raft", "get"}, time.Now())

	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.l.RLock()
	defer b.l.RUnlock()

	var valCopy []byte
	var found bool

	err := b.db.View(func(tx *bolt.Tx) error {

		value := tx.Bucket(bucketName).Get([]byte(path))
		if value != nil {
			found = true
			valCopy = make([]byte, len(value))
			copy(valCopy, value)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	return &physical.Entry{
		Key:   path,
		Value: valCopy,
	}, nil
}

func (b *FSM) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"raft", "put"}, time.Now())

	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.l.RLock()
	defer b.l.RUnlock()

	// Start a write transaction.
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(entry.Key), entry.Value)
	})
}

func (b *FSM) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"raft", "list"}, time.Now())

	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.l.RLock()
	defer b.l.RUnlock()

	var keys []string

	err := b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(bucketName).Cursor()

		prefixBytes := []byte(prefix)
		for k, _ := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, _ = c.Next() {
			key := string(k)
			key = strings.TrimPrefix(key, prefix)
			if i := strings.Index(key, "/"); i == -1 {
				// Add objects only from the current 'folder'
				keys = append(keys, key)
			} else if i != -1 {
				// Add truncated 'folder' paths
				keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
			}
		}

		return nil
	})

	return keys, err
}

func (b *FSM) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.l.RLock()
	defer b.l.RUnlock()

	// TODO: should this be a Batch?
	// Start a write transaction.
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for _, txn := range txns {
			var err error
			switch txn.Operation {
			case physical.PutOperation:
				err = b.Put([]byte(txn.Entry.Key), txn.Entry.Value)
			case physical.DeleteOperation:
				err = b.Delete([]byte(txn.Entry.Key))
			default:
				return fmt.Errorf("%q is not a supported transaction operation", txn.Operation)
			}
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (b *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshotter{
		f: func(sink raft.SnapshotSink) error {
			protoWriter := protoio.NewDelimitedWriter(sink)

			b.l.RLock()
			defer b.l.RUnlock()

			err := b.db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(bucketName)

				c := b.Cursor()

				for k, v := c.First(); k != nil; k, v = c.Next() {
					err := protoWriter.WriteMsg(&pb.StorageEntry{
						Key:   string(k),
						Value: v,
					})
					if err != nil {
						return err
					}
				}

				return nil
			})
			if err != nil {
				sink.Cancel()
				return err
			}

			sink.Close()
			return nil
		},
	}, nil
}

func (f *FSM) Apply(log *raft.Log) interface{} {
	command := &LogData{}
	err := proto.Unmarshal(log.Data, command)
	if err != nil {
		panic("error proto unmarshaling log data")
	}

	f.l.RLock()
	defer f.l.RUnlock()

	err = f.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for _, op := range command.Operations {
			var err error
			switch op.OpType {
			case putOp:
				err = b.Put([]byte(op.Key), op.Value)
			case deleteOp:
				err = b.Delete([]byte(op.Key))
			default:
				return fmt.Errorf("%q is not a supported transaction operation", op.OpType)
			}
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic("failed to store data")
	}

	return &FSMApplyResponse{
		Success: true,
	}
}

func (f *FSM) Restore(r io.ReadCloser) error {
	protoReader := protoio.NewDelimitedReader(r, math.MaxInt64)
	defer protoReader.Close()

	f.l.Lock()
	defer f.l.Unlock()

	// Start a write transaction.
	err := f.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil {
			return err
		}

		b, err := tx.CreateBucket(bucketName)
		if err != nil {
			return err
		}

		for {
			s := new(pb.StorageEntry)
			err := protoReader.ReadMsg(s)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			err = b.Put([]byte(s.Key), s.Value)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

type snapshotter struct {
	f func(sink raft.SnapshotSink) error
}

func (s *snapshotter) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()
	return s.f(sink)
}

func (s *snapshotter) Release() {}
