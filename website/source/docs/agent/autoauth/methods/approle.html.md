---
layout: "docs"
page_title: "Vault Agent Auto-Auth AppRole Method"
sidebar_title: "AppRole"
sidebar_current: "docs-agent-autoauth-methods-approle"
description: |-
  AppRole Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth AppRole Method

The `approle` method reads in a role-id/secret-id from a files and sends it to the [AppRole Auth
method](https://www.vaultproject.io/docs/auth/approle.html).

## Configuration

* `role_id_path` `(string: required)` - The path to the file with role-id

* `secret_id_path` `(string: required)` - The path to the file with secret-id
