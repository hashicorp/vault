// Copyright (c) 2015 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocql

import (
	"bytes"
	"math/big"
	"sort"
	"strconv"
	"testing"
)

// Tests of the murmur3Patitioner
func TestMurmur3Partitioner(t *testing.T) {
	token := murmur3Partitioner{}.ParseString("-1053604476080545076")

	if "-1053604476080545076" != token.String() {
		t.Errorf("Expected '-1053604476080545076' but was '%s'", token)
	}

	// at least verify that the partitioner
	// doesn't return nil
	pk, _ := marshalInt(nil, 1)
	token = murmur3Partitioner{}.Hash(pk)
	if token == nil {
		t.Fatal("token was nil")
	}
}

// Tests of the murmur3Token
func TestMurmur3Token(t *testing.T) {
	if murmur3Token(42).Less(murmur3Token(42)) {
		t.Errorf("Expected Less to return false, but was true")
	}
	if !murmur3Token(-42).Less(murmur3Token(42)) {
		t.Errorf("Expected Less to return true, but was false")
	}
	if murmur3Token(42).Less(murmur3Token(-42)) {
		t.Errorf("Expected Less to return false, but was true")
	}
}

// Tests of the orderedPartitioner
func TestOrderedPartitioner(t *testing.T) {
	// at least verify that the partitioner
	// doesn't return nil
	p := orderedPartitioner{}
	pk, _ := marshalInt(nil, 1)
	token := p.Hash(pk)
	if token == nil {
		t.Fatal("token was nil")
	}

	str := token.String()
	parsedToken := p.ParseString(str)

	if !bytes.Equal([]byte(token.(orderedToken)), []byte(parsedToken.(orderedToken))) {
		t.Errorf("Failed to convert to and from a string %s expected %x but was %x",
			str,
			[]byte(token.(orderedToken)),
			[]byte(parsedToken.(orderedToken)),
		)
	}
}

// Tests of the orderedToken
func TestOrderedToken(t *testing.T) {
	if orderedToken([]byte{0, 0, 4, 2}).Less(orderedToken([]byte{0, 0, 4, 2})) {
		t.Errorf("Expected Less to return false, but was true")
	}
	if !orderedToken([]byte{0, 0, 3}).Less(orderedToken([]byte{0, 0, 4, 2})) {
		t.Errorf("Expected Less to return true, but was false")
	}
	if orderedToken([]byte{0, 0, 4, 2}).Less(orderedToken([]byte{0, 0, 3})) {
		t.Errorf("Expected Less to return false, but was true")
	}
}

// Tests of the randomPartitioner
func TestRandomPartitioner(t *testing.T) {
	// at least verify that the partitioner
	// doesn't return nil
	p := randomPartitioner{}
	pk, _ := marshalInt(nil, 1)
	token := p.Hash(pk)
	if token == nil {
		t.Fatal("token was nil")
	}

	str := token.String()
	parsedToken := p.ParseString(str)

	if (*big.Int)(token.(*randomToken)).Cmp((*big.Int)(parsedToken.(*randomToken))) != 0 {
		t.Errorf("Failed to convert to and from a string %s expected %v but was %v",
			str,
			token,
			parsedToken,
		)
	}
}

// Tests of the randomToken
func TestRandomToken(t *testing.T) {
	if ((*randomToken)(big.NewInt(42))).Less((*randomToken)(big.NewInt(42))) {
		t.Errorf("Expected Less to return false, but was true")
	}
	if !((*randomToken)(big.NewInt(41))).Less((*randomToken)(big.NewInt(42))) {
		t.Errorf("Expected Less to return true, but was false")
	}
	if ((*randomToken)(big.NewInt(42))).Less((*randomToken)(big.NewInt(41))) {
		t.Errorf("Expected Less to return false, but was true")
	}
}

type intToken int

func (i intToken) String() string {
	return strconv.Itoa(int(i))
}

