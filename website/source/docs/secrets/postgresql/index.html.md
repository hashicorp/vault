---
layout: "docs"
page_title: "Secret Backend: PostgreSQL"
sidebar_current: "docs-secrets-postgresql"
description: |-
  The PostgreSQL secret backend for Vault generates database credentials to access PostgreSQL.
---

# PostgreSQL Secret Backend

Name: `postgresql`

The PostgreSQL secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Vault makes use both of its own internal revocation system as well as the
`VALID UNTIL` setting when creating PostgreSQL users to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault help` after mounting the backend.

## Quick Start

TODO
