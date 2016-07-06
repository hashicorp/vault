package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"errors"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
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
				Operation: logical.UpdateOperation,
				Path:      "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
				WrapTTL: 60 * time.Second,
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

		var expectedjson = new(JSONRequestEntry)
		if err := jsonutil.DecodeJSON([]byte(tc.Result), &expectedjson); err != nil {
			t.Fatalf("bad json: %s", err)
		}

		var actualjson = new(JSONRequestEntry)
		if err := jsonutil.DecodeJSON([]byte(buf.String()), &actualjson); err != nil {
			t.Fatalf("bad json: %s", err)
		}

		expectedjson.Time = actualjson.Time

		expectedBytes, err := json.Marshal(expectedjson)
		if err != nil {
			t.Fatalf("unable to marshal json: %s", err)
		}

		if strings.TrimSpace(buf.String()) != string(expectedBytes) {
			t.Fatalf(
				"bad: %s\nResult:\n\n'%s'\n\nExpected:\n\n'%s'",
				name, buf.String(), string(expectedBytes))
		}
	}
}

const testFormatJSONReqBasicStr = `{"time":"2015-08-05T13:45:46Z","type":"request","auth":{"display_name":"","policies":["root"],"metadata":null},"request":{"operation":"update","path":"/foo","data":null,"wrap_ttl":60,"remote_address":"127.0.0.1"},"error":"this is an error"}
`