func (i intToken) Less(token token) bool {
	return i < token.(intToken)
}

// Test of the token ring implementation based on example at the start of this
// page of documentation:
// http://www.datastax.com/docs/0.8/cluster_architecture/partitioning
func TestIntTokenRing(t *testing.T) {
	host0 := &HostInfo{}
	host25 := &HostInfo{}
	host50 := &HostInfo{}
	host75 := &HostInfo{}
	ring := &tokenRing{
		partitioner: nil,
		// these tokens and hosts are out of order to test sorting
		tokens: []token{
			intToken(0),
			intToken(50),
			intToken(75),
			intToken(25),
		},
		hosts: []*HostInfo{
			host0,
			host50,
			host75,
			host25,
		},
	}

	sort.Sort(ring)

	if ring.GetHostForToken(intToken(0)) != host0 {
		t.Error("Expected host 0 for token 0")
	}
	if ring.GetHostForToken(intToken(1)) != host25 {
		t.Error("Expected host 25 for token 1")
	}
	if ring.GetHostForToken(intToken(24)) != host25 {
		t.Error("Expected host 25 for token 24")
	}
	if ring.GetHostForToken(intToken(25)) != host25 {
		t.Error("Expected host 25 for token 25")
	}
	if ring.GetHostForToken(intToken(26)) != host50 {
		t.Error("Expected host 50 for token 26")
	}
	if ring.GetHostForToken(intToken(49)) != host50 {
		t.Error("Expected host 50 for token 49")
	}
	if ring.GetHostForToken(intToken(50)) != host50 {
		t.Error("Expected host 50 for token 50")
	}
	if ring.GetHostForToken(intToken(51)) != host75 {
		t.Error("Expected host 75 for token 51")
	}
	if ring.GetHostForToken(intToken(74)) != host75 {
		t.Error("Expected host 75 for token 74")
	}
	if ring.GetHostForToken(intToken(75)) != host75 {
		t.Error("Expected host 75 for token 75")
	}
	if ring.GetHostForToken(intToken(76)) != host0 {
		t.Error("Expected host 0 for token 76")
	}
	if ring.GetHostForToken(intToken(99)) != host0 {
		t.Error("Expected host 0 for token 99")
	}
	if ring.GetHostForToken(intToken(100)) != host0 {
		t.Error("Expected host 0 for token 100")
	}
}

// Test for the behavior of a nil pointer to tokenRing
func TestNilTokenRing(t *testing.T) {
	var ring *tokenRing = nil

	if ring.GetHostForToken(nil) != nil {
		t.Error("Expected nil for nil token ring")
	}
	if ring.GetHostForPartitionKey(nil) != nil {
		t.Error("Expected nil for nil token ring")
	}
}

// Test of the recognition of the partitioner class
func TestUnknownTokenRing(t *testing.T) {
	_, err := newTokenRing("UnknownPartitioner", nil)
	if err == nil {
		t.Error("Expected error for unknown partitioner value, but was nil")
	}
}

