---
layout: "docs"
page_title: "Vault Transit - Seals - Configuration"
sidebar_title: "Vault Transit"
sidebar_current: "docs-configuration-seal-transit"
description: |-
  The Transit seal configures Vault to use Vault's Transit Secret Engine as the
  autoseal mechanism.
---

# `transit` Seal

The Transit seal configures Vault to use Vault's Transit Secret Engine as the
autoseal mechanism.
The Transit seal is activated by one of the following:

* The presence of a `seal "transit"` block in Vault's configuration file
* The presence of the environment variable `VAULT_SEAL_TYPE` set to `transit`.

## `transit` Example

This example shows configuring Transit seal through the Vault configuration file
by providing all the required values:

```hcl
seal "transit" {
  address            = "https://vault:8200"
  token              = "s.Qf1s5zigZ4OX6akYjQXJC1jY"
  disable_renewal    = "false"

  // Key configuration
  key_name           = "transit_key_name"
  mount_path         = "transit/"
  namespace          = "ns1/"

  // TLS Configuration
  tls_ca_cert        = "/etc/vault/ca_cert.pem"
  tls_client_cert    = "/etc/vault/client_cert.pem"
  tls_client_key     = "/etc/vault/ca_cert.pem"
  tls_server_name    = "vault"
  tls_skip_verify    = "false"
}
```

## `tranist` Parameters

These parameters apply to the `seal` stanza in the Vault configuration file:

- `address` `(string: <required>)`: The full address to the Vault cluster.
  This may also be specified by the `VAULT_ADDR` environment variable.

- `token` `(string: <required>)`: The Vault token to use. This may also be
  specified by the `VAULT_TOKEN` environment variable.

- `key_name` `(string: <required>)`: The transit key to use for encryption and
  decryption.  This may also be supplied using the `VAULT_TRANSIT_SEAL_KEY_NAME`
  environment variable.

- `mount_path` `(string: <required>)`: The mount path to the transit secret engine.  
  This may also be supplied using the `VAULT_TRANSIT_SEAL_MOUNT_PATH` environment
  variable.

- `namespace` `(string: "")`: The namespace path to the transit secret engine.  
  This may also be supplied using the `VAULT_NAMESPACE` environment variable.

- `disable_renewal` `(string: "false")`: Disables the automatic renewal of the token
  in case the lifecyle of the token is managed with some other mechanism outside of
  Vault, such as Vault Agent.  This may also be specfied using the
  `VAULT_TRANSIT_SEAL_DISABLE_RENEWAL` environment variable.

- `tls_ca_cert` `(string: "")`: Specifies the path to the CA certificate file used
  for communication with the Vault server.  This may also be specified using the
  `VAULT_CA_CERT` environment variable.

- `tls_client_cert` `(string: "")`: Specifies the path to the client certificate
  for communication with the Vault server.  This may also be specified using the
  `VAULT_CLIENT_CERT` environment variable.

- `tls_client_key` `(string: "")`: Specifies the path to the private key for
  communication with the Vault server.  This may also be specified using the
  `VAULT_CLIENT_KEY` environment variable.

- `tls_server_name` `(string: "")`: Name to use as the SNI host when connecting
  to the Vault server via TLS.  This may also be specified via the
  `VAULT_TLS_SERVER_NAME` environment variable.

- `tls_skip_verify` `(bool: "false")`: Disable verification of TLS certificates.
  Using this option is highly discouraged and decreases the security of data
  transmissions to and from the Vault server.  This may also be specified using the
  `VAULT_TLS_SKIP_VERIFY` environment variable.

## Authentication

Authentication-related values must be provided, either as environment
variables or as configuration parameters.

~> **Note:** Although the configuration file allows you to pass in
`VAULT_TOKEN` as part of the seal's parameters, it is *strongly* recommended
to set these values via environment variables.

The Vault token used to authenticate needs the following permissions on the
transit key:

```hcl
path "<mount path>/encrypt/<key name>" {
  capabilities = ["update"]
}

path "<mount path>/decrypt/<key name>" {
  capabilities = ["update"]
}
```

## Key Rotation

This seal supports key rotation using the Transit Secret Engine's key rotation endpoints.  See
[doc](/api/secret/transit/index.html#rotate-key). Old keys must not be disabled or deleted and are
used to decrypt older data.
