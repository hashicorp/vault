// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// MongoDBX509 is the mechanism name for MongoDBX509.
const MongoDBX509 = "MONGODB-X509"

func newMongoDBX509Authenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	// TODO(GODRIVER-3309): Validate that cred.Source is either empty or
	// "$external" to make validation uniform with other auth mechanisms that
	// require Source to be "$external" (e.g. MONGODB-AWS, MONGODB-OIDC, etc).
	return &MongoDBX509Authenticator{User: cred.Username}, nil
}

// MongoDBX509Authenticator uses X.509 certificates over TLS to authenticate a connection.
type MongoDBX509Authenticator struct {
	User string
}

var _ SpeculativeAuthenticator = (*MongoDBX509Authenticator)(nil)

// x509 represents a X509 authentication conversation. This type implements the SpeculativeConversation interface so the
// conversation can be executed in multi-step speculative fashion.
type x509Conversation struct{}

var _ SpeculativeConversation = (*x509Conversation)(nil)

// FirstMessage returns the first message to be sent to the server.
func (c *x509Conversation) FirstMessage() (bsoncore.Document, error) {
	return createFirstX509Message(), nil
}

// createFirstX509Message creates the first message for the X509 conversation.
func createFirstX509Message() bsoncore.Document {
	elements := [][]byte{
		bsoncore.AppendInt32Element(nil, "authenticate", 1),
		bsoncore.AppendStringElement(nil, "mechanism", MongoDBX509),
	}

	return bsoncore.BuildDocument(nil, elements...)
}

// Finish implements the SpeculativeConversation interface and is a no-op because an X509 conversation only has one
// step.
func (c *x509Conversation) Finish(context.Context, *Config, bsoncore.Document) error {
	return nil
}

// CreateSpeculativeConversation creates a speculative conversation for X509 authentication.
func (a *MongoDBX509Authenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	return &x509Conversation{}, nil
}

// Auth authenticates the provided connection by conducting an X509 authentication conversation.
func (a *MongoDBX509Authenticator) Auth(ctx context.Context, cfg *Config) error {
	requestDoc := createFirstX509Message()
	authCmd := operation.
		NewCommand(requestDoc).
		Database(sourceExternal).
		Deployment(driver.SingleConnectionDeployment{cfg.Connection}).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	err := authCmd.Execute(ctx)
	if err != nil {
		return newAuthError("round trip error", err)
	}

	return nil
}

// Reauth reauthenticates the connection.
func (a *MongoDBX509Authenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("X509 does not support reauthentication", nil)
}
