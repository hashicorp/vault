//nolint:unused
package gocb

import (
	"context"
	"errors"
	"io"
	"time"

	"google.golang.org/grpc/status"

	"github.com/couchbase/gocbcore/v10"
	"github.com/couchbase/gocbcore/v10/memd"
	"github.com/couchbase/goprotostellar/genproto/kv_v1"
)

var _ kvProvider = &kvProviderPs{}

// // wraps kv and makes it compliant for gocb
type kvProviderPs struct {
	client kv_v1.KvServiceClient

	tracer RequestTracer
	meter  *meterWrapper
}

// this is uint8 to int32, need to check overflows etc
var memdToPsLookupinTranslation = map[memd.SubDocOpType]kv_v1.LookupInRequest_Spec_Operation{
	memd.SubDocOpGet:      kv_v1.LookupInRequest_Spec_OPERATION_GET,
	memd.SubDocOpExists:   kv_v1.LookupInRequest_Spec_OPERATION_EXISTS,
	memd.SubDocOpGetCount: kv_v1.LookupInRequest_Spec_OPERATION_COUNT,
}

var memdToPsMutateinTranslation = map[memd.SubDocOpType]kv_v1.MutateInRequest_Spec_Operation{
	memd.SubDocOpDictAdd:        kv_v1.MutateInRequest_Spec_OPERATION_INSERT,
	memd.SubDocOpDictSet:        kv_v1.MutateInRequest_Spec_OPERATION_UPSERT,
	memd.SubDocOpReplace:        kv_v1.MutateInRequest_Spec_OPERATION_REPLACE,
	memd.SubDocOpDelete:         kv_v1.MutateInRequest_Spec_OPERATION_REMOVE,
	memd.SubDocOpArrayPushFirst: kv_v1.MutateInRequest_Spec_OPERATION_ARRAY_PREPEND,
	memd.SubDocOpArrayPushLast:  kv_v1.MutateInRequest_Spec_OPERATION_ARRAY_APPEND,
	memd.SubDocOpArrayInsert:    kv_v1.MutateInRequest_Spec_OPERATION_ARRAY_INSERT,
	memd.SubDocOpArrayAddUnique: kv_v1.MutateInRequest_Spec_OPERATION_ARRAY_ADD_UNIQUE,
	memd.SubDocOpCounter:        kv_v1.MutateInRequest_Spec_OPERATION_COUNTER,
}

func (p *kvProviderPs) LookupIn(c *Collection, id string, ops []LookupInSpec, opts *LookupInOptions) (*LookupInResult, error) {
	opm := newKvOpManagerPs(c, "lookup_in", opts.ParentSpan, p)
	defer opm.Finish(opts.noMetrics)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)
	opm.SetIsIdempotent(true)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	lookUpInPSSpecs := make([]*kv_v1.LookupInRequest_Spec, len(ops))

	for i, op := range ops {

		if op.op == memd.SubDocOpGet && op.path == "" {
			if op.isXattr {
				return nil, makeInvalidArgumentsError("invalid xattr fetch with no path")
			}

			lookUpInPSSpecs[i] = &kv_v1.LookupInRequest_Spec{
				Operation: memdToPsLookupinTranslation[op.op],
				Path:      op.path,
			}

			continue
		}

		newOp, ok := memdToPsLookupinTranslation[op.op]
		if !ok {
			return nil, makeInvalidArgumentsError("unknown lookupin op")
		}

		isXattr := op.isXattr
		specFlag := &kv_v1.LookupInRequest_Spec_Flags{
			Xattr: &isXattr,
		}

		lookUpInPSSpecs[i] = &kv_v1.LookupInRequest_Spec{
			Operation: newOp,
			Path:      op.path,
			Flags:     specFlag,
		}
	}

	requestFlags := &kv_v1.LookupInRequest_Flags{}
	accessDeleted := opts.Internal.DocFlags|SubdocDocFlagAccessDeleted == 1
	if accessDeleted {
		requestFlags.AccessDeleted = &accessDeleted
	}

	req := &kv_v1.LookupInRequest{
		BucketName:     opm.BucketName(),
		ScopeName:      opm.ScopeName(),
		CollectionName: opm.CollectionName(),
		Key:            opm.DocumentID(),
		Specs:          lookUpInPSSpecs,
		Flags:          requestFlags,
	}

	res, err := wrapPSOp(opm, req, p.client.LookupIn)
	if err != nil {
		return nil, err
	}

	docOut := &LookupInResult{}
	docOut.cas = Cas(res.Cas)
	docOut.contents = make([]lookupInPartial, len(lookUpInPSSpecs))
	for i, opRes := range res.Specs {
		docOut.contents[i].op = ops[i].op
		if opRes.Status != nil && opRes.Status.Code != 0 {
			docOut.contents[i].err = opm.EnhanceErrorStatus(status.FromProto(opRes.Status), true)
		}
		docOut.contents[i].data = opRes.Content
	}

	return docOut, nil
}

