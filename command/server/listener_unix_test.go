package server

import (
	"net"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/mitchellh/cli"
)

func TestUnixListener(t *testing.T) {
	ln, _, _, err := unixListenerFactory(&configutil.Listener{
		Address: filepath.Join(t.TempDir(), "/vault.sock"),
	}, nil, cli.NewMockUi())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connFn := func(lnReal net.Listener) (net.Conn, error) {
		return net.Dial("unix", ln.Addr().String())
	}

	testListenerImpl(t, ln, connFn, "", 0, "", false)
}
