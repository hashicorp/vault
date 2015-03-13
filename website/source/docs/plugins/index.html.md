---
layout: "docs"
page_title: "Plugins"
sidebar_current: "docs-plugins"
description: |-
  Vault is built on a plugin-based architecture. All providers and provisioners that are used in Vault configurations are plugins, even the core types such as AWS and Heroku. Users of Vault are able to write new plugins in order to support new functionality in Vault.
---

# Plugins

Vault is built on a plugin-based architecture. All providers and
provisioners that are used in Vault configurations are plugins, even
the core types such as AWS and Heroku. Users of Vault are able to
write new plugins in order to support new functionality in Vault.

This section of the documentation gives a high-level overview of how
to write plugins for Vault. It does not hold your hand through the
process, however, and expects a relatively high level of understanding
of Go, provider semantics, Unix, etc.

~> **Advanced topic!** Plugin development is a highly advanced
topic in Vault, and is not required knowledge for day-to-day usage.
If you don't plan on writing any plugins, we recommend not reading
this section of the documentation.
