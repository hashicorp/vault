---
layout: "docs"
page_title: "Secret Backend: Transit"
sidebar_current: "docs-secrets-transit"
description: |-
  The transit secret backend for Vault encrypts/decrypts data in-transit. It doesn't store any secrets.
---

# Transit Secret Backend

Name: `transit`

The transit secret backend is used to encrypt/data in-transit. Vault doesn't
store the data sent to the backend. It can also be viewed as "encryption as
a service."

The primary use case for the transit backend is to encrypt data from
applications while still storing that encrypted data in some primary data
store. This relieves the burden of proper encryption/decryption from
application developers and pushes the burden onto the operators of Vault.
Operators of Vault generally include the security team at an organization,
which means they can ensure that data is encrypted/decrypted properly.

As of Vault 0.2, the transit backend also supports doing key derivation. This
allows data to be encrypted within a context such that the same context must be
used for decryption. This can be used to enable per transaction unique keys which
further increase the security of data at rest.

Additionally, since encrypt/decrypt operations must enter the audit log,
any decryption event is recorded.

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
Key        	Value
name        foo
cipher_mode aes-gcm
derived     false
````

We can read from the `raw/` endpoint to see the encryption key itself:

```
$ vault read transit/raw/foo
Key        	Value
name       	foo
cipher_mode	aes-gcm
key        	PhKFTALCmhAhVQfMBAH4+UwJ6J2gybapUH9BsrtIgR8=
derived     false
````

Here we can see that the randomly generated encryption key being used, as
well as the AES-GCM cipher mode. We don't need to know any of this to use
the key however.

Now, if we wanted to encrypt a piece of plain text, we use the encrypt
endpoint using our named key:

```
$ echo -n "the quick brown fox" | base64 | vault write transit/encrypt/foo plaintext=-
Key       	Value
ciphertext	vault:v0:czEwyKqGZY/limnuzDCUUe5AK0tbBObWqeZgFqxCuIqq7A84SeiOq3sKD0Y/KUvv
```

The encryption endpoint expects the plaintext to be provided as a base64 encoded
strings, so we must first convert it. Vault does not store the plaintext or the
ciphertext, but only handles it _in transit_ for processing. The application
is free to store the ciphertext in a database or file at rest.

To decrypt, we simply use the decrypt endpoint using the same named key:

```
$ vault write transit/decrypt/foo ciphertext=vault:v0:czEwyKqGZY/limnuzDCUUe5AK0tbBObWqeZgFqxCuIqq7A84SeiOq3sKD0Y/KUvv
Key      	Value
plaintext	dGhlIHF1aWNrIGJyb3duIGZveAo=

$ echo "dGhlIHF1aWNrIGJyb3duIGZveAo=" | base64 -D
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
    Creates a new named encryption key. This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/transit/keys/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">derived</span>
        <span class="param-flags">optional</span>
        Boolean flag indicating if key derivation MUST be used.
        If enabled, all encrypt/decrypt requests to this named key
        must provide a context which is used for key derivation.
        Defaults to false.
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
    Returns information about a named encryption key.
    This is a root protected endpoint.
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
          "name":        "foo",
          "cipher_mode": "aes-gcm",
          "derived":     "true",
          "kdf_mode":    "hmac-sha256-counter",
      }
    }
    ```

  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a named encryption key. This is a root protected endpoint.
    All data encrypted with the named key will no longer be decryptable.
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

### /transit/encrypt/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Encrypts the provided plaintext using the named key. If the named key
    does not already exist, it will be automatically generated for the given
    name with the default parameters.
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
        The plaintext to encrypt, provided as base64 encoded.
      </li>
      <li>
        <span class="param">context</span>
        <span class="param-flags">optional</span>
        The key derivation context, provided as base64 encoded.
        Must be provided if the derivation enabled.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
        "data": {
            "ciphertext": "vault:v0:abcdefgh"
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
    Decrypts the provided ciphertext using the named key.
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
        The key derivation context, provided as base64 encoded.
        Must be provided if the derivation enabled.
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

### /transit/raw/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns raw information about a named encryption key,
    Including the underlying encryption key. This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/transit/raw/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
          "name":        "foo",
          "cipher_mode": "aes-gcm",
          "key":         "PhKFTALCmhAhVQfMBAH4+UwJ6J2gybapUH9BsrtIgR8="
          "derived":     "true",
          "kdf_mode":    "hmac-sha256-counter",
      }
    }
    ```

  </dd>
</dl>


