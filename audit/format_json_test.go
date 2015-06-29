package audit

import (
	"bytes"
	"testing"

	"github.com/hashicorp/vault/logical"
	"errors"
)

func TestFormatJSON_formatRequest(t *testing.T) {
	cases := map[string]struct {
		Auth   *logical.Auth
		Req    *logical.Request
		Err    error
		Result string
	}{
		"auth, request": {
			&logical.Auth{ClientToken: "foo", Policies: []string{"root"}},
			&logical.Request{
				Operation: logical.WriteOperation,
				Path:      "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
			},
			errors.New("this is an error"),
			testFormatJSONReqBasicStr,
		},
	}

	for name, tc := range cases {
		var buf bytes.Buffer
		var format FormatJSON
		if err := format.FormatRequest(&buf, tc.Auth, tc.Req, tc.Err); err != nil {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}

		if buf.String() != tc.Result {
			t.Fatalf(
				"bad: %s\nResult:\n\n%s\n\nExpected:\n\n%s",
				name, buf.String(), tc.Result)
		}
	}
}

const testFormatJSONReqBasicStr = `{"type":"request","auth":{"display_name":"","policies":["root"],"metadata":null},"request":{"operation":"write","path":"/foo","data":null,"remote_address":"127.0.0.1"},"error":"this is an error"}
`
