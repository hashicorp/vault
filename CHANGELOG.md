## Previous versions
- [v1.10.0 - v1.15.16](CHANGELOG-v1.10-v1.15.md)
- [v1.0.0 - v1.9.10](CHANGELOG-pre-v1.10.md)
- [v0.11.6 and earlier](CHANGELOG-v0.md)

## 1.21.2
### January 07, 2026

CHANGES:

* auth/oci: bump plugin to v0.20.1
* core: Bump Go version to 1.25.5
* packaging: Container images are now exported using a compressed OCI image layout.
* packaging: UBI container images are now built on the UBI 10 minimal image.
* secrets/azure: Update plugin to v0.25.1+ent. Improves retry handling during Azure application and service principal creation to reduce transient failures.
* storage: Upgrade aerospike client library to v8.

IMPROVEMENTS:

* core: check rotation manager queue every 5 seconds instead of 10 seconds to improve responsiveness
* go: update to golang/x/crypto to v0.45.0 to resolve GHSA-f6x5-jh6r-wrfv, GHSA-j5w8-q4qc-rx2x, GO-2025-4134 and GO-2025-4135.
* rotation: Ensure rotations for shared paths only execute on the Primary cluster's active node. Ensure rotations for local paths execute on the cluster-local active node.
* sdk/rotation: Prevent rotation attempts on read-only storage.
* secrets-sync (enterprise): Added support for a boolean force_delete flag (default: false). When set to true, this flag allows deletion of a destination even if its associations cannot be unsynced. This option should be used only as a last-resort deletion mechanism, as any secrets already synced to the external provider will remain orphaned and require manual cleanup.
* secrets/pki: Avoid loading issuer information multiple times per leaf certificate signing.

BUG FIXES:

* core/activitylog (enterprise): Resolve a stability issue where Vault Enterprise could encounter a panic during month-end billing activity rollover.
* http: skip JSON limit parsing on cluster listener.
* quotas: Vault now protects plugins with ResolveRole operations from panicking on quota creation.
* replication (enterprise): fix rare panic due to race when enabling a secondary with Consul storage.
* rotation: Fix a bug where a performance secondary would panic if a write was made to a local mount.
* secret-sync (enterprise): Improved unsync error handling by treating cases where the destination no longer exists as successful.
* secrets-sync (enterprise): Corrected a bug where the deletion of the latest KV-V2 secret version caused the associated external secret to be deleted entirely. The sync job now implements a version fallback mechanism to find and sync the highest available active version, ensuring continuity and preventing the unintended deletion of the external secret resource.
* secrets-sync (enterprise): Fix issue where secrets were not properly un-synced after destination config changes.
* secrets-sync (enterprise): Fix issue where sync store deletion could be attempted when sync is disabled.
* ui/pki: Fix handling of values that contain commas in list fields like `crl_distribution_points`.

## 1.21.1
### November 19, 2025

SECURITY:

* auth/aws: fix an issue where a user may be able to bypass authentication to Vault due to incorrect caching of the AWS client
* ui: disable scarf analytics for ui builds

CHANGES:

