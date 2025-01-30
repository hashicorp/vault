// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"bufio"
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

	"github.com/coder/websocket"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cacheboltdb"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/exp/maps"
)

// Example write event (this does not contain all possible fields):
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

// Example event with namespaces for an undelete (this does not contain all possible fields):
// {
//  "id": "6c6b13fd-f133-f351-3cf0-b09ae6a417b1",
//  "source": "vault://hostname",
//  "specversion": "1.0",
//  "type": "*",
//  "data": {
//    "event": {
//      "id": "6c6b13fd-f133-f351-3cf0-b09ae6a417b1",
//      "metadata": {
//        "current_version": "3",
//        "destroyed_versions": "[2,3]",
//        "modified": "true",
//        "oldest_version": "0",
//        "operation": "destroy",
//        "path": "secret-v2/destroy/my-secret"
//      }
//    },
//    "event_type": "kv-v2/destroy",
//    "plugin_info": {
//      "mount_class": "secret",
//      "mount_accessor": "kv_b27b3cad",
//      "mount_path": "secret-v2/",
//      "plugin": "kv",
//      "version": "2"
//    }
//  },
//  "datacontentype": "application/cloudevents",
//  "time": "2024-08-27T12:46:01.373097-04:00"
//}

// StaticSecretCacheUpdater is a struct that utilizes
// the event system to keep the static secret cache up to date.
type StaticSecretCacheUpdater struct {
	client     *api.Client
	leaseCache *LeaseCache
	logger     hclog.Logger
	tokenSink  sink.Sink

	// allowForwardingViaHeaderDisabled is a bool that tracks if
	// allow_forwarding_via_header is disabled on the cluster we're talking to.
	// If we get an error back saying that it's disabled, we'll set this to true
	// and never try to forward again.
	allowForwardingViaHeaderDisabled bool
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
				// If data_path were in every event, we'd get that instead, but unfortunately it isn't.
				path, ok := metadata["path"].(string)
				if !ok {
					return fmt.Errorf("unexpected event format when decoding 'data_path' element, message: %s\nerror: %w", string(message), err)
				}
				namespace, ok := data["namespace"].(string)
				if ok {
					path = namespace + path
				}

				deletedOrDestroyedVersions, newPath := checkForDeleteOrDestroyEvent(messageMap)
				if len(deletedOrDestroyedVersions) > 0 {
					path = newPath
					err = updater.handleDeleteDestroyVersions(path, deletedOrDestroyedVersions)
					if err != nil {
						// While we are kind of 'missing' an event this way, re-calling this function will
						// result in the secret remaining up to date.
						return fmt.Errorf("error handling delete/destroy versions for static secret: path: %q, message: %s error: %w", path, message, err)
					}
				}

				// Note: For delete/destroy events, we continue through to updating the secret itself, too.
				// This means that if the latest version of the secret gets deleted, then the cache keeps
				// knowledge of which the latest version is.
				// One intricacy of e.g. destroyed events is that if the latest secret is destroyed, continuing
				// to update the secret will 404. This is consistent with other behaviour. For Proxy, this means
				// the secret may be evicted. That's okay.

				err = updater.updateStaticSecret(ctx, path)
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

// checkForDeleteOrDestroyEvent checks an event message for delete/destroy events and if there
// are any, returns the versions to be deleted or destroyed, as well as the path to
// If none can be found, returns empty array and empty string.
// We have to do this since events do not always return data_path for all events. If they did,
// we could rely on that instead of doing string manipulation.
// Example return value: [1, 2, 3], "secrets/data/my-secret".
func checkForDeleteOrDestroyEvent(eventMap map[string]interface{}) ([]int, string) {
	var versions []int

	data, ok := eventMap["data"].(map[string]interface{})
	if !ok {
		return versions, ""
	}

	event, ok := data["event"].(map[string]interface{})
	if !ok {
		return versions, ""
	}

	metadata, ok := event["metadata"].(map[string]interface{})
	if !ok {
		return versions, ""
	}

	// We should have only one of these:
	deletedVersions, ok := metadata["deleted_versions"].(string)
	if ok {
		err := json.Unmarshal([]byte(deletedVersions), &versions)
		if err != nil {
			return versions, ""
		}
	}

	destroyedVersions, ok := metadata["destroyed_versions"].(string)
	if ok {
		err := json.Unmarshal([]byte(destroyedVersions), &versions)
		if err != nil {
			return versions, ""
		}
	}

	undeletedVersions, ok := metadata["undeleted_versions"].(string)
	if ok {
		err := json.Unmarshal([]byte(undeletedVersions), &versions)
		if err != nil {
			return versions, ""
		}
	}

	// We have neither deleted_versions nor destroyed_versions, return early
	if len(versions) == 0 {
		return versions, ""
	}

	path, ok := metadata["path"].(string)
	if !ok {
		return versions, ""
	}

	namespace, ok := data["namespace"].(string)
	if ok {
		path = namespace + path
	}

	pluginInfo, ok := data["plugin_info"].(map[string]interface{})
	if !ok {
		return versions, ""
	}

	mountPath := pluginInfo["mount_path"].(string)
	if !ok {
		return versions, ""
	}

	// We get the path without the mount path for safety, just in case the namespace or mount path
	// have 'data' inside.
	namespaceMountPathOnly := namespace + mountPath
	pathWithoutMountPath := strings.TrimPrefix(path, namespaceMountPathOnly)

	// We need to trim destroy or delete to add the correct path for where the secret
	// is stored.
	trimmedPath := strings.TrimPrefix(pathWithoutMountPath, "delete")
	trimmedPath = strings.TrimPrefix(trimmedPath, "destroy")
	trimmedPath = strings.TrimPrefix(trimmedPath, "undelete")

	// This is how we form the ID of the cached secrets
	fixedPath := namespaceMountPathOnly + "data" + trimmedPath

	return versions, fixedPath
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

// handleDeleteDestroyVersions will handle calls to deleteVersions and destroyVersions for a given cached
// secret. The handling is simple: remove them from the cache. We do the same for undeletes, as this will
// also affect the cache, but we don't re-grab the secret for undeletes.
func (updater *StaticSecretCacheUpdater) handleDeleteDestroyVersions(path string, versions []int) error {
	indexId := hashStaticSecretIndex(path)
	// received delete/destroy versions request: path=secret-v2/delete/my-secret
	updater.logger.Debug("received delete/undelete/destroy versions request", "path", path, "indexId", indexId, "versions", versions)

	index, err := updater.leaseCache.db.Get(cachememdb.IndexNameID, indexId)
	if errors.Is(err, cachememdb.ErrCacheItemNotFound) {
		// This event doesn't correspond to a secret in our cache
		// so this is a no-op.
		return nil
	}
	if err != nil {
		return err
	}

	// Hold the lock as we're modifying the secret
	index.IndexLock.Lock()
	defer index.IndexLock.Unlock()

	for _, version := range versions {
		delete(index.Versions, version)
	}

	// Lastly, store the secret
	updater.logger.Debug("storing updated secret as result of delete/undelete/destroy", "path", path, "deletedVersions", versions)
	err = updater.leaseCache.db.Set(index)
	if err != nil {
		return err
	}

	return nil
}

// updateStaticSecret checks for updates for a static secret on the path given,
// and updates the cache if appropriate. For KVv2 secrets, we will also update
// the version at index.Versions[currentVersion] with the same data.
func (updater *StaticSecretCacheUpdater) updateStaticSecret(ctx context.Context, path string) error {
	// We clone the client, as we won't be using the same token.
	client, err := updater.client.Clone()
	if err != nil {
		return err
	}

	// Clear the client's header namespace since we'll be including the
	// namespace as part of the path.
	client.ClearNamespace()

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

		if !updater.allowForwardingViaHeaderDisabled {
			// Set this to always forward to active, since events could come before
			// replication, and if we're connected to the standby, then we will be
			// receiving events from the primary but otherwise getting old values from
			// the standby here. This makes sure that Proxy functions properly
			// even when its Vault address is set to a standby, since we cannot
			// currently receive events from a standby.
			// We only try this if updater.allowForwardingViaHeaderDisabled is false
			// and if we receive an error indicating that the config is set to false,
			// we will never set this header again.
			request.Headers.Set(api.HeaderForward, "active-node")
		}

		resp, err = client.RawRequestWithContext(ctx, request)
		if err != nil {
			if strings.Contains(err.Error(), "forwarding via header X-Vault-Forward disabled") {
				updater.logger.Info("allow_forwarding_via_header disabled, re-attempting update and no longer attempting to forward")
				updater.allowForwardingViaHeaderDisabled = true

				// Try again without the header
				request.Headers.Del(api.HeaderForward)
				resp, err = client.RawRequestWithContext(ctx, request)
			}
		}

		if err != nil {
			updater.logger.Trace("received error when trying to update cache", "path", path, "err", err, "token", token, "namespace", index.Namespace)
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

		index.Response = respBytes.Bytes()
		index.LastRenewed = time.Now().UTC()

		// For KVv2 secrets, let's also update index.Versions[version_of_secret]
		// with the response we received from the current version.
		// Instead of relying on current_version in the event, we should
		// check the message we received, since it's possible the secret
		// got updated between receipt of the event and when we received
		// the request for the secret.
		// First, re-read secret into response so that we can parse it again:
		reader := bufio.NewReader(bytes.NewReader(index.Response))
		resp, err := http.ReadResponse(reader, nil)
		if err != nil {
			// This shouldn't happen, but log just in case it does. There's
			// no real negative consequences of the following function though.
			updater.logger.Warn("failed to deserialize response", "error", err)
		}

		secret, err := api.ParseSecret(resp.Body)
		if err != nil {
			// This shouldn't happen, but log just in case it does. There's
			// no real negative consequences of the following function though.
			updater.logger.Warn("failed to serialize response", "error", err)
		}

		// In case of failures or KVv1 secrets, this function will simply fail silently,
		// which is fine (and expected) since this could be arbitrary JSON.
		updater.leaseCache.addToVersionListForCurrentVersionKVv2Secret(index, secret)

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
	query.Set("namespaces", "*")
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
