// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package types

import "io"

// VFS is the interface WAL needs to interact with the file system. In
// production it would normally be implemented by RealFS which interacts with
// the operating system FS using standard go os package. It's useful to allow
// testing both to run quicker (by being in memory only) and to make it easy to
// simulate all kinds of disk errors and failure modes without needing a more
// elaborate external test harness like ALICE.
type VFS interface {
	// ListDir returns a list of all files in the specified dir in lexicographical
	// order. If the dir doesn't exist, it must return an error. Empty array with
	// nil error is assumed to mean that the directory exists and was readable,
	// but contains no files.
	ListDir(dir string) ([]string, error)

	// Create creates a new file with the given name. If a file with the same name
	// already exists an error is returned. If a non-zero size is given,
	// implementations should make a best effort to pre-allocate the file to be
	// that size. The dir must already exist and be writable to the current
	// process.
	Create(dir, name string, size uint64) (WritableFile, error)

	// Delete indicates the file is no longer required. Typically it should be
	// deleted from the underlying system to free disk space.
	Delete(dir, name string) error

	// OpenReader opens an existing file in read-only mode. If the file doesn't
	// exist or permission is denied, an error is returned, otherwise no checks
	// are made about the well-formedness of the file, it may be empty, the wrong
	// size or corrupt in arbitrary ways.
	OpenReader(dir, name string) (ReadableFile, error)

	// OpenWriter opens a file in read-write mode. If the file doesn't exist or
	// permission is denied, an error is returned, otherwise no checks are made
	// about the well-formedness of the file, it may be empty, the wrong size or
	// corrupt in arbitrary ways.
	OpenWriter(dir, name string) (WritableFile, error)
}

// WritableFile provides random read-write access to a file as well as the
// ability to fsync it to disk.
type WritableFile interface {
	io.WriterAt
	io.ReaderAt
	io.Closer

	Sync() error
}

// ReadableFile provides random read access to a file.
type ReadableFile interface {
	io.ReaderAt
	io.Closer
}
