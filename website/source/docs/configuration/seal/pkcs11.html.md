---
layout: "docs"
page_title: "PKCS11 - Seals - Configuration"
sidebar_current: "docs-configuration-seal-pkcs11"
description: |-
  The PKCS11 seal configures Vault to use an HSM with PKCS11 as the seal
  wrapping mechanism.
---

# `pkcs11` Seal

The PKCS11 seal configures Vault to use an HSM with PKCS11 as the seal wrapping
mechanism. Vault Enterprise's HSM PKCS11 support is activated by one of the
following:

* The presence of a `seal "pkcs11"` block in Vault's configuration file
* The presence of the environment variable `VAULT_HSM_LIB` set to the library's
  path as well as `VAULT_HSM_TYPE` set to `pkcs11`. If enabling via environment
  variable, all other required values (i.e. `VAULT_HSM_SLOT`) must be also
  supplied.

**IMPORTANT**: Having Vault generate its own key is the easiest way to get up
and running, but for security, Vault marks the key as non-exportable. If your
HSM key backup strategy requires the key to be exportable, you should generate
the key yourself. The list of creation attributes that Vault uses to generate
the key are listed at the end of this document.

## Requirements

The following software packages are required for Vault Enterprise HSM:

* PKCS#11 compatible HSM intgration library
* `libtldl` library

## `pkcs11` Parameters

These parameters apply to the `seal` stanza in the Vault configuration file:

* `lib` `(string: <required>)`: The path to the PKCS#11 library shared object
  file. May also be specified by the `VAULT_HSM_LIB` environment variable.
  **Note:** Depending on your HSM, this may be either a binary or a dynamic
  library, and its use may require other libraries depending on which system the
  Vault binary is currently running on (e.g.: a Linux system may require other
  libraries to interpret Windows .dll files). 
* `slot` `(string: <required>)`: The slot number to use, specified as a string
  (e.g. `"0"`). May also be specified by the `VAULT_HSM_SLOT` environment
  variable.
* `pin` `(string: <required>)`: The PIN for login. May also be specified by the
  `VAULT_HSM_PIN` environment variable. _If set via the environment variable,
  Vault will obfuscate the environment variable after reading it, and it will
  need to be re-set if Vault is restarted._
* `key_label` `(string: <required>)`: The label of the key to use. If the key
  does not exist and generation is enabled, this is the label that will be given
  to the generated key. May also be specified by the `VAULT_HSM_KEY_LABEL`
  environment variable.
* `hmac_key_label` `(string: <required>)`: The label of the key to use for
  HMACing. This needs to be a suitable type; a good choice is an AES key marked
  as valid for signing and verifying. If the key does not exist and generation
  is enabled, this is the label that will be given to the generated key. May
  also be specified by the `VAULT_HSM_HMAC_KEY_LABEL` environment variable.
* `mechanism` `(string: "0x1082")`: The encryption/decryption mechanism to use,
  specified as a decimal or hexadecimal (prefixed by `0x`) string. Currently
  only `0x1082` (corresponding to `CKM_AES_CBC` from the specification) is
  supported. May also be specified by the `VAULT_HSM_MECHANISM` environment
  variable.
* `hmac_mechanism` `(string: "0x0251")`: The encryption/decryption mechanism to
  use, specified as a decimal or hexadecimal (prefixed by `0x`) string.
  Currently only `0x0251` (corresponding to `CKM_SHA256_HMAC` from the
  specification) is supported. May also be specified by the
  `VAULT_HSM_HMAC_MECHANISM` environment variable.
* `generate_key` `(string: "false")`: If no existing key with the label
  specified by `key_label` can be found at Vault initialization time, instructs
  Vault to generate a key. This is a boolean expressed as a string (e.g.
  `"true"`). May also be specified by the `VAULT_HSM_GENERATE_KEY` environment
  variable.
* `regenerate_key` `(string: "false")`: At Vault initialization time, force
  generation of a new key even if one with the given `key_label` already exists.
  This is a boolean expressed as a string (e.g. `"true"`). May also be specified
  by the `VAULT_HSM_REGENERATE_KEY` environment variable.

~> **Note:** Although the configuration file allows you to pass in
`VAULT_HSM_PIN` as part of the seal's parameters, it is *strongly* reccommended
to set this value via environment variables.

## `pkcs11` Environment Variables

Alternatively, the HSM seal can be activated by providing the following
environment variables:

```text
* `VAULT_HSM_LIB`
* `VAULT_HSM_TYPE`
* `VAULT_HSM_SLOT`
* `VAULT_HSM_PIN`
* `VAULT_HSM_KEY_LABEL`
* `VAULT_HSM_HMAC_KEY_LABEL`
* `VAULT_HSM_HMAC_KEY_LABEL`
* `VAULT_HSM_MECHANISM`
* `VAULT_HSM_HMAC_MECHANISM`
* `VAULT_HSM_GENERATE_KEY`
* `VAULT_HSM_REGENERATE_KEY`
```

## `pkcs11` Example

This example shows configuring HSM PKCS11 seal through the Vault configuration
file by providing all the required values:

```hcl
seal "pkcs11" {
  lib            = "/usr/vault/lib/libCryptoki2_64.so"
  slot           = "0"
  pin            = "AAAA-BBBB-CCCC-DDDD" 
  key_label      = "vault-hsm-key" 
  hmac_key_label = "vault-hsm-hmac-key"
}
```

## Vault Key Generation Attributes

If Vault generates the HSM key for you, the following is the list of attributes
it uses. These identifiers correspond to official PKCS#11 identifiers.

* `CKA_CLASS`: `CKO_SECRET_KEY` (It's a secret key)
* `CKA_KEY_TYPE`: `CKK_AES` (Key type is AES)
* `CKA_VALUE_LEN`: `32` (Key size is 256 bits)
* `CKA_LABEL`: Set to the key label set in Vault's configuration
* `CKA_ID`: Set to a random 32-bit unsigned integer
* `CKA_PRIVATE`: `true` (Key is private to this slot/token)
* `CKA_TOKEN`: `true` (Key persists to the slot/token rather than being for one
  session only)
* `CKA_SENSITIVE`: `true` (Key is a sensitive value)
* `CKA_ENCRYPT`: `true` (Key can be used for encryption)
* `CKA_DECRYPT`: `true` (Key can be used for decryption)
* `CKA_WRAP`: `true` (Key can be used for wrapping)
* `CKA_UNWRAP`: `true` (Key can be used for unwrapping)
* `CKA_EXTRACTABLE`: `false` (Key cannot be exported)
