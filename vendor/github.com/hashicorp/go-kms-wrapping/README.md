# Go-KMS-Wrapping - Go library for encrypting values through various KMS providers

[![Godoc](https://godoc.org/github.com/hashicorp/go-kms-wrapping?status.svg)](https://godoc.org/github.com/hashicorp/go-kms-wrapping)

*NOTE*: Currently no compatibility guarantees are provided for this library; we
expect tags to remain in the `0.x.y` range. Function signatures, interfaces,
etc. may change at any time.

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


- [Features](#features)
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
  * * Yandex.Cloud KMS (uses envelopes)
  * Transparently supports multiple decryption targets, allowing for key rotation
  * Supports Additional Authenticated Data (AAD) for all KMSes except Vault Transit.

A
[`multiwrapper`](https://github.com/hashicorp/go-kms-wrapping/tree/master/wrappers/multiwrapper)
KMS is also included, capable of encrypting to a specified wrapper and
decrypting using one of several wrappers switched on key ID. This can allow
easy key rotation for KMSes that do not natively support it.

The
[`structwrapping`](https://github.com/hashicorp/go-kms-wrapping/tree/master/structwrapping)
package allows for structs to have members encrypted and decrypted in a single
pass via a single wrapper. This can be used for workflows such as database
library callback functions to easily encrypt/decrypt data as it goes to/from
storage.

## Installation

Import like any other library; supports go modules. It has not been tested with
non-`go mod` vendoring tools.

## Overview

The library exports a `Wrapper` interface that is implemented by multiple
providers. Each of these providers may have some functions specific to them,
usually to pass configuration information. A normal workflow is to create the
provider directly, pass it any needed configuration via the provider-specific
methods, and then have the rest of your code use the `Wrapper` interface.

Some of the functions make use of option structs that are currently empty. This
is to allow options to be added later without breaking backwards compatibility.

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

wrapper := awskms.NewWrapper(nil)
_, err := wrapper.SetConfig(&map[string]string{
    "kms_key_id": "1234abcd-12ab-34cd-56ef-1234567890ab"
})
if err != nil {
    return err
}
blobInfo, err := wrapper.Encrypt(ctx, []byte{"foo"}, nil)
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
