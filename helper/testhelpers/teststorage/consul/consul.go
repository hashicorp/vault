package consul

import (
	realtesting "testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	physConsul "github.com/hashicorp/vault/physical/consul"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/go-testing-interface"
)

func MakeConsulBackend(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle {
	cleanup, config := consul.PrepareTestContainer(t.(*realtesting.T), "", false, true)

	consulConf := map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
		"max_parallel": "32",
	}
	consulBackend, err := physConsul.NewConsulBackend(consulConf, logger)
	if err != nil {
		t.Fatal(err)
	}
	return &vault.PhysicalBackendBundle{
		Backend: consulBackend,
		Cleanup: cleanup,
	}
}

func ConsulBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = teststorage.SharedPhysicalFactory(MakeConsulBackend)
}