func (p *kvProviderPs) LookupInAnyReplica(*Collection, string, []LookupInSpec, *LookupInAnyReplicaOptions) (*LookupInReplicaResult, error) {
	return nil, ErrFeatureNotAvailable
}

func (p *kvProviderPs) LookupInAllReplicas(*Collection, string, []LookupInSpec, *LookupInAllReplicaOptions) (*LookupInAllReplicasResult, error) {
	return nil, ErrFeatureNotAvailable
}

func (p *kvProviderPs) MutateIn(c *Collection, id string, ops []MutateInSpec, opts *MutateInOptions) (*MutateInResult, error) {
	opm := newKvOpManagerPs(c, "mutate_in", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	psSpecs := make([]*kv_v1.MutateInRequest_Spec, len(ops))
	memdDocFlags := memd.SubdocDocFlag(opts.Internal.DocFlags)

	var psAction kv_v1.MutateInRequest_StoreSemantic
	switch opts.StoreSemantic {
	case StoreSemanticsReplace:
		// this is default behavior
		psAction = kv_v1.MutateInRequest_STORE_SEMANTIC_REPLACE
		if opts.Expiry > 0 && opts.PreserveExpiry {
			return nil, makeInvalidArgumentsError("cannot use preserve expiry with expiry for replace store semantics")
		}
	case StoreSemanticsUpsert:
		memdDocFlags |= memd.SubdocDocFlagMkDoc
		psAction = kv_v1.MutateInRequest_STORE_SEMANTIC_UPSERT
	case StoreSemanticsInsert:
		memdDocFlags |= memd.SubdocDocFlagAddDoc
		psAction = kv_v1.MutateInRequest_STORE_SEMANTIC_INSERT
	default:
		return nil, makeInvalidArgumentsError("invalid StoreSemantics value provided")
	}

	for i, op := range ops {
		// does PS take care of this?
		if op.path == "" {
			switch op.op {
			case memd.SubDocOpDictAdd:
				return nil, makeInvalidArgumentsError("cannot specify a blank path with InsertSpec")
			case memd.SubDocOpDictSet:
				return nil, makeInvalidArgumentsError("cannot specify a blank path with UpsertSpec")
			default:
			}
		}
		etrace := opm.provider.StartKvOpTrace(opm.parent, "request_encoding", opm.TraceSpanContext(), true)
		bytes, flags, err := jsonMarshalMutateSpec(op)
		etrace.End()
		if err != nil {
			return nil, err
		}

		if flags&memd.SubdocFlagExpandMacros == memd.SubdocFlagExpandMacros {
			return nil, wrapError(ErrFeatureNotAvailable, "unsupported flag for  the couchbase2 protocol: macro expansion")
		}
		createPath := op.createPath
		isXattr := op.isXattr
		psmutateFlag := &kv_v1.MutateInRequest_Spec_Flags{
			CreatePath: &createPath,
			Xattr:      &isXattr,
		}

		psSpecs[i] = &kv_v1.MutateInRequest_Spec{
			Operation: memdToPsMutateinTranslation[op.op],
			Path:      op.path,
			Content:   bytes,
			Flags:     psmutateFlag,
		}
	}

	accessDeleted := memdDocFlags&memd.SubdocDocFlagAccessDeleted == memd.SubdocDocFlagAccessDeleted
	var mutateInRequestFlags *kv_v1.MutateInRequest_Flags
	if accessDeleted {
		mutateInRequestFlags = &kv_v1.MutateInRequest_Flags{
			AccessDeleted: &accessDeleted,
		}
	}

	var cas *uint64
	if opts.Cas > 0 {
		cas = (*uint64)(&opts.Cas)
	}

	// CNG interprets a nil expiry value as "preserve expiry"
	var expiry *kv_v1.MutateInRequest_ExpirySecs
	if !opts.PreserveExpiry {
		if opts.Expiry > 0 {
			expiry = &kv_v1.MutateInRequest_ExpirySecs{ExpirySecs: uint32(opts.Expiry.Seconds())}
		} else {
			expiry = &kv_v1.MutateInRequest_ExpirySecs{ExpirySecs: 0}
		}
	}

	request := &kv_v1.MutateInRequest{
		BucketName:      opm.BucketName(),
		ScopeName:       opm.ScopeName(),
		CollectionName:  opm.CollectionName(),
		Key:             opm.DocumentID(),
		Specs:           psSpecs,
		StoreSemantic:   &psAction,
		DurabilityLevel: opm.DurabilityLevel(),
		Cas:             cas,
		Flags:           mutateInRequestFlags,
		Expiry:          expiry,
	}

	res, err := wrapPSOp(opm, request, p.client.MutateIn)
	if err != nil {
		if kvErr, ok := err.(*GenericError); ok {
			if errors.Is(kvErr.InnerError, ErrCasMismatch) {
				kvErr.InnerError = ErrDocumentExists
			}
		}

		return nil, err
	}

	mutOut := &MutateInResult{}
	mutOut.cas = Cas(res.Cas)
	mutOut.mt = psMutToGoCbMut(res.MutationToken)
	mutOut.contents = make([]mutateInPartial, len(res.Specs))
	for i, op := range res.Specs {
		mutOut.contents[i] = mutateInPartial{data: op.Content}
	}

	return mutOut, nil
}

func (p *kvProviderPs) Scan(c *Collection, scanType ScanType, opts *ScanOptions) (*ScanResult, error) {
	return nil, ErrFeatureNotAvailable
}

func (p *kvProviderPs) Insert(c *Collection, id string, val interface{}, opts *InsertOptions) (*MutationResult, error) {
	opm := newKvOpManagerPs(c, "insert", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var expiry *kv_v1.InsertRequest_ExpirySecs
	if opts.Expiry > 0 {
		expiry = &kv_v1.InsertRequest_ExpirySecs{ExpirySecs: uint32(opts.Expiry.Seconds())}
	}

	request := &kv_v1.InsertRequest{
		BucketName:     opm.BucketName(),
		ScopeName:      opm.ScopeName(),
		CollectionName: opm.CollectionName(),

		Key:          opm.DocumentID(),
		ContentFlags: opm.ValueFlags(),

		Expiry:          expiry,
		DurabilityLevel: opm.DurabilityLevel(),
	}

	value, isCompressed := opm.Value()
	if isCompressed {
		request.Content = &kv_v1.InsertRequest_ContentCompressed{ContentCompressed: value}
	} else {
		request.Content = &kv_v1.InsertRequest_ContentUncompressed{ContentUncompressed: value}
	}

	res, err := wrapPSOp(opm, request, p.client.Insert)
	if err != nil {
		return nil, err
	}

	mt := psMutToGoCbMut(res.MutationToken)
	cas := res.Cas

	mutOut := MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(cas),
		},
	}
	return &mutOut, nil
}

