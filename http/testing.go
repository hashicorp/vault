package http

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestServer(t *testing.T, core *vault.Core) (net.Listener, string) {
	fail := func(format string, args ...interface{}) {
		panic(fmt.Sprintf(format, args...))
	}
	if t != nil {
		fail = t.Fatalf
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fail("err: %s", err)
	}
	addr := "http://" + ln.Addr().String()

	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: Handler(core),
	}
	go server.Serve(ln)

	return ln, addr
}
