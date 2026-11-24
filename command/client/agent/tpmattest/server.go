// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package tpmattest

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/client/agent/config"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
	"go.uber.org/atomic"
)

type ServerConfig struct {
	Logger      hclog.Logger
	AgentConfig *config.Config
}

type RunnableClient interface {
	Run(ctx context.Context, client *api.Client) error
	Stop()
}

const (
	defaultMaxRetries = 5
)

type Server struct {
	config *ServerConfig
	logger hclog.Logger

	managers       []*Manager
	managersByName map[string]*Manager

	mu sync.RWMutex

	doneCh  chan struct{}
	stopped *atomic.Bool
}

func NewServer(cfg *ServerConfig) (RunnableClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("server config cannot be nil")
	}

	server := &Server{
		config:         cfg,
		logger:         cfg.Logger,
		doneCh:         make(chan struct{}),
		stopped:        atomic.NewBool(false),
		managers:       make([]*Manager, 0, len(cfg.AgentConfig.TPMCerts)),
		managersByName: make(map[string]*Manager),
	}

	return server, nil
}

func (s *Server) Run(ctx context.Context, client *api.Client) error {
	s.logger.Info("starting tpm attestation server")
	defer s.logger.Info("tpm attestation server stopped")

	if len(s.config.AgentConfig.TPMCerts) == 0 {
		s.logger.Info("no tpm_cert blocks configured, tpm attestation server idle")
		<-ctx.Done()
		return nil
	}

	for idx, tpmConfig := range s.config.AgentConfig.TPMCerts {
		tpmClient, err := client.CloneWithHeaders()
		if err != nil {
			return fmt.Errorf("failed to clone client for tpm_cert[%d]: %w", idx, err)
		}

		tpmClient.SetMaxRetries(defaultMaxRetries)
		tpmClient.ClearToken()

		if tpmConfig.Namespace != "" {
			tpmClient.SetNamespace(tpmConfig.Namespace)
		}

		umask, err := parseUmask(tpmConfig.Destination.Umask)
		if err != nil {
			return fmt.Errorf("failed to parse umask for tpm_cert[%d]: %w", idx, err)
		}

		var devicePath string
		if tpmConfig.TPM != nil {
			devicePath = tpmConfig.TPM.DevicePath
		}

		var csrParams *sdktpm.CSRParams
		if tpmConfig.CSR != nil {
			csrParams = &sdktpm.CSRParams{
				CommonName: tpmConfig.CSR.CommonName,
				Country:    tpmConfig.CSR.Country,
				State:      tpmConfig.CSR.State,
				Locality:   tpmConfig.CSR.Locality,
				Org:        tpmConfig.CSR.Org,
				OrgUnit:    tpmConfig.CSR.OrgUnit,
			}
		}

		attestConfig := &sdktpm.AttestationConfig{
			Namespace: tpmConfig.Namespace,
			MountPath: tpmConfig.MountPath,
			RoleName:  tpmConfig.RoleName,
			CSRParams: csrParams,
			TPMConfig: sdktpm.TPMConfig{DevicePath: devicePath},
		}

		attest := func(ctx context.Context) (*sdktpm.AttestationResult, error) {
			return sdktpm.Attest(ctx, tpmClient, attestConfig)
		}

		validate := ValidateFunc(func(r *sdktpm.AttestationResult) error {
			return r.Validate(sdktpm.TPMConfig{DevicePath: devicePath})
		})

		renewPercent := *tpmConfig.Renew.PercentRenewBeforeExpiry

		manager, err := NewManager(
			attest,
			validate,
			tpmConfig.Destination.Path,
			umask,
			renewPercent,
			s.logger.Named(fmt.Sprintf("tpm[%d]", idx)),
		)
		if err != nil {
			return fmt.Errorf("failed to create manager for tpm_cert[%d]: %w", idx, err)
		}

		s.mu.Lock()
		s.managers = append(s.managers, manager)
		s.managersByName[tpmConfig.Name] = manager
		s.mu.Unlock()
	}

	s.mu.RLock()
	managers := make([]*Manager, len(s.managers))
	copy(managers, s.managers)
	s.mu.RUnlock()

	for idx, manager := range managers {
		go func(idx int, mgr *Manager) {
			mgr.Start(ctx)
		}(idx, manager)
	}

	<-ctx.Done()
	s.Stop()
	return nil
}

func (s *Server) Stop() {
	if s.stopped.CompareAndSwap(false, true) {
		s.logger.Info("stopping tpm attestation server")

		s.mu.RLock()
		for _, manager := range s.managers {
			manager.Stop()
		}
		s.mu.RUnlock()

		close(s.doneCh)
	}
}