// Test of the tokenRing with the Murmur3Partitioner
func TestMurmur3TokenRing(t *testing.T) {
	// Note, strings are parsed directly to int64, they are not murmur3 hashed
	var hosts []HostInfo = []HostInfo{
		HostInfo{
			Peer:   "0",
			Tokens: []string{"0"},
		},
		HostInfo{
			Peer:   "1",
			Tokens: []string{"25"},
		},
		HostInfo{
			Peer:   "2",
			Tokens: []string{"50"},
		},
		HostInfo{
			Peer:   "3",
			Tokens: []string{"75"},
		},
	}
	ring, err := newTokenRing("Murmur3Partitioner", hosts)
	if err != nil {
		t.Fatalf("Failed to create token ring due to error: %v", err)
	}

	p := murmur3Partitioner{}

	var actual *HostInfo
	actual = ring.GetHostForToken(p.ParseString("0"))
	if actual.Peer != "0" {
		t.Errorf("Expected peer 0 for token \"0\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("25"))
	if actual.Peer != "1" {
		t.Errorf("Expected peer 1 for token \"25\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("50"))
	if actual.Peer != "2" {
		t.Errorf("Expected peer 2 for token \"50\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("75"))
	if actual.Peer != "3" {
		t.Errorf("Expected peer 3 for token \"01\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("12"))
	if actual.Peer != "1" {
		t.Errorf("Expected peer 1 for token \"12\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("24324545443332"))
	if actual.Peer != "0" {
		t.Errorf("Expected peer 0 for token \"24324545443332\", but was %s", actual.Peer)
	}
}

// Test of the tokenRing with the OrderedPartitioner
func TestOrderedTokenRing(t *testing.T) {
	// Tokens here more or less are similar layout to the int tokens above due
	// to each numeric character translating to a consistently offset byte.
	var hosts []HostInfo = []HostInfo{
		HostInfo{
			Peer: "0",
			Tokens: []string{
				"00",
			},
		},
		HostInfo{
			Peer: "1",
			Tokens: []string{
				"25",
			},
		},
		HostInfo{
			Peer: "2",
			Tokens: []string{
				"50",
			},
		},
		HostInfo{
			Peer: "3",
			Tokens: []string{
				"75",
			},
		},
	}
	ring, err := newTokenRing("OrderedPartitioner", hosts)
	if err != nil {
		t.Fatalf("Failed to create token ring due to error: %v", err)
	}

	p := orderedPartitioner{}

	var actual *HostInfo
	actual = ring.GetHostForToken(p.ParseString("0"))
	if actual.Peer != "0" {
		t.Errorf("Expected peer 0 for token \"0\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("25"))
	if actual.Peer != "1" {
		t.Errorf("Expected peer 1 for token \"25\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("50"))
	if actual.Peer != "2" {
		t.Errorf("Expected peer 2 for token \"50\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("75"))
	if actual.Peer != "3" {
		t.Errorf("Expected peer 3 for token \"01\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("12"))
	if actual.Peer != "1" {
		t.Errorf("Expected peer 1 for token \"12\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("24324545443332"))
	if actual.Peer != "1" {
		t.Errorf("Expected peer 1 for token \"24324545443332\", but was %s", actual.Peer)
	}
}

// Test of the tokenRing with the RandomPartitioner
func TestRandomTokenRing(t *testing.T) {
	// String tokens are parsed into big.Int in base 10
	var hosts []HostInfo = []HostInfo{
		HostInfo{
			Peer: "0",
			Tokens: []string{
				"00",
			},
		},
		HostInfo{
			Peer: "1",
			Tokens: []string{
				"25",
			},
		},
		HostInfo{
			Peer: "2",
			Tokens: []string{
				"50",
			},
		},
		HostInfo{
			Peer: "3",
			Tokens: []string{
				"75",
			},
		},
	}
	ring, err := newTokenRing("RandomPartitioner", hosts)
	if err != nil {
		t.Fatalf("Failed to create token ring due to error: %v", err)
	}

	p := randomPartitioner{}

	var actual *HostInfo
	actual = ring.GetHostForToken(p.ParseString("0"))
	if actual.Peer != "0" {
		t.Errorf("Expected peer 0 for token \"0\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("25"))
	if actual.Peer != "1" {
		t.Errorf("Expected peer 1 for token \"25\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("50"))
	if actual.Peer != "2" {
		t.Errorf("Expected peer 2 for token \"50\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("75"))
	if actual.Peer != "3" {
		t.Errorf("Expected peer 3 for token \"01\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("12"))
	if actual.Peer != "1" {
		t.Errorf("Expected peer 1 for token \"12\", but was %s", actual.Peer)
	}

	actual = ring.GetHostForToken(p.ParseString("24324545443332"))
	if actual.Peer != "0" {
		t.Errorf("Expected peer 0 for token \"24324545443332\", but was %s", actual.Peer)
	}
}
