// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package watch

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/consul-template/config"
	dep "github.com/hashicorp/consul-template/dependency"
	"github.com/hashicorp/vault/api"
)

// VaultTokenWatcher monitors the vault token for updates
func VaultTokenWatcher(
	clients *dep.ClientSet, c *config.VaultConfig, doneCh chan struct{},
) (*Watcher, error) {
	// c.Vault.Token is populated by the config code from all places
	// vault tokens are supported. So if there is no token set here,
	// tokens are not being used.
	raw_token := strings.TrimSpace(config.StringVal(c.Token))
	if raw_token == "" {
		return nil, nil
	}

	unwrap := config.BoolVal(c.UnwrapToken)
	vault := clients.Vault()
	// get/set token once when kicked off, async after that..
	token, err := unpackToken(vault, raw_token, unwrap)
	if err != nil {
		return nil, fmt.Errorf("vaultwatcher: %w", err)
	}
	vault.SetToken(token)

	var once sync.Once
	var watcher *Watcher
	getWatcher := func() *Watcher {
		once.Do(func() {
			watcher = NewWatcher(&NewWatcherInput{
				Clients:        clients,
				RetryFuncVault: RetryFunc(c.Retry.RetryFunc()),
			})
		})
		return watcher
	}

	// Vault Agent Token File process //
	tokenFile := strings.TrimSpace(config.StringVal(c.VaultAgentTokenFile))
	if tokenFile != "" {
		w := getWatcher()
		watchLoop, err := watchTokenFile(w, tokenFile, raw_token, unwrap, doneCh)
		if err != nil {
			return nil, fmt.Errorf("vaultwatcher: %w", err)
		}
		go watchLoop()
	}

	// Vault Token Renewal process //
	renewVault := vault.Token() != "" && config.BoolVal(c.RenewToken)
	if renewVault {
		w := getWatcher()
		vt, err := dep.NewVaultTokenQuery(token)
		if err != nil {
			w.Stop()
			return nil, fmt.Errorf("vaultwatcher: %w", err)
		}
		if _, err := w.Add(vt); err != nil {
			w.Stop()
			return nil, fmt.Errorf("vaultwatcher: %w", err)
		}
	}

	return watcher, nil
}

func watchTokenFile(
	w *Watcher, tokenFile, raw_token string, unwrap bool, doneCh chan struct{},
) (func(), error) {
	// watcher, tokenFile, raw_token, unwrap, doneCh
	atf, err := dep.NewVaultAgentTokenQuery(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("vaultwatcher: %w", err)
	}
	if _, err := w.Add(atf); err != nil {
		w.Stop()
		return nil, fmt.Errorf("vaultwatcher: %w", err)
	}
	vault := w.clients.Vault()
	return func() {
		for {
			select {
			case v := <-w.DataCh():
				new_raw_token := strings.TrimSpace(v.Data().(string))
				if new_raw_token == raw_token {
					break
				}
				token, err := unpackToken(vault, new_raw_token, unwrap)
				switch err {
				case nil:
					raw_token = new_raw_token
					vault.SetToken(token)
				default:
					log.Printf("[INFO] %s", err)
				}
			case <-doneCh:
				return
			}
		}
	}, nil
}

type vaultClient interface {
	SetToken(string)
	Logical() *api.Logical
}

// unpackToken grabs the real token from raw_token (unwrap, etc)
func unpackToken(client vaultClient, token string, unwrap bool) (string, error) {
	// If vault agent specifies wrap_ttl for the token it is returned as
	// a SecretWrapInfo struct marshalled into JSON instead of the normal raw
	// token. This checks for that and pulls out the token if it is the case.
	var wrapinfo api.SecretWrapInfo
	if err := json.Unmarshal([]byte(token), &wrapinfo); err == nil {
		token = wrapinfo.Token
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	if unwrap {
		client.SetToken(token) // needs to be set to unwrap
		secret, err := client.Logical().Unwrap(token)
		switch {
		case err != nil:
			return token, fmt.Errorf("vault unwrap: %s", err)
		case secret == nil:
			return token, fmt.Errorf("vault unwrap: no secret")
		case secret.Auth == nil:
			return token, fmt.Errorf("vault unwrap: no secret auth")
		case secret.Auth.ClientToken == "":
			return token, fmt.Errorf("vault unwrap: no token returned")
		default:
			token = secret.Auth.ClientToken
		}
	}
	return token, nil
}
