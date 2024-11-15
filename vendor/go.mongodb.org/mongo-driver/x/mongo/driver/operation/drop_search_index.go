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

// DropSearchIndex performs an dropSearchIndex operation.
type DropSearchIndex struct {
	authenticator driver.Authenticator
	index         string
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	result        DropSearchIndexResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}

// DropSearchIndexResult represents a dropSearchIndex result returned by the server.
type DropSearchIndexResult struct {
	Ok int32
}

func buildDropSearchIndexResult(response bsoncore.Document) (DropSearchIndexResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DropSearchIndexResult{}, err
	}
	dsir := DropSearchIndexResult{}
	for _, element := range elements {
		if element.Key() == "ok" {
			var ok bool
			dsir.Ok, ok = element.Value().AsInt32OK()
			if !ok {
				return dsir, fmt.Errorf("response field 'ok' is type int32, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return dsir, nil
}

// NewDropSearchIndex constructs and returns a new DropSearchIndex.
func NewDropSearchIndex(index string) *DropSearchIndex {
	return &DropSearchIndex{
		index: index,
	}
}

// Result returns the result of executing this operation.
func (dsi *DropSearchIndex) Result() DropSearchIndexResult { return dsi.result }

func (dsi *DropSearchIndex) processResponse(info driver.ResponseInfo) error {
	var err error
	dsi.result, err = buildDropSearchIndexResult(info.ServerResponse)
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (dsi *DropSearchIndex) Execute(ctx context.Context) error {
	if dsi.deployment == nil {
		return errors.New("the DropSearchIndex operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         dsi.command,
		ProcessResponseFn: dsi.processResponse,
		Client:            dsi.session,
		Clock:             dsi.clock,
		CommandMonitor:    dsi.monitor,
		Crypt:             dsi.crypt,
		Database:          dsi.database,
		Deployment:        dsi.deployment,
		Selector:          dsi.selector,
		ServerAPI:         dsi.serverAPI,
		Timeout:           dsi.timeout,
		Authenticator:     dsi.authenticator,
	}.Execute(ctx)

}

func (dsi *DropSearchIndex) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "dropSearchIndex", dsi.collection)
	dst = bsoncore.AppendStringElement(dst, "name", dsi.index)
	return dst, nil
}

// Index specifies the name of the index to drop. If '*' is specified, all indexes will be dropped.
func (dsi *DropSearchIndex) Index(index string) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.index = index
	return dsi
}

// Session sets the session for this operation.
func (dsi *DropSearchIndex) Session(session *session.Client) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.session = session
	return dsi
}

// ClusterClock sets the cluster clock for this operation.
func (dsi *DropSearchIndex) ClusterClock(clock *session.ClusterClock) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.clock = clock
	return dsi
}

// Collection sets the collection that this command will run against.
func (dsi *DropSearchIndex) Collection(collection string) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.collection = collection
	return dsi
}

// CommandMonitor sets the monitor to use for APM events.
func (dsi *DropSearchIndex) CommandMonitor(monitor *event.CommandMonitor) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.monitor = monitor
	return dsi
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (dsi *DropSearchIndex) Crypt(crypt driver.Crypt) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.crypt = crypt
	return dsi
}

// Database sets the database to run this operation against.
func (dsi *DropSearchIndex) Database(database string) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.database = database
	return dsi
}

// Deployment sets the deployment to use for this operation.
func (dsi *DropSearchIndex) Deployment(deployment driver.Deployment) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.deployment = deployment
	return dsi
}

// ServerSelector sets the selector used to retrieve a server.
func (dsi *DropSearchIndex) ServerSelector(selector description.ServerSelector) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.selector = selector
	return dsi
}

// ServerAPI sets the server API version for this operation.
func (dsi *DropSearchIndex) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.serverAPI = serverAPI
	return dsi
}

// Timeout sets the timeout for this operation.
func (dsi *DropSearchIndex) Timeout(timeout *time.Duration) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.timeout = timeout
	return dsi
}

// Authenticator sets the authenticator to use for this operation.
func (dsi *DropSearchIndex) Authenticator(authenticator driver.Authenticator) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.authenticator = authenticator
	return dsi
}
