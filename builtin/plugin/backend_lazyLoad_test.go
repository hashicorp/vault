package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/logging"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin"
)

func TestBackend_startBackend(t *testing.T) {

	sysView := newTestSystemView()

	ctx := context.Background()
	config := &logical.BackendConfig{
		Logger: logging.NewVaultLogger(hclog.Trace),
		System: sysView,
		Config: map[string]string{
			"plugin_name": "test-plugin",
			"plugin_type": "secret",
		},
	}

	orig, err := plugin.NewBackend(ctx, "test-plugin", consts.PluginTypeSecrets, sysView, config, true)
	if err != nil {
		t.Fatal(err)
	}

	b := &PluginBackend{
		Backend: orig,
		config:  config,
	}

	err = b.startBackend(ctx, &logical.InmemStorage{})
	if err != nil {
		t.Fatal(err)
	}

	if !b.loaded {
		t.Fatalf("not loaded")
	}

	ob := orig.(*testBackend)
	if !ob.cleaned {
		t.Fatalf("not cleaned")
	}
	if ob.setup {
		t.Fatalf("setup")
	}
	if ob.initialized {
		t.Fatalf("initialized")
	}

	nb := b.Backend.(*testBackend)
	if nb.cleaned {
		t.Fatalf("cleaned")
	}
	if !nb.setup {
		t.Fatalf("not setup")
	}
	if !nb.initialized {
		t.Fatalf("not initialized")
	}
}

//------------------------------------------------------------------

type testBackend struct {
	cleaned     bool
	setup       bool
	initialized bool
}

var _ logical.Backend = (*testBackend)(nil)

func (b *testBackend) Cleanup(context.Context) {
	b.cleaned = true
}

func (b *testBackend) Setup(context.Context, *logical.BackendConfig) error {
	b.setup = true
	return nil
}

func (b *testBackend) Initialize(context.Context, *logical.InitializationRequest) error {
	b.initialized = true
	return nil
}

func (b *testBackend) Type() logical.BackendType {
	return logical.TypeLogical
}

func (b *testBackend) SpecialPaths() *logical.Paths {
	return &logical.Paths{
		Root: []string{"test-root"},
	}
}

func (b *testBackend) HandleRequest(context.Context, *logical.Request) (*logical.Response, error) {
	panic("not needed")
}
func (b *testBackend) System() logical.SystemView {
	panic("not needed")
}
func (b *testBackend) Logger() hclog.Logger {
	panic("not needed")
}
func (b *testBackend) HandleExistenceCheck(context.Context, *logical.Request) (bool, bool, error) {
	panic("not needed")
}
func (b *testBackend) InvalidateKey(context.Context, string) {
	panic("not needed")
}

//------------------------------------------------------------------

type testSystemView struct {

	// its probably not StaticSystemView's intended purpose to be embedded this
	// way, but we are doing it anyway for testing, so we don't have to define
	// a whole logical.SystemView
	logical.StaticSystemView

	factory logical.Factory
}

func newTestSystemView() testSystemView {
	return testSystemView{
		factory: func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
			return &testBackend{}, nil
		},
	}
}

func (v testSystemView) LookupPlugin(context.Context, string, consts.PluginType) (*pluginutil.PluginRunner, error) {

	return &pluginutil.PluginRunner{
		Name:    "test-plugin-runner",
		Builtin: true,
		BuiltinFactory: func() (interface{}, error) {
			return v.factory, nil
		},
	}, nil
}
