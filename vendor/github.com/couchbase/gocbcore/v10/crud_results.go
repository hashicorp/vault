package gocbcore

// ResourceUnitResult describes the number of compute units used by an operation.
// Internal: This should never be used and is not supported.
type ResourceUnitResult struct {
	ReadUnits  uint16
	WriteUnits uint16
}

// GetResult encapsulates the result of a GetEx operation.
type GetResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// GetAndTouchResult encapsulates the result of a GetAndTouchEx operation.
type GetAndTouchResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// GetAndLockResult encapsulates the result of a GetAndLockEx operation.
type GetAndLockResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// GetReplicaResult encapsulates the result of a GetReplica operation.
type GetReplicaResult struct {
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// TouchResult encapsulates the result of a TouchEx operation.
type TouchResult struct {
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// UnlockResult encapsulates the result of a UnlockEx operation.
type UnlockResult struct {
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// DeleteResult encapsulates the result of a DeleteEx operation.
type DeleteResult struct {
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// StoreResult encapsulates the result of a AddEx, SetEx or ReplaceEx operation.
type StoreResult struct {
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// AdjoinResult encapsulates the result of a AppendEx or PrependEx operation.
type AdjoinResult struct {
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// CounterResult encapsulates the result of a IncrementEx or DecrementEx operation.
type CounterResult struct {
	Value         uint64
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// GetRandomResult encapsulates the result of a GetRandomEx operation.
type GetRandomResult struct {
	Key      []byte
	Value    []byte
	Flags    uint32
	Datatype uint8
	Cas      Cas

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
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

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// SetMetaResult encapsulates the result of a SetMetaEx operation.
type SetMetaResult struct {
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// DeleteMetaResult encapsulates the result of a DeleteMetaEx operation.
type DeleteMetaResult struct {
	Cas           Cas
	MutationToken MutationToken

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}
