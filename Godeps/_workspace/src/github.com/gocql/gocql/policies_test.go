// Copyright (c) 2015 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocql

import (
	"fmt"
	"testing"

	"github.com/hailocab/go-hostpool"
)

// Tests of the round-robin host selection policy implementation
func TestRoundRobinHostPolicy(t *testing.T) {
	policy := RoundRobinHostPolicy()

	hosts := []HostInfo{
		HostInfo{HostId: "0"},
		HostInfo{HostId: "1"},
	}

	policy.SetHosts(hosts)

	// the first host selected is actually at [1], but this is ok for RR
	// interleaved iteration should always increment the host
	iterA := policy.Pick(nil)
	if actual := iterA(); actual.Info() != &hosts[1] {
		t.Errorf("Expected hosts[1] but was hosts[%s]", actual.Info().HostId)
	}
	iterB := policy.Pick(nil)
	if actual := iterB(); actual.Info() != &hosts[0] {
		t.Errorf("Expected hosts[0] but was hosts[%s]", actual.Info().HostId)
	}
	if actual := iterB(); actual.Info() != &hosts[1] {
		t.Errorf("Expected hosts[1] but was hosts[%s]", actual.Info().HostId)
	}
	if actual := iterA(); actual.Info() != &hosts[0] {
		t.Errorf("Expected hosts[0] but was hosts[%s]", actual.Info().HostId)
	}

	iterC := policy.Pick(nil)
	if actual := iterC(); actual.Info() != &hosts[1] {
		t.Errorf("Expected hosts[1] but was hosts[%s]", actual.Info().HostId)
	}
	if actual := iterC(); actual.Info() != &hosts[0] {
		t.Errorf("Expected hosts[0] but was hosts[%s]", actual.Info().HostId)
	}
}

// Tests of the token-aware host selection policy implementation with a
// round-robin host selection policy fallback.
func TestTokenAwareHostPolicy(t *testing.T) {
	policy := TokenAwareHostPolicy(RoundRobinHostPolicy())

	query := &Query{}

	iter := policy.Pick(nil)
	if iter == nil {
		t.Fatal("host iterator was nil")
	}
	actual := iter()
	if actual != nil {
		t.Fatalf("expected nil from iterator, but was %v", actual)
	}

	// set the hosts
	hosts := []HostInfo{
		HostInfo{Peer: "0", Tokens: []string{"00"}},
		HostInfo{Peer: "1", Tokens: []string{"25"}},
		HostInfo{Peer: "2", Tokens: []string{"50"}},
		HostInfo{Peer: "3", Tokens: []string{"75"}},
	}
	policy.SetHosts(hosts)

	// the token ring is not setup without the partitioner, but the fallback
	// should work
	if actual := policy.Pick(nil)(); actual.Info().Peer != "1" {
		t.Errorf("Expected peer 1 but was %s", actual.Info().Peer)
	}

	query.RoutingKey([]byte("30"))
	if actual := policy.Pick(query)(); actual.Info().Peer != "2" {
		t.Errorf("Expected peer 2 but was %s", actual.Info().Peer)
	}

	policy.SetPartitioner("OrderedPartitioner")

	// now the token ring is configured
	query.RoutingKey([]byte("20"))
	iter = policy.Pick(query)
	if actual := iter(); actual.Info().Peer != "1" {
		t.Errorf("Expected peer 1 but was %s", actual.Info().Peer)
	}
	// rest are round robin
	if actual := iter(); actual.Info().Peer != "3" {
		t.Errorf("Expected peer 3 but was %s", actual.Info().Peer)
	}
	if actual := iter(); actual.Info().Peer != "0" {
		t.Errorf("Expected peer 0 but was %s", actual.Info().Peer)
	}
	if actual := iter(); actual.Info().Peer != "2" {
		t.Errorf("Expected peer 2 but was %s", actual.Info().Peer)
	}
}

// Tests of the host pool host selection policy implementation
func TestHostPoolHostPolicy(t *testing.T) {
	policy := HostPoolHostPolicy(hostpool.New(nil))

	hosts := []HostInfo{
		HostInfo{HostId: "0", Peer: "0"},
		HostInfo{HostId: "1", Peer: "1"},
	}

	policy.SetHosts(hosts)

	// the first host selected is actually at [1], but this is ok for RR
	// interleaved iteration should always increment the host
	iter := policy.Pick(nil)
	actualA := iter()
	if actualA.Info().HostId != "0" {
		t.Errorf("Expected hosts[0] but was hosts[%s]", actualA.Info().HostId)
	}
	actualA.Mark(nil)

	actualB := iter()
	if actualB.Info().HostId != "1" {
		t.Errorf("Expected hosts[1] but was hosts[%s]", actualB.Info().HostId)
	}
	actualB.Mark(fmt.Errorf("error"))

	actualC := iter()
	if actualC.Info().HostId != "0" {
		t.Errorf("Expected hosts[0] but was hosts[%s]", actualC.Info().HostId)
	}
	actualC.Mark(nil)

	actualD := iter()
	if actualD.Info().HostId != "0" {
		t.Errorf("Expected hosts[0] but was hosts[%s]", actualD.Info().HostId)
	}
	actualD.Mark(nil)
}

// Tests of the round-robin connection selection policy implementation
func TestRoundRobinConnPolicy(t *testing.T) {
	policy := RoundRobinConnPolicy()()

	conn0 := &Conn{}
	conn1 := &Conn{}
	conn := []*Conn{
		conn0,
		conn1,
	}

	policy.SetConns(conn)

	// the first conn selected is actually at [1], but this is ok for RR
	if actual := policy.Pick(nil); actual != conn1 {
		t.Error("Expected conn1")
	}
	if actual := policy.Pick(nil); actual != conn0 {
		t.Error("Expected conn0")
	}
	if actual := policy.Pick(nil); actual != conn1 {
		t.Error("Expected conn1")
	}
}
