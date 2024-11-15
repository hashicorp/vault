package linodego

import (
	"net/http"
	"sync"
	"time"
)

// Client is a wrapper around the Resty client
//
//nolint:unused
type httpClient struct {
	//nolint:unused
	httpClient *http.Client
	//nolint:unused
	userAgent string
	//nolint:unused
	debug bool
	//nolint:unused
	retryConditionals []httpRetryConditional
	//nolint:unused
	retryAfter httpRetryAfter

	//nolint:unused
	pollInterval time.Duration

	//nolint:unused
	baseURL string
	//nolint:unused
	apiVersion string
	//nolint:unused
	apiProto string
	//nolint:unused
	selectedProfile string
	//nolint:unused
	loadedProfile string

	//nolint:unused
	configProfiles map[string]ConfigProfile

	// Fields for caching endpoint responses
	//nolint:unused
	shouldCache bool
	//nolint:unused
	cacheExpiration time.Duration
	//nolint:unused
	cachedEntries map[string]clientCacheEntry
	//nolint:unused
	cachedEntryLock *sync.RWMutex
	//nolint:unused
	logger httpLogger
	//nolint:unused
	onBeforeRequest []func(*http.Request) error
	//nolint:unused
	onAfterResponse []func(*http.Response) error
}
