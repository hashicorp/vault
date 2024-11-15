// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// SpeculativeConversation represents an authentication conversation that can be merged with the initial connection
// handshake.
//
// FirstMessage method returns the first message to be sent to the server. This message will be included in the initial
// hello command.
//
// Finish takes the server response to the initial message and conducts the remainder of the conversation to
// authenticate the provided connection.
type SpeculativeConversation interface {
	FirstMessage() (bsoncore.Document, error)
	Finish(ctx context.Context, cfg *Config, firstResponse bsoncore.Document) error
}

// SpeculativeAuthenticator represents an authenticator that supports speculative authentication.
type SpeculativeAuthenticator interface {
	CreateSpeculativeConversation() (SpeculativeConversation, error)
}
