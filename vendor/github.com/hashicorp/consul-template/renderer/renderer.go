package renderer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	// DefaultFilePerms are the default file permissions for files rendered onto
	// disk when a specific file permission has not already been specified.
	DefaultFilePerms = 0644
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

// Render atomically renders a file contents to disk, returning a result of
// whether it would have rendered and actually did render.
func Render(i *RenderInput) (*RenderResult, error) {
	existing, err := ioutil.ReadFile(i.Path)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "failed reading file")
	}

	if bytes.Equal(existing, i.Contents) {
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
func AtomicWrite(path string, createDestDirs bool, contents []byte, perms os.FileMode, backup bool) error {
	if path == "" {
		return ErrMissingDest
	}

	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		if createDestDirs {
			if err := os.MkdirAll(parent, 0755); err != nil {
				return err
			}
		} else {
			return ErrNoParentDir
		}
	}

	f, err := ioutil.TempFile(parent, "")
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
	if perms == 0 {
		currentInfo, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				perms = DefaultFilePerms
			} else {
				return err
			}
		} else {
			perms = currentInfo.Mode()

			// The file exists, so try to preserve the ownership as well.
			if err := preserveFilePermissions(f.Name(), currentInfo); err != nil {
				log.Printf("[WARN] (runner) could not preserve file permissions for %q: %v",
					f.Name(), err)
			}
		}
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
