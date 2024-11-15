// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// SearchIndexView is a type that can be used to create, drop, list and update search indexes on a collection. A SearchIndexView for
// a collection can be created by a call to Collection.SearchIndexes().
type SearchIndexView struct {
	coll *Collection
}

// SearchIndexModel represents a new search index to be created.
type SearchIndexModel struct {
	// A document describing the definition for the search index. It cannot be nil.
	// See https://www.mongodb.com/docs/atlas/atlas-search/create-index/ for reference.
	Definition interface{}

	// The search index options.
	Options *options.SearchIndexesOptions
}

// List executes a listSearchIndexes command and returns a cursor over the search indexes in the collection.
//
// The name parameter specifies the index name. A nil pointer matches all indexes.
//
// The opts parameter can be used to specify options for this operation (see the options.ListSearchIndexesOptions
// documentation).
func (siv SearchIndexView) List(
	ctx context.Context,
	searchIdxOpts *options.SearchIndexesOptions,
	opts ...*options.ListSearchIndexesOptions,
) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	index := bson.D{}
	if searchIdxOpts != nil && searchIdxOpts.Name != nil {
		index = bson.D{{"name", *searchIdxOpts.Name}}
	}

	aggregateOpts := make([]*options.AggregateOptions, len(opts))
	for i, opt := range opts {
		aggregateOpts[i] = opt.AggregateOpts
	}

	return siv.coll.Aggregate(ctx, Pipeline{{{"$listSearchIndexes", index}}}, aggregateOpts...)
}

// CreateOne executes a createSearchIndexes command to create a search index on the collection and returns the name of the new
// search index. See the SearchIndexView.CreateMany documentation for more information and an example.
func (siv SearchIndexView) CreateOne(
	ctx context.Context,
	model SearchIndexModel,
	opts ...*options.CreateSearchIndexesOptions,
) (string, error) {
	names, err := siv.CreateMany(ctx, []SearchIndexModel{model}, opts...)
	if err != nil {
		return "", err
	}

	return names[0], nil
}

// CreateMany executes a createSearchIndexes command to create multiple search indexes on the collection and returns
// the names of the new search indexes.
//
// For each SearchIndexModel in the models parameter, the index name can be specified.
//
// The opts parameter can be used to specify options for this operation (see the options.CreateSearchIndexesOptions
// documentation).
func (siv SearchIndexView) CreateMany(
	ctx context.Context,
	models []SearchIndexModel,
	_ ...*options.CreateSearchIndexesOptions,
) ([]string, error) {
	var indexes bsoncore.Document
	aidx, indexes := bsoncore.AppendArrayStart(indexes)

	for i, model := range models {
		if model.Definition == nil {
			return nil, fmt.Errorf("search index model definition cannot be nil")
		}

		definition, err := marshal(model.Definition, siv.coll.bsonOpts, siv.coll.registry)
		if err != nil {
			return nil, err
		}

		var iidx int32
		iidx, indexes = bsoncore.AppendDocumentElementStart(indexes, strconv.Itoa(i))
		if model.Options != nil && model.Options.Name != nil {
			indexes = bsoncore.AppendStringElement(indexes, "name", *model.Options.Name)
		}
		if model.Options != nil && model.Options.Type != nil {
			indexes = bsoncore.AppendStringElement(indexes, "type", *model.Options.Type)
		}
		indexes = bsoncore.AppendDocumentElement(indexes, "definition", definition)

		indexes, err = bsoncore.AppendDocumentEnd(indexes, iidx)
		if err != nil {
			return nil, err
		}
	}

	indexes, err := bsoncore.AppendArrayEnd(indexes, aidx)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)

	if sess == nil && siv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(siv.coll.client.sessionPool, siv.coll.client.id)
		defer sess.EndSession()
	}

	err = siv.coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	selector := makePinnedSelector(sess, siv.coll.writeSelector)

	op := operation.NewCreateSearchIndexes(indexes).
		Session(sess).CommandMonitor(siv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(siv.coll.client.clock).
		Collection(siv.coll.name).Database(siv.coll.db.name).
		Deployment(siv.coll.client.deployment).ServerAPI(siv.coll.client.serverAPI).
		Timeout(siv.coll.client.timeout).Authenticator(siv.coll.client.authenticator)

	err = op.Execute(ctx)
	if err != nil {
		_, err = processWriteError(err)
		return nil, err
	}

	indexesCreated := op.Result().IndexesCreated
	names := make([]string, 0, len(indexesCreated))
	for _, index := range indexesCreated {
		names = append(names, index.Name)
	}

	return names, nil
}

// DropOne executes a dropSearchIndexes operation to drop a search index on the collection.
//
// The name parameter should be the name of the search index to drop. If the name is "*", ErrMultipleIndexDrop will be returned
// without running the command because doing so would drop all search indexes.
//
// The opts parameter can be used to specify options for this operation (see the options.DropSearchIndexOptions
// documentation).
func (siv SearchIndexView) DropOne(
	ctx context.Context,
	name string,
	_ ...*options.DropSearchIndexOptions,
) error {
	if name == "*" {
		return ErrMultipleIndexDrop
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && siv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(siv.coll.client.sessionPool, siv.coll.client.id)
		defer sess.EndSession()
	}

	err := siv.coll.client.validSession(sess)
	if err != nil {
		return err
	}

	selector := makePinnedSelector(sess, siv.coll.writeSelector)

	op := operation.NewDropSearchIndex(name).
		Session(sess).CommandMonitor(siv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(siv.coll.client.clock).
		Collection(siv.coll.name).Database(siv.coll.db.name).
		Deployment(siv.coll.client.deployment).ServerAPI(siv.coll.client.serverAPI).
		Timeout(siv.coll.client.timeout).Authenticator(siv.coll.client.authenticator)

	err = op.Execute(ctx)
	if de, ok := err.(driver.Error); ok && de.NamespaceNotFound() {
		return nil
	}
	return err
}

// UpdateOne executes a updateSearchIndex operation to update a search index on the collection.
//
// The name parameter should be the name of the search index to update.
//
// The definition parameter is a document describing the definition for the search index. It cannot be nil.
//
// The opts parameter can be used to specify options for this operation (see the options.UpdateSearchIndexOptions
// documentation).
func (siv SearchIndexView) UpdateOne(
	ctx context.Context,
	name string,
	definition interface{},
	_ ...*options.UpdateSearchIndexOptions,
) error {
	if definition == nil {
		return fmt.Errorf("search index definition cannot be nil")
	}

	indexDefinition, err := marshal(definition, siv.coll.bsonOpts, siv.coll.registry)
	if err != nil {
		return err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && siv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(siv.coll.client.sessionPool, siv.coll.client.id)
		defer sess.EndSession()
	}

	err = siv.coll.client.validSession(sess)
	if err != nil {
		return err
	}

	selector := makePinnedSelector(sess, siv.coll.writeSelector)

	op := operation.NewUpdateSearchIndex(name, indexDefinition).
		Session(sess).CommandMonitor(siv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(siv.coll.client.clock).
		Collection(siv.coll.name).Database(siv.coll.db.name).
		Deployment(siv.coll.client.deployment).ServerAPI(siv.coll.client.serverAPI).
		Timeout(siv.coll.client.timeout).Authenticator(siv.coll.client.authenticator)

	return op.Execute(ctx)
}
