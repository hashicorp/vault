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
	"unsafe"
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
	h1 := murmur3H1(partitionKey)
	return murmur3Token(int64(h1))
}

// murmur3 little-endian, 128-bit hash, but returns only h1
func murmur3H1(data []byte) uint64 {
	length := len(data)

	var h1, h2, k1, k2 uint64

	const (
		c1 = 0x87c37b91114253d5
		c2 = 0x4cf5ad432745937f
	)

	// body
	nBlocks := length / 16
	for i := 0; i < nBlocks; i++ {
		block := (*[2]uint64)(unsafe.Pointer(&data[i*16]))

		k1 = block[0]
		k2 = block[1]

		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33) // ROTL64(k1, 31)
		k1 *= c2
		h1 ^= k1

		h1 = (h1 << 27) | (h1 >> 37) // ROTL64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31) // ROTL64(k2, 33)
		k2 *= c1
		h2 ^= k2

		h2 = (h2 << 31) | (h2 >> 33) // ROTL64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	// tail
	tail := data[nBlocks*16:]
	k1 = 0
	k2 = 0
	switch length & 15 {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8])

		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31) // ROTL64(k2, 33)
		k2 *= c1
		h2 ^= k2

		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0])

		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33) // ROTL64(k1, 31)
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= uint64(length)
	h2 ^= uint64(length)

	h1 += h2
	h2 += h1

	// finalizer
	const (
		fmix1 = 0xff51afd7ed558ccd
		fmix2 = 0xc4ceb9fe1a85ec53
	)

	// fmix64(h1)
	h1 ^= h1 >> 33
	h1 *= fmix1
	h1 ^= h1 >> 33
	h1 *= fmix2
	h1 ^= h1 >> 33

	// fmix64(h2)
	h2 ^= h2 >> 33
	h2 *= fmix1
	h2 ^= h2 >> 33
	h2 *= fmix2
	h2 ^= h2 >> 33

	h1 += h2
	// the following is extraneous since h2 is discarded
	// h2 += h1

	return h1
}

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
type orderedToken []byte

func (p orderedPartitioner) Name() string {
	return "OrderedPartitioner"
}

func (p orderedPartitioner) Hash(partitionKey []byte) token {
	// the partition key is the token
	return orderedToken(partitionKey)
}

func (p orderedPartitioner) ParseString(str string) token {
	return orderedToken([]byte(str))
}

func (o orderedToken) String() string {
	return string([]byte(o))
}

func (o orderedToken) Less(token token) bool {
	return -1 == bytes.Compare(o, token.(orderedToken))
}

// random partitioner and token
type randomPartitioner struct{}
type randomToken big.Int

func (r randomPartitioner) Name() string {
	return "RandomPartitioner"
}

func (p randomPartitioner) Hash(partitionKey []byte) token {
	hash := md5.New()
	sum := hash.Sum(partitionKey)

	val := new(big.Int)
	val = val.SetBytes(sum)
	val = val.Abs(val)

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

// a data structure for organizing the relationship between tokens and hosts
type tokenRing struct {
	partitioner partitioner
	tokens      []token
	hosts       []*HostInfo
}

func newTokenRing(partitioner string, hosts []HostInfo) (*tokenRing, error) {
	tokenRing := &tokenRing{
		tokens: []token{},
		hosts:  []*HostInfo{},
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

	for i := range hosts {
		host := &hosts[i]
		for _, strToken := range host.Tokens {
			token := tokenRing.partitioner.ParseString(strToken)
			tokenRing.tokens = append(tokenRing.tokens, token)
			tokenRing.hosts = append(tokenRing.hosts, host)
		}
	}

	sort.Sort(tokenRing)

	return tokenRing, nil
}

func (t *tokenRing) Len() int {
	return len(t.tokens)
}

func (t *tokenRing) Less(i, j int) bool {
	return t.tokens[i].Less(t.tokens[j])
}

func (t *tokenRing) Swap(i, j int) {
	t.tokens[i], t.hosts[i], t.tokens[j], t.hosts[j] =
		t.tokens[j], t.hosts[j], t.tokens[i], t.hosts[i]
}

func (t *tokenRing) String() string {

	buf := &bytes.Buffer{}
	buf.WriteString("TokenRing(")
	if t.partitioner != nil {
		buf.WriteString(t.partitioner.Name())
	}
	buf.WriteString("){")
	sep := ""
	for i := range t.tokens {
		buf.WriteString(sep)
		sep = ","
		buf.WriteString("\n\t[")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString("]")
		buf.WriteString(t.tokens[i].String())
		buf.WriteString(":")
		buf.WriteString(t.hosts[i].Peer)
	}
	buf.WriteString("\n}")
	return string(buf.Bytes())
}

func (t *tokenRing) GetHostForPartitionKey(partitionKey []byte) *HostInfo {
	if t == nil {
		return nil
	}

	token := t.partitioner.Hash(partitionKey)
	return t.GetHostForToken(token)
}

func (t *tokenRing) GetHostForToken(token token) *HostInfo {
	if t == nil {
		return nil
	}

	// find the primary replica
	ringIndex := sort.Search(
		len(t.tokens),
		func(i int) bool {
			return !t.tokens[i].Less(token)
		},
	)
	if ringIndex == len(t.tokens) {
		// wrap around to the first in the ring
		ringIndex = 0
	}
	host := t.hosts[ringIndex]
	return host
}
