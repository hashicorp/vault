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

// Delete performs a delete operation
type Delete struct {
	authenticator driver.Authenticator
	comment       bsoncore.Value
	deletes       []bsoncore.Document
	ordered       *bool
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	retry         *driver.RetryMode
	hint          *bool
	result        DeleteResult
	serverAPI     *driver.ServerAPIOptions
	let           bsoncore.Document
	timeout       *time.Duration
	logger        *logger.Logger
}

// DeleteResult represents a delete result returned by the server.
type DeleteResult struct {
	// Number of documents successfully deleted.
	N int64
}

func buildDeleteResult(response bsoncore.Document) (DeleteResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DeleteResult{}, err
	}
	dr := DeleteResult{}
	for _, element := range elements {
		if element.Key() == "n" {
			var ok bool
			dr.N, ok = element.Value().AsInt64OK()
			if !ok {
				return dr, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return dr, nil
}

// NewDelete constructs and returns a new Delete.
func NewDelete(deletes ...bsoncore.Document) *Delete {
	return &Delete{
		deletes: deletes,
	}
}

// Result returns the result of executing this operation.
func (d *Delete) Result() DeleteResult { return d.result }

func (d *Delete) processResponse(info driver.ResponseInfo) error {
	dr, err := buildDeleteResult(info.ServerResponse)
	d.result.N += dr.N
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (d *Delete) Execute(ctx context.Context) error {
	if d.deployment == nil {
		return errors.New("the Delete operation must have a Deployment set before Execute can be called")
	}
	batches := &driver.Batches{
		Identifier: "deletes",
		Documents:  d.deletes,
		Ordered:    d.ordered,
	}

	return driver.Operation{
		CommandFn:         d.command,
		ProcessResponseFn: d.processResponse,
		Batches:           batches,
		RetryMode:         d.retry,
		Type:              driver.Write,
		Client:            d.session,
		Clock:             d.clock,
		CommandMonitor:    d.monitor,
		Crypt:             d.crypt,
		Database:          d.database,
		Deployment:        d.deployment,
		Selector:          d.selector,
		WriteConcern:      d.writeConcern,
		ServerAPI:         d.serverAPI,
		Timeout:           d.timeout,
		Logger:            d.logger,
		Name:              driverutil.DeleteOp,
		Authenticator:     d.authenticator,
	}.Execute(ctx)

}

func (d *Delete) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "delete", d.collection)
	if d.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", d.comment)
	}
	if d.ordered != nil {
		dst = bsoncore.AppendBooleanElement(dst, "ordered", *d.ordered)
	}
	if d.hint != nil && *d.hint {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 5")
		}
		if !d.writeConcern.Acknowledged() {
			return nil, errUnacknowledgedHint
		}
	}
	if d.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", d.let)
	}
	return dst, nil
}

// Deletes adds documents to this operation that will be used to determine what documents to delete when this operation
// is executed. These documents should have the form {q: <query>, limit: <integer limit>, collation: <document>}. The
// collation field is optional. If limit is 0, there will be no limit on the number of documents deleted.
func (d *Delete) Deletes(deletes ...bsoncore.Document) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.deletes = deletes
	return d
}

// Ordered sets ordered. If true, when a write fails, the operation will return the error, when
// false write failures do not stop execution of the operation.
func (d *Delete) Ordered(ordered bool) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.ordered = &ordered
	return d
}

// Session sets the session for this operation.
func (d *Delete) Session(session *session.Client) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.session = session
	return d
}

// ClusterClock sets the cluster clock for this operation.
func (d *Delete) ClusterClock(clock *session.ClusterClock) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.clock = clock
	return d
}

// Collection sets the collection that this command will run against.
func (d *Delete) Collection(collection string) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.collection = collection
	return d
}

// Comment sets a value to help trace an operation.
func (d *Delete) Comment(comment bsoncore.Value) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.comment = comment
	return d
}

// CommandMonitor sets the monitor to use for APM events.
func (d *Delete) CommandMonitor(monitor *event.CommandMonitor) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.monitor = monitor
	return d
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (d *Delete) Crypt(crypt driver.Crypt) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.crypt = crypt
	return d
}

// Database sets the database to run this operation against.
func (d *Delete) Database(database string) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.database = database
	return d
}

// Deployment sets the deployment to use for this operation.
func (d *Delete) Deployment(deployment driver.Deployment) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.deployment = deployment
	return d
}

// ServerSelector sets the selector used to retrieve a server.
func (d *Delete) ServerSelector(selector description.ServerSelector) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.selector = selector
	return d
}

// WriteConcern sets the write concern for this operation.
func (d *Delete) WriteConcern(writeConcern *writeconcern.WriteConcern) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.writeConcern = writeConcern
	return d
}

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (d *Delete) Retry(retry driver.RetryMode) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.retry = &retry
	return d
}

// Hint is a flag to indicate that the update document contains a hint. Hint is only supported by
// servers >= 4.4. Older servers >= 3.4 will report an error for using the hint option. For servers <
// 3.4, the driver will return an error if the hint option is used.
func (d *Delete) Hint(hint bool) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.hint = &hint
	return d
}

// ServerAPI sets the server API version for this operation.
func (d *Delete) ServerAPI(serverAPI *driver.ServerAPIOptions) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.serverAPI = serverAPI
	return d
}

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (d *Delete) Let(let bsoncore.Document) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.let = let
	return d
}

// Timeout sets the timeout for this operation.
func (d *Delete) Timeout(timeout *time.Duration) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.timeout = timeout
	return d
}

// Logger sets the logger for this operation.
func (d *Delete) Logger(logger *logger.Logger) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.logger = logger

	return d
}

// Authenticator sets the authenticator to use for this operation.
func (d *Delete) Authenticator(authenticator driver.Authenticator) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.authenticator = authenticator
	return d
}
