---
layout: "docs"
page_title: "Vault Cookbook"
sidebar_current: "docs-cookbook"
description: |-
  Vault server how-to cookbook.
---

# Day-to-day tasks with Vault

## Generate a root token (when none exists)

It's considered [best practice](../concepts/tokens.html#root-tokens) not to keep root tokens around, as they are all-powerful. Instead, if one is absolutely needed, create it using vault's generate-root command:

1. Unseal the vault. You do not need to be authenticated (you do not need an existing root token).
2. Generate a one-time password with `vault generate-root -genotp`
3. Get the encoded root token: `vault generate-root -otp <generated_otp>` (Requires a quorum of unseal keys again, so needs to be done \<quorum\> times.)
4. Decode the encoded root token with `vault generate-root -otp <generated_otp> -decode=<encoded_root_token> `

(See `vault generate-root -h` for information on the alternate technique using a PGP key.)
