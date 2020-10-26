# Vault Plugins: MongoDB Atlas Secrets Engine and Database Secrets Engine for MongoDB Atlas plugin

**IMPORTANT: This plugin is currently under development.  Feel free to test it out following the instructions under the Developing section below, however consider this beta until it is verified by HashiCorp. Once verified and released versions will be documented in a CHANGELOG**

This contains two Secrets Engines specific to MongoDB Atlas for use with [Hashicorp Vault](https://github.com/hashicorp/vault).
The first is the MongoDB Atlas Secrets Engine which generates unique, ephemeral [Programmatic API](https://docs.atlas.mongodb.com/reference/api/apiKeys/) keys for MongoDB Atlas.
The second is an extension of the existing Database Secrets Engine and allows generation of unique, ephemeral
programmatic MongoDB [Database User](https://docs.atlas.mongodb.com/reference/api/database-users/) credentials in MongoDB Atlas, thus we refer to it as the Database Secrets
Engine for MongoDB Atlas.

The plugins are located in the following directories:
  - **MongoDB Atlas Secrets Engine:** `plugins/logical/mongodbatlas/`
  - **Database Secrets Engine for MongoDB Atlas plugin:** `plugins/database/mongodbatlas`

**Please note**: Hashicorp takes Vault's security and their users' trust very seriously, as does MongoDB.

If you believe you have found a security issue in Vault or with this plugin, _please responsibly disclose_ by
contacting us at [security@hashicorp.com](mailto:security@hashicorp.com) and contact MongoDB
directly via [security@mongodb.com](mailto:security@mongodb.com) or
[open a ticket](https://jira.mongodb.org/plugins/servlet/samlsso?redirectTo=%2Fbrowse%2FSECURITY) (link is external).

## Quick Links
- [Vault Website](https://www.vaultproject.io)
- [MongoDB Atlas Website](https://www.mongodb.com/cloud/atlas)
- [MongoDB Atlas Secrets Engine Docs](https://www.vaultproject.io/docs/secrets/mongodbatlas/index.html)
- [Database Secrets Engine for MongoDB Atlas](https://www.vaultproject.io/docs/secrets/databases/mongodbatlas.html)
- [Vault Github](https://www.github.com/hashicorp/vault)
- [Vault General Announcement List](https://groups.google.com/forum/#!forum/hashicorp-announce)
- [Vault Discussion List](https://groups.google.com/forum/#!forum/vault-tool)


## Usage

**The following will be accurate after review and approval by Hashicorp, which is in progress. Until then follow the instructions in the developing section that follows:**

These are a [Vault specific plugins (aka Secrets Engines/Backends)](https://www.vaultproject.io/docs/internals/plugins.html). This guide assumes you have already installed Vault
and have a basic understanding of how Vault works. Otherwise, first read this guide on
how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

If you are using Vault 11.0.1 or above, both plugins are packaged with Vault. The MongoDB Atlas Secrets Engine can be enabled by running:

The MongoDB Atlas Secrets Engine can be enabled by running:

 ```sh

 $ vault secrets enable mongodbatlas

 Success! Enabled the mongodbatlas secrets engine at: mongodbatlas/

 ```

 Then, write the configuration for the plugin and the lease, this is an example:

 ```sh

vault write mongodbatlas/config \
    public_key="a-public-key" \
    private_key="a-private-key"

vault write mongodbatlas/config/lease \
	ttl=300 \
	max_ttl=4800

 ```

The Database Secrets Engine for MongoDB Atlas can be enabled by running:

 ```sh

  $ vault secrets enable database

    Success! Enabled the database secrets engine at: database/

```

Then, write the configuration for the plugin, for example:

```sh
  $ vault write database/config/my-mongodbatlas-database \
      plugin_name=mongodbatlas-database-plugin \
      allowed_roles="my-role" \
      public_key="a-public-key" \
      private_key="a-private-key!" \
      project_id="a-project-id"

 ```

 If you are testing this plugin in an earlier version of Vault or
 want to develop, see the next section.

## Developing

If you wish to work on either plugin, you'll first need [Go](https://www.golang.org)
installed on your machine (whichever version is required by Vault).

Make sure Go is properly installed, including setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).

### Get Plugin

Clone this repository:

```

mkdir $GOPATH/src/github.com/hashicorp/vault-plugin-secrets-mongodbatlas`
cd $GOPATH/src/github.com/hashicorp/
git clone git@github.com:mongodb/vault-plugin-secrets-mongodbatlas.git
go mod download

```
(or use `go get github.com/mongodb/vault-plugin-secrets-mongodbatlas` ).

Then you can download any of the required tools to bootstrap your environment:

```sh
$ make bootstrap
```

To compile a development version of these plugins, run `make` or `make dev`.
This will put the plugin binaries in the `bin` and `$GOPATH/bin` folders. `dev`
mode will only generate binaries for your platform and is faster:

```sh
$ make
$ make dev
```

### Install Plugin in Vault

Put the plugin binaries into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://www.vaultproject.io/docs/configuration/index.html#plugin_directory)
in the Vault config used to start the server.

```hcl

plugin_directory = "path/to/plugin/directory"

```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.json ...
```

Once the server is started, register the plugins in the Vault server's [plugin catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog):

#### MongoDB Atlas Secrets Engine

To register the MongoDB Atlas Secrets Engine run the following:

```sh
$ vault write sys/plugins/catalog/vault-plugin-secrets-mongodbatlas \
        sha_256="$(shasum -a 256 path/to/plugin/directory/vault-plugin-secrets-mongodbatlas | cut -d " " -f1)" \
        command="vault-plugin-secrets-mongodbatlas"
```

Any name can be substituted for the plugin name "vault-plugin-secrets-mongodbatlas". This
name will be referenced in the next step, where we enable the secrets
plugin backend using the MongoDB Atlas Secrets Engine:

```sh
$ vault secrets enable --plugin-name='vault-plugin-secrets-mongodbatlas' --path="vault-plugin-secrets-mongodbatlas" plugin

```

#### Database Secrets Engine for MongoDB Atlas plugin

The following steps are required to register the Database Secrets Engine for MongoDB Atlas plugin:

```sh

vault write sys/plugins/catalog/database/mongodbatlas-database-plugin \
    sha256=$(shasum -a 256 mongodbatlas-database-plugin | cut -d' ' -f1) \
    command="mongodbatlas-database-plugin"

```

Then, you must enable the Vault's Database Secret Engine with Vault

```sh

vault secrets enable database

```

### Tests

This plugin has both integration tests, and acceptance tests.

The integration tests are run by `$ make test` and rather than firing real
API calls, they fire API calls at a local test server that returns expected
responses.

The acceptance tests fire real API calls, and are located in `plugins/logical/mongodbatlas/acceptance_test.go`
and `plugins/database/mongodbatlas/mongodbatlas_test.go`. These should be run
once as a final step before placing a PR. Please see `acceptance_test.go` and
`mongodbatlas_test.go` to learn the environment variables that will need to be set.

**Warning:** The acceptance tests create/destroy/modify *real resources*,
which may incur real costs in some cases. In the presence of a bug,
it is technically possible that broken backends could leave dangling
data behind. Therefore, please run the acceptance tests at your own risk.
At the very least, we recommend running them in their own private
account for whatever backend you're testing.

Before running the acceptance tests export the following environment variables:

- VAULT_ACC - Set to `1` to run the acceptance tests
- ATLAS_ORGANIZATION_ID - Your Organization ID
- ATLAS_PUBLIC_KEY and ATLAS_PRIVATE_KEY - Your Public and Private key with the correct permissions to run the tests
- ATLAS_PROJECT_ID - Your Project ID

To run the acceptance tests, after exporting the necessary environment variables,
from the home directory run `VAULT_ACC=1 make test`:

```sh
$ VAULT_ACC=1 make test
```

## Other Docs

**The following will be accurate after review and approval by Hashicorp, which is in progress. Until then read the docs within this repo for more information.**

See up-to-date **MongoDB Atlas Secrets Engine** [docs](https://www.vaultproject.io/docs/secrets/mongodbatlas/index.html),
 **Database Secrets Engine for MongoDB Atlas plugin** [docs](https://www.vaultproject.io/docs/secrets/databases/mongodbatlas.html)
and general [API docs](https://www.vaultproject.io/api/secret/mongodbatlas/index.html).
