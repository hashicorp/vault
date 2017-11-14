---
layout: "docs"
page_title: "Vault Enterprise Identity"
sidebar_current: "docs-vault-enterprise-identity"
description: |-
  Vault Enterprise has the foundations of the identity management system.
---

# Vault Enterprise Identity

In version 0.8, Vault introduced the foundations of identity management system.
The goal of identity in Vault is to associate a notion of caller identity to
the tokens used in Vault.

## Concepts

### Entities and Aliases

Each user will have multiple accounts with various identity providers. Users
can now be mapped as `Entities` and their corresponding accounts with
authentication providers can be mapped as `Aliases`. In essence, each entity
is made up of zero or more aliases.

### Entity Management

Entities in Vault **do not** automatically pull identity information from
anywhere. It needs to be explicitly managed by operators. This way, it is
flexible in terms of administratively controlling the number of entities to be
pulled in and pulled out of Vault, and in some sense Vault will serve as a
_cache_ of identities and not as the _source_ of identities.

### Entity Policies

Vault policies can be assigned to entities which will grant _additional_
permissions to the token on top of the existing policies on the token. If the
token presented on the API request contains an identifier for the entity and if
that entity has a set of policies on it, then the token will be capable of
performing actions allowed by the policies on the entity as well.

This is a paradigm shift in terms of _when_ the policies of the token get
evaluated. Before identity, the policy names on the token were immutable (not
the contents of those policies). But with entity policies, along with the
immutable set of policy names on the token, the evaluation of policies
applicable to the token through its identity will happen at request time. This
also adds enormous flexibility to control the behavior of already issued
tokens.

Its important to note that the policies on the entity are only a means to grant
_additional_ capabilities and not a replacement for the policies on the token,
and to know the full set of capabilities of the token with an associated entity
identifier, the policies on the token should be taken into account.

### Mount Bound Aliases

Vault supports multiple authentication backends and also allows enabling same
authentication backend on different mounts. The alias name of the user with
each identity provider will be unique within the provider. But Vault also needs
to uniquely distinguish between conflicting alias names across different
mounts of these identity providers. Hence the alias name, in combination with
the authentication backend mount's accessor serve as the unique identifier of a
alias.

### Implicit Entities

Operators can create entities for all the users of an auth mount
beforehand and assign policies to them, so that when users login, the desired
capabilities to the tokens via entities are already assigned. But if that's not
done, upon a successful user login from any of the authentication backends,
Vault will create a new entity and assign an alias against the login that was
successful.

Note that, tokens created using the token authentication backend will not have
an associated identity information. Logging in using the authentication
backends is the only way to create tokens that have a valid entity identifiers.

### Identity Auditing

If the token used to make API calls have an associated entity identifier, it will
be audit logged as well. This leaves a trail of actions performed by specific
users.

### API

Vault identity can be managed entirely over the HTTP API. Please see [Identity
API](/api/secret/identity/index.html) for more details.
