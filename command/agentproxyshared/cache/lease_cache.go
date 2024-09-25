// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cacheboltdb"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/helper/namespace"
	nshelper "github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	gocache "github.com/patrickmn/go-cache"
	"go.uber.org/atomic"
)

const (
	vaultPathTokenCreate         = "/v1/auth/token/create"
	vaultPathTokenRevoke         = "/v1/auth/token/revoke"
	vaultPathTokenRevokeSelf     = "/v1/auth/token/revoke-self"
	vaultPathTokenRevokeAccessor = "/v1/auth/token/revoke-accessor"
	vaultPathTokenRevokeOrphan   = "/v1/auth/token/revoke-orphan"
	vaultPathTokenLookup         = "/v1/auth/token/lookup"
	vaultPathTokenLookupSelf     = "/v1/auth/token/lookup-self"
	vaultPathTokenRenew          = "/v1/auth/token/renew"
	vaultPathTokenRenewSelf      = "/v1/auth/token/renew-self"
	vaultPathLeaseRevoke         = "/v1/sys/leases/revoke"
	vaultPathLeaseRevokeForce    = "/v1/sys/leases/revoke-force"
	vaultPathLeaseRevokePrefix   = "/v1/sys/leases/revoke-prefix"
)

var (
	contextIndexID  = contextIndex{}
	errInvalidType  = errors.New("invalid type provided")
	revocationPaths = []string{
		strings.TrimPrefix(vaultPathTokenRevoke, "/v1"),
		strings.TrimPrefix(vaultPathTokenRevokeSelf, "/v1"),
		strings.TrimPrefix(vaultPathTokenRevokeAccessor, "/v1"),
		strings.TrimPrefix(vaultPathTokenRevokeOrphan, "/v1"),
		strings.TrimPrefix(vaultPathLeaseRevoke, "/v1"),
		strings.TrimPrefix(vaultPathLeaseRevokeForce, "/v1"),
		strings.TrimPrefix(vaultPathLeaseRevokePrefix, "/v1"),
	}
)

type contextIndex struct{}

type cacheClearRequest struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Namespace string `json:"namespace"`
}

// LeaseCache is an implementation of Proxier that handles
// the caching of responses. It passes the incoming request
// to an underlying Proxier implementation.
type LeaseCache struct {
	client      *api.Client
	proxier     Proxier
	logger      hclog.Logger
	db          *cachememdb.CacheMemDB
	baseCtxInfo *cachememdb.ContextInfo
	l           *sync.RWMutex

	// userAgentToUse is the user agent to use when making independent requests
	// to Vault.
	userAgentToUse string

	// idLocks is used during cache lookup to ensure that identical requests made
	// in parallel won't trigger multiple renewal goroutines.
	idLocks []*locksutil.LockEntry

	// inflightCache keeps track of inflight requests
	inflightCache *gocache.Cache

	// ps is the persistent storage for tokens and leases
	ps *cacheboltdb.BoltStorage

	// shuttingDown is used to determine if cache needs to be evicted or not
	// when the context is cancelled
	shuttingDown atomic.Bool

	// cacheStaticSecrets is used to determine if the cache should also
	// cache static secrets, as well as dynamic secrets.
	cacheStaticSecrets bool

	// cacheDynamicSecrets is used to determine if the cache should
	// cache dynamic secrets
	cacheDynamicSecrets bool

	// capabilityManager is used when static secrets are enabled to
	// manage the capabilities of cached tokens.
	capabilityManager *StaticSecretCapabilityManager
}

// LeaseCacheConfig is the configuration for initializing a new
// LeaseCache.
type LeaseCacheConfig struct {
	Client              *api.Client
	BaseContext         context.Context
	Proxier             Proxier
	Logger              hclog.Logger
	UserAgentToUse      string
	Storage             *cacheboltdb.BoltStorage
	CacheStaticSecrets  bool
	CacheDynamicSecrets bool
}

type inflightRequest struct {
	// ch is closed by the request that ends up processing the set of
	// parallel request
	ch chan struct{}

	// remaining is the number of remaining inflight request that needs to
	// be processed before this object can be cleaned up
	remaining *atomic.Uint64
}

func newInflightRequest() *inflightRequest {
	return &inflightRequest{
		ch:        make(chan struct{}),
		remaining: atomic.NewUint64(0),
	}
}

// NewLeaseCache creates a new instance of a LeaseCache.
func NewLeaseCache(conf *LeaseCacheConfig) (*LeaseCache, error) {
	if conf == nil {
		return nil, errors.New("nil configuration provided")
	}

	if conf.Proxier == nil || conf.Logger == nil {
		return nil, fmt.Errorf("missing configuration required params: %v", conf)
	}

	if conf.Client == nil {
		return nil, fmt.Errorf("nil API client")
	}

	if conf.UserAgentToUse == "" {
		return nil, fmt.Errorf("no user agent specified -- see useragent.go")
	}

	db, err := cachememdb.New()
	if err != nil {
		return nil, err
	}

	// Create a base context for the lease cache layer
	baseCtxInfo := cachememdb.NewContextInfo(conf.BaseContext)

	return &LeaseCache{
		client:              conf.Client,
		proxier:             conf.Proxier,
		logger:              conf.Logger,
		userAgentToUse:      conf.UserAgentToUse,
		db:                  db,
		baseCtxInfo:         baseCtxInfo,
		l:                   &sync.RWMutex{},
		idLocks:             locksutil.CreateLocks(),
		inflightCache:       gocache.New(gocache.NoExpiration, gocache.NoExpiration),
		ps:                  conf.Storage,
		cacheStaticSecrets:  conf.CacheStaticSecrets,
		cacheDynamicSecrets: conf.CacheDynamicSecrets,
	}, nil
}

// SetCapabilityManager is a setter for CapabilityManager. If set, will manage capabilities
// for capability indexes.
func (c *LeaseCache) SetCapabilityManager(capabilityManager *StaticSecretCapabilityManager) {
	c.capabilityManager = capabilityManager
}

// SetShuttingDown is a setter for the shuttingDown field
func (c *LeaseCache) SetShuttingDown(in bool) {
	c.shuttingDown.Store(in)

	// Since we're shutting down, also stop the capability manager's jobs.
	// We can do this forcibly since no there's no reason to update
	// the cache when we're shutting down.
	if c.capabilityManager != nil {
		c.capabilityManager.Stop()
	}
}

// SetPersistentStorage is a setter for the persistent storage field in
// LeaseCache
func (c *LeaseCache) SetPersistentStorage(storageIn *cacheboltdb.BoltStorage) {
	c.ps = storageIn
}

// PersistentStorage is a getter for the persistent storage field in
// LeaseCache
func (c *LeaseCache) PersistentStorage() *cacheboltdb.BoltStorage {
	return c.ps
}

// checkCacheForDynamicSecretRequest checks the cache for a particular request based on its
// computed ID. It returns a non-nil *SendResponse if an entry is found.
func (c *LeaseCache) checkCacheForDynamicSecretRequest(id string) (*SendResponse, error) {
	c.logger.Trace("checking cache for dynamic secret request", "id", id)
	return c.checkCacheForRequest(id, nil)
}

// checkCacheForStaticSecretRequest checks the cache for a particular request based on its
// computed ID. It returns a non-nil *SendResponse if an entry is found.
// If a request is provided, it will validate that the token is allowed to retrieve this
// cache entry, and return nil if it isn't. It will also evict the cache if this is a non-GET
// request.
func (c *LeaseCache) checkCacheForStaticSecretRequest(id string, req *SendRequest) (*SendResponse, error) {
	c.logger.Trace("checking cache for static secret request", "id", id)
	return c.checkCacheForRequest(id, req)
}

