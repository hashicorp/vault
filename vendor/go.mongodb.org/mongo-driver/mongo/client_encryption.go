// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	mcopts "go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

// ClientEncryption is used to create data keys and explicitly encrypt and decrypt BSON values.
type ClientEncryption struct {
	crypt          driver.Crypt
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

	kmsProviders, err := marshal(ceo.KmsProviders, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS providers map: %w", err)
	}

	mc, err := mongocrypt.NewMongoCrypt(mcopts.MongoCrypt().
		SetKmsProviders(kmsProviders).
		// Explicitly disable loading the crypt_shared library for the Crypt used for
		// ClientEncryption because it's only needed for AutoEncryption and we don't expect users to
		// have the crypt_shared library installed if they're using ClientEncryption.
		SetCryptSharedLibDisabled(true).
		SetHTTPClient(ceo.HTTPClient))
	if err != nil {
		return nil, err
	}

	// create Crypt
	kr := keyRetriever{coll: ce.keyVaultColl}
	cir := collInfoRetriever{client: ce.keyVaultClient}
	ce.crypt = driver.NewCrypt(&driver.CryptOptions{
		MongoCrypt: mc,
		KeyFn:      kr.cryptKeys,
		CollInfoFn: cir.cryptCollInfo,
		TLSConfig:  ceo.TLSConfig,
	})

	return ce, nil
}

// CreateEncryptedCollection creates a new collection for Queryable Encryption with the help of automatic generation of new encryption data keys for null keyIds.
// It returns the created collection and the encrypted fields document used to create it.
func (ce *ClientEncryption) CreateEncryptedCollection(ctx context.Context,
	db *Database, coll string, createOpts *options.CreateCollectionOptions,
	kmsProvider string, masterKey interface{}) (*Collection, bson.M, error) {
	if createOpts == nil {
		return nil, nil, errors.New("nil CreateCollectionOptions")
	}
	ef := createOpts.EncryptedFields
	if ef == nil {
		return nil, nil, errors.New("no EncryptedFields defined for the collection")
	}

	efBSON, err := marshal(ef, db.bsonOpts, db.registry)
	if err != nil {
		return nil, nil, err
	}
	r := bsonrw.NewBSONDocumentReader(efBSON)
	dec, err := bson.NewDecoder(r)
	if err != nil {
		return nil, nil, err
	}
	var m bson.M
	err = dec.Decode(&m)
	if err != nil {
		return nil, nil, err
	}

	if v, ok := m["fields"]; ok {
		if fields, ok := v.(bson.A); ok {
			for _, field := range fields {
				if f, ok := field.(bson.M); !ok {
					continue
				} else if v, ok := f["keyId"]; ok && v == nil {
					dkOpts := options.DataKey()
					if masterKey != nil {
						dkOpts.SetMasterKey(masterKey)
					}
					keyid, err := ce.CreateDataKey(ctx, kmsProvider, dkOpts)
					if err != nil {
						createOpts.EncryptedFields = m
						return nil, m, err
					}
					f["keyId"] = keyid
				}
			}
			createOpts.EncryptedFields = m
		}
	}
	err = db.CreateCollection(ctx, coll, createOpts)
	if err != nil {
		return nil, m, err
	}
	return db.Collection(coll), m, nil
}

// AddKeyAltName adds a keyAltName to the keyAltNames array of the key document in the key vault collection with the
// given UUID (BSON binary subtype 0x04). Returns the previous version of the key document.
func (ce *ClientEncryption) AddKeyAltName(ctx context.Context, id primitive.Binary, keyAltName string) *SingleResult {
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	keyAltNameDoc := bsoncore.NewDocumentBuilder().AppendString("keyAltNames", keyAltName).Build()
	update := bsoncore.NewDocumentBuilder().AppendDocument("$addToSet", keyAltNameDoc).Build()
	return ce.keyVaultColl.FindOneAndUpdate(ctx, filter, update)
}

