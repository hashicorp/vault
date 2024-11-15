# Combined Database Engine
This package is how database plugins interact with Vault.

## Upgrading to Version 5

### Background
In Vault 1.6, a new Database interface was created that solved a number of issues with the
previous interface:

1. It could not use password policies because the database plugins were responsible for
   generating passwords.
2. There were significant inconsistencies between functions in the interface.
3. Several functions (`SetCredentials` and `RotateRootCredentials`) were doing the same operation.
4. It had a function that was no longer being used as it had been deprecated in a previous version
   but never removed.

Prior to Vault 1.6, the Database interface is version 4 (with other versions in older versions of
Vault). The new version introduced in Vault 1.6 is version 5. This distinction was not exposed in
previous iterations of the Database interface as the previous versions were additive to the
interface. Since version 5 is an overhaul of the interface, this distinction needed to be made.

We highly recommend that you upgrade any version 4 database plugins to version 5 as version 4 is
considered deprecated and support for it will be removed in a future release. Version 5 plugins
will not function with Vault prior to Vault 1.6.

The new interface is roughly modeled after a [gRPC](https://grpc.io/) interface. It has improved
future compatibility by not requiring changes to the interface definition to add additional data
in the requests or responses. It also simplifies the interface by merging several into a single
function call.

### Upgrading your custom database

Vault 1.6 supports both version 4 and version 5 database plugins. The support for version 4
plugins will be removed in a future release. Version 5 database plugins will not function with
Vault prior to version 1.6. If you upgrade your database plugins, ensure that you are only using
Vault 1.6 or later. To determine if a plugin is using version 4 or version 5, the following is a
list of changes in no particular order that you can check against your plugin to determine
the version:

1. The import path for version 4 is `github.com/hashicorp/vault/sdk/database/dbplugin`
   whereas the import path for version 5 is `github.com/hashicorp/vault/sdk/database/dbplugin/v5`
2. Version 4 has the following functions: `Initialize`, `Init`, `CreateUser`, `RenewUser`,
   `RevokeUser`, `SetCredentials`, `RotateRootCredentials`, `Type`, and `Close`. You can see the
   full function signatures in `sdk/database/dbplugin/plugin.go`.
3. Version 5 has the following functions: `Initialize`, `NewUser`, `UpdateUser`, `DeleteUser`,
   `Type`, and `Close`. You can see the full function signatures in
   `sdk/database/dbplugin/v5/database.go`.

If you are using a version 4 custom database plugin, the following are basic instructions
for upgrading to version 5.

-> In version 4, password generation was the responsibility of the plugin. This is no longer
   the case with version 5. Vault is responsible for generating passwords and passing them to
   the plugin via `NewUserRequest.Password` and `UpdateUserRequest.Password.NewPassword`.

1. Change the import path from `github.com/hashicorp/vault/sdk/database/dbplugin` to
   `github.com/hashicorp/vault/sdk/database/dbplugin/v5`. The package name is the same, so any
   references to `dbplugin` can remain as long as those symbols exist within the new package
   (such as the `Serve` function).
2. An easy way to see what functions need to be implemented is to put the following as a
   global variable within your package: `var _ dbplugin.Database = (*MyDatabase)(nil)`. This
   will fail to compile if the `MyDatabase` type does not adhere to the
   `dbplugin.Database` interface.
3. Replace `Init` and `Initialize` with the new `Initialize` function definition. The fields that
   `Init` was taking (`config` and `verifyConnection`) are now wrapped into `InitializeRequest`.
   The returned `map[string]interface{}` object is now wrapped into `InitializeResponse`.
   Only `Initialize` is needed to adhere to the `Database` interface.
4. Update `CreateUser` to `NewUser`. The `NewUserRequest` object contains the username and
   password of the user to be created. It also includes a list of statements for creating the
   user as well as several other fields that may or may not be applicable. Your custom plugin
   should use the password provided in the request, not generate one. If you generate a password
   instead, Vault will not know about it and will give the caller the wrong password.
5. `SetCredentials`, `RotateRootCredentials`, and `RenewUser` are combined into `UpdateUser`.
   The request object, `UpdateUserRequest` contains three parts: the username to change, a
   `ChangePassword` and a `ChangeExpiration` object. When one of the objects is not nil, this
   indicates that particular field (password or expiration) needs to change. For instance, if
   the `ChangePassword` field is not-nil, the user's password should be changed. This is
   equivalent to calling `SetCredentials`. If the `ChangeExpiration` field is not-nil, the
   user's expiration date should be changed. This is equivalent to calling `RenewUser`.
   Many databases don't need to do anything with the updated expiration.
6. Update `RevokeUser` to `DeleteUser`. This is the simplest change. The username to be
   deleted is enclosed in the `DeleteUserRequest` object.

