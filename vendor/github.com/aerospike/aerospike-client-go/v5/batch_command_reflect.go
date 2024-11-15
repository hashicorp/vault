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

// this function will only be set if the performance flag is not passed for build
func init() {
	multiObjectParser = batchParseObject
	prepareReflectionData = concretePrepareReflectionData
}

func concretePrepareReflectionData(cmd *baseMultiCommand) {
	// if a channel is assigned, assign its value type
	if cmd.recordset != nil && !cmd.recordset.objChan.IsNil() {
		// this channel must be of type chan *T
		cmd.resObjType = cmd.recordset.objChan.Type().Elem().Elem()
		cmd.resObjMappings = objectMappings.getMapping(cmd.recordset.objChan.Type().Elem().Elem())

		cmd.selectCases = []reflect.SelectCase{
			{Dir: reflect.SelectSend, Chan: cmd.recordset.objChan},
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cmd.recordset.cancelled)},
		}
	}
}

func batchParseObject(
	cmd *baseMultiCommand,
	obj reflect.Value,
	opCount int,
	fieldCount int,
	generation uint32,
	expiration uint32,
) Error {
	for i := 0; i < opCount; i++ {
		if err := cmd.readBytes(8); err != nil {
			err = newNodeError(cmd.node, err)
			return err
		}

		opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
		particleType := int(cmd.dataBuffer[5])
		nameSize := int(cmd.dataBuffer[7])

		if err := cmd.readBytes(nameSize); err != nil {
			err = newNodeError(cmd.node, err)
			return err
		}
		name := string(cmd.dataBuffer[:nameSize])

		particleBytesSize := opSize - (4 + nameSize)
		if err := cmd.readBytes(particleBytesSize); err != nil {
			err = newNodeError(cmd.node, err)
			return err
		}
		value, err := bytesToParticle(particleType, cmd.dataBuffer, 0, particleBytesSize)
		if err != nil {
			err = newNodeError(cmd.node, err)
			return err
		}

		iobj := indirect(obj)
		if err := setObjectField(cmd.resObjMappings, iobj, name, value); err != nil {
			return err
		}

		if err := setObjectMetaFields(obj, expiration, generation); err != nil {
			return err
		}
	}

	return nil
}
