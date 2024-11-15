// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

const (
	defaultKmsPort    = 443
	defaultKmsTimeout = 10 * time.Second
)

// CollectionInfoFn is a callback used to retrieve collection information.
type CollectionInfoFn func(ctx context.Context, db string, filter bsoncore.Document) (bsoncore.Document, error)

// KeyRetrieverFn is a callback used to retrieve keys from the key vault.
type KeyRetrieverFn func(ctx context.Context, filter bsoncore.Document) ([]bsoncore.Document, error)

// MarkCommandFn is a callback used to add encryption markings to a command.
type MarkCommandFn func(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error)

// CryptOptions specifies options to configure a Crypt instance.
type CryptOptions struct {
	MongoCrypt           *mongocrypt.MongoCrypt
	CollInfoFn           CollectionInfoFn
	KeyFn                KeyRetrieverFn
	MarkFn               MarkCommandFn
	TLSConfig            map[string]*tls.Config
	BypassAutoEncryption bool
	BypassQueryAnalysis  bool
}

// Crypt is an interface implemented by types that can encrypt and decrypt instances of
// bsoncore.Document.
//
// Users should rely on the driver's crypt type (used by default) for encryption and decryption
// unless they are perfectly confident in another implementation of Crypt.
type Crypt interface {
	// Encrypt encrypts the given command.
	Encrypt(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error)
	// Decrypt decrypts the given command response.
	Decrypt(ctx context.Context, cmdResponse bsoncore.Document) (bsoncore.Document, error)
	// CreateDataKey creates a data key using the given KMS provider and options.
	CreateDataKey(ctx context.Context, kmsProvider string, opts *options.DataKeyOptions) (bsoncore.Document, error)
	// EncryptExplicit encrypts the given value with the given options.
	EncryptExplicit(ctx context.Context, val bsoncore.Value, opts *options.ExplicitEncryptionOptions) (byte, []byte, error)
	// EncryptExplicitExpression encrypts the given expression with the given options.
	EncryptExplicitExpression(ctx context.Context, val bsoncore.Document, opts *options.ExplicitEncryptionOptions) (bsoncore.Document, error)
	// DecryptExplicit decrypts the given encrypted value.
	DecryptExplicit(ctx context.Context, subtype byte, data []byte) (bsoncore.Value, error)
	// Close cleans up any resources associated with the Crypt instance.
	Close()
	// BypassAutoEncryption returns true if auto-encryption should be bypassed.
	BypassAutoEncryption() bool
	// RewrapDataKey attempts to rewrap the document data keys matching the filter, preparing the re-wrapped documents
	// to be returned as a slice of bsoncore.Document.
	RewrapDataKey(ctx context.Context, filter []byte, opts *options.RewrapManyDataKeyOptions) ([]bsoncore.Document, error)
}

// crypt consumes the libmongocrypt.MongoCrypt type to iterate the mongocrypt state machine and perform encryption
// and decryption.
type crypt struct {
	mongoCrypt *mongocrypt.MongoCrypt
	collInfoFn CollectionInfoFn
	keyFn      KeyRetrieverFn
	markFn     MarkCommandFn
	tlsConfig  map[string]*tls.Config

	bypassAutoEncryption bool
}

// NewCrypt creates a new Crypt instance configured with the given AutoEncryptionOptions.
func NewCrypt(opts *CryptOptions) Crypt {
	c := &crypt{
		mongoCrypt:           opts.MongoCrypt,
		collInfoFn:           opts.CollInfoFn,
		keyFn:                opts.KeyFn,
		markFn:               opts.MarkFn,
		tlsConfig:            opts.TLSConfig,
		bypassAutoEncryption: opts.BypassAutoEncryption,
	}
	return c
}

// Encrypt encrypts the given command.
func (c *crypt) Encrypt(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error) {
	if c.bypassAutoEncryption {
		return cmd, nil
	}

	cryptCtx, err := c.mongoCrypt.CreateEncryptionContext(db, cmd)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, db)
}

