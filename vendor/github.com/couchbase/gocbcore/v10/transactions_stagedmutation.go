// Copyright 2021 Couchbase
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocbcore

import (
	"encoding/json"
	"errors"
)

// TransactionStagedMutationType represents the type of a mutation performed in a transaction.
type TransactionStagedMutationType int

const (
	// TransactionStagedMutationUnknown indicates an error has occured.
	TransactionStagedMutationUnknown = TransactionStagedMutationType(0)

	// TransactionStagedMutationInsert indicates the staged mutation was an insert operation.
	TransactionStagedMutationInsert = TransactionStagedMutationType(1)

	// TransactionStagedMutationReplace indicates the staged mutation was an replace operation.
	TransactionStagedMutationReplace = TransactionStagedMutationType(2)

	// TransactionStagedMutationRemove indicates the staged mutation was an remove operation.
	TransactionStagedMutationRemove = TransactionStagedMutationType(3)
)

func transactionStagedMutationTypeToString(mtype TransactionStagedMutationType) string {
	switch mtype {
	case TransactionStagedMutationInsert:
		return "INSERT"
	case TransactionStagedMutationReplace:
		return "REPLACE"
	case TransactionStagedMutationRemove:
		return "REMOVE"
	}
	return ""
}

func transactionStagedMutationTypeFromString(mtype string) (TransactionStagedMutationType, error) {
	switch mtype {
	case "INSERT":
		return TransactionStagedMutationInsert, nil
	case "REPLACE":
		return TransactionStagedMutationReplace, nil
	case "REMOVE":
		return TransactionStagedMutationRemove, nil
	}
	return TransactionStagedMutationUnknown, errors.New("invalid mutation type string")
}

// TransactionStagedMutation wraps all of the information about a mutation which has been staged
// as part of the transaction and which should later be unstaged when the transaction
// has been committed.
type TransactionStagedMutation struct {
	OpType         TransactionStagedMutationType
	BucketName     string
	ScopeName      string
	CollectionName string
	Key            []byte
	Cas            Cas
	Staged         json.RawMessage
}

type transactionStagedMutation struct {
	OpType         TransactionStagedMutationType
	Agent          *Agent
	OboUser        string
	ScopeName      string
	CollectionName string
	Key            []byte
	Cas            Cas
	Staged         json.RawMessage
}
