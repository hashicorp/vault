---
layout: "docs"
page_title: "Vault Agent Auto-Auth"
sidebar_current: "docs-agent-autoauth"
description: |-
  Vault Agent's Auto-Auth functionality allows easy and automatic
  authentication to Vault in a variety of environments.
---

# Vault Agent Auto-Auth

The Auto-Auth functionality of Vault Agent allows for easy authentication in a
wide variety of environments.

## Functionality

Auto-Auth consists of two parts: a Method, which is the authentication method
that should be used in the current environment; and one or more Sinks, which
are locations where the agent should write a token any time the current token
value has changed.

When the agent is started with Auto-Auth enabled, it will attempt to acquire a
Vault token using the configured Method. On failure, it will back off for a
short while (including some randomness to help prevent thundering herd
scenarios) and retry. On success, it will keep the resulting token renewed
until renewal is no longer allowed or fails, at which point it will attempt to
reauthenticate.

Every time an authentication is successful, the token is written to the
configured Sinks, subject to their configuration.

## Configuration

The top level `auto_auth` block has two configuration entries:

- `method` `(object)` - Configuration for the method

- `sinks` `(array of objects)` - Configuration for the sinks

### Configuration (Method)
