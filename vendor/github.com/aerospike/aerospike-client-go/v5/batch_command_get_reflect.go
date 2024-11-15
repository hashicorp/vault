// +build !as_performance

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
	"reflect"

	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

// if this file is included in the build, it will include this method
func init() {
	batchObjectParser = parseBatchObject
}

func parseBatchObject(
	cmd *batchCommandGet,
	offset int,
	opCount int,
	fieldCount int,
	generation uint32,
	expiration uint32,
) Error {
	if opCount > 0 {
		rv := *cmd.objects[offset]

		if rv.Kind() != reflect.Ptr {
			return ErrInvalidObjectType.err()
		}
		rv = rv.Elem()

		if !rv.CanAddr() {
			return ErrInvalidObjectType.err()
		}

		if rv.Kind() != reflect.Struct {
			return ErrInvalidObjectType.err()
		}

		// find the name based on tag mapping
		iobj := indirect(rv)
		mappings := objectMappings.getMapping(iobj.Type())

		if err := setObjectMetaFields(iobj, expiration, generation); err != nil {
			return err
		}

		for i := 0; i < opCount; i++ {
			if err := cmd.readBytes(8); err != nil {
				return err
			}
			opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
			particleType := int(cmd.dataBuffer[5])
			nameSize := int(cmd.dataBuffer[7])

			if err := cmd.readBytes(nameSize); err != nil {
				return err
			}
			name := string(cmd.dataBuffer[:nameSize])

			particleBytesSize := opSize - (4 + nameSize)
			if err := cmd.readBytes(particleBytesSize); err != nil {
				return err
			}
			value, err := bytesToParticle(particleType, cmd.dataBuffer, 0, particleBytesSize)
			if err != nil {
				return err
			}
			if err := setObjectField(mappings, iobj, name, value); err != nil {
				return err
			}
		}
	}

	return nil
}
