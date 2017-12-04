---
layout: "docs"
page_title: "Architecture"
sidebar_current: "docs-internals-architecture"
description: |-
  Learn about the internal architecture of Vault.
---

# Architecture

Vault is a complex system that has many different pieces. To help both users and developers of Vault
build a mental model of how it works, this page documents the system architecture.

~> **Advanced Topic!** This page covers technical details
of Vault. You don't need to understand these details to
effectively use Vault. The details are documented here for
those who wish to learn about them without having to go
spelunking through the source code. However, if you're an
operator of Vault, we recommend learning about the architecture
due to the importance of Vault in an environment.

# Glossary

Before describing the architecture, we provide a glossary of terms to help
clarify what is being discussed:

* **Storage Backend** - A storage backend is responsible for durable storage of _encrypted_ data.
  Backends are not trusted by Vault and are only expected to provide durability. The storage
  backend is configured when starting the Vault server.

* **Barrier** - The barrier is cryptographic steel and concrete around the Vault. All data that
  flows between Vault and the Storage Backend passes through the barrier. The barrier ensures
  that only encrypted data is written out, and that data is verified and decrypted on the way
  in. Much like a bank vault, the barrier must be "unsealed" before anything inside can be accessed.

* **Secret Backend** - A secret backend is responsible for managing secrets. Simple secret backends
  like the "kv" backend simply return the same secret when queried. Some backends support
  using policies to dynamically generate a secret each time they are queried. This allows for
  unique secrets to be used which allows Vault to do fine-grained revocation and policy updates.
  As an example, a MySQL backend could be configured with a "web" policy. When the "web" secret
  is read, a new MySQL user/password pair will be generated with a limited set of privileges
  for the web server.

* **Audit Backend** - An audit backend is responsible for managing audit logs. Every request to Vault
  and response from Vault goes through the configured audit backends. This provides a simple
  way to integrate Vault with multiple audit logging destinations of different types.

* **Auth Backend** - An auth backend is used to authenticate users or applications which
  are connecting to Vault. Once authenticated, the backend returns the list of applicable policies
  which should be applied. Vault takes an authenticated user and returns a client token that can
  be used for future requests. As an example, the `userpass` backend uses a username and password
  to authenticate the user. Alternatively, the `github` backend allows users to authenticate
  via GitHub.

* **Client Token** - A client token is a conceptually similar to a session cookie on a web site.
  Once a user authenticates, Vault returns a client token which is used for future requests.
  The token is used by Vault to verify the identity of the client and to enforce the applicable
  ACL policies. This token is passed via HTTP headers.

* **Secret** - A secret is the term for anything returned by Vault which contains confidential
  or cryptographic material. Not everything returned by Vault is a secret, for example
  system configuration, status information, or backend policies are not considered Secrets.
  Secrets always have an associated lease. This means clients cannot assume that the secret
  contents can be used indefinitely. Vault will revoke a secret at the end of the lease, and
  an operator may intervene to revoke the secret before the lease is over. This contract
  between Vault and its clients is critical, as it allows for changes in keys and policies
  without manual intervention.

* **Server** - Vault depends on a long-running instance which operates as a server.
  The Vault server provides an API which clients interact with and manages the
  interaction between all the backends, ACL enforcement, and secret lease revocation.
  Having a server based architecture decouples clients from the security keys and policies,
  enables centralized audit logging and simplifies administration for operators.

# High-Level Overview

A very high level overview of Vault looks like this:

[![Architecture Overview](/assets/images/layers.png)](/assets/images/layers.png)

Let's begin to break down this picture. There is a clear separation of components
that are inside or outside of the security barrier. Only the storage backend and
the HTTP API are outside, all other components are inside the barrier.

The storage backend is untrusted and is used to durably store encrypted data. When
the Vault server is started, it must be provided with a storage backend so that data
is available across restarts. The HTTP API similarly must be started by the Vault server
on start so that clients can interact with it.

Once started, the Vault is in a _sealed_ state. Before any operation can be performed
on the Vault it must be unsealed. This is done by providing the unseal keys. When
the Vault is initialized it generates an encryption key which is used to protect all the
data. That key is protected by a master key. By default, Vault uses a technique known
as [Shamir's secret sharing algorithm](https://en.wikipedia.org/wiki/Shamir's_Secret_Sharing)
to split the master key into 5 shares, any 3 of which are required to reconstruct the master
key.

[![Vault Shamir Secret Sharing Algorithm](/assets/images/vault-shamir-secret-sharing.svg)](/assets/images/vault-shamir-secret-sharing.svg)

The number of shares and the minimum threshold required can both be specified. Shamir's
technique can be disabled, and the master key used directly for unsealing. Once Vault
retrieves the encryption key, it is able to decrypt the data in the storage backend,
and enters the _unsealed_ state. Once unsealed, Vault loads all of the configured
audit, credential and secret backends.

The configuration of those backends must be stored in Vault since they are security
sensitive. Only users with the correct permissions should be able to modify them,
meaning they cannot be specified outside of the barrier. By storing them in Vault,
any changes to them are protected by the ACL system and tracked by audit logs.

After the Vault is unsealed, requests can be processed from the HTTP API to the Core.
The core is used to manage the flow of requests through the system, enforce ACLs,
and ensure audit logging is done.

When a client first connects to Vault, it needs to authenticate. Vault provides
configurable credential backends providing flexibility in the authentication mechanism
used. Human friendly mechanisms such as username/password or GitHub might be
used for operators, while applications may use public/private keys or tokens to authenticate.
An authentication request flows through core and into a credential backend, which determines
if the request is valid and returns a list of associated policies.

Policies are just a named ACL rule. For example, the "root" policy is built-in and
permits access to all resources. You can create any number of named policies with
fine-grained control over paths. Vault operates exclusively in a whitelist mode, meaning
that unless access is explicitly granted via a policy, the action is not allowed.
Since a user may have multiple policies associated, an action is allowed if any policy
permits it. Policies are stored and managed by an internal policy store. This internal store
is manipulated through the system backend, which is always mounted at `sys/`.

Once authentication takes place and a credential backend provides a set of applicable
policies, a new client token is generated and managed by the token store. This client token
is sent back to the client, and is used to make future requests. This is similar to
a cookie sent by a website after a user logs in. The client token may have a lease associated
with it depending on the credential backend configuration. This means the client token
may need to be periodically renewed to avoid invalidation.

Once authenticated, requests are made providing the client token. The token is used
to verify the client is authorized and to load the relevant policies. The policies
are used to authorize the client request. The request is then routed to the secret backend,
which is processed depending on the type of backend. If the backend returns a secret,
the core registers it with the expiration manager and attaches a lease ID.
The lease ID is used by clients to renew or revoke their secret. If a client allows the
lease to expire, the expiration manager automatically revokes the secret.

The core handles logging of requests and responses to the audit broker, which fans the
request out to all the configured audit backends. Outside of the request flow, the core
performs certain background activity. Lease management is critical, as it allows
expired client tokens or secrets to be revoked automatically. Additionally, Vault handles
certain partial failure cases by using write ahead logging with a rollback manager.
This is managed transparently within the core and is not user visible.

# Getting in Depth

This has been a brief high-level overview of the architecture of Vault. There
are more details available for each of the sub-systems.

For other details, either consult the code, ask in IRC or reach out to the mailing list.
