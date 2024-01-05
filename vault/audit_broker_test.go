// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/sha256"
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/audit/syslog"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// testAuditBackend will create an audit.Backend (which expects to use the eventlogger).
func testAuditBackend(t *testing.T, path string, config map[string]string) audit.Backend {
	t.Helper()

	headersCfg := &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
		view:    nil,
	}

	view := &logical.InmemStorage{}
	se := &logical.StorageEntry{Key: "salt", Value: []byte("juan")}
	err := view.Put(context.Background(), se)
	require.NoError(t, err)

	cfg := &audit.BackendConfig{
		SaltView: view,
		SaltConfig: &salt.Config{
			HMAC:     sha256.New,
			HMACType: "hmac-sha256",
		},
		Config:    config,
		MountPath: path,
	}

	be, err := syslog.Factory(context.Background(), cfg, true, headersCfg)
	require.NoError(t, err)
	require.NotNil(t, be)

	return be
}

// TestAuditBroker_Register_SuccessThresholdSinks tests that we are able to
// correctly identify what the required success threshold sinks value on the
// eventlogger broker should be set to.
// We expect:
// * 0 for only filtered backends
// * 1 for any other combination
func TestAuditBroker_Register_SuccessThresholdSinks(t *testing.T) {
	t.Parallel()
	l := corehelpers.NewTestLogger(t)
	a, err := NewAuditBroker(l, true)
	require.NoError(t, err)
	require.NotNil(t, a)

	filterBackend := testAuditBackend(t, "b1-filter", map[string]string{"filter": "foo == bar"})
	noFilterBackend := testAuditBackend(t, "b2-no-filter", map[string]string{})

	// Should be set to 0 for required sinks (and not found, as we've never registered before).
	res, ok := a.broker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.False(t, ok)
	require.Equal(t, 0, res)

	// Register the filtered backend first, this shouldn't change the
	// success threshold sinks to 1 as we can't guarantee any device yet.
	err = a.Register("b1-filter", filterBackend, false)
	require.NoError(t, err)

	// Check the SuccessThresholdSinks (we expect 0 still, but found).
	res, ok = a.broker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, ok)
	require.Equal(t, 0, res)

	// Register the non-filtered backend second, this should mean we
	// can rely on guarantees from the broker again.
	err = a.Register("b2-no-filter", noFilterBackend, false)
	require.NoError(t, err)

	// Check the SuccessThresholdSinks (we expect 1 now).
	res, ok = a.broker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, ok)
	require.Equal(t, 1, res)
}

// TestAuditBroker_Deregister_SuccessThresholdSinks tests that we are able to
// correctly identify what the required success threshold sinks value on the
// eventlogger broker should be set to when deregistering audit backends.
// We expect:
// * 0 for only filtered backends
// * 1 for any other combination
func TestAuditBroker_Deregister_SuccessThresholdSinks(t *testing.T) {
	t.Parallel()
	l := corehelpers.NewTestLogger(t)
	a, err := NewAuditBroker(l, true)
	require.NoError(t, err)
	require.NotNil(t, a)

	filterBackend := testAuditBackend(t, "b1-filter", map[string]string{"filter": "foo == bar"})
	noFilterBackend := testAuditBackend(t, "b2-no-filter", map[string]string{})

	err = a.Register("b1-filter", filterBackend, false)
	require.NoError(t, err)
	err = a.Register("b2-no-filter", noFilterBackend, false)
	require.NoError(t, err)

	// We have a mix of filtered and non-filtered backends, so the
	// successThresholdSinks should be 1.
	res, ok := a.broker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, ok)
	require.Equal(t, 1, res)

	// Deregister the non-filtered backend, there is one filtered backend left,
	// so the successThresholdSinks should be 0.
	err = a.Deregister(context.Background(), "b2-no-filter")
	require.NoError(t, err)
	res, ok = a.broker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, ok)
	require.Equal(t, 0, res)

	// Deregister the last backend, disabling audit. The value of
	// successThresholdSinks should still be 0.
	err = a.Deregister(context.Background(), "b1-filter")
	require.NoError(t, err)
	res, ok = a.broker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, ok)
	require.Equal(t, 0, res)

	// Re-register a backend that doesn't use filtering.
	err = a.Register("b2-no-filter", noFilterBackend, false)
	require.NoError(t, err)
	res, ok = a.broker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, ok)
	require.Equal(t, 1, res)
}

