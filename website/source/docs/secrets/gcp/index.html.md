---
layout: "docs"
page_title: "GCP - Secrets Engines"
sidebar_current: "docs-secrets-gcp"
description: |-
  The GCP secrets engine for Vault generates service account keys and OAuth tokens dynamically based 
  on IAM policies.
---

# GCP Secrets Engine

The Google Cloud Platform secrets engine dynamically generates IAM service account 
credentials based on IAM policies, allowing users to have access to GCP resources 
without having to manage a new identity (service account). 
      
## Quick Links:
* **[Things to Note](#things-to-note): If you are running into unexpected errors
    (or wish to avoid doing so), please read this section! In particular, if you have used other Vault secrets engines 
    (i.e. the AWS engine), this covers important differences in behavior.**
* Contributing/Issues: This engine was written as a Vault plugin and thus 
  exists outside the main Vault repo. Please report issues, request features, or 
  submit contributions to the [plugin-specific repo on Github](https://github.com/hashicorp/vault-plugin-secrets-gcp) 
  (all are welcome!)
* Docs Quick Links:
    * [GCP secrets engine HTTP API](/api/secret/gcp/index.html)
    * [Expected `roleset` bindings format](#roleset-bindings)


## Background

Benefits of using this Vault secrets engine to manage IAM service accounts and credentials:

* Automatic cleanup of long-lived IAM service account keys
    * Vault associates leases with generated IAM service account keys s.t. they will
      be automatically revoked when the lease expires. 
* Users don't need to create a new service-account per one-off or short-term access.
* Multi-cloud applications:
    * Users authenticate to Vault using some central identity (LDAP, AppRole, etc) 
      and can generate GCP credentials without having to create and manage a new 
      service account for that user. 


## Overview

*Note*: The following docs assume that the secrets engine has been mounted at `gcp`, 
the default mount path.  Adjust your calls accordingly if you mount the engine at a 
different path. All example calls in this doc will use the Vault CL tool - see [HTTP API](/api/secret/gcp/index.html)
for possible API calls.

### Setup

Initially, the secrets engine must be enabled in Vault.

```sh
$ vault secrets enable gcp
Success! Enabled the gcp secrets engine at: gcp/
```

You will need to set up the secrets engine with:

* **[Config](#config)**: General config, including credentials that Vault will need to make calls 
    to GCP APIs (either explicitly or using Application Default Credentials), lease defaults, etc. Example:

* **[Rolesets](#rolesets)**: Generated credentials will need to be associated with sets of 
    [IAM roles](https://cloud.google.com/iam/docs/understanding-roles) on specific GCP resources. Rolesets
    have an associated `secret_type` that determines a secret type that can be generated 
    under this role set. Example:

### Usage - Secret Generation

Once the secrets engine has been set up, you can generate two types of secrets in Vault:

* **[Access Tokens](#access-tokens)** (`secret_type`: `access_token`): 
    This endpoint generates non-renewable OAuth2 access tokens with a lifetime of an hour. 
    
* **[Service Account Keys](#service-account-keys)** (`secret_type`: `service_account_key`): This endpoint generates long-lived 
    [IAM service account keys](https://cloud.google.com/iam/docs/service-accounts#service_account_keys).


Each secret is associated with a [Vault lease](docs/concepts/lease.html) that
can be revoked and possibly renewed (see lease docs for how to do so). On revoking
a lease, the secret is no longer guaranteed to work.

Each secret is generated under a specified role set, which determine which permissions (IAM roles)
the generated credentials have on specific GCP resources. Note that a role set can only generate one
type of secret, specified at role set creation as `secret_type`


## Config

The `config/` endpoint is used to configure any information shared by the 
entire GCP secrets engine. You can 

Allowed Operations: `write`, `read`

Examples:

```bash
$ vault write gcp/config credentials="..." ttl=100 max_ttl=1000
$ vault read gcp/config
```

Parameters For Write:

* `credentials` (`string:""`) - JSON credentials (either file contents or '@path/to/file'). 
    See next sections to learn more about required permissions and alternative ways to provide these credentials.
* `ttl` (`int: 0 || string:"0s"`) – Specifies default config TTL for long-lived credentials 
    (i.e. service account keys). Accepts integer number of seconds or Go duration format string.
* `max_ttl` (`int: 0 || string:"0s"`)– Specifies default config TTL for long-lived credentials
    (i.e. service account keys). Accepts integer number of seconds or Go duration format string.

#### Passing Credentials To Vault
If you would rather not pass the IAM credentials using the payload `credentials` parameter,
there are multiple ways to pass IAM credentials to the Vault server. You can specify credentials
in the following ways (given in order of evaluation):

1. Static credential JSON provided to the API as a payload (i.e. `credentials` parameter)
2. Credentials in the environment variables `GOOGLE_CREDENTIALS` or `GOOGLE_CLOUD_KEYFILE_JSON`
3. Parse JSON file ~/.gcp/credentials
4. [Google Application Default Credentials](https://cloud.google.com/docs/authentication/production)


#### Required Permissions 

At present, this endpoint does not confirm that the provided GCP credentials
are for existing IAM users or have valid permissions. However, they will need
the following permissions:

```
// Service Account + Key Admin
    iam.serviceAccounts.create
    iam.serviceAccounts.delete
    iam.serviceAccounts.get
    iam.serviceAccounts.list
    iam.serviceAccountKeys.create
    iam.serviceAccountKeys.delete
    iam.serviceAccountKeys.get
    iam.serviceAccountKeys.list
    iam.serviceAccounts.update

// IAM Policy Changes
    $service.$resource.getIamPolicy
    $service.$resource.setIamPolicy
```

where `$service.$resource` refers to the GCP service/resources you wish to
assign credentials permissions on, e.g. 
`cloudresourcemanager.projects.get/setIamPolicy` for GCP projects

You can either create a [custom role](https://cloud.google.com/iam/docs/creating-custom-roles) 
with these permissions and assign it at a project-level to an IAM service account for Vault to use, or 
assign `roles/iam.serviceAccountAdmin` and `roles/iam.serviceAccountKeyAdmin` to get the permissions other 
than `get/setIamPolicy` permissions, and assign custom roles or global roles that include `get/setIamPolicy`
permissions for the resources that Vault will change policies for.


## Rolesets

Rolesets determine the permissions that service account credentials (secrets) 
generated by Vault will have on given GCP resources.

Endpoint:  `roleset/` (`rolesets/` for list)

Allowed Operations: `write`, `read`, `list`, `delete`, `rotate`, `rotate-key`

### Creating/Updating Rolesets
Each roleset will have an associated service account that is *created when the role set is created or updated*.

(see [things to note](#things-to-note) for background)

Parameters For Write:

* `name` (`string: <required>`): Required. Name of the role. Cannot be updated. Given as part of path.
* `secret_type` (`string: "access_token"`): Type of secret generated for this role set. 
    Accepted values: `access_token`, `service_account_key`. Cannot be updated.
* `project` (`string: <required>`): Name of the GCP project that this roleset's service account will belong to. 
    Cannot be updated. 
* `bindings` (`string: <required>`): Bindings configuration string (expected HCL or JSON string, raw or base64-encoded)
* `token_scopes` (`array: []`): List of OAuth scopes to assign to `access_token` secrets generated under 
    this role set (`access_token` role sets only)

Update is the same call but will error if non-updatable fields are given with different values.
If you update a roleset's bindings, note this will effectively revoke any secrets generated under 
this roleset.

Examples:

```bash
# secret type `access_token`
$ vault write gcp/roleset/my-token-roleset \
    project="mygcpproject" \
    secret_type="access_token"  \
    token_scopes="https://www.googleapis.com/auth/cloud-platform"
    bindings=@binds.hcl \
         
# secret type `service_account_key`       
$ vault write gcp/roleset/my-key-roleset \
    project="mygcpproject" \
    secret_type="service_account_key"  \
    bindings="<base64-encoded-hcl-string>"
```

#### Roleset Bindings

The rolesets `binding` argument accepts bindings in the following format:

```hcl
resource "path/to/my/resource" {
  roles = [
    "roles/viewer",
    "projects/X/roles/myprojCustomRole"
  ] 
}

resource "//service.googleapis.com/path/to/another/resource" {
  roles = [
    "organizations/Y/roles/myprojCustomRole"
  ] 
}
```

The top-level blocks are `resource` blocks, defined as `resource "a/resource/name" {...}`. Each block
define IAM policy information to be bound to this resource. 
    
The following resource path formats are supported:

1. Project-level Self-Link (Specify Service and Version)
    A URI with scheme and host. Generally the `self_link` attribute of some resource. 
    Must be resource with parent project (i.e. relative resource name `projects/$PROJECT/...`). 
    Examples:
       * Compute alpha zone: 
       
       ```
       https://www.googleapis.com/compute/alpha/projects/my-project/zones/us-central1-c
       ```

2. [Full resource name](https://cloud.google.com/apis/design/resource_names#full_resource_name) (Specify Service): 
    A scheme-less URI consisting of a DNS-compatible API service name and a resource path. 
    The resource path is the relative resource name (see next). Use to specify service but use 
    either the preferred service version or the only version for which this resource is IAM-enabled. 

    Examples:
    * Compute snapshot: 
    
    ```
    //compute.googleapis.com/project/X/snapshots/Y
    ``` 
    * Pubsub snapshot:
    
    ```
    //pubsub.googleapis.com/project/X/snapshots/Y
    ```
    
3. [Relative Resource Name](https://cloud.google.com/apis/design/resource_names#relative_resource_name): 
    A [path-noscheme](https://tools.ietf.org/html/rfc3986#appendix-A) URI path. Use if version/service are
    apparent from resource type (or you want to use only the preferred version of the service). General format is resource name
    as accepted by the corresponding REST API, but we've added some exceptions (namely Storage). 
    Examples: 
    * Storage bucket object: 
    
    ```
    b/bucketname/o/objectname
    buckets/bucketname/o/objectname 
    ```
    * Pubsub Topic:
     
     ```
     projects/X/topics/Y
     ```
            
Each `resource` block accepts the following arguments:

 * `roles`: An array of string names for [IAM roles](https://cloud.google.com/iam/docs/understanding-roles).
    Each string must be a global role name (`roles/roleFoo`), a project-level custom role 
    (`projects/myproj/roles/roleFoo`) or an organization-level custom role (`organizations/myorg/roles/roleFoo`)  
    
You can provide this as a plaintext string blob, the base64-encoded version of this string,
or using syntax `@path/to/bindings.hcl` to pass in a filename
    
#### Creation Workflow:

When an admin user creates a roleset or updates the role set bindings, Vault does the following:

1. Parses given bindings configuration file as described above.

2. Add [WAL](https://en.wikipedia.org/wiki/Write-ahead_logging) entries to clean-up:
    - Current service account, key, and bindings in roleset if updating and update succeeds
    - New service account, key, and bindings if create/update fails.
    
    
3. Attempt to create a new IAM service account.
    * The new service account email is `"vault$rolesetName-$timestamp@..."`, where the roleset name might be 
        truncated to fit the IAM service account character limit)
    * The new service account display name will contain the Vault roleset it belongs to. 

4. For each resource in the bindings
    
    a. call `getIamPolicy` method on the resource.
    
    b. For each role in the resource's role list, add a binding to this policy with the email of the new
        service account.
        
        { "role": "roles/theRoleToAssign", "member": "serviceAccount:$EMAIL" }
    
    c. call `setIamPolicy` method on the resource
    
    See [this section](#calling-iam-methods) to learn more about how we call IAM methods on various resources.  

    **Note: Because the service account and policies are assigned now, the
    resources given must exist at roleset creation time.**

5. If the roleset generates access tokens (`secret_type = "access_token"`), create a new service account key for
    this service account. 

6. Save the new role set.

    **Note:** If steps 2-6 fail, we return a error response and the new service account will not be saved. WAL entries
    will have been added s.t. any entities will eventually get cleaned up, but you may need to manually delete these
    if rollback fails and you want to free up quota immediately. 
    
    After step 6, the create/update operation will return a success, but we have a last cleanup step:

7. Try to delete the old service account (and any keys) and remove its old IAM policy bindings. 
    If any of these calls to GCP APIs fail, we added the WAL entries in step 2 so eventually these resources 
    will get cleaned up. We will return a warning in the response. In addition, we will try to delete the WAL entries
    for the now in-use service account that we just created s.t. it does not attempt to clean up these entities later.
     
    
#### WAL Cleanup

As mentioned, Vault will create WAL entries, both for newly created service accounts/bindings, 
which may not get used if update fails, and for old service account bindings, which may need to be deleted 
if immediate cleanup fail. 

Each WAL entry contains the name of the role set that was under update/creation. The WAL cleanup functions 
work as follows:

1. Vault will attempt to get the role set (which may be missing because it was deleted or was never created). 

2. For each type of WAL entry:
    * **Service Account**: We try to delete the service account saved in the WAL entry.
        * If the roleset is still using the service account: We do nothing and remove the WAL entry. 
          This happens if we failed to delete a WAL entry for a successful operation, or we preemptively 
          added a WAL entry to delete old service accounts for an update that ended up failing. 
    * **Service Account Key** (for `access_token` rolesets)
        * If a key name is included in the WAL entry:
            * We are attempting to clean up a previously created and saved key after a role set update.
            * If the roleset still exists and uses the key, we do nothing and remove the WAL entry.
              Otherwise, we delete the key.
        * If no key name is included:
            * We are attempting to clean up newly created keys for failed roleset updates. Because the key 
                name is generated by GCP, we do not know the key name when we create the WAL entry.
            * Instead of deleting one key, we list all the user-managed keys under this service account 
                and delete any that are not being used by the current roleset, if it exists. 
               Since this service account is being controlled by Vault and the secrets are access tokens, 
               there should not be any other user-managed keys for this service account.
    * **IAM Resource-Policy Bindings**: 
        * The entry contains the GCP resource name that may need policy bindings cleaned up, 
            and the bound roles and service account email that needs to be cleaned up. 
            We get the IAM policy, remove the necessary bindings, and set the IAM policy without these bindings.
        * If the roleset exists and is still using this service account email, we remove only IAM policy bindings 
            not included in the current roleset's bindings.

### Other Roleset Operations

#### Rotation

If you want to reset the service account for a given roleset, or rotate the key used for `access_token` rolesets,
the GCP secrets engine has two endpoints for this:

Examples:

```bash
$ vault write gcp/roleset/my-key-roleset/rotate \
      
$ vault write gcp/roleset/my-token-roleset/rotate
$ vault write gcp/roleset/my-token-roleset/rotate-key # `access_token` only
```

The `rotate/` endpoint works similar to roleset update (i.e. see the create/update workflow above).
It replaces the roleset's service account with a new service account and replaceds any bindings with the new
email as the IAM policy binding member. This will in effect revoke any secrets generated under this roleset.

The `rotate-key/` endpoint only applies to `access_token` rolesets and simply rotates the key
saved in the roleset used to generate access tokens, while keeping the same service account. 

Note that rotating the service account (`/rotate`) for an `access_token` roleset 
will also effectively rotate the key.

#### Read and Delete

Examples:

```bash
$ vault read gcp/roleset/my-token-roleset
         
$ vault delete gcp/roleset/my-key-roleset
```

On read, you can get the associated service account, token scopes for `access_token` rolesets,
and other information. 

On delete, the service account, bindings, and any keys associated with the roleset service account
will be deleted.


## Generating Secrets

### Access Tokens

You can only generate tokens under `access_token` rolesets.

Example Call: 

```bash
$ vault read gcp/token/my-token-roleset

Key                Value
---                -----
lease_id           gcp/token/my-token-roleset/some-uuid
lease_duration     59m59s
lease_renewable    false
token              ya29.c.Ell9...

```

This endpoint generates non-renewable OAuth2 access tokens with a lifetime of an hour. Vault will simply
return an access token (as `token` in output) that can be used in calls to GCP APIs as part of an
"Authorization: Bearer" header:

```bash

$ curl -H "Authorization: Bearer $TOKEN" ...

```    
### Service Account Keys

You can only generate tokens under `service_account_key` rolesets.

Example Call: 

```bash
$ vault read gcp/key/my-key-roleset

lease_id            gcp/key/my-key-roleset/some-uuid
lease_duration      768h
lease_renewable     true
key_algorithm       KEY_ALG_RSA_2048
key_type            TYPE_GOOGLE_CREDENTIALS_FILE
private_key_data    <base-64 encoded private key data>...
```

This endpoint generates [IAM service account keys](https://cloud.google.com/iam/docs/service-accounts#service_account_keys) 
associated with the role set's service account. These keys are by default long-lived in GCP, but Vault 
associates a lease (see output `lease_*` information). This lease can be renewed to extend the key lifetime,
or revoked. When the lease is revoked, the service account key will be deleted (thus freeing up quota).

As mentioned in [things to note](#things-to-note), there is a limit of 10 service account keys per
account and thus a limit of 10 secrets per `service_account_key` roleset. Read this section to learn about
how to get around this limit if you are running into issues.

To learn how to use service account keys, see 
[calling Google APIs](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#callinganapi),
[how to authenticate in gcloud](https://cloud.google.com/sdk/gcloud/reference/auth/activate-service-account), or
[how to authenticate in Google Cloud client libraries](https://cloud.google.com/docs/authentication/getting-started).


## Things to Note

### Access Tokens vs Service Account Keys

While in general the IAM team prefers that you use short-term access tokens over long-lived, hard-to-manage service account keys 
to authenticate calls to GCP APIs, this is currently impractical as several Google tools
require service account keys for authenticating calls (i.e. `gcloud` CLI tool, client libraries).
Thus, while we default to roleset `secret_type` = `access_token`, we also offer the option
to generate keys (set role set `secret_type` to `service_account_key`), but we caution that
these keys should still be carefully monitored even if they have been associated with
Vault leases. 

**NOTE**: `access_token` role sets will generate a service account key that is saved in Vault. Vault will be the only 
place this key can be accessed (i.e. console or API users cannot access the private key data) and this data will 
not be returned in output for reading the role set. However, **this key can and should be [rotated](#rotation) as 
necessary to avoid possible theft.** Note that rotation does effectively invalidate all secrets previously generated
under the roleset.

### Service account are created **on roleset update/creation** rather than per secret.

(**NOTE: This is different than AWS!!**)

There are a couple of reasons why we want service account creation to happen during role-set creation/updates (i.e.
during server set up) rather than on the fly, per secret generated: 

* IAM service account creation and permissions propagation can take up to a minute. Because the  service account is
    created and permissions are assigned during set-up, secrets can be used immediately instead of after a 
    possible delay. This can make automated workflows slightly less complicated/flaky.  
* **Service Account Quotas**: GCP projects by default have a limit on the number of IAM service accounts you can create 
    (currently 100, *including system-managed service accounts*). You can 
    [request additional quota](https://cloud.google.com/compute/quotas), which is better done
    during a set-up step.
    * If service accounts are created per secret, this limit would instead be on the number of issued secrets.
    
    
However, there are a couple of caveats that come from generating per-role-set as well:

* **Service account key limit**: 
    Because GCP IAM has set a hard limit (currently 10) on the number of service 
    account keys, attempts to generate more keys than the limit will result in an error - 
    Vault is notified *on the fly* when a user tries to create a new secret 
    and the GCP IAM service returns an error. Thus, if you are generating service account 
    keys, *for each `service_account_key` role set, there is a hard limit of of 
    10 secrets (keys)** 
    If you find yourself running into this limit, consider:
    * Having shorter TTLs or revoking access earlier: If you're not using service 
          account keys created earlier, consider rotating and freeing quota earlier.
    * Creating more role sets with the same set of permissions: Additional role-sets can 
            always be created with the same set of permissions, which adds service accounts and
            thus effectively increases the number of key secrets you can generate. 
    * Using `access_token` role sets: If you need several, very short-term accesses to GCP
            resources, consider instead requesting `access_token` secrets which have no limit and
            are naturally short-lived.
            
* **Resources must exist at role set creation time.** Because we set the bindings for the service account at this point,
    if the resource does not exist, calls to getIamPolicy will fail on the resource and the role set creation will fail.
    
* **Role-set creation might take a while and partially fail.** Because every service account creation, key creation, and 
    IAM policy change per resource is a GCP API call and we can't do transactions through GCP, if one call fails, the roleset
    creation fails and Vault will attempt to rollback all changes made. However, these rollback calls themselves are 
    API calls that also might fail. We add [WAL entries](#wal-cleanup) to ensure that any unused entities or bindings
    get cleaned up eventually, but you may run into possible situations where you hit quota limits or have unused 
    IAM service accounts that have not been cleaned up. If you manually do this cleanup, the WAL entry will just
    succeed without doing anything, so feel free to do this (carefully) if you need to immediately clean up issues.

### External Changes to Vault-owned IAM Accounts
While Vault will initially create
and assign permissions to IAM service accounts, it is possible that an external user deletes this service account 
and/or Vault-managed keys (for `access_token` role sets), or adjusts this service account's permissions. 
**We will deny secret generation in the first case until the role set has been rotated or updated, 
but the second case is hard to detect on the fly.**

* If your credentials has unexpected permissions, consider the cases that the service account has been
  either assigned new external permissions or has possibly inherited permissions from a parent resource policy.
    * Consider just rotating this service account periodically anyways to avoid tampering. 
* In general, you should not be changing these service accounts via console/API if you are using Vault. 
    Please warn your GCP project owners to avoid accidentally changing these Vault roleset service accounts.
    * Vault role set accounts have emails with the format `vault<roleset-prefix>-<creation UNIX timestamp>@...`, 
        where `roleset-prefix` is the roleset name, possibly truncated to fit the character limit on rolesets. The
        display name (description) will also have the full Vault roleset name. 
* If the service account/key has been deleted, you will need to regenerate the role set account/key using the 
    `gcp/$roleset/rotate` or `gcp/$roleset/rotate-key` endpoints. Updates to the role set bindings
     will also trigger service account recreation.
     

### Calling IAM Methods

An IAM-enabled resource (under an arbitrary GCP service) supports the following three IAM methods:

* `getIamPolicy`
* `setIamPolicy`
* `testIamPermissions`

In the case of this secrets engine, we need to call `getIamPolicy` and `setIamPolicy` on
an arbitrary resource under an arbitrary service, which would be difficult using
the [generated Go google APIs](https://github.com/google/google-api-go-client). Instead,
we autogenerated a library, using the [Google API Discovery Service](https://developers.google.com/discovery/)
to find IAM-enabled resources and configure HTTP calls on arbitrary services/resources for IAM.

For each binding config resource block (with a resource name), we attempt to find the resource type based on the 
relative resource name and match it to a service config as seen in this 
[autogenerated config file](https://github.com/hashicorp/vault-plugin-secrets-gcp/blob/master/plugin/iamutil/iam_resources_generated.go)

To re-generate this file, run: 

```
go generate github.com/hashicorp/vault-plugin-secrets-gcp/plugin/iamutil
```


In general, we try to make it so you can specify the resource as given in the HTTP API URL 
(between base API URL and get/setIamPolicy suffix). For some possibly non-standard APIs, we have also
 added exceptions to try to reach something more standard; a notable current example is the Google Cloud Storage API, 
 whose methods look like `https://www.googleapis.com/storage/v1/b/bucket/o/object` where we accept either 
 `b/bucket/o/object` or `buckets/bucket/objects/object` as valid relative resource names.

If you are having trouble during role set creation with errors suggesting the resource format is invalid or API calls
are failing for a resource you know exists, please [report any issues](https://github.com/hashicorp/vault-plugin-secrets-gcp/issues) 
you run into. It could be that the API is a non-standard form or we need to re-generate our config file.

## API
The GCP secrets engine has a full HTTP API. Please see the [GCP secrets engine API docs](/api/secret/gcp/index.html) 
for more details.