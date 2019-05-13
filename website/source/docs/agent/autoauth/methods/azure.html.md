---
layout: "docs"
page_title: "Vault Agent Auto-Auth Azure Method"
sidebar_title: "Azure"
sidebar_current: "docs-agent-autoauth-methods-azure"
description: |-
  Azure Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth Azure Method 

The `azure` method reads in Azure instance credentials and uses them to
authenticate with the [Azure Auth
method](https://www.vaultproject.io/docs/auth/azure.html). It reads most
parameters needed for authentication directly from instance information based
on the value of the `resource` parameter.

## Configuration

- `role` `(string: required)` - The role to authenticate against on Vault

- `resource` `(string: required)` - The resource name to use when getting instance information
