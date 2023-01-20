package meta

import (
	"context"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	scada "github.com/hashicorp/hcp-scada-provider"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/hcp_link/capabilities"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
	"github.com/hashicorp/vault/vault/hcp_link/proto/meta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type hcpLinkMetaHandler struct {
	meta.UnimplementedHCPLinkMetaServer

	wrappedCore   internal.WrappedCoreMeta
	scadaProvider scada.SCADAProvider
	logger        hclog.Logger

	l          sync.Mutex
	grpcServer *grpc.Server
	stopCh     chan struct{}
	running    bool
}

func NewHCPLinkMetaService(scadaProvider scada.SCADAProvider, c *vault.Core, baseLogger hclog.Logger) *hcpLinkMetaHandler {
	logger := baseLogger.Named(capabilities.MetaCapability)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 2 * time.Second,
		}),
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.MaxRecvMsgSize(math.MaxInt32),
	)

	handler := &hcpLinkMetaHandler{
		wrappedCore:   c,
		logger:        logger,
		grpcServer:    grpcServer,
		scadaProvider: scadaProvider,
	}

	meta.RegisterHCPLinkMetaServer(grpcServer, handler)
	reflection.Register(grpcServer)

	return handler
}

func (h *hcpLinkMetaHandler) Start() error {
	h.l.Lock()
	defer h.l.Unlock()

	if h.running {
		return nil
	}

	// Starting meta service
	metaListener, err := h.scadaProvider.Listen(capabilities.MetaCapability)
	if err != nil {
		return fmt.Errorf("failed to initialize meta capability listener: %w", err)
	}

	if metaListener == nil {
		return fmt.Errorf("no listener found for meta capability")
	}

	h.logger.Info("starting HCP meta capability")
	// Start the gRPC server
	go func() {
		err = h.grpcServer.Serve(metaListener)
		h.logger.Error("server closed", "error", err)
	}()

	h.running = true

	return nil
}

func (h *hcpLinkMetaHandler) Stop() error {
	h.l.Lock()
	defer h.l.Unlock()

	if !h.running {
		return nil
	}

	// Give some time for existing RPCs to drain.
	time.Sleep(cluster.ListenerAcceptDeadline)

	h.logger.Info("tearing down HCP meta capability")

	if h.stopCh != nil {
		close(h.stopCh)
		h.stopCh = nil
	}

	h.grpcServer.Stop()

	h.running = false

	return nil
}

func (h *hcpLinkMetaHandler) ListNamespaces(ctx context.Context, req *meta.ListNamespacesRequest) (retResp *meta.ListNamespacesResponse, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic serving list namespaces request", "error", r, "stacktrace", string(debug.Stack()))
			retErr = vault.ErrInternalError
		}
	}()

	children := h.wrappedCore.ListNamespaces(true)

	var namespaces []string
	for _, child := range children {
		namespaces = append(namespaces, child.Path)
	}

	return &meta.ListNamespacesResponse{
		Paths: namespaces,
	}, nil
}

func (h *hcpLinkMetaHandler) ListMounts(ctx context.Context, req *meta.ListMountsRequest) (retResp *meta.ListMountsResponse, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic serving list mounts request", "error", r, "stacktrace", string(debug.Stack()))
			retErr = vault.ErrInternalError
		}
	}()

	mountEntries, err := h.wrappedCore.ListMounts()
	if err != nil {
		return nil, fmt.Errorf("unable to list secret mounts: %w", err)
	}

	var mounts []*meta.Mount
	for _, entry := range mountEntries {
		nsID := entry.NamespaceID
		path := entry.Path

		if nsID != namespace.RootNamespaceID {
			ns, err := h.wrappedCore.NamespaceByID(ctx, entry.NamespaceID)
			if err != nil {
				return nil, fmt.Errorf("unable to get namespace associated with secret mount: %w", err)
			}

			path = ns.Path + path
		}

		mounts = append(mounts, &meta.Mount{
			Path:        path,
			Type:        entry.Type,
			Description: entry.Description,
		})
	}

	return &meta.ListMountsResponse{
		Mounts: mounts,
	}, nil
}

