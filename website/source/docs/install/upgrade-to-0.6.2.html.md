---
layout: "install"
page_title: "Upgrading to Vault 0.6.2"
sidebar_current: "docs-install-upgrade-to-0.6.2"
description: |-
  Learn how to upgrade to Vault 0.6.2
---

# Overview

This page contains the list of breaking changes for Vault 0.6.2. Please read it
carefully.

## AppRole Role Constraints

Creating or updating a role now requires at least one constraint to be enabled.
Currently there are only 2 constraints: `bind_secret_id` and `bound_cidr_list`.
`bind_secret_id` is enabled by default. Roles which had `bind_secret_id`
disabled and `bound_cidr_list` not set, will require a constraint to be
speficied during further updates.
