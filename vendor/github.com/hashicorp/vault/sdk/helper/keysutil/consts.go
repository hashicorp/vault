package keysutil

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

type HashType uint32

const (
	_                     = iota
	HashTypeSHA1 HashType = iota
	HashTypeSHA2224
	HashTypeSHA2256
	HashTypeSHA2384
	HashTypeSHA2512
)

type MarshalingType uint32

const (
	_                                 = iota
	MarshalingTypeASN1 MarshalingType = iota
	MarshalingTypeJWS
)

var (
	HashTypeMap = map[string]HashType{
		"sha1":     HashTypeSHA1,
		"sha2-224": HashTypeSHA2224,
		"sha2-256": HashTypeSHA2256,
		"sha2-384": HashTypeSHA2384,
		"sha2-512": HashTypeSHA2512,
	}

	HashFuncMap = map[HashType]func() hash.Hash{
		HashTypeSHA1:    sha1.New,
		HashTypeSHA2224: sha256.New224,
		HashTypeSHA2256: sha256.New,
		HashTypeSHA2384: sha512.New384,
		HashTypeSHA2512: sha512.New,
	}

	MarshalingTypeMap = map[string]MarshalingType{
		"asn1": MarshalingTypeASN1,
		"jws":  MarshalingTypeJWS,
	}
)