func (p *kvProviderPs) Upsert(c *Collection, id string, val interface{}, opts *UpsertOptions) (*MutationResult, error) {
	opm := newKvOpManagerPs(c, "upsert", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var expiry *kv_v1.UpsertRequest_ExpirySecs
	if opts.Expiry > 0 {
		expiry = &kv_v1.UpsertRequest_ExpirySecs{ExpirySecs: uint32(opts.Expiry.Seconds())}
	}

	var preserveExpiry *bool
	if opts.PreserveExpiry {
		preserveExpiry = &opts.PreserveExpiry
	} else {
		if opts.Expiry == 0 {
			expiry = &kv_v1.UpsertRequest_ExpirySecs{ExpirySecs: 0}
		}
	}

	request := &kv_v1.UpsertRequest{
		Key:            opm.DocumentID(),
		BucketName:     opm.BucketName(),
		ScopeName:      opm.ScopeName(),
		CollectionName: opm.CollectionName(),
		ContentFlags:   opm.ValueFlags(),

		PreserveExpiryOnExisting: preserveExpiry,
		Expiry:                   expiry,
		DurabilityLevel:          opm.DurabilityLevel(),
	}

	value, isCompressed := opm.Value()
	if isCompressed {
		request.Content = &kv_v1.UpsertRequest_ContentCompressed{ContentCompressed: value}
	} else {
		request.Content = &kv_v1.UpsertRequest_ContentUncompressed{ContentUncompressed: value}
	}

	res, err := wrapPSOp(opm, request, p.client.Upsert)
	if err != nil {
		return nil, err
	}

	mt := psMutToGoCbMut(res.MutationToken)
	cas := res.Cas

	mutOut := MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(cas),
		},
	}
	return &mutOut, nil
}

