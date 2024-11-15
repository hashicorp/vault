// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
	"go.mongodb.org/mongo-driver/internal/credproviders"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth/creds"
)

// MongoDBAWS is the mechanism name for MongoDBAWS.
const MongoDBAWS = "MONGODB-AWS"

func newMongoDBAWSAuthenticator(cred *Cred, httpClient *http.Client) (Authenticator, error) {
	if cred.Source != "" && cred.Source != sourceExternal {
		return nil, newAuthError("MONGODB-AWS source must be empty or $external", nil)
	}
	if httpClient == nil {
		return nil, errors.New("httpClient must not be nil")
	}
	return &MongoDBAWSAuthenticator{
		credentials: &credproviders.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     cred.Username,
				SecretAccessKey: cred.Password,
				SessionToken:    cred.Props["AWS_SESSION_TOKEN"],
			},
		},
		httpClient: httpClient,
	}, nil
}

// MongoDBAWSAuthenticator uses AWS-IAM credentials over SASL to authenticate a connection.
type MongoDBAWSAuthenticator struct {
	credentials *credproviders.StaticProvider
	httpClient  *http.Client
}

// Auth authenticates the connection.
func (a *MongoDBAWSAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	providers := creds.NewAWSCredentialProvider(a.httpClient, a.credentials)
	adapter := &awsSaslAdapter{
		conversation: &awsConversation{
			credentials: providers.Cred,
		},
	}
	err := ConductSaslConversation(ctx, cfg, sourceExternal, adapter)
	if err != nil {
		return newAuthError("sasl conversation error", err)
	}
	return nil
}

// Reauth reauthenticates the connection.
func (a *MongoDBAWSAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("AWS authentication does not support reauthentication", nil)
}

type awsSaslAdapter struct {
	conversation *awsConversation
}

var _ SaslClient = (*awsSaslAdapter)(nil)

func (a *awsSaslAdapter) Start() (string, []byte, error) {
	step, err := a.conversation.Step(nil)
	if err != nil {
		return MongoDBAWS, nil, err
	}
	return MongoDBAWS, step, nil
}

func (a *awsSaslAdapter) Next(_ context.Context, challenge []byte) ([]byte, error) {
	step, err := a.conversation.Step(challenge)
	if err != nil {
		return nil, err
	}
	return step, nil
}

func (a *awsSaslAdapter) Completed() bool {
	return a.conversation.Done()
}
