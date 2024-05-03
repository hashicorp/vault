// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cacheboltdb"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/exp/maps"
	"nhooyr.io/websocket"
)

// Example Event:
//{
//  "id": "a3be9fb1-b514-519f-5b25-b6f144a8c1ce",
//  "source": "https://vaultproject.io/",
//  "specversion": "1.0",
//  "type": "*",
//  "data": {
//    "event": {
//      "id": "a3be9fb1-b514-519f-5b25-b6f144a8c1ce",
//      "metadata": {
//        "current_version": "1",
//        "data_path": "secret/data/foo",
//        "modified": "true",
//        "oldest_version": "0",
//        "operation": "data-write",
//        "path": "secret/data/foo"
//      }
//    },
//    "event_type": "kv-v2/data-write",
//    "plugin_info": {
//      "mount_class": "secret",
//      "mount_accessor": "kv_5dc4d18e",
//      "mount_path": "secret/",
//      "plugin": "kv"
//    }
//  },
//  "datacontentype": "application/cloudevents",
//  "time": "2023-09-12T15:19:49.394915-07:00"
//}

// StaticSecretCacheUpdater is a struct that utilizes
// the event system to keep the static secret cache up to date.
type StaticSecretCacheUpdater struct {
	client     *api.Client
	leaseCache *LeaseCache
	logger     hclog.Logger
	tokenSink  sink.Sink
}

// StaticSecretCacheUpdaterConfig is the configuration for initializing a new
// StaticSecretCacheUpdater.
type StaticSecretCacheUpdaterConfig struct {
	Client     *api.Client
	LeaseCache *LeaseCache
	Logger     hclog.Logger
	// TokenSink is a token sync that will have the latest
	// token from auto-auth in it, to be used in event system
	// connections.
	TokenSink sink.Sink
}

// NewStaticSecretCacheUpdater creates a new instance of a StaticSecretCacheUpdater.
func NewStaticSecretCacheUpdater(conf *StaticSecretCacheUpdaterConfig) (*StaticSecretCacheUpdater, error) {
	if conf == nil {
		return nil, errors.New("nil configuration provided")
	}

	if conf.LeaseCache == nil {
		return nil, fmt.Errorf("nil Lease Cache (a required parameter): %v", conf)
	}

	if conf.Logger == nil {
		return nil, fmt.Errorf("nil Logger (a required parameter): %v", conf)
	}

	if conf.Client == nil {
		return nil, fmt.Errorf("nil API client (a required parameter): %v", conf)
	}

	if conf.TokenSink == nil {
		return nil, fmt.Errorf("nil token sink (a required parameter): %v", conf)
	}

	return &StaticSecretCacheUpdater{
		client:     conf.Client,
		leaseCache: conf.LeaseCache,
		logger:     conf.Logger,
		tokenSink:  conf.TokenSink,
	}, nil
}

// streamStaticSecretEvents streams static secret events and updates
// the cache when updates are notified. This method will return errors in cases
// of failed updates, malformed events, and other.
// For best results, the caller of this function should retry on error with backoff,
// if it is desired for the cache to always remain up to date.
func (updater *StaticSecretCacheUpdater) streamStaticSecretEvents(ctx context.Context) error {
	// First, ensure our token is up-to-date:
	updater.client.SetToken(updater.tokenSink.(sink.SinkReader).Token())
	conn, err := updater.openWebSocketConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	err = updater.preEventStreamUpdate(ctx)
	if err != nil {
		return fmt.Errorf("error when performing pre-event stream secret update: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, message, err := conn.Read(ctx)
			if err != nil {
				// The caller of this function should make the decision on if to retry. If it does, then
				// the websocket connection will be retried, and we will check for missed events.
				return fmt.Errorf("error when attempting to read from event stream, reopening websocket: %w", err)
			}
			updater.logger.Trace("received event", "message", string(message))
			messageMap := make(map[string]interface{})
			err = json.Unmarshal(message, &messageMap)
			if err != nil {
				return fmt.Errorf("error when unmarshaling event, message: %s\nerror: %w", string(message), err)
			}
			data, ok := messageMap["data"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("unexpected event format when decoding 'data' element, message: %s\nerror: %w", string(message), err)
			}
			event, ok := data["event"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("unexpected event format when decoding 'event' element, message: %s\nerror: %w", string(message), err)
			}
			metadata, ok := event["metadata"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("unexpected event format when decoding 'metadata' element, message: %s\nerror: %w", string(message), err)
			}
			modified, ok := metadata["modified"].(string)
			if ok && modified == "true" {
				path, ok := metadata["path"].(string)
				if !ok {
					return fmt.Errorf("unexpected event format when decoding 'path' element, message: %s\nerror: %w", string(message), err)
				}
				err := updater.updateStaticSecret(ctx, path)
				if err != nil {
					// While we are kind of 'missing' an event this way, re-calling this function will
					// result in the secret remaining up to date.
					return fmt.Errorf("error updating static secret: path: %q, message: %s error: %w", path, message, err)
				}
			} else {
				// This is an event we're not interested in, ignore it and
				// carry on.
				continue
			}
		}
	}

	return nil
}

