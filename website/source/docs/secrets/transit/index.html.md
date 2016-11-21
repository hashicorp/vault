---
layout: "docs"
page_title: "Secret Backend: Transit"
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

The first step to using the transit backend is to mount it. Unlike the `generic`
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
type                    aes256-gcm96
deletion_allowed        false
derived                 false
keys                    map[1:1.459861712e+09]
latest_version          1
min_decryption_version  1
name                    foo
````

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

### /transit/keys/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates a new named encryption key of the specified type. The values set
    here cannot be changed after key creation.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/keys/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">type</span>
        <span class="param-flags">required</span>
        The type of key to create. The currently-supported types are:
        <ul>
          <li>`aes256-gcm96`: AES-256 wrapped with GCM using a 12-byte nonce size (symmetric)</li>
          <li>`ecdsa-p256`: ECDSA using the P-256 elliptic curve (asymmetric)</li>
        </ul>
        Defaults to `aes256-gcm`.
      </li>
      <li>
        <span class="param">derived</span>
        <span class="param-flags">optional</span>
        Boolean flag indicating if key derivation MUST be used. If enabled, all
        encrypt/decrypt requests to this named key must provide a context
        which is used for key derivation. Defaults to false.
      </li>
      <li>
        <span class="param">convergent_encryption</span>
        <span class="param-flags">optional</span>
        If set, the key will support convergent encryption, where the same
        plaintext creates the same ciphertext. This requires _derived_ to be
        set to `true`. When enabled, each
        encryption(/decryption/rewrap/datakey) operation will derive a `nonce`
        value rather than randomly generate it. Note that while this is useful
        for particular situations, all nonce values used with a given context
        value **must be unique** or it will compromise the security of your
        key, and the key space for nonces is 96 bit -- not as large as the AES
        key itself. Defaults to false.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns information about a named encryption key. The `keys` object shows
    the creation time of each key version; the values are not the keys
    themselves. Depending on the type of key, different information may be
    returned, e.g. an asymmetric key will return its public key in a standard
    format for the type.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/transit/keys/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "type": "aes256-gcm96",
        "deletion_allowed": false,
        "derived": false,
        "keys": {
          "1": 1442851412
        },
        "min_decryption_version": 0,
        "name": "foo"
      }
    }
    ```

  </dd>
</dl>

#### LIST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns a list of keys. Only the key names are returned.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/transit/keys` (LIST) or `/transit/keys?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

  ```javascript
  {
    "data": {
      "keys": ["foo", "bar"]
    },
    "lease_duration": 0,
    "lease_id": "",
    "renewable": false
  }
  ```

  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a named encryption key.
    It will no longer be possible to decrypt any data encrypted with the
    named key. Because this is a potentially catastrophic operation, the
    `deletion_allowed` tunable must be set in the key's `/config` endpoint.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/transit/keys/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /transit/keys/config
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows tuning configuration values for a given key. (These values are
    returned during a read operation on the named key.)
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/keys/<name>/config`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">min_decryption_version</span>
        <span class="param-flags">optional</span>
        The minimum version of ciphertext allowed to be decrypted. Adjusting
        this as part of a key rotation policy can prevent old copies of
        ciphertext from being decrypted, should they fall into the wrong hands.
        For signatures, this value controls the minimum version of signature
        that can be verified against. For HMACs, this controls the minimum
        version of a key allowed to be used as the key for the HMAC function.
        Defaults to 0.
      </li>
      <li>
        <span class="param">allow_deletion</span>
        <span class="param-flags">optional</span>
        When set, the key is allowed to be deleted. Defaults to false.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /transit/keys/rotate/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Rotates the version of the named key. After rotation, new plaintext
    requests will be encrypted with the new version of the key. To upgrade
    ciphertext to be encrypted with the latest version of the key, use the
    `rewrap` endpoint. This is only supported with keys that support encryption
    and decryption operations.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/keys/<name>/rotate`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /transit/encrypt/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Encrypts the provided plaintext using the named key. Currently, this only
    supports symmetric keys. This path supports the `create` and `update`
    policy capabilities as follows: if the user has the `create` capability for
    this endpoint in their policies, and the key does not exist, it will be
    upserted with default values (whether the key requires derivation depends
    on whether the context parameter is empty or not). If the user only has
    `update` capability and the key does not exist, an error will be returned.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/encrypt/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">plaintext</span>
        <span class="param-flags">required</span>
        The plaintext to encrypt, provided as a base64-encoded string.
      </li>
      <li>
        <span class="param">context</span>
        <span class="param-flags">optional</span>
        The key derivation context, provided as a base64-encoded string.
        Must be provided if derivation is enabled.
      </li>
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">optional</span>
        The nonce value, provided as base64 encoded. Must be provided if
        convergent encryption is enabled for this key and the key was generated
        with Vault 0.6.1. Not required for keys created in 0.6.2+. The value
        must be exactly 96 bits (12 bytes) long and the user must ensure that
        for any given context (and thus, any given encryption key) this nonce
        value is **never reused**.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "ciphertext": "vault:v1:abcdefgh"
      }
    }
    ```

  </dd>
