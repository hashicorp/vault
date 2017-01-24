---
layout: "docs"
page_title: "Generate Root Tokens Using Unseal Keys"
sidebar_current: "docs-guides-generate-root"
description: |-
  Generate a new root token using a threshold of unseal keys.
---

# Generate Root Tokens Using Unseal Keys

It's considered [best practice](../concepts/tokens.html#root-tokens) not to
keep root tokens around, as they are all-powerful. Instead, if one is
absolutely needed, create it using Vault's `generate-root` command:

1. Unseal the vault. You do not need to be authenticated (you do not need an
 existing root token).
2. Generate a one-time password with `vault generate-root -genotp`.
3. Get the encoded root token with `vault generate-root -otp <generated_otp>`.
(Requires a quorum of unseal keys again, so needs to be done \<quorum\> times.)
4. Decode the encoded root token with
`vault generate-root -otp <generated_otp> -decode=<encoded_root_token>`.

See `vault generate-root -help` for information on the alternate technique
 using a PGP key.