// checkCacheForRequest checks the cache for a particular request based on its
// computed ID. It returns a non-nil *SendResponse if an entry is found.
// If a token is provided, it will validate that the token is allowed to retrieve this
// cache entry, and return nil if it isn't.
func (c *LeaseCache) checkCacheForRequest(id string, req *SendRequest) (*SendResponse, error) {
	index, err := c.db.Get(cachememdb.IndexNameID, id)
	if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	index.IndexLock.RLock()
	defer index.IndexLock.RUnlock()

	var token string
	if req != nil {
		// Req will be non-nil if we're checking for a static secret.
		// Token might still be "" if it's going to an unauthenticated
		// endpoint, or similar. For static secrets, we only care about
		// requests with tokens attached, as KV is authenticated.
		token = req.Token
	}

	if token != "" {
		// We are checking for a static secret. We need to ensure that this token
		// has previously demonstrated access to this static secret.
		// We could check the capabilities cache here, but since these
		// indexes should be in sync, this saves us an extra cache get.
		if _, ok := index.Tokens[token]; !ok {
			// We don't have access to this static secret, so
			// we do not return the cached response.
			return nil, nil
		}
	}

	var response []byte
	version := getStaticSecretVersionFromRequest(req)
	if version == 0 {
		response = index.Response
	} else {
		response = index.Versions[version]
	}

	// We don't have this response as either a current or older version.
	if response == nil {
		return nil, nil
	}

	// Cached request is found, deserialize the response
	reader := bufio.NewReader(bytes.NewReader(response))
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		c.logger.Error("failed to deserialize response", "error", err)
		return nil, err
	}

	sendResp, err := NewSendResponse(&api.Response{Response: resp}, response)
	if err != nil {
		c.logger.Error("failed to create new send response", "error", err)
		return nil, err
	}
	sendResp.CacheMeta.Hit = true

	respTime, err := http.ParseTime(resp.Header.Get("Date"))
	if err != nil {
		c.logger.Error("failed to parse cached response date", "error", err)
		return nil, err
	}
	sendResp.CacheMeta.Age = time.Now().Sub(respTime)

	return sendResp, nil
}

