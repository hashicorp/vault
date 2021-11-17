// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"time"

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
	CollInfoFn           CollectionInfoFn
	KeyFn                KeyRetrieverFn
	MarkFn               MarkCommandFn
	KmsProviders         bsoncore.Document
	SchemaMap            map[string]bsoncore.Document
	BypassAutoEncryption bool
}

// Crypt consumes the libmongocrypt.MongoCrypt type to iterate the mongocrypt state machine and perform encryption
// and decryption.
type Crypt struct {
	mongoCrypt *mongocrypt.MongoCrypt
	collInfoFn CollectionInfoFn
	keyFn      KeyRetrieverFn
	markFn     MarkCommandFn

	BypassAutoEncryption bool
}

// NewCrypt creates a new Crypt instance configured with the given AutoEncryptionOptions.
func NewCrypt(opts *CryptOptions) (*Crypt, error) {
	c := &Crypt{
		collInfoFn:           opts.CollInfoFn,
		keyFn:                opts.KeyFn,
		markFn:               opts.MarkFn,
		BypassAutoEncryption: opts.BypassAutoEncryption,
	}

	mongocryptOpts := options.MongoCrypt().SetKmsProviders(opts.KmsProviders).SetLocalSchemaMap(opts.SchemaMap)
	mc, err := mongocrypt.NewMongoCrypt(mongocryptOpts)
	if err != nil {
		return nil, err
	}

	c.mongoCrypt = mc
	return c, nil
}

// Encrypt encrypts the given command.
func (c *Crypt) Encrypt(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error) {
	if c.BypassAutoEncryption {
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
func (c *Crypt) Decrypt(ctx context.Context, cmdResponse bsoncore.Document) (bsoncore.Document, error) {
	cryptCtx, err := c.mongoCrypt.CreateDecryptionContext(cmdResponse)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
}

// CreateDataKey creates a data key using the given KMS provider and options.
func (c *Crypt) CreateDataKey(ctx context.Context, kmsProvider string, opts *options.DataKeyOptions) (bsoncore.Document, error) {
	cryptCtx, err := c.mongoCrypt.CreateDataKeyContext(kmsProvider, opts)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
}

// EncryptExplicit encrypts the given value with the given options.
func (c *Crypt) EncryptExplicit(ctx context.Context, val bsoncore.Value, opts *options.ExplicitEncryptionOptions) (byte, []byte, error) {
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

// DecryptExplicit decrypts the given encrypted value.
func (c *Crypt) DecryptExplicit(ctx context.Context, subtype byte, data []byte) (bsoncore.Value, error) {
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
func (c *Crypt) Close() {
	c.mongoCrypt.Close()
}

func (c *Crypt) executeStateMachine(ctx context.Context, cryptCtx *mongocrypt.Context, db string) (bsoncore.Document, error) {
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
			err = c.decryptKeys(ctx, cryptCtx)
		case mongocrypt.Ready:
			return cryptCtx.Finish()
		default:
			return nil, fmt.Errorf("invalid Crypt state: %v", state)
		}
		if err != nil {
			return nil, err
		}
	}
}

func (c *Crypt) collectionInfo(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error {
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

func (c *Crypt) markCommand(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error {
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

func (c *Crypt) retrieveKeys(ctx context.Context, cryptCtx *mongocrypt.Context) error {
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

func (c *Crypt) decryptKeys(ctx context.Context, cryptCtx *mongocrypt.Context) error {
	for {
		kmsCtx := cryptCtx.NextKmsContext()
		if kmsCtx == nil {
			break
		}

		if err := c.decryptKey(ctx, kmsCtx); err != nil {
			return err
		}
	}

	return cryptCtx.FinishKmsContexts()
}

func (c *Crypt) decryptKey(ctx context.Context, kmsCtx *mongocrypt.KmsContext) error {
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

	conn, err := tls.Dial("tcp", addr, &tls.Config{})
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
		if err != nil && err != io.EOF {
			return err
		}

		if err = kmsCtx.FeedResponse(res[:bytesRead]); err != nil {
			return err
		}
	}
}
