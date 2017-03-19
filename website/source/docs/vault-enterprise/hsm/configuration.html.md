---
layout: "docs"
page_title: "Vault Enterprise HSM Configuration"
sidebar_current: "docs-vault-enterprise-hsm-configuration"
description: |-
  Vault Enterprise HSM configuration details.

---

# Vault Enterprise HSM Configuration

Vault Enterprise's HSM support is activated by one of the following:

* The presence of an `hsm` block in Vault's configuration file
* Values set in both the `VAULT_HSM_LIB` and `VAULT_HSM_TYPE` environment
  variables

**IMPORTANT**: Having Vault generate its own key is the easiest way to get up
and running, but for security, Vault marks the key as non-exportable. If your
HSM key backup strategy requires the key to be exportable, you should generate
the key yourself. The list of creation attributes that Vault uses to generate
the key are listed at the end of this document.

## HSM Block Directives

Like the rest of Vault's configuration files, the `hsm` block is in
[HCL](https://github.com/hashicorp/hcl) format.

The key of the `hsm` block is the type of HSM:

```hcl
hsm "pkcs11" {
  ...
}
```

The type can also be set by the `VAULT_HSM_TYPE` environment variable.
Currently, only `pkcs11` is supported.

The following are the block directives and their effects. All parameters are
strings.

### Required Directives

 * `lib`: The path to the PKCS#11 library shared object file. May also be
   specified by the `VAULT_HSM_LIB` environment variable.
 * `slot`: The slot number to use, specified as a string (e.g. `"0"`). May also
   be specified by the `VAULT_HSM_SLOT` environment variable.
 * `pin`: The PIN for login. May also be specified by the `VAULT_HSM_PIN`
   environment variable. _If set via the environment variable, Vault will
   obfuscate the environment variable after reading it, and it will need to be
   re-set if Vault is restarted._
 * `key_label`: The label of the key to use. If the key does not exist and
   generation is enabled, this is the label that will be given to the generated
   key. May also be specified by the `VAULT_HSM_KEY_LABEL` environment
   variable.

### Optional Directives

 * `mechanism`: The encryption/decryption mechanism to use, specified as a
   decimal or hexadecimal (prefixed by `0x`) string. Currently only `0x1082`
   (corresponding to `CKM_AES_CBC` from the specification) is supported. May
   also be specified by the `VAULT_HSM_MECHANISM` environment variable.
 * `generate_key`: If no existing key with the label specified by `key_label`
   can be found at Vault initialization time, instructs Vault to generate a
   key. This is a boolean expressed as a string (e.g. `"true"`). May also be
   specified by the `VAULT_HSM_GENERATE_KEY` environment variable.
 * `regenerate_key`: At Vault initialization time, force generation of a new
   key even if one with the given `key_label` already exists. This is a boolean
   expressed as a string (e.g. `"true"`). May also be specified by the
   `VAULT_HSM_REGENERATE_KEY` environment variable.

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
