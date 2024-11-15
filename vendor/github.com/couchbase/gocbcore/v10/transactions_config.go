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
	"errors"
	"fmt"
	"time"
)

// TransactionDurabilityLevel specifies the durability level to use for a mutation.
type TransactionDurabilityLevel int

const (
	// TransactionDurabilityLevelUnknown indicates to use the default level.
	TransactionDurabilityLevelUnknown = TransactionDurabilityLevel(0)

	// TransactionDurabilityLevelNone indicates that no durability is needed.
	TransactionDurabilityLevelNone = TransactionDurabilityLevel(1)

	// TransactionDurabilityLevelMajority indicates the operation must be replicated to the majority.
	TransactionDurabilityLevelMajority = TransactionDurabilityLevel(2)

	// TransactionDurabilityLevelMajorityAndPersistToActive indicates the operation must be replicated
	// to the majority and persisted to the active server.
	TransactionDurabilityLevelMajorityAndPersistToActive = TransactionDurabilityLevel(3)

	// TransactionDurabilityLevelPersistToMajority indicates the operation must be persisted to the active server.
	TransactionDurabilityLevelPersistToMajority = TransactionDurabilityLevel(4)
)

func transactionDurabilityLevelToString(level TransactionDurabilityLevel) string {
	switch level {
	case TransactionDurabilityLevelUnknown:
		return "UNSET"
	case TransactionDurabilityLevelNone:
		return "NONE"
	case TransactionDurabilityLevelMajority:
		return "MAJORITY"
	case TransactionDurabilityLevelMajorityAndPersistToActive:
		return "MAJORITY_AND_PERSIST_TO_ACTIVE"
	case TransactionDurabilityLevelPersistToMajority:
		return "PERSIST_TO_MAJORITY"
	}
	return ""
}

func transactionDurabilityLevelFromString(level string) (TransactionDurabilityLevel, error) {
	switch level {
	case "UNSET":
		return TransactionDurabilityLevelUnknown, nil
	case "NONE":
		return TransactionDurabilityLevelNone, nil
	case "MAJORITY":
		return TransactionDurabilityLevelMajority, nil
	case "MAJORITY_AND_PERSIST_TO_ACTIVE":
		return TransactionDurabilityLevelMajorityAndPersistToActive, nil
	case "PERSIST_TO_MAJORITY":
		return TransactionDurabilityLevelPersistToMajority, nil
	}
	return TransactionDurabilityLevelUnknown, errors.New("invalid durability level string")
}

// TransactionATRLocation specifies a specific location where ATR entries should be
// placed when performing transactions.
type TransactionATRLocation struct {
	Agent          *Agent
	OboUser        string
	ScopeName      string
	CollectionName string
}

func (tlal TransactionATRLocation) build() string {
	if tlal.Agent == nil {
		return ""
	}

	scope := tlal.ScopeName
	if scope == "" {
		scope = "_default"
	}
	collection := tlal.CollectionName
	if collection == "" {
		collection = "_default"
	}

	return tlal.Agent.BucketName() + "." + scope + "." + collection
}

func (tlal TransactionATRLocation) String() string {
	if isLogRedactionLevelFull() || isLogRedactionLevelPartial() {
		return redactMetaData(tlal.build())
	}

	return tlal.build()
}

func (tlal TransactionATRLocation) redacted() interface{} {
	return redactMetaData(tlal.build())
}

// TransactionLostATRLocation specifies a specific location where lost transactions should
// attempt cleanup.
type TransactionLostATRLocation struct {
	BucketName     string
	ScopeName      string
	CollectionName string
}

func (tlal TransactionLostATRLocation) build() string {
	if tlal.BucketName == "" {
		return ""
	}

	scope := tlal.ScopeName
	if scope == "" {
		scope = "_default"
	}
	collection := tlal.CollectionName
	if collection == "" {
		collection = "_default"
	}

	return tlal.BucketName + "." + scope + "." + collection
}

func (tlal TransactionLostATRLocation) String() string {
	if isLogRedactionLevelFull() || isLogRedactionLevelPartial() {
		return redactMetaData(tlal.build())
	}

	return tlal.build()
}

func (tlal TransactionLostATRLocation) redacted() interface{} {
	return redactMetaData(tlal.build())
}

// TransactionsBucketAgentProviderFn is a function used to provide an agent for
// a particular bucket by name.
type TransactionsBucketAgentProviderFn func(bucketName string) (*Agent, string, error)

// TransactionsLostCleanupATRLocationProviderFn is a function used to provide a list of ATRLocations for
// lost transactions cleanup.
type TransactionsLostCleanupATRLocationProviderFn func() ([]TransactionLostATRLocation, error)