func (h *hcpLinkMetaHandler) ListAuths(ctx context.Context, req *meta.ListAuthsRequest) (retResp *meta.ListAuthResponse, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic serving list auths request", "error", r, "stacktrace", string(debug.Stack()))
			retErr = vault.ErrInternalError
		}
	}()

	authEntries, err := h.wrappedCore.ListAuths()
	if err != nil {
		return nil, fmt.Errorf("unable to list auth mounts: %w", err)
	}

	var auths []*meta.Auth
	for _, entry := range authEntries {
		nsID := entry.NamespaceID
		path := entry.Path

		if nsID != namespace.RootNamespaceID {
			ns, err := h.wrappedCore.NamespaceByID(ctx, entry.NamespaceID)
			if err != nil {
				return nil, fmt.Errorf("unable to get namespace associated with auth mount: %w", err)
			}

			path = ns.Path + path
		}

		auths = append(auths, &meta.Auth{
			Path:        path,
			Type:        entry.Type,
			Description: entry.Description,
		})
	}

	return &meta.ListAuthResponse{
		Auths: auths,
	}, nil
}

func (h *hcpLinkMetaHandler) GetClusterStatus(ctx context.Context, req *meta.GetClusterStatusRequest) (retResp *meta.GetClusterStatusResponse, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic serving cluster status request", "error", r, "stacktrace", string(debug.Stack()))
			retErr = vault.ErrInternalError
		}
	}()

	if h.wrappedCore.HAStateWithLock() != consts.Active {
		return nil, fmt.Errorf("node not active")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch hostname: %w", err)
	}

	haEnabled := h.wrappedCore.HAEnabled()
	haStatus := &meta.HAStatus{
		Enabled: haEnabled,
	}

	if haEnabled {
		leader := &meta.HANode{
			Hostname: hostname,
		}

		peers := h.wrappedCore.GetHAPeerNodesCached()

		haNodes := make([]*meta.HANode, len(peers)+1)
		haNodes[0] = leader

		for i, peerNode := range peers {
			haNodes[i+1] = &meta.HANode{
				Hostname: peerNode.Hostname,
			}
		}

		haStatus.Nodes = haNodes
	}

	raftStatus := &meta.RaftStatus{}
	raftConfig, err := h.wrappedCore.GetRaftConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get Raft configuration: %w", err)
	}

	if raftConfig != nil {
		raftServers := make([]*meta.RaftServer, len(raftConfig.Servers))

		var voterCount uint32
		for i, srv := range raftConfig.Servers {
			raftServers[i] = &meta.RaftServer{
				NodeID:          srv.NodeID,
				Address:         srv.Address,
				Voter:           srv.Voter,
				Leader:          srv.Leader,
				ProtocolVersion: srv.ProtocolVersion,
			}

			if srv.Voter {
				voterCount++
			}
		}

		raftStatus.RaftConfiguration = &meta.RaftConfiguration{
			Servers: raftServers,
		}

		evenVoterMessage := "Vault should have access to an odd number of voter nodes."
		largeClusterMessage := "Very large cluster detected."
		var quorumWarning string

		if voterCount == 1 {
			quorumWarning = "Only one server node found. Vault is not running in high availability mode."
		} else if voterCount%2 == 0 && voterCount > 7 {
			quorumWarning = evenVoterMessage + " " + largeClusterMessage
		} else if voterCount%2 == 0 {
			quorumWarning = evenVoterMessage
		} else if voterCount > 7 {
			quorumWarning = largeClusterMessage
		}

		raftStatus.QuorumWarning = quorumWarning
	}

	raftAutopilotState, err := h.wrappedCore.GetRaftAutopilotState(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get Raft Autopilot state: %w", err)
	}

	if raftAutopilotState != nil {
		autopilotStatus := &meta.AutopilotStatus{
			Healthy: raftAutopilotState.Healthy,
		}

		autopilotServers := make([]*meta.AutopilotServer, 0)
		for _, srv := range raftAutopilotState.Servers {
			autopilotServers = append(autopilotServers, &meta.AutopilotServer{
				ID:      srv.ID,
				Healthy: srv.Healthy,
			})
		}

		raftStatus.AutopilotStatus = autopilotStatus
	}

	resp := &meta.GetClusterStatusResponse{
		ClusterID:   h.wrappedCore.ClusterID(),
		HAStatus:    haStatus,
		RaftStatus:  raftStatus,
		StorageType: h.wrappedCore.StorageType(),
	}

	return resp, nil
}
