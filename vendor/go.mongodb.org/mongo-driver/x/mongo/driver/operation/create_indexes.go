// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// CreateIndexes performs a createIndexes operation.
type CreateIndexes struct {
	authenticator driver.Authenticator
	commitQuorum  bsoncore.Value
	indexes       bsoncore.Document
	maxTime       *time.Duration
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	result        CreateIndexesResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}

// CreateIndexesResult represents a createIndexes result returned by the server.
type CreateIndexesResult struct {
	// If the collection was created automatically.
	CreatedCollectionAutomatically bool
	// The number of indexes existing after this command.
	IndexesAfter int32
	// The number of indexes existing before this command.
	IndexesBefore int32
}

func buildCreateIndexesResult(response bsoncore.Document) (CreateIndexesResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return CreateIndexesResult{}, err
	}
	cir := CreateIndexesResult{}
	for _, element := range elements {
		switch element.Key() {
		case "createdCollectionAutomatically":
			var ok bool
			cir.CreatedCollectionAutomatically, ok = element.Value().BooleanOK()
			if !ok {
				return cir, fmt.Errorf("response field 'createdCollectionAutomatically' is type bool, but received BSON type %s", element.Value().Type)
			}
		case "indexesAfter":
			var ok bool
			cir.IndexesAfter, ok = element.Value().AsInt32OK()
			if !ok {
				return cir, fmt.Errorf("response field 'indexesAfter' is type int32, but received BSON type %s", element.Value().Type)
			}
		case "indexesBefore":
			var ok bool
			cir.IndexesBefore, ok = element.Value().AsInt32OK()
			if !ok {
				return cir, fmt.Errorf("response field 'indexesBefore' is type int32, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return cir, nil
}

// NewCreateIndexes constructs and returns a new CreateIndexes.
func NewCreateIndexes(indexes bsoncore.Document) *CreateIndexes {
	return &CreateIndexes{
		indexes: indexes,
	}
}

// Result returns the result of executing this operation.
func (ci *CreateIndexes) Result() CreateIndexesResult { return ci.result }

func (ci *CreateIndexes) processResponse(info driver.ResponseInfo) error {
	var err error
	ci.result, err = buildCreateIndexesResult(info.ServerResponse)
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (ci *CreateIndexes) Execute(ctx context.Context) error {
	if ci.deployment == nil {
		return errors.New("the CreateIndexes operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         ci.command,
		ProcessResponseFn: ci.processResponse,
		Client:            ci.session,
		Clock:             ci.clock,
		CommandMonitor:    ci.monitor,
		Crypt:             ci.crypt,
		Database:          ci.database,
		Deployment:        ci.deployment,
		MaxTime:           ci.maxTime,
		Selector:          ci.selector,
		WriteConcern:      ci.writeConcern,
		ServerAPI:         ci.serverAPI,
		Timeout:           ci.timeout,
		Name:              driverutil.CreateIndexesOp,
		Authenticator:     ci.authenticator,
	}.Execute(ctx)

}

func (ci *CreateIndexes) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "createIndexes", ci.collection)
	if ci.commitQuorum.Type != bsontype.Type(0) {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(9) {
			return nil, errors.New("the 'commitQuorum' command parameter requires a minimum server wire version of 9")
		}
		dst = bsoncore.AppendValueElement(dst, "commitQuorum", ci.commitQuorum)
	}
	if ci.indexes != nil {
		dst = bsoncore.AppendArrayElement(dst, "indexes", ci.indexes)
	}
	return dst, nil
}

// CommitQuorum specifies the number of data-bearing members of a replica set, including the primary, that must
// complete the index builds successfully before the primary marks the indexes as ready. This should either be a
// string or int32 value.
func (ci *CreateIndexes) CommitQuorum(commitQuorum bsoncore.Value) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.commitQuorum = commitQuorum
	return ci
}

// Indexes specifies an array containing index specification documents for the indexes being created.
func (ci *CreateIndexes) Indexes(indexes bsoncore.Document) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.indexes = indexes
	return ci
}

// MaxTime specifies the maximum amount of time to allow the query to run on the server.
func (ci *CreateIndexes) MaxTime(maxTime *time.Duration) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.maxTime = maxTime
	return ci
}

// Session sets the session for this operation.
func (ci *CreateIndexes) Session(session *session.Client) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.session = session
	return ci
}

// ClusterClock sets the cluster clock for this operation.
func (ci *CreateIndexes) ClusterClock(clock *session.ClusterClock) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.clock = clock
	return ci
}

// Collection sets the collection that this command will run against.
func (ci *CreateIndexes) Collection(collection string) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.collection = collection
	return ci
}

// CommandMonitor sets the monitor to use for APM events.
func (ci *CreateIndexes) CommandMonitor(monitor *event.CommandMonitor) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.monitor = monitor
	return ci
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (ci *CreateIndexes) Crypt(crypt driver.Crypt) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.crypt = crypt
	return ci
}

// Database sets the database to run this operation against.
func (ci *CreateIndexes) Database(database string) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.database = database
	return ci
}

// Deployment sets the deployment to use for this operation.
func (ci *CreateIndexes) Deployment(deployment driver.Deployment) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.deployment = deployment
	return ci
}

// ServerSelector sets the selector used to retrieve a server.
func (ci *CreateIndexes) ServerSelector(selector description.ServerSelector) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.selector = selector
	return ci
}

// WriteConcern sets the write concern for this operation.
func (ci *CreateIndexes) WriteConcern(writeConcern *writeconcern.WriteConcern) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.writeConcern = writeConcern
	return ci
}

// ServerAPI sets the server API version for this operation.
func (ci *CreateIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.serverAPI = serverAPI
	return ci
}

// Timeout sets the timeout for this operation.
func (ci *CreateIndexes) Timeout(timeout *time.Duration) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.timeout = timeout
	return ci
}

// Authenticator sets the authenticator to use for this operation.
func (ci *CreateIndexes) Authenticator(authenticator driver.Authenticator) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.authenticator = authenticator
	return ci
}
