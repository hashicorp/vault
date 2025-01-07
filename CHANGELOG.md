## Previous versions
- [v1.10.0 - v1.15.16](CHANGELOG-v1.10-v1.15.md)
- [v1.0.0 - v1.9.10](CHANGELOG-pre-v1.10.md)
- [v0.11.6 and earlier](CHANGELOG-v0.md)

## 1.18.3
### December 18, 2024

CHANGES:

* secrets/openldap: Update plugin to v0.14.4 [[GH-29131](https://github.com/hashicorp/vault/pull/29131)]
* secrets/pki: Enforce the issuer constraint extensions (extended key usage, name constraints, issuer name) when issuing or signing leaf certificates. For more information see [PKI considerations](https://developer.hashicorp.com/vault/docs/secrets/pki/considerations#issuer-constraints-enforcement) [[GH-29045](https://github.com/hashicorp/vault/pull/29045)]

IMPROVEMENTS:

* auth/okta: update to okta sdk v5 from v2. Transitively updates go-jose dependency to >=3.0.3 to resolve GO-2024-2631. See https://github.com/okta/okta-sdk-golang/blob/master/MIGRATING.md for details on changes. [[GH-28121](https://github.com/hashicorp/vault/pull/28121)]
* core: Added new `enable_post_unseal_trace` and `post_unseal_trace_directory` config options to generate Go traces during the post-unseal step for debug purposes. [[GH-28895](https://github.com/hashicorp/vault/pull/28895)]
* sdk: Add Vault build date to system view plugin environment response [[GH-29082](https://github.com/hashicorp/vault/pull/29082)]
* ui: Replace KVv2 json secret details view with Hds::CodeBlock component allowing users to search the full secret height. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]

BUG FIXES:

* autosnapshots (enterprise): Fix an issue where snapshot size metrics were not reported for cloud-based storage.
* core/metrics: Fix unlocked mounts read for usage reporting. [[GH-29091](https://github.com/hashicorp/vault/pull/29091)]
* core/seal (enterprise): Fix problem with nodes unable to join Raft clusters with Seal High Availability enabled. [[GH-29117](https://github.com/hashicorp/vault/pull/29117)]
* core: fix bug in seal unwrapper that caused high storage latency in Vault CE. For every storage read request, the
seal unwrapper was performing the read twice, and would also issue an unnecessary storage write. [[GH-29050](https://github.com/hashicorp/vault/pull/29050)]
* secret/db: Update static role rotation to generate a new password after 2 failed attempts. [[GH-28989](https://github.com/hashicorp/vault/pull/28989)]
* ui: Allow users to search the full json object within the json code-editor edit/create view. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]
* ui: Decode `connection_url` to fix database connection updates (i.e. editing connection config, deleting roles) failing when urls include template variables. [[GH-29114](https://github.com/hashicorp/vault/pull/29114)]
* ui: Fix Swagger explorer bug where requests with path params were not working. [[GH-28670](https://github.com/hashicorp/vault/issues/28670)]
* vault/diagnose: Fix time to expiration reporting within the TLS verification to not be a month off. [[GH-29128](https://github.com/hashicorp/vault/pull/29128)]

## 1.18.2
### November 21, 2024

SECURITY:

* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20241115202008-166203013d8e

CHANGES:

* auth/azure: Update plugin to v0.19.2 [[GH-28848](https://github.com/hashicorp/vault/pull/28848)]
* core/ha (enterprise): Failed attempts to become a performance standby node are now using an exponential backoff instead of a
10 second delay in between retries. The backoff starts at 2s and increases by a factor of two until reaching
the maximum of 16s. This should make unsealing of the node faster in some cases.
* login (enterprise): Return a 500 error during logins when performance standby nodes make failed gRPC requests to the active node. [[GH-28807](https://github.com/hashicorp/vault/pull/28807)]

FEATURES:

* **Product Usage Reporting**: Added product usage reporting, which collects anonymous, numerical, non-sensitive data about Vault secrets usage, and adds it to the existing utilization reports. See the [[docs](https://developer.hashicorp.com/vault/docs/enterprise/license/product-usage-reporting)] for more info [[GH-28858](https://github.com/hashicorp/vault/pull/28858)]

IMPROVEMENTS:

* secret/pki: Introduce a new value `always_enforce_err` within `leaf_not_after_behavior` to force the error in all circumstances such as CA issuance and ACME requests if requested TTL values are beyond the issuer's NotAfter. [[GH-28907](https://github.com/hashicorp/vault/pull/28907)]
* secrets-sync (enterprise): No longer attempt to unsync a random UUID secret name in GCP upon destination creation.
* ui: Adds navigation for LDAP hierarchical roles [[GH-28824](https://github.com/hashicorp/vault/pull/28824)]
* website/docs: changed outdated reference to consul-helm repository to consul-k8s repository. [[GH-28825](https://github.com/hashicorp/vault/pull/28825)]

BUG FIXES:

* auth/ldap: Fixed an issue where debug level logging was not emitted. [[GH-28881](https://github.com/hashicorp/vault/pull/28881)]
* core: Improved an internal helper function that sanitizes paths by adding a check for leading backslashes 
in addition to the existing check for leading slashes. [[GH-28878](https://github.com/hashicorp/vault/pull/28878)]
* secret/pki: Fix a bug that prevents PKI issuer field enable_aia_url_templating
to be set to false. [[GH-28832](https://github.com/hashicorp/vault/pull/28832)]
* secrets-sync (enterprise): Fixed issue where secret-key granularity destinations could sometimes cause a panic when loading a sync status.
* secrets/aws: Fix issue with static credentials not rotating after restart or leadership change. [[GH-28775](https://github.com/hashicorp/vault/pull/28775)]
* secrets/ssh: Return the flag `allow_empty_principals` in the read role api when key_type is "ca" [[GH-28901](https://github.com/hashicorp/vault/pull/28901)]
* secrets/transform (enterprise): Fix nil panic when accessing a partially setup database store.
* secrets/transit: Fix a race in which responses from the key update api could contain results from another subsequent update [[GH-28839](https://github.com/hashicorp/vault/pull/28839)]
* ui: Fixes rendering issues of LDAP dynamic and static roles with the same name [[GH-28824](https://github.com/hashicorp/vault/pull/28824)]

## 1.18.1
### October 30, 2024

SECURITY:
* core/raft: Add raft join limits [[GH-28790](https://github.com/hashicorp/vault/pull/28790), [HCSEC-2024-26](https://discuss.hashicorp.com/t/hcsec-2024-26-vault-vulnerable-to-denial-of-service-through-memory-exhaustion-when-processing-raft-cluster-join-requests)]

CHANGES:

* auth/azure: Update plugin to v0.19.1 [[GH-28712](https://github.com/hashicorp/vault/pull/28712)]
* secrets/azure: Update plugin to v0.20.1 [[GH-28699](https://github.com/hashicorp/vault/pull/28699)]
* secrets/openldap: Update plugin to v0.14.1 [[GH-28479](https://github.com/hashicorp/vault/pull/28479)]
* secrets/openldap: Update plugin to v0.14.2 [[GH-28704](https://github.com/hashicorp/vault/pull/28704)]
* secrets/openldap: Update plugin to v0.14.3 [[GH-28780](https://github.com/hashicorp/vault/pull/28780)]

IMPROVEMENTS:

* core: Add a mount tuneable that trims trailing slashes of request paths during POST.  Needed to support CMPv2 in PKI. [[GH-28752](https://github.com/hashicorp/vault/pull/28752)]
* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20241003195753-88fef418d705
* ui: Add button to copy secret path in kv v1 and v2 secrets engines [[GH-28629](https://github.com/hashicorp/vault/pull/28629)]
* ui: Adds copy button to identity entity, alias and mfa method IDs [[GH-28742](https://github.com/hashicorp/vault/pull/28742)]

BUG FIXES:

* agent: Fix chown error running agent on Windows with an auto-auth file sinks. [[GH-28748](https://github.com/hashicorp/vault/pull/28748)]
* audit: Prevent users from enabling multiple audit devices of file type with the same file_path to write to. [[GH-28751](https://github.com/hashicorp/vault/pull/28751)]
* cli: Fixed a CLI precedence issue where -agent-address didn't override VAULT_AGENT_ADDR as it should [[GH-28574](https://github.com/hashicorp/vault/pull/28574)]
* core/seal (enterprise): Fix bug that caused seal generation information to be replicated, which prevented disaster recovery and performance replication clusters from using their own seal high-availability configuration.
* core/seal: Fix an issue that could cause reading from sys/seal-backend-status to return stale information. [[GH-28631](https://github.com/hashicorp/vault/pull/28631)]
* core: Fixed panic seen when performing help requests without /v1/ in the URL. [[GH-28669](https://github.com/hashicorp/vault/pull/28669)]
* kmip (enterprise): Use the default KMIP port for IPv6 addresses missing a port, for the listen_addrs configuration field, in order to match the existing IPv4 behavior
* namespaces (enterprise): Fix issue where namespace patch requests to a performance secondary would not patch the namespace's metadata.
* proxy: Fix chown error running proxy on Windows with an auto-auth file sink. [[GH-28748](https://github.com/hashicorp/vault/pull/28748)]
* secrets/pki: Address issue with ACME HTTP-01 challenges failing for IPv6 IPs due to improperly formatted URLs [[GH-28718](https://github.com/hashicorp/vault/pull/28718)]
* ui: No longer running decodeURIComponent on KVv2 list view allowing percent encoded data-octets in path name. [[GH-28698](https://github.com/hashicorp/vault/pull/28698)]
  
## 1.18.0 
## October 9, 2024

SECURITY:

* secrets/identity: A privileged Vault operator with write permissions to the root namespace's identity endpoint could escalate their privileges to Vault's root policy (CVE-2024-9180) [HCSEC-2024-21](https://discuss.hashicorp.com/t/hcsec-2024-21-vault-operators-in-root-namespace-may-elevate-their-privileges/70565)

CHANGES:

* activity (enterprise): filter all fields in client count responses by the request namespace [[GH-27790](https://github.com/hashicorp/vault/pull/27790)]
* activity (enterprise): remove deprecated fields distinct_entities and non_entity_tokens [[GH-27830](https://github.com/hashicorp/vault/pull/27830)]
* activity log: Deprecated the field "default_report_months". Instead, the billing start time will be used to determine the start time
when querying the activity log endpoints. [[GH-27350](https://github.com/hashicorp/vault/pull/27350)]
* activity log: Deprecates the current_billing_period field for /sys/internal/counters/activity. The default start time
will automatically be set the billing period start date. [[GH-27426](https://github.com/hashicorp/vault/pull/27426)]
* activity: The [activity export API](https://developer.hashicorp.com/vault/api-docs/system/internal-counters#activity-export) now requires the `sudo` ACL capability. [[GH-27846](https://github.com/hashicorp/vault/pull/27846)]
* activity: The [activity export API](https://developer.hashicorp.com/vault/api-docs/system/internal-counters#activity-export) now responds with a status of 204 instead 400 when no data exists within the time range specified by `start_time` and `end_time`. [[GH-28064](https://github.com/hashicorp/vault/pull/28064)]
* activity: The startTime will be set to the start of the current billing period by default.
The endTime will be set to the end of the current month. This applies to /sys/internal/counters/activity,
/sys/internal/counters/activity/export, and the vault operator usage command that utilizes /sys/internal/counters/activity. [[GH-27379](https://github.com/hashicorp/vault/pull/27379)]
* api: Update backoff/v3 to backoff/v4.3.0 [[GH-26868](https://github.com/hashicorp/vault/pull/26868)]
* auth/alicloud: Update plugin to v0.19.0 [[GH-28263](https://github.com/hashicorp/vault/pull/28263)]
* auth/azure: Update plugin to v0.19.0 [[GH-28294](https://github.com/hashicorp/vault/pull/28294)]
* auth/cf: Update plugin to v0.18.0 [[GH-27724](https://github.com/hashicorp/vault/pull/27724)]
* auth/cf: Update plugin to v0.19.0 [[GH-28266](https://github.com/hashicorp/vault/pull/28266)]
* auth/gcp: Update plugin to v0.19.0 [[GH-28366](https://github.com/hashicorp/vault/pull/28366)]
* auth/jwt: Update plugin to v0.21.0 [[GH-27498](https://github.com/hashicorp/vault/pull/27498)]
* auth/jwt: Update plugin to v0.22.0 [[GH-28349](https://github.com/hashicorp/vault/pull/28349)]
* auth/kerberos: Update plugin to v0.13.0 [[GH-28264](https://github.com/hashicorp/vault/pull/28264)]
* auth/kubernetes: Update plugin to v0.20.0 [[GH-28289](https://github.com/hashicorp/vault/pull/28289)]
* auth/oci: Update plugin to v0.17.0 [[GH-28307](https://github.com/hashicorp/vault/pull/28307)]
* cli: The undocumented `-dev-three-node` and `-dev-four-cluster` CLI options have been removed. [[GH-27578](https://github.com/hashicorp/vault/pull/27578)]
* consul-template: updated to version 0.39.1 [[GH-27799](https://github.com/hashicorp/vault/pull/27799)]
* core(enterprise): Updated the following two control group related errors responses to respond with response code 400 instead of 500: `control group: could not find token`, and `control group: token is not a valid control group token`.
* core: Bump Go version to 1.22.7
* database/couchbase: Update plugin to v0.12.0 [[GH-28327](https://github.com/hashicorp/vault/pull/28327)]
* database/elasticsearch: Update plugin to v0.16.0 [[GH-28277](https://github.com/hashicorp/vault/pull/28277)]
* database/mongodbatlas: Update plugin to v0.13.0 [[GH-28268](https://github.com/hashicorp/vault/pull/28268)]
* database/redis-elasticache: Update plugin to v0.5.0 [[GH-28293](https://github.com/hashicorp/vault/pull/28293)]
* database/redis: Update plugin to v0.4.0 [[GH-28404](https://github.com/hashicorp/vault/pull/28404)]
* database/snowflake: Update plugin to v0.12.0 [[GH-28275](https://github.com/hashicorp/vault/pull/28275)]
* sdk: Upgrade to go-secure-stdlib/plugincontainer@v0.4.0, which also bumps github.com/docker/docker to v26.1.5+incompatible [[GH-28269](https://github.com/hashicorp/vault/pull/28269)]
* secrets/ad: Update plugin to v0.19.0 [[GH-28361](https://github.com/hashicorp/vault/pull/28361)]
* secrets/alicloud: Update plugin to v0.18.0 [[GH-28271](https://github.com/hashicorp/vault/pull/28271)]
* secrets/azure: Update plugin to v0.19.2 [[GH-27652](https://github.com/hashicorp/vault/pull/27652)]
* secrets/azure: Update plugin to v0.20.0 [[GH-28267](https://github.com/hashicorp/vault/pull/28267)]
* secrets/gcp: Update plugin to v0.20.0 [[GH-28324](https://github.com/hashicorp/vault/pull/28324)]
* secrets/gcpkms: Update plugin to v0.18.0 [[GH-28300](https://github.com/hashicorp/vault/pull/28300)]
* secrets/gcpkms: Update plugin to v0.19.0 [[GH-28360](https://github.com/hashicorp/vault/pull/28360)]
* secrets/kubernetes: Update plugin to v0.9.0 [[GH-28287](https://github.com/hashicorp/vault/pull/28287)]
* secrets/kv: Update plugin to v0.20.0 [[GH-28334](https://github.com/hashicorp/vault/pull/28334)]
* secrets/mongodbatlas: Update plugin to v0.13.0 [[GH-28348](https://github.com/hashicorp/vault/pull/28348)]
* secrets/openldap: Update plugin to v0.14.0 [[GH-28325](https://github.com/hashicorp/vault/pull/28325)]
* secrets/ssh: Add a flag, `allow_empty_principals` to allow keys or certs to apply to any user/principal. [[GH-28466](https://github.com/hashicorp/vault/pull/28466)]
* secrets/terraform: Update plugin to v0.10.0 [[GH-28312](https://github.com/hashicorp/vault/pull/28312)]
* secrets/terraform: Update plugin to v0.9.0 [[GH-28016](https://github.com/hashicorp/vault/pull/28016)]
* ui: Uses the internal/counters/activity/export endpoint for client count export data. [[GH-27455](https://github.com/hashicorp/vault/pull/27455)]

FEATURES:

* **AWS secrets engine STS session tags support**: Adds support for setting STS
session tags when generating temporary credentials using the AWS secrets
engine. [[GH-27620](https://github.com/hashicorp/vault/pull/27620)]
* **Adaptive Overload Protection (enterprise)**: Enables Adaptive Overload Protection
for write requests as a GA feature (enabled by default) for Integrated Storage.
* **Audit Entry Exclusion (enterprise)**: Audit devices support excluding fields from entries being written to them, with expression-based rules (powered by go-bexpr) to determine when the specific fields are excluded.
* **Workload Identity Federation UI for AWS (enterprise)**: Add WIF fields to AWS secrets engine. [[GH-28148](https://github.com/hashicorp/vault/pull/28148)]
* **KV v2 Patch/Subkey (enterprise)**: Adds GUI support to read the subkeys of a KV v2 secret and patch (partially update) secret data. [[GH-28212](https://github.com/hashicorp/vault/pull/28212)]
* **Self-Managed Static Roles**: Self-Managed Static Roles are now supported for the Postgres SQL database engine. Requires Vault Enterprise. [[GH-28199](https://github.com/hashicorp/vault/pull/28199)]
* **Vault Minimal Version**: Add the ability to build a minimal version of Vault
with only core features using the BUILD_MINIMAL environment variable. [[GH-27394](https://github.com/hashicorp/vault/pull/27394)]
* **Vault PKI 3GPP CMPv2 Server (Enterprise)**: Support for the PKI 3GPP CMPv2 certificate management protocol has been added to the Vault PKI Plugin. This allows standard CMPv2 clients to request certificates from a Vault server with no knowledge of Vault APIs.

IMPROVEMENTS:

* activity log: Changes how new client counts in the current month are estimated, in order to return more
visibly sensible totals. [[GH-27547](https://github.com/hashicorp/vault/pull/27547)]
* activity: The [activity export API](https://developer.hashicorp.com/vault/api-docs/system/internal-counters#activity-export) can now be called in non-root namespaces. Resulting records will be filtered to include the requested namespace (via `X-Vault-Namespace` header or within the path) and all child namespaces. [[GH-27846](https://github.com/hashicorp/vault/pull/27846)]
* activity: The [activity export API](https://developer.hashicorp.com/vault/api-docs/system/internal-counters#activity-export) now includes identity metadata about entity clients. [[GH-28064](https://github.com/hashicorp/vault/pull/28064)]
* activity: `/sys/internal/counters/activity` will now include a warning if the specified usage period contains estimated client counts. [[GH-28068](https://github.com/hashicorp/vault/pull/28068)]
* agent/sink: Allow configuration of the user and group ID of the file sink. [[GH-27123](https://github.com/hashicorp/vault/pull/27123)]
* agent: Add metric (vault.agent.authenticated) that is set to 1 when vault agent has a valid token and zero if it does not. [[GH-26570](https://github.com/hashicorp/vault/pull/26570)]
* agent: Add the ability to dump pprof to the filesystem using SIGUSR2 [[GH-27510](https://github.com/hashicorp/vault/pull/27510)]
* audit: Adds TRACE logging to log request/response under certain circumstances, and further improvements to the audit subsystem. [[GH-28056](https://github.com/hashicorp/vault/pull/28056)]
* audit: Ensure that any underyling errors from audit devices are logged even if we consider auditing to be a success. [[GH-27809](https://github.com/hashicorp/vault/pull/27809)]
* audit: Internal implementation changes to the audit subsystem which improve performance. [[GH-27952](https://github.com/hashicorp/vault/pull/27952)]
* audit: Internal implementation changes to the audit subsystem which improve relability. [[GH-28286](https://github.com/hashicorp/vault/pull/28286)]
* audit: sinks (file, socket, syslog) will attempt to log errors to the server operational
log before returning (if there are errors to log, and the context is done). [[GH-27859](https://github.com/hashicorp/vault/pull/27859)]
* auth/cert: Cache full list of role trust information separately to avoid
eviction, and avoid duplicate loading during multiple simultaneous logins on
the same role. [[GH-27902](https://github.com/hashicorp/vault/pull/27902)]
* cli: Add a `--dev-no-kv` flag to prevent auto mounting a key-value secret backend when running a dev server [[GH-16974](https://github.com/hashicorp/vault/pull/16974)]
* cli: Allow vault CLI HTTP headers to be specified using the JSON-encoded VAULT_HEADERS environment variable [[GH-21993](https://github.com/hashicorp/vault/pull/21993)]
* cli: `vault operator usage` will now include a warning if the specified usage period contains estimated client counts. [[GH-28068](https://github.com/hashicorp/vault/pull/28068)]
* core/activity: Ensure client count queries that include the current month return consistent results by sorting the clients before performing estimation [[GH-28062](https://github.com/hashicorp/vault/pull/28062)]
* core/cli: Example 'help' pages for vault read / write docs improved. [[GH-19064](https://github.com/hashicorp/vault/pull/19064)]
* core/identity: allow identity backend to be tuned using standard secrets backend tuning parameters. [[GH-14723](https://github.com/hashicorp/vault/pull/14723)]
* core/metrics: ensure core HA metrics are always output to Prometheus. [[GH-27966](https://github.com/hashicorp/vault/pull/27966)]
* core: log at level ERROR rather than INFO when all seals are unhealthy. [[GH-28564](https://github.com/hashicorp/vault/pull/28564)]
* core: make authLock and mountsLock in Core configurable via the detect_deadlocks configuration parameter. [[GH-27633](https://github.com/hashicorp/vault/pull/27633)]
* database/postgres: Add new fields to the plugin's config endpoint for client certificate authentication. [[GH-28024](https://github.com/hashicorp/vault/pull/28024)]
* db/cassandra: Add `disable_host_initial_lookup` option to backend, allowing the disabling of initial host lookup. [[GH-9733](https://github.com/hashicorp/vault/pull/9733)]
* identity: alias metadata is now returned when listing entity aliases [[GH-26073](https://github.com/hashicorp/vault/pull/26073)]
* license utilization reporting (enterprise): Auto-roll billing start date. [[GH-27656](https://github.com/hashicorp/vault/pull/27656)]
* physical/raft: Log when the MAP_POPULATE mmap flag gets disabled before opening the database. [[GH-28526](https://github.com/hashicorp/vault/pull/28526)]
* proxy/sink: Allow configuration of the user and group ID of the file sink. [[GH-27123](https://github.com/hashicorp/vault/pull/27123)]
* proxy: Add the ability to dump pprof to the filesystem using SIGUSR2 [[GH-27510](https://github.com/hashicorp/vault/pull/27510)]
* raft-snapshot (enterprise): add support for managed identity credentials for azure snapshots
* raft/autopilot: Persist Raft server versions so autopilot always knows the versions of all servers in the cluster. Include server versions in the Raft bootstrap challenge answer so autopilot immediately knows the versions of new nodes. [[GH-28186](https://github.com/hashicorp/vault/pull/28186)]
* sdk/helper: Allow setting environment variables when using NewTestDockerCluster [[GH-27457](https://github.com/hashicorp/vault/pull/27457)]
* secrets-sync (enterprise): add support for specifying the replication regions for secret storage within GCP Secret Manager destinations
* secrets-sync (enterprise): add support for syncing secrets to github environments within repositories
* secrets-sync (enterprise): add support for syncing secrets to github organizations (beta)
* secrets/database/hana: Update HANA db client to v1.10.1 [[GH-27950](https://github.com/hashicorp/vault/pull/27950)]
* secrets/database: Add support for GCP CloudSQL private IP's. [[GH-26828](https://github.com/hashicorp/vault/pull/26828)]
* secrets/pki: Key Usage can now be set on intermediate and root CAs, and CSRs generated by the PKI secret's engine. [[GH-28237](https://github.com/hashicorp/vault/pull/28237)]
* secrets/pki: Track the last time auto-tidy ran to address auto-tidy not running if the auto-tidy interval is longer than scheduled Vault restarts. [[GH-28488](https://github.com/hashicorp/vault/pull/28488)]
* serviceregistration: Added support for Consul ServiceMeta tags from config file from the new `service_meta` config field. [[GH-11084](https://github.com/hashicorp/vault/pull/11084)]
* storage/azure: Updated metadata endpoint to `GetMSIEndpoint`, which supports more than just the metadata service. [[GH-10624](https://github.com/hashicorp/vault/pull/10624)]
* storage/dynamodb: Speed up list and delete of large directories by only requesting keys from DynamoDB [[GH-21159](https://github.com/hashicorp/vault/pull/21159)]
* storage/etcd: Update etcd3 client to v3.5.13 to allow use of TLSv1.3. [[GH-26660](https://github.com/hashicorp/vault/pull/26660)]
* storage/raft: Bump raft to v1.7.0 which includes pre-vote. This should make clusters more stable during network partitions. [[GH-27605](https://github.com/hashicorp/vault/pull/27605)]
* storage/raft: Improve autopilot logging on startup to show config values clearly and avoid spurious logs [[GH-27464](https://github.com/hashicorp/vault/pull/27464)]
* ui/secrets-sync: Hide Secrets Sync from the sidebar nav if user does not have access to the feature. [[GH-27262](https://github.com/hashicorp/vault/pull/27262)]
* ui: AWS credentials form sets credential_type from backing role [[GH-27405](https://github.com/hashicorp/vault/pull/27405)]
* ui: Creates separate section for updating sensitive creds for Secrets sync create/edit view. [[GH-27538](https://github.com/hashicorp/vault/pull/27538)]
* ui: For AWS and SSH secret engines hide mount configuration details in toggle and display configuration details or cta. [[GH-27831](https://github.com/hashicorp/vault/pull/27831)]
* ui: Mask obfuscated fields when creating/editing a Secrets sync destination. [[GH-27348](https://github.com/hashicorp/vault/pull/27348)]
* ui: Move secret-engine configuration create/edit from routing `vault/settings/secrets/configure/<backend>` to  `vault/secrets/<backend>/configuration/edit` [[GH-27918](https://github.com/hashicorp/vault/pull/27918)]
* ui: Remove deprecated `current_billing_period` from dashboard activity log request [[GH-27559](https://github.com/hashicorp/vault/pull/27559)]
* ui: Update the client count dashboard to use API namespace filtering and other UX improvements [[GH-28036](https://github.com/hashicorp/vault/pull/28036)]
* ui: remove initial start/end parameters on the activity call for client counts dashboard. [[GH-27816](https://github.com/hashicorp/vault/pull/27816)]
* ui: simplify the date range editing experience in the client counts dashboard. [[GH-27796](https://github.com/hashicorp/vault/pull/27796)]
* website/docs: Added API documentation for Azure Secrets Engine delete role [[GH-27883](https://github.com/hashicorp/vault/pull/27883)]
* website/docs: corrected invalid json in sample payload for azure secrets engine create/update role [[GH-28076](https://github.com/hashicorp/vault/pull/28076)]

BUG FIXES:

* activity: The sys/internal/counters/activity endpoint will return current month data when the end_date parameter is set to a future date. [[GH-28042](https://github.com/hashicorp/vault/pull/28042)]
* agent: Fixed an issue causing excessive CPU usage during normal operation [[GH-27518](https://github.com/hashicorp/vault/pull/27518)]
* auth/appid, auth/cert, auth/github, auth/ldap, auth/okta, auth/radius, auth/userpass: fixed an issue with policy name normalization that would prevent a token associated with a policy containing an uppercase character to be renewed. [[GH-16484](https://github.com/hashicorp/vault/pull/16484)]
* auth/aws: fixes an issue where not supplying an external id was interpreted as an empty external id [[GH-27858](https://github.com/hashicorp/vault/pull/27858)]
* auth/cert: During certificate validation, OCSP requests are debug logged even if Vault's log level is above DEBUG. [[GH-28450](https://github.com/hashicorp/vault/pull/28450)]
* auth/cert: Merge error messages returned in login failures and include error when present [[GH-27202](https://github.com/hashicorp/vault/pull/27202)]
* auth/cert: Use subject's serial number, not issuer's within error message text in OCSP request errors [[GH-27696](https://github.com/hashicorp/vault/pull/27696)]
* auth/cert: When using ocsp_ca_certificates, an error was produced though extra certs validation succeeded. [[GH-28597](https://github.com/hashicorp/vault/pull/28597)]
* auth/cert: ocsp_ca_certificates field was not honored when validating OCSP responses signed by a CA that did not issue the certificate. [[GH-28309](https://github.com/hashicorp/vault/pull/28309)]
* auth/token: Fix token TTL calculation so that it uses `max_lease_ttl` tune value for tokens created via `auth/token/create`. [[GH-28498](https://github.com/hashicorp/vault/pull/28498)]
* auth/token: fixes an edge case bug that "identity_policies" is nil and causes cli vault login error [[GH-17007](https://github.com/hashicorp/vault/pull/17007)]
* auth: Updated error handling for missing login credentials in AppRole and UserPass auth methods to return a 400 error instead of a 500 error. [[GH-28441](https://github.com/hashicorp/vault/pull/28441)]
* cli: Fixed an erroneous warning appearing about `-address` not being set when it is. [[GH-27265](https://github.com/hashicorp/vault/pull/27265)]
* cli: Fixed issue with `vault hcp connect` where HCP resources with uppercase letters were inaccessible when entering the correct project name. [[GH-27694](https://github.com/hashicorp/vault/pull/27694)]
* command: The `vault secrets move` and `vault auth move` command will no longer attempt to write to storage on performance standby nodes. [[GH-28059](https://github.com/hashicorp/vault/pull/28059)]
* config: Vault TCP listener config now correctly supports the documented proxy_protocol_behavior
setting of 'deny_unauthorized' [[GH-27459](https://github.com/hashicorp/vault/pull/27459)]
* core (enterprise): Fix 500 errors that occurred querying `sys/internal/ui/mounts` for a mount prefixed by a namespace path when path filters are configured. [[GH-27939](https://github.com/hashicorp/vault/pull/27939)]
* core (enterprise): Fix HTTP redirects in namespaces to use the correct path and (in the case of event subscriptions) the correct URI scheme. [[GH-27660](https://github.com/hashicorp/vault/pull/27660)]
* core (enterprise): Fix deletion of MFA login-enforcement configurations on standby nodes
* core/audit: Audit logging a Vault request/response checks if the existing context
is cancelled and will now use a new context with a 5 second timeout.
If the existing context is cancelled a new context, will be used. [[GH-27531](https://github.com/hashicorp/vault/pull/27531)]
* core/config: fix issue when using `proxy_protocol_behavior` with `deny_unauthorized`,
which causes the Vault TCP listener to close after receiving an untrusted upstream proxy connection. [[GH-27589](https://github.com/hashicorp/vault/pull/27589)]
* core/identity: Fixed an issue where deleted/reassigned entity-aliases were not removed from in-memory database. [[GH-27750](https://github.com/hashicorp/vault/pull/27750)]
* core/seal (enterprise): Fix bug that caused seal generation information to be replicated, which prevented disaster recovery and performance replication clusters from using their own seal high-availability configuration.
* core: Fixed an issue where maximum request duration timeout was not being added to all requests containing strings sys/monitor and sys/events. With this change, timeout is now added to all requests except monitor and events endpoint. [[GH-28230](https://github.com/hashicorp/vault/pull/28230)]
* core: Fixed an issue with performance standbys not being able to handle rotate root requests. [[GH-27631](https://github.com/hashicorp/vault/pull/27631)]
* database/postgresql: Fix potential error revoking privileges in postgresql database secrets engine when a schema contains special characters [[GH-28519](https://github.com/hashicorp/vault/pull/28519)]
* databases: fix issue where local timezone was getting lost when using a rotation schedule cron [[GH-28509](https://github.com/hashicorp/vault/pull/28509)]
* helper/pkcs7: Fix parsing certain messages containing only certificates [[GH-27435](https://github.com/hashicorp/vault/pull/27435)]
* identity/oidc: prevent JWKS from being generated by multiple concurrent requests [[GH-27929](https://github.com/hashicorp/vault/pull/27929)]
* licensing (enterprise): fixed issue where billing start date might not be correctly updated on performance standbys
* proxy/cache (enterprise): Fixed a data race that could occur while tracking capabilities in Proxy's static secret cache. [[GH-28494](https://github.com/hashicorp/vault/pull/28494)]
* proxy/cache (enterprise): Fixed an issue where Proxy with static secret caching enabled would not correctly handle requests to older secret versions for KVv2 secrets. Proxy's static secret cache now properly handles all requests relating to older versions for KVv2 secrets. [[GH-28207](https://github.com/hashicorp/vault/pull/28207)]
* proxy/cache (enterprise): Fixed an issue where Proxy would not correctly update KV secrets when talking to a perf standby. Proxy will now attempt to forward requests to update secrets triggered by events to the active node. Note that this requires `allow_forwarding_via_header` to be configured on the cluster. [[GH-27891](https://github.com/hashicorp/vault/pull/27891)]
* proxy/cache (enterprise): Fixed an issue where cached static secrets could fail to update if the secrets belonged to a non-root namespace. [[GH-27730](https://github.com/hashicorp/vault/pull/27730)]
* proxy: Fixed an issue causing excessive CPU usage during normal operation [[GH-27518](https://github.com/hashicorp/vault/pull/27518)]
* raft/autopilot: Fixed panic that may occur during shutdown [[GH-27726](https://github.com/hashicorp/vault/pull/27726)]
* replication (enterprise): fix cache invalidation issue leading to namespace custom metadata not being shown correctly on performance secondaries
* secrets-sync (enterprise): Destination set/remove operations will no longer be blocked as "purge in progress" after a purge job ended in failure.
* secrets-sync (enterprise): Fix KV secret access sometimes being denied, due to a double forward-slash (`//`) in the mount path, when the token should otherwise have access.
* secrets-sync (enterprise): Normalize custom_tag keys and values for recoverable invalid characters.
* secrets-sync (enterprise): Normalize secret key names before storing the external_name in a secret association.
* secrets-sync (enterprise): Patching github sync destination credentials will properly update and save the new credentials.
* secrets-sync (enterprise): Properly remove tags from secrets in AWS when they are removed from the source association
* secrets-sync (enterprise): Return an error immediately on destination creation when providing invalid custom_tags based on destination type.
* secrets-sync (enterprise): Return more accurate error code for invalid connection details
* secrets-sync (enterprise): Secondary nodes in a cluster now properly check activation-flags values.
* secrets-sync (enterprise): Skip invalid GitHub repository names when creating destinations
* secrets-sync (enterprise): Validate corresponding GitHub app parameters `app_name` and `installation_id` are set
* secrets/database: Skip connection verification on reading existing DB connection configuration [[GH-28139](https://github.com/hashicorp/vault/pull/28139)]
* secrets/identity (enterprise): Fix a bug that can cause DR promotion to fail in rare cases where a PR secondary has inconsistent alias information in storage.
* secrets/pki: fix lack of serial number to a certificate read resulting in a server side error. [[GH-27681](https://github.com/hashicorp/vault/pull/27681)]
* secrets/transit (enterprise): Fix an issue that caused input data be returned as part of generated CMAC values.
* storage/azure: Fix invalid account name initialization bug [[GH-27563](https://github.com/hashicorp/vault/pull/27563)]
* storage/raft (enterprise): Fix issue with namespace cache not getting cleared on snapshot restore, resulting in namespaces not found in the snapshot being inaccurately represented by API responses. [[GH-27474](https://github.com/hashicorp/vault/pull/27474)]
* storage/raft: Fix auto_join not working with mDNS provider. [[GH-25080](https://github.com/hashicorp/vault/pull/25080)]
* sys: Fix a bug where mounts of external plugins that were registered before Vault v1.0.0 could not be tuned to
use versioned plugins. [[GH-27881](https://github.com/hashicorp/vault/pull/27881)]
* ui: Allow creation of session_token type roles for AWS secret backend [[GH-27424](https://github.com/hashicorp/vault/pull/27424)]
* ui: Display an error and force a timeout when TOTP passcode is incorrect [[GH-27574](https://github.com/hashicorp/vault/pull/27574)]
* ui: Ensure token expired banner displays when batch token expires [[GH-27479](https://github.com/hashicorp/vault/pull/27479)]
* ui: Fix UI improperly checking capabilities for enabling performance and dr replication [[GH-28371](https://github.com/hashicorp/vault/pull/28371)]
* ui: Fix cursor jump on KVv2 json editor that would occur after pressing ENTER. [[GH-27569](https://github.com/hashicorp/vault/pull/27569)]
* ui: fix `default_role` input missing from oidc auth method configuration form [[GH-28539](https://github.com/hashicorp/vault/pull/28539)]
* ui: fix issue where enabling then disabling "Tidy ACME" in PKI results in failed API call. [[GH-27742](https://github.com/hashicorp/vault/pull/27742)]
* ui: fix namespace picker not working when in small screen where the sidebar is collapsed by default. [[GH-27728](https://github.com/hashicorp/vault/pull/27728)]
* ui: fixes renew-self being called right after login for non-renewable tokens [[GH-28204](https://github.com/hashicorp/vault/pull/28204)]
* ui: fixes toast (flash) alert message saying "created" when deleting a kv v2 secret [[GH-28093](https://github.com/hashicorp/vault/pull/28093)]

## 1.17.10 Enterprise
### December 18, 2024

CHANGES:

* secrets/pki: Enforce the issuer constraint extensions (extended key usage, name constraints, issuer name) when issuing or signing leaf certificates. For more information see [PKI considerations](https://developer.hashicorp.com/vault/docs/secrets/pki/considerations#issuer-constraints-enforcement) [[GH-29045](https://github.com/hashicorp/vault/pull/29045)]

IMPROVEMENTS:

* auth/okta: update to okta sdk v5 from v2. Transitively updates go-jose dependency to >=3.0.3 to resolve GO-2024-2631. See https://github.com/okta/okta-sdk-golang/blob/master/MIGRATING.md for details on changes. [[GH-28121](https://github.com/hashicorp/vault/pull/28121)]
* core: Added new `enable_post_unseal_trace` and `post_unseal_trace_directory` config options to generate Go traces during the post-unseal step for debug purposes. [[GH-28895](https://github.com/hashicorp/vault/pull/28895)]
* sdk: Add Vault build date to system view plugin environment response [[GH-29082](https://github.com/hashicorp/vault/pull/29082)]
* ui: Replace KVv2 json secret details view with Hds::CodeBlock component allowing users to search the full secret height. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]

BUG FIXES:

* auth/ldap: Fixed an issue where debug level logging was not emitted. [[GH-28881](https://github.com/hashicorp/vault/pull/28881)]
* autosnapshots (enterprise): Fix an issue where snapshot size metrics were not reported for cloud-based storage.
* core/metrics: Fix unlocked mounts read for usage reporting. [[GH-29091](https://github.com/hashicorp/vault/pull/29091)]
* core/seal (enterprise): Fix decryption of the raft bootstrap challenge when using seal high availability. [[GH-29117](https://github.com/hashicorp/vault/pull/29117)]
* secret/db: Update static role rotation to generate a new password after 2 failed attempts. [[GH-28989](https://github.com/hashicorp/vault/pull/28989)]
* ui: Allow users to search the full json object within the json code-editor edit/create view. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]
* ui: Decode `connection_url` to fix database connection updates (i.e. editing connection config, deleting roles) failing when urls include template variables. [[GH-29114](https://github.com/hashicorp/vault/pull/29114)]
* ui: Fix Swagger explorer bug where requests with path params were not working. [[GH-28670](https://github.com/hashicorp/vault/issues/28670)]
* vault/diagnose: Fix time to expiration reporting within the TLS verification to not be a month off. [[GH-29128](https://github.com/hashicorp/vault/pull/29128)]

## 1.17.9 Enterprise
### November 21, 2024

SECURITY:

* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20241115202008-166203013d8e

CHANGES:

* activity log: Deprecated the field "default_report_months". Instead, the billing start time will be used to determine the start time
when querying the activity log endpoints. [[GH-27350](https://github.com/hashicorp/vault/pull/27350)]
* core/ha (enterprise): Failed attempts to become a performance standby node are now using an exponential backoff instead of a
10 second delay in between retries. The backoff starts at 2s and increases by a factor of two until reaching
the maximum of 16s. This should make unsealing of the node faster in some cases.
* login (enterprise): Return a 500 error during logins when performance standby nodes make failed gRPC requests to the active node. [[GH-28807](https://github.com/hashicorp/vault/pull/28807)]

FEATURES:

* **Product Usage Reporting**: Added product usage reporting, which collects anonymous, numerical, non-sensitive data about Vault secrets usage, and adds it to the existing utilization reports. See the [[docs](https://developer.hashicorp.com/vault/docs/enterprise/license/product-usage-reporting)] for more info [[GH-28858](https://github.com/hashicorp/vault/pull/28858)]

IMPROVEMENTS:

* secrets-sync (enterprise): No longer attempt to unsync a random UUID secret name in GCP upon destination creation.
* ui: Adds navigation for LDAP hierarchical roles [[GH-28824](https://github.com/hashicorp/vault/pull/28824)]

BUG FIXES:

* core: Improved an internal helper function that sanitizes paths by adding a check for leading backslashes 
in addition to the existing check for leading slashes. [[GH-28878](https://github.com/hashicorp/vault/pull/28878)]
* secret/pki: Fix a bug that prevents PKI issuer field enable_aia_url_templating
to be set to false. [[GH-28832](https://github.com/hashicorp/vault/pull/28832)]
* secrets-sync (enterprise): Fixed issue where secret-key granularity destinations could sometimes cause a panic when loading a sync status.
* secrets/aws: Fix issue with static credentials not rotating after restart or leadership change. [[GH-28775](https://github.com/hashicorp/vault/pull/28775)]
* secrets/ssh: Return the flag `allow_empty_principals` in the read role api when key_type is "ca" [[GH-28901](https://github.com/hashicorp/vault/pull/28901)]
* secrets/transform (enterprise): Fix nil panic when accessing a partially setup database store.
* secrets/transit: Fix a race in which responses from the key update api could contain results from another subsequent update [[GH-28839](https://github.com/hashicorp/vault/pull/28839)]
* ui: Fixes rendering issues of LDAP dynamic and static roles with the same name [[GH-28824](https://github.com/hashicorp/vault/pull/28824)]

## 1.17.8 Enterprise
### October 30, 2024

SECURITY:
* core/raft: Add raft join limits [[GH-28790](https://github.com/hashicorp/vault/pull/28790), [HCSEC-2024-26](https://discuss.hashicorp.com/t/hcsec-2024-26-vault-vulnerable-to-denial-of-service-through-memory-exhaustion-when-processing-raft-cluster-join-requests)]

CHANGES:

* auth/azure: Update plugin to v0.18.1
* secrets/openldap: update plugin to v0.13.2

IMPROVEMENTS:

* ui: Add button to copy secret path in kv v1 and v2 secrets engines [[GH-28629](https://github.com/hashicorp/vault/pull/28629)]
* ui: Adds copy button to identity entity, alias and mfa method IDs [[GH-28742](https://github.com/hashicorp/vault/pull/28742)]

BUG FIXES:

* audit: Prevent users from enabling multiple audit devices of file type with the same file_path to write to. [[GH-28751](https://github.com/hashicorp/vault/pull/28751)]
* core/seal (enterprise): Fix bug that caused seal generation information to be replicated, which prevented disaster recovery and performance replication clusters from using their own seal high-availability configuration.
* core/seal: Fix an issue that could cause reading from sys/seal-backend-status to return stale information. [[GH-28631](https://github.com/hashicorp/vault/pull/28631)]
* core: Fixed panic seen when performing help requests without /v1/ in the URL. [[GH-28669](https://github.com/hashicorp/vault/pull/28669)]
* namespaces (enterprise): Fix issue where namespace patch requests to a performance secondary would not patch the namespace's metadata.
* secrets/pki: Address issue with ACME HTTP-01 challenges failing for IPv6 IPs due to improperly formatted URLs [[GH-28718](https://github.com/hashicorp/vault/pull/28718)]
* ui: No longer running decodeURIComponent on KVv2 list view allowing percent encoded data-octets in path name. [[GH-28698](https://github.com/hashicorp/vault/pull/28698)]

## 1.17.7 Enterprise
### October 09, 2024

SECURITY:

* secrets/identity: A privileged Vault operator with write permissions to the root namespace's identity endpoint could escalate their privileges to Vault's root policy (CVE-2024-9180) [HCSEC-2024-21](https://discuss.hashicorp.com/t/hcsec-2024-21-vault-operators-in-root-namespace-may-elevate-their-privileges/70565)

IMPROVEMENTS:

* core: log at level ERROR rather than INFO when all seals are unhealthy. [[GH-28564](https://github.com/hashicorp/vault/pull/28564)]
* physical/raft: Log when the MAP_POPULATE mmap flag gets disabled before opening the database. [[GH-28526](https://github.com/hashicorp/vault/pull/28526)]
* secrets/pki: Track the last time auto-tidy ran to address auto-tidy not running if the auto-tidy interval is longer than scheduled Vault restarts. [[GH-28488](https://github.com/hashicorp/vault/pull/28488)]

BUG FIXES:

* auth/cert: When using ocsp_ca_certificates, an error was produced though extra certs validation succeeded. [[GH-28597](https://github.com/hashicorp/vault/pull/28597)]
* auth/token: Fix token TTL calculation so that it uses `max_lease_ttl` tune value for tokens created via `auth/token/create`. [[GH-28498](https://github.com/hashicorp/vault/pull/28498)]
* databases: fix issue where local timezone was getting lost when using a rotation schedule cron [[GH-28509](https://github.com/hashicorp/vault/pull/28509)]
* secrets-sync (enterprise): Fix KV secret access sometimes being denied, due to a double forward-slash (`//`) in the mount path, when the token should otherwise have access.

## 1.17.6
### September 25, 2024

SECURITY:
* secrets/ssh: require `valid_principals` to contain a value or `default_user` be set by default to guard against potentially insecure configurations. `allow_empty_principals` can be used for backwards compatibility [HCSEC-2024-20](https://discuss.hashicorp.com/t/hcsec-2024-20-vault-ssh-secrets-engine-configuration-did-not-restrict-valid-principals-by-default/70251)

CHANGES:

* core: Bump Go version to 1.22.7
* secrets/ldap: Update vault-plugin-secrets-openldap to v0.13.1 [[GH-28478](https://github.com/hashicorp/vault/pull/28478)]
* secrets/ssh: Add a flag, `allow_empty_principals` to allow keys or certs to apply to any user/principal. [[GH-28466](https://github.com/hashicorp/vault/pull/28466)]

IMPROVEMENTS:

* audit: Internal implementation changes to the audit subsystem which improve relability. [[GH-28286](https://github.com/hashicorp/vault/pull/28286)]
* ui: Remove deprecated `current_billing_period` from dashboard activity log request [[GH-27559](https://github.com/hashicorp/vault/pull/27559)]

BUG FIXES:

* secret/aws: Fixed potential panic after step-down and the queue has not repopulated. [[GH-28330](https://github.com/hashicorp/vault/pull/28330)]
* auth/cert: During certificate validation, OCSP requests are debug logged even if Vault's log level is above DEBUG. [[GH-28450](https://github.com/hashicorp/vault/pull/28450)]
* auth/cert: ocsp_ca_certificates field was not honored when validating OCSP responses signed by a CA that did not issue the certificate. [[GH-28309](https://github.com/hashicorp/vault/pull/28309)]
* auth: Updated error handling for missing login credentials in AppRole and UserPass auth methods to return a 400 error instead of a 500 error. [[GH-28441](https://github.com/hashicorp/vault/pull/28441)]
* core: Fixed an issue where maximum request duration timeout was not being added to all requests containing strings sys/monitor and sys/events. With this change, timeout is now added to all requests except monitor and events endpoint. [[GH-28230](https://github.com/hashicorp/vault/pull/28230)]
* proxy/cache (enterprise): Fixed a data race that could occur while tracking capabilities in Proxy's static secret cache. [[GH-28494](https://github.com/hashicorp/vault/pull/28494)]
* secrets-sync (enterprise): Secondary nodes in a cluster now properly check activation-flags values.
* secrets-sync (enterprise): Validate corresponding GitHub app parameters `app_name` and `installation_id` are set

## 1.17.5 
## August 30, 2024 

SECURITY: 

core/audit: fix regression where client tokens and token accessors were being
displayed in the audit log in plaintext [HCSEC-2024-18](https://discuss.hashicorp.com/t/hcsec-2024-18-vault-leaks-client-token-and-token-accessor-in-audit-devices)

BUG FIXES:

* proxy/cache (enterprise): Fixed an issue where Proxy with static secret caching enabled would not correctly handle requests to older secret versions for KVv2 secrets. Proxy's static secret cache now properly handles all requests relating to older versions for KVv2 secrets. [[GH-28207](https://github.com/hashicorp/vault/pull/28207)]
* ui: fixes renew-self being called right after login for non-renewable tokens [[GH-28204](https://github.com/hashicorp/vault/pull/28204)]

## 1.17.4
### August 29, 2024

CHANGES:

* activity (enterprise): filter all fields in client count responses by the request namespace [[GH-27790](https://github.com/hashicorp/vault/pull/27790)]
* core: Bump Go version to 1.22.6
* secrets/terraform: Update plugin to v0.9.0 [[GH-28016](https://github.com/hashicorp/vault/pull/28016)]

IMPROVEMENTS:

* activity log: Changes how new client counts in the current month are estimated, in order to return more
visibly sensible totals. [[GH-27547](https://github.com/hashicorp/vault/pull/27547)]
* activity: `/sys/internal/counters/activity` will now include a warning if the specified usage period contains estimated client counts. [[GH-28068](https://github.com/hashicorp/vault/pull/28068)]
* audit: Adds TRACE logging to log request/response under certain circumstances, and further improvements to the audit subsystem. [[GH-28056](https://github.com/hashicorp/vault/pull/28056)]
* cli: `vault operator usage` will now include a warning if the specified usage period contains estimated client counts. [[GH-28068](https://github.com/hashicorp/vault/pull/28068)]
* core/activity: Ensure client count queries that include the current month return consistent results by sorting the clients before performing estimation [[GH-28062](https://github.com/hashicorp/vault/pull/28062)]
* raft-snapshot (enterprise): add support for managed identity credentials for azure snapshots

BUG FIXES:

* activity: The sys/internal/counters/activity endpoint will return current month data when the end_date parameter is set to a future date. [[GH-28042](https://github.com/hashicorp/vault/pull/28042)]
* auth/aws: fixes an issue where not supplying an external id was interpreted as an empty external id [[GH-27858](https://github.com/hashicorp/vault/pull/27858)]
* command: The `vault secrets move` and `vault auth move` command will no longer attempt to write to storage on performance standby nodes. [[GH-28059](https://github.com/hashicorp/vault/pull/28059)]
* core (enterprise): Fix deletion of MFA login-enforcement configurations on standby nodes
* secrets/database: Skip connection verification on reading existing DB connection configuration [[GH-28139](https://github.com/hashicorp/vault/pull/28139)]
* ui: fixes toast (flash) alert message saying "created" when deleting a kv v2 secret [[GH-28093](https://github.com/hashicorp/vault/pull/28093)]

## 1.17.3
### August 07, 2024

CHANGES:

* auth/cf: Update plugin to v0.18.0 [[GH-27724](https://github.com/hashicorp/vault/pull/27724)]

IMPROVEMENTS:

* audit: Ensure that any underyling errors from audit devices are logged even if we consider auditing to be a success. [[GH-27809](https://github.com/hashicorp/vault/pull/27809)]
* audit: Internal implementation changes to the audit subsystem which improve performance. [[GH-27952](https://github.com/hashicorp/vault/pull/27952)]
* audit: sinks (file, socket, syslog) will attempt to log errors to the server operational 
log before returning (if there are errors to log, and the context is done). [[GH-27859](https://github.com/hashicorp/vault/pull/27859)]
* auth/cert: Cache full list of role trust information separately to avoid
eviction, and avoid duplicate loading during multiple simultaneous logins on
the same role. [[GH-27902](https://github.com/hashicorp/vault/pull/27902)]
* license utilization reporting (enterprise): Auto-roll billing start date. [[GH-27656](https://github.com/hashicorp/vault/pull/27656)]
* website/docs: Added API documentation for Azure Secrets Engine delete role [[GH-27883](https://github.com/hashicorp/vault/pull/27883)]

BUG FIXES:

* auth/cert: Use subject's serial number, not issuer's within error message text in OCSP request errors [[GH-27696](https://github.com/hashicorp/vault/pull/27696)]
* core (enterprise): Fix 500 errors that occurred querying `sys/internal/ui/mounts` for a mount prefixed by a namespace path when path filters are configured. [[GH-27939](https://github.com/hashicorp/vault/pull/27939)]
* core/identity: Fixed an issue where deleted/reassigned entity-aliases were not removed from in-memory database. [[GH-27750](https://github.com/hashicorp/vault/pull/27750)]
* proxy/cache (enterprise): Fixed an issue where Proxy would not correctly update KV secrets when talking to a perf standby. Proxy will now attempt to forward requests to update secrets triggered by events to the active node. Note that this requires `allow_forwarding_via_header` to be configured on the cluster. [[GH-27891](https://github.com/hashicorp/vault/pull/27891)]
* proxy/cache (enterprise): Fixed an issue where cached static secrets could fail to update if the secrets belonged to a non-root namespace. [[GH-27730](https://github.com/hashicorp/vault/pull/27730)]
* raft/autopilot: Fixed panic that may occur during shutdown [[GH-27726](https://github.com/hashicorp/vault/pull/27726)]
* secrets-sync (enterprise): Destination set/remove operations will no longer be blocked as "purge in progress" after a purge job ended in failure.
* secrets-sync (enterprise): Normalize custom_tag keys and values for recoverable invalid characters.
* secrets-sync (enterprise): Normalize secret key names before storing the external_name in a secret association.
* secrets-sync (enterprise): Patching github sync destination credentials will properly update and save the new credentials.
* secrets-sync (enterprise): Return an error immediately on destination creation when providing invalid custom_tags based on destination type.
* secrets/identity (enterprise): Fix a bug that can cause DR promotion to fail in rare cases where a PR secondary has inconsistent alias information in storage.
* sys: Fix a bug where mounts of external plugins that were registered before Vault v1.0.0 could not be tuned to
use versioned plugins. [[GH-27881](https://github.com/hashicorp/vault/pull/27881)]
* ui: Fix cursor jump on KVv2 json editor that would occur after pressing ENTER. [[GH-27569](https://github.com/hashicorp/vault/pull/27569)]
* ui: fix issue where enabling then disabling "Tidy ACME" in PKI results in failed API call. [[GH-27742](https://github.com/hashicorp/vault/pull/27742)]
* ui: fix namespace picker not working when in small screen where the sidebar is collapsed by default. [[GH-27728](https://github.com/hashicorp/vault/pull/27728)]


## 1.17.2
### July 10, 2024

CHANGES:

* core: Bump Go version to 1.22.5
* secrets/azure: Update plugin to v0.19.2 [[GH-27652](https://github.com/hashicorp/vault/pull/27652)]

FEATURES:

* **AWS secrets engine STS session tags support**: Adds support for setting STS
session tags when generating temporary credentials using the AWS secrets
engine. [[GH-27620](https://github.com/hashicorp/vault/pull/27620)]

BUG FIXES:

* cli: Fixed issue with `vault hcp connect` where HCP resources with uppercase letters were inaccessible when entering the correct project name. [[GH-27694](https://github.com/hashicorp/vault/pull/27694)]
* core (enterprise): Fix HTTP redirects in namespaces to use the correct path and (in the case of event subscriptions) the correct URI scheme. [[GH-27660](https://github.com/hashicorp/vault/pull/27660)]
* core/config: fix issue when using `proxy_protocol_behavior` with `deny_unauthorized`, 
which causes the Vault TCP listener to close after receiving an untrusted upstream proxy connection. [[GH-27589](https://github.com/hashicorp/vault/pull/27589)]
* core: Fixed an issue with performance standbys not being able to handle rotate root requests. [[GH-27631](https://github.com/hashicorp/vault/pull/27631)]
* secrets/transit (enterprise): Fix an issue that caused input data be returned as part of generated CMAC values.
* ui: Display an error and force a timeout when TOTP passcode is incorrect [[GH-27574](https://github.com/hashicorp/vault/pull/27574)]
* ui: Ensure token expired banner displays when batch token expires [[GH-27479](https://github.com/hashicorp/vault/pull/27479)]

## 1.17.1
### June 26, 2024

CHANGES:

* auth/jwt: Update plugin to v0.21.0 [[GH-27498](https://github.com/hashicorp/vault/pull/27498)]

IMPROVEMENTS:

* storage/raft: Improve autopilot logging on startup to show config values clearly and avoid spurious logs [[GH-27464](https://github.com/hashicorp/vault/pull/27464)]
* ui/secrets-sync: Hide Secrets Sync from the sidebar nav if user does not have access to the feature. [[GH-27262](https://github.com/hashicorp/vault/pull/27262)]

BUG FIXES:

* agent: Fixed an issue causing excessive CPU usage during normal operation [[GH-27518](https://github.com/hashicorp/vault/pull/27518)]
* config: Vault TCP listener config now correctly supports the documented proxy_protocol_behavior 
setting of 'deny_unauthorized' [[GH-27459](https://github.com/hashicorp/vault/pull/27459)]
* core/audit: Audit logging a Vault request/response checks if the existing context 
is cancelled and will now use a new context with a 5 second timeout.
If the existing context is cancelled a new context, will be used. [[GH-27531](https://github.com/hashicorp/vault/pull/27531)]
* helper/pkcs7: Fix parsing certain messages containing only certificates [[GH-27435](https://github.com/hashicorp/vault/pull/27435)]
* proxy: Fixed an issue causing excessive CPU usage during normal operation [[GH-27518](https://github.com/hashicorp/vault/pull/27518)]
* replication (enterprise): fix cache invalidation issue leading to namespace custom metadata not being shown correctly on performance secondaries
* secrets-sync (enterprise): Properly remove tags from secrets in AWS when they are removed from the source association
* secrets-sync (enterprise): Return more accurate error code for invalid connection details
* secrets-sync (enterprise): Skip invalid GitHub repository names when creating destinations
* storage/azure: Fix invalid account name initialization bug [[GH-27563](https://github.com/hashicorp/vault/pull/27563)]
* storage/raft (enterprise): Fix issue with namespace cache not getting cleared on snapshot restore, resulting in namespaces not found in the snapshot being inaccurately represented by API responses. [[GH-27474](https://github.com/hashicorp/vault/pull/27474)]
* ui: Allow creation of session_token type roles for AWS secret backend [[GH-27424](https://github.com/hashicorp/vault/pull/27424)]

## 1.17.0
### June 12, 2024

SECURITY:

* auth/jwt: Update plugin to v0.20.3 that resolves a security issue with validing JWTs [[GH-26890](https://github.com/hashicorp/vault/pull/26890), [HCSEC-2024-11](https://discuss.hashicorp.com/t/hcsec-2024-11-vault-incorrectly-validated-json-web-tokens-jwt-audience-claims/67770)]

CHANGES:

* api: Upgrade from github.com/go-jose/go-jose/v3 v3.0.3 to github.com/go-jose/go-jose/v4 v4.0.1. [[GH-26527](https://github.com/hashicorp/vault/pull/26527)]
* audit: breaking change - Vault now allows audit logs to contain 'correlation-id' and 'x-correlation-id' headers when they 
are present in the incoming request. By default they are not HMAC'ed (but can be configured to HMAC by Vault Operators). [[GH-26777](https://github.com/hashicorp/vault/pull/26777)]
* auth/alicloud: Update plugin to v0.18.0 [[GH-27133](https://github.com/hashicorp/vault/pull/27133)]
* auth/azure: Update plugin to v0.18.0 [[GH-27146](https://github.com/hashicorp/vault/pull/27146)]
* auth/centrify: Remove the deprecated Centrify auth method plugin [[GH-27130](https://github.com/hashicorp/vault/pull/27130)]
* auth/cf: Update plugin to v0.17.0 [[GH-27161](https://github.com/hashicorp/vault/pull/27161)]
* auth/gcp: Update plugin to v0.18.0 [[GH-27140](https://github.com/hashicorp/vault/pull/27140)]
* auth/jwt: Update plugin to v0.20.2 [[GH-26291](https://github.com/hashicorp/vault/pull/26291)]
* auth/kerberos: Update plugin to v0.12.0 [[GH-27177](https://github.com/hashicorp/vault/pull/27177)]
* auth/kubernetes: Update plugin to v0.19.0 [[GH-27186](https://github.com/hashicorp/vault/pull/27186)]
* auth/oci: Update plugin to v0.16.0 [[GH-27142](https://github.com/hashicorp/vault/pull/27142)]
* core (enterprise): Seal High Availability (HA) must be enabled by `enable_multiseal` in configuration.
* core/identity: improve performance for secondary nodes receiving identity related updates through replication [[GH-27184](https://github.com/hashicorp/vault/pull/27184)]
* core: Bump Go version to 1.22.4
* core: return an additional "invalid token" error message in 403 response when the provided request token is expired,
exceeded the number of uses, or is a bogus value [[GH-25953](https://github.com/hashicorp/vault/pull/25953)]
* database/couchbase: Update plugin to v0.11.0 [[GH-27145](https://github.com/hashicorp/vault/pull/27145)]
* database/elasticsearch: Update plugin to v0.15.0 [[GH-27136](https://github.com/hashicorp/vault/pull/27136)]
* database/mongodbatlas: Update plugin to v0.12.0 [[GH-27143](https://github.com/hashicorp/vault/pull/27143)]
* database/redis-elasticache: Update plugin to v0.4.0 [[GH-27139](https://github.com/hashicorp/vault/pull/27139)]
* database/redis: Update plugin to v0.3.0 [[GH-27117](https://github.com/hashicorp/vault/pull/27117)]
* database/snowflake: Update plugin to v0.11.0 [[GH-27132](https://github.com/hashicorp/vault/pull/27132)]
* sdk: String templates now have a maximum size of 100,000 characters. [[GH-26110](https://github.com/hashicorp/vault/pull/26110)]
* secrets/ad: Update plugin to v0.18.0 [[GH-27172](https://github.com/hashicorp/vault/pull/27172)]
* secrets/alicloud: Update plugin to v0.17.0 [[GH-27134](https://github.com/hashicorp/vault/pull/27134)]
* secrets/azure: Update plugin to v0.17.1 [[GH-26528](https://github.com/hashicorp/vault/pull/26528)]
* secrets/azure: Update plugin to v0.19.0 [[GH-27141](https://github.com/hashicorp/vault/pull/27141)]
* secrets/gcp: Update plugin to v0.19.0 [[GH-27164](https://github.com/hashicorp/vault/pull/27164)]
* secrets/gcpkms: Update plugin to v0.17.0 [[GH-27163](https://github.com/hashicorp/vault/pull/27163)]
* secrets/keymgmt (enterprise): Removed `namespace` label on the `vault.kmse.key.count` metric.
* secrets/kmip (enterprise): Update plugin to v0.15.0
* secrets/kubernetes: Update plugin to v0.8.0 [[GH-27187](https://github.com/hashicorp/vault/pull/27187)]
* secrets/kv: Update plugin to v0.18.0 [[GH-26877](https://github.com/hashicorp/vault/pull/26877)]
* secrets/kv: Update plugin to v0.19.0 [[GH-27159](https://github.com/hashicorp/vault/pull/27159)]
* secrets/mongodbatlas: Update plugin to v0.12.0 [[GH-27149](https://github.com/hashicorp/vault/pull/27149)]
* secrets/openldap: Update plugin to v0.13.0 [[GH-27137](https://github.com/hashicorp/vault/pull/27137)]
* secrets/pki: sign-intermediate API will truncate notAfter if calculated to go beyond the signing issuer's notAfter. Previously the notAfter was permitted to go beyond leading to invalid chains. [[GH-26796](https://github.com/hashicorp/vault/pull/26796)]
* secrets/terraform: Update plugin to v0.8.0 [[GH-27147](https://github.com/hashicorp/vault/pull/27147)]
* ui/kubernetes: Update the roles filter-input to use explicit search. [[GH-27178](https://github.com/hashicorp/vault/pull/27178)]
* ui: Update dependencies including D3 libraries [[GH-26346](https://github.com/hashicorp/vault/pull/26346)]
* ui: Upgrade Ember data from 4.11.3 to 4.12.4 [[GH-25272](https://github.com/hashicorp/vault/pull/25272)]
* ui: Upgrade Ember to version 5.4 [[GH-26708](https://github.com/hashicorp/vault/pull/26708)]
* ui: deleting a nested secret will no longer redirect you to the nearest path segment [[GH-26845](https://github.com/hashicorp/vault/pull/26845)]
* ui: flash messages render on right side of page [[GH-25459](https://github.com/hashicorp/vault/pull/25459)]

FEATURES:

* **PKI Certificate Metadata (enterprise)**:  Add Certificate Metadata Functionality to Record and Return Client Information about a Certificate.
* **Adaptive Overload Protection (enterprise)**: Adds Adaptive Overload Protection
for write requests as a Beta feature (disabled by default). This automatically
prevents overloads caused by too many write requests while maintaining optimal
throughput for the hardware configuration and workload.
* **Audit Filtering (enterprise)** : Audit devices support expression-based filter rules (powered by go-bexpr) to determine which entries are written to the audit log.
* **LDAP Secrets engine hierarchical path support**: Hierarchical path handling is now supported for role and set APIs. [[GH-27203](https://github.com/hashicorp/vault/pull/27203)]
* **Plugin Identity Tokens**: Adds secret-less configuration of AWS auth engine using web identity federation. [[GH-26507](https://github.com/hashicorp/vault/pull/26507)]
* **Plugin Workload Identity** (enterprise): Vault can generate identity tokens for plugins to use in workload identity federation auth flows.
* **Transit AES-CMAC (enterprise)**: Added support to create and verify AES backed cipher-based message authentication codes

IMPROVEMENTS:

* activity (enterprise): Change minimum retention window in activity log to 48 months
* agent: Added a new config option, `lease_renewal_threshold`, that controls the refresh rate of non-renewable leases in Agent's template engine. [[GH-25212](https://github.com/hashicorp/vault/pull/25212)]
* agent: Agent will re-trigger auto auth if token used for rendering templates has been revoked, has exceeded the number of uses, or is a bogus value. [[GH-26172](https://github.com/hashicorp/vault/pull/26172)]
* api: Move CLI token helper functions to importable packages in `api` module. [[GH-25744](https://github.com/hashicorp/vault/pull/25744)]
* audit: timestamps across multiple audit devices for an audit entry will now match. [[GH-26088](https://github.com/hashicorp/vault/pull/26088)]
* auth/aws: Add inferred_hostname metadata for IAM AWS authentication method. [[GH-25418](https://github.com/hashicorp/vault/pull/25418)]
* auth/aws: add canonical ARN as entity alias option [[GH-22460](https://github.com/hashicorp/vault/pull/22460)]
* auth/aws: add support for external_ids in AWS assume-role [[GH-26628](https://github.com/hashicorp/vault/pull/26628)]
* auth/cert: Adds support for TLS certificate authenticaion through a reverse proxy that terminates the SSL connection [[GH-17272](https://github.com/hashicorp/vault/pull/17272)]
* cli: Add events subscriptions commands
* command/server: Removed environment variable requirement to generate pprof 
files using SIGUSR2. Added CPU profile support. [[GH-25391](https://github.com/hashicorp/vault/pull/25391)]
* core (enterprise): persist seal rewrap status, so rewrap status API is consistent on secondary nodes.
* core/activity: Include ACME client metrics to precomputed queries [[GH-26519](https://github.com/hashicorp/vault/pull/26519)]
* core/activity: Include ACME clients in activity log responses [[GH-26020](https://github.com/hashicorp/vault/pull/26020)]
* core/activity: Include ACME clients in vault operator usage response [[GH-26525](https://github.com/hashicorp/vault/pull/26525)]
* core/config: reload service registration configuration on SIGHUP [[GH-17598](https://github.com/hashicorp/vault/pull/17598)]
* core: add deadlock detection in barrier and sealwrap
* license utilization reporting (enterprise): Add retention months to license utilization reports.
* proxy/cache (enterprise): Support new configuration parameter for static secret caching, `static_secret_token_capability_refresh_behavior`, to control the behavior when the capability refresh request receives an error from Vault.
* proxy: Proxy will re-trigger auto auth if the token used for requests has been revoked, has exceeded the number of uses,
or is an otherwise invalid value. [[GH-26307](https://github.com/hashicorp/vault/pull/26307)]
* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20221104090112-13395acd02c5
* replication (enterprise): Add replication heartbeat metric to telemetry
* replication (enterprise): Periodically write current time on the primary to storage, use that downstream to measure replication lag in time, expose that in health and replication status endpoints. [[GH-26406](https://github.com/hashicorp/vault/pull/26406)]
* sdk/decompression: DecompressWithCanary will now chunk the decompression in memory to prevent loading it all at once. [[GH-26464](https://github.com/hashicorp/vault/pull/26464)]
* sdk/helper/testcluster: add some new helpers, improve some error messages. [[GH-25329](https://github.com/hashicorp/vault/pull/25329)]
* sdk/helper/testhelpers: add namespace helpers [[GH-25270](https://github.com/hashicorp/vault/pull/25270)]
* secrets-sync (enterprise): Added global config path to the administrative namespace.
* secrets/pki (enterprise): Disable warnings about unknown parameters to the various CIEPS endpoints
* secrets/pki: Add a new ACME configuration parameter that allows increasing the maximum TTL for ACME leaf certificates [[GH-26797](https://github.com/hashicorp/vault/pull/26797)]
* secrets/transform (enterprise): Add delete by token and delete by plaintext operations to Tokenization.
* storage/azure: Perform validation on Azure account name and container name [[GH-26135](https://github.com/hashicorp/vault/pull/26135)]
* storage/raft (enterprise): add support for separate entry size limit for mount
and namespace table paths in storage to allow increased mount table size without
allowing other user storage entries to become larger. [[GH-25992](https://github.com/hashicorp/vault/pull/25992)]
* storage/raft: panic on unknown Raft operations [[GH-25991](https://github.com/hashicorp/vault/pull/25991)]
* ui (enterprise): Allow HVD users to access Secrets Sync. [[GH-26841](https://github.com/hashicorp/vault/pull/26841)]
* ui (enterprise): Update dashboard to make activity log query using the same start time as the metrics overview [[GH-26729](https://github.com/hashicorp/vault/pull/26729)]
* ui (enterprise): Update filters on the custom messages list view. [[GH-26653](https://github.com/hashicorp/vault/pull/26653)]
* ui: Allow users to wrap inputted data again instead of resetting form [[GH-27289](https://github.com/hashicorp/vault/pull/27289)]
* ui: Display ACME clients on a separate page in the UI. [[GH-26020](https://github.com/hashicorp/vault/pull/26020)]
* ui: Hide dashboard client count card if user does not have permission to view clients. [[GH-26848](https://github.com/hashicorp/vault/pull/26848)]
* ui: Show computed values from `sys/internal/ui/mounts` endpoint for auth mount configuration view [[GH-26663](https://github.com/hashicorp/vault/pull/26663)]
* ui: Update PGP display and show error for Generate Operation Token flow with PGP [[GH-26993](https://github.com/hashicorp/vault/pull/26993)]
* ui: Update language in Transit secret engine to reflect that not all keys are for encyryption [[GH-27346](https://github.com/hashicorp/vault/pull/27346)]
* ui: Update userpass user form to allow setting `password_hash` field. [[GH-26577](https://github.com/hashicorp/vault/pull/26577)]
* ui: fixes cases where inputs did not have associated labels [[GH-26263](https://github.com/hashicorp/vault/pull/26263)]
* ui: show banner instead of permission denied error when batch token is expired [[GH-26396](https://github.com/hashicorp/vault/pull/26396)]
* website/docs: Add note about eventual consietency with the MongoDB Atlas database secrets engine [[GH-24152](https://github.com/hashicorp/vault/pull/24152)]

DEPRECATIONS:

* Request Limiter Beta(enterprise): This Beta feature added in 1.16 has been
superseded by Adaptive Overload Protection and will be removed.
* secrets/azure: Deprecate field "password_policy" as we are not able to set it anymore with the new MS Graph API. [[GH-25637](https://github.com/hashicorp/vault/pull/25637)]

BUG FIXES:

* activity (enterprise): fix read-only storage error on upgrades
* agent: Correctly constructs kv-v2 secret paths in nested namespaces. [[GH-26863](https://github.com/hashicorp/vault/pull/26863)]
* agent: Fixes a high Vault load issue, by restarting the Conusl template server after backing off instead of immediately. [[GH-25497](https://github.com/hashicorp/vault/pull/25497)]
* agent: `vault.namespace` no longer gets incorrectly overridden by `auto_auth.namespace`, if set [[GH-26427](https://github.com/hashicorp/vault/pull/26427)]
* api: fixed a bug where LifetimeWatcher routines weren't respecting exponential backoff in the presence of unexpected errors [[GH-26383](https://github.com/hashicorp/vault/pull/26383)]
* audit: Operator changes to configured audit headers (via `/sys/config/auditing`) 
will now force invalidation and be reloaded from storage when data is replicated 
to other nodes.
* auth/ldap: Fix login error for group search anonymous bind. [[GH-26200](https://github.com/hashicorp/vault/pull/26200)]
* auth/ldap: Fix login error missing entity alias attribute value. [[GH-26200](https://github.com/hashicorp/vault/pull/26200)]
* auto-auth: Addressed issue where having no permissions to renew a renewable token caused auto-auth to attempt to renew constantly with no backoff [[GH-26844](https://github.com/hashicorp/vault/pull/26844)]
* cli/debug: Fix resource leak in CLI debug command. [[GH-26167](https://github.com/hashicorp/vault/pull/26167)]
* cli: fixed a bug where the Vault CLI would error out if 
HOME was not set. [[GH-26243](https://github.com/hashicorp/vault/pull/26243)]
* core (enterprise): Fix 403s returned when forwarding invalid token to active node from secondary.
* core (enterprise): Fix an issue that prevented the seal re-wrap status from reporting that a re-wrap is in progress for up to a second.
* core (enterprise): fix bug where raft followers disagree with the seal type after returning to one seal from two. [[GH-26523](https://github.com/hashicorp/vault/pull/26523)]
* core (enterprise): fix issue where the Seal HA rewrap system may remain running when an active node steps down.
* core/audit: Audit logging a Vault request/response will now use a minimum 5 second context timeout. 
If the existing context deadline occurs later than 5s in the future, it will be used, otherwise a 
new context, separate from the original will be used. [[GH-26616](https://github.com/hashicorp/vault/pull/26616)]
* core/metrics: store cluster name in unencrypted storage to prevent blank cluster name [[GH-26878](https://github.com/hashicorp/vault/pull/26878)]
* core/namespace (enterprise): Privileged namespace paths provided in the `administrative_namespace_path` config will now be canonicalized.
* core/seal: During a seal reload through SIGHUP, only write updated seal barrier on an active node [[GH-26381](https://github.com/hashicorp/vault/pull/26381)]
* core/seal: allow overriding of VAULT_GCPCKMS_SEAL_KEY_RING and VAULT_GCPCKMS_SEAL_CRYPTO_KEY environment keys in seal-ha
* core: Add missing field delegated_auth_accessors to GET /sys/mounts/:path API response [[GH-26876](https://github.com/hashicorp/vault/pull/26876)]
* core: Address a data race updating a seal's last seen healthy time attribute [[GH-27014](https://github.com/hashicorp/vault/pull/27014)]
* core: Fix `redact_version` listener parameter being ignored for some OpenAPI related endpoints. [[GH-26607](https://github.com/hashicorp/vault/pull/26607)]
* core: Only reload seal configuration when enable_multiseal is set to true. [[GH-26166](https://github.com/hashicorp/vault/pull/26166)]
* core: when listener configuration `chroot_namespace` is active, Vault will no longer report that the configuration is invalid when Vault is sealed
* events (enterprise): Fix bug preventing subscribing and receiving events within a namepace.
* events (enterprise): Terminate WebSocket connection when token is revoked.
* openapi: Fixing approle reponse duration types [[GH-25510](https://github.com/hashicorp/vault/pull/25510)]
* openapi: added the missing migrate parameter for the unseal endpoint in vault/logical_system_paths.go [[GH-25550](https://github.com/hashicorp/vault/pull/25550)]
* pki: Fix error in cross-signing using ed25519 keys [[GH-27093](https://github.com/hashicorp/vault/pull/27093)]
* plugin/wif: fix a bug where the namespace was not set for external plugins using workload identity federation [[GH-26384](https://github.com/hashicorp/vault/pull/26384)]
* replication (enterprise): fix "given mount path is not in the same namespace as the request" error that can occur when enabling replication for the first time on a secondary cluster
* replication (enterprise): fixed data integrity issue with the processing of identity aliases causing duplicates to occur in rare cases
* router: Fix missing lock in MatchingSystemView. [[GH-25191](https://github.com/hashicorp/vault/pull/25191)]
* secret/database: Fixed race condition where database mounts may leak connections [[GH-26147](https://github.com/hashicorp/vault/pull/26147)]
* secrets-sync (enterprise): Fixed an issue with syncing to target projects in GCP
* secrets/azure: Update vault-plugin-secrets-azure to 0.17.2 to include a bug fix for azure role creation [[GH-26896](https://github.com/hashicorp/vault/pull/26896)]
* secrets/pki (enterprise): cert_role parameter within authenticators.cert EST configuration handler could not be set
* secrets/pki: fixed validation bug which rejected ldap schemed URLs in crl_distribution_points. [[GH-26477](https://github.com/hashicorp/vault/pull/26477)]
* secrets/transform (enterprise): Fix a bug preventing the use of alternate schemas on PostgreSQL token stores.
* secrets/transit: Use 'hash_algorithm' parameter if present in HMAC verify requests. Otherwise fall back to deprecated 'algorithm' parameter. [[GH-27211](https://github.com/hashicorp/vault/pull/27211)]
* storage/raft (enterprise): Fix a bug where autopilot automated upgrades could fail due to using the wrong upgrade version
* storage/raft (enterprise): Fix a regression introduced in 1.15.8 that causes
autopilot to fail to discover new server versions and so not trigger an upgrade. [[GH-27277](https://github.com/hashicorp/vault/pull/27277)]
* storage/raft: prevent writes from impeding leader transfers, e.g. during automated upgrades [[GH-25390](https://github.com/hashicorp/vault/pull/25390)]
* transform (enterprise): guard against a panic looking up a token in exportable mode with barrier storage.
* ui: Do not show resultant-ACL banner when ancestor namespace grants wildcard access. [[GH-27263](https://github.com/hashicorp/vault/pull/27263)]
* ui: Fix KVv2 cursor jumping inside json editor after initial input. [[GH-27120](https://github.com/hashicorp/vault/pull/27120)]
* ui: Fix KVv2 json editor to allow null values. [[GH-27094](https://github.com/hashicorp/vault/pull/27094)]
* ui: Fix a bug where disabling TTL on the AWS credential form would still send TTL value [[GH-27366](https://github.com/hashicorp/vault/pull/27366)]
* ui: Fix broken help link in console for the web command. [[GH-26858](https://github.com/hashicorp/vault/pull/26858)]
* ui: Fix configuration link from Secret Engine list view for Ember engines. [[GH-27131](https://github.com/hashicorp/vault/pull/27131)]
* ui: Fix link to v2 generic secrets engine from secrets list page. [[GH-27019](https://github.com/hashicorp/vault/pull/27019)]
* ui: Prevent perpetual loading screen when Vault needs initialization [[GH-26985](https://github.com/hashicorp/vault/pull/26985)]
* ui: Refresh model within a namespace on the Secrets Sync overview page. [[GH-26790](https://github.com/hashicorp/vault/pull/26790)]
* ui: Remove possibility of returning an undefined timezone from date-format helper [[GH-26693](https://github.com/hashicorp/vault/pull/26693)]
* ui: Resolved accessibility issues with Web REPL. Associated label and help text with input, added a conditional to show the console/ui-panel only when toggled open, added keyboard focus trap. [[GH-26872](https://github.com/hashicorp/vault/pull/26872)]
* ui: fix issue where a month without new clients breaks the client count dashboard [[GH-27352](https://github.com/hashicorp/vault/pull/27352)]
* ui: fixed a bug where the replication pages did not update display when navigating between DR and performance [[GH-26325](https://github.com/hashicorp/vault/pull/26325)]
* ui: fixes undefined start time in filename for downloaded client count attribution csv [[GH-26485](https://github.com/hashicorp/vault/pull/26485)]

## 1.16.14 Enterprise
### December 18, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* secrets/pki: Enforce the issuer constraint extensions (extended key usage, name constraints, issuer name) when issuing or signing leaf certificates. For more information see [PKI considerations](https://developer.hashicorp.com/vault/docs/secrets/pki/considerations#issuer-constraints-enforcement) [[GH-29045](https://github.com/hashicorp/vault/pull/29045)]

IMPROVEMENTS:

* auth/okta: update to okta sdk v5 from v2. Transitively updates go-jose dependency to >=3.0.3 to resolve GO-2024-2631. See https://github.com/okta/okta-sdk-golang/blob/master/MIGRATING.md for details on changes. [[GH-28121](https://github.com/hashicorp/vault/pull/28121)]
* core: Added new `enable_post_unseal_trace` and `post_unseal_trace_directory` config options to generate Go traces during the post-unseal step for debug purposes. [[GH-28895](https://github.com/hashicorp/vault/pull/28895)]
* sdk: Add Vault build date to system view plugin environment response [[GH-29082](https://github.com/hashicorp/vault/pull/29082)]
* ui: Replace KVv2 json secret details view with Hds::CodeBlock component allowing users to search the full secret height. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]

BUG FIXES:

* autosnapshots (enterprise): Fix an issue where snapshot size metrics were not reported for cloud-based storage.
* core/metrics: Fix unlocked mounts read for usage reporting. [[GH-29091](https://github.com/hashicorp/vault/pull/29091)]
* core/seal (enterprise): Fix decryption of the raft bootstrap challenge when using seal high availability. [[GH-29117](https://github.com/hashicorp/vault/pull/29117)]
* secret/db: Update static role rotation to generate a new password after 2 failed attempts. [[GH-28989](https://github.com/hashicorp/vault/pull/28989)]
* ui: Allow users to search the full json object within the json code-editor edit/create view. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]
* ui: Decode `connection_url` to fix database connection updates (i.e. editing connection config, deleting roles) failing when urls include template variables. [[GH-29114](https://github.com/hashicorp/vault/pull/29114)]
* vault/diagnose: Fix time to expiration reporting within the TLS verification to not be a month off. [[GH-29128](https://github.com/hashicorp/vault/pull/29128)]

## 1.16.13 Enterprise
### November 21, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20241115202008-166203013d8e

CHANGES:

* activity log: Deprecated the field "default_report_months". Instead, the billing start time will be used to determine the start time
when querying the activity log endpoints. [[GH-27350](https://github.com/hashicorp/vault/pull/27350)]
* core/ha (enterprise): Failed attempts to become a performance standby node are now using an exponential backoff instead of a
10 second delay in between retries. The backoff starts at 2s and increases by a factor of two until reaching
the maximum of 16s. This should make unsealing of the node faster in some cases.
* login (enterprise): Return a 500 error during logins when performance standby nodes make failed gRPC requests to the active node. [[GH-28807](https://github.com/hashicorp/vault/pull/28807)]

FEATURES:

* **Product Usage Reporting**: Added product usage reporting, which collects anonymous, numerical, non-sensitive data about Vault secrets usage, and adds it to the existing utilization reports. See the [[docs](https://developer.hashicorp.com/vault/docs/enterprise/license/product-usage-reporting)] for more info [[GH-28858](https://github.com/hashicorp/vault/pull/28858)]

IMPROVEMENTS:

* raft-snapshot (enterprise): add support for managed identity credentials for azure snapshots
* secrets-sync (enterprise): No longer attempt to unsync a random UUID secret name in GCP upon destination creation.

BUG FIXES:

* auth/ldap: Fixed an issue where debug level logging was not emitted. [[GH-28881](https://github.com/hashicorp/vault/pull/28881)]
* core: Improved an internal helper function that sanitizes paths by adding a check for leading backslashes 
in addition to the existing check for leading slashes. [[GH-28878](https://github.com/hashicorp/vault/pull/28878)]
* secret/pki: Fix a bug that prevents PKI issuer field enable_aia_url_templating
to be set to false. [[GH-28832](https://github.com/hashicorp/vault/pull/28832)]
* secrets-sync (enterprise): Fixed issue where secret-key granularity destinations could sometimes cause a panic when loading a sync status.
* secrets/aws: Fix issue with static credentials not rotating after restart or leadership change. [[GH-28775](https://github.com/hashicorp/vault/pull/28775)]
* secrets/ssh: Return the flag `allow_empty_principals` in the read role api when key_type is "ca" [[GH-28901](https://github.com/hashicorp/vault/pull/28901)]
* secrets/transform (enterprise): Fix nil panic when accessing a partially setup database store.
* secrets/transit: Fix a race in which responses from the key update api could contain results from another subsequent update [[GH-28839](https://github.com/hashicorp/vault/pull/28839)]

## 1.16.12 Enterprise
### October 30, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:
* core/raft: Add raft join limits [[GH-28790](https://github.com/hashicorp/vault/pull/28790), [HCSEC-2024-26](https://discuss.hashicorp.com/t/hcsec-2024-26-vault-vulnerable-to-denial-of-service-through-memory-exhaustion-when-processing-raft-cluster-join-requests)]
CHANGES:

* auth/azure: Update plugin to v0.17.1
* secrets/openldap: Update plugin to v0.12.2

IMPROVEMENTS:

* ui: Add button to copy secret path in kv v1 and v2 secrets engines [[GH-28629](https://github.com/hashicorp/vault/pull/28629)]
* ui: Adds copy button to identity entity, alias and mfa method IDs [[GH-28742](https://github.com/hashicorp/vault/pull/28742)]

BUG FIXES:

* core/seal (enterprise): Fix bug that caused seal generation information to be replicated, which prevented disaster recovery and performance replication clusters from using their own seal high-availability configuration.
* core/seal: Fix an issue that could cause reading from sys/seal-backend-status to return stale information. [[GH-28631](https://github.com/hashicorp/vault/pull/28631)]
* core: Fixed panic seen when performing help requests without /v1/ in the URL. [[GH-28669](https://github.com/hashicorp/vault/pull/28669)]
* namespaces (enterprise): Fix issue where namespace patch requests to a performance secondary would not patch the namespace's metadata.
* secrets/pki: Address issue with ACME HTTP-01 challenges failing for IPv6 IPs due to improperly formatted URLs [[GH-28718](https://github.com/hashicorp/vault/pull/28718)]
* ui: No longer running decodeURIComponent on KVv2 list view allowing percent encoded data-octets in path name. [[GH-28698](https://github.com/hashicorp/vault/pull/28698)]
  
## 1.16.11 Enterprise
### October 09, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* secrets/identity: A privileged Vault operator with write permissions to the root namespace's identity endpoint could escalate their privileges to Vault's root policy (CVE-2024-9180) [HCSEC-2024-21](https://discuss.hashicorp.com/t/hcsec-2024-21-vault-operators-in-root-namespace-may-elevate-their-privileges/70565)

IMPROVEMENTS:

* core: log at level ERROR rather than INFO when all seals are unhealthy. [[GH-28564](https://github.com/hashicorp/vault/pull/28564)]
* physical/raft: Log when the MAP_POPULATE mmap flag gets disabled before opening the database. [[GH-28526](https://github.com/hashicorp/vault/pull/28526)]

BUG FIXES:

* auth/cert: When using ocsp_ca_certificates, an error was produced though extra certs validation succeeded. [[GH-28597](https://github.com/hashicorp/vault/pull/28597)]
* auth/token: Fix token TTL calculation so that it uses `max_lease_ttl` tune value for tokens created via `auth/token/create`. [[GH-28498](https://github.com/hashicorp/vault/pull/28498)]
* databases: fix issue where local timezone was getting lost when using a rotation schedule cron [[GH-28509](https://github.com/hashicorp/vault/pull/28509)]
* secrets-sync (enterprise): Fix KV secret access sometimes being denied, due to a double forward-slash (`//`) in the mount path, when the token should otherwise have access.

## 1.16.10 Enterprise
### September 25, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:
* secrets/ssh: require `valid_principals` to contain a value or `default_user` be set by default to guard against potentially insecure configurations. `allow_empty_principals` can be used for backwards compatibility [HCSEC-2024-20](https://discuss.hashicorp.com/t/hcsec-2024-20-vault-ssh-secrets-engine-configuration-did-not-restrict-valid-principals-by-default/7025

CHANGES:

* core: Bump Go version to 1.22.7.
* secrets/ssh: Add a flag, `allow_empty_principals` to allow keys or certs to apply to any user/principal. [[GH-28466](https://github.com/hashicorp/vault/pull/28466)]

IMPROVEMENTS:

* audit: Internal implementation changes to the audit subsystem which improve relability. [[GH-28286](https://github.com/hashicorp/vault/pull/28286)]
* ui: Remove deprecated `current_billing_period` from dashboard activity log request [[GH-27559](https://github.com/hashicorp/vault/pull/27559)]

BUG FIXES:

* secret/aws: Fixed potential panic after step-down and the queue has not repopulated. [[GH-28330](https://github.com/hashicorp/vault/pull/28330)]
* auth/cert: During certificate validation, OCSP requests are debug logged even if Vault's log level is above DEBUG. [[GH-28450](https://github.com/hashicorp/vault/pull/28450)]
* auth/cert: ocsp_ca_certificates field was not honored when validating OCSP responses signed by a CA that did not issue the certificate. [[GH-28309](https://github.com/hashicorp/vault/pull/28309)]
* auth: Updated error handling for missing login credentials in AppRole and UserPass auth methods to return a 400 error instead of a 500 error. [[GH-28441](https://github.com/hashicorp/vault/pull/28441)]
* core: Fixed an issue where maximum request duration timeout was not being added to all requests containing strings sys/monitor and sys/events. With this change, timeout is now added to all requests except monitor and events endpoint. [[GH-28230](https://github.com/hashicorp/vault/pull/28230)]
* proxy/cache (enterprise): Fixed a data race that could occur while tracking capabilities in Proxy's static secret cache. [[GH-28494](https://github.com/hashicorp/vault/pull/28494)]
* secrets-sync (enterprise): Validate corresponding GitHub app parameters `app_name` and `installation_id` are set

## 1.16.9 Enterprise
### August 30, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY: 

core/audit: fix regression where client tokens and token accessors were being
displayed in the audit log in plaintext [HCSEC-2024-18](https://discuss.hashicorp.com/t/hcsec-2024-18-vault-leaks-client-token-and-token-accessor-in-audit-devices)

BUG FIXES:

* proxy/cache (enterprise): Fixed an issue where Proxy with static secret caching enabled would not correctly handle requests to older secret versions for KVv2 secrets. Proxy's static secret cache now properly handles all requests relating to older versions for KVv2 secrets. [[GH-28207](https://github.com/hashicorp/vault/pull/28207)]
## 1.16.8 Enterprise
### August 29, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* activity (enterprise): filter all fields in client count responses by the request namespace [[GH-27790](https://github.com/hashicorp/vault/pull/27790)]
* core: Bump Go version to 1.22.6

IMPROVEMENTS:

* activity log: Changes how new client counts in the current month are estimated, in order to return more
visibly sensible totals. [[GH-27547](https://github.com/hashicorp/vault/pull/27547)]
* activity: `/sys/internal/counters/activity` will now include a warning if the specified usage period contains estimated client counts. [[GH-28068](https://github.com/hashicorp/vault/pull/28068)]
* audit: Adds TRACE logging to log request/response under certain circumstances, and further improvements to the audit subsystem. [[GH-28056](https://github.com/hashicorp/vault/pull/28056)]
* cli: `vault operator usage` will now include a warning if the specified usage period contains estimated client counts. [[GH-28068](https://github.com/hashicorp/vault/pull/28068)]
* core/activity: Ensure client count queries that include the current month return consistent results by sorting the clients before performing estimation [[GH-28062](https://github.com/hashicorp/vault/pull/28062)]

BUG FIXES:

* activity: The sys/internal/counters/activity endpoint will return current month data when the end_date parameter is set to a future date. [[GH-28042](https://github.com/hashicorp/vault/pull/28042)]
* command: The `vault secrets move` and `vault auth move` command will no longer attempt to write to storage on performance standby nodes. [[GH-28059](https://github.com/hashicorp/vault/pull/28059)]
* core (enterprise): Fix deletion of MFA login-enforcement configurations on standby nodes
* secrets/database: Skip connection verification on reading existing DB connection configuration [[GH-28139](https://github.com/hashicorp/vault/pull/28139)]
* ui: fixes toast (flash) alert message saying "created" when deleting a kv v2 secret [[GH-28093](https://github.com/hashicorp/vault/pull/28093)]

## 1.16.7 Enterprise
### August 07, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* auth/cf: Update plugin to v0.18.0 [[GH-27724](https://github.com/hashicorp/vault/pull/27724)]

IMPROVEMENTS:

* audit: Ensure that any underyling errors from audit devices are logged even if we consider auditing to be a success. [[GH-27809](https://github.com/hashicorp/vault/pull/27809)]
* audit: Internal implementation changes to the audit subsystem which improve performance. [[GH-27952](https://github.com/hashicorp/vault/pull/27952)]
* audit: sinks (file, socket, syslog) will attempt to log errors to the server operational 
log before returning (if there are errors to log, and the context is done). [[GH-27859](https://github.com/hashicorp/vault/pull/27859)]
* auth/cert: Cache full list of role trust information separately to avoid
eviction, and avoid duplicate loading during multiple simultaneous logins on
the same role. [[GH-27902](https://github.com/hashicorp/vault/pull/27902)]
* license utilization reporting (enterprise): Auto-roll billing start date. [[GH-27656](https://github.com/hashicorp/vault/pull/27656)]

BUG FIXES:

* auth/cert: Use subject's serial number, not issuer's within error message text in OCSP request errors [[GH-27696](https://github.com/hashicorp/vault/pull/27696)]
* cli: Fixed issue with `vault hcp connect` where HCP resources with uppercase letters were inaccessible when entering the correct project name. [[GH-27694](https://github.com/hashicorp/vault/pull/27694)]
* core (enterprise): Fix 500 errors that occurred querying `sys/internal/ui/mounts` for a mount prefixed by a namespace path when path filters are configured. [[GH-27939](https://github.com/hashicorp/vault/pull/27939)]
* core/identity: Fixed an issue where deleted/reassigned entity-aliases were not removed from in-memory database. [[GH-27750](https://github.com/hashicorp/vault/pull/27750)]
* proxy/cache (enterprise): Fixed an issue where Proxy would not correctly update KV secrets when talking to a perf standby. Proxy will now attempt to forward requests to update secrets triggered by events to the active node. Note that this requires `allow_forwarding_via_header` to be configured on the cluster. [[GH-27891](https://github.com/hashicorp/vault/pull/27891)]
* raft/autopilot: Fixed panic that may occur during shutdown [[GH-27726](https://github.com/hashicorp/vault/pull/27726)]
* secrets-sync (enterprise): Destination set/remove operations will no longer be blocked as "purge in progress" after a purge job ended in failure.
* secrets-sync (enterprise): Normalize custom_tag keys and values for recoverable invalid characters.
* secrets-sync (enterprise): Normalize secret key names before storing the external_name in a secret association.
* secrets-sync (enterprise): Patching github sync destination credentials will properly update and save the new credentials.
* secrets-sync (enterprise): Return an error immediately on destination creation when providing invalid custom_tags based on destination type.
* secrets/identity (enterprise): Fix a bug that can cause DR promotion to fail in rare cases where a PR secondary has inconsistent alias information in storage.
* sys: Fix a bug where mounts of external plugins that were registered before Vault v1.0.0 could not be tuned to
use versioned plugins. [[GH-27881](https://github.com/hashicorp/vault/pull/27881)]
* ui: Fix cursor jump on KVv2 json editor that would occur after pressing ENTER. [[GH-27569](https://github.com/hashicorp/vault/pull/27569)]
* ui: fix issue where enabling then disabling "Tidy ACME" in PKI results in failed API call. [[GH-27742](https://github.com/hashicorp/vault/pull/27742)]
* ui: fix namespace picker not working when in small screen where the sidebar is collapsed by default. [[GH-27728](https://github.com/hashicorp/vault/pull/27728)]


## 1.16.6 Enterprise
### July 10, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* core: Bump Go version to 1.22.5.
* auth/jwt: Revert [GH-295](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/295) which changed the way JWT `aud` claims were validated.

BUG FIXES:

* agent: Correctly constructs kv-v2 secret paths in nested namespaces. [[GH-26863](https://github.com/hashicorp/vault/pull/26863)]
* core (enterprise): Fix HTTP redirects in namespaces to use the correct path and (in the case of event subscriptions) the correct URI scheme. [[GH-27660](https://github.com/hashicorp/vault/pull/27660)]
* core/config: fix issue when using `proxy_protocol_behavior` with `deny_unauthorized`, 
which causes the Vault TCP listener to close after receiving an untrusted upstream proxy connection. [[GH-27589](https://github.com/hashicorp/vault/pull/27589)]
* core: Fixed an issue with performance standbys not being able to handle rotate root requests. [[GH-27631](https://github.com/hashicorp/vault/pull/27631)]
* ui: Display an error and force a timeout when TOTP passcode is incorrect [[GH-27574](https://github.com/hashicorp/vault/pull/27574)]
* ui: Ensure token expired banner displays when batch token expires [[GH-27479](https://github.com/hashicorp/vault/pull/27479)]
  
## 1.16.5 Enterprise
### June 26, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

BUG FIXES:

* cli/debug: Fix resource leak in CLI debug command. [[GH-26167](https://github.com/hashicorp/vault/pull/26167)]
* config: Vault TCP listener config now correctly supports the documented proxy_protocol_behavior 
setting of 'deny_unauthorized' [[GH-27459](https://github.com/hashicorp/vault/pull/27459)]
* core/audit: Audit logging a Vault request/response checks if the existing context 
is cancelled and will now use a new context with a 5 second timeout.
If the existing context is cancelled a new context, will be used. [[GH-27531](https://github.com/hashicorp/vault/pull/27531)]
* helper/pkcs7: Fix parsing certain messages containing only certificates [[GH-27435](https://github.com/hashicorp/vault/pull/27435)]
* replication (enterprise): fix cache invalidation issue leading to namespace custom metadata not being shown correctly on performance secondaries
* secrets-sync (enterprise): Properly remove tags from secrets in AWS when they are removed from the source association
* secrets-sync (enterprise): Return more accurate error code for invalid connection details
* secrets-sync (enterprise): Skip invalid GitHub repository names when creating destinations
* storage/raft (enterprise): Fix issue with namespace cache not getting cleared on snapshot restore, resulting in namespaces not found in the snapshot being inaccurately represented by API responses. [[GH-27474](https://github.com/hashicorp/vault/pull/27474)]
* ui: Allow creation of session_token type roles for AWS secret backend [[GH-27424](https://github.com/hashicorp/vault/pull/27424)]

## 1.16.4 Enterprise 
### June 12, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* core: Bump Go version to 1.22.4.
* ui/kubernetes: Update the roles filter-input to use explicit search. [[GH-27178](https://github.com/hashicorp/vault/pull/27178)]

IMPROVEMENTS:

* ui: Allow users to wrap inputted data again instead of resetting form [[GH-27289](https://github.com/hashicorp/vault/pull/27289)]
* ui: Update language in Transit secret engine to reflect that not all keys are for encyryption [[GH-27346](https://github.com/hashicorp/vault/pull/27346)]

BUG FIXES:

* secrets/transform (enterprise): Fix a bug preventing the use of alternate schemas on PostgreSQL token stores.
* storage/raft (enterprise): Fix a regression introduced in 1.15.8 that causes
autopilot to fail to discover new server versions and so not trigger an upgrade. [[GH-27277](https://github.com/hashicorp/vault/pull/27277)]
* ui: Do not show resultant-ACL banner when ancestor namespace grants wildcard access. [[GH-27263](https://github.com/hashicorp/vault/pull/27263)]
* ui: Fix a bug where disabling TTL on the AWS credential form would still send TTL value [[GH-27366](https://github.com/hashicorp/vault/pull/27366)]
* ui: fix issue where a month with total clients but no new clients breaks the client count dashboard [[GH-5962](https://github.com/hashicorp/vault/pull/5962)]

## 1.16.3 
### May 30, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* auth/jwt: Update plugin to v0.20.3 that resolves a security issue with validing JWTs [[GH-26890](https://github.com/hashicorp/vault/pull/26890), [HCSEC-2024-11](https://discuss.hashicorp.com/t/hcsec-2024-11-vault-incorrectly-validated-json-web-tokens-jwt-audience-claims/67770)]

CHANGES:

* core/identity: improve performance for secondary nodes receiving identity related updates through replication [[GH-27184](https://github.com/hashicorp/vault/pull/27184)]
* core: Bump Go version to 1.22.2.

IMPROVEMENTS:

* secrets/pki (enterprise): Disable warnings about unknown parameters to the various CIEPS endpoints
* ui: Update PGP display and show error for Generate Operation Token flow with PGP [[GH-26993](https://github.com/hashicorp/vault/pull/26993)]

BUG FIXES:

* activity (enterprise): fix read-only storage error on upgrades
* auto-auth: Addressed issue where having no permissions to renew a renewable token caused auto-auth to attempt to renew constantly with no backoff [[GH-26844](https://github.com/hashicorp/vault/pull/26844)]
* core (enterprise): Fix an issue that prevented the seal re-wrap status from reporting that a re-wrap is in progress for up to a second.
* core/audit: Audit logging a Vault request/response will now use a minimum 5 second context timeout. 
If the existing context deadline occurs later than 5s in the future, it will be used, otherwise a new context, separate from the original will be used. [[GH-26616](https://github.com/hashicorp/vault/pull/26616)]
* core: Add missing field delegated_auth_accessors to GET /sys/mounts/:path API response [[GH-26876](https://github.com/hashicorp/vault/pull/26876)]
* core: Address a data race updating a seal's last seen healthy time attribute [[GH-27014](https://github.com/hashicorp/vault/pull/27014)]
* core: Fix `redact_version` listener parameter being ignored for some OpenAPI related endpoints. [[GH-26607](https://github.com/hashicorp/vault/pull/26607)]
* events (enterprise): Fix bug preventing subscribing and receiving events within a namepace.
* pki: Fix error in cross-signing using ed25519 keys [[GH-27093](https://github.com/hashicorp/vault/pull/27093)]
* replication (enterprise): fix "given mount path is not in the same namespace as the request" error that can occur when enabling replication for the first time on a secondary cluster
* secrets-sync (enterprise): Secondary nodes in a cluster now properly check activation-flags values.
* secrets/azure: Update vault-plugin-secrets-azure to 0.17.2 to include a bug fix for azure role creation [[GH-26896](https://github.com/hashicorp/vault/pull/26896)]
* secrets/pki (enterprise): cert_role parameter within authenticators.cert EST configuration handler could not be set
* secrets/transit: Use 'hash_algorithm' parameter if present in HMAC verify requests. Otherwise fall back to deprecated 'algorithm' parameter. [[GH-27211](https://github.com/hashicorp/vault/pull/27211)]
* ui: Fix KVv2 cursor jumping inside json editor after initial input. [[GH-27120](https://github.com/hashicorp/vault/pull/27120)]
* ui: Fix KVv2 json editor to allow null values. [[GH-27094](https://github.com/hashicorp/vault/pull/27094)]
* ui: Fix broken help link in console for the web command. [[GH-26858](https://github.com/hashicorp/vault/pull/26858)]
* ui: Fix link to v2 generic secrets engine from secrets list page. [[GH-27019](https://github.com/hashicorp/vault/pull/27019)]
* ui: Prevent perpetual loading screen when Vault needs initialization [[GH-26985](https://github.com/hashicorp/vault/pull/26985)]
* ui: Refresh model within a namespace on the Secrets Sync overview page. [[GH-26790](https://github.com/hashicorp/vault/pull/26790)]

## 1.16.2
### April 24, 2024

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* auth/jwt: Update plugin to v0.20.2 [[GH-26291](https://github.com/hashicorp/vault/pull/26291)]
* core: Bump Go version to 1.21.9.
* secrets/azure: Update plugin to v0.17.1 [[GH-26528](https://github.com/hashicorp/vault/pull/26528)]
* ui: Update dependencies including D3 libraries [[GH-26346](https://github.com/hashicorp/vault/pull/26346)]

IMPROVEMENTS:

* activity (enterprise): Change minimum retention window in activity log to 48 months
* audit: timestamps across multiple audit devices for an audit entry will now match. [[GH-26088](https://github.com/hashicorp/vault/pull/26088)]
* license utilization reporting (enterprise): Add retention months to license utilization reports.
* sdk/decompression: DecompressWithCanary will now chunk the decompression in memory to prevent loading it all at once. [[GH-26464](https://github.com/hashicorp/vault/pull/26464)]
* ui: fixes cases where inputs did not have associated labels [[GH-26263](https://github.com/hashicorp/vault/pull/26263)]
* ui: show banner instead of permission denied error when batch token is expired [[GH-26396](https://github.com/hashicorp/vault/pull/26396)]

BUG FIXES:

* agent: `vault.namespace` no longer gets incorrectly overridden by `auto_auth.namespace`, if set [[GH-26427](https://github.com/hashicorp/vault/pull/26427)]
* api: fixed a bug where LifetimeWatcher routines weren't respecting exponential backoff in the presence of unexpected errors [[GH-26383](https://github.com/hashicorp/vault/pull/26383)]
* core (enterprise): fix bug where raft followers disagree with the seal type after returning to one seal from two. [[GH-26523](https://github.com/hashicorp/vault/pull/26523)]
* core/seal: During a seal reload through SIGHUP, only write updated seal barrier on an active node [[GH-26381](https://github.com/hashicorp/vault/pull/26381)]
* core/seal: allow overriding of VAULT_GCPCKMS_SEAL_KEY_RING and VAULT_GCPCKMS_SEAL_CRYPTO_KEY environment keys in seal-ha
* events (enterprise): Terminate WebSocket connection when token is revoked.
* secrets/pki: fixed validation bug which rejected ldap schemed URLs in crl_distribution_points. [[GH-26477](https://github.com/hashicorp/vault/pull/26477)]
* storage/raft (enterprise): Fix a bug where autopilot automated upgrades could fail due to using the wrong upgrade version
* ui: fixed a bug where the replication pages did not update display when navigating between DR and performance [[GH-26325](https://github.com/hashicorp/vault/pull/26325)]
* ui: fixes undefined start time in filename for downloaded client count attribution csv [[GH-26485](https://github.com/hashicorp/vault/pull/26485)]

## 1.16.1 
### April 04, 2024

**Please note that Vault 1.16.1 is the first Enterprise release of the Vault Enterprise 1.16 series.**

BUG FIXES:

* auth/ldap: Fix login error for group search anonymous bind. [[GH-26200](https://github.com/hashicorp/vault/pull/26200)]
* auth/ldap: Fix login error missing entity alias attribute value. [[GH-26200](https://github.com/hashicorp/vault/pull/26200)]
* cli: fixed a bug where the Vault CLI would error out if HOME was not set. [[GH-26243](https://github.com/hashicorp/vault/pull/26243)]
* core: Only reload seal configuration when enable_multiseal is set to true. [[GH-26166](https://github.com/hashicorp/vault/pull/26166)]
* secret/database: Fixed race condition where database mounts may leak connections [[GH-26147](https://github.com/hashicorp/vault/pull/26147)]

## 1.16.0
### March 26, 2024

SECURITY:

* auth/cert: compare public keys of trusted non-CA certificates with incoming
client certificates to prevent trusting certs with the same serial number
but not the same public/private key (CVE-2024-2048). [[GH-25649](https://github.com/hashicorp/vault/pull/25649), [HSEC-2024-05](https://discuss.hashicorp.com/t/hcsec-2024-05-vault-cert-auth-method-did-not-correctly-validate-non-ca-certificates/63382)]
* auth/cert: validate OCSP response was signed by the expected issuer and serial number matched request (CVE-2024-2660) [[GH-26091](https://github.com/hashicorp/vault/pull/26091), [HSEC-2024-07](https://discuss.hashicorp.com/t/hcsec-2024-07-vault-tls-cert-auth-method-did-not-correctly-validate-ocsp-responses/64573)]
* secrets/transit: fix a regression that was honoring nonces provided in non-convergent modes during encryption (CVE-2023-4680) [[GH-22852](https://github.com/hashicorp/vault/pull/22852), [HSEC-2023-28](https://discuss.hashicorp.com/t/hcsec-2023-28-vault-s-transit-secrets-engine-allowed-nonce-specified-without-convergent-encryption/58249)]

CHANGES:

* Upgrade grpc to v1.58.3 [[GH-23703](https://github.com/hashicorp/vault/pull/23703)]
* Upgrade x/net to v0.17.0 [[GH-23703](https://github.com/hashicorp/vault/pull/23703)]
* api: add the `enterprise` parameter to the `/sys/health` endpoint [[GH-24270](https://github.com/hashicorp/vault/pull/24270)]
* auth/alicloud: Update plugin to v0.16.1 [[GH-25014](https://github.com/hashicorp/vault/pull/25014)]
* auth/alicloud: Update plugin to v0.17.0 [[GH-25217](https://github.com/hashicorp/vault/pull/25217)]
* auth/approle: Normalized error response messages when invalid credentials are provided [[GH-23786](https://github.com/hashicorp/vault/pull/23786)]
* auth/azure: Update plugin to v0.16.1 [[GH-22795](https://github.com/hashicorp/vault/pull/22795)]
* auth/azure: Update plugin to v0.17.0 [[GH-25258](https://github.com/hashicorp/vault/pull/25258)]
* auth/cf: Update plugin to v0.16.0 [[GH-25196](https://github.com/hashicorp/vault/pull/25196)]
* auth/gcp: Update plugin to v0.16.2 [[GH-25233](https://github.com/hashicorp/vault/pull/25233)]
* auth/jwt: Update plugin to v0.19.0 [[GH-24972](https://github.com/hashicorp/vault/pull/24972)]
* auth/jwt: Update plugin to v0.20.0 [[GH-25326](https://github.com/hashicorp/vault/pull/25326)]
* auth/jwt: Update plugin to v0.20.1 [[GH-25937](https://github.com/hashicorp/vault/pull/25937)]
* auth/kerberos: Update plugin to v0.10.1 [[GH-22797](https://github.com/hashicorp/vault/pull/22797)]
* auth/kerberos: Update plugin to v0.11.0 [[GH-25232](https://github.com/hashicorp/vault/pull/25232)]
* auth/kubernetes: Update plugin to v0.18.0 [[GH-25207](https://github.com/hashicorp/vault/pull/25207)]
* auth/oci: Update plugin to v0.14.1 [[GH-22774](https://github.com/hashicorp/vault/pull/22774)]
* auth/oci: Update plugin to v0.15.1 [[GH-25245](https://github.com/hashicorp/vault/pull/25245)]
* cli: Using `vault plugin reload` with `-plugin` in the root namespace will now reload the plugin across all namespaces instead of just the root namespace. [[GH-24878](https://github.com/hashicorp/vault/pull/24878)]
* cli: `vault plugin info` and `vault plugin deregister` now require 2 positional arguments instead of accepting either 1 or 2. [[GH-24250](https://github.com/hashicorp/vault/pull/24250)]
* core (enterprise): Seal High Availability (HA) must be enabled by `enable_multiseal` in configuration.
* core: Bump Go version to 1.21.8.
* database/couchbase: Update plugin to v0.10.1 [[GH-25275](https://github.com/hashicorp/vault/pull/25275)]
* database/elasticsearch: Update plugin to v0.14.0 [[GH-25263](https://github.com/hashicorp/vault/pull/25263)]
* database/mongodbatlas: Update plugin to v0.11.0 [[GH-25264](https://github.com/hashicorp/vault/pull/25264)]
* database/redis-elasticache: Update plugin to v0.3.0 [[GH-25296](https://github.com/hashicorp/vault/pull/25296)]
* database/redis: Update plugin to v0.2.3 [[GH-25289](https://github.com/hashicorp/vault/pull/25289)]
* database/snowflake: Update plugin to v0.10.0 [[GH-25143](https://github.com/hashicorp/vault/pull/25143)]
* database/snowflake: Update plugin to v0.9.1 [[GH-25020](https://github.com/hashicorp/vault/pull/25020)]
* events: Remove event noficiations websocket endpoint in non-Enterprise [[GH-25640](https://github.com/hashicorp/vault/pull/25640)]
* events: Source URL is now `vault://{vault node}` [[GH-24201](https://github.com/hashicorp/vault/pull/24201)]
* identity (enterprise): POST requests to the `/identity/entity/merge` endpoint
are now always forwarded from standbys to the active node. [[GH-24325](https://github.com/hashicorp/vault/pull/24325)]
* plugins/database: Reading connection config at `database/config/:name` will now return a computed `running_plugin_version` field if a non-builtin version is running. [[GH-25105](https://github.com/hashicorp/vault/pull/25105)]
* plugins: Add a warning to the response from sys/plugins/reload/backend if no plugins were reloaded. [[GH-24512](https://github.com/hashicorp/vault/pull/24512)]
* plugins: By default, environment variables provided during plugin registration will now take precedence over system environment variables.
Use the environment variable `VAULT_PLUGIN_USE_LEGACY_ENV_LAYERING=true` to opt out and keep higher preference for system environment
variables. When this flag is set, Vault will check during unseal for conflicts and print warnings for any plugins with environment
variables that conflict with system environment variables. [[GH-25128](https://github.com/hashicorp/vault/pull/25128)]
* plugins: `/sys/plugins/runtimes/catalog` response will always include a list of "runtimes" in the response, even if empty. [[GH-24864](https://github.com/hashicorp/vault/pull/24864)]
* sdk: Upgrade dependent packages by sdk.
This includes github.com/docker/docker to v24.0.7+incompatible,
google.golang.org/grpc to  v1.57.2 and golang.org/x/net to v0.17.0. [[GH-23913](https://github.com/hashicorp/vault/pull/23913)]
* secrets/ad: Update plugin to v0.16.2 [[GH-25058](https://github.com/hashicorp/vault/pull/25058)]
* secrets/ad: Update plugin to v0.17.0 [[GH-25187](https://github.com/hashicorp/vault/pull/25187)]
* secrets/alicloud: Update plugin to v0.16.0 [[GH-25257](https://github.com/hashicorp/vault/pull/25257)]
* secrets/azure: Update plugin to v0.17.0 [[GH-25189](https://github.com/hashicorp/vault/pull/25189)]
* secrets/gcp: Update plugin to v0.18.0 [[GH-25173](https://github.com/hashicorp/vault/pull/25173)]
* secrets/gcpkms: Update plugin to v0.16.0 [[GH-25231](https://github.com/hashicorp/vault/pull/25231)]
* secrets/keymgmt: Update plugin to v0.10.0
* secrets/kubernetes: Update plugin to v0.7.0 [[GH-25204](https://github.com/hashicorp/vault/pull/25204)]
* secrets/kv: Update plugin to v0.16.2 [[GH-22790](https://github.com/hashicorp/vault/pull/22790)]
* secrets/kv: Update plugin to v0.17.0 [[GH-25277](https://github.com/hashicorp/vault/pull/25277)]
* secrets/mongodbatlas: Update plugin to v0.10.2 [[GH-23849](https://github.com/hashicorp/vault/pull/23849)]
* secrets/mongodbatlas: Update plugin to v0.11.0 [[GH-25253](https://github.com/hashicorp/vault/pull/25253)]
* secrets/openldap: Update plugin to v0.11.3 [[GH-25040](https://github.com/hashicorp/vault/pull/25040)]
* secrets/openldap: Update plugin to v0.12.0 [[GH-25251](https://github.com/hashicorp/vault/pull/25251)]
* secrets/openldap: Update plugin to v0.12.1 [[GH-25524](https://github.com/hashicorp/vault/pull/25524)]
* secrets/terraform: Update plugin to v0.7.5 [[GH-25288](https://github.com/hashicorp/vault/pull/25288)]
* telemetry: Seal wrap encrypt/decrypt metrics now differentiate between seals using a metrics label of seal name rather than separate metric names. [[GH-23837](https://github.com/hashicorp/vault/pull/23837)]
* ui: Update icons to use Flight icons where available. [[GH-24823](https://github.com/hashicorp/vault/pull/24823)]
* ui: add subnav for replication items [[GH-24283](https://github.com/hashicorp/vault/pull/24283)]

FEATURES:

* **Add Snapshot Inspector Tool**: Add CLI tool to inspect Vault snapshots [[GH-23457](https://github.com/hashicorp/vault/pull/23457)]
* **Audit Filtering**: Audit devices support expression-based filter rules (powered by go-bexpr) to determine which entries are written to the audit log. [[GH-24558](https://github.com/hashicorp/vault/pull/24558)]
* **Controlled Access to Unauthenticated Endpoints (enterprise)**: Gives admins more control over how unauthenticated endpoints in Vault can be accessed and in some cases what information they return. [[GH-23547](https://github.com/hashicorp/vault/pull/23547)] [[GH-23534](https://github.com/hashicorp/vault/pull/23534)] [[GH-23740](https://github.com/hashicorp/vault/pull/23740)]
* **Custom messages (enterprise)**: Introduces custom messages settings, allowing users to view, and operators to configure system-wide messages.
* **Database Event Notifications**: The database plugin now emits event notifications. [[GH-24718](https://github.com/hashicorp/vault/pull/24718)]
* **Default Lease Count Quota (enterprise)**: Apply a new global default lease count quota of 300k leases for all
new installs of Vault. [[GH-24382](https://github.com/hashicorp/vault/pull/24382)]
* **Experimental Raft-WAL Option**: Reduces risk of infinite snapshot loops for follower nodes in large-scale Integrated Storage deployments. [[GH-21460](https://github.com/hashicorp/vault/pull/21460)]
* **Manual License Utilization Reporting**: Added manual license
utilization reporting, which allows users to create manual exports of product-license [metering
data] to report to Hashicorp.
* **Plugin Identity Tokens**: Adds secret-less configuration of AWS secret engine using web identity federation. [[GH-24987](https://github.com/hashicorp/vault/pull/24987)]
* **Plugin Workload Identity** (enterprise): Vault can generate identity tokens for plugins to use in workload identity federation auth flows.
* **Quotas in Privileged Namespaces**: Enable creation/update/deletion of quotas from the privileged namespace
* **Reload seal configuration on SIGHUP**: Seal configuration is reloaded on SIGHUP so that seal configuration can
be changed without shutting down vault [[GH-23571](https://github.com/hashicorp/vault/pull/23571)]
* **Request Limiter (enterprise)**: Add adaptive concurrency limits to
write-based HTTP methods and special-case `pki/issue` requests to prevent
overloading the Vault server. [[GH-25093](https://github.com/hashicorp/vault/pull/25093)]
* **Rotate Root for LDAP auth**: Rotate root operations are now supported for the LDAP auth engine. [[GH-24099](https://github.com/hashicorp/vault/pull/24099)]
* **Seal High Availability (enterprise)**: Operators can configure more than one automatic seal for resilience against seal provider outages.
* **Secrets Sync UI (enterprise)**: Adds secret syncing for KV v2 secrets to external destinations using the UI. [[GH-23667](https://github.com/hashicorp/vault/pull/23667)]
* **Vault PKI EST Server (Enterprise/Beta)**: Beta support for the PKI Enrollment over Secure Transport (EST) certificate management protocol has been added to the Vault PKI Plugin. This allows standard EST clients to request certificates from a Vault server with no knowledge of Vault APIs.
* **Vault Proxy Static Secret Caching (enterprise)**: Adds support for static secret (KVv1 and KVv2) caching to Vault Proxy. [[GH-23621](https://github.com/hashicorp/vault/pull/23621)]
* **secrets-import (enterprise)**: Support importing secrets from external sources into KVv2
* **secrets/aws**: Support issuing an STS Session Token directly from the root credential. [[GH-23690](https://github.com/hashicorp/vault/pull/23690)]

IMPROVEMENTS:

* .release/linux: add LimitCORE=0 to vault.service [[GH-23272](https://github.com/hashicorp/vault/pull/23272)]
* agent/template: Added max_connections_per_host to limit total number of connections per Vault host. [[GH-24548](https://github.com/hashicorp/vault/pull/24548)]
* agent: Added new namespace top level configuration parameter, which can be used to make requests made by Agent to go to that namespace. [[GH-24667](https://github.com/hashicorp/vault/pull/24667)]
* agent: allow users to specify files for child process stdout/stderr [[GH-22812](https://github.com/hashicorp/vault/pull/22812)]
* api (enterprise): Enable the sys/license/features from any namespace
* api/plugins: add `tls-server-name` arg for plugin registration [[GH-23549](https://github.com/hashicorp/vault/pull/23549)]
* api: Add wrapper functions for GET /sys/mounts/:path and GET /sys/auth/:path [[GH-25499](https://github.com/hashicorp/vault/pull/25499)]
* api: Do not require sudo for API wrapper functions GetAuth and GetAuthWithContext [[GH-25968](https://github.com/hashicorp/vault/pull/25968)]
* api: added new API field to Vault responses, `mount_type`, returning mount information (e.g. `kv` for KVV1/KVV2) for mount when appropriate. [[GH-23047](https://github.com/hashicorp/vault/pull/23047)]
* api: sys/health and sys/ha-status now expose information about how long
the last heartbeat took, and the estimated clock skew between standby and
active node based on that heartbeat duration. [[GH-24343](https://github.com/hashicorp/vault/pull/24343)]
* auth/cert: Allow validation with OCSP responses with no NextUpdate time [[GH-25912](https://github.com/hashicorp/vault/pull/25912)]
* auth/cert: Cache trusted certs to reduce memory usage and improve performance of logins. [[GH-25421](https://github.com/hashicorp/vault/pull/25421)]
* auth/ldap: introduce cap/ldap.Client for LDAP authentication
auth/ldap: deprecates `connection_timeout` in favor of `request_timeout` for timeouts
sdk/ldaputil: deprecates Client in favor of cap/ldap.Client [[GH-22185](https://github.com/hashicorp/vault/pull/22185)]
* auth/saml: Update plugin to v0.2.0
* auto-auth/azure: Support setting the `authenticate_from_environment` variable to "true" and "false" string literals, too. [[GH-22996](https://github.com/hashicorp/vault/pull/22996)]
* cli: introduce new command group hcp which groups subcommands for authentication of users or machines to HCP using
either provided arguments or retrieved HCP token through browser login. [[GH-23897](https://github.com/hashicorp/vault/pull/23897)]
* cli: Improved error messages for `vault plugin` sub-commands. [[GH-24250](https://github.com/hashicorp/vault/pull/24250)]
* cli: adds plugin identity token to enable and tune commands for secret engines and auth methods [[GH-24980](https://github.com/hashicorp/vault/pull/24980)]
* cli: include secret syncs counts in the `vault operator usage` command output [[GH-25751](https://github.com/hashicorp/vault/pull/25751)]
* command/server: display logs on startup immediately if disable-gated-logs flag is set [[GH-24280](https://github.com/hashicorp/vault/pull/24280)]
* command/token-capabilities: allow using accessor when listing token capabilities on a path [[GH-24479](https://github.com/hashicorp/vault/pull/24479)]
* core (enterprise): Avoid seal rewrapping in some specific unnecessary cases.
* core (enterprise): Improve seal unwrap performance when in degraded mode with one or more unhealthy seals. [[GH-25171](https://github.com/hashicorp/vault/pull/25171)]
* core (enterprise): Speed up unseal when using namespaces
* core (enterprise): persist seal rewrap status, so rewrap status API is consistent on secondary nodes.
* core/activity: Include secret_syncs in activity log responses [[GH-24710](https://github.com/hashicorp/vault/pull/24710)]
* core/cli: Warning related to VAULT_ADDR & -address not set with CLI requests. [[GH-17076](https://github.com/hashicorp/vault/pull/17076)]
* core/metrics: add metrics for secret sync client count [[GH-25713](https://github.com/hashicorp/vault/pull/25713)]
* core: Added new `plugin_tmpdir` config option for containerized plugins, in addition to the existing `VAULT_PLUGIN_TMPDIR` environment variable. [[GH-24978](https://github.com/hashicorp/vault/pull/24978)]
* core: make the best effort timeout for encryption count tracking persistence configurable via an environment variable. [[GH-25636](https://github.com/hashicorp/vault/pull/25636)]
* core: update sys/seal-status (and CLI vault status) to report the type of
the seal when unsealed, as well as the type of the recovery seal if an
auto-seal. [[GH-23022](https://github.com/hashicorp/vault/pull/23022)]
* events: Add support for event subscription plugins, including SQS [[GH-24352](https://github.com/hashicorp/vault/pull/24352)]
* identity/tokens: adds plugin issuer with openid-configuration and keys APIs [[GH-24898](https://github.com/hashicorp/vault/pull/24898)]
* limits: Add a listener configuration option `disable_request_limiter` to allow
disabling the request limiter per-listener. [[GH-25098](https://github.com/hashicorp/vault/pull/25098)]
* limits: Introduce a reloadable opt-in configuration for the Request Limiter. [[GH-25095](https://github.com/hashicorp/vault/pull/25095)]
* oidc/provider: Adds `code_challenge_methods_supported` to OpenID Connect Metadata [[GH-24979](https://github.com/hashicorp/vault/pull/24979)]
* plugins: Add new pin version APIs to enforce all plugins of a specific type and name to run the same version. [[GH-25105](https://github.com/hashicorp/vault/pull/25105)]
* plugins: Containerized plugins can be run fully rootless with the runsc runtime. [[GH-24236](https://github.com/hashicorp/vault/pull/24236)]
* plugins: New API `sys/plugins/reload/:type/:name` available in the root namespace for reloading a specific plugin across all namespaces. [[GH-24878](https://github.com/hashicorp/vault/pull/24878)]
* proxy: Added new namespace top level configuration parameter, and prepend_configured_namespace API Proxy configuration parameter, which can be used to make requests made to Proxy get proxied to that namespace. [[GH-24667](https://github.com/hashicorp/vault/pull/24667)]
* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20221104090112-13395acd02c5
* replication (enterprise): Add last_upstream_remote_wal metric to telemetry and stop emitting last_remote_wal on standby nodes
* replication (enterprise): Add re-index status metric to telemetry
* replication: Add re-index status metric to telemetry [[GH-23160](https://github.com/hashicorp/vault/pull/23160)]
* sdk/plugin: Fix an issue where external plugins were not reporting logs below INFO level [[GH-23771](https://github.com/hashicorp/vault/pull/23771)]
* sdk: Add identity token helpers to consistently apply new plugin WIF fields across integrations. [[GH-24925](https://github.com/hashicorp/vault/pull/24925)]
* sdk: adds new method to system view to allow plugins to request identity tokens [[GH-24929](https://github.com/hashicorp/vault/pull/24929)]
* secrets-sync (enterprise): Add ability to turn the sync system on and off
* secrets-sync (enterprise): Add reconciliation loop
* secrets-sync (enterprise): Added PATCH request on the sync destinations API
* secrets-sync (enterprise): Added delete request to reset global config to factory defaults
* secrets-sync (enterprise): Added field to define custom tags to add on synced secrets
* secrets-sync (enterprise): Added global config path to the administrative namespace.
* secrets-sync (enterprise): Added telemetry on number of destinations and associations per type.
* secrets-sync (enterprise): Adds ability to set target GCP project ID to sync secrets with
* secrets-sync (enterprise): Adjusted associations list responses to be more CLI-friendly
* secrets-sync (enterprise): Adjusted destination list responses to be more CLI-friendly & added endpoint to list destinations by type.
* secrets-sync (enterprise): Clean up membdb tests
* secrets-sync (enterprise): Support AWS IAM assume role and external ID
* secrets-sync (enterprise): Support custom GitHub apps
* secrets-sync (enterprise): Support custom templating of external secret names
* secrets-sync (enterprise): Support granular secrets syncing
* secrets-sync (enterprise): add purge field to the destination delete endpoint
* secrets/database: Add new reload/:plugin_name API to reload database plugins by name for a specific mount. [[GH-24472](https://github.com/hashicorp/vault/pull/24472)]
* secrets/database: Support reloading named database plugins using the sys/plugins/reload/backend API endpoint. [[GH-24512](https://github.com/hashicorp/vault/pull/24512)]
* secrets/pki: do not check TLS validity on ACME requests redirected to https [[GH-22521](https://github.com/hashicorp/vault/pull/22521)]
* storage/etcd: etcd should only return keys when calling List() [[GH-23872](https://github.com/hashicorp/vault/pull/23872)]
* storage/raft (enterprise): Replication WAL batches may now contain up to 4096
entries rather than being limited to 62 like Consul is. Performance testing
shows improvements in throughput and latency under some write-heavy workloads.
* storage/raft: Add support for larger transactions when using raft storage. [[GH-24991](https://github.com/hashicorp/vault/pull/24991)]
* storage/raft: Upgrade to bbolt 1.3.8, along with an extra patch to reduce time scanning large freelist maps. [[GH-24010](https://github.com/hashicorp/vault/pull/24010)]
* sys (enterprise): Enable sys/config/group-application-policy in privileged namespace
* sys (enterprise): Adds the chroot_namespace field to this sys/internal/ui/resultant-acl endpoint, which exposes the value of the chroot namespace from the
listener config.
* sys: adds configuration of the key used to sign plugin identity tokens during mount enable and tune [[GH-24962](https://github.com/hashicorp/vault/pull/24962)]
* ui: Add `deletion_allowed` param to transformations and include `tokenization` as a type option [[GH-25436](https://github.com/hashicorp/vault/pull/25436)]
* ui: Add warning message to the namespace picker warning users about the behavior when logging in with a root token. [[GH-23277](https://github.com/hashicorp/vault/pull/23277)]
* ui: Adds a warning when whitespace is detected in a key of a KV secret [[GH-23702](https://github.com/hashicorp/vault/pull/23702)]
* ui: Adds allowed_response_headers, plugin_version and user_lockout_config params to auth method configuration [[GH-25646](https://github.com/hashicorp/vault/pull/25646)]
* ui: Adds toggle to KV secrets engine value download modal to optionally stringify value in downloaded file [[GH-23747](https://github.com/hashicorp/vault/pull/23747)]
* ui: Allow users in userpass auth mount to update their own password [[GH-23797](https://github.com/hashicorp/vault/pull/23797)]
* ui: Implement Helios Design System Breadcrumbs [[GH-24387](https://github.com/hashicorp/vault/pull/24387)]
* ui: Implement Helios Design System copy button component making copy buttons accessible [[GH-22333](https://github.com/hashicorp/vault/pull/22333)]
* ui: Implement Helios Design System footer component [[GH-24191](https://github.com/hashicorp/vault/pull/24191)]
* ui: Implement Helios Design System pagination component [[GH-23169](https://github.com/hashicorp/vault/pull/23169)]
* ui: Increase base font-size from 14px to 16px and update use of rem vs pixels for size variables [[GH-23994](https://github.com/hashicorp/vault/pull/23994)]
* ui: Makes modals accessible by implementing Helios Design System modal component [[GH-23382](https://github.com/hashicorp/vault/pull/23382)]
* ui: Replace inline confirm alert inside a popup-menu dropdown with confirm alert modal [[GH-21520](https://github.com/hashicorp/vault/pull/21520)]
* ui: Separates out client counts dashboard to overview and entity/non-entity tabs [[GH-24752](https://github.com/hashicorp/vault/pull/24752)]
* ui: Sort list view of entities and aliases alphabetically using the item name [[GH-24103](https://github.com/hashicorp/vault/pull/24103)]
* ui: Surface warning banner if UI has stopped auto-refreshing token [[GH-23143](https://github.com/hashicorp/vault/pull/23143)]
* ui: Update AlertInline component to use Helios Design System Alert component [[GH-24299](https://github.com/hashicorp/vault/pull/24299)]
* ui: Update flat, shell-quote and swagger-ui-dist packages. Remove swagger-ui styling overrides. [[GH-23700](https://github.com/hashicorp/vault/pull/23700)]
* ui: Update mount backend form to use selectable cards [[GH-14998](https://github.com/hashicorp/vault/pull/14998)]
* ui: Update sidebar Secrets engine to title case. [[GH-23964](https://github.com/hashicorp/vault/pull/23964)]
* ui: Use Hds::Dropdown component to replace list view popup menus [[GH-25321](https://github.com/hashicorp/vault/pull/25321)]
* ui: add error message when copy action fails [[GH-25479](https://github.com/hashicorp/vault/pull/25479)]
* ui: add granularity param to sync destinations [[GH-25500](https://github.com/hashicorp/vault/pull/25500)]
* ui: capabilities-self is always called in the user's root namespace [[GH-24168](https://github.com/hashicorp/vault/pull/24168)]
* ui: improve accessibility - color contrast, labels, and automatic testing [[GH-24476](https://github.com/hashicorp/vault/pull/24476)]
* ui: latest version of chrome does not automatically redirect back to the app after authentication unless triggered by the user, hence added a link to redirect back to the app. [[GH-18513](https://github.com/hashicorp/vault/pull/18513)]
* ui: obscure JSON values when KV v2 secret has nested objects [[GH-24530](https://github.com/hashicorp/vault/pull/24530)]
* ui: redirect back to current route after reauthentication when token expires [[GH-25335](https://github.com/hashicorp/vault/pull/25335)]
* ui: remove leading slash from KV version 2 secret paths [[GH-25874](https://github.com/hashicorp/vault/pull/25874)]
* ui: remove unnecessary OpenAPI calls for unmanaged auth methods [[GH-25364](https://github.com/hashicorp/vault/pull/25364)]
* ui: replace popup menu on list items (namespaces, auth items, KMIP, K8S, LDAP) [[GH-25588](https://github.com/hashicorp/vault/pull/25588)]
* ui: show banner when resultant-acl check fails due to permissions or wrong namespace. [[GH-23503](https://github.com/hashicorp/vault/pull/23503)]
* website/docs: Update references to Key Value secrets engine from 'K/V' to 'KV' [[GH-24529](https://github.com/hashicorp/vault/pull/24529)]
* website/docs: fix inaccuracies with unauthenticated_in_flight_requests_access parameter [[GH-23287](https://github.com/hashicorp/vault/pull/23287)]

BUG FIXES:

* Seal HA (enterprise/beta): Fix rejection of a seal configuration change
from two to one auto seal due to persistence of the previous seal type being
"multiseal". [[GH-23573](https://github.com/hashicorp/vault/pull/23573)]
* activity log (enterprise): De-duplicate client count estimates for license utilization reporting.
* agent/logging: Agent should now honor correct -log-format and -log-file settings in logs generated by the consul-template library. [[GH-24252](https://github.com/hashicorp/vault/pull/24252)]
* agent: Fix issue where Vault Agent was unable to render KVv2 secrets with delete_version_after set. [[GH-25387](https://github.com/hashicorp/vault/pull/25387)]
* agent: Fixed incorrect parsing of boolean environment variables for configuration. [[GH-24790](https://github.com/hashicorp/vault/pull/24790)]
* api/seal-status: Fix deadlock on calls to sys/seal-status with a namespace configured
on the request. [[GH-23861](https://github.com/hashicorp/vault/pull/23861)]
* api: Fix deadlock on calls to sys/leader with a namespace configured
on the request. [[GH-24256](https://github.com/hashicorp/vault/pull/24256)]
* api: sys/leader ActiveTime field no longer gets reset when we do an internal state change that doesn't change our active status. [[GH-24549](https://github.com/hashicorp/vault/pull/24549)]
* audit/socket: Provide socket based audit backends with 'prefix' configuration option when supplied. [[GH-25004](https://github.com/hashicorp/vault/pull/25004)]
* audit: Fix bug reopening 'file' audit devices on SIGHUP. [[GH-23598](https://github.com/hashicorp/vault/pull/23598)]
* audit: Fix bug where use of 'log_raw' option could result in other devices logging raw audit data [[GH-24968](https://github.com/hashicorp/vault/pull/24968)]
* audit: Handle a potential panic while formatting audit entries for an audit log [[GH-25605](https://github.com/hashicorp/vault/pull/25605)]
* audit: Operator changes to configured audit headers (via `/sys/config/auditing`)
will now force invalidation and be reloaded from storage when data is replicated
to other nodes.
* audit: Resolve potential race condition when auditing entries which use SSCT. [[GH-25443](https://github.com/hashicorp/vault/pull/25443)]
* auth/aws: Fixes a panic that can occur in IAM-based login when a [client config](https://developer.hashicorp.com/vault/api-docs/auth/aws#configure-client) does not exist. [[GH-23555](https://github.com/hashicorp/vault/pull/23555)]
* auth/cert: Address an issue in which OCSP query responses were not cached [[GH-25986](https://github.com/hashicorp/vault/pull/25986)]
* auth/cert: Allow cert auth login attempts if ocsp_fail_open is true and OCSP servers are unreachable [[GH-25982](https://github.com/hashicorp/vault/pull/25982)]
* auth/cert: Handle errors related to expired OCSP server responses [[GH-24193](https://github.com/hashicorp/vault/pull/24193)]
* auth/saml (enterprise): Fixes support for Microsoft Entra ID enterprise applications
* cap/ldap: Downgrade go-ldap client from v3.4.5 to v3.4.4 due to race condition found [[GH-23103](https://github.com/hashicorp/vault/pull/23103)]
* cassandra: Update Cassandra to set consistency prior to calling CreateSession, ensuring consistency setting is correct when opening connection. [[GH-24649](https://github.com/hashicorp/vault/pull/24649)]
* cli/kv: Undelete now properly handles KV-V2 mount paths that are more than one layer deep. [[GH-19811](https://github.com/hashicorp/vault/pull/19811)]
* cli: fixes plugin register CLI failure to error when plugin image doesn't exist [[GH-24990](https://github.com/hashicorp/vault/pull/24990)]
* command/server: Fix bug with sigusr2 where pprof files were not closed correctly [[GH-23636](https://github.com/hashicorp/vault/pull/23636)]
* core (Enterprise): Treat multiple disabled HA seals as a migration to Shamir.
* core (enterprise): Do not return an internal error when token policy type lookup fails, log it instead and continue.
* core (enterprise): Fix a deadlock that can occur on performance secondary clusters when there are many mounts and a mount is deleted or filtered [[GH-25448](https://github.com/hashicorp/vault/pull/25448)]
* core (enterprise): Fix a panic that can occur if only one seal exists but is unhealthy on the non-first restart of Vault.
* core (enterprise): fix a potential deadlock if an error is received twice from underlying storage for the same key
* core (enterprise): fix issue where the Seal HA rewrap system may remain running when an active node steps down.
* core/activity: Fixes segments fragment loss due to exceeding entry record size limit [[GH-23781](https://github.com/hashicorp/vault/pull/23781)]
* core/audit: Audit logging a Vault response will now use a 5 second context timeout, separate from the original request. [[GH-24238](https://github.com/hashicorp/vault/pull/24238)]
* core/config: Use correct HCL config value when configuring `log_requests_level`. [[GH-24056](https://github.com/hashicorp/vault/pull/24056)]
* core/ha: fix panic that can occur when an HA cluster contains an active node with version >=1.12.0 and another node with version <1.10 [[GH-24441](https://github.com/hashicorp/vault/pull/24441)]
* core/login: Fixed a potential deadlock when a login fails and user lockout is enabled. [[GH-25697](https://github.com/hashicorp/vault/pull/25697)]
* core/mounts: Fix reading an "auth" mount using "sys/internal/ui/mounts/" when filter paths are enforced returns 500 error code from the secondary [[GH-23802](https://github.com/hashicorp/vault/pull/23802)]
* core/quotas: Close rate-limit blocked client purge goroutines when sealing [[GH-24108](https://github.com/hashicorp/vault/pull/24108)]
* core/quotas: Deleting a namespace that contains a rate limit quota no longer breaks replication [[GH-25439](https://github.com/hashicorp/vault/pull/25439)]
* core: Fix a timeout initializing Vault by only using a short timeout persisting barrier keyring encryption counts. [[GH-24336](https://github.com/hashicorp/vault/pull/24336)]
* core: Fix an error that resulted in the wrong seal type being returned by sys/seal-status while
Vault is in seal migration mode. [[GH-24165](https://github.com/hashicorp/vault/pull/24165)]
* core: Skip unnecessary deriving of policies during Login MFA Check. [[GH-23894](https://github.com/hashicorp/vault/pull/23894)]
* core: fix bug where deadlock detection was always on for expiration and quotas.
These can now be configured individually with `detect_deadlocks`. [[GH-23902](https://github.com/hashicorp/vault/pull/23902)]
* core: fix policies with wildcards not matching list operations due to the policy path not having a trailing slash [[GH-23874](https://github.com/hashicorp/vault/pull/23874)]
* core: fix rare panic due to a race condition with metrics collection during seal [[GH-23906](https://github.com/hashicorp/vault/pull/23906)]
* core: upgrade github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault/v2 to
support azure workload identities. [[GH-24954](https://github.com/hashicorp/vault/pull/24954)]
* eventlogger: Update library to v0.2.7 to address race condition [[GH-24305](https://github.com/hashicorp/vault/pull/24305)]
* events: Ignore sending context to give more time for events to send [[GH-23500](https://github.com/hashicorp/vault/pull/23500)]
* expiration: Fix fatal error "concurrent map iteration and map write" when collecting metrics from leases. [[GH-24027](https://github.com/hashicorp/vault/pull/24027)]
* expiration: Prevent large lease loads from delaying state changes, e.g. becoming active or standby. [[GH-23282](https://github.com/hashicorp/vault/pull/23282)]
* fairshare: fix a race condition in JobManager.GetWorkerCounts [[GH-24616](https://github.com/hashicorp/vault/pull/24616)]
* helper/pkcs7: Fix slice out-of-bounds panic [[GH-24891](https://github.com/hashicorp/vault/pull/24891)]
* http: Include PATCH in the list of allowed CORS methods [[GH-24373](https://github.com/hashicorp/vault/pull/24373)]
* kmip (enterprise): Improve handling of failures due to storage replication issues.
* kmip (enterprise): Only return a Server Correlation Value to clients using KMIP version 1.4.
* kmip (enterprise): Return a structure in the response for query function Query Server Information.
* ldaputil: Disable tests for ARM64 [[GH-23118](https://github.com/hashicorp/vault/pull/23118)]
* mongo-db: allow non-admin database for root credential rotation [[GH-23240](https://github.com/hashicorp/vault/pull/23240)]
* openapi: Fixing response fields for rekey operations [[GH-25509](https://github.com/hashicorp/vault/pull/25509)]
* plugins: Fix panic when querying plugin runtimes from a performance secondary follower node.
* proxy: Fixed incorrect parsing of boolean environment variables for configuration. [[GH-24790](https://github.com/hashicorp/vault/pull/24790)]
* replication (enterprise): Fix a bug where undo logs would only get enabled on the initial node in a cluster.
* replication (enterprise): Fix a missing unlock when changing replication state
* replication (enterprise): disallow configuring paths filter for a mount path that does not exist
* replication (enterprise): fixed data integrity issue with the processing of identity aliases causing duplicates to occur in rare cases
* sdk: Return error when failure occurs setting up node 0 in NewDockerCluster, instead of ignoring it. [[GH-24136](https://github.com/hashicorp/vault/pull/24136)]
* secrets-sync (enterprise): Allow unsyncing secrets from an unmounted secrets engine
* secrets-sync (enterprise): Fix panic when setting usage_gauge_period to none
* secrets-sync (enterprise): Fixed an issue with syncing to target projects in GCP
* secrets-sync (enterprise): Fixed issue where we could sync a deleted secret
* secrets-sync (enterprise): Unsync secret when metadata is deleted
* secrets/aws: fix requeueing of rotation entry in cases where rotation fails [[GH-23673](https://github.com/hashicorp/vault/pull/23673)]
* secrets/aws: update credential rotation deadline when static role rotation period is updated [[GH-23528](https://github.com/hashicorp/vault/pull/23528)]
* secrets/consul: Fix revocations when Vault has an access token using specific namespace and admin partition policies [[GH-23010](https://github.com/hashicorp/vault/pull/23010)]
* secrets/pki: Do not set nextUpdate field in OCSP responses when ocsp_expiry is 0 [[GH-24192](https://github.com/hashicorp/vault/pull/24192)]
* secrets/pki: Stop processing in-flight ACME verifications when an active node steps down [[GH-23278](https://github.com/hashicorp/vault/pull/23278)]
* secrets/transit (enterprise): Address an issue using sign/verify operations with managed keys returning an error about it not containing a private key
* secrets/transit (enterprise): Address panic when using GCP,AWS,Azure managed keys for encryption operations. At this time all encryption operations for the cloud providers have been disabled, only signing operations are supported.
* secrets/transit (enterprise): Apply hashing arguments and defaults to managed key sign/verify operations
* secrets/transit: Do not allow auto rotation on managed_key key types [[GH-23723](https://github.com/hashicorp/vault/pull/23723)]
* secrets/transit: Fix a panic when attempting to export a public RSA key [[GH-24054](https://github.com/hashicorp/vault/pull/24054)]
* secrets/transit: When provided an invalid input with hash_algorithm=none, a lock was not released properly before reporting an error leading to deadlocks on a subsequent key configuration update. [[GH-25336](https://github.com/hashicorp/vault/pull/25336)]
* storage/consul: fix a bug where an active node in a specific sort of network
partition could continue to write data to Consul after a new leader is elected
potentially causing data loss or corruption for keys with many concurrent
writers. For Enterprise clusters this could cause corruption of the merkle trees
leading to failure to complete merkle sync without a full re-index. [[GH-23013](https://github.com/hashicorp/vault/pull/23013)]
* storage/file: Fixing spuriously deleting storage keys ending with .temp [[GH-25395](https://github.com/hashicorp/vault/pull/25395)]
* storage/raft: Fix a race whereby a new leader may present inconsistent node data to Autopilot. [[GH-24246](https://github.com/hashicorp/vault/pull/24246)]
* transform (enterprise): guard against a panic looking up a token in exportable mode with barrier storage.
* ui: Allows users to dismiss the resultant-acl banner. [[GH-25106](https://github.com/hashicorp/vault/pull/25106)]
* ui: Allows users to search within KV v2 directories from the Dashboard's quick action card. [[GH-25001](https://github.com/hashicorp/vault/pull/25001)]
* ui: Assumes version 1 for kv engines when options are null because no version is specified [[GH-23585](https://github.com/hashicorp/vault/pull/23585)]
* ui: Correctly handle directory redirects from pre 1.15.0 Kv v2 list view urls. [[GH-24281](https://github.com/hashicorp/vault/pull/24281)]
* ui: Correctly handle redirects from pre 1.15.0 Kv v2 edit, create, and show urls. [[GH-24339](https://github.com/hashicorp/vault/pull/24339)]
* ui: Decode the connection url for display on the connection details page [[GH-23695](https://github.com/hashicorp/vault/pull/23695)]
* ui: Do not disable JSON display toggle for KV version 2 secrets [[GH-25235](https://github.com/hashicorp/vault/pull/25235)]
* ui: Do not show resultant-acl banner on namespaces a user has access to [[GH-25256](https://github.com/hashicorp/vault/pull/25256)]
* ui: Fix AWS secret engine to allow empty policy_document field. [[GH-23470](https://github.com/hashicorp/vault/pull/23470)]
* ui: Fix JSON editor in KV V2 unable to handle pasted values [[GH-24224](https://github.com/hashicorp/vault/pull/24224)]
* ui: Fix PKI ca_chain display so value can be copied to clipboard [[GH-25399](https://github.com/hashicorp/vault/pull/25399)]
* ui: Fix bug where a change on OpenAPI added a double forward slash on some LIST endpoints. [[GH-23446](https://github.com/hashicorp/vault/pull/23446)]
* ui: Fix copy button not working on masked input when value is not a string [[GH-25269](https://github.com/hashicorp/vault/pull/25269)]
* ui: Fix error when tuning token auth configuration within namespace [[GH-24147](https://github.com/hashicorp/vault/pull/24147)]
* ui: Fix inconsistent empty state action link styles [[GH-25209](https://github.com/hashicorp/vault/pull/25209)]
* ui: Fix kubernetes auth method roles tab [[GH-25999](https://github.com/hashicorp/vault/pull/25999)]
* ui: Fix payload sent when disabling replication [[GH-24292](https://github.com/hashicorp/vault/pull/24292)]
* ui: Fix regression that broke the oktaNumberChallenge on the ui. [[GH-23565](https://github.com/hashicorp/vault/pull/23565)]
* ui: Fix the copy token button in the sidebar navigation window when in a collapsed state. [[GH-23331](https://github.com/hashicorp/vault/pull/23331)]
* ui: Fixed minor bugs with database secrets engine [[GH-24947](https://github.com/hashicorp/vault/pull/24947)]
* ui: Fixes input for jwks_ca_pem when configuring a JWT auth method [[GH-24697](https://github.com/hashicorp/vault/pull/24697)]
* ui: Fixes issue where you could not share the list view URL from the KV v2 secrets engine. [[GH-23620](https://github.com/hashicorp/vault/pull/23620)]
* ui: Fixes issue with no active tab when viewing transit encryption key [[GH-25614](https://github.com/hashicorp/vault/pull/25614)]
* ui: Fixes issue with sidebar navigation links disappearing when navigating to policies when a user is not authorized [[GH-23516](https://github.com/hashicorp/vault/pull/23516)]
* ui: Fixes issues displaying accurate TLS state in dashboard configuration details [[GH-23726](https://github.com/hashicorp/vault/pull/23726)]
* ui: Fixes policy input toolbar scrolling by default [[GH-23297](https://github.com/hashicorp/vault/pull/23297)]
* ui: The UI can now be used to create or update database roles by operator without permission on the database connection. [[GH-24660](https://github.com/hashicorp/vault/pull/24660)]
* ui: Update the KV secret data when you change the version you're viewing of a nested secret. [[GH-25152](https://github.com/hashicorp/vault/pull/25152)]
* ui: Updates OIDC/JWT login error handling to surface all role related errors [[GH-23908](https://github.com/hashicorp/vault/pull/23908)]
* ui: Upgrade HDS version to fix sidebar navigation issues when it collapses in smaller viewports. [[GH-23580](https://github.com/hashicorp/vault/pull/23580)]
* ui: When Kv v2 secret is an object, fix so details view defaults to readOnly JSON editor. [[GH-24290](https://github.com/hashicorp/vault/pull/24290)]
* ui: call resultant-acl without namespace header when user mounted at root namespace [[GH-25766](https://github.com/hashicorp/vault/pull/25766)]
* ui: fix KV v2 details view defaulting to JSON view when secret value includes `{` [[GH-24513](https://github.com/hashicorp/vault/pull/24513)]
* ui: fix broken GUI when accessing from listener with chroot_namespace defined [[GH-23942](https://github.com/hashicorp/vault/pull/23942)]
* ui: fix incorrectly calculated capabilities on PKI issuer endpoints [[GH-24686](https://github.com/hashicorp/vault/pull/24686)]
* ui: fix issue where kv v2 capabilities checks were not passing in the full secret path if secret was inside a directory. [[GH-24404](https://github.com/hashicorp/vault/pull/24404)]
* ui: fix navigation items shown to user when chroot_namespace configured [[GH-24492](https://github.com/hashicorp/vault/pull/24492)]
* ui: remove user_lockout_config settings for unsupported methods [[GH-25867](https://github.com/hashicorp/vault/pull/25867)]
* ui: show error from API when seal fails [[GH-23921](https://github.com/hashicorp/vault/pull/23921)]
