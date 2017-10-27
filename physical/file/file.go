package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/physical"
)

// FileBackend is a physical backend that stores data on disk
// at a given file path. It can be used for durable single server
// situations, or to develop locally where durability is not critical.
//
// WARNING: the file backend implementation is currently extremely unsafe
// and non-performant. It is meant mostly for local testing and development.
// It can be improved in the future.
type FileBackend struct {
	sync.RWMutex
	path       string
	logger     log.Logger
	permitPool *physical.PermitPool
}

type TransactionalFileBackend struct {
	FileBackend
}

// NewFileBackend constructs a FileBackend using the given directory
func NewFileBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	return &FileBackend{
		path:       path,
		logger:     logger,
		permitPool: physical.NewPermitPool(physical.DefaultParallelOperations),
	}, nil
}

func NewTransactionalFileBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	// Create a pool of size 1 so only one operation runs at a time
	return &TransactionalFileBackend{
		FileBackend: FileBackend{
			path:       path,
			logger:     logger,
			permitPool: physical.NewPermitPool(1),
		},
	}, nil
}

func (b *FileBackend) Delete(path string) error {
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.Lock()
	defer b.Unlock()

	return b.DeleteInternal(path)
}

func (b *FileBackend) DeleteInternal(path string) error {
	if path == "" {
		return nil
	}

	if err := b.validatePath(path); err != nil {
		return err
	}

	basePath, key := b.expandPath(path)
	fullPath := filepath.Join(basePath, key)

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to remove %q: %v", fullPath, err)
	}

	err = b.cleanupLogicalPath(path)

	return err
}

// cleanupLogicalPath is used to remove all empty nodes, begining with deepest
// one, aborting on first non-empty one, up to top-level node.
func (b *FileBackend) cleanupLogicalPath(path string) error {
	nodes := strings.Split(path, fmt.Sprintf("%c", os.PathSeparator))
	for i := len(nodes) - 1; i > 0; i-- {
		fullPath := filepath.Join(b.path, filepath.Join(nodes[:i]...))

		dir, err := os.Open(fullPath)
		if err != nil {
			if dir != nil {
				dir.Close()
			}
			if os.IsNotExist(err) {
				return nil
			} else {
				return err
			}
		}

		list, err := dir.Readdir(1)
		dir.Close()
		if err != nil && err != io.EOF {
			return err
		}

		// If we have no entries, it's an empty directory; remove it
		if err == io.EOF || list == nil || len(list) == 0 {
			err = os.Remove(fullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *FileBackend) Get(k string) (*physical.Entry, error) {
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.RLock()
	defer b.RUnlock()

	return b.GetInternal(k)
}

func (b *FileBackend) GetInternal(k string) (*physical.Entry, error) {
	if err := b.validatePath(k); err != nil {
		return nil, err
	}

	path, key := b.expandPath(k)
	path = filepath.Join(path, key)

	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var entry physical.Entry
	if err := jsonutil.DecodeJSONFromReader(f, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

func (b *FileBackend) Put(entry *physical.Entry) error {
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.Lock()
	defer b.Unlock()

	return b.PutInternal(entry)
}

func (b *FileBackend) PutInternal(entry *physical.Entry) error {
	if err := b.validatePath(entry.Key); err != nil {
		return err
	}

	path, key := b.expandPath(entry.Key)

	// Make the parent tree
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	// JSON encode the entry and write it
	f, err := os.OpenFile(
		filepath.Join(path, key),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	return enc.Encode(entry)
}

func (b *FileBackend) List(prefix string) ([]string, error) {
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.RLock()
	defer b.RUnlock()

	return b.ListInternal(prefix)
}

func (b *FileBackend) ListInternal(prefix string) ([]string, error) {
	if err := b.validatePath(prefix); err != nil {
		return nil, err
	}

	path := b.path
	if prefix != "" {
		path = filepath.Join(path, prefix)
	}

	// Read the directory contents
	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	for i, name := range names {
		fi, err := os.Stat(filepath.Join(path, name))
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			names[i] = name + "/"
		} else {
			if name[0] == '_' {
				names[i] = name[1:]
			}
		}
	}

	return names, nil
}

func (b *FileBackend) expandPath(k string) (string, string) {
	path := filepath.Join(b.path, k)
	key := filepath.Base(path)
	path = filepath.Dir(path)
	return path, "_" + key
}

func (b *FileBackend) validatePath(path string) error {
	switch {
	case strings.Contains(path, ".."):
		return consts.ErrPathContainsParentReferences
	}

	return nil
}

func (b *TransactionalFileBackend) Transaction(txns []*physical.TxnEntry) error {
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.Lock()
	defer b.Unlock()

	return physical.GenericTransactionHandler(b, txns)
}
