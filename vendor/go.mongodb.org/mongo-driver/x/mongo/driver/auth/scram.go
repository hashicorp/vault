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

	"github.com/xdg-go/scram"
	"github.com/xdg-go/stringprep"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	// SCRAMSHA1 holds the mechanism name "SCRAM-SHA-1"
	SCRAMSHA1 = "SCRAM-SHA-1"

	// SCRAMSHA256 holds the mechanism name "SCRAM-SHA-256"
	SCRAMSHA256 = "SCRAM-SHA-256"
)

var (
	// Additional options for the saslStart command to enable a shorter SCRAM conversation
	scramStartOptions bsoncore.Document = bsoncore.BuildDocumentFromElements(nil,
		bsoncore.AppendBooleanElement(nil, "skipEmptyExchange", true),
	)
)

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

var _ SpeculativeAuthenticator = (*ScramAuthenticator)(nil)

// Auth authenticates the provided connection by conducting a full SASL conversation.
func (a *ScramAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	err := ConductSaslConversation(ctx, cfg, a.source, a.createSaslClient())
	if err != nil {
		return newAuthError("sasl conversation error", err)
	}
	return nil
}

// CreateSpeculativeConversation creates a speculative conversation for SCRAM authentication.
func (a *ScramAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	return newSaslConversation(a.createSaslClient(), a.source, true), nil
}

func (a *ScramAuthenticator) createSaslClient() SaslClient {
	return &scramSaslAdapter{
		conversation: a.client.NewConversation(),
		mechanism:    a.mechanism,
	}
}

type scramSaslAdapter struct {
	mechanism    string
	conversation *scram.ClientConversation
}

var _ SaslClient = (*scramSaslAdapter)(nil)
var _ ExtraOptionsSaslClient = (*scramSaslAdapter)(nil)

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

func (*scramSaslAdapter) StartCommandOptions() bsoncore.Document {
	return scramStartOptions
}
