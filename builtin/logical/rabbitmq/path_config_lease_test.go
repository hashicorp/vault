// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testhelpers/observations"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestBackend_config_lease_RU(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or
	config.StorageView = &logical.InmemStorage{}
	b := Backend()
	if err = b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"ttl":     "10h",
		"max_ttl": "20h",
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/lease",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQConnectionConfigWrite))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialCreateFail))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialCreateSuccess))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialRenew))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialRevoke))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQLeaseConfigRead))
	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeRabbitMQLeaseConfigWrite))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQRoleDelete))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQRoleRead))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQRoleWrite))

	configReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQConnectionConfigWrite))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialCreateFail))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialCreateSuccess))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialRenew))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQCredentialRevoke))
	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeRabbitMQLeaseConfigRead))
	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeRabbitMQLeaseConfigWrite))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQRoleDelete))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQRoleRead))
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeRabbitMQRoleWrite))

	if resp.Data["ttl"].(time.Duration) != 36000 {
		t.Fatalf("bad: ttl: expected:36000 actual:%d", resp.Data["ttl"].(time.Duration))
	}
	if resp.Data["max_ttl"].(time.Duration) != 72000 {
		t.Fatalf("bad: ttl: expected:72000 actual:%d", resp.Data["ttl"].(time.Duration))
	}
}
