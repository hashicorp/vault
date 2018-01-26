---
layout: "guides"
page_title: "Generate Root Tokens using Unseal Keys - Guides"
sidebar_current: "guides-configuration-generate-root"
description: |-
  Generate a new root token using a threshold of unseal keys.
---

# Generate Root Tokens Using Unseal Keys

In a production Vault installation, the initial [root token][root-tokens] should only be used
for initial configuration.

The following command creates a token for an admin:

```shell
vault token-create -metadata "name=ADMIN_NAME" -display-name="ADMIN_USER_NAME" \
-orphan -no-default-policy
```

After a subset of administrators have sudo access,
almost all operations can be performed. However, for some system critical
operations, a root token may still be required.

It is generally considered a best practice to not persist [root
tokens][root-tokens]. Instead a root token should be generated using Vault's
`generate-root` command only when absolutely necessary. A quorum of unseal key
holders can generate a new root token. This enforces that there
is no single person has complete access to the system.

This guide demonstrates regenerating a root token using a one-time-password (OTP).

## Steps to Regenerate Root Tokens

1. Make sure that the Vault server is unsealed
2. Generate a one-time-password (OTP) to share
3. Each unseal key holder runs `generate-root` with the OTP
4. Decode the generated root token

### Step 1: Make sure that the Vault server is unsealed

First, verify the status:

```shell
$ vault status
```
The output should indicate that the Vault is unsealed (`Sealed: false`).

If the status indicates that the Vault server is sealed, unseal the vault using
the existing quorum of unseal keys. You do not need to be authenticated.

```shell
$ vault unseal
# ...
```

### Step 2: Generate a one-time-password (OTP)

Generate a one-time password:

```shell
$ vault generate-root -genotp
```

This generates the OTP to generate a new root token. The output would look like:

```shell
$ vault generate-root -genotp
OTP: +G07n16yukWxyn7nQbG0aw==
```

### Step 3: Each unseal key holder runs generate-root

Each unseal key holder runs the `generate-root` command with generated OTP:

```shell
$ vault generate-root -otp="<otp>"
```

Example:

```shell
$ vault generate-root -otp="+G07n16yukWxyn7nQbG0aw=="
Root generation operation nonce: abe86476-c6c5-9ca9-426e-bb6eba7fc987
Key (will be hidden):
Nonce: abe86476-c6c5-9ca9-426e-bb6eba7fc987
Started: true
Generate Root Progress: 1
Required Keys: 3
Complete: false
```

When the root key generation completes, an encoded new root token will be
provided.

The output would look like:

```shell
$ vault generate-root -otp="+G07n16yukWxyn7nQbG0aw=="
Root generation operation nonce: abe86476-c6c5-9ca9-426e-bb6eba7fc987
Key (will be hidden):
Nonce: abe86476-c6c5-9ca9-426e-bb6eba7fc987
Started: true
Generate Root Progress: 3
Required Keys: 3
Complete: true

Encoded root token: O7gIhugL3oHKeVmxpKGcYA==
```

### Step 4: Decode the generated root tokens

Run the `generate-root` command as follow:

```shell
$ vault generate-root -otp="<otp>" -decode="<encoded-token>"
```

Example:

```shell
$ vault generate-root -otp="+G07n16yukWxyn7nQbG0aw==" -decode="O7gIhugL3oHKeVmxpKGcYA=="
Root token: c3d53319-b6b9-64c4-7bb3-2756e510280b
```

## Additional References

Instead of using a shared OTP, you can pass a file on a disk containing a public
PGP key.

Example:

```shell
$ vault generate-root -pgp-key="keyname.asc"
```

Please see `vault generate-root -help` for more information about using PGP.

[root-tokens]: /docs/concepts/tokens.html#root-tokens
