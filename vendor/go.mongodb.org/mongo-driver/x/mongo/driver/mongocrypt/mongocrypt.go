// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// +build cse

package mongocrypt

// #cgo linux solaris darwin pkg-config: libmongocrypt
// #cgo windows CFLAGS: -I"c:/libmongocrypt/include"
// #cgo windows LDFLAGS: -lmongocrypt -Lc:/libmongocrypt/bin
// #include <mongocrypt.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

// These constants reprsent valid values for KmsProvider.
const (
	AwsProvider   = "aws"
	LocalProvider = "local"
)

// ErrInvalidProvider is returned when an invalid KMS provider is given.
var ErrInvalidProvider = errors.New("invalid KMS provider")

// MongoCrypt represents a mongocrypt_t handle.
type MongoCrypt struct {
	wrapped *C.mongocrypt_t
}

// NewMongoCrypt constructs a new MongoCrypt instance configured using the provided MongoCryptOptions.
func NewMongoCrypt(opts *options.MongoCryptOptions) (*MongoCrypt, error) {
	// create mongocrypt_t handle
	wrapped := C.mongocrypt_new()
	if wrapped == nil {
		return nil, errors.New("could not create new mongocrypt object")
	}
	crypt := &MongoCrypt{
		wrapped: wrapped,
	}

	// set options in mongocrypt
	if err := crypt.setLocalProviderOpts(opts.LocalProviderOpts); err != nil {
		return nil, err
	}
	if err := crypt.setAwsProviderOpts(opts.AwsProviderOpts); err != nil {
		return nil, err
	}
	if err := crypt.setLocalSchemaMap(opts.LocalSchemaMap); err != nil {
		return nil, err
	}

	// initialize handle
	if !C.mongocrypt_init(crypt.wrapped) {
		return nil, crypt.createErrorFromStatus()
	}
	return crypt, nil
}

// CreateEncryptionContext creates a Context to use for encryption.
func (m *MongoCrypt) CreateEncryptionContext(db string, cmd bsoncore.Document) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	cmdBinary := newBinaryFromBytes(cmd)
	defer cmdBinary.close()
	dbStr := C.CString(db)
	defer C.free(unsafe.Pointer(dbStr))

	if ok := C.mongocrypt_ctx_encrypt_init(ctx.wrapped, dbStr, C.int32_t(-1), cmdBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}

// CreateDecryptionContext creates a Context to use for decryption.
func (m *MongoCrypt) CreateDecryptionContext(cmd bsoncore.Document) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	cmdBinary := newBinaryFromBytes(cmd)
	defer cmdBinary.close()

	if ok := C.mongocrypt_ctx_decrypt_init(ctx.wrapped, cmdBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}

// lookupString returns a string for the value corresponding to the given key in the document.
// if the key does not exist or the value is not a string, the empty string is returned.
func lookupString(doc bsoncore.Document, key string) string {
	strVal, _ := doc.Lookup(key).StringValueOK()
	return strVal
}

func setAltName(ctx *Context, altName string) error {
	// create document {"keyAltName": keyAltName}
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendStringElement(doc, "keyAltName", altName)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	keyAltBinary := newBinaryFromBytes(doc)
	defer keyAltBinary.close()

	if ok := C.mongocrypt_ctx_setopt_key_alt_name(ctx.wrapped, keyAltBinary.wrapped); !ok {
		return ctx.createErrorFromStatus()
	}
	return nil
}

// CreateDataKeyContext creates a Context to use for creating a data key.
func (m *MongoCrypt) CreateDataKeyContext(kmsProvider string, opts *options.DataKeyOptions) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	var ok bool
	switch kmsProvider {
	case AwsProvider:
		// set region and key (required fields)
		region := C.CString(lookupString(opts.MasterKey, "region"))
		key := C.CString(lookupString(opts.MasterKey, "key"))
		defer C.free(unsafe.Pointer(region))
		defer C.free(unsafe.Pointer(key))
		ok = bool(C.mongocrypt_ctx_setopt_masterkey_aws(ctx.wrapped, region, -1, key, -1))
		if !ok {
			break
		}

		// set endpoint (not a required field)
		endpoint := lookupString(opts.MasterKey, "endpoint")
		if endpoint == "" {
			break
		}
		endpointCStr := C.CString(endpoint)
		defer C.free(unsafe.Pointer(endpointCStr))
		ok = bool(C.mongocrypt_ctx_setopt_masterkey_aws_endpoint(ctx.wrapped, endpointCStr, -1))
	case LocalProvider:
		ok = bool(C.mongocrypt_ctx_setopt_masterkey_local(ctx.wrapped))
	default:
		return nil, ErrInvalidProvider
	}
	if !ok {
		return nil, ctx.createErrorFromStatus()
	}

	for _, altName := range opts.KeyAltNames {
		if err := setAltName(ctx, altName); err != nil {
			return nil, err
		}
	}

	if ok := C.mongocrypt_ctx_datakey_init(ctx.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}

