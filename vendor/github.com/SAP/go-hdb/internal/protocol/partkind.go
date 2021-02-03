// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

//go:generate stringer -type=partKind

type partKind int8

const (
	pkNil                       partKind = 0
	pkCommand                   partKind = 3
	pkResultset                 partKind = 5
	pkError                     partKind = 6
	pkStatementID               partKind = 10
	pkTransactionID             partKind = 11
	pkRowsAffected              partKind = 12
	pkResultsetID               partKind = 13
	pkTopologyInformation       partKind = 15
	pkTableLocation             partKind = 16
	pkReadLobRequest            partKind = 17
	pkReadLobReply              partKind = 18
	pkAbapIStream               partKind = 25
	pkAbapOStream               partKind = 26
	pkCommandInfo               partKind = 27
	pkWriteLobRequest           partKind = 28
	pkClientContext             partKind = 29
	pkWriteLobReply             partKind = 30
	pkParameters                partKind = 32
	pkAuthentication            partKind = 33
	pkSessionContext            partKind = 34
	pkClientID                  partKind = 35
	pkProfile                   partKind = 38
	pkStatementContext          partKind = 39
	pkPartitionInformation      partKind = 40
	pkOutputParameters          partKind = 41
	pkConnectOptions            partKind = 42
	pkCommitOptions             partKind = 43
	pkFetchOptions              partKind = 44
	pkFetchSize                 partKind = 45
	pkParameterMetadata         partKind = 47
	pkResultMetadata            partKind = 48
	pkFindLobRequest            partKind = 49
	pkFindLobReply              partKind = 50
	pkItabSHM                   partKind = 51
	pkItabChunkMetadata         partKind = 53
	pkItabMetadata              partKind = 55
	pkItabResultChunk           partKind = 56
	pkClientInfo                partKind = 57
	pkStreamData                partKind = 58
	pkOStreamResult             partKind = 59
	pkFDARequestMetadata        partKind = 60
	pkFDAReplyMetadata          partKind = 61
	pkBatchPrepare              partKind = 62 //Reserved: do not use
	pkBatchExecute              partKind = 63 //Reserved: do not use
	pkTransactionFlags          partKind = 64
	pkRowSlotImageParamMetadata partKind = 65 //Reserved: do not use
	pkRowSlotImageResultset     partKind = 66 //Reserved: do not use
	pkDBConnectInfo             partKind = 67
	pkLobFlags                  partKind = 68
	pkResultsetOptions          partKind = 69
	pkXATransactionInfo         partKind = 70
	pkSessionVariable           partKind = 71
	pkWorkLoadReplayContext     partKind = 72
	pkSQLReplyOptions           partKind = 73
)
