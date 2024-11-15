package protocol

// FunctionCode represents a function code.
type FunctionCode int16

// FunctionCode constants.
const (
	fcNil                       FunctionCode = 0
	FcDDL                       FunctionCode = 1
	fcInsert                    FunctionCode = 2
	fcUpdate                    FunctionCode = 3
	fcDelete                    FunctionCode = 4
	fcSelect                    FunctionCode = 5
	fcSelectForUpdate           FunctionCode = 6
	fcExplain                   FunctionCode = 7
	fcDBProcedureCall           FunctionCode = 8
	fcDBProcedureCallWithResult FunctionCode = 9
	fcFetch                     FunctionCode = 10
	fcCommit                    FunctionCode = 11
	fcRollback                  FunctionCode = 12
	fcSavepoint                 FunctionCode = 13
	fcConnect                   FunctionCode = 14
	fcWriteLob                  FunctionCode = 15
	fcReadLob                   FunctionCode = 16
	fcPing                      FunctionCode = 17 //reserved: do not use
	fcDisconnect                FunctionCode = 18
	fcCloseCursor               FunctionCode = 19
	fcFindLob                   FunctionCode = 20
	fcAbapStream                FunctionCode = 21
	fcXAStart                   FunctionCode = 22
	fcXAJoin                    FunctionCode = 23
)

// IsProcedureCall returns true if the function code is a procedure call, false otherwise.
func (fc FunctionCode) IsProcedureCall() bool {
	return fc == fcDBProcedureCall
}
