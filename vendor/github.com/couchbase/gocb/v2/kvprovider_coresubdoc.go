package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/couchbase/gocbcore/v10"
	"github.com/couchbase/gocbcore/v10/memd"
)

func (p *kvProviderCore) LookupIn(c *Collection, id string, ops []LookupInSpec, opts *LookupInOptions) (docOut *LookupInResult, errOut error) {
	opm := newKvOpManagerCore(c, "lookup_in", opts.ParentSpan, p)
	defer opm.Finish(opts.noMetrics)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	return p.internalLookupIn(opm, ops, memd.SubdocDocFlag(opts.Internal.DocFlags), 0)
}

func (p *kvProviderCore) LookupInAllReplicas(c *Collection, id string, ops []LookupInSpec,
	opts *LookupInAllReplicaOptions) (docOut *LookupInAllReplicasResult, errOut error) {
	var tracectx RequestSpanContext
	if opts.ParentSpan != nil {
		tracectx = opts.ParentSpan.Context()
	}

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	span := p.StartKvOpTrace(c, "lookup_in_all_replicas", tracectx, false)

	// Timeout needs to be adjusted here, since we use it at the bottom of this
	// function, but the remaining options are all passed downwards and get handled
	// by those functions rather than us.
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = c.timeoutsConfig.KVTimeout
	}

	deadline := time.Now().Add(timeout)
	retryStrategy := opts.RetryStrategy

	snapshot, err := p.snapshotProvider.WaitForConfigSnapshot(ctx, deadline)
	if err != nil {
		return nil, err
	}

	supportStatus, err := c.bucket.Internal().bucketCapabilityStatus(gocbcore.BucketCapabilityReplicaRead)
	if err != nil {
		return nil, err
	}
	if supportStatus == CapabilityStatusUnsupported {
		return nil, ErrFeatureNotAvailable
	}

	var servers []int
	if opts.ReadPreference == ReadPreferenceSelectedServerGroup {
		serverGroups, err := snapshot.KeyToServersByServerGroup([]byte(id))
		if err != nil {
			return nil, err
		}

		for group, srvIdx := range serverGroups {
			if group == p.preferredServerGroup {
				servers = append(servers, srvIdx...)
			}
		}
	} else {
		numReplicas, err := snapshot.NumReplicas()
		if err != nil {
			return nil, err
		}

		numServers := numReplicas + 1
		for i := 0; i < numServers; i++ {
			servers = append(servers, i)
		}
	}

	outCh := make(chan interface{}, len(servers))
	cancelCh := make(chan struct{})

	recorder, err := p.meter.ValueRecorder(meterValueServiceKV, "lookup_in_all_replicas")
	if err != nil {
		logDebugf("Failed to create value recorder: %v", err)
	}

	repRes := &LookupInAllReplicasResult{
		res: &coreReplicasResult{
			totalRequests:       uint32(len(servers)),
			resCh:               outCh,
			cancelCh:            cancelCh,
			span:                span,
			childReqsCompleteCh: make(chan struct{}),
			valueRecorder:       recorder,
			startedTime:         time.Now(),
		},
	}

	if len(servers) == 0 {
		// This can happen when the selected server group does not exist, or has not been set
		close(repRes.res.resCh)
		close(repRes.res.childReqsCompleteCh)
		close(repRes.res.cancelCh)
		repRes.res.span.End()

		return repRes, nil
	}

	// Loop all the servers and populate the result object
	for _, replicaIdx := range servers {
		go func(replicaIdx int) {
			// This timeout value will cause the getOneReplica operation to timeout after our deadline has expired,
			// as the deadline has already begun. getOneReplica timing out before our deadline would cause inconsistent
			// behaviour.
			res, err := p.lookupInOneReplica(ctx, span, id, ops, replicaIdx, retryStrategy, cancelCh,
				timeout, opts.Internal.User, c, memd.SubdocDocFlag(opts.Internal.DocFlags))
			if err != nil {
				repRes.res.addFailed()
				logDebugf("Failed to fetch replica from replica %d: %s", replicaIdx, err)
			} else {
				repRes.res.addResult(res)
			}
		}(replicaIdx)
	}

	// Start a timer to close it after the deadline
	go func() {
		select {
		case <-time.After(time.Until(deadline)):
			// If we timeout, we should close the result
			err := repRes.Close()
			if err != nil {
				logDebugf("failed to close LookupInAllReplicas response: %s", err)
			}
		case <-cancelCh:
		// If the cancel channel closes, we are done
		case <-ctx.Done():
			err := repRes.Close()
			if err != nil {
				logDebugf("failed to close LookupInAllReplicas response: %s", err)
			}
		}
	}()

	return repRes, nil
}

