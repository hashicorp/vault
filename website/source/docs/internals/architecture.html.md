---
layout: "docs"
page_title: "Architecture"
sidebar_current: "docs-internals-architecture"
description: |-
  Vault builds a dependency graph from the Vault configurations, and walks this graph to generate plans, refresh state, and more. This page documents the details of what are contained in this graph, what types of nodes there are, and how the edges of the graph are determined.
---

# Architecture

Vault is a complex system that has many different pieces. To help both users and developers of Consul
build a mental model of how it works, this page documents the system architecture.

~> **Advanced Topic!** This page covers technical details
of Vault. You don't need to understand these details to
effectively use Vault. The details are documented here for
those who wish to learn about them without having to go
spelunking through the source code.

# Glossary

Before describing the architecture, we provide a glossary of terms to help clarify what is being discussed:

* Storage Backend - A storage backend is responsible for durable storage of _encrypted_ data.
  backends are not trusted by Vault and are only expected to provide durability. The storage
  backend is configured when starting the Vault server.

* Barrier - The barrier is cryptographic steel and concrete around the Vault. All data that
  flows between Vault and the Storage Backend passes through the barrier. The barrier ensures
  that only encrypted data is written out, and that data is verified and decrypted on the way
  in. Much like a bank vault, the barrier must be "unsealed" before anything inside can be accessed.

* Secret Backend - A secret backend is responsible for managing secrets. Simple secret backends
  like the "generic" backend simply return the same secret when queried. Some backends support
  using policies to dynamically generate a secret each time they are queried. This allows for
  unique secrets to be used which allows Vault to do fine-grained revocation and policy updates.
  As an example, a MySQL backend could be configured with a "web" policy. When the "web" secret
  is read, a new MySQL user/password pair will be generated with a limited set of privileges
  for the web server.

* Audit Backend - An audit backend is responsible for managing audit logs. Every request to Vault
  and response from Vault goes through the configured audit backends. This provides a simple
  way to integrate Vault with multiple audit logging destinations of different types.

* Credential Backend - A credential backend is used to authenticate users or applications which
  are connecting to Vault. Once authenticated, the backend returns the list of applicable policies
  which shoud be applied. Vault takes an authenticated user and returns a client token that can
  be used for future requests. As an example, the `user-password` backend uses a username and password
  to authenticate the user. Alternatively, the `github` backend allows users to authenticate
  via GitHub.

* Client Token - A client token is a conceptually similar to a session cookie on a web site.
  Once a user authenticates, Vault returns a client token which is used for future requests.
  The token is used by Vault to verify the identity of the client and to enforce the applicable
  ACL policies.

* Secret - A secret is the term for anything returned by Vault which contains confidential
  or cryptographic material. Not all everything returned by Vault is a secret, for example
  system configuration, status information, or backend policies are not considered Secrets.
  Secrets always have an associated lease. This means clients cannot assume that the secret
  contents can be used indefinitely. Vault will revoke a secret at the end of the lease, and
  an operator may intervene to revoke the secret before the lease is over. This contract
  between Vault and it's clients is critical, as it allows for changes in keys and policies
  without manual intervention.

* Server - Vault depends on a long-running instance which operates as a server.
  The Vault server provides an API which clients interact with and manages the
  interaction between all the backends, ACL enforcement, and secret lease revocation.
  Having a server based architecture decouples clients from the security keys and policies,
  enables centralized audit logging and simplifies administration for operators.

