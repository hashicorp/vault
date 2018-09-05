---
layout: "docs"
page_title: "Azure Key Vault - Seals - Configuration"
sidebar_current: "docs-configuration-seal-azurekeyvault"
description: |-
  The Azure Key Vault seal configures Vault to use Azure Key Vault as the seal wrapping
  mechanism.
---

# `azurekeyvault` Seal

The Azure Key Vault seal configures Vault to use Azure Key Vault as the seal
wrapping mechanism. Vault Enterprise's Azure Key Vault seal is activated by one of
the following:

* The presence of a `seal "azurekeyvault"` block in Vault's configuration file.
* The presence of the environment variable `VAULT_SEAL_TYPE` set to `azurekeyvault`.
  If enabling via environment variable, all other required values specific to
  Key Vault (i.e. `VAULT_AZUREKEYVAULT_VAULT_NAME`, etc.) must be also supplied, as
  well as all other Azure-related environment variables that lends to successful
  authentication (i.e. `AZURE_TENANT_ID`, etc.).

## `azurekeyvault` Example

This example shows configuring Azure Key Vault seal through the Vault
configuration file by providing all the required values:

```hcl
seal "azurekeyvault" {
  tenant_id      = "46646709-b63e-4747-be42-516edeaf1e14"
  client_id      = "03dc33fc-16d9-4b77-8152-3ec568f8af6e"
  client_secret  = "DUJDS3..."
  vault_name     = "hc-vault"
  key_name       = "vault_key"
}
```

## `azurekeyvault` Parameters

These parameters apply to the `seal` stanza in the Vault configuration file:

- `tenant_id` `(string: <required>)`: The tenant id for the Azure Active Directory organization. May 
  also be specified by the `AZURE_TENANT_ID` environment variable.

- `client_id` `(string: <required or MSI>)`: The client id for credentials to query the Azure APIs.
  May also be specified by the `AZURE_CLIENT_ID` environment variable.

- `client_secret` `(string: <required or MSI>)`: The client id for credentials to query the Azure APIs.
  May also be specified by the `AZURE_CLIENT_ID` environment variable.

- `environment` `(string: "AZUREPUBLICCLOUD")`: The Azure Cloud environment API endpoints to use.  May also 
  be specified by the `VAULT_AZUREKEYVAULT_VAULT_NAME` environment variable.

- `vault_name` `(string: <required>)`: The Key Vault vault to use the encryption keys for encryption and 
  decryption. May also be specified by the `VAULT_AZUREKEYVAULT_KEY_NAME` environment variable.

- `key_name` `(string: <required>)`: The Key Vault key to use for encryption and decryption. May also be specified by the
  `VAULT_AZUREKEYVAULT_KEY_NAME` environment variable.

## Authentication

Authentication-related values must be provided, either as environment
variables or as configuration parameters.

```text
Azure authentication values:

* `AZURE_TENANT_ID`
* `AZURE_CLIENT_ID`
* `AZURE_CLIENT_SECRET`
* `AZURE_ENVIRONMENT`
```

Note: If Vault is hosted on Azure, Vault can use Managed Service Identities (MSI) to access Azure instead of an environment and 
shared client id and secret.  MSI must be [enabled](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/qs-configure-portal-windows-vm) 
on the VMs hosting Vault. 


## `azurekeyvault` Environment Variables

Alternatively, the Azure Key Vault seal can be activated by providing the following
environment variables:

```text
* `VAULT_AZUREKEYVAULT_VAULT_NAME`
* `VAULT_AZUREKEYVAULT_KEY_NAME`
```

## Key Rotation

This seal supports rotating keys defined in Azure Key Vault.  Key metadata is stored with the 
encrypted data to ensure the correct key is used during decryption operations.