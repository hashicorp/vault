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
	"github.com/hashicorp/vault/helper/useragent"
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
}

// LeaseCacheConfig is the configuration for initializing a new
// LeaseCache.
type LeaseCacheConfig struct {
	Client             *api.Client
	BaseContext        context.Context
	Proxier            Proxier
	Logger             hclog.Logger
	Storage            *cacheboltdb.BoltStorage
	CacheStaticSecrets bool
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

	db, err := cachememdb.New()
	if err != nil {
		return nil, err
	}

	// Create a base context for the lease cache layer
	baseCtxInfo := cachememdb.NewContextInfo(conf.BaseContext)

	return &LeaseCache{
		client:             conf.Client,
		proxier:            conf.Proxier,
		logger:             conf.Logger,
		db:                 db,
		baseCtxInfo:        baseCtxInfo,
		l:                  &sync.RWMutex{},
		idLocks:            locksutil.CreateLocks(),
		inflightCache:      gocache.New(gocache.NoExpiration, gocache.NoExpiration),
		ps:                 conf.Storage,
		cacheStaticSecrets: conf.CacheStaticSecrets,
	}, nil
}

// SetShuttingDown is a setter for the shuttingDown field
func (c *LeaseCache) SetShuttingDown(in bool) {
	c.shuttingDown.Store(in)
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
	return c.checkCacheForRequest(id, nil)
}

// checkCacheForStaticSecretRequest checks the cache for a particular request based on its
// computed ID. It returns a non-nil *SendResponse if an entry is found.
// If a request is provided, it will validate that the token is allowed to retrieve this
// cache entry, and return nil if it isn't. It will also evict the cache if this is a non-GET
// request.
func (c *LeaseCache) checkCacheForStaticSecretRequest(id string, req *SendRequest) (*SendResponse, error) {
	return c.checkCacheForRequest(id, req)
}