// Send performs a cache lookup on the incoming request. If it's a cache hit,
// it will return the cached response, otherwise it will delegate to the
// underlying Proxier and cache the received response.
func (c *LeaseCache) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// Compute the index ID for both static and dynamic secrets.
	// The primary difference is that for dynamic secrets, the
	// Vault token forms part of the index.
	dynamicSecretCacheId, err := computeIndexID(req)
	if err != nil {
		c.logger.Error("failed to compute cache key", "error", err)
		return nil, err
	}
	staticSecretCacheId := computeStaticSecretCacheIndex(req)

	// Check the inflight cache to see if there are other inflight requests
	// of the same kind, based on the computed ID. If so, we increment a counter

	// Note: we lock both the dynamic secret cache ID and the static secret cache ID
	// as at this stage, we don't know what kind of secret it is.
	var inflight *inflightRequest

	defer func() {
		// Cleanup on the cache if there are no remaining inflight requests.
		// This is the last step, so we defer the call first
		if inflight != nil && inflight.remaining.Load() == 0 {
			c.inflightCache.Delete(dynamicSecretCacheId)
			if staticSecretCacheId != "" {
				c.inflightCache.Delete(staticSecretCacheId)
			}
		}
	}()

	idLockDynamicSecret := locksutil.LockForKey(c.idLocks, dynamicSecretCacheId)

	// Briefly grab an ID-based lock in here to emulate a load-or-store behavior
	// and prevent concurrent cacheable requests from being proxied twice if
	// they both miss the cache due to it being clean when peeking the cache
	// entry.
	idLockDynamicSecret.Lock()
	inflightRaw, found := c.inflightCache.Get(dynamicSecretCacheId)
	if found {
		idLockDynamicSecret.Unlock()
		inflight = inflightRaw.(*inflightRequest)
		inflight.remaining.Inc()
		defer inflight.remaining.Dec()

		// If found it means that there's an inflight request being processed.
		// We wait until that's finished before proceeding further.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-inflight.ch:
		}
	} else {
		if inflight == nil {
			inflight = newInflightRequest()
			inflight.remaining.Inc()
			defer inflight.remaining.Dec()
			defer close(inflight.ch)
		}

		c.inflightCache.Set(dynamicSecretCacheId, inflight, gocache.NoExpiration)
		idLockDynamicSecret.Unlock()
	}

	if staticSecretCacheId != "" {
		idLockStaticSecret := locksutil.LockForKey(c.idLocks, staticSecretCacheId)

		// Briefly grab an ID-based lock in here to emulate a load-or-store behavior
		// and prevent concurrent cacheable requests from being proxied twice if
		// they both miss the cache due to it being clean when peeking the cache
		// entry.
		idLockStaticSecret.Lock()
		inflightRaw, found = c.inflightCache.Get(staticSecretCacheId)
		if found {
			idLockStaticSecret.Unlock()
			inflight = inflightRaw.(*inflightRequest)
			inflight.remaining.Inc()
			defer inflight.remaining.Dec()

			// If found it means that there's an inflight request being processed.
			// We wait until that's finished before proceeding further.
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-inflight.ch:
			}
		} else {
			if inflight == nil {
				inflight = newInflightRequest()
				inflight.remaining.Inc()
				defer inflight.remaining.Dec()
				defer close(inflight.ch)
			}

			c.inflightCache.Set(staticSecretCacheId, inflight, gocache.NoExpiration)
			idLockStaticSecret.Unlock()
		}
	}

	// Check if the response for this request is already in the dynamic secret cache
	cachedResp, err := c.checkCacheForDynamicSecretRequest(dynamicSecretCacheId)
	if err != nil {
		return nil, err
	}
	if cachedResp != nil {
		c.logger.Debug("returning cached dynamic secret response", "path", req.Request.URL.Path)
		return cachedResp, nil
	}

	// Check if the response for this request is already in the static secret cache
	if staticSecretCacheId != "" && req.Request.Method == http.MethodGet && req.Token != "" {
		cachedResp, err = c.checkCacheForStaticSecretRequest(staticSecretCacheId, req)
		if err != nil {
			return nil, err
		}
		if cachedResp != nil {
			c.logger.Debug("returning cached static secret response", "id", staticSecretCacheId, "path", getStaticSecretPathFromRequest(req))
			return cachedResp, nil
		}
	}

	c.logger.Debug("forwarding request from cache", "method", req.Request.Method, "path", req.Request.URL.Path)

	// Pass the request down and get a response
	resp, err := c.proxier.Send(ctx, req)
	if err != nil {
		return resp, err
	}

	// If this is a non-2xx or if the returned response does not contain JSON payload,
	// we skip caching
	if resp.Response.StatusCode >= 300 || resp.Response.Header.Get("Content-Type") != "application/json" {
		return resp, err
	}

	// Get the namespace from the request header
	namespace := req.Request.Header.Get(consts.NamespaceHeaderName)
	// We need to populate an empty value since go-memdb will skip over indexes
	// that contain empty values.
	if namespace == "" {
		namespace = "root/"
	}

	// Build the index to cache based on the response received
	index := &cachememdb.Index{
		Namespace:   namespace,
		RequestPath: req.Request.URL.Path,
		LastRenewed: time.Now().UTC(),
	}

	secret, err := api.ParseSecret(bytes.NewReader(resp.ResponseBody))
	if err != nil {
		c.logger.Error("failed to parse response as secret", "error", err)
		return nil, err
	}

	isRevocation, err := c.handleRevocationRequest(ctx, req, resp)
	if err != nil {
		c.logger.Error("failed to process the response", "error", err)
		return nil, err
	}

	// If this is a revocation request, do not go through cache logic.
	if isRevocation {
		return resp, nil
	}

	// Fast path for responses with no secrets
	if secret == nil {
		c.logger.Debug("pass-through response; no secret in response", "method", req.Request.Method, "path", req.Request.URL.Path)
		return resp, nil
	}

	// There shouldn't be a situation where secret.MountType == "kv" and
	// staticSecretCacheId == "", but just in case.
	// We restrict this to GETs as those are all we want to cache.
	if c.cacheStaticSecrets && secret.MountType == "kv" &&
		staticSecretCacheId != "" && req.Request.Method == http.MethodGet {
		index.Type = cacheboltdb.StaticSecretType
		index.ID = staticSecretCacheId
		// We set the request path to be the canonical static secret path, so that
		// two differently shaped (but equivalent) requests to the same path
		// will be the same.
		// This differs slightly from dynamic secrets, where the /v1/ will be
		// included in the request path.
		index.RequestPath = getStaticSecretPathFromRequest(req)

		c.logger.Trace("attempting to cache static secret with following request path", "request path", index.RequestPath, "version", getStaticSecretVersionFromRequest(req))
		err := c.cacheStaticSecret(ctx, req, resp, index, secret)
		if err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		// Since it's not a static secret, set the ID to be the dynamic id
		index.ID = dynamicSecretCacheId
	}

	// Short-circuit if we've been configured to not cache dynamic secrets
	if !c.cacheDynamicSecrets {
		return resp, nil
	}

	// Short-circuit if the secret is not renewable
	tokenRenewable, err := secret.TokenIsRenewable()
	if err != nil {
		c.logger.Error("failed to parse renewable param", "error", err)
		return nil, err
	}
	if !secret.Renewable && !tokenRenewable {
		c.logger.Debug("pass-through response; secret not renewable", "method", req.Request.Method, "path", req.Request.URL.Path)
		return resp, nil
	}

	var renewCtxInfo *cachememdb.ContextInfo
	switch {
	case secret.LeaseID != "":
		c.logger.Debug("processing lease response", "method", req.Request.Method, "path", req.Request.URL.Path)
		entry, err := c.db.Get(cachememdb.IndexNameToken, req.Token)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			// If the lease belongs to a token that is not managed by the lease cache,
			// return the response without caching it.
			c.logger.Debug("pass-through lease response; token not managed by lease cache", "method", req.Request.Method, "path", req.Request.URL.Path)
			return resp, nil
		}
		if err != nil {
			return nil, err
		}

		// Derive a context for renewal using the token's context
		renewCtxInfo = cachememdb.NewContextInfo(entry.RenewCtxInfo.Ctx)

		index.Lease = secret.LeaseID
		index.LeaseToken = req.Token

		index.Type = cacheboltdb.LeaseType

	case secret.Auth != nil:
		c.logger.Debug("processing auth response", "method", req.Request.Method, "path", req.Request.URL.Path)

		// Check if this token creation request resulted in a non-orphan token, and if so
		// correctly set the parentCtx to the request's token context.
		var parentCtx context.Context
		if !secret.Auth.Orphan {
			entry, err := c.db.Get(cachememdb.IndexNameToken, req.Token)
			if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
				// If the lease belongs to a token that is not managed by the lease cache,
				// return the response without caching it.
				c.logger.Debug("pass-through lease response; parent token not managed by lease cache", "method", req.Request.Method, "path", req.Request.URL.Path)
				return resp, nil
			}
			if err != nil {
				return nil, err
			}

			c.logger.Debug("setting parent context", "method", req.Request.Method, "path", req.Request.URL.Path)
			parentCtx = entry.RenewCtxInfo.Ctx

			index.TokenParent = req.Token
		}

		renewCtxInfo = c.createCtxInfo(parentCtx)
		index.Token = secret.Auth.ClientToken
		index.TokenAccessor = secret.Auth.Accessor

		index.Type = cacheboltdb.LeaseType

	default:
		// We shouldn't be hitting this, but will err on the side of caution and
		// simply proxy.
		c.logger.Debug("pass-through response; secret without lease and token", "method", req.Request.Method, "path", req.Request.URL.Path)
		return resp, nil
	}

	// Serialize the response to store it in the cached index
	var respBytes bytes.Buffer
	err = resp.Response.Write(&respBytes)
	if err != nil {
		c.logger.Error("failed to serialize response", "error", err)
		return nil, err
	}

	// Reset the response body for upper layers to read
	if resp.Response.Body != nil {
		resp.Response.Body.Close()
	}
	resp.Response.Body = io.NopCloser(bytes.NewReader(resp.ResponseBody))

	// Set the index's Response
	index.Response = respBytes.Bytes()

	// Store the index ID in the lifetimewatcher context
	renewCtx := context.WithValue(renewCtxInfo.Ctx, contextIndexID, index.ID)

	// Store the lifetime watcher context in the index
	index.RenewCtxInfo = &cachememdb.ContextInfo{
		Ctx:        renewCtx,
		CancelFunc: renewCtxInfo.CancelFunc,
		DoneCh:     renewCtxInfo.DoneCh,
	}

	// Add extra information necessary for restoring from persisted cache
	index.RequestMethod = req.Request.Method
	index.RequestToken = req.Token
	index.RequestHeader = req.Request.Header

	if index.Type != cacheboltdb.StaticSecretType {
		// Store the index in the cache
		c.logger.Debug("storing dynamic secret response into the cache", "method", req.Request.Method, "path", req.Request.URL.Path, "id", index.ID)
		err = c.Set(ctx, index)
		if err != nil {
			c.logger.Error("failed to cache the proxied response", "error", err)
			return nil, err
		}

		// Start renewing the secret in the response
		go c.startRenewing(renewCtx, index, req, secret)
	}

	return resp, nil
}

