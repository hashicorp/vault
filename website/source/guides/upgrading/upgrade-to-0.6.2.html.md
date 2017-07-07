---
layout: "guides"
page_title: "Upgrading to Vault 0.6.2 - Guides"
sidebar_current: "guides-upgrading-to-0.6.2"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.6.2. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.6.2. Please read it carefully.

## Request Forwarding On By Default

In 0.6.1 this feature was in beta and required opting-in, but is now enabled by
default. This can be disabled via the `"disable_clustering"` parameter in
Vault's [config](/docs/configuration/index.html), or
per-request with the `X-Vault-No-Request-Forwarding` header.

## AppRole Role Constraints

Creating or updating a role now requires at least one constraint to be enabled,
whereas previously it was sufficient to require only the role ID by itself.
Currently there are two constraints: `bind_secret_id` and `bound_cidr_list`.
`bind_secret_id` is enabled by default. Roles which were previously using only
the role ID for authentication will continue to work but will require a
constraint to be specified if updated.

## Convergent Encryption v2

New keys in `transit` using convergent mode will use a new nonce derivation
mechanism rather than require the user to supply a nonce. While not explicitly
increasing security, it minimizes the likelihood that a user will use the mode
improperly and impact the security of their keys. Keys in convergent mode that
were created in 0.6.1 will continue to work with the same mechanism
(user-supplied nonce).

## `etcd` HA Off By Default

Following in the footsteps of `dynamodb`, the `etcd` storage backend now
requires that `ha_enabled` be explicitly specified in the configuration file.
The backend currently has known broken HA behavior, so this flag discourages
use by default without explicitly enabling it. If you are using this
functionality, when upgrading, you should set `ha_enabled` to `"true"` *before*
starting the new versions of Vault.

## Reading Wrapped Responses From `cubbyhole/response` Is Deprecated

The `sys/wrapping/unwrap` endpoint should be used instead as it provides
additional security, auditing, and other benefits. The ability to read directly
will be removed in a future release.

## Default/Max Lease/Token TTLs Now 32 Days

In previous versions of Vault the default was 30 days, but changing it to 32
days allows some operations (e.g. reauthenticating, renewing, etc.) to be
performed via a monthly cron job.

## AppRole Secret ID Endpoints Changed

Secret ID and Secret ID accessors are no longer part of request URLs. The `GET`
and `DELETE` operations are now moved to new endpoints (`/lookup` and
`/destroy`) which consumes the input from the body via `POST` (or `PUT`) and
not the URL.

## Behavior Change for `bound_iam_role_arn` in AWS-EC2 Backend

In prior versions a bug caused the `bound_iam_role_arn` value in the `aws-ec2`
authentication backend to actually use the instance profile ARN.  This has been
corrected, but as a result there is a behavior change. To match using the
instance profile ARN, a new parameter `bound_iam_instance_profile_arn` has been
added. Existing roles will automatically transfer the value over to the correct
parameter, but the next time the role is updated, the new meanings will take
effect.
