// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"net/url"
)

// Keyring is used to access the Variables keyring.
type Keyring struct {
	client *Client
}

// Keyring returns a handle to the Keyring endpoint
func (c *Client) Keyring() *Keyring {
	return &Keyring{client: c}
}

// EncryptionAlgorithm chooses which algorithm is used for
// encrypting / decrypting entries with this key
type EncryptionAlgorithm string

const (
	EncryptionAlgorithmAES256GCM EncryptionAlgorithm = "aes256-gcm"
)

// RootKeyMeta is the metadata used to refer to a RootKey.
type RootKeyMeta struct {
	KeyID       string // UUID
	Algorithm   EncryptionAlgorithm
	CreateTime  int64
	CreateIndex uint64
	ModifyIndex uint64
	State       RootKeyState
	PublishTime int64
}

// RootKeyState enum describes the lifecycle of a root key.
type RootKeyState string

const (
	RootKeyStateInactive     RootKeyState = "inactive"
	RootKeyStateActive                    = "active"
	RootKeyStateRekeying                  = "rekeying"
	RootKeyStateDeprecated                = "deprecated"
	RootKeyStatePrepublished              = "prepublished"
)

// List lists all the keyring metadata
func (k *Keyring) List(q *QueryOptions) ([]*RootKeyMeta, *QueryMeta, error) {
	var resp []*RootKeyMeta
	qm, err := k.client.query("/v1/operator/keyring/keys", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Delete deletes a specific inactive key from the keyring
func (k *Keyring) Delete(opts *KeyringDeleteOptions, w *WriteOptions) (*WriteMeta, error) {
	wm, err := k.client.delete(fmt.Sprintf("/v1/operator/keyring/key/%v",
		url.PathEscape(opts.KeyID)), nil, nil, w)
	return wm, err
}

// KeyringDeleteOptions are parameters for the Delete API
type KeyringDeleteOptions struct {
	KeyID string // UUID
}

// Rotate requests a key rotation
func (k *Keyring) Rotate(opts *KeyringRotateOptions, w *WriteOptions) (*RootKeyMeta, *WriteMeta, error) {
	qp := url.Values{}
	if opts != nil {
		if opts.Algorithm != "" {
			qp.Set("algo", string(opts.Algorithm))
		}
		if opts.Full {
			qp.Set("full", "true")
		}
		if opts.PublishTime > 0 {
			qp.Set("publish_time", fmt.Sprintf("%d", opts.PublishTime))
		}
	}
	resp := &struct{ Key *RootKeyMeta }{}
	wm, err := k.client.put("/v1/operator/keyring/rotate?"+qp.Encode(), nil, resp, w)
	return resp.Key, wm, err
}

// KeyringRotateOptions are parameters for the Rotate API
type KeyringRotateOptions struct {
	Full        bool
	Algorithm   EncryptionAlgorithm
	PublishTime int64
}
