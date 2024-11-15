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

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// DropCollection performs a drop operation.
type DropCollection struct {
	authenticator driver.Authenticator
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	result        DropCollectionResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}

// DropCollectionResult represents a dropCollection result returned by the server.
type DropCollectionResult struct {
	// The number of indexes in the dropped collection.
	NIndexesWas int32
	// The namespace of the dropped collection.
	Ns string
}

func buildDropCollectionResult(response bsoncore.Document) (DropCollectionResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DropCollectionResult{}, err
	}
	dcr := DropCollectionResult{}
	for _, element := range elements {
		switch element.Key() {
		case "nIndexesWas":
			var ok bool
			dcr.NIndexesWas, ok = element.Value().AsInt32OK()
			if !ok {
				return dcr, fmt.Errorf("response field 'nIndexesWas' is type int32, but received BSON type %s", element.Value().Type)
			}
		case "ns":
			var ok bool
			dcr.Ns, ok = element.Value().StringValueOK()
			if !ok {
				return dcr, fmt.Errorf("response field 'ns' is type string, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return dcr, nil
}

// NewDropCollection constructs and returns a new DropCollection.
func NewDropCollection() *DropCollection {
	return &DropCollection{}
}

// Result returns the result of executing this operation.
func (dc *DropCollection) Result() DropCollectionResult { return dc.result }

func (dc *DropCollection) processResponse(info driver.ResponseInfo) error {
	var err error
	dc.result, err = buildDropCollectionResult(info.ServerResponse)
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (dc *DropCollection) Execute(ctx context.Context) error {
	if dc.deployment == nil {
		return errors.New("the DropCollection operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         dc.command,
		ProcessResponseFn: dc.processResponse,
		Client:            dc.session,
		Clock:             dc.clock,
		CommandMonitor:    dc.monitor,
		Crypt:             dc.crypt,
		Database:          dc.database,
		Deployment:        dc.deployment,
		Selector:          dc.selector,
		WriteConcern:      dc.writeConcern,
		ServerAPI:         dc.serverAPI,
		Timeout:           dc.timeout,
		Name:              driverutil.DropOp,
		Authenticator:     dc.authenticator,
	}.Execute(ctx)

}

func (dc *DropCollection) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "drop", dc.collection)
	return dst, nil
}

// Session sets the session for this operation.
func (dc *DropCollection) Session(session *session.Client) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.session = session
	return dc
}

// ClusterClock sets the cluster clock for this operation.
func (dc *DropCollection) ClusterClock(clock *session.ClusterClock) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.clock = clock
	return dc
}

// Collection sets the collection that this command will run against.
func (dc *DropCollection) Collection(collection string) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.collection = collection
	return dc
}

// CommandMonitor sets the monitor to use for APM events.
func (dc *DropCollection) CommandMonitor(monitor *event.CommandMonitor) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.monitor = monitor
	return dc
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (dc *DropCollection) Crypt(crypt driver.Crypt) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.crypt = crypt
	return dc
}

// Database sets the database to run this operation against.
func (dc *DropCollection) Database(database string) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.database = database
	return dc
}

// Deployment sets the deployment to use for this operation.
func (dc *DropCollection) Deployment(deployment driver.Deployment) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.deployment = deployment
	return dc
}

// ServerSelector sets the selector used to retrieve a server.
func (dc *DropCollection) ServerSelector(selector description.ServerSelector) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.selector = selector
	return dc
}

// WriteConcern sets the write concern for this operation.
func (dc *DropCollection) WriteConcern(writeConcern *writeconcern.WriteConcern) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.writeConcern = writeConcern
	return dc
}

// ServerAPI sets the server API version for this operation.
func (dc *DropCollection) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.serverAPI = serverAPI
	return dc
}

// Timeout sets the timeout for this operation.
func (dc *DropCollection) Timeout(timeout *time.Duration) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.timeout = timeout
	return dc
}

// Authenticator sets the authenticator to use for this operation.
func (dc *DropCollection) Authenticator(authenticator driver.Authenticator) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.authenticator = authenticator
	return dc
}
