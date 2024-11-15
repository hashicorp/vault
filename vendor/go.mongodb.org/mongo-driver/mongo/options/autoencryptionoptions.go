// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"crypto/tls"
	"net/http"

	"go.mongodb.org/mongo-driver/internal/httputil"
)

// AutoEncryptionOptions represents options used to configure auto encryption/decryption behavior for a mongo.Client
// instance.
//
// Automatic encryption is an enterprise only feature that only applies to operations on a collection. Automatic
// encryption is not supported for operations on a database or view, and operations that are not bypassed will result
// in error. Too bypass automatic encryption for all operations, set BypassAutoEncryption=true.
//
// Auto encryption requires the authenticated user to have the listCollections privilege action.
//
// If automatic encryption fails on an operation, use a MongoClient configured with bypassAutoEncryption=true and use
// ClientEncryption.encrypt() to manually encrypt values.
//
// Enabling Client Side Encryption reduces the maximum document and message size (using a maxBsonObjectSize of 2MiB and
// maxMessageSizeBytes of 6MB) and may have a negative performance impact.
type AutoEncryptionOptions struct {
	KeyVaultClientOptions *ClientOptions
	KeyVaultNamespace     string
	KmsProviders          map[string]map[string]interface{}
	SchemaMap             map[string]interface{}
	BypassAutoEncryption  *bool
	ExtraOptions          map[string]interface{}
	TLSConfig             map[string]*tls.Config
	HTTPClient            *http.Client
	EncryptedFieldsMap    map[string]interface{}
	BypassQueryAnalysis   *bool
}

// AutoEncryption creates a new AutoEncryptionOptions configured with default values.
func AutoEncryption() *AutoEncryptionOptions {
	return &AutoEncryptionOptions{
		HTTPClient: httputil.DefaultHTTPClient,
	}
}

// SetKeyVaultClientOptions specifies options for the client used to communicate with the key vault collection.
//
// If this is set, it is used to create an internal mongo.Client.
// Otherwise, if the target mongo.Client being configured has an unlimited connection pool size (i.e. maxPoolSize=0),
// it is reused to interact with the key vault collection.
// Otherwise, if the target mongo.Client has a limited connection pool size, a separate internal mongo.Client is used
// (and created if necessary). The internal mongo.Client may be shared during automatic encryption (if
// BypassAutomaticEncryption is false). The internal mongo.Client is configured with the same options as the target
// mongo.Client except minPoolSize is set to 0 and AutoEncryptionOptions is omitted.
func (a *AutoEncryptionOptions) SetKeyVaultClientOptions(opts *ClientOptions) *AutoEncryptionOptions {
	a.KeyVaultClientOptions = opts
	return a
}

// SetKeyVaultNamespace specifies the namespace of the key vault collection. This is required.
func (a *AutoEncryptionOptions) SetKeyVaultNamespace(ns string) *AutoEncryptionOptions {
	a.KeyVaultNamespace = ns
	return a
}

// SetKmsProviders specifies options for KMS providers. This is required.
func (a *AutoEncryptionOptions) SetKmsProviders(providers map[string]map[string]interface{}) *AutoEncryptionOptions {
	a.KmsProviders = providers
	return a
}

// SetSchemaMap specifies a map from namespace to local schema document. Schemas supplied in the schemaMap only apply
// to configuring automatic encryption for client side encryption. Other validation rules in the JSON schema will not
// be enforced by the driver and will result in an error.
//
// Supplying a schemaMap provides more security than relying on JSON Schemas obtained from the server. It protects
// against a malicious server advertising a false JSON Schema, which could trick the client into sending unencrypted
// data that should be encrypted.
func (a *AutoEncryptionOptions) SetSchemaMap(schemaMap map[string]interface{}) *AutoEncryptionOptions {
	a.SchemaMap = schemaMap
	return a
}

// SetBypassAutoEncryption specifies whether or not auto encryption should be done.
//
// If this is unset or false and target mongo.Client being configured has an unlimited connection pool size
// (i.e. maxPoolSize=0), it is reused in the process of auto encryption.
// Otherwise, if the target mongo.Client has a limited connection pool size, a separate internal mongo.Client is used
// (and created if necessary). The internal mongo.Client may be shared for key vault operations (if KeyVaultClient is
// unset). The internal mongo.Client is configured with the same options as the target mongo.Client except minPoolSize
// is set to 0 and AutoEncryptionOptions is omitted.
func (a *AutoEncryptionOptions) SetBypassAutoEncryption(bypass bool) *AutoEncryptionOptions {
	a.BypassAutoEncryption = &bypass
	return a
}

