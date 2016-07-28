---
layout: "docs"
page_title: "Auth Backends"
sidebar_current: "docs-auth"
description: |-
  Auth backends are mountable backends that perform authentication for Vault.
---

# Auth Backends

Auth backends are the components in Vault that perform authentication
and are responsible for assigning identity and a set of policies to a
user.

Having multiple auth backends enables you to use an auth backend
that makes the sense for your use case of Vault and your organization.

For example, on developer machines, the [GitHub auth backend](/docs/auth/github.html)
is easiest to use. But for servers the [AppRole](/docs/auth/approle.html)
backend is the recommended choice.

To learn more about authentication, see the
[authentication concepts page](/docs/concepts/auth.html).

## Enabling/Disabling Auth Backends

Auth backends can be enabled/disabled using the CLI or the API.

When enabled, auth backends are similar to [secret backends](/docs/secrets/index.html):
they are mounted within the Vault mount table and can be accessed
and configured using the standard read/write API. The only difference
is that all auth backends are mounted underneath the `auth/` prefix.

By default, auth backends are mounted to `auth/<type>`. For example,
if you enable "github", then you can interact with it at `auth/github`.
However, this path is customizable, allowing users with advanced use
cases to mount a single auth backend multiple times.

When an auth backend is disabled, all users authenticated via that
backend are automatically logged out.