func (p *kvProviderCore) LookupInAnyReplica(c *Collection, id string, ops []LookupInSpec,
	opts *LookupInAnyReplicaOptions) (docOut *LookupInReplicaResult, errOut error) {
	start := time.Now()
	defer p.meter.ValueRecord("kv", "lookup_in_any_replica", start)

	var tracectx RequestSpanContext
	if opts.ParentSpan != nil {
		tracectx = opts.ParentSpan.Context()
	}

	span := p.StartKvOpTrace(c, "lookup_in_any_replica", tracectx, false)
	defer span.End()

	repRes, err := p.LookupInAllReplicas(c, id, ops, &LookupInAllReplicaOptions{
		Timeout:        opts.Timeout,
		RetryStrategy:  opts.RetryStrategy,
		Internal:       opts.Internal,
		ParentSpan:     span,
		Context:        opts.Context,
		ReadPreference: opts.ReadPreference,
	})
	if err != nil {
		return nil, err
	}

	// Try to fetch at least one result
	res := repRes.Next()
	if res == nil {
		return nil, &KeyValueError{
			InnerError:     ErrDocumentUnretrievable,
			BucketName:     c.bucketName(),
			ScopeName:      c.scope,
			CollectionName: c.collectionName,
		}
	}

	// Close the results channel since we don't care about any of the
	// remaining result objects at this point.
	err = repRes.Close()
	if err != nil {
		logDebugf("failed to close LookupInAnyReplica response: %s", err)
	}

	return res, nil
}

