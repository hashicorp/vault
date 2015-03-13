---
layout: "docs"
page_title: "Configuration"
sidebar_current: "docs-config"
description: |-
  Vault uses text files to describe infrastructure and to set variables. These text files are called Vault _configurations_ and end in `.tf`. This section talks about the format of these files as well as how they're loaded.
---

# Configuration

Vault uses text files to describe infrastructure and to set variables.
These text files are called Vault _configurations_ and end in
`.tf`. This section talks about the format of these files as well as
how they're loaded.

The format of the configuration files are able to be in two formats:
Vault format and JSON. The Vault format is more human-readable,
supports comments, and is the generally recommended format for most
Vault files. The JSON format is meant for machines to create,
modify, and update, but can also be done by Vault operators if
you prefer. Vault format ends in `.tf` and JSON format ends in
`.tf.json`.

Click a sub-section in the navigation to the left to learn more about
Vault configuration.
