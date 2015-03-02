package physical

import "github.com/hashicorp/consul/api"

// ConsulBackend is a physical backend that stores data at specific
// prefix within Consul. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ConsulBackend struct {
}

// NewConsulBackend constructs a Consul backend using the given API client
// and the prefix in the KV store.
func NewConsulBackend(client *api.Client, prefix string) (*ConsulBackend, error) {
	// TODO
	c := &ConsulBackend{}
	return c, nil
}
