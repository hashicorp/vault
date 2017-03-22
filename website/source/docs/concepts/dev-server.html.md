---
layout: "docs"
page_title: "Dev Server Mode"
sidebar_current: "docs-concepts-devserver"
description: |-
  The dev server in Vault can be used for development or to experiment with Vault.
---

# "Dev" Server Mode

You can start Vault as a server in "dev" mode like so: `vault server -dev`.
This dev-mode server requires no further setup, and your local `vault` CLI will
be authenticated to talk to it. This makes it easy to experiment with Vault or
start a Vault instance for development. Every feature of Vault is available in
"dev" mode. The `-dev` flag just short-circuits a lot of setup to insecure
defaults.

~> **Warning:** Never, ever, ever run a "dev" mode server in production.
It is insecure and will lose data on every restart (since it stores data
in-memory). It is only made for development or experimentation.

## Properties

The properties of the dev server:

  * **Initialized and unsealed** - The server will be automatically initialized
    and unsealed. You don't need to use `vault unseal`. It is ready for use
    immediately.

  * **In-memory storage** - All data is stored (encrypted) in-memory. Vault
    server doesn't require any file permissions.

  * **Bound to local address without TLS** - The server is listening on
    `127.0.0.1:8200` (the default server address) _without_ TLS.

  * **Automatically Authenticated** - The server stores your root access
    token so `vault` CLI access is ready to go. If you are accessing Vault
    via the API, you'll need to authenticate using the token printed out.

  * **Single unseal key** - The server is initialized with a single unseal
    key. The Vault is already unsealed, but if you want to experiment with
    seal/unseal, then only the single outputted key is required.

## Use Case

The dev server should be used for experimentation with Vault features, such
as different authentication backends, secret backends, audit backends, etc.
If you're new to Vault, you may want to pick up with [Your First
Secret](/intro/getting-started/first-secret.html) in
our getting started guide.

In addition to experimentation, the dev server is very easy to automate
for development environments.
