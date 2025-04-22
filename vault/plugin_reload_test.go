package vault

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestReload vets the reload functionality of Core by adding plugins of varying kinds
// and finds out if they reload as expected.
func TestReload(t *testing.T) {
	ctx := context.Background()
	log := corehelpers.NewTestLogger(t)
	core := &Core{
		mountsLock: &sync.RWMutex{},
		authLock:   &sync.RWMutex{},
		router:     NewRouter(),
		logger:     log,
	}

	cases := []struct {
		name     string
		backends map[string]logical.Factory
	}{
		{
			name: "reload one",
			backends: map[string]logical.Factory{
				"plugin": plugin.Factory,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			core.configureCredentialsBackends(c.backends, log)
			for k := range core.credentialBackends {
				fmt.Println(k)
			}
			if len(core.credentialBackends) != len(c.backends)+1 { // registering backends adds a plugin called "token" by default
				t.Fatalf("backend lengths didn't match: %d vs %d", len(core.credentialBackends), len(c.backends))
			}
			err := core.reloadMatchingPluginMounts(namespace.RootContext(ctx), namespace.RootNamespace, []string{"auth/plugin"})
			if err != nil {
				t.Fatalf("%s", err)
			}
		})
	}
}
