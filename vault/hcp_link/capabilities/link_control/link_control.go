// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package link_control

import (
	"context"
	"fmt"
	"math"
	"runtime/debug"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	scada "github.com/hashicorp/hcp-scada-provider"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
	"github.com/hashicorp/vault/vault/hcp_link/proto/link_control"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type purgePolicyFunc func()

type hcpLinkControlHandler struct {
	link_control.UnimplementedHCPLinkControlServer

	purgeFunc     purgePolicyFunc
	wrappedCore   internal.WrappedCoreStandbyStates
	scadaProvider scada.SCADAProvider
	logger        hclog.Logger

	l          sync.Mutex
	grpcServer *grpc.Server
	running    bool
}

func NewHCPLinkControlService(scadaProvider scada.SCADAProvider, core *vault.Core, policyPurger purgePolicyFunc, baseLogger hclog.Logger) *hcpLinkControlHandler {
	logger := baseLogger.Named(capabilities.LinkControlCapability)
	logger.Trace("initializing HCP Link Control capability")

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 2 * time.Second,
		}),
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.MaxRecvMsgSize(math.MaxInt32),
	)

	handler := &hcpLinkControlHandler{
		purgeFunc:     policyPurger,
		logger:        logger,
		grpcServer:    grpcServer,
		scadaProvider: scadaProvider,
		wrappedCore:   core,
	}

	link_control.RegisterHCPLinkControlServer(grpcServer, handler)
	reflection.Register(grpcServer)

	return handler
}

func (h *hcpLinkControlHandler) Start() error {
	h.l.Lock()
	defer h.l.Unlock()

	if h.running {
		return nil
	}

	// Starting link-control service
	linkControlListener, err := h.scadaProvider.Listen(capabilities.LinkControlCapability)
	if err != nil {
		return fmt.Errorf("failed to initialize link-control capability listener: %w", err)
	}

	if linkControlListener == nil {
		return fmt.Errorf("no listener found for link-control capability")
	}

	// Start the gRPC server
	go func() {
		err = h.grpcServer.Serve(linkControlListener)
		h.logger.Error("server closed", "error", err)
	}()

	h.running = true

	h.logger.Trace("started HCP Link Control capability")
	return nil
}

func (h *hcpLinkControlHandler) Stop() error {
	h.l.Lock()
	defer h.l.Unlock()

	if !h.running {
		return nil
	}

	// Give some time for existing RPCs to drain.
	time.Sleep(cluster.ListenerAcceptDeadline)

	h.logger.Info("Tearing down HCP Link Control")

	h.grpcServer.Stop()

	h.running = false

	return nil
}

func (h *hcpLinkControlHandler) PurgePolicy(ctx context.Context, req *link_control.PurgePolicyRequest) (retResp *link_control.PurgePolicyResponse, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic serving purge policy request", "error", r, "stacktrace", string(debug.Stack()))
			retErr = vault.ErrInternalError
		}
	}()

	standby, perfStandby := h.wrappedCore.StandbyStates()
	// only purging an active node, perf/standby nodes should purge
	// automatically
	if standby || perfStandby {
		h.logger.Debug("cannot purge the policy on a non-active node")
	} else {
		h.purgeFunc()
		h.logger.Debug("Purged token and policy")
	}

	return &link_control.PurgePolicyResponse{}, nil
}