// preEventStreamUpdate is called after successful connection to the event system, but before
// we process any events, to ensure we don't miss any updates.
// In some cases, this will result in multiple processing of the same updates, but
// this ensures that we don't lose any updates to secrets that might have been sent
// while the connection is forming.
func (updater *StaticSecretCacheUpdater) preEventStreamUpdate(ctx context.Context) error {
	indexes, err := updater.leaseCache.db.GetByPrefix(cachememdb.IndexNameID)
	if err != nil {
		return err
	}

	updater.logger.Debug("starting pre-event stream update of static secrets")

	var errs *multierror.Error
	for _, index := range indexes {
		if index.Type != cacheboltdb.StaticSecretType {
			continue
		}
		err = updater.updateStaticSecret(ctx, index.RequestPath)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	updater.logger.Debug("finished pre-event stream update of static secrets")

	return errs.ErrorOrNil()
}

// updateStaticSecret checks for updates for a static secret on the path given,
// and updates the cache if appropriate
func (updater *StaticSecretCacheUpdater) updateStaticSecret(ctx context.Context, path string) error {
	// We clone the client, as we won't be using the same token.
	client, err := updater.client.Clone()
	if err != nil {
		return err
	}

	indexId := hashStaticSecretIndex(path)

	updater.logger.Debug("received update static secret request", "path", path, "indexId", indexId)

	index, err := updater.leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
		// This event doesn't correspond to a secret in our cache
		// so this is a no-op.
		return nil
	}
	if err != nil {
		return err
	}

	// We use a raw request so that we can store all the
	// request information, just like we do in the Proxier Send methods.
	request := client.NewRequest(http.MethodGet, "/v1/"+path)
	if request.Headers == nil {
		request.Headers = make(http.Header)
	}
	request.Headers.Set("User-Agent", useragent.ProxyString())

	var resp *api.Response
	var tokensToRemove []string
	var successfulAttempt bool
	for _, token := range maps.Keys(index.Tokens) {
		client.SetToken(token)
		request.Headers.Set(api.AuthHeaderName, token)
		resp, err = client.RawRequestWithContext(ctx, request)
		if err != nil {
			updater.logger.Trace("received error when trying to update cache", "path", path, "err", err, "token", token)
			// We cannot access this secret with this token for whatever reason,
			// so token for removal.
			tokensToRemove = append(tokensToRemove, token)
			continue
		} else {
			// We got our updated secret!
			successfulAttempt = true
			break
		}
	}

	if successfulAttempt {
		// We need to update the index, so first, hold the lock.
		index.IndexLock.Lock()
		defer index.IndexLock.Unlock()

		// First, remove the tokens we noted couldn't access the secret from the token index
		for _, token := range tokensToRemove {
			delete(index.Tokens, token)
		}

		sendResponse, err := NewSendResponse(resp, nil)
		if err != nil {
			return err
		}

		// Serialize the response to store it in the cached index
		var respBytes bytes.Buffer
		err = sendResponse.Response.Write(&respBytes)
		if err != nil {
			updater.logger.Error("failed to serialize response", "error", err)
			return err
		}

		// Set the index's Response
		index.Response = respBytes.Bytes()
		index.LastRenewed = time.Now().UTC()

		// Lastly, store the secret
		updater.logger.Debug("storing response into the cache due to update", "path", path)
		err = updater.leaseCache.db.Set(index)
		if err != nil {
			return err
		}
	} else {
		// No token could successfully update the secret, or secret was deleted.
		// We should evict the cache instead of re-storing the secret.
		updater.logger.Debug("evicting response from cache", "path", path)
		err = updater.leaseCache.db.Evict(cachememdb.IndexNameID, indexId)
		if err != nil {
			return err
		}
	}

	return nil
}

