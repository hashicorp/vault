## 0.11.6 (December 14th, 2018)

This release contains the three security fixes from 1.0.0 and 1.0.1 and the
following bug fixes from 1.0.0/1.0.1:

 * namespaces: Correctly reload the proper mount when tuning or reloading the
   mount [[GH-5937](https://github.com/hashicorp/vault/pull/5937)]
 * replication/perfstandby: Fix audit table upgrade on standbys [[GH-5811](https://github.com/hashicorp/vault/pull/5811)]
 * replication/perfstandby: Fix redirect on approle update [[GH-5820](https://github.com/hashicorp/vault/pull/5820)]
 * secrets/kv: Fix issue where storage version would get incorrectly downgraded
   [[GH-5809](https://github.com/hashicorp/vault/pull/5809)]

It is otherwise identical to 0.11.5.

## 0.11.5 (November 13th, 2018)

BUG FIXES:

 * agent: Fix issue when specifying two file sinks [[GH-5610](https://github.com/hashicorp/vault/pull/5610)]
 * auth/userpass: Fix minor timing issue that could leak the presence of a
   username [[GH-5614](https://github.com/hashicorp/vault/pull/5614)]
 * autounseal/alicloud: Fix issue interacting with the API (Enterprise)
 * autounseal/azure: Fix key version tracking (Enterprise)
 * cli: Fix panic that could occur if parameters were not provided [[GH-5603](https://github.com/hashicorp/vault/pull/5603)]
 * core: Fix buggy behavior if trying to remount into a namespace
 * identity: Fix duplication of entity alias entity during alias transfer
   between entities [[GH-5733](https://github.com/hashicorp/vault/pull/5733)]
 * namespaces: Fix tuning of auth mounts in a namespace
 * ui: Fix bug where editing secrets as JSON doesn't save properly [[GH-5660](https://github.com/hashicorp/vault/pull/5660)]
 * ui: Fix issue where IE 11 didn't render the UI and also had a broken form
   when trying to use tool/hash [[GH-5714](https://github.com/hashicorp/vault/pull/5714)]

## 0.11.4 (October 23rd, 2018)

CHANGES:

 * core: HA lock file is no longer copied during `operator migrate` [[GH-5503](https://github.com/hashicorp/vault/pull/5503)].
   We've categorized this as a change, but generally this can be considered
   just a bug fix, and no action is needed.

FEATURES:

 * **Transit Key Trimming**: Keys in transit secret engine can now be trimmed to
   remove older unused key versions
 * **Web UI support for KV Version 2**: Browse, delete, undelete and destroy
   individual secret versions in the UI
 * **Azure Existing Service Principal Support**: Credentials can now be generated
   against an existing service principal

IMPROVEMENTS:

 * core: Add last WAL in leader/health output for easier debugging [[GH-5523](https://github.com/hashicorp/vault/pull/5523)]
 * identity: Identity names will now be handled case insensitively by default.
   This includes names of entities, aliases and groups [[GH-5404](https://github.com/hashicorp/vault/pull/5404)]
 * secrets/aws: Added role-option max_sts_ttl to cap TTL for AWS STS
   credentials [[GH-5500](https://github.com/hashicorp/vault/pull/5500)]
 * secret/database: Allow Cassandra user to be non-superuser so long as it has
   role creation permissions [[GH-5402](https://github.com/hashicorp/vault/pull/5402)]
 * secret/radius: Allow setting the NAS Identifier value in the generated
   packet [[GH-5465](https://github.com/hashicorp/vault/pull/5465)]
 * secret/ssh: Allow usage of JSON arrays when setting zero addresses [[GH-5528](https://github.com/hashicorp/vault/pull/5528)]
 * secret/transit: Allow trimming unused keys [[GH-5388](https://github.com/hashicorp/vault/pull/5388)]
 * ui: Support KVv2 [[GH-5547](https://github.com/hashicorp/vault/pull/5547)], [[GH-5563](https://github.com/hashicorp/vault/pull/5563)]
 * ui: Allow viewing and updating Vault license via the UI
 * ui: Onboarding will now display your progress through the chosen tutorials
 * ui: Dynamic secret backends obfuscate sensitive data by default and
   visibility is toggleable

BUG FIXES:

 * agent: Fix potential hang during agent shutdown [[GH-5026](https://github.com/hashicorp/vault/pull/5026)]
 * auth/ldap: Fix listing of users/groups that contain slashes [[GH-5537](https://github.com/hashicorp/vault/pull/5537)]
 * core: Fix memory leak during some expiration calls [[GH-5505](https://github.com/hashicorp/vault/pull/5505)]
 * core: Fix generate-root operations requiring empty `otp` to be provided
   instead of an empty body [[GH-5495](https://github.com/hashicorp/vault/pull/5495)]
 * identity: Remove lookup check during alias removal from entity [[GH-5524](https://github.com/hashicorp/vault/pull/5524)]
 * secret/pki: Fix TTL/MaxTTL check when using `sign-verbatim` [[GH-5549](https://github.com/hashicorp/vault/pull/5549)]
 * secret/pki: Fix regression in 0.11.2+ causing the NotBefore value of
   generated certificates to be set to the Unix epoch if the role value was not
   set, instead of using the default of 30 seconds [[GH-5481](https://github.com/hashicorp/vault/pull/5481)]
 * storage/mysql: Use `varbinary` instead of `varchar` when creating HA tables
   [[GH-5529](https://github.com/hashicorp/vault/pull/5529)]

## 0.11.3 (October 8th, 2018)

SECURITY:

 * Revocation: A regression in 0.11.2 (OSS) and 0.11.0 (Enterprise) caused
   lease IDs containing periods (`.`) to not be revoked properly. Upon startup
   when revocation is tried again these should now revoke successfully.

IMPROVEMENTS:

 * auth/ldap: Listing of users and groups return absolute paths [[GH-5537](https://github.com/hashicorp/vault/pull/5537)]
 * secret/pki: OID SANs can now specify `*` to allow any value [[GH-5459](https://github.com/hashicorp/vault/pull/5459)]

BUG FIXES:

 * auth/ldap: Fix panic if specific values were given to be escaped [[GH-5471](https://github.com/hashicorp/vault/pull/5471)]
 * cli/auth: Fix panic if `vault auth` was given no parameters [[GH-5473](https://github.com/hashicorp/vault/pull/5473)]
 * secret/database/mongodb: Fix panic that could occur at high load [[GH-5463](https://github.com/hashicorp/vault/pull/5463)]
 * secret/pki: Fix CA generation not allowing OID SANs [[GH-5459](https://github.com/hashicorp/vault/pull/5459)]

## 0.11.2 (October 2nd, 2018)

CHANGES:

 * `sys/seal-status` now includes an `initialized` boolean in the output. If
   Vault is not initialized, it will return a `200` with this value set `false`
   instead of a `400`.
 * `passthrough_request_headers` will now deny certain headers from being
   provided to backends based on a global denylist.
 * Token Format: Tokens are now represented as a base62 value; tokens in
   namespaces will have the namespace identifier appended. (This appeared in
   Enterprise in 0.11.0, but is only in OSS in 0.11.2.)

FEATURES:

 * **AWS Secret Engine Root Credential Rotation**: The credential used by the AWS
   secret engine can now be rotated, to ensure that only Vault knows the
   credentials it is using [[GH-5140](https://github.com/hashicorp/vault/pull/5140)]
 * **Storage Backend Migrator**: A new `operator migrate` command allows offline
   migration of data between two storage backends
 * **AliCloud KMS Auto Unseal and Seal Wrap Support (Enterprise)**: AliCloud KMS can now be used a support seal for
   Auto Unseal and Seal Wrapping

BUG FIXES:

 * auth/okta: Fix reading deprecated `token` parameter if a token was
   previously set in the configuration [[GH-5409](https://github.com/hashicorp/vault/pull/5409)]
 * core: Re-add deprecated capabilities information for now [[GH-5360](https://github.com/hashicorp/vault/pull/5360)]
 * core: Fix handling of cyclic token relationships [[GH-4803](https://github.com/hashicorp/vault/pull/4803)]
 * storage/mysql: Fix locking on MariaDB [[GH-5343](https://github.com/hashicorp/vault/pull/5343)]
 * replication: Fix DR API when using a token [[GH-5398](https://github.com/hashicorp/vault/pull/5398)]
 * identity: Ensure old group alias is removed when a new one is written [[GH-5350](https://github.com/hashicorp/vault/pull/5350)]
 * storage/alicloud: Don't call uname on package init [[GH-5358](https://github.com/hashicorp/vault/pull/5358)]
 * secrets/jwt: Fix issue where request context would be canceled too early
 * ui: fix need to have update for aws iam creds generation [GF-5294]
 * ui: fix calculation of token expiry [[GH-5435](https://github.com/hashicorp/vault/pull/5435)]

IMPROVEMENTS:

 * auth/aws: The identity alias name can now configured to be either IAM unique
   ID of the IAM Principal, or ARN of the caller identity [[GH-5247](https://github.com/hashicorp/vault/pull/5247)]
 * auth/cert: Add allowed_organizational_units support [[GH-5252](https://github.com/hashicorp/vault/pull/5252)]
 * cli: Format TTLs for non-secret responses [[GH-5367](https://github.com/hashicorp/vault/pull/5367)]
 * identity: Support operating on entities and groups by their names [[GH-5355](https://github.com/hashicorp/vault/pull/5355)]
 * plugins: Add `env` parameter when registering plugins to the catalog to allow
   operators to include environment variables during plugin execution. [[GH-5359](https://github.com/hashicorp/vault/pull/5359)]
 * secrets/aws: WAL Rollback improvements [[GH-5202](https://github.com/hashicorp/vault/pull/5202)]
 * secrets/aws: Allow specifying STS role-default TTLs [[GH-5138](https://github.com/hashicorp/vault/pull/5138)]
 * secrets/pki: Add configuration support for setting NotBefore [[GH-5325](https://github.com/hashicorp/vault/pull/5325)]
 * core: Support for passing the Vault token via an Authorization Bearer header [[GH-5397](https://github.com/hashicorp/vault/pull/5397)]
 * replication: Reindex process now runs in the background and does not block other
   vault operations
 * storage/zookeeper: Enable TLS based communication with Zookeeper [[GH-4856](https://github.com/hashicorp/vault/pull/4856)]
 * ui: you can now init a cluster with a seal config [[GH-5428](https://github.com/hashicorp/vault/pull/5428)]
 * ui: added the option to force promote replication clusters [[GH-5438](https://github.com/hashicorp/vault/pull/5438)]
 * replication: Allow promotion of a secondary when data is syncing with a "force" flag

## 0.11.1.1 (September 17th, 2018) (Enterprise Only)

BUG FIXES:

 * agent: Fix auth handler-based wrapping of output tokens [[GH-5316](https://github.com/hashicorp/vault/pull/5316)]
 * core: Properly store the replication checkpoint file if it's larger than the
   storage engine's per-item limit
 * core: Improve WAL deletion rate
 * core: Fix token creation on performance standby nodes
 * core: Fix unwrapping inside a namespace
 * core: Always forward tidy operations from performance standby nodes

IMPROVEMENTS:

 * auth/aws: add support for key/value pairs or JSON values for
   `iam_request_headers` with IAM auth method [[GH-5320](https://github.com/hashicorp/vault/pull/5320)]
 * auth/aws, secret/aws: Throttling errors from the AWS API will now be
   reported as 502 errors by Vault, along with the original error [[GH-5270](https://github.com/hashicorp/vault/pull/5270)]
 * replication: Start fetching during a sync from where it previously errored

## 0.11.1 (September 6th, 2018)

SECURITY:

 * Random Byte Reading in Barrier: Prior to this release, Vault was not
   properly checking the error code when reading random bytes for the IV for
   AES operations in its cryptographic barrier. Specifically, this means that
   such an IV could potentially be zero multiple times, causing nonce re-use
   and weakening the security of the key. On most platforms this should never
   happen because reading from kernel random sources is non-blocking and always
   successful, but there may be platform-specific behavior that has not been
   accounted for. (Vault has tests to check exactly this, and the tests have
   never seen nonce re-use.)

FEATURES:

 * AliCloud Agent Support: Vault Agent can now authenticate against the
   AliCloud auth method.
 * UI: Enable AliCloud auth method and Azure secrets engine via the UI.

IMPROVEMENTS:

 * core: Logging level for most logs (not including secrets/auth plugins) can
   now be changed on-the-fly via `SIGHUP`, reading the desired value from
   Vault's config file [[GH-5280](https://github.com/hashicorp/vault/pull/5280)]

BUG FIXES:

 * core: Ensure we use a background context when stepping down [[GH-5290](https://github.com/hashicorp/vault/pull/5290)]
 * core: Properly check error return from random byte reading [[GH-5277](https://github.com/hashicorp/vault/pull/5277)]
 * core: Re-add `sys/` top-route injection for now [[GH-5241](https://github.com/hashicorp/vault/pull/5241)]
 * core: Policies stored in minified JSON would return an error [[GH-5229](https://github.com/hashicorp/vault/pull/5229)]
 * core: Evaluate templated policies in capabilities check [[GH-5250](https://github.com/hashicorp/vault/pull/5250)]
 * identity: Update MemDB with identity group alias while loading groups [[GH-5289](https://github.com/hashicorp/vault/pull/5289)]
 * secrets/database: Fix nil pointer when revoking some leases [[GH-5262](https://github.com/hashicorp/vault/pull/5262)]
 * secrets/pki: Fix sign-verbatim losing extra Subject attributes [[GH-5245](https://github.com/hashicorp/vault/pull/5245)]
 * secrets/pki: Remove certificates from store when tidying revoked
   certificates and simplify API [[GH-5231](https://github.com/hashicorp/vault/pull/5231)]
 * ui: JSON editor will not coerce input to an object, and will now show an
   error about Vault expecting an object [[GH-5271](https://github.com/hashicorp/vault/pull/5271)]
 * ui: authentication form will now default to any methods that have been tuned
   to show up for unauthenticated users [[GH-5281](https://github.com/hashicorp/vault/pull/5281)]


## 0.11.0 (August 28th, 2018)

DEPRECATIONS/CHANGES:

 * Request Timeouts: A default request timeout of 90s is now enforced. This
   setting can be overwritten in the config file. If you anticipate requests
   taking longer than 90s this setting should be updated before upgrading.
 * (NOTE: will be re-added into 0.11.1 as it broke more than anticipated. There
   will be some further guidelines around when this will be removed again.)
   * `sys/` Top Level Injection: For the last two years for backwards
   compatibility data for various `sys/` routes has been injected into both the
   Secret's Data map and into the top level of the JSON response object.
   However, this has some subtle issues that pop up from time to time and is
   becoming increasingly complicated to maintain, so it's finally being
   removed.
 * Path Fallback for List Operations: For a very long time Vault has
   automatically adjusted `list` operations to always end in a `/`, as list
   operations operates on prefixes, so all list operations by definition end
   with `/`. This was done server-side so affects all clients. However, this
   has also led to a lot of confusion for users writing policies that assume
   that the path that they use in the CLI is the path used internally. Starting
   in 0.11, ACL policies gain a new fallback rule for listing: they will use a
   matching path ending in `/` if available, but if not found, they will look
   for the same path without a trailing `/`. This allows putting `list`
   capabilities in the same path block as most other capabilities for that
   path, while not providing any extra access if `list` wasn't actually
   provided there.
 * Performance Standbys On By Default: If you flavor/license of Vault
   Enterprise supports Performance Standbys, they are on by default. You can
   disable this behavior per-node with the `disable_performance_standby`
   configuration flag.
 * AWS Secret Engine Roles: The AWS Secret Engine roles are now explicit about
   the type of AWS credential they are generating; this reduces reduce
   ambiguity that existed previously as well as enables new features for
   specific credential types. Writing role data and generating credentials
   remain backwards compatible; however, the data returned when reading a
   role's configuration has changed in backwards-incompatible ways. Anything
   that depended on reading role data from the AWS secret engine will break
   until it is updated to work with the new format.
 * Token Format (Enterprise): Tokens are now represented as a base62 value;
   tokens in namespaces will have the namespace identifier appended.

FEATURES:

 * **Namespaces (Enterprise)**: A set of features within Vault Enterprise
   that allows Vault environments to support *Secure Multi-tenancy* within a
   single Vault Enterprise infrastructure. Through namespaces, Vault
   administrators can support tenant isolation for teams and individuals as
   well as empower those individuals to self-manage their own tenant
   environment.
 * **Performance Standbys (Enterprise)**: Standby nodes can now service
   requests that do not modify storage. This provides near-horizontal scaling
   of a cluster in some workloads, and is the intra-cluster analogue of
   the existing Performance Replication feature, which replicates to distinct
   clusters in other datacenters, geos, etc.
 * **AliCloud OSS Storage**: AliCloud OSS can now be used for Vault storage.
 * **AliCloud Auth Plugin**: AliCloud's identity services can now be used to
   grant access to Vault. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-auth-alicloud) for
   more information.
 * **Azure Secrets Plugin**: There is now a plugin (pulled in to Vault) that
   allows generating credentials to allow access to Azure. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-secrets-azure) for
   more information.
 * **HA Support for MySQL Storage**: MySQL storage now supports HA.
 * **ACL Templating**: ACL policies can now be templated using identity Entity,
   Groups, and Metadata.
 * **UI Onboarding wizards**: The Vault UI can provide contextual help and
   guidance, linking out to relevant links or guides on vaultproject.io for
   various workflows in Vault.

IMPROVEMENTS:

 * agent: Add `exit_after_auth` to be able to use the Agent for a single
   authentication [[GH-5013](https://github.com/hashicorp/vault/pull/5013)]
 * auth/approle: Add ability to set token bound CIDRs on individual Secret IDs
   [[GH-5034](https://github.com/hashicorp/vault/pull/5034)]
 * cli: Add support for passing parameters to `vault read` operations [[GH-5093](https://github.com/hashicorp/vault/pull/5093)]
 * secrets/aws: Make credential types more explicit [[GH-4360](https://github.com/hashicorp/vault/pull/4360)]
 * secrets/nomad: Support for longer token names [[GH-5117](https://github.com/hashicorp/vault/pull/5117)]
 * secrets/pki: Allow disabling CRL generation [[GH-5134](https://github.com/hashicorp/vault/pull/5134)]
 * storage/azure: Add support for different Azure environments [[GH-4997](https://github.com/hashicorp/vault/pull/4997)]
 * storage/file: Sort keys in list responses [[GH-5141](https://github.com/hashicorp/vault/pull/5141)]
 * storage/mysql: Support special characters in database and table names.

BUG FIXES:

 * auth/jwt: Always validate `aud` claim even if `bound_audiences` isn't set
   (IOW, error in this case)
 * core: Prevent Go's HTTP library from interspersing logs in a different
   format and/or interleaved [[GH-5135](https://github.com/hashicorp/vault/pull/5135)]
 * identity: Properly populate `mount_path` and `mount_type` on group lookup
   [[GH-5074](https://github.com/hashicorp/vault/pull/5074)]
 * identity: Fix persisting alias metadata [[GH-5188](https://github.com/hashicorp/vault/pull/5188)]
 * identity: Fix carryover issue from previously fixed race condition that
   could cause Vault not to start up due to two entities referencing the same
   alias. These entities are now merged. [[GH-5000](https://github.com/hashicorp/vault/pull/5000)]
 * replication: Fix issue causing some pages not to flush to storage
 * secrets/database: Fix inability to update custom SQL statements on
   database roles. [[GH-5080](https://github.com/hashicorp/vault/pull/5080)]
 * secrets/pki: Disallow putting the CA's serial on its CRL. While technically
   legal, doing so inherently means the CRL can't be trusted anyways, so it's
   not useful and easy to footgun. [[GH-5134](https://github.com/hashicorp/vault/pull/5134)]
 * storage/gcp,spanner: Fix data races [[GH-5081](https://github.com/hashicorp/vault/pull/5081)]

## 0.10.4 (July 25th, 2018)

SECURITY:

 * Control Groups: The associated Identity entity with a request was not being
   properly persisted. As a result, the same authorizer could provide more than
   one authorization.

DEPRECATIONS/CHANGES:

 * Revocations of dynamic secrets leases are now queued/asynchronous rather
   than synchronous. This allows Vault to take responsibility for revocation
   even if the initial attempt fails. The previous synchronous behavior can be
   attained via the `-sync` CLI flag or `sync` API parameter. When in
   synchronous mode, if the operation results in failure it is up to the user
   to retry.
 * CLI Retries: The CLI will no longer retry commands on 5xx errors. This was a
   source of confusion to users as to why Vault would "hang" before returning a
   5xx error. The Go API client still defaults to two retries.
 * Identity Entity Alias metadata: You can no longer manually set metadata on
   entity aliases. All alias data (except the canonical entity ID it refers to)
   is intended to be managed by the plugin providing the alias information, so
   allowing it to be set manually didn't make sense.

FEATURES:

 * **JWT/OIDC Auth Method**: The new `jwt` auth method accepts JWTs and either
   validates signatures locally or uses OIDC Discovery to fetch the current set
   of keys for signature validation. Various claims can be specified for
   validation (in addition to the cryptographic signature) and a user and
   optional groups claim can be used to provide Identity information.
 * **FoundationDB Storage**: You can now use FoundationDB for storing Vault
   data.
 * **UI Control Group Workflow (enterprise)**: The UI will now detect control
   group responses and provides a workflow to view the status of the request
   and to authorize requests.
 * **Vault Agent (Beta)**: Vault Agent is a daemon that can automatically
   authenticate for you across a variety of authentication methods, provide
   tokens to clients, and keep the tokens renewed, reauthenticating as
   necessary.

IMPROVEMENTS:

 * auth/azure: Add support for virtual machine scale sets
 * auth/gcp: Support multiple bindings for region, zone, and instance group
 * cli: Add subcommands for interacting with the plugin catalog [[GH-4911](https://github.com/hashicorp/vault/pull/4911)]
 * cli: Add a `-description` flag to secrets and auth tune subcommands to allow
   updating an existing secret engine's or auth method's description. This
   change also allows the description to be unset by providing an empty string.
 * core: Add config flag to disable non-printable character check [[GH-4917](https://github.com/hashicorp/vault/pull/4917)]
 * core: A `max_request_size` parameter can now be set per-listener to adjust
   the maximum allowed size per request [[GH-4824](https://github.com/hashicorp/vault/pull/4824)]
 * core: Add control group request endpoint to default policy [[GH-4904](https://github.com/hashicorp/vault/pull/4904)]
 * identity: Identity metadata is now passed through to plugins [[GH-4967](https://github.com/hashicorp/vault/pull/4967)]
 * replication: Add additional saftey checks and logging when replication is
   in a bad state
 * secrets/kv: Add support for using `-field=data` to KVv2 when using `vault
   kv` [[GH-4895](https://github.com/hashicorp/vault/pull/4895)]
 * secrets/pki: Add the ability to tidy revoked but unexpired certificates
   [[GH-4916](https://github.com/hashicorp/vault/pull/4916)]
 * secrets/ssh: Allow Vault to work with single-argument SSH flags [[GH-4825](https://github.com/hashicorp/vault/pull/4825)]
 * secrets/ssh: SSH executable path can now be configured in the CLI [[GH-4937](https://github.com/hashicorp/vault/pull/4937)]
 * storage/swift: Add additional configuration options [[GH-4901](https://github.com/hashicorp/vault/pull/4901)]
 * ui: Choose which auth methods to show to unauthenticated users via
   `listing_visibility` in the auth method edit forms [[GH-4854](https://github.com/hashicorp/vault/pull/4854)]
 * ui: Authenticate users automatically by passing a wrapped token to the UI via
   the new `wrapped_token` query parameter [[GH-4854](https://github.com/hashicorp/vault/pull/4854)]

BUG FIXES:

 * api: Fix response body being cleared too early [[GH-4987](https://github.com/hashicorp/vault/pull/4987)]
 * auth/approle: Fix issue with tidy endpoint that would unnecessarily remove
   secret accessors [[GH-4981](https://github.com/hashicorp/vault/pull/4981)]
 * auth/aws: Fix updating `max_retries` [[GH-4980](https://github.com/hashicorp/vault/pull/4980)]
 * auth/kubernetes: Trim trailing whitespace when sending JWT
 * cli: Fix parsing of environment variables for integer flags [[GH-4925](https://github.com/hashicorp/vault/pull/4925)]
 * core: Fix returning 500 instead of 503 if a rekey is attempted when Vault is
   sealed [[GH-4874](https://github.com/hashicorp/vault/pull/4874)]
 * core: Fix issue releasing the leader lock in some circumstances [[GH-4915](https://github.com/hashicorp/vault/pull/4915)]
 * core: Fix a panic that could happen if the server was shut down while still
   starting up
 * core: Fix deadlock that would occur if a leadership loss occurs at the same
   time as a seal operation [[GH-4932](https://github.com/hashicorp/vault/pull/4932)]
 * core: Fix issue with auth mounts failing to renew tokens due to policies
   changing [[GH-4960](https://github.com/hashicorp/vault/pull/4960)]
 * auth/radius: Fix issue where some radius logins were being canceled too early
   [[GH-4941](https://github.com/hashicorp/vault/pull/4941)]
 * core: Fix accidental seal of vault of we lose leadership during startup
   [[GH-4924](https://github.com/hashicorp/vault/pull/4924)]
 * core: Fix standby not being able to forward requests larger than 4MB
   [[GH-4844](https://github.com/hashicorp/vault/pull/4844)]
 * core: Avoid panic while processing group memberships [[GH-4841](https://github.com/hashicorp/vault/pull/4841)]
 * identity: Fix a race condition creating aliases [[GH-4965](https://github.com/hashicorp/vault/pull/4965)]
 * plugins: Fix being unable to send very large payloads to or from plugins
   [[GH-4958](https://github.com/hashicorp/vault/pull/4958)]
 * physical/azure: Long list responses would sometimes be truncated [[GH-4983](https://github.com/hashicorp/vault/pull/4983)]
 * replication: Allow replication status requests to be processed while in
   merkle sync
 * replication: Ensure merkle reindex flushes all changes to storage immediately
 * replication: Fix a case where a network interruption could cause a secondary
   to be unable to reconnect to a primary
 * secrets/pki: Fix permitted DNS domains performing improper validation
   [[GH-4863](https://github.com/hashicorp/vault/pull/4863)]
 * secrets/database: Fix panic during DB creds revocation [[GH-4846](https://github.com/hashicorp/vault/pull/4846)]
 * ui: Fix usage of cubbyhole backend in the UI [[GH-4851](https://github.com/hashicorp/vault/pull/4851)]
 * ui: Fix toggle state when a secret is JSON-formatted [[GH-4913](https://github.com/hashicorp/vault/pull/4913)]
 * ui: Fix coercion of falsey values to empty string when editing secrets as
   JSON [[GH-4977](https://github.com/hashicorp/vault/pull/4977)]

## 0.10.3 (June 20th, 2018)

DEPRECATIONS/CHANGES:

 * In the audit log and in client responses, policies are now split into three
   parameters: policies that came only from tokens, policies that came only
   from Identity, and the combined set. Any previous location of policies via
   the API now contains the full, combined set.
 * When a token is tied to an Identity entity and the entity is deleted, the
   token will no longer be usable, regardless of the validity of the token
   itself.
 * When authentication succeeds but no policies were defined for that specific
   user, most auth methods would allow a token to be generated but a few would
   reject the authentication, namely `ldap`, `okta`, and `radius`. Since the
   `default` policy is added by Vault's core, this would incorrectly reject
   valid authentications before they would in fact be granted policies. This
   inconsistency has been addressed; valid authentications for these methods
   now succeed even if no policy was specifically defined in that method for
   that user.

FEATURES:

 * Root Rotation for Active Directory: You can now command Vault to rotate the
   configured root credentials used in the AD secrets engine, to ensure that
   only Vault knows the credentials it's using.
 * URI SANs in PKI: You can now configure URI Subject Alternate Names in the
   `pki` backend. Roles can limit which SANs are allowed via globbing.
 * `kv rollback` Command: You can now use `vault kv rollback` to roll a KVv2
   path back to a previous non-deleted/non-destroyed version. The previous
   version becomes the next/newest version for the path.
 * Token Bound CIDRs in AppRole: You can now add CIDRs to which a token
   generated from AppRole will be bound.

IMPROVEMENTS:

 * approle: Return 404 instead of 202 on invalid role names during POST
   operations [[GH-4778](https://github.com/hashicorp/vault/pull/4778)]
 * core: Add idle and initial header read/TLS handshake timeouts to connections
   to ensure server resources are cleaned up [[GH-4760](https://github.com/hashicorp/vault/pull/4760)]
 * core: Report policies in token, identity, and full sets [[GH-4747](https://github.com/hashicorp/vault/pull/4747)]
 * secrets/databases: Add `create`/`update` distinction for connection
   configurations [[GH-3544](https://github.com/hashicorp/vault/pull/3544)]
 * secrets/databases: Add `create`/`update` distinction for role configurations
   [[GH-3544](https://github.com/hashicorp/vault/pull/3544)]
 * secrets/databases: Add best-effort revocation logic for use when a role has
   been deleted [[GH-4782](https://github.com/hashicorp/vault/pull/4782)]
 * secrets/kv: Add `kv rollback` [[GH-4774](https://github.com/hashicorp/vault/pull/4774)]
 * secrets/pki: Add URI SANs support [[GH-4675](https://github.com/hashicorp/vault/pull/4675)]
 * secrets/ssh: Allow standard SSH command arguments to be used, without
   requiring username@hostname syntax [[GH-4710](https://github.com/hashicorp/vault/pull/4710)]
 * storage/consul: Add context support so that requests are cancelable
   [[GH-4739](https://github.com/hashicorp/vault/pull/4739)]
 * sys: Added `hidden` option to `listing_visibility` field on `sys/mounts`
   API [[GH-4827](https://github.com/hashicorp/vault/pull/4827)]
 * ui: Secret values are obfuscated by default and visibility is toggleable [[GH-4422](https://github.com/hashicorp/vault/pull/4422)]

BUG FIXES:

 * auth/approle: Fix panic due to metadata being nil [[GH-4719](https://github.com/hashicorp/vault/pull/4719)]
 * auth/aws: Fix delete path for tidy operations [[GH-4799](https://github.com/hashicorp/vault/pull/4799)]
 * core: Optimizations to remove some speed regressions due to the
   security-related changes in 0.10.2
 * storage/dynamodb: Fix errors seen when reading existing DynamoDB data [[GH-4721](https://github.com/hashicorp/vault/pull/4721)]
 * secrets/database: Fix default MySQL root rotation statement [[GH-4748](https://github.com/hashicorp/vault/pull/4748)]
 * secrets/gcp: Fix renewal for GCP account keys
 * secrets/kv: Fix writing to the root of a KVv2 mount from `vault kv` commands
   incorrectly operating on a root+mount path instead of being an error
   [[GH-4726](https://github.com/hashicorp/vault/pull/4726)]
 * seal/pkcs11: Add `CKK_SHA256_HMAC` to the search list when finding HMAC
   keys, fixing lookup on some Thales devices
 * replication: Fix issue enabling replication when a non-auth mount and auth
   mount have the same name
 * auth/kubernetes: Fix issue verifying ECDSA signed JWTs
 * ui: add missing edit mode for auth method configs [[GH-4770](https://github.com/hashicorp/vault/pull/4770)]

## 0.10.2 (June 6th, 2018)

SECURITY:

 * Tokens: A race condition was identified that could occur if a token's
   lease expired while Vault was not running. In this case, when Vault came
   back online, sometimes it would properly revoke the lease but other times it
   would not, leading to a Vault token that no longer had an expiration and had
   essentially unlimited lifetime. This race was per-token, not all-or-nothing
   for all tokens that may have expired during Vault's downtime. We have fixed
   the behavior and put extra checks in place to help prevent any similar
   future issues. In addition, the logic we have put in place ensures that such
   lease-less tokens can no longer be used (unless they are root tokens that
   never had an expiration to begin with).
 * Convergent Encryption: The version 2 algorithm used in `transit`'s
   convergent encryption feature is susceptible to offline
   plaintext-confirmation attacks. As a result, we are introducing a version 3
   algorithm that mitigates this. If you are currently using convergent
   encryption, we recommend upgrading, rotating your encryption key (the new
   key version will use the new algorithm), and rewrapping your data (the
   `rewrap` endpoint can be used to allow a relatively non-privileged user to
   perform the rewrapping while never divulging the plaintext).
 * AppRole case-sensitive role name secret-id leaking: When using a mixed-case
   role name via AppRole, deleting a secret-id via accessor or other operations
   could end up leaving the secret-id behind and valid but without an accessor.
   This has now been fixed, and we have put checks in place to prevent these
   secret-ids from being used.

DEPRECATIONS/CHANGES:

 * PKI duration return types: The PKI backend now returns durations (e.g. when
   reading a role) as an integer number of seconds instead of a Go-style
   string, in line with how the rest of Vault's API returns durations.

FEATURES:

 * Active Directory Secrets Engine: A new `ad` secrets engine has been created
   which allows Vault to rotate and provide credentials for configured AD
   accounts.
 * Rekey Verification: Rekey operations can now require verification. This
   turns on a two-phase process where the existing key shares authorize
   generating a new master key, and a threshold of the new, returned key shares
   must be provided to verify that they have been successfully received in
   order for the actual master key to be rotated.
 * CIDR restrictions for `cert`, `userpass`, and `kubernetes` auth methods:
   You can now limit authentication to specific CIDRs; these will also be
   encoded in resultant tokens to limit their use.
 * Vault UI Browser CLI: The UI now supports usage of read/write/list/delete
   commands in a CLI that can be accessed from the nav bar. Complex inputs such
   as JSON files are not currently supported. This surfaces features otherwise
   unsupported in Vault's UI.
 * Azure Key Vault Auto Unseal/Seal Wrap Support (Enterprise): Azure Key Vault
   can now be used a support seal for Auto Unseal and Seal Wrapping.

IMPROVEMENTS:

 * api: Close renewer's doneCh when the renewer is stopped, so that programs
   expecting a final value through doneCh behave correctly [[GH-4472](https://github.com/hashicorp/vault/pull/4472)]
 * auth/cert: Break out `allowed_names` into component parts and add
   `allowed_uri_sans` [[GH-4231](https://github.com/hashicorp/vault/pull/4231)]
 * auth/ldap: Obfuscate error messages pre-bind for greater security [[GH-4700](https://github.com/hashicorp/vault/pull/4700)]
 * cli: `vault login` now supports a `-no-print` flag to suppress printing
   token information but still allow storing into the token helper [[GH-4454](https://github.com/hashicorp/vault/pull/4454)]
 * core/pkcs11 (enterprise): Add support for CKM_AES_CBC_PAD, CKM_RSA_PKCS, and
   CKM_RSA_PKCS_OAEP mechanisms
 * core/pkcs11 (enterprise): HSM slots can now be selected by token label
   instead of just slot number
 * core/token: Optimize token revocation by removing unnecessary list call
   against the storage backend when calling revoke-orphan on tokens [[GH-4465](https://github.com/hashicorp/vault/pull/4465)]
 * core/token: Refactor token revocation logic to not block on the call when
   underlying leases are pending revocation by moving the expiration logic to
   the expiration manager [[GH-4512](https://github.com/hashicorp/vault/pull/4512)]
 * expiration: Allow revoke-prefix and revoke-force to work on single leases as
   well as prefixes [[GH-4450](https://github.com/hashicorp/vault/pull/4450)]
 * identity: Return parent group info when reading a group [[GH-4648](https://github.com/hashicorp/vault/pull/4648)]
 * identity: Provide more contextual key information when listing entities,
   groups, and aliases
 * identity: Passthrough EntityID to backends [[GH-4663](https://github.com/hashicorp/vault/pull/4663)]
 * identity: Adds ability to request entity information through system view
   [GH_4681]
 * secret/pki: Add custom extended key usages [[GH-4667](https://github.com/hashicorp/vault/pull/4667)]
 * secret/pki: Add custom PKIX serial numbers [[GH-4694](https://github.com/hashicorp/vault/pull/4694)]
 * secret/ssh: Use hostname instead of IP in OTP mode, similar to CA mode
   [[GH-4673](https://github.com/hashicorp/vault/pull/4673)]
 * storage/file: Attempt in some error conditions to do more cleanup [[GH-4684](https://github.com/hashicorp/vault/pull/4684)]
 * ui: wrapping lookup now distplays the path [[GH-4644](https://github.com/hashicorp/vault/pull/4644)]
 * ui: Identity interface now has more inline actions to make editing and adding
   aliases to an entity or group easier [[GH-4502](https://github.com/hashicorp/vault/pull/4502)]
 * ui: Identity interface now lists groups by name [[GH-4655](https://github.com/hashicorp/vault/pull/4655)]
 * ui: Permission denied errors still render the sidebar in the Access section
   [[GH-4658](https://github.com/hashicorp/vault/pull/4658)]
 * replication: Improve performance of index page flushes and WAL garbage
   collecting

BUG FIXES:

 * auth/approle: Make invalid role_id a 400 error instead of 500 [[GH-4470](https://github.com/hashicorp/vault/pull/4470)]
 * auth/cert: Fix Identity alias using serial number instead of common name
   [[GH-4475](https://github.com/hashicorp/vault/pull/4475)]
 * cli: Fix panic running `vault token capabilities` with multiple paths
   [[GH-4552](https://github.com/hashicorp/vault/pull/4552)]
 * core: When using the `use_always` option with PROXY protocol support, do not
   require `authorized_addrs` to be set [[GH-4065](https://github.com/hashicorp/vault/pull/4065)]
 * core: Fix panic when certain combinations of policy paths and allowed/denied
   parameters were used [[GH-4582](https://github.com/hashicorp/vault/pull/4582)]
 * secret/gcp: Make `bound_region` able to use short names
 * secret/kv: Fix response wrapping for KV v2 [[GH-4511](https://github.com/hashicorp/vault/pull/4511)]
 * secret/kv: Fix address flag not being honored correctly [[GH-4617](https://github.com/hashicorp/vault/pull/4617)]
 * secret/pki: Fix `safety_buffer` for tidy being allowed to be negative,
   clearing all certs [[GH-4641](https://github.com/hashicorp/vault/pull/4641)]
 * secret/pki: Fix `key_type` not being allowed to be set to `any` [[GH-4595](https://github.com/hashicorp/vault/pull/4595)]
 * secret/pki: Fix path length parameter being ignored when using
   `use_csr_values` and signing an intermediate CA cert [[GH-4459](https://github.com/hashicorp/vault/pull/4459)]
 * secret/ssh: Only append UserKnownHostsFile to args when configured with a
   value [[GH-4674](https://github.com/hashicorp/vault/pull/4674)]
 * storage/dynamodb: Fix listing when one child is left within a nested path
   [[GH-4570](https://github.com/hashicorp/vault/pull/4570)]
 * storage/gcs: Fix swallowing an error on connection close [[GH-4691](https://github.com/hashicorp/vault/pull/4691)]
 * ui: Fix HMAC algorithm in transit [[GH-4604](https://github.com/hashicorp/vault/pull/4604)]
 * ui: Fix unwrap of auth responses via the UI's unwrap tool [[GH-4611](https://github.com/hashicorp/vault/pull/4611)]
 * ui (enterprise): Fix parsing of version string that blocked some users from seeing
   enterprise-specific pages in the UI [[GH-4547](https://github.com/hashicorp/vault/pull/4547)]
 * ui: Fix incorrect capabilities path check when viewing policies [[GH-4566](https://github.com/hashicorp/vault/pull/4566)]
 * replication: Fix error while running plugins on a newly created replication
   secondary
 * replication: Fix issue with token store lookups after a secondary's mount table
   is invalidated.
 * replication: Improve startup time when a large merkle index is in use.
 * replication: Fix panic when storage becomes unreachable during unseal.

## 0.10.1/0.9.7 (April 25th, 2018)

The following two items are in both 0.9.7 and 0.10.1. They only affect
Enterprise, and as such 0.9.7 is an Enterprise-only release:

SECURITY:

 * EGPs: A regression affecting 0.9.6 and 0.10.0 causes EGPs to not be applied
   correctly if an EGP is updated in a running Vault after initial write or
   after it is loaded on unseal. This has been fixed.

BUG FIXES:

 * Fixed an upgrade issue affecting performance secondaries when migrating from
   a version that did not include Identity to one that did.

All other content in this release is for 0.10.1 only.

DEPRECATIONS/CHANGES:

 * `vault kv` and Vault versions: In 0.10.1 some issues with `vault kv` against
   v1 K/V engine mounts are fixed. However, using 0.10.1 for both the server
   and CLI versions is required.
 * Mount information visibility: Users that have access to any path within a
   mount can now see information about that mount, such as its type and
   options, via some API calls.
 * Identity and Local Mounts: Local mounts would allow creating Identity
   entities but these would not be able to be used successfully (even locally)
   in replicated scenarios. We have now disallowed entities and groups from
   being created for local mounts in the first place.

FEATURES:

 * X-Forwarded-For support: `X-Forwarded-For` headers can now be used to set the
   client IP seen by Vault. See the [TCP listener configuration
   page](https://www.vaultproject.io/docs/configuration/listener/tcp.html) for
   details.
 * CIDR IP Binding for Tokens: Tokens now support being bound to specific
   CIDR(s) for usage. Currently this is implemented in Token Roles; usage can be
   expanded to other authentication backends over time.
 * `vault kv patch` command: A new `kv patch` helper command that allows
   modifying only some values in existing data at a K/V path, but uses
   check-and-set to ensure that this modification happens safely.
 * AppRole Local Secret IDs: Roles can now be configured to generate secret IDs
   local to the cluster. This enables performance secondaries to generate and
   consume secret IDs without contacting the primary.
 * AES-GCM Support for PKCS#11 [BETA] (Enterprise): For supporting HSMs,
   AES-GCM can now be used in lieu of AES-CBC/HMAC-SHA256. This has currently
   only been fully tested on AWS CloudHSM.
 * Auto Unseal/Seal Wrap Key Rotation Support (Enterprise): Auto Unseal
   mechanisms, including PKCS#11 HSMs, now support rotation of encryption keys,
   and migration between key and encryption types, such as from AES-CBC to
   AES-GCM, can be performed at the same time (where supported).

IMPROVEMENTS:

 * auth/approle: Support for cluster local secret IDs. This enables secondaries
   to generate secret IDs without contacting the primary [[GH-4427](https://github.com/hashicorp/vault/pull/4427)]
 * auth/token: Add to the token lookup response, the policies inherited due to
   identity associations [[GH-4366](https://github.com/hashicorp/vault/pull/4366)]
 * auth/token: Add CIDR binding to token roles [[GH-815](https://github.com/hashicorp/vault/pull/815)]
 * cli: Add `vault kv patch` [[GH-4432](https://github.com/hashicorp/vault/pull/4432)]
 * core: Add X-Forwarded-For support [[GH-4380](https://github.com/hashicorp/vault/pull/4380)]
 * core: Add token CIDR-binding support [[GH-815](https://github.com/hashicorp/vault/pull/815)]
 * identity: Add the ability to disable an entity. Disabling an entity does not
   revoke associated tokens, but while the entity is disabled they cannot be
   used. [[GH-4353](https://github.com/hashicorp/vault/pull/4353)]
 * physical/consul: Allow tuning of session TTL and lock wait time [[GH-4352](https://github.com/hashicorp/vault/pull/4352)]
 * replication: Dynamically adjust WAL cleanup over a period of time based on
   the rate of writes committed
 * secret/ssh: Update dynamic key install script to use shell locking to avoid
   concurrent modifications [[GH-4358](https://github.com/hashicorp/vault/pull/4358)]
 * ui: Access to `sys/mounts` is no longer needed to use the UI - the list of
   engines will show you the ones you implicitly have access to (because you have
   access to to secrets in those engines) [[GH-4439](https://github.com/hashicorp/vault/pull/4439)]

BUG FIXES:

 * cli: Fix `vault kv` backwards compatibility with KV v1 engine mounts
   [[GH-4430](https://github.com/hashicorp/vault/pull/4430)]
 * identity: Persist entity memberships in external identity groups across
   mounts [[GH-4365](https://github.com/hashicorp/vault/pull/4365)]
 * identity: Fix error preventing authentication using local mounts on
   performance secondary replication clusters [[GH-4407](https://github.com/hashicorp/vault/pull/4407)]
 * replication: Fix issue causing secondaries to not connect properly to a
   pre-0.10 primary until the primary was upgraded
 * secret/gcp: Fix panic on rollback when a roleset wasn't created properly
   [[GH-4344](https://github.com/hashicorp/vault/pull/4344)]
 * secret/gcp: Fix panic on renewal
 * ui: Fix IE11 form submissions in a few parts of the application [[GH-4378](https://github.com/hashicorp/vault/pull/4378)]
 * ui: Fix IE file saving on policy pages and init screens [[GH-4376](https://github.com/hashicorp/vault/pull/4376)]
 * ui: Fixed an issue where the AWS secret backend would show the wrong menu
   [[GH-4371](https://github.com/hashicorp/vault/pull/4371)]
 * ui: Fixed an issue where policies with commas would not render in the
   interface properly [[GH-4398](https://github.com/hashicorp/vault/pull/4398)]
 * ui: Corrected the saving of mount tune ttls for auth methods [[GH-4431](https://github.com/hashicorp/vault/pull/4431)]
 * ui: Credentials generation no longer checks capabilities before making
   api calls. This should fix needing "update" capabilites to read IAM
   credentials in the AWS secrets engine [[GH-4446](https://github.com/hashicorp/vault/pull/4446)]

## 0.10.0 (April 10th, 2018)

SECURITY:

 * Log sanitization for Combined Database Secret Engine: In certain failure
   scenarios with incorrectly formatted connection urls, the raw connection
   errors were being returned to the user with the configured database
   credentials. Errors are now sanitized before being returned to the user.

DEPRECATIONS/CHANGES:

 * Database plugin compatibility: The database plugin interface was enhanced to
   support some additional functionality related to root credential rotation
   and supporting templated URL strings. The changes were made in a
   backwards-compatible way and all builtin plugins were updated with the new
   features. Custom plugins not built into Vault will need to be upgraded to
   support templated URL strings and root rotation. Additionally, the
   Initialize method was deprecated in favor of a new Init method that supports
   configuration modifications that occur in the plugin back to the primary
   data store.
 * Removal of returned secret information: For a long time Vault has returned
   configuration given to various secret engines and auth methods with secret
   values (such as secret API keys or passwords) still intact, and with a
   warning to the user on write that anyone with read access could see the
   secret. This was mostly done to make it easy for tools like Terraform to
   judge whether state had drifted. However, it also feels quite un-Vault-y to
   do this and we've never felt very comfortable doing so. In 0.10 we have gone
   through and removed this behavior from the various backends; fields which
   contained secret values are simply no longer returned on read. We are
   working with the Terraform team to make changes to their provider to
   accommodate this as best as possible, and users of other tools may have to
   make adjustments, but in the end we felt that the ends did not justify the
   means and we needed to prioritize security over operational convenience.
 * LDAP auth method case sensitivity: We now treat usernames and groups
   configured locally for policy assignment in a case insensitive fashion by
   default. Existing configurations will continue to work as they do now;
   however, the next time a configuration is written `case_sensitive_names`
   will need to be explicitly set to `true`.
 * TTL handling within core: All lease TTL handling has been centralized within
   the core of Vault to ensure consistency across all backends. Since this was
   previously delegated to individual backends, there may be some slight
   differences in TTLs generated from some backends.
 * Removal of default `secret/` mount: In 0.12 we will stop mounting `secret/`
   by default at initialization time (it will still be available in `dev`
   mode).

FEATURES:

 * OSS UI: The Vault UI is now fully open-source. Similarly to the CLI, some
   features are only available with a supporting version of Vault, but the code
   base is entirely open.
 * Versioned K/V: The `kv` backend has been completely revamped, featuring
   flexible versioning of values, check-and-set protections, and more. A new
   `vault kv` subcommand allows friendly interactions with it. Existing mounts
   of the `kv` backend can be upgraded to the new versioned mode (downgrades
   are not currently supported). The old "passthrough" mode is still the
   default for new mounts; versioning can be turned on by setting the
   `-version=2` flag for the `vault secrets enable` command.
 * Database Root Credential Rotation: Database configurations can now rotate
   their own configured admin/root credentials, allowing configured credentials
   for a database connection to be rotated immediately after sending them into
   Vault, invalidating the old credentials and ensuring only Vault knows the
   actual valid values.
 * Azure Authentication Plugin: There is now a plugin (pulled in to Vault) that
   allows authenticating Azure machines to Vault using Azure's Managed Service
   Identity credentials. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-auth-azure) for more
   information.
 * GCP Secrets Plugin: There is now a plugin (pulled in to Vault) that allows
   generating secrets to allow access to GCP. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-secrets-gcp) for more
   information.
 * Selective Audit HMACing of Request and Response Data Keys: HMACing in audit
   logs can be turned off for specific keys in the request input map and
   response `data` map on a per-mount basis.
 * Passthrough Request Headers: Request headers can now be selectively passed
   through to backends on a per-mount basis. This is useful in various cases
   when plugins are interacting with external services.
 * HA for Google Cloud Storage: The GCS storage type now supports HA.
 * UI support for identity: Add and edit entities, groups, and their associated
   aliases.
 * UI auth method support: Enable, disable, and configure all of the built-in
   authentication methods.
 * UI (Enterprise): View and edit Sentinel policies.

IMPROVEMENTS:

 * core: Centralize TTL generation for leases in core [[GH-4230](https://github.com/hashicorp/vault/pull/4230)]
 * identity: API to update group-alias by ID [[GH-4237](https://github.com/hashicorp/vault/pull/4237)]
 * secret/cassandra: Update Cassandra storage delete function to not use batch
   operations [[GH-4054](https://github.com/hashicorp/vault/pull/4054)]
 * storage/mysql: Allow setting max idle connections and connection lifetime
   [[GH-4211](https://github.com/hashicorp/vault/pull/4211)]
 * storage/gcs: Add HA support [[GH-4226](https://github.com/hashicorp/vault/pull/4226)]
 * ui: Add Nomad to the list of available secret engines
 * ui: Adds ability to set static headers to be returned by the UI

BUG FIXES:

 * api: Fix retries not working [[GH-4322](https://github.com/hashicorp/vault/pull/4322)]
 * auth/gcp: Invalidate clients on config change
 * auth/token: Revoke-orphan and tidy operations now correctly cleans up the
   parent prefix entry in the underlying storage backend. These operations also
   mark corresponding child tokens as orphans by removing the parent/secondary
   index from the entries. [[GH-4193](https://github.com/hashicorp/vault/pull/4193)]
 * command: Re-add `-mfa` flag and migrate to OSS binary [[GH-4223](https://github.com/hashicorp/vault/pull/4223)]
 * core: Fix issue occurring from mounting two auth backends with the same path
   with one mount having `auth/` in front [[GH-4206](https://github.com/hashicorp/vault/pull/4206)]
 * mfa: Invalidation of MFA configurations (Enterprise)
 * replication: Fix a panic on some non-64-bit platforms
 * replication: Fix invalidation of policies on performance secondaries
 * secret/pki: When tidying if a value is unexpectedly nil, delete it and move
   on [[GH-4214](https://github.com/hashicorp/vault/pull/4214)]
 * storage/s3: Fix panic if S3 returns no Content-Length header [[GH-4222](https://github.com/hashicorp/vault/pull/4222)]
 * ui: Fixed an issue where the UI was checking incorrect paths when operating
   on transit keys. Capabilities are now checked when attempting to encrypt /
   decrypt, etc.
 * ui: Fixed IE 11 layout issues and JS errors that would stop the application
   from running.
 * ui: Fixed the link that gets rendered when a user doesn't have permissions
   to view the root of a secret engine. The link now sends them back to the list
   of secret engines.
 * replication: Fix issue with DR secondaries when using mount specified local
   paths.
 * cli: Fix an issue where generating a dr operation token would not output the
   token [[GH-4328](https://github.com/hashicorp/vault/pull/4328)]

## 0.9.6 (March 20th, 2018)

DEPRECATIONS/CHANGES:

 * The AWS authentication backend now allows binds for inputs as either a
   comma-delimited string or a string array. However, to keep consistency with
   input and output, when reading a role the binds will now be returned as
   string arrays rather than strings.
 * In order to prefix-match IAM role and instance profile ARNs in AWS auth
   backend, you now must explicitly opt-in by adding a `*` to the end of the
   ARN. Existing configurations will be upgraded automatically, but when
   writing a new role configuration the updated behavior will be used.

FEATURES:

 * Replication Activation Enhancements: When activating a replication
   secondary, a public key can now be fetched first from the target cluster.
   This public key can be provided to the primary when requesting the
   activation token. If provided, the public key will be used to perform a
   Diffie-Hellman key exchange resulting in a shared key that encrypts the
   contents of the activation token. The purpose is to protect against
   accidental disclosure of the contents of the token if unwrapped by the wrong
   party, given that the contents of the token are highly sensitive. If
   accidentally unwrapped, the contents of the token are not usable by the
   unwrapping party. It is important to note that just as a malicious operator
   could unwrap the contents of the token, a malicious operator can pretend to
   be a secondary and complete the Diffie-Hellman exchange on their own; this
   feature provides defense in depth but still requires due diligence around
   replication activation, including multiple eyes on the commands/tokens and
   proper auditing.

IMPROVEMENTS:

 * api: Update renewer grace period logic. It no longer is static, but rather
   dynamically calculates one based on the current lease duration after each
   renew. [[GH-4090](https://github.com/hashicorp/vault/pull/4090)]
 * auth/approle: Allow array input for bound_cidr_list [4078]
 * auth/aws: Allow using lists in role bind parameters [[GH-3907](https://github.com/hashicorp/vault/pull/3907)]
 * auth/aws: Allow binding by EC2 instance IDs [[GH-3816](https://github.com/hashicorp/vault/pull/3816)]
 * auth/aws: Allow non-prefix-matched IAM role and instance profile ARNs
   [[GH-4071](https://github.com/hashicorp/vault/pull/4071)]
 * auth/ldap: Set a very large size limit on queries [[GH-4169](https://github.com/hashicorp/vault/pull/4169)]
 * core: Log info notifications of revoked leases for all leases/reasons, not
   just expirations [[GH-4164](https://github.com/hashicorp/vault/pull/4164)]
 * physical/couchdb: Removed limit on the listing of items [[GH-4149](https://github.com/hashicorp/vault/pull/4149)]
 * secret/pki: Support certificate policies [[GH-4125](https://github.com/hashicorp/vault/pull/4125)]
 * secret/pki: Add ability to have CA:true encoded into intermediate CSRs, to
   improve compatibility with some ADFS scenarios [[GH-3883](https://github.com/hashicorp/vault/pull/3883)]
 * secret/transit: Allow selecting signature algorithm as well as hash
   algorithm when signing/verifying [[GH-4018](https://github.com/hashicorp/vault/pull/4018)]
 * server: Make sure `tls_disable_client_cert` is actually a true value rather
   than just set [[GH-4049](https://github.com/hashicorp/vault/pull/4049)]
 * storage/dynamodb: Allow specifying max retries for dynamo client [[GH-4115](https://github.com/hashicorp/vault/pull/4115)]
 * storage/gcs: Allow specifying chunk size for transfers, which can reduce
   memory utilization [[GH-4060](https://github.com/hashicorp/vault/pull/4060)]
 * sys/capabilities: Add the ability to use multiple paths for capability
   checking [[GH-3663](https://github.com/hashicorp/vault/pull/3663)]

BUG FIXES:

 * auth/aws: Fix honoring `max_ttl` when a corresponding role `ttl` is not also
   set [[GH-4107](https://github.com/hashicorp/vault/pull/4107)]
 * auth/okta: Fix honoring configured `max_ttl` value [[GH-4110](https://github.com/hashicorp/vault/pull/4110)]
 * auth/token: If a periodic token being issued has a period greater than the
   max_lease_ttl configured on the token store mount, truncate it. This matches
   renewal behavior; before it was inconsistent between issuance and renewal.
   [[GH-4112](https://github.com/hashicorp/vault/pull/4112)]
 * cli: Improve error messages around `vault auth help` when there is no CLI
   helper for a particular method [[GH-4056](https://github.com/hashicorp/vault/pull/4056)]
 * cli: Fix autocomplete installation when using Fish as the shell [[GH-4094](https://github.com/hashicorp/vault/pull/4094)]
 * secret/database: Properly honor mount-tuned max TTL [[GH-4051](https://github.com/hashicorp/vault/pull/4051)]
 * secret/ssh: Return `key_bits` value when reading a role [[GH-4098](https://github.com/hashicorp/vault/pull/4098)]
 * sys: When writing policies on a performance replication secondary, properly
   forward requests to the primary [[GH-4129](https://github.com/hashicorp/vault/pull/4129)]

## 0.9.5 (February 26th, 2018)

IMPROVEMENTS:

 * auth: Allow sending default_lease_ttl and max_lease_ttl values when enabling
   auth methods. [[GH-4019](https://github.com/hashicorp/vault/pull/4019)]
 * secret/database: Add list functionality to `database/config` endpoint
   [[GH-4026](https://github.com/hashicorp/vault/pull/4026)]
 * physical/consul: Allow setting a specific service address [[GH-3971](https://github.com/hashicorp/vault/pull/3971)]
 * replication: When bootstrapping a new secondary, if the initial cluster
   connection fails, Vault will attempt to roll back state so that
   bootstrapping can be tried again, rather than having to recreate the
   downstream cluster. This will still require fetching a new secondary
   activation token.

BUG FIXES:

 * auth/aws: Update libraries to fix regression verifying PKCS#7 identity
   documents [[GH-4014](https://github.com/hashicorp/vault/pull/4014)]
 * listener: Revert to Go 1.9 for now to allow certificates with non-DNS names
   in their DNS SANs to be used for Vault's TLS connections [[GH-4028](https://github.com/hashicorp/vault/pull/4028)]
 * replication: Fix issue with a performance secondary/DR primary node losing
   its DR primary status when performing an update-primary operation
 * replication: Fix issue where performance secondaries could be unable to
   automatically connect to a performance primary after that performance
   primary has been promoted to a DR primary from a DR secondary
 * ui: Fix behavior when a value contains a `.`

## 0.9.4 (February 20th, 2018)

SECURITY:

 * Role Tags used with the EC2 style of AWS auth were being improperly parsed;
   as a result they were not being used to properly restrict values.
   Implementations following our suggestion of using these as defense-in-depth
   rather than the only source of restriction should not have significant
   impact.

FEATURES:

 * **ChaCha20-Poly1305 support in `transit`**: You can now encrypt and decrypt
   with ChaCha20-Poly1305 in `transit`. Key derivation and convergent
   encryption is also supported.
 * **Okta Push support in Okta Auth Backend**: If a user account has MFA
   required within Okta, an Okta Push MFA flow can be used to successfully
   finish authentication.
 * **PKI Improvements**: Custom OID subject alternate names can now be set,
   subject to allow restrictions that support globbing. Additionally, Country,
   Locality, Province, Street Address, and Postal Code can now be set in
   certificate subjects.
 * **Manta Storage**: Joyent Triton Manta can now be used for Vault storage
 * **Google Cloud Spanner Storage**: Google Cloud Spanner can now be used for
   Vault storage

IMPROVEMENTS:

 * auth/centrify: Add CLI helper
 * audit: Always log failure metrics, even if zero, to ensure the values appear
   on dashboards [[GH-3937](https://github.com/hashicorp/vault/pull/3937)]
 * cli: Disable color when output is not a TTY [[GH-3897](https://github.com/hashicorp/vault/pull/3897)]
 * cli: Add `-format` flag to all subcommands [[GH-3897](https://github.com/hashicorp/vault/pull/3897)]
 * cli: Do not display deprecation warnings when the format is not table
   [[GH-3897](https://github.com/hashicorp/vault/pull/3897)]
 * core: If over a predefined lease count (256k), log a warning not more than
   once a minute. Too many leases can be problematic for many of the storage
   backends and often this number of leases is indicative of a need for
   workflow improvements. [[GH-3957](https://github.com/hashicorp/vault/pull/3957)]
 * secret/nomad: Have generated ACL tokens cap out at 64 characters [[GH-4009](https://github.com/hashicorp/vault/pull/4009)]
 * secret/pki: Country, Locality, Province, Street Address, and Postal Code can
   now be set on certificates [[GH-3992](https://github.com/hashicorp/vault/pull/3992)]
 * secret/pki: UTF-8 Other Names can now be set in Subject Alternate Names in
   issued certs; allowed values can be set per role and support globbing
   [[GH-3889](https://github.com/hashicorp/vault/pull/3889)]
 * secret/pki: Add a flag to make the common name optional on certs [[GH-3940](https://github.com/hashicorp/vault/pull/3940)]
 * secret/pki: Ensure only DNS-compatible names go into DNS SANs; additionally,
   properly handle IDNA transformations for these DNS names [[GH-3953](https://github.com/hashicorp/vault/pull/3953)]
 * secret/ssh: Add `valid-principles` flag to CLI for CA mode [[GH-3922](https://github.com/hashicorp/vault/pull/3922)]
 * storage/manta: Add Manta storage [[GH-3270](https://github.com/hashicorp/vault/pull/3270)]
 * ui (Enterprise): Support for ChaCha20-Poly1305 keys in the transit engine.

BUG FIXES:
 * api/renewer: Honor increment value in renew auth calls [[GH-3904](https://github.com/hashicorp/vault/pull/3904)]
 * auth/approle: Fix inability to use limited-use-count secret IDs on
   replication performance secondaries
 * auth/approle: Cleanup of secret ID accessors during tidy and removal of
   dangling accessor entries [[GH-3924](https://github.com/hashicorp/vault/pull/3924)]
 * auth/aws-ec2: Avoid masking of role tag response [[GH-3941](https://github.com/hashicorp/vault/pull/3941)]
 * auth/cert: Verify DNS SANs in the authenticating certificate [[GH-3982](https://github.com/hashicorp/vault/pull/3982)]
 * auth/okta: Return configured durations as seconds, not nanoseconds [[GH-3871](https://github.com/hashicorp/vault/pull/3871)]
 * auth/okta: Get all okta groups for a user vs. default 200 limit [[GH-4034](https://github.com/hashicorp/vault/pull/4034)]
 * auth/token: Token creation via the CLI no longer forces periodic token
   creation. Passing an explicit zero value for the period no longer create
   periodic tokens. [[GH-3880](https://github.com/hashicorp/vault/pull/3880)]
 * command: Fix interpreted formatting directives when printing raw fields
   [[GH-4005](https://github.com/hashicorp/vault/pull/4005)]
 * command: Correctly format output when using -field and -format flags at the
   same time [[GH-3987](https://github.com/hashicorp/vault/pull/3987)]
 * command/rekey: Re-add lost `stored-shares` parameter [[GH-3974](https://github.com/hashicorp/vault/pull/3974)]
 * command/ssh: Create and reuse the api client [[GH-3909](https://github.com/hashicorp/vault/pull/3909)]
 * command/status: Fix panic when status returns 500 from leadership lookup
   [[GH-3998](https://github.com/hashicorp/vault/pull/3998)]
 * identity: Fix race when creating entities [[GH-3932](https://github.com/hashicorp/vault/pull/3932)]
 * plugin/gRPC: Fixed an issue with list requests and raw responses coming from
   plugins using gRPC transport [[GH-3881](https://github.com/hashicorp/vault/pull/3881)]
 * plugin/gRPC: Fix panic when special paths are not set [[GH-3946](https://github.com/hashicorp/vault/pull/3946)]
 * secret/pki: Verify a name is a valid hostname before adding to DNS SANs
   [[GH-3918](https://github.com/hashicorp/vault/pull/3918)]
 * secret/transit: Fix auditing when reading a key after it has been backed up
   or restored [[GH-3919](https://github.com/hashicorp/vault/pull/3919)]
 * secret/transit: Fix storage/memory consistency when persistence fails
   [[GH-3959](https://github.com/hashicorp/vault/pull/3959)]
 * storage/consul: Validate that service names are RFC 1123 compliant [[GH-3960](https://github.com/hashicorp/vault/pull/3960)]
 * storage/etcd3: Fix memory ballooning with standby instances [[GH-3798](https://github.com/hashicorp/vault/pull/3798)]
 * storage/etcd3: Fix large lists (like token loading at startup) not being
   handled [[GH-3772](https://github.com/hashicorp/vault/pull/3772)]
 * storage/postgresql: Fix compatibility with versions using custom string
   version tags [[GH-3949](https://github.com/hashicorp/vault/pull/3949)]
 * storage/zookeeper: Update vendoring to fix freezing issues [[GH-3896](https://github.com/hashicorp/vault/pull/3896)]
 * ui (Enterprise): Decoding the replication token should no longer error and
   prevent enabling of a secondary replication cluster via the ui.
 * plugin/gRPC: Add connection info to the request object [[GH-3997](https://github.com/hashicorp/vault/pull/3997)]

## 0.9.3 (January 28th, 2018)

A regression from a feature merge disabled the Nomad secrets backend in 0.9.2.
This release re-enables the Nomad secrets backend; it is otherwise identical to
0.9.2.

## 0.9.2 (January 26th, 2018)

SECURITY:

 * Okta Auth Backend: While the Okta auth backend was successfully verifying
   usernames and passwords, it was not checking the returned state of the
   account, so accounts that had been marked locked out could still be used to
   log in. Only accounts in SUCCESS or PASSWORD_WARN states are now allowed.
 * Periodic Tokens: A regression in 0.9.1 meant that periodic tokens created by
   the AppRole, AWS, and Cert auth backends would expire when the max TTL for
   the backend/mount/system was hit instead of their stated behavior of living
   as long as they are renewed. This is now fixed; existing tokens do not have
   to be reissued as this was purely a regression in the renewal logic.
 * Seal Wrapping: During certain replication states values written marked for
   seal wrapping may not be wrapped on the secondaries. This has been fixed,
   and existing values will be wrapped on next read or write. This does not
   affect the barrier keys.

DEPRECATIONS/CHANGES:

 * `sys/health` DR Secondary Reporting: The `replication_dr_secondary` bool
   returned by `sys/health` could be misleading since it would be `false` both
   when a cluster was not a DR secondary but also when the node is a standby in
   the cluster and has not yet fully received state from the active node. This
   could cause health checks on LBs to decide that the node was acceptable for
   traffic even though DR secondaries cannot handle normal Vault traffic. (In
   other words, the bool could only convey "yes" or "no" but not "not sure
   yet".) This has been replaced by `replication_dr_mode` and
   `replication_perf_mode` which are string values that convey the current
   state of the node; a value of `disabled` indicates that replication is
   disabled or the state is still being discovered. As a result, an LB check
   can positively verify that the node is both not `disabled` and is not a DR
   secondary, and avoid sending traffic to it if either is true.
 * PKI Secret Backend Roles parameter types: For `ou` and `organization`
   in role definitions in the PKI secret backend, input can now be a
   comma-separated string or an array of strings. Reading a role will
   now return arrays for these parameters.
 * Plugin API Changes: The plugin API has been updated to utilize golang's
   context.Context package. Many function signatures now accept a context
   object as the first parameter. Existing plugins will need to pull in the
   latest Vault code and update their function signatures to begin using
   context and the new gRPC transport.

FEATURES:

 * **gRPC Backend Plugins**: Backend plugins now use gRPC for transport,
   allowing them to be written in other languages.
 * **Brand New CLI**: Vault has a brand new CLI interface that is significantly
   streamlined, supports autocomplete, and is almost entirely backwards
   compatible.
 * **UI: PKI Secret Backend (Enterprise)**: Configure PKI secret backends,
   create and browse roles and certificates, and issue and sign certificates via
   the listed roles.

IMPROVEMENTS:

 * auth/aws: Handle IAM headers produced by clients that formulate numbers as
   ints rather than strings [[GH-3763](https://github.com/hashicorp/vault/pull/3763)]
 * auth/okta: Support JSON lists when specifying groups and policies [[GH-3801](https://github.com/hashicorp/vault/pull/3801)]
 * autoseal/hsm: Attempt reconnecting to the HSM on certain kinds of issues,
   including HA scenarios for some Gemalto HSMs.
   (Enterprise)
 * cli: Output password prompts to stderr to make it easier to pipe an output
   token to another command [[GH-3782](https://github.com/hashicorp/vault/pull/3782)]
 * core: Report replication status in `sys/health` [[GH-3810](https://github.com/hashicorp/vault/pull/3810)]
 * physical/s3: Allow using paths with S3 for non-AWS deployments [[GH-3730](https://github.com/hashicorp/vault/pull/3730)]
 * physical/s3: Add ability to disable SSL for non-AWS deployments [[GH-3730](https://github.com/hashicorp/vault/pull/3730)]
 * plugins: Args for plugins can now be specified separately from the command,
   allowing the same output format and input format for plugin information
   [[GH-3778](https://github.com/hashicorp/vault/pull/3778)]
 * secret/pki: `ou` and `organization` can now be specified as a
   comma-separated string or an array of strings [[GH-3804](https://github.com/hashicorp/vault/pull/3804)]
 * plugins: Plugins will fall back to using netrpc as the communication protocol
   on older versions of Vault [[GH-3833](https://github.com/hashicorp/vault/pull/3833)]

BUG FIXES:

 * auth/(approle,aws,cert): Fix behavior where periodic tokens generated by
   these backends could not have their TTL renewed beyond the system/mount max
   TTL value [[GH-3803](https://github.com/hashicorp/vault/pull/3803)]
 * auth/aws: Fix error returned if `bound_iam_principal_arn` was given to an
   existing role update [[GH-3843](https://github.com/hashicorp/vault/pull/3843)]
 * core/sealwrap: Speed improvements and bug fixes (Enterprise)
 * identity: Delete group alias when an external group is deleted [[GH-3773](https://github.com/hashicorp/vault/pull/3773)]
 * legacymfa/duo: Fix intermittent panic when Duo could not be reached
   [[GH-2030](https://github.com/hashicorp/vault/pull/2030)]
 * secret/database: Fix a location where a lock could potentially not be
   released, leading to deadlock [[GH-3774](https://github.com/hashicorp/vault/pull/3774)]
 * secret/(all databases) Fix behavior where if a max TTL was specified but no
   default TTL was specified the system/mount default TTL would be used but not
   be capped by the local max TTL [[GH-3814](https://github.com/hashicorp/vault/pull/3814)]
 * secret/database: Fix an issue where plugins were not closed properly if they
   failed to initialize [[GH-3768](https://github.com/hashicorp/vault/pull/3768)]
 * ui: mounting a secret backend will now properly set `max_lease_ttl` and
   `default_lease_ttl` when specified - previously both fields set
   `default_lease_ttl`.

## 0.9.1 (December 21st, 2017)

DEPRECATIONS/CHANGES:

 * AppRole Case Sensitivity: In prior versions of Vault, `list` operations
   against AppRole roles would require preserving case in the role name, even
   though most other operations within AppRole are case-insensitive with
   respect to the role name. This has been fixed; existing roles will behave as
   they have in the past, but new roles will act case-insensitively in these
   cases.
 * Token Auth Backend Roles parameter types: For `allowed_policies` and
   `disallowed_policies` in role definitions in the token auth backend, input
   can now be a comma-separated string or an array of strings. Reading a role
   will now return arrays for these parameters.
 * Transit key exporting: You can now mark a key in the `transit` backend as
   `exportable` at any time, rather than just at creation time; however, once
   this value is set, it still cannot be unset.
 * PKI Secret Backend Roles parameter types: For `allowed_domains` and
   `key_usage` in role definitions in the PKI secret backend, input
   can now be a comma-separated string or an array of strings. Reading a role
   will now return arrays for these parameters.
 * SSH Dynamic Keys Method Defaults to 2048-bit Keys: When using the dynamic
   key method in the SSH backend, the default is now to use 2048-bit keys if no
   specific key bit size is specified.
 * Consul Secret Backend lease handling: The `consul` secret backend can now
   accept both strings and integer numbers of seconds for its lease value. The
   value returned on a role read will be an integer number of seconds instead
   of a human-friendly string.
 * Unprintable characters not allowed in API paths: Unprintable characters are
   no longer allowed in names in the API (paths and path parameters), with an
   extra restriction on whitespace characters. Allowed characters are those
   that are considered printable by Unicode plus spaces.

FEATURES:

 * **Transit Backup/Restore**: The `transit` backend now supports a backup
   operation that can export a given key, including all key versions and
   configuration, as well as a restore operation allowing import into another
   Vault.
 * **gRPC Database Plugins**: Database plugins now use gRPC for transport,
   allowing them to be written in other languages.
 * **Nomad Secret Backend**: Nomad ACL tokens can now be generated and revoked
   using Vault.
 * **TLS Cert Auth Backend Improvements**: The `cert` auth backend can now
   match against custom certificate extensions via exact or glob matching, and
   additionally supports max_ttl and periodic token toggles.

IMPROVEMENTS:

 * auth/cert: Support custom certificate constraints [[GH-3634](https://github.com/hashicorp/vault/pull/3634)]
 * auth/cert: Support setting `max_ttl` and `period` [[GH-3642](https://github.com/hashicorp/vault/pull/3642)]
 * audit/file: Setting a file mode of `0000` will now disable Vault from
   automatically `chmod`ing the log file [[GH-3649](https://github.com/hashicorp/vault/pull/3649)]
 * auth/github: The legacy MFA system can now be used with the GitHub auth
   backend [[GH-3696](https://github.com/hashicorp/vault/pull/3696)]
 * auth/okta: The legacy MFA system can now be used with the Okta auth backend
   [[GH-3653](https://github.com/hashicorp/vault/pull/3653)]
 * auth/token: `allowed_policies` and `disallowed_policies` can now be specified
   as a comma-separated string or an array of strings [[GH-3641](https://github.com/hashicorp/vault/pull/3641)]
 * command/server: The log level can now be specified with `VAULT_LOG_LEVEL`
   [[GH-3721](https://github.com/hashicorp/vault/pull/3721)]
 * core: Period values from auth backends will now be checked and applied to the
   TTL value directly by core on login and renewal requests [[GH-3677](https://github.com/hashicorp/vault/pull/3677)]
 * database/mongodb: Add optional `write_concern` parameter, which can be set
   during database configuration. This establishes a session-wide [write
   concern](https://docs.mongodb.com/manual/reference/write-concern/) for the
   lifecycle of the mount [[GH-3646](https://github.com/hashicorp/vault/pull/3646)]
 * http: Request path containing non-printable characters will return 400 - Bad
   Request [[GH-3697](https://github.com/hashicorp/vault/pull/3697)]
 * mfa/okta: Filter a given email address as a login filter, allowing operation
   when login email and account email are different
 * plugins: Make Vault more resilient when unsealing when plugins are
   unavailable [[GH-3686](https://github.com/hashicorp/vault/pull/3686)]
 * secret/pki: `allowed_domains` and `key_usage` can now be specified
   as a comma-separated string or an array of strings [[GH-3642](https://github.com/hashicorp/vault/pull/3642)]
 * secret/ssh: Allow 4096-bit keys to be used in dynamic key method [[GH-3593](https://github.com/hashicorp/vault/pull/3593)]
 * secret/consul: The Consul secret backend now uses the value of `lease` set
   on the role, if set, when renewing a secret. [[GH-3796](https://github.com/hashicorp/vault/pull/3796)]
 * storage/mysql: Don't attempt database creation if it exists, which can help
   under certain permissions constraints [[GH-3716](https://github.com/hashicorp/vault/pull/3716)]

BUG FIXES:

 * api/status (enterprise): Fix status reporting when using an auto seal
 * auth/approle: Fix case-sensitive/insensitive comparison issue [[GH-3665](https://github.com/hashicorp/vault/pull/3665)]
 * auth/cert: Return `allowed_names` on role read [[GH-3654](https://github.com/hashicorp/vault/pull/3654)]
 * auth/ldap: Fix incorrect control information being sent [[GH-3402](https://github.com/hashicorp/vault/pull/3402)] [[GH-3496](https://github.com/hashicorp/vault/pull/3496)]
   [[GH-3625](https://github.com/hashicorp/vault/pull/3625)] [[GH-3656](https://github.com/hashicorp/vault/pull/3656)]
 * core: Fix seal status reporting when using an autoseal
 * core: Add creation path to wrap info for a control group token
 * core: Fix potential panic that could occur using plugins when a node
   transitioned from active to standby [[GH-3638](https://github.com/hashicorp/vault/pull/3638)]
 * core: Fix memory ballooning when a connection would connect to the cluster
   port and then go away -- redux! [[GH-3680](https://github.com/hashicorp/vault/pull/3680)]
 * core: Replace recursive token revocation logic with depth-first logic, which
   can avoid hitting stack depth limits in extreme cases [[GH-2348](https://github.com/hashicorp/vault/pull/2348)]
 * core: When doing a read on configured audited-headers, properly handle case
   insensitivity [[GH-3701](https://github.com/hashicorp/vault/pull/3701)]
 * core/pkcs11 (enterprise): Fix panic when PKCS#11 library is not readable
 * database/mysql: Allow the creation statement to use commands that are not yet
   supported by the prepare statement protocol [[GH-3619](https://github.com/hashicorp/vault/pull/3619)]
 * plugin/auth-gcp: Fix IAM roles when using `allow_gce_inference` [VPAG-19]

## 0.9.0.1 (November 21st, 2017) (Enterprise Only)

IMPROVEMENTS:

 * auth/gcp: Support seal wrapping of configuration parameters
 * auth/kubernetes: Support seal wrapping of configuration parameters

BUG FIXES:

 * Fix an upgrade issue with some physical backends when migrating from legacy
   HSM stored key support to the new Seal Wrap mechanism (Enterprise)
 * mfa: Add the 'mfa' flag that was removed by mistake [[GH-4223](https://github.com/hashicorp/vault/pull/4223)]

## 0.9.0 (November 14th, 2017)

DEPRECATIONS/CHANGES:

 * HSM config parameter requirements: When using Vault with an HSM, a new
   parameter is required: `hmac_key_label`.  This performs a similar function to
   `key_label` but for the HMAC key Vault will use. Vault will generate a
   suitable key if this value is specified and `generate_key` is set true.
 * API HTTP client behavior: When calling `NewClient` the API no longer
   modifies the provided client/transport. In particular this means it will no
   longer enable redirection limiting and HTTP/2 support on custom clients. It
   is suggested that if you want to make changes to an HTTP client that you use
   one created by `DefaultConfig` as a starting point.
 * AWS EC2 client nonce behavior: The client nonce generated by the backend
   that gets returned along with the authentication response will be audited in
   plaintext. If this is undesired, the clients can choose to supply a custom
   nonce to the login endpoint. The custom nonce set by the client will from
   now on, not be returned back with the authentication response, and hence not
   audit logged.
 * AWS Auth role options: The API will now error when trying to create or
   update a role with the mutually-exclusive options
   `disallow_reauthentication` and `allow_instance_migration`.
 * SSH CA role read changes: When reading back a role from the `ssh` backend,
   the TTL/max TTL values will now be an integer number of seconds rather than
   a string. This better matches the API elsewhere in Vault.
 * SSH role list changes: When listing roles from the `ssh` backend via the API,
   the response data will additionally return a `key_info` map that will contain
   a map of each key with a corresponding object containing the `key_type`.
 * More granularity in audit logs: Audit request and response entries are still
   in RFC3339 format but now have a granularity of nanoseconds.
 * High availability related values have been moved out of the `storage` and
   `ha_storage` stanzas, and into the top-level configuration. `redirect_addr`
   has been renamed to `api_addr`. The stanzas still support accepting
   HA-related values to maintain backward compatibility, but top-level values
   will take precedence.
 * A new `seal` stanza has been added to the configuration file, which is
   optional and enables configuration of the seal type to use for additional
   data protection, such as using HSM or Cloud KMS solutions to encrypt and
   decrypt data.

FEATURES:

 * **RSA Support for Transit Backend**: Transit backend can now generate RSA
   keys which can be used for encryption and signing. [[GH-3489](https://github.com/hashicorp/vault/pull/3489)]
 * **Identity System**: Now in open source and with significant enhancements,
   Identity is an integrated system for understanding users across tokens and
   enabling easier management of users directly and via groups.
 * **External Groups in Identity**: Vault can now automatically assign users
   and systems to groups in Identity based on their membership in external
   groups.
 * **Seal Wrap / FIPS 140-2 Compatibility (Enterprise)**: Vault can now take
   advantage of FIPS 140-2-certified HSMs to ensure that Critical Security
   Parameters are protected in a compliant fashion. Vault's implementation has
   received a statement of compliance from Leidos.
 * **Control Groups (Enterprise)**: Require multiple members of an Identity
   group to authorize a requested action before it is allowed to run.
 * **Cloud Auto-Unseal (Enterprise)**: Automatically unseal Vault using AWS KMS
   and GCP CKMS.
 * **Sentinel Integration (Enterprise)**: Take advantage of HashiCorp Sentinel
   to create extremely flexible access control policies -- even on
   unauthenticated endpoints.
 * **Barrier Rekey Support for Auto-Unseal (Enterprise)**: When using auto-unsealing
   functionality, the `rekey` operation is now supported; it uses recovery keys
   to authorize the master key rekey.
 * **Operation Token for Disaster Recovery Actions (Enterprise)**: When using
   Disaster Recovery replication, a token can be created that can be used to
   authorize actions such as promotion and updating primary information, rather
   than using recovery keys.
 * **Trigger Auto-Unseal with Recovery Keys (Enterprise)**: When using
   auto-unsealing, a request to unseal Vault can be triggered by a threshold of
   recovery keys, rather than requiring the Vault process to be restarted.
 * **UI Redesign (Enterprise)**: All new experience for the Vault Enterprise
   UI. The look and feel has been completely redesigned to give users a better
   experience and make managing secrets fast and easy.
 * **UI: SSH Secret Backend (Enterprise)**: Configure an SSH secret backend,
   create and browse roles. And use them to sign keys or generate one time
   passwords.
 * **UI: AWS Secret Backend (Enterprise)**: You can now configure the AWS
   backend via the Vault Enterprise UI. In addition you can create roles,
   browse the roles and Generate IAM Credentials from them in the UI.

IMPROVEMENTS:

 * api: Add ability to set custom headers on each call [[GH-3394](https://github.com/hashicorp/vault/pull/3394)]
 * command/server: Add config option to disable requesting client certificates
   [[GH-3373](https://github.com/hashicorp/vault/pull/3373)]
 * auth/aws: Max retries can now be customized for the AWS client [[GH-3965](https://github.com/hashicorp/vault/pull/3965)]
 * core: Disallow mounting underneath an existing path, not just over [[GH-2919](https://github.com/hashicorp/vault/pull/2919)]
 * physical/file: Use `700` as permissions when creating directories. The files
   themselves were `600` and are all encrypted, but this doesn't hurt.
 * secret/aws: Add ability to use custom IAM/STS endpoints [[GH-3416](https://github.com/hashicorp/vault/pull/3416)]
 * secret/aws: Max retries can now be customized for the AWS client [[GH-3965](https://github.com/hashicorp/vault/pull/3965)]
 * secret/cassandra: Work around Cassandra ignoring consistency levels for a
   user listing query [[GH-3469](https://github.com/hashicorp/vault/pull/3469)]
 * secret/pki: Private keys can now be marshalled as PKCS#8 [[GH-3518](https://github.com/hashicorp/vault/pull/3518)]
 * secret/pki: Allow entering URLs for `pki` as both comma-separated strings and JSON
   arrays [[GH-3409](https://github.com/hashicorp/vault/pull/3409)]
 * secret/ssh: Role TTL/max TTL can now be specified as either a string or an
   integer [[GH-3507](https://github.com/hashicorp/vault/pull/3507)]
 * secret/transit: Sign and verify operations now support a `none` hash
   algorithm to allow signing/verifying pre-hashed data [[GH-3448](https://github.com/hashicorp/vault/pull/3448)]
 * secret/database: Add the ability to glob allowed roles in the Database Backend [[GH-3387](https://github.com/hashicorp/vault/pull/3387)]
 * ui (enterprise): Support for RSA keys in the transit backend
 * ui (enterprise): Support for DR Operation Token generation, promoting, and
   updating primary on DR Secondary clusters

BUG FIXES:

 * api: Fix panic when setting a custom HTTP client but with a nil transport
   [[GH-3435](https://github.com/hashicorp/vault/pull/3435)] [[GH-3437](https://github.com/hashicorp/vault/pull/3437)]
 * api: Fix authing to the `cert` backend when the CA for the client cert is
   not known to the server's listener [[GH-2946](https://github.com/hashicorp/vault/pull/2946)]
 * auth/approle: Create role ID index during read if a role is missing one [[GH-3561](https://github.com/hashicorp/vault/pull/3561)]
 * auth/aws: Don't allow mutually exclusive options [[GH-3291](https://github.com/hashicorp/vault/pull/3291)]
 * auth/radius: Fix logging in in some situations [[GH-3461](https://github.com/hashicorp/vault/pull/3461)]
 * core: Fix memleak when a connection would connect to the cluster port and
   then go away [[GH-3513](https://github.com/hashicorp/vault/pull/3513)]
 * core: Fix panic if a single-use token is used to step-down or seal [[GH-3497](https://github.com/hashicorp/vault/pull/3497)]
 * core: Set rather than add headers to prevent some duplicated headers in
   responses when requests were forwarded to the active node [[GH-3485](https://github.com/hashicorp/vault/pull/3485)]
 * physical/etcd3: Fix some listing issues due to how etcd3 does prefix
   matching [[GH-3406](https://github.com/hashicorp/vault/pull/3406)]
 * physical/etcd3: Fix case where standbys can lose their etcd client lease
   [[GH-3031](https://github.com/hashicorp/vault/pull/3031)]
 * physical/file: Fix listing when underscores are the first component of a
   path [[GH-3476](https://github.com/hashicorp/vault/pull/3476)]
 * plugins: Allow response errors to be returned from backend plugins [[GH-3412](https://github.com/hashicorp/vault/pull/3412)]
 * secret/transit: Fix panic if the length of the input ciphertext was less
   than the expected nonce length [[GH-3521](https://github.com/hashicorp/vault/pull/3521)]
 * ui (enterprise): Reinstate support for generic secret backends - this was
   erroneously removed in a previous release

## 0.8.3 (September 19th, 2017)

CHANGES:

 * Policy input/output standardization: For all built-in authentication
   backends, policies can now be specified as a comma-delimited string or an
   array if using JSON as API input; on read, policies will be returned as an
   array; and the `default` policy will not be forcefully added to policies
   saved in configurations. Please note that the `default` policy will continue
   to be added to generated tokens, however, rather than backends adding
   `default` to the given set of input policies (in some cases, and not in
   others), the stored set will reflect the user-specified set.
 * `sign-self-issued` modifies Issuer in generated certificates: In 0.8.2 the
   endpoint would not modify the Issuer in the generated certificate, leaving
   the output self-issued. Although theoretically valid, in practice crypto
   stacks were unhappy validating paths containing such certs. As a result,
   `sign-self-issued` now encodes the signing CA's Subject DN into the Issuer
   DN of the generated certificate.
 * `sys/raw` requires enabling: While the `sys/raw` endpoint can be extremely
   useful in break-glass or support scenarios, it is also extremely dangerous.
   As of now, a configuration file option `raw_storage_endpoint` must be set in
   order to enable this API endpoint. Once set, the available functionality has
   been enhanced slightly; it now supports listing and decrypting most of
   Vault's core data structures, except for the encryption keyring itself.
 * `generic` is now `kv`: To better reflect its actual use, the `generic`
   backend is now `kv`. Using `generic` will still work for backwards
   compatibility.

FEATURES:

 * **GCE Support for GCP Auth**: GCE instances can now authenticate to Vault
   using machine credentials.
 * **Support for Kubernetes Service Account Auth**: Kubernetes Service Accounts
   can now authenticate to vault using JWT tokens.

IMPROVEMENTS:

 * configuration: Provide a config option to store Vault server's process ID
   (PID) in a file [[GH-3321](https://github.com/hashicorp/vault/pull/3321)]
 * mfa (Enterprise): Add the ability to use identity metadata in username format
 * mfa/okta (Enterprise): Add support for configuring base_url for API calls
 * secret/pki: `sign-intermediate` will now allow specifying a `ttl` value
   longer than the signing CA certificate's NotAfter value. [[GH-3325](https://github.com/hashicorp/vault/pull/3325)]
 * sys/raw: Raw storage access is now disabled by default [[GH-3329](https://github.com/hashicorp/vault/pull/3329)]

BUG FIXES:

 * auth/okta: Fix regression that removed the ability to set base_url [[GH-3313](https://github.com/hashicorp/vault/pull/3313)]
 * core: Fix panic while loading leases at startup on ARM processors
   [[GH-3314](https://github.com/hashicorp/vault/pull/3314)]
 * secret/pki: Fix `sign-self-issued` encoding the wrong subject public key
   [[GH-3325](https://github.com/hashicorp/vault/pull/3325)]

## 0.8.2.1 (September 11th, 2017) (Enterprise Only)

BUG FIXES:

 * Fix an issue upgrading to 0.8.2 for Enterprise customers.

## 0.8.2 (September 5th, 2017)

SECURITY:

* In prior versions of Vault, if authenticating via AWS IAM and requesting a
  periodic token, the period was not properly respected. This could lead to
  tokens expiring unexpectedly, or a token lifetime being longer than expected.
  Upon token renewal with Vault 0.8.2 the period will be properly enforced.

DEPRECATIONS/CHANGES:

* `vault ssh` users should supply `-mode` and `-role` to reduce the number of
  API calls. A future version of Vault will mark these optional values are
  required. Failure to supply `-mode` or `-role` will result in a warning.
* Vault plugins will first briefly run a restricted version of the plugin to
  fetch metadata, and then lazy-load the plugin on first request to prevent
  crash/deadlock of Vault during the unseal process. Plugins will need to be
  built with the latest changes in order for them to run properly.

FEATURES:

* **Lazy Lease Loading**: On startup, Vault will now load leases from storage
  in a lazy fashion (token checks and revocation/renewal requests still force
  an immediate load). For larger installations this can significantly reduce
  downtime when switching active nodes or bringing Vault up from cold start.
* **SSH CA Login with `vault ssh`**: `vault ssh` now supports the SSH CA
  backend for authenticating to machines. It also supports remote host key
  verification through the SSH CA backend, if enabled.
* **Signing of Self-Issued Certs in PKI**: The `pki` backend now supports
  signing self-issued CA certs. This is useful when switching root CAs.

IMPROVEMENTS:

 * audit/file: Allow specifying `stdout` as the `file_path` to log to standard
   output [[GH-3235](https://github.com/hashicorp/vault/pull/3235)]
 * auth/aws: Allow wildcards in `bound_iam_principal_arn` [[GH-3213](https://github.com/hashicorp/vault/pull/3213)]
 * auth/okta: Compare groups case-insensitively since Okta is only
   case-preserving [[GH-3240](https://github.com/hashicorp/vault/pull/3240)]
 * auth/okta: Standardize Okta configuration APIs across backends [[GH-3245](https://github.com/hashicorp/vault/pull/3245)]
 * cli: Add subcommand autocompletion that can be enabled with
   `vault -autocomplete-install` [[GH-3223](https://github.com/hashicorp/vault/pull/3223)]
 * cli: Add ability to handle wrapped responses when using `vault auth`. What
   is output depends on the other given flags; see the help output for that
   command for more information. [[GH-3263](https://github.com/hashicorp/vault/pull/3263)]
 * core: TLS cipher suites used for cluster behavior can now be set via
   `cluster_cipher_suites` in configuration [[GH-3228](https://github.com/hashicorp/vault/pull/3228)]
 * core: The `plugin_name` can now either be specified directly as part of the
   parameter or within the `config` object when mounting a secret or auth backend
   via `sys/mounts/:path` or `sys/auth/:path` respectively [[GH-3202](https://github.com/hashicorp/vault/pull/3202)]
 * core: It is now possible to update the `description` of a mount when
   mount-tuning, although this must be done through the HTTP layer [[GH-3285](https://github.com/hashicorp/vault/pull/3285)]
 * secret/databases/mongo: If an EOF is encountered, attempt reconnecting and
   retrying the operation [[GH-3269](https://github.com/hashicorp/vault/pull/3269)]
 * secret/pki: TTLs can now be specified as a string or an integer number of
   seconds [[GH-3270](https://github.com/hashicorp/vault/pull/3270)]
 * secret/pki: Self-issued certs can now be signed via
   `pki/root/sign-self-issued` [[GH-3274](https://github.com/hashicorp/vault/pull/3274)]
 * storage/gcp: Use application default credentials if they exist [[GH-3248](https://github.com/hashicorp/vault/pull/3248)]

BUG FIXES:

 * auth/aws: Properly use role-set period values for IAM-derived token renewals
   [[GH-3220](https://github.com/hashicorp/vault/pull/3220)]
 * auth/okta: Fix updating organization/ttl/max_ttl after initial setting
   [[GH-3236](https://github.com/hashicorp/vault/pull/3236)]
 * core: Fix PROXY when underlying connection is TLS [[GH-3195](https://github.com/hashicorp/vault/pull/3195)]
 * core: Policy-related commands would sometimes fail to act case-insensitively
   [[GH-3210](https://github.com/hashicorp/vault/pull/3210)]
 * storage/consul: Fix parsing TLS configuration when using a bare IPv6 address
   [[GH-3268](https://github.com/hashicorp/vault/pull/3268)]
 * plugins: Lazy-load plugins to prevent crash/deadlock during unseal process.
   [[GH-3255](https://github.com/hashicorp/vault/pull/3255)]
 * plugins: Skip mounting plugin-based secret and credential mounts when setting
   up mounts if the plugin is no longer present in the catalog. [[GH-3255](https://github.com/hashicorp/vault/pull/3255)]

## 0.8.1 (August 16th, 2017)

DEPRECATIONS/CHANGES:

 * PKI Root Generation: Calling `pki/root/generate` when a CA cert/key already
   exists will now return a `204` instead of overwriting an existing root. If
   you want to recreate the root, first run a delete operation on `pki/root`
   (requires `sudo` capability), then generate it again.

FEATURES:

 * **Oracle Secret Backend**: There is now an external plugin to support leased
   credentials for Oracle databases (distributed separately).
 * **GCP IAM Auth Backend**: There is now an authentication backend that allows
   using GCP IAM credentials to retrieve Vault tokens. This is available as
   both a plugin and built-in to Vault.
 * **PingID Push Support for Path-Based MFA (Enterprise)**: PingID Push can
   now be used for MFA with the new path-based MFA introduced in Vault
   Enterprise 0.8.
 * **Permitted DNS Domains Support in PKI**: The `pki` backend now supports
   specifying permitted DNS domains for CA certificates, allowing you to
   narrowly scope the set of domains for which a CA can issue or sign child
   certificates.
 * **Plugin Backend Reload Endpoint**: Plugin backends can now be triggered to
   reload using the `sys/plugins/reload/backend` endpoint and providing either
   the plugin name or the mounts to reload.
 * **Self-Reloading Plugins**: The plugin system will now attempt to reload a
   crashed or stopped plugin, once per request.

IMPROVEMENTS:

 * auth/approle: Allow array input for policies in addition to comma-delimited
   strings [[GH-3163](https://github.com/hashicorp/vault/pull/3163)]
 * plugins: Send logs through Vault's logger rather than stdout [[GH-3142](https://github.com/hashicorp/vault/pull/3142)]
 * secret/pki: Add `pki/root` delete operation [[GH-3165](https://github.com/hashicorp/vault/pull/3165)]
 * secret/pki: Don't overwrite an existing root cert/key when calling generate
   [[GH-3165](https://github.com/hashicorp/vault/pull/3165)]

BUG FIXES:

 * aws: Don't prefer a nil HTTP client over an existing one [[GH-3159](https://github.com/hashicorp/vault/pull/3159)]
 * core: If there is an error when checking for create/update existence, return
   500 instead of 400 [[GH-3162](https://github.com/hashicorp/vault/pull/3162)]
 * secret/database: Avoid creating usernames that are too long for legacy MySQL
   [[GH-3138](https://github.com/hashicorp/vault/pull/3138)]

## 0.8.0 (August 9th, 2017)

SECURITY:

 * We've added a note to the docs about the way the GitHub auth backend works
   as it may not be readily apparent that GitHub personal access tokens, which
   are used by the backend, can be used for unauthorized access if they are
   stolen from third party services and access to Vault is public.

DEPRECATIONS/CHANGES:

 * Database Plugin Backends: Passwords generated for these backends now
   enforce stricter password requirements, as opposed to the previous behavior
   of returning a randomized UUID. Passwords are of length 20, and have a `A1a-`
   characters prepended to ensure stricter requirements. No regressions are
   expected from this change. (For database backends that were previously
   substituting underscores for hyphens in passwords, this will remain the
   case.)
 * Lease Endpoints: The endpoints `sys/renew`, `sys/revoke`, `sys/revoke-prefix`,
   `sys/revoke-force` have been deprecated and relocated under `sys/leases`.
   Additionally, the deprecated path `sys/revoke-force` now requires the `sudo`
   capability.
 * Response Wrapping Lookup Unauthenticated: The `sys/wrapping/lookup` endpoint
   is now unauthenticated. This allows introspection of the wrapping info by
   clients that only have the wrapping token without then invalidating the
   token. Validation functions/checks are still performed on the token.

FEATURES:

 * **Cassandra Storage**: Cassandra can now be used for Vault storage
 * **CockroachDB Storage**: CockroachDB can now be used for Vault storage
 * **CouchDB Storage**: CouchDB can now be used for Vault storage
 * **SAP HANA Database Plugin**: The `databases` backend can now manage users
   for SAP HANA databases
 * **Plugin Backends**: Vault now supports running secret and auth backends as
   plugins. Plugins can be mounted like normal backends and can be developed
   independently from Vault.
 * **PROXY Protocol Support** Vault listeners can now be configured to honor
   PROXY protocol v1 information to allow passing real client IPs into Vault. A
   list of authorized addresses (IPs or subnets) can be defined and
   accept/reject behavior controlled.
 * **Lease Lookup and Browsing in the Vault Enterprise UI**: Vault Enterprise UI
   now supports lookup and listing of leases and the associated actions from the
   `sys/leases` endpoints in the API. These are located in the new top level
   navigation item "Leases".
 * **Filtered Mounts for Performance Mode Replication**: Whitelists or
   blacklists of mounts can be defined per-secondary to control which mounts
   are actually replicated to that secondary. This can allow targeted
   replication of specific sets of data to specific geolocations/datacenters.
 * **Disaster Recovery Mode Replication (Enterprise Only)**: There is a new
   replication mode, Disaster Recovery (DR), that performs full real-time
   replication (including tokens and leases) to DR secondaries. DR secondaries
   cannot handle client requests, but can be promoted to primary as needed for
   failover.
 * **Manage New Replication Features in the Vault Enterprise UI**: Support for
   Replication features in Vault Enterprise UI has expanded to include new DR
   Replication mode and management of Filtered Mounts in Performance Replication
   mode.
 * **Vault Identity (Enterprise Only)**: Vault's new Identity system allows
   correlation of users across tokens. At present this is only used for MFA,
   but will be the foundation of many other features going forward.
 * **Duo Push, Okta Push, and TOTP MFA For All Authenticated Paths (Enterprise
   Only)**: A brand new MFA system built on top of Identity allows MFA
   (currently Duo Push, Okta Push, and TOTP) for any authenticated path within
   Vault. MFA methods can be configured centrally, and TOTP keys live within
   the user's Identity information to allow using the same key across tokens.
   Specific MFA method(s) required for any given path within Vault can be
   specified in normal ACL path statements.

IMPROVEMENTS:

 * api: Add client method for a secret renewer background process [[GH-2886](https://github.com/hashicorp/vault/pull/2886)]
 * api: Add `RenewTokenAsSelf` [[GH-2886](https://github.com/hashicorp/vault/pull/2886)]
 * api: Client timeout can now be adjusted with the `VAULT_CLIENT_TIMEOUT` env
   var or with a new API function [[GH-2956](https://github.com/hashicorp/vault/pull/2956)]
 * api/cli: Client will now attempt to look up SRV records for the given Vault
   hostname [[GH-3035](https://github.com/hashicorp/vault/pull/3035)]
 * audit/socket: Enhance reconnection logic and don't require the connection to
   be established at unseal time [[GH-2934](https://github.com/hashicorp/vault/pull/2934)]
 * audit/file: Opportunistically try re-opening the file on error [[GH-2999](https://github.com/hashicorp/vault/pull/2999)]
 * auth/approle: Add role name to token metadata [[GH-2985](https://github.com/hashicorp/vault/pull/2985)]
 * auth/okta: Allow specifying `ttl`/`max_ttl` inside the mount [[GH-2915](https://github.com/hashicorp/vault/pull/2915)]
 * cli: Client timeout can now be adjusted with the `VAULT_CLIENT_TIMEOUT` env
   var [[GH-2956](https://github.com/hashicorp/vault/pull/2956)]
 * command/auth: Add `-token-only` flag to `vault auth` that returns only the
   token on stdout and does not store it via the token helper [[GH-2855](https://github.com/hashicorp/vault/pull/2855)]
 * core: CORS allowed origins can now be configured [[GH-2021](https://github.com/hashicorp/vault/pull/2021)]
 * core: Add metrics counters for audit log failures [[GH-2863](https://github.com/hashicorp/vault/pull/2863)]
 * cors: Allow setting allowed headers via the API instead of always using
   wildcard [[GH-3023](https://github.com/hashicorp/vault/pull/3023)]
 * secret/ssh: Allow specifying the key ID format using template values for CA
   type [[GH-2888](https://github.com/hashicorp/vault/pull/2888)]
 * server: Add `tls_client_ca_file` option for specifying a CA file to use for
   client certificate verification when `tls_require_and_verify_client_cert` is
   enabled [[GH-3034](https://github.com/hashicorp/vault/pull/3034)]
 * storage/cockroachdb: Add CockroachDB storage backend [[GH-2713](https://github.com/hashicorp/vault/pull/2713)]
 * storage/couchdb: Add CouchDB storage backend [[GH-2880](https://github.com/hashicorp/vault/pull/2880)]
 * storage/mssql: Add `max_parallel` [[GH-3026](https://github.com/hashicorp/vault/pull/3026)]
 * storage/postgresql: Add `max_parallel` [[GH-3026](https://github.com/hashicorp/vault/pull/3026)]
 * storage/postgresql: Improve listing speed [[GH-2945](https://github.com/hashicorp/vault/pull/2945)]
 * storage/s3: More efficient paging when an object has a lot of subobjects
   [[GH-2780](https://github.com/hashicorp/vault/pull/2780)]
 * sys/wrapping: Make `sys/wrapping/lookup` unauthenticated [[GH-3084](https://github.com/hashicorp/vault/pull/3084)]
 * sys/wrapping: Wrapped tokens now store the original request path of the data
   [[GH-3100](https://github.com/hashicorp/vault/pull/3100)]
 * telemetry: Add support for DogStatsD [[GH-2490](https://github.com/hashicorp/vault/pull/2490)]

BUG FIXES:

 * api/health: Don't treat standby `429` codes as an error [[GH-2850](https://github.com/hashicorp/vault/pull/2850)]
 * api/leases: Fix lease lookup returning lease properties at the top level
 * audit: Fix panic when audit logging a read operation on an asymmetric
   `transit` key [[GH-2958](https://github.com/hashicorp/vault/pull/2958)]
 * auth/approle: Fix panic when secret and cidr list not provided in role
   [[GH-3075](https://github.com/hashicorp/vault/pull/3075)]
 * auth/aws: Look up proper account ID on token renew [[GH-3012](https://github.com/hashicorp/vault/pull/3012)]
 * auth/aws: Store IAM header in all cases when it changes [[GH-3004](https://github.com/hashicorp/vault/pull/3004)]
 * auth/ldap: Verify given certificate is PEM encoded instead of failing
   silently [[GH-3016](https://github.com/hashicorp/vault/pull/3016)]
 * auth/token: Don't allow using the same token ID twice when manually
   specifying [[GH-2916](https://github.com/hashicorp/vault/pull/2916)]
 * cli: Fix issue with parsing keys that start with special characters [[GH-2998](https://github.com/hashicorp/vault/pull/2998)]
 * core: Relocated `sys/leases/renew` returns same payload as original
   `sys/leases` endpoint [[GH-2891](https://github.com/hashicorp/vault/pull/2891)]
 * secret/ssh: Fix panic when signing with incorrect key type [[GH-3072](https://github.com/hashicorp/vault/pull/3072)]
 * secret/totp: Ensure codes can only be used once. This makes some automated
   workflows harder but complies with the RFC. [[GH-2908](https://github.com/hashicorp/vault/pull/2908)]
 * secret/transit: Fix locking when creating a key with unsupported options
   [[GH-2974](https://github.com/hashicorp/vault/pull/2974)]

## 0.7.3 (June 7th, 2017)

SECURITY:

 * Cert auth backend now checks validity of individual certificates: In
   previous versions of Vault, validity (e.g. expiration) of individual leaf
   certificates added for authentication was not checked. This was done to make
   it easier for administrators to control lifecycles of individual
   certificates added to the backend, e.g. the authentication material being
   checked was access to that specific certificate's private key rather than
   all private keys signed by a CA. However, this behavior is often unexpected
   and as a result can lead to insecure deployments, so we are now validating
   these certificates as well.
 * App-ID path salting was skipped in 0.7.1/0.7.2: A regression in 0.7.1/0.7.2
   caused the HMACing of any App-ID information stored in paths (including
   actual app-IDs and user-IDs) to be unsalted and written as-is from the API.
   In 0.7.3 any such paths will be automatically changed to salted versions on
   access (e.g. login or read); however, if you created new app-IDs or user-IDs
   in 0.7.1/0.7.2, you may want to consider whether any users with access to
   Vault's underlying data store may have intercepted these values, and
   revoke/roll them.

DEPRECATIONS/CHANGES:

 * Step-Down is Forwarded: When a step-down is issued against a non-active node
   in an HA cluster, it will now forward the request to the active node.

FEATURES:

 * **ed25519 Signing/Verification in Transit with Key Derivation**: The
   `transit` backend now supports generating
   [ed25519](https://ed25519.cr.yp.to/) keys for signing and verification
   functionality. These keys support derivation, allowing you to modify the
   actual encryption key used by supplying a `context` value.
 * **Key Version Specification for Encryption in Transit**: You can now specify
   the version of a key you use to wish to generate a signature, ciphertext, or
   HMAC. This can be controlled by the `min_encryption_version` key
   configuration property.
 * **Replication Primary Discovery (Enterprise)**: Replication primaries will
   now advertise the addresses of their local HA cluster members to replication
   secondaries. This helps recovery if the primary active node goes down and
   neither service discovery nor load balancers are in use to steer clients.

IMPROVEMENTS:

 * api/health: Add Sys().Health() [[GH-2805](https://github.com/hashicorp/vault/pull/2805)]
 * audit: Add auth information to requests that error out [[GH-2754](https://github.com/hashicorp/vault/pull/2754)]
 * command/auth: Add `-no-store` option that prevents the auth command from
   storing the returned token into the configured token helper [[GH-2809](https://github.com/hashicorp/vault/pull/2809)]
 * core/forwarding: Request forwarding now heartbeats to prevent unused
   connections from being terminated by firewalls or proxies
 * plugins/databases: Add MongoDB as an internal database plugin [[GH-2698](https://github.com/hashicorp/vault/pull/2698)]
 * storage/dynamodb: Add a method for checking the existence of children,
   speeding up deletion operations in the DynamoDB storage backend [[GH-2722](https://github.com/hashicorp/vault/pull/2722)]
 * storage/mysql: Add max_parallel parameter to MySQL backend [[GH-2760](https://github.com/hashicorp/vault/pull/2760)]
 * secret/databases: Support listing connections [[GH-2823](https://github.com/hashicorp/vault/pull/2823)]
 * secret/databases: Support custom renewal statements in Postgres database
   plugin [[GH-2788](https://github.com/hashicorp/vault/pull/2788)]
 * secret/databases: Use the role name as part of generated credentials
   [[GH-2812](https://github.com/hashicorp/vault/pull/2812)]
 * ui (Enterprise): Transit key and secret browsing UI handle large lists better
 * ui (Enterprise): root tokens are no longer persisted
 * ui (Enterprise): support for mounting Database and TOTP secret backends

BUG FIXES:

 * auth/app-id: Fix regression causing loading of salts to be skipped
 * auth/aws: Improve EC2 describe instances performance [[GH-2766](https://github.com/hashicorp/vault/pull/2766)]
 * auth/aws: Fix lookup of some instance profile ARNs [[GH-2802](https://github.com/hashicorp/vault/pull/2802)]
 * auth/aws: Resolve ARNs to internal AWS IDs which makes lookup at various
   points (e.g. renewal time) more robust [[GH-2814](https://github.com/hashicorp/vault/pull/2814)]
 * auth/aws: Properly honor configured period when using IAM authentication
   [[GH-2825](https://github.com/hashicorp/vault/pull/2825)]
 * auth/aws: Check that a bound IAM principal is not empty (in the current
   state of the role) before requiring it match the previously authenticated
   client [[GH-2781](https://github.com/hashicorp/vault/pull/2781)]
 * auth/cert: Fix panic on renewal [[GH-2749](https://github.com/hashicorp/vault/pull/2749)]
 * auth/cert: Certificate verification for non-CA certs [[GH-2761](https://github.com/hashicorp/vault/pull/2761)]
 * core/acl: Prevent race condition when compiling ACLs in some scenarios
   [[GH-2826](https://github.com/hashicorp/vault/pull/2826)]
 * secret/database: Increase wrapping token TTL; in a loaded scenario it could
   be too short
 * secret/generic: Allow integers to be set as the value of `ttl` field as the
   documentation claims is supported [[GH-2699](https://github.com/hashicorp/vault/pull/2699)]
 * secret/ssh: Added host key callback to ssh client config [[GH-2752](https://github.com/hashicorp/vault/pull/2752)]
 * storage/s3: Avoid a panic when some bad data is returned [[GH-2785](https://github.com/hashicorp/vault/pull/2785)]
 * storage/dynamodb: Fix list functions working improperly on Windows [[GH-2789](https://github.com/hashicorp/vault/pull/2789)]
 * storage/file: Don't leak file descriptors in some error cases
 * storage/swift: Fix pre-v3 project/tenant name reading [[GH-2803](https://github.com/hashicorp/vault/pull/2803)]

## 0.7.2 (May 8th, 2017)

BUG FIXES:

 * audit: Fix auditing entries containing certain kinds of time values
   [[GH-2689](https://github.com/hashicorp/vault/pull/2689)]

## 0.7.1 (May 5th, 2017)

DEPRECATIONS/CHANGES:

 * LDAP Auth Backend: Group membership queries will now run as the `binddn`
   user when `binddn`/`bindpass` are configured, rather than as the
   authenticating user as was the case previously.

FEATURES:

 * **AWS IAM Authentication**: IAM principals can get Vault tokens
   automatically, opening AWS-based authentication to users, ECS containers,
   Lambda instances, and more. Signed client identity information retrieved
   using the AWS API `sts:GetCallerIdentity` is validated against the AWS STS
   service before issuing a Vault token. This backend is unified with the
   `aws-ec2` authentication backend under the name `aws`, and allows additional
   EC2-related restrictions to be applied during the IAM authentication; the
   previous EC2 behavior is also still available. [[GH-2441](https://github.com/hashicorp/vault/pull/2441)]
 * **MSSQL Physical Backend**: You can now use Microsoft SQL Server as your
   Vault physical data store [[GH-2546](https://github.com/hashicorp/vault/pull/2546)]
 * **Lease Listing and Lookup**: You can now introspect a lease to get its
   creation and expiration properties via `sys/leases/lookup`; with `sudo`
   capability you can also list leases for lookup, renewal, or revocation via
   that endpoint. Various lease functions (renew, revoke, revoke-prefix,
   revoke-force) have also been relocated to `sys/leases/`, but they also work
   at the old paths for compatibility. Reading (but not listing) leases via
   `sys/leases/lookup` is now a part of the current `default` policy. [[GH-2650](https://github.com/hashicorp/vault/pull/2650)]
 * **TOTP Secret Backend**: You can now store multi-factor authentication keys
   in Vault and use the API to retrieve time-based one-time use passwords on
   demand. The backend can also be used to generate a new key and validate
   passwords generated by that key. [[GH-2492](https://github.com/hashicorp/vault/pull/2492)]
 * **Database Secret Backend & Secure Plugins (Beta)**: This new secret backend
   combines the functionality of the MySQL, PostgreSQL, MSSQL, and Cassandra
   backends. It also provides a plugin interface for extendability through
   custom databases. [[GH-2200](https://github.com/hashicorp/vault/pull/2200)]

IMPROVEMENTS:

 * auth/cert: Support for constraints on subject Common Name and DNS/email
   Subject Alternate Names in certificates [[GH-2595](https://github.com/hashicorp/vault/pull/2595)]
 * auth/ldap: Use the binding credentials to search group membership rather
   than the user credentials [[GH-2534](https://github.com/hashicorp/vault/pull/2534)]
 * cli/revoke: Add `-self` option to allow revoking the currently active token
   [[GH-2596](https://github.com/hashicorp/vault/pull/2596)]
 * core: Randomize x coordinate in Shamir shares [[GH-2621](https://github.com/hashicorp/vault/pull/2621)]
 * replication: Fix a bug when enabling `approle` on a primary before
   secondaries were connected
 * replication: Add heartbeating to ensure firewalls don't kill connections to
   primaries
 * secret/pki: Add `no_store` option that allows certificates to be issued
   without being stored. This removes the ability to look up and/or add to a
   CRL but helps with scaling to very large numbers of certificates. [[GH-2565](https://github.com/hashicorp/vault/pull/2565)]
 * secret/pki: If used with a role parameter, the `sign-verbatim/<role>`
   endpoint honors the values of `generate_lease`, `no_store`, `ttl` and
   `max_ttl` from the given role [[GH-2593](https://github.com/hashicorp/vault/pull/2593)]
 * secret/pki: Add role parameter `allow_glob_domains` that enables defining
   names in `allowed_domains` containing `*` glob patterns [[GH-2517](https://github.com/hashicorp/vault/pull/2517)]
 * secret/pki: Update certificate storage to not use characters that are not
   supported on some filesystems [[GH-2575](https://github.com/hashicorp/vault/pull/2575)]
 * storage/etcd3: Add `discovery_srv` option to query for SRV records to find
   servers [[GH-2521](https://github.com/hashicorp/vault/pull/2521)]
 * storage/s3: Support `max_parallel` option to limit concurrent outstanding
   requests [[GH-2466](https://github.com/hashicorp/vault/pull/2466)]
 * storage/s3: Use pooled transport for http client [[GH-2481](https://github.com/hashicorp/vault/pull/2481)]
 * storage/swift: Allow domain values for V3 authentication [[GH-2554](https://github.com/hashicorp/vault/pull/2554)]
 * tidy: Improvements to `auth/token/tidy` and `sys/leases/tidy` to handle more
   cleanup cases [[GH-2452](https://github.com/hashicorp/vault/pull/2452)]

BUG FIXES:

 * api: Respect a configured path in Vault's address [[GH-2588](https://github.com/hashicorp/vault/pull/2588)]
 * auth/aws-ec2: New bounds added as criteria to allow role creation [[GH-2600](https://github.com/hashicorp/vault/pull/2600)]
 * auth/ldap: Don't lowercase groups attached to users [[GH-2613](https://github.com/hashicorp/vault/pull/2613)]
 * cli: Don't panic if `vault write` is used with the `force` flag but no path
   [[GH-2674](https://github.com/hashicorp/vault/pull/2674)]
 * core: Help operations should request forward since standbys may not have
   appropriate info [[GH-2677](https://github.com/hashicorp/vault/pull/2677)]
 * replication: Fix enabling secondaries when certain mounts already existed on
   the primary
 * secret/mssql: Update mssql driver to support queries with colons [[GH-2610](https://github.com/hashicorp/vault/pull/2610)]
 * secret/pki: Don't lowercase O/OU values in certs [[GH-2555](https://github.com/hashicorp/vault/pull/2555)]
 * secret/pki: Don't attempt to validate IP SANs if none are provided [[GH-2574](https://github.com/hashicorp/vault/pull/2574)]
 * secret/ssh: Don't automatically lowercase principles in issued SSH certs
   [[GH-2591](https://github.com/hashicorp/vault/pull/2591)]
 * storage/consul: Properly handle state events rather than timing out
   [[GH-2548](https://github.com/hashicorp/vault/pull/2548)]
 * storage/etcd3: Ensure locks are released if client is improperly shut down
   [[GH-2526](https://github.com/hashicorp/vault/pull/2526)]

## 0.7.0 (March 21th, 2017)

SECURITY:

 * Common name not being validated when `exclude_cn_from_sans` option used in
   `pki` backend: When using a role in the `pki` backend that specified the
   `exclude_cn_from_sans` option, the common name would not then be properly
   validated against the role's constraints. This has been fixed. We recommend
   any users of this feature to upgrade to 0.7 as soon as feasible.

DEPRECATIONS/CHANGES:

 * List Operations Always Use Trailing Slash: Any list operation, whether via
   the `GET` or `LIST` HTTP verb, will now internally canonicalize the path to
   have a trailing slash. This makes policy writing more predictable, as it
   means clients will no longer work or fail based on which client they're
   using or which HTTP verb they're using. However, it also means that policies
   allowing `list` capability must be carefully checked to ensure that they
   contain a trailing slash; some policies may need to be split into multiple
   stanzas to accommodate.
 * PKI Defaults to Unleased Certificates: When issuing certificates from the
   PKI backend, by default, no leases will be issued. If you want to manually
   revoke a certificate, its serial number can be used with the `pki/revoke`
   endpoint. Issuing leases is still possible by enabling the `generate_lease`
   toggle in PKI role entries (this will default to `true` for upgrades, to
   keep existing behavior), which will allow using lease IDs to revoke
   certificates. For installations issuing large numbers of certificates (tens
   to hundreds of thousands, or millions), this will significantly improve
   Vault startup time since leases associated with these certificates will not
   have to be loaded; however note that it also means that revocation of a
   token used to issue certificates will no longer add these certificates to a
   CRL. If this behavior is desired or needed, consider keeping leases enabled
   and ensuring lifetimes are reasonable, and issue long-lived certificates via
   a different role with leases disabled.

FEATURES:

 * **Replication (Enterprise)**: Vault Enterprise now has support for creating
   a multi-datacenter replication set between clusters. The current replication
   offering is based on an asynchronous primary/secondary (1:N) model that
   replicates static data while keeping dynamic data (leases, tokens)
   cluster-local, focusing on horizontal scaling for high-throughput and
   high-fanout deployments.
 * **Response Wrapping & Replication in the Vault Enterprise UI**: Vault
   Enterprise UI now supports looking up and rotating response wrapping tokens,
   as well as creating tokens with arbitrary values inside. It also now
   supports replication functionality, enabling the configuration of a
   replication set in the UI.
 * **Expanded Access Control Policies**: Access control policies can now
   specify allowed and denied parameters -- and, optionally, their values -- to
   control what a client can and cannot submit during an API call. Policies can
   also specify minimum/maximum response wrapping TTLs to both enforce the use
   of response wrapping and control the duration of resultant wrapping tokens.
   See the [policies concepts
   page](https://www.vaultproject.io/docs/concepts/policies.html) for more
   information.
 * **SSH Backend As Certificate Authority**: The SSH backend can now be
   configured to sign host and user certificates. Each mount of the backend
   acts as an independent signing authority. The CA key pair can be configured
   for each mount and the public key is accessible via an unauthenticated API
   call; additionally, the backend can generate a public/private key pair for
   you. We recommend using separate mounts for signing host and user
   certificates.

IMPROVEMENTS:

 * api/request: Passing username and password information in API request
   [GH-2469]
 * audit: Logging the token's use count with authentication response and
   logging the remaining uses of the client token with request [GH-2437]
 * auth/approle: Support for restricting the number of uses on the tokens
   issued [GH-2435]
 * auth/aws-ec2: AWS EC2 auth backend now supports constraints for VPC ID,
   Subnet ID and Region [GH-2407]
 * auth/ldap: Use the value of the `LOGNAME` or `USER` env vars for the
   username if not explicitly set on the command line when authenticating
   [GH-2154]
 * audit: Support adding a configurable prefix (such as `@cee`) before each
   line [GH-2359]
 * core: Canonicalize list operations to use a trailing slash [GH-2390]
 * core: Add option to disable caching on a per-mount level [GH-2455]
 * core: Add ability to require valid client certs in listener config [GH-2457]
 * physical/dynamodb: Implement a session timeout to avoid having to use
   recovery mode in the case of an unclean shutdown, which makes HA much safer
   [GH-2141]
 * secret/pki: O (Organization) values can now be set to role-defined values
   for issued/signed certificates [GH-2369]
 * secret/pki: Certificates issued/signed from PKI backend do not generate
   leases by default [GH-2403]
 * secret/pki: When using DER format, still return the private key type
   [GH-2405]
 * secret/pki: Add an intermediate to the CA chain even if it lacks an
   authority key ID [GH-2465]
 * secret/pki: Add role option to use CSR SANs [GH-2489]
 * secret/ssh: SSH backend as CA to sign user and host certificates [GH-2208]
 * secret/ssh: Support reading of SSH CA public key from `config/ca` endpoint
   and also return it when CA key pair is generated [GH-2483]

BUG FIXES:

 * audit: When auditing headers use case-insensitive comparisons [GH-2362]
 * auth/aws-ec2: Return role period in seconds and not nanoseconds [GH-2374]
 * auth/okta: Fix panic if user had no local groups and/or policies set
   [GH-2367]
 * command/server: Fix parsing of redirect address when port is not mentioned
   [GH-2354]
 * physical/postgresql: Fix listing returning incorrect results if there were
   multiple levels of children [GH-2393]

## 0.6.5 (February 7th, 2017)

FEATURES:

 * **Okta Authentication**: A new Okta authentication backend allows you to use
   Okta usernames and passwords to authenticate to Vault. If provided with an
   appropriate Okta API token, group membership can be queried to assign
   policies; users and groups can be defined locally as well.
 * **RADIUS Authentication**: A new RADIUS authentication backend allows using
   a RADIUS server to authenticate to Vault. Policies can be configured for
   specific users or for any authenticated user.
 * **Exportable Transit Keys**: Keys in `transit` can now be marked as
   `exportable` at creation time. This allows a properly ACL'd user to retrieve
   the associated signing key, encryption key, or HMAC key. The `exportable`
   value is returned on a key policy read and cannot be changed, so if a key is
   marked `exportable` it will always be exportable, and if it is not it will
   never be exportable.
 * **Batch Transit Operations**: `encrypt`, `decrypt` and `rewrap` operations
   in the transit backend now support processing multiple input items in one
   call, returning the output of each item in the response.
 * **Configurable Audited HTTP Headers**: You can now specify headers that you
   want to have included in each audit entry, along with whether each header
   should be HMAC'd or kept plaintext. This can be useful for adding additional
   client or network metadata to the audit logs.
 * **Transit Backend UI (Enterprise)**: Vault Enterprise UI now supports the transit
   backend, allowing creation, viewing and editing of named keys as well as using
   those keys to perform supported transit operations directly in the UI.
 * **Socket Audit Backend** A new socket audit backend allows audit logs to be sent
   through TCP, UDP, or UNIX Sockets.

IMPROVEMENTS:

 * auth/aws-ec2: Add support for cross-account auth using STS [GH-2148]
 * auth/aws-ec2: Support issuing periodic tokens [GH-2324]
 * auth/github: Support listing teams and users [GH-2261]
 * auth/ldap: Support adding policies to local users directly, in addition to
   local groups [GH-2152]
 * command/server: Add ability to select and prefer server cipher suites
   [GH-2293]
 * core: Add a nonce to unseal operations as a check (useful mostly for
   support, not as a security principle) [GH-2276]
 * duo: Added ability to supply extra context to Duo pushes [GH-2118]
 * physical/consul: Add option for setting consistency mode on Consul gets
   [GH-2282]
 * physical/etcd: Full v3 API support; code will autodetect which API version
   to use. The v3 code path is significantly less complicated and may be much
   more stable. [GH-2168]
 * secret/pki: Allow specifying OU entries in generated certificate subjects
   [GH-2251]
 * secret mount ui (Enterprise): the secret mount list now shows all mounted
   backends even if the UI cannot browse them. Additional backends can now be
   mounted from the UI as well.

BUG FIXES:

 * auth/token: Fix regression in 0.6.4 where using token store roles as a
   blacklist (with only `disallowed_policies` set) would not work in most
   circumstances [GH-2286]
 * physical/s3: Page responses in client so list doesn't truncate [GH-2224]
 * secret/cassandra: Stop a connection leak that could occur on active node
   failover [GH-2313]
 * secret/pki: When using `sign-verbatim`, don't require a role and use the
   CSR's common name [GH-2243]

## 0.6.4 (December 16, 2016)

SECURITY:

Further details about these security issues can be found in the 0.6.4 upgrade
guide.

 * `default` Policy Privilege Escalation: If a parent token did not have the
   `default` policy attached to its token, it could still create children with
   the `default` policy. This is no longer allowed (unless the parent has
   `sudo` capability for the creation path). In most cases this is low severity
   since the access grants in the `default` policy are meant to be access
   grants that are acceptable for all tokens to have.
 * Leases Not Expired When Limited Use Token Runs Out of Uses: When using
   limited-use tokens to create leased secrets, if the limited-use token was
   revoked due to running out of uses (rather than due to TTL expiration or
   explicit revocation) it would fail to revoke the leased secrets. These
   secrets would still be revoked when their TTL expired, limiting the severity
   of this issue. An endpoint has been added (`auth/token/tidy`) that can
   perform housekeeping tasks on the token store; one of its tasks can detect
   this situation and revoke the associated leases.

FEATURES:

  * **Policy UI (Enterprise)**: Vault Enterprise UI now supports viewing,
    creating, and editing policies.

IMPROVEMENTS:

 * http: Vault now sets a `no-store` cache control header to make it more
   secure in setups that are not end-to-end encrypted [GH-2183]

BUG FIXES:

 * auth/ldap: Don't panic if dialing returns an error and starttls is enabled;
   instead, return the error [GH-2188]
 * ui (Enterprise): Submitting an unseal key now properly resets the
   form so a browser refresh isn't required to continue.

## 0.6.3 (December 6, 2016)

DEPRECATIONS/CHANGES:

 * Request size limitation: A maximum request size of 32MB is imposed to
   prevent a denial of service attack with arbitrarily large requests [GH-2108]
 * LDAP denies passwordless binds by default: In new LDAP mounts, or when
   existing LDAP mounts are rewritten, passwordless binds will be denied by
   default. The new `deny_null_bind` parameter can be set to `false` to allow
   these. [GH-2103]
 * Any audit backend activated satisfies conditions: Previously, when a new
   Vault node was taking over service in an HA cluster, all audit backends were
   required to be loaded successfully to take over active duty. This behavior
   now matches the behavior of the audit logging system itself: at least one
   audit backend must successfully be loaded. The server log contains an error
   when this occurs. This helps keep a Vault HA cluster working when there is a
   misconfiguration on a standby node. [GH-2083]

FEATURES:

 * **Web UI (Enterprise)**: Vault Enterprise now contains a built-in web UI
   that offers access to a number of features, including init/unsealing/sealing,
   authentication via userpass or LDAP, and K/V reading/writing. The capability
   set of the UI will be expanding rapidly in further releases. To enable it,
   set `ui = true` in the top level of Vault's configuration file and point a
   web browser at your Vault address.
 * **Google Cloud Storage Physical Backend**: You can now use GCS for storing
   Vault data [GH-2099]

IMPROVEMENTS:

 * auth/github: Policies can now be assigned to users as well as to teams
   [GH-2079]
 * cli: Set the number of retries on 500 down to 0 by default (no retrying). It
   can be very confusing to users when there is a pause while the retries
   happen if they haven't explicitly set it. With request forwarding the need
   for this is lessened anyways. [GH-2093]
 * core: Response wrapping is now allowed to be specified by backend responses
   (requires backends gaining support) [GH-2088]
 * physical/consul: When announcing service, use the scheme of the Vault server
   rather than the Consul client [GH-2146]
 * secret/consul: Added listing functionality to roles [GH-2065]
 * secret/postgresql: Added `revocation_sql` parameter on the role endpoint to
   enable customization of user revocation SQL statements [GH-2033]
 * secret/transit: Add listing of keys [GH-1987]

BUG FIXES:

 * api/unwrap, command/unwrap: Increase compatibility of `unwrap` command with
   Vault 0.6.1 and older [GH-2014]
 * api/unwrap, command/unwrap: Fix error when no client token exists [GH-2077]
 * auth/approle: Creating the index for the role_id properly [GH-2004]
 * auth/aws-ec2: Handle the case of multiple upgrade attempts when setting the
   instance-profile ARN [GH-2035]
 * auth/ldap: Avoid leaking connections on login [GH-2130]
 * command/path-help: Use the actual error generated by Vault rather than
   always using 500 when there is a path help error [GH-2153]
 * command/ssh: Use temporary file for identity and ensure its deletion before
   the command returns [GH-2016]
 * cli: Fix error printing values with `-field` if the values contained
   formatting directives [GH-2109]
 * command/server: Don't say mlock is supported on OSX when it isn't. [GH-2120]
 * core: Fix bug where a failure to come up as active node (e.g. if an audit
   backend failed) could lead to deadlock [GH-2083]
 * physical/mysql: Fix potential crash during setup due to a query failure
   [GH-2105]
 * secret/consul: Fix panic on user error [GH-2145]

## 0.6.2 (October 5, 2016)

DEPRECATIONS/CHANGES:

 * Convergent Encryption v2: New keys in `transit` using convergent mode will
   use a new nonce derivation mechanism rather than require the user to supply
   a nonce. While not explicitly increasing security, it minimizes the
   likelihood that a user will use the mode improperly and impact the security
   of their keys. Keys in convergent mode that were created in v0.6.1 will
   continue to work with the same mechanism (user-supplied nonce).
 * `etcd` HA off by default: Following in the footsteps of `dynamodb`, the
   `etcd` storage backend now requires that `ha_enabled` be explicitly
   specified in the configuration file. The backend currently has known broken
   HA behavior, so this flag discourages use by default without explicitly
   enabling it. If you are using this functionality, when upgrading, you should
   set `ha_enabled` to `"true"` *before* starting the new versions of Vault.
 * Default/Max lease/token TTLs are now 32 days: In previous versions of Vault
   the default was 30 days, but moving it to 32 days allows some operations
   (e.g. reauthenticating, renewing, etc.) to be performed via a monthly cron
   job.
 * AppRole Secret ID endpoints changed: Secret ID and Secret ID accessors are
   no longer part of request URLs. The GET and DELETE operations are now moved
   to new endpoints (`/lookup` and `/destroy`) which consumes the input from
   the body and not the URL.
 * AppRole requires at least one constraint: previously it was sufficient to
   turn off all AppRole authentication constraints (secret ID, CIDR block) and
   use the role ID only. It is now required that at least one additional
   constraint is enabled. Existing roles are unaffected, but any new roles or
   updated roles will require this.
 * Reading wrapped responses from `cubbyhole/response` is deprecated. The
   `sys/wrapping/unwrap` endpoint should be used instead as it provides
   additional security, auditing, and other benefits. The ability to read
   directly will be removed in a future release.
 * Request Forwarding is now on by default: in 0.6.1 this required toggling on,
   but is now enabled by default. This can be disabled via the
   `"disable_clustering"` parameter in Vault's
   [config](https://www.vaultproject.io/docs/config/index.html), or per-request
   with the `X-Vault-No-Request-Forwarding` header.
 * In prior versions a bug caused the `bound_iam_role_arn` value in the
   `aws-ec2` authentication backend to actually use the instance profile ARN.
   This has been corrected, but as a result there is a behavior change. To
   match using the instance profile ARN, a new parameter
   `bound_iam_instance_profile_arn` has been added. Existing roles will
   automatically transfer the value over to the correct parameter, but the next
   time the role is updated, the new meanings will take effect.

FEATURES:

 * **Secret ID CIDR Restrictions in `AppRole`**: Secret IDs generated under an
   approle can now specify a list of CIDR blocks from where the requests to
   generate secret IDs should originate from. If an approle already has CIDR
   restrictions specified, the CIDR restrictions on the secret ID should be a
   subset of those specified on the role [GH-1910]
 * **Initial Root Token PGP Encryption**: Similar to `generate-root`, the root
   token created at initialization time can now be PGP encrypted [GH-1883]
 * **Support Chained Intermediate CAs in `pki`**: The `pki` backend now allows,
   when a CA cert is being supplied as a signed root or intermediate, a trust
   chain of arbitrary length. The chain is returned as a parameter at
   certificate issue/sign time and is retrievable independently as well.
   [GH-1694]
 * **Response Wrapping Enhancements**: There are new endpoints to look up
   response wrapped token parameters; wrap arbitrary values; rotate wrapping
   tokens; and unwrap with enhanced validation. In addition, list operations
   can now be response-wrapped. [GH-1927]
 * **Transit Features**: The `transit` backend now supports generating random
   bytes and SHA sums; HMACs; and signing and verification functionality using
   EC keys (P-256 curve)

IMPROVEMENTS:

 * api: Return error when an invalid (as opposed to incorrect) unseal key is
   submitted, rather than ignoring it [GH-1782]
 * api: Add method to call `auth/token/create-orphan` endpoint [GH-1834]
 * api: Rekey operation now redirects from standbys to master [GH-1862]
 * audit/file: Sending a `SIGHUP` to Vault now causes Vault to close and
   re-open the log file, making it easier to rotate audit logs [GH-1953]
 * auth/aws-ec2: EC2 instances can get authenticated by presenting the identity
   document and its SHA256 RSA digest [GH-1961]
 * auth/aws-ec2: IAM bound parameters on the aws-ec2 backend will perform a
   prefix match instead of exact match [GH-1943]
 * auth/aws-ec2: Added a new constraint `bound_iam_instance_profile_arn` to
   refer to IAM instance profile ARN and fixed the earlier `bound_iam_role_arn`
   to refer to IAM role ARN instead of the instance profile ARN [GH-1913]
 * auth/aws-ec2: Backend generates the nonce by default and clients can
   explicitly disable reauthentication by setting empty nonce [GH-1889]
 * auth/token: Added warnings if tokens and accessors are used in URLs [GH-1806]
 * command/format: The `format` flag on select CLI commands takes `yml` as an
   alias for `yaml` [GH-1899]
 * core: Allow the size of the read cache to be set via the config file, and
   change the default value to 1MB (from 32KB) [GH-1784]
 * core: Allow single and two-character path parameters for most places
   [GH-1811]
 * core: Allow list operations to be response-wrapped [GH-1814]
 * core: Provide better protection against timing attacks in Shamir code
   [GH-1877]
 * core: Unmounting/disabling backends no longer returns an error if the mount
   didn't exist. This is line with elsewhere in Vault's API where `DELETE` is
   an idempotent operation. [GH-1903]
 * credential/approle: At least one constraint is required to be enabled while
   creating and updating a role [GH-1882]
 * secret/cassandra: Added consistency level for use with roles [GH-1931]
 * secret/mysql: SQL for revoking user can be configured on the role [GH-1914]
 * secret/transit: Use HKDF (RFC 5869) as the key derivation function for new
   keys [GH-1812]
 * secret/transit: Empty plaintext values are now allowed [GH-1874]

BUG FIXES:

 * audit: Fix panic being caused by some values logging as underlying Go types
   instead of formatted strings [GH-1912]
 * auth/approle: Fixed panic on deleting approle that doesn't exist [GH-1920]
 * auth/approle: Not letting secret IDs and secret ID accessors to get logged
   in plaintext in audit logs [GH-1947]
 * auth/aws-ec2: Allow authentication if the underlying host is in a bad state
   but the instance is running [GH-1884]
 * auth/token: Fixed metadata getting missed out from token lookup response by
   gracefully handling token entry upgrade [GH-1924]
 * cli: Don't error on newline in token file [GH-1774]
 * core: Pass back content-type header for forwarded requests [GH-1791]
 * core: Fix panic if the same key was given twice to `generate-root` [GH-1827]
 * core: Fix potential deadlock on unmount/remount [GH-1793]
 * physical/file: Remove empty directories from the `file` storage backend [GH-1821]
 * physical/zookeeper: Remove empty directories from the `zookeeper` storage
   backend and add a fix to the `file` storage backend's logic [GH-1964]
 * secret/aws: Added update operation to `aws/sts` path to consider `ttl`
   parameter [39b75c6]
 * secret/aws: Mark STS secrets as non-renewable [GH-1804]
 * secret/cassandra: Properly store session for re-use [GH-1802]
 * secret/ssh: Fix panic when revoking SSH dynamic keys [GH-1781]

## 0.6.1 (August 22, 2016)

DEPRECATIONS/CHANGES:

 * Once the active node is 0.6.1, standby nodes must also be 0.6.1 in order to
   connect to the HA cluster. We recommend following our [general upgrade
   instructions](https://www.vaultproject.io/docs/install/upgrade.html) in
   addition to 0.6.1-specific upgrade instructions to ensure that this is not
   an issue.
 * Status codes for sealed/uninitialized Vaults have changed to `503`/`501`
   respectively. See the [version-specific upgrade
   guide](https://www.vaultproject.io/docs/install/upgrade-to-0.6.1.html) for
   more details.
 * Root tokens (tokens with the `root` policy) can no longer be created except
   by another root token or the `generate-root` endpoint.
 * Issued certificates from the `pki` backend against new roles created or
   modified after upgrading will contain a set of default key usages.
 * The `dynamodb` physical data store no longer supports HA by default. It has
   some non-ideal behavior around failover that was causing confusion. See the
   [documentation](https://www.vaultproject.io/docs/config/index.html#ha_enabled)
   for information on enabling HA mode. It is very important that this
   configuration is added _before upgrading_.
 * The `ldap` backend no longer searches for `memberOf` groups as part of its
   normal flow. Instead, the desired group filter must be specified. This fixes
   some errors and increases speed for directories with different structures,
   but if this behavior has been relied upon, ensure that you see the upgrade
   notes _before upgrading_.
 * `app-id` is now deprecated with the addition of the new AppRole backend.
   There are no plans to remove it, but we encourage using AppRole whenever
   possible, as it offers enhanced functionality and can accommodate many more
   types of authentication paradigms.

FEATURES:

 * **AppRole Authentication Backend**: The `approle` backend is a
   machine-oriented authentication backend that provides a similar concept to
   App-ID while adding many missing features, including a pull model that
   allows for the backend to generate authentication credentials rather than
   requiring operators or other systems to push credentials in. It should be
   useful in many more situations than App-ID. The inclusion of this backend
   deprecates App-ID. [GH-1426]
 * **Request Forwarding**: Vault servers can now forward requests to each other
   rather than redirecting clients. This feature is off by default in 0.6.1 but
   will be on by default in the next release. See the [HA concepts
   page](https://www.vaultproject.io/docs/concepts/ha.html) for information on
   enabling and configuring it. [GH-443]
 * **Convergent Encryption in `Transit`**: The `transit` backend now supports a
   convergent encryption mode where the same plaintext will produce the same
   ciphertext. Although very useful in some situations, this has potential
   security implications, which are mostly mitigated by requiring the use of
   key derivation when convergent encryption is enabled. See [the `transit`
   backend
   documentation](https://www.vaultproject.io/docs/secrets/transit/index.html)
   for more details. [GH-1537]
 * **Improved LDAP Group Filters**: The `ldap` auth backend now uses templates
   to define group filters, providing the capability to support some
   directories that could not easily be supported before (especially specific
   Active Directory setups with nested groups). [GH-1388]
 * **Key Usage Control in `PKI`**: Issued certificates from roles created or
   modified after upgrading contain a set of default key usages for increased
   compatibility with OpenVPN and some other software. This set can be changed
   when writing a role definition. Existing roles are unaffected. [GH-1552]
 * **Request Retrying in the CLI and Go API**: Requests that fail with a `5xx`
   error code will now retry after a backoff. The maximum total number of
   retries (including disabling this functionality) can be set with an
   environment variable. See the [environment variable
   documentation](https://www.vaultproject.io/docs/commands/environment.html)
   for more details. [GH-1594]
 * **Service Discovery in `vault init`**: The new `-auto` option on `vault init`
   will perform service discovery using Consul. When only one node is discovered,
   it will be initialized and when more than one node is discovered, they will
   be output for easy selection. See `vault init --help` for more details. [GH-1642]
 * **MongoDB Secret Backend**: Generate dynamic unique MongoDB database
   credentials based on configured roles. Sponsored by
   [CommerceHub](http://www.commercehub.com/). [GH-1414]
 * **Circonus Metrics Integration**: Vault can now send metrics to
   [Circonus](http://www.circonus.com/). See the [configuration
   documentation](https://www.vaultproject.io/docs/config/index.html) for
   details. [GH-1646]

IMPROVEMENTS:

 * audit: Added a unique identifier to each request which will also be found in
   the request portion of the response. [GH-1650]
 * auth/aws-ec2: Added a new constraint `bound_account_id` to the role
   [GH-1523]
 * auth/aws-ec2: Added a new constraint `bound_iam_role_arn` to the role
   [GH-1522]
 * auth/aws-ec2: Added `ttl` field for the role [GH-1703]
 * auth/ldap, secret/cassandra, physical/consul: Clients with `tls.Config`
   have the minimum TLS version set to 1.2 by default. This is configurable.
 * auth/token: Added endpoint to list accessors [GH-1676]
 * auth/token: Added `disallowed_policies` option to token store roles [GH-1681]
 * auth/token: `root` or `sudo` tokens can now create periodic tokens via
   `auth/token/create`; additionally, the same token can now be periodic and
   have an explicit max TTL [GH-1725]
 * build: Add support for building on Solaris/Illumos [GH-1726]
 * cli: Output formatting in the presence of warnings in the response object
   [GH-1533]
 * cli: `vault auth` command supports a `-path` option to take in the path at
   which the auth backend is enabled, thereby allowing authenticating against
   different paths using the command options [GH-1532]
 * cli: `vault auth -methods` will now display the config settings of the mount
   [GH-1531]
 * cli: `vault read/write/unwrap -field` now allows selecting token response
   fields [GH-1567]
 * cli: `vault write -field` now allows selecting wrapped response fields
   [GH-1567]
 * command/status: Version information and cluster details added to the output
   of `vault status` command [GH-1671]
 * core: Response wrapping is now enabled for login endpoints [GH-1588]
 * core: The duration of leadership is now exported via events through
   telemetry [GH-1625]
 * core: `sys/capabilities-self` is now accessible as part of the `default`
   policy [GH-1695]
 * core: `sys/renew` is now accessible as part of the `default` policy [GH-1701]
 * core: Unseal keys will now be returned in both hex and base64 forms, and
   either can be used [GH-1734]
 * core: Responses from most `/sys` endpoints now return normal `api.Secret`
   structs in addition to the values they carried before. This means that
   response wrapping can now be used with most authenticated `/sys` operations
   [GH-1699]
 * physical/etcd: Support `ETCD_ADDR` env var for specifying addresses [GH-1576]
 * physical/consul: Allowing additional tags to be added to Consul service
   registration via `service_tags` option [GH-1643]
 * secret/aws: Listing of roles is supported now  [GH-1546]
 * secret/cassandra: Add `connect_timeout` value for Cassandra connection
   configuration [GH-1581]
 * secret/mssql,mysql,postgresql: Reading of connection settings is supported
   in all the sql backends [GH-1515]
 * secret/mysql: Added optional maximum idle connections value to MySQL
   connection configuration [GH-1635]
 * secret/mysql: Use a combination of the role name and token display name in
   generated user names and allow the length to be controlled [GH-1604]
 * secret/{cassandra,mssql,mysql,postgresql}: SQL statements can now be passed
   in via one of four ways: a semicolon-delimited string, a base64-delimited
   string, a serialized JSON string array, or a base64-encoded serialized JSON
   string array [GH-1686]
 * secret/ssh: Added `allowed_roles` to vault-ssh-helper's config and returning
   role name as part of response of `verify` API
 * secret/ssh: Added passthrough of command line arguments to `ssh` [GH-1680]
 * sys/health: Added version information to the response of health status
   endpoint [GH-1647]
 * sys/health: Cluster information isbe returned as part of health status when
   Vault is unsealed [GH-1671]
 * sys/mounts: MountTable data is compressed before serializing to accommodate
   thousands of mounts [GH-1693]
 * website: The [token
   concepts](https://www.vaultproject.io/docs/concepts/tokens.html) page has
   been completely rewritten [GH-1725]

BUG FIXES:

 * auth/aws-ec2: Added a nil check for stored whitelist identity object
   during renewal [GH-1542]
 * auth/cert: Fix panic if no client certificate is supplied [GH-1637]
 * auth/token: Don't report that a non-expiring root token is renewable, as
   attempting to renew it results in an error [GH-1692]
 * cli: Don't retry a command when a redirection is received [GH-1724]
 * core: Fix regression causing status codes to be `400` in most non-5xx error
   cases [GH-1553]
 * core: Fix panic that could occur during a leadership transition [GH-1627]
 * physical/postgres: Remove use of prepared statements as this causes
   connection multiplexing software to break [GH-1548]
 * physical/consul: Multiple Vault nodes on the same machine leading to check ID
   collisions were resulting in incorrect health check responses [GH-1628]
 * physical/consul: Fix deregistration of health checks on exit [GH-1678]
 * secret/postgresql: Check for existence of role before attempting deletion
   [GH-1575]
 * secret/postgresql: Handle revoking roles that have privileges on sequences
   [GH-1573]
 * secret/postgresql(,mysql,mssql): Fix incorrect use of database over
   transaction object which could lead to connection exhaustion [GH-1572]
 * secret/pki: Fix parsing CA bundle containing trailing whitespace [GH-1634]
 * secret/pki: Fix adding email addresses as SANs [GH-1688]
 * secret/pki: Ensure that CRL values are always UTC, per RFC [GH-1727]
 * sys/seal-status: Fixed nil Cluster object while checking seal status [GH-1715]

## 0.6.0 (June 14th, 2016)

SECURITY:

 * Although `sys/revoke-prefix` was intended to revoke prefixes of secrets (via
   lease IDs, which incorporate path information) and
   `auth/token/revoke-prefix` was intended to revoke prefixes of tokens (using
   the tokens' paths and, since 0.5.2, role information), in implementation
   they both behaved exactly the same way since a single component in Vault is
   responsible for managing lifetimes of both, and the type of the tracked
   lifetime was not being checked. The end result was that either endpoint
   could revoke both secret leases and tokens. We consider this a very minor
   security issue as there are a number of mitigating factors: both endpoints
   require `sudo` capability in addition to write capability, preventing
   blanket ACL path globs from providing access; both work by using the prefix
   to revoke as a part of the endpoint path, allowing them to be properly
   ACL'd; and both are intended for emergency scenarios and users should
   already not generally have access to either one. In order to prevent
   confusion, we have simply removed `auth/token/revoke-prefix` in 0.6, and
   `sys/revoke-prefix` will be meant for both leases and tokens instead.

DEPRECATIONS/CHANGES:

 * `auth/token/revoke-prefix` has been removed. See the security notice for
   details. [GH-1280]
 * Vault will now automatically register itself as the `vault` service when
   using the `consul` backend and will perform its own health checks.  See
   the Consul backend documentation for information on how to disable
   auto-registration and service checks.
 * List operations that do not find any keys now return a `404` status code
   rather than an empty response object [GH-1365]
 * CA certificates issued from the `pki` backend no longer have associated
   leases, and any CA certs already issued will ignore revocation requests from
   the lease manager. This is to prevent CA certificates from being revoked
   when the token used to issue the certificate expires; it was not be obvious
   to users that they need to ensure that the token lifetime needed to be at
   least as long as a potentially very long-lived CA cert.

FEATURES:

 * **AWS EC2 Auth Backend**: Provides a secure introduction mechanism for AWS
   EC2 instances allowing automated retrieval of Vault tokens. Unlike most
   Vault authentication backends, this backend does not require first deploying
   or provisioning security-sensitive credentials (tokens, username/password,
   client certificates, etc). Instead, it treats AWS as a Trusted Third Party
   and uses the cryptographically signed dynamic metadata information that
   uniquely represents each EC2 instance. [Vault
   Enterprise](https://www.hashicorp.com/vault.html) customers have access to a
   turnkey client that speaks the backend API and makes access to a Vault token
   easy.
 * **Response Wrapping**: Nearly any response within Vault can now be wrapped
   inside a single-use, time-limited token's cubbyhole, taking the [Cubbyhole
   Authentication
   Principles](https://www.hashicorp.com/blog/vault-cubbyhole-principles.html)
   mechanism to its logical conclusion. Retrieving the original response is as
   simple as a single API command or the new `vault unwrap` command. This makes
   secret distribution easier and more secure, including secure introduction.
 * **Azure Physical Backend**: You can now use Azure blob object storage as
   your Vault physical data store [GH-1266]
 * **Swift Physical Backend**: You can now use Swift blob object storage as
   your Vault physical data store [GH-1425]
 * **Consul Backend Health Checks**: The Consul backend will automatically
   register a `vault` service and perform its own health checking. By default
   the active node can be found at `active.vault.service.consul` and all with
   standby nodes are `standby.vault.service.consul`. Sealed vaults are marked
   critical and are not listed by default in Consul's service discovery.  See
   the documentation for details. [GH-1349]
 * **Explicit Maximum Token TTLs**: You can now set explicit maximum TTLs on
   tokens that do not honor changes in the system- or mount-set values. This is
   useful, for instance, when the max TTL of the system or the `auth/token`
   mount must be set high to accommodate certain needs but you want more
   granular restrictions on tokens being issued directly from the Token
   authentication backend at `auth/token`. [GH-1399]
 * **Non-Renewable Tokens**: When creating tokens directly through the token
   authentication backend, you can now specify in both token store roles and
   the API whether or not a token should be renewable, defaulting to `true`.
 * **RabbitMQ Secret Backend**: Vault can now generate credentials for
   RabbitMQ. Vhosts and tags can be defined within roles. [GH-788]

IMPROVEMENTS:

 * audit: Add the DisplayName value to the copy of the Request object embedded
   in the associated Response, to match the original Request object [GH-1387]
 * audit: Enable auditing of the `seal` and `step-down` commands [GH-1435]
 * backends: Remove most `root`/`sudo` paths in favor of normal ACL mechanisms.
   A particular exception are any current MFA paths. A few paths in `token` and
   `sys` also require `root` or `sudo`. [GH-1478]
 * command/auth: Restore the previous authenticated token if the `auth` command
   fails to authenticate the provided token [GH-1233]
 * command/write: `-format` and `-field` can now be used with the `write`
   command [GH-1228]
 * core: Add `mlock` support for FreeBSD, OpenBSD, and Darwin [GH-1297]
 * core: Don't keep lease timers around when tokens are revoked [GH-1277]
 * core: If using the `disable_cache` option, caches for the policy store and
   the `transit` backend are now disabled as well [GH-1346]
 * credential/cert: Renewal requests are rejected if the set of policies has
   changed since the token was issued [GH-477]
 * credential/cert: Check CRLs for specific non-CA certs configured in the
   backend [GH-1404]
 * credential/ldap: If `groupdn` is not configured, skip searching LDAP and
   only return policies for local groups, plus a warning [GH-1283]
 * credential/ldap: `vault list` support for users and groups [GH-1270]
 * credential/ldap: Support for the `memberOf` attribute for group membership
   searching [GH-1245]
 * credential/userpass: Add list support for users [GH-911]
 * credential/userpass: Remove user configuration paths from requiring sudo, in
   favor of normal ACL mechanisms [GH-1312]
 * credential/token: Sanitize policies and add `default` policies in appropriate
   places [GH-1235]
 * credential/token: Setting the renewable status of a token is now possible
   via `vault token-create` and the API. The default is true, but tokens can be
   specified as non-renewable. [GH-1499]
 * secret/aws: Use chain credentials to allow environment/EC2 instance/shared
   providers [GH-307]
 * secret/aws: Support for STS AssumeRole functionality [GH-1318]
 * secret/consul: Reading consul access configuration supported. The response
   will contain non-sensitive information only [GH-1445]
 * secret/pki: Added `exclude_cn_from_sans` field to prevent adding the CN to
   DNS or Email Subject Alternate Names [GH-1220]
 * secret/pki: Added list support for certificates [GH-1466]
 * sys/capabilities: Enforce ACL checks for requests that query the capabilities
   of a token on a given path [GH-1221]
 * sys/health: Status information can now be retrieved with `HEAD` [GH-1509]

BUG FIXES:

 * command/read: Fix panic when using `-field` with a non-string value [GH-1308]
 * command/token-lookup: Fix TTL showing as 0 depending on how a token was
   created. This only affected the value shown at lookup, not the token
   behavior itself. [GH-1306]
 * command/various: Tell the JSON decoder to not convert all numbers to floats;
   fixes some various places where numbers were showing up in scientific
   notation
 * command/server: Prioritized `devRootTokenID` and `devListenAddress` flags
   over their respective env vars [GH-1480]
 * command/ssh: Provided option to disable host key checking. The automated
   variant of `vault ssh` command uses `sshpass` which was failing to handle
   host key checking presented by the `ssh` binary. [GH-1473]
 * core: Properly persist mount-tuned TTLs for auth backends [GH-1371]
 * core: Don't accidentally crosswire SIGINT to the reload handler [GH-1372]
 * credential/github: Make organization comparison case-insensitive during
   login [GH-1359]
 * credential/github: Fix panic when renewing a token created with some earlier
   versions of Vault [GH-1510]
 * credential/github: The token used to log in via `vault auth` can now be
   specified in the `VAULT_AUTH_GITHUB_TOKEN` environment variable [GH-1511]
 * credential/ldap: Fix problem where certain error conditions when configuring
   or opening LDAP connections would cause a panic instead of return a useful
   error message [GH-1262]
 * credential/token: Fall back to normal parent-token semantics if
   `allowed_policies` is empty for a role. Using `allowed_policies` of
   `default` resulted in the same behavior anyways. [GH-1276]
 * credential/token: Fix issues renewing tokens when using the "suffix"
   capability of token roles [GH-1331]
 * credential/token: Fix lookup via POST showing the request token instead of
   the desired token [GH-1354]
 * credential/various: Fix renewal conditions when `default` policy is not
   contained in the backend config [GH-1256]
 * physical/s3: Don't panic in certain error cases from bad S3 responses [GH-1353]
 * secret/consul: Use non-pooled Consul API client to avoid leaving files open
   [GH-1428]
 * secret/pki: Don't check whether a certificate is destined to be a CA
   certificate if sign-verbatim endpoint is used [GH-1250]

## 0.5.3 (May 27th, 2016)

SECURITY:

 * Consul ACL Token Revocation: An issue was reported to us indicating that
   generated Consul ACL tokens were not being properly revoked. Upon
   investigation, we found that this behavior was reproducible in a specific
   scenario: when a generated lease for a Consul ACL token had been renewed
   prior to revocation. In this case, the generated token was not being
   properly persisted internally through the renewal function, leading to an
   error during revocation due to the missing token. Unfortunately, this was
   coded as a user error rather than an internal error, and the revocation
   logic was expecting internal errors if revocation failed. As a result, the
   revocation logic believed the revocation to have succeeded when it in fact
   failed, causing the lease to be dropped while the token was still valid
   within Consul. In this release, the Consul backend properly persists the
   token through renewals, and the revocation logic has been changed to
   consider any error type to have been a failure to revoke, causing the lease
   to persist and attempt to be revoked later.

We have written an example shell script that searches through Consul's ACL
tokens and looks for those generated by Vault, which can be used as a template
for a revocation script as deemed necessary for any particular security
response. The script is available at
https://gist.github.com/jefferai/6233c2963f9407a858d84f9c27d725c0

Please note that any outstanding leases for Consul tokens produced prior to
0.5.3 that have been renewed will continue to exhibit this behavior. As a
result, we recommend either revoking all tokens produced by the backend and
issuing new ones, or if needed, a more advanced variant of the provided example
could use the timestamp embedded in each generated token's name to decide which
tokens are too old and should be deleted. This could then be run periodically
up until the maximum lease time for any outstanding pre-0.5.3 tokens has
expired.

This is a security-only release. There are no other code changes since 0.5.2.
The binaries have one additional change: they are built against Go 1.6.1 rather
than Go 1.6, as Go 1.6.1 contains two security fixes to the Go programming
language itself.

## 0.5.2 (March 16th, 2016)

FEATURES:

 * **MSSQL Backend**: Generate dynamic unique MSSQL database credentials based
   on configured roles [GH-998]
 * **Token Accessors**: Vault now provides an accessor with each issued token.
   This accessor is an identifier that can be used for a limited set of
   actions, notably for token revocation. This value can be logged in
   plaintext to audit logs, and in combination with the plaintext metadata
   logged to audit logs, provides a searchable and straightforward way to
   revoke particular users' or services' tokens in many cases. To enable
   plaintext audit logging of these accessors, set `hmac_accessor=false` when
   enabling an audit backend.
 * **Token Credential Backend Roles**: Roles can now be created in the `token`
   credential backend that allow modifying token behavior in ways that are not
   otherwise exposed or easily delegated. This allows creating tokens with a
   fixed set (or subset) of policies (rather than a subset of the calling
   token's), periodic tokens with a fixed TTL but no expiration, specified
   prefixes, and orphans.
 * **Listener Certificate Reloading**: Vault's configured listeners now reload
   their TLS certificate and private key when the Vault process receives a
   SIGHUP.

IMPROVEMENTS:

 * auth/token: Endpoints optionally accept tokens from the HTTP body rather
   than just from the URLs [GH-1211]
 * auth/token,sys/capabilities: Added new endpoints
   `auth/token/lookup-accessor`, `auth/token/revoke-accessor` and
   `sys/capabilities-accessor`, which enables performing the respective actions
   with just the accessor of the tokens, without having access to the actual
   token [GH-1188]
 * core: Ignore leading `/` in policy paths [GH-1170]
 * core: Ignore leading `/` in mount paths [GH-1172]
 * command/policy-write: Provided HCL is now validated for format violations
   and provides helpful information around where the violation occurred
   [GH-1200]
 * command/server: The initial root token ID when running in `-dev` mode can
   now be specified via `-dev-root-token-id` or the environment variable
   `VAULT_DEV_ROOT_TOKEN_ID` [GH-1162]
 * command/server: The listen address when running in `-dev` mode can now be
   specified via `-dev-listen-address` or the environment variable
   `VAULT_DEV_LISTEN_ADDRESS` [GH-1169]
 * command/server: The configured listeners now reload their TLS
   certificates/keys when Vault is SIGHUP'd [GH-1196]
 * command/step-down: New `vault step-down` command and API endpoint to force
   the targeted node to give up active status, but without sealing. The node
   will wait ten seconds before attempting to grab the lock again. [GH-1146]
 * command/token-renew: Allow no token to be passed in; use `renew-self` in
   this case. Change the behavior for any token being passed in to use `renew`.
   [GH-1150]
 * credential/app-id: Allow `app-id` parameter to be given in the login path;
   this causes the `app-id` to be part of the token path, making it easier to
   use with `revoke-prefix` [GH-424]
 * credential/cert: Non-CA certificates can be used for authentication. They
   must be matched exactly (issuer and serial number) for authentication, and
   the certificate must carry the client authentication or 'any' extended usage
   attributes. [GH-1153]
 * credential/cert: Subject and Authority key IDs are output in metadata; this
   allows more flexible searching/revocation in the audit logs [GH-1183]
 * credential/cert: Support listing configured certs [GH-1212]
 * credential/userpass: Add support for `create`/`update` capability
   distinction in user path, and add user-specific endpoints to allow changing
   the password and policies [GH-1216]
 * credential/token: Add roles [GH-1155]
 * secret/mssql: Add MSSQL backend [GH-998]
 * secret/pki: Add revocation time (zero or Unix epoch) to `pki/cert/SERIAL`
   endpoint [GH-1180]
 * secret/pki: Sanitize serial number in `pki/revoke` endpoint to allow some
   other formats [GH-1187]
 * secret/ssh: Added documentation for `ssh/config/zeroaddress` endpoint.
   [GH-1154]
 * sys: Added new endpoints `sys/capabilities` and `sys/capabilities-self` to
   fetch the capabilities of a token on a given path [GH-1171]
 * sys: Added `sys/revoke-force`, which enables a user to ignore backend errors
   when revoking a lease, necessary in some emergency/failure scenarios
   [GH-1168]
 * sys: The return codes from `sys/health` can now be user-specified via query
   parameters [GH-1199]

BUG FIXES:

 * logical/cassandra: Apply hyphen/underscore replacement to the entire
   generated username, not just the UUID, in order to handle token display name
   hyphens [GH-1140]
 * physical/etcd: Output actual error when cluster sync fails [GH-1141]
 * vault/expiration: Not letting the error responses from the backends to skip
   during renewals [GH-1176]

## 0.5.1 (February 25th, 2016)

DEPRECATIONS/CHANGES:

 * RSA keys less than 2048 bits are no longer supported in the PKI backend.
   1024-bit keys are considered unsafe and are disallowed in the Internet PKI.
   The `pki` backend has enforced SHA256 hashes in signatures from the
   beginning, and software that can handle these hashes should be able to
   handle larger key sizes. [GH-1095]
 * The PKI backend now does not automatically delete expired certificates,
   including from the CRL. Doing so could lead to a situation where a time
   mismatch between the Vault server and clients could result in a certificate
   that would not be considered expired by a client being removed from the CRL.
   The new `pki/tidy` endpoint can be used to trigger expirations. [GH-1129]
 * The `cert` backend now performs a variant of channel binding at renewal time
   for increased security. In order to not overly burden clients, a notion of
   identity is used. This functionality can be disabled. See the 0.5.1 upgrade
   guide for more specific information [GH-1127]

FEATURES:

 * **Codebase Audit**: Vault's 0.5 codebase was audited by iSEC. (The terms of
   the audit contract do not allow us to make the results public.) [GH-220]

IMPROVEMENTS:

 * api: The `VAULT_TLS_SERVER_NAME` environment variable can be used to control
   the SNI header during TLS connections [GH-1131]
 * api/health: Add the server's time in UTC to health responses [GH-1117]
 * command/rekey and command/generate-root: These now return the status at
   attempt initialization time, rather than requiring a separate fetch for the
   nonce [GH-1054]
 * credential/cert: Don't require root/sudo tokens for the `certs/` and `crls/`
   paths; use normal ACL behavior instead [GH-468]
 * credential/github: The validity of the token used for login will be checked
   at renewal time [GH-1047]
 * credential/github: The `config` endpoint no longer requires a root token;
   normal ACL path matching applies
 * deps: Use the standardized Go 1.6 vendoring system
 * secret/aws: Inform users of AWS-imposed policy restrictions around STS
   tokens if they attempt to use an invalid policy [GH-1113]
 * secret/mysql: The MySQL backend now allows disabling verification of the
   `connection_url` [GH-1096]
 * secret/pki: Submitted CSRs are now verified to have the correct key type and
   minimum number of bits according to the role. The exception is intermediate
   CA signing and the `sign-verbatim` path [GH-1104]
 * secret/pki: New `tidy` endpoint to allow expunging expired certificates.
   [GH-1129]
 * secret/postgresql: The PostgreSQL backend now allows disabling verification
   of the `connection_url` [GH-1096]
 * secret/ssh: When verifying an OTP, return 400 if it is not valid instead of
   204 [GH-1086]
 * credential/app-id: App ID backend will check the validity of app-id and user-id
   during renewal time [GH-1039]
 * credential/cert: TLS Certificates backend, during renewal, will now match the
   client identity with the client identity used during login [GH-1127]

BUG FIXES:

 * credential/ldap: Properly escape values being provided to search filters
   [GH-1100]
 * secret/aws: Capping on length of usernames for both IAM and STS types
   [GH-1102]
 * secret/pki: If a cert is not found during lookup of a serial number,
   respond with a 400 rather than a 500 [GH-1085]
 * secret/postgresql: Add extra revocation statements to better handle more
   permission scenarios [GH-1053]
 * secret/postgresql: Make connection_url work properly [GH-1112]

## 0.5.0 (February 10, 2016)

SECURITY:

 * Previous versions of Vault could allow a malicious user to hijack the rekey
   operation by canceling an operation in progress and starting a new one. The
   practical application of this is very small. If the user was an unseal key
   owner, they could attempt to do this in order to either receive unencrypted
   reseal keys or to replace the PGP keys used for encryption with ones under
   their control. However, since this would invalidate any rekey progress, they
   would need other unseal key holders to resubmit, which would be rather
   suspicious during this manual operation if they were not also the original
   initiator of the rekey attempt. If the user was not an unseal key holder,
   there is no benefit to be gained; the only outcome that could be attempted
   would be a denial of service against a legitimate rekey operation by sending
   cancel requests over and over. Thanks to Josh Snyder for the report!

DEPRECATIONS/CHANGES:

 * `s3` physical backend: Environment variables are now preferred over
   configuration values. This makes it behave similar to the rest of Vault,
   which, in increasing order of preference, uses values from the configuration
   file, environment variables, and CLI flags. [GH-871]
 * `etcd` physical backend: `sync` functionality is now supported and turned on
   by default. This can be disabled. [GH-921]
 * `transit`: If a client attempts to encrypt a value with a key that does not
   yet exist, what happens now depends on the capabilities set in the client's
   ACL policies. If the client has `create` (or `create` and `update`)
   capability, the key will upsert as in the past. If the client has `update`
   capability, they will receive an error. [GH-1012]
 * `token-renew` CLI command: If the token given for renewal is the same as the
   client token, the `renew-self` endpoint will be used in the API. Given that
   the `default` policy (by default) allows all clients access to the
   `renew-self` endpoint, this makes it much more likely that the intended
   operation will be successful. [GH-894]
 * Token `lookup`: the `ttl` value in the response now reflects the actual
   remaining TTL rather than the original TTL specified when the token was
   created; this value is now located in `creation_ttl` [GH-986]
 * Vault no longer uses grace periods on leases or token TTLs. Uncertainty
   about the length grace period for any given backend could cause confusion
   and uncertainty. [GH-1002]
 * `rekey`: Rekey now requires a nonce to be supplied with key shares. This
   nonce is generated at the start of a rekey attempt and is unique for that
   attempt.
 * `status`: The exit code for the `status` CLI command is now `2` for an
   uninitialized Vault instead of `1`. `1` is returned for errors. This better
   matches the rest of the CLI.

FEATURES:

 * **Split Data/High Availability Physical Backends**: You can now configure
   two separate physical backends: one to be used for High Availability
   coordination and another to be used for encrypted data storage. See the
   [configuration
   documentation](https://vaultproject.io/docs/config/index.html) for details.
   [GH-395]
 * **Fine-Grained Access Control**: Policies can now use the `capabilities` set
   to specify fine-grained control over operations allowed on a path, including
   separation of `sudo` privileges from other privileges. These can be mixed
   and matched in any way desired. The `policy` value is kept for backwards
   compatibility. See the [updated policy
   documentation](https://vaultproject.io/docs/concepts/policies.html) for
   details. [GH-914]
 * **List Support**: Listing is now supported via the API and the new `vault
   list` command. This currently supports listing keys in the `generic` and
   `cubbyhole` backends and a few other places (noted in the IMPROVEMENTS
   section below). Different parts of the API and backends will need to
   implement list capabilities in ways that make sense to particular endpoints,
   so further support will appear over time. [GH-617]
 * **Root Token Generation via Unseal Keys**: You can now use the
   `generate-root` CLI command to generate new orphaned, non-expiring root
   tokens in case the original is lost or revoked (accidentally or
   purposefully). This requires a quorum of unseal key holders. The output
   value is protected via any PGP key of the initiator's choosing or a one-time
   pad known only to the initiator (a suitable pad can be generated via the
   `-genotp` flag to the command. [GH-915]
 * **Unseal Key Archiving**: You can now optionally have Vault store your
   unseal keys in your chosen physical store for disaster recovery purposes.
   This option is only available when the keys are encrypted with PGP. [GH-907]
 * **Keybase Support for PGP Encryption Keys**: You can now specify Keybase
   users when passing in PGP keys to the `init`, `rekey`, and `generate-root`
   CLI commands.  Public keys for these users will be fetched automatically.
   [GH-901]
 * **DynamoDB HA Physical Backend**: There is now a new, community-supported
   HA-enabled physical backend using Amazon DynamoDB. See the [configuration
   documentation](https://vaultproject.io/docs/config/index.html) for details.
   [GH-878]
 * **PostgreSQL Physical Backend**: There is now a new, community-supported
   physical backend using PostgreSQL. See the [configuration
   documentation](https://vaultproject.io/docs/config/index.html) for details.
   [GH-945]
 * **STS Support in AWS Secret Backend**: You can now use the AWS secret
   backend to fetch STS tokens rather than IAM users. [GH-927]
 * **Speedups in the transit backend**: The `transit` backend has gained a
   cache, and now loads only the working set of keys (e.g. from the
   `min_decryption_version` to the current key version) into its working set.
   This provides large speedups and potential memory savings when the `rotate`
   feature of the backend is used heavily.

IMPROVEMENTS:

 * cli: Output secrets sorted by key name [GH-830]
 * cli: Support YAML as an output format [GH-832]
 * cli: Show an error if the output format is incorrect, rather than falling
   back to an empty table [GH-849]
 * cli: Allow setting the `advertise_addr` for HA via the
   `VAULT_ADVERTISE_ADDR` environment variable [GH-581]
 * cli/generate-root: Add generate-root and associated functionality [GH-915]
 * cli/init: Add `-check` flag that returns whether Vault is initialized
   [GH-949]
 * cli/server: Use internal functions for the token-helper rather than shelling
   out, which fixes some problems with using a static binary in Docker or paths
   with multiple spaces when launching in `-dev` mode [GH-850]
 * cli/token-lookup: Add token-lookup command [GH-892]
 * command/{init,rekey}: Allow ASCII-armored keychain files to be arguments for
   `-pgp-keys` [GH-940]
 * conf: Use normal bool values rather than empty/non-empty for the
   `tls_disable` option [GH-802]
 * credential/ldap: Add support for binding, both anonymously (to discover a
   user DN) and via a username and password [GH-975]
 * credential/token: Add `last_renewal_time` to token lookup calls [GH-896]
 * credential/token: Change `ttl` to reflect the current remaining TTL; the
   original value is in `creation_ttl` [GH-1007]
 * helper/certutil: Add ability to parse PKCS#8 bundles [GH-829]
 * logical/aws: You can now get STS tokens instead of IAM users [GH-927]
 * logical/cassandra: Add `protocol_version` parameter to set the CQL proto
   version [GH-1005]
 * logical/cubbyhole: Add cubbyhole access to default policy [GH-936]
 * logical/mysql: Add list support for roles path [GH-984]
 * logical/pki: Fix up key usages being specified for CAs [GH-989]
 * logical/pki: Add list support for roles path [GH-985]
 * logical/pki: Allow `pem_bundle` to be specified as the format, which
   provides a concatenated PEM bundle of returned values [GH-1008]
 * logical/pki: Add 30 seconds of slack to the validity start period to
   accommodate some clock skew in machines [GH-1036]
 * logical/postgres: Add `max_idle_connections` parameter [GH-950]
 * logical/postgres: Add list support for roles path
 * logical/ssh: Add list support for roles path [GH-983]
 * logical/transit: Keys are archived and only keys between the latest version
   and `min_decryption_version` are loaded into the working set. This can
   provide a very large speed increase when rotating keys very often. [GH-977]
 * logical/transit: Keys are now cached, which should provide a large speedup
   in most cases [GH-979]
 * physical/cache: Use 2Q cache instead of straight LRU [GH-908]
 * physical/etcd: Support basic auth [GH-859]
 * physical/etcd: Support sync functionality and enable by default [GH-921]

BUG FIXES:

 * api: Correct the HTTP verb used in the LookupSelf method [GH-887]
 * api: Fix the output of `Sys().MountConfig(...)` to return proper values
   [GH-1017]
 * command/read: Fix panic when an empty argument was given [GH-923]
 * command/ssh: Fix panic when username lookup fails [GH-886]
 * core: When running in standalone mode, don't advertise that we are active
   until post-unseal setup completes [GH-872]
 * core: Update go-cleanhttp dependency to ensure idle connections aren't
   leaked [GH-867]
 * core: Don't allow tokens to have duplicate policies [GH-897]
 * core: Fix regression in `sys/renew` that caused information stored in the
   Secret part of the response to be lost [GH-912]
 * physical: Use square brackets when setting an IPv6-based advertise address
   as the auto-detected advertise address [GH-883]
 * physical/s3: Use an initialized client when using IAM roles to fix a
   regression introduced against newer versions of the AWS Go SDK [GH-836]
 * secret/pki: Fix a condition where unmounting could fail if the CA
   certificate was not properly loaded [GH-946]
 * secret/ssh: Fix a problem where SSH connections were not always closed
   properly [GH-942]

MISC:

 * Clarified our stance on support for community-derived physical backends.
   See the [configuration
   documentation](https://vaultproject.io/docs/config/index.html) for details.
 * Add `vault-java` to libraries [GH-851]
 * Various minor documentation fixes and improvements [GH-839] [GH-854]
   [GH-861] [GH-876] [GH-899] [GH-900] [GH-904] [GH-923] [GH-924] [GH-958]
   [GH-959] [GH-981] [GH-990] [GH-1024] [GH-1025]

BUILD NOTE:

 * The HashiCorp-provided binary release of Vault 0.5.0 is built against a
   patched version of Go 1.5.3 containing two specific bug fixes affecting TLS
   certificate handling. These fixes are in the Go 1.6 tree and were
   cherry-picked on top of stock Go 1.5.3. If you want to examine the way in
   which the releases were built, please look at our [cross-compilation
   Dockerfile](https://github.com/hashicorp/vault/blob/v0.5.0/scripts/cross/Dockerfile-patched-1.5.3).

## 0.4.1 (January 13, 2016)

SECURITY:

  * Build against Go 1.5.3 to mitigate a security vulnerability introduced in
    Go 1.5. For more information, please see
    https://groups.google.com/forum/#!topic/golang-dev/MEATuOi_ei4

This is a security-only release; other than the version number and building
against Go 1.5.3, there are no changes from 0.4.0.

## 0.4.0 (December 10, 2015)

DEPRECATIONS/CHANGES:

 * Policy Name Casing: Policy names are now normalized to lower-case on write,
   helping prevent accidental case mismatches. For backwards compatibility,
   policy names are not currently normalized when reading or deleting. [GH-676]
 * Default etcd port number: the default connection string for the `etcd`
   physical store uses port 2379 instead of port 4001, which is the port used
   by the supported version 2.x of etcd. [GH-753]
 * As noted below in the FEATURES section, if your Vault installation contains
   a policy called `default`, new tokens created will inherit this policy
   automatically.
 * In the PKI backend there have been a few minor breaking changes:
   * The token display name is no longer a valid option for providing a base
   domain for issuance. Since this name is prepended with the name of the
   authentication backend that issued it, it provided a faulty use-case at best
   and a confusing experience at worst. We hope to figure out a better
   per-token value in a future release.
   * The `allowed_base_domain` parameter has been changed to `allowed_domains`,
   which accepts a comma-separated list of domains. This allows issuing
   certificates with DNS subjects across multiple domains. If you had a
   configured `allowed_base_domain` parameter, it will be migrated
   automatically when the role is read (either via a normal read, or via
   issuing a certificate).

FEATURES:

 * **Significantly Enhanced PKI Backend**: The `pki` backend can now generate
   and sign root CA certificates and intermediate CA CSRs. It can also now sign
   submitted client CSRs, as well as a significant number of other
   enhancements. See the updated documentation for the full API. [GH-666]
 * **CRL Checking for Certificate Authentication**: The `cert` backend now
   supports pushing CRLs into the mount and using the contained serial numbers
   for revocation checking. See the documentation for the `cert` backend for
   more info. [GH-330]
 * **Default Policy**: Vault now ensures that a policy named `default` is added
   to every token. This policy cannot be deleted, but it can be modified
   (including to an empty policy). There are three endpoints allowed in the
   default `default` policy, related to token self-management: `lookup-self`,
   which allows a token to retrieve its own information, and `revoke-self` and
   `renew-self`, which are self-explanatory. If your existing Vault
   installation contains a policy called `default`, it will not be overridden,
   but it will be added to each new token created. You can override this
   behavior when using manual token creation (i.e. not via an authentication
   backend) by setting the "no_default_policy" flag to true. [GH-732]

IMPROVEMENTS:

 * api: API client now uses a 60 second timeout instead of indefinite [GH-681]
 * api: Implement LookupSelf, RenewSelf, and RevokeSelf functions for auth
   tokens [GH-739]
 * api: Standardize environment variable reading logic inside the API; the CLI
   now uses this but can still override via command-line parameters [GH-618]
 * audit: HMAC-SHA256'd client tokens are now stored with each request entry.
   Previously they were only displayed at creation time; this allows much
   better traceability of client actions. [GH-713]
 * audit: There is now a `sys/audit-hash` endpoint that can be used to generate
   an HMAC-SHA256'd value from provided data using the given audit backend's
   salt [GH-784]
 * core: The physical storage read cache can now be disabled via
   "disable_cache" [GH-674]
 * core: The unsealing process can now be reset midway through (this feature
   was documented before, but not enabled) [GH-695]
 * core: Tokens can now renew themselves [GH-455]
 * core: Base64-encoded PGP keys can be used with the CLI for `init` and
   `rekey` operations [GH-653]
 * core: Print version on startup [GH-765]
 * core: Access to `sys/policy` and `sys/mounts` now uses the normal ACL system
   instead of requiring a root token [GH-769]
 * credential/token: Display whether or not a token is an orphan in the output
   of a lookup call [GH-766]
 * logical: Allow `.` in path-based variables in many more locations [GH-244]
 * logical: Responses now contain a "warnings" key containing a list of
   warnings returned from the server. These are conditions that did not require
   failing an operation, but of which the client should be aware. [GH-676]
 * physical/(consul,etcd): Consul and etcd now use a connection pool to limit
   the number of outstanding operations, improving behavior when a lot of
   operations must happen at once [GH-677] [GH-780]
 * physical/consul: The `datacenter` parameter was removed; It could not be
   effective unless the Vault node (or the Consul node it was connecting to)
   was in the datacenter specified, in which case it wasn't needed [GH-816]
 * physical/etcd: Support TLS-encrypted connections and use a connection pool
   to limit the number of outstanding operations [GH-780]
 * physical/s3: The S3 endpoint can now be configured, allowing using
   S3-API-compatible storage solutions [GH-750]
 * physical/s3: The S3 bucket can now be configured with the `AWS_S3_BUCKET`
   environment variable [GH-758]
 * secret/consul: Management tokens can now be created [GH-714]

BUG FIXES:

 * api: API client now checks for a 301 response for redirects. Vault doesn't
   generate these, but in certain conditions Go's internal HTTP handler can
   generate them, leading to client errors.
 * cli: `token-create` now supports the `ttl` parameter in addition to the
   deprecated `lease` parameter. [GH-688]
 * core: Return data from `generic` backends on the last use of a limited-use
   token [GH-615]
 * core: Fix upgrade path for leases created in `generic` prior to 0.3 [GH-673]
 * core: Stale leader entries will now be reaped [GH-679]
 * core: Using `mount-tune` on the auth/token path did not take effect.
   [GH-688]
 * core: Fix a potential race condition when (un)sealing the vault with metrics
   enabled [GH-694]
 * core: Fix an error that could happen in some failure scenarios where Vault
   could fail to revert to a clean state [GH-733]
 * core: Ensure secondary indexes are removed when a lease is expired [GH-749]
 * core: Ensure rollback manager uses an up-to-date mounts table [GH-771]
 * everywhere: Don't use http.DefaultClient, as it shares state implicitly and
   is a source of hard-to-track-down bugs [GH-700]
 * credential/token: Allow creating orphan tokens via an API path [GH-748]
 * secret/generic: Validate given duration at write time, not just read time;
   if stored durations are not parseable, return a warning and the default
   duration rather than an error [GH-718]
 * secret/generic: Return 400 instead of 500 when `generic` backend is written
   to with no data fields [GH-825]
 * secret/postgresql: Revoke permissions before dropping a user or revocation
   may fail [GH-699]

MISC:

 * Various documentation fixes and improvements [GH-685] [GH-688] [GH-697]
   [GH-710] [GH-715] [GH-831]

## 0.3.1 (October 6, 2015)

SECURITY:

 * core: In certain failure scenarios, the full values of requests and
   responses would be logged [GH-665]

FEATURES:

 * **Settable Maximum Open Connections**: The `mysql` and `postgresql` backends
   now allow setting the number of maximum open connections to the database,
   which was previously capped to 2. [GH-661]
 * **Renewable Tokens for GitHub**: The `github` backend now supports
   specifying a TTL, enabling renewable tokens. [GH-664]

BUG FIXES:

 * dist: linux-amd64 distribution was dynamically linked [GH-656]
 * credential/github: Fix acceptance tests [GH-651]

MISC:

 * Various minor documentation fixes and improvements [GH-649] [GH-650]
   [GH-654] [GH-663]

## 0.3.0 (September 28, 2015)

DEPRECATIONS/CHANGES:

Note: deprecations and breaking changes in upcoming releases are announced
ahead of time on the "vault-tool" mailing list.

 * **Cookie Authentication Removed**: As of 0.3 the only way to authenticate is
   via the X-Vault-Token header. Cookie authentication was hard to properly
   test, could result in browsers/tools/applications saving tokens in plaintext
   on disk, and other issues. [GH-564]
 * **Terminology/Field Names**: Vault is transitioning from overloading the
   term "lease" to mean both "a set of metadata" and "the amount of time the
   metadata is valid". The latter is now being referred to as TTL (or
   "lease_duration" for backwards-compatibility); some parts of Vault have
   already switched to using "ttl" and others will follow in upcoming releases.
   In particular, the "token", "generic", and "pki" backends accept both "ttl"
   and "lease" but in 0.4 only "ttl" will be accepted. [GH-528]
 * **Downgrade Not Supported**: Due to enhancements in the storage subsystem,
   values written by Vault 0.3+ will not be able to be read by prior versions
   of Vault. There are no expected upgrade issues, however, as with all
   critical infrastructure it is recommended to back up Vault's physical
   storage before upgrading.

FEATURES:

 * **SSH Backend**: Vault can now be used to delegate SSH access to machines,
   via a (recommended) One-Time Password approach or by issuing dynamic keys.
   [GH-385]
 * **Cubbyhole Backend**: This backend works similarly to the "generic" backend
   but provides a per-token workspace. This enables some additional
   authentication workflows (especially for containers) and can be useful to
   applications to e.g. store local credentials while being restarted or
   upgraded, rather than persisting to disk. [GH-612]
 * **Transit Backend Improvements**: The transit backend now allows key
   rotation and datakey generation. For rotation, data encrypted with previous
   versions of the keys can still be decrypted, down to a (configurable)
   minimum previous version; there is a rewrap function for manual upgrades of
   ciphertext to newer versions. Additionally, the backend now allows
   generating and returning high-entropy keys of a configurable bitsize
   suitable for AES and other functions; this is returned wrapped by a named
   key, or optionally both wrapped and plaintext for immediate use. [GH-626]
 * **Global and Per-Mount Default/Max TTL Support**: You can now set the
   default and maximum Time To Live for leases both globally and per-mount.
   Per-mount settings override global settings. Not all backends honor these
   settings yet, but the maximum is a hard limit enforced outside the backend.
   See the documentation for "/sys/mounts/" for details on configuring
   per-mount TTLs.  [GH-469]
 * **PGP Encryption for Unseal Keys**: When initializing or rotating Vault's
   master key, PGP/GPG public keys can now be provided. The output keys will be
   encrypted with the given keys, in order. [GH-570]
 * **Duo Multifactor Authentication Support**: Backends that support MFA can
   now use Duo as the mechanism. [GH-464]
 * **Performance Improvements**: Users of the "generic" backend will see a
   significant performance improvement as the backend no longer creates leases,
   although it does return TTLs (global/mount default, or set per-item) as
   before.  [GH-631]
 * **Codebase Audit**: Vault's codebase was audited by iSEC. (The terms of the
   audit contract do not allow us to make the results public.) [GH-220]

IMPROVEMENTS:

 * audit: Log entries now contain a time field [GH-495]
 * audit: Obfuscated audit entries now use hmac-sha256 instead of sha1 [GH-627]
 * backends: Add ability for a cleanup function to be called on backend unmount
   [GH-608]
 * config: Allow specifying minimum acceptable TLS version [GH-447]
 * core: If trying to mount in a location that is already mounted, be more
   helpful about the error [GH-510]
 * core: Be more explicit on failure if the issue is invalid JSON [GH-553]
 * core: Tokens can now revoke themselves [GH-620]
 * credential/app-id: Give a more specific error when sending a duplicate POST
   to sys/auth/app-id [GH-392]
 * credential/github: Support custom API endpoints (e.g. for Github Enterprise)
   [GH-572]
 * credential/ldap: Add per-user policies and option to login with
   userPrincipalName [GH-420]
 * credential/token: Allow root tokens to specify the ID of a token being
   created from CLI [GH-502]
 * credential/userpass: Enable renewals for login tokens [GH-623]
 * scripts: Use /usr/bin/env to find Bash instead of hardcoding [GH-446]
 * scripts: Use godep for build scripts to use same environment as tests
   [GH-404]
 * secret/mysql: Allow reading configuration data [GH-529]
 * secret/pki: Split "allow_any_name" logic to that and "enforce_hostnames", to
   allow for non-hostname values (e.g. for client certificates) [GH-555]
 * storage/consul: Allow specifying certificates used to talk to Consul
   [GH-384]
 * storage/mysql: Allow SSL encrypted connections [GH-439]
 * storage/s3: Allow using temporary security credentials [GH-433]
 * telemetry: Put telemetry object in configuration to allow more flexibility
   [GH-419]
 * testing: Disable mlock for testing of logical backends so as not to require
   root [GH-479]

BUG FIXES:

 * audit/file: Do not enable auditing if file permissions are invalid [GH-550]
 * backends: Allow hyphens in endpoint patterns (fixes AWS and others) [GH-559]
 * cli: Fixed missing setup of client TLS certificates if no custom CA was
   provided
 * cli/read: Do not include a carriage return when using raw field output
   [GH-624]
 * core: Bad input data could lead to a panic for that session, rather than
   returning an error [GH-503]
 * core: Allow SHA2-384/SHA2-512 hashed certificates [GH-448]
 * core: Do not return a Secret if there are no uses left on a token (since it
   will be unable to be used) [GH-615]
 * core: Code paths that called lookup-self would decrement num_uses and
   potentially immediately revoke a token [GH-552]
 * core: Some /sys/ paths would not properly redirect from a standby to the
   leader [GH-499] [GH-551]
 * credential/aws: Translate spaces in a token's display name to avoid making
   IAM unhappy [GH-567]
 * credential/github: Integration failed if more than ten organizations or
   teams [GH-489]
 * credential/token: Tokens with sudo access to "auth/token/create" can now use
   root-only options [GH-629]
 * secret/cassandra: Work around backwards-incompatible change made in
   Cassandra 2.2 preventing Vault from properly setting/revoking leases
   [GH-549]
 * secret/mysql: Use varbinary instead of varchar to avoid InnoDB/UTF-8 issues
   [GH-522]
 * secret/postgres: Explicitly set timezone in connections [GH-597]
 * storage/etcd: Renew semaphore periodically to prevent leadership flapping
   [GH-606]
 * storage/zk: Fix collisions in storage that could lead to data unavailability
   [GH-411]

MISC:

 * Various documentation fixes and improvements [GH-412] [GH-474] [GH-476]
   [GH-482] [GH-483] [GH-486] [GH-508] [GH-568] [GH-574] [GH-586] [GH-590]
   [GH-591] [GH-592] [GH-595] [GH-613] [GH-637]
 * Less "armon" in stack traces [GH-453]
 * Sourcegraph integration [GH-456]

## 0.2.0 (July 13, 2015)

FEATURES:

 * **Key Rotation Support**: The `rotate` command can be used to rotate the
   master encryption key used to write data to the storage (physical) backend.
   [GH-277]
 * **Rekey Support**: Rekey can be used to rotate the master key and change the
   configuration of the unseal keys (number of shares, threshold required).
   [GH-277]
 * **New secret backend: `pki`**: Enable Vault to be a certificate authority
   and generate signed TLS certificates. [GH-310]
 * **New secret backend: `cassandra`**: Generate dynamic credentials for
   Cassandra [GH-363]
 * **New storage backend: `etcd`**: store physical data in etcd [GH-259]
   [GH-297]
 * **New storage backend: `s3`**: store physical data in S3. Does not support
   HA. [GH-242]
 * **New storage backend: `MySQL`**: store physical data in MySQL. Does not
   support HA. [GH-324]
 * `transit` secret backend supports derived keys for per-transaction unique
   keys [GH-399]

IMPROVEMENTS:

 * cli/auth: Enable `cert` method [GH-380]
 * cli/auth: read input from stdin [GH-250]
 * cli/read: Ability to read a single field from a secret [GH-257]
 * cli/write: Adding a force flag when no input required
 * core: allow time duration format in place of seconds for some inputs
 * core: audit log provides more useful information [GH-360]
 * core: graceful shutdown for faster HA failover
 * core: **change policy format** to use explicit globbing [GH-400] Any
   existing policy in Vault is automatically upgraded to avoid issues.  All
   policy files must be updated for future writes. Adding the explicit glob
   character `*` to the path specification is all that is required.
 * core: policy merging to give deny highest precedence [GH-400]
 * credential/app-id: Protect against timing attack on app-id
 * credential/cert: Record the common name in the metadata [GH-342]
 * credential/ldap: Allow TLS verification to be disabled [GH-372]
 * credential/ldap: More flexible names allowed [GH-245] [GH-379] [GH-367]
 * credential/userpass: Protect against timing attack on password
 * credential/userpass: Use bcrypt for password matching
 * http: response codes improved to reflect error [GH-366]
 * http: the `sys/health` endpoint supports `?standbyok` to return 200 on
   standby [GH-389]
 * secret/app-id: Support deleting AppID and UserIDs [GH-200]
 * secret/consul: Fine grained lease control [GH-261]
 * secret/transit: Decouple raw key from key management endpoint [GH-355]
 * secret/transit: Upsert named key when encrypt is used [GH-355]
 * storage/zk: Support for HA configuration [GH-252]
 * storage/zk: Changing node representation. **Backwards incompatible**.
   [GH-416]

BUG FIXES:

 * audit/file: file removing TLS connection state
 * audit/syslog: fix removing TLS connection state
 * command/*: commands accepting `k=v` allow blank values
 * core: Allow building on FreeBSD [GH-365]
 * core: Fixed various panics when audit logging enabled
 * core: Lease renewal does not create redundant lease
 * core: fixed leases with negative duration [GH-354]
 * core: token renewal does not create child token
 * core: fixing panic when lease increment is null [GH-408]
 * credential/app-id: Salt the paths in storage backend to avoid information
   leak
 * credential/cert: Fixing client certificate not being requested
 * credential/cert: Fixing panic when no certificate match found [GH-361]
 * http: Accept PUT as POST for sys/auth
 * http: Accept PUT as POST for sys/mounts [GH-349]
 * http: Return 503 when sealed [GH-225]
 * secret/postgres: Username length is capped to exceeding limit
 * server: Do not panic if backend not configured [GH-222]
 * server: Explicitly check value of tls_diable [GH-201]
 * storage/zk: Fixed issues with version conflicts [GH-190]

MISC:

 * cli/path-help: renamed from `help` to avoid confusion

## 0.1.2 (May 11, 2015)

FEATURES:

  * **New physical backend: `zookeeper`**: store physical data in Zookeeper.
    HA not supported yet.
  * **New credential backend: `ldap`**: authenticate using LDAP credentials.

IMPROVEMENTS:

  * core: Auth backends can store internal data about auth creds
  * audit: display name for auth is shown in logs [GH-176]
  * command/*: `-insecure` has been renamed to `-tls-skip-verify` [GH-130]
  * command/*: `VAULT_TOKEN` overrides local stored auth [GH-162]
  * command/server: environment variables are copy-pastable
  * credential/app-id: hash of app and user ID are in metadata [GH-176]
  * http: HTTP API accepts `X-Vault-Token` as auth header [GH-124]
  * logical/*: Generate help output even if no synopsis specified

BUG FIXES:

  * core: login endpoints should never return secrets
  * core: Internal data should never be returned from core endpoints
  * core: defer barrier initialization to as late as possible to avoid error
    cases during init that corrupt data (no data loss)
  * core: guard against invalid init config earlier
  * audit/file: create file if it doesn't exist [GH-148]
  * command/*: ignore directories when traversing CA paths [GH-181]
  * credential/*: all policy mapping keys are case insensitive [GH-163]
  * physical/consul: Fixing path for locking so HA works in every case

## 0.1.1 (May 2, 2015)

SECURITY CHANGES:

  * physical/file: create the storge with 0600 permissions [GH-102]
  * token/disk: write the token to disk with 0600 perms

IMPROVEMENTS:

  * core: Very verbose error if mlock fails [GH-59]
  * command/*: On error with TLS oversized record, show more human-friendly
    error message. [GH-123]
  * command/read: `lease_renewable` is now outputted along with the secret to
    show whether it is renewable or not
  * command/server: Add configuration option to disable mlock
  * command/server: Disable mlock for dev mode so it works on more systems

BUG FIXES:

  * core: if token helper isn't absolute, prepend with path to Vault
    executable, not "vault" (which requires PATH) [GH-60]
  * core: Any "mapping" routes allow hyphens in keys [GH-119]
  * core: Validate `advertise_addr` is a valid URL with scheme [GH-106]
  * command/auth: Using an invalid token won't crash [GH-75]
  * credential/app-id: app and user IDs can have hyphens in keys [GH-119]
  * helper/password: import proper DLL for Windows to ask password [GH-83]

## 0.1.0 (April 28, 2015)

  * Initial release
