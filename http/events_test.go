// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/experiments"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/stretchr/testify/assert"
	"nhooyr.io/websocket"
)

// TestEventsSubscribe tests the websocket endpoint for subscribing to events
// by generating some events.
func TestEventsSubscribe(t *testing.T) {
	core := vault.TestCoreWithConfig(t, &vault.CoreConfig{
		Experiments: []string{experiments.VaultExperimentEventsAlpha1},
	})

	ln, addr := TestServer(t, core)
	defer ln.Close()

	// unseal the core
	keys, token := vault.TestCoreInit(t, core)
	for _, key := range keys {
		_, err := core.Unseal(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	stop := atomic.Bool{}

	const eventType = "abc"

	// send some events
	go func() {
		for !stop.Load() {
			id, err := uuid.GenerateUUID()
			if err != nil {
				core.Logger().Info("Error generating UUID, exiting sender", "error", err)
			}
			pluginInfo := &logical.EventPluginInfo{
				MountPath: "secret",
			}
			err = core.Events().SendEventInternal(namespace.RootContext(context.Background()), namespace.RootNamespace, pluginInfo, logical.EventType(eventType), &logical.EventData{
				Id:        id,
				Metadata:  nil,
				EntityIds: nil,
				Note:      "testing",
			})
			if err != nil {
				core.Logger().Info("Error sending event, exiting sender", "error", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	t.Cleanup(func() {
		stop.Store(true)
	})

	ctx := context.Background()
	wsAddr := strings.Replace(addr, "http", "ws", 1)

	testCases := []struct {
		json bool
	}{{true}, {false}}

	for _, testCase := range testCases {
		url := fmt.Sprintf("%s/v1/sys/events/subscribe/%s?namespaces=ns1&namespaces=ns*&json=%v", wsAddr, eventType, testCase.json)
		conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
			HTTPHeader: http.Header{"x-vault-token": []string{token}},
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			conn.Close(websocket.StatusNormalClosure, "")
		})

		_, msg, err := conn.Read(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if testCase.json {
			event := map[string]interface{}{}
			err = json.Unmarshal(msg, &event)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(msg))
			data := event["data"].(map[string]interface{})
			if actualType := data["event_type"].(string); actualType != eventType {
				t.Fatalf("Expeced event type %s, got %s", eventType, actualType)
			}
			pluginInfo, ok := data["plugin_info"].(map[string]interface{})
			if !ok || pluginInfo == nil {
				t.Fatalf("No plugin_info object: %v", data)
			}
			mountPath, ok := pluginInfo["mount_path"].(string)
			if !ok || mountPath != "secret" {
				t.Fatalf("Wrong mount_path: %v", data)
			}
			innerEvent := data["event"].(map[string]interface{})
			if innerEvent["id"].(string) != event["id"].(string) {
				t.Fatalf("IDs don't match, expected %s, got %s", innerEvent["id"].(string), event["id"].(string))
			}
			if innerEvent["note"].(string) != "testing" {
				t.Fatalf("Expected 'testing', got %s", innerEvent["note"].(string))
			}

			checkRequiredCloudEventsFields(t, event)
		}
	}
}

func TestNamespacePrepend(t *testing.T) {
	testCases := []struct {
		requestNs string
		patterns  []string
		result    []string
	}{
		{"", []string{"ns*"}, []string{"", "ns*"}},
		{"ns1", []string{"ns*"}, []string{"ns1", "ns1/ns*"}},
		{"ns1", []string{"ns1*"}, []string{"ns1", "ns1/ns1*"}},
		{"ns1", []string{"ns1/*"}, []string{"ns1", "ns1/ns1/*"}},
		{"", []string{"ns1/ns13", "ns1/other"}, []string{"", "ns1/ns13", "ns1/other"}},
		{"ns1", []string{"ns1/ns13", "ns1/other"}, []string{"ns1", "ns1/ns1/ns13", "ns1/ns1/other"}},
		{"", []string{""}, []string{""}},
		{"", nil, []string{""}},
		{"ns1", []string{""}, []string{"ns1"}},
		{"ns1", []string{"", ""}, []string{"ns1"}},
		{"ns1", []string{"ns1"}, []string{"ns1", "ns1/ns1"}},
		{"", []string{"*"}, []string{"", "*"}},
		{"ns1", []string{"*"}, []string{"ns1", "ns1/*"}},
		{"", []string{"ns1/ns13*", "ns2"}, []string{"", "ns1/ns13*", "ns2"}},
		{"ns1", []string{"ns1/ns13*", "ns2"}, []string{"ns1", "ns1/ns1/ns13*", "ns1/ns2"}},
		{"", []string{"ns*", "ns1"}, []string{"", "ns*", "ns1"}},
		{"ns1", []string{"ns*", "ns1"}, []string{"ns1", "ns1/ns*", "ns1/ns1"}},
		{"ns1", []string{"ns1*", "ns1"}, []string{"ns1", "ns1/ns1*", "ns1/ns1"}},
		{"ns1", []string{"ns1/*", "ns1"}, []string{"ns1", "ns1/ns1/*", "ns1/ns1"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.requestNs+" "+strings.Join(testCase.patterns, " "), func(t *testing.T) {
			result := prependNamespacePatterns(testCase.patterns, &namespace.Namespace{ID: testCase.requestNs, Path: testCase.requestNs})
			assert.Equal(t, testCase.result, result)
		})
	}
}

func checkRequiredCloudEventsFields(t *testing.T, event map[string]interface{}) {
	t.Helper()
	for _, attr := range []string{"id", "source", "specversion", "type"} {
		if v, ok := event[attr]; !ok {
			t.Errorf("Missing attribute %s", attr)
		} else if str, ok := v.(string); !ok {
			t.Errorf("Expected %s to be string but got %T", attr, v)
		} else if str == "" {
			t.Errorf("%s was empty string", attr)
		}
	}
}

// TestEventsSubscribeAuth tests that unauthenticated and unauthorized subscriptions
// fail correctly.
func TestEventsSubscribeAuth(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	// unseal the core
	keys, root := vault.TestCoreInit(t, core)
	for _, key := range keys {
		_, err := core.Unseal(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	var nonPrivilegedToken string
	// Fetch a valid non privileged token.
	{
		config := api.DefaultConfig()
		config.Address = addr

		client, err := api.NewClient(config)
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(root)

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{Policies: []string{"default"}})
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken == "" {
			t.Fatal("Failed to fetch a non privileged token")
		}
		nonPrivilegedToken = secret.Auth.ClientToken
	}

	ctx := context.Background()
	wsAddr := strings.Replace(addr, "http", "ws", 1)

	// Get a 403 with no token.
	_, resp, err := websocket.Dial(ctx, wsAddr+"/v1/sys/events/subscribe/abc", nil)
	if err == nil {
		t.Error("Expected websocket error but got none")
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %+v", resp)
	}

	// Get a 403 with a non privileged token.
	_, resp, err = websocket.Dial(ctx, wsAddr+"/v1/sys/events/subscribe/abc", &websocket.DialOptions{
		HTTPHeader: http.Header{"x-vault-token": []string{nonPrivilegedToken}},
	})
	if err == nil {
		t.Error("Expected websocket error but got none")
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %+v", resp)
	}
}

func TestCanForwardEventConnections(t *testing.T) {
	// Run again with in-memory network
	inmemCluster, err := cluster.NewInmemLayerCluster("inmem-cluster", 3, hclog.New(&hclog.LoggerOptions{
		Mutex: &sync.Mutex{},
		Level: hclog.Trace,
		Name:  "inmem-cluster",
	}))
	if err != nil {
		t.Fatal(err)
	}
	testCluster := vault.NewTestCluster(t, &vault.CoreConfig{
		Experiments: []string{experiments.VaultExperimentEventsAlpha1},
		AuditBackends: map[string]audit.Factory{
			"nop": corehelpers.NoopAuditFactory(nil),
		},
	}, &vault.TestClusterOptions{
		ClusterLayers: inmemCluster,
	})
	cores := testCluster.Cores
	testCluster.Start()
	defer testCluster.Cleanup()

	rootToken := testCluster.RootToken

	// Wait for core to become active
	vault.TestWaitActiveForwardingReady(t, cores[0].Core)

	// Test forwarding a request. Since we're going directly from core to core
	// with no fallback we know that if it worked, request handling is working
	c := cores[1]
	standby, err := c.Standby()
	if err != nil {
		t.Fatal(err)
	}
	if !standby {
		t.Fatal("expected core to be standby")
	}

	// We need to call Leader as that refreshes the connection info
	isLeader, _, _, err := c.Leader()
	if err != nil {
		t.Fatal(err)
	}
	if isLeader {
		t.Fatal("core should not be leader")
	}
	corehelpers.RetryUntil(t, 5*time.Second, func() error {
		state := c.ActiveNodeReplicationState()
		if state == 0 {
			return fmt.Errorf("heartbeats have not yet returned a valid active node replication state: %d", state)
		}
		return nil
	})

	req, err := http.NewRequest("GET", "https://pushit.real.good:9281/v1/sys/events/subscribe/xyz?json=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(namespace.RootContext(req.Context()))
	req.Header.Add(consts.AuthHeaderName, rootToken)

	resp := httptest.NewRecorder()
	forwardRequest(cores[1].Core, resp, req)

	header := resp.Header()
	if header == nil {
		t.Fatal("err: expected at least a Location header")
	}
	if !strings.HasPrefix(header.Get("Location"), "wss://") {
		t.Fatalf("bad location: %s", header.Get("Location"))
	}

	// test forwarding requests to each core
	handled := 0
	forwarded := 0
	for _, c := range cores {
		resp := httptest.NewRecorder()
		fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handled++
		})
		handleRequestForwarding(c.Core, fakeHandler).ServeHTTP(resp, req)
		header := resp.Header()
		if header == nil {
			continue
		}
		if strings.HasPrefix(header.Get("Location"), "wss://") {
			forwarded++
		}
	}
	if handled != 1 && forwarded != 2 {
		t.Fatalf("Expected 1 core to handle the request and 2 to forward")
	}
}