// SetExtraOptions specifies a map of options to configure the mongocryptd process or mongo_crypt shared library.
//
// # Supported Extra Options
//
// "mongocryptdURI" - The mongocryptd URI. Allows setting a custom URI used to communicate with the
// mongocryptd process. The default is "mongodb://localhost:27020", which works with the default
// mongocryptd process spawned by the Client. Must be a string.
//
// "mongocryptdBypassSpawn" - If set to true, the Client will not attempt to spawn a mongocryptd
// process. Must be a bool.
//
// "mongocryptdSpawnPath" - The path used when spawning mongocryptd.
// Defaults to empty string and spawns mongocryptd from system path. Must be a string.
//
// "mongocryptdSpawnArgs" - Command line arguments passed when spawning mongocryptd.
// Defaults to ["--idleShutdownTimeoutSecs=60"]. Must be an array of strings.
//
// "cryptSharedLibRequired" - If set to true, Client creation will return an error if the
// crypt_shared library is not loaded. If unset or set to false, Client creation will not return an
// error if the crypt_shared library is not loaded. The default is unset. Must be a bool.
//
// "cryptSharedLibPath" - The crypt_shared library override path. This must be the path to the
// crypt_shared dynamic library file (for example, a .so, .dll, or .dylib file), not the directory
// that contains it. If the override path is a relative path, it will be resolved relative to the
// working directory of the process. If the override path is a relative path and the first path
// component is the literal string "$ORIGIN", the "$ORIGIN" component will be replaced by the
// absolute path to the directory containing the linked libmongocrypt library. Setting an override
// path disables the default system library search path. If an override path is specified but the
// crypt_shared library cannot be loaded, Client creation will return an error. Must be a string.
func (a *AutoEncryptionOptions) SetExtraOptions(extraOpts map[string]interface{}) *AutoEncryptionOptions {
	a.ExtraOptions = extraOpts
	return a
}

// SetTLSConfig specifies tls.Config instances for each KMS provider to use to configure TLS on all connections created
// to the KMS provider.
//
// This should only be used to set custom TLS configurations. By default, the connection will use an empty tls.Config{} with MinVersion set to tls.VersionTLS12.
func (a *AutoEncryptionOptions) SetTLSConfig(tlsOpts map[string]*tls.Config) *AutoEncryptionOptions {
	tlsConfigs := make(map[string]*tls.Config)
	for provider, config := range tlsOpts {
		// use TLS min version 1.2 to enforce more secure hash algorithms and advanced cipher suites
		if config.MinVersion == 0 {
			config.MinVersion = tls.VersionTLS12
		}
		tlsConfigs[provider] = config
	}
	a.TLSConfig = tlsConfigs
	return a
}

// SetEncryptedFieldsMap specifies a map from namespace to local EncryptedFieldsMap document.
// EncryptedFieldsMap is used for Queryable Encryption.
func (a *AutoEncryptionOptions) SetEncryptedFieldsMap(ef map[string]interface{}) *AutoEncryptionOptions {
	a.EncryptedFieldsMap = ef
	return a
}

// SetBypassQueryAnalysis specifies whether or not query analysis should be used for automatic encryption.
// Use this option when using explicit encryption with Queryable Encryption.
func (a *AutoEncryptionOptions) SetBypassQueryAnalysis(bypass bool) *AutoEncryptionOptions {
	a.BypassQueryAnalysis = &bypass
	return a
}

// MergeAutoEncryptionOptions combines the argued AutoEncryptionOptions in a last-one wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeAutoEncryptionOptions(opts ...*AutoEncryptionOptions) *AutoEncryptionOptions {
	aeo := AutoEncryption()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.KeyVaultClientOptions != nil {
			aeo.KeyVaultClientOptions = opt.KeyVaultClientOptions
		}
		if opt.KeyVaultNamespace != "" {
			aeo.KeyVaultNamespace = opt.KeyVaultNamespace
		}
		if opt.KmsProviders != nil {
			aeo.KmsProviders = opt.KmsProviders
		}
		if opt.SchemaMap != nil {
			aeo.SchemaMap = opt.SchemaMap
		}
		if opt.BypassAutoEncryption != nil {
			aeo.BypassAutoEncryption = opt.BypassAutoEncryption
		}
		if opt.ExtraOptions != nil {
			aeo.ExtraOptions = opt.ExtraOptions
		}
		if opt.TLSConfig != nil {
			aeo.TLSConfig = opt.TLSConfig
		}
		if opt.EncryptedFieldsMap != nil {
			aeo.EncryptedFieldsMap = opt.EncryptedFieldsMap
		}
		if opt.BypassQueryAnalysis != nil {
			aeo.BypassQueryAnalysis = opt.BypassQueryAnalysis
		}
		if opt.HTTPClient != nil {
			aeo.HTTPClient = opt.HTTPClient
		}
	}

	return aeo
}
