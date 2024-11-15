// Copyright (C) MongoDB, Inc. 2023-present.
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

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// UpdateSearchIndex performs a updateSearchIndex operation.
type UpdateSearchIndex struct {
	authenticator driver.Authenticator
	index         string
	definition    bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	result        UpdateSearchIndexResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}

// UpdateSearchIndexResult represents a single index in the updateSearchIndexResult result.
type UpdateSearchIndexResult struct {
	Ok int32
}

func buildUpdateSearchIndexResult(response bsoncore.Document) (UpdateSearchIndexResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return UpdateSearchIndexResult{}, err
	}
	usir := UpdateSearchIndexResult{}
	for _, element := range elements {
		if element.Key() == "ok" {
			var ok bool
			usir.Ok, ok = element.Value().AsInt32OK()
			if !ok {
				return usir, fmt.Errorf("response field 'ok' is type int32, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return usir, nil
}

// NewUpdateSearchIndex constructs and returns a new UpdateSearchIndex.
func NewUpdateSearchIndex(index string, definition bsoncore.Document) *UpdateSearchIndex {
	return &UpdateSearchIndex{
		index:      index,
		definition: definition,
	}
}

// Result returns the result of executing this operation.
func (usi *UpdateSearchIndex) Result() UpdateSearchIndexResult { return usi.result }

func (usi *UpdateSearchIndex) processResponse(info driver.ResponseInfo) error {
	var err error
	usi.result, err = buildUpdateSearchIndexResult(info.ServerResponse)
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (usi *UpdateSearchIndex) Execute(ctx context.Context) error {
	if usi.deployment == nil {
		return errors.New("the UpdateSearchIndex operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         usi.command,
		ProcessResponseFn: usi.processResponse,
		Client:            usi.session,
		Clock:             usi.clock,
		CommandMonitor:    usi.monitor,
		Crypt:             usi.crypt,
		Database:          usi.database,
		Deployment:        usi.deployment,
		Selector:          usi.selector,
		ServerAPI:         usi.serverAPI,
		Timeout:           usi.timeout,
		Authenticator:     usi.authenticator,
	}.Execute(ctx)

}

func (usi *UpdateSearchIndex) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "updateSearchIndex", usi.collection)
	dst = bsoncore.AppendStringElement(dst, "name", usi.index)
	dst = bsoncore.AppendDocumentElement(dst, "definition", usi.definition)
	return dst, nil
}

// Index specifies the index of the document being updated.
func (usi *UpdateSearchIndex) Index(name string) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.index = name
	return usi
}

// Definition specifies the definition for the document being created.
func (usi *UpdateSearchIndex) Definition(definition bsoncore.Document) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.definition = definition
	return usi
}

// Session sets the session for this operation.
func (usi *UpdateSearchIndex) Session(session *session.Client) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.session = session
	return usi
}

// ClusterClock sets the cluster clock for this operation.
func (usi *UpdateSearchIndex) ClusterClock(clock *session.ClusterClock) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.clock = clock
	return usi
}

// Collection sets the collection that this command will run against.
func (usi *UpdateSearchIndex) Collection(collection string) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.collection = collection
	return usi
}

// CommandMonitor sets the monitor to use for APM events.
func (usi *UpdateSearchIndex) CommandMonitor(monitor *event.CommandMonitor) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.monitor = monitor
	return usi
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (usi *UpdateSearchIndex) Crypt(crypt driver.Crypt) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.crypt = crypt
	return usi
}

// Database sets the database to run this operation against.
func (usi *UpdateSearchIndex) Database(database string) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.database = database
	return usi
}

// Deployment sets the deployment to use for this operation.
func (usi *UpdateSearchIndex) Deployment(deployment driver.Deployment) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.deployment = deployment
	return usi
}

// ServerSelector sets the selector used to retrieve a server.
func (usi *UpdateSearchIndex) ServerSelector(selector description.ServerSelector) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.selector = selector
	return usi
}

// ServerAPI sets the server API version for this operation.
func (usi *UpdateSearchIndex) ServerAPI(serverAPI *driver.ServerAPIOptions) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.serverAPI = serverAPI
	return usi
}

// Timeout sets the timeout for this operation.
func (usi *UpdateSearchIndex) Timeout(timeout *time.Duration) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.timeout = timeout
	return usi
}

// Authenticator sets the authenticator to use for this operation.
func (usi *UpdateSearchIndex) Authenticator(authenticator driver.Authenticator) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.authenticator = authenticator
	return usi
}
