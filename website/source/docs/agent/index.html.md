---
layout: "docs"
page_title: "Vault Agent"
sidebar_current: "docs-agent"
description: |-
  Vault Agent is a client-side daemon that can be used to perform some Vault
  functionality automatically.
---

# Vault Agent

Vault Agent is a client daemon that can perform useful tasks.

To get help, run:

```text
$ vault agent -h
```
## Auto-Auth

Vault Agent allows for easy authentication to Vault in a wide variety of
environments. Please see the [Auto-Auth docs](/docs/agent/autoauth/index.html)
for information.

Auto-Auth functionality takes place within an `auto_auth` configuration stanza.

## Configuration

There is one currently-available general configuration option:

- `pid_file` `(string: "")` - Path to the file in which the agent's Process ID
  (PID) should be stored.
