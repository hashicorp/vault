// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consts

// ProxyPathCacheClear is the path that the proxy will use as its cache-clear
// endpoint.
const ProxyPathCacheClear = "/proxy/v1/cache-clear"

// ProxyPathMetrics is the path the proxy will use to expose its internal
// metrics.
const ProxyPathMetrics = "/proxy/v1/metrics"

// ProxyPathQuit is the path that the proxy will use to trigger stopping it.
const ProxyPathQuit = "/proxy/v1/quit"
