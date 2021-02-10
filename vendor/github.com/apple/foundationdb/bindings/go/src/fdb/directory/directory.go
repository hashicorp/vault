/*
 * directory.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2018 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// FoundationDB Go Directory Layer

// Package directory provides a tool for managing related subspaces. Directories
// are a recommended approach for administering applications. Each application
// should create or open at least one directory to manage its subspaces.
//
// For general guidance on directory usage, see the Directories section of the
// Developer Guide
// (https://apple.github.io/foundationdb/developer-guide.html#directories).
//
// Directories are identified by hierarchical paths analogous to the paths in a
// Unix-like file system. A path is represented as a slice of strings. Each
// directory has an associated subspace used to store its content. The directory
// layer maps each path to a short prefix used for the corresponding
// subspace. In effect, directories provide a level of indirection for access to
// subspaces.
//
// Directory operations are transactional. A byte slice layer option is used as
// a metadata identifier when opening a directory.
package directory

import (
	"errors"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
)

const (
	_SUBDIRS int = 0

	// []int32{1,0,0} by any other name
	_MAJORVERSION int32 = 1
	_MINORVERSION int32 = 0
	_MICROVERSION int32 = 0
)

var (
	// ErrDirAlreadyExists is returned when trying to create a directory while it already exists.
	ErrDirAlreadyExists = errors.New("the directory already exists")

	// ErrDirNotExists is returned when opening or listing a directory that does not exist.
	ErrDirNotExists = errors.New("the directory does not exist")

	// ErrParentDirDoesNotExist is returned when opening a directory and one or more
	// parent directories in the path do not exist.
	ErrParentDirDoesNotExist = errors.New("the parent directory does not exist")
)

// Directory represents a subspace of keys in a FoundationDB database,
// identified by a hierarchical path.
type Directory interface {
	// CreateOrOpen opens the directory specified by path (relative to this
	// Directory), and returns the directory and its contents as a
	// DirectorySubspace. If the directory does not exist, it is created
	// (creating parent directories if necessary).
	//
	// If the byte slice layer is specified and the directory is new, it is
	// recorded as the layer; if layer is specified and the directory already
	// exists, it is compared against the layer specified when the directory was
	// created, and an error is returned if they differ.
	CreateOrOpen(t fdb.Transactor, path []string, layer []byte) (DirectorySubspace, error)

	// Open opens the directory specified by path (relative to this Directory),
	// and returns the directory and its contents as a DirectorySubspace (or ErrDirNotExists
	// error if the directory does not exist, or ErrParentDirDoesNotExist if one of the parent
	// directories in the path does not exist).
	//
	// If the byte slice layer is specified, it is compared against the layer
	// specified when the directory was created, and an error is returned if
	// they differ.
	Open(rt fdb.ReadTransactor, path []string, layer []byte) (DirectorySubspace, error)

	// Create creates a directory specified by path (relative to this
	// Directory), and returns the directory and its contents as a
	// DirectorySubspace (or ErrDirAlreadyExists if the directory already exists).
	//
	// If the byte slice layer is specified, it is recorded as the layer and
	// will be checked when opening the directory in the future.
	Create(t fdb.Transactor, path []string, layer []byte) (DirectorySubspace, error)

	// CreatePrefix behaves like Create, but uses a manually specified byte
	// slice prefix to physically store the contents of this directory, rather
	// than an automatically allocated prefix.
	//
	// If this Directory was created in a root directory that does not allow
	// manual prefixes, CreatePrefix will return an error. The default root
	// directory does not allow manual prefixes.
	CreatePrefix(t fdb.Transactor, path []string, layer []byte, prefix []byte) (DirectorySubspace, error)

	// Move moves the directory at oldPath to newPath (both relative to this
	// Directory), and returns the directory (at its new location) and its
	// contents as a DirectorySubspace. Move will return an error if a directory
	// does not exist at oldPath, a directory already exists at newPath, or the
	// parent directory of newPath does not exist.
	//
	// There is no effect on the physical prefix of the given directory or on
	// clients that already have the directory open.
	Move(t fdb.Transactor, oldPath []string, newPath []string) (DirectorySubspace, error)

	// MoveTo moves this directory to newAbsolutePath (relative to the root
	// directory of this Directory), and returns the directory (at its new
	// location) and its contents as a DirectorySubspace. MoveTo will return an
	// error if a directory already exists at newAbsolutePath or the parent
	// directory of newAbsolutePath does not exist.
	//
	// There is no effect on the physical prefix of the given directory or on
	// clients that already have the directory open.
	MoveTo(t fdb.Transactor, newAbsolutePath []string) (DirectorySubspace, error)

	// Remove removes the directory at path (relative to this Directory), its
	// content, and all subdirectories. Remove returns true if a directory
	// existed at path and was removed, and false if no directory exists at
	// path.
	//
	// Note that clients that have already opened this directory might still
	// insert data into its contents after removal.
	Remove(t fdb.Transactor, path []string) (bool, error)

	// Exists returns true if the directory at path (relative to this Directory)
	// exists, and false otherwise.
	Exists(rt fdb.ReadTransactor, path []string) (bool, error)

	// List returns the names of the immediate subdirectories of the directory
	// at path (relative to this Directory) as a slice of strings. Each string
	// is the name of the last component of a subdirectory's path.
	List(rt fdb.ReadTransactor, path []string) ([]string, error)

	// GetLayer returns the layer specified when this Directory was created.
	GetLayer() []byte

	// GetPath returns the path with which this Directory was opened.
	GetPath() []string
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func moveTo(t fdb.Transactor, dl directoryLayer, path, newAbsolutePath []string) (DirectorySubspace, error) {
	partition_len := len(dl.path)

	if !stringsEqual(newAbsolutePath[:partition_len], dl.path) {
		return nil, errors.New("cannot move between partitions")
	}

	return dl.Move(t, path[partition_len:], newAbsolutePath[partition_len:])
}

var root = NewDirectoryLayer(subspace.FromBytes([]byte{0xFE}), subspace.AllKeys(), false)

// CreateOrOpen opens the directory specified by path (resolved relative to the
// default root directory), and returns the directory and its contents as a
// DirectorySubspace. If the directory does not exist, it is created (creating
// parent directories if necessary).
//
// If the byte slice layer is specified and the directory is new, it is recorded
// as the layer; if layer is specified and the directory already exists, it is
// compared against the layer specified when the directory was created, and an
// error is returned if they differ.
func CreateOrOpen(t fdb.Transactor, path []string, layer []byte) (DirectorySubspace, error) {
	return root.CreateOrOpen(t, path, layer)
}

// Open opens the directory specified by path (resolved relative to the default
// root directory), and returns the directory and its contents as a
// DirectorySubspace (or an error if the directory does not exist).
//
// If the byte slice layer is specified, it is compared against the layer
// specified when the directory was created, and an error is returned if they
// differ.
func Open(rt fdb.ReadTransactor, path []string, layer []byte) (DirectorySubspace, error) {
	return root.Open(rt, path, layer)
}

// Create creates a directory specified by path (resolved relative to the
// default root directory), and returns the directory and its contents as a
// DirectorySubspace (or an error if the directory already exists).
//
// If the byte slice layer is specified, it is recorded as the layer and will be
// checked when opening the directory in the future.
func Create(t fdb.Transactor, path []string, layer []byte) (DirectorySubspace, error) {
	return root.Create(t, path, layer)
}

// Move moves the directory at oldPath to newPath (both resolved relative to the
// default root directory), and returns the directory (at its new location) and
// its contents as a DirectorySubspace. Move will return an error if a directory
// does not exist at oldPath, a directory already exists at newPath, or the
// parent directory of newPath does not exit.
//
// There is no effect on the physical prefix of the given directory or on
// clients that already have the directory open.
func Move(t fdb.Transactor, oldPath []string, newPath []string) (DirectorySubspace, error) {
	return root.Move(t, oldPath, newPath)
}

// Exists returns true if the directory at path (relative to the default root
// directory) exists, and false otherwise.
func Exists(rt fdb.ReadTransactor, path []string) (bool, error) {
	return root.Exists(rt, path)
}

// List returns the names of the immediate subdirectories of the default root
// directory as a slice of strings. Each string is the name of the last
// component of a subdirectory's path.
func List(rt fdb.ReadTransactor, path []string) ([]string, error) {
	return root.List(rt, path)
}

// Root returns the default root directory. Any attempt to move or remove the
// root directory will return an error.
//
// The default root directory stores directory layer metadata in keys beginning
// with 0xFE, and allocates newly created directories in (unused) prefixes
// starting with 0x00 through 0xFD. This is appropriate for otherwise empty
// databases, but may conflict with other formal or informal partitionings of
// keyspace. If you already have other content in your database, you may wish to
// use NewDirectoryLayer to construct a non-standard root directory to control
// where metadata and keys are stored.
//
// As an alternative to Root, you may use the package-level functions
// CreateOrOpen, Open, Create, CreatePrefix, Move, Exists and List to operate
// directly on the default DirectoryLayer.
func Root() Directory {
	return root
}
