package audit

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"errors"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

func TestFormatJSONx_formatRequest(t *testing.T) {
	cases := map[string]struct {
		Auth     *logical.Auth
		Req      *logical.Request
		Err      error
		Result   string
		Expected string
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
			"",
			`<json:object name="auth"><json:string name="accessor"></json:string><json:string name="client_token"></json:string><json:string name="display_name"></json:string><json:null name="metadata" /><json:array name="policies"><json:string>root</json:string></json:array></json:object><json:string name="error">this is an error</json:string><json:object name="request"><json:string name="client_token"></json:string><json:null name="data" /><json:string name="id"></json:string><json:string name="operation">update</json:string><json:string name="path">/foo</json:string><json:string name="remote_address">127.0.0.1</json:string><json:number name="wrap_ttl">60</json:number></json:object><json:string name="type">request</json:string>`,
		},
	}

	for name, tc := range cases {
		var buf bytes.Buffer
		formatter := AuditFormatter{
			AuditFormatWriter: &JSONxFormatWriter{},
		}
		salter, _ := salt.NewSalt(nil, nil)
		config := FormatterConfig{
			Salt:     salter,
			OmitTime: true,
		}
		if err := formatter.FormatRequest(&buf, config, tc.Auth, tc.Req, tc.Err); err != nil {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}

		if strings.TrimSpace(buf.String()) != string(tc.Expected) {
			t.Fatalf(
				"bad: %s\nResult:\n\n'%s'\n\nExpected:\n\n'%s'",
				name, strings.TrimSpace(buf.String()), string(tc.Expected))
		}
	}
}
