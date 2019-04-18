---
layout: "docs"
page_title: "TCP - Listeners - Configuration"
sidebar_title: "TCP"
sidebar_current: "docs-configuration-listener-tcp"
description: |-
  The TCP listener configures Vault to listen on the specified TCP address and
  port.
---

# `tcp` Listener

The TCP listener configures Vault to listen on a TCP address/port.

```hcl
listener "tcp" {
  address = "127.0.0.1:8200"
}
```

The `listener` stanza may be specified more than once to make Vault listen on
multiple interfaces. If you configure multiple listeners you also need to
specify [`api_addr`][api-addr] and [`cluster_addr`][cluster-addr] so Vault will
advertise the correct address to other nodes.

## `tcp` Listener Parameters

- `address` `(string: "127.0.0.1:8200")` – Specifies the address to bind to for
  listening.

- `cluster_address` `(string: "127.0.0.1:8201")` – Specifies the address to bind
  to for cluster server-to-server requests. This defaults to one port higher
  than the value of `address`. This does not usually need to be set, but can be
  useful in case Vault servers are isolated from each other in such a way that
  they need to hop through a TCP load balancer or some other scheme in order to
  talk.

- `max_request_size` `(int: 33554432)` – Specifies a hard maximum allowed
  request size, in bytes. Defaults to 32 MB. Specifying a number less than or
  equal to `0` turns off limiting altogether.

- `max_request_duration` `(string: "90s")` – Specifies the maximum
  request duration allowed before Vault cancels the request. This overrides
  `default_max_request_duration` for this listener.

- `proxy_protocol_behavior` `(string: "")` – When specified, enables a PROXY
  protocol version 1 behavior for the listener.
  Accepted Values:
  - *use_always* - The client's IP address will always be used.
  - *allow_authorized* - If the source IP address is in the
  `proxy_protocol_authorized_addrs` list, the client's IP address will be used.
  If the source IP is not in the list, the source IP address will be used.
  - *deny_unauthorized* - The traffic will be rejected if the source IP
  address is not in the `proxy_protocol_authorized_addrs` list.

- `proxy_protocol_authorized_addrs` `(string: <required-if-enabled> or array: <required-if-enabled> )` –
  Specifies the list of allowed source IP addresses to be used with the PROXY protocol.
  Not required if `proxy_protocol_behavior` is set to `use_always`. Source IPs should
  be comma-delimited if provided as a string. At least one source IP must be provided,
  `proxy_protocol_authorized_addrs` cannot be an empty array or string.

- `tls_disable` `(string: "false")` – Specifies if TLS will be disabled. Vault
  assumes TLS by default, so you must explicitly disable TLS to opt-in to
  insecure communication.

- `tls_cert_file` `(string: <required-if-enabled>, reloads-on-SIGHUP)` –
  Specifies the path to the certificate for TLS. To configure the listener to
  use a CA certificate, concatenate the primary certificate and the CA
  certificate together. The primary certificate should appear first in the
  combined file. On `SIGHUP`, the path set here _at Vault startup_ will be used
  for reloading the certificate; modifying this value while Vault is running
  will have no effect for `SIGHUP`s.

- `tls_key_file` `(string: <required-if-enabled>, reloads-on-SIGHUP)` –
  Specifies the path to the private key for the certificate. If the key file
  is encrypted, you will be prompted to enter the passphrase on server startup.
  The passphrase must stay the same between key files when reloading your
  configuration using `SIGHUP`. On `SIGHUP`, the path set here _at Vault
  startup_ will be used for reloading the certificate; modifying this value
  while Vault is running will have no effect for `SIGHUP`s.

- `tls_min_version` `(string: "tls12")` – Specifies the minimum supported
  version of TLS. Accepted values are "tls10", "tls11" or "tls12".

    ~> **Warning**: TLS 1.1 and lower are generally considered insecure.

- `tls_cipher_suites` `(string: "")` – Specifies the list of supported
  ciphersuites as a comma-separated-list. The list of all available ciphersuites
  is available in the [Golang TLS documentation][golang-tls].

- `tls_prefer_server_cipher_suites` `(string: "false")` – Specifies to prefer the
  server's ciphersuite over the client ciphersuites.

- `tls_require_and_verify_client_cert` `(string: "false")` – Turns on client
  authentication for this listener; the listener will require a presented
  client cert that successfully validates against system CAs.

- `tls_client_ca_file` `(string: "")` – PEM-encoded Certificate Authority file
  used for checking the authenticity of client.

- `tls_disable_client_certs` `(string: "false")` – Turns off client
  authentication for this listener. The default behavior (when this is false)
  is for Vault to request client certificates when available.

- `x_forwarded_for_authorized_addrs` `(string: <required-to-enable>)` –
  Specifies the list of source IP CIDRs for which an X-Forwarded-For header
  will be trusted. Comma-separated list or JSON array. This turns on
  X-Forwarded-For support.

- `x_forwarded_for_hop_skips` `(string: "0")` – The number of addresses that will be
  skipped from the *rear* of the set of hops. For instance, for a header value
  of `1.2.3.4, 2.3.4.5, 3.4.5.6`, if this value is set to `"1"`, the address that
  will be used as the originating client IP is `2.3.4.5`.

- `x_forwarded_for_reject_not_authorized` `(string: "true")` – If set false,
  if there is an X-Forwarded-For header in a connection from an unauthorized
  address, the header will be ignored and the client connection used as-is,
  rather than the client connection rejected.

- `x_forwarded_for_reject_not_present` `(string: "true")` – If set false, if
  there is no X-Forwarded-For header or it is empty, the client address will be
  used as-is, rather than the client connection rejected.

## `tcp` Listener Examples

### Configuring TLS

This example shows enabling a TLS listener.

```hcl
listener "tcp" {
  tls_cert_file = "/etc/certs/vault.crt"
  tls_key_file  = "/etc/certs/vault.key"
}
```

### Listening on Multiple Interfaces

This example shows Vault listening on a private interface, as well as localhost.

```hcl
listener "tcp" {
  address = "127.0.0.1:8200"
}

listener "tcp" {
  address = "10.0.0.5:8200"
}

# Advertise the non-loopback interface
api_addr = "https://10.0.0.5:8200"
cluster_addr = "https://10.0.0.5:8201"
```

[golang-tls]: https://golang.org/src/crypto/tls/cipher_suites.go
[api-addr]: /docs/configuration/index.html#api_addr
[cluster-addr]: /docs/configuration/index.html#cluster_addr
