---
layout: "install"
page_title: "Upgrading to Vault 0.6.1"
sidebar_current: "docs-install-upgrade-to-0.6.1"
description: |-
  Learn how to upgrade to Vault 0.6.1
---

# Overview

This page contains the list of breaking changes for Vault 0.6.1. Please read it
carefully.

## PKI Backend Certificates Will Contain Default Key Usages

Issued certificates from the `pki` backend against roles created or modified
after upgrading will contain a set of default key usages. This increases
compatibility with some software that requires strict adherence to RFCs, such
as OpenVPN.

This behavior is fully adjustable; see the [PKI backend
documentation](https://www.vaultproject.io/docs/secrets/pki/index.html) for
details.

## DynamoDB Does Not Support HA By Default

If using DynamoDB and want to use HA support, you will need to explicitly
enable it in Vault's configuration; see the
[documentation](https://www.vaultproject.io/docs/config/index.html#ha_enabled)
for details.

If you are already using DynamoDB in an HA fashion and wish to keep doing so,
it is *very important* that you set this option before upgrading your Vault
instances. Without doing so, each Vault instance will believe that it is
standalone and there will be consistency issues.
