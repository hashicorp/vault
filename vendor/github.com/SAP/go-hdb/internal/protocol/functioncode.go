/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

//go:generate stringer -type=functionCode

type functionCode int16

const (
	fcNil                       functionCode = 0
	fcDDL                       functionCode = 1
	fcInsert                    functionCode = 2
	fcUpdate                    functionCode = 3
	fcDelete                    functionCode = 4
	fcSelect                    functionCode = 5
	fcSelectForUpdate           functionCode = 6
	fcExplain                   functionCode = 7
	fcDBProcedureCall           functionCode = 8
	fcDBProcedureCallWithResult functionCode = 9
	fcFetch                     functionCode = 10
	fcCommit                    functionCode = 11
	fcRollback                  functionCode = 12
	fcSavepoint                 functionCode = 13
	fcConnect                   functionCode = 14
	fcWriteLob                  functionCode = 15
	fcReadLob                   functionCode = 16
	fcPing                      functionCode = 17 //reserved: do not use
	fcDisconnect                functionCode = 18
	fcCloseCursor               functionCode = 19
	fcFindLob                   functionCode = 20
	fcAbapStream                functionCode = 21
	fcXAStart                   functionCode = 22
	fcXAJoin                    functionCode = 23
)

func (k functionCode) queryType() QueryType {

	switch k {
	default:
		return QtNone
	case fcSelect, fcSelectForUpdate:
		return QtSelect
	case fcDBProcedureCall:
		return QtProcedureCall
	}
}
