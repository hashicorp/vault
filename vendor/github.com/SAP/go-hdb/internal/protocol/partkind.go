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

//go:generate stringer -type=partKind

type partKind int8

const (
	pkNil                  partKind = 0
	pkCommand              partKind = 3
	pkResultset            partKind = 5
	pkError                partKind = 6
	pkStatementID          partKind = 10
	pkTransactionID        partKind = 11
	pkRowsAffected         partKind = 12
	pkResultsetID          partKind = 13
	pkTopologyInformation  partKind = 15
	pkTableLocation        partKind = 16
	pkReadLobRequest       partKind = 17
	pkReadLobReply         partKind = 18
	pkAbapIStream          partKind = 25
	pkAbapOStream          partKind = 26
	pkCommandInfo          partKind = 27
	pkWriteLobRequest      partKind = 28
	pkWriteLobReply        partKind = 30
	pkParameters           partKind = 32
	pkAuthentication       partKind = 33
	pkSessionContext       partKind = 34
	pkStatementContext     partKind = 39
	pkPartitionInformation partKind = 40
	pkOutputParameters     partKind = 41
	pkConnectOptions       partKind = 42
	pkCommitOptions        partKind = 43
	pkFetchOptions         partKind = 44
	pkFetchSize            partKind = 45
	pkParameterMetadata    partKind = 47
	pkResultMetadata       partKind = 48
	pkFindLobRequest       partKind = 49
	pkFindLobReply         partKind = 50
	pkItabSHM              partKind = 51
	pkItabChunkMetadata    partKind = 53
	pkItabMetadata         partKind = 55
	pkItabResultChunk      partKind = 56
	pkClientInfo           partKind = 57
	pkStreamData           partKind = 58
	pkOStreamResult        partKind = 59
	pkFDARequestMetadata   partKind = 60
	pkFDAReplyMetadata     partKind = 61
	pkTransactionFlags     partKind = 64
)