func (c *LeaseCache) cacheStaticSecret(ctx context.Context, req *SendRequest, resp *SendResponse, index *cachememdb.Index, secret *api.Secret) error {
	// If a cached version of this secret exists, we now have access, so
	// we don't need to re-cache, just update index.Tokens
	indexFromCache, err := c.db.Get(cachememdb.IndexNameID, index.ID)
	if err != nil && !errors.Is(err, cachememdb.ErrCacheItemNotFound) {
		return err
	}

	version := getStaticSecretVersionFromRequest(req)

	// The index already exists, so all we need to do is add our token
	// to the index's allowed token list, and if necessary, the new version,
	// then re-store it.
	if indexFromCache != nil {
		// We must hold a lock for the index while it's being updated.
		// We keep the two locking mechanisms distinct, so that it's only writes
		// that have to be serial.
		indexFromCache.IndexLock.Lock()
		defer indexFromCache.IndexLock.Unlock()
		indexFromCache.Tokens[req.Token] = struct{}{}

		// Are we looking for a version that's already cached?
		haveVersion := false
		if version != 0 {
			_, ok := indexFromCache.Versions[version]
			if ok {
				haveVersion = true
			}
		} else {
			if indexFromCache.Response != nil {
				haveVersion = true
			}
		}

		if !haveVersion {
			var respBytes bytes.Buffer
			err = resp.Response.Write(&respBytes)
			if err != nil {
				c.logger.Error("failed to serialize response", "error", err)
				return err
			}

			// Reset the response body for upper layers to read
			if resp.Response.Body != nil {
				resp.Response.Body.Close()
			}
			resp.Response.Body = io.NopCloser(bytes.NewReader(resp.ResponseBody))

			// Set the index's Response
			if version == 0 {
				indexFromCache.Response = respBytes.Bytes()
				// For current KVv2 secrets, see if we can add the version that the secret is
				// to the versions map, too. If we got the latest version and the version is #2,
				// also update Versions[2]
				c.addToVersionListForCurrentVersionKVv2Secret(indexFromCache, secret)
			} else {
				indexFromCache.Versions[version] = respBytes.Bytes()
			}
		}

		return c.storeStaticSecretIndex(ctx, req, indexFromCache)
	}

	// Serialize the response to store it in the cached index
	var respBytes bytes.Buffer
	err = resp.Response.Write(&respBytes)
	if err != nil {
		c.logger.Error("failed to serialize response", "error", err)
		return err
	}

	// Reset the response body for upper layers to read
	if resp.Response.Body != nil {
		resp.Response.Body.Close()
	}
	resp.Response.Body = io.NopCloser(bytes.NewReader(resp.ResponseBody))

	// Initialize the versions
	index.Versions = map[int][]byte{}

	// Set the index's Response
	if version == 0 {
		index.Response = respBytes.Bytes()
		// For current KVv2 secrets, see if we can add the version that the secret is
		// to the versions map, too. If we got the latest version and the version is #2,
		// also update Versions[2]
		c.addToVersionListForCurrentVersionKVv2Secret(index, secret)
	} else {
		index.Versions[version] = respBytes.Bytes()
	}

	// Initialize the token map and add this token to it.
	index.Tokens = map[string]struct{}{req.Token: {}}

	// Set the index type
	index.Type = cacheboltdb.StaticSecretType

	// Store the index:
	return c.storeStaticSecretIndex(ctx, req, index)
}

// addToVersionListForCurrentVersionKVv2Secret takes a secret index and, if it's
// a KVv2 secret, adds the given response to the corresponding version for it.
// This function fails silently, as we could be parsing arbitrary JSON.
// This function can store a version for a KVv1 secret iff:
// - It has 'data' in the path
// - It has a numerical 'metadata.version' field
// However, this risk seems very small, and the negatives of such a secret being
// stored in the cache aren't worth additional mitigations to check if it's a KVv1
// or KVv2 mount (such as doing a 'preflight' request like the CLI).
// There's no way to access it and it's just a couple of extra bytes, in the
// case that this does happen to a KVv1 secret.
func (c *LeaseCache) addToVersionListForCurrentVersionKVv2Secret(index *cachememdb.Index, secret *api.Secret) {
	if secret != nil {
		// First do an imperfect but lightweight check. This saves parsing the secret in the case that the secret isn't KVv2.
		// KVv2 secrets always contain /data/, but KVv1 secrets can too, so we can't rely on this.
		if strings.Contains(index.RequestPath, "/data/") {
			metadata, ok := secret.Data["metadata"]
			if ok {
				metaDataAsMap, ok := metadata.(map[string]interface{})
				if ok {
					versionJson, ok := metaDataAsMap["version"].(json.Number)
					if ok {
						versionInt64, err := versionJson.Int64()
						if err == nil {
							version := int(versionInt64)
							c.logger.Trace("adding response for current KVv2 secret to index's Versions map", "path", index.RequestPath, "version", version)

							if index.Versions == nil {
								index.Versions = map[int][]byte{}
							}

							index.Versions[version] = index.Response
						}
					}
				}
			}
		}
	}
}

func (c *LeaseCache) storeStaticSecretIndex(ctx context.Context, req *SendRequest, index *cachememdb.Index) error {
	// Store the index in the cache
	c.logger.Debug("storing static secret response into the cache", "path", index.RequestPath, "id", index.ID)
	err := c.Set(ctx, index)
	if err != nil {
		c.logger.Error("failed to cache the proxied response", "error", err)
		return err
	}

	capabilitiesIndex, created, err := c.retrieveOrCreateTokenCapabilitiesEntry(req.Token)
	if err != nil {
		c.logger.Error("failed to cache the proxied response", "error", err)
		return err
	}

	path := getStaticSecretPathFromRequest(req)

	capabilitiesIndex.IndexLock.Lock()
	// Extra caution -- avoid potential nil
	if capabilitiesIndex.ReadablePaths == nil {
		capabilitiesIndex.ReadablePaths = make(map[string]struct{})
	}

	// update the index with the new capability:
	capabilitiesIndex.ReadablePaths[path] = struct{}{}
	capabilitiesIndex.IndexLock.Unlock()

	err = c.SetCapabilitiesIndex(ctx, capabilitiesIndex)
	if err != nil {
		c.logger.Error("failed to cache token capabilities as part of caching the proxied response", "error", err)
		return err
	}

	// Lastly, ensure that we start renewing this index, if it's new.
	// We require the 'created' check so that we don't renew the same
	// index multiple times.
	if c.capabilityManager != nil && created {
		c.capabilityManager.StartRenewingCapabilities(capabilitiesIndex)
	}

	return nil
}

// retrieveOrCreateTokenCapabilitiesEntry will either retrieve the token
// capabilities entry from the cache, or create a new, empty one.
// The bool represents if a new token capability has been created.
func (c *LeaseCache) retrieveOrCreateTokenCapabilitiesEntry(token string) (*cachememdb.CapabilitiesIndex, bool, error) {
	// The index ID is a hash of the token.
	indexId := hashStaticSecretIndex(token)
	indexFromCache, err := c.db.GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	if err != nil && !errors.Is(err, cachememdb.ErrCacheItemNotFound) {
		return nil, false, err
	}

	if indexFromCache != nil {
		return indexFromCache, false, nil
	}

	// Build the index to cache based on the response received
	index := &cachememdb.CapabilitiesIndex{
		ID:            indexId,
		Token:         token,
		ReadablePaths: make(map[string]struct{}),
	}

	return index, true, nil
}

func (c *LeaseCache) createCtxInfo(ctx context.Context) *cachememdb.ContextInfo {
	if ctx == nil {
		c.l.RLock()
		ctx = c.baseCtxInfo.Ctx
		c.l.RUnlock()
	}
	return cachememdb.NewContextInfo(ctx)
}

func (c *LeaseCache) startRenewing(ctx context.Context, index *cachememdb.Index, req *SendRequest, secret *api.Secret) {
	defer func() {
		id := ctx.Value(contextIndexID).(string)
		if c.shuttingDown.Load() {
			c.logger.Trace("not evicting index from cache during shutdown", "id", id, "method", req.Request.Method, "path", req.Request.URL.Path)
			return
		}
		c.logger.Debug("evicting index from cache", "id", id, "method", req.Request.Method, "path", req.Request.URL.Path)
		err := c.Evict(index)
		if err != nil {
			c.logger.Error("failed to evict index", "id", id, "error", err)
			return
		}
	}()

	client, err := c.client.Clone()
	if err != nil {
		c.logger.Error("failed to create API client in the lifetime watcher", "error", err)
		return
	}
	client.SetToken(req.Token)

	headers := client.Headers()
	if headers == nil {
		headers = make(http.Header)
	}

	// We do not preserve any initial User-Agent here since these requests are from
	// the proxy subsystem, but are made by the lease cache's lifetime watcher,
	// not triggered by a specific request.
	headers.Set("User-Agent", c.userAgentToUse)
	client.SetHeaders(headers)

	watcher, err := client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret: secret,
	})
	if err != nil {
		c.logger.Error("failed to create secret lifetime watcher", "error", err)
		return
	}

	c.logger.Debug("initiating renewal", "method", req.Request.Method, "path", req.Request.URL.Path)
	go watcher.Start()
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			// This is the case which captures context cancellations from token
			// and leases. Since all the contexts are derived from the agent's
			// context, this will also cover the shutdown scenario.
			c.logger.Debug("context cancelled; stopping lifetime watcher", "path", req.Request.URL.Path)
			return
		case err := <-watcher.DoneCh():
			// This case covers renewal completion and renewal errors
			if err != nil {
				c.logger.Error("failed to renew secret", "error", err)
				return
			}
			c.logger.Debug("renewal halted; evicting from cache", "path", req.Request.URL.Path)
			return
		case <-watcher.RenewCh():
			c.logger.Debug("secret renewed", "path", req.Request.URL.Path)
			if c.ps != nil {
				if err := c.updateLastRenewed(ctx, index, time.Now().UTC()); err != nil {
					c.logger.Warn("not able to update lastRenewed time for cached index", "id", index.ID)
				}
			}
		case <-index.RenewCtxInfo.DoneCh:
			// This case indicates the renewal process to shutdown and evict
			// the cache entry. This is triggered when a specific secret
			// renewal needs to be killed without affecting any of the derived
			// context renewals.
			c.logger.Debug("done channel closed")
			return
		}
	}
}

