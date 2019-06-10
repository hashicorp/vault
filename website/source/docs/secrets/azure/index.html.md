---
layout: "docs"
page_title: "Azure - Secrets Engine"
sidebar_title: "Azure"
sidebar_current: "docs-secrets-azure"
description: |-
  The Azure Vault secrets engine dynamically generates Azure
  service principals and role assignments.
---

# Azure Secrets Engine

The Azure secrets engine dynamically generates Azure service principals and role
assignments.  Vault roles can be mapped to one or more Azure roles, providing a
simple, flexible way to manage the permissions granted to generated service
principals.

Each service principal is associated with a Vault lease. When the lease expires
(either during normal revocation or through early revocation), the service
principal is automatically deleted.

If an existing service principal is specified as part of the role configuration,
a new password will be dynamically generated instead of a new service principal.
The password will be deleted when the lease is revoked.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the Azure secrets engine:

    ```text
    $ vault secrets enable azure
    Success! Enabled the azure secrets engine at: azure/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Configure the secrets engine with account credentials:

    ```text
    $ vault write azure/config \
    subscription_id=$AZURE_SUBSCRIPTION_ID \
    tenant_id=$AZURE_TENANT_ID \
    client_id=$AZURE_CLIENT_ID \
    client_secret=$AZURE_CLIENT_SECRET

    Success! Data written to: azure/config
    ```

    If you are running Vault inside an Azure VM with MSI enabled, `client_id` and
    `client_secret` may be omitted. For more information on authentication, see the [authentication](#authentication) section below.

1. Configure a role. A role may be set up with either an existing service principal, or
a set of Azure roles that will be assigned to a dynamically created service principal.

To configure a role called "my-role" with an existing service principal:

    ```text
    $ vault write azure/roles/my-role application_object_id=<existing_app_obj_id> ttl=1h
    ```

Alternatively, to configure the role to create a new service principal with Azure roles:

    ```text
    $ vault write azure/roles/my-role ttl=1h azure_roles=-<<EOF
        [
       	    {
       	        "role_name": "Contributor",
       	        "scope":  "/subscriptions/<uuid>/resourceGroups/Website"
       	    }
        ]
    EOF
    ```

Roles may also have their own TTL configuration that is separate from the mount's
TTL. For more information on roles see the [roles](#roles) section below.


## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permissions, it can generate credentials. The usage pattern is the same
whether an existing or dynamic service principal is used.

To generate a credential using the "my-role" role:

```text
$ vault read azure/creds/my-role

Key                Value
---                -----
lease_id           azure/creds/sp_role/1afd0969-ad23-73e2-f974-962f7ac1c2b4
lease_duration     60m
lease_renewable    true
client_id          408bf248-dd4e-4be5-919a-7f6207a307ab
client_secret      ad06228a-2db9-4e0a-8a5d-e047c7f32594
```

This endpoint generates a renewable set of credentials. The application can login
using the `client_id`/`client_secret` and will have access provided by configured service
principal or the Azure roles set in the "my-role" configuration.


## Roles

Vault roles let you configure either an existing service principal or a set of Azure roles, along with
role-specific TTL parameters. If an existing service principal is not provided, the configured Azure
roles will be assigned to a newly created service principal. The Vault role may optionally specify
role-specific `ttl` and/or `max_ttl` values. When the lease is created, the more restrictive of the
mount or role TTL value will be used.

### Application Object IDs
If an existing service principal is to be used, the Application Object ID must be set on the Vault role.
This ID can be found by inspecting the desired Application with the `az` CLI tool, or via the Azure Portal. Note
that the Application **Object** ID must be provided, not the Application ID.

### Azure Roles
If dynamic service principals are used, Azure roles must be configured on the Vault role.
Azure roles are provided as a JSON list, with each element describing an Azure role and scope to be assigned.
Azure roles may be specified using the `role_name` parameter ("Owner"), or `role_id`
("/subscriptions/.../roleDefinitions/...").
`role_id` is the definitive ID that's used during Vault operation; `role_name` is a convenience during
role management operations. All roles *must exist* when the configuration is written or the operation will fail. The role lookup priority is:

1. If `role_id` is provided, it validated and the corresponding `role_name` updated.
1. If only `role_name` is provided, a case-insensitive search-by-name is made, succeeding
only if *exactly one* matching role is found. The `role_id` field will updated with the matching role ID.

The `scope` must be provided for every role assignment.

Example of role configuration:

```text
$ vault write azure/roles/my-role ttl=1h max_ttl=24h azure_roles=-<<EOF
  [
    {
        "role_name": "Contributor",
    	"scope":  "/subscriptions/<uuid>/resourceGroups/Website"
    },
    {
        "role_id": "/subscriptions/<uuid>/providers/Microsoft.Authorization/roleDefinitions/<uuid>",
    	"scope":  "/subscriptions/<uuid>"
    },
    {
   	    "role_name": "This won't matter as it will be overwritten",
   	    "role_id": "/subscriptions/<uuid>/providers/Microsoft.Authorization/roleDefinitions/<uuid>",
   	    "scope":  "/subscriptions/<uuid>/resourceGroups/Database"
    }
  ]
