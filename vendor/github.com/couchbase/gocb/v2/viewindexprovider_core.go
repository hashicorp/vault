package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/url"
	"strings"
	"time"
)

// View represents a Couchbase view within a design document.
type jsonView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

// DesignDocument represents a Couchbase design document containing multiple views.
type jsonDesignDocument struct {
	Views map[string]jsonView `json:"views,omitempty"`
}

type viewIndexProviderCore struct {
	mgmtProvider mgmtProvider
	bucketName   string
	tracer       RequestTracer
	meter        *meterWrapper
}

func (vm *viewIndexProviderCore) tryParseErrorMessage(req mgmtRequest, resp *mgmtResponse) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read view index manager response body: %s", err)
		return nil
	}

	if resp.StatusCode == 404 {
		if strings.Contains(strings.ToLower(string(b)), "not_found") {
			return makeGenericMgmtError(ErrDesignDocumentNotFound, &req, resp, string(b))
		}

		return makeGenericMgmtError(errors.New(string(b)), &req, resp, string(b))
	}

	if err := checkForRateLimitError(resp.StatusCode, string(b)); err != nil {
		return makeGenericMgmtError(err, &req, resp, string(b))
	}

	var mgrErr bucketMgrErrorResp
	err = json.Unmarshal(b, &mgrErr)
	if err != nil {
		logDebugf("Failed to unmarshal error body: %s", err)
		return makeGenericMgmtError(errors.New(string(b)), &req, resp, string(b))
	}

	var bodyErr error
	var firstErr string
	for _, err := range mgrErr.Errors {
		firstErr = strings.ToLower(err)
		break
	}

	if strings.Contains(firstErr, "bucket with given name already exists") {
		bodyErr = ErrBucketExists
	} else {
		bodyErr = errors.New(firstErr)
	}

	return makeGenericMgmtError(bodyErr, &req, resp, string(b))
}

func (vm *viewIndexProviderCore) doMgmtRequest(ctx context.Context, req mgmtRequest) (*mgmtResponse, error) {
	resp, err := vm.mgmtProvider.executeMgmtRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (vm *viewIndexProviderCore) ddocName(name string, namespace DesignDocumentNamespace) string {
	if namespace == DesignDocumentNamespaceProduction {
		if strings.HasPrefix(name, "dev_") {
			name = strings.TrimLeft(name, "dev_")
		}
	} else {
		if !strings.HasPrefix(name, "dev_") {
			name = "dev_" + name
		}
	}

	return name
}

// GetDesignDocument retrieves a single design document for the given bucket.
func (vm *viewIndexProviderCore) GetDesignDocument(name string, namespace DesignDocumentNamespace, opts *GetDesignDocumentOptions) (*DesignDocument, error) {
	if opts == nil {
		opts = &GetDesignDocumentOptions{}
	}

	start := time.Now()
	defer vm.meter.ValueRecord(meterValueServiceManagement, "manager_views_get_design_document", start)

	return vm.getDesignDocument(name, namespace, time.Now(), opts)
}

func (vm *viewIndexProviderCore) getDesignDocument(name string, namespace DesignDocumentNamespace,
	startTime time.Time, opts *GetDesignDocumentOptions) (*DesignDocument, error) {

	name = vm.ddocName(name, namespace)

	span := createSpan(vm.tracer, opts.ParentSpan, "manager_views_get_design_document", "management")
	span.SetAttribute("db.operation", "GET "+fmt.Sprintf("/_design/%s", name))
	span.SetAttribute("db.name", vm.bucketName)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeViews,
		Path:          fmt.Sprintf("/_design/%s", url.PathEscape(name)),
		Method:        "GET",
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
		UniqueID:      uuid.New().String(),
	}
	resp, err := vm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		vwErr := vm.tryParseErrorMessage(req, resp)
		if vwErr != nil {
			return nil, vwErr
		}

		return nil, makeMgmtBadStatusError("failed to get design document", &req, resp)
	}

	var ddocData jsonDesignDocument
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&ddocData)
	if err != nil {
		return nil, err
	}

	ddocName := strings.TrimPrefix(name, "dev_")

	var ddoc DesignDocument
	err = ddoc.fromData(ddocData, ddocName)
	if err != nil {
		return nil, err
	}

	return &ddoc, nil
}

