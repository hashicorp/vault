---
layout: "docs"
page_title: "Transit - Secrets Engines"
sidebar_current: "docs-secrets-transit"
description: |-
  The transit secrets engine for Vault encrypts/decrypts data in-transit. It doesn't store any secrets.
---

# Transit Secrets Engine

The transit secrets engine handles cryptographic functions on data in-transit.
Vault doesn't store the data sent to the secrets engine. It can also be viewed
as "cryptography as a service" or "encryption as a service". The transit secrets
engine can also sign and verify data; generate hashes and HMACs of data; and act
as a source of random bytes.

The primary use case for `transit` is to encrypt data from applications while
still storing that encrypted data in some primary data store. This relieves the
burden of proper encryption/decryption from application developers and pushes
the burden onto the operators of Vault.

Key derivation is supported, which allows the same key to be used for multiple
purposes by deriving a new key based on a user-supplied context value. In this
mode, convergent encryption can optionally be supported, which allows the same
input values to produce the same ciphertext.

Datakey generation allows processes to request a high-entropy key of a given
bit length be returned to them, encrypted with the named key. Normally this will
also return the key in plaintext to allow for immediate use, but this can be
disabled to accommodate auditing requirements.

## Working Set Management

This secrets engine does not currently delete keys. Keys that are out of the
working set (earlier than a key's specified `min_decryption_version` are
instead archived. This is a performance consideration to keep key loading fast,
as well as a security consideration: by disallowing decryption of old versions
of keys, found ciphertext corresponding to obsolete (but sensitive) data can
not be decrypted by most users, but in an emergency the
`min_decryption_version` can be moved back to allow for legitimate decryption.

Currently this archive is stored in a single storage entry. With some storage
backends, notably those using Raft or Paxos for HA capabilities, frequent
rotation may lead to a storage entry size for the archive that is larger than
the storage backend can handle. For frequent rotation needs, using named keys
that correspond to time bounds (e.g. five-minute periods floored to the closest
multiple of five) may provide a good alternative, allowing for several keys to
be live at once and a deterministic way to decide which key to use at any given
time.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the Transit secrets engine:

    ```text
    $ vault secrets enable transit
    Success! Enabled the transit secrets engine at: transit/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Create a named encryption key ring:

    ```text
    $ vault write -f transit/keys/my-key
    Success! Data written to: transit/keys/my-key
    ```

    Usually each application has its own encryption key ring.

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can use this secrets engine.

1. Encrypt some plaintext data using the `/encrypt` endpoint with a named key:

    ```text
    $ vault write transit/encrypt/my-key plaintext=$(base64 <<< "my secret data")

    Key           Value
    ---           -----
    ciphertext    vault:v1:8SDd3WHDOjf7mq69CyCqYjBXAiQQAVZRkFM13ok481zoCmHnSeDX9vyf7w==
    ```

    All plaintext data **must be base64-encoded**. The reason for this
    requirement is that Vault does not require that the plaintext is "text". It
    could be a binary file such as a PDF or image. The easiest safe transport
    mechanism for this data as part of a JSON payload is to base64-encode it.

    Note that Vault does not _store_ any of this data. The caller is responsible
    for storing the encrypted ciphertext. When the caller wants the plaintext,
    it must provide the ciphertext back to Vault to decrypt the value.

1. Decrypt a piece of data using the `/decrypt` endpoint with a named key:

    ```text
    $ vault write transit/decrypt/my-key ciphertext=vault:v1:8SDd3WHDOjf7mq69CyCqYjBXAiQQAVZRkFM13ok481zoCmHnSeDX9vyf7w==

    Key          Value
    ---          -----
    plaintext    bXkgc2VjcmV0IGRhdGEK
    ```

    The resulting data is base64-encoded (see the note above for details on
    why). Decode it to get the raw plaintext:

    ```text
    $ base64 --decode <<< "bXkgc2VjcmV0IGRhdGEK"
    my secret data
    ```

    It is also possible to script this decryption using some clever shell
    scripting in one command:

    ```text
    $ vault write -field=plaintext transit/decrypt/my-key ciphertext=... | base64 --decode
    my secret data
    ```

    Using ACLs, it is possible to restrict using the transit secrets engine such
    that trusted operators can manage the named keys, and applications can only
    encrypt or decrypt using the named keys they need access to.

1. Rotate the underlying encryption key. This will generate a new encryption key
and add it to the keyring for the named key:

    ```text
    $ vault write -f transit/keys/my-key/rotate
    Success! Data written to: transit/keys/my-key/rotate
    ```

    Future encryptions will use this new key. Old data can still be decrypted
    due to the use of a key ring.

1. Upgrade already-encrypted data to a new key. Vault will decrypt the value
using the appropriate key in the keyring and then encrypted the resulting
plaintext with the newest key in the keyring.

    ```text
    $ vault write transit/rewrap/my-key ciphertext=vault:v1:8SDd3WHDOjf7mq69CyCqYjBXAiQQAVZRkFM13ok481zoCmHnSeDX9vyf7w==

    Key           Value
    ---           -----
    ciphertext    vault:v2:0VHTTBb2EyyNYHsa3XiXsvXOQSLKulH+NqS4eRZdtc2TwQCxqJ7PUipvqQ==
    ```

    This process **does not** reveal the plaintext data. As such, a Vault policy
    could grant almost an untrusted process the ability to "rewrap" encrypted
    data, since the process would not be able to get access to the plaintext
    data.

## API

The Transit secrets engine has a full HTTP API. Please see the
[Transit secrets engine API](/api/secret/transit/index.html) for more
details.
