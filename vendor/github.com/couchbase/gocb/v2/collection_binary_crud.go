package gocb

import (
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

// BinaryCollection is a set of binary operations.
type BinaryCollection struct {
	collection *Collection
}

// AppendOptions are the options available to the Append operation.
type AppendOptions struct {
	Timeout         time.Duration
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	Cas             Cas
	RetryStrategy   RetryStrategy
}

func (c *Collection) binaryAppend(id string, val []byte, opts *AppendOptions) (mutOut *MutationResult, errOut error) {
	if opts == nil {
		opts = &AppendOptions{}
	}

	opm := c.newKvOpManager("Append", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Append(gocbcore.AdjoinOptions{
		Key:                    opm.DocumentID(),
		Value:                  val,
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(opts.Cas),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.AdjoinResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// Append appends a byte value to a document.
func (c *BinaryCollection) Append(id string, val []byte, opts *AppendOptions) (mutOut *MutationResult, errOut error) {
	return c.collection.binaryAppend(id, val, opts)
}

// PrependOptions are the options available to the Prepend operation.
type PrependOptions struct {
	Timeout         time.Duration
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	Cas             Cas
	RetryStrategy   RetryStrategy
}

func (c *Collection) binaryPrepend(id string, val []byte, opts *PrependOptions) (mutOut *MutationResult, errOut error) {
	if opts == nil {
		opts = &PrependOptions{}
	}

	opm := c.newKvOpManager("Prepend", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Prepend(gocbcore.AdjoinOptions{
		Key:                    opm.DocumentID(),
		Value:                  val,
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(opts.Cas),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.AdjoinResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// Prepend prepends a byte value to a document.
func (c *BinaryCollection) Prepend(id string, val []byte, opts *PrependOptions) (mutOut *MutationResult, errOut error) {
	return c.collection.binaryPrepend(id, val, opts)
}

// IncrementOptions are the options available to the Increment operation.
type IncrementOptions struct {
	Timeout time.Duration
	// Expiry is the length of time that the document will be stored in Couchbase.
	// A value of 0 will set the document to never expire.
	Expiry time.Duration
	// Initial, if non-negative, is the `initial` value to use for the document if it does not exist.
	// If present, this is the value that will be returned by a successful operation.
	Initial int64
	// Delta is the value to use for incrementing/decrementing if Initial is not present.
	Delta           uint64
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	Cas             Cas
	RetryStrategy   RetryStrategy
}

func (c *Collection) binaryIncrement(id string, opts *IncrementOptions) (countOut *CounterResult, errOut error) {
	if opts == nil {
		opts = &IncrementOptions{}
	}

	opm := c.newKvOpManager("Increment", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if opts.Initial >= 0 {
		realInitial = uint64(opts.Initial)
	}

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Increment(gocbcore.CounterOptions{
		Key:                    opm.DocumentID(),
		Delta:                  opts.Delta,
		Initial:                realInitial,
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(opts.Cas),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.CounterResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		countOut = &CounterResult{}
		countOut.cas = Cas(res.Cas)
		countOut.mt = opm.EnhanceMt(res.MutationToken)
		countOut.content = res.Value

		opm.Resolve(countOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// Increment performs an atomic addition for an integer document. Passing a
// non-negative `initial` value will cause the document to be created if it did not
// already exist.
func (c *BinaryCollection) Increment(id string, opts *IncrementOptions) (countOut *CounterResult, errOut error) {
	return c.collection.binaryIncrement(id, opts)
}

// DecrementOptions are the options available to the Decrement operation.
type DecrementOptions struct {
	Timeout time.Duration
	// Expiry is the length of time that the document will be stored in Couchbase.
	// A value of 0 will set the document to never expire.
	Expiry time.Duration
	// Initial, if non-negative, is the `initial` value to use for the document if it does not exist.
	// If present, this is the value that will be returned by a successful operation.
	Initial int64
	// Delta is the value to use for incrementing/decrementing if Initial is not present.
	Delta           uint64
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	Cas             Cas
	RetryStrategy   RetryStrategy
}

func (c *Collection) binaryDecrement(id string, opts *DecrementOptions) (countOut *CounterResult, errOut error) {
	if opts == nil {
		opts = &DecrementOptions{}
	}

	opm := c.newKvOpManager("Decrement", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if opts.Initial >= 0 {
		realInitial = uint64(opts.Initial)
	}

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Decrement(gocbcore.CounterOptions{
		Key:                    opm.DocumentID(),
		Delta:                  opts.Delta,
		Initial:                realInitial,
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(opts.Cas),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.CounterResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		countOut = &CounterResult{}
		countOut.cas = Cas(res.Cas)
		countOut.mt = opm.EnhanceMt(res.MutationToken)
		countOut.content = res.Value

		opm.Resolve(countOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// Decrement performs an atomic subtraction for an integer document. Passing a
// non-negative `initial` value will cause the document to be created if it did not
// already exist.
func (c *BinaryCollection) Decrement(id string, opts *DecrementOptions) (countOut *CounterResult, errOut error) {
	return c.collection.binaryDecrement(id, opts)
}
