---
layout: "docs"
page_title: "OCI KMS - Seals - Configuration"
sidebar_title: "OCI KMS"
sidebar_current: "docs-configuration-seal-ocikms"
description: |-
  The OCI KMS seal configures Vault to use OCI KMS as the seal wrapping
  mechanism.
---

# `ocikms` Seal

The OCI KMS seal configures Vault to use OCI KMS as the seal wrapping mechanism.
The OCI KMS seal is activated by one of the following:

* The presence of a `seal "ocikms"` block in Vault's configuration file
* The presence of the environment variable `VAULT_SEAL_TYPE` set to `ocikms`. If
  enabling via environment variable, all other required values specific to OCI
  KMS (i.e. `VAULT_OCIKMS_SEAL_KEY_ID`, `VAULT_OCIKMS_CRYPTO_ENDPOINT` `VAULT_OCIKMS_MANAGEMENT_ENDPOINT`) must be also supplied, as well as all
  other OCI-related [environment variables][oci-sdk] that lends to successful
  authentication. 
  
## `ocikms` Example

This example shows configuring OCI KMS seal through the Vault configuration file
by providing all the required values:

```hcl
seal "ocikms" {
    keyID               = "ocid1.key.oc1.iad.afnxza26aag4s.abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx"
    cryptoEndpoint      = "https://afnxza26aag4s-crypto.kms.us-ashburn-1.oraclecloud.com"
    managementEndpoint  = "https://afnxza26aag4s-management.kms.us-ashburn-1.oraclecloud.com"
    authTypeAPIKey      = "true"
}
```

## `ocikms` Parameters

These parameters apply to the `seal` stanza in the Vault configuration file:

- `keyID` `(string: <required>)`: The OCI KMS key ID to use. May also be
  specified by the `VAULT_OCIKMS_SEAL_KEY_ID` environment variable.
- `cryptoEndpoint` `(string: <required>)`: The OCI KMS cryptographic endpoint (or data plane endpoint) 
  to be used to make OCI KMS encryption/decryption requests. May also be specified by the `VAULT_OCIKMS_CRYPTO_ENDPOINT` environment
  variable.
- `managementEndpoint` `(string: <required>)`: The OCI KMS management endpoint (or control plane endpoint) 
  to be used to make OCI KMS key management requests. May also be specified by the `VAULT_OCIKMS_MANAGEMENT_ENDPOINT` environment
  variable.
- `authTypeAPIKey` `(boolean: false)`: Specifies if using API key to authenticate to OCI KMS service.
  If it is `false`, Vault authenticates using instance principal from compute instance. See Authentication section for details. Default is `false`. 

## Authentication

Authentication-related values must be provided, either as environment
variables or as configuration parameters.

1. If you want to use Instance Principal, add section configuration below and add further configuration settings as detailed at https://www.vaultproject.io/docs/configuration/.
    ```hcl
    seal "ocikms" {
        cryptoEndpoint = "<kms-crypto-endpoint>"
        managementEndpoint = "<kms-management-endpoint>"
        keyID = "<kms-key-id>"
    }
    # Notes:
    # cryptoEndpoint can be replaced by VAULT_OCIKMS_CRYPTO_ENDPOINT environment var
    # managementEndpoint can be replaced by VAULT_OCIKMS_MANAGEMENT_ENDPOINT environment var
    # keyID can be replaced by VAULT_OCIKMS_SEAL_KEY_ID environment var
    ```
1. If you want to use User Principal, the plugin will take API key you defined for OCI SDK, often under `~/.oci/config`.
    ```
    seal "ocikms" {
        authTypeAPIKey = true
        cryptoEndpoint = "<kms-crypto-endpoint>"
        managementEndpoint = "<kms-management-endpoint>"
        keyID = "<kms-key-id>"
    }
    ```

To grant permission for a compute instance to use OCI KMS service, write policies for KMS access.

- Create a [Dynamic Group][oci-dg] in your OCI tenancy.
- Create a policy that allows the Dynamic Group to use or manage keys from OCI KMS. There are multiple ways to write these policies. The [OCI Identity Policy][oci-id] can be used as a reference or starting point.

The most common policy allows a dynamic group of tenant A to use KMS's keys in tenant B:
1. Policy for tenant A
    ```text
    define tenancy tenantB as <tenantB-ocid>
     
    endorse dynamic-group <dynamic-group-name> to use keys in tenancy tenantB
 
    ```
1. Policy for tenant B
   ```text
   define tenancy tenantA as <tenantA-ocid>
    
   define dynamic-group <dynamic-group-name> as <dynamic-group-ocid>

   admit dynamic-group <dynamic-group-name> of tenancy tenantA to use keys in compartment <key-compartment>


   ```
   
## `ocikms` Rotate OCI KMS Master Key

For the [OCI KMS key rotation feature][oci-kms-rotation], OCI KMS will create a new version of key internally. This process is independent from the vault, vault still uses the same `keyID` without any interruption.

If you want to change the `keyID`, migrate to Shamir, change `keyID`, and then migrate to OCI KMS with the new `keyID`.

[oci-sdk]: https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm
[oci-dg]:  https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/managingdynamicgroups.htm
[oci-id]: https://docs.cloud.oracle.com/iaas/Content/Identity/Concepts/policies.htm
[oci-kms-rotation]: https://docs.cloud.oracle.com/iaas/Content/KeyManagement/Tasks/managingkeys.htm


