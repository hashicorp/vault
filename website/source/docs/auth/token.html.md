---
layout: "docs"
page_title: "Auth Backend: Token"
sidebar_current: "docs-auth-token"
description: |-
  The token store auth backend is used to authenticate using tokens.
---

# Auth Backend: Token

The token backend is the only auth backend that is built-in and
automatically available at `/auth/token` as well as with first-class
built-in CLI methods such as `vault token-create`. It allows users to
authenticate using a token, as well to create new tokens, revoke
secrets by token, and more.

When any other auth backend returns an identity, Vault core invokes the
token backend to create a new unique token for that identity.

The token store can also be used to bypass any other auth backend:
you can create tokens directly, as well as perform a variety of other
operations on tokens such as renewal and revocation.

Please see the [token concepts](/docs/concepts/tokens.html) page dedicated
to tokens.

## Authentication

### Via the CLI

```
$ vault auth <token>
...
```

### Via the API

The token is set directly as a header for the HTTP API. The name
of the header should be "X-Vault-Token" and the value should be the token.

## API

The Token authentication backend has a full HTTP API. Please see the
[Token auth backend API](/api/auth/token/index.html) for more
details.
