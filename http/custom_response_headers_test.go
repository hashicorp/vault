package http

import (
	"testing"

	"github.com/hashicorp/vault/vault"
)

var defaultCustomHeaders = map[string]string {
	"Strict-Transport-Security": "max-age=1; domains",
	"Content-Security-Policy": "default-src 'others'",
	"X-Custom-Header": "Custom header value default",
	"X-Frame-Options": "Deny",
	"X-Content-Type-Options": "nosniff",
	"Content-Type": "application/json",
	"X-XSS-Protection": "1; mode=block",
}

var customHeader200 = map[string]string {
	"Someheader-200": "200",
	"X-Custom-Header": "Custom header value 200",
}

var customHeader3xx = map[string]string {
	"X-Custom-Header": "Custom header value 3xx",
	"X-Vault-Ignored-3xx": "Ignored 3xx",
}

var customHeader2xx = map[string]string {
	"X-Custom-Header": "Custom header value 2xx",
}

var customHeader400 = map[string]string {
	"Someheader-400": "400",
}

var customHeader4xx = map[string]string {
	"Someheader-4xx": "4xx",
}

var customHeader405 = map[string]string {
	"Someheader-405": "405",
}

func TestCustomResponseHeaders(t *testing.T) {
	core, _, token := vault.TestCoreWithCustomResponseHeaderAndUI(t, true)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/sys/raw/")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/v1/sys/generate-recovery-token/attempt")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/v1/sys/generate-recovery-token/update")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/v1/sys/config/state/")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/v1/sys/seal")
	testResponseStatus(t, resp, 405)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)
	testResponseHeader(t, resp, customHeader405)

	resp = testHttpGet(t, token, addr+"/v1/sys/step-down")
	testResponseStatus(t, resp, 405)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)
	testResponseHeader(t, resp, customHeader405)

	resp = testHttpGet(t, token, addr+"/v1/sys/unseal")
	testResponseStatus(t, resp, 405)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)
	testResponseHeader(t, resp, customHeader405)

	resp = testHttpGet(t, token, addr+"/v1/sys/leader")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/health")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/generate-root/attempt")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/generate-root/update")
	testResponseStatus(t, resp, 400)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)
	testResponseHeader(t, resp, customHeader400)

	resp = testHttpGet(t, token, addr+"/v1/sys/rekey/init")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/rekey/update")
	testResponseStatus(t, resp, 400)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)
	testResponseHeader(t, resp, customHeader400)

	resp = testHttpGet(t, token, addr+"/v1/sys/rekey/verify")
	testResponseStatus(t, resp, 400)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)
	testResponseHeader(t, resp, customHeader400)

	resp = testHttpGet(t, token, addr+"/v1/sys/")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/v1/sys")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/v1/")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/v1")
	testResponseStatus(t, resp, 404)
	testResponseHeader(t, resp, defaultCustomHeaders)
	testResponseHeader(t, resp, customHeader4xx)

	resp = testHttpGet(t, token, addr+"/")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/host-info")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/init")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/seal-status")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/v1/sys/auth")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/ui")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpGet(t, token, addr+"/ui/")
	testResponseStatus(t, resp, 200)
	testResponseHeader(t, resp, customHeader200)

	resp = testHttpPost(t, token, addr+"/v1/sys/auth/foo", map[string]interface{}{
		"type":        "noop",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)
	testResponseHeader(t, resp, customHeader2xx)

}
