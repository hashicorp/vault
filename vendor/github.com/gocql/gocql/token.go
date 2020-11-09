// Copyright (c) 2015 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocql

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/gocql/gocql/internal/murmur"
)

// a token partitioner
type partitioner interface {
	Name() string
	Hash([]byte) token
	ParseString(string) token
}

// a token
type token interface {
	fmt.Stringer
	Less(token) bool
}

// murmur3 partitioner and token
type murmur3Partitioner struct{}
type murmur3Token int64

func (p murmur3Partitioner) Name() string {
	return "Murmur3Partitioner"
}

func (p murmur3Partitioner) Hash(partitionKey []byte) token {
	h1 := murmur.Murmur3H1(partitionKey)
	return murmur3Token(h1)
}

// murmur3 little-endian, 128-bit hash, but returns only h1
func (p murmur3Partitioner) ParseString(str string) token {
	val, _ := strconv.ParseInt(str, 10, 64)
	return murmur3Token(val)
}

func (m murmur3Token) String() string {
	return strconv.FormatInt(int64(m), 10)
}

func (m murmur3Token) Less(token token) bool {
	return m < token.(murmur3Token)
}

// order preserving partitioner and token
type orderedPartitioner struct{}
type orderedToken string

func (p orderedPartitioner) Name() string {
	return "OrderedPartitioner"
}

func (p orderedPartitioner) Hash(partitionKey []byte) token {
	// the partition key is the token
	return orderedToken(partitionKey)
}

func (p orderedPartitioner) ParseString(str string) token {
	return orderedToken(str)
}

func (o orderedToken) String() string {
	return string(o)
}

func (o orderedToken) Less(token token) bool {
	return o < token.(orderedToken)
}

// random partitioner and token
type randomPartitioner struct{}
type randomToken big.Int

func (r randomPartitioner) Name() string {
	return "RandomPartitioner"
}

// 2 ** 128
var maxHashInt, _ = new(big.Int).SetString("340282366920938463463374607431768211456", 10)

func (p randomPartitioner) Hash(partitionKey []byte) token {
	sum := md5.Sum(partitionKey)
	val := new(big.Int)
	val.SetBytes(sum[:])
	if sum[0] > 127 {
		val.Sub(val, maxHashInt)
		val.Abs(val)
	}

	return (*randomToken)(val)
}

func (p randomPartitioner) ParseString(str string) token {
	val := new(big.Int)
	val.SetString(str, 10)
	return (*randomToken)(val)
}

func (r *randomToken) String() string {
	return (*big.Int)(r).String()
}

func (r *randomToken) Less(token token) bool {
	return -1 == (*big.Int)(r).Cmp((*big.Int)(token.(*randomToken)))
}

type hostToken struct {
	token token
	host  *HostInfo
}

func (ht hostToken) String() string {
	return fmt.Sprintf("{token=%v host=%v}", ht.token, ht.host.HostID())
}

// a data structure for organizing the relationship between tokens and hosts
type tokenRing struct {
	partitioner partitioner

	// tokens map token range to primary replica.
	// The elements in tokens are sorted by token ascending.
	// The range for a given item in tokens starts after preceding range and ends with the token specified in
	// token. The end token is part of the range.
	// The lowest (i.e. index 0) range wraps around the ring (its preceding range is the one with largest index).
	tokens []hostToken

	hosts []*HostInfo
}

func newTokenRing(partitioner string, hosts []*HostInfo) (*tokenRing, error) {
	tokenRing := &tokenRing{
		hosts: hosts,
	}

	if strings.HasSuffix(partitioner, "Murmur3Partitioner") {
		tokenRing.partitioner = murmur3Partitioner{}
	} else if strings.HasSuffix(partitioner, "OrderedPartitioner") {
		tokenRing.partitioner = orderedPartitioner{}
	} else if strings.HasSuffix(partitioner, "RandomPartitioner") {
		tokenRing.partitioner = randomPartitioner{}
	} else {
		return nil, fmt.Errorf("Unsupported partitioner '%s'", partitioner)
	}

	for _, host := range hosts {
		for _, strToken := range host.Tokens() {
			token := tokenRing.partitioner.ParseString(strToken)
			tokenRing.tokens = append(tokenRing.tokens, hostToken{token, host})
		}
	}

	sort.Sort(tokenRing)

	return tokenRing, nil
}

func (t *tokenRing) Len() int {
	return len(t.tokens)
}

func (t *tokenRing) Less(i, j int) bool {
	return t.tokens[i].token.Less(t.tokens[j].token)
}

func (t *tokenRing) Swap(i, j int) {
	t.tokens[i], t.tokens[j] = t.tokens[j], t.tokens[i]
}

func (t *tokenRing) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("TokenRing(")
	if t.partitioner != nil {
		buf.WriteString(t.partitioner.Name())
	}
	buf.WriteString("){")
	sep := ""
	for i, th := range t.tokens {
		buf.WriteString(sep)
		sep = ","
		buf.WriteString("\n\t[")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString("]")
		buf.WriteString(th.token.String())
		buf.WriteString(":")
		buf.WriteString(th.host.ConnectAddress().String())
	}
	buf.WriteString("\n}")
	return string(buf.Bytes())
}

func (t *tokenRing) GetHostForToken(token token) (host *HostInfo, endToken token) {
	if t == nil || len(t.tokens) == 0 {
		return nil, nil
	}

	// find the primary replica
	p := sort.Search(len(t.tokens), func(i int) bool {
		return !t.tokens[i].token.Less(token)
	})

	if p == len(t.tokens) {
		// wrap around to the first in the ring
		p = 0
	}

	v := t.tokens[p]
	return v.host, v.token
}