// CreateExplicitEncryptionContext creates a Context to use for explicit encryption.
func (m *MongoCrypt) CreateExplicitEncryptionContext(doc bsoncore.Document, opts *options.ExplicitEncryptionOptions) (*Context, error) {

	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	if opts.KeyID != nil {
		keyIDBinary := newBinaryFromBytes(opts.KeyID.Data)
		defer keyIDBinary.close()

		if ok := C.mongocrypt_ctx_setopt_key_id(ctx.wrapped, keyIDBinary.wrapped); !ok {
			return nil, ctx.createErrorFromStatus()
		}
	}
	if opts.KeyAltName != nil {
		if err := setAltName(ctx, *opts.KeyAltName); err != nil {
			return nil, err
		}
	}

	algoStr := C.CString(opts.Algorithm)
	defer C.free(unsafe.Pointer(algoStr))
	if ok := C.mongocrypt_ctx_setopt_algorithm(ctx.wrapped, algoStr, -1); !ok {
		return nil, ctx.createErrorFromStatus()
	}

	docBinary := newBinaryFromBytes(doc)
	defer docBinary.close()
	if ok := C.mongocrypt_ctx_explicit_encrypt_init(ctx.wrapped, docBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}

	return ctx, nil
}

// CreateExplicitDecryptionContext creates a Context to use for explicit decryption.
func (m *MongoCrypt) CreateExplicitDecryptionContext(doc bsoncore.Document) (*Context, error) {
	ctx := newContext(C.mongocrypt_ctx_new(m.wrapped))
	if ctx.wrapped == nil {
		return nil, m.createErrorFromStatus()
	}

	docBinary := newBinaryFromBytes(doc)
	defer docBinary.close()

	if ok := C.mongocrypt_ctx_explicit_decrypt_init(ctx.wrapped, docBinary.wrapped); !ok {
		return nil, ctx.createErrorFromStatus()
	}
	return ctx, nil
}

// Close cleans up any resources associated with the given MongoCrypt instance.
func (m *MongoCrypt) Close() {
	C.mongocrypt_destroy(m.wrapped)
}

// setLocalProviderOpts sets options for the local KMS provider in mongocrypt.
func (m *MongoCrypt) setLocalProviderOpts(opts *options.LocalKmsProviderOptions) error {
	if opts == nil {
		return nil
	}

	keyBinary := newBinaryFromBytes(opts.MasterKey)
	defer keyBinary.close()

	if ok := C.mongocrypt_setopt_kms_provider_local(m.wrapped, keyBinary.wrapped); !ok {
		return m.createErrorFromStatus()
	}
	return nil
}

// setAwsProviderOpts sets options for the AWS KMS provider in mongocrypt.
func (m *MongoCrypt) setAwsProviderOpts(opts *options.AwsKmsProviderOptions) error {
	if opts == nil {
		return nil
	}

	// create C strings for function params
	accessKeyID := C.CString(opts.AccessKeyID)
	secretAccessKey := C.CString(opts.SecretAccessKey)
	defer C.free(unsafe.Pointer(accessKeyID))
	defer C.free(unsafe.Pointer(secretAccessKey))

	if ok := C.mongocrypt_setopt_kms_provider_aws(m.wrapped, accessKeyID, C.int32_t(-1), secretAccessKey, C.int32_t(-1)); !ok {
		return m.createErrorFromStatus()
	}
	return nil
}

// setLocalSchemaMap sets the local schema map in mongocrypt.
func (m *MongoCrypt) setLocalSchemaMap(schemaMap map[string]bsoncore.Document) error {
	if len(schemaMap) == 0 {
		return nil
	}

	// convert schema map to BSON document
	midx, mdoc := bsoncore.AppendDocumentStart(nil)
	for key, doc := range schemaMap {
		mdoc = bsoncore.AppendDocumentElement(mdoc, key, doc)
	}
	mdoc, _ = bsoncore.AppendDocumentEnd(mdoc, midx)

	schemaMapBinary := newBinaryFromBytes(mdoc)
	defer schemaMapBinary.close()

	if ok := C.mongocrypt_setopt_schema_map(m.wrapped, schemaMapBinary.wrapped); !ok {
		return m.createErrorFromStatus()
	}
	return nil
}

// createErrorFromStatus creates a new Error based on the status of the MongoCrypt instance.
func (m *MongoCrypt) createErrorFromStatus() error {
	status := C.mongocrypt_status_new()
	defer C.mongocrypt_status_destroy(status)
	C.mongocrypt_status(m.wrapped, status)
	return errorFromStatus(status)
}