func (c *LeaseCache) updateLastRenewed(ctx context.Context, index *cachememdb.Index, t time.Time) error {
	idLock := locksutil.LockForKey(c.idLocks, index.ID)
	idLock.Lock()
	defer idLock.Unlock()

	getIndex, err := c.db.Get(cachememdb.IndexNameID, index.ID)
	if err != nil && err != cachememdb.ErrCacheItemNotFound {
		return err
	}
	index.LastRenewed = t
	if err := c.Set(ctx, getIndex); err != nil {
		return err
	}
	return nil
}

// computeIndexID results in a value that uniquely identifies a request
// received by the agent. It does so by SHA256 hashing the serialized request
// object containing the request path, query parameters and body parameters.
func computeIndexID(req *SendRequest) (string, error) {
	var b bytes.Buffer

	cloned := req.Request.Clone(context.Background())
	cloned.Header.Del(vaulthttp.VaultIndexHeaderName)
	cloned.Header.Del(vaulthttp.VaultForwardHeaderName)
	cloned.Header.Del(vaulthttp.VaultInconsistentHeaderName)
	// Serialize the request
	if err := cloned.Write(&b); err != nil {
		return "", fmt.Errorf("failed to serialize request: %v", err)
	}

	// Reset the request body after it has been closed by Write
	req.Request.Body = io.NopCloser(bytes.NewReader(req.RequestBody))

	// Append req.Token into the byte slice. This is needed since auto-auth'ed
	// requests sets the token directly into SendRequest.Token
	if _, err := b.WriteString(req.Token); err != nil {
		return "", fmt.Errorf("failed to write token to hash input: %w", err)
	}

	return hex.EncodeToString(cryptoutil.Blake2b256Hash(string(b.Bytes()))), nil
}

// canonicalizeStaticSecretPath takes an API request path such as
// /v1/foo/bar and a namespace, and turns it into a canonical representation
// of the secret's path in Vault.
// We opt for this form as namespace.Canonicalize returns a namespace in the
// form of "ns1/", so we keep consistent with path canonicalization.
func canonicalizeStaticSecretPath(requestPath string, ns string) string {
	// /sys/capabilities accepts both requests that look like foo/bar
	// and /foo/bar but not /v1/foo/bar.
	// We trim the /v1/ from the start of the URL to get the foo/bar form.
	// This means that we can use the paths we retrieve from the
	// /sys/capabilities endpoint to access this index
	// without having to re-add the /v1/
	path := strings.TrimPrefix(requestPath, "/v1/")
	// Trim any leading slashes, as we never want those.
	// This ensures /foo/bar gets turned to foo/bar
	path = strings.TrimPrefix(path, "/")

	// If a namespace was provided in a way that wasn't directly in the path,
	// it must be added to the path.
	path = namespace.Canonicalize(ns) + path

	return path
}

// getStaticSecretVersionFromRequest gets the version of a secret
// from a request. For the latest secret and for KVv1 secrets,
// this will return 0.
func getStaticSecretVersionFromRequest(req *SendRequest) int {
	if req == nil || req.Request == nil {
		return 0
	}
	version := req.Request.FormValue("version")
	if version == "" {
		return 0
	}
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		// It's not a valid version.
		return 0
	}
	return versionInt
}

// getStaticSecretPathFromRequest gets the canonical path for a
// request, taking into account intricacies relating to /v1/ and namespaces
// in the header.
// Returns a path like foo/bar or ns1/foo/bar.
// We opt for this form as namespace.Canonicalize returns a namespace in the
// form of "ns1/", so we keep consistent with path canonicalization.
func getStaticSecretPathFromRequest(req *SendRequest) string {
	path := req.Request.URL.Path
	// Static secrets always have /v1 as a prefix. This enables us to
	// enable a pass-through and never attempt to cache or view-from-cache
	// any request without the /v1 prefix.
	if !strings.HasPrefix(path, "/v1") {
		return ""
	}
	var namespace string
	if header := req.Request.Header; header != nil {
		namespace = header.Get(api.NamespaceHeaderName)
	}
	return canonicalizeStaticSecretPath(path, namespace)
}

// hashStaticSecretIndex is a simple function that hashes the path into
// a function. This is kept as a helper function for ease of use by downstream functions.
func hashStaticSecretIndex(unhashedIndex string) string {
	return hex.EncodeToString(cryptoutil.Blake2b256Hash(unhashedIndex))
}

// computeStaticSecretCacheIndex results in a value that uniquely identifies a static
// secret's cached ID. Notably, we intentionally ignore headers (for example,
// the X-Vault-Token header) to remain agnostic to which token is being
// used in the request. We care only about the path.
// This will return "" if the index does not have a /v1 prefix, and therefore
// cannot be a static secret.
func computeStaticSecretCacheIndex(req *SendRequest) string {
	path := getStaticSecretPathFromRequest(req)
	if path == "" {
		return path
	}

	return hashStaticSecretIndex(path)
}

// HandleCacheClear returns a handlerFunc that can perform cache clearing operations.
func (c *LeaseCache) HandleCacheClear(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the cache is not enabled, return a 200
		if c == nil {
			return
		}

		// Only handle POST/PUT requests
		switch r.Method {
		case http.MethodPost:
		case http.MethodPut:
		default:
			return
		}

		req := new(cacheClearRequest)
		if err := jsonutil.DecodeJSONFromReader(r.Body, req); err != nil {
			if err == io.EOF {
				err = errors.New("empty JSON provided")
			}
			logical.RespondError(w, http.StatusBadRequest, fmt.Errorf("failed to parse JSON input: %w", err))
			return
		}

		c.logger.Debug("received cache-clear request", "type", req.Type, "namespace", req.Namespace, "value", req.Value)

		in, err := parseCacheClearInput(req)
		if err != nil {
			c.logger.Error("unable to parse clear input", "error", err)
			logical.RespondError(w, http.StatusBadRequest, fmt.Errorf("failed to parse clear input: %w", err))
			return
		}

		if err := c.handleCacheClear(ctx, in); err != nil {
			// Default to 500 on error, unless the user provided an invalid type,
			// which would then be a 400.
			httpStatus := http.StatusInternalServerError
			if errors.Is(err, errInvalidType) {
				httpStatus = http.StatusBadRequest
			}
			logical.RespondError(w, httpStatus, fmt.Errorf("failed to clear cache: %w", err))
			return
		}

		return
	})
}

