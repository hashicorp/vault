// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Copyright (C) MongoDB, Inc. 2018-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"fmt"

	"github.com/xdg/scram"
	"github.com/xdg/stringprep"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
)

// SCRAMSHA1 holds the mechanism name "SCRAM-SHA-1"
const SCRAMSHA1 = "SCRAM-SHA-1"

// SCRAMSHA256 holds the mechanism name "SCRAM-SHA-256"
const SCRAMSHA256 = "SCRAM-SHA-256"

func newScramSHA1Authenticator(cred *Cred) (Authenticator, error) {
	passdigest := mongoPasswordDigest(cred.Username, cred.Password)
	client, err := scram.SHA1.NewClientUnprepped(cred.Username, passdigest, "")
	if err != nil {
		return nil, newAuthError("error initializing SCRAM-SHA-1 client", err)
	}
	client.WithMinIterations(4096)
	return &ScramAuthenticator{
		mechanism: SCRAMSHA1,
		source:    cred.Source,
		client:    client,
	}, nil
}

func newScramSHA256Authenticator(cred *Cred) (Authenticator, error) {
	passprep, err := stringprep.SASLprep.Prepare(cred.Password)
	if err != nil {
		return nil, newAuthError(fmt.Sprintf("error SASLprepping password '%s'", cred.Password), err)
	}
	client, err := scram.SHA256.NewClientUnprepped(cred.Username, passprep, "")
	if err != nil {
		return nil, newAuthError("error initializing SCRAM-SHA-256 client", err)
	}
	client.WithMinIterations(4096)
	return &ScramAuthenticator{
		mechanism: SCRAMSHA256,
		source:    cred.Source,
		client:    client,
	}, nil
}

// ScramAuthenticator uses the SCRAM algorithm over SASL to authenticate a connection.
type ScramAuthenticator struct {
	mechanism string
	source    string
	client    *scram.Client
}

// Auth authenticates the connection.
func (a *ScramAuthenticator) Auth(ctx context.Context, _ description.Server, conn driver.Connection) error {
	adapter := &scramSaslAdapter{conversation: a.client.NewConversation(), mechanism: a.mechanism}
	err := ConductSaslConversation(ctx, conn, a.source, adapter)
	if err != nil {
		return newAuthError("sasl conversation error", err)
	}
	return nil
}

type scramSaslAdapter struct {
	mechanism    string
	conversation *scram.ClientConversation
}

func (a *scramSaslAdapter) Start() (string, []byte, error) {
	step, err := a.conversation.Step("")
	if err != nil {
		return a.mechanism, nil, err
	}
	return a.mechanism, []byte(step), nil
}

func (a *scramSaslAdapter) Next(challenge []byte) ([]byte, error) {
	step, err := a.conversation.Step(string(challenge))
	if err != nil {
		return nil, err
	}
	return []byte(step), nil
}

func (a *scramSaslAdapter) Completed() bool {
	return a.conversation.Done()
}