</dl>

### /transit/decrypt/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Decrypts the provided ciphertext using the named key. Currently, this only
    supports symmetric keys.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/decrypt/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ciphertext</span>
        <span class="param-flags">required</span>
        The ciphertext to decrypt, provided as returned by encrypt.
      </li>
      <li>
        <span class="param">context</span>
        <span class="param-flags">optional</span>
        The key derivation context, provided as a base64-encoded string.
        Must be provided if derivation is enabled.
      </li>
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">optional</span>
        The nonce value used during encryption, provided as base64 encoded.
        Must be provided if convergent encryption is enabled for this key and
        the key was created with Vault 0.6.1. Not required for keys created in
        0.6.2+.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveAo="
      }
    }
    ```

  </dd>
</dl>

### /transit/rewrap/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Rewrap the provided ciphertext using the latest version of the named key.
    Because this never returns plaintext, it is possible to delegate this
    functionality to untrusted users or scripts.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/rewrap/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ciphertext</span>
        <span class="param-flags">required</span>
        The ciphertext to decrypt, provided as returned by encrypt.
      </li>
      <li>
        <span class="param">context</span>
        <span class="param-flags">optional</span>
        The key derivation context, provided as base64-encoded string.
        Must be provided if derivation is enabled.
      </li>
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">optional</span>
        The nonce value used during encryption, provided as base64 encoded.
        Must be provided if convergent encryption is enabled for this key and
        the key was created with Vault 0.6.1. Not required for keys created in
        0.6.2+.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "ciphertext": "vault:v2:abcdefgh"
      }
    }
    ```

  </dd>
</dl>

