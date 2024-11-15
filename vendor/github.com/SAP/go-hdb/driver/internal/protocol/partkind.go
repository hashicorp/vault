package protocol

// PartKind represents the part kind.
type PartKind int8

// PartKind constants.
const (
	pkNil                       PartKind = 0
	PkCommand                   PartKind = 3
	PkResultset                 PartKind = 5
	pkError                     PartKind = 6
	PkStatementID               PartKind = 10
	pkTransactionID             PartKind = 11
	pkRowsAffected              PartKind = 12
	PkResultsetID               PartKind = 13
	PkTopologyInformation       PartKind = 15
	pkTableLocation             PartKind = 16
	PkReadLobRequest            PartKind = 17
	PkReadLobReply              PartKind = 18
	pkAbapIStream               PartKind = 25
	pkAbapOStream               PartKind = 26
	pkCommandInfo               PartKind = 27
	PkWriteLobRequest           PartKind = 28
	PkClientContext             PartKind = 29
	PkWriteLobReply             PartKind = 30
	PkParameters                PartKind = 32
	PkAuthentication            PartKind = 33
	pkSessionContext            PartKind = 34
	PkClientID                  PartKind = 35
	pkProfile                   PartKind = 38
	PkStatementContext          PartKind = 39
	pkPartitionInformation      PartKind = 40
	PkOutputParameters          PartKind = 41
	PkConnectOptions            PartKind = 42
	pkCommitOptions             PartKind = 43
	pkFetchOptions              PartKind = 44
	PkFetchSize                 PartKind = 45
	PkParameterMetadata         PartKind = 47
	PkResultMetadata            PartKind = 48
	pkFindLobRequest            PartKind = 49
	pkFindLobReply              PartKind = 50
	pkItabSHM                   PartKind = 51
	pkItabChunkMetadata         PartKind = 53
	pkItabMetadata              PartKind = 55
	pkItabResultChunk           PartKind = 56
	PkClientInfo                PartKind = 57
	pkStreamData                PartKind = 58
	pkOStreamResult             PartKind = 59
	pkFDARequestMetadata        PartKind = 60
	pkFDAReplyMetadata          PartKind = 61
	pkBatchPrepare              PartKind = 62 //Reserved: do not use
	pkBatchExecute              PartKind = 63 //Reserved: do not use
	PkTransactionFlags          PartKind = 64
	pkRowSlotImageParamMetadata PartKind = 65 //Reserved: do not use
	pkRowSlotImageResultset     PartKind = 66 //Reserved: do not use
	PkDBConnectInfo             PartKind = 67
	pkLobFlags                  PartKind = 68
	pkResultsetOptions          PartKind = 69
	pkXATransactionInfo         PartKind = 70
	pkSessionVariable           PartKind = 71
	pkWorkLoadReplayContext     PartKind = 72
	pkSQLReplyOptions           PartKind = 73
)
