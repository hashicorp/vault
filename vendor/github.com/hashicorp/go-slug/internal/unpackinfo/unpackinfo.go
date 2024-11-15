// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package unpackinfo

import (
	"archive/tar"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UnpackInfo stores information about the file (or directory, or symlink) being
// unpacked. UnpackInfo ensures certain malicious tar files are not unpacked.
// The information can be used later to restore the original permissions
// and timestamps based on the type of entry the info represents.
type UnpackInfo struct {
	Path               string
	OriginalAccessTime time.Time
	OriginalModTime    time.Time
	OriginalMode       fs.FileMode
	Typeflag           byte
}

// NewUnpackInfo returns an UnpackInfo based on a destination root and a tar header.
// It will return an error if the header represents an illegal symlink extraction
// or if the entry type is not supported by go-slug.
func NewUnpackInfo(dst string, header *tar.Header) (UnpackInfo, error) {
	// Get rid of absolute paths.
	path := header.Name

	if path[0] == '/' {
		path = path[1:]
	}
	path = filepath.Join(dst, path)

	// Check for paths outside our directory, they are forbidden
	target := filepath.Clean(path)
	if !strings.HasPrefix(target, dst) {
		return UnpackInfo{}, errors.New("invalid filename, traversal with \"..\" outside of current directory")
	}

	// Ensure the destination is not through any symlinks. This prevents
	// any files from being deployed through symlinks defined in the slug.
	// There are malicious cases where this could be used to escape the
	// slug's boundaries (zipslip), and any legitimate use is questionable
	// and likely indicates a hand-crafted tar file, which we are not in
	// the business of supporting here.
	//
	// The strategy is to Lstat each path  component from dst up to the
	// immediate parent directory of the file name in the tarball, checking
	// the mode on each to ensure we wouldn't be passing through any
	// symlinks.
	currentPath := dst // Start at the root of the unpacked tarball.
	components := strings.Split(header.Name, "/")

	for i := 0; i < len(components)-1; i++ {
		currentPath = filepath.Join(currentPath, components[i])
		fi, err := os.Lstat(currentPath)
		if os.IsNotExist(err) {
			// Parent directory structure is incomplete. Technically this
			// means from here upward cannot be a symlink, so we cancel the
			// remaining path tests.
			break
		}
		if err != nil {
			return UnpackInfo{}, fmt.Errorf("failed to evaluate path %q: %w", header.Name, err)
		}
		if fi.Mode()&fs.ModeSymlink != 0 {
			return UnpackInfo{}, fmt.Errorf("cannot extract %q through symlink", header.Name)
		}
	}

	result := UnpackInfo{
		Path:               path,
		OriginalAccessTime: header.AccessTime,
		OriginalModTime:    header.ModTime,
		OriginalMode:       header.FileInfo().Mode(),
		Typeflag:           header.Typeflag,
	}

	if !result.IsDirectory() && !result.IsSymlink() && !result.IsRegular() && !result.IsTypeX() {
		return UnpackInfo{}, fmt.Errorf("failed creating %q, unsupported file type %c", path, result.Typeflag)
	}

	return result, nil
}

// IsSymlink describes whether the file being unpacked is a symlink
func (i UnpackInfo) IsSymlink() bool {
	return i.Typeflag == tar.TypeSymlink
}

// IsDirectory describes whether the file being unpacked is a directory
func (i UnpackInfo) IsDirectory() bool {
	return i.Typeflag == tar.TypeDir
}

// IsTypeX describes whether the file being unpacked is a special TypeXHeader that can
// be ignored by go-slug
func (i UnpackInfo) IsTypeX() bool {
	return i.Typeflag == tar.TypeXGlobalHeader || i.Typeflag == tar.TypeXHeader
}

// IsRegular describes whether the file being unpacked is a regular file
func (i UnpackInfo) IsRegular() bool {
	return i.Typeflag == tar.TypeReg || i.Typeflag == tar.TypeRegA
}

// RestoreInfo changes the file mode and timestamps for the given UnpackInfo data
func (i UnpackInfo) RestoreInfo() error {
	switch {
	case i.IsDirectory():
		return i.restoreDirectory()
	case i.IsSymlink():
		if CanMaintainSymlinkTimestamps() {
			return i.restoreSymlink()
		}
		return nil
	default: // Normal file
		return i.restoreNormal()
	}
}

func (i UnpackInfo) restoreDirectory() error {
	if err := os.Chmod(i.Path, i.OriginalMode); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed setting permissions on directory %q: %w", i.Path, err)
	}

	if err := os.Chtimes(i.Path, i.OriginalAccessTime, i.OriginalModTime); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed setting times on directory %q: %w", i.Path, err)
	}
	return nil
}

func (i UnpackInfo) restoreSymlink() error {
	if err := i.Lchtimes(); err != nil {
		return fmt.Errorf("failed setting times on symlink %q: %w", i.Path, err)
	}
	return nil
}

func (i UnpackInfo) restoreNormal() error {
	if err := os.Chmod(i.Path, i.OriginalMode); err != nil {
		return fmt.Errorf("failed setting permissions on %q: %w", i.Path, err)
	}

	if err := os.Chtimes(i.Path, i.OriginalAccessTime, i.OriginalModTime); err != nil {
		return fmt.Errorf("failed setting times on %q: %w", i.Path, err)
	}
	return nil
}
