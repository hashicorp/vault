/*
 * directoryPartition.go
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
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

type directoryPartition struct {
	directoryLayer
	parentDirectoryLayer directoryLayer
}

func (dp directoryPartition) Sub(el ...tuple.TupleElement) subspace.Subspace {
	panic("cannot open subspace in the root of a directory partition")
}

func (dp directoryPartition) Bytes() []byte {
	panic("cannot get key for the root of a directory partition")
}

func (dp directoryPartition) Pack(t tuple.Tuple) fdb.Key {
	panic("cannot pack keys using the root of a directory partition")
}

func (dp directoryPartition) Unpack(k fdb.KeyConvertible) (tuple.Tuple, error) {
	panic("cannot unpack keys using the root of a directory partition")
}

func (dp directoryPartition) Contains(k fdb.KeyConvertible) bool {
	panic("cannot check whether a key belongs to the root of a directory partition")
}

func (dp directoryPartition) FDBKey() fdb.Key {
	panic("cannot use the root of a directory partition as a key")
}

func (dp directoryPartition) FDBRangeKeys() (fdb.KeyConvertible, fdb.KeyConvertible) {
	panic("cannot get range for the root of a directory partition")
}

func (dp directoryPartition) FDBRangeKeySelectors() (fdb.Selectable, fdb.Selectable) {
	panic("cannot get range for the root of a directory partition")
}

func (dp directoryPartition) GetLayer() []byte {
	return []byte("partition")
}

func (dp directoryPartition) getLayerForPath(path []string) directoryLayer {
	if len(path) == 0 {
		return dp.parentDirectoryLayer
	}
	return dp.directoryLayer
}

func (dp directoryPartition) MoveTo(t fdb.Transactor, newAbsolutePath []string) (DirectorySubspace, error) {
	return moveTo(t, dp.parentDirectoryLayer, dp.path, newAbsolutePath)
}

func (dp directoryPartition) Remove(t fdb.Transactor, path []string) (bool, error) {
	dl := dp.getLayerForPath(path)
	return dl.Remove(t, dl.partitionSubpath(dp.path, path))
}

func (dp directoryPartition) Exists(rt fdb.ReadTransactor, path []string) (bool, error) {
	dl := dp.getLayerForPath(path)
	return dl.Exists(rt, dl.partitionSubpath(dp.path, path))
}
