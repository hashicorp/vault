// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

var (
	firstBatchIdentifier     = "firstBatch"
	nextBatchIdentifier      = "nextBatch"
	listCollectionsNamespace = "system.namespaces"
	listIndexesNamespace     = "system.indexes"

	// ErrFilterType is returned when the filter for a legacy list collections operation is of the wrong type.
	ErrFilterType = errors.New("filter for list collections operation must be a string")
)

func (op Operation) getFullCollectionName(coll string) string {
	return op.Database + "." + coll
}

func (op Operation) legacyFind(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error {
	wm, startedInfo, collName, err := op.createLegacyFindWireMessage(dst, desc)
	if err != nil {
		return err
	}
	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation{
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	}

	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, firstBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil {
		return finishedInfo.cmdErr
	}

	if op.ProcessResponseFn != nil {
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo{
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		}
		return op.ProcessResponseFn(info)
	}
	return nil
}

// returns wire message, collection name, error
func (op Operation) createLegacyFindWireMessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) {
	info := startedInformation{
		requestID: wiremessage.NextRequestID(),
		cmdName:   "find",
	}

	// call CommandFn on an empty slice rather than dst because the options will need to be converted to legacy
	var cmdDoc bsoncore.Document
	var cmdIndex int32
	var err error

	cmdIndex, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil {
		return dst, info, "", err
	}
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIndex)
	// for monitoring legacy events, the upconverted document should be captured rather than the legacy one
	info.cmd = cmdDoc

	cmdElems, err := cmdDoc.Elements()
	if err != nil {
		return dst, info, "", err
	}

	// take each option from the non-legacy command and convert it
	// build options as a byte slice of elements rather than a bsoncore.Document because they will be appended
	// to another document with $query
	var optsElems []byte
	flags := op.slaveOK(desc)
	var numToSkip, numToReturn, batchSize, limit int32 // numToReturn calculated from batchSize and limit
	var filter, returnFieldsSelector bsoncore.Document
	var collName string
	var singleBatch bool
	for _, elem := range cmdElems {
		switch elem.Key() {
		case "find":
			collName = elem.Value().StringValue()
		case "filter":
			filter = elem.Value().Data
		case "sort":
			optsElems = bsoncore.AppendValueElement(optsElems, "$orderby", elem.Value())
		case "hint":
			optsElems = bsoncore.AppendValueElement(optsElems, "$hint", elem.Value())
		case "comment":
			optsElems = bsoncore.AppendValueElement(optsElems, "$comment", elem.Value())
		case "max":
			optsElems = bsoncore.AppendValueElement(optsElems, "$max", elem.Value())
		case "min":
			optsElems = bsoncore.AppendValueElement(optsElems, "$min", elem.Value())
		case "returnKey":
			optsElems = bsoncore.AppendValueElement(optsElems, "$returnKey", elem.Value())
		case "showRecordId":
			optsElems = bsoncore.AppendValueElement(optsElems, "$showDiskLoc", elem.Value())
		case "maxTimeMS":
			optsElems = bsoncore.AppendValueElement(optsElems, "$maxTimeMS", elem.Value())
		case "snapshot":
			optsElems = bsoncore.AppendValueElement(optsElems, "$snapshot", elem.Value())
		case "projection":
			returnFieldsSelector = elem.Value().Data
		case "skip":
			// CRUD spec declares skip as int64 but numToSkip is int32 in OP_QUERY
			numToSkip = int32(elem.Value().Int64())
		case "batchSize":
			batchSize = elem.Value().Int32()
			// Not possible to use batchSize = 1 because cursor will be closed on first batch
			if batchSize == 1 {
				batchSize = 2
			}
		case "limit":
			// CRUD spec declares limit as int64 but numToReturn is int32 in OP_QUERY
			limit = int32(elem.Value().Int64())
		case "singleBatch":
			singleBatch = elem.Value().Boolean()
		case "tailable":
			flags |= wiremessage.TailableCursor
		case "awaitData":
			flags |= wiremessage.AwaitData
		case "oplogReplay":
			flags |= wiremessage.OplogReplay
		case "noCursorTimeout":
			flags |= wiremessage.NoCursorTimeout
		case "allowPartialResults":
			flags |= wiremessage.Partial
		}
	}

	// for non-legacy servers, a negative limit is implemented as a positive limit + singleBatch = true
	if singleBatch {
		limit = limit * -1
	}
	numToReturn = op.calculateNumberToReturn(limit, batchSize)

	// add read preference if needed
	rp, err := op.createReadPref(desc.Server.Kind, desc.Kind, true)
	if err != nil {
		return dst, info, "", err
	}
	if len(rp) > 0 {
		optsElems = bsoncore.AppendDocumentElement(optsElems, "$readPreference", rp)
	}

	if len(filter) == 0 {
		var fidx int32
		fidx, filter = bsoncore.AppendDocumentStart(filter)
		filter, _ = bsoncore.AppendDocumentEnd(filter, fidx)
	}

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, flags)
	dst = wiremessage.AppendQueryFullCollectionName(dst, op.getFullCollectionName(collName))
	dst = wiremessage.AppendQueryNumberToSkip(dst, numToSkip)
	dst = wiremessage.AppendQueryNumberToReturn(dst, numToReturn)
	dst = op.appendLegacyQueryDocument(dst, filter, optsElems)
	if len(returnFieldsSelector) != 0 {
		// returnFieldsSelector is optional
		dst = append(dst, returnFieldsSelector...)
	}

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, collName, nil
}

