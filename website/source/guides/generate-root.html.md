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
  be authenticated to generate a new root token, but the Vault must be unsealed
  and a quorum of unseal keys must be available.

    ```shell
    $ vault operator unseal
    # ...
    ```

### Using OTP

In this method, an OTP is XORed with the generated token on final output.

1. Generate a one-time password (OTP) to use for XORing the resulting token:

    ```text
    $ vault operator generate-root -generate-otp
    mOXx7iVimjE6LXQ2Zna6NA==
    ```

    Save this OTP because you will need it to get the decoded final root token.

1. Initialize a root token generation, providing the OTP code from the step
   above:

    ```text
    $ vault operator generate-root -init -otp=mOXx7iVimjE6LXQ2Zna6NA==
    Nonce              f67f4da3-4ae4-68fb-4716-91da6b609c3e
    Started            true
    Progress           0/5
    Complete           false
    ```

    The nonce value should be distributed to all unseal key holders.

1. Each unseal key holder provides their unseal key:

    ```text
    $ vault operator generate-root
    Root generation operation nonce: f67f4da3-4ae4-68fb-4716-91da6b609c3e
    Unseal Key (will be hidden): ...
    ```

    If there is a tty, Vault will prompt for the key and automatically
    complete the nonce value. If there is no tty, or if the value is piped
    from stdin, the user must specify the nonce value from the `-init`
    operation.

    ```text
    $ echo $UNSEAL_KEY | vault operator generate-root -nonce=f67f4da3... -
    ```

1. When the quorum of unseal keys are supplied, the final user will also get
   the encoded root token.

    ```text
    $ vault operator generate-root
    Root generation operation nonce: f67f4da3-4ae4-68fb-4716-91da6b609c3e
    Unseal Key (will be hidden):

    Nonce         f67f4da3-4ae4-68fb-4716-91da6b609c3e
    Started       true
    Progress      5/5
    Complete      true
    Root Token    IxJpyqxn3YafOGhqhvP6cQ==
    ```

1. Decode the encoded token using the OTP:

    ```text
    $ vault operator generate-root \
        -decode=IxJpyqxn3YafOGhqhvP6cQ== \
        -otp=mOXx7iVimjE6LXQ2Zna6NA==

    24bde68f-3df3-e137-cf4d-014fe9ebc43f
    ```

### Using PGP

1. Initialize a root token generation, providing the path to a GPG public key
   or keybase username of a user to encrypted the resulting token.

    ```text
    $ vault operator generate-root -init -pgp-key=keybase:sethvargo
    Nonce              e24dec5e-f1ea-2dfe-ecce-604022006976
    Started            true
    Progress           0/5
    Complete           false
    PGP Fingerprint    e2f8e2974623ba2a0e933a59c921994f9c27e0ff
    ```

    The nonce value should be distributed to all unseal key holders.

1. Each unseal key holder providers their unseal key:

    ```text
    $ vault operator generate-root
    Root generation operation nonce: e24dec5e-f1ea-2dfe-ecce-604022006976
    Unseal Key (will be hidden): ...
    ```

    If there is a tty, Vault will prompt for the key and automatically
    complete the nonce value. If there is no tty, or if the value is piped
    from stdin, the user must specify the nonce value from the `-init`
    operation.

    ```text
    $ echo $UNSEAL_KEY | vault generate-root -nonce=f67f4da3... -
    ```

1. When the quorum of unseal keys are supplied, the final user will also get
   the encoded root token.

    ```text
    $ vault operator generate-root
    Root generation operation nonce: e24dec5e-f1ea-2dfe-ecce-604022006976
    Unseal Key (will be hidden):

    Nonce              e24dec5e-f1ea-2dfe-ecce-604022006976
    Started            true
    Progress           1/1
    Complete           true
    PGP Fingerprint    e2f8e2974623ba2a0e933a59c921994f9c27e0ff
    Root Token         wcFMA0RVkFtoqzRlARAAI3Ux8kdSpfgXdF9mg...
    ```

1. Decrypt the encrypted token using associated private key:

    ```text
    $ echo "wcFMA0RVkFtoqzRlARAAI3Ux8kdSpfgXdF9mg..." | base64 --decode | gpg --decrypt

    d0f71e9b-ebff-6d8a-50ae-b8859f2e5671
    ```

    or via keybase:

    ```text
    $ echo "wcFMA0RVkFtoqzRlARAAI3Ux8kdSpfgXdF9mg..." | base64 --decode | keybase pgp decrypt

    d0f71e9b-ebff-6d8a-50ae-b8859f2e5671
    ```

[root-tokens]: /docs/concepts/tokens.html#root-tokens
