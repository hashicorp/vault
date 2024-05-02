// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package puller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
)

var _ pluginSource = (*httpPluginSource)(nil)

type httpPluginSource struct {
	httpClient *http.Client
	baseURL    string
}

func newHTTPPluginSource(baseURL string) *httpPluginSource {
	return &httpPluginSource{
		httpClient: cleanhttp.DefaultClient(),
		baseURL:    baseURL,
	}
}

func (p *httpPluginSource) listMetadata(ctx context.Context, plugin string) ([]metadata, error) {
	url := fmt.Sprintf("%s/v1/releases/%s", p.baseURL, plugin)
	return httpGet[[]metadata](ctx, p.httpClient, url)
}

func (p *httpPluginSource) getMetadata(ctx context.Context, plugin, version string) (metadata, error) {
	url := fmt.Sprintf("%s/v1/releases/%s/%s", p.baseURL, plugin, version)
	return httpGet[metadata](ctx, p.httpClient, url)
}

func httpGet[T any](ctx context.Context, httpClient *http.Client, url string) (T, error) {
	var t T
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return t, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return t, errNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return t, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
	}

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&t); err != nil {
		return t, err
	}

	return t, nil
}

func (p *httpPluginSource) getContentReader(ctx context.Context, url string) (reader io.ReadCloser, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
	}

	return resp.Body, nil
}
