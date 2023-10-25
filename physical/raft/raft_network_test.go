// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !race

// The raft networking layer tends to reset the TLS keyring, which triggers
// the race detector even though it should be a no-op.

package raft

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/cluster"
)

// TestRaftNetworkClusterWithMultipleTimeEncodingsSet tests that Raft nodes
// with different msgpack time.Time encodings set will still cluster together.
// However, with go-msgpack 2.1.0+, the decoder is tolerant of both encodings,
// so this could only fail if the decoder drastically changes in the future.
func TestRaftNetworkClusterWithMultipleTimeEncodingsSet(t *testing.T) {
	// Create raft node
	cipherSuites := []uint16{
		// 1.3
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
	}

	port1 := freeport.GetOne(t)
	port2 := freeport.GetOne(t)
	addr1, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port1))
	if err != nil {
		t.Fatal(err)
	}
	addr2, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port2))
	if err != nil {
		t.Fatal(err)
	}
	key1, err := GenerateTLSKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	key2, err := GenerateTLSKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	logger1 := hclog.New(&hclog.LoggerOptions{
		Name: "raft1",
	})
	logger2 := hclog.New(&hclog.LoggerOptions{
		Name: "raft2",
	})
	listener1 := cluster.NewListener(
		cluster.NewTCPLayer([]*net.TCPAddr{addr1}, logger1), cipherSuites, logger1, time.Minute)
	listener2 := cluster.NewListener(
		cluster.NewTCPLayer([]*net.TCPAddr{addr2}, logger2), cipherSuites, logger2, time.Minute)
	go listener1.Run(context.Background())
	go listener2.Run(context.Background())
	t.Cleanup(listener1.Stop)
	t.Cleanup(listener2.Stop)

	raft1, dir1 := GetRaftWithOpts(t, true, true, SetupOpts{
		TLSKeyring: &TLSKeyring{
			Keys:        []*TLSKey{key1, key2},
			ActiveKeyID: key1.ID,
		},
		ClusterListener: listener1,
	})

	overrideTimeFormatFalse := false
	setupOpts2 := SetupOpts{
		TLSKeyring: &TLSKeyring{
			Keys:        []*TLSKey{key2, key1},
			ActiveKeyID: key2.ID,
		},
		ClusterListener:                 listener2,
		overrideMsgpackUseNewTimeFormat: &overrideTimeFormatFalse,
	}
	raft2, dir2 := GetRaftWithOpts(t, false, true, setupOpts2)
	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)

	// Add raft2 to the cluster
	addNetworkPeer(t, raft1, raft2, addr2, setupOpts2)

	for i := 0; i < 100; i++ {
		err = raft1.Put(context.Background(), &physical.Entry{
			Key:   "test",
			Value: []byte{byte(i)},
		})
		if err != nil {
			t.Error(err)
		}
	}
	for raft2.AppliedIndex() != raft1.AppliedIndex() {
		time.Sleep(1 * time.Millisecond)
	}
	entry, err := raft2.Get(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("Entry from raft secondary is nil")
	}
	if !bytes.Equal(entry.Value, []byte{99}) {
		t.Errorf("Expected {99} but got %+v", entry.Value)
	}
}

func addNetworkPeer(t *testing.T, leader, follower *RaftBackend, followerAddr *net.TCPAddr, setupOpts SetupOpts) {
	t.Helper()
	if err := leader.AddPeer(context.Background(), follower.NodeID(), followerAddr.String()); err != nil {
		t.Fatal(err)
	}

	peers, err := leader.Peers(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	err = follower.Bootstrap(peers)
	if err != nil {
		t.Fatal(err)
	}

	err = follower.SetupCluster(context.Background(), setupOpts)
	if err != nil {
		t.Fatal(err)
	}
}
