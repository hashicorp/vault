---
layout: "docs"
page_title: "Vault Agent Auto-Auth PCF Method"
sidebar_title: "PCF"
sidebar_current: "docs-agent-autoauth-methods-pcf"
description: |-
  PCF Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth PCF Method

The `pcf` method performs authentication against the [PCF Auth 
method] (https://www.vaultproject.io/docs/auth/pcf.html).

## Credentials

The Vault agent will use the `CF_INSTANCE_CERT` and `CF_INSTANCE_KEY` env variables to
construct a valid login call for PCF.

## Configuration

- `role` `(string: required)` - The role to authenticate against on Vault.
