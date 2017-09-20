---
layout: "docs"
page_title: "path-help - Command"
sidebar_current: "docs-commands-path-help"
description: |-
  The "path-help" command retrieves API help for paths. All endpoints in Vault
  provide built-in help in markdown format. This includes system paths, secret
  engines, and auth methods.
---

# path-help

The `path-help` command retrieves API help for paths. All endpoints in Vault
provide built-in help in markdown format. This includes system paths, secret
engines, and auth methods.

The help system is the easiest way to learn how to use the various systems
in Vault, and also allows you to discover new paths.

Before using `path-help`, it is important to understand "paths" within Vault.
Paths are the parameters used for `vault read`, `vault write`, etc. An example
path is `secret/foo`, or `aws/config/root`. The paths available depend on the
secrets engines in use. Because of this, the interactive help is an
indispensable tool to finding what paths are supported.

To discover what paths are supported, use `vault path-help PATH`. For example,
if you enabled the AWS secrets engine, you can use `vault path-help aws` to find
the paths supported by that backend. The paths are shown with regular
expressions, which can make them hard to parse, but they are also extremely
exact.

## Examples

Get help output for the KV secrets engine:

```text
$ vault path-help secret
## DESCRIPTION

The KV backend reads and writes arbitrary secrets to the backend.
The secrets are encrypted/decrypted by Vault: they are never stored
unencrypted in the backend and the backend never has an opportunity to
see the unencrypted value.

Leases can be set on a per-secret basis. These leases will be sent down
when that secret is read, and it is assumed that some outside process will
revoke and/or replace the secret at that path.

## PATHS

The following paths are supported by this backend. To view help for
any of the paths below, use the help command with any route matching
the path pattern. Note that depending on the policy of your auth token,
you may or may not be able to access certain paths.

    ^.*$
        Pass-through secret storage to the storage backend, allowing you to
        read/write arbitrary data into secret storage.
```

Once you've found a path you like, you can learn more about it by using `vault
path-help <path>` where "path" is a path that matches one of the regular
expressions from the backend help.

```text
$ vault path-help secret/password
Request:        password
Matching Route: ^.*$

Pass-through secret storage to the storage backend, allowing you to
read/write arbitrary data into secret storage.

## PARAMETERS

    lease (string)
        Lease time for this key when read. Ex: 1h

## DESCRIPTION

The pass-through backend reads and writes arbitrary data into secret storage,
encrypting it along the way.

A lease can be specified when writing with the "lease" field. If given, then
when the secret is read, Vault will report a lease with that duration. It
is expected that the consumer of this backend properly writes renewed keys
before the lease is up. In addition, revocation must be handled by the
user of this backend.
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
