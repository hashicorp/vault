// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package scram

import (
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"
)

type serverState int

const (
	serverFirst serverState = iota
	serverFinal
	serverDone
)

// ServerConversation implements the server-side of an authentication
// conversation with a client.  A new conversation must be created for
// each authentication attempt.
type ServerConversation struct {
	nonceGen     NonceGeneratorFcn
	hashGen      HashGeneratorFcn
	credentialCB CredentialLookup
	state        serverState
	credential   StoredCredentials
	valid        bool
	gs2Header    string
	username     string
	authzID      string
	nonce        string
	c1b          string
	s1           string
}

// Step takes a string provided from a client and attempts to move the
// authentication conversation forward.  It returns a string to be sent to the
// client or an error if the client message is invalid.  Calling Step after a
// conversation completes is also an error.
func (sc *ServerConversation) Step(challenge string) (response string, err error) {
	switch sc.state {
	case serverFirst:
		sc.state = serverFinal
		response, err = sc.firstMsg(challenge)
	case serverFinal:
		sc.state = serverDone
		response, err = sc.finalMsg(challenge)
	default:
		response, err = "", errors.New("Conversation already completed")
	}
	return
}

// Done returns true if the conversation is completed or has errored.
func (sc *ServerConversation) Done() bool {
	return sc.state == serverDone
}

// Valid returns true if the conversation successfully authenticated the
// client.
func (sc *ServerConversation) Valid() bool {
	return sc.valid
}

// Username returns the client-provided username.  This is valid to call
// if the first conversation Step() is successful.
func (sc *ServerConversation) Username() string {
	return sc.username
}

// AuthzID returns the (optional) client-provided authorization identity, if
// any.  If one was not provided, it returns the empty string.  This is valid
// to call if the first conversation Step() is successful.
func (sc *ServerConversation) AuthzID() string {
	return sc.authzID
}

func (sc *ServerConversation) firstMsg(c1 string) (string, error) {
	msg, err := parseClientFirst(c1)
	if err != nil {
		sc.state = serverDone
		return "", err
	}

	sc.gs2Header = msg.gs2Header
	sc.username = msg.username
	sc.authzID = msg.authzID

	sc.credential, err = sc.credentialCB(msg.username)
	if err != nil {
		sc.state = serverDone
		return "e=unknown-user", err
	}

	sc.nonce = msg.nonce + sc.nonceGen()
	sc.c1b = msg.c1b
	sc.s1 = fmt.Sprintf("r=%s,s=%s,i=%d",
		sc.nonce,
		base64.StdEncoding.EncodeToString([]byte(sc.credential.Salt)),
		sc.credential.Iters,
	)

	return sc.s1, nil
}

// For errors, returns server error message as well as non-nil error.  Callers
// can choose whether to send server error or not.
func (sc *ServerConversation) finalMsg(c2 string) (string, error) {
	msg, err := parseClientFinal(c2)
	if err != nil {
		return "", err
	}

	// Check channel binding matches what we expect; in this case, we expect
	// just the gs2 header we received as we don't support channel binding
	// with a data payload.  If we add binding, we need to independently
	// compute the header to match here.
	if string(msg.cbind) != sc.gs2Header {
		return "e=channel-bindings-dont-match", fmt.Errorf("channel binding received '%s' doesn't match expected '%s'", msg.cbind, sc.gs2Header)
	}

	// Check nonce received matches what we sent
	if msg.nonce != sc.nonce {
		return "e=other-error", errors.New("nonce received did not match nonce sent")
	}

	// Create auth message
	authMsg := sc.c1b + "," + sc.s1 + "," + msg.c2wop

	// Retrieve ClientKey from proof and verify it
	clientSignature := computeHMAC(sc.hashGen, sc.credential.StoredKey, []byte(authMsg))
	clientKey := xorBytes([]byte(msg.proof), clientSignature)
	storedKey := computeHash(sc.hashGen, clientKey)

	// Compare with constant-time function
	if !hmac.Equal(storedKey, sc.credential.StoredKey) {
		return "e=invalid-proof", errors.New("challenge proof invalid")
	}

	sc.valid = true

	// Compute and return server verifier
	serverSignature := computeHMAC(sc.hashGen, sc.credential.ServerKey, []byte(authMsg))
	return "v=" + base64.StdEncoding.EncodeToString(serverSignature), nil
}