func (c *LeaseCache) handleCacheClear(ctx context.Context, in *cacheClearInput) error {
	if in == nil {
		return errors.New("no value(s) provided to clear corresponding cache entries")
	}

	switch in.Type {
	case "request_path":
		// For this particular case, we need to ensure that there are 2 provided
		// indexers for the proper lookup.
		if in.RequestPath == "" {
			return errors.New("request path not provided")
		}

		// The first value provided for this case will be the namespace, but if it's
		// an empty value we need to overwrite it with "root/" to ensure proper
		// cache lookup.
		if in.Namespace == "" {
			in.Namespace = "root/"
		}

		// Find all the cached entries which has the given request path and
		// cancel the contexts of all the respective lifetime watchers
		indexes, err := c.db.GetByPrefix(cachememdb.IndexNameRequestPath, in.Namespace, in.RequestPath)
		if err != nil {
			return err
		}
		for _, index := range indexes {
			// If it's a static secret, we must remove directly, as there
			// is no renew func to cancel.
			if index.Type == cacheboltdb.StaticSecretType {
				err = c.db.Evict(cachememdb.IndexNameID, index.ID)
				if err != nil {
					return err
				}
			} else {
				if index.RenewCtxInfo != nil {
					if index.RenewCtxInfo.CancelFunc != nil {
						index.RenewCtxInfo.CancelFunc()
					}
				}
			}
		}

	case "token":
		if in.Token == "" {
			return errors.New("token not provided")
		}

		// Get the context for the given token and cancel its context
		index, err := c.db.Get(cachememdb.IndexNameToken, in.Token)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		c.logger.Debug("canceling context of index attached to token")

		index.RenewCtxInfo.CancelFunc()

	case "token_accessor":
		if in.TokenAccessor == "" && in.Type != cacheboltdb.StaticSecretType {
			return errors.New("token accessor not provided")
		}

		// Get the cached index and cancel the corresponding lifetime watcher
		// context
		index, err := c.db.Get(cachememdb.IndexNameTokenAccessor, in.TokenAccessor)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		c.logger.Debug("canceling context of index attached to accessor")

		index.RenewCtxInfo.CancelFunc()

	case "lease":
		if in.Lease == "" {
			return errors.New("lease not provided")
		}

		// Get the cached index and cancel the corresponding lifetime watcher
		// context
		index, err := c.db.Get(cachememdb.IndexNameLease, in.Lease)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		c.logger.Debug("canceling context of index attached to accessor")

		index.RenewCtxInfo.CancelFunc()

	case "all":
		// Cancel the base context which triggers all the goroutines to
		// stop and evict entries from cache.
		c.logger.Debug("canceling base context")
		c.l.Lock()
		c.baseCtxInfo.CancelFunc()
		// Reset the base context
		baseCtx, baseCancel := context.WithCancel(ctx)
		c.baseCtxInfo = &cachememdb.ContextInfo{
			Ctx:        baseCtx,
			CancelFunc: baseCancel,
		}
		c.l.Unlock()

		// Reset the memdb instance (and persistent storage if enabled)
		if err := c.Flush(); err != nil {
			return err
		}

	default:
		return errInvalidType
	}

	c.logger.Debug("successfully cleared matching cache entries")

	return nil
}

// handleRevocationRequest checks whether the originating request is a
// revocation request, and if so perform applicable cache cleanups.
// Returns true is this is a revocation request.
func (c *LeaseCache) handleRevocationRequest(ctx context.Context, req *SendRequest, resp *SendResponse) (bool, error) {
	// Lease and token revocations return 204's on success. Fast-path if that's
	// not the case.
	if resp.Response.StatusCode != http.StatusNoContent {
		return false, nil
	}

	_, path := deriveNamespaceAndRevocationPath(req)

	switch {
	case path == vaultPathTokenRevoke:
		// Get the token from the request body
		jsonBody := map[string]interface{}{}
		if err := json.Unmarshal(req.RequestBody, &jsonBody); err != nil {
			return false, err
		}
		tokenRaw, ok := jsonBody["token"]
		if !ok {
			return false, fmt.Errorf("failed to get token from request body")
		}
		token, ok := tokenRaw.(string)
		if !ok {
			return false, fmt.Errorf("expected token in the request body to be string")
		}

		// Clear the cache entry associated with the token and all the other
		// entries belonging to the leases derived from this token.
		in := &cacheClearInput{
			Type:  "token",
			Token: token,
		}
		if err := c.handleCacheClear(ctx, in); err != nil {
			return false, err
		}

	case path == vaultPathTokenRevokeSelf:
		// Clear the cache entry associated with the token and all the other
		// entries belonging to the leases derived from this token.
		in := &cacheClearInput{
			Type:  "token",
			Token: req.Token,
		}
		if err := c.handleCacheClear(ctx, in); err != nil {
			return false, err
		}

	case path == vaultPathTokenRevokeAccessor:
		jsonBody := map[string]interface{}{}
		if err := json.Unmarshal(req.RequestBody, &jsonBody); err != nil {
			return false, err
		}
		accessorRaw, ok := jsonBody["accessor"]
		if !ok {
			return false, fmt.Errorf("failed to get accessor from request body")
		}
		accessor, ok := accessorRaw.(string)
		if !ok {
			return false, fmt.Errorf("expected accessor in the request body to be string")
		}

		in := &cacheClearInput{
			Type:          "token_accessor",
			TokenAccessor: accessor,
		}
		if err := c.handleCacheClear(ctx, in); err != nil {
			return false, err
		}

	case path == vaultPathTokenRevokeOrphan:
		jsonBody := map[string]interface{}{}
		if err := json.Unmarshal(req.RequestBody, &jsonBody); err != nil {
			return false, err
		}
		tokenRaw, ok := jsonBody["token"]
		if !ok {
			return false, fmt.Errorf("failed to get token from request body")
		}
		token, ok := tokenRaw.(string)
		if !ok {
			return false, fmt.Errorf("expected token in the request body to be string")
		}

		// Kill the lifetime watchers of all the leases attached to the revoked
		// token
		indexes, err := c.db.GetByPrefix(cachememdb.IndexNameLeaseToken, token)
		if err != nil {
			return false, err
		}
		for _, index := range indexes {
			index.RenewCtxInfo.CancelFunc()
		}

		// Kill the lifetime watchers of the revoked token
		index, err := c.db.Get(cachememdb.IndexNameToken, token)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			return true, nil
		}
		if err != nil {
			return false, err
		}

		// Indicate the lifetime watcher goroutine for this index to return.
		// This will not affect the child tokens because the context is not
		// getting cancelled.
		close(index.RenewCtxInfo.DoneCh)

		// Clear the parent references of the revoked token in the entries
		// belonging to the child tokens of the revoked token.
		indexes, err = c.db.GetByPrefix(cachememdb.IndexNameTokenParent, token)
		if err != nil {
			return false, err
		}
		for _, index := range indexes {
			index.TokenParent = ""
			err = c.db.Set(index)
			if err != nil {
				c.logger.Error("failed to persist index", "error", err)
				return false, err
			}
		}

	case path == vaultPathLeaseRevoke:
		// TODO: Should lease present in the URL itself be considered here?
		// Get the lease from the request body
		jsonBody := map[string]interface{}{}
		if err := json.Unmarshal(req.RequestBody, &jsonBody); err != nil {
			return false, err
		}
		leaseIDRaw, ok := jsonBody["lease_id"]
		if !ok {
			return false, fmt.Errorf("failed to get lease_id from request body")
		}
		leaseID, ok := leaseIDRaw.(string)
		if !ok {
			return false, fmt.Errorf("expected lease_id the request body to be string")
		}
		in := &cacheClearInput{
			Type:  "lease",
			Lease: leaseID,
		}
		if err := c.handleCacheClear(ctx, in); err != nil {
			return false, err
		}

	case strings.HasPrefix(path, vaultPathLeaseRevokeForce):
		// Trim the URL path to get the request path prefix
		prefix := strings.TrimPrefix(path, vaultPathLeaseRevokeForce)
		// Get all the cache indexes that use the request path containing the
		// prefix and cancel the lifetime watcher context of each.
		indexes, err := c.db.GetByPrefix(cachememdb.IndexNameLease, prefix)
		if err != nil {
			return false, err
		}

		_, tokenNSID := namespace.SplitIDFromString(req.Token)
		for _, index := range indexes {
			_, leaseNSID := namespace.SplitIDFromString(index.Lease)
			// Only evict leases that match the token's namespace
			if tokenNSID == leaseNSID {
				index.RenewCtxInfo.CancelFunc()
			}
		}

	case strings.HasPrefix(path, vaultPathLeaseRevokePrefix):
		// Trim the URL path to get the request path prefix
		prefix := strings.TrimPrefix(path, vaultPathLeaseRevokePrefix)
		// Get all the cache indexes that use the request path containing the
		// prefix and cancel the lifetime watcher context of each.
		indexes, err := c.db.GetByPrefix(cachememdb.IndexNameLease, prefix)
		if err != nil {
			return false, err
		}

		_, tokenNSID := namespace.SplitIDFromString(req.Token)
		for _, index := range indexes {
			_, leaseNSID := namespace.SplitIDFromString(index.Lease)
			// Only evict leases that match the token's namespace
			if tokenNSID == leaseNSID {
				index.RenewCtxInfo.CancelFunc()
			}
		}

	default:
		return false, nil
	}

	c.logger.Debug("triggered caching eviction from revocation request")

	return true, nil
}