func (p *kvProviderPs) Replace(c *Collection, id string, val interface{}, opts *ReplaceOptions) (*MutationResult, error) {
	opm := newKvOpManagerPs(c, "replace", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var cas *uint64
	if opts.Cas > 0 {
		cas = (*uint64)(&opts.Cas)
	}

	// CNG interprets a nil expiry value as "preserve expiry"
	var expiry *kv_v1.ReplaceRequest_ExpirySecs
	if !opts.PreserveExpiry {
		if opts.Expiry > 0 {
			expiry = &kv_v1.ReplaceRequest_ExpirySecs{ExpirySecs: uint32(opts.Expiry.Seconds())}
		} else {
			expiry = &kv_v1.ReplaceRequest_ExpirySecs{ExpirySecs: 0}
		}
	}

	request := &kv_v1.ReplaceRequest{
		Key:          opm.DocumentID(),
		ContentFlags: opm.ValueFlags(),

		Cas:            cas,
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		BucketName:     opm.BucketName(),

		Expiry:          expiry,
		DurabilityLevel: opm.DurabilityLevel(),
	}

	value, isCompressed := opm.Value()
	if isCompressed {
		request.Content = &kv_v1.ReplaceRequest_ContentCompressed{ContentCompressed: value}
	} else {
		request.Content = &kv_v1.ReplaceRequest_ContentUncompressed{ContentUncompressed: value}
	}

	res, err := wrapPSOp(opm, request, p.client.Replace)
	if err != nil {
		return nil, err
	}

	mt := psMutToGoCbMut(res.MutationToken)
	outCas := res.Cas

	mutOut := MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(outCas),
		},
	}
	return &mutOut, nil
}

func (p *kvProviderPs) Get(c *Collection, id string, opts *GetOptions) (*GetResult, error) {
	opm := newKvOpManagerPs(c, "get", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)
	opm.SetIsIdempotent(true)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	request := &kv_v1.GetRequest{
		Key: opm.DocumentID(),

		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		BucketName:     opm.BucketName(),
		Project:        opts.Project,

		Compression: opm.CompressionEnabled(),
	}

	res, err := wrapPSOp(opm, request, p.client.Get)
	if err != nil {
		return nil, err
	}

	content, err := c.compressor.Decompress(res)
	if err != nil {
		return nil, err
	}

	resOut := GetResult{
		Result: Result{Cas(res.Cas)},

		contents: content,
		flags:    res.ContentFlags,

		transcoder: opm.Transcoder(),
	}
	if res.Expiry != nil {
		t := res.Expiry.AsTime()
		resOut.expiryTime = &t
	}

	return &resOut, nil

}

