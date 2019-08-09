---
layout: "docs"
page_title: "OCI Object Storage - Storage Backends - Configuration"
sidebar_title: "OCI Object Storage"
sidebar_current: "docs-configuration-storage-ociobjectstorage"
description: |-
  The OCI Object Storage backend is used to persist Vault's data in OCI Object Storage service.
---

# OCI Object Storage Backend
The open source OCI Storage plugin for HashiCorp (HC) Vault provides a highly available storage backend for Vault, using OCI Object Storage. Read more about high availability in Vault at [High Availability Mode (HA)][ha-docs]. In a typical high available vault setup, Vault runs in multiple OCI compute instances with the same storage backend configuration. This plugin has been tested with version 0.11.4 of HC Vault and requires Go.

After setting up the OCI Storage plugin, you must periodically rotate unseal keys and vault encryption keys for security purposes.

### Configure Your OCI Tenancy to Run Vault with an Object Storage Backend
1. In your OCI tenancy, launch the [compute instances][compute-docs] that will run vault server.
    *  For regional high availability in vault, launch at least one compute instance per availability domain. 
    *  For high availability within an availability domain, launch at least one compute instance per [fault domain][fault-domain] in that availability domain.
1.  Create two Object Storage buckets in the same region where the compute instances are launched, one for vault data and one for leader lock. 
1. Write policies for bucket access.
    * Create a Dynamic Group in your OCI tenancy. For more information on dynamic groups, see https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/managingdynamicgroups.htm.
    * Create a policy that allows the Dynamic Group to read and write objects in these two buckets. There are multiple ways to write these policies. The [OCI Identity Policy][oci-identity-policy] can be used as a reference or starting point.
    * Create a [private load balancer][oci-lb] with backend sets pointing to the IP addresses of the launched compute instances.
        *  The backend set port of this load balancer is the port that the vault will listen to in the compute instances.
        *  In the Traffic Distribution Policy, do not enable USE SSL or USE SESSION PERSISTENCE. The goal is to create a TCP load balancer, so SSL is not applicable.
        *  For the health check, use protocol HTTP, port 443, a URL path of /v1/sys/health, and Status Code 200.
    * Edit the backends of the backend set to add the compute instance Oracle Cloud Identifier (OCID). The port for these backends is the port that vault will listen to in the compute instances. 
    * Add the appropriate security rules as required to allow connection from the load balancer to the compute instances.
    * For the listener, use protocol TCP and port 443. Do not enable USE SSL. For the backend set select the set previously created.
    * If you are using a DNS record, set the TTL of the record to a low value (such as 60 seconds).
1. Ensure your virtual cloud network (VCN) is configured correctly.
    * Allow a connection from the load balancer to the compute instances, based on the backend set ports.
    * Allow a connection between the compute instances in the port that the vault will listen to.
    * Allow a connection from the clients that will call the vault to the load balancer IP address in the listener port.
    

### Configure the OCI Vault Storage Plugin
1. Set the vault API address (api_addr) based on the host IP address and port, as seen by the other compute instance in the high availability setup. Typically, the IP address will be the private IP address of the compute instance, as shown in the console. For more information, see https://www.vaultproject.io/docs/configuration/.
    ```hcl
    export VAULT_API_ADDR=https://10.100.55.135:443
    ```
1. Create the configuration. Use the code below and add further configuration settings as detailed at https://www.vaultproject.io/docs/configuration/.
    ```hcl
    listener "tcp" {
     #narrow down the scope of the listener as per your security requirements
     address = "0.0.0.0:<port>"
     tls_disable = "false"
     tls_cert_file = "<cert_path or cert_chain_path>"
     tls_key_file = "<key_path>"
    }
      
    storage "oci_objectstorage" {
        namespaceName = "<object_storage_namespace_name>"
        bucketName = "<vault_data_bucket_name>"
        ha_enabled = "true"
        lockBucketName = "<leader_lock_bucket_name>"
    }
    ```
1. Start the vault server on each of the compute instances, passing their individual configuration file as the parameter. For more information, see https://learn.hashicorp.com/vault/getting-started/deploy#starting-the-server.
1. [Initialize vault][vault-init] on one of the compute instances.
    * The initialization process distributes an initial root token, which should be revoked after first use. Do not store the root token.
1. [Unseal the vault][vault-unseal] on all compute instances using the unseal key that was generated after starting vault.
    
    **Note:** The unseal keys should be encrypted and securely backed up outside of vault. If you lose the unseal keys and vault is sealed again (for example when the host is rebooted), you will be unable to unseal vault and the data stored in the vault will be inaccessible.
1. [Configure the vault audits and logs.][vault-config] Audits have sensitive information and must be secured.
1. Test that the high availability setup is working.
    * Stop vault on one of the compute instances.
    * Perform an API operation on vault, such as listing secrets, through the load balancer's listener IP address or the DNS record that has all the client IP addresses.
1. Follow the [production hardening recommendations][vault-hardening] for running vault securely.
1. Activate the OCI HashiCorp Vault Auth plugin.

[ha-docs]: https://www.vaultproject.io/docs/concepts/ha.html
[compute-docs]: https://docs.cloud.oracle.com/iaas/Content/Compute/Tasks/launchinginstance.htm
[fault-domain]: https://blogs.oracle.com/cloud-infrastructure/introducing-fault-domains-for-virtual-machine-and-bare-metal-instances
[oci-identity-policy]: https://docs.cloud.oracle.com/iaas/Content/Identity/Concepts/policies.htm
[oci-lb]: https://cloud.oracle.com/load-balancing/faq
[vault-init]: https://www.vaultproject.io/docs/commands/operator/init.html
[vault-unseal]: https://www.vaultproject.io/docs/concepts/seal.html
[vault-config]: https://www.vaultproject.io/docs/audit/syslog.html
[vault-hardening]: https://www.vaultproject.io/guides/operations/production