// Decrypt decrypts the given command response.
func (c *crypt) Decrypt(ctx context.Context, cmdResponse bsoncore.Document) (bsoncore.Document, error) {
	cryptCtx, err := c.mongoCrypt.CreateDecryptionContext(cmdResponse)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
}

// CreateDataKey creates a data key using the given KMS provider and options.
func (c *crypt) CreateDataKey(ctx context.Context, kmsProvider string, opts *options.DataKeyOptions) (bsoncore.Document, error) {
	cryptCtx, err := c.mongoCrypt.CreateDataKeyContext(kmsProvider, opts)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
}

// RewrapDataKey attempts to rewrap the document data keys matching the filter, preparing the re-wrapped documents to
// be returned as a slice of bsoncore.Document.
func (c *crypt) RewrapDataKey(ctx context.Context, filter []byte,
	opts *options.RewrapManyDataKeyOptions) ([]bsoncore.Document, error) {

	cryptCtx, err := c.mongoCrypt.RewrapDataKeyContext(filter, opts)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	rewrappedBSON, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return nil, err
	}
	if rewrappedBSON == nil {
		return nil, nil
	}

	// mongocrypt_ctx_rewrap_many_datakey_init wraps the documents in a BSON of the form { "v": [(BSON document), ...] }
	// where each BSON document in the slice is a document containing a rewrapped datakey.
	rewrappedDocumentBytes, err := rewrappedBSON.LookupErr("v")
	if err != nil {
		return nil, err
	}

	// Parse the resulting BSON as individual documents.
	rewrappedDocsArray, ok := rewrappedDocumentBytes.ArrayOK()
	if !ok {
		return nil, fmt.Errorf("expected results from mongocrypt_ctx_rewrap_many_datakey_init to be an array")
	}

	rewrappedDocumentValues, err := rewrappedDocsArray.Values()
	if err != nil {
		return nil, err
	}

	rewrappedDocuments := []bsoncore.Document{}
	for _, rewrappedDocumentValue := range rewrappedDocumentValues {
		if rewrappedDocumentValue.Type != bsontype.EmbeddedDocument {
			// If a value in the document's array returned by mongocrypt is anything other than an embedded document,
			// then something is wrong and we should terminate the routine.
			return nil, fmt.Errorf("expected value of type %q, got: %q",
				bsontype.EmbeddedDocument.String(),
				rewrappedDocumentValue.Type.String())
		}
		rewrappedDocuments = append(rewrappedDocuments, rewrappedDocumentValue.Document())
	}
	return rewrappedDocuments, nil
}

// EncryptExplicit encrypts the given value with the given options.
func (c *crypt) EncryptExplicit(ctx context.Context, val bsoncore.Value, opts *options.ExplicitEncryptionOptions) (byte, []byte, error) {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendValueElement(doc, "v", val)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitEncryptionContext(doc, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return 0, nil, err
	}

	sub, data := res.Lookup("v").Binary()
	return sub, data, nil
}

// EncryptExplicitExpression encrypts the given expression with the given options.
func (c *crypt) EncryptExplicitExpression(ctx context.Context, expr bsoncore.Document, opts *options.ExplicitEncryptionOptions) (bsoncore.Document, error) {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendDocumentElement(doc, "v", expr)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitEncryptionExpressionContext(doc, opts)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return nil, err
	}

	encryptedExpr := res.Lookup("v").Document()
	return encryptedExpr, nil
}

// DecryptExplicit decrypts the given encrypted value.
func (c *crypt) DecryptExplicit(ctx context.Context, subtype byte, data []byte) (bsoncore.Value, error) {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendBinaryElement(doc, "v", subtype, data)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitDecryptionContext(doc)
	if err != nil {
		return bsoncore.Value{}, err
	}
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return bsoncore.Value{}, err
	}

	return res.Lookup("v"), nil
}

// Close cleans up any resources associated with the Crypt instance.
func (c *crypt) Close() {
	c.mongoCrypt.Close()
}

func (c *crypt) BypassAutoEncryption() bool {
	return c.bypassAutoEncryption
}

