---
layout: "docs"
page_title: "Namespaces - Vault Enterprise"
sidebar_current: "docs-vault-enterprise-namespaces"
description: |-
  Vault Enterprise has support for Namespaces, a feature to enable Secure Multi-tenancy (SMT) and self-management. 

---

# Vault Enterprise Namespaces

## Overview

Many organizations implement *Vault as a Service* (or "VaaS"), providing centralized 
management to a security or ops team while ensuring that separate teams within that 
organization operate within self-contained environments known as "*tenants*." 

There are two common challenges when implementing this architecture in Vault:

**Tenant Isolation**

Frequently teams within a VaaS environment require strong isolation from other
users in their policies, secrets, and sometimes even their own identity entities 
and groups. Frequently tenant isolation is a result of regulations such as [GDPR](https://www.eugdpr.org/),
though it may be necessitated by corporate or organizational infosec requirements as 
well.

**Self-Management**

As new tenants are added, there is an additional human cost in the management 
overhead for teams. Given that tenants will likely have different policies and
request changes at a different rate, managing a multi-tenant environment can
become very difficult for a single team as the number of tenants within that
environment grow.

'Namespaces' is a set of features within Vault Enterprise that allows Vault
environments to support *Secure Multi-tenancy* (or *SMT*) within a single Vault Enterprise
infrastructure. Through namespaces, Vault administrators can support tenant isolation
for teams and individuals as well as empower those individuals to self-manage their
own tenant environment. 

## Architecture

Namespaces are isolated environments that functionally exist as "Vaults within a Vault."
They have separate login paths and support creating and managing data isolated to a namespace
including the following:

- Secret Engine Mounts
- Policies
- Identities (Entities, Groups)
- Tokens

Namespaces can also be configured to inherit all of this data from a higher *parent* namespace.
This simplifies the deployment of new namespaces, and can be combined with sentinel policies 
to prescribe organization-wide infosec policies on tenants. 

## Example Implementation



## Setup and Best Practices

A [deployment guide](/guides/operations/replication.html) is
available to help you get started, and contains examples on namespace architecture.

## API

Namespaces supports a full HTTP API. Please see the
[Vault Namespace API](/api/system/replication.html) for more
details.
