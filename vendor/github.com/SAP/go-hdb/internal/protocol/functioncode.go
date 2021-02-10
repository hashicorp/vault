// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

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

func (fc functionCode) isProcedureCall() bool {
	return fc == fcDBProcedureCall
}
