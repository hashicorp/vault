// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package managed_key

import (
	"context"
	"crypto"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"

	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ssh"
)

type ManagedKeyInfo struct {
	publicKey ssh.PublicKey
	Name      NameKey
	Uuid      UUIDKey
	ctx       context.Context
	mkv       SSHManagedKeyView
}

type managedKeyId interface {
	String() string
}

type SSHManagedKeyView interface {
	BackendUUID() string
	GetManagedKeyView() (logical.ManagedKeySystemView, error)
	GetRandomReader() io.Reader
}

type (
	UUIDKey string
	NameKey string
)

type ecdsaSignature struct {
	R, S *big.Int
}

var HashFunctions = map[string]crypto.Hash{
	ssh.KeyAlgoECDSA256:  crypto.SHA256,
	ssh.KeyAlgoECDSA384:  crypto.SHA384,
	ssh.KeyAlgoECDSA521:  crypto.SHA512,
	ssh.KeyAlgoRSA:       crypto.SHA1,
	ssh.KeyAlgoRSASHA256: crypto.SHA256,
	ssh.KeyAlgoRSASHA512: crypto.SHA512,
}

func (u UUIDKey) String() string {
	return string(u)
}

func (n NameKey) String() string {
	return string(n)
}

func (m ManagedKeyInfo) PublicKey() ssh.PublicKey {
	return m.publicKey
}

func (m ManagedKeyInfo) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return m.SignWithAlgorithm(rand, data, "")
}

func (m ManagedKeyInfo) SignWithAlgorithm(rand io.Reader, input []byte, algorithm string) (*ssh.Signature, error) {
	keyType := m.publicKey.Type()
	if keyType != ssh.KeyAlgoRSA && (algorithm == ssh.KeyAlgoRSASHA256 || algorithm == ssh.KeyAlgoRSASHA512) {
		return nil, fmt.Errorf("invalid algorithm %s for key type %s", algorithm, keyType)
	}

	if algorithm == "default" || algorithm == "" {
		if keyType == ssh.KeyAlgoRSA {
			algorithm = ssh.KeyAlgoRSASHA256
		} else {
			algorithm = keyType
		}
	}

	opts, ok := HashFunctions[algorithm]
	if !ok {
		return nil, fmt.Errorf("invalid hash algorithm %s", algorithm)
	}

	hashFcn := opts.New()
	hashFcn.Write(input)
	digest := hashFcn.Sum(nil)

	sig, err := Sign(m.ctx, m.mkv, m.Uuid, rand, digest, opts)
	if err != nil {
		return nil, err
	}

	var sigBytes []byte
	switch keyType {
	case ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521:
		var ecdsaSig ecdsaSignature
		_, err = asn1.Unmarshal(sig, &ecdsaSig)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling signature: %w", err)
		}
		sigBytes = ssh.Marshal(ecdsaSig)
	case ssh.KeyAlgoRSA:
		sigBytes = sig
	}

	return &ssh.Signature{
		Format: algorithm,
		Blob:   sigBytes,
	}, nil
}
