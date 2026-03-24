// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pkiexternalca

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"go.uber.org/atomic"
)

// Server is a CE stub; PKI external CA is an enterprise-only feature.
type Server struct {
	DoneCh  chan struct{}
	stopped *atomic.Bool
}

// NewServer returns a stub server for CE builds.
func NewServer(cfg *ServerConfig) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("server config cannot be nil")
	}
	return &Server{
		DoneCh:  make(chan struct{}),
		stopped: atomic.NewBool(false),
	}, nil
}

// Run waits for context cancellation; PKI external CA is never active in CE builds.
func (s *Server) Run(ctx context.Context, _ chan string, _ *api.Client) error {
	<-ctx.Done()
	return nil
}

// Stop closes DoneCh idempotently.
func (s *Server) Stop() {
	if s.stopped.CAS(false, true) {
		close(s.DoneCh)
	}
}

// CertIssuedCh returns nil in CE builds.
func (s *Server) CertIssuedCh() <-chan struct{} { return nil }

// TemplatePEMByName returns an error in CE builds.
func (s *Server) TemplatePEMByName(_ string) (any, error) {
	return nil, fmt.Errorf("pki_external_ca is not supported in this build")
}
