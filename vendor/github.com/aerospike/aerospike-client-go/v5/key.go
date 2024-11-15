// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"bytes"
	"fmt"

	"github.com/aerospike/aerospike-client-go/v5/types"
	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

// Key is the unique record identifier. Records can be identified using a specified namespace,
// an optional set name, and a user defined key which must be unique within a set.
// Records can also be identified by namespace/digest which is the combination used
// on the server.
type Key struct {
	// namespace. Equivalent to database name.
	namespace string

	// Optional set name. Equivalent to database table.
	setName string

	// Unique server hash value generated from set name and user key.
	digest [20]byte

	// Original user key. This key is immediately converted to a hash digest.
	// This key is not used or returned by the server by default. If the user key needs
	// to persist on the server, use one of the following methods:
	//
	// Set "WritePolicy.sendKey" to true. In this case, the key will be sent to the server for storage on writes
	// and retrieved on multi-record scans and queries.
	// Explicitly store and retrieve the key in a bin.
	userKey Value

	keyWriter keyWriter
}

// Namespace returns key's namespace.
func (ky *Key) Namespace() string {
	return ky.namespace
}

// SetName returns key's set name.
func (ky *Key) SetName() string {
	return ky.setName
}

// Value returns key's value.
func (ky *Key) Value() Value {
	return ky.userKey
}

// SetValue sets the Key's value and recompute's its digest without allocating new memory.
// This allows the keys to be reusable.
func (ky *Key) SetValue(val Value) Error {
	ky.userKey = val
	return ky.computeDigest()
}

// Digest returns key digest.
func (ky *Key) Digest() []byte {
	return ky.digest[:]
}

// Equals uses key digests to compare key equality.
func (ky *Key) Equals(other *Key) bool {
	return bytes.Equal(ky.digest[:], other.digest[:])
}

// String implements Stringer interface and returns string representation of key.
func (ky *Key) String() string {
	if ky == nil {
		return ""
	}

	if ky.userKey != nil {
		return fmt.Sprintf("%s:%s:%s:%v", ky.namespace, ky.setName, ky.userKey.String(), Buffer.BytesToHexString(ky.digest[:]))
	}
	return fmt.Sprintf("%s:%s::%v", ky.namespace, ky.setName, Buffer.BytesToHexString(ky.digest[:]))
}

// NewKey initializes a key from namespace, optional set name and user key.
// The set name and user defined key are converted to a digest before sending to the server.
// The server handles record identifiers by digest only.
func NewKey(namespace string, setName string, key interface{}) (*Key, Error) {
	newKey := &Key{
		namespace: namespace,
		setName:   setName,
		userKey:   NewValue(key),
	}

	if err := newKey.computeDigest(); err != nil {
		return nil, err
	}

	return newKey, nil
}

// NewKeyWithDigest initializes a key from namespace, optional set name and user key.
// The server handles record identifiers by digest only.
func NewKeyWithDigest(namespace string, setName string, key interface{}, digest []byte) (*Key, Error) {
	newKey := &Key{
		namespace: namespace,
		setName:   setName,
		userKey:   NewValue(key),
	}

	if err := newKey.SetDigest(digest); err != nil {
		return nil, err
	}
	return newKey, nil
}

// SetDigest sets a custom hash
func (ky *Key) SetDigest(digest []byte) Error {
	if len(digest) != 20 {
		return newError(types.PARAMETER_ERROR, "Invalid digest: Digest is required to be exactly 20 bytes.")
	}
	copy(ky.digest[:], digest)
	return nil
}

// Generate unique server hash value from set name, key type and user defined key.
// The hash function is RIPEMD-160 (a 160 bit hash).
func (ky *Key) computeDigest() Error {
	// With custom changes to the ripemd160 package,
	// now the following line does not allocate on the heap anymore/.
	ky.keyWriter.hash.Reset()

	if _, err := ky.keyWriter.Write([]byte(ky.setName)); err != nil {
		return err
	}

	if _, err := ky.keyWriter.Write([]byte{byte(ky.userKey.GetType())}); err != nil {
		return err
	}

	if err := ky.keyWriter.writeKey(ky.userKey); err != nil {
		return err
	}

	// With custom changes to the ripemd160 package,
	// the following line does not allocate on he heap anymore.
	ky.keyWriter.hash.Sum(ky.digest[:])
	return nil
}

// PartitionId returns the partition that the key belongs to.
func (ky *Key) PartitionId() int {
	// CAN'T USE MOD directly - mod will give negative numbers.
	// First AND makes positive and negative correctly, then mod.
	return int(Buffer.LittleBytesToInt32(ky.digest[:], 0)&0xFFFF) & (_PARTITIONS - 1)
}
