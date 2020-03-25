# Vault Plugin: OpenLDAP Secrets Backend

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin provides OpenLDAP functionality to Vault.

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links
- Vault Website: https://www.vaultproject.io
- OpenLDAP Docs: https://www.vaultproject.io/docs/secrets/openldap/index.html
- Main Project Github: https://www.github.com/hashicorp/vault

## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Usage

Please see [documentation for the plugin](https://www.vaultproject.io/docs/secrets/openldap/index.html)
on the Vault website.

This plugin is currently built into Vault and by default is accessed
at `openldap`. To enable this in a running Vault server:

```sh
$ vault secrets enable openldap
Success! Enabled the openldap secrets engine at: openldap/
```