// TransactionsConfig specifies various tunable options related to transactions.
type TransactionsConfig struct {
	// CustomATRLocation specifies a specific location to place meta-data.
	CustomATRLocation TransactionATRLocation

	// ExpirationTime sets the maximum time that transactions created
	// by this TransactionsManager object can run for, before expiring.
	ExpirationTime time.Duration

	// DurabilityLevel specifies the durability level that should be used
	// for all write operations performed by this TransactionsManager object.
	DurabilityLevel TransactionDurabilityLevel

	// KeyValueTimeout specifies the default timeout used for all KV writes.
	KeyValueTimeout time.Duration

	// CleanupWindow specifies how often to the cleanup process runs
	// attempting to garbage collection transactions that have failed but
	// were not cleaned up by the previous client.
	CleanupWindow time.Duration

	// CleanupClientAttempts controls where any transaction attempts made
	// by this client are automatically removed.
	CleanupClientAttempts bool

	// CleanupLostAttempts controls where a background process is created
	// to cleanup any ‘lost’ transaction attempts.
	CleanupLostAttempts bool

	// CleanupQueueSize controls the maximum queue size for the cleanup thread.
	CleanupQueueSize uint32

	// BucketAgentProvider provides a function which returns an agent for
	// a particular bucket by name.
	BucketAgentProvider TransactionsBucketAgentProviderFn

	// LostCleanupATRLocationProvider provides a function which returns a list of LostATRLocations
	// for use in lost transaction cleanup.
	LostCleanupATRLocationProvider TransactionsLostCleanupATRLocationProviderFn

	// CleanupWatchATRs is *NOT* used within the codebase, it is *only* here to provide API level backward
	// compatibility.
	// This should *never* be used.
	CleanupWatchATRs bool

	// Internal specifies a set of options for internal use.
	// Internal: This should never be used and is not supported.
	Internal struct {
		Hooks                   TransactionHooks
		CleanUpHooks            TransactionCleanUpHooks
		ClientRecordHooks       TransactionClientRecordHooks
		EnableNonFatalGets      bool
		EnableParallelUnstaging bool
		EnableExplicitATRs      bool
		EnableMutationCaching   bool
		NumATRs                 int
	}
}

func (config *TransactionsConfig) String() string {
	if config == nil {
		return "<nil>"
	}

	return fmt.Sprintf("CustomATRLocation:%s ExpirationTime:%s DurabilityLevel:%s KeyValueTimeout:%s CleanupWindow:%s "+
		"CleanupClientAttempts:%t CleanupLostAttempts:%t CleanupQueueSize:%d BucketAgentProvider:%p LostCleanupATRLocationProvider:%p "+
		"Internal:{EnableNonFatalGets:%t EnableParallelUnstaging:%t "+"EnableExplicitATRs:%t EnableMutationCaching:%t NumATRs:%d}",
		config.CustomATRLocation, config.ExpirationTime, transactionDurabilityLevelToString(config.DurabilityLevel),
		config.KeyValueTimeout, config.CleanupWindow, config.CleanupClientAttempts, config.CleanupLostAttempts, config.CleanupQueueSize,
		config.BucketAgentProvider, config.LostCleanupATRLocationProvider, config.Internal.EnableNonFatalGets,
		config.Internal.EnableParallelUnstaging, config.Internal.EnableExplicitATRs, config.Internal.EnableMutationCaching,
		config.Internal.NumATRs)
}

// TransactionOptions specifies options which can be overridden on a per transaction basis.
type TransactionOptions struct {
	// CustomATRLocation specifies a specific location to place meta-data.
	CustomATRLocation TransactionATRLocation

	// ExpirationTime sets the maximum time that this transaction will
	// run for, before expiring.
	ExpirationTime time.Duration

	// DurabilityLevel specifies the durability level that should be used
	// for all write operations performed by this transaction.
	DurabilityLevel TransactionDurabilityLevel

	// KeyValueTimeout specifies the timeout used for all KV writes.
	KeyValueTimeout time.Duration

	// BucketAgentProvider provides a function which returns an agent for
	// a particular bucket by name.
	BucketAgentProvider TransactionsBucketAgentProviderFn

	// TransactionLogger is the logger to use with this transaction.
	// Uncommitted: This API may change in the future.
	TransactionLogger TransactionLogger

	// Internal specifies a set of options for internal use.
	// Internal: This should never be used and is not supported.
	Internal struct {
		Hooks                TransactionHooks
		ResourceUnitCallback func(result *ResourceUnitResult)
	}
}

func (opts *TransactionOptions) String() string {
	if opts == nil {
		return "<nil>"
	}

	return fmt.Sprintf("CustomATRLocation:%s ExpirationTime:%s DurabilityLevel:%s KeyValueTimeout:%s "+
		"BucketAgentProvider:%p TransactionLogger:%p ",
		opts.CustomATRLocation, opts.ExpirationTime, transactionDurabilityLevelToString(opts.DurabilityLevel),
		opts.KeyValueTimeout, opts.BucketAgentProvider, opts.TransactionLogger)
}
