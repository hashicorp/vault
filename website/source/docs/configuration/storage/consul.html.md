---
layout: "docs"
page_title: "Consul - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-consul"
description: |-
  The Consul storage backend is used to persist Vault's data in Consul's
  key-value store. In addition to providing durable storage, inclusion of this
  backend will also register Vault as a service in Consul with a default health
  check.
---

# Consul Storage Backend

The Consul storage backend is used to persist Vault's data in [Consul's][consul]
key-value store. In addition to providing durable storage, inclusion of this
backend will also register Vault as a service in Consul with a default health
check.

- **High Availability** – the Consul storage backend supports high availability.

- **HashiCorp Supported** – the Consul storage backend is officially supported
  by HashiCorp.

```hcl
storage "consul" {
  address = "127.0.0.1:8500"
  path    = "vault"
}
```

Once properly configured, an unsealed Vault installation should be available and
accessible at:

```text
active.vault.service.consul
```

Unsealed Vault instances in standby mode are available at:

```text
standby.vault.service.consul
```

All unsealed Vault instances are available as healthy at:

```text
vault.service.consul
```

Sealed Vault instances will mark themselves as unhealthy to avoid being returned
at Consul's service discovery layer.


## `consul` Parameters

- `address` `(string: "127.0.0.1:8500")` – Specifies the address of the Consul
  agent to communicate with. This can be an IP address, DNS record, or unix
  socket. It is recommended that you communicate with a local Consul agent; do
  not communicate directly with a server.

- `check_timeout` `(string: "5s")` – Specifies the check interval used to send
  health check information back to Consul. This is specified using a label
  suffix like `"30s"` or `"1h"`.

- `consistency_mode` `(string: "default")` – Specifies the Consul
  [consistency mode][consul-consistency]. Possible values are `"default"` or
  `"strong"`.

- `disable_registration` `(bool: false)` – Specifies whether Vault should
  register itself with Consul.

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests to Consul.

- `path` `(string: "vault/")` – Specifies the path in Consul's key-value store
  where Vault data will be stored.

- `scheme` `(string: "http")` – Specifies the scheme to use when communicating
  with Consul. This can be set to "http" or "https". It is highly recommended
  you communicate with Consul over https over non-local connections. When
  communicating over a unix socket, this option is ignored.

- `service` `(string: "vault")` – Specifies the name of the service to register
  in Consul.

- `service_tags` `(string: "")` – Specifies a comma-separated list of tags to
  attach to the service registration in Consul.

- `token` `(string: "")` – Specifies the [Consul ACL token][consul-acl] with
  permission to read and write from the `path` in Consul's key-value store.
  This is **not** a Vault token. See the ACL section below for help.

The following settings apply when communicating with Consul via an encrypted
connection. You can read more about encrypting Consul connections on the
[Consul encryption page][consul-encryption].

- `tls_ca_file` `(string: "")` – Specifies the path to the CA certificate used
  for Consul communication. This defaults to system bundle if not specified.
  This should be set according to the
  [`ca_file`](https://www.consul.io/docs/agent/options.html#ca_file) setting in
  Consul.

- `tls_cert_file` `(string: "")` (optional) – Specifies the path to the
  certificate for Consul communication. This should be set according to the
  [`cert_file`](https://www.consul.io/docs/agent/options.html#cert_file) setting
  in Consul.

- `tls_key_file` `(string: "")` – Specifies the path to the private key for
  Consul communication. This should be set according to the
  [`key_file`](https://www.consul.io/docs/agent/options.html#key_file) setting
  in Consul.

- `tls_min_version` `(string: "tls12")` – Specifies the minimum TLS version to
  use. Accepted values are `"tls10"`, `"tls11"` or `"tls12"`.

- `tls_skip_verify` `(bool: false)` – Specifies if the TLS host verification
  should be disabled. It is highly discouraged that you disable this option.

## ACLs

If using ACLs in Consul, you'll need appropriate permissions. For Consul 0.8,
the following will work for most use-cases, assuming that your service name is
`vault` and the prefix being used is `vault/`:

```json
{
  "key": {
    "vault/": {
      "policy": "write"
    }
  },
  "node": {
    "": {
      "policy": "write"
    }
  },
  "service": {
    "vault": {
      "policy": "write"
    }
  },
  "agent": {
    "": {
      "policy": "write"
    }

  },
  "session": {
    "": {
      "policy": "write"
    }
  }
}
```

## `consul` Examples

### Local Agent

This example shows a sample physical backend configuration which communicates
with a local Consul agent running on `127.0.0.1:8500`.

```hcl
storage "consul" {}
```

### Detailed Customization

This example shows communicating with Consul on a custom address with an ACL
token.

```hcl
storage "consul" {
  address = "10.5.7.92:8194"
  token   = "abcd1234"
}
```

### Custom Storage Path

This example shows storing data at a custom path in Consul's key-value store.
This path must be readable and writable by the Consul ACL token, if Consul
configured to use ACLs.

```hcl
storage "consul" {
  path = "vault/"
}
```

### Consul via Unix Socket

This example shows communicating with Consul over a local unix socket.

```hcl
storage "consul" {
  address = "unix:///tmp/.consul.http.sock"
}
```

### Custom TLS

This example shows using a custom CA, certificate, and key file to securely
communicate with Consul over TLS.

```hcl
storage "consul" {
  scheme        = "https"
  tls_ca_file   = "/etc/pem/vault.ca"
  tls_cert_file = "/etc/pem/vault.cert"
  tls_key_file  = "/etc/pem/vault.key"
}
```

[consul]: https://www.consul.io/ "Consul by HashiCorp"
[consul-acl]: https://www.consul.io/docs/guides/acl.html "Consul ACLs"
[consul-consistency]: https://www.consul.io/api/index.html#consistency-modes "Consul Consistency Modes"
[consul-encryption]: https://www.consul.io/docs/agent/encryption.html "Consul Encryption"