// checkCacheForRequest checks the cache for a particular request based on its
// computed ID. It returns a non-nil *SendResponse if an entry is found.
// If a token is provided, it will validate that the token is allowed to retrieve this
// cache entry, and return nil if it isn't.
func (c *LeaseCache) checkCacheForRequest(id string, req *SendRequest) (*SendResponse, error) {
	var token string
	if req != nil {
		token = req.Token
		if req.Request.Method != http.MethodGet {
			// This isn't a GET, so we should short-circuit and invalidate the cache
			// as we know the cache is now stale.
			c.logger.Debug("evicting index from cache, as non-GET received", "id", id, "method", req.Request.Method, "path", req.Request.URL.Path)
			err := c.db.Evict(cachememdb.IndexNameID, id)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}

	index, err := c.db.Get(cachememdb.IndexNameID, id)
	if err != nil {
		return nil, err
	}

	if index == nil {
		return nil, nil
	}

	if token != "" {
		// This is a static secret check. We need to ensure that this token
		// has previously demonstrated access to this static secret.
		// We could check the capabilities cache here, but since these
		// indexes should be in sync, this saves us an extra cache get.
		if _, ok := index.Tokens[token]; !ok {
			// We don't have access to this static secret, so
			// we do not return the cached response.
			return nil, nil
		}
	}

	// Cached request is found, deserialize the response
	reader := bufio.NewReader(bytes.NewReader(index.Response))
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		c.logger.Error("failed to deserialize response", "error", err)
		return nil, err
	}

	sendResp, err := NewSendResponse(&api.Response{Response: resp}, index.Response)
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
			c.inflightCache.Delete(staticSecretCacheId)
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

	// Check if the response for this request is already in the dynamic secret cache
	cachedResp, err := c.checkCacheForDynamicSecretRequest(dynamicSecretCacheId)
	if err != nil {
		return nil, err
	}
	if cachedResp != nil {
		c.logger.Debug("returning cached response", "path", req.Request.URL.Path)
		return cachedResp, nil
	}

	// Check if the response for this request is already in the static secret cache
	cachedResp, err = c.checkCacheForStaticSecretRequest(staticSecretCacheId, req)
	if err != nil {
		return nil, err
	}
	if cachedResp != nil {
		c.logger.Debug("returning cached response", "id", staticSecretCacheId, "path", req.Request.URL.Path)
		return cachedResp, nil
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

	// TODO: if secret.MountType == "kvv1" || secret.MountType == "kvv2"
	if c.cacheStaticSecrets && secret != nil {
		index.Type = cacheboltdb.StaticSecretType
		index.ID = staticSecretCacheId
		err := c.cacheStaticSecret(ctx, req, resp, index)
		if err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		// Since it's not a static secret, set the ID to be the dynamic id
		index.ID = dynamicSecretCacheId
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
		if err != nil {
			return nil, err
		}
		// If the lease belongs to a token that is not managed by the agent,
		// return the response without caching it.
		if entry == nil {
			c.logger.Debug("pass-through lease response; token not managed by agent", "method", req.Request.Method, "path", req.Request.URL.Path)
			return resp, nil
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
			if err != nil {
				return nil, err
			}
			// If parent token is not managed by the agent, child shouldn't be
			// either.
			if entry == nil {
				c.logger.Debug("pass-through auth response; parent token not managed by agent", "method", req.Request.Method, "path", req.Request.URL.Path)
				return resp, nil
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
		c.logger.Debug("storing response into the cache", "method", req.Request.Method, "path", req.Request.URL.Path)
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

func (c *LeaseCache) cacheStaticSecret(ctx context.Context, req *SendRequest, resp *SendResponse, index *cachememdb.Index) error {
	// We must hold a lock for the index while it's being updated.
	// We prepend "index/" to this lock so that it's distinct to the lock
	// being held for inflight requests.
	// We keep the two locking mechanisms distinct, so that it's only writes
	// that have to be serial.
	lock := locksutil.LockForKey(c.idLocks, "index/"+index.ID)
	lock.Lock()
	defer lock.Unlock()

	// If a cached version of this secret exists, we now have access, so
	// we don't need to re-cache, just update index.Tokens
	indexFromCache, err := c.db.Get(cachememdb.IndexNameID, index.ID)
	if err != nil {
		return err
	}

	// The index already exists, so all we need to do is add our token
	// to the index's allowed token list, then re-store it
	if indexFromCache != nil {
		indexFromCache.Tokens[req.Token] = struct{}{}

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

	// Set the index's Response
	index.Response = respBytes.Bytes()

	// Initialize the token map and add this token to it.
	index.Tokens = map[string]struct{}{req.Token: {}}

	// Set the index type
	index.Type = cacheboltdb.StaticSecretType

	return c.storeStaticSecretIndex(ctx, req, index)
}

func (c *LeaseCache) storeStaticSecretIndex(ctx context.Context, req *SendRequest, index *cachememdb.Index) error {
	// Store the index in the cache
	c.logger.Debug("storing response into the cache", "method", req.Request.Method, "path", req.Request.URL.Path)
	err := c.Set(ctx, index)
	if err != nil {
		c.logger.Error("failed to cache the proxied response", "error", err)
		return err
	}

	capabilitiesIndex, err := c.retrieveOrCreateTokenCapabilitiesEntry(req.Token)
	if err != nil {
		c.logger.Error("failed to cache the proxied response", "error", err)
		return err
	}

	// /sys/capabilities accepts both requests that look like foo/bar
	// and /foo/bar but not /v1/foo/bar.
	// We trim the /v1 from the start of the URL to get the /foo/bar form.
	path := strings.TrimPrefix(req.Request.URL.Path, "/v1")

	// Extra caution -- avoid potential nil
	if capabilitiesIndex.Capabilities == nil {
		capabilitiesIndex.Capabilities = make(map[string]struct{})
	}

	// update the index with the new capability:
	capabilitiesIndex.Capabilities[path] = struct{}{}

	err = c.Set(ctx, capabilitiesIndex)
	if err != nil {
		c.logger.Error("failed to cache the proxied response", "error", err)
		return err
	}

	return nil
}

// retrieveOrCreateTokenCapabilitiesEntry will either retrieve the token
// capabilities entry from the cache, or create a new, empty one.
func (c *LeaseCache) retrieveOrCreateTokenCapabilitiesEntry(token string) (*cachememdb.Index, error) {
	// The index ID is a hash of the token.
	indexId := hex.EncodeToString(cryptoutil.Blake2b256Hash(token))
	indexFromCache, err := c.db.Get(cachememdb.IndexNameID, indexId)
	if err != nil {
		return nil, err
	}

	if indexFromCache != nil {
		return indexFromCache, nil
	}

	// Build the index to cache based on the response received
	index := &cachememdb.Index{
		ID:           indexId,
		Token:        token,
		Type:         cacheboltdb.TokenCapabilitiesType,
		Capabilities: make(map[string]struct{}),
	}

	return index, nil
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

	// We do not preserve the initial User-Agent here (i.e. use
	// AgentProxyStringWithProxiedUserAgent) since these requests are from
	// the proxy subsystem, but are made by Agent's lifetime watcher,
	// not triggered by a specific request.
	headers.Set("User-Agent", useragent.AgentProxyString())
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
	if err != nil {
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

// computeStaticSecretCacheIndex results in a value that uniquely identifies a static
// secret's cached ID. Notably, we intentionally ignore headers (for example,
// the X-Vault-Token header) to remain agnostic to which token is being
// used in the request. We care only about the path.
func computeStaticSecretCacheIndex(req *SendRequest) string {
	// /sys/capabilities accepts both requests that look like foo/bar
	// and /foo/bar but not /v1/foo/bar.
	// We trim the /v1 from the start of the URL to get the /foo/bar form.
	// This means that we can use the paths we retrieve from the
	// /sys/capabilities endpoint to access this index
	// without having to re-add the /v1
	path := strings.TrimPrefix(req.Request.URL.Path, "/v1")
	return hex.EncodeToString(cryptoutil.Blake2b256Hash(path))
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
			if err == errInvalidType {
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
			if index.RenewCtxInfo != nil {
				if index.RenewCtxInfo.CancelFunc != nil {
					index.RenewCtxInfo.CancelFunc()
				}
			}
		}

	case "token":
		if in.Token == "" {
			return errors.New("token not provided")
		}

		// Get the context for the given token and cancel its context
		index, err := c.db.Get(cachememdb.IndexNameToken, in.Token)
		if err != nil {
			return err
		}
		if index == nil {
			return nil
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
		if err != nil {
			return err
		}
		if index == nil {
			return nil
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
		if err != nil {
			return err
		}
		if index == nil {
			return nil
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
		if err != nil {
			return false, err
		}
		if index == nil {
			return true, nil
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
		if err != nil {
			return err
		}

		if entry == nil {
			return fmt.Errorf("could not find parent Token %s for req path %s", index.RequestToken, index.RequestPath)
		}

		// Derive a context for renewal using the token's context
		renewCtxInfo = cachememdb.NewContextInfo(entry.RenewCtxInfo.Ctx)

	case secret.Auth != nil:
		var parentCtx context.Context
		if !secret.Auth.Orphan {
			entry, err := c.db.Get(cachememdb.IndexNameToken, index.RequestToken)
			if err != nil {
				return err
			}
			// If parent token is not managed by the agent, child shouldn't be
			// either.
			if entry == nil {
				return fmt.Errorf("could not find parent Token %s for req path %s", index.RequestToken, index.RequestPath)
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
	if err != nil {
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