func (op Operation) calculateNumberToReturn(limit, batchSize int32) int32 {
	var numToReturn int32

	if limit < 0 {
		numToReturn = limit
	} else if limit == 0 {
		numToReturn = batchSize
	} else if batchSize == 0 {
		numToReturn = limit
	} else if limit < batchSize {
		numToReturn = limit
	} else {
		numToReturn = batchSize
	}

	return numToReturn
}

func (op Operation) legacyGetMore(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error {
	wm, startedInfo, collName, err := op.createLegacyGetMoreWiremessage(dst, desc)
	if err != nil {
		return err
	}

	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation{
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	}
	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, nextBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil {
		return finishedInfo.cmdErr
	}

	if op.ProcessResponseFn != nil {
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo{
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		}
		return op.ProcessResponseFn(info)
	}
	return nil
}

func (op Operation) createLegacyGetMoreWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) {
	info := startedInformation{
		requestID: wiremessage.NextRequestID(),
		cmdName:   "getMore",
	}

	var cmdDoc bsoncore.Document
	var cmdIdx int32
	var err error

	cmdIdx, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil {
		return dst, info, "", err
	}
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIdx)
	info.cmd = cmdDoc

	cmdElems, err := cmdDoc.Elements()
	if err != nil {
		return dst, info, "", err
	}

	var cursorID int64
	var numToReturn int32
	var collName string
	for _, elem := range cmdElems {
		switch elem.Key() {
		case "getMore":
			cursorID = elem.Value().Int64()
		case "collection":
			collName = elem.Value().StringValue()
		case "batchSize":
			numToReturn = elem.Value().Int32()
		}
	}

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpGetMore)
	dst = wiremessage.AppendGetMoreZero(dst)
	dst = wiremessage.AppendGetMoreFullCollectionName(dst, op.getFullCollectionName(collName))
	dst = wiremessage.AppendGetMoreNumberToReturn(dst, numToReturn)
	dst = wiremessage.AppendGetMoreCursorID(dst, cursorID)

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, collName, nil
}

func (op Operation) legacyKillCursors(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error {
	wm, startedInfo, _, err := op.createLegacyKillCursorsWiremessage(dst, desc)
	if err != nil {
		return err
	}

	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	// skip startTime because OP_KILL_CURSORS does not return a response
	finishedInfo := finishedInformation{
		cmdName:   "killCursors",
		requestID: startedInfo.requestID,
		connID:    startedInfo.connID,
	}

	err = conn.WriteWireMessage(ctx, wm)
	if err != nil {
		err = Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}}
		if ep, ok := srvr.(ErrorProcessor); ok {
			_ = ep.ProcessError(err, conn)
		}

		finishedInfo.cmdErr = err
		op.publishFinishedEvent(ctx, finishedInfo)
		return err
	}

	ridx, response := bsoncore.AppendDocumentStart(nil)
	response = bsoncore.AppendInt32Element(response, "ok", 1)
	response = bsoncore.AppendArrayElement(response, "cursorsUnknown", startedInfo.cmd.Lookup("cursors").Array())
	response, _ = bsoncore.AppendDocumentEnd(response, ridx)

	finishedInfo.response = response
	op.publishFinishedEvent(ctx, finishedInfo)
	return nil
}

func (op Operation) createLegacyKillCursorsWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) {
	info := startedInformation{
		cmdName:   "killCursors",
		requestID: wiremessage.NextRequestID(),
	}

	var cmdDoc bsoncore.Document
	var cmdIdx int32
	var err error

	cmdIdx, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil {
		return nil, info, "", err
	}
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIdx)
	info.cmd = cmdDoc

	cmdElems, err := cmdDoc.Elements()
	if err != nil {
		return nil, info, "", err
	}

	var collName string
	var cursors bsoncore.Array
	for _, elem := range cmdElems {
		switch elem.Key() {
		case "killCursors":
			collName = elem.Value().StringValue()
		case "cursors":
			cursors = elem.Value().Array()
		}
	}

	var cursorIDs []int64
	if cursors != nil {
		cursorValues, err := cursors.Values()
		if err != nil {
			return nil, info, "", err
		}

		for _, cursorVal := range cursorValues {
			cursorIDs = append(cursorIDs, cursorVal.Int64())
		}
	}

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpKillCursors)
	dst = wiremessage.AppendKillCursorsZero(dst)
	dst = wiremessage.AppendKillCursorsNumberIDs(dst, int32(len(cursorIDs)))
	dst = wiremessage.AppendKillCursorsCursorIDs(dst, cursorIDs)

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, collName, nil
}

func (op Operation) legacyListCollections(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error {
	wm, startedInfo, collName, err := op.createLegacyListCollectionsWiremessage(dst, desc)
	if err != nil {
		return err
	}
	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation{
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	}

	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, firstBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil {
		return finishedInfo.cmdErr
	}

	if op.ProcessResponseFn != nil {
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo{
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		}
		return op.ProcessResponseFn(info)
	}
	return nil
}

func (op Operation) createLegacyListCollectionsWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) {
	info := startedInformation{
		cmdName:   "find",
		requestID: wiremessage.NextRequestID(),
	}

	var cmdDoc bsoncore.Document
	var cmdIdx int32
	var err error

	cmdIdx, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	if cmdDoc, err = op.CommandFn(cmdDoc, desc); err != nil {
		return dst, info, "", err
	}
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIdx)
	info.cmd, err = op.convertCommandToFind(cmdDoc, listCollectionsNamespace)
	if err != nil {
		return nil, info, "", err
	}

	// lookup filter directly instead of calling cmdDoc.Elements() because the only listCollections option is nameOnly,
	// which doesn't apply to legacy servers
	var originalFilter bsoncore.Document
	if filterVal, err := cmdDoc.LookupErr("filter"); err == nil {
		originalFilter = filterVal.Document()
	}

	var optsElems []byte
	filter, err := op.transformListCollectionsFilter(originalFilter)
	if err != nil {
		return dst, info, "", err
	}
	rp, err := op.createReadPref(desc.Server.Kind, desc.Kind, true)
	if err != nil {
		return dst, info, "", err
	}
	if len(rp) > 0 {
		optsElems = bsoncore.AppendDocumentElement(optsElems, "$readPreference", rp)
	}

	var batchSize int32
	if val, ok := cmdDoc.Lookup("cursor", "batchSize").AsInt32OK(); ok {
		batchSize = val
	}

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, op.slaveOK(desc))
	dst = wiremessage.AppendQueryFullCollectionName(dst, op.getFullCollectionName(listCollectionsNamespace))
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, batchSize)
	dst = op.appendLegacyQueryDocument(dst, filter, optsElems)
	// leave out returnFieldsSelector because it is optional

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, listCollectionsNamespace, nil
}

