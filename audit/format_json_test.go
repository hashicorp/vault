package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"errors"

	"fmt"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

func TestFormatJSON_formatRequest(t *testing.T) {
	salter, err := salt.NewSalt(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	saltFunc := func(context.Context) (*salt.Salt, error) {
		return salter, nil
	}

	expectedResultStr := fmt.Sprintf(testFormatJSONReqBasicStrFmt, salter.GetIdentifiedHMAC("foo"))

	cases := map[string]struct {
		Auth        *logical.Auth
		Req         *logical.Request
		Err         error
		Prefix      string
		ExpectedStr string
	}{
		"auth, request": {
			&logical.Auth{
				ClientToken: "foo",
				Accessor:    "bar",
				DisplayName: "testtoken",
				Policies:    []string{"root"},
				TokenType:   logical.TokenTypeService,
			},
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
			expectedResultStr,
		},
		"auth, request with prefix": {
			&logical.Auth{
				ClientToken: "foo",
				Accessor:    "bar",
				DisplayName: "testtoken",
				Policies:    []string{"root"},
				TokenType:   logical.TokenTypeService,
			},
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
			"@cee: ",
			expectedResultStr,
		},
	}

	for name, tc := range cases {
		var buf bytes.Buffer
		formatter := AuditFormatter{
			AuditFormatWriter: &JSONFormatWriter{
				Prefix:   tc.Prefix,
				SaltFunc: saltFunc,
			},
		}
		config := FormatterConfig{
			HMACAccessor: false,
		}
		in := &LogInput{
			Auth:     tc.Auth,
			Request:  tc.Req,
			OuterErr: tc.Err,
		}
		if err := formatter.FormatRequest(namespace.RootContext(nil), &buf, config, in); err != nil {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}

		if !strings.HasPrefix(buf.String(), tc.Prefix) {
			t.Fatalf("no prefix: %s \n log: %s\nprefix: %s", name, expectedResultStr, tc.Prefix)
		}

		var expectedjson = new(AuditRequestEntry)

		if err := jsonutil.DecodeJSON([]byte(expectedResultStr), &expectedjson); err != nil {
			t.Fatalf("bad json: %s", err)
		}
		expectedjson.Request.Namespace = AuditNamespace{ID: "root"}

		var actualjson = new(AuditRequestEntry)
		if err := jsonutil.DecodeJSON([]byte(buf.String())[len(tc.Prefix):], &actualjson); err != nil {
			t.Fatalf("bad json: %s", err)
		}

		expectedjson.Time = actualjson.Time

		expectedBytes, err := json.Marshal(expectedjson)
		if err != nil {
			t.Fatalf("unable to marshal json: %s", err)
		}

		if !strings.HasSuffix(strings.TrimSpace(buf.String()), string(expectedBytes)) {
			t.Fatalf(
				"bad: %s\nResult:\n\n'%s'\n\nExpected:\n\n'%s'",
				name, buf.String(), string(expectedBytes))
		}
	}
}

const testFormatJSONReqBasicStrFmt = `{"time":"2015-08-05T13:45:46Z","type":"request","auth":{"client_token":"%s","accessor":"bar","display_name":"testtoken","policies":["root"],"metadata":null,"entity_id":"","token_type":"service"},"request":{"operation":"update","path":"/foo","data":null,"wrap_ttl":60,"remote_address":"127.0.0.1","headers":{"foo":["bar"]}},"error":"this is an error"}
`
