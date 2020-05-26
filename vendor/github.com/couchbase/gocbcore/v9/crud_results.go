package gocbcore

// GetResult encapsulates the result of a GetEx operation.
type GetResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas
}

// GetAndTouchResult encapsulates the result of a GetAndTouchEx operation.
type GetAndTouchResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas
}

// GetAndLockResult encapsulates the result of a GetAndLockEx operation.
type GetAndLockResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas
}

// GetReplicaResult encapsulates the result of a GetReplica operation.
type GetReplicaResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas
}

// TouchResult encapsulates the result of a TouchEx operation.
type TouchResult struct {
	Cas           Cas
	MutationToken MutationToken
}

// UnlockResult encapsulates the result of a UnlockEx operation.
type UnlockResult struct {
	Cas           Cas
	MutationToken MutationToken
}

// DeleteResult encapsulates the result of a DeleteEx operation.
type DeleteResult struct {
	Cas           Cas
	MutationToken MutationToken
}

// StoreResult encapsulates the result of a AddEx, SetEx or ReplaceEx operation.
type StoreResult struct {
	Cas           Cas
	MutationToken MutationToken
}

// AdjoinResult encapsulates the result of a AppendEx or PrependEx operation.
type AdjoinResult struct {
	Cas           Cas
	MutationToken MutationToken
}

// CounterResult encapsulates the result of a IncrementEx or DecrementEx operation.
type CounterResult struct {
	Value         uint64
	Cas           Cas
	MutationToken MutationToken
}

// GetRandomResult encapsulates the result of a GetRandomEx operation.
type GetRandomResult struct {
	Key      []byte
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas
}

// GetMetaResult encapsulates the result of a GetMetaEx operation.
type GetMetaResult struct {
	Value    []byte
	Flags    uint32
	Cas      Cas
	Expiry   uint32
	SeqNo    SeqNo
	Datatype uint8
	Deleted  uint32
}

// SetMetaResult encapsulates the result of a SetMetaEx operation.
type SetMetaResult struct {
	Cas           Cas
	MutationToken MutationToken
}

// DeleteMetaResult encapsulates the result of a DeleteMetaEx operation.
type DeleteMetaResult struct {
	Cas           Cas
	MutationToken MutationToken
}
