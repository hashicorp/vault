# HashiCorp Vault Database Secrets Engine - MongoDB Atlas plugin

MongoDB Atlas is one of the supported plugins for the HashiCorp Vault Database Secrets Engine and allows for the programmatic generation of unique, ephemeral MongoDB [Database User](https://docs.atlas.mongodb.com/reference/api/database-users/) credentials in MongoDB Atlas Projects.

**The plugin is included in version 1.4 of Vault.**

## Support, Bugs and Feature Requests
Support for the HashiCorp Vault Database Secrets Engine - MongoDB Atlas is provided under MongoDB Atlas support plans. Please submit support questions within the Atlas UI.  Vault support is via HashiCorp.

Bugs should be filed under the Issues section of this repo.

Feature requests can be submitted in the Issues section or directly to MongoDB at https://feedback.mongodb.com/forums/924145-atlas - just select the Vault plugin as the category or vote for an already suggested feature.

## Quick Links
- [Database Secrets Engine for MongoDB Atlas - Docs](https://www.vaultproject.io/docs/secrets/databases/mongodbatlas)
- [Database Secrets Engine for MongoDB Atlas - API Docs](https://www.vaultproject.io/api-docs/secret/databases/mongodbatlas/)
- [MongoDB Atlas Website](https://www.mongodb.com/cloud/atlas)
- [Vault Website](https://www.vaultproject.io)

**Please note**: Hashicorp takes Vault's security and their users' trust very seriously, as does MongoDB.

If you believe you have found a security issue in Vault or with this plugin, _please responsibly disclose_ by
contacting HashiCorp at [security@hashicorp.com](mailto:security@hashicorp.com) and contact MongoDB
directly via [security@mongodb.com](mailto:security@mongodb.com) or
[open a ticket](https://jira.mongodb.org/plugins/servlet/samlsso?redirectTo=%2Fbrowse%2FSECURITY) (link is external).

## Acceptance Testing

In order to perform acceptance testing, you need to set the environment
variable `VAULT_ACC=1` as well as provide all the of necessary information to
connect to a MongoDB Atlas Project. All `ATLAS_*` environment variables must be
provided in order for the acceptance tests to run properly. A cluster must be
available during the test. A 
[free tier cluster](https://docs.atlas.mongodb.com/tutorial/deploy-free-tier-cluster/) 
can be provisioned manually to test.

| Environment variable | Description                                                   |
|----------------------|---------------------------------------------------------------|
| ATLAS_PUBLIC_KEY     | The Atlas API public key                                      |
| ATLAS_PRIVATE_KEY    | The Atlas API private key                                     |
| ATLAS_PROJECT_ID     | The desired project ID or group ID                            |
| ATLAS_CONN_URL       | The desired cluster's connection URL within the project       |
| ATLAS_ALLOWLIST_IP   | The public IP of the machine that the test is being performed |