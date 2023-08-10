// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api_capability

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	scada "github.com/hashicorp/hcp-scada-provider"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities"
)

type APICapability struct {
	l             sync.Mutex
	logger        hclog.Logger
	scadaProvider scada.SCADAProvider
	scadaServer   *http.Server
	tokenManager  *HCPLinkTokenManager
	running       bool
}

var _ capabilities.Capability = &APICapability{}

func NewAPICapability(scadaConfig *scada.Config, scadaProvider scada.SCADAProvider, core *vault.Core, logger hclog.Logger) (*APICapability, error) {
	apiLogger := logger.Named(capabilities.APICapability)
	tokenManager, err := NewHCPLinkTokenManager(scadaConfig, core, apiLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to start HCP Link token manager")
	}

	linkHandler := injectBatchTokenHandler(vaulthttp.Handler.Handler(&vault.HandlerProperties{Core: core}), tokenManager)

	apiLogger.Trace("initializing HCP Link API capability")

	// server defaults
	server := &http.Server{
		Handler:           linkHandler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ErrorLog:          apiLogger.StandardLogger(nil),
	}
	return &APICapability{
		logger:        apiLogger,
		scadaProvider: scadaProvider,
		scadaServer:   server,
		tokenManager:  tokenManager,
	}, nil
}

func (c *APICapability) Start() error {
	c.l.Lock()
	defer c.l.Unlock()

	if c.running {
		return nil
	}

	// Start listening on a SCADA capability
	listener, err := c.scadaProvider.Listen(capabilities.APICapability)
	if err != nil {
		return fmt.Errorf("failed to start listening on a capability: %w", err)
	}

	go func() {
		err = c.scadaServer.Serve(listener)
		c.logger.Error("server closed", "error", err)
	}()

	c.running = true
	c.logger.Info("started HCP Link API capability")

	return nil
}

func (c *APICapability) Stop() error {
	c.l.Lock()
	defer c.l.Unlock()

	if !c.running {
		return nil
	}

	var retErr *multierror.Error

	c.logger.Info("Tearing down HCP Link API capability")

	err := c.scadaServer.Shutdown(context.Background())
	if err != nil {
		retErr = multierror.Append(err, fmt.Errorf("failed to shutdown scada provider HTTP server %w", err))
	}
	c.scadaServer = nil

	c.tokenManager.Shutdown()
	c.tokenManager = nil

	c.running = false

	return retErr.ErrorOrNil()
}

func (c *APICapability) PurgePolicy() {
	if c.tokenManager == nil {
		return
	}
	c.tokenManager.ForgetTokenPolicy()

	return
}

func injectBatchTokenHandler(handler http.Handler, tokenManager *HCPLinkTokenManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenManager.logger.Debug("received request", "method", r.Method, "path", r.URL.Path)

		// Only the hcp link token should be used
		r.Header.Del(consts.AuthHeaderName)

		// for Standby or perfStandby return 412
		standby, perfStandby := tokenManager.wrappedCore.StandbyStates()

		hcpLinkToken := tokenManager.HandleTokenPolicy(r.Context(), !standby && !perfStandby)

		if standby || perfStandby {
			logical.RespondError(w, http.StatusPreconditionFailed, fmt.Errorf("API capability is inactive in non-Active nodes"))
			return
		}

		r.Header.Set(consts.AuthHeaderName, hcpLinkToken)

		handler.ServeHTTP(w, r)
		return
	})
}
