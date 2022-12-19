package meta

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	scada "github.com/hashicorp/hcp-scada-provider"
	"github.com/hashicorp/vault/helper/namespace"
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

	wrappedCore   internal.WrappedCoreListNamespacesMounts
	scadaProvider scada.SCADAProvider
	logger        hclog.Logger

	l          sync.Mutex
	grpcServer *grpc.Server
	stopCh     chan struct{}
	running    bool
}

func NewHCPLinkMetaService(scadaProvider scada.SCADAProvider, c *vault.Core, baseLogger hclog.Logger) *hcpLinkMetaHandler {
	logger := baseLogger.Named(capabilities.MetaCapability)
	logger.Info("Setting up HCP Link Meta Service")

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

	h.logger.Info("starting HCP Link Meta Service")
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

	h.logger.Info("Tearing down HCP Link Meta Service")

	if h.stopCh != nil {
		close(h.stopCh)
		h.stopCh = nil
	}

	h.grpcServer.Stop()

	h.running = false

	return nil
}

func (h *hcpLinkMetaHandler) ListNamespaces(ctx context.Context, req *meta.ListNamespacesRequest) (*meta.ListNamespacesResponse, error) {
	children := h.wrappedCore.ListNamespaces(true)

	var namespaces []string
	for _, child := range children {
		namespaces = append(namespaces, child.Path)
	}

	return &meta.ListNamespacesResponse{
		Paths: namespaces,
	}, nil
}

func (h *hcpLinkMetaHandler) ListMounts(ctx context.Context, req *meta.ListMountsRequest) (*meta.ListMountsResponse, error) {
	mountEntries, err := h.wrappedCore.ListMounts()
	if err != nil {
		return nil, err
	}

	var mounts []*meta.Mount
	for _, entry := range mountEntries {
		nsID := entry.NamespaceID
		path := entry.Path

		if nsID != namespace.RootNamespaceID {
			ns, err := h.wrappedCore.NamespaceByID(ctx, entry.NamespaceID)
			if err != nil {
				return nil, err
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

func (h *hcpLinkMetaHandler) ListAuths(ctx context.Context, req *meta.ListAuthsRequest) (*meta.ListAuthResponse, error) {
	authEntries, err := h.wrappedCore.ListAuths()
	if err != nil {
		return nil, err
	}

	var auths []*meta.Auth
	for _, entry := range authEntries {
		nsID := entry.NamespaceID
		path := entry.Path

		if nsID != namespace.RootNamespaceID {
			ns, err := h.wrappedCore.NamespaceByID(ctx, entry.NamespaceID)
			if err != nil {
				return nil, err
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
