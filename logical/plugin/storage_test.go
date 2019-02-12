package plugin

import (
	"testing"

	"google.golang.org/grpc"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
)

func TestStorage_impl(t *testing.T) {
	var _ logical.Storage = new(GRPCStorageClient)
}

func TestStorage_GRPC(t *testing.T) {
	storage := &logical.InmemStorage{}
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterStorageServer(s, &GRPCStorageServer{
			impl: storage,
		})
	})
	defer client.Close()

	testStorage := &GRPCStorageClient{client: pb.NewStorageClient(client)}

	logical.TestStorage(t, testStorage)

}
