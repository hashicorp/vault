// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// keyRetriever gets keys from the key vault collection.
type keyRetriever struct {
	coll *Collection
}

func (kr *keyRetriever) cryptKeys(ctx context.Context, filter bsoncore.Document) ([]bsoncore.Document, error) {
	// Remove the explicit session from the context if one is set.
	// The explicit session may be from a different client.
	ctx = NewSessionContext(ctx, nil)
	cursor, err := kr.coll.Find(ctx, filter)
	if err != nil {
		return nil, EncryptionKeyVaultError{Wrapped: err}
	}
	defer cursor.Close(ctx)

	var results []bsoncore.Document
	for cursor.Next(ctx) {
		cur := make([]byte, len(cursor.Current))
		copy(cur, cursor.Current)
		results = append(results, cur)
	}
	if err = cursor.Err(); err != nil {
		return nil, EncryptionKeyVaultError{Wrapped: err}
	}

	return results, nil
}

// collInfoRetriever gets info for collections from a database.
type collInfoRetriever struct {
	client *Client
}

func (cir *collInfoRetriever) cryptCollInfo(ctx context.Context, db string, filter bsoncore.Document) (bsoncore.Document, error) {
	// Remove the explicit session from the context if one is set.
	// The explicit session may be from a different client.
	ctx = NewSessionContext(ctx, nil)
	cursor, err := cir.client.Database(db).ListCollections(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, cursor.Err()
	}

	res := make([]byte, len(cursor.Current))
	copy(res, cursor.Current)
	return res, nil
}
