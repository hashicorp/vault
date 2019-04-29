---
layout: "docs"
page_title: "Vault Agent Auto-Auth Cert Method"
sidebar_title: "Cert"
sidebar_current: "docs-agent-autoauth-methods-cert"
description: |-
  Cert Method for Vault Agent Auto-Auth
---

# Vault Agent Cert AppRole Method

The `cert` method reads uses the configured TLS cerificates from the agent config
and takes an optional `name` parameters. There is no option to use certificates
which differ from those used in the `vault` block.
See [Agent Overview](https://vaultproject.io/docs/agent/index.html)

## Configuration

* `name` `(string: optional)` - The trusted cert `name` which should be used
  when authenticating with TLS.