// Set stores the index in the cachememdb, and also stores it in the persistent
// cache (if enabled)
func (c *LeaseCache) Set(ctx context.Context, index *cachememdb.Index) error {
	if err := c.db.Set(index); err != nil {
		return err
	}

	if c.ps != nil {
		plaintext, err := index.Serialize()
		if err != nil {
			return err
		}

		if err := c.ps.Set(ctx, index.ID, plaintext, index.Type); err != nil {
			return err
		}
		c.logger.Trace("set entry in persistent storage", "type", index.Type, "path", index.RequestPath, "id", index.ID)
	}

	return nil
}

// SetCapabilitiesIndex stores the capabilities index in the cachememdb, and also stores it in the persistent
// cache (if enabled)
func (c *LeaseCache) SetCapabilitiesIndex(ctx context.Context, index *cachememdb.CapabilitiesIndex) error {
	if err := c.db.SetCapabilitiesIndex(index); err != nil {
		return err
	}

	if c.ps != nil {
		plaintext, err := index.SerializeCapabilitiesIndex()
		if err != nil {
			return err
		}

		if err := c.ps.Set(ctx, index.ID, plaintext, cacheboltdb.TokenCapabilitiesType); err != nil {
			return err
		}
		c.logger.Trace("set entry in persistent storage", "type", cacheboltdb.TokenCapabilitiesType, "id", index.ID)
	}

	return nil
}

// Evict removes an Index from the cachememdb, and also removes it from the
// persistent cache (if enabled)
func (c *LeaseCache) Evict(index *cachememdb.Index) error {
	if err := c.db.Evict(cachememdb.IndexNameID, index.ID); err != nil {
		return err
	}

	if c.ps != nil {
		if err := c.ps.Delete(index.ID, index.Type); err != nil {
			return err
		}
		c.logger.Trace("deleted item from persistent storage", "id", index.ID)
	}

	return nil
}

// Flush the cachememdb and persistent cache (if enabled)
func (c *LeaseCache) Flush() error {
	if err := c.db.Flush(); err != nil {
		return err
	}

	if c.ps != nil {
		c.logger.Trace("clearing persistent storage")
		return c.ps.Clear()
	}

	return nil
}

// Restore loads the cachememdb from the persistent storage passed in. Loads
// tokens first, since restoring a lease's renewal context and watcher requires
// looking up the token in the cachememdb.
// Restore also restarts any capability management for managed static secret
// tokens.
func (c *LeaseCache) Restore(ctx context.Context, storage *cacheboltdb.BoltStorage) error {
	var errs *multierror.Error

	// Process tokens first
	tokens, err := storage.GetByType(ctx, cacheboltdb.TokenType)
	if err != nil {
		errs = multierror.Append(errs, err)
	} else {
		if err := c.restoreTokens(tokens); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	// Then process leases
	leases, err := storage.GetByType(ctx, cacheboltdb.LeaseType)
	if err != nil {
		errs = multierror.Append(errs, err)
	} else {
		for _, lease := range leases {
			newIndex, err := cachememdb.Deserialize(lease)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}

			c.logger.Trace("restoring lease", "id", newIndex.ID, "path", newIndex.RequestPath)

			// Check if this lease has already expired
			expired, err := c.hasExpired(time.Now().UTC(), newIndex)
			if err != nil {
				c.logger.Warn("failed to check if lease is expired", "id", newIndex.ID, "error", err)
			}
			if expired {
				continue
			}

			if err := c.restoreLeaseRenewCtx(newIndex); err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			if err := c.db.Set(newIndex); err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			c.logger.Trace("restored lease", "id", newIndex.ID, "path", newIndex.RequestPath)
		}
	}

	// Then process static secrets and their capabilities
	if c.cacheStaticSecrets {
		staticSecrets, err := storage.GetByType(ctx, cacheboltdb.StaticSecretType)
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			for _, staticSecret := range staticSecrets {
				newIndex, err := cachememdb.Deserialize(staticSecret)
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}

				c.logger.Trace("restoring static secret index", "id", newIndex.ID, "path", newIndex.RequestPath)
				if err := c.db.Set(newIndex); err != nil {
					errs = multierror.Append(errs, err)
					continue
				}
			}
		}

		capabilityIndexes, err := storage.GetByType(ctx, cacheboltdb.TokenCapabilitiesType)
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			for _, capabilityIndex := range capabilityIndexes {
				newIndex, err := cachememdb.DeserializeCapabilitiesIndex(capabilityIndex)
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}

				c.logger.Trace("restoring capability index", "id", newIndex.ID)
				if err := c.db.SetCapabilitiesIndex(newIndex); err != nil {
					errs = multierror.Append(errs, err)
					continue
				}

				if c.capabilityManager != nil {
					c.capabilityManager.StartRenewingCapabilities(newIndex)
				}
			}
		}
	}

	return errs.ErrorOrNil()
}

func (c *LeaseCache) restoreTokens(tokens [][]byte) error {
	var errors *multierror.Error

	for _, token := range tokens {
		newIndex, err := cachememdb.Deserialize(token)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		newIndex.RenewCtxInfo = c.createCtxInfo(nil)
		if err := c.db.Set(newIndex); err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		c.logger.Trace("restored token", "id", newIndex.ID)
	}

	return errors.ErrorOrNil()
}

