package pluginutil

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PluginMultiplexingServerImpl struct {
	UnimplementedPluginMultiplexingServer

	Supported bool
}

func (pm PluginMultiplexingServerImpl) MultiplexingSupport(_ context.Context, _ *MultiplexingSupportRequest) (*MultiplexingSupportResponse, error) {
	return &MultiplexingSupportResponse{
		Supported: pm.Supported,
	}, nil
}

func MultiplexingSupported(ctx context.Context, cc grpc.ClientConnInterface, name string) (bool, error) {
	if cc == nil {
		return false, fmt.Errorf("client connection is nil")
	}

	out := strings.Split(os.Getenv(PluginMultiplexingOptOut), ",")
	if strutil.StrListContains(out, name) {
		return false, nil
	}

	req := new(MultiplexingSupportRequest)
	resp, err := NewPluginMultiplexingClient(cc).MultiplexingSupport(ctx, req)
	if err != nil {

		// If the server does not implement the multiplexing server then we can
		// assume it is not multiplexed
		if status.Code(err) == codes.Unimplemented {
			return false, nil
		}

		return false, err
	}
	if resp == nil {
		// Somehow got a nil response, assume not multiplexed
		return false, nil
	}

	return resp.Supported, nil
}

func GetMultiplexIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing plugin multiplexing metadata")
	}

	multiplexIDs := md[MultiplexingCtxKey]
	if len(multiplexIDs) != 1 {
		return "", fmt.Errorf("unexpected number of IDs in metadata: (%d)", len(multiplexIDs))
	}

	multiplexID := multiplexIDs[0]
	if multiplexID == "" {
		return "", fmt.Errorf("empty multiplex ID in metadata")
	}

	return multiplexID, nil
}
