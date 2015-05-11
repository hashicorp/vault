package audit

import (
	"bytes"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestFormatJSON_formatRequest(t *testing.T) {
	cases := map[string]struct {
		Auth   *logical.Auth
		Req    *logical.Request
		Result string
	}{
		"auth, request": {
			&logical.Auth{ClientToken: "foo", Policies: []string{"root"}},
			&logical.Request{
				Operation: logical.WriteOperation,
				Path:      "/foo",
			},
			testFormatJSONReqBasicStr,
		},
	}

	for name, tc := range cases {
		var buf bytes.Buffer
		var format FormatJSON
		if err := format.FormatRequest(&buf, tc.Auth, tc.Req); err != nil {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}

		if buf.String() != tc.Result {
			t.Fatalf(
				"bad: %s\nResult:\n\n%s\n\nExpected:\n\n%s",
				name, buf.String(), tc.Result)
		}
	}
}

const testFormatJSONReqBasicStr = `{"type":"request","auth":{"display_name":"","policies":["root"],"metadata":null},"request":{"operation":"write","path":"/foo","data":null}}
`
