---
layout: "docs"
page_title: "Namespaces - Vault Enterprise"
sidebar_title: "Namespaces"
sidebar_current: "docs-vault-enterprise-namespaces"
description: |-
  Vault Enterprise has support for Namespaces, a feature to enable Secure Multi-tenancy (SMT) and self-management. 

---

# Vault Enterprise Namespaces

## Overview

Many organizations implement Vault as a "service", providing centralized 
management for teams within an organization while ensuring that those teams
operate within isolated environments known as *tenants*. 

There are two common challenges when implementing this architecture in Vault:

**Tenant Isolation**

Frequently teams within a VaaS environment require strong isolation from other
users in their policies, secrets, and identities. Tenant isolation is typically a 
result of compliance regulations such as [GDPR](https://www.eugdpr.org/), though it may 
be necessitated by corporate or organizational infosec requirements.

**Self-Management**

As new tenants are added, there is an additional human cost in the management 
overhead for teams. Given that tenants will likely have different policies and
request changes at a different rate, managing a multi-tenant environment can
become very difficult for a single team as the number of tenants within that
organization grow.

'Namespaces' is a set of features within Vault Enterprise that allows Vault
environments to support *Secure Multi-tenancy* (or *SMT*) within a single Vault
infrastructure. Through namespaces, Vault administrators can support tenant isolation
for teams and individuals as well as empower delegated administrators to manage their
own tenant environment.

## Usage

API operations performed under a namespace can be done by providing the relative
request path along with the namespace path using the `X-Vault-Namespace` header.
Similarly, the namespace header value can be provided in full or partially when
reaching into nested namespaces. When provided partially, the remaining
namespace path must be provided in the request path in order to reach into the
desired nested namespace.

Alternatively, the fully qualified path can be provided without using the
`X-Vault-Namespace` header. In either scenario, Vault will construct the fully
qualified path from these two sources to correctly route the request to the
appropriate namespace.

For example, these three requests are equivalent:
1. Path: `ns1/ns2/secret/foo`
2. Path: `secret/foo`, Header: `X-Vault-Namespace: ns1/ns2/`
3. Path: `ns2/secret/foo`, Header: `X-Vault-Namespace: ns1/`

## Architecture

Namespaces are isolated environments that functionally exist as "Vaults within a Vault."
They have separate login paths and support creating and managing data isolated to their
namespace. This data includes the following: 

- Secret Engines
- Auth Methods
- Policies
- Identities (Entities, Groups)
- Tokens

Rather than rely on Vault system admins, namespaces can be managed by delegated admins who
can be prescribed administration rights for their namespace. These delegated admins can also
create their own child namespaces, thereby prescribing admin rights on a subordinate group 
of delegate admins. 

Child namespaces can share policies from their parent namespaces. For example, a child namespace
may refer to parent identities (entities and groups) when writing policies that function only
within that child namespace. Similarly, a parent namespace can have policies asserted on child
identities. 

## Setup and Best Practices

A [deployment guide](/guides/operations/multi-tenant.html) is available to help guide you
through the deployment and administration of namespaces, and contains examples on architecture
for using namespaces to implement SMT across your organization. 