func (p *kvProviderCore) lookupInOneReplica(
	ctx context.Context,
	span RequestSpan,
	id string,
	ops []LookupInSpec,
	replicaIdx int,
	retryStrategy RetryStrategy,
	cancelCh chan struct{},
	timeout time.Duration,
	user string,
	c *Collection,
	flags memd.SubdocDocFlag,
) (*LookupInReplicaResult, error) {
	opm := newKvOpManagerCore(c, "lookup_in", span, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(retryStrategy)
	opm.SetTimeout(timeout)
	opm.SetImpersonate(user)
	opm.SetContext(ctx)
	opm.SetCancelCh(cancelCh)

	if replicaIdx == 0 {
		res, err := p.internalLookupIn(opm, ops, flags, 0)
		if err != nil {
			return nil, err
		}

		docOut := &LookupInReplicaResult{}
		docOut.LookupInResult = res

		return docOut, nil
	}

	newFlags := memd.SubdocDocFlagReplicaRead | flags
	res, err := p.internalLookupIn(opm, ops, newFlags, replicaIdx)
	if err != nil {
		return nil, err
	}

	docOut := &LookupInReplicaResult{}
	docOut.LookupInResult = res
	docOut.isReplica = true

	return docOut, nil
}

func (p *kvProviderCore) internalLookupIn(
	opm *kvOpManagerCore,
	ops []LookupInSpec,
	flags memd.SubdocDocFlag,
	replicaIdx int,
) (*LookupInResult, error) {
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

	var docOut *LookupInResult
	var errOut error
	err := opm.Wait(p.agent.LookupIn(gocbcore.LookupInOptions{
		Key:            opm.DocumentID(),
		Ops:            subdocs,
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		Flags:          flags,
		User:           opm.Impersonate(),
		ReplicaIdx:     replicaIdx,
	}, func(res *gocbcore.LookupInResult, err error) {
		if err != nil && res == nil {
			var kvErr *gocbcore.KeyValueError
			if errors.As(err, &kvErr) {
				if errors.Is(err, gocbcore.ErrMemdSubDocBadCombo) {
					kvErr.InnerError = ErrInvalidArgument
				}
			}
			errOut = opm.EnhanceErr(err)
		}

		if res != nil {
			docOut = &LookupInResult{}
			docOut.cas = Cas(res.Cas)
			docOut.contents = make([]lookupInPartial, len(subdocs))
			for i, opRes := range res.Ops {
				docOut.contents[i].op = ops[i].op
				if ops[i].op == memd.SubDocOpExists {
					if opRes.Err == nil {
						docOut.contents[i].data = []byte("true")
					} else if errors.Is(opRes.Err, ErrPathNotFound) {
						docOut.contents[i].data = []byte("false")
					} else {
						docOut.contents[i].err = opm.EnhanceErr(opRes.Err)
					}
				} else {
					docOut.contents[i].err = opm.EnhanceErr(opRes.Err)
					docOut.contents[i].data = opRes.Value
				}
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
	return docOut, errOut
}

func (p *kvProviderCore) MutateIn(c *Collection, id string, ops []MutateInSpec, opts *MutateInOptions) (mutOut *MutateInResult, errOut error) {
	opm := newKvOpManagerCore(c, "mutate_in", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)
	opm.SetPreserveExpiry(opts.PreserveExpiry)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	return p.internalMutateIn(opm, opts.StoreSemantic, opts.Expiry, opts.Cas, ops, memd.SubdocDocFlag(opts.Internal.DocFlags))
}

func (p *kvProviderCore) internalMutateIn(
	opm *kvOpManagerCore,
	action StoreSemantics,
	expiry time.Duration,
	cas Cas,
	ops []MutateInSpec,
	docFlags memd.SubdocDocFlag,
) (*MutateInResult, error) {
	preserveTTL := opm.PreserveExpiry()
	if action == StoreSemanticsReplace {
		// this is the default behaviour
		if expiry > 0 && preserveTTL {
			return nil, makeInvalidArgumentsError("cannot use preserve expiry with expiry for replace store semantics")
		}
	} else if action == StoreSemanticsUpsert {
		docFlags |= memd.SubdocDocFlagMkDoc
	} else if action == StoreSemanticsInsert {
		if preserveTTL {
			return nil, makeInvalidArgumentsError("cannot use preserve ttl with insert store semantics")
		}
		docFlags |= memd.SubdocDocFlagAddDoc
	} else {
		return nil, makeInvalidArgumentsError("invalid StoreSemantics value provided")
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
				op.op = memd.SubDocOpDeleteDoc
			case memd.SubDocOpReplace:
				op.op = memd.SubDocOpSetDoc
			default:
			}
		}

		etrace := opm.kv.StartKvOpTrace(opm.parent, "request_encoding", opm.TraceSpanContext(), true)
		bytes, flags, err := jsonMarshalMutateSpec(op)
		etrace.End()
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

	var mutOut *MutateInResult
	var errOut error
	err := opm.Wait(p.agent.MutateIn(gocbcore.MutateInOptions{
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
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
		PreserveExpiry:         preserveTTL,
	}, func(res *gocbcore.MutateInResult, err error) {
		if err != nil {
			var kvErr *gocbcore.KeyValueError
			if errors.As(err, &kvErr) {
				if errors.Is(kvErr.InnerError, ErrCasMismatch) {
					// GOCBC-1019: Due to a previous bug in gocbcore we need to convert cas mismatch back to exists.
					kvErr.InnerError = ErrDocumentExists
				} else if errors.Is(err, gocbcore.ErrMemdSubDocBadCombo) {
					kvErr.InnerError = ErrInvalidArgument
				}
			}
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
	return mutOut, errOut
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
		// If the mutation is to write, then this is a json `null` value
		switch op.op {
		case memd.SubDocOpDictAdd,
			memd.SubDocOpDictSet,
			memd.SubDocOpReplace,
			memd.SubDocOpArrayPushLast,
			memd.SubDocOpArrayPushFirst,
			memd.SubDocOpArrayInsert,
			memd.SubDocOpArrayAddUnique,
			memd.SubDocOpSetDoc,
			memd.SubDocOpAddDoc:
			return []byte("null"), memd.SubdocFlagNone, nil
		}

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
