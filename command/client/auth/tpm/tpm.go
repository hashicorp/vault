// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package tpm

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/client/auth"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
)

type tpmMethod struct {
	logger    hclog.Logger
	mountPath string
	roleName  string

	tpmCertDestination string
	tpmCertDevice      string
	reload             bool

	client *api.Client

	stopCh          chan struct{}
	doneCh          chan struct{}
	credSuccessGate chan struct{}
	ticker          *time.Ticker
	once            *sync.Once
	credsFound      chan struct{}
	latestHash      []byte
	tlsCfg          *api.TLSConfig
}

var _ auth.AuthMethodWithClient = &tpmMethod{}

// NewTPMAuthMethod creates a new TPM auth method handler.  Besides the
// config that the user can specify in the autoauth config stanza, we
// also depend on configuration derived from the tpm_cert stanza:
// role_name, tmp_cert_destination, and (optionally) tpm_cert_device.
func NewTPMAuthMethod(conf *auth.AuthConfig, tlsCfg *api.TLSConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}

	if conf.Config == nil {
		return nil, errors.New("missing config")
	}

	t := &tpmMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,

		stopCh:          make(chan struct{}),
		doneCh:          make(chan struct{}),
		credSuccessGate: make(chan struct{}),
		once:            new(sync.Once),
		credsFound:      make(chan struct{}),
		tlsCfg:          tlsCfg,
	}

	keys := slices.Collect(maps.Keys(conf.Config))
	conf.Logger.Debug("TPM auth method config received", "config_keys", keys)

	roleNameRaw, ok := conf.Config["role_name"]
	if !ok {
		return nil, fmt.Errorf("missing required 'role_name' config value")
	}
	t.roleName, ok = roleNameRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role_name' config value to string")
	}
	if t.roleName == "" {
		return nil, errors.New("'role_name' config value cannot be empty")
	}

	tpmCertDestRaw, ok := conf.Config["tpm_cert_destination"]
	if !ok {
		return nil, errors.New("missing required 'tpm_cert_destination' config value")
	}
	t.tpmCertDestination, ok = tpmCertDestRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'tpm_cert_destination' config value to string")
	}
	if t.tpmCertDestination == "" {
		return nil, errors.New("'tpm_cert_destination' config value cannot be empty")
	}

	tpmCertDevRaw, ok := conf.Config["tpm_cert_device"]
	if ok {
		t.tpmCertDevice, ok = tpmCertDevRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'tpm_cert_device' config value to string")
		}
	}

	reload, ok := conf.Config["reload"]
	if ok {
		t.reload, ok = reload.(bool)
		if !ok {
			return nil, errors.New("could not convert 'reload' config value to bool")
		}
	}
	conf.Logger.Debug("TPM auth method initialized", "role_name", t.roleName, "mount_path", t.mountPath,
		"tpm_cert_destination", t.tpmCertDestination)

	if t.reload {
		reloadPeriod := time.Minute
		if reloadPeriodRaw, ok := conf.Config["reload_period"]; ok {
			period, err := parseutil.ParseDurationSecond(reloadPeriodRaw)
			if err != nil {
				return nil, fmt.Errorf("error parsing 'reload_period' value: %w", err)
			}
			reloadPeriod = period
		}
		t.ticker = time.NewTicker(reloadPeriod)

		go t.runWatcher()
	}

	return t, nil
}

func (t *tpmMethod) WaitForCreds(ctx context.Context) error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		result, err := sdktpm.LoadFromDirectory(t.tpmCertDestination)
		if err == nil {
			err = result.Validate(sdktpm.TPMConfig{DevicePath: t.tpmCertDevice})
			if err == nil {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("no valid TPM credentials found: %w", err)
		case <-t.stopCh:
			return fmt.Errorf("no valid TPM credentials found: %w", err)
		case <-ticker.C:
		}
	}
}

func (t *tpmMethod) Authenticate(_ context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	t.logger.Trace("beginning authentication")

	authMap := map[string]interface{}{
		"role_name": t.roleName,
	}

	t.logger.Debug("TPM authenticate", "mount_path", t.mountPath, "auth_map", authMap)

	return fmt.Sprintf("%s/login", t.mountPath), nil, authMap, nil
}

