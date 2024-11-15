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
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// FindAndModify performs a findAndModify operation.
type FindAndModify struct {
	authenticator            driver.Authenticator
	arrayFilters             bsoncore.Array
	bypassDocumentValidation *bool
	collation                bsoncore.Document
	comment                  bsoncore.Value
	fields                   bsoncore.Document
	maxTime                  *time.Duration
	newDocument              *bool
	query                    bsoncore.Document
	remove                   *bool
	sort                     bsoncore.Document
	update                   bsoncore.Value
	upsert                   *bool
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	database                 string
	deployment               driver.Deployment
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	crypt                    driver.Crypt
	hint                     bsoncore.Value
	serverAPI                *driver.ServerAPIOptions
	let                      bsoncore.Document
	timeout                  *time.Duration

	result FindAndModifyResult
}

// LastErrorObject represents information about updates and upserts returned by the server.
type LastErrorObject struct {
	// True if an update modified an existing document
	UpdatedExisting bool
	// Object ID of the upserted document.
	Upserted interface{}
}

// FindAndModifyResult represents a findAndModify result returned by the server.
type FindAndModifyResult struct {
	// Either the old or modified document, depending on the value of the new parameter.
	Value bsoncore.Document
	// Contains information about updates and upserts.
	LastErrorObject LastErrorObject
}

func buildFindAndModifyResult(response bsoncore.Document) (FindAndModifyResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return FindAndModifyResult{}, err
	}
	famr := FindAndModifyResult{}
	for _, element := range elements {
		switch element.Key() {
		case "value":
			var ok bool
			famr.Value, ok = element.Value().DocumentOK()

			// The 'value' field returned by a FindAndModify can be null in the case that no document was found.
			if element.Value().Type != bsontype.Null && !ok {
				return famr, fmt.Errorf("response field 'value' is type document or null, but received BSON type %s", element.Value().Type)
			}
		case "lastErrorObject":
			valDoc, ok := element.Value().DocumentOK()
			if !ok {
				return famr, fmt.Errorf("response field 'lastErrorObject' is type document, but received BSON type %s", element.Value().Type)
			}

			var leo LastErrorObject
			if err = bson.Unmarshal(valDoc, &leo); err != nil {
				return famr, err
			}
			famr.LastErrorObject = leo
		}
	}
	return famr, nil
}

// NewFindAndModify constructs and returns a new FindAndModify.
func NewFindAndModify(query bsoncore.Document) *FindAndModify {
	return &FindAndModify{
		query: query,
	}
}

// Result returns the result of executing this operation.
func (fam *FindAndModify) Result() FindAndModifyResult { return fam.result }

func (fam *FindAndModify) processResponse(info driver.ResponseInfo) error {
	var err error

	fam.result, err = buildFindAndModifyResult(info.ServerResponse)
	return err

}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (fam *FindAndModify) Execute(ctx context.Context) error {
	if fam.deployment == nil {
		return errors.New("the FindAndModify operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         fam.command,
		ProcessResponseFn: fam.processResponse,

		RetryMode:      fam.retry,
		Type:           driver.Write,
		Client:         fam.session,
		Clock:          fam.clock,
		CommandMonitor: fam.monitor,
		Database:       fam.database,
		Deployment:     fam.deployment,
		MaxTime:        fam.maxTime,
		Selector:       fam.selector,
		WriteConcern:   fam.writeConcern,
		Crypt:          fam.crypt,
		ServerAPI:      fam.serverAPI,
		Timeout:        fam.timeout,
		Name:           driverutil.FindAndModifyOp,
		Authenticator:  fam.authenticator,
	}.Execute(ctx)

}

func (fam *FindAndModify) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "findAndModify", fam.collection)
	if fam.arrayFilters != nil {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(6) {
			return nil, errors.New("the 'arrayFilters' command parameter requires a minimum server wire version of 6")
		}
		dst = bsoncore.AppendArrayElement(dst, "arrayFilters", fam.arrayFilters)
	}
	if fam.bypassDocumentValidation != nil {

		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *fam.bypassDocumentValidation)
	}
	if fam.collation != nil {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		}
		dst = bsoncore.AppendDocumentElement(dst, "collation", fam.collation)
	}
	if fam.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", fam.comment)
	}
	if fam.fields != nil {

		dst = bsoncore.AppendDocumentElement(dst, "fields", fam.fields)
	}
	if fam.newDocument != nil {

		dst = bsoncore.AppendBooleanElement(dst, "new", *fam.newDocument)
	}
	if fam.query != nil {

		dst = bsoncore.AppendDocumentElement(dst, "query", fam.query)
	}
	if fam.remove != nil {

		dst = bsoncore.AppendBooleanElement(dst, "remove", *fam.remove)
	}
	if fam.sort != nil {

		dst = bsoncore.AppendDocumentElement(dst, "sort", fam.sort)
	}
	if fam.update.Data != nil {
		dst = bsoncore.AppendValueElement(dst, "update", fam.update)
	}
	if fam.upsert != nil {

		dst = bsoncore.AppendBooleanElement(dst, "upsert", *fam.upsert)
	}
	if fam.hint.Type != bsontype.Type(0) {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(8) {
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 8")
		}
		if !fam.writeConcern.Acknowledged() {
			return nil, errUnacknowledgedHint
		}
		dst = bsoncore.AppendValueElement(dst, "hint", fam.hint)
	}
	if fam.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", fam.let)
	}

	return dst, nil
}

