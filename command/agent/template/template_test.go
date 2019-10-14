package template

import (
	"context"
	"testing"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

// TestNewServer is a simple test to make sure NewServer returns a Server and
// channel
func TestNewServer(t *testing.T) {
	ts, ch := NewServer(&ServerConfig{})
	if ts == nil {
		t.Fatal("nil server returned")
	}
	if ch == nil {
		t.Fatal("nil blocking channel returned")
	}
}

func TestServerRun(t *testing.T) {
	// create http test server

	templateTokenCh := make(chan string, 1)
	ctx, _ := context.WithCancel(context.Background())
	sc := ServerConfig{
		Logger: logging.NewVaultLogger(hclog.Trace),
		VaultConf: &config.Vault{
			Address: "http://127.0.0.1:8200", // replace with test server address
		},
	}
	// var tsDoneCh, unblockCh chan struct{}
	var unblockCh chan struct{}
	ts, unblockCh := NewServer(&sc)
	if ts == nil {
		t.Fatal("nil server returned")
	}
	if unblockCh == nil {
		t.Fatal("nil blocking channel returned")
	}
	var tcs []*ctconfig.TemplateConfig
	ts.Run(ctx, templateTokenCh, tcs)

	// Unblock should close immediately b/c there are no templates to render
	select {
	case <-ctx.Done():
	case <-unblockCh:
	}
}