func (t *tpmMethod) NewCreds() chan struct{} {
	return t.credsFound
}

func (t *tpmMethod) CredSuccess() {
	t.once.Do(func() {
		close(t.credSuccessGate)
	})
}

func (t *tpmMethod) Shutdown() {
	if t.reload {
		t.ticker.Stop()
		close(t.stopCh)
		<-t.doneCh
	}
}

// AuthClient uses the existing client's address and returns a new client with
// the auto-auth method's certificate information if that's provided in its
// config map.
func (t *tpmMethod) AuthClient(client *api.Client) (*api.Client, error) {
	t.logger.Trace("deriving auth client to use")

	clientToAuth := client

	// Return cached client if present
	if t.client != nil && !t.reload {
		return t.client, nil
	}

	config := api.DefaultConfig()
	if config.Error != nil {
		return nil, config.Error
	}
	config.Address = client.Address()

	result, err := sdktpm.LoadFromDirectory(t.tpmCertDestination)
	if err != nil {
		return nil, fmt.Errorf("error reading attestation results from %q: %w", t.tpmCertDestination, err)
	}

	tlsConf, err := sdktpm.BuildTLSConfig(result, t.tlsCfg.CACert, t.tlsCfg.CAPath, sdktpm.TPMConfig{
		DevicePath: t.tpmCertDevice,
	})
	if err != nil {
		return nil, fmt.Errorf("error configuring TLS with attestation results from %q: %w", t.tpmCertDestination, err)
	}

	transport, ok := config.HttpClient.Transport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("unsupported HTTPClient transport type %T", config.HttpClient.Transport)
	}

	// Propagate InsecureSkipVerify
	tlsConf.InsecureSkipVerify = transport.TLSClientConfig.InsecureSkipVerify || t.tlsCfg.Insecure

	transport.TLSClientConfig = tlsConf

	// set last hash if load it successfully
	t.latestHash = result.Hash()

	clientToAuth, err = api.NewClient(config)
	if err != nil {
		return nil, err
	}
	if ns := client.Headers().Get(consts.NamespaceHeaderName); ns != "" {
		clientToAuth.SetNamespace(ns)
	}

	// Cache the client for future use
	t.client = clientToAuth

	return clientToAuth, nil
}

// runWatcher uses polling instead of inotify to sense the changes on the cert/key/ca files.
// The reason not to use inotify:
// 1. To not miss any changes, we need to watch the directory instead of files when using inotify.
// 2. These files are not frequently changed/renewed, and they don't need to be reloaded immediately after renewal.
// 3. Some network based filesystem and FUSE don't support inotify.
func (t *tpmMethod) runWatcher() {
	defer close(t.doneCh)

	select {
	case <-t.stopCh:
		return

	case <-t.credSuccessGate:
		// We only start the next loop once we're initially successful,
		// since at startup Authenticate will be called, and we don't want
		// to end up immediately re-authenticating by having found a new
		// value
	}

	for {
		changed := false
		select {
		case <-t.stopCh:
			return

		case <-t.ticker.C:
			t.logger.Debug("Checking if files changed", "path", t.tpmCertDestination)
			result, err := sdktpm.LoadFromDirectory(t.tpmCertDestination)
			// ignore errors in watcher; we could encounter an in-progress write,
			// or a failed partial write that will be retried later.
			if err == nil {
				newHash := result.Hash()
				t.logger.Debug("hash before/after", "new", hex.EncodeToString(newHash), "old", hex.EncodeToString(t.latestHash))
				changed = !bytes.Equal(t.latestHash, newHash)
				if changed {
					t.logger.Info("The cert/key files changed", "path", t.tpmCertDestination,
						"notAfter", result.NotAfter.Format(time.RFC1123Z))
				}
			}
		}

		if changed {
			select {
			case t.credsFound <- struct{}{}:
			case <-t.stopCh:
			}
		}
	}
}
