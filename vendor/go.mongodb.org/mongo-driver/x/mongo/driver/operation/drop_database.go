// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// DropDatabase performs a dropDatabase operation
type DropDatabase struct {
	authenticator driver.Authenticator
	session       *session.Client
	clock         *session.ClusterClock
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	serverAPI     *driver.ServerAPIOptions
}

// NewDropDatabase constructs and returns a new DropDatabase.
func NewDropDatabase() *DropDatabase {
	return &DropDatabase{}
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (dd *DropDatabase) Execute(ctx context.Context) error {
	if dd.deployment == nil {
		return errors.New("the DropDatabase operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:      dd.command,
		Client:         dd.session,
		Clock:          dd.clock,
		CommandMonitor: dd.monitor,
		Crypt:          dd.crypt,
		Database:       dd.database,
		Deployment:     dd.deployment,
		Selector:       dd.selector,
		WriteConcern:   dd.writeConcern,
		ServerAPI:      dd.serverAPI,
		Name:           driverutil.DropDatabaseOp,
		Authenticator:  dd.authenticator,
	}.Execute(ctx)

}

func (dd *DropDatabase) command(dst []byte, _ description.SelectedServer) ([]byte, error) {

	dst = bsoncore.AppendInt32Element(dst, "dropDatabase", 1)
	return dst, nil
}

// Session sets the session for this operation.
func (dd *DropDatabase) Session(session *session.Client) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.session = session
	return dd
}

// ClusterClock sets the cluster clock for this operation.
func (dd *DropDatabase) ClusterClock(clock *session.ClusterClock) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.clock = clock
	return dd
}

// CommandMonitor sets the monitor to use for APM events.
func (dd *DropDatabase) CommandMonitor(monitor *event.CommandMonitor) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.monitor = monitor
	return dd
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (dd *DropDatabase) Crypt(crypt driver.Crypt) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.crypt = crypt
	return dd
}

// Database sets the database to run this operation against.
func (dd *DropDatabase) Database(database string) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.database = database
	return dd
}

// Deployment sets the deployment to use for this operation.
func (dd *DropDatabase) Deployment(deployment driver.Deployment) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.deployment = deployment
	return dd
}

// ServerSelector sets the selector used to retrieve a server.
func (dd *DropDatabase) ServerSelector(selector description.ServerSelector) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.selector = selector
	return dd
}

// WriteConcern sets the write concern for this operation.
func (dd *DropDatabase) WriteConcern(writeConcern *writeconcern.WriteConcern) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.writeConcern = writeConcern
	return dd
}

// ServerAPI sets the server API version for this operation.
func (dd *DropDatabase) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.serverAPI = serverAPI
	return dd
}

// Authenticator sets the authenticator to use for this operation.
func (dd *DropDatabase) Authenticator(authenticator driver.Authenticator) *DropDatabase {
	if dd == nil {
		dd = new(DropDatabase)
	}

	dd.authenticator = authenticator
	return dd
}