// CreateDataKey creates a new key document and inserts into the key vault collection. Returns the _id of the created
// document as a UUID (BSON binary subtype 0x04).
func (ce *ClientEncryption) CreateDataKey(ctx context.Context, kmsProvider string,
	opts ...*options.DataKeyOptions) (primitive.Binary, error) {

	// translate opts to mcopts.DataKeyOptions
	dko := options.MergeDataKeyOptions(opts...)
	co := mcopts.DataKey().SetKeyAltNames(dko.KeyAltNames)
	if dko.MasterKey != nil {
		keyDoc, err := marshal(
			dko.MasterKey,
			ce.keyVaultClient.bsonOpts,
			ce.keyVaultClient.registry)
		if err != nil {
			return primitive.Binary{}, err
		}
		co.SetMasterKey(keyDoc)
	}
	if dko.KeyMaterial != nil {
		co.SetKeyMaterial(dko.KeyMaterial)
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

// transformExplicitEncryptionOptions creates explicit encryption options to be passed to libmongocrypt.
func transformExplicitEncryptionOptions(opts ...*options.EncryptOptions) *mcopts.ExplicitEncryptionOptions {
	eo := options.MergeEncryptOptions(opts...)
	transformed := mcopts.ExplicitEncryption()
	if eo.KeyID != nil {
		transformed.SetKeyID(*eo.KeyID)
	}
	if eo.KeyAltName != nil {
		transformed.SetKeyAltName(*eo.KeyAltName)
	}
	transformed.SetAlgorithm(eo.Algorithm)
	transformed.SetQueryType(eo.QueryType)

	if eo.ContentionFactor != nil {
		transformed.SetContentionFactor(*eo.ContentionFactor)
	}

	if eo.RangeOptions != nil {
		var transformedRange mcopts.ExplicitRangeOptions
		if eo.RangeOptions.Min != nil {
			transformedRange.Min = &bsoncore.Value{Type: eo.RangeOptions.Min.Type, Data: eo.RangeOptions.Min.Value}
		}
		if eo.RangeOptions.Max != nil {
			transformedRange.Max = &bsoncore.Value{Type: eo.RangeOptions.Max.Type, Data: eo.RangeOptions.Max.Value}
		}
		if eo.RangeOptions.Precision != nil {
			transformedRange.Precision = eo.RangeOptions.Precision
		}
		if eo.RangeOptions.Sparsity != nil {
			transformedRange.Sparsity = eo.RangeOptions.Sparsity
		}
		if eo.RangeOptions.TrimFactor != nil {
			transformedRange.TrimFactor = eo.RangeOptions.TrimFactor
		}
		transformed.SetRangeOptions(transformedRange)
	}
	return transformed
}

// Encrypt encrypts a BSON value with the given key and algorithm. Returns an encrypted value (BSON binary of subtype 6).
func (ce *ClientEncryption) Encrypt(ctx context.Context, val bson.RawValue,
	opts ...*options.EncryptOptions) (primitive.Binary, error) {

	transformed := transformExplicitEncryptionOptions(opts...)
	subtype, data, err := ce.crypt.EncryptExplicit(ctx, bsoncore.Value{Type: val.Type, Data: val.Value}, transformed)
	if err != nil {
		return primitive.Binary{}, err
	}
	return primitive.Binary{Subtype: subtype, Data: data}, nil
}

// EncryptExpression encrypts an expression to query a range index.
// On success, `result` is populated with the resulting BSON document.
// `expr` is expected to be a BSON document of one of the following forms:
// 1. A Match Expression of this form:
// {$and: [{<field>: {$gt: <value1>}}, {<field>: {$lt: <value2> }}]}
// 2. An Aggregate Expression of this form:
// {$and: [{$gt: [<fieldpath>, <value1>]}, {$lt: [<fieldpath>, <value2>]}]
// $gt may also be $gte. $lt may also be $lte.
// Only supported for queryType "range"
func (ce *ClientEncryption) EncryptExpression(ctx context.Context, expr interface{}, result interface{}, opts ...*options.EncryptOptions) error {
	transformed := transformExplicitEncryptionOptions(opts...)

	exprDoc, err := marshal(expr, nil, nil)
	if err != nil {
		return err
	}

	encryptedExprDoc, err := ce.crypt.EncryptExplicitExpression(ctx, exprDoc, transformed)
	if err != nil {
		return err
	}
	if raw, ok := result.(*bson.Raw); ok {
		// Avoid the cost of Unmarshal.
		*raw = bson.Raw(encryptedExprDoc)
		return nil
	}
	err = bson.Unmarshal([]byte(encryptedExprDoc), result)
	if err != nil {
		return err
	}
	return nil
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

// DeleteKey removes the key document with the given UUID (BSON binary subtype 0x04) from the key vault collection.
// Returns the result of the internal deleteOne() operation on the key vault collection.
func (ce *ClientEncryption) DeleteKey(ctx context.Context, id primitive.Binary) (*DeleteResult, error) {
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	return ce.keyVaultColl.DeleteOne(ctx, filter)
}

// GetKeyByAltName returns a key document in the key vault collection with the given keyAltName.
func (ce *ClientEncryption) GetKeyByAltName(ctx context.Context, keyAltName string) *SingleResult {
	filter := bsoncore.NewDocumentBuilder().AppendString("keyAltNames", keyAltName).Build()
	return ce.keyVaultColl.FindOne(ctx, filter)
}

// GetKey finds a single key document with the given UUID (BSON binary subtype 0x04). Returns the result of the
// internal find() operation on the key vault collection.
func (ce *ClientEncryption) GetKey(ctx context.Context, id primitive.Binary) *SingleResult {
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	return ce.keyVaultColl.FindOne(ctx, filter)
}

// GetKeys finds all documents in the key vault collection. Returns the result of the internal find() operation on the
// key vault collection.
func (ce *ClientEncryption) GetKeys(ctx context.Context) (*Cursor, error) {
	return ce.keyVaultColl.Find(ctx, bson.D{})
}

// RemoveKeyAltName removes a keyAltName from the keyAltNames array of the key document in the key vault collection with
// the given UUID (BSON binary subtype 0x04). Returns the previous version of the key document.
func (ce *ClientEncryption) RemoveKeyAltName(ctx context.Context, id primitive.Binary, keyAltName string) *SingleResult {
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	update := bson.A{bson.D{{"$set", bson.D{{"keyAltNames", bson.D{{"$cond", bson.A{bson.D{{"$eq",
		bson.A{"$keyAltNames", bson.A{keyAltName}}}}, "$$REMOVE", bson.D{{"$filter",
		bson.D{{"input", "$keyAltNames"}, {"cond", bson.D{{"$ne", bson.A{"$$this", keyAltName}}}}}}}}}}}}}}}
	return ce.keyVaultColl.FindOneAndUpdate(ctx, filter, update)
}

// setRewrapManyDataKeyWriteModels will prepare the WriteModel slice for a bulk updating rewrapped documents.
func setRewrapManyDataKeyWriteModels(rewrappedDocuments []bsoncore.Document, writeModels *[]WriteModel) error {
	const idKey = "_id"
	const keyMaterial = "keyMaterial"
	const masterKey = "masterKey"

	if writeModels == nil {
		return fmt.Errorf("writeModels pointer not set for location referenced")
	}

	// Append a slice of WriteModel with the update document per each rewrappedDoc _id filter.
	for _, rewrappedDocument := range rewrappedDocuments {
		// Prepare the new master key for update.
		masterKeyValue, err := rewrappedDocument.LookupErr(masterKey)
		if err != nil {
			return err
		}
		masterKeyDoc := masterKeyValue.Document()

		// Prepare the new material key for update.
		keyMaterialValue, err := rewrappedDocument.LookupErr(keyMaterial)
		if err != nil {
			return err
		}
		keyMaterialSubtype, keyMaterialData := keyMaterialValue.Binary()
		keyMaterialBinary := primitive.Binary{Subtype: keyMaterialSubtype, Data: keyMaterialData}

		// Prepare the _id filter for documents to update.
		id, err := rewrappedDocument.LookupErr(idKey)
		if err != nil {
			return err
		}

		idSubtype, idData, ok := id.BinaryOK()
		if !ok {
			return fmt.Errorf("expected to assert %q as binary, got type %T", idKey, id)
		}
		binaryID := primitive.Binary{Subtype: idSubtype, Data: idData}

		// Append the mutable document to the slice for bulk update.
		*writeModels = append(*writeModels, NewUpdateOneModel().
			SetFilter(bson.D{{idKey, binaryID}}).
			SetUpdate(
				bson.D{
					{"$set", bson.D{{keyMaterial, keyMaterialBinary}, {masterKey, masterKeyDoc}}},
					{"$currentDate", bson.D{{"updateDate", true}}},
				},
			))
	}
	return nil
}

// RewrapManyDataKey decrypts and encrypts all matching data keys with a possibly new masterKey value. For all
// matching documents, this method will overwrite the "masterKey", "updateDate", and "keyMaterial". On error, some
// matching data keys may have been rewrapped.
// libmongocrypt 1.5.2 is required. An error is returned if the detected version of libmongocrypt is less than 1.5.2.
func (ce *ClientEncryption) RewrapManyDataKey(ctx context.Context, filter interface{},
	opts ...*options.RewrapManyDataKeyOptions) (*RewrapManyDataKeyResult, error) {

	// libmongocrypt versions 1.5.0 and 1.5.1 have a severe bug in RewrapManyDataKey.
	// Check if the version string starts with 1.5.0 or 1.5.1. This accounts for pre-release versions, like 1.5.0-rc0.
	libmongocryptVersion := mongocrypt.Version()
	if strings.HasPrefix(libmongocryptVersion, "1.5.0") || strings.HasPrefix(libmongocryptVersion, "1.5.1") {
		return nil, fmt.Errorf("RewrapManyDataKey requires libmongocrypt 1.5.2 or newer. Detected version: %v", libmongocryptVersion)
	}

	rmdko := options.MergeRewrapManyDataKeyOptions(opts...)
	if ctx == nil {
		ctx = context.Background()
	}

	// Transfer rmdko options to /x/ package options to publish the mongocrypt feed.
	co := mcopts.RewrapManyDataKey()
	if rmdko.MasterKey != nil {
		keyDoc, err := marshal(
			rmdko.MasterKey,
			ce.keyVaultClient.bsonOpts,
			ce.keyVaultClient.registry)
		if err != nil {
			return nil, err
		}
		co.SetMasterKey(keyDoc)
	}
	if rmdko.Provider != nil {
		co.SetProvider(*rmdko.Provider)
	}

	// Prepare the filters and rewrap the data key using mongocrypt.
	filterdoc, err := marshal(filter, ce.keyVaultClient.bsonOpts, ce.keyVaultClient.registry)
	if err != nil {
		return nil, err
	}

	rewrappedDocuments, err := ce.crypt.RewrapDataKey(ctx, filterdoc, co)
	if err != nil {
		return nil, err
	}
	if len(rewrappedDocuments) == 0 {
		// If there are no documents to rewrap, then do nothing.
		return new(RewrapManyDataKeyResult), nil
	}

	// Prepare the WriteModel slice for bulk updating the rewrapped data keys.
	models := []WriteModel{}
	if err := setRewrapManyDataKeyWriteModels(rewrappedDocuments, &models); err != nil {
		return nil, err
	}

	bulkWriteResults, err := ce.keyVaultColl.BulkWrite(ctx, models)
	return &RewrapManyDataKeyResult{BulkWriteResult: bulkWriteResults}, err
}

// splitNamespace takes a namespace in the form "database.collection" and returns (database name, collection name)
func splitNamespace(ns string) (string, string) {
	firstDot := strings.Index(ns, ".")
	if firstDot == -1 {
		return "", ns
	}

	return ns[:firstDot], ns[firstDot+1:]
}
