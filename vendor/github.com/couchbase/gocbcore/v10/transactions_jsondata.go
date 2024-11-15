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
)

type jsonAtrState string

const (
	jsonAtrStateUnknown    = jsonAtrState("")
	jsonAtrStatePending    = jsonAtrState("PENDING")
	jsonAtrStateCommitted  = jsonAtrState("COMMITTED")
	jsonAtrStateCompleted  = jsonAtrState("COMPLETED")
	jsonAtrStateAborted    = jsonAtrState("ABORTED")
	jsonAtrStateRolledBack = jsonAtrState("ROLLED_BACK")
)

type jsonMutationType string

const (
	jsonMutationInsert  = jsonMutationType("insert")
	jsonMutationReplace = jsonMutationType("replace")
	jsonMutationRemove  = jsonMutationType("remove")
)

type jsonAtrMutation struct {
	BucketName     string `json:"bkt,omitempty"`
	ScopeName      string `json:"scp,omitempty"`
	CollectionName string `json:"col,omitempty"`
	DocID          string `json:"id,omitempty"`
}

type jsonAtrAttempt struct {
	TransactionID string `json:"tid,omitempty"`
	ExpiryTime    int    `json:"exp,omitempty"`
	State         string `json:"st,omitempty"`

	PendingCAS    string `json:"tst,omitempty"`
	CommitCAS     string `json:"tsc,omitempty"`
	CompletedCAS  string `json:"tsco,omitempty"`
	AbortCAS      string `json:"tsrs,omitempty"`
	RolledBackCAS string `json:"tsrc,omitempty"`

	Inserts  []jsonAtrMutation `json:"ins,omitempty"`
	Replaces []jsonAtrMutation `json:"rep,omitempty"`
	Removes  []jsonAtrMutation `json:"rem,omitempty"`

	DurabilityLevel string `json:"d,omitempty"`

	ForwardCompat map[string][]jsonForwardCompatibilityEntry `json:"fc,omitempty"`
}

type jsonTxnXattrID struct {
	Transaction string `json:"txn,omitempty"`
	Attempt     string `json:"atmpt,omitempty"`
}

type jsonTxnXattrATR struct {
	DocID          string `json:"id,omitempty"`
	BucketName     string `json:"bkt,omitempty"`
	CollectionName string `json:"coll,omitempty"`
	ScopeName      string `json:"scp,omitempty"`
}

type jsonTxnXattrOp struct {
	Type   jsonMutationType `json:"type,omitempty"`
	Staged json.RawMessage  `json:"stgd,omitempty"`
	CRC32  string           `json:"crc32,omitempty"`
}

type jsonTxnXattrRestore struct {
	OriginalCAS string `json:"CAS,omitempty"`
	ExpiryTime  uint   `json:"exptime"`
	RevID       string `json:"revid,omitempty"`
}

type jsonTxnXattr struct {
	ID            jsonTxnXattrID                             `json:"id,omitempty"`
	ATR           jsonTxnXattrATR                            `json:"atr,omitempty"`
	Operation     jsonTxnXattrOp                             `json:"op,omitempty"`
	Restore       *jsonTxnXattrRestore                       `json:"restore,omitempty"`
	ForwardCompat map[string][]jsonForwardCompatibilityEntry `json:"fc,omitempty"`
}

type transactionDocMeta struct {
	Cas        string `json:"CAS"`
	RevID      string `json:"revid"`
	Expiration uint   `json:"exptime"`
	CRC32      string `json:"value_crc32c,omitempty"`
}

type transactionGetDoc struct {
	Body    []byte
	TxnMeta *jsonTxnXattr
	DocMeta *transactionDocMeta
	Cas     Cas
	Deleted bool
}

type jsonForwardCompatibilityEntry struct {
	ProtocolVersion   string `json:"p,omitempty"`
	ProtocolExtension string `json:"e,omitempty"`
	Behaviour         string `json:"b,omitempty"`
	RetryInterval     int    `json:"ra,omitempty"`
}
