// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/physical"
)

// Verify FileBackend satisfies the correct interfaces
var (
	_ physical.Backend             = (*FileBackend)(nil)
	_ physical.Transactional       = (*TransactionalFileBackend)(nil)
	_ physical.PseudoTransactional = (*FileBackend)(nil)
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
	permitPool *permitpool.Pool
}

type TransactionalFileBackend struct {
	FileBackend
}

type fileEntry struct {
	Value []byte
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
		permitPool: permitpool.New(physical.DefaultParallelOperations),
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
			permitPool: permitpool.New(1),
		},
	}, nil
}

func (b *FileBackend) Delete(ctx context.Context, path string) error {
	if err := b.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer b.permitPool.Release()

	b.Lock()
	defer b.Unlock()

	return b.DeleteInternal(ctx, path)
}

func (b *FileBackend) DeleteInternal(ctx context.Context, path string) error {
	if path == "" {
		return nil
	}

	if err := b.validatePath(path); err != nil {
		return err
	}

	basePath, key := b.expandPath(path)
	fullPath := filepath.Join(basePath, key)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return errwrap.Wrapf(fmt.Sprintf("failed to remove %q: {{err}}", fullPath), err)
	}

	err = b.cleanupLogicalPath(path)

	return err
}

// cleanupLogicalPath is used to remove all empty nodes, beginning with deepest
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

func (b *FileBackend) Get(ctx context.Context, k string) (*physical.Entry, error) {
	if err := b.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer b.permitPool.Release()

	b.RLock()
	defer b.RUnlock()

	return b.GetInternal(ctx, k)
}

func (b *FileBackend) GetInternal(ctx context.Context, k string) (*physical.Entry, error) {
	if err := b.validatePath(k); err != nil {
		return nil, err
	}

	path, key := b.expandPath(k)
	path = filepath.Join(path, key)

	// If we stat it and it exists but is size zero, it may be left from some
	// previous FS error like out-of-space. No Vault entry will ever be zero
	// length, so simply remove it and return nil.
	fi, err := os.Stat(path)
	if err == nil {
		if fi.Size() == 0 {
			// Best effort, ignore errors
			os.Remove(path)
			return nil, nil
		}
	}

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

	var entry fileEntry
	if err := jsonutil.DecodeJSONFromReader(f, &entry); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return &physical.Entry{
		Key:   k,
		Value: entry.Value,
	}, nil
}

func (b *FileBackend) Put(ctx context.Context, entry *physical.Entry) error {
	if err := b.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer b.permitPool.Release()

	b.Lock()
	defer b.Unlock()

	return b.PutInternal(ctx, entry)
}

func (b *FileBackend) PutInternal(ctx context.Context, entry *physical.Entry) error {
	if err := b.validatePath(entry.Key); err != nil {
		return err
	}
	path, key := b.expandPath(entry.Key)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Make the parent tree
	if err := os.MkdirAll(path, 0o700); err != nil {
		return err
	}

	// JSON encode the entry and write it
	fullPath := filepath.Join(path, key)
	f, err := os.CreateTemp(path, key)
	if err != nil {
		if f != nil {
			f.Close()
		}
		return err
	}

	if err = os.Chmod(f.Name(), 0o600); err != nil {
		if f != nil {
			f.Close()
		}
		return err
	}

	if f == nil {
		return errors.New("could not successfully get a file handle")
	}

	enc := json.NewEncoder(f)
	encErr := enc.Encode(&fileEntry{
		Value: entry.Value,
	})
	f.Close()
	if encErr == nil {
		err = os.Rename(f.Name(), fullPath)
		if err != nil {
			return err
		}
		return nil
	}

	// Everything below is best-effort and will result in encErr being returned

	// See if we ended up with a zero-byte file and if so delete it, might be a
	// case of disk being full but the file info is in metadata that is
	// reserved.
	fi, err := os.Stat(f.Name())
	if err != nil {
		return encErr
	}
	if fi == nil {
		return encErr
	}
	if fi.Size() == 0 {
		os.Remove(f.Name())
	}
	return encErr
}

func (b *FileBackend) List(ctx context.Context, prefix string) ([]string, error) {
	if err := b.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer b.permitPool.Release()

	b.RLock()
	defer b.RUnlock()

	return b.ListInternal(ctx, prefix)
}

func (b *FileBackend) ListInternal(ctx context.Context, prefix string) ([]string, error) {
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

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(names) > 0 {
		sort.Strings(names)
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

func (b *TransactionalFileBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	if err := b.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer b.permitPool.Release()

	b.Lock()
	defer b.Unlock()

	return physical.GenericTransactionHandler(ctx, b, txns)
}
