package pprof

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"golang.org/x/net/http2"
)

func TestSysPprof(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	cases := []struct {
		name    string
		path    string
		seconds string
	}{
		{
			"index",
			"/v1/sys/pprof/",
			"",
		},
		{
			"cmdline",
			"/v1/sys/pprof/cmdline",
			"",
		},
		{
			"goroutine",
			"/v1/sys/pprof/goroutine",
			"",
		},
		{
			"heap",
			"/v1/sys/pprof/heap",
			"",
		},
		{
			"profile",
			"/v1/sys/pprof/profile",
			"1",
		},
		{
			"symbol",
			"/v1/sys/pprof/symbol",
			"",
		},
		{
			"trace",
			"/v1/sys/pprof/trace",
			"1",
		},
	}

	pprofRequest := func(path string, seconds string) {
		req := client.NewRequest("GET", path)
		if seconds != "" {
			req.Params.Set("seconds", seconds)
		}
		httpReq, err := req.ToHTTP()
		if err != nil {
			t.Fatal(err)
		}
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		httpRespBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		// Make sure that we don't get a error response
		httpResp := make(map[string]interface{})

		// Skip this error check since some endpoints return binary blobs, we
		// only care about the ok check right after as an existence check.
		_ = json.Unmarshal(httpRespBody, &httpResp)

		if _, ok := httpResp["errors"]; ok {
			t.Fatalf("unexpected error response: %v", httpResp["errors"])
		}

		if len(httpRespBody) == 0 {
			t.Fatal("no pprof index returned")
		}
		t.Log(string(httpRespBody))
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pprofRequest(tc.path, tc.seconds)
		})
	}
}
