---
layout: "docs"
page_title: "Identity - Secrets Engines"
sidebar_current: "docs-secrets-identity"
description: |-
  The Identity secrets engine for Vault manages client identities.
---

# Identity Secrets Engine

Name: `identity`

The Identity secrets engine is the identity management solution for Vault. It
internally maintains the clients who are recognized by Vault. Each client is
internally termed as an `Entity`. An entity can have multiple `Aliases`. For
example, a single user who has accounts in both Github and LDAP, can be mapped
to a single entity in Vault that has 2 aliases, one of type Github and one of
type LDAP. When a client authenticates via any of the credential backend
(except the Token backend), Vault creates a new entity and attaches a new
alias to it, if a corresponding entity doesn't already exist. The entity identifier will
be tied to the authenticated token. When such tokens are put to use, their
entity identifiers are audit logged, marking a trail of actions performed by
specific users.

Identity store allows operators to **manage** the entities in Vault. Entities
can be created and aliases can be tied to entities, via the ACL'd API. There
can be policies set on the entities which adds capabilities to the tokens that
are tied to entity identifiers. The capabilities granted to tokens via the
entities are **an addition** to the existing capabilities of the token and
**not** a replacement. The capabilities of the token that get inherited from
entities are computed dynamically at request time. This provides flexibility in
controlling the access of tokens that are already issued.

This secrets engine will be mounted by default. This secrets engine cannot be
disabled or moved.

## API

The Identity secrets engine has a full HTTP API. Please see the
[Identity secrets engine API](/api/secret/identity/index.html) for more
details.
