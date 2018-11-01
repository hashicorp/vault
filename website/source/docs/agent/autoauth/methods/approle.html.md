---
layout: "docs"
page_title: "Vault Agent Auto-Auth AppRole Method"
sidebar_title: "AppRole"
sidebar_current: "docs-agent-autoauth-methods-approle"
description: |-
  AppRole Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth AppRole Method

The `approle` method reads in a role ID and a secret ID from files and sends
the values to the [AppRole Auth
method](https://www.vaultproject.io/docs/auth/approle.html).

The method caches values and it is safe to delete the role ID/secret ID files
after they have been read. In fact, by default, after reading the secret ID,
the agent will delete the file. New files or values written at the expected
locations will be used on next authentication and the new values will be
cached.

## Configuration

* `role_id_file_path` `(string: required)` - The path to the file with role ID

* `secret_id_file_path` `(string: required)` - The path to the file with secret
  ID

* `remove_secret_id_file_after_reading` `(bool: optional, defaults to true)` -
  This can be set to `false` to disable the default behavior of removing the
  secret ID file after it's been read.

* `secret_id_response_wrapping_path` `(string: optional)` - If set, the value
  at `secret_id_file_path` will be expected to be a [Response-Wrapping
  Token](https://www.vaultproject.io/docs/concepts/response-wrapping.html)
  containing the output of the secret ID retrieval endpoint for the role (e.g.
  `auth/approle/role/webservers/secret-id`) and the creation path for the
  response-wrapping token must match the value set here.