// GetAllDesignDocuments will retrieve all design documents for the given bucket.
func (vm *viewIndexProviderCore) GetAllDesignDocuments(namespace DesignDocumentNamespace, opts *GetAllDesignDocumentsOptions) ([]DesignDocument, error) {
	if opts == nil {
		opts = &GetAllDesignDocumentsOptions{}
	}

	start := time.Now()
	defer vm.meter.ValueRecord(meterValueServiceManagement, "manager_views_get_all_design_documents", start)

	path := fmt.Sprintf("/pools/default/buckets/%s/ddocs", vm.bucketName)
	span := createSpan(vm.tracer, opts.ParentSpan, "manager_views_get_all_design_documents", "management")
	span.SetAttribute("db.operation", "GET "+path)
	span.SetAttribute("db.name", vm.bucketName)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          path,
		Method:        "GET",
		IsIdempotent:  true,
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		parentSpanCtx: span.Context(),
		UniqueID:      uuid.New().String(),
	}
	resp, err := vm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		vwErr := vm.tryParseErrorMessage(req, resp)
		if vwErr != nil {
			return nil, vwErr
		}

		return nil, makeMgmtBadStatusError("failed to get design documents", &req, resp)
	}

	var ddocsResp struct {
		Rows []struct {
			Doc struct {
				Meta struct {
					ID string `json:"id"`
				}
				JSON jsonDesignDocument `json:"json"`
			} `json:"doc"`
		} `json:"rows"`
	}
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&ddocsResp)
	if err != nil {
		return nil, err
	}

	var ddocs []DesignDocument
	for _, ddocData := range ddocsResp.Rows {
		if len(ddocData.Doc.Meta.ID) <= 8 {
			logErrorf("Design document name was less than 9 characters long: %s", ddocData.Doc.Meta.ID)
			continue
		}
		ddocName := ddocData.Doc.Meta.ID[8:]
		isDevDoc := strings.HasPrefix(ddocName, "dev_")
		switch namespace {
		case DesignDocumentNamespaceProduction:
			if isDevDoc {
				continue
			}
		case DesignDocumentNamespaceDevelopment:
			if !isDevDoc {
				continue
			}

			ddocName = strings.TrimPrefix(ddocName, "dev_")
		default:
			return nil, makeInvalidArgumentsError("design document namespace unknown")
		}

		var ddoc DesignDocument
		err := ddoc.fromData(ddocData.Doc.JSON, ddocName)
		if err != nil {
			return nil, err
		}
		ddocs = append(ddocs, ddoc)
	}

	return ddocs, nil
}

// UpsertDesignDocument will insert a design document to the given bucket, or update
// an existing design document with the same name.
func (vm *viewIndexProviderCore) UpsertDesignDocument(ddoc DesignDocument, namespace DesignDocumentNamespace, opts *UpsertDesignDocumentOptions) error {
	if opts == nil {
		opts = &UpsertDesignDocumentOptions{}
	}

	start := time.Now()
	defer vm.meter.ValueRecord(meterValueServiceManagement, "manager_views_upsert_design_document", start)

	return vm.upsertDesignDocument(ddoc, namespace, time.Now(), opts)
}

