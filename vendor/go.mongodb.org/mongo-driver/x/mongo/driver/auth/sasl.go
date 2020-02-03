// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// SaslClient is the client piece of a sasl conversation.
type SaslClient interface {
	Start() (string, []byte, error)
	Next(challenge []byte) ([]byte, error)
	Completed() bool
}

// SaslClientCloser is a SaslClient that has resources to clean up.
type SaslClientCloser interface {
	SaslClient
	Close()
}

// ConductSaslConversation handles running a sasl conversation with MongoDB.
func ConductSaslConversation(ctx context.Context, conn driver.Connection, db string, client SaslClient) error {

	if db == "" {
		db = defaultAuthDB
	}

	if closer, ok := client.(SaslClientCloser); ok {
		defer closer.Close()
	}

	mech, payload, err := client.Start()
	if err != nil {
		return newError(err, mech)
	}

	doc := bsoncore.BuildDocumentFromElements(nil,
		bsoncore.AppendInt32Element(nil, "saslStart", 1),
		bsoncore.AppendStringElement(nil, "mechanism", mech),
		bsoncore.AppendBinaryElement(nil, "payload", 0x00, payload),
	)
	saslStartCmd := operation.NewCommand(doc).Database(db).Deployment(driver.SingleConnectionDeployment{conn})

	type saslResponse struct {
		ConversationID int    `bson:"conversationId"`
		Code           int    `bson:"code"`
		Done           bool   `bson:"done"`
		Payload        []byte `bson:"payload"`
	}

	var saslResp saslResponse

	err = saslStartCmd.Execute(ctx)
	if err != nil {
		return newError(err, mech)
	}
	rdr := saslStartCmd.Result()

	err = bson.Unmarshal(rdr, &saslResp)
	if err != nil {
		return newAuthError("unmarshall error", err)
	}

	cid := saslResp.ConversationID

	for {
		if saslResp.Code != 0 {
			return newError(err, mech)
		}

		if saslResp.Done && client.Completed() {
			return nil
		}

		payload, err = client.Next(saslResp.Payload)
		if err != nil {
			return newError(err, mech)
		}

		if saslResp.Done && client.Completed() {
			return nil
		}

		doc := bsoncore.BuildDocumentFromElements(nil,
			bsoncore.AppendInt32Element(nil, "saslContinue", 1),
			bsoncore.AppendInt32Element(nil, "conversationId", int32(cid)),
			bsoncore.AppendBinaryElement(nil, "payload", 0x00, payload),
		)
		saslContinueCmd := operation.NewCommand(doc).Database(db).Deployment(driver.SingleConnectionDeployment{conn})

		err = saslContinueCmd.Execute(ctx)
		if err != nil {
			return newError(err, mech)
		}
		rdr = saslContinueCmd.Result()

		err = bson.Unmarshal(rdr, &saslResp)
		if err != nil {
			return newAuthError("unmarshal error", err)
		}
	}
}
