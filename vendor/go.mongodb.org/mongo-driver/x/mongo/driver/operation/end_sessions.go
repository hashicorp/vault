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
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// EndSessions performs an endSessions operation.
type EndSessions struct {
	authenticator driver.Authenticator
	sessionIDs    bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	serverAPI     *driver.ServerAPIOptions
}

// NewEndSessions constructs and returns a new EndSessions.
func NewEndSessions(sessionIDs bsoncore.Document) *EndSessions {
	return &EndSessions{
		sessionIDs: sessionIDs,
	}
}

func (es *EndSessions) processResponse(driver.ResponseInfo) error {
	var err error
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (es *EndSessions) Execute(ctx context.Context) error {
	if es.deployment == nil {
		return errors.New("the EndSessions operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         es.command,
		ProcessResponseFn: es.processResponse,
		Client:            es.session,
		Clock:             es.clock,
		CommandMonitor:    es.monitor,
		Crypt:             es.crypt,
		Database:          es.database,
		Deployment:        es.deployment,
		Selector:          es.selector,
		ServerAPI:         es.serverAPI,
		Name:              driverutil.EndSessionsOp,
		Authenticator:     es.authenticator,
	}.Execute(ctx)

}

func (es *EndSessions) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	if es.sessionIDs != nil {
		dst = bsoncore.AppendArrayElement(dst, "endSessions", es.sessionIDs)
	}
	return dst, nil
}

// SessionIDs specifies the sessions to be expired.
func (es *EndSessions) SessionIDs(sessionIDs bsoncore.Document) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.sessionIDs = sessionIDs
	return es
}

// Session sets the session for this operation.
func (es *EndSessions) Session(session *session.Client) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.session = session
	return es
}

// ClusterClock sets the cluster clock for this operation.
func (es *EndSessions) ClusterClock(clock *session.ClusterClock) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.clock = clock
	return es
}

// CommandMonitor sets the monitor to use for APM events.
func (es *EndSessions) CommandMonitor(monitor *event.CommandMonitor) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.monitor = monitor
	return es
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (es *EndSessions) Crypt(crypt driver.Crypt) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.crypt = crypt
	return es
}

// Database sets the database to run this operation against.
func (es *EndSessions) Database(database string) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.database = database
	return es
}

// Deployment sets the deployment to use for this operation.
func (es *EndSessions) Deployment(deployment driver.Deployment) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.deployment = deployment
	return es
}

// ServerSelector sets the selector used to retrieve a server.
func (es *EndSessions) ServerSelector(selector description.ServerSelector) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.selector = selector
	return es
}

// ServerAPI sets the server API version for this operation.
func (es *EndSessions) ServerAPI(serverAPI *driver.ServerAPIOptions) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.serverAPI = serverAPI
	return es
}

// Authenticator sets the authenticator to use for this operation.
func (es *EndSessions) Authenticator(authenticator driver.Authenticator) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.authenticator = authenticator
	return es
}
