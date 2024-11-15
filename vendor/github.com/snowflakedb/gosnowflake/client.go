// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// InternalClient is implemented by HTTPClient
type InternalClient interface {
	Get(context.Context, *url.URL, map[string]string, time.Duration) (*http.Response, error)
	Post(context.Context, *url.URL, map[string]string, []byte, time.Duration, currentTimeProvider) (*http.Response, error)
}

type httpClient struct {
	sr *snowflakeRestful
}

func (cli *httpClient) Get(
	ctx context.Context,
	url *url.URL,
	headers map[string]string,
	timeout time.Duration) (*http.Response, error) {
	return cli.sr.FuncGet(ctx, cli.sr, url, headers, timeout)
}

func (cli *httpClient) Post(
	ctx context.Context,
	url *url.URL,
	headers map[string]string,
	body []byte,
	timeout time.Duration,
	currentTimeProvider currentTimeProvider) (*http.Response, error) {
	return cli.sr.FuncPost(ctx, cli.sr, url, headers, body, timeout, currentTimeProvider, nil)
}
