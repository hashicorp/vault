# Vault Plugin: Centrify Identity Platform Auth Backend

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin allows for Centrify Identity Platform users accounts to authenticate with Vault.

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links
    - Vault Website: https://www.vaultproject.io
    - Main Project Github: https://www.github.com/hashicorp/vault

## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Security Model

The current authentication model requires providing Vault with an OAuth2 Client ID and Secret, which can be used to make authenticated calls to the Centrify Identity Platform API.  This token is scoped to allow only the required APIs for Vault integration, and cannot be used for interactive login directly. 

## Usage

This plugin is currently built into Vault and by default is accessed
at `auth/centrify`. To enable this in a running Vault server:

```sh
$ vault auth-enable centrify
Successfully enabled 'centrify' at 'centrify'!
```

Before the plugin can authenticate users, both the plugin and your cloud service tenant must be configured correctly.  To configure your cloud tenant, sign in as an administrator and perform the following actions.  Please note that this plugin requires the Centrify Cloud Identity Service version 17.11 or newer.

### Create an OAuth2 Confidential Client

An OAuth2 Confidentical Client is a Centrify Directory User.

- Users -> Add User
  - Login Name: vault_integration@<yoursuffix>
  - Display Name: Vault Integration Confidential Client
  - Check the "Is OAuth confidentical client" box
  - Password Type: Generated (be sure to copy the value, you will need it later)
  - Create User

### Create a Role

To scope the users who can authenticate to vault, and to allow our Confidential Client access, we will create a role.

- Roles -> Add Role
  - Name: Vault Integration
  - Members -> Add
    - Search for and add the vault_integration@<yoursuffix> user
    - Additionally add any roles/groups/users who should be able to authenticate to vault
  - Save

### Create an OAuth2 Client Application
- Apps -> Add Web Apps -> Custom -> OAuth2 Client
- Configure the added application
  - Description:
    - Application ID: "vault_io_integration" 
    - Application Name: "Vault Integration"
  - General Usage:
    - Client ID Type -> Confidential (must be OAuth client)
  - Tokens:
    - Token Type: JwtRS256
    - Auth methods: Client Creds + Resource Owner    
  - Scope
    - Add a single scope named "vault_io_integration" with the following regexes:
      - usermgmt/getusersrolesandadministrativerights
      - security/whoami
  - User Access
    - Add the previously created "Vault Integration" role    
  - Save

### Configuring the Vault Plugin

As an administrative vault user, you can read/write the centrify plugin configuration using the /auth/centrify/config path:

```sh
$ vault write auth/centrify/config service_url=https://<tenantid>.my.centrify.com client_id=vault_integration@<yoursuffix> client_secret=<password copied earlier> app_id=vault_io_integration scope=vault_io_integration
```

### Authenticating

As a valid user of your tenant, in the appropriate role for accessing the Vault Integration app, you can now authenticate to the vault:

```sh
$ vault auth -method=centrify username=<your username>
```

Your vault token will be valid for the length of time defined in the app's token lifetime configuration (default 5 hours).

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine
(version 1.9+ is *required*).

For local dev first make sure Go is properly installed, including
setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).
Next, clone this repository into
`$GOPATH/src/github.com/hashicorp/vault-plugin-auth-centrify`.
You can then download any required build tools by bootstrapping your
environment:

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

Put the plugin binary into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://www.vaultproject.io/docs/configuration/index.html#plugin_directory)
in the Vault config used to start the server.

```json
...
plugin_directory = "path/to/plugin/directory"
...
```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.json ...
...
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog):

```sh
$ vault write sys/plugins/catalog/centrify \
        sha_256=<expected SHA256 Hex value of the plugin binary> \
        command="vault-plugin-auth-centrify"
...
Success! Data written to: sys/plugins/catalog/centrify
```

Note you should generate a new sha256 checksum if you have made changes
to the plugin. Example using openssl:

```sh
openssl dgst -sha256 $GOPATH/vault-plugin-auth-centrify
...
SHA256(.../go/bin/vault-plugin-auth-centrify)= 896c13c0f5305daed381952a128322e02bc28a57d0c862a78cbc2ea66e8c6fa1
```

Enable the auth plugin backend using the Centrify auth plugin:

```sh
$ vault auth-enable -plugin-name='centrify' plugin
...

Successfully enabled 'plugin' at 'centrify'!
```