// openWebSocketConnection opens a websocket connection to the event system for
// the events that the static secret cache updater is interested in.
func (updater *StaticSecretCacheUpdater) openWebSocketConnection(ctx context.Context) (*websocket.Conn, error) {
	// We parse this into a URL object to get the specific host and scheme
	// information without nasty string parsing.
	vaultURL, err := url.Parse(updater.client.Address())
	if err != nil {
		return nil, err
	}
	vaultHost := vaultURL.Host
	// If we're using https, use wss, otherwise ws
	scheme := "wss"
	if vaultURL.Scheme == "http" {
		scheme = "ws"
	}

	webSocketURL := url.URL{
		Path:   "/v1/sys/events/subscribe/kv*",
		Host:   vaultHost,
		Scheme: scheme,
	}
	query := webSocketURL.Query()
	query.Set("json", "true")
	webSocketURL.RawQuery = query.Encode()

	updater.client.AddHeader(api.AuthHeaderName, updater.client.Token())
	updater.client.AddHeader(api.NamespaceHeaderName, updater.client.Namespace())

	// Populate these now to avoid recreating them in the upcoming for loop.
	headers := updater.client.Headers()
	wsURL := webSocketURL.String()
	httpClient := updater.client.CloneConfig().HttpClient

	// We do ten attempts, to ensure we follow forwarding to the leader.
	var conn *websocket.Conn
	var resp *http.Response
	for attempt := 0; attempt < 10; attempt++ {
		conn, resp, err = websocket.Dial(ctx, wsURL, &websocket.DialOptions{
			HTTPClient: httpClient,
			HTTPHeader: headers,
		})
		if err == nil {
			break
		}

		switch {
		case resp == nil:
			break
		case resp.StatusCode == http.StatusTemporaryRedirect:
			wsURL = resp.Header.Get("Location")
			continue
		default:
			break
		}
	}

	if err != nil {
		errMessage := err.Error()
		if resp != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil, fmt.Errorf("received 404 when opening web socket to %s, ensure Vault is Enterprise version 1.16 or above", wsURL)
			}
			if resp.StatusCode == http.StatusForbidden {
				var errBytes []byte
				errBytes, err = io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					return nil, fmt.Errorf("error occured when attempting to read error response from Vault server")
				}
				errMessage = string(errBytes)
			}
		}
		return nil, fmt.Errorf("error returned when opening event stream web socket to %s, ensure auto-auth token"+
			" has correct permissions and Vault is Enterprise version 1.16 or above: %s", wsURL, errMessage)
	}

	if conn == nil {
		return nil, errors.New(fmt.Sprintf("too many redirects as part of establishing web socket connection to %s", wsURL))
	}

	return conn, nil
}

// Run is intended to be the method called by Vault Proxy, that runs the subsystem.
// Once a token is provided to the sink, we will start the websocket and start consuming
// events and updating secrets.
// Run will shut down gracefully when the context is cancelled.
func (updater *StaticSecretCacheUpdater) Run(ctx context.Context, authRenewalInProgress *atomic.Bool, invalidTokenErrCh chan error) error {
	updater.logger.Info("starting static secret cache updater subsystem")
	defer func() {
		updater.logger.Info("static secret cache updater subsystem stopped")
	}()

tokenLoop:
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Wait for the auto-auth token to be populated...
			if updater.tokenSink.(sink.SinkReader).Token() != "" {
				break tokenLoop
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	shouldBackoff := false
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// If we're erroring and the context isn't done, we should add
			// a little backoff to make sure we don't accidentally overload
			// Vault or similar.
			if shouldBackoff {
				time.Sleep(10 * time.Second)
			}
			err := updater.streamStaticSecretEvents(ctx)
			if err != nil {
				updater.logger.Error("error occurred during streaming static secret cache update events", "err", err)
				shouldBackoff = true
				if strings.Contains(err.Error(), logical.ErrInvalidToken.Error()) && !authRenewalInProgress.Load() {
					// Drain the channel in case there is an error that has already been sent but not received
					select {
					case <-invalidTokenErrCh:
					default:
					}
					updater.logger.Error("received invalid token error while opening websocket")
					invalidTokenErrCh <- err
				}
				continue
			}
		}
	}
}
