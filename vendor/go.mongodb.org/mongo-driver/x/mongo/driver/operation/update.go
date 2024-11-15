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

	"go.mongodb.org/mongo-driver/bson"
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

// Update performs an update operation.
type Update struct {
	authenticator            driver.Authenticator
	bypassDocumentValidation *bool
	comment                  bsoncore.Value
	ordered                  *bool
	updates                  []bsoncore.Document
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	database                 string
	deployment               driver.Deployment
	hint                     *bool
	arrayFilters             *bool
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	result                   UpdateResult
	crypt                    driver.Crypt
	serverAPI                *driver.ServerAPIOptions
	let                      bsoncore.Document
	timeout                  *time.Duration
	logger                   *logger.Logger
}

// Upsert contains the information for an upsert in an Update operation.
type Upsert struct {
	Index int64
	ID    interface{} `bson:"_id"`
}

// UpdateResult contains information for the result of an Update operation.
type UpdateResult struct {
	// Number of documents matched.
	N int64
	// Number of documents modified.
	NModified int64
	// Information about upserted documents.
	Upserted []Upsert
}

func buildUpdateResult(response bsoncore.Document) (UpdateResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return UpdateResult{}, err
	}
	ur := UpdateResult{}
	for _, element := range elements {
		switch element.Key() {
		case "nModified":
			var ok bool
			ur.NModified, ok = element.Value().AsInt64OK()
			if !ok {
				return ur, fmt.Errorf("response field 'nModified' is type int32 or int64, but received BSON type %s", element.Value().Type)
			}
		case "n":
			var ok bool
			ur.N, ok = element.Value().AsInt64OK()
			if !ok {
				return ur, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			}
		case "upserted":
			arr, ok := element.Value().ArrayOK()
			if !ok {
				return ur, fmt.Errorf("response field 'upserted' is type array, but received BSON type %s", element.Value().Type)
			}

			var values []bsoncore.Value
			values, err = arr.Values()
			if err != nil {
				break
			}

			for _, val := range values {
				valDoc, ok := val.DocumentOK()
				if !ok {
					return ur, fmt.Errorf("upserted value is type document, but received BSON type %s", val.Type)
				}
				var upsert Upsert
				if err = bson.Unmarshal(valDoc, &upsert); err != nil {
					return ur, err
				}
				ur.Upserted = append(ur.Upserted, upsert)
			}
		}
	}
	return ur, nil
}

// NewUpdate constructs and returns a new Update.
func NewUpdate(updates ...bsoncore.Document) *Update {
	return &Update{
		updates: updates,
	}
}

// Result returns the result of executing this operation.
func (u *Update) Result() UpdateResult { return u.result }

func (u *Update) processResponse(info driver.ResponseInfo) error {
	ur, err := buildUpdateResult(info.ServerResponse)

	u.result.N += ur.N
	u.result.NModified += ur.NModified
	if info.CurrentIndex > 0 {
		for ind := range ur.Upserted {
			ur.Upserted[ind].Index += int64(info.CurrentIndex)
		}
	}
	u.result.Upserted = append(u.result.Upserted, ur.Upserted...)
	return err

}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (u *Update) Execute(ctx context.Context) error {
	if u.deployment == nil {
		return errors.New("the Update operation must have a Deployment set before Execute can be called")
	}
	batches := &driver.Batches{
		Identifier: "updates",
		Documents:  u.updates,
		Ordered:    u.ordered,
	}

	return driver.Operation{
		CommandFn:         u.command,
		ProcessResponseFn: u.processResponse,
		Batches:           batches,
		RetryMode:         u.retry,
		Type:              driver.Write,
		Client:            u.session,
		Clock:             u.clock,
		CommandMonitor:    u.monitor,
		Database:          u.database,
		Deployment:        u.deployment,
		Selector:          u.selector,
		WriteConcern:      u.writeConcern,
		Crypt:             u.crypt,
		ServerAPI:         u.serverAPI,
		Timeout:           u.timeout,
		Logger:            u.logger,
		Name:              driverutil.UpdateOp,
		Authenticator:     u.authenticator,
	}.Execute(ctx)

}

func (u *Update) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "update", u.collection)
	if u.bypassDocumentValidation != nil &&
		(desc.WireVersion != nil && desc.WireVersion.Includes(4)) {

		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *u.bypassDocumentValidation)
	}
	if u.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", u.comment)
	}
	if u.ordered != nil {

		dst = bsoncore.AppendBooleanElement(dst, "ordered", *u.ordered)
	}
	if u.hint != nil && *u.hint {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 5")
		}
		if !u.writeConcern.Acknowledged() {
			return nil, errUnacknowledgedHint
		}
	}
	if u.arrayFilters != nil && *u.arrayFilters {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(6) {
			return nil, errors.New("the 'arrayFilters' command parameter requires a minimum server wire version of 6")
		}
	}
	if u.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", u.let)
	}

	return dst, nil
}

