# Vault Plugin: AliCloud Platform Secrets Backend

This is a backend plugin to be used with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin generates unique, ephemeral API keys and STS credentials.

**Please note**: We take Vault's security and our users' trust very seriously. 
If you believe you have found a security issue in Vault or with this plugin, 
_please responsibly disclose_ by 
contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links
- [Vault Website](https://www.vaultproject.io)
- [AliCloud Secrets Docs](https://developer.hashicorp.com/vault/docs/secrets/alicloud)
- [Vault Github](https://www.github.com/hashicorp/vault)
- [General Announcement List](https://groups.google.com/forum/#!forum/hashicorp-announce)
- [Discussion List](https://groups.google.com/forum/#!forum/vault-tool)


## Usage

This is a [Vault plugin](https://developer.hashicorp.com/vault/docs/plugins)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works. Otherwise, first read this guide on 
how to [get started with Vault](https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install).

If you are using Vault 11.0.1 or above, this plugin is packaged with Vault
and by default can be enabled by running:
 ```sh
 
 $ vault secrets enable alicloud
 
 Success! Enabled the alicloud secrets engine at: alicloud/
 
 ```
 
 If you are testing this plugin in an earlier version of Vault or 
 want to develop, see the next section. 

## Developing

If you wish to work on this plugin, you'll first need [Go](https://www.golang.org) 
installed on your machine (whichever version is required by Vault).

Make sure Go is properly installed, including setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).

### Get Plugin 
Clone this repository: 

```

mkdir $GOPATH/src/github.com/hashicorp/vault-plugin-secrets-alicloud`
cd $GOPATH/src/github.com/hashicorp/
git clone https://github.com/hashicorp/vault-plugin-secrets-alicloud.git

```
(or use `go get github.com/hashicorp/vault-plugin-secrets-alicloud` ).

You can then download any required build tools by bootstrapping your environment:

```sh
$ make bootstrap
```

To compile a development version of this plugin, run `make` or `make dev`.
This will put the plugin binary in the `bin` and `$GOPATH/bin` folders. `dev`
mode will only generate the binary for your platform and is faster:

```sh
$ make
$ make dev
```

### Install Plugin in Vault

Put the plugin binary into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://developer.hashicorp.com/vault/docs/configuration#plugin_directory)
in the Vault config used to start the server.

```hcl

plugin_directory = "path/to/plugin/directory"

```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.json ...
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://developer.hashicorp.com/vault/docs/plugins/plugin-architecture#plugin-catalog):

```sh
$ vault plugin register \
        -sha256="$(shasum -a 256 path/to/plugin/directory/vault-plugin-secrets-alicloud | cut -d " " -f1)" \
        -command="vault-plugin-secrets-alicloud" \
        secret \
        alicloudsecrets
```

Any name can be substituted for the plugin name "alicloudsecrets". This
name will be referenced in the next step, where we enable the secrets
plugin backend using the AliCloud secrets plugin:

```sh
$ vault secrets enable --plugin-name='alicloudsecrets' --path="alicloud" plugin

```

### Tests

This plugin has both integration tests and acceptance tests.

The integration tests fire API calls at a local test server that returns expected
responses rather than firing real API calls. They are executed by the following:

```sh
$ make test
```

The acceptance tests fire real API calls, and are located in `acceptance_test.go`.
These should be run once as a final step before placing a PR. The following environment
variables will need to be set for the acceptance tests to run:
* `VAULT_ACC=1`
* `VAULT_SECRETS_ACC_TEST_ROLE_ARN`
* `VAULT_SECRETS_ACC_TEST_ACCESS_KEY_ID`
* `VAULT_SECRETS_ACC_TEST_SECRET_KEY`

**Warning:** The acceptance tests create/destroy/modify *real resources*,
which may incur real costs in some cases. In the presence of a bug,
it is technically possible that broken backends could leave dangling
data behind. Therefore, please run the acceptance tests at your own risk.
At the very least, we recommend running them in their own private
account for whatever backend you're testing.

To run the acceptance tests, after exporting the necessary environment variables,
execute the following from the home directory:

```sh
$ make testacc
```

Or to execute only the acceptance tests:

```sh
./scripts/run_acceptance.sh
```

## Other Docs

See up-to-date [docs](https://developer.hashicorp.com/vault/docs/secrets/alicloud)
and general [API docs](https://developer.hashicorp.com/vault/api-docs/secret/alicloud).