func (op Operation) transformListCollectionsFilter(filter bsoncore.Document) (bsoncore.Document, error) {
	// filter out results containing $ because those represent indexes
	var regexFilter bsoncore.Document
	var ridx int32
	ridx, regexFilter = bsoncore.AppendDocumentStart(regexFilter)
	regexFilter = bsoncore.AppendRegexElement(regexFilter, "name", "^[^$]*$", "")
	regexFilter, _ = bsoncore.AppendDocumentEnd(regexFilter, ridx)

	if len(filter) == 0 {
		return regexFilter, nil
	}

	convertedIdx, convertedFilter := bsoncore.AppendDocumentStart(nil)
	elems, err := filter.Elements()
	if err != nil {
		return nil, err
	}

	for _, elem := range elems {
		if elem.Key() != "name" {
			convertedFilter = append(convertedFilter, elem...)
			continue
		}

		// the name value in a filter for legacy list collections must be a string and has to be prepended
		// with the database name
		nameVal := elem.Value()
		if nameVal.Type != bsontype.String {
			return nil, ErrFilterType
		}
		convertedFilter = bsoncore.AppendStringElement(convertedFilter, "name", op.getFullCollectionName(nameVal.StringValue()))
	}
	convertedFilter, _ = bsoncore.AppendDocumentEnd(convertedFilter, convertedIdx)

	// combine regexFilter and convertedFilter with $and
	var combinedFilter bsoncore.Document
	var cidx, aidx int32
	cidx, combinedFilter = bsoncore.AppendDocumentStart(combinedFilter)
	aidx, combinedFilter = bsoncore.AppendArrayElementStart(combinedFilter, "$and")
	combinedFilter = bsoncore.AppendDocumentElement(combinedFilter, "0", regexFilter)
	combinedFilter = bsoncore.AppendDocumentElement(combinedFilter, "1", convertedFilter)
	combinedFilter, _ = bsoncore.AppendArrayEnd(combinedFilter, aidx)
	combinedFilter, _ = bsoncore.AppendDocumentEnd(combinedFilter, cidx)

	return combinedFilter, nil
}

func (op Operation) legacyListIndexes(ctx context.Context, dst []byte, srvr Server, conn Connection, desc description.SelectedServer) error {
	wm, startedInfo, collName, err := op.createLegacyListIndexesWiremessage(dst, desc)
	if err != nil {
		return err
	}
	startedInfo.connID = conn.ID()
	op.publishStartedEvent(ctx, startedInfo)

	finishedInfo := finishedInformation{
		cmdName:   startedInfo.cmdName,
		requestID: startedInfo.requestID,
		startTime: time.Now(),
		connID:    startedInfo.connID,
	}

	finishedInfo.response, finishedInfo.cmdErr = op.roundTripLegacyCursor(ctx, wm, srvr, conn, collName, firstBatchIdentifier)
	op.publishFinishedEvent(ctx, finishedInfo)

	if finishedInfo.cmdErr != nil {
		return finishedInfo.cmdErr
	}

	if op.ProcessResponseFn != nil {
		// CurrentIndex is always 0 in this mode.
		info := ResponseInfo{
			ServerResponse:        finishedInfo.response,
			Server:                srvr,
			Connection:            conn,
			ConnectionDescription: desc.Server,
		}
		return op.ProcessResponseFn(info)
	}
	return nil
}

func (op Operation) createLegacyListIndexesWiremessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, string, error) {
	info := startedInformation{
		cmdName:   "find",
		requestID: wiremessage.NextRequestID(),
	}

	var cmdDoc bsoncore.Document
	var cmdIndex int32
	var err error

	cmdIndex, cmdDoc = bsoncore.AppendDocumentStart(cmdDoc)
	cmdDoc, err = op.CommandFn(cmdDoc, desc)
	if err != nil {
		return dst, info, "", err
	}
	cmdDoc, _ = bsoncore.AppendDocumentEnd(cmdDoc, cmdIndex)
	info.cmd, err = op.convertCommandToFind(cmdDoc, listIndexesNamespace)
	if err != nil {
		return nil, info, "", err
	}

	cmdElems, err := cmdDoc.Elements()
	if err != nil {
		return nil, info, "", err
	}

	var filterCollName string
	var batchSize int32
	var optsElems []byte // options elements
	for _, elem := range cmdElems {
		switch elem.Key() {
		case "listIndexes":
			filterCollName = elem.Value().StringValue()
		case "cursor":
			// the batchSize option is embedded in a cursor subdocument
			cursorDoc := elem.Value().Document()
			if val, err := cursorDoc.LookupErr("batchSize"); err == nil {
				batchSize = val.Int32()
			}
		case "maxTimeMS":
			optsElems = bsoncore.AppendValueElement(optsElems, "$maxTimeMS", elem.Value())
		}
	}

	// always filter with {ns: db.collection}
	fidx, filter := bsoncore.AppendDocumentStart(nil)
	filter = bsoncore.AppendStringElement(filter, "ns", op.getFullCollectionName(filterCollName))
	filter, _ = bsoncore.AppendDocumentEnd(filter, fidx)

	rp, err := op.createReadPref(desc.Server.Kind, desc.Kind, true)
	if err != nil {
		return dst, info, "", err
	}
	if len(rp) > 0 {
		optsElems = bsoncore.AppendDocumentElement(optsElems, "$readPreference", rp)
	}

	var wmIdx int32
	wmIdx, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, op.slaveOK(desc))
	dst = wiremessage.AppendQueryFullCollectionName(dst, op.getFullCollectionName(listIndexesNamespace))
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, batchSize)
	dst = op.appendLegacyQueryDocument(dst, filter, optsElems)
	// leave out returnFieldsSelector because it is optional

	return bsoncore.UpdateLength(dst, wmIdx, int32(len(dst[wmIdx:]))), info, listIndexesNamespace, nil
}

