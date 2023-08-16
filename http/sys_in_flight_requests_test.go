package http

import (
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func TestInFlightRequestUnauthenticated(t *testing.T) {
	conf := &vault.CoreConfig{}
	core, _, token := vault.TestCoreUnsealedWithConfig(t, conf)
	ln, addr := TestServer(t, core)
	TestServerAuth(t, addr, token)

	// Default: Only authenticated access
	resp := testHttpGet(t, "", addr+"/v1/sys/in-flight-req")
	testResponseStatus(t, resp, 403)
	resp = testHttpGet(t, token, addr+"/v1/sys/in-flight-req")
	testResponseStatus(t, resp, 200)

	// Close listener
	ln.Close()

	// Setup new custom listener with unauthenticated metrics access
	ln, addr = TestListener(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			InFlightRequestLogging: configutil.ListenerInFlightRequestLogging{
				UnauthenticatedInFlightAccess: true,
			},
		},
	}
	TestServerWithListenerAndProperties(t, ln, addr, core, props)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Test without token
	resp = testHttpGet(t, "", addr+"/v1/sys/in-flight-req")
	testResponseStatus(t, resp, 200)

	// Should also work with token
	resp = testHttpGet(t, token, addr+"/v1/sys/in-flight-req")
	testResponseStatus(t, resp, 200)
}