// restoreLeaseRenewCtx re-creates a RenewCtx for an index object and starts
// the watcher go routine
func (c *LeaseCache) restoreLeaseRenewCtx(index *cachememdb.Index) error {
	if index.Response == nil {
		return fmt.Errorf("cached response was nil for %s", index.ID)
	}

	// Parse the secret to determine which type it is
	reader := bufio.NewReader(bytes.NewReader(index.Response))
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		c.logger.Error("failed to deserialize response", "error", err)
		return err
	}
	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		c.logger.Error("failed to parse response as secret", "error", err)
		return err
	}

	var renewCtxInfo *cachememdb.ContextInfo
	switch {
	case secret.LeaseID != "":
		entry, err := c.db.Get(cachememdb.IndexNameToken, index.RequestToken)
		if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
			return fmt.Errorf("could not find parent Token %s for req path %s", index.RequestToken, index.RequestPath)
		}
		if err != nil {
			return err
		}

		// Derive a context for renewal using the token's context
		renewCtxInfo = cachememdb.NewContextInfo(entry.RenewCtxInfo.Ctx)

	case secret.Auth != nil:
		var parentCtx context.Context
		if !secret.Auth.Orphan {
			entry, err := c.db.Get(cachememdb.IndexNameToken, index.RequestToken)
			if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
				// If parent token is not managed by the cache, child shouldn't be
				// either.
				if entry == nil {
					return fmt.Errorf("could not find parent Token %s for req path %s", index.RequestToken, index.RequestPath)
				}
			}
			if err != nil {
				return err
			}

			c.logger.Debug("setting parent context", "method", index.RequestMethod, "path", index.RequestPath)
			parentCtx = entry.RenewCtxInfo.Ctx
		}
		renewCtxInfo = c.createCtxInfo(parentCtx)
	default:
		// This isn't a renewable cache entry, i.e. a static secret cache entry.
		// We return, because there's nothing to do.
		return nil
	}

	renewCtx := context.WithValue(renewCtxInfo.Ctx, contextIndexID, index.ID)
	index.RenewCtxInfo = &cachememdb.ContextInfo{
		Ctx:        renewCtx,
		CancelFunc: renewCtxInfo.CancelFunc,
		DoneCh:     renewCtxInfo.DoneCh,
	}

	sendReq := &SendRequest{
		Token: index.RequestToken,
		Request: &http.Request{
			Header: index.RequestHeader,
			Method: index.RequestMethod,
			URL: &url.URL{
				Path: index.RequestPath,
			},
		},
	}
	go c.startRenewing(renewCtx, index, sendReq, secret)

	return nil
}

// deriveNamespaceAndRevocationPath returns the namespace and relative path for
// revocation paths.
//
// If the path contains a namespace, but it's not a revocation path, it will be
// returned as-is, since there's no way to tell where the namespace ends and
// where the request path begins purely based off a string.
//
// Case 1: /v1/ns1/leases/revoke  -> ns1/, /v1/leases/revoke
// Case 2: ns1/ /v1/leases/revoke -> ns1/, /v1/leases/revoke
// Case 3: /v1/ns1/foo/bar  -> root/, /v1/ns1/foo/bar
// Case 4: ns1/ /v1/foo/bar -> ns1/, /v1/foo/bar
func deriveNamespaceAndRevocationPath(req *SendRequest) (string, string) {
	namespace := "root/"
	nsHeader := req.Request.Header.Get(consts.NamespaceHeaderName)
	if nsHeader != "" {
		namespace = nsHeader
	}

	fullPath := req.Request.URL.Path
	nonVersionedPath := strings.TrimPrefix(fullPath, "/v1")

	for _, pathToCheck := range revocationPaths {
		// We use strings.Contains here for paths that can contain
		// vars in the path, e.g. /v1/lease/revoke-prefix/:prefix
		i := strings.Index(nonVersionedPath, pathToCheck)
		// If there's no match, move on to the next check
		if i == -1 {
			continue
		}

		// If the index is 0, this is a relative path with no namespace preppended,
		// so we can break early
		if i == 0 {
			break
		}

		// We need to turn /ns1 into ns1/, this makes it easy
		namespaceInPath := nshelper.Canonicalize(nonVersionedPath[:i])

		// If it's root, we replace, otherwise we join
		if namespace == "root/" {
			namespace = namespaceInPath
		} else {
			namespace = namespace + namespaceInPath
		}

		return namespace, fmt.Sprintf("/v1%s", nonVersionedPath[i:])
	}

	return namespace, fmt.Sprintf("/v1%s", nonVersionedPath)
}

// RegisterAutoAuthToken adds the provided auto-token into the cache. This is
// primarily used to register the auto-auth token and should only be called
// within a sink's WriteToken func.
func (c *LeaseCache) RegisterAutoAuthToken(token string) error {
	// Get the token from the cache
	oldIndex, err := c.db.Get(cachememdb.IndexNameToken, token)
	if err != nil && err != cachememdb.ErrCacheItemNotFound {
		return err
	}

	// If the index is found, just keep it in the cache and ignore the incoming
	// token (since they're the same)
	if oldIndex != nil {
		c.logger.Trace("auto-auth token already exists in cache; no need to store it again")
		return nil
	}

	// The following randomly generated values are required for index stored by
	// the cache, but are not actually used. We use random values to prevent
	// accidental access.
	id, err := base62.Random(5)
	if err != nil {
		return err
	}
	namespace, err := base62.Random(5)
	if err != nil {
		return err
	}
	requestPath, err := base62.Random(5)
	if err != nil {
		return err
	}

	index := &cachememdb.Index{
		ID:          id,
		Token:       token,
		Namespace:   namespace,
		RequestPath: requestPath,
		Type:        cacheboltdb.TokenType,
	}

	// Derive a context off of the lease cache's base context
	ctxInfo := c.createCtxInfo(nil)

	index.RenewCtxInfo = &cachememdb.ContextInfo{
		Ctx:        ctxInfo.Ctx,
		CancelFunc: ctxInfo.CancelFunc,
		DoneCh:     ctxInfo.DoneCh,
	}

	// Store the index in the cache
	c.logger.Debug("storing auto-auth token into the cache")
	err = c.Set(c.baseCtxInfo.Ctx, index)
	if err != nil {
		c.logger.Error("failed to cache the auto-auth token", "error", err)
		return err
	}

	return nil
}

type cacheClearInput struct {
	Type string

	RequestPath   string
	Namespace     string
	Token         string
	TokenAccessor string
	Lease         string
}

func parseCacheClearInput(req *cacheClearRequest) (*cacheClearInput, error) {
	if req == nil {
		return nil, errors.New("nil request options provided")
	}

	if req.Type == "" {
		return nil, errors.New("no type provided")
	}

	in := &cacheClearInput{
		Type:      req.Type,
		Namespace: req.Namespace,
	}

	switch req.Type {
	case "request_path":
		in.RequestPath = req.Value
	case "token":
		in.Token = req.Value
	case "token_accessor":
		in.TokenAccessor = req.Value
	case "lease":
		in.Lease = req.Value
	}

	return in, nil
}

func (c *LeaseCache) hasExpired(currentTime time.Time, index *cachememdb.Index) (bool, error) {
	reader := bufio.NewReader(bytes.NewReader(index.Response))
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize response: %w", err)
	}
	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to parse response as secret: %w", err)
	}

	elapsed := currentTime.Sub(index.LastRenewed)
	var leaseDuration int
	switch {
	case secret.LeaseID != "":
		leaseDuration = secret.LeaseDuration
	case secret.Auth != nil:
		leaseDuration = secret.Auth.LeaseDuration
	default:
		return false, errors.New("secret without lease encountered in expiration check")
	}

	if int(elapsed.Seconds()) > leaseDuration {
		c.logger.Trace("secret has expired", "id", index.ID, "elapsed", elapsed, "lease duration", leaseDuration)
		return true, nil
	}
	return false, nil
}
