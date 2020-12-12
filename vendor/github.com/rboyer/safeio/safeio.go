// Package safeio provides functions to perform atomic, fsync-safe disk
// operations.
package safeio

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// WriteToFile consumes the provided io.Reader and writes it to a temp
// file in the provided directory.
func WriteToFile(src io.Reader, path string, perm os.FileMode) (written int64, err error) {
	tempName, written, err := writeToTempFile(src, path, perm)

	if err == nil {
		err = Rename(tempName, path)
	}

	return written, err
}

// writeToTempFile consumes the provided io.Reader and writes it to a
// temp file in the same directory as path.
func writeToTempFile(src io.Reader, path string, perm os.FileMode) (tempName string, written int64, err error) {
	dir := filepath.Dir(path)
	name := filepath.Base(path)

	f, err := ioutil.TempFile(dir, name+".tmp")
	if err != nil {
		return "", 0, err
	}

	tempName = f.Name()

	cleanup := func(written int64, err error) (string, int64, error) {
		_ = f.Close()
		_ = os.Remove(tempName)
		return "", written, err
	}

	if err = f.Chmod(perm); err != nil {
		return cleanup(0, err)
	}

	written, err = io.Copy(f, src)
	if err != nil {
		return cleanup(written, err)
	}

	if err := f.Sync(); err != nil {
		return cleanup(written, err)
	}

	if err := f.Close(); err != nil {
		return cleanup(written, err)
	}

	return tempName, written, nil
}

// Remove is just like os.Remove, except this also calls sync on the
// parent directory.
func Remove(fn string) error {
	err := os.Remove(fn)
	if err != nil {
		return err
	}

	// fsync the dir
	return syncParentDir(fn)
}

// Rename renames the file using os.Rename and fsyncs the NEW parent
// directory. It should only be used if both oldname and newname are in
// the same directory.
func Rename(oldname, newname string) error {
	err := os.Rename(oldname, newname)
	if err != nil {
		return err
	}

	// fsync the dir
	return syncParentDir(newname)
}

func syncParentDir(name string) error {
	f, err := os.Open(filepath.Dir(name))
	if err != nil {
		return err
	}
	defer f.Close()

	return f.Sync()
}
