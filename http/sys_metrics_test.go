// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func TestSysMetricsUnauthenticated(t *testing.T) {
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)
	conf := &vault.CoreConfig{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		MetricsHelper:   metricsutil.NewMetricsHelper(inm, true),
	}
	core, _, token := vault.TestCoreUnsealedWithConfig(t, conf)
	ln, addr := TestServer(t, core)
	TestServerAuth(t, addr, token)

	// Default: Only authenticated access
	resp := testHttpGet(t, "", addr+"/v1/sys/metrics")
	testResponseStatus(t, resp, 403)
	resp = testHttpGet(t, token, addr+"/v1/sys/metrics")
	testResponseStatus(t, resp, 200)

	// Close listener
	ln.Close()

	// Setup new custom listener with unauthenticated metrics access
	ln, addr = TestListener(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			Telemetry: configutil.ListenerTelemetry{
				UnauthenticatedMetricsAccess: true,
			},
		},
	}
	TestServerWithListenerAndProperties(t, ln, addr, core, props)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Test without token
	resp = testHttpGet(t, "", addr+"/v1/sys/metrics")
	testResponseStatus(t, resp, 200)

	// Should also work with token
	resp = testHttpGet(t, token, addr+"/v1/sys/metrics")
	testResponseStatus(t, resp, 200)

	// Test if prometheus response is correct
	resp = testHttpGet(t, "", addr+"/v1/sys/metrics?format=prometheus")
	testResponseStatus(t, resp, 200)
}

func TestSysPProfUnauthenticated(t *testing.T) {
	conf := &vault.CoreConfig{}
	core, _, token := vault.TestCoreUnsealedWithConfig(t, conf)
	ln, addr := TestServer(t, core)
	TestServerAuth(t, addr, token)

	// Default: Only authenticated access
	resp := testHttpGet(t, "", addr+"/v1/sys/pprof/cmdline")
	testResponseStatus(t, resp, 403)
	resp = testHttpGet(t, token, addr+"/v1/sys/pprof/cmdline")
	testResponseStatus(t, resp, 200)

	// Close listener
	ln.Close()

	// Setup new custom listener with unauthenticated metrics access
	ln, addr = TestListener(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			Profiling: configutil.ListenerProfiling{
				UnauthenticatedPProfAccess: true,
			},
		},
	}
	TestServerWithListenerAndProperties(t, ln, addr, core, props)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Test without token
	resp = testHttpGet(t, "", addr+"/v1/sys/pprof/cmdline")
	testResponseStatus(t, resp, 200)

	// Should also work with token
	resp = testHttpGet(t, token, addr+"/v1/sys/pprof/cmdline")
	testResponseStatus(t, resp, 200)
}