func (vm *viewIndexProviderCore) upsertDesignDocument(
	ddoc DesignDocument,
	namespace DesignDocumentNamespace,
	startTime time.Time,
	opts *UpsertDesignDocumentOptions,
) error {
	ddocData, ddocName, err := ddoc.toData()
	if err != nil {
		return err
	}

	ddocName = vm.ddocName(ddocName, namespace)

	span := createSpan(vm.tracer, opts.ParentSpan, "manager_views_upsert_design_document", "management")
	span.SetAttribute("db.operation", "PUT "+fmt.Sprintf("/_design/%s", ddocName))
	span.SetAttribute("db.name", vm.bucketName)
	defer span.End()

	espan := createSpan(vm.tracer, span, "request_encoding", "")
	data, err := json.Marshal(&ddocData)
	espan.End()
	if err != nil {
		return err
	}

	req := mgmtRequest{
		Service:       ServiceTypeViews,
		Path:          fmt.Sprintf("/_design/%s", url.PathEscape(ddocName)),
		Method:        "PUT",
		Body:          data,
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		parentSpanCtx: span.Context(),
		UniqueID:      uuid.New().String(),
	}
	resp, err := vm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 201 {
		vwErr := vm.tryParseErrorMessage(req, resp)
		if vwErr != nil {
			return vwErr
		}

		return makeMgmtBadStatusError("failed to upsert design document", &req, resp)
	}

	return nil
}

// DropDesignDocument will remove a design document from the given bucket.
func (vm *viewIndexProviderCore) DropDesignDocument(name string, namespace DesignDocumentNamespace, opts *DropDesignDocumentOptions) error {
	if opts == nil {
		opts = &DropDesignDocumentOptions{}
	}

	start := time.Now()
	defer vm.meter.ValueRecord(meterValueServiceManagement, "manager_views_drop_design_document", start)

	span := createSpan(vm.tracer, opts.ParentSpan, "manager_views_drop_design_document", "management")
	span.SetAttribute("db.operation", "DELETE "+fmt.Sprintf("/_design/%s", name))
	span.SetAttribute("db.name", vm.bucketName)
	defer span.End()

	return vm.dropDesignDocument(span.Context(), name, namespace, time.Now(), opts)
}

func (vm *viewIndexProviderCore) dropDesignDocument(tracectx RequestSpanContext, name string, namespace DesignDocumentNamespace,
	startTime time.Time, opts *DropDesignDocumentOptions) error {

	name = vm.ddocName(name, namespace)

	req := mgmtRequest{
		Service:       ServiceTypeViews,
		Path:          fmt.Sprintf("/_design/%s", url.PathEscape(name)),
		Method:        "DELETE",
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		parentSpanCtx: tracectx,
		UniqueID:      uuid.New().String(),
	}
	resp, err := vm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		vwErr := vm.tryParseErrorMessage(req, resp)
		if vwErr != nil {
			return vwErr
		}

		return makeMgmtBadStatusError("failed to drop design document", &req, resp)
	}

	return nil
}

// PublishDesignDocument publishes a design document to the given bucket.
func (vm *viewIndexProviderCore) PublishDesignDocument(name string, opts *PublishDesignDocumentOptions) error {
	startTime := time.Now()
	if opts == nil {
		opts = &PublishDesignDocumentOptions{}
	}

	start := time.Now()
	defer vm.meter.ValueRecord(meterValueServiceManagement, "manager_views_publish_design_document", start)

	span := createSpan(vm.tracer, opts.ParentSpan, "manager_views_publish_design_document", "management")
	span.SetAttribute("db.name", vm.bucketName)
	defer span.End()

	devdoc, err := vm.getDesignDocument(
		name,
		DesignDocumentNamespaceDevelopment,
		startTime,
		&GetDesignDocumentOptions{
			RetryStrategy: opts.RetryStrategy,
			Timeout:       opts.Timeout,
			ParentSpan:    span,
			Context:       opts.Context,
		})
	if err != nil {
		return err
	}

	err = vm.upsertDesignDocument(
		*devdoc,
		DesignDocumentNamespaceProduction,
		startTime,
		&UpsertDesignDocumentOptions{
			RetryStrategy: opts.RetryStrategy,
			Timeout:       opts.Timeout,
			ParentSpan:    span,
			Context:       opts.Context,
		})
	if err != nil {
		return err
	}

	return nil
}
