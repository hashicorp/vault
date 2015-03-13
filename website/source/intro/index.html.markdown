---
layout: "intro"
page_title: "Introduction"
sidebar_current: "what"
description: |-
  Welcome to the intro guide to Vault! This guide is the best place to start with Vault. We cover what Vault is, what problems it can solve, how it compares to existing software, and contains a quick start for using Vault.
---

# Introduction to Vault

Welcome to the intro guide to Vault! This guide is the best
place to start with Vault. We cover what Vault is, what
problems it can solve, how it compares to existing software,
and contains a quick start for using Vault.

If you are already familiar with the basics of Vault, the
[documentation](/docs/index.html) provides a better reference
guide for all available features as well as internals.

## What is Vault?

Vault is a tool for building, changing, and versioning infrastructure
safely and efficiently. Vault can manage existing and popular service
providers as well as custom in-house solutions.

Configuration files describe to Vault the components needed to
run a single application or your entire datacenter.
Vault generates an execution plan describing
what it will do to reach the desired state, and then executes it to build the
described infrastructure. As the configuration changes, Vault is able
to determine what changed and create incremental execution plans which
can be applied.

The infrastructure Vault can manage includes
low-level components such as
compute instances, storage, and networking, as well as high-level
components such as DNS entries, SaaS features, etc.

Examples work best to showcase Vault. Please see the
[use cases](/intro/use-cases.html).

The key features of Vault are:

* **Infrastructure as Code**: Infrastructure is described using a high-level
  configuration syntax. This allows a blueprint of your datacenter to be
  versioned and treated as you would any other code. Additionally,
  infrastructure can be shared and re-used.

* **Execution Plans**: Vault has a "planning" step where it generates
  an _execution plan_. The execution plan shows what Vault will do when
  you call apply. This lets you avoid any surprises when Vault
  manipulates infrastructure.

* **Resource Graph**: Vault builds a graph of all your resources,
  and parallelizes the creation and modification of any non-dependent
  resources. Because of this, Vault builds infrastructure as efficiently
  as possible, and operators get insight into dependencies in their
  infrastructure.

* **Change Automation**: Complex changesets can be applied to
  your infrastructure with minimal human interaction.
  With the previously mentioned execution
  plan and resource graph, you know exactly what Vault will change
  and in what order, avoiding many possible human errors.

## Next Steps

See the page on [Vault use cases](/intro/use-cases.html) to see the
multiple ways Vault can be used. Then see
[how Vault compares to other software](/intro/vs/index.html)
to see how it fits into your existing infrastructure. Finally, continue onwards with
the [getting started guide](/intro/getting-started/install.html) to use
Vault to manage real infrastructure and to see how it works.
