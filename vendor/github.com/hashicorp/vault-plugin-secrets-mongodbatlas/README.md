# HashiCorp Vault MongoDB Atlas Secrets Engine

The MongoDB Atlas Secrets Engine is a plugin for HashiCorp Vault which generates unique, ephemeral [Programmatic API](https://docs.atlas.mongodb.com/reference/api/apiKeys/) keys for MongoDB Atlas.

**The plugin is included in version 1.4 of Vault.**

## Support, Bugs and Feature Requests
Support for the HashiCorp Vault MongoDB Atlas Secrets Engine is provided under MongoDB Atlas support plans. Please submit support questions within the Atlas UI.  Vault support is via HashiCorp.

Bugs should be filed under the Issues section of this repo.

Feature requests can be submitted in the Issues section or directly with MongoDB at https://feedback.mongodb.com/forums/924145-atlas - just select the Vault plugin as the category or vote for an already suggested feature.

## Quick Links
- [MongoDB Atlas Secrets Engine - Docs](https://www.vaultproject.io/docs/secrets/mongodbatlas)
- [MongoDB Atlas Secrets Engine - API Docs](https://www.vaultproject.io/api-docs/secret/mongodbatlas/)
- [MongoDB Atlas Website](https://www.mongodb.com/cloud/atlas)
- [Vault Website](https://www.vaultproject.io)

**Please note**: Hashicorp takes Vault's security and their users' trust very seriously, as does MongoDB.

If you believe you have found a security issue in Vault or with this plugin, _please responsibly disclose_ by
contacting HashiCorp at [security@hashicorp.com](mailto:security@hashicorp.com) and contact MongoDB
directly via [security@mongodb.com](mailto:security@mongodb.com) or
[open a ticket](https://jira.mongodb.org/plugins/servlet/samlsso?redirectTo=%2Fbrowse%2FSECURITY) (link is external).

## Running tests

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
through the Organization access manager view on the web UI.
