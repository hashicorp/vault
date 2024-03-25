// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package kafka

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/plugins/event"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKafka_SendOneMessage tests the basic happy path with a dev single-node
// Kafka cluster within a docker container.
func TestKafka_SendOneMessage(t *testing.T) {
	addr := prepapreKafkaTestContainer(t)

	ctx := context.Background()
	plugin, err := New(ctx)
	require.NoError(t, err)

	subID, err := uuid.GenerateUUID()
	require.NoError(t, err)

	err = plugin.Subscribe(ctx, &event.SubscribeRequest{
		SubscriptionID: subID,
		Config: map[string]interface{}{
			"addresses":         []string{addr},
			"topic":             "test-topic",
			"auto_create_topic": true,
			"tls_disabled":      true,
		},
	})
	require.NoError(t, err)

	// Re-subscribe - should not error.
	err = plugin.Subscribe(ctx, &event.SubscribeRequest{
		SubscriptionID: subID + "2",
		Config: map[string]interface{}{
			"addresses":         []string{addr},
			"topic":             "test-topic",
			"auto_create_topic": true,
			"tls_disabled":      true,
		},
	})
	require.NoError(t, err)

	err = plugin.Send(ctx, &event.SendRequest{
		SubscriptionID: subID,
		EventJSON:      "{}",
	})
	require.NoError(t, err)

	// Now read the event back from kafka.
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{addr},
		Topic:   "test-topic",
	})

	msg, err := reader.ReadMessage(ctx)
	require.NoError(t, err)
	assert.Equal(t, []byte("{}"), msg.Value)

	err = plugin.Unsubscribe(ctx, &event.UnsubscribeRequest{
		SubscriptionID: subID,
	})
	require.NoError(t, err)
}

func prepapreKafkaTestContainer(t *testing.T) string {
	// The cluster advertises its own address to clients, so we need to find a
	// host port to provide as config input _before_ starting the container.
	hostPort := findFreePort(t)

	// Configure a single node dev container which will act as both the controller
	// and the broker.
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "bitnami/kafka",
		ImageTag:      "3.3.2",
		ContainerName: "kafka",
		Ports:         []string{"9092/tcp"},
		PortBindings: nat.PortMap{
			"9092/tcp": []nat.PortBinding{{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(hostPort),
			}},
		},
		Env: []string{
			"KAFKA_CFG_NODE_ID=0",
			"KAFKA_CFG_PROCESS_ROLES=controller,broker",
			"KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093",
			fmt.Sprintf("KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:%d", hostPort),
			"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
			"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@localhost:9093",
			"KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		var connected bool
		for i := 0; i < 100; i++ {
			conn, err := kafka.DialLeader(ctx, "tcp", fmt.Sprintf("%s:%d", host, port), "test-topic", 0)
			if err == nil {
				connected = true
				conn.Close()
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		if !connected {
			return nil, errors.New("failed to connect to kafka")
		}

		return docker.NewServiceHostPort(host, port), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(svc.Cleanup)

	return svc.Config.Address()
}

func findFreePort(t *testing.T) int {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := ln.Close(); err != nil {
		t.Fatal(err)
	}
	return ln.Addr().(*net.TCPAddr).Port
}
