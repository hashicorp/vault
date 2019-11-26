---
layout: "docs"
page_title: "Vault Agent Auto-Auth"
sidebar_title: "Auto-Auth"
sidebar_current: "docs-agent-autoauth"
description: |-
  Vault Agent's Auto-Auth functionality allows easy and automatic
  authentication to Vault in a variety of environments.
---

# Vault Agent Auto-Auth

The Auto-Auth functionality of Vault Agent allows for easy authentication in a
wide variety of environments.

## Functionality

Auto-Auth consists of two parts: a Method, which is the authentication method
that should be used in the current environment; and one or more Sinks, which
are locations where the agent should write a token any time the current token
value has changed.

When the agent is started with Auto-Auth enabled, it will attempt to acquire a
Vault token using the configured Method. On failure, it will back off for a
short while (including some randomness to help prevent thundering herd
scenarios) and retry. On success, unless the auth method is configured to wrap
the tokens, it will keep the resulting token renewed until renewal is no longer
allowed or fails, at which point it will attempt to reauthenticate.

Every time an authentication is successful, the token is written to the
configured Sinks, subject to their configuration.

## Advanced Functionality

Sinks support some advanced features, including the ability for the written
values to be encrypted or
[response-wrapped](/docs/concepts/response-wrapping.html).

Both mechanisms can be used concurrently; in this case, the value will be
response-wrapped, then encrypted.

### Response-Wrapping Tokens

There are two ways that tokens can be response-wrapped by the agent:

1. By the auth method. This allows the end client to introspect the
   `creation_path` of the token, helping prevent Man-In-The-Middle (MITM)
   attacks. However, because the agent cannot then unwrap the token and rewrap
   it without modifying the `creation_path`, the agent is not able to renew the
   token; it is up to the end client to renew the token. The agent stays
   daemonized in this mode since some auth methods allow for reauthentication
   on certain events.

2. By any of the token sinks. Because more than one sink can be configured, the
   token must be wrapped after it is fetched, rather than wrapped by Vault as
   it's being returned. As a result, the `creation_path` will always be
   `sys/wrapping/wrap`, and validation of this field cannot be used as
   protection against MITM attacks. However, this mode allows the agent to keep
   the token renewed for the end client and automatically reauthenticate when
   it expires.

### Encrypting Tokens

 ~> This is experimental; if input/output formats change we will make every
 effort to provide backwards compatibility.

Tokens can be encrypted, using a Diffie-Hellman exchange to generate an
ephemeral key. In this mechanism, the client receiving the token writes a
generated public key to a file. The sink responsible for writing the token to
that client looks for this public key and uses it to compute a shared secret
key, which is then used to encrypt the token via AES-GCM. The nonce, encrypted
payload, and the sink's public key are then written to the output file, where
the client can compute the shared secret and decrypt the token value.

~> NOTE: This is not a protection against MITM attacks! The purpose of this
feature is for forward-secrecy and coverage against bare token values being
persisted. A MITM that can write to the sink's output and/or client public-key
input files could attack this exchange.

To help mitigate MITM attacks, additional authenticated data (AAD) can be
provided to the agent. This data is written as part of the AES-GCM tag and must
match on both the agent and the client. This of course means that protecting
this AAD becomes important, but it provides another layer for an attacker to
have to overcome. For instance, if the attacker has access to the file system
where the token is being written, but not to read agent configuration or read
environment variables, this AAD can be generated and passed to the agent and
the client in ways that would be difficult for the attacker to find.

When using AAD, it is always a good idea for this to be as fresh as possible;
generate a value and pass it to your client and agent on startup. Additionally,
agent uses a Trust On First Use model; after it finds a generated public key,
it will reuse that public key instead of looking for new values that have been
written.

If writing a client that uses this feature, it will likely be helpful to look
at the
[dhutil](https://github.com/hashicorp/vault/blob/master/helper/dhutil/dhutil.go)
library. This shows the expected format of the public key input and envelope
output formats.

## Configuration

The top level `auto_auth` block has two configuration entries:

- `method` `(object: required)` - Configuration for the method

- `sinks` `(array of objects: required)` - Configuration for the sinks

### Configuration (Method)

These are common configuration values that live within the `method` block:

- `type` `(string: required)` - The type of the method to use, e.g. `aws`,
  `gcp`, `azure`, etc. *Note*: when using HCL this can be used as the key for
  the block, e.g. `method "aws" {...}`.

- `mount_path` `(string: optional)` - The mount path of the method. If not
  specified, defaults to a value of `auth/<method type>`.
  
- `namespace` `(string: optional)` - The default namespace path for the mount.
  If not specified, defaults to the root namespace. 

- `wrap_ttl` `(string or integer: optional)` - If specified, the written token
  will be response-wrapped by the agent. This is more secure than wrapping by
  sinks, but does not allow the agent to keep the token renewed or
  automatically reauthenticate when it expires. Rather than a simple string,
  the written value will be a JSON-encoded
  [SecretWrapInfo](https://godoc.org/github.com/hashicorp/vault/api#SecretWrapInfo)
  structure. Values can be an integer number of seconds or a stringish value
  like `5m`.

- `config` `(object: required)` - Configuration of the method itself. See the
  sidebar for information about each method.

### Configuration (Sinks)

These configuration values are common to all Sinks:

- `type` `(string: required)` - The type of the method to use, e.g. `file`.
  *Note*: when using HCL this can be used as the key for the block, e.g. `sink
  "file" {...}`.

- `wrap_ttl` `(string or integer: optional)` - If specified, the written token
  will be response-wrapped by the sink. This is less secure than wrapping by
  the method, but allows the agent to keep the token renewed and automatically
  reauthenticate when it expires. Rather than a simple string, the written
  value will be a JSON-encoded
  [SecretWrapInfo](https://godoc.org/github.com/hashicorp/vault/api#SecretWrapInfo)
  structure. Values can be an integer number of seconds or a stringish value
  like `5m`.

- `dh_type` `(string: optional)` - If specified, the type of Diffie-Hellman exchange to
  perform, meaning, which ciphers and/or curves. Currently only `curve25519` is
  supported.

- `dh_path` `(string: required if dh_type is set)` - The path from which the
  agent should read the client's initial parameters (e.g. curve25519 public
  key).

- `aad` `(string: optional)` - If specified, additional authenticated data to
  use with the AES-GCM encryption of the token. Can be any string, including
  serialized data.

- `aad_env_var` `(string: optional)` - If specified, AAD will be read from the
  given environment variable rather than a value in the configuration file.

- `config` `(object: required)` - Configuration of the sink itself. See the
  sidebar for information about each sink.
