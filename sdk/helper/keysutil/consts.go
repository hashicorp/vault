// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keysutil

import (
	"crypto"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"golang.org/x/crypto/sha3"
)

type HashType uint32

const (
	HashTypeNone HashType = iota
	HashTypeSHA1
	HashTypeSHA2224
	HashTypeSHA2256
	HashTypeSHA2384
	HashTypeSHA2512
	HashTypeSHA3224
	HashTypeSHA3256
	HashTypeSHA3384
	HashTypeSHA3512
)

//go:generate enumer -type=MarshalingType -trimprefix=MarshalingType -transform=snake
type MarshalingType uint32

const (
	_ MarshalingType = iota
	MarshalingTypeASN1
	MarshalingTypeJWS
)

var (
	HashTypeMap = map[string]HashType{
		"none":     HashTypeNone,
		"sha1":     HashTypeSHA1,
		"sha2-224": HashTypeSHA2224,
		"sha2-256": HashTypeSHA2256,
		"sha2-384": HashTypeSHA2384,
		"sha2-512": HashTypeSHA2512,
		"sha3-224": HashTypeSHA3224,
		"sha3-256": HashTypeSHA3256,
		"sha3-384": HashTypeSHA3384,
		"sha3-512": HashTypeSHA3512,
	}

	HashFuncMap = map[HashType]func() hash.Hash{
		HashTypeNone:    nil,
		HashTypeSHA1:    sha1.New,
		HashTypeSHA2224: sha256.New224,
		HashTypeSHA2256: sha256.New,
		HashTypeSHA2384: sha512.New384,
		HashTypeSHA2512: sha512.New,
		HashTypeSHA3224: sha3.New224,
		HashTypeSHA3256: sha3.New256,
		HashTypeSHA3384: sha3.New384,
		HashTypeSHA3512: sha3.New512,
	}

	CryptoHashMap = map[HashType]crypto.Hash{
		HashTypeNone:    0,
		HashTypeSHA1:    crypto.SHA1,
		HashTypeSHA2224: crypto.SHA224,
		HashTypeSHA2256: crypto.SHA256,
		HashTypeSHA2384: crypto.SHA384,
		HashTypeSHA2512: crypto.SHA512,
		HashTypeSHA3224: crypto.SHA3_224,
		HashTypeSHA3256: crypto.SHA3_256,
		HashTypeSHA3384: crypto.SHA3_384,
		HashTypeSHA3512: crypto.SHA3_512,
	}

	MarshalingTypeMap = _MarshalingTypeNameToValueMap
)
