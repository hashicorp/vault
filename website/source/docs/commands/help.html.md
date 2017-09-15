---
layout: "docs"
page_title: "Path Help"
sidebar_current: "docs-commands-path-help"
description: |-
  The Vault CLI has a built-in help system that can be used to get help for not only the CLI itself, but also any paths that the CLI can be used with within Vault.
---

# Help

In addition to standard CLI help using the `-h` or `-help` flag for
commands, Vault has a built-in `path-help` command that can be used to get
help for specific paths within Vault. These paths are used with the
API or `read, write, delete` commands in order to interact with Vault.

The help system is the easiest way to learn how to use the various systems
in Vault, and also allows you to discover new paths.

-> **Important!** The help system is incredibly important in day-to-day
use of Vault. As a beginner or experienced user of Vault, you'll be using
the help command a lot to remember how to use different components of
Vault. Note that the Vault Server must be running and the client configured
properly to execute this command to look up paths.

## Discovering Paths

Before using `path-help`, it is important to understand "paths" within Vault.
Paths are the parameters used for `vault read`, `vault write`, etc. An
example path is `secret/foo`, or `aws/config/root`. The paths available
depend on the mounted secret backends. Because of this, the interactive
help is an indispensable tool to finding what paths are supported.

To discover what paths are supported, use `vault path-help <mount point>`.
For example, if you mounted the AWS secret backend, you can use
`vault path-help aws` to find the paths supported by that backend. The paths
will be shown with regular expressions, which can make them hard to
parse, but they're also extremely exact.

You can try it right away with any Vault with `vault path-help secret`, since
`secret` is always mounted initially. The output from this command is shown
below and contains both a description of what that backend is for, along with
the paths it supports.

```
$ vault path-help secret
## DESCRIPTION

The key/value backend reads and writes arbitrary secrets to the backend.
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

## Single Path

Once you've found a path you like, you can learn more about it by
using `vault path-help <path>` where "path" is a path that matches one of the
regular expressions from the backend help.

Or, if you saw an example online with `vault write` or some similar
command, you can plug that directly into `vault path-help` to learn about it
(assuming you have the proper backends mounted!).

For example, below we get the help for a single secret in the `secret/`
mount point. The help shows the operations that that path supports, the
parameters it takes (for write), and a description of that specific path.

```
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
