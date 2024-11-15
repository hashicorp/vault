package protocol

// MessageType represents the message type.
type MessageType int8

// MessageType constants.
const (
	mtNil             MessageType = 0
	MtExecuteDirect   MessageType = 2
	MtPrepare         MessageType = 3
	mtAbapStream      MessageType = 4
	mtXAStart         MessageType = 5
	mtXAJoin          MessageType = 6
	MtExecute         MessageType = 13
	MtWriteLob        MessageType = 16
	MtReadLob         MessageType = 17
	mtFindLob         MessageType = 18
	MtAuthenticate    MessageType = 65
	MtConnect         MessageType = 66
	MtCommit          MessageType = 67
	MtRollback        MessageType = 68
	MtCloseResultset  MessageType = 69
	MtDropStatementID MessageType = 70
	MtFetchNext       MessageType = 71
	mtFetchAbsolute   MessageType = 72
	mtFetchRelative   MessageType = 73
	mtFetchFirst      MessageType = 74
	mtFetchLast       MessageType = 75
	MtDisconnect      MessageType = 77
	mtExecuteITab     MessageType = 78
	mtFetchNextITab   MessageType = 79
	mtInsertNextITab  MessageType = 80
	mtBatchPrepare    MessageType = 81
	MtDBConnectInfo   MessageType = 82
	mtXopenXAStart    MessageType = 83
	mtXopenXAEnd      MessageType = 84
	mtXopenXAPrepare  MessageType = 85
	mtXopenXACommit   MessageType = 86
	mtXopenXARollback MessageType = 87
	mtXopenXARecover  MessageType = 88
	mtXopenXAForget   MessageType = 89
)

// ClientInfoSupported returns true if message does support client info, false otherwise.
func (mt MessageType) ClientInfoSupported() bool {
	/*
		mtConnect is only supported since 2.00.042
		As server version is only available after connect we do not use it
		to support especially version 1.00.122 until maintenance
		will end in sommer 2021

		return mt == mtConnect || mt == mtPrepare || mt == mtExecuteDirect || mt == mtExecute
	*/
	return mt == MtPrepare || mt == MtExecuteDirect || mt == MtExecute
}