func (c *crypt) executeStateMachine(ctx context.Context, cryptCtx *mongocrypt.Context, db string) (bsoncore.Document, error) {
	var err error
	for {
		state := cryptCtx.State()
		switch state {
		case mongocrypt.NeedMongoCollInfo:
			err = c.collectionInfo(ctx, cryptCtx, db)
		case mongocrypt.NeedMongoMarkings:
			err = c.markCommand(ctx, cryptCtx, db)
		case mongocrypt.NeedMongoKeys:
			err = c.retrieveKeys(ctx, cryptCtx)
		case mongocrypt.NeedKms:
			err = c.decryptKeys(cryptCtx)
		case mongocrypt.Ready:
			return cryptCtx.Finish()
		case mongocrypt.Done:
			return nil, nil
		case mongocrypt.NeedKmsCredentials:
			err = c.provideKmsProviders(ctx, cryptCtx)
		default:
			return nil, fmt.Errorf("invalid Crypt state: %v", state)
		}
		if err != nil {
			return nil, err
		}
	}
}

func (c *crypt) collectionInfo(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error {
	op, err := cryptCtx.NextOperation()
	if err != nil {
		return err
	}

	collInfo, err := c.collInfoFn(ctx, db, op)
	if err != nil {
		return err
	}
	if collInfo != nil {
		if err = cryptCtx.AddOperationResult(collInfo); err != nil {
			return err
		}
	}

	return cryptCtx.CompleteOperation()
}

func (c *crypt) markCommand(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error {
	op, err := cryptCtx.NextOperation()
	if err != nil {
		return err
	}

	markedCmd, err := c.markFn(ctx, db, op)
	if err != nil {
		return err
	}
	if err = cryptCtx.AddOperationResult(markedCmd); err != nil {
		return err
	}

	return cryptCtx.CompleteOperation()
}

func (c *crypt) retrieveKeys(ctx context.Context, cryptCtx *mongocrypt.Context) error {
	op, err := cryptCtx.NextOperation()
	if err != nil {
		return err
	}

	keys, err := c.keyFn(ctx, op)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err = cryptCtx.AddOperationResult(key); err != nil {
			return err
		}
	}

	return cryptCtx.CompleteOperation()
}

func (c *crypt) decryptKeys(cryptCtx *mongocrypt.Context) error {
	for {
		kmsCtx := cryptCtx.NextKmsContext()
		if kmsCtx == nil {
			break
		}

		if err := c.decryptKey(kmsCtx); err != nil {
			return err
		}
	}

	return cryptCtx.FinishKmsContexts()
}

func (c *crypt) decryptKey(kmsCtx *mongocrypt.KmsContext) error {
	host, err := kmsCtx.HostName()
	if err != nil {
		return err
	}
	msg, err := kmsCtx.Message()
	if err != nil {
		return err
	}

	// add a port to the address if it's not already present
	addr := host
	if idx := strings.IndexByte(host, ':'); idx == -1 {
		addr = fmt.Sprintf("%s:%d", host, defaultKmsPort)
	}

	kmsProvider := kmsCtx.KMSProvider()
	tlsCfg := c.tlsConfig[kmsProvider]
	if tlsCfg == nil {
		tlsCfg = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	if err = conn.SetWriteDeadline(time.Now().Add(defaultKmsTimeout)); err != nil {
		return err
	}
	if _, err = conn.Write(msg); err != nil {
		return err
	}

	for {
		bytesNeeded := kmsCtx.BytesNeeded()
		if bytesNeeded == 0 {
			return nil
		}

		res := make([]byte, bytesNeeded)
		bytesRead, err := conn.Read(res)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		if err = kmsCtx.FeedResponse(res[:bytesRead]); err != nil {
			return err
		}
	}
}

func (c *crypt) provideKmsProviders(ctx context.Context, cryptCtx *mongocrypt.Context) error {
	kmsProviders, err := c.mongoCrypt.GetKmsProviders(ctx)
	if err != nil {
		return err
	}
	return cryptCtx.ProvideKmsProviders(kmsProviders)
}
