/*
 * directorySubspace.go
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

package directory

import (
	"fmt"
	"strings"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
)

// DirectorySubspace represents a Directory that may also be used as a Subspace
// to store key/value pairs. Subdirectories of a root directory (as returned by
// Root or NewDirectoryLayer) are DirectorySubspaces, and provide all methods of
// the Directory and subspace.Subspace interfaces.
type DirectorySubspace interface {
	subspace.Subspace
	Directory
}

type directorySubspace struct {
	subspace.Subspace
	dl    directoryLayer
	path  []string
	layer []byte
}

// String implements the fmt.Stringer interface and returns human-readable
// string representation of this object.
func (ds directorySubspace) String() string {
	var path string
	if len(ds.path) > 0 {
		path = "(" + strings.Join(ds.path, ",") + ")"
	} else {
		path = "nil"
	}
	return fmt.Sprintf("DirectorySubspace(%s, %s)", path, fdb.Printable(ds.Bytes()))
}

func (d directorySubspace) CreateOrOpen(t fdb.Transactor, path []string, layer []byte) (DirectorySubspace, error) {
	return d.dl.CreateOrOpen(t, d.dl.partitionSubpath(d.path, path), layer)
}

func (d directorySubspace) Create(t fdb.Transactor, path []string, layer []byte) (DirectorySubspace, error) {
	return d.dl.Create(t, d.dl.partitionSubpath(d.path, path), layer)
}

func (d directorySubspace) CreatePrefix(t fdb.Transactor, path []string, layer []byte, prefix []byte) (DirectorySubspace, error) {
	return d.dl.CreatePrefix(t, d.dl.partitionSubpath(d.path, path), layer, prefix)
}

func (d directorySubspace) Open(rt fdb.ReadTransactor, path []string, layer []byte) (DirectorySubspace, error) {
	return d.dl.Open(rt, d.dl.partitionSubpath(d.path, path), layer)
}

func (d directorySubspace) MoveTo(t fdb.Transactor, newAbsolutePath []string) (DirectorySubspace, error) {
	return moveTo(t, d.dl, d.path, newAbsolutePath)
}

func (d directorySubspace) Move(t fdb.Transactor, oldPath []string, newPath []string) (DirectorySubspace, error) {
	return d.dl.Move(t, d.dl.partitionSubpath(d.path, oldPath), d.dl.partitionSubpath(d.path, newPath))
}

func (d directorySubspace) Remove(t fdb.Transactor, path []string) (bool, error) {
	return d.dl.Remove(t, d.dl.partitionSubpath(d.path, path))
}

func (d directorySubspace) Exists(rt fdb.ReadTransactor, path []string) (bool, error) {
	return d.dl.Exists(rt, d.dl.partitionSubpath(d.path, path))
}

func (d directorySubspace) List(rt fdb.ReadTransactor, path []string) (subdirs []string, e error) {
	return d.dl.List(rt, d.dl.partitionSubpath(d.path, path))
}

func (d directorySubspace) GetLayer() []byte {
	return d.layer
}

func (d directorySubspace) GetPath() []string {
	return d.path
}