func (p *kvProviderPs) GetAndTouch(c *Collection, id string, expiry time.Duration, opts *GetAndTouchOptions) (*GetResult, error) {
	opm := newKvOpManagerPs(c, "get_and_touch", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	reqExpiry := &kv_v1.GetAndTouchRequest_ExpirySecs{ExpirySecs: uint32(expiry.Seconds())}

	request := &kv_v1.GetAndTouchRequest{
		Key: opm.DocumentID(),

		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		BucketName:     opm.BucketName(),

		Expiry:      reqExpiry,
		Compression: opm.CompressionEnabled(),
	}

	res, err := wrapPSOp(opm, request, p.client.GetAndTouch)
	if err != nil {
		return nil, err
	}

	content, err := c.compressor.Decompress(res)
	if err != nil {
		return nil, err
	}

	resOut := GetResult{
		Result:     Result{Cas(res.Cas)},
		transcoder: opm.Transcoder(),

		contents: content,
		flags:    res.ContentFlags,
	}
	if res.Expiry != nil {
		t := res.Expiry.AsTime()
		resOut.expiryTime = &t
	}

	return &resOut, nil
}

func (p *kvProviderPs) GetAndLock(c *Collection, id string, lockTime time.Duration, opts *GetAndLockOptions) (*GetResult, error) {
	opm := newKvOpManagerPs(c, "get_and_lock", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	request := &kv_v1.GetAndLockRequest{
		Key: opm.DocumentID(),

		BucketName:     opm.BucketName(),
		ScopeName:      opm.ScopeName(),
		CollectionName: opm.CollectionName(),
		LockTime:       uint32(lockTime.Seconds()),

		Compression: opm.CompressionEnabled(),
	}

	res, err := wrapPSOp(opm, request, p.client.GetAndLock)
	if err != nil {
		return nil, err
	}

	content, err := c.compressor.Decompress(res)
	if err != nil {
		return nil, err
	}

	resOut := GetResult{
		Result:     Result{Cas(res.Cas)},
		transcoder: opm.Transcoder(),

		contents: content,
		flags:    res.ContentFlags,
	}
	if res.Expiry != nil {
		t := res.Expiry.AsTime()
		resOut.expiryTime = &t
	}

	return &resOut, nil
}

func (p *kvProviderPs) Exists(c *Collection, id string, opts *ExistsOptions) (*ExistsResult, error) {
	opm := newKvOpManagerPs(c, "exists", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTimeout(opts.Timeout)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetContext(opts.Context)
	opm.SetIsIdempotent(true)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	request := &kv_v1.ExistsRequest{
		Key: opm.DocumentID(),

		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		BucketName:     opm.BucketName(),
	}

	res, err := wrapPSOp(opm, request, p.client.Exists)
	if err != nil {
		return nil, err
	}

	resOut := ExistsResult{
		Result: Result{
			Cas(res.Cas),
		},
		docExists: res.Result,
	}

	return &resOut, nil
}

func (p *kvProviderPs) Remove(c *Collection, id string, opts *RemoveOptions) (*MutationResult, error) {
	opm := newKvOpManagerPs(c, "remove", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var cas *uint64
	if opts.Cas > 0 {
		cas = (*uint64)(&opts.Cas)
	}

	request := &kv_v1.RemoveRequest{
		BucketName:      opm.BucketName(),
		ScopeName:       opm.ScopeName(),
		CollectionName:  opm.CollectionName(),
		Key:             opm.DocumentID(),
		Cas:             cas,
		DurabilityLevel: opm.DurabilityLevel(),
	}

	res, err := wrapPSOp(opm, request, p.client.Remove)
	if err != nil {
		return nil, err
	}

	mt := psMutToGoCbMut(res.MutationToken)
	outCas := res.Cas

	mutOut := MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(outCas),
		},
	}

	return &mutOut, nil
}

func (p *kvProviderPs) Unlock(c *Collection, id string, cas Cas, opts *UnlockOptions) error {
	opm := newKvOpManagerPs(c, "unlock", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return err
	}

	request := &kv_v1.UnlockRequest{
		BucketName:     opm.BucketName(),
		ScopeName:      opm.ScopeName(),
		CollectionName: opm.CollectionName(),
		Key:            opm.DocumentID(),
		Cas:            (uint64)(cas),
	}

	_, err := wrapPSOp(opm, request, p.client.Unlock)
	if err != nil {
		return err
	}

	return nil
}

func (p *kvProviderPs) Touch(c *Collection, id string, expiry time.Duration, opts *TouchOptions) (*MutationResult, error) {
	opm := newKvOpManagerPs(c, "touch", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	request := &kv_v1.TouchRequest{
		BucketName:     opm.BucketName(),
		ScopeName:      opm.ScopeName(),
		CollectionName: opm.CollectionName(),
		Key:            opm.DocumentID(),
		Expiry:         &kv_v1.TouchRequest_ExpirySecs{ExpirySecs: uint32(expiry.Seconds())},
	}

	res, err := wrapPSOp(opm, request, p.client.Touch)
	if err != nil {
		return nil, err
	}

	mt := psMutToGoCbMut(res.MutationToken)
	outCas := res.Cas

	mutOut := MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(outCas),
		},
	}

	return &mutOut, nil
}

type psReplicasResult struct {
	cli        kv_v1.KvService_GetAllReplicasClient
	cancelFunc context.CancelFunc
	err        error

	transcoder Transcoder
}

func (r *psReplicasResult) Next() interface{} {
	res, err := r.cli.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			r.finishWithoutError()
			return nil
		}
		r.finishWithError(err)
		return nil
	}

	return &GetReplicaResult{
		GetResult: GetResult{
			Result: Result{
				cas: Cas(res.Cas),
			},
			transcoder: r.transcoder,
			flags:      res.ContentFlags,
			contents:   res.Content,
		},
		isReplica: res.IsReplica,
	}
}