// convertCommandToFind takes a non-legacy command document for a command that needs to be run as a find on legacy
// servers and converts it to a find command document for APM.
func (op Operation) convertCommandToFind(cmdDoc bsoncore.Document, collName string) (bsoncore.Document, error) {
	cidx, converted := bsoncore.AppendDocumentStart(nil)
	elems, err := cmdDoc.Elements()
	if err != nil {
		return nil, err
	}

	converted = bsoncore.AppendStringElement(converted, "find", collName)
	// skip the first element because that will have the old command name
	for i := 1; i < len(elems); i++ {
		converted = bsoncore.AppendValueElement(converted, elems[i].Key(), elems[i].Value())
	}

	converted, _ = bsoncore.AppendDocumentEnd(converted, cidx)
	return converted, nil
}

// appendLegacyQueryDocument takes a filter and a list of options elements for a legacy find operation, creates
// a query document, and appends it to dst.
func (op Operation) appendLegacyQueryDocument(dst []byte, filter bsoncore.Document, opts []byte) []byte {
	if len(opts) == 0 {
		dst = append(dst, filter...)
		return dst
	}

	// filter must be wrapped in $query if other $-modifiers are used
	var qidx int32
	qidx, dst = bsoncore.AppendDocumentStart(dst)
	dst = bsoncore.AppendDocumentElement(dst, "$query", filter)
	dst = append(dst, opts...)
	dst, _ = bsoncore.AppendDocumentEnd(dst, qidx)
	return dst
}

// roundTripLegacyCursor sends a wiremessage for an operation expecting a cursor result and converts the legacy
// document sequence into a cursor document.
func (op Operation) roundTripLegacyCursor(ctx context.Context, wm []byte, srvr Server, conn Connection, collName, identifier string) (bsoncore.Document, error) {
	wm, err := op.roundTripLegacy(ctx, conn, wm)
	if ep, ok := srvr.(ErrorProcessor); ok {
		_ = ep.ProcessError(err, conn)
	}
	if err != nil {
		return nil, err
	}

	return op.upconvertCursorResponse(wm, identifier, collName)
}

// roundTripLegacy handles writing a wire message and reading the response.
func (op Operation) roundTripLegacy(ctx context.Context, conn Connection, wm []byte) ([]byte, error) {
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil {
		return nil, Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}, Wrapped: err}
	}

	wm, err = conn.ReadWireMessage(ctx, wm[:0])
	if err != nil {
		err = Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}, Wrapped: err}
	}
	return wm, err
}

func (op Operation) upconvertCursorResponse(wm []byte, batchIdentifier string, collName string) (bsoncore.Document, error) {
	reply := op.decodeOpReply(wm, true)
	if reply.err != nil {
		return nil, reply.err
	}

	cursorIdx, cursorDoc := bsoncore.AppendDocumentStart(nil)
	// convert reply documents to BSON array
	var arrIdx int32
	arrIdx, cursorDoc = bsoncore.AppendArrayElementStart(cursorDoc, batchIdentifier)
	for i, doc := range reply.documents {
		cursorDoc = bsoncore.AppendDocumentElement(cursorDoc, strconv.Itoa(i), doc)
	}
	cursorDoc, _ = bsoncore.AppendArrayEnd(cursorDoc, arrIdx)

	cursorDoc = bsoncore.AppendInt64Element(cursorDoc, "id", reply.cursorID)
	cursorDoc = bsoncore.AppendStringElement(cursorDoc, "ns", op.getFullCollectionName(collName))
	cursorDoc, _ = bsoncore.AppendDocumentEnd(cursorDoc, cursorIdx)

	resIdx, resDoc := bsoncore.AppendDocumentStart(nil)
	resDoc = bsoncore.AppendInt32Element(resDoc, "ok", 1)
	resDoc = bsoncore.AppendDocumentElement(resDoc, "cursor", cursorDoc)
	resDoc, _ = bsoncore.AppendDocumentEnd(resDoc, resIdx)

	return resDoc, nil
}
