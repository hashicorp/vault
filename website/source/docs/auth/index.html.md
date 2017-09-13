---
layout: "docs"
page_title: "Auth Methods"
sidebar_current: "docs-auth"
description: |-
  Auth methods are mountable methods that perform authentication for Vault.
---

# Auth Methods

Auth methods are the components in Vault that perform authentication and are
responsible for assigning identity and a set of policies to a user.

Having multiple auth methods enables you to use an auth method that makes the
sense for your use case of Vault and your organization.

For example, on developer machines, the [GitHub auth method](/docs/auth/github.html)
is easiest to use. But for servers the [AppRole](/docs/auth/approle.html)
method is the recommended choice.

To learn more about authentication, see the
[authentication concepts page](/docs/concepts/auth.html).

## Enabling/Disabling Auth Methods

Auth methods can be enabled/disabled using the CLI or the API.

```text
$ vault auth enable userpass
```

When enabled, auth methods are similar to [secrets engines](/docs/secrets/index.html):
they are mounted within the Vault mount table and can be accessed
and configured using the standard read/write API. All auth methods are mounted underneath the `auth/` prefix.

By default, auth methods are mounted to `auth/<type>`. For example, if you
enable "github", then you can interact with it at `auth/github`. However, this
path is customizable, allowing users with advanced use cases to mount a single
auth method multiple times.

```text
$ vault auth enable -path=my-login userpass
```

When an auth method is disabled, all users authenticated via that method are
automatically logged out.