func (r *psReplicasResult) Close() error {
	if r.err != nil {
		return r.err
	}
	// if the client is nil then we must be closed already.
	if r.cli == nil {
		return nil
	}
	r.cancelFunc()
	err := r.cli.CloseSend()
	r.cli = nil
	return err
}

func (r *psReplicasResult) finishWithoutError() {
	r.cancelFunc()
	// Close the stream now that we are done with it
	err := r.cli.CloseSend()
	if err != nil {
		logWarnf("replicas stream close failed after results: %s", err)
	}

	r.cli = nil
}

func (r *psReplicasResult) finishWithError(err error) {
	// Lets record the error that happened
	r.err = err
	r.cancelFunc()

	// Lets Close the underlying stream
	closeErr := r.cli.CloseSend()
	if closeErr != nil {
		// We log this at debug level, but its almost always going to be an
		// error since thats the most likely reason we are in finishWithError
		logDebugf("replicas stream close failed after error: %s", closeErr)
	}

	// Our client is invalidated as soon as an error occurs
	r.cli = nil
}

func (p *kvProviderPs) GetAllReplicas(c *Collection, id string, opts *GetAllReplicaOptions) (*GetAllReplicasResult, error) {
	opm := newKvOpManagerPs(c, "get_all_replicas", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTimeout(opts.Timeout)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	request := &kv_v1.GetAllReplicasRequest{
		BucketName:     opm.BucketName(),
		ScopeName:      opm.ScopeName(),
		CollectionName: opm.CollectionName(),
		Key:            opm.DocumentID(),
	}

	ctx, cancel := p.newOpCtx(opts.Context, opm.getTimeout())

	res, err := wrapPSOpCtx(ctx, opm, request, p.client.GetAllReplicas)
	if err != nil {
		return nil, err
	}

	return &GetAllReplicasResult{
		res: &psReplicasResult{
			cli:        res,
			cancelFunc: cancel,

			transcoder: opm.Transcoder(),
		},
	}, nil
}

func (p *kvProviderPs) GetAnyReplica(c *Collection, id string, opts *GetAnyReplicaOptions) (*GetReplicaResult, error) {
	opm := newKvOpManagerPs(c, "get_any_replica", opts.ParentSpan, p)
	defer opm.Finish(false)

	res, err := p.GetAllReplicas(c, id, &GetAllReplicaOptions{
		Transcoder:    opts.Transcoder,
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		ParentSpan:    opm.TraceSpan(),
		Context:       opts.Context,
		Internal:      opts.Internal,
		noMetrics:     false,
	})
	if err != nil {
		return nil, opm.EnhanceErr(err, false)
	}

	recv := res.Next()
	if recv == nil {
		return nil, opm.EnhanceErr(ErrDocumentUnretrievable, true)
	}

	return recv, nil
}

func (p *kvProviderPs) Prepend(c *Collection, id string, val []byte, opts *PrependOptions) (*MutationResult, error) {
	opm := newKvOpManagerPs(c, "prepend", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)
	opm.SetValue(val)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var cas *uint64
	if opts.Cas > 0 {
		cas = (*uint64)(&opts.Cas)
	}

	request := &kv_v1.PrependRequest{
		BucketName:      opm.BucketName(),
		ScopeName:       opm.ScopeName(),
		CollectionName:  opm.CollectionName(),
		Key:             opm.DocumentID(),
		Content:         opm.ValueBytes(),
		Cas:             cas,
		DurabilityLevel: opm.DurabilityLevel(),
	}

	res, err := wrapPSOp(opm, request, p.client.Prepend)
	if err != nil {
		return nil, err
	}

	mt := psMutToGoCbMut(res.MutationToken)
	outCas := res.Cas
	mutOut := &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(outCas),
		},
	}

	return mutOut, nil
}

