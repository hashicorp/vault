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

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

type testOptMarshaler struct {
	input string
}

func (t *testOptMarshaler) MarshalJSONWithOptions(opts *logical.MarshalOptions) ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, opts.ValueHasher(t.input))), nil
}

func TestFormatJSON_formatRequest(t *testing.T) {
	salter, err := salt.NewSalt(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	saltFunc := func(context.Context) (*salt.Salt, error) {
		return salter, nil
	}

	expectedResultStr := fmt.Sprintf(testFormatJSONReqBasicStrFmt, salter.GetIdentifiedHMAC("foo"))

	genreq := &testOptMarshaler{"generic"}
	expectedGenreqResultStr := fmt.Sprintf(testFormatJSONReqGenericStrFmt, salter.GetIdentifiedHMAC("generic"))

	cases := map[string]struct {
		Auth        *logical.Auth
		Req         interface{}
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
		"generic request": {
			nil,
			genreq,
			errors.New("this is an error"),
			"",
			expectedGenreqResultStr,
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
		in := &logical.LogInput{
			Auth:     tc.Auth,
			Request:  tc.Req,
			OuterErr: tc.Err,
		}
		if err := formatter.FormatRequest(namespace.RootContext(nil), &buf, config, in); err != nil {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}

		if !strings.HasPrefix(buf.String(), tc.Prefix) {
			t.Fatalf("no prefix: %s \n log: %s\nprefix: %s", name, tc.ExpectedStr, tc.Prefix)
		}

		var expected = new(AuditRequestEntry)

		if err := jsonutil.DecodeJSON([]byte(tc.ExpectedStr), &expected); err != nil {
			t.Fatalf("bad json: %s", err)
		}
		if _, ok := tc.Req.(*logical.Request); ok {
			expected.Request.Namespace = AuditNamespace{ID: "root"}
		}

		var actual = new(AuditRequestEntry)
		if err := jsonutil.DecodeJSON([]byte(buf.String())[len(tc.Prefix):], &actual); err != nil {
			t.Fatalf("bad json: %s", err)
		}

		expected.Time = actual.Time

		expectedBytes, err := json.Marshal(expected)
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
const testFormatJSONReqGenericStrFmt = `{"time":"2015-08-05T13:45:46Z","type":"","auth":{"client_token":""},"request":{"data":"%s"},"error":"this is an error"}
`
