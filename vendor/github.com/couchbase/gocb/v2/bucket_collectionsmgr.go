package gocb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/couchbase/gocbcore/v9"
)

// CollectionSpec describes the specification of a collection.
type CollectionSpec struct {
	Name      string
	ScopeName string
	MaxExpiry time.Duration
}

// ScopeSpec describes the specification of a scope.
type ScopeSpec struct {
	Name        string
	Collections []CollectionSpec
}

// These 3 types are temporary. They are necessary for now as the server beta was released with ns_server returning
// a different jsonManifest format to what it will return in the future.
type jsonManifest struct {
	UID    uint64                       `json:"uid"`
	Scopes map[string]jsonManifestScope `json:"scopes"`
}

type jsonManifestScope struct {
	UID         uint32                            `json:"uid"`
	Collections map[string]jsonManifestCollection `json:"collections"`
}

type jsonManifestCollection struct {
	UID uint32 `json:"uid"`
}

// CollectionManager provides methods for performing collections management.
type CollectionManager struct {
	mgmtProvider mgmtProvider
	bucketName   string
	tracer       requestTracer
}

func (cm *CollectionManager) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logDebugf("failed to read http body: %s", err)
		return nil
	}

	errText := strings.ToLower(string(b))

	if resp.StatusCode == 404 {
		if strings.Contains(errText, "not found") && strings.Contains(errText, "scope") {
			return makeGenericMgmtError(ErrScopeNotFound, req, resp)
		} else if strings.Contains(errText, "not found") && strings.Contains(errText, "scope") {
			return makeGenericMgmtError(ErrScopeNotFound, req, resp)
		}
	}

	if strings.Contains(errText, "already exists") && strings.Contains(errText, "collection") {
		return makeGenericMgmtError(ErrCollectionExists, req, resp)
	} else if strings.Contains(errText, "already exists") && strings.Contains(errText, "scope") {
		return makeGenericMgmtError(ErrScopeExists, req, resp)
	}

	return makeGenericMgmtError(errors.New(errText), req, resp)
}

