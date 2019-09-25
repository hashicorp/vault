---
layout: "docs"
page_title: "Vault Agent Auto-Auth Cert Method"
sidebar_title: "Cert"
sidebar_current: "docs-agent-autoauth-methods-cert"
description: |-
  Cert Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth Cert Method

The `cert` method uses the configured TLS certificates from the `vault` stanza of
the agent configuration and takes an optional `name` parameter. There is no option
to use certificates which differ from those used in the `vault` stanza.

See TLS settings in the [`vault` Stanza](https://vaultproject.io/docs/agent/index.html#vault-stanza)

## Configuration

* `name` `(string: optional)` - The trusted certificate role which should be used
  when authenticating with TLS. If a `name` is not specified, the auth method will
  try to authenticate against [all trusted certificates](https://www.vaultproject.io/docs/auth/cert.html#authentication).
