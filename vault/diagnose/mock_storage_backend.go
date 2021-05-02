package diagnose

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
)

<<<<<<< HEAD
const timeoutCallRead string = "lag Read"
const timeoutCallWrite string = "lag Write"
const timeoutCallDelete string = "lag Delete"
const errCallWrite string = "err Write"
const errCallDelete string = "err Delete"
const errCallRead string = "err Read"
const badReadCall string = "bad Read"
const storageErrStringWrite string = "storage error on write"
const storageErrStringRead string = "storage error on read"
const storageErrStringDelete string = "storage error on delete"
const readOp string = "read"
const writeOp string = "write"
const deleteOp string = "delete"
=======
const (
	timeoutCallRead        string = "lag Read"
	timeoutCallWrite       string = "lag Write"
	timeoutCallDelete      string = "lag Delete"
	errCallWrite           string = "err Write"
	errCallDelete          string = "err Delete"
	errCallRead            string = "err Read"
	badReadCall            string = "bad Read"
	storageErrStringWrite  string = "storage error on write"
	storageErrStringRead   string = "storage error on read"
	storageErrStringDelete string = "storage error on delete"
	readOp                 string = "read"
	writeOp                string = "write"
	deleteOp               string = "delete"
)
>>>>>>> master

var goodEntry physical.Entry = physical.Entry{Key: secretKey, Value: []byte(secretVal)}
var badEntry physical.Entry = physical.Entry{}

type mockStorageBackend struct {
	callType string
}

func (m mockStorageBackend) storageLogicGeneralInternal(op string) error {
	if (m.callType == timeoutCallRead && op == readOp) || (m.callType == timeoutCallWrite && op == writeOp) ||
		(m.callType == timeoutCallDelete && op == deleteOp) {
		time.Sleep(25 * time.Second)
	} else if m.callType == errCallWrite && op == writeOp {
		return fmt.Errorf(storageErrStringWrite)
	} else if m.callType == errCallDelete && op == deleteOp {
		return fmt.Errorf(storageErrStringDelete)
	} else if m.callType == errCallRead && op == readOp {
		return fmt.Errorf(storageErrStringRead)
	}

	return nil
}

// Put is used to insert or update an entry
func (m mockStorageBackend) Put(ctx context.Context, entry *physical.Entry) error {
	return m.storageLogicGeneralInternal(writeOp)
}

// Get is used to fetch an entry
func (m mockStorageBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	if m.callType == errCallRead || m.callType == timeoutCallRead {
		return nil, m.storageLogicGeneralInternal(readOp)
	}
	if m.callType == badReadCall {
		return &badEntry, nil
	}
	return &goodEntry, nil

}

// Delete is used to permanently delete an entry
func (m mockStorageBackend) Delete(ctx context.Context, key string) error {

	return m.storageLogicGeneralInternal(deleteOp)
}

// List is not used in a mock.
func (m mockStorageBackend) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, fmt.Errorf("method not implemented")
}
