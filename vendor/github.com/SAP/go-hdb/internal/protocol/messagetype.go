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

//go:generate stringer -type=messageType

type messageType int8

const (
	mtNil             messageType = 0
	mtExecuteDirect   messageType = 2
	mtPrepare         messageType = 3
	mtAbapStream      messageType = 4
	mtXAStart         messageType = 5
	mtXAJoin          messageType = 6
	mtExecute         messageType = 13
	mtWriteLob        messageType = 16
	mtReadLob         messageType = 17
	mtFindLob         messageType = 18
	mtAuthenticate    messageType = 65
	mtConnect         messageType = 66
	mtCommit          messageType = 67
	mtRollback        messageType = 68
	mtCloseResultset  messageType = 69
	mtDropStatementID messageType = 70
	mtFetchNext       messageType = 71
	mtFetchAbsolute   messageType = 72
	mtFetchRelative   messageType = 73
	mtFetchFirst      messageType = 74
	mtFetchLast       messageType = 75
	mtDisconnect      messageType = 77
	mtExecuteITab     messageType = 78
	mtFetchNextITab   messageType = 79
	mtInsertNextITab  messageType = 80
)
