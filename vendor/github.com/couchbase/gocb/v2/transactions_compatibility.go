package gocb

import (
	"github.com/couchbase/gocbcore/v10"
)

// TransactionsProtocolVersion returns the protocol version that this library supports.
func TransactionsProtocolVersion() string {
	return gocbcore.TransactionsProtocolVersion()
}

// TransactionsProtocolExtensions returns a list strings representing the various features
// that this specific version of the library supports within its protocol version.
func TransactionsProtocolExtensions() []string {
	return gocbcore.TransactionsProtocolExtensions()
}
