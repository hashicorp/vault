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

func BenchmarkRaftWithNetwork(b *testing.B) {
	b.StopTimer()
	raft1, _ := createRaftNetworkCluster(b, true, false)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
			Key:   "test",
			Value: []byte{byte(i)},
		})
		if err != nil {
			b.Fatal(err)
		}
		snapOut, err := os.CreateTemp(b.TempDir(), "bench")
		if err != nil {
			b.Fatal(err)
		}
		b.Cleanup(func() {
			_ = snapOut.Close()
			_ = os.Remove(snapOut.Name())
		})
		err = raft1.Snapshot(snapOut, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestRaftNetworkClusterWithMultipleTimeEncodingsSet tests that Raft nodes
// with different msgpack time.Time encodings set will still cluster together.
// However, with go-msgpack 2.1.0+, the decoder is tolerant of both encodings,
// so this could only fail if the decoder drastically changes in the future.
func TestRaftNetworkClusterWithMultipleTimeEncodingsSet(t *testing.T) {
	raft1, raft2 := createRaftNetworkCluster(t, true, false)
	for i := 0; i < 10; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
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
	if !bytes.Equal(entry.Value, []byte{9}) {
		t.Errorf("Expected {9} but got %+v", entry.Value)
	}
}

func createRaftNetworkCluster(tb testing.TB, overrideTimeFormat1, overrideTimeFormat2 bool) (*RaftBackend, *RaftBackend) {
	cipherSuites := []uint16{
		// 1.3
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
	}

	port1 := freeport.GetOne(tb)
	port2 := freeport.GetOne(tb)
	addr1, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port1))
	if err != nil {
		tb.Fatal(err)
	}
	addr2, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port2))
	if err != nil {
		tb.Fatal(err)
	}
	key1, err := GenerateTLSKey(rand.Reader)
	if err != nil {
		tb.Fatal(err)
	}
	key2, err := GenerateTLSKey(rand.Reader)
	if err != nil {
		tb.Fatal(err)
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
	tb.Cleanup(listener1.Stop)
	tb.Cleanup(listener2.Stop)

	raft1, dir1 := GetRaftWithOpts(tb, true, true, SetupOpts{
		TLSKeyring: &TLSKeyring{
			Keys:        []*TLSKey{key1, key2},
			ActiveKeyID: key1.ID,
		},
		ClusterListener:                 listener1,
		overrideMsgpackUseNewTimeFormat: &overrideTimeFormat1,
	})

	setupOpts2 := SetupOpts{
		TLSKeyring: &TLSKeyring{
			Keys:        []*TLSKey{key2, key1},
			ActiveKeyID: key2.ID,
		},
		ClusterListener:                 listener2,
		overrideMsgpackUseNewTimeFormat: &overrideTimeFormat2,
	}
	raft2, dir2 := GetRaftWithOpts(tb, false, true, setupOpts2)
	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)

	// Add raft2 to the cluster
	addNetworkPeer(tb, raft1, raft2, addr2, setupOpts2)

	return raft1, raft2
}

func addNetworkPeer(tb testing.TB, leader, follower *RaftBackend, followerAddr *net.TCPAddr, setupOpts SetupOpts) {
	tb.Helper()
	if err := leader.AddPeer(context.Background(), follower.NodeID(), followerAddr.String()); err != nil {
		tb.Fatal(err)
	}

	peers, err := leader.Peers(context.Background())
	if err != nil {
		tb.Fatal(err)
	}

	err = follower.Bootstrap(peers)
	if err != nil {
		tb.Fatal(err)
	}

	err = follower.SetupCluster(context.Background(), setupOpts)
	if err != nil {
		tb.Fatal(err)
	}
}
