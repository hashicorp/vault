---
layout: "docs"
page_title: "GCP Cloud KMS - Seals - Configuration"
sidebar_current: "docs-configuration-seal-gcpckms"
description: |-
  The GCP Cloud KMS seal configures Vault to use GCP Cloud KMS as the seal wrapping
  mechanism.
---

# `gcpckms` Seal

The GCP Cloud KMS seal configures Vault to use GCP Cloud KMS as the seal
wrapping mechanism. Vault Enterprise's GCP Cloud KMS seal is activated by one of
the following:

* The presence of a `seal "gcpckms"` block in Vault's configuration file.
* The presence of the environment variable `VAULT_SEAL_TYPE` set to `gcpckms`.
  If enabling via environment variable, all other required values specific to
  Cloud KMS (i.e. `VAULT_GCPCKMS_SEAL_KEY_RING`, etc.) must be also supplied, as
  well as all other GCP-related environment variables that lends to successful
  authentication (i.e. `GOOGLE_PROJECT`, etc.).

## Authentication

Authentication-related values must be provided, either as enviroment
variables or as configuration parameters.

```text
GCP authentication values:

* `GOOGLE_CREDENTIALS` or `GOOGLE_APPLICATION_CREDENTIALS`
* `GOOGLE_PROJECT`
* `GOOGLE_REGION`
```

Note: The client uses the official Google SDK and will use the specified
credentials, environment credentials, or [application default
credentials](https://developers.google.com/identity/protocols/application-default-credentials)
in that order, if the above GCP specific values are not provided.

## `gcpckms` Parameters

These parameters apply to the `seal` stanza in the Vault configuration file:

* `credentials` `(string: <required>)`: The path to the credentials JSON file
  to use. May be also specified by the `GOOGLE_CREDENTIALS` or
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable or set automatically if
  running under Google App Engine, Google Compute Engine or Google Container
  Engine.
* `project` `(string: <required>)`: The GCP project ID to use. May also be
  specified by the `GOOGLE_PROJECT` environment variable.
* `region` `(string: "us-east-1")`: The GCP region/location where the key ring
  lives. May also be specified by the `GOOGLE_REGION` environment variable.
* `key_ring` `(string: <required>)`: The GCP CKMS key ring to use. May also be
  specified by the `VAULT_GCPCKMS_SEAL_KEY_RING` environment variable.
* `crypto_key` `(string: <required>)`: The GCP CKMS crypto key to use for
  encryption and decryption. May also be specified by the
  `VAULT_GCPCKMS_SEAL_CRYPTO_KEY` environment variable.

## `gcpckms` Environment Variables

Alternatively, the GCP Cloud KMS seal can be activated by providing the following
environment variables:

```text
* `VAULT_SEAL_TYPE`
* `VAULT_GCPCKMS_SEAL_KEY_RING`
* `VAULT_GCPCKMS_SEAL_CRYPTO_KEY`
```

## `gcpckms` Example

This example shows configuring GCP Cloud KMS seal through the Vault
configuration file by providing all the required values:

```hcl
seal "gcpckms" {
  credentials = "/usr/vault/vault-project-user-creds.json" 
  project     = "vault-project"
  region      = "global"
  key_ring    = "vault-keyring"
  crypto_key  = "vault-key"
}
```