---
layout: "docs"
page_title: "auth disable - Command"
sidebar_current: "docs-commands-auth-disable"
description: |-
  The "auth disable" command disables an auth method at a given path, if one
  exists. This command is idempotent, meaning it succeeds even if no auth method
  is enabled at the path.
---

# auth disable

The `auth disable` command disables an auth method at a given path, if one
exists. This command is idempotent, meaning it succeeds even if no auth method
is enabled at the path.

Once an auth method is disabled, it can no longer be used for authentication.
**All access tokens generated via the disabled auth method are immediately
revoked.** This command will block until all tokens are revoked.

## Examples

Disable the auth method enabled at "userpass/":

```text
$ vault auth disable userpass/
Success! Disabled the auth method (if it existed) at: userpass/
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
