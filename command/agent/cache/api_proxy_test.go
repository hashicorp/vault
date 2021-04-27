package cache

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/go-test/deep"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestAPIProxy(t *testing.T) {
	cleanup, client, _, _ := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client: client,
		Logger: logging.NewVaultLogger(hclog.Trace),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := client.NewRequest("GET", "/v1/sys/health")
	req, err := r.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: req,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result api.HealthResponse
	err = jsonutil.DecodeJSONFromReader(resp.Response.Body, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Initialized || result.Sealed || result.Standby {
		t.Fatalf("bad sys/health response: %#v", result)
	}
}

func TestAPIProxy_queryParams(t *testing.T) {
	// Set up an agent that points to a standby node for this particular test
	// since it needs to proxy a /sys/health?standbyok=true request to a standby
	cleanup, client, _, _ := setupClusterAndAgentOnStandby(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client: client,
		Logger: logging.NewVaultLogger(hclog.Trace),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := client.NewRequest("GET", "/v1/sys/health")
	req, err := r.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}

	// Add a query parameter for testing
	q := req.URL.Query()
	q.Add("standbyok", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: req,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result api.HealthResponse
	err = jsonutil.DecodeJSONFromReader(resp.Response.Body, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Initialized || result.Sealed || !result.Standby {
		t.Fatalf("bad sys/health response: %#v", result)
	}

	if resp.Response.StatusCode != http.StatusOK {
		t.Fatalf("exptected standby to return 200, got: %v", resp.Response.StatusCode)
	}
}

func TestMergeStates(t *testing.T) {
	type testCase struct {
		name     string
		old      []string
		new      string
		expected []string
	}

	testCases := []testCase{
		{
			name:     "empty-old",
			old:      nil,
			new:      "v1:cid:1:0:",
			expected: []string{"v1:cid:1:0:"},
		},
		{
			name:     "old-smaller",
			old:      []string{"v1:cid:1:0:"},
			new:      "v1:cid:2:0:",
			expected: []string{"v1:cid:2:0:"},
		},
		{
			name:     "old-bigger",
			old:      []string{"v1:cid:2:0:"},
			new:      "v1:cid:1:0:",
			expected: []string{"v1:cid:2:0:"},
		},
		{
			name:     "mixed-single",
			old:      []string{"v1:cid:1:0:"},
			new:      "v1:cid:0:1:",
			expected: []string{"v1:cid:0:1:", "v1:cid:1:0:"},
		},
		{
			name:     "mixed-single-alt",
			old:      []string{"v1:cid:0:1:"},
			new:      "v1:cid:1:0:",
			expected: []string{"v1:cid:0:1:", "v1:cid:1:0:"},
		},
		{
			name:     "mixed-double",
			old:      []string{"v1:cid:0:1:", "v1:cid:1:0:"},
			new:      "v1:cid:2:0:",
			expected: []string{"v1:cid:0:1:", "v1:cid:2:0:"},
		},
		{
			name:     "newer-both",
			old:      []string{"v1:cid:0:1:", "v1:cid:1:0:"},
			new:      "v1:cid:2:1:",
			expected: []string{"v1:cid:2:1:"},
		},
	}

	b64enc := func(ss []string) []string {
		var ret []string
		for _, s := range ss {
			ret = append(ret, base64.StdEncoding.EncodeToString([]byte(s)))
		}
		return ret
	}
	b64dec := func(ss []string) []string {
		var ret []string
		for _, s := range ss {
			d, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				t.Fatal(err)
			}
			ret = append(ret, string(d))
		}
		return ret
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := b64dec(mergeStates(b64enc(tc.old), base64.StdEncoding.EncodeToString([]byte(tc.new))))
			if diff := deep.Equal(out, tc.expected); len(diff) != 0 {
				t.Errorf("got=%v, expected=%v, diff=%v", out, tc.expected, diff)
			}
		})
	}
}
