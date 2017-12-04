---
layout: "docs"
page_title: "Transit Secret Backend"
sidebar_current: "docs-secrets-transit"
description: |-
  The transit secret backend for Vault encrypts/decrypts data in-transit. It doesn't store any secrets.
---

# Transit Secret Backend

Name: `transit`

The transit secret backend handles cryptographic functions on data in-transit.
Vault doesn't store the data sent to the backend. It can also be viewed as
"cryptography as a service."

The primary use case for `transit` is to encrypt data from applications while
still storing that encrypted data in some primary data store. This relieves the
burden of proper encryption/decryption from application developers and pushes
the burden onto the operators of Vault.  Operators of Vault generally include
the security team at an organization, which means they can ensure that data is
encrypted/decrypted properly. Additionally, since encrypt/decrypt operations
must enter the audit log, any decryption event is recorded.

`transit` can also sign and verify data; generate hashes and HMACs of data; and
act as a source of random bytes.

Due to Vault's flexible ACLs, other interesting use-cases are possible. For
instance, one set of Internet-facing servers can be given permission to encrypt
with a named key but not decrypt with it; a separate set of servers not
directly connected to the Internet can then perform decryption, reducing the
data's attack surface.

Key derivation is supported, which allows the same key to be used for multiple
purposes by deriving a new key based on a user-supplied context value. In this
mode, convergent encryption can optionally be supported, which allows the same
input values to produce the same ciphertext.

The backend also supports key rotation, which allows a new version of the named
key to be generated. All data encrypted with the key will use the newest
version of the key; previously encrypted data can be decrypted using old
versions of the key. Administrators can control which previous versions of a
key are available for decryption, to prevent an attacker gaining an old copy of
ciphertext to be able to successfully decrypt it. At any time, a legitimate
user can "rewrap" the data, providing an old version of the ciphertext and
receiving a new version encrypted with the latest key. Because rewrapping does
not expose the plaintext, using Vault's ACL system, this can even be safely
performed by unprivileged users or cron jobs.

Datakey generation allows processes to request a high-entropy key of a given
bit length be returned to them, encrypted with the named key. Normally this will
also return the key in plaintext to allow for immediate use, but this can be
disabled to accommodate auditing requirements.

N.B.: As part of adding rotation support, the initial version of a named key
produces ciphertext starting with version 1, i.e. containing `:v1:`. Keys from
very old versions of Vault, when rotated, will jump to version 2 despite their
previous ciphertext output containing `:v0:`. Decryption, however, treats
version 0 and version 1 the same, so old ciphertext will still work.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the transit backend is to mount it. Unlike the `kv`
backend, the `transit` backend is not mounted by default.

```
$ vault mount transit
Successfully mounted 'transit' at 'transit'!
```

The next step is to create a named encryption key. A named key is used so that
many different applications can use the transit backend with independent keys.
This is done by doing a write against the backend:

```
$ vault write -f transit/keys/foo
Success! Data written to: transit/keys/foo
```

This will create the "foo" named key in the transit backend. We can inspect
the settings of the "foo" key by reading it:

```
$ vault read transit/keys/foo
Key                     Value
deletion_allowed      	false
derived               	false
exportable            	false
keys                  	map[1:1484070923]
latest_version        	1
min_decryption_version	1
name                  	foo
supports_decryption   	true
supports_derivation   	true
supports_encryption   	true
supports_signing      	false
type                  	aes256-gcm96
```

Now, if we wanted to encrypt a piece of plain text, we use the encrypt
endpoint using our named key:

```
$ echo -n "the quick brown fox" | base64 | vault write transit/encrypt/foo plaintext=-
Key       	Value
ciphertext	vault:v1:czEwyKqGZY/limnuzDCUUe5AK0tbBObWqeZgFqxCuIqq7A84SeiOq3sKD0Y/KUvv
```

The encryption endpoint expects the plaintext to be provided as a base64 encoded
strings, so we must first convert it. Vault does not store the plaintext or the
ciphertext, but only handles it _in transit_ for processing. The application
is free to store the ciphertext in a database or file at rest.

To decrypt, we simply use the decrypt endpoint using the same named key:

```
$ vault write transit/decrypt/foo ciphertext=vault:v1:czEwyKqGZY/limnuzDCUUe5AK0tbBObWqeZgFqxCuIqq7A84SeiOq3sKD0Y/KUvv
Key      	Value
plaintext	dGhlIHF1aWNrIGJyb3duIGZveAo=

$ echo "dGhlIHF1aWNrIGJyb3duIGZveAo=" | base64 -d
the quick brown fox
```

Using ACLs, it is possible to restrict using the transit backend such
that trusted operators can manage the named keys, and applications can
only encrypt or decrypt using the named keys they need access to.

## API

The Transit secret backend has a full HTTP API. Please see the
[Transit secret backend API](/api/secret/transit/index.html) for more
details.
