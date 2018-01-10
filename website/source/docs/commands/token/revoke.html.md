---
layout: "docs"
page_title: "token revoke - Command"
sidebar_current: "docs-commands-token-revoke"
description: |-
  The "token revoke" revokes authentication tokens and their children. If a
  TOKEN is not provided, the locally authenticated token is used.
---

# token revoke

The `token revoke` revokes authentication tokens and their children. If a TOKEN
is not provided, the locally authenticated token is used. The `-mode` flag can
be used to control the behavior of the revocation.

## Examples

Revoke a token and all the token's children:

```text
$ vault token revoke 96ddf4bc-d217-f3ba-f9bd-017055595017
Success! Revoked token (if it existed)
```

Revoke a token leaving the token's children:

```text
$ vault token revoke -mode=orphan 96ddf4bc-d217-f3ba-f9bd-017055595017
Success! Revoked token (if it existed)
```

Revoke a token by accessor:

```text
$ vault token revoke -accessor 9793c9b3-e04a-46f3-e7b8-748d7da248da
Success! Revoked token (if it existed)
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-accessor` `(bool: false)` - Treat the argument as an accessor instead of a
  token.

- `-mode` `(string: "")` - Type of revocation to perform. If unspecified, Vault
  will revoke the token and all of the token's children. If "orphan", Vault will
  revoke only the token, leaving the children as orphans. If "path", tokens
  created from the given authentication path prefix are deleted along with their
  children.

- `-self` -  Perform the revocation on the currently authenticated token.
