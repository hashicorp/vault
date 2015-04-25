---
layout: "docs"
page_title: "Secret Backend: MySQL"
sidebar_current: "docs-secrets-mysql"
description: |-
  The MySQL secret backend for Vault generates database credentials to access MySQL.
---

# MySQL Secret Backend

Name: `mysql`

The MySQL secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Vault makes use of its own internal revocation system to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault help` after mounting the backend.

## Quick Start

TODO
