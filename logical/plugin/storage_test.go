package plugin

import (
	"testing"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
)

func TestStorage_impl(t *testing.T) {
	var _ logical.Storage = new(StorageClient)
}

func TestStorage_operations(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	storage := &logical.InmemStorage{}

	server.RegisterName("Plugin", &StorageServer{
		impl: storage,
	})

	testStorage := &StorageClient{client: client}

	logical.TestStorage(t, testStorage)
}
