// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package scram

import (
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

// Client implements the client side of SCRAM authentication.  It holds
// configuration values needed to initialize new client-side conversations for
// a specific username, password and authorization ID tuple.  Client caches
// the computationally-expensive parts of a SCRAM conversation as described in
// RFC-5802.  If repeated authentication conversations may be required for a
// user (e.g. disconnect/reconnect), the user's Client should be preserved.
//
// For security reasons, Clients have a default minimum PBKDF2 iteration count
// of 4096.  If a server requests a smaller iteration count, an authentication
// conversation will error.
//
// A Client can also be used by a server application to construct the hashed
// authentication values to be stored for a new user.  See StoredCredentials()
// for more.
type Client struct {
	sync.RWMutex
	username string
	password string
	authzID  string
	minIters int
	nonceGen NonceGeneratorFcn
	hashGen  HashGeneratorFcn
	cache    map[KeyFactors]derivedKeys
}

func newClient(username, password, authzID string, fcn HashGeneratorFcn) *Client {
	return &Client{
		username: username,
		password: password,
		authzID:  authzID,
		minIters: 4096,
		nonceGen: defaultNonceGenerator,
		hashGen:  fcn,
		cache:    make(map[KeyFactors]derivedKeys),
	}
}

// WithMinIterations changes minimum required PBKDF2 iteration count.
func (c *Client) WithMinIterations(n int) *Client {
	c.Lock()
	defer c.Unlock()
	c.minIters = n
	return c
}

// WithNonceGenerator replaces the default nonce generator (base64 encoding of
// 24 bytes from crypto/rand) with a custom generator.  This is provided for
// testing or for users with custom nonce requirements.
func (c *Client) WithNonceGenerator(ng NonceGeneratorFcn) *Client {
	c.Lock()
	defer c.Unlock()
	c.nonceGen = ng
	return c
}

// NewConversation constructs a client-side authentication conversation.
// Conversations cannot be reused, so this must be called for each new
// authentication attempt.
func (c *Client) NewConversation() *ClientConversation {
	c.RLock()
	defer c.RUnlock()
	return &ClientConversation{
		client:   c,
		nonceGen: c.nonceGen,
		hashGen:  c.hashGen,
		minIters: c.minIters,
	}
}

func (c *Client) getDerivedKeys(kf KeyFactors) derivedKeys {
	dk, ok := c.getCache(kf)
	if !ok {
		dk = c.computeKeys(kf)
		c.setCache(kf, dk)
	}
	return dk
}

// GetStoredCredentials takes a salt and iteration count structure and
// provides the values that must be stored by a server to authentication a
// user.  These values are what the Server credential lookup function must
// return for a given username.
func (c *Client) GetStoredCredentials(kf KeyFactors) StoredCredentials {
	dk := c.getDerivedKeys(kf)
	return StoredCredentials{
		KeyFactors: kf,
		StoredKey:  dk.StoredKey,
		ServerKey:  dk.ServerKey,
	}
}

func (c *Client) computeKeys(kf KeyFactors) derivedKeys {
	h := c.hashGen()
	saltedPassword := pbkdf2.Key([]byte(c.password), []byte(kf.Salt), kf.Iters, h.Size(), c.hashGen)
	clientKey := computeHMAC(c.hashGen, saltedPassword, []byte("Client Key"))

	return derivedKeys{
		ClientKey: clientKey,
		StoredKey: computeHash(c.hashGen, clientKey),
		ServerKey: computeHMAC(c.hashGen, saltedPassword, []byte("Server Key")),
	}
}

func (c *Client) getCache(kf KeyFactors) (derivedKeys, bool) {
	c.RLock()
	defer c.RUnlock()
	dk, ok := c.cache[kf]
	return dk, ok
}

func (c *Client) setCache(kf KeyFactors, dk derivedKeys) {
	c.Lock()
	defer c.Unlock()
	c.cache[kf] = dk
	return
}
