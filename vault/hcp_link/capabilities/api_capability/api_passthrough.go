// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api_capability

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	scada "github.com/hashicorp/hcp-scada-provider"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
)

type APIPassThroughCapability struct {
	l             sync.Mutex
	logger        hclog.Logger
	scadaProvider scada.SCADAProvider
	scadaServer   *http.Server
	running       bool
}

var _ capabilities.Capability = &APIPassThroughCapability{}

func NewAPIPassThroughCapability(scadaProvider scada.SCADAProvider, core *vault.Core, logger hclog.Logger) (*APIPassThroughCapability, error) {
	apiLogger := logger.Named(capabilities.APIPassThroughCapability)

	linkHandler := requestHandler(vaulthttp.Handler.Handler(&vault.HandlerProperties{Core: core}), core, apiLogger)

	apiLogger.Trace("initializing HCP Link API PassThrough capability")

	// server defaults
	server := &http.Server{
		Handler:           linkHandler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ErrorLog:          apiLogger.StandardLogger(nil),
	}
	return &APIPassThroughCapability{
		logger:        apiLogger,
		scadaProvider: scadaProvider,
		scadaServer:   server,
	}, nil
}

func (p *APIPassThroughCapability) Start() error {
	p.l.Lock()
	defer p.l.Unlock()

	if p.running {
		return nil
	}

	// Start listening on a SCADA capability
	listener, err := p.scadaProvider.Listen(capabilities.APIPassThroughCapability)
	if err != nil {
		return fmt.Errorf("failed to start listening on a capability: %w", err)
	}

	go func() {
		err = p.scadaServer.Serve(listener)
		p.logger.Error("server closed", "error", err)
	}()

	p.running = true
	p.logger.Info("started HCP Link API PassThrough capability")

	return nil
}

func (p *APIPassThroughCapability) Stop() error {
	p.l.Lock()
	defer p.l.Unlock()

	if !p.running {
		return nil
	}

	p.logger.Info("Tearing down HCP Link API passthrough capability")

	var retErr error
	err := p.scadaServer.Shutdown(context.Background())
	if err != nil {
		retErr = fmt.Errorf("failed to shutdown scada provider HTTP server %w", err)
	}
	p.scadaServer = nil
	p.running = false

	return retErr
}

func requestHandler(handler http.Handler, wrappedCore internal.WrappedCoreStandbyStates, logger hclog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("received a request in HCP link API passthrough", "method", r.Method, "path", r.URL.Path)

		handler.ServeHTTP(w, r)
		return
	})
}
