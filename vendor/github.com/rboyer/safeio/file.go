package safeio

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var errClosed = errors.New("file is already closed")

// OpenFile is the incremental version of WriteToFile.  It opens a temp
// file and proxies writes through to the underlying file.
//
// If Close is called before Commit, the temp file is closed and erased.
//
// If Commit is called before Close, the temp file is closed, fsynced,
// and atomically renamed to the desired final name.
func OpenFile(path string, perm os.FileMode) (*File, error) {
	dir := filepath.Dir(path)
	name := filepath.Base(path)

	f, err := ioutil.TempFile(dir, name+".tmp")
	if err != nil {
		return nil, err
	}

	return &File{
		name:     path,
		tempName: f.Name(),
		perm:     perm,
		file:     f,
	}, nil
}

// File is an implementation detail of OpenFile.
type File struct {
	name     string // track desired filename
	tempName string // track actual filename
	perm     os.FileMode
	file     *os.File
	closed   bool
	err      error // the first error encountered
}

// Write is a thin proxy to *os.File#Write.
//
// If Close or Commit were called, this immediately exits with an error.
func (f *File) Write(p []byte) (n int, err error) {
	if f.closed {
		return 0, errClosed
	} else if f.err != nil {
		return 0, f.err
	}

	n, err = f.file.Write(p)
	if err != nil {
		f.err = err
	}

	return n, err
}

// Commit causes the current temp file to be safely persisted to disk and atomically renamed to the desired final filename.
//
// It is safe to call Close after commit, so you can defer Close as
// usual without worries about write-safey.
func (f *File) Commit() error {
	if f.closed {
		return errClosed
	} else if f.err != nil {
		return f.err
	}

	if err := f.file.Sync(); err != nil {
		return f.cleanup(err)
	}

	if err := f.file.Chmod(f.perm); err != nil {
		return f.cleanup(err)
	}

	if err := f.file.Close(); err != nil {
		return f.cleanup(err)
	}

	if err := Rename(f.tempName, f.name); err != nil {
		return f.cleanup(err)
	}

	f.closed = true

	return nil
}

// Close closes the current file and erases it, unless Commit was
// previously called.  In that case it does nothing.
//
// Close is idempotent.
//
// After Close is called, Write and Commit will fail.
func (f *File) Close() error {
	if !f.closed {
		_ = f.cleanup(nil)
		f.closed = true
	}
	return f.err
}

func (f *File) cleanup(err error) error {
	_ = f.file.Close()
	_ = os.Remove(f.tempName)

	if f.err == nil {
		f.err = err
	}
	return f.err
}

// setErr is only used during testing to simulate os.File errors
func (f *File) setErr(err error) {
	f.err = err
}
