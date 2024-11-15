// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package renderer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	// DefaultFilePerms are the default file permissions for files rendered onto
	// disk when a specific file permission has not already been specified.
	DefaultFilePerms = 0o644
)

var (
	// ErrNoParentDir is the error returned with the parent directory is missing
	// and the user disabled it.
	ErrNoParentDir = errors.New("parent directory is missing")

	// ErrMissingDest is the error returned with the destination is empty.
	ErrMissingDest = errors.New("missing destination")
)

// RenderInput is used as input to the render function.
type RenderInput struct {
	Backup         bool
	Contents       []byte
	CreateDestDirs bool
	Dry            bool
	DryStream      io.Writer
	Path           string
	Perms          os.FileMode
	User, Group    string
}

// RenderResult is returned and stored. It contains the status of the render
// operation.
type RenderResult struct {
	// DidRender indicates if the template rendered to disk. This will be false in
	// the event of an error, but it will also be false in dry mode or when the
	// template on disk matches the new result.
	DidRender bool

	// WouldRender indicates if the template would have rendered to disk. This
	// will return false in the event of an error, but will return true in dry
	// mode or when the template on disk matches the new result.
	WouldRender bool

	// Contents are the actual contents of the resulting template from the render
	// operation.
	Contents []byte
}

type Renderer func(*RenderInput) (*RenderResult, error)

// Render atomically renders a file contents to disk, returning a result of
// whether it would have rendered and actually did render.
func Render(i *RenderInput) (*RenderResult, error) {
	existing, err := os.ReadFile(i.Path)
	fileExists := !os.IsNotExist(err)
	if err != nil && fileExists {
		return nil, errors.Wrap(err, "failed reading file")
	}

	uid, err := lookupUser(i.User)
	if err != nil {
		return nil, errors.Wrap(err, "failed looking up user")
	}
	gid, err := lookupGroup(i.Group)
	if err != nil {
		return nil, errors.Wrap(err, "failed looking up group")
	}

	var chownNeeded bool

	if fileExists {
		chownNeeded, err = isChownNeeded(i.Path, uid, gid)
		if err != nil {
			log.Printf("[WARN] (runner) could not determine existing output file's permissions")
			chownNeeded = true
		}
	}

	if bytes.Equal(existing, i.Contents) && fileExists && !chownNeeded {
		return &RenderResult{
			DidRender:   false,
			WouldRender: true,
			Contents:    existing,
		}, nil
	}

	if i.Dry {
		fmt.Fprintf(i.DryStream, "> %s\n%s", i.Path, i.Contents)
	} else {
		if err := AtomicWrite(i.Path, i.CreateDestDirs, i.Contents, i.Perms, i.Backup); err != nil {
			return nil, errors.Wrap(err, "failed writing file")
		}

		if err = setFileOwnership(i.Path, uid, gid); err != nil {
			return nil, errors.Wrap(err, "failed setting file ownership")
		}
	}

	return &RenderResult{
		DidRender:   true,
		WouldRender: true,
		Contents:    i.Contents,
	}, nil
}

// AtomicWrite accepts a destination path and the template contents. It writes
// the template contents to a TempFile on disk, returning if any errors occur.
//
// If the parent destination directory does not exist, it will be created
// automatically with permissions 0755. To use a different permission, create
// the directory first or use `chmod` in a Command.
//
// If the destination path exists, all attempts will be made to preserve the
// existing file permissions. If those permissions cannot be read, an error is
// returned. If the file does not exist, it will be created automatically with
// permissions 0644. To use a different permission, create the destination file
// first or use `chmod` in a Command.
//
// If no errors occur, the Tempfile is "renamed" (moved) to the destination
// path.
//
// Please note that this is only atomic on POSIX systems. It is not atomic on
// Windows and it is impossible to rename atomically on Windows. For more on
// this see: https://github.com/golang/go/issues/22397#issuecomment-498856679
func AtomicWrite(path string, createDestDirs bool, contents []byte, perms os.FileMode, backup bool) error {
	if path == "" {
		return ErrMissingDest
	}

	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		if createDestDirs {
			if err := os.MkdirAll(parent, 0o755); err != nil {
				return err
			}
		} else {
			return ErrNoParentDir
		}
	}

	f, err := os.CreateTemp(parent, "")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(contents); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	// If the user did not explicitly set permissions, attempt to lookup the
	// current permissions on the file. If the file does not exist, fall back to
	// the default. Otherwise, inherit the current permissions.
	var existingPerms os.FileMode = DefaultFilePerms
	currentInfo, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		existingPerms = currentInfo.Mode()

		// The file exists, so try to preserve the ownership as well.
		if err := preserveFilePermissions(f.Name(), currentInfo); err != nil {
			log.Printf("[WARN] (runner) could not preserve file permissions for %q: %v",
				f.Name(), err)
		}
	}

	if perms == 0 {
		perms = existingPerms
	}

	if err := os.Chmod(f.Name(), perms); err != nil {
		return err
	}

	// If we got this far, it means we are about to save the file. Copy the
	// current file so we have a backup. Note that os.Link preserves the Mode.
	if backup {
		bak, old := path+".bak", path+".old.bak"
		os.Rename(bak, old) // ignore error
		if err := os.Link(path, bak); err != nil {
			log.Printf("[WARN] (runner) could not backup %q: %v", path, err)
		} else {
			os.Remove(old) // ignore error
		}
	}

	if err := os.Rename(f.Name(), path); err != nil {
		return err
	}

	return nil
}

// intPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}
