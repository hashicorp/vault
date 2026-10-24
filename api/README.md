Vault API
=================

This provides the `github.com/hashicorp/vault/api` package which contains code useful for interacting with a Vault server.

For examples of how to use this module, see the [vault-examples](https://github.com/hashicorp/vault-examples) repo.
For a step-by-step walkthrough on using these client libraries, see the [developer quickstart](https://developer.hashicorp.com/vault/docs/get-started/developer-qs).

[![GoDoc](https://godoc.org/github.com/hashicorp/vault/api?status.png)](https://godoc.org/github.com/hashicorp/vault/api)

## Checking the Vault server certificate against a CRL

By default, as with Go's `crypto/tls`, the client does not check whether the
Vault server's certificate has been revoked. To enable revocation checking,
supply a certificate revocation list through `TLSConfig`:

```go
config := api.DefaultConfig()
err := config.ConfigureTLS(&api.TLSConfig{
    CACert: "/etc/vault/ca.pem",
    CRL:    "/etc/vault/vault.crl", // or CRLBytes, or CRLPath for a directory
})
```

The equivalent environment variables are `VAULT_CRL`, `VAULT_CRL_BYTES` and
`VAULT_CRL_PATH`, and the CLI accepts `-crl` and `-crl-path`.

CRLs given as `CRL` or `CRLPath` are re-read when they change on disk, so an
externally refreshed CRL is picked up without restarting the process. Files are
checked at most once per `CRLRefreshInterval` (default 5 minutes) and only while
connections are being made; there is no background goroutine, and no disk I/O is
performed while other handshakes wait. If a refresh fails, the last successfully
loaded CRLs continue to be used and the refresh is retried shortly afterwards.

`CRLPath` is read recursively, matching `CAPath`.

To source revocation data from somewhere other than the filesystem, implement
`CRLProvider` and set it on `TLSConfig`. `CRLs` is called during the TLS
handshake, so implementations should return cached data.

### Turning revocation checking off

TLS settings can reach a client through more than one `ConfigureTLS` call (the
CLI, for instance, applies environment settings and then flag settings), so a
`TLSConfig` that simply does not mention CRLs leaves any existing revocation
checking in place rather than silently removing it. To turn it off deliberately,
set `DisableCRL`:

```go
err := config.ConfigureTLS(&api.TLSConfig{DisableCRL: true})
```

### What is and is not checked

Revocation checking is strict about the server certificate and lenient about the
rest of the chain:

- The Vault server's certificate **must** be covered by a valid, current CRL
  from its issuer. If no such CRL is available, the connection fails, rather
  than silently succeeding without a revocation check.
- Certificates further up the chain are checked when a CRL for their issuer is
  available, and skipped when it is not. This means supplying a CRL for one CA
  in a chain does not break every connection.
- A CRL that is past its `NextUpdate` does not count as coverage, because it
  says nothing about revocations issued after it was published. Its existing
  revocation entries are still honored, since revocation is permanent. A CRL
  with no `NextUpdate` at all, which X.509 permits, counts as
  coverage indefinitely.
- A delta CRL, which lists only the changes since a base CRL, does not count as
  coverage on its own, because it cannot show that a certificate is unrevoked.
  Its revocation entries are honored.
- A CRL whose issuer name matches but whose signature does not verify, which is
  what a stale CRL from before a CA key rotation looks like, is ignored rather
  than treated as an error.
- A self-signed server certificate is exempt from the coverage requirement,
  since it is its own issuer and requiring a CRL from itself would make
  revocation checking and self-signed certificates mutually exclusive. A CRL it
  does publish is still honored, so a self-signed CA can revoke its own
  certificate.

Revocation checking cannot be combined with `Insecure`: with chain verification
disabled there is no verified chain to check a CRL against, so `ConfigureTLS`
returns an error instead of appearing to check something it cannot. If
`InsecureSkipVerify` reaches the underlying `tls.Config` by another route, the
check fails the connection rather than validating an unverified chain.

### TLS session resumption

The check is installed as [`tls.Config.VerifyConnection`][verifyconnection]
rather than `VerifyPeerCertificate`, because the latter is not invoked on
resumed TLS connections. Using `VerifyConnection` means the revocation check
runs on every connection, including resumed ones, so it cannot be bypassed by
session resumption and there is no need to disable session tickets.

This matters only for callers who supply their own `http.Client`: Go disables
client-side resumption when `tls.Config.ClientSessionCache` is nil, and
`DefaultConfig` never sets one.

If you set your own `VerifyConnection` before calling `ConfigureTLS`, it is
preserved and runs before the CRL check.

[verifyconnection]: https://pkg.go.dev/crypto/tls#Config.VerifyConnection