EOF
```


## Authentication

The Azure secrets backend must have sufficient permissions to read Azure role information and manage
service principals. The authentication parameters can be set in the backend configuration or as environment variables. Environment variables will take precedence.
 The individual parameters are described in the [configuration][config] section of the API docs.

If the client ID or secret are not present and Vault is running on and Azure VM, Vault will attempt to use
[Managed Service Identity (MSI)](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview) to access Azure. Note that when MSI is used, tenant and subscription IDs must still be explicitly provided in the configuration or environment variables.

The following Azure roles and Azure Active Directory (AAD) permissions are required, regardless of which authentication method is used:

- "Owner" role for the subscription scope
- "Read and write all applications" permission in AAD

These permissions can be configured through the Azure Portal, CLI tool, or PowerShell.
In your Azure subscription, your account must have `Microsoft.Authorization/*/Write`
access to assign an AD app to a role. This action is granted through the [Owner](https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#owner) role or
[User Access Administrator](https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#user-access-administrator) role. If your account is assigned to the Contributor role, you
don't have adequate permission. You will receive an error when attempting to assign the service
principal to a role.

## Choosing between dynamic or existing service principals

Dynamic service principals are preferred if the desired Azure resources can be provided
via the RBAC system and Azure roles defined in the Vault role. This form of credential is
completely decoupled from any other clients, is not subject to permission changes after
issuance, and offers the best audit granularity.

Access to some Azure services cannot be provided with the RBAC system, however. In these
cases, an existing service principal can be set up with the necessary access, and Vault
can create new passwords for this service principal. Any changes to the service principal
permissions affect all clients. Furthermore, Azure does not provide any logging with
regard to _which_ credential was used for an operation.

An important limitation when using an existing service principal is that Azure limits the
number of passwords for a single Application. This limit is based on Application object
size and isn't firmly specified, but in practice hundreds of passwords can be issued per
Application. An error will be returned if the object size is reached. This limit can be
managed by reducing the role TTL, or by creating another Vault role against a different
Azure service principal configured with the same permissions.


## Additional Notes

-  **If a referenced Azure role doesn't exist, a credential will not be generated.**
  Service principals will only be generated if *all* role assignments are successful.
  This is important to note if you're using custom Azure role definitions that might be deleted
  at some point.

- Azure roles are assigned only once, when the service principal is created. If the
  Vault role changes the list of Azure roles, these changes will not be reflected in
  any existing service principal, even after token renewal.

- The time required to issue a credential is roughly proportional to the number of
  Azure roles that must be assigned. This operation make take some time (10s of seconds
  are common, and over a minute has been seen).

- Service principal credential timeouts are not used. Vault will revoke access by
  deleting the service principal.

- The Application Name for dynamic service principals will be prefixed with `vault-`. Similarly
  the `keyId` of any passwords added to an existing service principal will begin with
  `ffffff`. These may be used to search for Vault-created credentials using the `az` tool
  or Portal.

## Help &amp; Support

The Azure secrets engine is written as an external Vault plugin and
thus exists outside the main Vault repository. It is automatically bundled with
Vault releases, but the code is managed separately.

Please report issues, add feature requests, and submit contributions to the
[vault-plugin-secrets-azure repo][repo] on GitHub.


## API
The Azure secrets engine has a full HTTP API. Please see the [Azure secrets engine API docs][api]
for more details.

[api]: /api/secret/azure/index.html
[config]: /api/secret/azure/index.html#configure
[repo]: https://github.com/hashicorp/vault-plugin-secrets-azure
