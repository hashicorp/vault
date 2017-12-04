---
layout: "guides"
page_title: "Generate Root Tokens using Unseal Keys - Guides"
sidebar_current: "guides-generate-root"
description: |-
  Generate a new root token using a threshold of unseal keys.
---

# Generate Root Tokens Using Unseal Keys

It is generally considered a best practice to not persist
[root tokens][root-tokens]. Instead a root token should be generated using
Vault's `generate-root` command only when absolutely necessary. This guide
demonstrates regenerating a root token.

1. Unseal the vault using the existing quorum of unseal keys. You do not need to
  be authenticated.

    ```shell
    $ vault unseal
    # ...
    ```

2. Generate a one-time password:

    ```shell
    $ vault generate-root -genotp
    ```

3. Get the encoded root token:

    ```shell
    $ vault generate-root -otp="<otp>"
    ```

    This will require a quorum of unseal keys. This will then output an encoded
    root token.

4. Decode the encoded root token:

    ```shell
    $ vault generate-root -otp="<otp>" -decode="<encoded-token>"
    ```

Please see `vault generate-root -help` for information on the alternate
technique using a PGP key.

[root-tokens]: /docs/concepts/tokens.html#root-tokens
