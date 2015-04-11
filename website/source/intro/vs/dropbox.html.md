---
layout: "intro"
page_title: "Vault vs. Dropbox"
sidebar_current: "vs-other-dropbox"
description: |-
  Comparison between Vault and attempting to store secrets with Dropbox.
---

# Vault vs. Dropbox

It is an unfortunate truth that many organizations, big and small,
often use Dropbox as a mechanism for storing secrets. It is so common
that we've decided to make a special section for it instead of throwing
it under the "custom solutions" header.

Dropbox is not made for storing secrets. Even if you're using something
such as an encrypted disk image within Dropbox, it is subpar versus a
real secret storage server.

A real secret management tool such as Vault has a stronger security
model, integrates with many different authentication services, stores
audit logs, can generate dynamic secrets, and more.

And, due to `vault` CLI, using `vault` on a developer machine is
simple!
