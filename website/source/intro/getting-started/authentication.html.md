---
layout: "intro"
page_title: "Authentication"
sidebar_current: "gettingstarted-auth"
description: |-
  Authentication to Vault gives a user access to use Vault. Vault can authenticate using multiple methods.
---

# Authentication

Now that we know how to use the basics of Vault, it is important to understand
how to authenticate to Vault itself. Up to this point, we haven't had to
authenticate because starting the Vault sever in dev mode automatically logs
us in as root. In practice, you'll almost always have to manually authenticate.

On this page, we'll talk specifically about _authentication_. On the next
page, we talk about _authorization_.
Authentication is the mechanism of assigning an identity to a Vault user.
The access control and permissions associated with an identity are
authorization, and will not covered on this page.

Vault has pluggable authentication backends, making it easy to authenticate
with Vault using whatever form works best for your organization. On this page
we'll use the token backend as well as the GitHub backend.

## Tokens

We'll first explain token authentication before going over any other
authentication backends. Token authentication is enabled by default in
Vault and cannot be disabled. It is also what we've been using up to this
point.

TODO
