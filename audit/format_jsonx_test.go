package audit

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"errors"

	"fmt"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestFormatJSONx_formatRequest(t *testing.T) {
	salter, err := salt.NewSalt(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	saltFunc := func(context.Context) (*salt.Salt, error) {
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
			&logical.Auth{
				ClientToken:     "foo",
				Accessor:        "bar",
				DisplayName:     "testtoken",
				EntityID:        "foobarentity",
				NoDefaultPolicy: true,
				Policies:        []string{"root"},
				TokenType:       logical.TokenTypeService,
			},
			&logical.Request{
				ID:                  "request",
				ClientToken:         "foo",
				ClientTokenAccessor: "bar",
				Operation:           logical.UpdateOperation,
				Path:                "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
				WrapInfo: &logical.RequestWrapInfo{
					TTL: 60 * time.Second,
				},
				Headers: map[string][]string{
					"foo": []string{"bar"},
				},
				PolicyOverride: true,
			},
			errors.New("this is an error"),
			"",
			"",
			fmt.Sprintf(`<json:object name="auth"><json:string name="accessor">bar</json:string><json:string name="client_token">%s</json:string><json:string name="display_name">testtoken</json:string><json:string name="entity_id">foobarentity</json:string><json:boolean name="no_default_policy">true</json:boolean><json:array name="policies"><json:string>root</json:string></json:array><json:string name="token_type">service</json:string></json:object><json:string name="error">this is an error</json:string><json:object name="request"><json:string name="client_token">%s</json:string><json:string name="client_token_accessor">bar</json:string><json:object name="headers"><json:array name="foo"><json:string>bar</json:string></json:array></json:object><json:string name="id">request</json:string><json:object name="namespace"><json:string name="id">root</json:string></json:object><json:string name="operation">update</json:string><json:string name="path">/foo</json:string><json:boolean name="policy_override">true</json:boolean><json:string name="remote_address">127.0.0.1</json:string><json:number name="wrap_ttl">60</json:number></json:object><json:string name="type">request</json:string>`,
				fooSalted, fooSalted),
		},
		"auth, request with prefix": {
			&logical.Auth{
				ClientToken:     "foo",
				Accessor:        "bar",
				DisplayName:     "testtoken",
				NoDefaultPolicy: true,
				EntityID:        "foobarentity",
				Policies:        []string{"root"},
				TokenType:       logical.TokenTypeService,
			},
			&logical.Request{
				ID:                  "request",
				ClientToken:         "foo",
				ClientTokenAccessor: "bar",
				Operation:           logical.UpdateOperation,
				Path:                "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
				WrapInfo: &logical.RequestWrapInfo{
					TTL: 60 * time.Second,
				},
				Headers: map[string][]string{
					"foo": []string{"bar"},
				},
				PolicyOverride: true,
			},
			errors.New("this is an error"),
			"",
			"@cee: ",
			fmt.Sprintf(`<json:object name="auth"><json:string name="accessor">bar</json:string><json:string name="client_token">%s</json:string><json:string name="display_name">testtoken</json:string><json:string name="entity_id">foobarentity</json:string><json:boolean name="no_default_policy">true</json:boolean><json:array name="policies"><json:string>root</json:string></json:array><json:string name="token_type">service</json:string></json:object><json:string name="error">this is an error</json:string><json:object name="request"><json:string name="client_token">%s</json:string><json:string name="client_token_accessor">bar</json:string><json:object name="headers"><json:array name="foo"><json:string>bar</json:string></json:array></json:object><json:string name="id">request</json:string><json:object name="namespace"><json:string name="id">root</json:string></json:object><json:string name="operation">update</json:string><json:string name="path">/foo</json:string><json:boolean name="policy_override">true</json:boolean><json:string name="remote_address">127.0.0.1</json:string><json:number name="wrap_ttl">60</json:number></json:object><json:string name="type">request</json:string>`,
				fooSalted, fooSalted),
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
		in := &logical.LogInput{
			Auth:     tc.Auth,
			Request:  tc.Req,
			OuterErr: tc.Err,
		}
		if err := formatter.FormatRequest(namespace.RootContext(nil), &buf, config, in); err != nil {
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
