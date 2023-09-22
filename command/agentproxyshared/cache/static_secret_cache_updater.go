// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/vault/helper/useragent"

	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/api"

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

const StaticSecretBackoff = 10 * time.Second

// StaticSecretCacheUpdater is a struct that utilizes
// the event system to keep the static secret cache up to date.
type StaticSecretCacheUpdater struct {
	client     *api.Client
	leaseCache *LeaseCache
	logger     hclog.Logger
}

// StaticSecretCacheUpdaterConfig is the configuration for initializing a new
// StaticSecretCacheUpdater.
type StaticSecretCacheUpdaterConfig struct {
	Client     *api.Client
	LeaseCache *LeaseCache
	Logger     hclog.Logger
}

// NewStaticSecretCacheUpdater creates a new instance of a StaticSecretCacheUpdater.
func NewStaticSecretCacheUpdater(conf *StaticSecretCacheUpdaterConfig) (*StaticSecretCacheUpdater, error) {
	if conf == nil {
		return nil, errors.New("nil configuration provided")
	}

	if conf.LeaseCache == nil || conf.Logger == nil {
		return nil, fmt.Errorf("missing configuration required params: %v", conf)
	}

	if conf.Client == nil {
		return nil, fmt.Errorf("nil API client")
	}

	return &StaticSecretCacheUpdater{
		client:     conf.Client,
		leaseCache: conf.LeaseCache,
		logger:     conf.Logger,
	}, nil
}

// streamStaticSecretEvents streams static secret events and updates
// the cache when updates are notified. This method will return errors in cases
// of failed updates, malformed events, and other.
// For best results, the caller of this function should retry on error with backoff,
// if it is desired for the cache to always remain up to date.
func (updater *StaticSecretCacheUpdater) streamStaticSecretEvents(ctx context.Context) error {
	var conn *websocket.Conn
	for {
		var err error
		conn, err = updater.openWebSocketConnection(ctx)
		if err != nil {
			updater.logger.Error("error when opening event stream:", err)

			// We backoff in case of Vault downtime etc
			time.Sleep(StaticSecretBackoff)
			continue
		} else {
			break
		}
	}

	defer conn.Close(websocket.StatusNormalClosure, "")

	// before we check for events, update all of our cached
	// kv secrets, in case we missed any events
	// TODO: to be implemented in a future PR

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			// The caller of this function should make the decision on if to retry. If it does, then
			// the websocket connection will be retried, and we will check for missed events.
			return fmt.Errorf("error when attempting to read from event stream, reopening websocket: %w", err)
		}
		messageMap := make(map[string]interface{})
		err = json.Unmarshal(message, &messageMap)
		if err != nil {
			return fmt.Errorf("error when unmarshaling event, message: %s\nerror: %w", string(message), err)
		}
		data, ok := messageMap["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected event format, message: %s\nerror: %w", string(message), err)
		}
		event, ok := data["event"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected event format, message: %s\nerror: %w", string(message), err)
		}
		metadata, ok := event["metadata"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected event format, message: %s\nerror: %w", string(message), err)
		}
		modified, ok := metadata["modified"].(bool)
		if ok && modified {
			path, ok := metadata["path"].(string)
			if !ok {
				return fmt.Errorf("unexpected event format, message: %s\nerror: %w", string(message), err)
			}
			err := updater.updateStaticSecret(ctx, path)
			if err != nil {
				// While we are kind of 'missing' an event this way, re-calling this function will
				// result in the secret remaining up to date.
				return fmt.Errorf("unexpected event format, message: %s\nerror: %w", string(message), err)
			}
		} else {
			// This is an event we're not interested in, ignore it and
			// carry on.
			continue
		}
	}

	return nil
}

func (updater *StaticSecretCacheUpdater) updateStaticSecret(ctx context.Context, path string) error {
	// We clone the client, as we won't be using the same token.
	client, err := updater.client.Clone()
	if err != nil {
		return err
	}

	// TODO: avoid using req here make a new method
	// that takes path and returns index
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: path,
			},
		},
	}
	index, err := updater.leaseCache.db.Get(cachememdb.IndexNameID, computeStaticSecretCacheIndex(req))
	if err != nil {
		return err
	}
	if index == nil {
		// This event doesn't correspond to a secret in our cache
		// so this is a no-op.
		return nil
	}

	// We use a raw request so that we can store all the
	// request information, just like we do in the Proxier Send methods.
	request := client.NewRequest(http.MethodGet, path)
	if request.Headers == nil {
		request.Headers = make(http.Header)
	}
	request.Headers.Set("User-Agent", useragent.ProxyString())

	var resp *api.Response
	var tokensToRemove []string
	for _, token := range index.Tokens {
		client.SetToken(token)
		resp, err = client.RawRequestWithContext(ctx, request)
		if err != nil {
			// We cannot access this secret with this token for whatever reason,
			// so token for removal.
			tokensToRemove = append(tokensToRemove, token)
			continue
		} else {
			// We got our updated secret!
			break
		}
	}

	if resp != nil {
		// We need to update the index, so first, hold the lock.
		index.IndexLock.Lock()
		defer index.IndexLock.Unlock()

		// First, remove the tokens we noted couldn't access the secret from the token index
		for _, token := range tokensToRemove {
			// TODO, fix this once we get the map in this branch
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

		// Lastly, store the secret
		updater.logger.Debug("storing response into the cache due to event update", "method", req.Request.Method, "path", req.Request.URL.Path)
		err = updater.leaseCache.db.Set(index)
		if err != nil {
			return err
		}
	} else {
		// No token could successfully update the secret, or secret was deleted.
		// We should evict the cache instead of re-storing the secret.
		err = updater.leaseCache.db.Evict(cachememdb.IndexNameID, computeStaticSecretCacheIndex(req))
		if err != nil {
			return err
		}
	}

	return nil
}

func (updater *StaticSecretCacheUpdater) openWebSocketConnection(ctx context.Context) (*websocket.Conn, error) {
	wsLocation := fmt.Sprintf("ws://%s/v1/sys/events/subscribe/kv*?json=true", updater.client.Address())
	updater.client.AddHeader("X-Vault-Token", updater.client.Token())
	updater.client.AddHeader("X-Vault-Namespace", updater.client.Namespace())
	conn, _, err := websocket.Dial(ctx, wsLocation, &websocket.DialOptions{
		HTTPClient: updater.client.CloneConfig().HttpClient,
		HTTPHeader: updater.client.Headers(),
	})
	if err != nil {
		return nil, fmt.Errorf("error returned when opening event stream web socket, ensure auto-auth token"+
			" has correct permissions and Vault is version 1.16 or above: %w", err)
	}
	return conn, nil
}
