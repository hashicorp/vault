---
layout: "docs"
page_title: "Token - Auth Methods"
sidebar_current: "docs-auth-token"
description: |-
  The token store auth method is used to authenticate using tokens.
---

# Token Auth Method

The `token` method is built-in and automatically available at `/auth/token`. It
allows users to authenticate using a token, as well to create new tokens, revoke
secrets by token, and more.

When any other auth method returns an identity, Vault core invokes the
token method to create a new unique token for that identity.

The token store can also be used to bypass any other auth method:
you can create tokens directly, as well as perform a variety of other
operations on tokens such as renewal and revocation.

Please see the [token concepts](/docs/concepts/tokens.html) page dedicated
to tokens.

## Authentication

### Via the CLI

```text
$ vault login token=<token>
```

### Via the API

The token is set directly as a header for the HTTP API. The name
of the header should be "X-Vault-Token" and the value should be the token.

## API

The Token auth method has a full HTTP API. Please see the
[Token auth method API](/api/auth/token/index.html) for more
details.
