package audit

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"errors"

	"fmt"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

func TestFormatJSONx_formatRequest(t *testing.T) {
	salter, err := salt.NewSalt(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	saltFunc := func() (*salt.Salt, error) {
		return salter, nil
	}

	fooSalted := salter.GetIdentifiedHMAC("foo")

	cases := map[string]struct {
		Auth        *logical.Auth
		Req         *logical.Request
		Err         error
		Prefix      string
		Result      string
		ExpectedStr string
	}{
		"auth, request": {
			&logical.Auth{ClientToken: "foo", Accessor: "bar", DisplayName: "testtoken", Policies: []string{"root"}},
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
				WrapInfo: &logical.RequestWrapInfo{
					TTL: 60 * time.Second,
				},
				Headers: map[string][]string{
					"foo": []string{"bar"},
				},
			},
			errors.New("this is an error"),
			"",
			"",
			fmt.Sprintf(`<json:object name="auth"><json:string name="accessor">bar</json:string><json:string name="client_token">%s</json:string><json:string name="display_name">testtoken</json:string><json:string name="entity_id"></json:string><json:null name="metadata" /><json:array name="policies"><json:string>root</json:string></json:array></json:object><json:string name="error">this is an error</json:string><json:object name="request"><json:string name="client_token"></json:string><json:string name="client_token_accessor"></json:string><json:null name="data" /><json:object name="headers"><json:array name="foo"><json:string>bar</json:string></json:array></json:object><json:string name="id"></json:string><json:string name="operation">update</json:string><json:string name="path">/foo</json:string><json:boolean name="policy_override">false</json:boolean><json:string name="remote_address">127.0.0.1</json:string><json:number name="wrap_ttl">60</json:number></json:object><json:string name="type">request</json:string>`,
				fooSalted),
		},
		"auth, request with prefix": {
			&logical.Auth{ClientToken: "foo", Accessor: "bar", DisplayName: "testtoken", Policies: []string{"root"}},
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
				WrapInfo: &logical.RequestWrapInfo{
					TTL: 60 * time.Second,
				},
				Headers: map[string][]string{
					"foo": []string{"bar"},
				},
			},
			errors.New("this is an error"),
			"",
			"@cee: ",
			fmt.Sprintf(`<json:object name="auth"><json:string name="accessor">bar</json:string><json:string name="client_token">%s</json:string><json:string name="display_name">testtoken</json:string><json:string name="entity_id"></json:string><json:null name="metadata" /><json:array name="policies"><json:string>root</json:string></json:array></json:object><json:string name="error">this is an error</json:string><json:object name="request"><json:string name="client_token"></json:string><json:string name="client_token_accessor"></json:string><json:null name="data" /><json:object name="headers"><json:array name="foo"><json:string>bar</json:string></json:array></json:object><json:string name="id"></json:string><json:string name="operation">update</json:string><json:string name="path">/foo</json:string><json:boolean name="policy_override">false</json:boolean><json:string name="remote_address">127.0.0.1</json:string><json:number name="wrap_ttl">60</json:number></json:object><json:string name="type">request</json:string>`,
				fooSalted),
		},
	}

	for name, tc := range cases {
		var buf bytes.Buffer
		formatter := AuditFormatter{
			AuditFormatWriter: &JSONxFormatWriter{
				Prefix:   tc.Prefix,
				SaltFunc: saltFunc,
			},
		}
		config := FormatterConfig{
			OmitTime:     true,
			HMACAccessor: false,
		}
		if err := formatter.FormatRequest(&buf, config, tc.Auth, tc.Req, tc.Err); err != nil {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}

		if !strings.HasPrefix(buf.String(), tc.Prefix) {
			t.Fatalf("no prefix: %s \n log: %s\nprefix: %s", name, tc.Result, tc.Prefix)
		}

		if !strings.HasSuffix(strings.TrimSpace(buf.String()), string(tc.ExpectedStr)) {
			t.Fatalf(
				"bad: %s\nResult:\n\n'%s'\n\nExpected:\n\n'%s'",
				name, strings.TrimSpace(buf.String()), string(tc.ExpectedStr))
		}
	}
}
