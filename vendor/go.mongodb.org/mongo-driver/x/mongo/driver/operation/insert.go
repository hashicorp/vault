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
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Insert performs an insert operation.
type Insert struct {
	authenticator            driver.Authenticator
	bypassDocumentValidation *bool
	comment                  bsoncore.Value
	documents                []bsoncore.Document
	ordered                  *bool
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	crypt                    driver.Crypt
	database                 string
	deployment               driver.Deployment
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	result                   InsertResult
	serverAPI                *driver.ServerAPIOptions
	timeout                  *time.Duration
	logger                   *logger.Logger
}

// InsertResult represents an insert result returned by the server.
type InsertResult struct {
	// Number of documents successfully inserted.
	N int64
}

func buildInsertResult(response bsoncore.Document) (InsertResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return InsertResult{}, err
	}
	ir := InsertResult{}
	for _, element := range elements {
		if element.Key() == "n" {
			var ok bool
			ir.N, ok = element.Value().AsInt64OK()
			if !ok {
				return ir, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return ir, nil
}

// NewInsert constructs and returns a new Insert.
func NewInsert(documents ...bsoncore.Document) *Insert {
	return &Insert{
		documents: documents,
	}
}

// Result returns the result of executing this operation.
func (i *Insert) Result() InsertResult { return i.result }

func (i *Insert) processResponse(info driver.ResponseInfo) error {
	ir, err := buildInsertResult(info.ServerResponse)
	i.result.N += ir.N
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (i *Insert) Execute(ctx context.Context) error {
	if i.deployment == nil {
		return errors.New("the Insert operation must have a Deployment set before Execute can be called")
	}
	batches := &driver.Batches{
		Identifier: "documents",
		Documents:  i.documents,
		Ordered:    i.ordered,
	}

	return driver.Operation{
		CommandFn:         i.command,
		ProcessResponseFn: i.processResponse,
		Batches:           batches,
		RetryMode:         i.retry,
		Type:              driver.Write,
		Client:            i.session,
		Clock:             i.clock,
		CommandMonitor:    i.monitor,
		Crypt:             i.crypt,
		Database:          i.database,
		Deployment:        i.deployment,
		Selector:          i.selector,
		WriteConcern:      i.writeConcern,
		ServerAPI:         i.serverAPI,
		Timeout:           i.timeout,
		Logger:            i.logger,
		Name:              driverutil.InsertOp,
		Authenticator:     i.authenticator,
	}.Execute(ctx)

}

func (i *Insert) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "insert", i.collection)
	if i.bypassDocumentValidation != nil && (desc.WireVersion != nil && desc.WireVersion.Includes(4)) {
		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *i.bypassDocumentValidation)
	}
	if i.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", i.comment)
	}
	if i.ordered != nil {
		dst = bsoncore.AppendBooleanElement(dst, "ordered", *i.ordered)
	}
	return dst, nil
}

// BypassDocumentValidation allows the operation to opt-out of document level validation. Valid
// for server versions >= 3.2. For servers < 3.2, this setting is ignored.
func (i *Insert) BypassDocumentValidation(bypassDocumentValidation bool) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.bypassDocumentValidation = &bypassDocumentValidation
	return i
}

// Comment sets a value to help trace an operation.
func (i *Insert) Comment(comment bsoncore.Value) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.comment = comment
	return i
}

// Documents adds documents to this operation that will be inserted when this operation is
// executed.
func (i *Insert) Documents(documents ...bsoncore.Document) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.documents = documents
	return i
}

// Ordered sets ordered. If true, when a write fails, the operation will return the error, when
// false write failures do not stop execution of the operation.
func (i *Insert) Ordered(ordered bool) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.ordered = &ordered
	return i
}

// Session sets the session for this operation.
func (i *Insert) Session(session *session.Client) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.session = session
	return i
}

// ClusterClock sets the cluster clock for this operation.
func (i *Insert) ClusterClock(clock *session.ClusterClock) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.clock = clock
	return i
}

// Collection sets the collection that this command will run against.
func (i *Insert) Collection(collection string) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.collection = collection
	return i
}

// CommandMonitor sets the monitor to use for APM events.
func (i *Insert) CommandMonitor(monitor *event.CommandMonitor) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.monitor = monitor
	return i
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (i *Insert) Crypt(crypt driver.Crypt) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.crypt = crypt
	return i
}

// Database sets the database to run this operation against.
func (i *Insert) Database(database string) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.database = database
	return i
}

// Deployment sets the deployment to use for this operation.
func (i *Insert) Deployment(deployment driver.Deployment) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.deployment = deployment
	return i
}

// ServerSelector sets the selector used to retrieve a server.
func (i *Insert) ServerSelector(selector description.ServerSelector) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.selector = selector
	return i
}

// WriteConcern sets the write concern for this operation.
func (i *Insert) WriteConcern(writeConcern *writeconcern.WriteConcern) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.writeConcern = writeConcern
	return i
}

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (i *Insert) Retry(retry driver.RetryMode) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.retry = &retry
	return i
}

// ServerAPI sets the server API version for this operation.
func (i *Insert) ServerAPI(serverAPI *driver.ServerAPIOptions) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.serverAPI = serverAPI
	return i
}

// Timeout sets the timeout for this operation.
func (i *Insert) Timeout(timeout *time.Duration) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.timeout = timeout
	return i
}

// Logger sets the logger for this operation.
func (i *Insert) Logger(logger *logger.Logger) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.logger = logger
	return i
}

// Authenticator sets the authenticator to use for this operation.
func (i *Insert) Authenticator(authenticator driver.Authenticator) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.authenticator = authenticator
	return i
}
