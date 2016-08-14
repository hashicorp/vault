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

## Standby Nodes Must Be 0.6.1 As Well

Once an active node is running 0.6.1, only standby nodes running 0.6.1+ will be
able to form an HA cluster. If following our [general upgrade
instructions](https://www.vaultproject.io/docs/install/upgrade.html) this will
not be an issue.

## Root Token Creation Restrictions

Root tokens (tokens with the `root` policy) can no longer be created except by
another root token or the `generate-root` endpoint.

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

## LDAP Auth Backend Does Not Search `memberOf`

The LDAP backend went from a model where all permutations of storing and
filtering groups were tried in all cases to one where specific filters are
defined by the administrator. This vastly increases overall directory
compatibility, especially with Active Directory when using nested groups, but
unfortunately has the side effect that `memberOf` is no longer searched for by
default, which is a breaking change for many existing setups. 

`Scenario 2` in the [updated
documentation](https://github.com/hashicorp/vault/blob/master/website/source/docs/auth/ldap.html.md)
shows an example of configuring the backend to query `memberOf`. It is
recommended that a test Vault server be set up and that successful
authentication can be performed using the new configuration before upgrading a
primary or production Vault instance.

In addition, if LDAP is relied upon for authentication, operators should ensure
that they have valid tokens with policies allowing modification of LDAP
parameters before upgrading, so that once an upgrade is performed, the new
configuration can be specified successfully.

## App-ID is Deprecated

With the addition of of the new [AppRole
backend](https://www.vaultproject.io/docs/auth/approle.html), App-ID is
deprecated. There are no plans to remove it, but we encourage using AppRole
whenever possible, as it offers enhanced functionality and can accommodate many
more types of authentication paradigms.