// ArrayFilters specifies an array of filter documents that determines which array elements to modify for an update operation on an array field.
func (fam *FindAndModify) ArrayFilters(arrayFilters bsoncore.Array) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.arrayFilters = arrayFilters
	return fam
}

// BypassDocumentValidation specifies if document validation can be skipped when executing the operation.
func (fam *FindAndModify) BypassDocumentValidation(bypassDocumentValidation bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.bypassDocumentValidation = &bypassDocumentValidation
	return fam
}

// Collation specifies a collation to be used.
func (fam *FindAndModify) Collation(collation bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.collation = collation
	return fam
}

// Comment sets a value to help trace an operation.
func (fam *FindAndModify) Comment(comment bsoncore.Value) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.comment = comment
	return fam
}

// Fields specifies a subset of fields to return.
func (fam *FindAndModify) Fields(fields bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.fields = fields
	return fam
}

// MaxTime specifies the maximum amount of time to allow the operation to run on the server.
func (fam *FindAndModify) MaxTime(maxTime *time.Duration) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.maxTime = maxTime
	return fam
}

// NewDocument specifies whether to return the modified document or the original. Defaults to false (return original).
func (fam *FindAndModify) NewDocument(newDocument bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.newDocument = &newDocument
	return fam
}

// Query specifies the selection criteria for the modification.
func (fam *FindAndModify) Query(query bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.query = query
	return fam
}

// Remove specifies that the matched document should be removed. Defaults to false.
func (fam *FindAndModify) Remove(remove bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.remove = &remove
	return fam
}

// Sort determines which document the operation modifies if the query matches multiple documents.The first document matched by the sort order will be modified.
func (fam *FindAndModify) Sort(sort bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.sort = sort
	return fam
}

// Update specifies the update document to perform on the matched document.
func (fam *FindAndModify) Update(update bsoncore.Value) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.update = update
	return fam
}

// Upsert specifies whether or not to create a new document if no documents match the query when doing an update. Defaults to false.
func (fam *FindAndModify) Upsert(upsert bool) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.upsert = &upsert
	return fam
}

// Session sets the session for this operation.
func (fam *FindAndModify) Session(session *session.Client) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.session = session
	return fam
}

// ClusterClock sets the cluster clock for this operation.
func (fam *FindAndModify) ClusterClock(clock *session.ClusterClock) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.clock = clock
	return fam
}

// Collection sets the collection that this command will run against.
func (fam *FindAndModify) Collection(collection string) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.collection = collection
	return fam
}

// CommandMonitor sets the monitor to use for APM events.
func (fam *FindAndModify) CommandMonitor(monitor *event.CommandMonitor) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.monitor = monitor
	return fam
}

// Database sets the database to run this operation against.
func (fam *FindAndModify) Database(database string) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.database = database
	return fam
}

// Deployment sets the deployment to use for this operation.
func (fam *FindAndModify) Deployment(deployment driver.Deployment) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.deployment = deployment
	return fam
}

// ServerSelector sets the selector used to retrieve a server.
func (fam *FindAndModify) ServerSelector(selector description.ServerSelector) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.selector = selector
	return fam
}

// WriteConcern sets the write concern for this operation.
func (fam *FindAndModify) WriteConcern(writeConcern *writeconcern.WriteConcern) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.writeConcern = writeConcern
	return fam
}

// Retry enables retryable writes for this operation. Retries are not handled automatically,
// instead a boolean is returned from Execute and SelectAndExecute that indicates if the
// operation can be retried. Retrying is handled by calling RetryExecute.
func (fam *FindAndModify) Retry(retry driver.RetryMode) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.retry = &retry
	return fam
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (fam *FindAndModify) Crypt(crypt driver.Crypt) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.crypt = crypt
	return fam
}

// Hint specifies the index to use.
func (fam *FindAndModify) Hint(hint bsoncore.Value) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.hint = hint
	return fam
}

// ServerAPI sets the server API version for this operation.
func (fam *FindAndModify) ServerAPI(serverAPI *driver.ServerAPIOptions) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.serverAPI = serverAPI
	return fam
}

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (fam *FindAndModify) Let(let bsoncore.Document) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.let = let
	return fam
}

// Timeout sets the timeout for this operation.
func (fam *FindAndModify) Timeout(timeout *time.Duration) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.timeout = timeout
	return fam
}

// Authenticator sets the authenticator to use for this operation.
func (fam *FindAndModify) Authenticator(authenticator driver.Authenticator) *FindAndModify {
	if fam == nil {
		fam = new(FindAndModify)
	}

	fam.authenticator = authenticator
	return fam
}
