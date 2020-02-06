// Package etype provides the Kerberos Encryption Type interface
package etype

import "hash"

// EType is the interface defining the Encryption Type.
type EType interface {
	GetETypeID() int32
	GetHashID() int32
	GetKeyByteSize() int
	GetKeySeedBitLength() int                                   // key-generation seed length, k
	GetDefaultStringToKeyParams() string                        // default string-to-key parameters (s2kparams)
	StringToKey(string, salt, s2kparams string) ([]byte, error) // string-to-key (UTF-8 string, UTF-8 string, opaque)->(protocol-key)
	RandomToKey(b []byte) []byte                                // random-to-key (bitstring[K])->(protocol-key)
	GetHMACBitLength() int                                      // HMAC output size, h
	GetMessageBlockByteSize() int                               // message block size, m
	EncryptData(key, data []byte) ([]byte, []byte, error)       // E function - encrypt (specific-key, state, octet string)->(state, octet string)
	EncryptMessage(key, message []byte, usage uint32) ([]byte, []byte, error)
	DecryptData(key, data []byte) ([]byte, error) // D function
	DecryptMessage(key, ciphertext []byte, usage uint32) ([]byte, error)
	GetCypherBlockBitLength() int                           // cipher block size, c
	GetConfounderByteSize() int                             // This is the same as the cipher block size but in bytes.
	DeriveKey(protocolKey, usage []byte) ([]byte, error)    // DK key-derivation (protocol-key, integer)->(specific-key)
	DeriveRandom(protocolKey, usage []byte) ([]byte, error) // DR pseudo-random (protocol-key, octet-string)->(octet-string)
	VerifyIntegrity(protocolKey, ct, pt []byte, usage uint32) bool
	GetChecksumHash(protocolKey, data []byte, usage uint32) ([]byte, error)
	VerifyChecksum(protocolKey, data, chksum []byte, usage uint32) bool
	GetHashFunc() func() hash.Hash
}
