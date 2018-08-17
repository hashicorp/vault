package http

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestListener(tb testing.TB) (net.Listener, string) {
	fail := func(format string, args ...interface{}) {
		panic(fmt.Sprintf(format, args...))
	}
	if tb != nil {
		fail = tb.Fatalf
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fail("err: %s", err)
	}
	addr := "http://" + ln.Addr().String()
	return ln, addr
}

func TestServerWithListenerAndProperties(tb testing.TB, ln net.Listener, addr string, core *vault.Core, props *vault.HandlerProperties) {
	// Create a muxer to handle our requests so that we can authenticate
	// for tests.
	mux := http.NewServeMux()
	mux.Handle("/_test/auth", http.HandlerFunc(testHandleAuth))
	mux.Handle("/", Handler(props))

	server := &http.Server{
		Addr:     ln.Addr().String(),
		Handler:  mux,
		ErrorLog: core.Logger().StandardLogger(nil),
	}
	go server.Serve(ln)
}

func TestServerWithListener(tb testing.TB, ln net.Listener, addr string, core *vault.Core) {
	// Create a muxer to handle our requests so that we can authenticate
	// for tests.
	props := &vault.HandlerProperties{
		Core:           core,
		MaxRequestSize: DefaultMaxRequestSize,
	}
	TestServerWithListenerAndProperties(tb, ln, addr, core, props)
}

func TestServer(tb testing.TB, core *vault.Core) (net.Listener, string) {
	ln, addr := TestListener(tb)
	TestServerWithListener(tb, ln, addr, core)
	return ln, addr
}

func TestServerAuth(tb testing.TB, addr string, token string) {
	if _, err := http.Get(addr + "/_test/auth?token=" + token); err != nil {
		tb.Fatalf("error authenticating: %s", err)
	}
}

func testHandleAuth(w http.ResponseWriter, req *http.Request) {
	respondOk(w, nil)
}