// GetAllScopesOptions is the set of options available to the GetAllScopes operation.
type GetAllScopesOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetAllScopes gets all scopes from the bucket.
func (cm *CollectionManager) GetAllScopes(opts *GetAllScopesOptions) ([]ScopeSpec, error) {
	if opts == nil {
		opts = &GetAllScopesOptions{}
	}

	span := cm.tracer.StartSpan("GetAllScopes", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s/collections", cm.bucketName),
		Method:        "GET",
		RetryStrategy: opts.RetryStrategy,
		IsIdempotent:  true,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := cm.mgmtProvider.executeMgmtRequest(req)
	if err != nil {
		colErr := cm.tryParseErrorMessage(&req, resp)
		if colErr != nil {
			return nil, colErr
		}
		return nil, makeMgmtBadStatusError("failed to get all scopes", &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		return nil, makeMgmtBadStatusError("failed to get all scopes", &req, resp)
	}

	var scopes []ScopeSpec
	var mfest gocbcore.Manifest
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&mfest)
	if err == nil {
		for _, scope := range mfest.Scopes {
			var collections []CollectionSpec
			for _, col := range scope.Collections {
				collections = append(collections, CollectionSpec{
					Name:      col.Name,
					ScopeName: scope.Name,
				})
			}
			scopes = append(scopes, ScopeSpec{
				Name:        scope.Name,
				Collections: collections,
			})
		}
	} else {
		// Temporary support for older server version
		var oldMfest jsonManifest
		jsonDec := json.NewDecoder(resp.Body)
		err = jsonDec.Decode(&oldMfest)
		if err != nil {
			return nil, err
		}

		for scopeName, scope := range oldMfest.Scopes {
			var collections []CollectionSpec
			for colName := range scope.Collections {
				collections = append(collections, CollectionSpec{
					Name:      colName,
					ScopeName: scopeName,
				})
			}
			scopes = append(scopes, ScopeSpec{
				Name:        scopeName,
				Collections: collections,
			})
		}
	}

	return scopes, nil
}

// CreateCollectionOptions is the set of options available to the CreateCollection operation.
type CreateCollectionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// CreateCollection creates a new collection on the bucket.
func (cm *CollectionManager) CreateCollection(spec CollectionSpec, opts *CreateCollectionOptions) error {
	if spec.Name == "" {
		return makeInvalidArgumentsError("collection name cannot be empty")
	}

	if spec.ScopeName == "" {
		return makeInvalidArgumentsError("scope name cannot be empty")
	}

	if opts == nil {
		opts = &CreateCollectionOptions{}
	}

	span := cm.tracer.StartSpan("CreateCollection", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	posts := url.Values{}
	posts.Add("name", spec.Name)

	if spec.MaxExpiry > 0 {
		posts.Add("maxTTL", fmt.Sprintf("%d", int(spec.MaxExpiry.Seconds())))
	}

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s/collections/%s", cm.bucketName, spec.ScopeName),
		Method:        "POST",
		Body:          []byte(posts.Encode()),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := cm.mgmtProvider.executeMgmtRequest(req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		colErr := cm.tryParseErrorMessage(&req, resp)
		if colErr != nil {
			return colErr
		}
		return makeMgmtBadStatusError("failed to create collection", &req, resp)
	}

	err = resp.Body.Close()
	if err != nil {
		logDebugf("Failed to close socket (%s)", err)
	}

	return nil
}

// DropCollectionOptions is the set of options available to the DropCollection operation.
type DropCollectionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// DropCollection removes a collection.
func (cm *CollectionManager) DropCollection(spec CollectionSpec, opts *DropCollectionOptions) error {
	if spec.Name == "" {
		return makeInvalidArgumentsError("collection name cannot be empty")
	}

	if spec.ScopeName == "" {
		return makeInvalidArgumentsError("scope name cannot be empty")
	}

	if opts == nil {
		opts = &DropCollectionOptions{}
	}

	span := cm.tracer.StartSpan("DropCollection", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s/collections/%s/%s", cm.bucketName, spec.ScopeName, spec.Name),
		Method:        "DELETE",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := cm.mgmtProvider.executeMgmtRequest(req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		colErr := cm.tryParseErrorMessage(&req, resp)
		if colErr != nil {
			return colErr
		}
		return makeMgmtBadStatusError("failed to drop collection", &req, resp)
	}

	err = resp.Body.Close()
	if err != nil {
		logDebugf("Failed to close socket (%s)", err)
	}

	return nil
}

// CreateScopeOptions is the set of options available to the CreateScope operation.
type CreateScopeOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// CreateScope creates a new scope on the bucket.
func (cm *CollectionManager) CreateScope(scopeName string, opts *CreateScopeOptions) error {
	if scopeName == "" {
		return makeInvalidArgumentsError("scope name cannot be empty")
	}

	if opts == nil {
		opts = &CreateScopeOptions{}
	}

	span := cm.tracer.StartSpan("CreateScope", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	posts := url.Values{}
	posts.Add("name", scopeName)

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s/collections", cm.bucketName),
		Method:        "POST",
		Body:          []byte(posts.Encode()),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := cm.mgmtProvider.executeMgmtRequest(req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		colErr := cm.tryParseErrorMessage(&req, resp)
		if colErr != nil {
			return colErr
		}
		return makeMgmtBadStatusError("failed to create scope", &req, resp)
	}

	err = resp.Body.Close()
	if err != nil {
		logDebugf("Failed to close socket (%s)", err)
	}

	return nil
}

// DropScopeOptions is the set of options available to the DropScope operation.
type DropScopeOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// DropScope removes a scope.
func (cm *CollectionManager) DropScope(scopeName string, opts *DropScopeOptions) error {
	if opts == nil {
		opts = &DropScopeOptions{}
	}

	span := cm.tracer.StartSpan("DropScope", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s/collections/%s", cm.bucketName, scopeName),
		Method:        "DELETE",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := cm.mgmtProvider.executeMgmtRequest(req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		colErr := cm.tryParseErrorMessage(&req, resp)
		if colErr != nil {
			return colErr
		}
		return makeMgmtBadStatusError("failed to drop scope", &req, resp)
	}

	err = resp.Body.Close()
	if err != nil {
		logDebugf("Failed to close socket (%s)", err)
	}

	return nil
}
