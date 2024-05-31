// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pprof

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func SysPprof_Test(t *testing.T, cluster testcluster.VaultCluster) {
	nodes := cluster.Nodes()
	client := nodes[0].APIClient()

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = nodes[0].TLSConfig()
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

	pprofRequest := func(t *testing.T, path string, seconds string) {
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

		httpResp := make(map[string]interface{})

		// Skip this error check since some endpoints return binary blobs, we
		// only care about the ok check right after as an existence check.
		_ = json.Unmarshal(httpRespBody, &httpResp)

		// Make sure that we don't get a error response
		if _, ok := httpResp["errors"]; ok {
			t.Fatalf("unexpected error response: %v", httpResp["errors"])
		}

		if len(httpRespBody) == 0 {
			t.Fatal("no pprof index returned")
		}
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pprofRequest(t, tc.path, tc.seconds)
		})
	}
}

func SysPprof_Standby_Test(t *testing.T, cluster testcluster.VaultCluster) {
	pprof := func(client *api.Client) (string, error) {
		req := client.NewRequest("GET", "/v1/sys/pprof/cmdline")
		resp, err := client.RawRequestWithContext(context.Background(), req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		return string(data), err
	}

	cmdline, err := pprof(cluster.Nodes()[0].APIClient())
	require.Nil(t, err)
	require.NotEmpty(t, cmdline)
	t.Log(cmdline)

	cmdline, err = pprof(cluster.Nodes()[1].APIClient())
	require.Nil(t, err)
	require.NotEmpty(t, cmdline)
	t.Log(cmdline)
}