### /transit/datakey/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generate a new high-entropy key and the value encrypted with the named
    key. Optionally return the plaintext of the key as well. Whether plaintext
    is returned depends on the path; as a result, you can use Vault ACL
    policies to control whether a user is allowed to retrieve the plaintext
    value of a key. This is useful if you want an untrusted user or operation
    to generate keys that are then made available to trusted users.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/datakey/<plaintext|wrapped>/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">plaintext|wrapped (path parameter)</span>
        <span class="param-flags">required</span>
        If `plaintext`, the plaintext key will be returned along with the
        ciphertext. If `wrapped`, only the ciphertext value will be returned.
      </li>
      <li>
        <span class="param">context</span>
        <span class="param-flags">optional</span>
        The key derivation context, provided as a base64-encoded string.
        Must be provided if derivation is enabled.
      </li>
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">optional</span>
        The nonce value, provided as base64 encoded. Must be provided if
        convergent encryption is enabled for this key and the key was generated
        with Vault 0.6.1. Not required for keys created in 0.6.2+. The value
        must be exactly 96 bits (12 bytes) long and the user must ensure that
        for any given context (and thus, any given encryption key) this nonce
        value is **never reused**.
      </li>
      <li>
        <span class="param">bits</span>
        <span class="param-flags">optional</span>
        The number of bits in the desired key. Can be 128, 256, or 512; if not
        given, defaults to 256.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveAo=",
        "ciphertext": "vault:v1:abcdefgh"
      }
    }
    ```

  </dd>
</dl>

### /transit/random
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Return high-quality random bytes of the specified length.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/random(/<bytes>)`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">bytes</span>
        <span class="param-flags">optional</span>
        The number of bytes to return. Defaults to 32 (256 bits). This value
        can be specified either in the request body, or as a part of the URL
        with a format like `/transit/random/48`.
      </li>
      <li>
        <span class="param">format</span>
        <span class="param-flags">optional</span>
        The output encoding; can be either `hex` or `base64`. Defaults to
        `base64`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "random_bytes": "dGhlIHF1aWNrIGJyb3duIGZveAo="
      }
    }
    ```

  </dd>
</dl>

### /transit/hash
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns the hash of given data using the specified algorithm. The algorithm
    can be specified as part of the URL or given via a parameter; the URL value
    takes precedence if both are set.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/hash(/<algorithm>)`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">input</span>
        <span class="param-flags">required</span>
        The base64-encoded input data.
      </li>
      <li>
        <span class="param">algorithm</span>
        <span class="param-flags">optional</span>
        The hash algorithm to use. This can also be specified in the URL.
        Currently-supported algorithms are:
        <ul>
          <li>`sha2-224`</li>
          <li>`sha2-256`</li>
          <li>`sha2-384`</li>
          <li>`sha2-512`</li>
        </ul>
        Defaults to `sha2-256`.
      </li>
      <li>
        <span class="param">format</span>
        <span class="param-flags">optional</span>
        The output encoding; can be either `hex` or `base64`. Defaults to
        `hex`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "sum": "dGhlIHF1aWNrIGJyb3duIGZveAo="
      }
    }
    ```

  </dd>
</dl>

### /transit/hmac/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns the digest of given data using the specified hash algorithm and the
    named key. The key can be of any type supported by `transit`; the raw key
    will be marshalled into bytes to be used for the HMAC function. If the key
    is of a type that supports rotation, the latest (current) version will be
    used.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/hmac/<name>(/<algorithm>)`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">input</span>
        <span class="param-flags">required</span>
        The base64-encoded input data.
      </li>
      <li>
        <span class="param">algorithm</span>
        <span class="param-flags">optional</span>
        The hash algorithm to use. This can also be specified in the URL.
        Currently-supported algorithms are:
        <ul>
          <li>`sha2-224`</li>
          <li>`sha2-256`</li>
          <li>`sha2-384`</li>
          <li>`sha2-512`</li>
        </ul>
        Defaults to `sha2-256`.
      </li>
      <li>
        <span class="param">format</span>
        <span class="param-flags">optional</span>
        The output encoding; can be either `hex` or `base64`. Defaults to
        `hex`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "hmac": "dGhlIHF1aWNrIGJyb3duIGZveAo="
      }
    }
    ```

  </dd>
</dl>

### /transit/sign/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns the cryptographic signature of the given data using the named key
    and the specified hash algorithm. The key must be of a type that supports
    signing.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/sign/<name>(/<algorithm>)`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">input</span>
        <span class="param-flags">required</span>
        The base64-encoded input data.
      </li>
      <li>
        <span class="param">algorithm</span>
        <span class="param-flags">optional</span>
        The hash algorithm to use. This can also be specified in the URL.
        Currently-supported algorithms are:
        <ul>
          <li>`sha2-224`</li>
          <li>`sha2-256`</li>
          <li>`sha2-384`</li>
          <li>`sha2-512`</li>
        </ul>
        Defaults to `sha2-256`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "signature": "vault:v1:MEUCIQCyb869d7KWuA0hBM9b5NJrmWzMW3/pT+0XYCM9VmGR+QIgWWF6ufi4OS2xo1eS2V5IeJQfsi59qeMWtgX0LipxEHI="
      }
    }
    ```

  </dd>
</dl>

### /transit/verify/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns whether the provided signature is valid for the given data.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/verify/<name>(/<algorithm>)`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">input</span>
        <span class="param-flags">required</span>
        The base64-encoded input data.
      </li>
      <li>
        <span class="param">signature</span>
        <span class="param-flags">required</span>
        The signature output from the `/transit/sign` function. Either this must be supplied or `hmac` must be supplied.
      </li>
      <li>
        <span class="param">hmac</span>
        <span class="param-flags">required</span>
        The signature output from the `/transit/hmac` function. Either this must be supplied or `signature` must be supplied.
      </li>
      <li>
        <span class="param">algorithm</span>
        <span class="param-flags">optional</span>
        The hash algorithm to use. This can also be specified in the URL.
        Currently-supported algorithms are:
        <ul>
          <li>`sha2-224`</li>
          <li>`sha2-256`</li>
          <li>`sha2-384`</li>
          <li>`sha2-512`</li>
        </ul>
        Defaults to `sha2-256`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "valid": true
      }
    }
    ```

  </dd>
</dl>
