package gocb

import (
	"time"
)

type kvProvider interface {
	Insert(*Collection, string, interface{}, *InsertOptions) (*MutationResult, error)   // Done
	Upsert(*Collection, string, interface{}, *UpsertOptions) (*MutationResult, error)   // Done
	Replace(*Collection, string, interface{}, *ReplaceOptions) (*MutationResult, error) // Done
	Remove(*Collection, string, *RemoveOptions) (*MutationResult, error)                // Done

	Get(*Collection, string, *GetOptions) (*GetResult, error)                                // Done
	Exists(*Collection, string, *ExistsOptions) (*ExistsResult, error)                       // Done
	GetAndTouch(*Collection, string, time.Duration, *GetAndTouchOptions) (*GetResult, error) // Done
	GetAndLock(*Collection, string, time.Duration, *GetAndLockOptions) (*GetResult, error)   // Done
	Unlock(*Collection, string, Cas, *UnlockOptions) error                                   // Done
	Touch(*Collection, string, time.Duration, *TouchOptions) (*MutationResult, error)        // Done

	GetAnyReplica(c *Collection, id string, opts *GetAnyReplicaOptions) (*GetReplicaResult, error)
	GetAllReplicas(*Collection, string, *GetAllReplicaOptions) (*GetAllReplicasResult, error)

	LookupIn(*Collection, string, []LookupInSpec, *LookupInOptions) (*LookupInResult, error)
	LookupInAnyReplica(*Collection, string, []LookupInSpec, *LookupInAnyReplicaOptions) (*LookupInReplicaResult, error)
	LookupInAllReplicas(*Collection, string, []LookupInSpec, *LookupInAllReplicaOptions) (*LookupInAllReplicasResult, error)
	MutateIn(*Collection, string, []MutateInSpec, *MutateInOptions) (*MutateInResult, error)

	Increment(*Collection, string, *IncrementOptions) (*CounterResult, error)      // Done
	Decrement(*Collection, string, *DecrementOptions) (*CounterResult, error)      // Done
	Append(*Collection, string, []byte, *AppendOptions) (*MutationResult, error)   // Done
	Prepend(*Collection, string, []byte, *PrependOptions) (*MutationResult, error) // Done

	Scan(*Collection, ScanType, *ScanOptions) (*ScanResult, error)

	StartKvOpTrace(*Collection, string, RequestSpanContext, bool) RequestSpan
}

type kvBulkProvider interface {
	Do(*Collection, []BulkOp, *BulkOpOptions) error
}
