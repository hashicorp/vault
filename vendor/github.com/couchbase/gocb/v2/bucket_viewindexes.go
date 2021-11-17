package gocb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// DesignDocumentNamespace represents which namespace a design document resides in.
type DesignDocumentNamespace uint

const (
	// DesignDocumentNamespaceProduction means that a design document resides in the production namespace.
	DesignDocumentNamespaceProduction DesignDocumentNamespace = iota

	// DesignDocumentNamespaceDevelopment means that a design document resides in the development namespace.
	DesignDocumentNamespaceDevelopment
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

// View represents a Couchbase view within a design document.
type View struct {
	Map    string
	Reduce string
}

func (v *View) fromData(data jsonView) error {
	v.Map = data.Map
	v.Reduce = data.Reduce

	return nil
}

func (v *View) toData() (jsonView, error) {
	var data jsonView

	data.Map = v.Map
	data.Reduce = v.Reduce

	return data, nil
}

// DesignDocument represents a Couchbase design document containing multiple views.
type DesignDocument struct {
	Name  string
	Views map[string]View
}

func (dd *DesignDocument) fromData(data jsonDesignDocument, name string) error {
	dd.Name = name

	views := make(map[string]View)
	for viewName, viewData := range data.Views {
		var view View
		err := view.fromData(viewData)
		if err != nil {
			return err
		}

		views[viewName] = view
	}
	dd.Views = views

	return nil
}

func (dd *DesignDocument) toData() (jsonDesignDocument, string, error) {
	var data jsonDesignDocument

	views := make(map[string]jsonView)
	for viewName, view := range dd.Views {
		viewData, err := view.toData()
		if err != nil {
			return jsonDesignDocument{}, "", err
		}

		views[viewName] = viewData
	}
	data.Views = views

	return data, dd.Name, nil
}

// ViewIndexManager provides methods for performing View management.
type ViewIndexManager struct {
	mgmtProvider mgmtProvider
	bucketName   string

	tracer RequestTracer
	meter  *meterWrapper
}

func (vm *ViewIndexManager) tryParseErrorMessage(req mgmtRequest, resp *mgmtResponse) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read view index manager response body: %s", err)
		return nil
	}

	if resp.StatusCode == 404 {
		if strings.Contains(strings.ToLower(string(b)), "not_found") {
			return makeGenericMgmtError(ErrDesignDocumentNotFound, &req, resp)
		}

		return makeGenericMgmtError(errors.New(string(b)), &req, resp)
	}

	var mgrErr bucketMgrErrorResp
	err = json.Unmarshal(b, &mgrErr)
	if err != nil {
		logDebugf("Failed to unmarshal error body: %s", err)
		return makeGenericMgmtError(errors.New(string(b)), &req, resp)
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

	return makeGenericMgmtError(bodyErr, &req, resp)
}

func (vm *ViewIndexManager) doMgmtRequest(ctx context.Context, req mgmtRequest) (*mgmtResponse, error) {
	resp, err := vm.mgmtProvider.executeMgmtRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetDesignDocumentOptions is the set of options available to the ViewIndexManager GetDesignDocument operation.
type GetDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

func (vm *ViewIndexManager) ddocName(name string, namespace DesignDocumentNamespace) string {
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
func (vm *ViewIndexManager) GetDesignDocument(name string, namespace DesignDocumentNamespace, opts *GetDesignDocumentOptions) (*DesignDocument, error) {
	if opts == nil {
		opts = &GetDesignDocumentOptions{}
	}

	start := time.Now()
	defer vm.meter.ValueRecord(meterValueServiceManagement, "manager_views_get_design_document", start)

	return vm.getDesignDocument(name, namespace, time.Now(), opts)
}

func (vm *ViewIndexManager) getDesignDocument(name string, namespace DesignDocumentNamespace,
	startTime time.Time, opts *GetDesignDocumentOptions) (*DesignDocument, error) {

	name = vm.ddocName(name, namespace)

	span := createSpan(vm.tracer, opts.ParentSpan, "manager_views_get_design_document", "management")
	span.SetAttribute("db.operation", "GET "+fmt.Sprintf("/_design/%s", name))
	span.SetAttribute("db.name", vm.bucketName)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeViews,
		Path:          fmt.Sprintf("/_design/%s", name),
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

		return nil, makeGenericMgmtError(errors.New("failed to get design document"), &req, resp)
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

// GetAllDesignDocumentsOptions is the set of options available to the ViewIndexManager GetAllDesignDocuments operation.
type GetAllDesignDocumentsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllDesignDocuments will retrieve all design documents for the given bucket.
func (vm *ViewIndexManager) GetAllDesignDocuments(namespace DesignDocumentNamespace, opts *GetAllDesignDocumentsOptions) ([]DesignDocument, error) {
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

		return nil, makeGenericMgmtError(errors.New("failed to get design documents"), &req, resp)
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

// UpsertDesignDocumentOptions is the set of options available to the ViewIndexManager UpsertDesignDocument operation.
type UpsertDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UpsertDesignDocument will insert a design document to the given bucket, or update
// an existing design document with the same name.
func (vm *ViewIndexManager) UpsertDesignDocument(ddoc DesignDocument, namespace DesignDocumentNamespace, opts *UpsertDesignDocumentOptions) error {
	if opts == nil {
		opts = &UpsertDesignDocumentOptions{}
	}

	start := time.Now()
	defer vm.meter.ValueRecord(meterValueServiceManagement, "manager_views_upsert_design_document", start)

	return vm.upsertDesignDocument(ddoc, namespace, time.Now(), opts)
}

func (vm *ViewIndexManager) upsertDesignDocument(
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
		Path:          fmt.Sprintf("/_design/%s", ddocName),
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

		return makeGenericMgmtError(errors.New("failed to upsert design document"), &req, resp)
	}

	return nil
}

// DropDesignDocumentOptions is the set of options available to the ViewIndexManager Upsert operation.
type DropDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropDesignDocument will remove a design document from the given bucket.
func (vm *ViewIndexManager) DropDesignDocument(name string, namespace DesignDocumentNamespace, opts *DropDesignDocumentOptions) error {
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

func (vm *ViewIndexManager) dropDesignDocument(tracectx RequestSpanContext, name string, namespace DesignDocumentNamespace,
	startTime time.Time, opts *DropDesignDocumentOptions) error {

	name = vm.ddocName(name, namespace)

	req := mgmtRequest{
		Service:       ServiceTypeViews,
		Path:          fmt.Sprintf("/_design/%s", name),
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

		return makeGenericMgmtError(errors.New("failed to drop design document"), &req, resp)
	}

	return nil
}

// PublishDesignDocumentOptions is the set of options available to the ViewIndexManager PublishDesignDocument operation.
type PublishDesignDocumentOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// PublishDesignDocument publishes a design document to the given bucket.
func (vm *ViewIndexManager) PublishDesignDocument(name string, opts *PublishDesignDocumentOptions) error {
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
