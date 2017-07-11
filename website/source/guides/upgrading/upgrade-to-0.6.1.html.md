---
layout: "guides"
page_title: "Upgrading to Vault 0.6.1 - Guides"
sidebar_current: "guides-upgrading-to-0.6.1"
description: |-
  This page contains the list of breaking changes for Vault 0.6.1. Please read
  it carefully.
---

# Overview

This page contains the list of breaking changes for Vault 0.6.1. Please read it
carefully.

## Standby Nodes Must Be 0.6.1 As Well

Once an active node is running 0.6.1, only standby nodes running 0.6.1+ will be
able to form an HA cluster. If following our [general upgrade
instructions](/guides/upgrading/index.html) this will
not be an issue.

## Health Endpoint Status Code Changes

Prior to 0.6.1, the health endpoint would return a `500` (Internal Server
Error) for both a sealed and uninitialized state. In both states this was
confusing, since it was hard to tell, based on the status code, an actual
internal error from Vault from a Vault that was simply uninitialized or sealed,
not to mention differentiating between those two states.

In 0.6.1, a sealed Vault will return a `503` (Service Unavailable) status code.
As before, this can be adjusted with the `sealedcode` query parameter. An
uninitialized Vault will return a `501` (Not Implemented) status code. This can
be adjusted with the `uninitcode` query parameter.

This removes ambiguity/confusion and falls more in line with the intention of
each status code (including `500`).

## Root Token Creation Restrictions

Root tokens (tokens with the `root` policy) can no longer be created except by
another root token or the
[`generate-root`](/api/system/generate-root.html)
endpoint or CLI command.

## PKI Backend Certificates Will Contain Default Key Usages

Issued certificates from the `pki` backend against roles created or modified
after upgrading will contain a set of default key usages. This increases
compatibility with some software that requires strict adherence to RFCs, such
as OpenVPN.

This behavior is fully adjustable; see the [PKI backend
documentation](/docs/secrets/pki/index.html) for
details.

## DynamoDB Does Not Support HA By Default

If using DynamoDB and want to use HA support, you will need to explicitly
enable it in Vault's configuration; see the
[documentation](/docs/configuration/index.html#ha_enabled)
for details.

If you are already using DynamoDB in an HA fashion and wish to keep doing so,
it is *very important* that you set this option **before** upgrading your Vault
instances. Without doing so, each Vault instance will believe that it is
standalone and there could be consistency issues.

## LDAP Auth Backend Forgets Bind Password and Insecure TLS Settings

Due to a bug, these two settings are forgotten if they have been configured in
the LDAP backend prior to 0.6.1. If you are using these settings with LDAP,
please be sure to re-submit your LDAP configuration to Vault after the upgrade,
so ensure that you have a valid token to do so before upgrading if you are
relying on LDAP authentication for permissions to modify the backend itself.

## LDAP Auth Backend Does Not Search `memberOf`

The LDAP backend went from a model where all permutations of storing and
filtering groups were tried in all cases to one where specific filters are
defined by the administrator. This vastly increases overall directory
compatibility, especially with Active Directory when using nested groups, but
unfortunately has the side effect that `memberOf` is no longer searched for by
default, which is a breaking change for many existing setups.

`Scenario 2` in the [updated
documentation](/docs/auth/ldap.html) shows an
example of configuring the backend to query `memberOf`. It is recommended that
a test Vault server be set up and that successful authentication can be
performed using the new configuration before upgrading a primary or production
Vault instance.

In addition, if LDAP is relied upon for authentication, operators should ensure
that they have valid tokens with policies allowing modification of LDAP
parameters before upgrading, so that once an upgrade is performed, the new
configuration can be specified successfully.

## App-ID is Deprecated

With the addition of of the new [AppRole
backend](/docs/auth/approle.html), App-ID is
deprecated. There are no current plans to remove it, but we encourage using
AppRole whenever possible, as it offers enhanced functionality and can
accommodate many more types of authentication paradigms. App-ID will receive
security-related fixes only.
