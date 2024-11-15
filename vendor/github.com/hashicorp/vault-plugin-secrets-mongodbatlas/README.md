# HashiCorp Vault MongoDB Atlas Secrets Engine

The MongoDB Atlas Secrets Engine is a plugin for HashiCorp Vault which generates unique, ephemeral [Programmatic API](https://docs.atlas.mongodb.com/reference/api/apiKeys/) keys for MongoDB Atlas.

**The plugin is included in version 1.4 of Vault.**

## Support, Bugs and Feature Requests
Support for the HashiCorp Vault MongoDB Atlas Secrets Engine is provided under MongoDB Atlas support plans. Please submit support questions within the Atlas UI.  Vault support is via HashiCorp.

Bugs should be filed under the Issues section of this repo.

Feature requests can be submitted in the Issues section or directly with MongoDB at https://feedback.mongodb.com/forums/924145-atlas - just select the Vault plugin as the category or vote for an already suggested feature.

## Quick Links
- [MongoDB Atlas Secrets Engine - Docs](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas)
- [MongoDB Atlas Secrets Engine - API Docs](https://developer.hashicorp.com/vault/api-docs/secret/mongodbatlas)
- [MongoDB Atlas Website](https://www.mongodb.com/cloud/atlas)
- [Vault Website](https://www.vaultproject.io)

**Please note**: Hashicorp takes Vault's security and their users' trust very seriously, as does MongoDB.

If you believe you have found a security issue in Vault or with this plugin, _please responsibly disclose_ by
contacting HashiCorp at [security@hashicorp.com](mailto:security@hashicorp.com) and contact MongoDB
directly via [security@mongodb.com](mailto:security@mongodb.com) or
[open a ticket](https://jira.mongodb.org/plugins/servlet/samlsso?redirectTo=%2Fbrowse%2FSECURITY) (link is external).

## Usage

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works. Otherwise, first read this guide on
how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

If you are just interested in using this plugin with Vault, it is packaged with Vault and
by default can be enabled by running:

```sh
$ vault secrets enable mongodbatlas
Success! Enabled the mongodbatlas secrets engine at: mongodbatlas/
```

## Developing

If you wish to work on this plugin, you'll first need [Go](https://www.golang.org)
installed on your machine (whichever version is required by Vault).

Make sure Go is properly installed, including setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).

### Build Plugin
```sh
make dev
```

## Running tests

### Unit tests
```sh
make test
```

### Acceptance tests
```sh
make testacc
```

To run the acceptance tests, you need to set the following environment variables:

```bash
VAULT_ACC=1
ATLAS_PRIVATE_KEY=...
ATLAS_PUBLIC_KEY=...
ATLAS_PROJECT_ID=...
ATLAS_ORGANIZATION_ID=...
```

The programmatic API key provided must be an "Organization Owner", and must have
your public IP address set as an allowed address. You can manage both of these
through the Organization access manager view on the web UI. See
https://www.mongodb.com/docs/atlas/configure-api-access for details.