func (p *kvProviderPs) Append(c *Collection, id string, val []byte, opts *AppendOptions) (*MutationResult, error) {
	opm := newKvOpManagerPs(c, "append", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetContext(opts.Context)
	opm.SetValue(val)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var cas *uint64
	if opts.Cas > 0 {
		cas = (*uint64)(&opts.Cas)
	}

	request := &kv_v1.AppendRequest{
		BucketName:      opm.BucketName(),
		ScopeName:       opm.ScopeName(),
		CollectionName:  opm.CollectionName(),
		Key:             opm.DocumentID(),
		Content:         opm.ValueBytes(),
		Cas:             cas,
		DurabilityLevel: opm.DurabilityLevel(),
	}

	res, err := wrapPSOp(opm, request, p.client.Append)
	if err != nil {
		return nil, err
	}

	mt := psMutToGoCbMut(res.MutationToken)
	outCas := res.Cas
	mutOut := &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(outCas),
		},
	}

	return mutOut, nil
}

func (p *kvProviderPs) Increment(c *Collection, id string, opts *IncrementOptions) (*CounterResult, error) {
	if opts.Cas > 0 {
		return nil, makeInvalidArgumentsError("cas is not supported for the increment operation")
	}

	opm := newKvOpManagerPs(c, "increment", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetTimeout(opts.Timeout)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetContext(opts.Context)

	var expiry *kv_v1.IncrementRequest_ExpirySecs
	if opts.Expiry > 0 {
		expiry = &kv_v1.IncrementRequest_ExpirySecs{ExpirySecs: uint32(opts.Expiry.Seconds())}
	}

	request := &kv_v1.IncrementRequest{
		BucketName:      opm.BucketName(),
		ScopeName:       opm.ScopeName(),
		CollectionName:  opm.CollectionName(),
		Key:             opm.DocumentID(),
		Delta:           opts.Delta,
		Expiry:          expiry,
		Initial:         &opts.Initial,
		DurabilityLevel: opm.DurabilityLevel(),
	}

	res, err := wrapPSOp(opm, request, p.client.Increment)
	if err != nil {
		return nil, err
	}

	countOut := &CounterResult{}
	countOut.cas = Cas(res.Cas)
	countOut.mt = psMutToGoCbMut(res.MutationToken)
	countOut.content = uint64(res.Content)

	return countOut, nil

}
func (p *kvProviderPs) Decrement(c *Collection, id string, opts *DecrementOptions) (*CounterResult, error) {
	if opts.Cas > 0 {
		return nil, makeInvalidArgumentsError("cas is not supported for the decrement operation")
	}

	opm := newKvOpManagerPs(c, "decrement", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.DurabilityLevel)
	opm.SetTimeout(opts.Timeout)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetContext(opts.Context)

	var expiry *kv_v1.DecrementRequest_ExpirySecs
	if opts.Expiry > 0 {
		expiry = &kv_v1.DecrementRequest_ExpirySecs{ExpirySecs: uint32(opts.Expiry.Seconds())}
	}

	request := &kv_v1.DecrementRequest{
		BucketName:      opm.BucketName(),
		ScopeName:       opm.ScopeName(),
		CollectionName:  opm.CollectionName(),
		Key:             opm.DocumentID(),
		Delta:           opts.Delta,
		Expiry:          expiry,
		Initial:         &opts.Initial,
		DurabilityLevel: opm.DurabilityLevel(),
	}

	res, err := wrapPSOp(opm, request, p.client.Decrement)
	if err != nil {
		return nil, err
	}

	countOut := &CounterResult{}
	countOut.cas = Cas(res.Cas)
	countOut.mt = psMutToGoCbMut(res.MutationToken)
	countOut.content = uint64(res.Content)

	return countOut, nil

}

func (p *kvProviderPs) StartKvOpTrace(c *Collection, operationName string, tracectx RequestSpanContext, noAttributes bool) RequestSpan {
	return c.startKvOpTrace(operationName, tracectx, p.tracer, noAttributes)
}

func psMutToGoCbMut(in *kv_v1.MutationToken) *MutationToken {
	if in != nil {
		return &MutationToken{
			bucketName: in.BucketName,
			token: gocbcore.MutationToken{
				VbID:   uint16(in.VbucketId),
				VbUUID: gocbcore.VbUUID(in.VbucketUuid),
				SeqNo:  gocbcore.SeqNo(in.SeqNo),
			},
		}
	}

	return nil
}

func (p *kvProviderPs) newOpCtx(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithTimeout(ctx, timeout)
}
