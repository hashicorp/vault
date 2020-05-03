// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// MongoDBX509 is the mechanism name for MongoDBX509.
const MongoDBX509 = "MONGODB-X509"

func newMongoDBX509Authenticator(cred *Cred) (Authenticator, error) {
	return &MongoDBX509Authenticator{User: cred.Username}, nil
}

// MongoDBX509Authenticator uses X.509 certificates over TLS to authenticate a connection.
type MongoDBX509Authenticator struct {
	User string
}

// Auth implements the Authenticator interface.
func (a *MongoDBX509Authenticator) Auth(ctx context.Context, desc description.Server, conn driver.Connection) error {
	requestDoc := bsoncore.AppendInt32Element(nil, "authenticate", 1)
	requestDoc = bsoncore.AppendStringElement(requestDoc, "mechanism", MongoDBX509)

	if desc.WireVersion == nil || desc.WireVersion.Max < 5 {
		requestDoc = bsoncore.AppendStringElement(requestDoc, "user", a.User)
	}

	authCmd := operation.
		NewCommand(bsoncore.BuildDocument(nil, requestDoc)).
		Database("$external").
		Deployment(driver.SingleConnectionDeployment{conn})
	err := authCmd.Execute(ctx)
	if err != nil {
		return newAuthError("round trip error", err)
	}

	return nil
}