// BypassDocumentValidation allows the operation to opt-out of document level validation. Valid
// for server versions >= 3.2. For servers < 3.2, this setting is ignored.
func (u *Update) BypassDocumentValidation(bypassDocumentValidation bool) *Update {
	if u == nil {
		u = new(Update)
	}

	u.bypassDocumentValidation = &bypassDocumentValidation
	return u
}

// Hint is a flag to indicate that the update document contains a hint. Hint is only supported by
// servers >= 4.2. Older servers >= 3.4 will report an error for using the hint option. For servers <
// 3.4, the driver will return an error if the hint option is used.
func (u *Update) Hint(hint bool) *Update {
	if u == nil {
		u = new(Update)
	}

	u.hint = &hint
	return u
}

// ArrayFilters is a flag to indicate that the update document contains an arrayFilters field. This option is only
// supported on server versions 3.6 and higher. For servers < 3.6, the driver will return an error.
func (u *Update) ArrayFilters(arrayFilters bool) *Update {
	if u == nil {
		u = new(Update)
	}

	u.arrayFilters = &arrayFilters
	return u
}

// Ordered sets ordered. If true, when a write fails, the operation will return the error, when
// false write failures do not stop execution of the operation.
func (u *Update) Ordered(ordered bool) *Update {
	if u == nil {
		u = new(Update)
	}

	u.ordered = &ordered
	return u
}

// Updates specifies an array of update statements to perform when this operation is executed.
// Each update document must have the following structure:
// {q: <query>, u: <update>, multi: <boolean>, collation: Optional<Document>, arrayFitlers: Optional<Array>, hint: Optional<string/Document>}.
func (u *Update) Updates(updates ...bsoncore.Document) *Update {
	if u == nil {
		u = new(Update)
	}

	u.updates = updates
	return u
}

// Session sets the session for this operation.
func (u *Update) Session(session *session.Client) *Update {
	if u == nil {
		u = new(Update)
	}

	u.session = session
	return u
}

// ClusterClock sets the cluster clock for this operation.
func (u *Update) ClusterClock(clock *session.ClusterClock) *Update {
	if u == nil {
		u = new(Update)
	}

	u.clock = clock
	return u
}

// Collection sets the collection that this command will run against.
func (u *Update) Collection(collection string) *Update {
	if u == nil {
		u = new(Update)
	}

	u.collection = collection
	return u
}

// CommandMonitor sets the monitor to use for APM events.
func (u *Update) CommandMonitor(monitor *event.CommandMonitor) *Update {
	if u == nil {
		u = new(Update)
	}

	u.monitor = monitor
	return u
}

// Comment sets a value to help trace an operation.
func (u *Update) Comment(comment bsoncore.Value) *Update {
	if u == nil {
		u = new(Update)
	}

	u.comment = comment
	return u
}

// Database sets the database to run this operation against.
func (u *Update) Database(database string) *Update {
	if u == nil {
		u = new(Update)
	}

	u.database = database
	return u
}

// Deployment sets the deployment to use for this operation.
func (u *Update) Deployment(deployment driver.Deployment) *Update {
	if u == nil {
		u = new(Update)
	}

	u.deployment = deployment
	return u
}

// ServerSelector sets the selector used to retrieve a server.
func (u *Update) ServerSelector(selector description.ServerSelector) *Update {
	if u == nil {
		u = new(Update)
	}

	u.selector = selector
	return u
}

// WriteConcern sets the write concern for this operation.
func (u *Update) WriteConcern(writeConcern *writeconcern.WriteConcern) *Update {
	if u == nil {
		u = new(Update)
	}

	u.writeConcern = writeConcern
	return u
}

// Retry enables retryable writes for this operation. Retries are not handled automatically,
// instead a boolean is returned from Execute and SelectAndExecute that indicates if the
// operation can be retried. Retrying is handled by calling RetryExecute.
func (u *Update) Retry(retry driver.RetryMode) *Update {
	if u == nil {
		u = new(Update)
	}

	u.retry = &retry
	return u
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (u *Update) Crypt(crypt driver.Crypt) *Update {
	if u == nil {
		u = new(Update)
	}

	u.crypt = crypt
	return u
}

// ServerAPI sets the server API version for this operation.
func (u *Update) ServerAPI(serverAPI *driver.ServerAPIOptions) *Update {
	if u == nil {
		u = new(Update)
	}

	u.serverAPI = serverAPI
	return u
}

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (u *Update) Let(let bsoncore.Document) *Update {
	if u == nil {
		u = new(Update)
	}

	u.let = let
	return u
}

// Timeout sets the timeout for this operation.
func (u *Update) Timeout(timeout *time.Duration) *Update {
	if u == nil {
		u = new(Update)
	}

	u.timeout = timeout
	return u
}

// Logger sets the logger for this operation.
func (u *Update) Logger(logger *logger.Logger) *Update {
	if u == nil {
		u = new(Update)
	}

	u.logger = logger
	return u
}

// Authenticator sets the authenticator to use for this operation.
func (u *Update) Authenticator(authenticator driver.Authenticator) *Update {
	if u == nil {
		u = new(Update)
	}

	u.authenticator = authenticator
	return u
}
