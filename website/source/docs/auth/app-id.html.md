---
layout: "docs"
page_title: "AppID - Auth Methods"
sidebar_current: "docs-auth-appid"
description: |-
  The AppID auth method is a mechanism for machines to authenticate with Vault.
---

# AppID Auth Method

~> **DEPRECATED!** As of Vault 0.6.1, AppID is deprecated in favor of
[AppRole](/docs/auth/approle.html). AppRole can accommodate the same workflow as
AppID while enabling much more secure and flexible management and other types
of authentication workflows. No new features or enhancements are planned for App
ID, and new users should use AppRole instead of AppID.

The AppID auth method is a mechanism for machines to authenticate with Vault.
It works by requiring two hard-to-guess unique pieces of information: a unique
app ID, and a unique user ID.

The goal of this auth method is to allow elastic users (dynamic
machines, containers, etc.) to authenticate with Vault without having to store
passwords outside of Vault. It is a single method of solving the
chicken-and-egg problem of setting up Vault access on a machine.  With this
provider, nobody except the machine itself has access to both pieces of
information necessary to authenticate. For example: configuration management
will have the app IDs, but the machine itself will detect its user ID based on
some unique machine property such as a MAC address (or a hash of it with some
salt).

An example, real world process for using this provider:

1. Create unique app IDs (UUIDs work well) and map them to policies.  (Path:
   map/app-id/<app-id>)

2. Store the app IDs within configuration management systems.

3. An out-of-band process run by security operators map unique user IDs to
   these app IDs. Example: when an instance is launched, a cloud-init system
   tells security operators a unique ID for this machine. This process can be
   scripted, but the key is that it is out-of-band and out of reach of
   configuration management.  (Path: map/user-id/<user-id>)

4. A new server is provisioned. Configuration management configures the app
   ID, the server itself detects its user ID. With both of these pieces of
   information, Vault can be accessed according to the policy set by the app
   ID.

More details on this process follow:

- The app ID is a unique ID that maps to a set of policies. This ID is generated
  by an operator and configured into the method. The ID itself is usually a
  UUID, but any hard-to-guess unique value can be used.

- After creating app IDs, an operator authorizes a fixed set of user IDs with
  each app ID. When a valid {app ID, user ID} tuple is given to the "login"
  path, then the user is authenticated with the configured app ID policies.

- The user ID can be any value (just like the app ID), however it is generally a
  value unique to a machine, such as a MAC address or instance ID, or a value
  hashed from these unique values.


## Authentication

Via the CLI:

```text
$ vault write auth/app-id/login/:app-id user_id=:user_id
```

Via the API:

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --method POST \
    --data '{"user_id": ":user_id"}' \
    https://vault.rocks/v1/auth/app-id/login/:app_id
```

## Configuration

Auth methods must be configured in advance before users or machines can
authenticate. These steps are usually completed by an operator or configuration
management tool.

1. Enable the AppID auth method:

    ```text
    $ vault auth enable app-id
    ```

1. Configure it with the set of AppIDs, user IDs, and the mapping between them:

    ```text
    $ vault write auth/app-id/map/app-id/foo value=admins display_name=foo
    # ...

    $ vault write auth/app-id/map/user-id/bar value=foo cidr_block=10.0.0.0/16
    # ...
    ```

    This will create an AppID "foo" that associates with the policy "admins".
    The `display_name` sets the display name for audit logs and secrets. Next,
    we configure the user ID "bar" and say that the user ID bar can be paired
    with "foo" but only if the client is in the "10.0.0.0/16" CIDR block. The
    `cidr_block` configuration is optional.

    This means that if a client authenticates and provide both "foo" and "bar",
    then the app ID will authenticate that client with the policy "admins".

    In practice, both the user and app ID are likely hard-to-guess UUID-like
    values.

    Note that it is possible to authorize multiple app IDs with each user ID by
    writing them as comma-separated values to the user ID mapping:

    ```text
    $ vault write auth/app-id/map/user-id/bar value=foo,baz cidr_block=10.0.0.0/16
    # ...
    ```

## API

The AppID auth method has a full HTTP API. Please see the
[AppID auth method API](/api/auth/app-id/index.html) for more
details.