// TestAuditBroker_Register_Fallback ensures we can register a fallback device.
func TestAuditBroker_Register_Fallback(t *testing.T) {
	t.Parallel()

	l := corehelpers.NewTestLogger(t)
	a, err := NewAuditBroker(l, true)
	require.NoError(t, err)
	require.NotNil(t, a)

	path := "juan/"
	fallbackBackend := testAuditBackend(t, path, map[string]string{"fallback": "true"})
	err = a.Register(path, fallbackBackend, false)
	require.NoError(t, err)
	require.True(t, a.fallbackBroker.IsAnyPipelineRegistered(eventlogger.EventType(event.AuditType.String())))
	require.Equal(t, path, a.fallbackName)
	threshold, found := a.fallbackBroker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, found)
	require.Equal(t, 1, threshold)
}

// TestAuditBroker_Register_FallbackMultiple tests that trying to register more
// than a single fallback device results in the correct error.
func TestAuditBroker_Register_FallbackMultiple(t *testing.T) {
	t.Parallel()

	l := corehelpers.NewTestLogger(t)
	a, err := NewAuditBroker(l, true)
	require.NoError(t, err)
	require.NotNil(t, a)

	path1 := "juan1/"
	fallbackBackend1 := testAuditBackend(t, path1, map[string]string{"fallback": "true"})
	err = a.Register(path1, fallbackBackend1, false)
	require.NoError(t, err)
	require.True(t, a.fallbackBroker.IsAnyPipelineRegistered(eventlogger.EventType(event.AuditType.String())))
	require.Equal(t, path1, a.fallbackName)

	path2 := "juan2/"
	fallbackBackend2 := testAuditBackend(t, path2, map[string]string{"fallback": "true"})
	err = a.Register(path1, fallbackBackend2, false)
	require.Error(t, err)
	require.EqualError(t, err, "vault.(AuditBroker).Register: backend already registered 'juan1/'")
	require.True(t, a.fallbackBroker.IsAnyPipelineRegistered(eventlogger.EventType(event.AuditType.String())))
	require.Equal(t, path1, a.fallbackName)
}

// TestAuditBroker_Deregister_Fallback ensures that we can deregister a fallback
// device successfully.
func TestAuditBroker_Deregister_Fallback(t *testing.T) {
	t.Parallel()

	l := corehelpers.NewTestLogger(t)
	a, err := NewAuditBroker(l, true)
	require.NoError(t, err)
	require.NotNil(t, a)

	path := "juan/"
	fallbackBackend := testAuditBackend(t, path, map[string]string{"fallback": "true"})
	err = a.Register(path, fallbackBackend, false)
	require.NoError(t, err)
	require.True(t, a.fallbackBroker.IsAnyPipelineRegistered(eventlogger.EventType(event.AuditType.String())))
	require.Equal(t, path, a.fallbackName)

	threshold, found := a.fallbackBroker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, found)
	require.Equal(t, 1, threshold)

	err = a.Deregister(context.Background(), path)
	require.NoError(t, err)
	require.False(t, a.fallbackBroker.IsAnyPipelineRegistered(eventlogger.EventType(event.AuditType.String())))
	require.Equal(t, "", a.fallbackName)

	threshold, found = a.fallbackBroker.SuccessThresholdSinks(eventlogger.EventType(event.AuditType.String()))
	require.True(t, found)
	require.Equal(t, 0, threshold)
}

// TestAuditBroker_Deregister_Multiple ensures that we can call deregister multiple
// times without issue if is no matching backend registered.
func TestAuditBroker_Deregister_Multiple(t *testing.T) {
	t.Parallel()

	l := corehelpers.NewTestLogger(t)
	a, err := NewAuditBroker(l, true)
	require.NoError(t, err)
	require.NotNil(t, a)

	err = a.Deregister(context.Background(), "foo")
	require.NoError(t, err)

	err = a.Deregister(context.Background(), "foo2")
	require.NoError(t, err)
}

// TestAuditBroker_Register_MultipleFails checks for failure when we try to
// re-register an audit backend.
func TestAuditBroker_Register_MultipleFails(t *testing.T) {
	t.Parallel()

	l := corehelpers.NewTestLogger(t)
	a, err := NewAuditBroker(l, true)
	require.NoError(t, err)
	require.NotNil(t, a)

	path := "b2-no-filter"
	noFilterBackend := testAuditBackend(t, path, map[string]string{})

	err = a.Register(path, noFilterBackend, false)
	require.NoError(t, err)

	err = a.Register(path, noFilterBackend, false)
	require.Error(t, err)
	require.EqualError(t, err, "vault.(AuditBroker).Register: backend already registered 'b2-no-filter'")
}
