# Go-KMS-Wrapping - Go library for encrypting values through various KMS providers

[![Go Reference](https://godoc.org/github.com/hashicorp/go-kms-wrapping/v2?status.svg)](https://godoc.org/github.com/hashicorp/go-kms-wrapping/v2)

*NOTE*: This is version 2 of the library. The `v0` branch contains version 0,
which may be needed for legacy applications or while transitioning to version 2.

Go-KMS-Wrapping is a library that can be used to encrypt things through various
KMS providers -- public clouds, Vault's Transit plugin, etc. It is similar in
concept to various other cryptosystems (like NaCl) but focuses on using third
party KMSes. This library is the underpinning of Vault's auto-unseal
functionality, and should be ready to use for many other applications.

For KMS providers that do not support encrypting arbitrarily large values, the
library will generate an envelope data encryption key (DEK), encrypt the value
with it using an authenticated cipher, and use the KMS to encrypt the DEK.

The key being used by a given implementation can change; the library stores
information about which key was actually used to encrypt a given value as part
of the returned data, and this key will be used for decryption. By extension,
this means that users should be careful not to delete keys in KMS systems
simply because they're not configured to be used by this library _currently_,
as they may have been used for past encryption operations.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Go-KMS-Wrapping - Go library for encrypting values through various KMS providers](#go-kms-wrapping---go-library-for-encrypting-values-through-various-kms-providers)
  - [Features](#features)
  - [Extras](#extras)
  - [Installation](#installation)
  - [Overview](#overview)
  - [Usage](#usage)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Features

  * Supports many KMSes:
  * * AEAD using AES-GCM and a provided key
  * * Alibaba Cloud KMS (uses envelopes)
  * * AWS KMS (uses envelopes)
  * * Azure KeyVault (uses envelopes)
  * * GCP CKMS (uses envelopes)
  * * Huawei Cloud KMS (uses envelopes)
  * * OCI KMS (uses envelopes)
  * * Tencent Cloud KMS (uses envelopes)
  * * Vault Transit mount
  * Transparently supports multiple decryption targets, allowing for key rotation
  * Supports Additional Authenticated Data (AAD) for all KMSes except Vault Transit.

## Extras

There are several **extra**(s) packages included which build upon the base
go-kms-wrapping features to provide "extra" capabilities.  

* The
[`multi`](https://github.com/hashicorp/go-kms-wrapping/tree/main/extras/multi)
package is capable of encrypting to a specified wrapper and
decrypting using one of several wrappers switched on key ID. This can allow
easy key rotation for KMSes that do not natively support it.

* The
[`structwrapping`](https://github.com/hashicorp/go-kms-wrapping/tree/main/extras/structwrapping)
package allows for structs to have members encrypted and decrypted in a single
pass via a single wrapper. This can be used for workflows such as database
library callback functions to easily encrypt/decrypt data as it goes to/from
storage.

* The [`kms`](https://github.com/hashicorp/go-kms-wrapping/tree/main/extras/kms)
  package provides key management system features for wrappers
  including scoped [KEKs](https://en.wikipedia.org/wiki/Glossary_of_cryptographic_keys)
  and [DEKs](https://en.wikipedia.org/wiki/Glossary_of_cryptographic_keys) which
  are wrapped with an external KMS when stored in sqlite or postgres. 

* The [`crypto`](https://github.com/hashicorp/go-kms-wrapping/tree/main/extras/crypto) package provides additional operations like HMAC-SHA256 and a
  derived reader from which keys can be read.

## Installation

`go get github.com/hashicorp/go-kms-wrapping/v2`

## Overview

The library exports a `Wrapper` interface that is implemented by multiple
providers. For each provider, the standard flow is as follows:

1. Create a wrapper using the New method
1. Call `SetConfig` to pass either wrapper-specific options or use the
`wrapping.WithConfigMap` option to pass a configuration map
1. Use the wrapper as needed

It is possible, in `v2` of this library, to instantiate a wrapper as a
[`plugin`](https://github.com/hashicorp/go-kms-wrapping/tree/main/plugin). This
allows avoiding pulling dependencies of the wrapper directly into another
system's process space. See the [`example plugin-cli`](examples/plugin-cli/) for
a complete example on how to do build wrapper plugins and use them in an application or the [`test
plugins`](https://github.com/hashicorp/go-kms-wrapping/tree/main/plugin/testplugins)
for guidance in how to build a plugin; in this case, you'll definitely want to use
`wrapping.WithConfigMap` to pass configuration to avoid pulling in
package-specific options.

The best place to find the currently available set of configuration options
supported by each provider is its code, but it can also be found in [Vault's
seal configuration
documentation](https://www.vaultproject.io/docs/configuration/seal/index.html).
All environment variables noted there also work in this library, however,
non-Vault-specific variants of the environment variables are also available for
each provider. See the code/comments in each given provider for the currently
allowed env vars.

## Usage

Following is an example usage of the AWS KMS provider. 

```go
// Context used in this library is passed to various underlying provider
// libraries; how it's used is dependent on the provider libraries
ctx := context.Background()

wrapper := awskms.NewWrapper()
_, err := wrapper.SetConfig(ctx, wrapping.WithConfigMap(map[string]string{
    "kms_key_id": "1234abcd-12ab-34cd-56ef-1234567890ab",
}))
if err != nil {
    return err
}
blobInfo, err := wrapper.Encrypt(ctx, []byte("foo"))
if err != nil {
    return err
}

//
// Do some things...
//

plaintext, err := wrapper.Decrypt(ctx, blobInfo)
if err != nil {
    return err
}
if string(plaintext) != "foo" {
    return errors.New("mismatch between input and output")
}
```
