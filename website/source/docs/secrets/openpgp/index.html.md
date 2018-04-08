---
layout: "docs"
page_title: "OpenPGP - Secrets Engines"
sidebar_current: "docs-secrets-openpgp"
description: |-
  The OpenPGP secrets engine for Vault does PGP operations on data in-transit.
---

# OpenPGP Secrets Engine

Name: `openpgp`

The transit secrets engine handles PGP operations on data in-transit.
Vault doesn't store the data sent to the secrets engine. If you do not
have a strong requirement on PGP, the
[transit secrets engine](/docs/secrets/transit/index.html) should be
preferred for handling cryptographic functions on data in-transit.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the OpenPGP secrets engine:

    ```text
    $ vault secrets enable openpgp
    Success! Enabled the openpgp secrets engine at: openpgp/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.


1. Create a new PGP key:

    ```text
    $ vault write -f openpgp/keys/my-key
    Success! Data written to: openpgp/keys/my-key
    ```

1. Retrieving the public part of a PGP key:

    ```text
    $ vault read openpgp/keys/my-key
      Key            Value
      ---            -----
      exportable     false
      fingerprint    677138a1cbfe1e0e49f4c58784e9759b896915dd
      public_key     -----BEGIN PGP PUBLIC KEY BLOCK-----

      xsBNBFrKEcUBCADUdV3HxWN7jA+QSoB7A6l1pkmTvYnDxUO8cbSXebKeoSKawqVN
      vN8LiG+ttjtQ0DcNVNSJw6jJMImizJ6fh1Jr9X2UXzVXjPO7e1Na+VJsQst5L4SP
      ozm0MW+AqqbFo9fDSJb/2Ic7piozFW5h99v5lYk5qINP4LbdHvgSCpuCkgajd9RD
      5NW1s5DOPM3XTcEXeHD7augrNNBGed5cMVPB/WQeK7BnatH54qbB8QcFgKKme3Ug
      N4m3IDOsS4SB5MxhHB1PVy/MsFwYoWSpBDUc1Axfaq85XNTVHjF1cPmA2F/5N7GU
      BfF3m/bdd6Px8iPQUCDP/Km4k/SE5WO2J3UZABEBAAHNAMLAYgQTAQgAFgUCWsoR
      xQkQhOl1m4lpFd0CGwMCGQEAAAKeCACd1lxWuC1EwCkZ1KnxCYpi9ahOqoyf3IyU
      EIulqop06xRDRaDtycEq7dEGd303b2DFtLDAVE5ow2hhUgcDXZ7V48hPcDV4ebk7
      e6thtpMDHce5SPSWsVHwqGgIABkInPEilLdyIKFAMhHPm7ItAhCuAPs7A9d++O7p
      VrOlqmxW6Qakn3ABwrT6sy0Nm9X0cdSD7bDUi8bpxu6m86HqvX6DuF7ctmvQCs4M
      oLHWfPEjtymVmJo60PAIokUp0pMW9b0fdYcFFTuVFRV/b+JtisT6yie8bKhsKAfb
      So09bzB9/swvFehGVSmqOnGpMzknBgzL2TxPah0Fk2gflTAX/n9GzsBNBFrKEcUB
      CADMaWz3GEySa8paJI+T0a97djWWW0Pnpdjy1+izw1JVYpef3YZZpnUfHXLTaquc
      +aXM4dMijZ2QLu2vlmCrKO6Q/Uh5E0HB7kv7Sb2hOm3Q7hSVArSgH9M7jA0h9Gff
      JPmDnQnb/K2dMnLO3oINRb8SS5tZdGaPlgLx6Gvb+TklqkV1QW7Zr0DOQx66XgGI
      TwLIW3w76AskNbGwcVL9ZGZg4q6DnEuPufCx00LfdJLULT/Y2TaSCvqNW/3JUeje
      y5dWBQp9eTGGSiO5Fs46bu0xVTkS9VdNu4ctLCliga0c/RBLM5nlL6nqB/ZOWV7Y
      iUfjxvPgQxB377/Z3m7Obtk/ABEBAAHCwF8EGAEIABMFAlrKEcUJEITpdZuJaRXd
      AhsMAAANUQgAouKbFWv0dWFOtfP/J4/g/pAxuCOInVn1YwNp6sW5fRU2SdSRr7rR
      0ZLCjNG0FyOtY4DnmlIkVA5ENJG8TVQyJDgSTT4HPObEEUf5Y/XdMkHXgUdkypS0
      6padU/pb5lXodcwwqwQjVFufNyO3moMvl8Mx5JeBPyoevZ6A7l25IW/low8nQA9Q
      j76QkSjZjCiNZs2MSRGsc3JolQUhxxO+x9aWAOMbsCqNApkH/4vW9AiFAQ7qtwRN
      ncfMkc4Oy6XEUAP8RqWQJDkbGzXpMZC7ScuIW0jTW/CpUyQYueOjsW1lNsUG9rXc
      ArmjhwBbOdJBTbQLqqK7Tle9YFqX0nz1zw==
      =mn53
      -----END PGP PUBLIC KEY BLOCK-----
    ```

1. Signing some plaintext data using the `/sign` endpoint with a named key:

   ```text
   $ vault write openpgp/sign/my-key input=$(base64 <<< "my data")
     Key          Value
     ---          -----
     signature    wsBcBAABCAAQBQJayhPHCRCE6XWbiWkV3QAAR8AIAEmpvwtvPv8akWsVCcMhRtOgxKJqg9kekWx8s9Ki36FM2ozth2FQA40hAID5JbX+1ju19BtoXAeVkPxepRpwG6D9bTFj3TPHbW8i81oXXTqr5CWcLI74KwrWbzPdEnPn61tFWt0czWXDzfeZ+7EmuzPBkfDhCIz0G0y/Pw+/N5EjZiiGwoQix/rgC8HkVxcmb3hOSdxx71LueCwNqzCEHD2hT78vedwVacGS1h7HfTDg7PYp1CNEcaFQu6UHQMqp3rza9pUDe0xc15vHnViIvUICXXZ7tNSl9ENdvcbkk7iSOtlb0aheDnnKf44KWfJMkiKRarxgc2ZgJnv6Uc5iLC8=
   ```

1. Verify if the signature is valid for the given data using the `/verify`
endpoint with a named key:

    ```text
    $ vault write openpgp/verify/my-key input=$(base64 <<< "my data") signature='wsBcBAABCAAQBQJayhXbCRCE6XWbiWkV3QAA9y8IAIcLLHdXqFBhWmY8lRLwSYiiZtWJ0cmbVJl0JWhV0nptRwiuXzNGWYkF8aB9NaM6Q85yZzKPH1soe7/anE8nBJSR1TSawHD8Ph03jcJpdbrTG0KUS5GgJ4fu3d+L+AroNWE8Nu4ZRAV948reEOXuBQKl4/JU3c+j/Dki6RVwP6o2shAdLtrkfptk2r7jz+VihRe6jPuJc85Tnji3k+8I6MAvRztUQrODgbCZ5/58gu6HpmJGcbTHvWYfS05D7nqiZPsVjZiZHBnSyNNs+nDZ+X4ITIQiztQMip6u+IvLaRh26q5s3eoQ+eFG/kXrTNVZIn8w2306dSqRaTGo6hWY5vU='
      Key      Value
      ---      -----
      valid    true
    ```

1. Decrypt a piece of data using the `/decrypt` endpoint with a named key:

    ```text
    $ vault write openpgp/decrypt/my-key ciphertext=$(base64 --wrap=0 encrypted-data.pgp)
      Key          Value
      ---          -----
      plaintext    bXkgZGF0YQo=

    ```

## API

The Nomad secret backend has a full HTTP API. Please see the
[OpenPGP secret backend API](/api/secret/openpgp/index.html) for more
details.
