---
layout: "docs"
page_title: "Vault Agent Auto-Auth JWT Method"
sidebar_title: "JWT"
sidebar_current: "docs-agent-autoauth-methods-jwt"
description: |-
  JWT Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth JWT Method

The `jwt` method reads in a JWT from a file and sends it to the [JWT Auth
method](https://www.vaultproject.io/docs/auth/jwt.html). Since JWTs often have
limited lifetime, it constantly watches for a new JWT to be written, and when
found it will immediately ingress this value, delete the file, and use the new
JWT to perform a reauthentication.

## Configuration

* `path` `(string: required)` - The path to the JWT file

* `role` `(string: required)` - The role to authenticate against on Vault
