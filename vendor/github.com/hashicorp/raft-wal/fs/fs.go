// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package fs

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/hashicorp/raft-wal/types"
)

// FS implements the wal.VFS interface using GO's built in OS Filesystem (and a
// few helpers).
//
// TODO if we changed the interface to be Dir centric we could cache the open
// dir handle and save some time opening it on each Create in order to fsync.
type FS struct {
}

func New() *FS {
	return &FS{}
}

// ListDir returns a list of all files in the specified dir in lexicographical
// order. If the dir doesn't exist, it must return an error. Empty array with
// nil error is assumed to mean that the directory exists and was readable,
// but contains no files.
func (fs *FS) ListDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(files))
	for i, f := range files {
		if f.IsDir() {
			continue
		}
		names[i] = f.Name()
	}
	return names, nil
}

// Create creates a new file with the given name. If a file with the same name
// already exists an error is returned. If a non-zero size is given,
// implementations should make a best effort to pre-allocate the file to be
// that size. The dir must already exist and be writable to the current
// process.
func (fs *FS) Create(dir string, name string, size uint64) (types.WritableFile, error) {
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_EXCL|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return nil, err
	}
	// We just created the file. Preallocate it's size.
	if size > 0 {
		if size > math.MaxInt32 {
			return nil, fmt.Errorf("maximum file size is %d bytes", math.MaxInt32)
		}
		if err := fileutil.Preallocate(f, int64(size), true); err != nil {
			f.Close()
			return nil, err
		}
	}
	// We don't fsync here for performance reasons. Technically we need to fsync
	// the file itself to make sure it is really persisted to disk, and you always
	// need to fsync its parent dir after a creation because fsync doesn't ensure
	// the directory entry is persisted - a crash could make the file appear to be
	// missing as there is no directory entry.
	//
	// BUT, it doesn't actually matter if this file is crash safe, right up to the
	// point where we actually commit log data. Since we always fsync the file
	// when we commit logs, we don't need to again here. That does however leave
	// the parent dir fsync which must be done after the first fsync to a newly
	// created file to ensure it survives a crash.
	//
	// To handle that, we return a wrapped io.File that will fsync the parent dir
	// as well the first time Sync is called (and only the first time),
	fi := &File{
		new:  0,
		dir:  dir,
		File: *f,
	}
	return fi, nil
}

// Delete indicates the file is no longer required. Typically it should be
// deleted from the underlying system to free disk space.
func (fs *FS) Delete(dir string, name string) error {
	if err := os.Remove(filepath.Join(dir, name)); err != nil {
		return err
	}
	// Make sure parent directory metadata is fsynced too before we call this
	// "done".
	return syncDir(dir)
}

// OpenReader opens an existing file in read-only mode. If the file doesn't
// exist or permission is denied, an error is returned, otherwise no checks
// are made about the well-formedness of the file, it may be empty, the wrong
// size or corrupt in arbitrary ways.
func (fs *FS) OpenReader(dir string, name string) (types.ReadableFile, error) {
	return os.OpenFile(filepath.Join(dir, name), os.O_RDONLY, os.FileMode(0644))
}

// OpenWriter opens a file in read-write mode. If the file doesn't exist or
// permission is denied, an error is returned, otherwise no checks are made
// about the well-formedness of the file, it may be empty, the wrong size or
// corrupt in arbitrary ways.
func (fs *FS) OpenWriter(dir string, name string) (types.WritableFile, error) {
	return os.OpenFile(filepath.Join(dir, name), os.O_RDWR, os.FileMode(0644))
}

func syncDir(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	err = f.Sync()
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}
