---
layout: "docs"
page_title: "Custom - Database - Secrets Engines"
sidebar_title: "Custom"
sidebar_current: "docs-secrets-databases-custom"
description: |-
  The database secrets engine allows new functionality to be added through a
  plugin interface without needing to modify vault's core code. This allows you
  write your own code to generate credentials in any database you wish. It also
  allows databases that require dynamically linked libraries to be used as
  plugins while keeping Vault itself statically linked.
---

# Custom Database Secrets Engines

The database secrets engine allows new functionality to be added through a
plugin interface without needing to modify vault's core code. This allows you
write your own code to generate credentials in any database you wish. It also
allows databases that require dynamically linked libraries to be used as plugins
while keeping Vault itself statically linked.

~> **Advanced topic!** Plugin development is a highly advanced topic in Vault,
and is not required knowledge for day-to-day usage. If you don't plan on writing
any plugins, we recommend not reading this section of the documentation.

Please read the [Plugins internals](/docs/internals/plugins.html) docs for more
information about the plugin system before getting started building your
Database plugin.

## Plugin Interface

All plugins for the database secrets engine must implement the same simple interface.

```go
type Database interface {
	Type() (string, error)
	CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error)
	RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) error
	RevokeUser(ctx context.Context, statements Statements, username string) error
	RotateRootCredentials(ctx context.Context, statements []string) (config map[string]interface{}, err error)
	Init(ctx context.Context, config map[string]interface{}, verifyConnection bool) (saveConfig map[string]interface{}, err error)
	Close() error
}
```

You'll notice the first parameter to a number of those functions is a
`Statements` struct. This struct is used to pass the Role's configured
statements to the plugin on function call. The struct is defined as:

```go
type Statements struct {
	Creation   []string
	Revocation []string
	Rollback   []string
	Renewal    []string
}
```

It is up to your plugin to replace the `{{name}}`, `{{password}}`, and
`{{expiration}}` in these statements with the proper values.

The `Initialize` function is passed a map of keys to values, this data is what the
user specified as the configuration for the plugin. Your plugin should use this
data to make connections to the database. It is also passed a boolean value
specifying whether or not your plugin should return an error if it is unable to
connect to the database.

## Serving your plugin

Once your plugin is built you should pass it to vault's `plugins` package by
calling the `Serve` method:

```go
package main

import (
    "github.com/hashicorp/vault/plugins"
)

func main() {
    plugins.Serve(new(MyPlugin), nil)
}
```

Replacing `MyPlugin` with the actual implementation of your plugin.

The second parameter to `Serve` takes in an optional vault `api.TLSConfig` for
configuring the plugin to communicate with vault for the initial unwrap call.
This is useful if your vault setup requires client certificate checks. This
config wont be used once the plugin unwraps its own TLS cert and key.

## Running your plugin

The above main package, once built, will supply you with a binary of your
plugin. We also recommend if you are planning on distributing your plugin to
build with [gox](https://github.com/mitchellh/gox) for cross platform builds.

To use your plugin with the database secrets engine you need to place the binary in the
plugin directory as specified in the [plugin internals](/docs/internals/plugins.html) docs.

You should now be able to register your plugin into the vault catalog. To do
this your token will need sudo permissions.

```text
$ vault write sys/plugins/catalog/database/myplugin-database-plugin \
    sha256="..." \
    command="myplugin"
Success! Data written to: sys/plugins/catalog/database/myplugin-database-plugin
```

Now you should be able to configure your plugin like any other:

```text
$ vault write database/config/myplugin \
    plugin_name=myplugin-database-plugin \
    allowed_roles="readonly" \
    myplugins_connection_details="..."
```
