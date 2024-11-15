// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	// Ignore gosec warning "Blocklisted import crypto/md5: weak cryptographic primitive". We need
	// to use MD5 here to implement the MONGODB-CR specification.
	/* #nosec G501 */
	"crypto/md5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// MONGODBCR is the mechanism name for MONGODB-CR.
//
// The MONGODB-CR authentication mechanism is deprecated in MongoDB 3.6 and removed in
// MongoDB 4.0.
const MONGODBCR = "MONGODB-CR"

func newMongoDBCRAuthenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	source := cred.Source
	if source == "" {
		source = "admin"
	}
	return &MongoDBCRAuthenticator{
		DB:       source,
		Username: cred.Username,
		Password: cred.Password,
	}, nil
}

// MongoDBCRAuthenticator uses the MONGODB-CR algorithm to authenticate a connection.
//
// The MONGODB-CR authentication mechanism is deprecated in MongoDB 3.6 and removed in
// MongoDB 4.0.
type MongoDBCRAuthenticator struct {
	DB       string
	Username string
	Password string
}

// Auth authenticates the connection.
//
// The MONGODB-CR authentication mechanism is deprecated in MongoDB 3.6 and removed in
// MongoDB 4.0.
func (a *MongoDBCRAuthenticator) Auth(ctx context.Context, cfg *Config) error {

	db := a.DB
	if db == "" {
		db = defaultAuthDB
	}

	doc := bsoncore.BuildDocumentFromElements(nil, bsoncore.AppendInt32Element(nil, "getnonce", 1))
	cmd := operation.NewCommand(doc).
		Database(db).
		Deployment(driver.SingleConnectionDeployment{cfg.Connection}).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	err := cmd.Execute(ctx)
	if err != nil {
		return newError(err, MONGODBCR)
	}
	rdr := cmd.Result()

	var getNonceResult struct {
		Nonce string `bson:"nonce"`
	}

	err = bson.Unmarshal(rdr, &getNonceResult)
	if err != nil {
		return newAuthError("unmarshal error", err)
	}

	doc = bsoncore.BuildDocumentFromElements(nil,
		bsoncore.AppendInt32Element(nil, "authenticate", 1),
		bsoncore.AppendStringElement(nil, "user", a.Username),
		bsoncore.AppendStringElement(nil, "nonce", getNonceResult.Nonce),
		bsoncore.AppendStringElement(nil, "key", a.createKey(getNonceResult.Nonce)),
	)
	cmd = operation.NewCommand(doc).
		Database(db).
		Deployment(driver.SingleConnectionDeployment{cfg.Connection}).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	err = cmd.Execute(ctx)
	if err != nil {
		return newError(err, MONGODBCR)
	}

	return nil
}

// Reauth reauthenticates the connection.
func (a *MongoDBCRAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("MONGODB-CR does not support reauthentication", nil)
}

func (a *MongoDBCRAuthenticator) createKey(nonce string) string {
	// Ignore gosec warning "Use of weak cryptographic primitive". We need to use MD5 here to
	// implement the MONGODB-CR specification.
	/* #nosec G401 */
	h := md5.New()

	_, _ = io.WriteString(h, nonce)
	_, _ = io.WriteString(h, a.Username)
	_, _ = io.WriteString(h, mongoPasswordDigest(a.Username, a.Password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
