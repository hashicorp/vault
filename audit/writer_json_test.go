// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestFormatJSON_formatRequest(t *testing.T) {
	ss := newStaticSalt(t)

	expectedResultStr := fmt.Sprintf(testFormatJSONReqBasicStrFmt, ss.salt.GetIdentifiedHMAC("foo"))

	issueTime, _ := time.Parse(time.RFC3339, "2020-05-28T13:40:18-05:00")
	cases := map[string]struct {
		Auth        *logical.Auth
		Req         *logical.Request
		Err         error
		Prefix      string
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
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Hour * 4,
					IssueTime: issueTime,
				},
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
					"foo": {"bar"},
				},
			},
			errors.New("this is an error"),
			"",
			expectedResultStr,
		},
		"auth, request with prefix": {
			&logical.Auth{
				ClientToken:     "foo",
				Accessor:        "bar",
				EntityID:        "foobarentity",
				DisplayName:     "testtoken",
				NoDefaultPolicy: true,
				Policies:        []string{"root"},
				TokenType:       logical.TokenTypeService,
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Hour * 4,
					IssueTime: issueTime,
				},
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
					"foo": {"bar"},
				},
			},
			errors.New("this is an error"),
			"@cee: ",
			expectedResultStr,
		},
	}

	for name, tc := range cases {
		var buf bytes.Buffer
		cfg, err := NewFormatterConfig(WithHMACAccessor(false))
		require.NoError(t, err)
		f, err := NewEntryFormatter(cfg, ss)
		require.NoError(t, err)
		formatter := EntryFormatterWriter{
			Formatter: f,
			Writer: &JSONWriter{
				Prefix: tc.Prefix,
			},
			config: cfg,
		}

		in := &logical.LogInput{
			Auth:     tc.Auth,
			Request:  tc.Req,
			OuterErr: tc.Err,
		}

		err = formatter.FormatAndWriteRequest(namespace.RootContext(nil), &buf, in)
		require.NoErrorf(t, err, "bad: %s\nerr: %s", name, err)

		if !strings.HasPrefix(buf.String(), tc.Prefix) {
			t.Fatalf("no prefix: %s \n log: %s\nprefix: %s", name, expectedResultStr, tc.Prefix)
		}

		expectedJSON := new(RequestEntry)

		if err := jsonutil.DecodeJSON([]byte(expectedResultStr), &expectedJSON); err != nil {
			t.Fatalf("bad json: %s", err)
		}
		expectedJSON.Request.Namespace = &Namespace{ID: "root"}

		actualjson := new(RequestEntry)
		if err := jsonutil.DecodeJSON([]byte(buf.String())[len(tc.Prefix):], &actualjson); err != nil {
			t.Fatalf("bad json: %s", err)
		}

		expectedJSON.Time = actualjson.Time

		expectedBytes, err := json.Marshal(expectedJSON)
		if err != nil {
			t.Fatalf("unable to marshal json: %s", err)
		}

		if !strings.HasSuffix(strings.TrimSpace(buf.String()), string(expectedBytes)) {
			t.Fatalf(
				"bad: %s\nResult:\n\n%q\n\nExpected:\n\n%q",
				name, buf.String(), string(expectedBytes))
		}
	}
}

const testFormatJSONReqBasicStrFmt = `{"time":"2015-08-05T13:45:46Z","type":"request","auth":{"client_token":"%s","accessor":"bar","display_name":"testtoken","policies":["root"],"no_default_policy":true,"metadata":null,"entity_id":"foobarentity","token_type":"service", "token_ttl": 14400, "token_issue_time": "2020-05-28T13:40:18-05:00"},"request":{"operation":"update","path":"/foo","data":null,"wrap_ttl":60,"remote_address":"127.0.0.1","headers":{"foo":["bar"]}},"error":"this is an error"}
`
