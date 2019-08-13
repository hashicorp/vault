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
  KMS (i.e. `VAULT_OCIKMS_SEAL_KEY_ID`, `VAULT_OCIKMS_CRYPTO_ENDPOINT`) must be also supplied, as well as all
  other OCI-related [environment variables][oci-sdk] that lends to successful
  authentication. 
## `ocikms` Configure Your OCI Tenancy to Run Vault with KMS Auto Unseal Plugin
1. In your OCI tenancy, [launch the compute instances][oci-compute] that will run vault server. 
    * For regional high availability in vault, launch at least one compute instance per availability domain. 
    * For high availability within an availability domain, launch at least one compute instance per [fault domain][oci-fd] in that availability domain.
2.  Write policies for KMS access.
    * Create a [Dynamic Group][oci-dg] in your OCI tenancy.
    * Create a policy that allows the Dynamic Group to use keys from KMS. There are multiple ways to write these policies. The [OCI Identity Policy][oci-id] can be used as a reference or starting point.

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
   
## `ocikms` Configure the OCI Vault KMS Auto Unseal Plugin
1. Set the vault API address (api_addr) based on the host IP address and port, as seen by the other compute instance in the high availability setup. Typically, the IP address will be the private IP address of the compute instance, as shown in the console.
    ```hcl
    export VAULT_API_ADDR=https://10.100.55.135:443
    ```
1. If you want to use Instance Principal, add section configuration below and add further configuration settings as detailed at https://www.vaultproject.io/docs/configuration/.
    ```hcl
    seal "ocikms" {
        cryptoEndpoint = "<kms-crypto-endpoint>"
        keyID = "<kms-key-id>"
    }
    # Notes:
    # cryptoEndpoint can be replaced by VAULT_OCIKMS_CRYPTO_ENDPOINT environment var
    # keyID can be replaced by VAULT_OCIKMS_SEAL_KEY_ID environment var
    ```
1. If you want to use User Principal, the plugin will take API key you defined for OCI SDK, often under `~/.oci/config`.
    ```
    seal "ocikms" {
        authTypeAPIKey = true
        cryptoEndpoint = "<kms-crypto-endpoint>"
        keyID = "<kms-key-id>"
    }
    ```
1. Start the vault server on each of the compute instances, passing their individual configuration file as the parameter.
1. Initialize vault on one of the compute instances. This will also auto unseal your instance.
    * The initialization process distributes a set of Recovery Keys and an initial root token, which should be revoked after first use. Do not store the root token.
    * When using auto unseal there are certain operations in the vault that still require a quorum of users to perform an operation, such as generating a root token. During the initialization process, a set of Shamir keys are generated that are called Recovery Keys and are used for these operations.
    * Recovery Keys also become Unseal Key when the vault is migrated back to Shamir from auto unseal.

## `ocikms` Rotate OCI KMS Master Key
For the [KMS key rotation feature][oci-kms-rotation], KMS will create a new version of key. This process is independent from the vault; vault still uses the same keyID without any break.

If you want to change the keyID, migrate to Shamir, change keyID, and then migrate to OCI KMS with the new keyID.


[oci-sdk]: https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm
[oci-compute]: https://docs.cloud.oracle.com/iaas/Content/Compute/Tasks/launchinginstance.htm
[oci-fd]: https://blogs.oracle.com/cloud-infrastructure/introducing-fault-domains-for-virtual-machine-and-bare-metal-instances
[oci-dg]:  https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/managingdynamicgroups.htm
[oci-id]: https://docs.cloud.oracle.com/iaas/Content/Identity/Concepts/policies.htm
[oci-kms-rotation]: https://docs.cloud.oracle.com/iaas/Content/KeyManagement/Tasks/managingkeys.htm


