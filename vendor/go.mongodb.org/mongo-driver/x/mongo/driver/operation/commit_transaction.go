// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// CommitTransaction attempts to commit a transaction.
type CommitTransaction struct {
	authenticator driver.Authenticator
	maxTime       *time.Duration
	recoveryToken bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	retry         *driver.RetryMode
	serverAPI     *driver.ServerAPIOptions
}

// NewCommitTransaction constructs and returns a new CommitTransaction.
func NewCommitTransaction() *CommitTransaction {
	return &CommitTransaction{}
}

func (ct *CommitTransaction) processResponse(driver.ResponseInfo) error {
	var err error
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (ct *CommitTransaction) Execute(ctx context.Context) error {
	if ct.deployment == nil {
		return errors.New("the CommitTransaction operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         ct.command,
		ProcessResponseFn: ct.processResponse,
		RetryMode:         ct.retry,
		Type:              driver.Write,
		Client:            ct.session,
		Clock:             ct.clock,
		CommandMonitor:    ct.monitor,
		Crypt:             ct.crypt,
		Database:          ct.database,
		Deployment:        ct.deployment,
		MaxTime:           ct.maxTime,
		Selector:          ct.selector,
		WriteConcern:      ct.writeConcern,
		ServerAPI:         ct.serverAPI,
		Name:              driverutil.CommitTransactionOp,
		Authenticator:     ct.authenticator,
	}.Execute(ctx)

}

func (ct *CommitTransaction) command(dst []byte, _ description.SelectedServer) ([]byte, error) {

	dst = bsoncore.AppendInt32Element(dst, "commitTransaction", 1)
	if ct.recoveryToken != nil {
		dst = bsoncore.AppendDocumentElement(dst, "recoveryToken", ct.recoveryToken)
	}
	return dst, nil
}

// MaxTime specifies the maximum amount of time to allow the query to run on the server.
func (ct *CommitTransaction) MaxTime(maxTime *time.Duration) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.maxTime = maxTime
	return ct
}

// RecoveryToken sets the recovery token to use when committing or aborting a sharded transaction.
func (ct *CommitTransaction) RecoveryToken(recoveryToken bsoncore.Document) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.recoveryToken = recoveryToken
	return ct
}

// Session sets the session for this operation.
func (ct *CommitTransaction) Session(session *session.Client) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.session = session
	return ct
}

// ClusterClock sets the cluster clock for this operation.
func (ct *CommitTransaction) ClusterClock(clock *session.ClusterClock) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.clock = clock
	return ct
}

// CommandMonitor sets the monitor to use for APM events.
func (ct *CommitTransaction) CommandMonitor(monitor *event.CommandMonitor) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.monitor = monitor
	return ct
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (ct *CommitTransaction) Crypt(crypt driver.Crypt) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.crypt = crypt
	return ct
}

// Database sets the database to run this operation against.
func (ct *CommitTransaction) Database(database string) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.database = database
	return ct
}

// Deployment sets the deployment to use for this operation.
func (ct *CommitTransaction) Deployment(deployment driver.Deployment) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.deployment = deployment
	return ct
}

// ServerSelector sets the selector used to retrieve a server.
func (ct *CommitTransaction) ServerSelector(selector description.ServerSelector) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.selector = selector
	return ct
}

// WriteConcern sets the write concern for this operation.
func (ct *CommitTransaction) WriteConcern(writeConcern *writeconcern.WriteConcern) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.writeConcern = writeConcern
	return ct
}

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (ct *CommitTransaction) Retry(retry driver.RetryMode) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.retry = &retry
	return ct
}

// ServerAPI sets the server API version for this operation.
func (ct *CommitTransaction) ServerAPI(serverAPI *driver.ServerAPIOptions) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.serverAPI = serverAPI
	return ct
}

// Authenticator sets the authenticator to use for this operation.
func (ct *CommitTransaction) Authenticator(authenticator driver.Authenticator) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.authenticator = authenticator
	return ct
}
