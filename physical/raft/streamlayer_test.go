// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/vault/vault/cluster"
)

type mockClusterHook struct {
	address net.Addr
}

func (*mockClusterHook) AddClient(alpn string, client cluster.Client)       {}
func (*mockClusterHook) RemoveClient(alpn string)                           {}
func (*mockClusterHook) AddHandler(alpn string, handler cluster.Handler)    {}
func (*mockClusterHook) StopHandler(alpn string)                            {}
func (*mockClusterHook) TLSConfig(ctx context.Context) (*tls.Config, error) { return nil, nil }
func (m *mockClusterHook) Addr() net.Addr                                   { return m.address }
func (*mockClusterHook) GetDialerFunc(ctx context.Context, alpnProto string) func(string, time.Duration) (net.Conn, error) {
	return func(string, time.Duration) (net.Conn, error) {
		return nil, nil
	}
}

func TestStreamLayer_UnspecifiedIP(t *testing.T) {
	m := &mockClusterHook{
		address: &cluster.NetAddr{
			Host: "0.0.0.0:8200",
		},
	}

	raftTLSKey, err := GenerateTLSKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	raftTLS := &TLSKeyring{
		Keys:        []*TLSKey{raftTLSKey},
		ActiveKeyID: raftTLSKey.ID,
	}

	layer, err := NewRaftLayer(nil, raftTLS, m)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "cannot use unspecified IP with raft storage: 0.0.0.0:8200" {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if layer != nil {
		t.Fatal("expected nil layer")
	}

	m.address.(*cluster.NetAddr).Host = "10.0.0.1:8200"

	layer, err = NewRaftLayer(nil, raftTLS, m)
	if err != nil {
		t.Fatal(err)
	}

	if layer == nil {
		t.Fatal("nil layer")
	}
}
