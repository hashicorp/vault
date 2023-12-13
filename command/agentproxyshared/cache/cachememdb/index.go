// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cachememdb

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Index holds the response to be cached along with multiple other values that
// serve as pointers to refer back to this index.
type Index struct {
	// ID is a value that uniquely represents the request held by this
	// index. This is computed by serializing and hashing the response object.
	// Required: true, Unique: true
	ID string

	// Token is the token that fetched the response held by this index
	// Required: true, Unique: true
	Token string

	// Tokens is a set of tokens that can access this cached response,
	// which is used for static secret caching, and enabling multiple
	// tokens to be able to access the same cache entry for static secrets.
	// Implemented as a map so that all values are unique.
	// Required: false, Unique: false
	Tokens map[string]struct{}

	// TokenParent is the parent token of the token held by this index
	// Required: false, Unique: false
	TokenParent string

	// TokenAccessor is the accessor of the token being cached in this index
	// Required: true, Unique: true
	TokenAccessor string

	// Namespace is the namespace that was provided in the request path as the
	// Vault namespace to query
	Namespace string

	// RequestPath is the path of the request that resulted in the response
	// held by this index.
	// For dynamic secrets, this will be the actual path sent to the request,
	// e.g. /v1/foo/bar (which will not include the namespace if it was included
	// in the headers).
	// For static secrets, this will be the canonical path to the secret (i.e.
	// after calling getStaticSecretPathFromRequest--see its godocs for more
	// information).
	// Required: true, Unique: false
	RequestPath string

	// Lease is the identifier of the lease in Vault, that belongs to the
	// response held by this index.
	// Required: false, Unique: true
	Lease string

	// LeaseToken is the identifier of the token that created the lease held by
	// this index.
	// Required: false, Unique: false
	LeaseToken string

	// Response is the serialized response object that the agent is caching.
	Response []byte

	// RenewCtxInfo holds the context and the corresponding cancel func for the
	// goroutine that manages the renewal of the secret belonging to the
	// response in this index.
	RenewCtxInfo *ContextInfo

	// RequestMethod is the HTTP method of the request
	RequestMethod string

	// RequestToken is the token used in the request
	RequestToken string

	// RequestHeader is the header used in the request
	RequestHeader http.Header

	// LastRenewed is the timestamp of last renewal
	LastRenewed time.Time

	// Type is the index type (token, auth-lease, secret-lease, static-secret)
	Type string

	// IndexLock is a lock held for some indexes to prevent data
	// races upon update.
	IndexLock sync.RWMutex
}

// CapabilitiesIndex holds the capabilities for cached static secrets.
// This type of index does not represent a response.
type CapabilitiesIndex struct {
	// ID is a value that uniquely represents the request held by this
	// index. This is computed by hashing the token that this capabilities
	// index represents the capabilities of.
	// Required: true, Unique: true
	ID string

	// Token is the token that fetched the response held by this index
	// Required: true, Unique: true
	Token string

	// ReadablePaths is a set of paths with read capabilities for the given token.
	// Implemented as a map for uniqueness. The key to the map is a path (such as
	// `foo/bar` that we've demonstrated we can read.
	ReadablePaths map[string]struct{}

	// IndexLock is a lock held for some indexes to prevent data
	// races upon update.
	IndexLock sync.RWMutex
}

type IndexName uint32

const (
	// IndexNameID is the ID of the index constructed from the serialized request.
	IndexNameID = "id"

	// IndexNameLease is the lease of the index.
	IndexNameLease = "lease"

	// IndexNameRequestPath is the request path of the index.
	IndexNameRequestPath = "request_path"

	// IndexNameToken is the token of the index.
	IndexNameToken = "token"

	// IndexNameTokenAccessor is the token accessor of the index.
	IndexNameTokenAccessor = "token_accessor"

	// IndexNameTokenParent is the token parent of the index.
	IndexNameTokenParent = "token_parent"

	// IndexNameLeaseToken is the token that created the lease.
	IndexNameLeaseToken = "lease_token"

	// CapabilitiesIndexNameID is the ID of the capabilities index.
	CapabilitiesIndexNameID = "id"
)

func validIndexName(indexName string) bool {
	switch indexName {
	case IndexNameID:
	case IndexNameLease:
	case IndexNameRequestPath:
	case IndexNameToken:
	case IndexNameTokenAccessor:
	case IndexNameTokenParent:
	case IndexNameLeaseToken:
	default:
		return false
	}
	return true
}

func validCapabilitiesIndexName(indexName string) bool {
	switch indexName {
	case CapabilitiesIndexNameID:
	default:
		return false
	}
	return true
}

type ContextInfo struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	DoneCh     chan struct{}
}

func NewContextInfo(ctx context.Context) *ContextInfo {
	if ctx == nil {
		return nil
	}

	ctxInfo := new(ContextInfo)
	ctxInfo.Ctx, ctxInfo.CancelFunc = context.WithCancel(ctx)
	ctxInfo.DoneCh = make(chan struct{})
	return ctxInfo
}

// Serialize returns a json marshal'ed Index object, without the RenewCtxInfo
func (i Index) Serialize() ([]byte, error) {
	i.RenewCtxInfo = nil

	indexBytes, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return indexBytes, nil
}

// Deserialize converts json bytes to an Index object
// Note: RenewCtxInfo will need to be reconstructed elsewhere.
func Deserialize(indexBytes []byte) (*Index, error) {
	index := new(Index)
	if err := json.Unmarshal(indexBytes, index); err != nil {
		return nil, err
	}
	return index, nil
}

// SerializeCapabilitiesIndex returns a json marshal'ed CapabilitiesIndex object
func (i CapabilitiesIndex) SerializeCapabilitiesIndex() ([]byte, error) {
	indexBytes, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return indexBytes, nil
}

// DeserializeCapabilitiesIndex converts json bytes to an CapabilitiesIndex object
func DeserializeCapabilitiesIndex(indexBytes []byte) (*CapabilitiesIndex, error) {
	index := new(CapabilitiesIndex)
	if err := json.Unmarshal(indexBytes, index); err != nil {
		return nil, err
	}
	return index, nil
}
