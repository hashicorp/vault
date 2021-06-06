package gocb

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/couchbase/gocbcore/v9/memd"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

// LookupInOptions are the set of options available to LookupIn.
type LookupInOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy

	// Internal: This should never be used and is not supported.
	Internal struct {
		AccessDeleted bool
	}
}

// LookupIn performs a set of subdocument lookup operations on the document identified by id.
func (c *Collection) LookupIn(id string, ops []LookupInSpec, opts *LookupInOptions) (docOut *LookupInResult, errOut error) {
	if opts == nil {
		opts = &LookupInOptions{}
	}

	opm := c.newKvOpManager("LookupIn", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	return c.internalLookupIn(opm, ops, opts.Internal.AccessDeleted)
}

func (c *Collection) internalLookupIn(
	opm *kvOpManager,
	ops []LookupInSpec,
	accessDeleted bool,
) (docOut *LookupInResult, errOut error) {
	var subdocs []gocbcore.SubDocOp
	for _, op := range ops {
		if op.op == memd.SubDocOpGet && op.path == "" {
			if op.isXattr {
				return nil, errors.New("invalid xattr fetch with no path")
			}

			subdocs = append(subdocs, gocbcore.SubDocOp{
				Op:    memd.SubDocOpGetDoc,
				Flags: memd.SubdocFlag(SubdocFlagNone),
			})
			continue
		} else if op.op == memd.SubDocOpDictSet && op.path == "" {
			if op.isXattr {
				return nil, errors.New("invalid xattr set with no path")
			}

			subdocs = append(subdocs, gocbcore.SubDocOp{
				Op:    memd.SubDocOpSetDoc,
				Flags: memd.SubdocFlag(SubdocFlagNone),
			})
			continue
		}

		flags := memd.SubdocFlagNone
		if op.isXattr {
			flags |= memd.SubdocFlagXattrPath
		}

		subdocs = append(subdocs, gocbcore.SubDocOp{
			Op:    op.op,
			Path:  op.path,
			Flags: flags,
		})
	}

	var flags memd.SubdocDocFlag
	if accessDeleted {
		flags = memd.SubdocDocFlagAccessDeleted
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}

	err = opm.Wait(agent.LookupIn(gocbcore.LookupInOptions{
		Key:            opm.DocumentID(),
		Ops:            subdocs,
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan(),
		Deadline:       opm.Deadline(),
		Flags:          flags,
	}, func(res *gocbcore.LookupInResult, err error) {
		if err != nil && res == nil {
			errOut = opm.EnhanceErr(err)
		}

		if res != nil {
			docOut = &LookupInResult{}
			docOut.cas = Cas(res.Cas)
			docOut.contents = make([]lookupInPartial, len(subdocs))
			for i, opRes := range res.Ops {
				docOut.contents[i].err = opm.EnhanceErr(opRes.Err)
				docOut.contents[i].data = json.RawMessage(opRes.Value)
			}
		}

		if err == nil {
			opm.Resolve(nil)
		} else {
			opm.Reject()
		}
	}))
	if err != nil {
		errOut = err
	}
	return
}

// StoreSemantics is used to define the document level action to take during a MutateIn operation.
type StoreSemantics uint8

const (
	// StoreSemanticsReplace signifies to Replace the document, and fail if it does not exist.
	// This is the default action
	StoreSemanticsReplace StoreSemantics = iota

	// StoreSemanticsUpsert signifies to replace the document or create it if it doesn't exist.
	StoreSemanticsUpsert

	// StoreSemanticsInsert signifies to create the document, and fail if it exists.
	StoreSemanticsInsert
)

// MutateInOptions are the set of options available to MutateIn.
type MutateInOptions struct {
	Expiry          time.Duration
	Cas             Cas
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	StoreSemantic   StoreSemantics
	Timeout         time.Duration
	RetryStrategy   RetryStrategy

	// Internal: This should never be used and is not supported.
	Internal struct {
		AccessDeleted bool
	}
}

// MutateIn performs a set of subdocument mutations on the document specified by id.
func (c *Collection) MutateIn(id string, ops []MutateInSpec, opts *MutateInOptions) (mutOut *MutateInResult, errOut error) {
	if opts == nil {
		opts = &MutateInOptions{}
	}

	opm := c.newKvOpManager("MutateIn", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	return c.internalMutateIn(opm, opts.StoreSemantic, opts.Expiry, opts.Cas, ops, opts.Internal.AccessDeleted)
}

func jsonMarshalMultiArray(in interface{}) ([]byte, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	// Assert first character is a '['
	if len(out) < 2 || out[0] != '[' {
		return nil, makeInvalidArgumentsError("not a JSON array")
	}

	out = out[1 : len(out)-1]
	return out, nil
}

func jsonMarshalMutateSpec(op MutateInSpec) ([]byte, memd.SubdocFlag, error) {
	if op.value == nil {
		return nil, memd.SubdocFlagNone, nil
	}

	if macro, ok := op.value.(MutationMacro); ok {
		return []byte(macro), memd.SubdocFlagExpandMacros | memd.SubdocFlagXattrPath, nil
	}

	if op.multiValue {
		bytes, err := jsonMarshalMultiArray(op.value)
		return bytes, memd.SubdocFlagNone, err
	}

	bytes, err := json.Marshal(op.value)
	return bytes, memd.SubdocFlagNone, err
}

func (c *Collection) internalMutateIn(
	opm *kvOpManager,
	action StoreSemantics,
	expiry time.Duration,
	cas Cas,
	ops []MutateInSpec,
	accessDeleted bool,
) (mutOut *MutateInResult, errOut error) {
	var docFlags memd.SubdocDocFlag
	if action == StoreSemanticsReplace {
		// this is the default behaviour
	} else if action == StoreSemanticsUpsert {
		docFlags |= memd.SubdocDocFlagMkDoc
	} else if action == StoreSemanticsInsert {
		docFlags |= memd.SubdocDocFlagAddDoc
	} else {
		return nil, makeInvalidArgumentsError("invalid StoreSemantics value provided")
	}

	if accessDeleted {
		docFlags |= memd.SubdocDocFlagAccessDeleted
	}

	var subdocs []gocbcore.SubDocOp
	for _, op := range ops {
		if op.path == "" {
			switch op.op {
			case memd.SubDocOpDictAdd:
				return nil, makeInvalidArgumentsError("cannot specify a blank path with InsertSpec")
			case memd.SubDocOpDictSet:
				return nil, makeInvalidArgumentsError("cannot specify a blank path with UpsertSpec")
			case memd.SubDocOpDelete:
				return nil, makeInvalidArgumentsError("cannot specify a blank path with DeleteSpec")
			case memd.SubDocOpReplace:
				op.op = memd.SubDocOpSetDoc
			default:
			}
		}

		etrace := c.startKvOpTrace("encode", opm.TraceSpan())
		bytes, flags, err := jsonMarshalMutateSpec(op)
		etrace.Finish()
		if err != nil {
			return nil, err
		}

		if op.createPath {
			flags |= memd.SubdocFlagMkDirP
		}

		if op.isXattr {
			flags |= memd.SubdocFlagXattrPath
		}

		subdocs = append(subdocs, gocbcore.SubDocOp{
			Op:    op.op,
			Flags: flags,
			Path:  op.path,
			Value: bytes,
		})
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.MutateIn(gocbcore.MutateInOptions{
		Key:                    opm.DocumentID(),
		Flags:                  docFlags,
		Cas:                    gocbcore.Cas(cas),
		Ops:                    subdocs,
		Expiry:                 durationToExpiry(expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.MutateInResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutateInResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)
		mutOut.contents = make([]mutateInPartial, len(res.Ops))
		for i, op := range res.Ops {
			mutOut.contents[i] = mutateInPartial{data: op.Value}
		}

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}
