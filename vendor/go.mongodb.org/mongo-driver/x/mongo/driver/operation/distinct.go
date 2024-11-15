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

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Distinct performs a distinct operation.
type Distinct struct {
	authenticator  driver.Authenticator
	collation      bsoncore.Document
	key            *string
	maxTime        *time.Duration
	query          bsoncore.Document
	session        *session.Client
	clock          *session.ClusterClock
	collection     string
	comment        bsoncore.Value
	monitor        *event.CommandMonitor
	crypt          driver.Crypt
	database       string
	deployment     driver.Deployment
	readConcern    *readconcern.ReadConcern
	readPreference *readpref.ReadPref
	selector       description.ServerSelector
	retry          *driver.RetryMode
	result         DistinctResult
	serverAPI      *driver.ServerAPIOptions
	timeout        *time.Duration
}

// DistinctResult represents a distinct result returned by the server.
type DistinctResult struct {
	// The distinct values for the field.
	Values bsoncore.Value
}

func buildDistinctResult(response bsoncore.Document) (DistinctResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DistinctResult{}, err
	}
	dr := DistinctResult{}
	for _, element := range elements {
		if element.Key() == "values" {
			dr.Values = element.Value()
		}
	}
	return dr, nil
}

// NewDistinct constructs and returns a new Distinct.
func NewDistinct(key string, query bsoncore.Document) *Distinct {
	return &Distinct{
		key:   &key,
		query: query,
	}
}

// Result returns the result of executing this operation.
func (d *Distinct) Result() DistinctResult { return d.result }

func (d *Distinct) processResponse(info driver.ResponseInfo) error {
	var err error
	d.result, err = buildDistinctResult(info.ServerResponse)
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (d *Distinct) Execute(ctx context.Context) error {
	if d.deployment == nil {
		return errors.New("the Distinct operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         d.command,
		ProcessResponseFn: d.processResponse,
		RetryMode:         d.retry,
		Type:              driver.Read,
		Client:            d.session,
		Clock:             d.clock,
		CommandMonitor:    d.monitor,
		Crypt:             d.crypt,
		Database:          d.database,
		Deployment:        d.deployment,
		MaxTime:           d.maxTime,
		ReadConcern:       d.readConcern,
		ReadPreference:    d.readPreference,
		Selector:          d.selector,
		ServerAPI:         d.serverAPI,
		Timeout:           d.timeout,
		Name:              driverutil.DistinctOp,
		Authenticator:     d.authenticator,
	}.Execute(ctx)

}

func (d *Distinct) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "distinct", d.collection)
	if d.collation != nil {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		}
		dst = bsoncore.AppendDocumentElement(dst, "collation", d.collation)
	}
	if d.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", d.comment)
	}
	if d.key != nil {
		dst = bsoncore.AppendStringElement(dst, "key", *d.key)
	}
	if d.query != nil {
		dst = bsoncore.AppendDocumentElement(dst, "query", d.query)
	}
	return dst, nil
}

// Collation specifies a collation to be used.
func (d *Distinct) Collation(collation bsoncore.Document) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.collation = collation
	return d
}

// Key specifies which field to return distinct values for.
func (d *Distinct) Key(key string) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.key = &key
	return d
}

// MaxTime specifies the maximum amount of time to allow the query to run on the server.
func (d *Distinct) MaxTime(maxTime *time.Duration) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.maxTime = maxTime
	return d
}

// Query specifies which documents to return distinct values from.
func (d *Distinct) Query(query bsoncore.Document) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.query = query
	return d
}

// Session sets the session for this operation.
func (d *Distinct) Session(session *session.Client) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.session = session
	return d
}

// ClusterClock sets the cluster clock for this operation.
func (d *Distinct) ClusterClock(clock *session.ClusterClock) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.clock = clock
	return d
}

// Collection sets the collection that this command will run against.
func (d *Distinct) Collection(collection string) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.collection = collection
	return d
}

// Comment sets a value to help trace an operation.
func (d *Distinct) Comment(comment bsoncore.Value) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.comment = comment
	return d
}

// CommandMonitor sets the monitor to use for APM events.
func (d *Distinct) CommandMonitor(monitor *event.CommandMonitor) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.monitor = monitor
	return d
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (d *Distinct) Crypt(crypt driver.Crypt) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.crypt = crypt
	return d
}

// Database sets the database to run this operation against.
func (d *Distinct) Database(database string) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.database = database
	return d
}

// Deployment sets the deployment to use for this operation.
func (d *Distinct) Deployment(deployment driver.Deployment) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.deployment = deployment
	return d
}

// ReadConcern specifies the read concern for this operation.
func (d *Distinct) ReadConcern(readConcern *readconcern.ReadConcern) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.readConcern = readConcern
	return d
}

// ReadPreference set the read preference used with this operation.
func (d *Distinct) ReadPreference(readPreference *readpref.ReadPref) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.readPreference = readPreference
	return d
}

// ServerSelector sets the selector used to retrieve a server.
func (d *Distinct) ServerSelector(selector description.ServerSelector) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.selector = selector
	return d
}

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (d *Distinct) Retry(retry driver.RetryMode) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.retry = &retry
	return d
}

// ServerAPI sets the server API version for this operation.
func (d *Distinct) ServerAPI(serverAPI *driver.ServerAPIOptions) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.serverAPI = serverAPI
	return d
}

// Timeout sets the timeout for this operation.
func (d *Distinct) Timeout(timeout *time.Duration) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.timeout = timeout
	return d
}

// Authenticator sets the authenticator to use for this operation.
func (d *Distinct) Authenticator(authenticator driver.Authenticator) *Distinct {
	if d == nil {
		d = new(Distinct)
	}

	d.authenticator = authenticator
	return d
}
