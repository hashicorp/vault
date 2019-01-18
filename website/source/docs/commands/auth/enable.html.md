---
layout: "docs"
page_title: "auth enable - Command"
sidebar_title: "<code>enable</code>"
sidebar_current: "docs-commands-auth-enable"
description: |-
  The "auth enable" command enables an auth method at a given path. If an auth
  method already exists at the given path, an error is returned. After the auth
  method is enabled, it usually needs configuration.
---

# auth enable

The `auth enable` command enables an auth method at a given path. If an auth
method already exists at the given path, an error is returned. After the auth
method is enabled, it usually needs configuration. The configuration varies by
auth method.

An auth method is responsible for authenticating users or machines and assigning
them policies and a token with which they can access Vault. Authentication is
usually mapped to policy. Please see the [policies
concepts](/docs/concepts/policies.html) page for more information.

## Examples

Enable the auth method "userpass" enabled at "userpass/":

```text
$ vault auth enable userpass
Success! Enabled the userpass auth method at: userpass/
```

Create a user:

```text
$ vault write auth/userpass/users/sethvargo password=secret
Success! Data written to: auth/userpass/users/sethvargo
```

For more information on the specific configuration options and paths, please see
the [auth method](/docs/auth/index.html) documentation.

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-description` `(string: "")` - Human-friendly description for the purpose of
  this auth method.

- `-local` `(bool: false)` - Mark the auth method as local-only. Local auth
  methods are not replicated nor removed by replication.

- `-path` `(string: "")` - Place where the auth method will be accessible. This
  must be unique across all auth methods. This defaults to the "type" of the
  auth method. The auth method will be accessible at `/auth/<path>`.
