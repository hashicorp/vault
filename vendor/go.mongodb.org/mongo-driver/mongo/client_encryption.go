// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	cryptOpts "go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

// ClientEncryption is used to create data keys and explicitly encrypt and decrypt BSON values.
type ClientEncryption struct {
	crypt          *driver.Crypt
	keyVaultClient *Client
	keyVaultColl   *Collection
}

// NewClientEncryption creates a new ClientEncryption instance configured with the given options.
func NewClientEncryption(keyVaultClient *Client, opts ...*options.ClientEncryptionOptions) (*ClientEncryption, error) {
	if keyVaultClient == nil {
		return nil, errors.New("keyVaultClient must not be nil")
	}

	ce := &ClientEncryption{
		keyVaultClient: keyVaultClient,
	}
	ceo := options.MergeClientEncryptionOptions(opts...)

	// create keyVaultColl
	db, coll := splitNamespace(ceo.KeyVaultNamespace)
	ce.keyVaultColl = ce.keyVaultClient.Database(db).Collection(coll, keyVaultCollOpts)

	kmsProviders, err := transformBsoncoreDocument(bson.DefaultRegistry, ceo.KmsProviders, true, "kmsProviders")
	if err != nil {
		return nil, fmt.Errorf("error creating KMS providers map: %v", err)
	}

	// create Crypt
	kr := keyRetriever{coll: ce.keyVaultColl}
	cir := collInfoRetriever{client: ce.keyVaultClient}
	ce.crypt, err = driver.NewCrypt(&driver.CryptOptions{
		KeyFn:        kr.cryptKeys,
		CollInfoFn:   cir.cryptCollInfo,
		KmsProviders: kmsProviders,
	})
	if err != nil {
		return nil, err
	}

	return ce, nil
}

// CreateDataKey creates a new key document and inserts it into the key vault collection. Returns the _id of the
// created document.
func (ce *ClientEncryption) CreateDataKey(ctx context.Context, kmsProvider string, opts ...*options.DataKeyOptions) (primitive.Binary, error) {
	// translate opts to cryptOpts.DataKeyOptions
	dko := options.MergeDataKeyOptions(opts...)
	co := cryptOpts.DataKey().SetKeyAltNames(dko.KeyAltNames)
	if dko.MasterKey != nil {
		keyDoc, err := transformBsoncoreDocument(ce.keyVaultClient.registry, dko.MasterKey, true, "masterKey")
		if err != nil {
			return primitive.Binary{}, err
		}

		co.SetMasterKey(keyDoc)
	}

	// create data key document
	dataKeyDoc, err := ce.crypt.CreateDataKey(ctx, kmsProvider, co)
	if err != nil {
		return primitive.Binary{}, err
	}

	// insert key into key vault
	_, err = ce.keyVaultColl.InsertOne(ctx, dataKeyDoc)
	if err != nil {
		return primitive.Binary{}, err
	}

	subtype, data := bson.Raw(dataKeyDoc).Lookup("_id").Binary()
	return primitive.Binary{Subtype: subtype, Data: data}, nil
}

// Encrypt encrypts a BSON value with the given key and algorithm. Returns an encrypted value (BSON binary of subtype 6).
func (ce *ClientEncryption) Encrypt(ctx context.Context, val bson.RawValue, opts ...*options.EncryptOptions) (primitive.Binary, error) {
	eo := options.MergeEncryptOptions(opts...)
	transformed := cryptOpts.ExplicitEncryption()
	if eo.KeyID != nil {
		transformed.SetKeyID(*eo.KeyID)
	}
	if eo.KeyAltName != nil {
		transformed.SetKeyAltName(*eo.KeyAltName)
	}
	transformed.SetAlgorithm(eo.Algorithm)

	subtype, data, err := ce.crypt.EncryptExplicit(ctx, bsoncore.Value{Type: val.Type, Data: val.Value}, transformed)
	if err != nil {
		return primitive.Binary{}, err
	}
	return primitive.Binary{Subtype: subtype, Data: data}, nil
}

// Decrypt decrypts an encrypted value (BSON binary of subtype 6) and returns the original BSON value.
func (ce *ClientEncryption) Decrypt(ctx context.Context, val primitive.Binary) (bson.RawValue, error) {
	decrypted, err := ce.crypt.DecryptExplicit(ctx, val.Subtype, val.Data)
	if err != nil {
		return bson.RawValue{}, err
	}

	return bson.RawValue{Type: decrypted.Type, Value: decrypted.Data}, nil
}

// Close cleans up any resources associated with the ClientEncryption instance. This includes disconnecting the
// key-vault Client instance.
func (ce *ClientEncryption) Close(ctx context.Context) error {
	ce.crypt.Close()
	return ce.keyVaultClient.Disconnect(ctx)
}

// splitNamespace takes a namespace in the form "database.collection" and returns (database name, collection name)
func splitNamespace(ns string) (string, string) {
	firstDot := strings.Index(ns, ".")
	if firstDot == -1 {
		return "", ns
	}

	return ns[:firstDot], ns[firstDot+1:]
}