* auth/kubernetes: Update plugin to [v0.23.1](https://github.com/hashicorp/vault-plugin-auth-kubernetes/releases/tag/v0.23.1)
* auth/saml: Update plugin to [v0.7.1](https://github.com/hashicorp/vault-plugin-auth-saml/releases/tag/v0.7.1), which adds the environment variable VAULT_SAML_DENY_INTERNAL_URLS to allow prevention of idp_metadata_url, idp_sso_url, or acs_urls fields from containing URLs that resolve to internal IP addresses
* core: Bump Go version to 1.25.4
* secrets/azure (enterprise): Update plugin to v0.25.0+ent
* secrets/pki: sign-verbatim endpoints no longer ignore basic constraints extension in CSRs, using them in generated certificates if isCA=false or returning an error if isCA=true

IMPROVEMENTS:

* Update github.com/dvsekhvalnov/jose2go to fix security vulnerability CVE-2025-63811.
* api: Added sudo-permissioned `sys/reporting/scan` endpoint which will output a set of files containing information about Vault state to the location specified by the `reporting_scan_directory` config item.
* auth/ldap: Require non-empty passwords on login command to prevent unauthenticated access to Vault.
* core/metrics: Reading and listing from a snapshot are now tracked via the `vault.route.read-snapshot.{mount_point}` and `vault.route.list-snapshot.{mount_point}` metrics.
* license utilization reporting (enterprise): Add metrics for the number of issued PKI certificates.
* policies: add warning about list comparison when using allowed_parameters or denied_parameters
* secret-sync: add parallelization support to sync and unsync operations for secret-key granularity associations
* secrets/pki: Include the certificate's AuthorityKeyID in response fields for API endpoints that issue, sign, or fetch certs.
* sys (enterprise): Add sys/billing/certificates API endpoint to retrieve the number of issued PKI certificates.
* ui/activity (enterprise): Add clarifying text to explain the "Initial Usage" column will only have timestamps for clients initially used after upgrading to version 1.21
* ui/activity (enterprise): Allow manual querying of client usage if there is a problem retrieving the license start time.
* ui/activity (enterprise): Reduce requests to the activity export API by only fetching new data when the dashboard initially loads or is manually refreshed.
* ui/activity (enterprise): Support filtering months dropdown by ISO timestamp or display value.
* ui/activity: Display total instead of new monthly clients for HCP managed clusters
* ui/pki: Adds support to configure `server_flag`, `client_flag`, `code_signing_flag`, and `email_protection_flag` parameters for creating/updating a role.

BUG FIXES:

* activity (enterprise): sys/internal/counters/activity outputs the correct mount type when called from a non root namespace
* auth/approle (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/aws (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/cert (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/github (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/ldap (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/okta (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/radius (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/scep (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/userpass (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth: fixed panic when supplying integer as a lease_id in renewal.
* core/rotation: avoid shifting timezones by ignoring cron.SpecSchedule
* core: interpret all new rotation manager rotation_schedules as UTC to avoid inadvertent use of tz-local
* secrets/azure: Ensure proper installation of the Azure enterprise secrets plugin.
* secrets/pki: Return error when issuing/signing certs whose NotAfter is before NotBefore or whose validity period isn't contained by the CA's.
* ui (enterprise): Fix KV v2 not displaying secrets in namespaces.
* ui (enterprise): Fixes login form so input renders correctly when token is a preferred login method for a namespace.
* ui/pki: Fixes certificate parsing of the `key_usage` extension so details accurately reflect certificate values.
* ui/pki: Fixes creating and updating a role so `basic_constraints_valid_for_non_ca` is correctly set.
* ui: Fix KV v2 metadata list request failing for policies without a trailing slash in the path.
* ui: Resolved a regression that prevented users with create and update permissions on KV v1 secrets from opening the edit view. The UI now correctly recognizes these capabilities and allows editing without requiring full read access.
* ui: Update LDAP accounts checked-in table to display hierarchical LDAP libraries
* ui: Update LDAP library count to reflect the total number of nodes instead of number of directories

## 1.21.0
### October 22, 2025

SECURITY:

* auth/aws: fix an issue where a user may be able to bypass authentication to Vault due to incorrect caching of the AWS client
* auth/ldap: fix MFA/TOTP enforcement bypass when username_as_alias is enabled.
* core: Update github.com/hashicorp/go-getter to fix security vulnerability GHSA-wjrx-6529-hcj3.
* core: Update github.com/ulikunitz/xz to fix security vulnerability GHSA-25xm-hr59-7c27.
* ui: disable scarf analytics for ui builds

CHANGES:

* Secrets Recovery (enterprise): Deprecate the `recover_snapshot_id` query parameter to pass the snapshot ID for recover operations, in favor of a `X-Vault-Recover-Snapshot-Id` header. Vault will still accept the query parameter for backward compatibility. Also support setting the HTTP method to `RECOVER` for recover operations, in addition to `POST` and `PUT`.
* activity: Renamed `timestamp` in export API response to `token_creation_time`.
* auth/alicloud: Update plugin to [v0.22.0](https://github.com/hashicorp/vault-plugin-auth-alicloud/releases/tag/v0.22.0)
* auth/azure: Update plugin to [v0.22.0](https://github.com/hashicorp/vault-plugin-auth-azure/releases/tag/v0.22.0)
* auth/cf: Update plugin to [v0.22.0](https://github.com/hashicorp/vault-plugin-auth-cf/releases/tag/v0.22.0)
* auth/gcp: Update plugin to [v0.22.0](https://github.com/hashicorp/vault-plugin-auth-gcp/releases/tag/v0.22.0)
* auth/jwt: Update plugin to [v0.25.0](https://github.com/hashicorp/vault-plugin-auth-jwt/releases/tag/v0.25.0)
* auth/kerberos: Update plugin to [v0.16.0](https://github.com/hashicorp/vault-plugin-auth-kerberos/releases/tag/v0.16.0)
* auth/kubernetes: Update plugin to [v0.23.0](https://github.com/hashicorp/vault-plugin-auth-kubernetes/releases/tag/v0.23.0)
* auth/oci: Update plugin to [v0.20.0](https://github.com/hashicorp/vault-plugin-auth-oci/releases/tag/v0.20.0)
* auth/saml: Update plugin to [v0.7.0](https://github.com/hashicorp/vault-plugin-auth-saml/releases/tag/v0.7.0)
* core: Updates post-install script to print updated license information
* database/couchbase: Update plugin to [v0.15.0](https://github.com/hashicorp/vault-plugin-database-couchbase/releases/tag/v0.15.0)
* database/elasticsearch: Update plugin to [v0.19.0](https://github.com/hashicorp/vault-plugin-database-elasticsearch/releases/tag/v0.19.0)
* database/mongodbatlas: Update plugin to [v0.16.0](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/releases/tag/v0.16.0)
* database/redis-elasticache: Update plugin to [v0.8.0](https://github.com/hashicorp/vault-plugin-database-redis-elasticache/releases/tag/v0.8.0)
* database/redis: Update plugin to [v0.7.0](https://github.com/hashicorp/vault-plugin-database-redis/releases/tag/v0.7.0)
* database/snowflake: Update plugin to [v0.15.0](https://github.com/hashicorp/vault-plugin-database-snowflake/releases/tag/v0.15.0)
* http: Add JSON configurable limits to HTTP handling for JSON payloads: `max_json_depth`, `max_json_string_value_length`, `max_json_object_entry_count`, `max_json_array_element_count`.
* http: Evaluate rate limit quotas before checking JSON limits during request handling.
* policies: change list comparison to allowed_parameters and denied_parameters from "exact match" to "contains all"
* sdk: Upgrade to go-secure-stdlib/plugincontainer@v0.4.2, which also bumps github.com/docker/docker to v28.3.3+incompatible
* secrets/alicloud: Update plugin to [v0.21.0](https://github.com/hashicorp/vault-plugin-secrets-alicloud/releases/tag/v0.21.0)
* secrets/azure: Update azure enterprise secrets plugin to include static roles.
* secrets/azure: Update plugin to [v0.23.0](https://github.com/hashicorp/vault-plugin-secrets-azure/releases/tag/v0.23.0)
* secrets/gcp: Update plugin to [v0.23.0](https://github.com/hashicorp/vault-plugin-secrets-gcp/releases/tag/v0.23.0)
* secrets/kubernetes: Update plugin to [v0.12.0](https://github.com/hashicorp/vault-plugin-secrets-kubernetes/releases/tag/v0.12.0)
* secrets/kv: Update plugin to [v0.25.0](https://github.com/hashicorp/vault-plugin-secrets-kv/releases/tag/v0.25.0)
* secrets/mongodbatlas: Update plugin to [v0.16.0](https://github.com/hashicorp/vault-plugin-secrets-mongodbatlas/releases/tag/v0.16.0)
* secrets/openldap: Update plugin to [v0.17.0](https://github.com/hashicorp/vault-plugin-secrets-openldap/releases/tag/v0.17.0)
* secrets/terraform: Update plugin to [v0.13.0](https://github.com/hashicorp/vault-plugin-secrets-terraform/releases/tag/v0.13.0)
* ui/client-counts: removes tabs for each client count type and adds split view for counts per type in overview stacked bar chart
* ui: Add client count attribution for the full billing period to the client counts overview table
* ui: Remove namespace context filter for activity in client count dashboard

FEATURES:

* **AES-CBC in Transit** (Enterprise): Add support for encryption and decryption with AES-CBC in the Transit Secrets Engine.
* **KV v2 Version Attribution**: Vault now includes attribution metadata for
versioned KV secrets. This allows lookup of attribution information for each
version of KV v2 secrets from CLI and API.
* **Login MFA TOTP Self-Enrollment (Enterprise)**: Simplify creation of login MFA TOTP credentials for users, allowing them to self-enroll MFA TOTP using a QR code (TOTP secret) generated during login. The new functionality is configurable on the TOTP login MFA method configuration screen and via the `enable_self_enrollment` parameter in the API.
* **Plugin Downloads**: Support automatically downloading official HashiCorp secret and auth plugins from releases.hashicorp.com (beta)
* **Post-Quantum Cryptography Support**: Experimental support for PQC signatures with ML-DSA in Transit.
* **Post-Quantum Cryptography Support**: Experimental support for PQC signatures with SLH-DSA in Transit.
* **SPIFFE Authentication Plugin (enterprise)**: Add support to authenticate to Vault using JWT and x509 based SPIFFE IDs.
* **SSH Key Signing Improvements ** (Enterprise): Add support for using managed keys to sign SSH keys in the SSH secrets engine.
* **Secret Recovery from Snapshot (enterprise)**: Adds a framework to load an integrated storage snapshot into Vault and read, list, and recover KV v1 and cubbyhole secrets from the snapshot.
* **UI Client List Explorer (Enterprise)**: Adds ability to view and filter client IDs and metadata by namespace, mount path, or mount type for a billing period.
* **UI Secrets Recovery (Enterprise)**: Allows end users to recover single KV v1 secrets, Cubbyhole secrets, or Database static roles from a loaded snapshot if the secrets were changed or deleted in error. Automatic snapshot configurations can now automatically load the snapshot to Vault itself, making it available for recovery. Snapshot management permissions are separated from recovery permissions so that recovery operations can be delegated but controlled.
* **UI: Secret Engine Tune Support**: Add support for updating secret engine mount configuration via the Tune endpoint
* **Vault PKI SCEP Server (Enterprise)**: Support for the Simple Certificate Enrollment Protocol (SCEP) has been added to the Vault PKI Plugin. This allows standard SCEP clients to request certificates from a Vault server with no knowledge of Vault APIs.

IMPROVEMENTS:

* ui/activity: Updates running total stats to be displayed via a donut chart.
* Plugin Downloads (enterprise): add CLI `-download` option for plugin register (beta)
* Raft: Auto-join will now allow you to enforce IPv4 on networks that allow IPv6 and dual-stack enablement, which is on by default in certain regions.
* Secrets Recovery (enterprise): Support recovering items from a snapshot to a new path in the live cluster. By calling the `vault recover` command with a `-from` flag, users can specify the path of the item in the snapshot.
* Secrets Sync (enterprise): add `enterprise_url` field to enable support for self-hosted GitHub Enterprise Server instances.
* activity (enterprise): Add a cumulative namespace client count API at `sys/internal/counters/activity/cumulative`. For each namespace in the response it returns the sum of its own client counts and that of all its child namespaces.
* activity: The [activity export API](https://developer.hashicorp.com/vault/api-docs/system/internal-counters#activity-export) response now includes a new timestamp that denotes the first time the client was used within the specified query period.
* api (sys/utilization-report): Added namespace filter and more granularity for secret sync data in the response.
* api: Add new logical client request interfaces for read, write, delete, list operations.
* audit: Add additional verifications to the target of file audit sinks.
* auth/approle (enterprise): Add ability to specify custom alias metadata via new role creation parameter `alias_metadata`.
* auth/aws (enterprise): Add ability to specify custom alias metadata via new role creation parameter `alias_metadata`.
* auth/cert (enterprise): Add ability to specify custom alias metadata via new CA certificate role creation parameter `alias_metadata`.
* auth/cert: Add allowed_organizations support
* auth/cert: Support RFC 9440 colon-wrapped Base64 certificates in `x_forwarded_for_client_cert_header`, to fix TLS certificate auth errors with Google Cloud Application Load Balancer.
* auth/cert: test non-CA cert equality on login matching instead of individual fields.
* auth/github (enterprise): Add ability to specify custom alias metadata via new configuration parameter `alias_metadata`.
* auth/ldap (enterprise): Add ability to specify custom alias metadata via new role configuration parameter `alias_metadata`.
* auth/ldap: Introduces an option to connect to an alternative LDAP URL for root credential rotation, in cases where it differs from the configured LDAP URL.
* auth/ldap: add explicit logging to rotations in ldap
* auth/okta (enterprise): Add ability to specify custom alias metadata via new role configuration parameter `alias_metadata`.
* auth/radius (enterprise): Add ability to specify custom alias metadata via new role configuration parameter `alias_metadata`.
* auth/scep (enterprise): Add ability to specify custom alias metadata via new role creation parameter `alias_metadata`.
* auth/userpass (enterprise): Add ability to specify custom alias metadata via new user creation parameter `alias_metadata`.
* cli (enterprise): Add a `-force` flag to `vault operator raft snapshot unload` command to force deletion of a loaded snapshot.
* core (enterprise): Allow setting of an entropy source on password generation
policies, and with it the selection of "seal" to use entropy augmentation.
* core (enterprise): add ability to get time remaining until rotation from rotation manager
* core (enterprise): add support for new pki-only license feature
* core (enterprise): improve rotation manager logging to include specific lines for rotation success and failure
* core/metrics: Reading and listing from a snapshot are now tracked via the `vault.route.read-snapshot.{mount_point}` and `vault.route.list-snapshot.{mount_point}` metrics.
* core/snapshot-load (enterprise): Add a `force` query parameter to the `DELETE sys/storage/raft/snapshot-load/{snapshot_id}` endpoint to allow for forced deletion of snapshots. This is useful when the snapshot is in a state that prevents normal deletion, such as being in the process of loading.
* license utilization reporting (enterprise): Add metrics for the number of issued PKI certificates.
* openapi: Add OpenAPI support for secret recovery operations.
* openapi: Add openapi response definitions to `sys/internal/counters/activity/*` endpoints.
* plugins: Clarify usage of sha256, command, and version for plugin registration of binary or artifact with API and CLI. Introduce new RegisterPluginDetailed and RegisterPluginWtihContextDetailed functions to API client to propagate response along with error, and mark RegisterPlugin and RegisterPluginWithContext as deprecated.
* proxy/cache (enterprise): Vault Proxy will now use vault_index on events to be able to update cached static secrets from performance secondaries without needing to be forwarded. This will take precedence over attempting to forward the request to the primary.
* sdk: add stub code for retrieving rotation schedule information
* secrets/database (enterprise): Add support for reading, listing, and recovering static roles from a loaded snapshot. Also add support for reading static credentials from a loaded snapshot.
* secrets/database: Add PSC support for GCP CloudSQL MySQL and Postgresql
* secrets/database: Add PrivateIP support for MySQL
* secrets/database: Add root rotation support for Snowflake database secrets engines using key-pair credentials.
* secrets/database: log password rotation success (info) and failure (error). Some relevant log lines have been updated to include "path" fields.
* secrets/kmip (enterprise): Update various third party dependencies.
* secrets/pki (enterprise): add integrations/guardium configuration endpoint.
* secrets/pki (enterprise): add new batch/certs endpoint to allow multiple certificates to be fetched at once.
* secrets/pki (enterprise): enable separately-configured logging for SCEP-enrollment.
* secrets/pki: Add the digest OID when logging SCEP digest mismatch errors.
* secrets/ssh: Add support for recovering the SSH plugin CA from a loaded snapshot (enterprise only).
* secrets/transform (enterprise): Update various third party dependencies.
* secrets/transit: add logging on both success and failure of key rotation
* storage/raft (enterprise): Add `autoload_enabled` option to raft automated snapshot configurations. When enabled, this option will automatically load raft snapshots into Vault, which can then be used for recovery operations.
* sys (enterprise): Add sys/billing/certificates API endpoint to retrieve the number of issued PKI certificates.
* ui/activity (enterprise): Add clarifying text to explain the "Initial Usage" column will only have timestamps for clients initially used after upgrading to version 1.21
* ui/activity (enterprise): Reduce requests to the activity export API by only fetching new data when the dashboard initially loads or is manually refreshed.
* ui/activity (enterprise): Support filtering months dropdown by ISO timestamp or display value.
* ui/activity: Adds filtering by month to the Client Count dashboard to link client counts to specific client IDs from the export API
* ui/auth: the role field on the OIDC login form now auto-fills from the `role` URL query string parameter
* ui/auth: the role field on the SAML login form now auto-fills from the `role` URL query string parameter
* ui/secrets: Display the plugin version on the secret engine list view. Move KV's version to a tooltip that appears when hovering over the engine's name.
* ui/secrets: Updated filters on secret engines list to sort by path, engine type and version
* ui: Add `namespace_path`, `mount_path` and `mount_type` filters to attribution table
* ui: Enhanced secret engine selection dynamically displays all available plugins from the plugin catalog.
* ui: Format multiline API error messages to render as bulleted lists.
* ui: Use the Helios Design System Code Block component for all readonly code editors and use its Code Editor component for all other code editors

DEPRECATIONS:

* core: disallow usage of duplicate attributes in HCL configuration files and policy definitions, which were already deprecated. For now those errors can be suppressed back to warnings by setting the environment variable VAULT_ALLOW_PENDING_REMOVAL_DUPLICATE_HCL_ATTRIBUTES.

BUG FIXES:

* activity (enterprise): Fix `development_cluster` setting being overwritten on performance secondaries upon cluster reload.
* activity (enterprise): sys/internal/counters/activity outputs the correct mount type when called from a non root namespace
* agent/template: Fixed issue where templates would not render correctly if namespaces was provided by config, and the namespace and mount path of the secret were the same.
* auth/approle (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/aws (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/cert (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/cert: Recover from partially populated caches of trusted certificates if one or more certificates fails to load.
* auth/github (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/ldap (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/okta (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/radius (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/scep (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth/scep (enterprise): enforce the token_bound_cidrs role parameter within SCEP roles
* auth/spiffe: Address an issue updating a role with overlapping workload_id_pattern values it previously contained.
* auth/userpass (enterprise): Role parameter `alias_metadata` now populates alias custom metadata field instead of alias metadata.
* auth: fixed panic when supplying integer as a lease_id in renewal.
* auth: update alias lookahead to respect username case for LDAP and username/password
* auto-reporting (enterprise): Clarify debug logs to accurately reflect when automated license utilization reporting is enabled or disabled, especially since manual reporting is always initialized.
* core (enterprise): Avoid duplicate seal rewrapping, and ensure that cluster secondaries rewrap after a seal migration.
* core (enterprise): fix a bug where issuing a token in a namespace used root auth configuration instead of namespace auth configuration
* core/activitylog (enterprise): Fix nil panic when trying reload census manager before activity log is setup.
* core/metrics: Add service name prefix for core HA metrics to avoid duplicate, zero-value metrics.
* core/seal (enterprise): Fix a bug that caused the seal rewrap process to abort in the presence of partially sealed entries.
* core/seal: When Seal-HA is enabled, make it an error to persist the barrier keyring when not all seals are healthy.  This prevents the possibility of failing to unseal when a different subset of seals are healthy than were healthy at last write.
* core: Fixed issue where under certain circumstances the rotation manager would spawn goroutines indefinitely.
* core: Role based quotas now work for cert auth
* core: interpret all new rotation manager rotation_schedules as UTC to avoid inadvertent use of tz-local
* core: resultant-acl now merges segment-wildcard (`+`) paths with existing prefix rules in `glob_paths`, so clients receive a complete view of glob-style permissions. This unblocks UI sidebar navigation checks and namespace access banners.
* default-auth: Fix bug where listing default-auth configurations caused panic during auditing.
* gcs: fix failed locking due to updated library error checks
* identity/mfa: revert cache entry change from #31217 and document cache entry values
* kmip (enterprise): Fix a panic that can happen when a KMIP client makes a request before the Vault server has finished unsealing.
* mongodb: fix mongodb connection issue when using TLS client + username/password authentication
* plugins: Fix panics that can occur when a plugin audits a request or response before the Vault server has finished unsealing.
* product usage reporting (enterprise): Clarify debug logs to accurately reflect when anonymous product usage reporting is enabled or disabled, especially since manual reporting is always initialized.
* raft (enterprise): auto-join will now work in regions that do not support dual-stack
* raft/autopilot: Fixes an issue with enterprise redundancy zones where, if the leader was in a redundancy zone and that leader becomes unavailable, the node would become an unzoned voter. This can artificially inflate the required number of nodes for quorum, leading to a situation where the cluster cannot recover if another leader subsequently becomes unavailable. Vault will now keep an unavailable node in its last known redundancy zone as a non-voter.
* replication (enterprise): Fix bug where group updates fail when processed on a standby node in a PR secondary cluster.
* replication (enterprise): Fix bug with mount invalidations consuming excessive memory.
* secrets-sync (enterprise): GCP locational KMS keys are no longer incorrectly removed when the location name is all lowercase.
* secrets-sync (enterprise): Unsyncing secret-key granularity associations will no longer give a misleading error about a failed unsync operation that did indeed succeed.
* secrets/azure: Ensure proper installation of the Azure enterprise secrets plugin.
* secrets/database/postgresql: Support for multiline statements in the `rotation_statements` field.
* secrets/database: respect the escaping/disable_escaping state when using self-managed static roles
* secrets/transit: Fix error when using ed25519 keys that were imported with derivation enabled
* sentinel (enterprise): Fix a Sentinel bug, where the soft-mandatory policy override would not work in overriding request denial. Now, Vault correctly allows requests when the policy override flag is set. Previously, requests were denied even if an override was explicitly set. Error messaging for denied requests is now clearer and more actionable.
* sys/mounts: enable unsetting allowed_response_headers
* ui (enterprise): Fixes login form so input renders correctly when token is a preferred login method for a namespace.
* ui: Fix DR secondary view from not loading/transitioning.
* ui: Fix kv v2 overview page from erroring if a user does not have access to the /subkeys endpoint and the policy check fails.
* ui: Fix page loading error when users navigate away from identity entities and groups list views.
* ui: Fix regression in 1.20.0 to properly set namespace context for capabilities checks
* ui: Fix selecting multiple namespaces in the namespace picker when the path contains matching nodes
* ui: Fixes UI login settings list page which was not rendering rules with an underscore in the name.
* ui: Include user's root namespace in the namespace picker if it's a namespace other than the actual root ("")
* ui: Revert camelizing of parameters returned from `sys/internal/ui/mounts` so mount paths match serve value
* ui: Fixes permissions for hiding and showing sidebar navigation items for policies that include special characters: `+`, `*`

## 1.20.7 Enterprise
### January 07, 2026

CHANGES:

* auth/oci: bump plugin to v0.19.1
* go: bump go version to 1.25.5
* packaging: Container images are now exported using a compressed OCI image layout.
* packaging: UBI container images are now built on the UBI 10 minimal image.
* secrets/azure: Update plugin to [v0.22.1](https://github.com/hashicorp/vault-plugin-secrets-azure/releases/tag/v0.22.1). Improves retry handling during Azure application and service principal creation to reduce transient failures.
* storage: Upgrade aerospike client library to v8.

IMPROVEMENTS:

* core: check rotation manager queue every 5 seconds instead of 10 seconds to improve responsiveness.
* go: update to golang/x/crypto to v0.45.0 to resolve GHSA-f6x5-jh6r-wrfv, GHSA-j5w8-q4qc-rx2x, GO-2025-4134 and GO-2025-4135.
* rotation: Ensure rotations for shared paths only execute on the Primary cluster's active node. Ensure rotations for local paths execute on the cluster-local active node.
* sdk/rotation: Prevent rotation attempts on read-only storage
* secrets-sync (enterprise): Added support for a boolean force_delete flag (default: false). 
When set to true, this flag allows deletion of a destination even if its associations cannot be unsynced. 
This option should be used only as a last-resort deletion mechanism, as any secrets already synced to the external provider will remain orphaned and require manual cleanup.

BUG FIXES:

* auth/approle (enterprise): Fixed bug that prevented periodic tidy running on performance secondary.
* core/activitylog (enterprise): Resolve a stability issue where Vault Enterprise could encounter a panic during month-end billing activity rollover.
* http: skip JSON limit parsing on cluster listener.
* quotas: Vault now protects plugins with ResolveRole operations from panicking.
on quota creation.
* replication (enterprise): fix rare panic due to race when enabling a secondary with Consul storage.
* rotation: Fix a bug where a performance secondary would panic if a write was made to a local mount.
* secret-sync (enterprise): Improved unsync error handling by treating cases where the destination no longer exists as successful.
* secrets-sync (enterprise): Corrected a bug where the deletion of the latest KV-V2 secret version caused the associated external secret to be deleted entirely. The sync job now implements a version fallback mechanism to find and sync the highest available active version, ensuring continuity and preventing the unintended deletion of the external secret resource.
* ui/kvv2 (enterprise): Fixes listing stale secrets when switching between namespaces that have KV v2 engines with the same mount path.
* ui/pki: Fix handling of values that contain commas in list fields like `crl_distribution_points`.

## 1.20.6 Enterprise
### November 19, 2025

CHANGES:

* auth/kubernetes: Update plugin to [v0.22.5](https://github.com/hashicorp/vault-plugin-auth-kubernetes/releases/tag/v0.22.5)
* core: Bump Go version to 1.24.10
* policies: add VAULT_NEW_PER_ELEMENT_MATCHING_ON_LIST env var to adopt new "contains all" list matching behavior on
allowed_parameters and denied_parameters

IMPROVEMENTS:

* Update github.com/dvsekhvalnov/jose2go to fix security vulnerability CVE-2025-63811.
* auth/ldap: Require non-empty passwords on login command to prevent unauthenticated access to Vault.
* ui/pki: Adds support to configure `server_flag`, `client_flag`, `code_signing_flag`, and `email_protection_flag` parameters for creating/updating a role.

BUG FIXES:

* core/activitylog (enterprise): Fix nil panic when trying reload census manager before activity log is setup.
* core/rotation: avoid shifting timezones by ignoring cron.SpecSchedule
* ui/pki: Fixes certificate parsing of the `key_usage` extension so details accurately reflect certificate values.
* ui/pki: Fixes creating and updating a role so `basic_constraints_valid_for_non_ca` is correctly set.
* ui: Resolved a regression that prevented users with create and update permissions on KV v1 secrets from opening the edit view. The UI now correctly recognizes these capabilities and allows editing without requiring full read access.
* ui: Update LDAP accounts checked-in table to display hierarchical LDAP libraries
* ui: Update LDAP library count to reflect the total number of nodes instead of number of directories
* ui: remove unnecessary 'credential type' form input when generating AWS secrets

## 1.20.5 Enterprise
### October 22, 2025

SECURITY:

* auth/aws: fix an issue where a user may be able to bypass authentication to Vault due to incorrect caching of the AWS client
* ui: disable scarf analytics for ui builds

CHANGES:

* core: Bump Go version to 1.24.9.
* http: Evaluate rate limit quotas before checking JSON limits during request handling.

IMPROVEMENTS:

* core/metrics: Reading and listing from a snapshot are now tracked via the `vault.route.read-snapshot.{mount_point}` and `vault.route.list-snapshot.{mount_point}` metrics.
* secrets/database: Add root rotation support for Snowflake database secrets engines using key-pair credentials.

BUG FIXES:

* activity (enterprise): sys/internal/counters/activity outputs the correct mount type when called from a non root namespace
* auth: fixed panic when suppling integer as a lease_id in renewal.
* core (enterprise): Avoid duplicate seal rewrapping, and ensure that cluster secondaries rewrap after a seal migration.
* core: interpret all new rotation manager rotation_schedules as UTC to avoid inadvertent use of tz-local
* core: resultant-acl now merges segment-wildcard (`+`) paths with existing prefix rules in `glob_paths`, so clients receive a complete view of glob-style permissions. This unblocks UI sidebar navigation checks and namespace access banners.
* secrets/database: respect the escaping/disable_escaping state when using self-managed static roles
* sentinel (enterprise): Fix a Sentinel bug, where the soft-mandatory policy override would not work in overriding request denial. Now, Vault correctly allows requests when the policy override flag is set. Previously, requests were denied even if an override was explicitly set. Error messaging for denied requests is now clearer and more actionable.
* ui (enterprise): Fixes login form so input renders correctly when token is a preferred login method for a namespace.
* ui: Fixes permissions for hiding and showing sidebar navigation items for policies that include special characters: `+`, `*`

## 1.20.4
### September 24, 2025

SECURITY:

* core: Update github.com/ulikunitz/xz to fix security vulnerability GHSA-25xm-hr59-7c27. ([ce4b4264](https://github.com/hashicorp/vault/commit/ce4b42642403f30370dde0a39e9a04991c387291))

CHANGES:

* database/snowflake: Update plugin to [v0.14.2](https://github.com/hashicorp/vault-plugin-database-snowflake/releases/tag/v0.14.2) ([9f06df77](https://github.com/hashicorp/vault/commit/9f06df77b48024350fbc67b7d9a7eaaa5d0022fa))

IMPROVEMENTS:

* Raft: Auto-join will now allow you to enforce IPv4 on networks that allow IPv6 and dual-stack enablement, which is on by default in certain regions. ([1fd38796](https://github.com/hashicorp/vault/commit/1fd38796639ed861fc2fa1b58f138caad8fc0950))
* auth/cert: Support RFC 9440 colon-wrapped Base64 certificates in `x_forwarded_for_client_cert_header`, to fix TLS certificate auth errors with Google Cloud Application Load Balancer. [[GH-31501](https://github.com/hashicorp/vault/pull/31501)]
* secrets/database (enterprise): Add support for reading, listing, and recovering static roles from a loaded snapshot. Also add support for reading static credentials from a loaded snapshot. ([24cd1aa5](https://github.com/hashicorp/vault/commit/24cd1aa5961bfbed396251aebd3490dcfc7a106f))
* secrets/ssh: Add support for recovering the SSH plugin CA from a loaded snapshot (enterprise only). ([0087af9d](https://github.com/hashicorp/vault/commit/0087af9da59692351e7a3c3f5269af9de082a52e))

BUG FIXES:

* auth/cert: Recover from partially populated caches of trusted certificates if one or more certificates fails to load. [[GH-31438](https://github.com/hashicorp/vault/pull/31438)]
* core: Role based quotas now work for cert auth ([fc775dea](https://github.com/hashicorp/vault/commit/fc775deacee3ab2dd956b8ab7ab64847601c1685))
* sys/mounts: enable unsetting allowed_response_headers [[GH-31555](https://github.com/hashicorp/vault/pull/31555)]
* ui: Fix page loading error when users navigate away from identity entities and groups list views. ([81170963](https://github.com/hashicorp/vault/commit/8117096364d5fbb541124a6e057cd11b24eaa6f3))


## 1.20.3
### August 28, 2025

FEATURES:

* **IBM RACF Static Role Password Phrase Management (Enterprise)**: Add support for static role password phrase management to the LDAP secrets engine. (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/184)

SECURITY:

* core: Update github.com/hashicorp/go-getter to fix security vulnerability GHSA-wjrx-6529-hcj3. ([8b3a9ce1](https://github.com/hashicorp/vault/commit/8b3a9ce1f651932559a129a7889243d24127cee2))

CHANGES:

* core: Bump Go version to 1.24.6. ([ce56e14e](https://github.com/hashicorp/vault/commit/ce56e14e7466ae80e05d11a83c8f41db0f4653be))
* http: Add JSON configurable limits to HTTP handling for JSON payloads: `max_json_depth`, `max_json_string_value_length`, `max_json_object_entry_count`, `max_json_array_element_count`. [[GH-31069](https://github.com/hashicorp/vault/pull/31069)]
* sdk: Upgrade to go-secure-stdlib/plugincontainer@v0.4.2, which also bumps github.com/docker/docker to v28.3.3+incompatible ([8f172169](https://github.com/hashicorp/vault/commit/8f1721697bba123117f4f98dae4154ef9fe614e5))
* secrets/openldap (enterprise): update plugin to v0.16.1

IMPROVEMENTS:

* auth/ldap: add explicit logging to rotations in ldap [[GH-31401](https://github.com/hashicorp/vault/pull/31401)]
* core (enterprise): improve rotation manager logging to include specific lines for rotation success and failure
* secrets/database: log password rotation success (info) and failure (error). Some relevant log lines have been updated to include "path" fields. [[GH-31402](https://github.com/hashicorp/vault/pull/31402)]
* secrets/transit: add logging on both success and failure of key rotation [[GH-31420](https://github.com/hashicorp/vault/pull/31420)]
* ui: Use the Helios Design System Code Block component for all readonly code editors and use its Code Editor component for all other code editors [[GH-30188](https://github.com/hashicorp/vault/pull/30188)]

BUG FIXES:

* core (enterprise): fix a bug where issuing a token in a namespace used root auth configuration instead of namespace auth configuration
* core/metrics: Add service name prefix for core HA metrics to avoid duplicate, zero-value metrics. ([91e5f443](https://github.com/hashicorp/vault/commit/91e5f44315fb52c37b54e8b0eece1b4390665cc3))
* core/seal: When Seal-HA is enabled, make it an error to persist the barrier
keyring when not all seals are healthy.  This prevents the possibility of
failing to unseal when a different subset of seals are healthy than were
healthy at last write. ([bbe64227](https://github.com/hashicorp/vault/commit/bbe64227c586cb34f73d9ae8025398f24aa7e12d))
* raft (enterprise): auto-join will now work in regions that do not support dual-stack ([c66baf5e](https://github.com/hashicorp/vault/commit/c66baf5ee1ee9320daa6af5528cb2f250f2a0f3a))
* raft/autopilot: Fixes an issue with enterprise redundancy zones where, if the leader was in a redundancy zone and that leader becomes unavailable, the node would become an unzoned voter. This can artificially inflate the required number of nodes for quorum, leading to a situation where the cluster cannot recover if another leader subsequently becomes unavailable. Vault will now keep an unavailable node in its last known redundancy zone as a non-voter. [[GH-31443](https://github.com/hashicorp/vault/pull/31443)]
* replication (enterprise): Fix bug where group updates fail when processed on a
standby node in a PR secondary cluster.
* secrets-sync (enterprise): GCP locational KMS keys are no longer incorrectly removed when the location name is all lowercase.
* secrets/database/postgresql: Support for multiline statements in the `rotation_statements` field. [[GH-31442](https://github.com/hashicorp/vault/pull/31442)]
* ui: Fix DR secondary view from not loading/transitioning. [[GH-31478](https://github.com/hashicorp/vault/pull/31478)]

## 1.20.2
### August 06, 2025

SECURITY:

* auth/ldap: fix MFA/TOTP enforcement bypass when username_as_alias is enabled [[GH-31427](https://github.com/hashicorp/vault/pull/31427),[HCSEC-2025-20](https://discuss.hashicorp.com/t/hcsec-2025-20-vault-ldap-mfa-enforcement-bypass-when-using-username-as-alias/76092)].

BUG FIXES:

* agent/template: Fixed issue where templates would not render correctly if namespaces was provided by config, and the namespace and mount path of the secret were the same. [[GH-31392](https://github.com/hashicorp/vault/pull/31392)]
* identity/mfa: revert cache entry change from #31217 and document cache entry values [[GH-31421](https://github.com/hashicorp/vault/pull/31421)]

## 1.20.1
### July 25, 2025

SECURITY:

* audit: **breaking change** privileged vault operator may execute code on the underlying host (CVE-2025-6000). Vault will not unseal if the only configured file audit device has executable permissions (e.g., 0777, 0755). See recent [breaking change](https://developer.hashicorp.com/vault/docs/updates/important-changes#breaking-changes) docs for more details. [[GH-31211](https://github.com/hashicorp/vault/pull/31211),[HCSEC-2025-14](https://discuss.hashicorp.com/t/hcsec-2025-14-privileged-vault-operator-may-execute-code-on-the-underlying-host/76033)]
* auth/userpass: timing side-channel in vault's userpass auth method (CVE-2025-6011)[HCSEC-2025-15](https://discuss.hashicorp.com/t/hcsec-2025-15-timing-side-channel-in-vault-s-userpass-auth-method/76034)
* core/login: vault userpass and ldap user lockout bypass (CVE-2025-6004). update alias lookahead to respect username case for LDAP and username/password. [[GH-31352](https://github.com/hashicorp/vault/pull/31352),[HCSEC-2025-16](https://discuss.hashicorp.com/t/hcsec-2025-16-vault-userpass-and-ldap-user-lockout-bypass/76035)]
* secrets/totp: vault totp secrets engine code reuse (CVE-2025-6014) [[GH-31246](https://github.com/hashicorp/vault/pull/31246),[HCSEC-2025-17](https://discuss.hashicorp.com/t/hcsec-2025-17-vault-totp-secrets-engine-code-reuse/76036)]
* auth/cert: vault certificate auth method did not validate common name for non-ca certificates (CVE-2025-6037). test non-CA cert equality on login matching instead of individual fields. [[GH-31210](https://github.com/hashicorp/vault/pull/31210),[HCSEC-2025-18](https://discuss.hashicorp.com/t/hcsec-2025-18-vault-certificate-auth-method-did-not-validate-common-name-for-non-ca-certificates/76037)]
* core/mfa: vault login mfa bypass of rate limiting and totp token reuse (CVE-2025-6015) [[GH-31217](https://github.com/hashicorp/vault/pull/31297),[HCSEC-2025-19](https://discuss.hashicorp.com/t/hcsec-2025-19-vault-login-mfa-bypass-of-rate-limiting-and-totp-token-reuse/76038)]

FEATURES:

* **Post-Quantum Cryptography Support**: Experimental support for PQC signatures with SLH-DSA in Transit.

IMPROVEMENTS:

* Plugin Downloads (enterprise): add CLI `-download` option for plugin register (beta)
* openapi: Add OpenAPI support for secret recovery operations. [[GH-31331](https://github.com/hashicorp/vault/pull/31331)]
* plugins: Clarify usage of sha256, command, and version for plugin registration of binary or artifact with API and CLI. Introduce new RegisterPluginDetailed and RegisterPluginWtihContextDetailed functions to API client to propagate response along with error, and mark RegisterPlugin and RegisterPluginWithContext as deprecated. [[GH-30811](https://github.com/hashicorp/vault/pull/30811)]
* secrets/pki (enterprise): enable separately-configured logging for SCEP-enrollment.
* secrets/pki: Add the digest OID when logging SCEP digest mismatch errors. [[GH-31232](https://github.com/hashicorp/vault/pull/31232)]

BUG FIXES:

* activity (enterprise): Fix `development_cluster` setting being overwritten on performance secondaries upon cluster reload. [[GH-31223](https://github.com/hashicorp/vault/pull/31223)]
* auth/scep (enterprise): enforce the token_bound_cidrs role parameter within SCEP roles
* auto-reporting (enterprise): Clarify debug logs to accurately reflect when automated license utilization reporting is enabled or disabled, especially since manual reporting is always initialized.
* core/seal (enterprise): Fix a bug that caused the seal rewrap process to abort in the presence of partially sealed entries.
* kmip (enterprise): Fix a panic that can happen when a KMIP client makes a request before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* plugins: Fix panics that can occur when a plugin audits a request or response before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* product usage reporting (enterprise): Clarify debug logs to accurately reflect when anonymous product usage reporting is enabled or disabled, especially since manual reporting is always initialized.
* replication (enterprise): Fix bug with mount invalidations consuming excessive memory.
* secrets-sync (enterprise): Unsyncing secret-key granularity associations will no longer give a misleading error about a failed unsync operation that did indeed succeed.
* secrets/gcp: Update to vault-plugin-secrets-gcp@v0.22.1 to address more eventual consistency issues [[GH-31350](https://github.com/hashicorp/vault/pull/31350)]
* ui: Fix capability checks for api resources with underscores to properly hide actions and dropdown items a user cannot perform [[GH-31271](https://github.com/hashicorp/vault/pull/31271)]
* ui: Fix kv v2 overview page from erroring if a user does not have access to the /subkeys endpoint and the policy check fails. [[GH-31136](https://github.com/hashicorp/vault/pull/31136)]
* ui: Fix mutation of unwrapped data when keys contain underscores [[GH-31287](https://github.com/hashicorp/vault/pull/31287)]
* ui: Fix regression in 1.20.0 to properly set namespace context for capabilities checks [[GH-31276](https://github.com/hashicorp/vault/pull/31276)]
* ui: Fix selecting multiple namespaces in the namespace picker when the path contains matching nodes [[GH-31326](https://github.com/hashicorp/vault/pull/31326)]
* ui: Fixes UI login settings list page which was not rendering rules with an underscore in the name. [[GH-31150](https://github.com/hashicorp/vault/pull/31150)]
* ui: Include user's root namespace in the namespace picker if it's a namespace other than the actual root ("") [[GH-31300](https://github.com/hashicorp/vault/pull/31300)]
* ui: Revert camelizing of parameters returned from `sys/internal/ui/mounts` so mount paths match serve value [[GH-31094](https://github.com/hashicorp/vault/pull/31094)]

## 1.20.0
### June 25, 2025

SECURITY:

* core: require a nonce when cancelling a rekey operation that was initiated within the last 10 minutes. [[GH-30794](https://github.com/hashicorp/vault/pull/30794)],[[HCSEC-2025-11](https://discuss.hashicorp.com/t/hcsec-2025-11-vault-vulnerable-to-recovery-key-cancellation-denial-of-service/75570)]
* core/identity: vault root namespace operator may elevate privileges (CVE-2025-5999). Fix string contains check in Identity APIs to be case-insensitive. [[GH-31045](https://github.com/hashicorp/vault/pull/31045),[HCSEC-2025-13](https://discuss.hashicorp.com/t/hcsec-2025-13-vault-root-namespace-operator-may-elevate-token-privileges/76032)]

CHANGES:

* UI: remove outdated and unneeded js string extensions [[GH-29834](https://github.com/hashicorp/vault/pull/29834)]
* activity (enterprise): The sys/internal/counters/activity endpoint will return actual values for new clients in the current month.
* activity (enterprise): provided values for `start_time` and `end_time` in `sys/internal/counters/activity` are aligned to the corresponding billing period.
* activity: provided value for `end_time` in `sys/internal/counters/activity` is now capped at the end of the last completed month. [[GH-30164](https://github.com/hashicorp/vault/pull/30164)]
* api: Update the default API client to check for the `Retry-After` header and, if it exists, wait for the specified duration before retrying the request. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* auth/alicloud: Update plugin to v0.21.0 [[GH-30810](https://github.com/hashicorp/vault/pull/30810)]
* auth/azure: Update plugin to v0.20.2. Login requires `resource_group_name`, `vm_name`, and `vmss_name` to match token claims [[GH-30052](https://github.com/hashicorp/vault/pull/30052)]
* auth/azure: Update plugin to v0.20.3 [[GH-30082](https://github.com/hashicorp/vault/pull/30082)]
* auth/azure: Update plugin to v0.20.4 [[GH-30543](https://github.com/hashicorp/vault/pull/30543)]
* auth/azure: Update plugin to v0.21.0 [[GH-30872](https://github.com/hashicorp/vault/pull/30872)]
* auth/azure: Update plugin to v0.21.1 [[GH-31010](https://github.com/hashicorp/vault/pull/31010)]
* auth/cf: Update plugin to v0.20.1 [[GH-30583](https://github.com/hashicorp/vault/pull/30583)]
* auth/cf: Update plugin to v0.21.0 [[GH-30842](https://github.com/hashicorp/vault/pull/30842)]
* auth/gcp: Update plugin to v0.20.2 [[GH-30081](https://github.com/hashicorp/vault/pull/30081)]
* auth/jwt: Update plugin to v0.23.2 [[GH-30431](https://github.com/hashicorp/vault/pull/30431)]
* auth/jwt: Update plugin to v0.24.1 [[GH-30876](https://github.com/hashicorp/vault/pull/30876)]
* auth/kerberos: Update plugin to v0.15.0 [[GH-30845](https://github.com/hashicorp/vault/pull/30845)]
* auth/kubernetes: Update plugin to v0.22.1 [[GH-30910](https://github.com/hashicorp/vault/pull/30910)]
* auth/oci: Update plugin to v0.19.0 [[GH-30841](https://github.com/hashicorp/vault/pull/30841)]
* auth/saml: Update plugin to v0.6.0
* core: Bump Go version to 1.24.4.
* core: Verify that the client IP address extracted from an X-Forwarded-For header is a valid IPv4 or IPv6 address [[GH-29774](https://github.com/hashicorp/vault/pull/29774)]
* database/couchbase: Update plugin to v0.14.0 [[GH-30836](https://github.com/hashicorp/vault/pull/30836)]
* database/elasticsearch: Update plugin to v0.18.0 [[GH-30796](https://github.com/hashicorp/vault/pull/30796)]
* database/mongodbatlas: Update plugin to v0.15.0 [[GH-30856](https://github.com/hashicorp/vault/pull/30856)]
* database/redis-elasticache: Update plugin to v0.7.0 [[GH-30785](https://github.com/hashicorp/vault/pull/30785)]
* database/redis: Update plugin to v0.6.0 [[GH-30797](https://github.com/hashicorp/vault/pull/30797)]
* database/snowflake: Update plugin to v0.14.0 [[GH-30748](https://github.com/hashicorp/vault/pull/30748)]
* database/snowflake: Update plugin to v0.14.1 [[GH-30868](https://github.com/hashicorp/vault/pull/30868)]
* logical/system: add ent stub for plugin catalog handling [[GH-30890](https://github.com/hashicorp/vault/pull/30890)]
* quotas/rate-limit: Round up the `Retry-After` value to the nearest second when calculating the retry delay. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* secrets/ad: Update plugin to v0.21.0 [[GH-30819](https://github.com/hashicorp/vault/pull/30819)]
* secrets/alicloud: Update plugin to v0.20.0 [[GH-30809](https://github.com/hashicorp/vault/pull/30809)]
* secrets/azure: Update plugin to v0.21.2 [[GH-30037](https://github.com/hashicorp/vault/pull/30037)]
* secrets/azure: Update plugin to v0.21.3 [[GH-30083](https://github.com/hashicorp/vault/pull/30083)]
* secrets/azure: Update plugin to v0.22.0 [[GH-30832](https://github.com/hashicorp/vault/pull/30832)]
* secrets/gcp: Update plugin to v0.21.2 [[GH-29970](https://github.com/hashicorp/vault/pull/29970)]
* secrets/gcp: Update plugin to v0.21.3 [[GH-30080](https://github.com/hashicorp/vault/pull/30080)]
* secrets/gcp: Update plugin to v0.22.0 [[GH-30846](https://github.com/hashicorp/vault/pull/30846)]
* secrets/gcpkms: Update plugin to v0.21.0 [[GH-30835](https://github.com/hashicorp/vault/pull/30835)]
* secrets/kubernetes: Update plugin to v0.11.0 [[GH-30855](https://github.com/hashicorp/vault/pull/30855)]
* secrets/kv: Update plugin to v0.24.0 [[GH-30826](https://github.com/hashicorp/vault/pull/30826)]
* secrets/mongodbatlas: Update plugin to v0.15.0 [[GH-30860](https://github.com/hashicorp/vault/pull/30860)]
* secrets/openldap: Update plugin to v0.15.2 [[GH-30079](https://github.com/hashicorp/vault/pull/30079)]
* secrets/openldap: Update plugin to v0.15.4 [[GH-30279](https://github.com/hashicorp/vault/pull/30279)]
* secrets/openldap: Update plugin to v0.16.0 [[GH-30844](https://github.com/hashicorp/vault/pull/30844)]
* secrets/terraform: Update plugin to v0.12.0 [[GH-30905](https://github.com/hashicorp/vault/pull/30905)]
* server: disable_mlock configuration option is now required for integrated storage and no longer has a default. If you are using the default value with integrated storage, you must now explicitly set disable_mlock to true or false or Vault server will fail to start. [[GH-29974](https://github.com/hashicorp/vault/pull/29974)]
* ui/activity: Replaces mount and namespace attribution charts with a table to allow sorting 
client count data by `namespace`, `mount_path`, `mount_type` or number of clients for 
a selected month. [[GH-30678](https://github.com/hashicorp/vault/pull/30678)]
* ui: Client count side nav link 'Vault Usage Metrics' renamed to 'Client Usage' [[GH-30765](https://github.com/hashicorp/vault/pull/30765)]
* ui: Client counting "running total" charts now reflect new clients only [[GH-30506](https://github.com/hashicorp/vault/pull/30506)]
* ui: Removed `FormError` component (not used) [[GH-34699](https://github.com/hashicorp/vault/pull/34699)]
* ui: Selecting a different method in the login form no longer updates the `/vault/auth?with=` query parameter [[GH-30500](https://github.com/hashicorp/vault/pull/30500)]
* ui: `/vault/auth?with=` query parameter now exclusively refers to the auth mount path and renders a simplified form [[GH-30500](https://github.com/hashicorp/vault/pull/30500)]

FEATURES:

* **Auto Irrevocable Lease Removal (Enterprise)**: Add the Vault Enterprise configuration param, `remove_irrevocable_lease_after`. When set to a non-zero value, this will automatically delete irrevocable leases after the configured duration exceeds the lease's expire time. The minimum duration allowed for this field is two days. [[GH-30703](https://github.com/hashicorp/vault/pull/30703)]
* **Development Cluster Configuration (Enterprise)**: Added `development_cluster` as a field to Vault's utilization reports.
The field is configurable via HCL and indicates whether the cluster is being used in a development environment, defaults to false if not set. [[GH-30659](https://github.com/hashicorp/vault/pull/30659)]
* **Entity-based and collective rate limit quotas (Enterprise)**: Add new `group_by` field to the rate limit quota API to support different grouping modes.
* **Login form customization (Enterprise)**: Adds support to choose a default and/or backup auth methods for the web UI login form to streamline the web UI login experience. [[GH-30700](https://github.com/hashicorp/vault/pull/30700)]
* **Plugin Downloads**: Support automatically downloading official HashiCorp secret and auth plugins from releases.hashicorp.com (beta)
* **SSH Key Signing Improvements (Enterprise)**: Add support for using managed keys to sign SSH keys in the SSH secrets engine.
* **Secret Recovery from Snapshot (Enterprise)**: Adds a framework to load an integrated storage 
snapshot into Vault and read, list, and recover KV v1 and cubbyhole secrets from the snapshot. [[GH-30739](https://github.com/hashicorp/vault/pull/30739)]
* **UI Secrets Engines**: TOTP secrets engine is now supported. [[GH-29751](https://github.com/hashicorp/vault/pull/29751)]
* **UI Telemetry**: Add Posthog for UI telemetry tracking on Vault Dedicated managed clusters [[GH-30425](https://github.com/hashicorp/vault/pull/30425)]
* **Vault Namespace Picker**: Updating the Vault Namespace Picker to enable search functionality, allow direct navigation to nested namespaces and improve accessibility. [[GH-30490](https://github.com/hashicorp/vault/pull/30490)]
* **Vault PKI SCEP Server (Enterprise)**: Support for the Simple Certificate Enrollment Protocol (SCEP) has been added to the Vault PKI Plugin. This allows standard SCEP clients to request certificates from a Vault server with no knowledge of Vault APIs.

IMPROVEMENTS:

* activity (enterprise): Added vault.client.billing_period.activity telemetry metric to emit information about the total number of distinct clients used in the current billing period.
* activity: mount_type was added to the API response of sys/internal/counters/activity [[GH-30071](https://github.com/hashicorp/vault/pull/30071)]
* activity: mount_type was added to the API response of sys/internal/counters/activity
* api (enterprise): Added a new API, `/sys/utilization-report`, giving a snapshot overview of Vault's utilization at a high level.
* api/client: Add Cert auth method support. This allows the client to authenticate using a client certificate. [[GH-29546](https://github.com/hashicorp/vault/pull/29546)]
* core (enterprise): Updated code and documentation to support FIPS 140-3 compliant algorithms.
* core (enterprise): allow a root token to relock a namespace locked by the Namespace API Lock feature.
* core (enterprise): report errors from the underlying seal when getting entropy.
* core (enterprise): update to FIPS 140-3 cryptographic module in the FIPS builds.
* core/metrics: added a new telemetry metric, `vault.core.response_status_code`, with two labels, `code`, and `type`, detailing the status codes of all responses to requests that Vault handles. [[GH-30354](https://github.com/hashicorp/vault/pull/30354)]
* core: Improve memory use of path management for namespaces, auth methods, and secrets engines. Now Vault should handle larger numbers of namespaces and multiple instances of the same secrets engine or auth method more efficiently. [[GH-31022](https://github.com/hashicorp/vault/pull/31022)]
* core: Updated code and documentation to support FIPS 140-3 compliant algorithms. [[GH-30576](https://github.com/hashicorp/vault/pull/30576)]
* core: support for X25519MLKEM768 (post quantum key agreement) in the Go TLS stack. [[GH-30603](https://github.com/hashicorp/vault/pull/30603)]
* events: Add `vault_index` to an event's metadata if the metadata contains `modified=true`, to support client consistency controls when reading from Vault in response to an event where storage was modified. [[GH-30725](https://github.com/hashicorp/vault/pull/30725)]
* physical/postgres: Adds support to authenticate with the PostgreSQL Backend server with cloud based identities (AWS IAM, Azure MSI and GCP IAM) [[GH-30681](https://github.com/hashicorp/vault/pull/30681)]
* plugins: Support registration of CE plugins with extracted artifact directory. [[GH-30673](https://github.com/hashicorp/vault/pull/30673)]
* secrets/aws: Add LIST endpoint to the AWS secrets engine static roles. [[GH-29842](https://github.com/hashicorp/vault/pull/29842)]
* secrets/pki: Add Delta (Freshest) CRL support to AIA information (both mount-level and issuer configured) [[GH-30319](https://github.com/hashicorp/vault/pull/30319)]
* secrets/transit (enterprise): enable the use of 192-bit keys for AES CMAC
* storage/mysql: Added support for getting mysql backend username and password from the environment variables `VAULT_MYSQL_USERNAME` and `VAULT_MYSQL_PASSWORD`. [[GH-30136](https://github.com/hashicorp/vault/pull/30136)]
* storage/raft: Upgrade hashicorp/raft library to v1.7.3 which includes additional logging on the leader when opening and sending a snapshot to a follower. [[GH-29976](https://github.com/hashicorp/vault/pull/29976)]
* transit: Exclude the partial wrapping key path from the transit/keys LIST operation. [[GH-30728](https://github.com/hashicorp/vault/pull/30728)]
* ui (enterprise): Replace date selector in client count usage page with fixed start and end dates that align with billing periods in order to return more relevant client counting data. [[GH-30349](https://github.com/hashicorp/vault/pull/30349)]
* ui/database: Adding input field for setting skip static role password rotation for database connection config, updating static role skip field to use toggle button [[GH-29820](https://github.com/hashicorp/vault/pull/29820)]
* ui/database: Adding password input field for creating a static role [[GH-30275](https://github.com/hashicorp/vault/pull/30275)]
* ui/database: Adding warning modal pop up when creating a static role that will be rotated immediately [[GH-30119](https://github.com/hashicorp/vault/pull/30119)]
* ui/database: Glimmerizing and adding validations to role create [[GH-29754](https://github.com/hashicorp/vault/pull/29754)]
* ui/database: Updating toggle buttons for skip_rotation_import to reverse polarity of values that get displayed versus whats sent to api [[GH-30055](https://github.com/hashicorp/vault/pull/30055)]
* ui: Add 'Refresh list' button to the namespace list page. [[GH-30692](https://github.com/hashicorp/vault/pull/30692)]
* ui: Enable search for a namespace on the namespace list page. [[GH-30680](https://github.com/hashicorp/vault/pull/30680)]
* ui: Hide "Other" tab when mounts are configured with `listing_visibility="unauth"`; all methods can be accessed via the "Sign in with other methods" link [[GH-30500](https://github.com/hashicorp/vault/pull/30500)]
* ui: Improve accessibility of login form to meet a11y standards [[GH-30500](https://github.com/hashicorp/vault/pull/30500)]
* ui: Replaces all instances of the deprecated event.keyCode with event.key [[GH-30493](https://github.com/hashicorp/vault/pull/30493)]
* ui: Update date selector in client count usage page to disable current month selection for Vault clusters without a license. [[GH-30488](https://github.com/hashicorp/vault/pull/30488)]
* ui: Use Hds::CodeBlock component to replace readonly JsonEditor instances [[GH-29720](https://github.com/hashicorp/vault/pull/29720)]
* ui: adds key value pair string inputs as optional form for wrap tool [[GH-29677](https://github.com/hashicorp/vault/pull/29677)]
* ui: remove ember-svg-jar dependency [[GH-30181](https://github.com/hashicorp/vault/pull/30181)]

DEPRECATIONS:

* api: Deprecated the `/sys/internal/counters/tokens` endpoint. Attempting to call this endpoint will return a 403 "unsupported path" exception. [[GH-30561](https://github.com/hashicorp/vault/pull/30561)]
* core: deprecate duplicate attributes in HCL configuration files and policy definitions [[GH-30386](https://github.com/hashicorp/vault/pull/30386)]

BUG FIXES:

* api/tokenhelper: Exec token_helper without a shell [[GH-29653](https://github.com/hashicorp/vault/pull/29653)]
* auth/aws: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* auth/ldap: Fix a bug that does not properly delete users and groups by first converting their names to lowercase when case senstivity option is off. [[GH-29922](https://github.com/hashicorp/vault/pull/29922)]
* auth/ldap: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* aws/secrets: Prevent vault from rejecting secret role configurations where no regions or endpoints are set [[GH-29996](https://github.com/hashicorp/vault/pull/29996)]
* core (enterprise): add nil check before attempting to use Rotation Manager operations.
* core (enterprise): fix a bug where plugin automated root rotations would stop after seal/unseal operations
* core (enterprise): fix issue with errors being swallowed on failed HSM logins. 
core/managed-keys (enterprise): fix RSA encryption/decryption with OAEP on managed keys.
* core: Fix a bug that prevents certain loggers from writing to a log file. [[GH-29917](https://github.com/hashicorp/vault/pull/29917)]
* core: Omit automatic version control information of the main module from compiled Vault binaries [[GH-30926](https://github.com/hashicorp/vault/pull/30926)]
* database: Prevent static roles created in versions prior to 1.15.0 from rotating on backend restart. [[GH-30320](https://github.com/hashicorp/vault/pull/30320)]
* database: no longer incorrectly add an "unrecognized parameters" warning for certain SQL database secrets config operations when another warning is returned [[GH-30327](https://github.com/hashicorp/vault/pull/30327)]
* identity: Fix non-deterministic merge behavior when two entities have
conflicting local aliases. [[GH-30390](https://github.com/hashicorp/vault/pull/30390)]
* identity: reintroduce RPC functionality for group creates, allowing performance standbys to handle external group changes during login and token renewal [[GH-30069](https://github.com/hashicorp/vault/pull/30069)]
* plugins (enterprise): Fix an issue where Enterprise plugins can't run on a standby node
when it becomes active because standby nodes don't extract the artifact when the plugin
is registered. Remove extracting from Vault and require the operator to place
the extracted artifact in the plugin directory before registration.
* plugins (enterprise): Fix plugin registration with artifact when a binary for the same plugin is already present in the plugin directory.
* plugins: plugin registration should honor the `plugin_tmpdir` config [[GH-29978](https://github.com/hashicorp/vault/pull/29978)]
* plugins: plugin registration should honor the `plugin_tmpdir` config
* raft/retry_join: Fix decoding `auto_join` configurations that include escape characters [[GH-29874](https://github.com/hashicorp/vault/pull/29874)]
* secrets/aws: fix a bug where environment and shared credential providers were overriding the WIF configuration [[GH-29982](https://github.com/hashicorp/vault/pull/29982)]
* secrets/aws: fix a case where GovCloud wasn't taken into account; fix a case where the region setting wasn't respected [[GH-30312](https://github.com/hashicorp/vault/pull/30312)]
* secrets/aws: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* secrets/database: Fix a bug where a global database plugin reload exits if any of the database connections are not available [[GH-29519](https://github.com/hashicorp/vault/pull/29519)]
* secrets/database: Treat all rotation_schedule values as UTC to ensure consistent behavior. [[GH-30606](https://github.com/hashicorp/vault/pull/30606)]
* secrets/db: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* secrets/openldap: Prevent static role rotation on upgrade when `NextVaultRotation` is nil.
Fixes an issue where static roles were unexpectedly rotated after upgrade due to a missing `NextVaultRotation` value. 
Now sets it to either `LastVaultRotation + RotationPeriod` or `now + RotationPeriod`. [[GH-30265](https://github.com/hashicorp/vault/pull/30265)]
* secrets/pki (enterprise): Address a parsing bug that rejected CMPv2 requests containing a validity field.
* secrets/pki: Fix a bug that prevents enabling automatic tidying of the CMPv2 nonce store. [[GH-29852](https://github.com/hashicorp/vault/pull/29852)]
* secrets/pki: fix a bug where key_usage was ignored when generating root certificates, and signing certain
intermediate certificates. [[GH-30034](https://github.com/hashicorp/vault/pull/30034)]
* secrets/transit (enterprise): ensure verify endpoint always returns valid field in batch_results with CMAC
* secrets/transit (enterprise): fixed encryption/decryption with RSA against PKCS#11 managed keys
* secrets/transit: ensure verify endpoint always returns valid field in batch_results with HMAC [[GH-30852](https://github.com/hashicorp/vault/pull/30852)]
* secrets/transit: fix a panic when rotating on a managed key returns an error [[GH-30214](https://github.com/hashicorp/vault/pull/30214)]
* ui/database: Added input field for setting 'skip_import_rotation' when creating a static role [[GH-29633](https://github.com/hashicorp/vault/pull/29633)]
* ui/kmip: Fixes KMIP credentials view and displays `private_key` after generating [[GH-30778](https://github.com/hashicorp/vault/pull/30778)]
* ui: Automatically refresh namespace list inside the namespace picker after creating or deleting a namespace in the UI. [[GH-30737](https://github.com/hashicorp/vault/pull/30737)]
* ui: Fix broken link to Hashicorp Vault developer site in the Web REPL help. [[GH-30670](https://github.com/hashicorp/vault/pull/30670)]
* ui: Fix initial setting of form toggle inputs for parameters nested within the `config` block [[GH-30960](https://github.com/hashicorp/vault/pull/30960)]
* ui: Fix refresh namespace list after deleting a namespace. [[GH-30680](https://github.com/hashicorp/vault/pull/30680)]
* ui: MFA methods now display the namespace path instead of the namespace id. [[GH-29588](https://github.com/hashicorp/vault/pull/29588)]
* ui: Redirect users authenticating with Vault as an OIDC provider to log in again when token expires. [[GH-30838](https://github.com/hashicorp/vault/pull/30838)]

## 1.19.13 Enterprise
### January 07, 2026

CHANGES:

* auth/oci: bump plugin to v0.18.1
* go: bump go version to 1.25.5
* packaging: Container images are now exported using a compressed OCI image layout.
* packaging: UBI container images are now built on the UBI 10 minimal image.
* secrets/azure: Update plugin to [v0.21.5](https://github.com/hashicorp/vault-plugin-secrets-azure/releases/tag/v0.21.5). Improves retry handling during Azure application and service principal creation to reduce transient failures.
* storage: Upgrade aerospike client library to v8.

IMPROVEMENTS:

* core: check rotation manager queue every 5 seconds instead of 10 seconds to improve responsiveness.
* go: update to golang/x/crypto to v0.45.0 to resolve GHSA-f6x5-jh6r-wrfv, GHSA-j5w8-q4qc-rx2x, GO-2025-4134 and GO-2025-4135.
* rotation: Ensure rotations for shared paths only execute on the Primary cluster's active node. Ensure rotations for local paths execute on the cluster-local active node.
* sdk/rotation: Prevent rotation attempts on read-only storage.
* secrets-sync (enterprise): Added support for a boolean force_delete flag (default: false). 
When set to true, this flag allows deletion of a destination even if its associations cannot be unsynced. 
This option should be used only as a last-resort deletion mechanism, as any secrets already synced to the external provider will remain orphaned and require manual cleanup.

BUG FIXES:

* auth/approle (enterprise): Fixed bug that prevented periodic tidy running on performance secondary.
* http: skip JSON limit parsing on cluster listener.
* quotas: Vault now protects plugins with ResolveRole operations from panicking on quota creation.
* replication (enterprise): fix rare panic due to race when enabling a secondary with Consul storage.
* rotation: Fix a bug where a performance secondary would panic if a write was made to a local mount.
* secret-sync (enterprise): Improved unsync error handling by treating cases where the destination no longer exists as successful.
* secrets-sync (enterprise): Corrected a bug where the deletion of the latest KV-V2 secret version caused the associated external secret to be deleted entirely. The sync job now implements a version fallback mechanism to find and sync the highest available active version, ensuring continuity and preventing the unintended deletion of the external secret resource.
* ui/pki: Fix handling of values that contain commas in list fields like `crl_distribution_points`.


## 1.19.12 Enterprise
### November 19, 2025

CHANGES:

* core: Bump Go version to 1.24.10
* policies: add VAULT_NEW_PER_ELEMENT_MATCHING_ON_LIST env var to adopt new "contains all" list matching behavior on
allowed_parameters and denied_parameters

IMPROVEMENTS:

* Update github.com/dvsekhvalnov/jose2go to fix security vulnerability CVE-2025-63811.
* auth/ldap: Require non-empty passwords on login command to prevent unauthenticated access to Vault.
* ui/pki: Adds support to configure `server_flag`, `client_flag`, `code_signing_flag`, and `email_protection_flag` parameters for creating/updating a role.

BUG FIXES:

* core/rotation: avoid shifting timezones by ignoring cron.SpecSchedule
* ui/pki: Fixes certificate parsing of the `key_usage` extension so details accurately reflect certificate values.
* ui/pki: Fixes creating and updating a role so `basic_constraints_valid_for_non_ca` is correctly set.
* ui: Resolved a regression that prevented users with create and update permissions on KV v1 secrets from opening the edit view. The UI now correctly recognizes these capabilities and allows editing without requiring full read access.
* ui: remove unnecessary 'credential type' form input when generating AWS secrets

## 1.19.11 Enterprise
### October 22, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* auth/aws: fix an issue where a user may be able to bypass authentication to Vault due to incorrect caching of the AWS client
* ui: disable scarf analytics for ui builds

CHANGES:

* auth/alicloud: Update plugin to [v0.20.1](https://github.com/hashicorp/vault-plugin-auth-alicloud/releases/tag/v0.20.1)
* core: Bump Go version to 1.24.9.
* http: Evaluate rate limit quotas before checking JSON limits during request handling.

IMPROVEMENTS:

* secrets/database: Add root rotation support for Snowflake database secrets engines using key-pair credentials.

BUG FIXES:

* auth: fixed panic when suppling integer as a lease_id in renewal.
* core (enterprise): Avoid duplicate seal rewrapping, and ensure that cluster secondaries rewrap after a seal migration.
* core: interpret all new rotation manager rotation_schedules as UTC to avoid inadvertent use of tz-local
* core: resultant-acl now merges segment-wildcard (`+`) paths with existing prefix rules in `glob_paths`, so clients receive a complete view of glob-style permissions. This unblocks UI sidebar navigation checks and namespace access banners.
* secrets/database: respect the escaping/disable_escaping state when using self-managed static roles
* ui: Fixes permissions for hiding and showing sidebar navigation items for policies that include special characters: `+`, `*`

## 1.19.10 Enterprise
### September 24, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* core: Update github.com/ulikunitz/xz to fix security vulnerability GHSA-25xm-hr59-7c27.

CHANGES:

* core: Updates post-install script to print updated license information
* database/snowflake: Update plugin to [v0.13.3](https://github.com/hashicorp/vault-plugin-database-snowflake/releases/tag/v0.13.3)

IMPROVEMENTS:

* Raft: Auto-join will now allow you to enforce IPv4 on networks that allow IPv6 and dual-stack enablement, which is on by default in certain regions.
* auth/cert: Support RFC 9440 colon-wrapped Base64 certificates in `x_forwarded_for_client_cert_header`, to fix TLS certificate auth errors with Google Cloud Application Load Balancer.

BUG FIXES:

* auth/cert: Recover from partially populated caches of trusted certificates if one or more certificates fails to load.
* core: Fixed issue where under certain circumstances the rotation manager would spawn goroutines indefinitely.
* core: Role based quotas now work for cert auth
* secrets/transit: Fix error when using ed25519 keys that were imported with derivation enabled
* sys/mounts: enable unsetting allowed_response_headers

## 1.19.9 Enterprise
### August 28, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

FEATURES:

* **IBM RACF Static Role Password Phrase Management**: Add support for static role password phrase management to the LDAP secrets engine. (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/184)

SECURITY:

* core: Update github.com/hashicorp/go-getter to fix security vulnerability GHSA-wjrx-6529-hcj3.

CHANGES:

* core: Bump Go version to 1.24.6.
* http: Add JSON configurable limits to HTTP handling for JSON payloads: `max_json_depth`, `max_json_string_value_length`, `max_json_object_entry_count`, `max_json_array_element_count`.
* sdk: Upgrade to go-secure-stdlib/plugincontainer@v0.4.2, which also bumps github.com/docker/docker to v28.3.3+incompatible
* secrets/openldap: update plugin to v0.15.5

IMPROVEMENTS:

* auth/ldap: add explicit logging to rotations in ldap
* core (enterprise): improve rotation manager logging to include specific lines for rotation success and failure
* secrets/database: log password rotation success (info) and failure (error). Some relevant log lines have been updated to include "path" fields.
* secrets/transit: add logging on both success and failure of key rotation
* ui: Use the Helios Design System Code Block component for all readonly code editors and use its Code Editor component for all other code editors

BUG FIXES:

* core (enterprise): fix a bug where issuing a token in a namespace used root auth configuration instead of namespace auth configuration
* core/metrics: Add service name prefix for core HA metrics to avoid duplicate, zero-value metrics.
* core/seal: When Seal-HA is enabled, make it an error to persist the barrier
keyring when not all seals are healthy.  This prevents the possibility of
failing to unseal when a different subset of seals are healthy than were
healthy at last write.
* raft (enterprise): auto-join will now work in regions that do not support dual-stack
* raft/autopilot: Fixes an issue with enterprise redundancy zones where, if the leader was in a redundancy zone and that leader becomes unavailable, the node would become an unzoned voter. This can artificially inflate the required number of nodes for quorum, leading to a situation where the cluster cannot recover if another leader subsequently becomes unavailable. Vault will now keep an unavailable node in its last known redundancy zone as a non-voter.
* replication (enterprise): Fix bug where group updates fail when processed on a
standby node in a PR secondary cluster.
* secrets-sync (enterprise): GCP locational KMS keys are no longer incorrectly removed when the location name is all lowercase.
* secrets/database/postgresql: Support for multiline statements in the `rotation_statements` field.

## 1.19.8 Enterprise
### August 06, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* auth/ldap: fix MFA/TOTP enforcement bypass when username_as_alias is enabled [[GH-31427](https://github.com/hashicorp/vault/pull/31427),[HCSEC-2025-20](https://discuss.hashicorp.com/t/hcsec-2025-20-vault-ldap-mfa-enforcement-bypass-when-using-username-as-alias/76092)].

BUG FIXES:

* agent/template: Fixed issue where templates would not render correctly if namespaces was provided by config, and the namespace and mount path of the secret were the same. [[GH-31392](https://github.com/hashicorp/vault/pull/31392)]
* identity/mfa: revert cache entry change from #31217 and document cache entry values [[GH-31421](https://github.com/hashicorp/vault/pull/31421)]

## 1.19.7 Enterprise
### July 25, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY: 

* audit: **breaking change** privileged vault operator may execute code on the underlying host (CVE-2025-6000). Vault will not unseal if the only configured file audit device has executable permissions (e.g., 0777, 0755). See recent [breaking change](https://developer.hashicorp.com/vault/docs/updates/important-changes#breaking-changes) docs for more details. [[GH-31211](https://github.com/hashicorp/vault/pull/31211),[HCSEC-2025-14](https://discuss.hashicorp.com/t/hcsec-2025-14-privileged-vault-operator-may-execute-code-on-the-underlying-host/76033)]
* auth/userpass: timing side-channel in vault's userpass auth method (CVE-2025-6011)[HCSEC-2025-15](https://discuss.hashicorp.com/t/hcsec-2025-15-timing-side-channel-in-vault-s-userpass-auth-method/76034)
* core/login: vault userpass and ldap user lockout bypass (CVE-2025-6004). update alias lookahead to respect username case for LDAP and username/password. [[GH-31352](https://github.com/hashicorp/vault/pull/31352),[HCSEC-2025-16](https://discuss.hashicorp.com/t/hcsec-2025-16-vault-userpass-and-ldap-user-lockout-bypass/76035)]
* secrets/totp: vault totp secrets engine code reuse (CVE-2025-6014) [[GH-31246](https://github.com/hashicorp/vault/pull/31246),[HCSEC-2025-17](https://discuss.hashicorp.com/t/hcsec-2025-17-vault-totp-secrets-engine-code-reuse/76036)]
* auth/cert: vault certificate auth method did not validate common name for non-ca certificates (CVE-2025-6037). test non-CA cert equality on login matching instead of individual fields. [[GH-31210](https://github.com/hashicorp/vault/pull/31210),[HCSEC-2025-18](https://discuss.hashicorp.com/t/hcsec-2025-18-vault-certificate-auth-method-did-not-validate-common-name-for-non-ca-certificates/76037)]
* core/mfa: vault login mfa bypass of rate limiting and totp token reuse (CVE-2025-6015) [[GH-31217](https://github.com/hashicorp/vault/pull/31297),[HCSEC-2025-19](https://discuss.hashicorp.com/t/hcsec-2025-19-vault-login-mfa-bypass-of-rate-limiting-and-totp-token-reuse/76038)]

BUG FIXES:

* auto-reporting (enterprise): Clarify debug logs to accurately reflect when automated license utilization reporting is enabled or disabled, especially since manual reporting is always initialized.
* core/seal (enterprise): Fix a bug that caused the seal rewrap process to abort in the presence of partially sealed entries.
* kmip (enterprise): Fix a panic that can happen when a KMIP client makes a request before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* plugins: Fix panics that can occur when a plugin audits a request or response before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* product usage reporting (enterprise): Clarify debug logs to accurately reflect when anonymous product usage reporting is enabled or disabled, especially since manual reporting is always initialized.
* replication (enterprise): Fix bug with mount invalidations consuming excessive memory.
* secrets-sync (enterprise): Unsyncing secret-key granularity associations will no longer give a misleading error about a failed unsync operation that did indeed succeed.
* secrets/gcp: Update to vault-plugin-secrets-gcp@v0.21.4 to address more eventual consistency issues

## 1.19.6 Enterprise
### June 25, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* core: require a nonce when cancelling a rekey operation that was initiated within the last 10 minutes. [[GH-30794](https://github.com/hashicorp/vault/pull/30794)],[[HCSEC-2025-11](https://discuss.hashicorp.com/t/hcsec-2025-11-vault-vulnerable-to-recovery-key-cancellation-denial-of-service/75570)]
* core/identity: vault root namespace operator may elevate privileges (CVE-2025-5999). Fix string contains check in Identity APIs to be case-insensitive. [[GH-31045](https://github.com/hashicorp/vault/pull/31045),[HCSEC-2025-13](https://discuss.hashicorp.com/t/hcsec-2025-13-vault-root-namespace-operator-may-elevate-token-privileges/76032)]

CHANGES:

* api: Update the default API client to check for the `Retry-After` header and, if it exists, wait for the specified duration before retrying the request. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* auth/azure: Update plugin to v0.20.5
* core: Bump Go version to 1.24.4.
* quotas/rate-limit: Round up the `Retry-After` value to the nearest second when calculating the retry delay. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* secrets/azure: Update plugin to v0.21.4 [[GH-30833](https://github.com/hashicorp/vault/pull/30833)]
* secrets/database: Update vault-plugin-database-snowflake to v0.13.2 [[GH-30867](https://github.com/hashicorp/vault/pull/30867)]

IMPROVEMENTS:

* core: Improve memory use of path management for namespaces, auth methods, and secrets engines. Now Vault should handle larger numbers of namespaces and multiple instances of the same secrets engine or auth method more efficiently. [[GH-31022](https://github.com/hashicorp/vault/pull/31022)]

DEPRECATIONS:

* core: deprecate duplicate attributes in HCL configuration files and policy definitions [[GH-30386](https://github.com/hashicorp/vault/pull/30386)]

BUG FIXES:

* core: Omit automatic version control information of the main module from compiled Vault binaries [[GH-30926](https://github.com/hashicorp/vault/pull/30926)]
* secrets/database: Treat all rotation_schedule values as UTC to ensure consistent behavior. [[GH-30606](https://github.com/hashicorp/vault/pull/30606)]
* secrets/transit (enterprise): ensure verify endpoint always returns valid field in batch_results with CMAC
* secrets/transit: ensure verify endpoint always returns valid field in batch_results with HMAC [[GH-30852](https://github.com/hashicorp/vault/pull/30852)]
* ui/kmip: Fixes KMIP credentials view and displays `private_key` after generating [[GH-30778](https://github.com/hashicorp/vault/pull/30778)]
* ui: Redirect users authenticating with Vault as an OIDC provider to log in again when token expires. [[GH-30838](https://github.com/hashicorp/vault/pull/30838)]

## 1.19.5
### May 30, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* database/snowflake: Update plugin to v0.13.1 [[GH-30775](https://github.com/hashicorp/vault/pull/30775)]

IMPROVEMENTS:

* plugins: Support registration of CE plugins with extracted artifact directory. [[GH-30673](https://github.com/hashicorp/vault/pull/30673)]

BUG FIXES:

* ui: Fix broken link to Hashicorp Vault developer site in the Web REPL help. [[GH-30670](https://github.com/hashicorp/vault/pull/30670)]

## 1.19.4
### May 16, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* Update vault-plugin-auth-cf to v0.20.1 [[GH-30586](https://github.com/hashicorp/vault/pull/30586)]
* auth/azure: Update plugin to v0.20.4 [[GH-30543](https://github.com/hashicorp/vault/pull/30543)]
* core: Bump Go version to 1.24.3.

IMPROVEMENTS:

* Namespaces (enterprise): allow a root token to relock a namespace
* core (enterprise): update to FIPS 140-3 cryptographic module in the FIPS builds.
* core: Updated code and documentation to support FIPS 140-3 compliant algorithms. [[GH-30576](https://github.com/hashicorp/vault/pull/30576)]
* core: support for X25519MLKEM768 (post quantum key agreement) in the Go TLS stack. [[GH-30603](https://github.com/hashicorp/vault/pull/30603)]
* ui: Replaces all instances of the deprecated event.keyCode with event.key [[GH-30493](https://github.com/hashicorp/vault/pull/30493)]

BUG FIXES:

* core (enterprise): fix a bug where plugin automated root rotations would stop after seal/unseal operations
* plugins (enterprise): Fix an issue where Enterprise plugins can't run on a standby node
when it becomes active because standby nodes don't extract the artifact when the plugin
is registered. Remove extracting from Vault and require the operator to place
the extracted artifact in the plugin directory before registration.

## 1.19.3
### April 30, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* core: vault may expose sensitive information in error logs when processing malformed data with the kv v2 plugin[[GH-30388](https://github.com/hashicorp/vault/pull/30388), [HCSEC-2025-09](https://discuss.hashicorp.com/t/hcsec-2025-09-vault-may-expose-sensitive-information-in-error-logs-when-processing-malformed-data-with-the-kv-v2-plugin/74717)]

CHANGES:

* auth/jwt: Update plugin to v0.23.2 [[GH-30434](https://github.com/hashicorp/vault/pull/30434)]

BUG FIXES:

* core (enterprise): fix issue with errors being swallowed on failed HSM logins.
* database: Prevent static roles created in versions prior to 1.15.0 from rotating on backend restart. [[GH-30320](https://github.com/hashicorp/vault/pull/30320)]
* database: no longer incorrectly add an "unrecognized parameters" warning for certain SQL database secrets config operations when another warning is returned [[GH-30327](https://github.com/hashicorp/vault/pull/30327)]
* identity: Fix non-deterministic merge behavior when two entities have conflicting local aliases. [[GH-30390](https://github.com/hashicorp/vault/pull/30390)]
* plugins: plugin registration should honor the `plugin_tmpdir` config [[GH-29978](https://github.com/hashicorp/vault/pull/29978)]
* secrets/aws: fix a case where GovCloud wasn't taken into account; fix a case where the region setting wasn't respected [[GH-30312](https://github.com/hashicorp/vault/pull/30312)]

## 1.19.2
### April 18, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* core: Bump Go version to 1.23.7
* core: Bump Go version to 1.23.8
* secrets/openldap: Update plugin to v0.15.4 [[GH-30279](https://github.com/hashicorp/vault/pull/30279)]

BUG FIXES:

* secrets/openldap: Prevent static role rotation on upgrade when `NextVaultRotation` is nil. Fixes an issue where static roles were unexpectedly rotated after upgrade due to a missing `NextVaultRotation` value. Now sets it to either `LastVaultRotation + RotationPeriod` or `now + RotationPeriod`. [[GH-30265](https://github.com/hashicorp/vault/pull/30265)]
* secrets/pki (enterprise): Address a parsing bug that rejected CMPv2 requests containing a validity field.
* secrets/pki: fix a bug where key_usage was ignored when generating root certificates, and signing certain intermediate certificates. [[GH-30034](https://github.com/hashicorp/vault/pull/30034)]
* secrets/transit: fix a panic when rotating on a managed key returns an error [[GH-30214](https://github.com/hashicorp/vault/pull/30214)]

## 1.19.1
### April 4, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* auth/azure: Update plugin to v0.20.2. Login requires `resource_group_name`, `vm_name`, and `vmss_name` to match token claims [[HCSEC-2025-07](https://discuss.hashicorp.com/t/hcsec-2025-07-vault-s-azure-authentication-method-bound-location-restriction-could-be-bypassed-on-login/74716) [GH-30052](https://github.com/hashicorp/vault/pull/30052)].

CHANGES:

* UI: remove outdated and unneeded js string extensions [[GH-29834](https://github.com/hashicorp/vault/pull/29834)]
* auth/azure: Update plugin to v0.20.3 [[GH-30082](https://github.com/hashicorp/vault/pull/30082)]
* auth/gcp: Update plugin to v0.20.2 [[GH-30081](https://github.com/hashicorp/vault/pull/30081)]
* core: Verify that the client IP address extracted from an X-Forwarded-For header is a valid IPv4 or IPv6 address [[GH-29774](https://github.com/hashicorp/vault/pull/29774)]
* secrets/azure: Update plugin to v0.21.2 [[GH-30037](https://github.com/hashicorp/vault/pull/30037)]
* secrets/azure: Update plugin to v0.21.3 [[GH-30083](https://github.com/hashicorp/vault/pull/30083)]
* secrets/gcp: Update plugin to v0.21.2 [[GH-29970](https://github.com/hashicorp/vault/pull/29970)]
* secrets/gcp: Update plugin to v0.21.3 [[GH-30080](https://github.com/hashicorp/vault/pull/30080)]
* secrets/openldap: Update plugin to v0.15.2 [[GH-30079](https://github.com/hashicorp/vault/pull/30079)]

IMPROVEMENTS:

* activity: mount_type was added to the API response of sys/internal/counters/activity [[GH-30071](https://github.com/hashicorp/vault/pull/30071)]
* core (enterprise): report errors from the underlying seal when getting entropy.
* storage/raft: Upgrade hashicorp/raft library to v1.7.3 which includes additional logging on the leader when opening and sending a snapshot to a follower. [[GH-29976](https://github.com/hashicorp/vault/pull/29976)]

BUG FIXES:

* auth/aws: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* auth/ldap: Fix a bug that does not properly delete users and groups by first converting their names to lowercase when case senstivity option is off. [[GH-29922](https://github.com/hashicorp/vault/pull/29922)]
* auth/ldap: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* aws/secrets: Prevent vault from rejecting secret role configurations where no regions or endpoints are set [[GH-29996](https://github.com/hashicorp/vault/pull/29996)]
* core (enterprise): add nil check before attempting to use Rotation Manager operations.
* core: Fix a bug that prevents certain loggers from writing to a log file. [[GH-29917](https://github.com/hashicorp/vault/pull/29917)]
* core/raft: Fix decoding `auto_join` configurations that include escape characters. [[GH-29874](https://github.com/hashicorp/vault/pull/29874)]
* identity: reintroduce RPC functionality for group creates, allowing performance standbys to handle external group changes during login and token renewal [[GH-30069](https://github.com/hashicorp/vault/pull/30069)]
* plugins (enterprise): Fix plugin registration with artifact when a binary for the same plugin is already present in the plugin directory.
* secrets/aws: fix a bug where environment and shared credential providers were overriding the WIF configuration [[GH-29982](https://github.com/hashicorp/vault/pull/29982)]
* secrets/aws: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* secrets/db: fix a panic when a performance standby node attempts to write/update config. [[GH-30039](https://github.com/hashicorp/vault/pull/30039)]
* secrets/pki: Fix a bug that prevents enabling automatic tidying of the CMPv2 nonce store. [[GH-29852](https://github.com/hashicorp/vault/pull/29852)]

## 1.19.0
### March 5, 2025

**Enterprise LTS:** Vault Enterprise 1.19 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20241115202008-166203013d8e
* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.2.0

CHANGES:

* agent/config: Configuration values including IPv6 addresses will be automatically translated and displayed conformant to RFC-5952 §4. [[GH-29517](https://github.com/hashicorp/vault/pull/29517)]
* api: Add to sys/health whether the node has been removed from the HA cluster. If the node has been removed, return code 530 by default or the value of the `removedcode` query parameter. [[GH-28991](https://github.com/hashicorp/vault/pull/28991)]
* api: Add to sys/health whether the standby node has been able to successfully send heartbeats to the active node and the time in milliseconds since the last heartbeat. If the standby has been unable to send a heartbeat, return code 474 by default or the value of the `haunhealthycode` query parameter. [[GH-28991](https://github.com/hashicorp/vault/pull/28991)]
* auth/alicloud: Update plugin to v0.20.0 [[GH-29613](https://github.com/hashicorp/vault/pull/29613)]
* auth/azure: Update plugin to v0.19.1 [[GH-28712](https://github.com/hashicorp/vault/pull/28712)]
* auth/azure: Update plugin to v0.19.2 [[GH-28848](https://github.com/hashicorp/vault/pull/28848)]
* auth/azure: Update plugin to v0.20.0 [[GH-29606](https://github.com/hashicorp/vault/pull/29606)]
* auth/azure: Update plugin to v0.20.1 [[GH-29728](https://github.com/hashicorp/vault/pull/29728)]
* auth/cf: Update plugin to v0.19.1 [[GH-29295](https://github.com/hashicorp/vault/pull/29295)]
* auth/cf: Update plugin to v0.20.0 [[GH-29528](https://github.com/hashicorp/vault/pull/29528)]
* auth/gcp: Update plugin to v0.20.0 [[GH-29591](https://github.com/hashicorp/vault/pull/29591)]
* auth/gcp: Update plugin to v0.20.1 [[GH-29736](https://github.com/hashicorp/vault/pull/29736)]
* auth/jwt: Update plugin to v0.23.0 [[GH-29553](https://github.com/hashicorp/vault/pull/29553)]
* auth/kerberos: Update plugin to v0.14.0 [[GH-29617](https://github.com/hashicorp/vault/pull/29617)]
* auth/kubernetes: Update plugin to v0.21.0 [[GH-29619](https://github.com/hashicorp/vault/pull/29619)]
* auth/ldap: An error will now be returned on login if the number of entries returned from the user DN LDAP search is more than one. [[GH-29302](https://github.com/hashicorp/vault/pull/29302)]
* auth/ldap: No longer return authentication warnings to client. [[GH-29134](https://github.com/hashicorp/vault/pull/29134)]
* auth/oci: Update plugin to v0.18.0 [[GH-29620](https://github.com/hashicorp/vault/pull/29620)]
* core (enterprise): Add tracking of performance standbys by their HA node ID so that RPC connections can be more easily cleaned up when nodes are removed. [[GH-29303](https://github.com/hashicorp/vault/pull/29303)]
* core/ha (enterprise): Failed attempts to become a performance standby node are now using an exponential backoff instead of a
10 second delay in between retries. The backoff starts at 2s and increases by a factor of two until reaching
the maximum of 16s. This should make unsealing of the node faster in some cases.
* core/raft: Return an error on sys/storage/raft/join if a node that has been removed from raft cluster attempts to re-join when it still has existing raft data on disk. [[GH-29090](https://github.com/hashicorp/vault/pull/29090)]
* core: Bump Go version to 1.23.6.
* database/couchbase: Update plugin to v0.13.0 [[GH-29543](https://github.com/hashicorp/vault/pull/29543)]
* database/elasticsearch: Update plugin to v0.17.0 [[GH-29542](https://github.com/hashicorp/vault/pull/29542)]
* database/mongodbatlas: Update plugin to v0.14.0 [[GH-29584](https://github.com/hashicorp/vault/pull/29584)]
* database/redis-elasticache: Update plugin to v0.6.0 [[GH-29594](https://github.com/hashicorp/vault/pull/29594)]
* database/redis: Update plugin to v0.5.0 [[GH-29597](https://github.com/hashicorp/vault/pull/29597)]
* database/snowflake: Update plugin to v0.13.0 [[GH-29554](https://github.com/hashicorp/vault/pull/29554)]
* kmip (enterprise): RSA key generation now enforces key sizes of 2048 or higher
* login (enterprise): Return a 500 error during logins when performance standby nodes make failed gRPC requests to the active node. [[GH-28807](https://github.com/hashicorp/vault/pull/28807)]
* proxy/config: Configuration values including IPv6 addresses will be automatically translated and displayed conformant to RFC-5952 §4. [[GH-29517](https://github.com/hashicorp/vault/pull/29517)]
* raft/autopilot (enterprise): Alongside the CE autopilot update, update raft-autopilot-enterprise library to v0.3.0 and add enterprise-specific regression testing.
* sdk: Upgrade to go-secure-stdlib/plugincontainer@v0.4.1, which also bumps github.com/docker/docker to v27.2.1+incompatible [[GH-28456](https://github.com/hashicorp/vault/pull/28456)]
* secrets/ad: Update plugin to v0.20.1 [[GH-29648](https://github.com/hashicorp/vault/pull/29648)]
* secrets/alicloud: Update plugin to v0.19.0 [[GH-29512](https://github.com/hashicorp/vault/pull/29512)]
* secrets/aws: The AWS Secrets engine now persists entries to storage between writes. This enables users
to not have to pass every required field on each write and to make individual updates as necessary.
Note: in order to zero out a value that is previously configured, users must now explicitly set the
field to its zero value on an update. [[GH-29497](https://github.com/hashicorp/vault/pull/29497)]
* secrets/azure: Update plugin to v0.20.1 [[GH-28699](https://github.com/hashicorp/vault/pull/28699)]
* secrets/azure: Update plugin to v0.21.0 [[GH-29639](https://github.com/hashicorp/vault/pull/29639)]
* secrets/azure: Update plugin to v0.21.1 [[GH-29729](https://github.com/hashicorp/vault/pull/29729)]
* secrets/gcp: Update plugin to v0.21.0 [[GH-29598](https://github.com/hashicorp/vault/pull/29598)]
* secrets/gcp: Update plugin to v0.21.1 [[GH-29747](https://github.com/hashicorp/vault/pull/29747)]
* secrets/gcpkms: Update plugin to v0.20.0 [[GH-29612](https://github.com/hashicorp/vault/pull/29612)]
* secrets/kubernetes: Update plugin to v0.10.0 [[GH-29592](https://github.com/hashicorp/vault/pull/29592)]
* secrets/kv: Update plugin to v0.21.0 [[GH-29614](https://github.com/hashicorp/vault/pull/29614)]
* secrets/mongodbatlas: Update plugin to v0.14.0 [[GH-29583](https://github.com/hashicorp/vault/pull/29583)]
* secrets/openldap: Update plugin to v0.14.1 [[GH-28479](https://github.com/hashicorp/vault/pull/28479)]
* secrets/openldap: Update plugin to v0.14.2 [[GH-28704](https://github.com/hashicorp/vault/pull/28704)]
* secrets/openldap: Update plugin to v0.14.3 [[GH-28780](https://github.com/hashicorp/vault/pull/28780)]
* secrets/openldap: Update plugin to v0.14.5 [[GH-29551](https://github.com/hashicorp/vault/pull/29551)]
* secrets/openldap: Update plugin to v0.15.0 [[GH-29605](https://github.com/hashicorp/vault/pull/29605)]
* secrets/openldap: Update plugin to v0.15.1 [[GH-29727](https://github.com/hashicorp/vault/pull/29727)]
* secrets/pki: Enforce the issuer constraint extensions (extended key usage, name constraints, issuer name) when issuing or signing leaf certificates. For more information see [PKI considerations](https://developer.hashicorp.com/vault/docs/secrets/pki/considerations#issuer-constraints-enforcement) [[GH-29045](https://github.com/hashicorp/vault/pull/29045)]
* secrets/terraform: Update plugin to v0.11.0 [[GH-29541](https://github.com/hashicorp/vault/pull/29541)]
* server/config: Configuration values including IPv6 addresses will be automatically translated and displayed conformant to RFC-5952 §4. [[GH-29228](https://github.com/hashicorp/vault/pull/29228)]
* storage/raft: Do not allow nodes that have been removed from the raft cluster configuration to respond to requests. Shutdown and seal raft nodes when they are removed. [[GH-28875](https://github.com/hashicorp/vault/pull/28875)]
* ui: Partially reverts #20431 and removes ability to download unencrypted kv v2 secret data [[GH-29290](https://github.com/hashicorp/vault/pull/29290)]
* ui: Upgrade Ember data to v5.3.2 (and minor upgrade of ember-cli, ember-source to v5.8.0) [[GH-28798](https://github.com/hashicorp/vault/pull/28798)]

FEATURES:

* **AWS Secrets Cross-Account Management Support** (enterprise): Add support for cross-account management of static roles in AWS secrets engine.
* **Automated Root Rotation**: A schedule or ttl can be defined for automated rotation of the root credential. [[GH-29535](https://github.com/hashicorp/vault/pull/29535)]
* **Automated Root Rotation**: Adds Automated Root Rotation capabilities to the AWS Auth and AWS Secrets
plugins. This allows plugin users to automate their root credential rotations based on configurable
schedules/periods via the Rotation Manager. Note: Enterprise only. [[GH-29497](https://github.com/hashicorp/vault/pull/29497)]
* **Automated Root Rotation**: Adds Automated Root Rotation capabilities to the DB Secrets plugin.
This allows plugin users to automate their root credential rotations based on configurable
schedules/periods via the Rotation Manager. Note: Enterprise only. [[GH-29557](https://github.com/hashicorp/vault/pull/29557)]
* **Automated Root Rotation**: Adds Automated Root Rotation capabilities to the GCP Auth plugin.
This allows plugin users to automate their root credential rotations based on configurable
schedules/periods via the Rotation Manager. Note: Enterprise only. [[GH-29591](https://github.com/hashicorp/vault/pull/29591)]
* **Automated Root Rotation**: Adds Automated Root Rotation capabilities to the GCP Secrets plugin.
This allows plugin users to automate their root credential rotations based on configurable
schedules/periods via the Rotation Manager. Note: Enterprise only. [[GH-29598](https://github.com/hashicorp/vault/pull/29598)]
* **Identity De-duplication**: Vault can now automatically resolve duplicate
Entities and Groups by renaming them. This feature is disabled by default and
can be enabled through the `force_identity_deduplication` activation flag. [[GH-29356](https://github.com/hashicorp/vault/pull/29356)]
* **Plugins**: Allow Enterprise plugins to run externally on Vault Enterprise only.
* **Post-Quantum Cryptography Support**: Experimental support for PQC signatures with ML-DSA in Transit.
* **Product Usage Reporting**: Added product usage reporting, which collects anonymous, numerical, non-sensitive data about Vault feature usage, and adds it to the existing utilization reports. [[GH-28858](https://github.com/hashicorp/vault/pull/28858)]
* **Rotation Manager**: Add Rotation Manager to Vault Enterprise Core. The Rotation Manager enables
plugin users to automate their root credential rotations based on configurable schedules/periods.
* **Skip auto import rotation of static roles (enterprise)**: The Database secrets engine now allows skipping the automatic rotation of static roles during import.
* **Transit Ed25519ph and Ed25519ctx support (Enterprise)**: Support for signing and verifying Ed25519ph and Ed25519ctx signatures types.

IMPROVEMENTS:

* CLI: adds an optional flag (--fail-if-not-fulfilled) to the renew command, which lets the renew command fail on unfulfillable requests and allows command chaining to allow further executions. [[GH-29060](https://github.com/hashicorp/vault/pull/29060)]
* audit: Audit logs will contain User-Agent headers when they are present in the incoming request. They are not
HMAC'ed by default but can be configured to be via the `/sys/config/auditing/request-headers/user-agent` endpoint. [[GH-28596](https://github.com/hashicorp/vault/pull/28596)]
* auth/approle: seal wrap approle secrets if seal wrap is enabled. [[GH-28703](https://github.com/hashicorp/vault/pull/28703)]
* auth/cert: Add new configuration option `enable_metadata_on_failures` to add client cert metadata on login failures to audit log and response [[GH-29044](https://github.com/hashicorp/vault/pull/29044)]
* auth/ldap: Adds an option to enable sAMAccountname logins when upndomain is set. [[GH-29118](https://github.com/hashicorp/vault/pull/29118)]
* auth/okta: update to okta sdk v5 from v2. Transitively updates go-jose dependency to >=3.0.3 to resolve GO-2024-2631. See https://github.com/okta/okta-sdk-golang/blob/master/MIGRATING.md for details on changes. [[GH-28121](https://github.com/hashicorp/vault/pull/28121)]
* auto-auth/cert: support watching changes on certificate/key files and notifying the auth handler when `enable_reauth_on_new_credentials` is enabled. [[GH-28126](https://github.com/hashicorp/vault/pull/28126)]
* auto-auth: support new config option `enable_reauth_on_new_credentials`, supporting re-authentication when receiving new credential on certain auto-auth types [[GH-28126](https://github.com/hashicorp/vault/pull/28126)]
* command/server: Add support for dumping pprof files during startup using CLI option `pprof-dump-dir` [[GH-27033](https://github.com/hashicorp/vault/pull/27033)]
* core/identity: Improve performance of loading entities when unsealing by batching updates, caching local alias storage reads, and doing more work in parallel. [[GH-29326](https://github.com/hashicorp/vault/pull/29326)]
* core: Add `removed_from_cluster` field to sys/seal-status and vault status output to indicate whether the node has been removed from the HA cluster. [[GH-28938](https://github.com/hashicorp/vault/pull/28938)]
* core: Add a mount tuneable that trims trailing slashes of request paths during POST.  Needed to support CMPv2 in PKI. [[GH-28752](https://github.com/hashicorp/vault/pull/28752)]
* core: Add activation flags. A mechanism for users to opt in to new functionality at a convenient time. Previously used only in Enterprise for SecretSync, activation flags are now available in CE for future features to use. [[GH-29237](https://github.com/hashicorp/vault/pull/29237)]
* core: Added new `enable_post_unseal_trace` and `post_unseal_trace_directory` config options to generate Go traces during the post-unseal step for debug purposes. [[GH-28895](https://github.com/hashicorp/vault/pull/28895)]
* core: Config reloading on SIGHUP now includes some Raft settings, which are now also present in `/sys/config/state/sanitized` output. [[GH-29485](https://github.com/hashicorp/vault/pull/29485)]
* core: add support for reading certain sensitive seal wrap and managed key (enterprise) configuration values from the environment or files. [[GH-29402](https://github.com/hashicorp/vault/pull/29402)]
* events (enterprise): Send events downstream to a performance standby node only when there is a subscriber on the standby node with a filter matching the events. [[GH-29618](https://github.com/hashicorp/vault/pull/29618)]
* events (enterprise): Send events downstream to performance standby nodes in a cluster, removing the need to redirect client event subscriptions to the active node. [[GH-29470](https://github.com/hashicorp/vault/pull/29470)]
* events (enterprise): Use the `path` event metadata field when authorizing a client's `subscribe` capability for consuming an event, instead of requiring `data_path` to be present in the event metadata.
* identity: Added reporting in Vault logs during unseal to help identify any
duplicate identify resources in storage. [[GH-29325](https://github.com/hashicorp/vault/pull/29325)]
* physical/dynamodb: Allow Vault to modify its DynamoDB table and use per-per-request billing mode. [[GH-29371](https://github.com/hashicorp/vault/pull/29371)]
* raft/autopilot: We've updated the autopilot reconciliation logic (by updating the raft-autopilot dependency to v0.3.0) to avoid artificially increasing the quorum in presence of an unhealthy node. Now autopilot will start the reconciliation process by attempting to demote a failed voter node before any promotions, fixing the issue where Vault would initially increase quorum when faced with a failure of a voter node. In certain configurations, especially when using Vault Enterprise Redundancy Zones and losing a voter then a non-voter in quick succession, this would lead to a loss of quorum and cluster failure. [[GH-29306](https://github.com/hashicorp/vault/pull/29306)]
* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.0.0-20241003195753-88fef418d705
* sdk/helper: utitilize a randomly seeded cryptographic determinstic random bit generator for
RSA key generation when using slow random sources, speeding key generation
considerably. [[GH-29020](https://github.com/hashicorp/vault/pull/29020)]
* sdk: Add Vault build date to system view plugin environment response [[GH-29082](https://github.com/hashicorp/vault/pull/29082)]
* sdk: Add helpers and CE stubs for plugins to communicate with Rotation Manager (Enterprise). [[GH-29273](https://github.com/hashicorp/vault/pull/29273)]
* secret/pki: Introduce a new value `always_enforce_err` within `leaf_not_after_behavior` to force the error in all circumstances such as CA issuance and ACME requests if requested TTL values are beyond the issuer's NotAfter. [[GH-28907](https://github.com/hashicorp/vault/pull/28907)]
* secrets(pki): Error if attempt to set a manual chain on an issuer that can't issue any certificate. [[GH-29473](https://github.com/hashicorp/vault/pull/29473)]
* secrets-sync (enterprise): No longer attempt to unsync a random UUID secret name in GCP upon destination creation.
* secrets-sync (enterprise): add support for user-managed encryption keys in GCP secrets sync destinations.
* secrets/aws: add fallback endpoint and region parameters to sts configuration [[GH-29051](https://github.com/hashicorp/vault/pull/29051)]
* secrets/pki (enterprise): Add issuer configuration fields which allow disabling specific validations on certificate chains.
* secrets/pki: Add ACME error types to errors encountered during challenge validation. [[GH-28678](https://github.com/hashicorp/vault/pull/28678)]
* secrets/pki: Add `serial_number_source` option to PKI roles to control the source for the subject serial number. [[GH-29369](https://github.com/hashicorp/vault/pull/29369)]
* secrets/pki: Add a CRL entry limit to prevent runaway revocations from overloading Vault, reconfigurable with max_crl_entries on the CRL config. [[GH-28654](https://github.com/hashicorp/vault/pull/28654)]
* secrets/pki: Add a new set of APIs that allow listing ACME account key ids, retrieving ACME account information along with the associated order and certificate information and updating an ACME account's status [[GH-29173](https://github.com/hashicorp/vault/pull/29173)]
* secrets/pki: Add a warning when issuers are updated with validations that cause the issuer to be non-functional.
* secrets/pki: Add necessary validation configuration fields to CMPv2 to enable customers with different clients.
* secrets/pki: Complete the set of name constraints parameters by adding permitted_email_addresses, permitted_ip_ranges, permitted_uri_domains, excluded_dns_domains, excluded_email_addresses, excluded_ip_ranges, and excluded_uri_domains; this makes it possible for the name constraints extension to be fully specified when creating root and intermediate CA certificates. [[GH-29245](https://github.com/hashicorp/vault/pull/29245)]
* secrets/transit: Add support for RSA padding scheme pkcs1v15 for encryption [[GH-25486](https://github.com/hashicorp/vault/pull/25486)]
* storage/dynamodb: Pass context to AWS SDK calls [[GH-27927](https://github.com/hashicorp/vault/pull/27927)]
* storage/s3: Pass context to AWS SDK calls [[GH-27927](https://github.com/hashicorp/vault/pull/27927)]
* ui (enterprise): Allow WIF configuration on the Azure secrets engine. [[GH-29047](https://github.com/hashicorp/vault/pull/29047)]
* ui (enterprise): Allow WIF configuration on the GCP secrets engine. [[GH-29423](https://github.com/hashicorp/vault/pull/29423)]
* ui: Add button to copy secret path in kv v1 and v2 secrets engines [[GH-28629](https://github.com/hashicorp/vault/pull/28629)]
* ui: Add identity_token_key to mount view for the GCP and Azure Secret engines. [[GH-28822](https://github.com/hashicorp/vault/pull/28822)]
* ui: Add support for the name constraints extension to be fully specified when creating root and intermediate CA certificates. [[GH-29263](https://github.com/hashicorp/vault/pull/29263)]
* ui: Adds ability to edit, create, and view the Azure secrets engine configuration. [[GH-29047](https://github.com/hashicorp/vault/pull/29047)]
* ui: Adds ability to edit, create, and view the GCP secrets engine configuration. [[GH-29423](https://github.com/hashicorp/vault/pull/29423)]
* ui: Adds copy button to identity entity, alias and mfa method IDs [[GH-28742](https://github.com/hashicorp/vault/pull/28742)]
* ui: Adds navigation for LDAP hierarchical libraries [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]
* ui: Adds navigation for LDAP hierarchical roles [[GH-28824](https://github.com/hashicorp/vault/pull/28824)]
* ui: Adds params to postgresql database to improve editing a connection in the web browser. [[GH-29200](https://github.com/hashicorp/vault/pull/29200)]
* ui: Application static breadcrumbs should be formatted in title case. [[GH-29206](https://github.com/hashicorp/vault/pull/29206)]
* ui: Replace KVv2 json secret details view with Hds::CodeBlock component allowing users to search the full secret height. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]
* website/docs: changed outdated reference to consul-helm repository to consul-k8s repository. [[GH-28825](https://github.com/hashicorp/vault/pull/28825)]

BUG FIXES:

* UI: Fix missing Client Count card when running as a Vault Dedicated cluster [[GH-29241](https://github.com/hashicorp/vault/pull/29241)]
* activity: Include activity records from clients created by deleted or disabled auth mounts in Export API response. [[GH-29376](https://github.com/hashicorp/vault/pull/29376)]
* activity: Show activity records from clients created in deleted namespaces when activity log is queried from admin namespace. [[GH-29432](https://github.com/hashicorp/vault/pull/29432)]
* agent: Fix chown error running agent on Windows with an auto-auth file sinks. [[GH-28748](https://github.com/hashicorp/vault/pull/28748)]
* agent: Fixed an issue where giving the agent multiple config files could cause the merged config to be incorrect
when `template_config` is set in one of the config files. [[GH-29680](https://github.com/hashicorp/vault/pull/29680)]
* audit: Fixing TestAudit_enableAudit_fallback_two test failure.
* audit: Prevent users from enabling multiple audit devices of file type with the same file_path to write to. [[GH-28751](https://github.com/hashicorp/vault/pull/28751)]
* auth/ldap: Fixed an issue where debug level logging was not emitted. [[GH-28881](https://github.com/hashicorp/vault/pull/28881)]
* auth/radius: Fixed an issue where usernames with upper case characters where not honored [[GH-28884](https://github.com/hashicorp/vault/pull/28884)]
* autosnapshots (enterprise): Fix an issue where snapshot size metrics were not reported for cloud-based storage.
* cli: Fixed a CLI precedence issue where -agent-address didn't override VAULT_AGENT_ADDR as it should [[GH-28574](https://github.com/hashicorp/vault/pull/28574)]
* core/api: Added missing LICENSE files to API sub-modules to ensure Go module tooling recognizes MPL-2.0 license. [[GH-27920](https://github.com/hashicorp/vault/pull/27920)]
* core/managed-keys (enterprise): Allow mechanism numbers above 32 bits in PKCS#11 managed keys.
* core/metrics: Fix unlocked mounts read for usage reporting. [[GH-29091](https://github.com/hashicorp/vault/pull/29091)]
* core/seal (enterprise): Fix bug that caused seal generation information to be replicated, which prevented disaster recovery and performance replication clusters from using their own seal high-availability configuration.
* core/seal (enterprise): Fix problem with nodes unable to join Raft clusters with Seal High Availability enabled. [[GH-29117](https://github.com/hashicorp/vault/pull/29117)]
* core/seal: Azure seals required client_secret, preventing use of managed service identities and user assigned identities. [[GH-29499](https://github.com/hashicorp/vault/pull/29499)]
* core/seal: Fix an issue that could cause reading from sys/seal-backend-status to return stale information. [[GH-28631](https://github.com/hashicorp/vault/pull/28631)]
* core: Fix Azure authentication for seal/managed keys to work for both federated workload identity and managed user identities.  Fixes regression for federated workload identities. [[GH-29792](https://github.com/hashicorp/vault/pull/29792)]
* core: Fix an issue where duplicate identity aliases in storage could be merged
inconsistently during different unseal events or on different servers. [[GH-28867](https://github.com/hashicorp/vault/pull/28867)]
* core: Fix bug when if failing to persist the barrier keyring to track encryption counts, the number of outstanding encryptions remains added to the count, overcounting encryptions. [[GH-29506](https://github.com/hashicorp/vault/pull/29506)]
* core: Fixed panic seen when performing help requests without /v1/ in the URL. [[GH-28669](https://github.com/hashicorp/vault/pull/28669)]
* core: Improved an internal helper function that sanitizes paths by adding a check for leading backslashes
in addition to the existing check for leading slashes. [[GH-28878](https://github.com/hashicorp/vault/pull/28878)]
* core: Prevent integer overflows of the barrier key counter on key rotation requests [[GH-29176](https://github.com/hashicorp/vault/pull/29176)]
* core: fix bug in seal unwrapper that caused high storage latency in Vault CE. For every storage read request, the
seal unwrapper was performing the read twice, and would also issue an unnecessary storage write. [[GH-29050](https://github.com/hashicorp/vault/pull/29050)]
* core: fix issue when attempting to re-bootstrap HA when using Raft as HA but not storage [[GH-18615](https://github.com/hashicorp/vault/pull/18615)]
* database/mssql: Fix a bug where contained databases would silently fail root rotation if a custom root rotation statement was not provided. [[GH-29399](https://github.com/hashicorp/vault/pull/29399)]
* database: Fix a bug where static role passwords are erroneously rotated across backend restarts when using skip import rotation. [[GH-29537](https://github.com/hashicorp/vault/pull/29537)]
* export API: Normalize the start_date parameter to the start of the month as is done in the sys/counters API to keep the results returned from both of the API's consistent. [[GH-29562](https://github.com/hashicorp/vault/pull/29562)]
* identity/oidc (enterprise): Fix delays in rotation and invalidation of OIDC keys when there are too many namespaces.
The Cache-Control header returned by the identity/oidc/.well-known/keys endpoint now depends only on the named keys for
the queried namespace. [[GH-29312](https://github.com/hashicorp/vault/pull/29312)]
* kmip (enterprise): Use the default KMIP port for IPv6 addresses missing a port, for the listen_addrs configuration field, in order to match the existing IPv4 behavior
* namespaces (enterprise): Fix issue where namespace patch requests to a performance secondary would not patch the namespace's metadata.
* plugins: Fix a bug that causes zombie dbus-daemon processes on certain systems. [[GH-29334](https://github.com/hashicorp/vault/pull/29334)]
* proxy: Fix chown error running proxy on Windows with an auto-auth file sink. [[GH-28748](https://github.com/hashicorp/vault/pull/28748)]
* sdk/database: Fix a bug where slow database connections can cause goroutines to be blocked. [[GH-29097](https://github.com/hashicorp/vault/pull/29097)]
* secret/aws: Fixed potential panic after step-down and the queue has not repopulated. [[GH-28330](https://github.com/hashicorp/vault/pull/28330)]
* secret/db: Update static role rotation to generate a new password after 2 failed attempts.
Unblocks customers that were stuck in a failing loop when attempting to rotate static role passwords. [[GH-28989](https://github.com/hashicorp/vault/pull/28989)]
* secret/pki: Fix a bug that prevents PKI issuer field enable_aia_url_templating
to be set to false. [[GH-28832](https://github.com/hashicorp/vault/pull/28832)]
* secrets-sync (enterprise): Add new parameters for destination configs to specify allowlists for IP's and ports.
* secrets-sync (enterprise): Fixed issue where secret-key granularity destinations could sometimes cause a panic when loading a sync status.
* secrets/aws: Add sts_region parameter to root config for STS API calls. [[GH-22726](https://github.com/hashicorp/vault/pull/22726)]
* secrets/aws: Fix issue with static credentials not rotating after restart or leadership change. [[GH-28775](https://github.com/hashicorp/vault/pull/28775)]
* secrets/database: Fix a bug where a global database plugin reload exits if any of the database connections are not available [[GH-29519](https://github.com/hashicorp/vault/pull/29519)]
* secrets/openldap: Update static role rotation to generate a new password after 2 failed attempts.
Unblocks customers that were stuck in a failing loop when attempting to rotate static role passwords. [[GH-29131](https://github.com/hashicorp/vault/pull/29131)]
* secrets/pki: Address issue with ACME HTTP-01 challenges failing for IPv6 IPs due to improperly formatted URLs [[GH-28718](https://github.com/hashicorp/vault/pull/28718)]
* secrets/pki: Fix a bug that prevented the full CA chain to be used when enforcing name constraints. [[GH-29255](https://github.com/hashicorp/vault/pull/29255)]
* secrets/pki: fixes issue #28749 requiring all chains to be single line of authority. [[GH-29342](https://github.com/hashicorp/vault/pull/29342)]
* secrets/ssh: Return the flag `allow_empty_principals` in the read role api when key_type is "ca" [[GH-28901](https://github.com/hashicorp/vault/pull/28901)]
* secrets/transform (enterprise): Fix nil panic when accessing a partially setup database store.
* secrets/transit: Fix a race in which responses from the key update api could contain results from another subsequent update [[GH-28839](https://github.com/hashicorp/vault/pull/28839)]
* sentinel (enterprise): No longer report inaccurate log messages for when failing an advisory policy.
* ui (enterprise): Fixes login to web UI when MFA is enabled for SAML auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui (enterprise): Fixes token renewal to ensure capability checks are performed in the relevant namespace, resolving 'Not authorized' errors for resources that users have permission to access. [[GH-29416](https://github.com/hashicorp/vault/pull/29416)]
* ui/database: Fixes 'cannot update static username' error when updating static role's rotation period [[GH-29498](https://github.com/hashicorp/vault/pull/29498)]
* ui: Allow users to search the full json object within the json code-editor edit/create view. [[GH-28808](https://github.com/hashicorp/vault/pull/28808)]
* ui: Decode `connection_url` to fix database connection updates (i.e. editing connection config, deleting roles) failing when urls include template variables. [[GH-29114](https://github.com/hashicorp/vault/pull/29114)]
* ui: Fixes login to web UI when MFA is enabled for OIDC (i.e. azure, auth0) and Okta auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui: Fixes navigation for quick actions in LDAP roles' popup menu [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]
* ui: Fixes rendering issues of LDAP dynamic and static roles with the same name [[GH-28824](https://github.com/hashicorp/vault/pull/28824)]
* ui: Fixes text overflow on Secrets engines and Auth Engines list views for long names & descriptions [[GH-29430](https://github.com/hashicorp/vault/pull/29430)]
* ui: MFA methods now display the namespace path instead of the namespace id. [[GH-29588](https://github.com/hashicorp/vault/pull/29588)]
* ui: No longer running decodeURIComponent on KVv2 list view allowing percent encoded data-octets in path name. [[GH-28698](https://github.com/hashicorp/vault/pull/28698)]
* vault/diagnose: Fix time to expiration reporting within the TLS verification to not be a month off. [[GH-29128](https://github.com/hashicorp/vault/pull/29128)]

## 1.18.15 Enterprise
### September 24, 2025

SECURITY:

* core: Update github.com/hashicorp/go-getter to fix security vulnerability GHSA-wjrx-6529-hcj3.
* core: Update github.com/ulikunitz/xz to fix security vulnerability GHSA-25xm-hr59-7c27.

CHANGES:

* core: Bump Go version to 1.24.7.
* core: Updates post-install script to print updated license information
* database/snowflake: Update plugin to [v0.12.3](https://github.com/hashicorp/vault-plugin-database-snowflake/releases/tag/v0.12.3)
* sdk: Upgrade to go-secure-stdlib/plugincontainer@v0.4.2, which also bumps github.com/docker/docker to v28.3.3+incompatible

IMPROVEMENTS:

* Raft: Auto-join will now allow you to enforce IPv4 on networks that allow IPv6 and dual-stack enablement, which is on by default in certain regions.
* auth/cert: Support RFC 9440 colon-wrapped Base64 certificates in `x_forwarded_for_client_cert_header`, to fix TLS certificate auth errors with Google Cloud Application Load Balancer.
* core (enterprise): Updated code to support FIPS 140-3 compliant algorithms.

BUG FIXES:

* auth/cert: Recover from partially populated caches of trusted certificates if one or more certificates fails to load.
* secrets/transit: Fix error when using ed25519 keys that were imported with derivation enabled
* sys/mounts: enable unsetting allowed_response_headers

## 1.18.14 Enterprise
### August 28, 2025

FEATURES:

* **IBM RACF Static Role Password Phrase Management**: Add support for static role password phrase management to the LDAP secrets engine. (https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/184)

CHANGES:

* core: Bump Go version to 1.23.12.
* http: Add JSON configurable limits to HTTP handling for JSON payloads: `max_json_depth`, `max_json_string_value_length`, `max_json_object_entry_count`, `max_json_array_element_count`.
* secrets/openldap: update plugin to v0.14.7

BUG FIXES:

* core (enterprise): fix a bug where issuing a token in a namespace used root auth configuration instead of namespace auth configuration
* core/metrics: Add service name prefix for core HA metrics to avoid duplicate, zero-value metrics.
* core/seal: When Seal-HA is enabled, make it an error to persist the barrier
keyring when not all seals are healthy.  This prevents the possibility of
failing to unseal when a different subset of seals are healthy than were
healthy at last write.
* raft/autopilot: Fixes an issue with enterprise redundancy zones where, if the leader was in a redundancy zone and that leader becomes unavailable, the node would become an unzoned voter. This can artificially inflate the required number of nodes for quorum, leading to a situation where the cluster cannot recover if another leader subsequently becomes unavailable. Vault will now keep an unavailable node in its last known redundancy zone as a non-voter.
* secrets/database/postgresql: Support for multiline statements in the `rotation_statements` field.

## 1.18.13 Enterprise
### August 06, 2025

SECURITY:

* auth/ldap: fix MFA/TOTP enforcement bypass when username_as_alias is enabled [[GH-31427](https://github.com/hashicorp/vault/pull/31427),[HCSEC-2025-20](https://discuss.hashicorp.com/t/hcsec-2025-20-vault-ldap-mfa-enforcement-bypass-when-using-username-as-alias/76092)].

BUG FIXES:

* identity/mfa: revert cache entry change from #31217 and document cache entry values [[GH-31421](https://github.com/hashicorp/vault/pull/31421)]

## 1.18.12 Enterprise
### July 25, 2025

SECURITY:

* audit: **breaking change** privileged vault operator may execute code on the underlying host (CVE-2025-6000). Vault will not unseal if the only configured file audit device has executable permissions (e.g., 0777, 0755). See recent [breaking change](https://developer.hashicorp.com/vault/docs/updates/important-changes#breaking-changes) docs for more details. [[GH-31211](https://github.com/hashicorp/vault/pull/31211),[HCSEC-2025-14](https://discuss.hashicorp.com/t/hcsec-2025-14-privileged-vault-operator-may-execute-code-on-the-underlying-host/76033)]
* auth/userpass: timing side-channel in vault's userpass auth method (CVE-2025-6011)[HCSEC-2025-15](https://discuss.hashicorp.com/t/hcsec-2025-15-timing-side-channel-in-vault-s-userpass-auth-method/76034)
* core/login: vault userpass and ldap user lockout bypass (CVE-2025-6004). update alias lookahead to respect username case for LDAP and username/password. [[GH-31352](https://github.com/hashicorp/vault/pull/31352),[HCSEC-2025-16](https://discuss.hashicorp.com/t/hcsec-2025-16-vault-userpass-and-ldap-user-lockout-bypass/76035)]
* secrets/totp: vault totp secrets engine code reuse (CVE-2025-6014) [[GH-31246](https://github.com/hashicorp/vault/pull/31246),[HCSEC-2025-17](https://discuss.hashicorp.com/t/hcsec-2025-17-vault-totp-secrets-engine-code-reuse/76036)]
* auth/cert: vault certificate auth method did not validate common name for non-ca certificates (CVE-2025-6037). test non-CA cert equality on login matching instead of individual fields. [[GH-31210](https://github.com/hashicorp/vault/pull/31210),[HCSEC-2025-18](https://discuss.hashicorp.com/t/hcsec-2025-18-vault-certificate-auth-method-did-not-validate-common-name-for-non-ca-certificates/76037)]
* core/mfa: vault login mfa bypass of rate limiting and totp token reuse (CVE-2025-6015) [[GH-31217](https://github.com/hashicorp/vault/pull/31297),[HCSEC-2025-19](https://discuss.hashicorp.com/t/hcsec-2025-19-vault-login-mfa-bypass-of-rate-limiting-and-totp-token-reuse/76038)]

BUG FIXES:

* core/seal (enterprise): Fix a bug that caused the seal rewrap process to abort in the presence of partially sealed entries.
* kmip (enterprise): Fix a panic that can happen when a KMIP client makes a request before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* plugins: Fix panics that can occur when a plugin audits a request or response before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* replication (enterprise): Fix bug with mount invalidations consuming excessive memory.
* secrets-sync (enterprise): Unsyncing secret-key granularity associations will no longer give a misleading error about a failed unsync operation that did indeed succeed.

## 1.18.11 Enterprise
### June 25, 2025

SECURITY:

* core: require a nonce when cancelling a rekey operation that was initiated within the last 10 minutes. [[GH-30794](https://github.com/hashicorp/vault/pull/30794)],[[HCSEC-2025-11](https://discuss.hashicorp.com/t/hcsec-2025-11-vault-vulnerable-to-recovery-key-cancellation-denial-of-service/75570)]
* core/identity: vault root namespace operator may elevate privileges (CVE-2025-5999). Fix string contains check in Identity APIs to be case-insensitive. [[GH-31045](https://github.com/hashicorp/vault/pull/31045),[HCSEC-2025-13](https://discuss.hashicorp.com/t/hcsec-2025-13-vault-root-namespace-operator-may-elevate-token-privileges/76032)]

CHANGES:

* api: Update the default API client to check for the `Retry-After` header and, if it exists, wait for the specified duration before retrying the request. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* auth/azure: Update plugin to v0.19.5
* core: Bump Go version to 1.23.10.
* quotas/rate-limit: Round up the `Retry-After` value to the nearest second when calculating the retry delay. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* secrets/azure: Update plugin to v0.20.3
* secrets/database: Update vault-plugin-database-snowflake to v0.12.2

BUG FIXES:

* secrets/database: Treat all rotation_schedule values as UTC to ensure consistent behavior. [[GH-30606](https://github.com/hashicorp/vault/pull/30606)]
* secrets/transit (enterprise): ensure verify endpoint always returns valid field in batch_results with CMAC
* secrets/transit: ensure verify endpoint always returns valid field in batch_results with HMAC [[GH-30852](https://github.com/hashicorp/vault/pull/30852)]
* ui: Redirect users authenticating with Vault as an OIDC provider to log in again when token expires. [[GH-30838](https://github.com/hashicorp/vault/pull/30838)]

## 1.18.10 Enterprise
### May 30, 2025

CHANGES:

* Update vault-plugin-auth-cf to v0.19.2
* auth/azure: Upgrade plugin to v0.19.4
* database/snowflake: Update plugin to v0.12.1

IMPROVEMENTS:

* ui: Replaces all instances of the deprecated event.keyCode with event.key [[GH-30493](https://github.com/hashicorp/vault/pull/30493)]

BUG FIXES:

* plugins (enterprise): Fix an issue where Enterprise plugins can't run on a standby node
when it becomes active because standby nodes don't extract the artifact when the plugin
is registered. Remove extracting from Vault and require the operator to place
the extracted artifact in the plugin directory before registration.

## 1.18.9 Enterprise
### April 30, 2025

SECURITY:

* core: vault may expose sensitive information in error logs when processing malformed data with the kv v2 plugin[[GH-30388](https://github.com/hashicorp/vault/pull/30388), [HCSEC-2025-09](https://discuss.hashicorp.com/t/hcsec-2025-09-vault-may-expose-sensitive-information-in-error-logs-when-processing-malformed-data-with-the-kv-v2-plugin/74717)]

BUG FIXES:

* core (enterprise): fix issue with errors being swallowed on failed HSM logins.
* database: Prevent static roles created in versions prior to 1.15.0 from rotating on backend restart. [[GH-30320](https://github.com/hashicorp/vault/pull/30320)]
* database: no longer incorrectly add an "unrecognized parameters" warning for certain SQL database secrets config operations when another warning is returned [[GH-30327](https://github.com/hashicorp/vault/pull/30327)]

## 1.18.8 Enterprise
### April 18, 2025

CHANGES:

* core: Bump Go version to 1.23.7
* core: Bump Go version to 1.23.8

BUG FIXES:

* secrets/openldap: Prevent static role rotation on upgrade when `NextVaultRotation` is nil. Fixes an issue where static roles were unexpectedly rotated after upgrade due to a missing `NextVaultRotation` value. Now sets it to either `LastVaultRotation + RotationPeriod` or `now + RotationPeriod`. [[GH-30265](https://github.com/hashicorp/vault/pull/30265)]
* secrets/pki (enterprise): Address a parsing bug that rejected CMPv2 requests containing a validity field.
* secrets/pki: fix a bug where key_usage was ignored when generating root certificates, and signing certain intermediate certificates. [[GH-30034](https://github.com/hashicorp/vault/pull/30034)]
* secrets/transit: fix a panic when rotating on a managed key returns an error [[GH-30214](https://github.com/hashicorp/vault/pull/30214)]

## 1.18.7 Enterprise
### April 4, 2025

SECURITY:

* auth/azure: Update plugin to v0.19.3. Login requires `resource_group_name`, `vm_name`, and `vmss_name` to match token claims [[HCSEC-2025-07](https://discuss.hashicorp.com/t/hcsec-2025-07-vault-s-azure-authentication-method-bound-location-restriction-could-be-bypassed-on-login/74716)].

CHANGES:

* core: Verify that the client IP address extracted from an X-Forwarded-For header is a valid IPv4 or IPv6 address [[GH-29774](https://github.com/hashicorp/vault/pull/29774)]

IMPROVEMENTS:

* core (enterprise): report errors from the underlying seal when getting entropy.
* storage/raft: Upgrade hashicorp/raft library to v1.7.3 which includes additional logging on the leader when opening and sending a snapshot to a follower. [[GH-29976](https://github.com/hashicorp/vault/pull/29976)]

BUG FIXES:

* auth/ldap: Fix a bug that does not properly delete users and groups by first converting their names to lowercase when case senstivity option is off. [[GH-29922](https://github.com/hashicorp/vault/pull/29922)]
* core: Fix Azure authentication for seal/managed keys to work for both federated workload identity and managed user identities.  Fixes regression for federated workload identities. [[GH-29792](https://github.com/hashicorp/vault/pull/29792)]
* core: Fix a bug that prevents certain loggers from writing to a log file. [[GH-29917](https://github.com/hashicorp/vault/pull/29917)]
* plugins (enterprise): Fix plugin registration with artifact when a binary for the same plugin is already present in the plugin directory.
* plugins: plugin registration should honor the `plugin_tmpdir` config [[GH-29978](https://github.com/hashicorp/vault/pull/29978)]
* secrets/azure: Upgrade plugin to v0.20.2 which reverts role name changes to no longer be a GUID.
* secrets/pki: Fix a bug that prevents enabling automatic tidying of the CMPv2 nonce store. [[GH-29852](https://github.com/hashicorp/vault/pull/29852)]

## 1.18.6 Enterprise
### March 5, 2025

IMPROVEMENTS:

* secrets/pki: Add necessary validation configuration fields to CMPv2 to enable customers with different clients.

BUG FIXES:

* agent: Fixed an issue where giving the agent multiple config files could cause the merged config to be incorrect
when `template_config` is set in one of the config files. [[GH-29680](https://github.com/hashicorp/vault/pull/29680)]
* secrets/database: Fix a bug where a global database plugin reload exits if any of the database connections are not available [[GH-29519](https://github.com/hashicorp/vault/pull/29519)]
* core: Fix Azure authentication for seal/managed keys to work for both federated workload identity and managed user identities.  Fixes regression for federated workload identities. [[GH-29792](https://github.com/hashicorp/vault/pull/29792)]

## 1.18.5
### February 25, 2025

SECURITY:

* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.2.0

CHANGES:

* core: Bump Go version to 1.23.6
* raft/autopilot (enterprise): Alongside the CE autopilot update, update raft-autopilot-enterprise library to v0.3.0 and add enterprise-specific regression testing.
* secrets/openldap: Update plugin to v0.14.5 [[GH-29551](https://github.com/hashicorp/vault/pull/29551)]

FEATURES:

* **Plugins**: Allow Enterprise plugins to run externally on Vault Enterprise only.

IMPROVEMENTS:

* raft/autopilot: We've updated the autopilot reconciliation logic (by updating the raft-autopilot dependency to v0.3.0) to avoid artificially increasing the quorum in presence of an unhealthy node. Now autopilot will start the reconciliation process by attempting to demote a failed voter node before any promotions, fixing the issue where Vault would initially increase quorum when faced with a failure of a voter node. In certain configurations, especially when using Vault Enterprise Redundancy Zones and losing a voter then a non-voter in quick succession, this would lead to a loss of quorum and cluster failure. [[GH-29306](https://github.com/hashicorp/vault/pull/29306)]
* ui: Application static breadcrumbs should be formatted in title case. [[GH-29206](https://github.com/hashicorp/vault/pull/29206)]

BUG FIXES:

* activity: Show activity records from clients created in deleted namespaces when activity log is queried from admin namespace. [[GH-29432](https://github.com/hashicorp/vault/pull/29432)]
* core/managed-keys (enterprise): Allow mechanism numbers above 32 bits in PKCS#11 managed keys.
* core: Fix bug when if failing to persist the barrier keyring to track encryption counts, the number of outstanding encryptions remains added to the count, overcounting encryptions. [[GH-29506](https://github.com/hashicorp/vault/pull/29506)]
* database: Fix a bug where static role passwords are erroneously rotated across backend restarts when using skip import rotation. [[GH-29537](https://github.com/hashicorp/vault/pull/29537)]
* export API: Normalize the start_date parameter to the start of the month as is done in the sys/counters API to keep the results returned from both of the API's consistent. [[GH-29562](https://github.com/hashicorp/vault/pull/29562)]
* identity/oidc (enterprise): Fix delays in rotation and invalidation of OIDC keys when there are too many namespaces.
The Cache-Control header returned by the identity/oidc/.well-known/keys endpoint now depends only on the named keys for
the queried namespace. [[GH-29312](https://github.com/hashicorp/vault/pull/29312)]
* secrets-sync (enterprise): Add new parameters for destination configs to specify allowlists for IP's and ports.
* secrets/pki: fixes issue #28749 requiring all chains to be single line of authority. [[GH-29342](https://github.com/hashicorp/vault/pull/29342)]
* ui (enterprise): Fixes token renewal to ensure capability checks are performed in the relevant namespace, resolving 'Not authorized' errors for resources that users have permission to access. [[GH-29416](https://github.com/hashicorp/vault/pull/29416)]
* ui/database: Fixes 'cannot update static username' error when updating static role's rotation period [[GH-29498](https://github.com/hashicorp/vault/pull/29498)]
* ui: Fixes text overflow on Secrets engines and Auth Engines list views for long names & descriptions [[GH-29430](https://github.com/hashicorp/vault/pull/29430)]
* ui: MFA methods now display the namespace path instead of the namespace id. [[GH-29588](https://github.com/hashicorp/vault/pull/29588)]

## 1.18.4
### January 30, 2025

CHANGES:

* auth/cf: Update plugin to v0.19.1 [[GH-29295](https://github.com/hashicorp/vault/pull/29295)]
* sdk: Updated golang and dependency versions to be consistent across core, API, SDK to address [[GO-2024-3333](https://pkg.go.dev/vuln/GO-2024-3333)] and ensure version consistency [[GH-29422](https://github.com/hashicorp/vault/pull/29422)]

IMPROVEMENTS:

* plugins (enterprise): The Database secrets engine now allows skipping the automatic rotation of static roles during import.
* events (enterprise): Use the `path` event metadata field when authorizing a client's `subscribe` capability for consuming an event, instead of requiring `data_path` to be present in the event metadata.
* ui: Adds navigation for LDAP hierarchical libraries [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]
* ui: Adds params to postgresql database to improve editing a connection in the web browser. [[GH-29200](https://github.com/hashicorp/vault/pull/29200)]

BUG FIXES:

* activity: Include activity records from clients created by deleted or disabled auth mounts in Export API response. [[GH-29376](https://github.com/hashicorp/vault/pull/29376)]
* core: Prevent integer overflows of the barrier key counter on key rotation requests [[GH-29176](https://github.com/hashicorp/vault/pull/29176)]
* database/mssql: Fix a bug where contained databases would silently fail root rotation if a custom root rotation statement was not provided. [[GH-29399](https://github.com/hashicorp/vault/pull/29399)]
* plugins: Fix a bug that causes zombie dbus-daemon processes on certain systems. [[GH-29334](https://github.com/hashicorp/vault/pull/29334)]
* sdk/database: Fix a bug where slow database connections can cause goroutines to be blocked. [[GH-29097](https://github.com/hashicorp/vault/pull/29097)]
* secrets/pki: Fix a bug that prevented the full CA chain to be used when enforcing name constraints. [[GH-29255](https://github.com/hashicorp/vault/pull/29255)]
* sentinel (enterprise): No longer report inaccurate log messages for when failing an advisory policy.
* ui (enterprise): Fixes login to web UI when MFA is enabled for SAML auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui: Fixes login to web UI when MFA is enabled for OIDC (i.e. azure, auth0) and Okta auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui: Fixes navigation for quick actions in LDAP roles' popup menu [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]


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

## 1.17.18 Enterprise
### June 25, 2025

SECURITY:

* core: require a nonce when cancelling a rekey operation that was initiated within the last 10 minutes. [[GH-30794](https://github.com/hashicorp/vault/pull/30794)],[[HCSEC-2025-11](https://discuss.hashicorp.com/t/hcsec-2025-11-vault-vulnerable-to-recovery-key-cancellation-denial-of-service/75570)]
* core/identity: vault root namespace operator may elevate privileges (CVE-2025-5999). Fix string contains check in Identity APIs to be case-insensitive. [[GH-31045](https://github.com/hashicorp/vault/pull/31045),[HCSEC-2025-13](https://discuss.hashicorp.com/t/hcsec-2025-13-vault-root-namespace-operator-may-elevate-token-privileges/76032)]

CHANGES:

* api: Update the default API client to check for the `Retry-After` header and, if it exists, wait for the specified duration before retrying the request. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* auth/azure: Update plugin to v0.18.4
* core: Bump Go version to 1.23.10
* quotas/rate-limit: Round up the `Retry-After` value to the nearest second when calculating the retry delay. [[GH-30887](https://github.com/hashicorp/vault/pull/30887)]
* secrets/azure: Update plugin to v0.19.4
* secrets/database: Update vault-plugin-database-snowflake to v0.11.2

BUG FIXES:

* secrets/database: Treat all rotation_schedule values as UTC to ensure consistent behavior. [[GH-30606](https://github.com/hashicorp/vault/pull/30606)]
* secrets/transit (enterprise): ensure verify endpoint always returns valid field in batch_results with CMAC
* secrets/transit: ensure verify endpoint always returns valid field in batch_results with HMAC [[GH-30852](https://github.com/hashicorp/vault/pull/30852)]
* ui: Redirect users authenticating with Vault as an OIDC provider to log in again when token expires. [[GH-30838](https://github.com/hashicorp/vault/pull/30838)]


## 1.17.17 Enterprise
### May 30, 2025

CHANGES:

* Update vault-plugin-auth-cf to v0.18.2
* auth/azure: Upgrade plugin to v0.18.3
* database/snowflake: Update plugin to v0.11.1

IMPROVEMENTS:

* ui: Replaces all instances of the deprecated event.keyCode with event.key [[GH-30493](https://github.com/hashicorp/vault/pull/30493)]

BUG FIXES:

* plugins (enterprise): Fix an issue where Enterprise plugins can't run on a standby node
when it becomes active because standby nodes don't extract the artifact when the plugin
is registered. Remove extracting from Vault and require the operator to place
the extracted artifact in the plugin directory before registration.

## 1.17.16 Enterprise
### April 30, 2025

SECURITY:

* core: vault may expose sensitive information in error logs when processing malformed data with the kv v2 plugin[[GH-30388](https://github.com/hashicorp/vault/pull/30388), [HCSEC-2025-09](https://discuss.hashicorp.com/t/hcsec-2025-09-vault-may-expose-sensitive-information-in-error-logs-when-processing-malformed-data-with-the-kv-v2-plugin/74717)]

BUG FIXES:

* core (enterprise): fix issue with errors being swallowed on failed HSM logins.
* database: Prevent static roles created in versions prior to 1.15.0 from rotating on backend restart. [[GH-30320](https://github.com/hashicorp/vault/pull/30320)]
* database: no longer incorrectly add an "unrecognized parameters" warning for certain SQL database secrets config operations when another warning is returned [[GH-30327](https://github.com/hashicorp/vault/pull/30327)]

## 1.17.15 Enterprise
### April 18, 2025

CHANGES:

* core: Bump Go version to 1.23.7
* core: Bump Go version to 1.23.8

BUG FIXES:

* secrets/openldap: Prevent static role rotation on upgrade when `NextVaultRotation` is nil. Fixes an issue where static roles were unexpectedly rotated after upgrade due to a missing `NextVaultRotation` value. Now sets it to either `LastVaultRotation + RotationPeriod` or `now + RotationPeriod`. [[GH-30265](https://github.com/hashicorp/vault/pull/30265)]
* secrets/transit: fix a panic when rotating on a managed key returns an error [[GH-30214](https://github.com/hashicorp/vault/pull/30214)]

## 1.17.14 Enterprise
### April 04, 2025

SECURITY:

* auth/azure: Update plugin to v0.18.2. Login requires `resource_group_name`, `vm_name`, and `vmss_name` to match token claims [[HCSEC-2025-07](https://discuss.hashicorp.com/t/hcsec-2025-07-vault-s-azure-authentication-method-bound-location-restriction-could-be-bypassed-on-login/74716)].

CHANGES:

* core: Verify that the client IP address extracted from an X-Forwarded-For header is a valid IPv4 or IPv6 address [[GH-29774](https://github.com/hashicorp/vault/pull/29774)]

IMPROVEMENTS:

* core (enterprise): report errors from the underlying seal when getting entropy.

BUG FIXES:

* auth/ldap: Fix a bug that does not properly delete users and groups by first converting their names to lowercase when case senstivity option is off. [[GH-29922](https://github.com/hashicorp/vault/pull/29922)]
* core: Fix Azure authentication for seal/managed keys to work for both federated workload identity and managed user identities.  Fixes regression for federated workload identities. [[GH-29792](https://github.com/hashicorp/vault/pull/29792)]
* core: Fix a bug that prevents certain loggers from writing to a log file. [[GH-29917](https://github.com/hashicorp/vault/pull/29917)]
* export API: Normalize the start_date parameter to the start of the month as is done in the sys/counters API to keep the results returned from both of the API's consistent. [[GH-29562](https://github.com/hashicorp/vault/pull/29562)]
* plugins (enterprise): Fix plugin registration with artifact when a binary for the same plugin is already present in the plugin directory.
* plugins: plugin registration should honor the `plugin_tmpdir` config [[GH-29978](https://github.com/hashicorp/vault/pull/29978)]
* secrets/azure: Upgrade plugin to v0.19.3 which reverts role name changes to no longer be a GUID.
* secrets/database: Fix a bug where a global database plugin reload exits if any of the database connections are not available [[GH-29519](https://github.com/hashicorp/vault/pull/29519)]

## 1.17.13 Enterprise
### March 5, 2025

BUG FIXES:

* agent: Fixed an issue where giving the agent multiple config files could cause the merged config to be incorrect
when `template_config` is set in one of the config files. [[GH-29680](https://github.com/hashicorp/vault/pull/29680)]
* ui: MFA methods now display the namespace path instead of the namespace id. [[GH-29588](https://github.com/hashicorp/vault/pull/29588)]
* core: Fix Azure authentication for seal/managed keys to work for both federated workload identity and managed user identities.  Fixes regression for federated workload identities. [[GH-29792](https://github.com/hashicorp/vault/pull/29792)]

## 1.17.12 Enterprise
### February 25, 2025

SECURITY:

* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.2.0

CHANGES:

* core: Bump Go version to 1.23.6
* raft/autopilot (enterprise): Alongside the CE autopilot update, update raft-autopilot-enterprise library to v0.3.0 and add enterprise-specific regression testing.
* secrets/openldap: Update plugin to v0.13.5

FEATURES:

* **Plugins**: Allow Enterprise plugins to run externally on Vault Enterprise only.

IMPROVEMENTS:

* raft/autopilot: We've updated the autopilot reconciliation logic (by updating the raft-autopilot dependency to v0.3.0) to avoid artificially increasing the quorum in presence of an unhealthy node. Now autopilot will start the reconciliation process by attempting to demote a failed voter node before any promotions, fixing the issue where Vault would initially increase quorum when faced with a failure of a voter node. In certain configurations, especially when using Vault Enterprise Redundancy Zones and losing a voter then a non-voter in quick succession, this would lead to a loss of quorum and cluster failure. [[GH-29306](https://github.com/hashicorp/vault/pull/29306)]
* ui: Application static breadcrumbs should be formatted in title case. [[GH-29206](https://github.com/hashicorp/vault/pull/29206)]

BUG FIXES:

* activity: Show activity records from clients created in deleted namespaces when activity log is queried from admin namespace. [[GH-29432](https://github.com/hashicorp/vault/pull/29432)]
* core/managed-keys (enterprise): Allow mechanism numbers above 32 bits in PKCS#11 managed keys.
* core: Fix bug when if failing to persist the barrier keyring to track encryption counts, the number of outstanding encryptions remains added to the count, overcounting encryptions. [[GH-29506](https://github.com/hashicorp/vault/pull/29506)]
* identity/oidc (enterprise): Fix delays in rotation and invalidation of OIDC keys when there are too many namespaces.
The Cache-Control header returned by the identity/oidc/.well-known/keys endpoint now depends only on the named keys for
the queried namespace. [[GH-29312](https://github.com/hashicorp/vault/pull/29312)]
* secrets-sync (enterprise): Add new parameters for destination configs to specify allowlists for IP's and ports.
* secrets/pki: fixes issue #28749 requiring all chains to be single line of authority. [[GH-29342](https://github.com/hashicorp/vault/pull/29342)]
* ui (enterprise): Fixes token renewal to ensure capability checks are performed in the relevant namespace, resolving 'Not authorized' errors for resources that users have permission to access. [[GH-29416](https://github.com/hashicorp/vault/pull/29416)]
* ui/database: Fixes 'cannot update static username' error when updating static role's rotation period [[GH-29498](https://github.com/hashicorp/vault/pull/29498)]
* ui: Fixes text overflow on Secrets engines and Auth Engines list views for long names & descriptions [[GH-29430](https://github.com/hashicorp/vault/pull/29430)]

## 1.17.11 Enterprise
### January 30, 2025

CHANGES:

* auth/cf: Update plugin to v0.19.1 [[GH-29295](https://github.com/hashicorp/vault/pull/29295)]
* sdk: Updated golang and dependency versions to be consistent across core, API, SDK to address [[GO-2024-3333](https://pkg.go.dev/vuln/GO-2024-3333)] and ensure version consistency [[GH-29422](https://github.com/hashicorp/vault/pull/29422)]

IMPROVEMENTS:

* plugins (enterprise): The Database secrets engine now allows skipping the automatic rotation of static roles during import.
* events (enterprise): Use the `path` event metadata field when authorizing a client's `subscribe` capability for consuming an event, instead of requiring `data_path` to be present in the event metadata.
* ui: Adds navigation for LDAP hierarchical libraries [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]
* ui: Adds params to postgresql database to improve editing a connection in the web browser. [[GH-29200](https://github.com/hashicorp/vault/pull/29200)]

BUG FIXES:

* activity: Include activity records from clients created by deleted or disabled auth mounts in Export API response. [[GH-29376](https://github.com/hashicorp/vault/pull/29376)]
* core: Prevent integer overflows of the barrier key counter on key rotation requests [[GH-29176](https://github.com/hashicorp/vault/pull/29176)]
* database/mssql: Fix a bug where contained databases would silently fail root rotation if a custom root rotation statement was not provided. [[GH-29399](https://github.com/hashicorp/vault/pull/29399)]
* plugins: Fix a bug that causes zombie dbus-daemon processes on certain systems. [[GH-29334](https://github.com/hashicorp/vault/pull/29334)]
* sdk/database: Fix a bug where slow database connections can cause goroutines to be blocked. [[GH-29097](https://github.com/hashicorp/vault/pull/29097)]
* secrets/pki: Fix a bug that prevented the full CA chain to be used when enforcing name constraints. [[GH-29255](https://github.com/hashicorp/vault/pull/29255)]
* sentinel (enterprise): No longer report inaccurate log messages for when failing an advisory policy.
* ui (enterprise): Fixes login to web UI when MFA is enabled for SAML auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui: Fixes login to web UI when MFA is enabled for OIDC (i.e. azure, auth0) and Okta auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui: Fixes navigation for quick actions in LDAP roles' popup menu [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]

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

## 1.16.29 Enterprise
### January 07, 2026

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* core: Bump Go version to 1.24.11
* packaging: Container images are now exported using a compressed OCI image layout.
* packaging: UBI container images are now built on the UBI 10 minimal image.
* storage: Upgrade aerospike client library to v8.

IMPROVEMENTS:

* go: update to golang/x/crypto to v0.45.0 to resolve GHSA-f6x5-jh6r-wrfv, GHSA-j5w8-q4qc-rx2x, GO-2025-4134 and GO-2025-4135.
* secrets-sync (enterprise): Added support for a boolean force_delete flag (default: false). When set to true, this flag allows deletion of a destination even if its associations cannot be unsynced. This option should be used only as a last-resort deletion mechanism, as any secrets already synced to the external provider will remain orphaned and require manual cleanup.

BUG FIXES:

* http: skip JSON limit parsing on cluster listener
* secret-sync (enterprise): Improved unsync error handling by treating cases where the destination no longer exists as successful.
* secrets-sync (enterprise): Corrected a bug where the deletion of the latest KV-V2 secret version caused the associated external secret to be deleted entirely. The sync job now implements a version fallback mechanism to find and sync the highest available active version, ensuring continuity and preventing the unintended deletion of the external secret resource.
* ui/pki: Fix handling of values that contain commas in list fields like `crl_distribution_points`.


## 1.16.28 Enterprise
### November 19, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* core: Bump Go version to 1.24.10
* policies: add VAULT_NEW_PER_ELEMENT_MATCHING_ON_LIST env var to adopt new "contains all" list matching behavior on
allowed_parameters and denied_parameters

IMPROVEMENTS:

* Update github.com/dvsekhvalnov/jose2go to fix security vulnerability CVE-2025-63811.
* auth/ldap: Require non-empty passwords on login command to prevent unauthenticated access to Vault.

BUG FIXES:

* core: resultant-acl now merges segment-wildcard (`+`) paths with existing prefix rules in `glob_paths`, so clients receive a complete view of glob-style permissions. This unblocks UI sidebar navigation checks and namespace access banners.
* ui: Fixes permissions for hiding and showing sidebar navigation items for policies that include special characters: `+`, `*`


## 1.16.27 Enterprise
### October 23, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* auth/aws: fix an issue where a user may be able to bypass authentication to Vault due to incorrect caching of the AWS client
* ui: disable scarf analytics for ui builds

CHANGES:

* core: Bump Go version to 1.24.9.
* http: Evaluate rate limit quotas before checking JSON limits during request handling.

IMPROVEMENTS:

* secrets/database: Add root rotation support for Snowflake database secrets engines using key-pair credentials.

BUG FIXES:

* core (enterprise): Avoid duplicate seal rewrapping, and ensure that cluster secondaries rewrap after a seal migration.

## 1.16.26 Enterprise
### September 24, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* core: Update github.com/hashicorp/go-getter to fix security vulnerability GHSA-wjrx-6529-hcj3.
* core: Update github.com/ulikunitz/xz to fix security vulnerability GHSA-25xm-hr59-7c27.

CHANGES:

* core: Bump Go version to 1.24.7.
* core: Updates post-install script to print updated license information
* database/snowflake: Update plugin to [v0.10.4](https://github.com/hashicorp/vault-plugin-database-snowflake/releases/tag/v0.10.4)
* sdk: Upgrade to go-secure-stdlib/plugincontainer@v0.4.2, which also bumps github.com/docker/docker to v28.3.3+incompatible

IMPROVEMENTS:

* core (enterprise): Updated code to support FIPS 140-3 compliant algorithms.

BUG FIXES:

* auth/cert: Recover from partially populated caches of trusted certificates if one or more certificates fails to load.
* secrets/transit: Fix error when using ed25519 keys that were imported with derivation enabled
* sys/mounts: enable unsetting allowed_response_headers


## 1.16.25 Enterprise
### August 28, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* core: Bump Go version to 1.23.12
* http: Add JSON configurable limits to HTTP handling for JSON payloads: `max_json_depth`, `max_json_string_value_length`, `max_json_object_entry_count`, `max_json_array_element_count`.

BUG FIXES:

* core (enterprise): fix a bug where issuing a token in a namespace used root auth configuration instead of namespace auth configuration
* core/seal: When Seal-HA is enabled, make it an error to persist the barrier
keyring when not all seals are healthy.  This prevents the possibility of
failing to unseal when a different subset of seals are healthy than were
healthy at last write.

## 1.16.24 Enterprise
### August 06, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* auth/ldap: fix MFA/TOTP enforcement bypass when username_as_alias is enabled [[GH-31427](https://github.com/hashicorp/vault/pull/31427),[HCSEC-2025-20](https://discuss.hashicorp.com/t/hcsec-2025-20-vault-ldap-mfa-enforcement-bypass-when-using-username-as-alias/76092)].

BUG FIXES:

* identity/mfa: revert cache entry change from #31217 and document cache entry values [[GH-31421](https://github.com/hashicorp/vault/pull/31421)]

## 1.16.23 Enterprise
### July 25, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* audit: **breaking change** privileged vault operator may execute code on the underlying host (CVE-2025-6000). Vault will not unseal if the only configured file audit device has executable permissions (e.g., 0777, 0755). See recent [breaking change](https://developer.hashicorp.com/vault/docs/updates/important-changes#breaking-changes) docs for more details. [[GH-31211](https://github.com/hashicorp/vault/pull/31211),[HCSEC-2025-14](https://discuss.hashicorp.com/t/hcsec-2025-14-privileged-vault-operator-may-execute-code-on-the-underlying-host/76033)]
* auth/userpass: timing side-channel in vault's userpass auth method (CVE-2025-6011)[HCSEC-2025-15](https://discuss.hashicorp.com/t/hcsec-2025-15-timing-side-channel-in-vault-s-userpass-auth-method/76034)
* core/login: vault userpass and ldap user lockout bypass (CVE-2025-6004). update alias lookahead to respect username case for LDAP and username/password. [[GH-31352](https://github.com/hashicorp/vault/pull/31352),[HCSEC-2025-16](https://discuss.hashicorp.com/t/hcsec-2025-16-vault-userpass-and-ldap-user-lockout-bypass/76035)]
* secrets/totp: vault totp secrets engine code reuse (CVE-2025-6014) [[GH-31246](https://github.com/hashicorp/vault/pull/31246),[HCSEC-2025-17](https://discuss.hashicorp.com/t/hcsec-2025-17-vault-totp-secrets-engine-code-reuse/76036)]
* auth/cert: vault certificate auth method did not validate common name for non-ca certificates (CVE-2025-6037). test non-CA cert equality on login matching instead of individual fields. [[GH-31210](https://github.com/hashicorp/vault/pull/31210),[HCSEC-2025-18](https://discuss.hashicorp.com/t/hcsec-2025-18-vault-certificate-auth-method-did-not-validate-common-name-for-non-ca-certificates/76037)]
* core/mfa: vault login mfa bypass of rate limiting and totp token reuse (CVE-2025-6015) [[GH-31217](https://github.com/hashicorp/vault/pull/31297),[HCSEC-2025-19](https://discuss.hashicorp.com/t/hcsec-2025-19-vault-login-mfa-bypass-of-rate-limiting-and-totp-token-reuse/76038)]

BUG FIXES:

* core/seal (enterprise): Fix a bug that caused the seal rewrap process to abort in the presence of partially sealed entries.
* kmip (enterprise): Fix a panic that can happen when a KMIP client makes a request before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* plugins: Fix panics that can occur when a plugin audits a request or response before the Vault server has finished unsealing. [[GH-31266](https://github.com/hashicorp/vault/pull/31266)]
* secrets-sync (enterprise): Unsyncing secret-key granularity associations will no longer give a misleading error about a failed unsync operation that did indeed succeed.

## 1.16.22 Enterprise
### June 25, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* core: require a nonce when cancelling a rekey operation that was initiated within the last 10 minutes. [[GH-30794](https://github.com/hashicorp/vault/pull/30794)],[[HCSEC-2025-11](https://discuss.hashicorp.com/t/hcsec-2025-11-vault-vulnerable-to-recovery-key-cancellation-denial-of-service/75570)]
* core/identity: vault root namespace operator may elevate privileges (CVE-2025-5999). Fix string contains check in Identity APIs to be case-insensitive. [[GH-31045](https://github.com/hashicorp/vault/pull/31045),[HCSEC-2025-13](https://discuss.hashicorp.com/t/hcsec-2025-13-vault-root-namespace-operator-may-elevate-token-privileges/76032)]

CHANGES:

* auth/azure: Update plugin to v0.17.5
* core: Bump Go version to 1.23.10
* secrets/azure: Update plugin to v0.17.5
* secrets/database: Update vault-plugin-database-snowflake to v0.10.3

BUG FIXES:

* secrets/database: Treat all rotation_schedule values as UTC to ensure consistent behavior. [[GH-30606](https://github.com/hashicorp/vault/pull/30606)]

## 1.16.21 Enterprise
### May 30, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.


CHANGES:

* Update vault-plugin-auth-cf to v0.18.2
* auth/azure: Upgrade plugin to v0.17.4
* database/snowflake: Update plugin to v0.10.2

BUG FIXES:

* plugins (enterprise): Fix an issue where Enterprise plugins can't run on a standby node
when it becomes active because standby nodes don't extract the artifact when the plugin
is registered. Remove extracting from Vault and require the operator to place
the extracted artifact in the plugin directory before registration.

## 1.16.20 Enterprise
### April 30, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* core: vault may expose sensitive information in error logs when processing malformed data with the kv v2 plugin[[GH-30388](https://github.com/hashicorp/vault/pull/30388), [HCSEC-2025-09](https://discuss.hashicorp.com/t/hcsec-2025-09-vault-may-expose-sensitive-information-in-error-logs-when-processing-malformed-data-with-the-kv-v2-plugin/74717)]

BUG FIXES:

* core (enterprise): fix issue with errors being swallowed on failed HSM logins.
* database: Prevent static roles created in versions prior to 1.15.0 from rotating on backend restart. [[GH-30320](https://github.com/hashicorp/vault/pull/30320)]
* database: no longer incorrectly add an "unrecognized parameters" warning for certain SQL database secrets config operations when another warning is returned [[GH-30327](https://github.com/hashicorp/vault/pull/30327)]

## 1.16.19 Enterprise
### April 18, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* core: Bump Go version to 1.23.7
* core: Bump Go version to 1.23.8

BUG FIXES:

* secrets/openldap: Prevent static role rotation on upgrade when `NextVaultRotation` is nil. Fixes an issue where static roles were unexpectedly rotated after upgrade due to a missing `NextVaultRotation` value. Now sets it to either `LastVaultRotation + RotationPeriod` or `now + RotationPeriod`. [[GH-30265](https://github.com/hashicorp/vault/pull/30265)]
* secrets/transit: fix a panic when rotating on a managed key returns an error [[GH-30214](https://github.com/hashicorp/vault/pull/30214)]

## 1.16.18 Enterprise
### April 4, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* auth/azure: Update plugin to v0.17.3. Login requires `resource_group_name`, `vm_name`, and `vmss_name` to match token claims [[HCSEC-2025-07](https://discuss.hashicorp.com/t/hcsec-2025-07-vault-s-azure-authentication-method-bound-location-restriction-could-be-bypassed-on-login/74716)].

IMPROVEMENTS:

* core (enterprise): report errors from the underlying seal when getting entropy.

BUG FIXES:

* auth/ldap: Fix a bug that does not properly delete users and groups by first converting their names to lowercase when case senstivity option is off. [[GH-29922](https://github.com/hashicorp/vault/pull/29922)]
* core: Fix Azure authentication for seal/managed keys to work for both federated workload identity and managed user identities.  Fixes regression for federated workload identities. [[GH-29792](https://github.com/hashicorp/vault/pull/29792)]
* core: Fix a bug that prevents certain loggers from writing to a log file. [[GH-29917](https://github.com/hashicorp/vault/pull/29917)]
* export API: Normalize the start_date parameter to the start of the month as is done in the sys/counters API to keep the results returned from both of the API's consistent. [[GH-29562](https://github.com/hashicorp/vault/pull/29562)]
* plugins (enterprise): Fix plugin registration with artifact when a binary for the same plugin is already present in the plugin directory.
* plugins: plugin registration should honor the `plugin_tmpdir` config [[GH-29978](https://github.com/hashicorp/vault/pull/29978)]
* secrets/azure: Upgrade plugin to v0.17.4 which reverts role name changes to no longer be a GUID.
* secrets/database: Fix a bug where a global database plugin reload exits if any of the database connections are not available [[GH-29519](https://github.com/hashicorp/vault/pull/29519)]

## 1.16.17 Enterprise
### March 5, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

BUG FIXES:

* core: Fix Azure authentication for seal/managed keys to work for both federated workload identity and managed user identities.  Fixes regression for federated workload identities. [[GH-29792](https://github.com/hashicorp/vault/pull/29792)]

## 1.16.16 Enterprise
### February 25, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

SECURITY:

* raft/snapshotagent (enterprise): upgrade raft-snapshotagent to v0.2.0

CHANGES:

* build: Drop `netbsd/386` and `netbsd/arm` builds as downstream modules no longer support them.
* core: Bump Go version to 1.23.6.
* raft/autopilot (enterprise): Alongside the CE autopilot update, update raft-autopilot-enterprise library to v0.3.0 and add enterprise-specific regression testing.
* secrets/openldap: Update plugin to v0.12.4

FEATURES:

* **Plugins**: Allow Enterprise plugins to run externally on Vault Enterprise only.

IMPROVEMENTS:

* raft/autopilot: We've updated the autopilot reconciliation logic (by updating the raft-autopilot dependency to v0.3.0) to avoid artificially increasing the quorum in presence of an unhealthy node. Now autopilot will start the reconciliation process by attempting to demote a failed voter node before any promotions, fixing the issue where Vault would initially increase quorum when faced with a failure of a voter node. In certain configurations, especially when using Vault Enterprise Redundancy Zones and losing a voter then a non-voter in quick succession, this would lead to a loss of quorum and cluster failure. [[GH-29306](https://github.com/hashicorp/vault/pull/29306)]

BUG FIXES:

* activity: Show activity records from clients created in deleted namespaces when activity log is queried from admin namespace. [[GH-29432](https://github.com/hashicorp/vault/pull/29432)]
* core/managed-keys (enterprise): Allow mechanism numbers above 32 bits in PKCS#11 managed keys.
* core: Fix bug when if failing to persist the barrier keyring to track encryption counts, the number of outstanding encryptions remains added to the count, overcounting encryptions. [[GH-29506](https://github.com/hashicorp/vault/pull/29506)]
* secrets-sync (enterprise): Add new parameters for destination configs to specify allowlists for IP's and ports.
* secrets/pki: fixes issue #28749 requiring all chains to be single line of authority. [[GH-29342](https://github.com/hashicorp/vault/pull/29342)]
* ui/database: Fixes 'cannot update static username' error when updating static role's rotation period [[GH-29498](https://github.com/hashicorp/vault/pull/29498)]

## 1.16.15 Enterprise
### January 30, 2025

**Enterprise LTS:** Vault Enterprise 1.16 is a [Long-Term Support (LTS)](https://developer.hashicorp.com/vault/docs/enterprise/lts) release.

CHANGES:

* auth/cf: Update plugin to v0.19.1 [[GH-29295](https://github.com/hashicorp/vault/pull/29295)]
* sdk: Updated golang and dependency versions to be consistent across core, API, SDK to address [[GO-2024-3333](https://pkg.go.dev/vuln/GO-2024-3333)] and ensure version consistency [[GH-29422](https://github.com/hashicorp/vault/pull/29422)]

IMPROVEMENTS:

* plugins (enterprise): The Database secrets engine now allows skipping the automatic rotation of static roles during import.
* events (enterprise): Use the `path` event metadata field when authorizing a client's `subscribe` capability for consuming an event, instead of requiring `data_path` to be present in the event metadata.
* ui: Adds navigation for LDAP hierarchical libraries [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]
* ui: Adds params to postgresql database to improve editing a connection in the web browser. [[GH-29200](https://github.com/hashicorp/vault/pull/29200)]

BUG FIXES:

* activity: Include activity records from clients created by deleted or disabled auth mounts in Export API response. [[GH-29376](https://github.com/hashicorp/vault/pull/29376)]
* core: Prevent integer overflows of the barrier key counter on key rotation requests [[GH-29176](https://github.com/hashicorp/vault/pull/29176)]
* database/mssql: Fix a bug where contained databases would silently fail root rotation if a custom root rotation statement was not provided. [[GH-29399](https://github.com/hashicorp/vault/pull/29399)]
* plugins: Fix a bug that causes zombie dbus-daemon processes on certain systems. [[GH-29334](https://github.com/hashicorp/vault/pull/29334)]
* sdk/database: Fix a bug where slow database connections can cause goroutines to be blocked. [[GH-29097](https://github.com/hashicorp/vault/pull/29097)]
* secrets/pki: Fix a bug that prevented the full CA chain to be used when enforcing name constraints. [[GH-29255](https://github.com/hashicorp/vault/pull/29255)]
* sentinel (enterprise): No longer report inaccurate log messages for when failing an advisory policy.
* ui (enterprise): Fixes login to web UI when MFA is enabled for SAML auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui: Fixes login to web UI when MFA is enabled for OIDC (i.e. azure, auth0) and Okta auth methods [[GH-28873](https://github.com/hashicorp/vault/pull/28873)]
* ui: Fixes navigation for quick actions in LDAP roles' popup menu [[GH-29293](https://github.com/hashicorp/vault/pull/29293)]


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
* **Audit Filtering (enterprise)**: Audit devices support expression-based filter rules (powered by go-bexpr) to determine which entries are written to the audit log. [[GH-24558](https://github.com/hashicorp/vault/pull/24558)]
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
