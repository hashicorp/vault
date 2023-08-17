## 1.13.0
### Unreleased

SECURITY:

* secrets/ssh: removal of the deprecated dynamic keys mode. **When any remaining dynamic key leases expire**, an error stating `secret is unsupported by this backend` will be thrown by the lease manager. [[GH-18874](https://github.com/hashicorp/vault/pull/18874)]

CHANGES:

* auth/alicloud: require the `role` field on login [[GH-19005](https://github.com/hashicorp/vault/pull/19005)]
* auth/approle: Add maximum length of 4096 for approle role_names, as this value results in HMAC calculation [[GH-17768](https://github.com/hashicorp/vault/pull/17768)]
* auth: Returns invalid credentials for ldap, userpass and approle when wrong credentials are provided for existent users.
This will only be used internally for implementing user lockout. [[GH-17104](https://github.com/hashicorp/vault/pull/17104)]
* core: Bump Go version to 1.20.
* core: Vault version has been moved out of sdk and into main vault module.
Plugins using sdk/useragent.String must instead use sdk/useragent.PluginString. [[GH-14229](https://github.com/hashicorp/vault/pull/14229)]
* logging: Removed legacy environment variable for log format ('LOGXI_FORMAT'), should use 'VAULT_LOG_FORMAT' instead [[GH-17822](https://github.com/hashicorp/vault/pull/17822)]
* plugins: Mounts can no longer be pinned to a specific _builtin_ version. Mounts previously pinned to a specific builtin version will now automatically upgrade to the latest builtin version, and may now be overridden if an unversioned plugin of the same name and type is registered. Mounts using plugin versions without `builtin` in their metadata remain unaffected. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]
* plugins: `GET /database/config/:name` endpoint now returns an additional `plugin_version` field in the response data. [[GH-16982](https://github.com/hashicorp/vault/pull/16982)]
* plugins: `GET /sys/auth/:path/tune` and `GET /sys/mounts/:path/tune` endpoints may now return an additional `plugin_version` field in the response data if set. [[GH-17167](https://github.com/hashicorp/vault/pull/17167)]
* plugins: `GET` for `/sys/auth`, `/sys/auth/:path`, `/sys/mounts`, and `/sys/mounts/:path` paths now return additional `plugin_version`, `running_plugin_version` and `running_sha256` fields in the response data for each mount. [[GH-17167](https://github.com/hashicorp/vault/pull/17167)]
* secrets/aws: do not create leases for non-renewable/non-revocable STS credentials to reduce storage calls [[GH-15869](https://github.com/hashicorp/vault/pull/15869)]
* sys/internal/inspect: Turns of this endpoint by default. A SIGHUP can now be used to reload the configs and turns this endpoint on.
* ui: Upgrade Ember to version 4.4.0 [[GH-17086](https://github.com/hashicorp/vault/pull/17086)]

FEATURES:

* **User lockout**: Ignore repeated bad credentials from the same user for a configured period of time. Enabled by default.
* **New PKI UI**: Add beta support for new and improved PKI UI [[GH-18842](https://github.com/hashicorp/vault/pull/18842)]
* **Server UDS Listener**: Adding listener to Vault server to serve http request via unix domain socket [[GH-18227](https://github.com/hashicorp/vault/pull/18227)]
* **Transit managed keys**: The transit secrets engine now supports configuring and using managed keys
* ui: Adds Kubernetes secrets engine [[GH-17893](https://github.com/hashicorp/vault/pull/17893)]

IMPROVEMENTS:

* **Redis ElastiCache DB Engine**: Renamed configuration parameters for disambiguation; old parameters still supported for compatibility. [[GH-18752](https://github.com/hashicorp/vault/pull/18752)]
* Reduced binary size [[GH-17678](https://github.com/hashicorp/vault/pull/17678)]
* agent/config: Allow config directories to be specified with -config, and allow multiple -configs to be supplied. [[GH-18403](https://github.com/hashicorp/vault/pull/18403)]
* agent: Add note in logs when starting Vault Agent indicating if the version differs to the Vault Server. [[GH-18684](https://github.com/hashicorp/vault/pull/18684)]
* agent: Added `token_file` auto-auth configuration to allow using a pre-existing token for Vault Agent. [[GH-18740](https://github.com/hashicorp/vault/pull/18740)]
* agent: Agent listeners can now be to be the `metrics_only` role, serving only metrics, as part of the listener's new top level `role` option. [[GH-18101](https://github.com/hashicorp/vault/pull/18101)]
* agent: Configured Vault Agent listeners now listen without the need for caching to be configured. [[GH-18137](https://github.com/hashicorp/vault/pull/18137)]
* agent: allows some parts of config to be reloaded without requiring a restart. [[GH-18638](https://github.com/hashicorp/vault/pull/18638)]
* agent: fix incorrectly used loop variables in parallel tests and when finalizing seals [[GH-16872](https://github.com/hashicorp/vault/pull/16872)]
* api: Remove dependency on sdk module. [[GH-18962](https://github.com/hashicorp/vault/pull/18962)]
* api: Support VAULT_DISABLE_REDIRECTS environment variable (and --disable-redirects flag) to disable default client behavior and prevent the client following any redirection responses. [[GH-17352](https://github.com/hashicorp/vault/pull/17352)]
* audit: Add `elide_list_responses` option, providing a countermeasure for a common source of oversized audit log entries [[GH-18128](https://github.com/hashicorp/vault/pull/18128)]
* audit: Include stack trace when audit logging recovers from a panic. [[GH-18121](https://github.com/hashicorp/vault/pull/18121)]
* auth/alicloud: upgrades dependencies [[GH-18021](https://github.com/hashicorp/vault/pull/18021)]
* auth/azure: Adds support for authentication with Managed Service Identity (MSI) from a 
Virtual Machine Scale Set (VMSS) in flexible orchestration mode. [[GH-17540](https://github.com/hashicorp/vault/pull/17540)]
* auth/azure: upgrades dependencies [[GH-17857](https://github.com/hashicorp/vault/pull/17857)]
* auth/cert: Add configurable support for validating client certs with OCSP. [[GH-17093](https://github.com/hashicorp/vault/pull/17093)]
* auth/cert: Support listing provisioned CRLs within the mount. [[GH-18043](https://github.com/hashicorp/vault/pull/18043)]
* auth/gcp: Upgrades dependencies [[GH-17858](https://github.com/hashicorp/vault/pull/17858)]
* auth/token (enterprise): Allow batch token creation in perfStandby nodes
* auth: Allow naming login MFA methods and using those names instead of IDs in satisfying MFA requirement for requests.
Make passcode arguments consistent across login MFA method types. [[GH-18610](https://github.com/hashicorp/vault/pull/18610)]
* auth: Provide an IP address of the requests from Vault to a Duo challenge after successful authentication. [[GH-18811](https://github.com/hashicorp/vault/pull/18811)]
* autopilot: Update version to v.0.2.0 to add better support for respecting min quorum [[GH-17848](https://github.com/hashicorp/vault/pull/17848)]
* autopilot: Update version to v.0.2.0 to add better support for respecting min quorum
* cli/kv: improve kv CLI to remove data or custom metadata using kv patch [[GH-18067](https://github.com/hashicorp/vault/pull/18067)]
* cli/pki: Add List-Intermediates functionality to pki client. [[GH-18463](https://github.com/hashicorp/vault/pull/18463)]
* cli/pki: Add health-check subcommand to evaluate the health of a PKI instance. [[GH-17750](https://github.com/hashicorp/vault/pull/17750)]
* cli/pki: Add pki issue command, which creates a CSR, has a vault mount sign it, then reimports it. [[GH-18467](https://github.com/hashicorp/vault/pull/18467)]
* cli: Add support for creating requests to existing non-KVv2 PATCH-capable endpoints. [[GH-17650](https://github.com/hashicorp/vault/pull/17650)]
* cli: Add transit import key helper commands for BYOK to Transit/Transform. [[GH-18887](https://github.com/hashicorp/vault/pull/18887)]
* cli: Support the -format=raw option, to read non-JSON Vault endpoints and original response bodies. [[GH-14945](https://github.com/hashicorp/vault/pull/14945)]
* cli: updated `vault operator rekey` prompts to describe recovery keys when `-target=recovery` [[GH-18892](https://github.com/hashicorp/vault/pull/18892)]
* client/pki: Add a new command verify-sign which checks the relationship between two certificates. [[GH-18437](https://github.com/hashicorp/vault/pull/18437)]
* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* core/identity: Add machine-readable output to body of response upon alias clash during entity merge [[GH-17459](https://github.com/hashicorp/vault/pull/17459)]
* core/server: Added an environment variable to write goroutine stacktraces to a 
temporary file for SIGUSR2 signals. [[GH-17929](https://github.com/hashicorp/vault/pull/17929)]
* core: Add RPCs to read and update userFailedLoginInfo map
* core: Add experiments system and `events.alpha1` experiment. [[GH-18682](https://github.com/hashicorp/vault/pull/18682)]
* core: Add read support to `sys/loggers` and `sys/loggers/:name` endpoints [[GH-17979](https://github.com/hashicorp/vault/pull/17979)]
* core: Add user lockout field to config and configuring this for auth mount using auth tune to prevent brute forcing in auth methods [[GH-17338](https://github.com/hashicorp/vault/pull/17338)]
* core: Add vault.core.locked_users telemetry metric to emit information about total number of locked users. [[GH-18718](https://github.com/hashicorp/vault/pull/18718)]
* core: Added sys/locked-users endpoint to list locked users. Changed api endpoint from
sys/lockedusers/[mount_accessor]/unlock/[alias_identifier] to sys/locked-users/[mount_accessor]/unlock/[alias_identifier]. [[GH-18675](https://github.com/hashicorp/vault/pull/18675)]
* core: Added sys/lockedusers/[mount_accessor]/unlock/[alias_identifier] endpoint to unlock an user
with given mount_accessor and alias_identifier if locked [[GH-18279](https://github.com/hashicorp/vault/pull/18279)]
* core: Added warning to /sys/seal-status and vault status command if potentially dangerous behaviour overrides are being used. [[GH-17855](https://github.com/hashicorp/vault/pull/17855)]
* core: Implemented background thread to update locked user entries every 15 minutes to prevent brute forcing in auth methods. [[GH-18673](https://github.com/hashicorp/vault/pull/18673)]
* core: License location is no longer cache exempt, meaning sys/health will not contribute as greatly to storage load when using consul as a storage backend. [[GH-17265](https://github.com/hashicorp/vault/pull/17265)]
* core: Update protoc from 3.21.5 to 3.21.7 [[GH-17499](https://github.com/hashicorp/vault/pull/17499)]
* core: add `detect_deadlocks` config to optionally detect core state deadlocks [[GH-18604](https://github.com/hashicorp/vault/pull/18604)]
* core: added changes for user lockout workflow. [[GH-17951](https://github.com/hashicorp/vault/pull/17951)]
* core: parallelize backend initialization to improve startup time for large numbers of mounts. [[GH-18244](https://github.com/hashicorp/vault/pull/18244)]
* database/postgres: Support multiline strings for revocation statements. [[GH-18632](https://github.com/hashicorp/vault/pull/18632)]
* database/snowflake: Allow parallel requests to Snowflake [[GH-17593](https://github.com/hashicorp/vault/pull/17593)]
* hcp/connectivity: Add foundational OSS support for opt-in secure communication between self-managed Vault nodes and [HashiCorp Cloud Platform](https://cloud.hashicorp.com) [[GH-18228](https://github.com/hashicorp/vault/pull/18228)]
* hcp/connectivity: Include HCP organization, project, and resource ID in server startup logs [[GH-18315](https://github.com/hashicorp/vault/pull/18315)]
* hcp/connectivity: Only update SCADA session metadata if status changes [[GH-18585](https://github.com/hashicorp/vault/pull/18585)]
* hcp/status: Add cluster-level status information [[GH-18351](https://github.com/hashicorp/vault/pull/18351)]
* hcp/status: Expand node-level status information [[GH-18302](https://github.com/hashicorp/vault/pull/18302)]
* logging: Vault Agent supports logging to a specified file path via environment variable, CLI or config [[GH-17841](https://github.com/hashicorp/vault/pull/17841)]
* logging: Vault agent and server commands support log file and log rotation. [[GH-18031](https://github.com/hashicorp/vault/pull/18031)]
* migration: allow parallelization of key migration for `vault operator migrate` in order to speed up a migration. [[GH-18817](https://github.com/hashicorp/vault/pull/18817)]
* namespaces (enterprise): Add new API, `sys/config/group-policy-application`, to allow group policies to be configurable
to apply to a group in `any` namespace. The default, `within_namespace_hierarchy`, is the current behaviour.
* openapi: Add default values to thing_mount_path parameters [[GH-18935](https://github.com/hashicorp/vault/pull/18935)]
* openapi: Add logic to generate openapi response structures [[GH-18192](https://github.com/hashicorp/vault/pull/18192)]
* openapi: Add openapi response definitions to approle/path_login.go & approle/path_tidy_user_id.go [[GH-18772](https://github.com/hashicorp/vault/pull/18772)]
* openapi: Add openapi response definitions to approle/path_role.go [[GH-18198](https://github.com/hashicorp/vault/pull/18198)]
* openapi: Change gen_openapi.sh to generate schema with generic mount paths [[GH-18934](https://github.com/hashicorp/vault/pull/18934)]
* openapi: Mark request body objects as required [[GH-17909](https://github.com/hashicorp/vault/pull/17909)]
* openapi: add openapi response defintions to /sys/audit endpoints [[GH-18456](https://github.com/hashicorp/vault/pull/18456)]
* openapi: generic_mount_paths: Move implementation fully into server, rather than partially in plugin framework; recognize all 4 singleton mounts (auth/token, cubbyhole, identity, system) rather than just 2; change parameter from `{mountPath}` to `{<type>_mount_path}` [[GH-18663](https://github.com/hashicorp/vault/pull/18663)]
* plugins: Add plugin version information to key plugin lifecycle log lines. [[GH-17430](https://github.com/hashicorp/vault/pull/17430)]
* plugins: Allow selecting builtin plugins by their reported semantic version of the form `vX.Y.Z+builtin` or `vX.Y.Z+builtin.vault`. [[GH-17289](https://github.com/hashicorp/vault/pull/17289)]
* plugins: Let Vault unseal and mount deprecated builtin plugins in a
deactivated state if this is not the first unseal after an upgrade. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* plugins: Mark app-id auth method Removed and remove the plugin code. [[GH-18039](https://github.com/hashicorp/vault/pull/18039)]
* plugins: Mark logical database plugins Removed and remove the plugin code. [[GH-18039](https://github.com/hashicorp/vault/pull/18039)]
* sdk/ldap: Added support for paging when searching for groups using group filters [[GH-17640](https://github.com/hashicorp/vault/pull/17640)]
* sdk: Add response schema validation method framework/FieldData.ValidateStrict and two test helpers (ValidateResponse, ValidateResponseData) [[GH-18635](https://github.com/hashicorp/vault/pull/18635)]
* sdk: Adding FindResponseSchema test helper to assist with response schema validation in tests [[GH-18636](https://github.com/hashicorp/vault/pull/18636)]
* secrets/aws: Update dependencies [[PR-17747](https://github.com/hashicorp/vault/pull/17747)] [[GH-17747](https://github.com/hashicorp/vault/pull/17747)]
* secrets/azure: upgrades dependencies [[GH-17964](https://github.com/hashicorp/vault/pull/17964)]
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* secrets/gcp: Upgrades dependencies [[GH-17871](https://github.com/hashicorp/vault/pull/17871)]
* secrets/kubernetes: Add /check endpoint to determine if environment variables are set [[GH-18](https://github.com/hashicorp/vault-plugin-secrets-kubernetes/pull/18)] [[GH-18587](https://github.com/hashicorp/vault/pull/18587)]
* secrets/kv: new KVv2 mounts and KVv1 mounts without any keys will upgrade synchronously, allowing for instant use [[GH-17406](https://github.com/hashicorp/vault/pull/17406)]
* secrets/pki: Add a new API that returns the serial numbers of revoked certificates on the local cluster [[GH-17779](https://github.com/hashicorp/vault/pull/17779)]
* secrets/pki: Add support to specify signature bits when generating CSRs through intermediate/generate apis [[GH-17388](https://github.com/hashicorp/vault/pull/17388)]
* secrets/pki: Added a new API that allows external actors to craft a CRL through JSON parameters [[GH-18040](https://github.com/hashicorp/vault/pull/18040)]
* secrets/pki: Allow UserID Field (https://www.rfc-editor.org/rfc/rfc1274#section-9.3.1) to be set on Certificates when
allowed by role [[GH-18397](https://github.com/hashicorp/vault/pull/18397)]
* secrets/pki: Allow issuer creation, import to change default issuer via `default_follows_latest_issuer`. [[GH-17824](https://github.com/hashicorp/vault/pull/17824)]
* secrets/pki: Allow templating performance replication cluster- and issuer-specific AIA URLs. [[GH-18199](https://github.com/hashicorp/vault/pull/18199)]
* secrets/pki: Allow tidying of expired issuer certificates. [[GH-17823](https://github.com/hashicorp/vault/pull/17823)]
* secrets/pki: Allow tidying of the legacy ca_bundle, improving startup on post-migrated, seal-wrapped PKI mounts. [[GH-18645](https://github.com/hashicorp/vault/pull/18645)]
* secrets/pki: Respond with written data to `config/auto-tidy`, `config/crl`, and `roles/:role`. [[GH-18222](https://github.com/hashicorp/vault/pull/18222)]
* secrets/pki: Return issuer_id and issuer_name on /issuer/:issuer_ref/json endpoint. [[GH-18482](https://github.com/hashicorp/vault/pull/18482)]
* secrets/pki: Return new fields revocation_time_rfc3339 and issuer_id to existing certificate serial lookup api if it is revoked [[GH-17774](https://github.com/hashicorp/vault/pull/17774)]
* secrets/ssh: Allow removing SSH host keys from the dynamic keys feature. [[GH-18939](https://github.com/hashicorp/vault/pull/18939)]
* secrets/ssh: Evaluate ssh validprincipals user template before splitting [[GH-16622](https://github.com/hashicorp/vault/pull/16622)]
* secrets/transit: Add an optional reference field to batch operation items
which is repeated on batch responses to help more easily correlate inputs with outputs. [[GH-18243](https://github.com/hashicorp/vault/pull/18243)]
* secrets/transit: Add associated_data parameter for additional authenticated data in AEAD ciphers [[GH-17638](https://github.com/hashicorp/vault/pull/17638)]
* secrets/transit: Add support for PKCSv1_5_NoOID RSA signatures [[GH-17636](https://github.com/hashicorp/vault/pull/17636)]
* secrets/transit: Allow configuring whether upsert of keys is allowed. [[GH-18272](https://github.com/hashicorp/vault/pull/18272)]
* storage/raft: Add `retry_join_as_non_voter` config option. [[GH-18030](https://github.com/hashicorp/vault/pull/18030)]
* storage/raft: add additional raft metrics relating to applied index and heartbeating; also ensure OSS standbys emit periodic metrics. [[GH-12166](https://github.com/hashicorp/vault/pull/12166)]
* sys/internal/inspect: Creates an endpoint to look to inspect internal subsystems. [[GH-17789](https://github.com/hashicorp/vault/pull/17789)]
* sys/internal/inspect: Creates an endpoint to look to inspect internal subsystems.
* ui: Add algorithm-signer as a SSH Secrets Engine UI field [[GH-10299](https://github.com/hashicorp/vault/pull/10299)]
* ui: Add inline policy creation when creating an identity entity or group [[GH-17749](https://github.com/hashicorp/vault/pull/17749)]
* ui: Added JWT authentication warning message about blocked pop-up windows and web browser settings. [[GH-18787](https://github.com/hashicorp/vault/pull/18787)]
* ui: Enable typescript for future development [[GH-17927](https://github.com/hashicorp/vault/pull/17927)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]
* ui: consolidate all <a> tag usage [[GH-17866](https://github.com/hashicorp/vault/pull/17866)]
* ui: mfa: use proper request id generation [[GH-17835](https://github.com/hashicorp/vault/pull/17835)]
* ui: update DocLink component to use new host url: developer.hashicorp.com [[GH-18374](https://github.com/hashicorp/vault/pull/18374)]
* ui: update TTL picker for consistency [[GH-18114](https://github.com/hashicorp/vault/pull/18114)]
* ui: use the combined activity log (partial + historic) API for client count dashboard and remove use of monthly endpoint [[GH-17575](https://github.com/hashicorp/vault/pull/17575)]
* vault/diagnose: Upgrade `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk`, `go.opentelemetry.io/otel/trace` to v1.11.2 [[GH-18589](https://github.com/hashicorp/vault/pull/18589)]

BUG FIXES:

* api: Remove timeout logic from ReadRaw functions and add ReadRawWithContext [[GH-18708](https://github.com/hashicorp/vault/pull/18708)]
* auth/alicloud: fix regression in vault login command that caused login to fail [[GH-19005](https://github.com/hashicorp/vault/pull/19005)]
* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* auth/cert: Address a race condition accessing the loaded crls without a lock [[GH-18945](https://github.com/hashicorp/vault/pull/18945)]
* auth/kubernetes: Ensure a consistent TLS configuration for all k8s API requests [[#173](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/173)] [[GH-18716](https://github.com/hashicorp/vault/pull/18716)]
* auth/okta: fix a panic for AuthRenew in Okta [[GH-18011](https://github.com/hashicorp/vault/pull/18011)]
* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* cli: Fix issue preventing kv commands from executing properly when the mount path provided by `-mount` flag and secret key path are the same. [[GH-17679](https://github.com/hashicorp/vault/pull/17679)]
* cli: Fix vault read handling to return raw data as secret.Data when there is no top-level data object from api response. [[GH-17913](https://github.com/hashicorp/vault/pull/17913)]
* cli: Remove empty table heading for `vault secrets list -detailed` output. [[GH-17577](https://github.com/hashicorp/vault/pull/17577)]
* command/namespace: Fix vault cli namespace patch examples in help text. [[GH-18143](https://github.com/hashicorp/vault/pull/18143)]
* core (enterprise): Fix missing quotation mark in error message
* core (enterprise): Fix panic that could occur with SSCT alongside invoking external plugins for revocation.
* core (enterprise): Supported storage check in `vault server` command will no longer prevent startup. Instead, a warning will be logged if configured to use storage backend other than `raft` or `consul`.
* core/activity: add namespace breakdown for new clients when date range spans multiple months, including the current month. [[GH-18766](https://github.com/hashicorp/vault/pull/18766)]
* core/activity: de-duplicate namespaces when historical and current month data are mixed [[GH-18452](https://github.com/hashicorp/vault/pull/18452)]
* core/activity: fix the end_date returned from the activity log endpoint when partial counts are computed [[GH-17856](https://github.com/hashicorp/vault/pull/17856)]
* core/activity: include mount counts when de-duplicating current and historical month data [[GH-18598](https://github.com/hashicorp/vault/pull/18598)]
* core/activity: report mount paths (rather than mount accessors) in current month activity log counts and include deleted mount paths in precomputed queries. [[GH-18916](https://github.com/hashicorp/vault/pull/18916)]
* core/activity: return partial month counts when querying a historical date range and no historical data exists. [[GH-17935](https://github.com/hashicorp/vault/pull/17935)]
* core/auth: Return a 403 instead of a 500 for wrapping requests when token is not provided [[GH-18859](https://github.com/hashicorp/vault/pull/18859)]
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): Switch to using hash length as PSS Salt length within the test/sign api for better PKCS#11 compatibility
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core/seal: Fix regression handling of the key_id parameter in seal configuration HCL. [[GH-17612](https://github.com/hashicorp/vault/pull/17612)]
* core: Fix panic caused in Vault Agent when rendering certificate templates [[GH-17419](https://github.com/hashicorp/vault/pull/17419)]
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: Fix spurious `permission denied` for all HelpOperations on sudo-protected paths [[GH-18568](https://github.com/hashicorp/vault/pull/18568)]
* core: Fix vault operator init command to show the right curl string with -output-curl-string and right policy hcl with -output-policy [[GH-17514](https://github.com/hashicorp/vault/pull/17514)]
* core: Fixes spurious warnings being emitted relating to "unknown or unsupported fields" for JSON config [[GH-17660](https://github.com/hashicorp/vault/pull/17660)]
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* core: Refactor lock grabbing code to simplify stateLock deadlock investigations [[GH-17187](https://github.com/hashicorp/vault/pull/17187)]
* core: fix GPG encryption to support subkeys. [[GH-16224](https://github.com/hashicorp/vault/pull/16224)]
* core: fix a start up race condition where performance standbys could go into a 
mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: fix race when using SystemView.ReplicationState outside of a request context [[GH-17186](https://github.com/hashicorp/vault/pull/17186)]
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* core: trying to unseal with the wrong key now returns HTTP 400 [[GH-17836](https://github.com/hashicorp/vault/pull/17836)]
* credential/cert: adds error message if no tls connection is found during the AliasLookahead operation [[GH-17904](https://github.com/hashicorp/vault/pull/17904)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* expiration: Prevent panics on perf standbys when an irrevocable lease gets deleted. [[GH-18401](https://github.com/hashicorp/vault/pull/18401)]
* licensing (enterprise): update autoloaded license cache after reload
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* openapi: Fix many incorrect details in generated API spec, by using better techniques to parse path regexps [[GH-18554](https://github.com/hashicorp/vault/pull/18554)]
* openapi: fix gen_openapi.sh script to correctly load vault plugins [[GH-17752](https://github.com/hashicorp/vault/pull/17752)]
* plugins/kv: KV v2 returns 404 instead of 500 for request paths that incorrectly include a trailing slash. [[GH-17339](https://github.com/hashicorp/vault/pull/17339)]
* plugins: Allow running external plugins which override deprecated builtins. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* plugins: Corrected the path to check permissions on when the registered plugin name does not match the plugin binary's filename. [[GH-17340](https://github.com/hashicorp/vault/pull/17340)]
* plugins: Listing all plugins while audit logging is enabled will no longer result in an internal server error. [[GH-18173](https://github.com/hashicorp/vault/pull/18173)]
* plugins: Only report deprecation status for builtin plugins. [[GH-17816](https://github.com/hashicorp/vault/pull/17816)]
* plugins: Skip loading but still mount data associated with missing plugins on unseal. [[GH-18189](https://github.com/hashicorp/vault/pull/18189)]
* plugins: Vault upgrades will no longer fail if a mount has been created using an explicit builtin plugin version. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]
* sdk: Don't panic if system view or storage methods called during plugin setup. [[GH-18210](https://github.com/hashicorp/vault/pull/18210)]
* secret/pki: fix bug with initial legacy bundle migration (from < 1.11 into 1.11+) and missing issuers from ca_chain [[GH-17772](https://github.com/hashicorp/vault/pull/17772)]
* secrets/pki: Address nil panic when an empty POST request is sent to the OCSP handler [[GH-18184](https://github.com/hashicorp/vault/pull/18184)]
* secrets/pki: Allow patching issuer to set an empty issuer name. [[GH-18466](https://github.com/hashicorp/vault/pull/18466)]
* secrets/pki: Do not read revoked certificates from backend when CRL is disabled [[GH-17385](https://github.com/hashicorp/vault/pull/17385)]
* secrets/pki: Fix upgrade of missing expiry, delta_rebuild_interval by setting them to the default. [[GH-17693](https://github.com/hashicorp/vault/pull/17693)]
* secrets/pki: Fixes duplicate otherName in certificates created by the sign-verbatim endpoint. [[GH-16700](https://github.com/hashicorp/vault/pull/16700)]
* secrets/pki: OCSP GET request parameter was not being URL unescaped before processing. [[GH-18938](https://github.com/hashicorp/vault/pull/18938)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* secrets/pki: consistently use UTC for CA's notAfter exceeded error message [[GH-18984](https://github.com/hashicorp/vault/pull/18984)]
* secrets/pki: fix race between tidy's cert counting and tidy status reporting. [[GH-18899](https://github.com/hashicorp/vault/pull/18899)]
* secrets/transit: Do not warn about unrecognized parameter 'batch_input' [[GH-18299](https://github.com/hashicorp/vault/pull/18299)]
* secrets/transit: Honor `partial_success_response_code` on decryption failures. [[GH-18310](https://github.com/hashicorp/vault/pull/18310)]
* storage/raft (enterprise): An already joined node can rejoin by wiping storage
and re-issueing a join request, but in doing so could transiently become a
non-voter.  In some scenarios this resulted in loss of quorum. [[GH-18263](https://github.com/hashicorp/vault/pull/18263)]
* storage/raft: Don't panic on unknown raft ops [[GH-17732](https://github.com/hashicorp/vault/pull/17732)]
* storage/raft: Fix race with follower heartbeat tracker during teardown. [[GH-18704](https://github.com/hashicorp/vault/pull/18704)]
* ui/keymgmt: Sets the defaultValue for type when creating a key. [[GH-17407](https://github.com/hashicorp/vault/pull/17407)]
* ui: Fixes issue with not being able to download raft snapshot via service worker [[GH-17769](https://github.com/hashicorp/vault/pull/17769)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]
* ui: Remove default value of 30 to TtlPicker2 if no value is passed in. [[GH-17376](https://github.com/hashicorp/vault/pull/17376)]
* ui: cleanup unsaved auth method ember data record when navigating away from mount backend form [[GH-18651](https://github.com/hashicorp/vault/pull/18651)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]
                                                                       
## 1.12.3
### February 6, 2023
                                                                       
CHANGES:

* core: Bump Go version to 1.19.4.

IMPROVEMENTS:

* audit: Include stack trace when audit logging recovers from a panic. [[GH-18121](https://github.com/hashicorp/vault/pull/18121)]
* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* core: Add read support to `sys/loggers` and `sys/loggers/:name` endpoints [[GH-17979](https://github.com/hashicorp/vault/pull/17979)]
* plugins: Let Vault unseal and mount deprecated builtin plugins in a
deactivated state if this is not the first unseal after an upgrade. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* secrets/kv: new KVv2 mounts and KVv1 mounts without any keys will upgrade synchronously, allowing for instant use [[GH-17406](https://github.com/hashicorp/vault/pull/17406)]
* storage/raft: add additional raft metrics relating to applied index and heartbeating; also ensure OSS standbys emit periodic metrics. [[GH-12166](https://github.com/hashicorp/vault/pull/12166)]
* ui: Added JWT authentication warning message about blocked pop-up windows and web browser settings. [[GH-18787](https://github.com/hashicorp/vault/pull/18787)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]

BUG FIXES:

* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* auth/cert: Address a race condition accessing the loaded crls without a lock [[GH-18945](https://github.com/hashicorp/vault/pull/18945)]
* auth/kubernetes: Ensure a consistent TLS configuration for all k8s API requests [[#173](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/173)] [[GH-18716](https://github.com/hashicorp/vault/pull/18716)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* command/namespace: Fix vault cli namespace patch examples in help text. [[GH-18143](https://github.com/hashicorp/vault/pull/18143)]
* core (enterprise): Fix a race condition resulting in login errors to PKCS#11 modules under high concurrency.
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core/seal: Fix regression handling of the key_id parameter in seal configuration HCL. [[GH-17612](https://github.com/hashicorp/vault/pull/17612)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* expiration: Prevent panics on perf standbys when an irrevocable release gets deleted. [[GH-18401](https://github.com/hashicorp/vault/pull/18401)]
* kmip (enterprise): Fix Destroy operation response that omitted Unique Identifier on some batched responses.
* kmip (enterprise): Fix Locate operation response incompatibility with clients using KMIP versions prior to 1.3.
* kmip (enterprise): Fix Query operation response that omitted streaming capability and supported profiles.
* licensing (enterprise): update autoloaded license cache after reload
* plugins: Allow running external plugins which override deprecated builtins. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* plugins: Listing all plugins while audit logging is enabled will no longer result in an internal server error. [[GH-18173](https://github.com/hashicorp/vault/pull/18173)]
* plugins: Skip loading but still mount data associated with missing plugins on unseal. [[GH-18189](https://github.com/hashicorp/vault/pull/18189)]
* sdk: Don't panic if system view or storage methods called during plugin setup. [[GH-18210](https://github.com/hashicorp/vault/pull/18210)]
* secrets/pki: Address nil panic when an empty POST request is sent to the OCSP handler [[GH-18184](https://github.com/hashicorp/vault/pull/18184)]
* secrets/pki: Allow patching issuer to set an empty issuer name. [[GH-18466](https://github.com/hashicorp/vault/pull/18466)]
* secrets/pki: OCSP GET request parameter was not being URL unescaped before processing. [[GH-18938](https://github.com/hashicorp/vault/pull/18938)]
* secrets/pki: fix race between tidy's cert counting and tidy status reporting. [[GH-18899](https://github.com/hashicorp/vault/pull/18899)]
* secrets/transit: Do not warn about unrecognized parameter 'batch_input' [[GH-18299](https://github.com/hashicorp/vault/pull/18299)]
* secrets/transit: Honor `partial_success_response_code` on decryption failures. [[GH-18310](https://github.com/hashicorp/vault/pull/18310)]
* storage/raft (enterprise): An already joined node can rejoin by wiping storage
and re-issueing a join request, but in doing so could transiently become a
non-voter.  In some scenarios this resulted in loss of quorum. [[GH-18263](https://github.com/hashicorp/vault/pull/18263)]
* storage/raft: Don't panic on unknown raft ops [[GH-17732](https://github.com/hashicorp/vault/pull/17732)]
* ui: cleanup unsaved auth method ember data record when navigating away from mount backend form [[GH-18651](https://github.com/hashicorp/vault/pull/18651)]
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]                                                                                                                                              
## 1.12.2
### November 30, 2022

CHANGES:

* core: Bump Go version to 1.19.3.
* plugins: Mounts can no longer be pinned to a specific _builtin_ version. Mounts previously pinned to a specific builtin version will now automatically upgrade to the latest builtin version, and may now be overridden if an unversioned plugin of the same name and type is registered. Mounts using plugin versions without `builtin` in their metadata remain unaffected. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]

IMPROVEMENTS:

* secrets/pki: Allow issuer creation, import to change default issuer via `default_follows_latest_issuer`. [[GH-17824](https://github.com/hashicorp/vault/pull/17824)]
* storage/raft: Add `retry_join_as_non_voter` config option. [[GH-18030](https://github.com/hashicorp/vault/pull/18030)]

BUG FIXES:

* auth/okta: fix a panic for AuthRenew in Okta [[GH-18011](https://github.com/hashicorp/vault/pull/18011)]
* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* cli: Fix issue preventing kv commands from executing properly when the mount path provided by `-mount` flag and secret key path are the same. [[GH-17679](https://github.com/hashicorp/vault/pull/17679)]
* core (enterprise): Supported storage check in `vault server` command will no longer prevent startup. Instead, a warning will be logged if configured to use storage backend other than `raft` or `consul`.
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: fix a start up race condition where performance standbys could go into a
  mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* plugins: Only report deprecation status for builtin plugins. [[GH-17816](https://github.com/hashicorp/vault/pull/17816)]
* plugins: Vault upgrades will no longer fail if a mount has been created using an explicit builtin plugin version. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]
* secret/pki: fix bug with initial legacy bundle migration (from < 1.11 into 1.11+) and missing issuers from ca_chain [[GH-17772](https://github.com/hashicorp/vault/pull/17772)]
* secrets/azure: add WAL to clean up role assignments if errors occur [[GH-18086](https://github.com/hashicorp/vault/pull/18086)]
* secrets/gcp: Fixes duplicate service account key for rotate root on standby or secondary [[GH-18111](https://github.com/hashicorp/vault/pull/18111)]
* secrets/pki: Fix upgrade of missing expiry, delta_rebuild_interval by setting them to the default. [[GH-17693](https://github.com/hashicorp/vault/pull/17693)]
* ui: Fixes issue with not being able to download raft snapshot via service worker [[GH-17769](https://github.com/hashicorp/vault/pull/17769)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]

## 1.12.1
### November 2, 2022

IMPROVEMENTS:

* api: Support VAULT_DISABLE_REDIRECTS environment variable (and --disable-redirects flag) to disable default client behavior and prevent the client following any redirection responses. [[GH-17352](https://github.com/hashicorp/vault/pull/17352)]
* database/snowflake: Allow parallel requests to Snowflake [[GH-17593](https://github.com/hashicorp/vault/pull/17593)]
* plugins: Add plugin version information to key plugin lifecycle log lines. [[GH-17430](https://github.com/hashicorp/vault/pull/17430)]
* sdk/ldap: Added support for paging when searching for groups using group filters [[GH-17640](https://github.com/hashicorp/vault/pull/17640)]

BUG FIXES:

* cli: Remove empty table heading for `vault secrets list -detailed` output. [[GH-17577](https://github.com/hashicorp/vault/pull/17577)]
* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): Switch to using hash length as PSS Salt length within the test/sign api for better PKCS#11 compatibility
* core: Fix panic caused in Vault Agent when rendering certificate templates [[GH-17419](https://github.com/hashicorp/vault/pull/17419)]
* core: Fixes spurious warnings being emitted relating to "unknown or unsupported fields" for JSON config [[GH-17660](https://github.com/hashicorp/vault/pull/17660)]
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* kmip (enterprise): Fix a problem in the handling of attributes that caused Import operations to fail.
* kmip (enterprise): Fix selection of Cryptographic Parameters for Encrypt/Decrypt operations.
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]

## 1.12.0
### October 13, 2022

CHANGES:

* api: Exclusively use `GET /sys/plugins/catalog` endpoint for listing plugins, and add `details` field to list responses. [[GH-17347](https://github.com/hashicorp/vault/pull/17347)]
* auth: `GET /sys/auth/:name` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* auth: `GET /sys/auth` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* auth: `POST /sys/auth/:type` endpoint response contains a warning for `Deprecated` auth methods. [[GH-17058](https://github.com/hashicorp/vault/pull/17058)]
* auth: `auth enable` returns an error and `POST /sys/auth/:type` endpoint reports an error for `Pending Removal` auth methods. [[GH-17005](https://github.com/hashicorp/vault/pull/17005)]
* core/entities: Fixed stranding of aliases upon entity merge, and require explicit selection of which aliases should be kept when some must be deleted [[GH-16539](https://github.com/hashicorp/vault/pull/16539)]
* core: Bump Go version to 1.19.2.
* core: Validate input parameters for vault operator init command. Vault 1.12 CLI version is needed to run operator init now. [[GH-16379](https://github.com/hashicorp/vault/pull/16379)]
* identity: a request to `/identity/group` that includes `member_group_ids` that contains a cycle will now be responded to with a 400 rather than 500 [[GH-15912](https://github.com/hashicorp/vault/pull/15912)]
* licensing (enterprise): Terminated licenses will no longer result in shutdown. Instead, upgrades will not be allowed if the license expiration time is before the build date of the binary.
* plugins: Add plugin version to auth register, list, and mount table [[GH-16856](https://github.com/hashicorp/vault/pull/16856)]
* plugins: `GET /sys/plugins/catalog/:type/:name` endpoint contains deprecation status for builtin plugins. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* plugins: `GET /sys/plugins/catalog/:type/:name` endpoint now returns an additional `version` field in the response data. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* plugins: `GET /sys/plugins/catalog/` endpoint contains deprecation status in `detailed` list. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* plugins: `GET /sys/plugins/catalog` endpoint now returns an additional `detailed` field in the response data with a list of additional plugin metadata. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* plugins: `plugin info` displays deprecation status for builtin plugins. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* plugins: `plugin list` now accepts a `-detailed` flag, which display deprecation status and version info. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* secrets/azure: Removed deprecated AAD graph API support from the secrets engine. [[GH-17180](https://github.com/hashicorp/vault/pull/17180)]
* secrets: All database-specific (standalone DB) secrets engines are now marked `Pending Removal`. [[GH-17038](https://github.com/hashicorp/vault/pull/17038)]
* secrets: `GET /sys/mounts/:name` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* secrets: `GET /sys/mounts` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* secrets: `POST /sys/mounts/:type` endpoint response contains a warning for `Deprecated` secrets engines. [[GH-17058](https://github.com/hashicorp/vault/pull/17058)]
* secrets: `secrets enable` returns an error and `POST /sys/mount/:type` endpoint reports an error for `Pending Removal` secrets engines. [[GH-17005](https://github.com/hashicorp/vault/pull/17005)]

FEATURES:

* **GCP Cloud KMS support for managed keys**: Managed keys now support using GCP Cloud KMS keys
* **LDAP Secrets Engine**: Adds the `ldap` secrets engine with service account check-out functionality for all supported schemas. [[GH-17152](https://github.com/hashicorp/vault/pull/17152)]
* **OCSP Responder**: PKI mounts now have an OCSP responder that implements a subset of RFC6960, answering single serial number OCSP requests for a specific cluster's revoked certificates in a mount. [[GH-16723](https://github.com/hashicorp/vault/pull/16723)]
* **Redis DB Engine**: Adding the new Redis database engine that supports the generation of static and dynamic user roles and root credential rotation on a stand alone Redis server. [[GH-17070](https://github.com/hashicorp/vault/pull/17070)]
* **Redis ElastiCache DB Plugin**: Added Redis ElastiCache as a built-in plugin. [[GH-17075](https://github.com/hashicorp/vault/pull/17075)]
* **Secrets/auth plugin multiplexing**: manage multiple plugin configurations with a single plugin process [[GH-14946](https://github.com/hashicorp/vault/pull/14946)]
* **Transform Key Import (BYOK)**: The transform secrets engine now supports importing keys for tokenization and FPE transformations
* HCP (enterprise): Adding foundational support for self-managed vault nodes to securely communicate with [HashiCorp Cloud Platform](https://cloud.hashicorp.com) as an opt-in feature
* ui: UI support for Okta Number Challenge. [[GH-15998](https://github.com/hashicorp/vault/pull/15998)]
* **Plugin Versioning**: Vault supports registering, managing, and running plugins with semantic versions specified.

IMPROVEMENTS:

* core/managed-keys (enterprise): Allow operators to specify PSS signatures and/or hash algorithm for the test/sign api
* activity (enterprise): Added new clients unit tests to test accuracy of estimates
* agent/auto-auth: Add `exit_on_err` which when set to true, will cause Agent to exit if any errors are encountered during authentication. [[GH-17091](https://github.com/hashicorp/vault/pull/17091)]
* agent: Added `disable_idle_connections` configuration to disable leaving idle connections open in auto-auth, caching and templating. [[GH-15986](https://github.com/hashicorp/vault/pull/15986)]
* agent: Added `disable_keep_alives` configuration to disable keep alives in auto-auth, caching and templating. [[GH-16479](https://github.com/hashicorp/vault/pull/16479)]
* agent: JWT auto auth now supports a `remove_jwt_after_reading` config option which defaults to true. [[GH-11969](https://github.com/hashicorp/vault/pull/11969)]
* agent: Send notifications to systemd on start and stop. [[GH-9802](https://github.com/hashicorp/vault/pull/9802)]
* api/mfa: Add namespace path to the MFA read/list endpoint [[GH-16911](https://github.com/hashicorp/vault/pull/16911)]
* api: Add a sentinel error for missing KV secrets [[GH-16699](https://github.com/hashicorp/vault/pull/16699)]
* auth/alicloud: Enables AliCloud roles to be compatible with Vault's role based quotas. [[GH-17251](https://github.com/hashicorp/vault/pull/17251)]
* auth/approle: SecretIDs can now be generated with an per-request specified TTL and num_uses.
When either the ttl and num_uses fields are not specified, the role's configuration is used. [[GH-14474](https://github.com/hashicorp/vault/pull/14474)]
* auth/aws: PKCS7 signatures will now use SHA256 by default in prep for Go 1.18 [[GH-16455](https://github.com/hashicorp/vault/pull/16455)]
* auth/azure: Enables Azure roles to be compatible with Vault's role based quotas. [[GH-17194](https://github.com/hashicorp/vault/pull/17194)]
* auth/cert: Add metadata to identity-alias [[GH-14751](https://github.com/hashicorp/vault/pull/14751)]
* auth/cert: Operators can now specify a CRL distribution point URL, in which case the cert auth engine will fetch and use the CRL from that location rather than needing to push CRLs directly to auth/cert. [[GH-17136](https://github.com/hashicorp/vault/pull/17136)]
* auth/cf: Enables CF roles to be compatible with Vault's role based quotas. [[GH-17196](https://github.com/hashicorp/vault/pull/17196)]
* auth/gcp: Add support for GCE regional instance groups [[GH-16435](https://github.com/hashicorp/vault/pull/16435)]
* auth/gcp: Updates dependencies: `google.golang.org/api@v0.83.0`, `github.com/hashicorp/go-gcp-common@v0.8.0`. [[GH-17160](https://github.com/hashicorp/vault/pull/17160)]
* auth/jwt: Adds support for Microsoft US Gov L4 to the Azure provider for groups fetching. [[GH-16525](https://github.com/hashicorp/vault/pull/16525)]
* auth/jwt: Improves detection of Windows Subsystem for Linux (WSL) for CLI-based logins. [[GH-16525](https://github.com/hashicorp/vault/pull/16525)]
* auth/kerberos: add `add_group_aliases` config to include LDAP groups in Vault group aliases [[GH-16890](https://github.com/hashicorp/vault/pull/16890)]
* auth/kerberos: add `remove_instance_name` parameter to the login CLI and the Kerberos config in Vault. This removes any instance names found in the keytab service principal name. [[GH-16594](https://github.com/hashicorp/vault/pull/16594)]
* auth/kubernetes: Role resolution for K8S Auth [[GH-156](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/156)] [[GH-17161](https://github.com/hashicorp/vault/pull/17161)]
* auth/oci: Add support for role resolution. [[GH-17212](https://github.com/hashicorp/vault/pull/17212)]
* auth/oidc: Adds support for group membership parsing when using SecureAuth as an OIDC provider. [[GH-16274](https://github.com/hashicorp/vault/pull/16274)]
* cli: CLI commands will print a warning if flags will be ignored because they are passed after positional arguments. [[GH-16441](https://github.com/hashicorp/vault/pull/16441)]
* cli: `auth` and `secrets` list `-detailed` commands now show Deprecation Status for builtin plugins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* cli: `vault plugin list` now has a `details` field in JSON format, and version and type information in table format. [[GH-17347](https://github.com/hashicorp/vault/pull/17347)]
* command/audit: Improve missing type error message [[GH-16409](https://github.com/hashicorp/vault/pull/16409)]
* command/server: add `-dev-tls` and `-dev-tls-cert-dir` subcommands to create a Vault dev server with generated certificates and private key. [[GH-16421](https://github.com/hashicorp/vault/pull/16421)]
* command: Fix shell completion for KV v2 mounts [[GH-16553](https://github.com/hashicorp/vault/pull/16553)]
* core (enterprise): Add HTTP PATCH support for namespaces with an associated `namespace patch` CLI command
* core (enterprise): Add check to `vault server` command to ensure configured storage backend is supported.
* core (enterprise): Add custom metadata support for namespaces
* core/activity: generate hyperloglogs containing clientIds for each month during precomputation [[GH-16146](https://github.com/hashicorp/vault/pull/16146)]
* core/activity: refactor activity log api to reuse partial api functions in activity endpoint when current month is specified [[GH-16162](https://github.com/hashicorp/vault/pull/16162)]
* core/activity: use monthly hyperloglogs to calculate new clients approximation for current month [[GH-16184](https://github.com/hashicorp/vault/pull/16184)]
* core/quotas (enterprise): Added ability to add path suffixes for lease-count resource quotas
* core/quotas (enterprise): Added ability to add role information for lease-count resource quotas, to limit login requests on auth mounts made using that role
* core/quotas: Added ability to add path suffixes for rate-limit resource quotas [[GH-15989](https://github.com/hashicorp/vault/pull/15989)]
* core/quotas: Added ability to add role information for rate-limit resource quotas, to limit login requests on auth mounts made using that role [[GH-16115](https://github.com/hashicorp/vault/pull/16115)]
* core: Activity log goroutine management improvements to allow tests to be more deterministic. [[GH-17028](https://github.com/hashicorp/vault/pull/17028)]
* core: Add `sys/loggers` and `sys/loggers/:name` endpoints to provide ability to modify logging verbosity [[GH-16111](https://github.com/hashicorp/vault/pull/16111)]
* core: Handle and log deprecated builtin mounts. Introduces `VAULT_ALLOW_PENDING_REMOVAL_MOUNTS` to override shutdown and error when attempting to mount `Pending Removal` builtin plugins. [[GH-17005](https://github.com/hashicorp/vault/pull/17005)]
* core: Limit activity log client count usage by namespaces [[GH-16000](https://github.com/hashicorp/vault/pull/16000)]
* core: Upgrade github.com/hashicorp/raft [[GH-16609](https://github.com/hashicorp/vault/pull/16609)]
* core: remove gox [[GH-16353](https://github.com/hashicorp/vault/pull/16353)]
* docs: Clarify the behaviour of local mounts in the context of DR replication [[GH-16218](https://github.com/hashicorp/vault/pull/16218)]
* identity/oidc: Adds support for detailed listing of clients and providers. [[GH-16567](https://github.com/hashicorp/vault/pull/16567)]
* identity/oidc: Adds the `client_secret_post` token endpoint authentication method. [[GH-16598](https://github.com/hashicorp/vault/pull/16598)]
* identity/oidc: allows filtering the list providers response by an allowed_client_id [[GH-16181](https://github.com/hashicorp/vault/pull/16181)]
* identity: Prevent possibility of data races on entity creation. [[GH-16487](https://github.com/hashicorp/vault/pull/16487)]
* physical/postgresql: pass context to queries to propagate timeouts and cancellations on requests. [[GH-15866](https://github.com/hashicorp/vault/pull/15866)]
* plugins/multiplexing: Added multiplexing support to database plugins if run as external plugins [[GH-16995](https://github.com/hashicorp/vault/pull/16995)]
* plugins: Add Deprecation Status method to builtinregistry. [[GH-16846](https://github.com/hashicorp/vault/pull/16846)]
* plugins: Added environment variable flag to opt-out specific plugins from multiplexing [[GH-16972](https://github.com/hashicorp/vault/pull/16972)]
* plugins: Adding version to plugin GRPC interface [[GH-17088](https://github.com/hashicorp/vault/pull/17088)]
* plugins: Plugin catalog supports registering and managing plugins with semantic version information. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* replication (enterprise): Fix race in merkle sync that can prevent streaming by returning key value matching provided hash if found in log shipper buffer.
* secret/nomad: allow reading CA and client auth certificate from /nomad/config/access [[GH-15809](https://github.com/hashicorp/vault/pull/15809)]
* secret/pki: Add RSA PSS signature support for issuing certificates, signing CRLs [[GH-16519](https://github.com/hashicorp/vault/pull/16519)]
* secret/pki: Add signature_bits to sign-intermediate, sign-verbatim endpoints [[GH-16124](https://github.com/hashicorp/vault/pull/16124)]
* secret/pki: Allow issuing certificates with non-domain, non-email Common Names from roles, sign-verbatim, and as issuers (`cn_validations`). [[GH-15996](https://github.com/hashicorp/vault/pull/15996)]
* secret/pki: Allow specifying SKID for cross-signed issuance from older Vault versions. [[GH-16494](https://github.com/hashicorp/vault/pull/16494)]
* secret/transit: Allow importing Ed25519 keys from PKCS#8 with inner RFC 5915 ECPrivateKey blobs (NSS-wrapped keys). [[GH-15742](https://github.com/hashicorp/vault/pull/15742)]
* secrets/ad: set config default length only if password_policy is missing [[GH-16140](https://github.com/hashicorp/vault/pull/16140)]
* secrets/azure: Adds option to permanently delete AzureAD objects created by Vault. [[GH-17045](https://github.com/hashicorp/vault/pull/17045)]
* secrets/database/hana: Add ability to customize dynamic usernames [[GH-16631](https://github.com/hashicorp/vault/pull/16631)]
* secrets/database/snowflake: Add multiplexing support [[GH-17159](https://github.com/hashicorp/vault/pull/17159)]
* secrets/gcp: Updates dependencies: `google.golang.org/api@v0.83.0`, `github.com/hashicorp/go-gcp-common@v0.8.0`. [[GH-17174](https://github.com/hashicorp/vault/pull/17174)]
* secrets/gcpkms: Update dependencies: google.golang.org/api@v0.83.0. [[GH-17199](https://github.com/hashicorp/vault/pull/17199)]
* secrets/kubernetes: upgrade to v0.2.0 [[GH-17164](https://github.com/hashicorp/vault/pull/17164)]
* secrets/pki/tidy: Add another pair of metrics counting certificates not deleted by the tidy operation. [[GH-16702](https://github.com/hashicorp/vault/pull/16702)]
* secrets/pki: Add a new flag to issue/sign APIs which can filter out root CAs from the returned ca_chain field [[GH-16935](https://github.com/hashicorp/vault/pull/16935)]
* secrets/pki: Add a warning to any successful response when the requested TTL is overwritten by MaxTTL [[GH-17073](https://github.com/hashicorp/vault/pull/17073)]
* secrets/pki: Add ability to cancel tidy operations, control tidy resource usage. [[GH-16958](https://github.com/hashicorp/vault/pull/16958)]
* secrets/pki: Add ability to periodically rebuild CRL before expiry [[GH-16762](https://github.com/hashicorp/vault/pull/16762)]
* secrets/pki: Add ability to periodically run tidy operations to remove expired certificates. [[GH-16900](https://github.com/hashicorp/vault/pull/16900)]
* secrets/pki: Add support for per-issuer Authority Information Access (AIA) URLs [[GH-16563](https://github.com/hashicorp/vault/pull/16563)]
* secrets/pki: Add support to specify signature bits when generating CSRs through intermediate/generate apis [[GH-17388](https://github.com/hashicorp/vault/pull/17388)]
* secrets/pki: Added gauge metrics "secrets.pki.total_revoked_certificates_stored" and "secrets.pki.total_certificates_stored" to track the number of certificates in storage. [[GH-16676](https://github.com/hashicorp/vault/pull/16676)]
* secrets/pki: Allow revocation of certificates with explicitly provided certificate (bring your own certificate / BYOC). [[GH-16564](https://github.com/hashicorp/vault/pull/16564)]
* secrets/pki: Allow revocation via proving possession of certificate's private key [[GH-16566](https://github.com/hashicorp/vault/pull/16566)]
* secrets/pki: Allow tidy to associate revoked certs with their issuers for OCSP performance [[GH-16871](https://github.com/hashicorp/vault/pull/16871)]
* secrets/pki: Honor If-Modified-Since header on CA, CRL fetch; requires passthrough_request_headers modification on the mount point. [[GH-16249](https://github.com/hashicorp/vault/pull/16249)]
* secrets/pki: Improve stability of association of revoked cert with its parent issuer; when an issuer loses crl-signing usage, do not place certs on default issuer's CRL. [[GH-16874](https://github.com/hashicorp/vault/pull/16874)]
* secrets/pki: Support generating delta CRLs for up-to-date CRLs when auto-building is enabled. [[GH-16773](https://github.com/hashicorp/vault/pull/16773)]
* secrets/ssh: Add allowed_domains_template to allow templating of allowed_domains. [[GH-16056](https://github.com/hashicorp/vault/pull/16056)]
* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]
* secrets/ssh: Allow the use of Identity templates in the `default_user` field [[GH-16351](https://github.com/hashicorp/vault/pull/16351)]
* secrets/transit: Add a dedicated HMAC key type, which can be used with key import. [[GH-16668](https://github.com/hashicorp/vault/pull/16668)]
* secrets/transit: Added a parameter to encrypt/decrypt batch operations to allow the caller to override the HTTP response code in case of partial user-input failures. [[GH-17118](https://github.com/hashicorp/vault/pull/17118)]
* secrets/transit: Allow configuring the possible salt lengths for RSA PSS signatures. [[GH-16549](https://github.com/hashicorp/vault/pull/16549)]
* ssh: Addition of an endpoint `ssh/issue/:role` to allow the creation of signed key pairs [[GH-15561](https://github.com/hashicorp/vault/pull/15561)]
* storage/cassandra: tuning parameters for clustered environments `connection_timeout`, `initial_connection_timeout`, `simple_retry_policy_retries`. [[GH-10467](https://github.com/hashicorp/vault/pull/10467)]
* storage/gcs: Add documentation explaining how to configure the gcs backend using environment variables instead of options in the configuration stanza [[GH-14455](https://github.com/hashicorp/vault/pull/14455)]
* ui: Changed the tokenBoundCidrs tooltip content to clarify that comma separated values are not accepted in this field. [[GH-15852](https://github.com/hashicorp/vault/pull/15852)]
* ui: Prevents requests to /sys/internal/ui/resultant-acl endpoint when unauthenticated [[GH-17139](https://github.com/hashicorp/vault/pull/17139)]
* ui: Removed deprecated version of core-js 2.6.11 [[GH-15898](https://github.com/hashicorp/vault/pull/15898)]
* ui: Renamed labels under Tools for wrap, lookup, rewrap and unwrap with description. [[GH-16489](https://github.com/hashicorp/vault/pull/16489)]
* ui: Replaces non-inclusive terms [[GH-17116](https://github.com/hashicorp/vault/pull/17116)]
* ui: redirect_to param forwards from auth route when authenticated [[GH-16821](https://github.com/hashicorp/vault/pull/16821)]
* website/docs: API generate-recovery-token documentation. [[GH-16213](https://github.com/hashicorp/vault/pull/16213)]
* website/docs: Add documentation around the expensiveness of making lots of lease count quotas in a short period [[GH-16950](https://github.com/hashicorp/vault/pull/16950)]
* website/docs: Removes mentions of unauthenticated from internal ui resultant-acl doc [[GH-17139](https://github.com/hashicorp/vault/pull/17139)]
* website/docs: Update replication docs to mention Integrated Storage [[GH-16063](https://github.com/hashicorp/vault/pull/16063)]
* website/docs: changed to echo for all string examples instead of (<<<) here-string. [[GH-9081](https://github.com/hashicorp/vault/pull/9081)]  

BUG FIXES:

* agent/template: Fix parsing error for the exec stanza [[GH-16231](https://github.com/hashicorp/vault/pull/16231)]
* agent: Agent will now respect `max_retries` retry configuration even when caching is set. [[GH-16970](https://github.com/hashicorp/vault/pull/16970)]
* agent: Update consul-template for pkiCert bug fixes [[GH-16087](https://github.com/hashicorp/vault/pull/16087)]
* api/sys/internal/specs/openapi: support a new "dynamic" query parameter to generate generic mountpaths [[GH-15835](https://github.com/hashicorp/vault/pull/15835)]
* api: Fixed erroneous warnings of unrecognized parameters when unwrapping data. [[GH-16794](https://github.com/hashicorp/vault/pull/16794)]
* api: Fixed issue with internal/ui/mounts and internal/ui/mounts/(?P<path>.+) endpoints where it was not properly handling /auth/ [[GH-15552](https://github.com/hashicorp/vault/pull/15552)]
* api: properly handle switching to/from unix domain socket when changing client address [[GH-11904](https://github.com/hashicorp/vault/pull/11904)]
* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* auth/kerberos: Maintain headers set by the client [[GH-16636](https://github.com/hashicorp/vault/pull/16636)]
* auth/kubernetes: Restore support for JWT signature algorithm ES384 [[GH-160](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/160)] [[GH-17161](https://github.com/hashicorp/vault/pull/17161)]
* auth/token: Fix ignored parameter warnings for valid parameters on token create [[GH-16938](https://github.com/hashicorp/vault/pull/16938)]
* command/debug: fix bug where monitor was not honoring configured duration [[GH-16834](https://github.com/hashicorp/vault/pull/16834)]
* core (enterprise): Fix bug where wrapping token lookup does not work within namespaces. [[GH-15583](https://github.com/hashicorp/vault/pull/15583)]
* core (enterprise): Fix creation of duplicate entities via alias metadata changes on local auth mounts.
* core/auth: Return a 403 instead of a 500 for a malformed SSCT [[GH-16112](https://github.com/hashicorp/vault/pull/16112)]
* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/license (enterprise): Always remove stored license and allow unseal to complete when license cleanup fails
* core/managed-keys (enterprise): fix panic when having `cache_disable` true
* core/quotas (enterprise): Fixed issue with improper counting of leases if lease count quota created after leases
* core/quotas: Added globbing functionality on the end of path suffix quota paths [[GH-16386](https://github.com/hashicorp/vault/pull/16386)]
* core/quotas: Fix goroutine leak caused by the seal process not fully cleaning up Rate Limit Quotas. [[GH-17281](https://github.com/hashicorp/vault/pull/17281)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* core: Fix panic when the plugin catalog returns neither a plugin nor an error. [[GH-17204](https://github.com/hashicorp/vault/pull/17204)]
* core: Fixes parsing boolean values for ha_storage backends in config [[GH-15900](https://github.com/hashicorp/vault/pull/15900)]
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* core: Prevent two or more DR failovers from invalidating SSCT tokens generated on the previous primaries. [[GH-16956](https://github.com/hashicorp/vault/pull/16956)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* debug: Fix panic when capturing debug bundle on Windows [[GH-14399](https://github.com/hashicorp/vault/pull/14399)]
* debug: Remove extra empty lines from vault.log when debug command is run [[GH-16714](https://github.com/hashicorp/vault/pull/16714)]
* identity (enterprise): Fix a data race when creating an entity for a local alias.
* identity/oidc: Adds `claims_supported` to discovery document. [[GH-16992](https://github.com/hashicorp/vault/pull/16992)]
* identity/oidc: Change the `state` parameter of the Authorization Endpoint to optional. [[GH-16599](https://github.com/hashicorp/vault/pull/16599)]
* identity/oidc: Detect invalid `redirect_uri` values sooner in validation of the Authorization Endpoint. [[GH-16601](https://github.com/hashicorp/vault/pull/16601)]
* identity/oidc: Fixes validation of the `request` and `request_uri` parameters. [[GH-16600](https://github.com/hashicorp/vault/pull/16600)]
* openapi: Fixed issue where information about /auth/token endpoints was not present with explicit policy permissions [[GH-15552](https://github.com/hashicorp/vault/pull/15552)]
* plugin/multiplexing: Fix panic when id doesn't exist in connection map [[GH-16094](https://github.com/hashicorp/vault/pull/16094)]
* plugin/secrets/auth: Fix a bug with aliased backends such as aws-ec2 or generic [[GH-16673](https://github.com/hashicorp/vault/pull/16673)]
* plugins: Corrected the path to check permissions on when the registered plugin name does not match the plugin binary's filename. [[GH-17340](https://github.com/hashicorp/vault/pull/17340)]
* quotas/lease-count: Fix lease-count quotas on mounts not properly being enforced when the lease generating request is a read [[GH-15735](https://github.com/hashicorp/vault/pull/15735)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* replication (enterprise): Fix data race in saveCheckpoint.
* replication (enterprise): Fix possible data race during merkle diff/sync
* secret/pki: Do not fail validation with a legacy key_bits default value and key_type=any when signing CSRs [[GH-16246](https://github.com/hashicorp/vault/pull/16246)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* secrets/gcp: Fixes duplicate static account key creation from performance secondary clusters. [[GH-16534](https://github.com/hashicorp/vault/pull/16534)]
* secrets/kv: Fix `kv get` issue preventing the ability to read a secret when providing a leading slash [[GH-16443](https://github.com/hashicorp/vault/pull/16443)]
* secrets/pki: Allow import of issuers without CRLSign KeyUsage; prohibit setting crl-signing usage on such issuers [[GH-16865](https://github.com/hashicorp/vault/pull/16865)]
* secrets/pki: Do not ignore provided signature bits value when signing intermediate and leaf certificates with a managed key [[GH-17328](https://github.com/hashicorp/vault/pull/17328)]
* secrets/pki: Do not read revoked certificates from backend when CRL is disabled [[GH-17385](https://github.com/hashicorp/vault/pull/17385)]
* secrets/pki: Fix migration to properly handle mounts that contain only keys, no certificates [[GH-16813](https://github.com/hashicorp/vault/pull/16813)]
* secrets/pki: Ignore EC PARAMETER PEM blocks during issuer import (/config/ca, /issuers/import/*, and /intermediate/set-signed) [[GH-16721](https://github.com/hashicorp/vault/pull/16721)]
* secrets/pki: LIST issuers endpoint is now unauthenticated. [[GH-16830](https://github.com/hashicorp/vault/pull/16830)]
* secrets/transform (enterprise): Fix an issue loading tokenization transform configuration after a specific sequence of reconfigurations.
* secrets/transform (enterprise): Fix persistence problem with tokenization store credentials.
* storage/raft (enterprise): Fix some storage-modifying RPCs used by perf standbys that weren't returning the resulting WAL state.
* storage/raft (enterprise): Prevent unauthenticated voter status change with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* storage/raft: Fix retry_join initialization failure [[GH-16550](https://github.com/hashicorp/vault/pull/16550)]
* storage/raft: Nodes no longer get demoted to nonvoter if we don't know their version due to missing heartbeats. [[GH-17019](https://github.com/hashicorp/vault/pull/17019)]
* ui/keymgmt: Sets the defaultValue for type when creating a key. [[GH-17407](https://github.com/hashicorp/vault/pull/17407)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix info tooltip submitting form [[GH-16659](https://github.com/hashicorp/vault/pull/16659)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]
* ui: Fix naming of permitted_dns_domains form parameter on CA creation (root generation and sign intermediate). [[GH-16739](https://github.com/hashicorp/vault/pull/16739)]
* ui: Fixed bug where red spellcheck underline appears in sensitive/secret kv values when it should not appear [[GH-15681](https://github.com/hashicorp/vault/pull/15681)]
* ui: Fixes secret version and status menu links transitioning to auth screen [[GH-16983](https://github.com/hashicorp/vault/pull/16983)]
* ui: OIDC login type uses localStorage instead of sessionStorage [[GH-16170](https://github.com/hashicorp/vault/pull/16170)]
* vault: Fix a bug where duplicate policies could be added to an identity group. [[GH-15638](https://github.com/hashicorp/vault/pull/15638)]

## 1.11.7
### February 6, 2023
  
CHANGES:

* core: Bump Go version to 1.19.4.

IMPROVEMENTS:

* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]

BUG FIXES:

* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* core (enterprise): Fix a race condition resulting in login errors to PKCS#11 modules under high concurrency.
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* identity (enterprise): Fix a data race when creating an entity for a local alias.
* kmip (enterprise): Fix Destroy operation response that omitted Unique Identifier on some batched responses.
* kmip (enterprise): Fix Locate operation response incompatibility with clients using KMIP versions prior to 1.3.
* kmip (enterprise): Fix Query operation response that omitted streaming capability and supported profiles.
* licensing (enterprise): update autoloaded license cache after reload
* secrets/pki: Allow patching issuer to set an empty issuer name. [[GH-18466](https://github.com/hashicorp/vault/pull/18466)]
* secrets/transit: Do not warn about unrecognized parameter 'batch_input' [[GH-18299](https://github.com/hashicorp/vault/pull/18299)]
* storage/raft (enterprise): An already joined node can rejoin by wiping storage
and re-issueing a join request, but in doing so could transiently become a
non-voter.  In some scenarios this resulted in loss of quorum. [[GH-18263](https://github.com/hashicorp/vault/pull/18263)]
* storage/raft (enterprise): Fix some storage-modifying RPCs used by perf standbys that weren't returning the resulting WAL state.
* storage/raft: Don't panic on unknown raft ops [[GH-17732](https://github.com/hashicorp/vault/pull/17732)]
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]

## 1.11.6
### November 30, 2022

IMPROVEMENTS:

* secrets/pki: Allow issuer creation, import to change default issuer via `default_follows_latest_issuer`. [[GH-17824](https://github.com/hashicorp/vault/pull/17824)]

BUG FIXES:

* auth/okta: fix a panic for AuthRenew in Okta [[GH-18011](https://github.com/hashicorp/vault/pull/18011)]
* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* cli: Fix issue preventing kv commands from executing properly when the mount path provided by `-mount` flag and secret key path are the same. [[GH-17679](https://github.com/hashicorp/vault/pull/17679)]
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: fix a start up race condition where performance standbys could go into a
  mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* secret/pki: fix bug with initial legacy bundle migration (from < 1.11 into 1.11+) and missing issuers from ca_chain [[GH-17772](https://github.com/hashicorp/vault/pull/17772)]
* secrets/azure: add WAL to clean up role assignments if errors occur [[GH-18085](https://github.com/hashicorp/vault/pull/18085)]
* secrets/gcp: Fixes duplicate service account key for rotate root on standby or secondary [[GH-18110](https://github.com/hashicorp/vault/pull/18110)]
* ui: Fixes issue with not being able to download raft snapshot via service worker [[GH-17769](https://github.com/hashicorp/vault/pull/17769)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]

## 1.11.5
### November 2, 2022

IMPROVEMENTS:

* database/snowflake: Allow parallel requests to Snowflake [[GH-17594](https://github.com/hashicorp/vault/pull/17594)]
* sdk/ldap: Added support for paging when searching for groups using group filters [[GH-17640](https://github.com/hashicorp/vault/pull/17640)]

BUG FIXES:

* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): fix panic when having `cache_disable` true
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* kmip (enterprise): Fix a problem in the handling of attributes that caused Import operations to fail.
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* secrets/pki: Do not ignore provided signature bits value when signing intermediate and leaf certificates with a managed key [[GH-17328](https://github.com/hashicorp/vault/pull/17328)]
* secrets/pki: Do not read revoked certificates from backend when CRL is disabled [[GH-17384](https://github.com/hashicorp/vault/pull/17384)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* ui/keymgmt: Sets the defaultValue for type when creating a key. [[GH-17407](https://github.com/hashicorp/vault/pull/17407)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]
  
## 1.11.4
### September 30, 2022

IMPROVEMENTS:

* agent/auto-auth: Add `exit_on_err` which when set to true, will cause Agent to exit if any errors are encountered during authentication. [[GH-17091](https://github.com/hashicorp/vault/pull/17091)]
* agent: Send notifications to systemd on start and stop. [[GH-9802](https://github.com/hashicorp/vault/pull/9802)]

BUG FIXES:

* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* auth/kubernetes: Restore support for JWT signature algorithm ES384 [[GH-160](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/160)] [[GH-17162](https://github.com/hashicorp/vault/pull/17162)]
* auth/token: Fix ignored parameter warnings for valid parameters on token create [[GH-16938](https://github.com/hashicorp/vault/pull/16938)]
* core/quotas: Fix goroutine leak caused by the seal process not fully cleaning up Rate Limit Quotas. [[GH-17281](https://github.com/hashicorp/vault/pull/17281)]
* core: Prevent two or more DR failovers from invalidating SSCT tokens generated on the previous primaries. [[GH-16956](https://github.com/hashicorp/vault/pull/16956)]
* identity/oidc: Adds `claims_supported` to discovery document. [[GH-16992](https://github.com/hashicorp/vault/pull/16992)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* secrets/transform (enterprise): Fix an issue loading tokenization transform configuration after a specific sequence of reconfigurations.
* secrets/transform (enterprise): Fix persistence problem with tokenization store credentials.
* ui: Fixes secret version and status menu links transitioning to auth screen [[GH-16983](https://github.com/hashicorp/vault/pull/16983)]  
* ui: Fixes secret version and status menu links transitioning to auth screen [[GH-16983](https://github.com/hashicorp/vault/pull/16983)]

## 1.11.3
### August 31, 2022

CHANGES:

* core: Bump Go version to 1.17.13.

IMPROVEMENTS:

* auth/kerberos: add `add_group_aliases` config to include LDAP groups in Vault group aliases [[GH-16890](https://github.com/hashicorp/vault/pull/16890)]
* auth/kerberos: add `remove_instance_name` parameter to the login CLI and the 
Kerberos config in Vault. This removes any instance names found in the keytab 
service principal name. [[GH-16594](https://github.com/hashicorp/vault/pull/16594)]
* identity/oidc: Adds the `client_secret_post` token endpoint authentication method. [[GH-16598](https://github.com/hashicorp/vault/pull/16598)]
* storage/gcs: Add documentation explaining how to configure the gcs backend using environment variables instead of options in the configuration stanza [[GH-14455](https://github.com/hashicorp/vault/pull/14455)]

BUG FIXES:

* api: Fixed erroneous warnings of unrecognized parameters when unwrapping data. [[GH-16794](https://github.com/hashicorp/vault/pull/16794)]
* auth/gcp: Fixes the ability to reset the configuration's credentials to use application default credentials. [[GH-16523](https://github.com/hashicorp/vault/pull/16523)]
* auth/kerberos: Maintain headers set by the client [[GH-16636](https://github.com/hashicorp/vault/pull/16636)]
* command/debug: fix bug where monitor was not honoring configured duration [[GH-16834](https://github.com/hashicorp/vault/pull/16834)]
* core/license (enterprise): Always remove stored license and allow unseal to complete when license cleanup fails
* database/elasticsearch: Fixes a bug in boolean parsing for initialize [[GH-16526](https://github.com/hashicorp/vault/pull/16526)]
* identity/oidc: Change the `state` parameter of the Authorization Endpoint to optional. [[GH-16599](https://github.com/hashicorp/vault/pull/16599)]
* identity/oidc: Detect invalid `redirect_uri` values sooner in validation of the 
Authorization Endpoint. [[GH-16601](https://github.com/hashicorp/vault/pull/16601)]
* identity/oidc: Fixes validation of the `request` and `request_uri` parameters. [[GH-16600](https://github.com/hashicorp/vault/pull/16600)]
* plugin/secrets/auth: Fix a bug with aliased backends such as aws-ec2 or generic [[GH-16673](https://github.com/hashicorp/vault/pull/16673)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* secrets/gcp: Fixes duplicate static account key creation from performance secondary clusters. [[GH-16534](https://github.com/hashicorp/vault/pull/16534)]
* secrets/pki: Fix migration to properly handle mounts that contain only keys, no certificates [[GH-16813](https://github.com/hashicorp/vault/pull/16813)]
* secrets/pki: Ignore EC PARAMETER PEM blocks during issuer import (/config/ca, /issuers/import/*, and /intermediate/set-signed) [[GH-16721](https://github.com/hashicorp/vault/pull/16721)]
* secrets/pki: LIST issuers endpoint is now unauthenticated. [[GH-16830](https://github.com/hashicorp/vault/pull/16830)]
* storage/raft: Fix retry_join initialization failure [[GH-16550](https://github.com/hashicorp/vault/pull/16550)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix info tooltip submitting form [[GH-16659](https://github.com/hashicorp/vault/pull/16659)]
* ui: Fix naming of permitted_dns_domains form parameter on CA creation (root generation and sign intermediate). [[GH-16739](https://github.com/hashicorp/vault/pull/16739)]

SECURITY:

* identity/entity: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. [[HCSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

## 1.11.2
### August 2, 2022

IMPROVEMENTS:

* agent: Added `disable_keep_alives` configuration to disable keep alives in auto-auth, caching and templating. [[GH-16479](https://github.com/hashicorp/vault/pull/16479)]

BUG FIXES:

* core/auth: Return a 403 instead of a 500 for a malformed SSCT [[GH-16112](https://github.com/hashicorp/vault/pull/16112)]
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* secrets/kv: Fix `kv get` issue preventing the ability to read a secret when providing a leading slash [[GH-16443](https://github.com/hashicorp/vault/pull/16443)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]

## 1.11.1
### July 21, 2022

CHANGES:

* core: Bump Go version to 1.17.12.

IMPROVEMENTS:

* agent: Added `disable_idle_connections` configuration to disable leaving idle connections open in auto-auth, caching and templating. [[GH-15986](https://github.com/hashicorp/vault/pull/15986)]
* core: Add `sys/loggers` and `sys/loggers/:name` endpoints to provide ability to modify logging verbosity [[GH-16111](https://github.com/hashicorp/vault/pull/16111)]
* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]

BUG FIXES:

* agent/template: Fix parsing error for the exec stanza [[GH-16231](https://github.com/hashicorp/vault/pull/16231)]
* agent: Update consul-template for pkiCert bug fixes [[GH-16087](https://github.com/hashicorp/vault/pull/16087)]
* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* kmip (enterprise): Return SecretData as supported Object Type.
* plugin/multiplexing: Fix panic when id doesn't exist in connection map [[GH-16094](https://github.com/hashicorp/vault/pull/16094)]
* secret/pki: Do not fail validation with a legacy key_bits default value and key_type=any when signing CSRs [[GH-16246](https://github.com/hashicorp/vault/pull/16246)]
* storage/raft (enterprise): Prevent unauthenticated voter status change with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* transform (enterprise): Fix a bug in the handling of nested or unmatched capture groups in FPE transformations.
* ui: OIDC login type uses localStorage instead of sessionStorage [[GH-16170](https://github.com/hashicorp/vault/pull/16170)]

SECURITY:

* storage/raft (enterprise): Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HCSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]

## 1.11.0
### June 20, 2022

CHANGES:

* auth/aws: Add RoleSession to DisplayName when using assumeRole for authentication [[GH-14954](https://github.com/hashicorp/vault/pull/14954)]
* auth/kubernetes: If `kubernetes_ca_cert` is unset, and there is no pod-local CA available, an error will be surfaced when writing config instead of waiting for login. [[GH-15584](https://github.com/hashicorp/vault/pull/15584)]
* auth: Remove support for legacy MFA
(https://www.vaultproject.io/docs/v1.10.x/auth/mfa) [[GH-14869](https://github.com/hashicorp/vault/pull/14869)]
* core/fips: Disable and warn about entropy augmentation in FIPS 140-2 Inside mode [[GH-15858](https://github.com/hashicorp/vault/pull/15858)]
* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.17.11. [[GH-go-ver-1110](https://github.com/hashicorp/vault/pull/go-ver-1110)]
* database & storage: Change underlying driver library from [lib/pq](https://github.com/lib/pq) to [pgx](https://github.com/jackc/pgx). This change affects Redshift & Postgres database secrets engines, and CockroachDB & Postgres storage engines [[GH-15343](https://github.com/hashicorp/vault/pull/15343)]
* licensing (enterprise): Remove support for stored licenses and associated `sys/license` and `sys/license/signed`
endpoints in favor of [autoloaded licenses](https://www.vaultproject.io/docs/enterprise/license/autoloading).
* replication (enterprise): The `/sys/replication/performance/primary/mount-filter` endpoint has been removed. Please use [Paths Filter](https://www.vaultproject.io/api-docs/system/replication/replication-performance#create-paths-filter) instead.
* secret/pki: Remove unused signature_bits parameter from intermediate CSR generation; this parameter doesn't control the final certificate's signature algorithm selection as that is up to the signing CA [[GH-15478](https://github.com/hashicorp/vault/pull/15478)]
* secrets/kubernetes: Split `additional_metadata` into `extra_annotations` and `extra_labels` parameters [[GH-15655](https://github.com/hashicorp/vault/pull/15655)]
* secrets/pki: A new aliased api path (/pki/issuer/:issuer_ref/sign-self-issued)
providing the same functionality as the existing API(/pki/root/sign-self-issued)
does not require sudo capabilities but the latter still requires it in an
effort to maintain backwards compatibility. [[GH-15211](https://github.com/hashicorp/vault/pull/15211)]
* secrets/pki: Err on unknown role during sign-verbatim. [[GH-15543](https://github.com/hashicorp/vault/pull/15543)]
* secrets/pki: Existing CRL API (/pki/crl) now returns an X.509 v2 CRL instead
of a v1 CRL. [[GH-15100](https://github.com/hashicorp/vault/pull/15100)]
* secrets/pki: The `ca_chain` response field within issuing (/pki/issue/:role)
and signing APIs will now include the root CA certificate if the mount is
aware of it. [[GH-15155](https://github.com/hashicorp/vault/pull/15155)]
* secrets/pki: existing Delete Root API (pki/root) will now delete all issuers
and keys within the mount path. [[GH-15004](https://github.com/hashicorp/vault/pull/15004)]
* secrets/pki: existing Generate Root (pki/root/generate/:type),
Set Signed Intermediate (/pki/intermediate/set-signed) APIs will
add new issuers/keys to a mount instead of warning that an existing CA exists [[GH-14975](https://github.com/hashicorp/vault/pull/14975)]
* secrets/pki: the signed CA certificate from the sign-intermediate api will now appear within the ca_chain
response field along with the issuer's ca chain. [[GH-15524](https://github.com/hashicorp/vault/pull/15524)]
* ui: Upgrade Ember to version 3.28 [[GH-14763](https://github.com/hashicorp/vault/pull/14763)]

FEATURES:

* **Autopilot Improvements (Enterprise)**: Autopilot on Vault Enterprise now supports automated upgrades and redundancy zones when using integrated storage.
* **KeyMgmt UI**: Add UI support for managing the Key Management Secrets Engine [[GH-15523](https://github.com/hashicorp/vault/pull/15523)]
* **Kubernetes Secrets Engine**: This new secrets engine generates Kubernetes service account tokens, service accounts, role bindings, and roles dynamically. [[GH-15551](https://github.com/hashicorp/vault/pull/15551)]
* **Non-Disruptive Intermediate/Root Certificate Rotation**: This allows
import, generation and configuration of any number of keys and/or issuers
within a PKI mount, providing operators the ability to rotate certificates
in place without affecting existing client configurations. [[GH-15277](https://github.com/hashicorp/vault/pull/15277)]
* **Print minimum required policy for any command**: The global CLI flag `-output-policy` can now be used with any command to print out the minimum required policy HCL for that operation, including whether the given path requires the "sudo" capability. [[GH-14899](https://github.com/hashicorp/vault/pull/14899)]
* **Snowflake Database Plugin**: Adds ability to manage RSA key pair credentials for dynamic and static Snowflake users. [[GH-15376](https://github.com/hashicorp/vault/pull/15376)]
* **Transit BYOK**: Allow import of externally-generated keys into the Transit secrets engine. [[GH-15414](https://github.com/hashicorp/vault/pull/15414)]
* nomad: Bootstrap Nomad ACL system if no token is provided [[GH-12451](https://github.com/hashicorp/vault/pull/12451)]
* storage/dynamodb: Added `AWS_DYNAMODB_REGION` environment variable. [[GH-15054](https://github.com/hashicorp/vault/pull/15054)]

IMPROVEMENTS:

* activity: return nil response months in activity log API when no month data exists [[GH-15420](https://github.com/hashicorp/vault/pull/15420)]
* agent/auto-auth: Add `min_backoff` to the method stanza for configuring initial backoff duration. [[GH-15204](https://github.com/hashicorp/vault/pull/15204)]
* agent: Update consul-template to v0.29.0 [[GH-15293](https://github.com/hashicorp/vault/pull/15293)]
* agent: Upgrade hashicorp/consul-template version for sprig template functions and improved writeTo function [[GH-15092](https://github.com/hashicorp/vault/pull/15092)]
* api/monitor: Add log_format option to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* api: Add ability to pass certificate as PEM bytes to api.Client. [[GH-14753](https://github.com/hashicorp/vault/pull/14753)]
* api: Add context-aware functions to vault/api for each API wrapper function. [[GH-14388](https://github.com/hashicorp/vault/pull/14388)]
* api: Added MFALogin() for handling MFA flow when using login helpers. [[GH-14900](https://github.com/hashicorp/vault/pull/14900)]
* api: If the parameters supplied over the API payload are ignored due to not
being what the endpoints were expecting, or if the parameters supplied get
replaced by the values in the endpoint's path itself, warnings will be added to
the non-empty responses listing all the ignored and replaced parameters. [[GH-14962](https://github.com/hashicorp/vault/pull/14962)]
* api: KV helper methods to simplify the common use case of reading and writing KV secrets [[GH-15305](https://github.com/hashicorp/vault/pull/15305)]
* api: Provide a helper method WithNamespace to create a cloned client with a new NS [[GH-14963](https://github.com/hashicorp/vault/pull/14963)]
* api: Support VAULT_PROXY_ADDR environment variable to allow overriding the Vault client's HTTP proxy. [[GH-15377](https://github.com/hashicorp/vault/pull/15377)]
* api: Use the context passed to the api/auth Login helpers. [[GH-14775](https://github.com/hashicorp/vault/pull/14775)]
* api: make ListPlugins parse only known plugin types [[GH-15434](https://github.com/hashicorp/vault/pull/15434)]
* audit: Add a policy_results block into the audit log that contains the set of
policies that granted this request access. [[GH-15457](https://github.com/hashicorp/vault/pull/15457)]
* audit: Include mount_accessor in audit request and response logs [[GH-15342](https://github.com/hashicorp/vault/pull/15342)]
* audit: added entity_created boolean to audit log, set when login operations create an entity [[GH-15487](https://github.com/hashicorp/vault/pull/15487)]
* auth/aws: Add rsa2048 signature type to API [[GH-15719](https://github.com/hashicorp/vault/pull/15719)]
* auth/gcp: Enable the Google service endpoints used by the underlying client to be customized [[GH-15592](https://github.com/hashicorp/vault/pull/15592)]
* auth/gcp: Vault CLI now infers the service account email when running on Google Cloud [[GH-15592](https://github.com/hashicorp/vault/pull/15592)]
* auth/jwt: Adds ability to use JSON pointer syntax for the `user_claim` value. [[GH-15593](https://github.com/hashicorp/vault/pull/15593)]
* auth/okta: Add support for Google provider TOTP type in the Okta auth method [[GH-14985](https://github.com/hashicorp/vault/pull/14985)]
* auth/okta: Add support for performing [the number
challenge](https://help.okta.com/en-us/Content/Topics/Mobile/ov-admin-config.htm?cshid=csh-okta-verify-number-challenge-v1#enable-number-challenge)
during an Okta Verify push challenge [[GH-15361](https://github.com/hashicorp/vault/pull/15361)]
* auth: Globally scoped Login MFA method Get/List endpoints [[GH-15248](https://github.com/hashicorp/vault/pull/15248)]
* auth: enforce a rate limit for TOTP passcode validation attempts [[GH-14864](https://github.com/hashicorp/vault/pull/14864)]
* auth: forward cached MFA auth response to the leader using RPC instead of forwarding all login requests [[GH-15469](https://github.com/hashicorp/vault/pull/15469)]
* cli/debug: added support for retrieving metrics from DR clusters if `unauthenticated_metrics_access` is enabled [[GH-15316](https://github.com/hashicorp/vault/pull/15316)]
* cli/vault: warn when policy name contains upper-case letter [[GH-14670](https://github.com/hashicorp/vault/pull/14670)]
* cli: Alternative flag-based syntax for KV to mitigate confusion from automatically appended /data [[GH-14807](https://github.com/hashicorp/vault/pull/14807)]
* cockroachdb: add high-availability support [[GH-12965](https://github.com/hashicorp/vault/pull/12965)]
* command/debug: Add log_format flag to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* command: Support optional '-log-level' flag to be passed to 'operator migrate' command (defaults to info). Also support VAULT_LOG_LEVEL env var. [[GH-15405](https://github.com/hashicorp/vault/pull/15405)]
* command: Support the optional '-detailed' flag to be passed to 'vault list' command to show ListResponseWithInfo data. Also supports the VAULT_DETAILED env var. [[GH-15417](https://github.com/hashicorp/vault/pull/15417)]
* core (enterprise): Include `termination_time` in `sys/license/status` response
* core (enterprise): Include termination time in `license inspect` command output
* core,transit: Allow callers to choose random byte source including entropy augmentation sources for the sys/tools/random and transit/random endpoints. [[GH-15213](https://github.com/hashicorp/vault/pull/15213)]
* core/activity: Order month data in ascending order of timestamps [[GH-15259](https://github.com/hashicorp/vault/pull/15259)]
* core/activity: allow client counts to be precomputed and queried on non-contiguous chunks of data [[GH-15352](https://github.com/hashicorp/vault/pull/15352)]
* core/managed-keys (enterprise): Allow configuring the number of parallel operations to PKCS#11 managed keys.
* core: Add an export API for historical activity log data [[GH-15586](https://github.com/hashicorp/vault/pull/15586)]
* core: Add new DB methods that do not prepare statements. [[GH-15166](https://github.com/hashicorp/vault/pull/15166)]
* core: check uid and permissions of config dir, config file, plugin dir and plugin binaries [[GH-14817](https://github.com/hashicorp/vault/pull/14817)]
* core: Fix some identity data races found by Go race detector (no known impact yet). [[GH-15123](https://github.com/hashicorp/vault/pull/15123)]
* core: Include build date in `sys/seal-status` and `sys/version-history` endpoints. [[GH-14957](https://github.com/hashicorp/vault/pull/14957)]
* core: Upgrade github.org/x/crypto/ssh [[GH-15125](https://github.com/hashicorp/vault/pull/15125)]
* kmip (enterprise): Implement operations Query, Import, Encrypt and Decrypt. Improve operations Locate, Add Attribute, Get Attributes and Get Attribute List to handle most supported attributes.
* mfa/okta: migrate to use official Okta SDK [[GH-15355](https://github.com/hashicorp/vault/pull/15355)]
* sdk: Change OpenAPI code generator to extract request objects into /components/schemas and reference them by name. [[GH-14217](https://github.com/hashicorp/vault/pull/14217)]
* secrets/consul: Add support for Consul node-identities and service-identities [[GH-15295](https://github.com/hashicorp/vault/pull/15295)]
* secrets/consul: Vault is now able to automatically bootstrap the Consul ACL system. [[GH-10751](https://github.com/hashicorp/vault/pull/10751)]
* secrets/database/elasticsearch: Use the new /_security base API path instead of /_xpack/security when managing elasticsearch. [[GH-15614](https://github.com/hashicorp/vault/pull/15614)]
* secrets/pki: Add not_before_duration to root CA generation, intermediate CA signing paths. [[GH-14178](https://github.com/hashicorp/vault/pull/14178)]
* secrets/pki: Add support for CPS URLs and User Notice to Policy Information [[GH-15751](https://github.com/hashicorp/vault/pull/15751)]
* secrets/pki: Allow operators to control the issuing certificate behavior when
the requested TTL is beyond the NotAfter value of the signing certificate [[GH-15152](https://github.com/hashicorp/vault/pull/15152)]
* secrets/pki: Always return CRLs, URLs configurations, even if using the default value. [[GH-15470](https://github.com/hashicorp/vault/pull/15470)]
* secrets/pki: Enable Patch Functionality for Roles and Issuers (API only) [[GH-15510](https://github.com/hashicorp/vault/pull/15510)]
* secrets/pki: Have pki/sign-verbatim use the not_before_duration field defined in the role [[GH-15429](https://github.com/hashicorp/vault/pull/15429)]
* secrets/pki: Warn on empty Subject field during issuer generation (root/generate and root/sign-intermediate). [[GH-15494](https://github.com/hashicorp/vault/pull/15494)]
* secrets/pki: Warn on missing AIA access information when generating issuers (config/urls). [[GH-15509](https://github.com/hashicorp/vault/pull/15509)]
* secrets/pki: Warn when `generate_lease` and `no_store` are both set to `true` on requests. [[GH-14292](https://github.com/hashicorp/vault/pull/14292)]
* secrets/ssh: Add connection timeout of 1 minute for outbound SSH connection in deprecated Dynamic SSH Keys mode. [[GH-15440](https://github.com/hashicorp/vault/pull/15440)]
* secrets/ssh: Support for `add_before_duration` in SSH [[GH-15250](https://github.com/hashicorp/vault/pull/15250)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer
* storage/raft: Use larger timeouts at startup to reduce likelihood of inducing elections. [[GH-15042](https://github.com/hashicorp/vault/pull/15042)]
* ui: Allow namespace param to be parsed from state queryParam [[GH-15378](https://github.com/hashicorp/vault/pull/15378)]
* ui: Default auto-rotation period in transit is 30 days [[GH-15474](https://github.com/hashicorp/vault/pull/15474)]
* ui: Parse schema refs from OpenAPI [[GH-14508](https://github.com/hashicorp/vault/pull/14508)]
* ui: Remove stored license references [[GH-15513](https://github.com/hashicorp/vault/pull/15513)]
* ui: Remove storybook. [[GH-15074](https://github.com/hashicorp/vault/pull/15074)]
* ui: Replaces the IvyCodemirror wrapper with a custom ember modifier. [[GH-14659](https://github.com/hashicorp/vault/pull/14659)]
* website/docs: Add usage documentation for Kubernetes Secrets Engine [[GH-15527](https://github.com/hashicorp/vault/pull/15527)]
* website/docs: added a link to an Enigma secret plugin. [[GH-14389](https://github.com/hashicorp/vault/pull/14389)]

DEPRECATIONS:

* docs: Document removal of X.509 certificates with signatures who use SHA-1 in Vault 1.12 [[GH-15581](https://github.com/hashicorp/vault/pull/15581)]
* secrets/consul: Deprecate old parameters "token_type" and "policy" [[GH-15550](https://github.com/hashicorp/vault/pull/15550)]
* secrets/consul: Deprecate parameter "policies" in favor of "consul_policies" for consistency [[GH-15400](https://github.com/hashicorp/vault/pull/15400)]

BUG FIXES:

* Fixed panic when adding or modifying a Duo MFA Method in Enterprise
* agent: Fix log level mismatch between ERR and ERROR [[GH-14424](https://github.com/hashicorp/vault/pull/14424)]
* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* api: Fixes bug where OutputCurlString field was unintentionally being copied over during client cloning [[GH-14968](https://github.com/hashicorp/vault/pull/14968)]
* api: Respect increment value in grace period calculations in LifetimeWatcher [[GH-14836](https://github.com/hashicorp/vault/pull/14836)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* auth/kubernetes: Fix error code when using the wrong service account [[GH-15584](https://github.com/hashicorp/vault/pull/15584)]
* auth/ldap: The logic for setting the entity alias when `username_as_alias` is set
has been fixed. The previous behavior would make a request to the LDAP server to
get `user_attr` before discarding it and using the username instead. This would
make it impossible for a user to connect if this attribute was missing or had
multiple values, even though it would not be used anyway. This has been fixed
and the username is now used without making superfluous LDAP searches. [[GH-15525](https://github.com/hashicorp/vault/pull/15525)]
* auth: Fixed erroneous success message when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed erroneous token information being displayed when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed two-phase MFA information missing from table format when using vault login [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Prevent deleting a valid MFA method ID using the endpoint for a different MFA method type [[GH-15482](https://github.com/hashicorp/vault/pull/15482)]
* auth: forward requests subject to login MFA from perfStandby to Active node [[GH-15009](https://github.com/hashicorp/vault/pull/15009)]
* auth: load login MFA configuration upon restart [[GH-15261](https://github.com/hashicorp/vault/pull/15261)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* cli: kv get command now honors trailing spaces to retrieve secrets [[GH-15188](https://github.com/hashicorp/vault/pull/15188)]
* command: do not report listener and storage types as key not found warnings [[GH-15383](https://github.com/hashicorp/vault/pull/15383)]
* core (enterprise): Allow local alias create RPCs to persist alias metadata
* core (enterprise): Fix overcounting of lease count quota usage at startup.
* core (enterprise): Fix some races in merkle index flushing code found in testing
* core (enterprise): Handle additional edge cases reinitializing PKCS#11 libraries after login errors.
* core/config: Only ask the system about network interfaces when address configs contain a template having the format: {{ ... }} [[GH-15224](https://github.com/hashicorp/vault/pull/15224)]
* core/managed-keys (enterprise): Allow PKCS#11 managed keys to use 0 as a slot number
* core/metrics: Fix incorrect table size metric for local mounts [[GH-14755](https://github.com/hashicorp/vault/pull/14755)]
* core: Fix double counting for "route" metrics [[GH-12763](https://github.com/hashicorp/vault/pull/12763)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: Limit SSCT WAL checks on perf standbys to raft backends only [[GH-15879](https://github.com/hashicorp/vault/pull/15879)]
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: Prevent metrics generation from causing deadlocks. [[GH-15693](https://github.com/hashicorp/vault/pull/15693)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]
* core: renaming the environment variable VAULT_DISABLE_FILE_PERMISSIONS_CHECK to VAULT_ENABLE_FILE_PERMISSIONS_CHECK and adjusting the logic [[GH-15452](https://github.com/hashicorp/vault/pull/15452)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* identity: deduplicate policies when creating/updating identity groups [[GH-15055](https://github.com/hashicorp/vault/pull/15055)]
* mfa/okta: disable client side rate limiting causing delays in push notifications [[GH-15369](https://github.com/hashicorp/vault/pull/15369)]
* plugin: Fix a bug where plugin reload would falsely report success in certain scenarios. [[GH-15579](https://github.com/hashicorp/vault/pull/15579)]
* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* raft: Ensure initialMmapSize is set to 0 on Windows [[GH-14977](https://github.com/hashicorp/vault/pull/14977)]
* replication (enterprise): fix panic due to missing entity during invalidation of local aliases. [[GH-14622](https://github.com/hashicorp/vault/pull/14622)]
* sdk/cidrutil: Only check if cidr contains remote address for IP addresses [[GH-14487](https://github.com/hashicorp/vault/pull/14487)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]
* sdk: Fix OpenApi spec generator to remove duplicate sha_256 parameter [[GH-15163](https://github.com/hashicorp/vault/pull/15163)]
* secrets/database: Ensure that a `connection_url` password is redacted in all cases. [[GH-14744](https://github.com/hashicorp/vault/pull/14744)]
* secrets/kv: Fix issue preventing the ability to reset the `delete_version_after` key metadata field to 0s via HTTP `PATCH`. [[GH-15792](https://github.com/hashicorp/vault/pull/15792)]
* secrets/pki: CRLs on performance secondary clusters are now automatically
rebuilt upon changes to the list of issuers. [[GH-15179](https://github.com/hashicorp/vault/pull/15179)]
* secrets/pki: Fix handling of "any" key type with default zero signature bits value. [[GH-14875](https://github.com/hashicorp/vault/pull/14875)]
* secrets/pki: Fixed bug where larger SHA-2 hashes were truncated with shorter ECDSA CA certificates [[GH-14943](https://github.com/hashicorp/vault/pull/14943)]
* secrets/ssh: Convert role field not_before_duration to seconds before returning it [[GH-15559](https://github.com/hashicorp/vault/pull/15559)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* storage/raft: Forward autopilot state requests on perf standbys to active node. [[GH-15493](https://github.com/hashicorp/vault/pull/15493)]
* storage/raft: joining a node to a cluster now ignores any VAULT_NAMESPACE environment variable set on the server process [[GH-15519](https://github.com/hashicorp/vault/pull/15519)]
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not accepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fix KV secret showing in the edit form after a user creates a new version but doesn't have read capabilities [[GH-14794](https://github.com/hashicorp/vault/pull/14794)]
* ui: Fix inconsistent behavior in client count calendar widget [[GH-15789](https://github.com/hashicorp/vault/pull/15789)]
* ui: Fix issue where metadata tab is hidden even though policy grants access [[GH-15824](https://github.com/hashicorp/vault/pull/15824)]
* ui: Fix issue with KV not recomputing model when you changed versions. [[GH-14941](https://github.com/hashicorp/vault/pull/14941)]
* ui: Fixed client count timezone for start and end months [[GH-15167](https://github.com/hashicorp/vault/pull/15167)]
* ui: Fixed unsupported revocation statements field for DB roles [[GH-15573](https://github.com/hashicorp/vault/pull/15573)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging in with OIDC from a listed auth mounts tab [[GH-14916](https://github.com/hashicorp/vault/pull/14916)]
* ui: Revert using localStorage in favor of sessionStorage [[GH-15769](https://github.com/hashicorp/vault/pull/15769)]
* ui: Updated `leasId` to `leaseId` in the "Copy Credentials" section of "Generate AWS Credentials" [[GH-15685](https://github.com/hashicorp/vault/pull/15685)]
* ui: fix firefox inability to recognize file format of client count csv export [[GH-15364](https://github.com/hashicorp/vault/pull/15364)]
* ui: fix form validations ignoring default values and disabling submit button [[GH-15560](https://github.com/hashicorp/vault/pull/15560)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]

## 1.10.10
### February 6, 2023
  
CHANGES:

* core: Bump Go version to 1.19.4.

IMPROVEMENTS:

* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]

BUG FIXES:

* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* auth/token: Fix ignored parameter warnings for valid parameters on token create [[GH-16938](https://github.com/hashicorp/vault/pull/16938)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* core (enterprise): Fix a race condition resulting in login errors to PKCS#11 modules under high concurrency.
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* identity (enterprise): Fix a data race when creating an entity for a local alias.
* kmip (enterprise): Fix Destroy operation response that omitted Unique Identifier on some batched responses.
* kmip (enterprise): Fix Locate operation response incompatibility with clients using KMIP versions prior to 1.3.
* licensing (enterprise): update autoloaded license cache after reload
* storage/raft (enterprise): Fix some storage-modifying RPCs used by perf standbys that weren't returning the resulting WAL state.
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]
  
## 1.10.9
### November 30, 2022

BUG FIXES:

* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: fix a start up race condition where performance standbys could go into a
  mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* secrets/azure: add WAL to clean up role assignments if errors occur [[GH-18084](https://github.com/hashicorp/vault/pull/18084)]
* secrets/gcp: Fixes duplicate service account key for rotate root on standby or secondary [[GH-18109](https://github.com/hashicorp/vault/pull/18109)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]

## 1.10.8
### November 2, 2022
  
BUG FIXES:

* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): fix panic when having `cache_disable` true
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* secrets/pki: Do not ignore provided signature bits value when signing intermediate and leaf certificates with a managed key [[GH-17328](https://github.com/hashicorp/vault/pull/17328)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]

## 1.10.7
### September 30, 2022

BUG FIXES:

* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* core/quotas: Fix goroutine leak caused by the seal process not fully cleaning up Rate Limit Quotas. [[GH-17281](https://github.com/hashicorp/vault/pull/17281)]
* core: Prevent two or more DR failovers from invalidating SSCT tokens generated on the previous primaries. [[GH-16956](https://github.com/hashicorp/vault/pull/16956)]
* identity/oidc: Adds `claims_supported` to discovery document. [[GH-16992](https://github.com/hashicorp/vault/pull/16992)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* secrets/transform (enterprise): Fix an issue loading tokenization transform configuration after a specific sequence of reconfigurations.
* secrets/transform (enterprise): Fix persistence problem with tokenization store credentials.
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]

## 1.10.6
### August 31, 2022

CHANGES:

* core: Bump Go version to 1.17.13.

IMPROVEMENTS:

* identity/oidc: Adds the `client_secret_post` token endpoint authentication method. [[GH-16598](https://github.com/hashicorp/vault/pull/16598)]

BUG FIXES:

* auth/gcp: Fixes the ability to reset the configuration's credentials to use application default credentials. [[GH-16524](https://github.com/hashicorp/vault/pull/16524)]
* command/debug: fix bug where monitor was not honoring configured duration [[GH-16834](https://github.com/hashicorp/vault/pull/16834)]
* core/auth: Return a 403 instead of a 500 for a malformed SSCT [[GH-16112](https://github.com/hashicorp/vault/pull/16112)]
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* identity/oidc: Change the `state` parameter of the Authorization Endpoint to optional. [[GH-16599](https://github.com/hashicorp/vault/pull/16599)]
* identity/oidc: Detect invalid `redirect_uri` values sooner in validation of the 
Authorization Endpoint. [[GH-16601](https://github.com/hashicorp/vault/pull/16601)]
* identity/oidc: Fixes validation of the `request` and `request_uri` parameters. [[GH-16600](https://github.com/hashicorp/vault/pull/16600)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* secrets/gcp: Fixes duplicate static account key creation from performance secondary clusters. [[GH-16534](https://github.com/hashicorp/vault/pull/16534)]
* storage/raft: Fix retry_join initialization failure [[GH-16550](https://github.com/hashicorp/vault/pull/16550)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]
* ui: Fix naming of permitted_dns_domains form parameter on CA creation (root generation and sign intermediate). [[GH-16739](https://github.com/hashicorp/vault/pull/16739)]

SECURITY:

* identity/entity: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. [[HCSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]
  
## 1.10.5
### July 21, 2022

CHANGES:

* core/fips: Disable and warn about entropy augmentation in FIPS 140-2 Inside mode [[GH-15858](https://github.com/hashicorp/vault/pull/15858)]
* core: Bump Go version to 1.17.12.

IMPROVEMENTS:

* core: Add `sys/loggers` and `sys/loggers/:name` endpoints to provide ability to modify logging verbosity [[GH-16111](https://github.com/hashicorp/vault/pull/16111)]
* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]

BUG FIXES:

* agent/template: Fix parsing error for the exec stanza [[GH-16231](https://github.com/hashicorp/vault/pull/16231)]
* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* core: Limit SSCT WAL checks on perf standbys to raft backends only [[GH-15879](https://github.com/hashicorp/vault/pull/15879)]
* plugin/multiplexing: Fix panic when id doesn't exist in connection map [[GH-16094](https://github.com/hashicorp/vault/pull/16094)]
* secret/pki: Do not fail validation with a legacy key_bits default value and key_type=any when signing CSRs [[GH-16246](https://github.com/hashicorp/vault/pull/16246)]
* storage/raft (enterprise): Prevent unauthenticated voter status with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* transform (enterprise): Fix a bug in the handling of nested or unmatched capture groups in FPE transformations.
* ui: Fix issue where metadata tab is hidden even though policy grants access [[GH-15824](https://github.com/hashicorp/vault/pull/15824)]
* ui: Revert using localStorage in favor of sessionStorage [[GH-16169](https://github.com/hashicorp/vault/pull/16169)]
* ui: Updated `leasId` to `leaseId` in the "Copy Credentials" section of "Generate AWS Credentials" [[GH-15685](https://github.com/hashicorp/vault/pull/15685)]

## 1.10.4
### June 10, 2022

CHANGES:

* core: Bump Go version to 1.17.11. [[GH-go-ver-1104](https://github.com/hashicorp/vault/pull/go-ver-1104)]

IMPROVEMENTS:

* api/monitor: Add log_format option to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* auth: Globally scoped Login MFA method Get/List endpoints [[GH-15248](https://github.com/hashicorp/vault/pull/15248)]
* auth: forward cached MFA auth response to the leader using RPC instead of forwarding all login requests [[GH-15469](https://github.com/hashicorp/vault/pull/15469)]
* cli/debug: added support for retrieving metrics from DR clusters if `unauthenticated_metrics_access` is enabled [[GH-15316](https://github.com/hashicorp/vault/pull/15316)]
* command/debug: Add log_format flag to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* core: Fix some identity data races found by Go race detector (no known impact yet). [[GH-15123](https://github.com/hashicorp/vault/pull/15123)]
* storage/raft: Use larger timeouts at startup to reduce likelihood of inducing elections. [[GH-15042](https://github.com/hashicorp/vault/pull/15042)]
* ui: Allow namespace param to be parsed from state queryParam [[GH-15378](https://github.com/hashicorp/vault/pull/15378)]

BUG FIXES:

* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* auth/kubernetes: Fix error code when using the wrong service account [[GH-15585](https://github.com/hashicorp/vault/pull/15585)]
* auth/ldap: The logic for setting the entity alias when `username_as_alias` is set
has been fixed. The previous behavior would make a request to the LDAP server to
get `user_attr` before discarding it and using the username instead. This would
make it impossible for a user to connect if this attribute was missing or had
multiple values, even though it would not be used anyway. This has been fixed
and the username is now used without making superfluous LDAP searches. [[GH-15525](https://github.com/hashicorp/vault/pull/15525)]
* auth: Fixed erroneous success message when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed erroneous token information being displayed when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed two-phase MFA information missing from table format when using vault login [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Prevent deleting a valid MFA method ID using the endpoint for a different MFA method type [[GH-15482](https://github.com/hashicorp/vault/pull/15482)]
* core (enterprise): Fix overcounting of lease count quota usage at startup.
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: Prevent metrics generation from causing deadlocks. [[GH-15693](https://github.com/hashicorp/vault/pull/15693)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* mfa/okta: disable client side rate limiting causing delays in push notifications [[GH-15369](https://github.com/hashicorp/vault/pull/15369)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* transform (enterprise): Fix non-overridable column default value causing tokenization tokens to expire prematurely when using the MySQL storage backend.
* ui: Fix inconsistent behavior in client count calendar widget [[GH-15789](https://github.com/hashicorp/vault/pull/15789)]
* ui: Fixed client count timezone for start and end months [[GH-15167](https://github.com/hashicorp/vault/pull/15167)]
* ui: fix firefox inability to recognize file format of client count csv export [[GH-15364](https://github.com/hashicorp/vault/pull/15364)]

## 1.10.3
### May 11, 2022

SECURITY:
* auth: A vulnerability was identified in Vault and Vault Enterprise (“Vault”) from 1.10.0 to 1.10.2 where MFA may not be enforced on user logins after a server restart. This vulnerability, CVE-2022-30689, was fixed in Vault 1.10.3.

BUG FIXES:

* auth: load login MFA configuration upon restart [[GH-15261](https://github.com/hashicorp/vault/pull/15261)]
* core/config: Only ask the system about network interfaces when address configs contain a template having the format: {{ ... }} [[GH-15224](https://github.com/hashicorp/vault/pull/15224)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]

## 1.10.2
### April 29, 2022

BUG FIXES:

* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]

## 1.10.1
### April 22, 2022

CHANGES:

* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.17.9. [[GH-15044](https://github.com/hashicorp/vault/pull/15044)]

IMPROVEMENTS:

* agent: Upgrade hashicorp/consul-template version for sprig template functions and improved writeTo function [[GH-15092](https://github.com/hashicorp/vault/pull/15092)]
* auth: enforce a rate limit for TOTP passcode validation attempts [[GH-14864](https://github.com/hashicorp/vault/pull/14864)]
* cli/vault: warn when policy name contains upper-case letter [[GH-14670](https://github.com/hashicorp/vault/pull/14670)]
* cockroachdb: add high-availability support [[GH-12965](https://github.com/hashicorp/vault/pull/12965)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer

BUG FIXES:

* Fixed panic when adding or modifying a Duo MFA Method in Enterprise
* agent: Fix log level mismatch between ERR and ERROR [[GH-14424](https://github.com/hashicorp/vault/pull/14424)]
* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* api: Respect increment value in grace period calculations in LifetimeWatcher [[GH-14836](https://github.com/hashicorp/vault/pull/14836)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* auth: forward requests subject to login MFA from perfStandby to Active node [[GH-15009](https://github.com/hashicorp/vault/pull/15009)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* core (enterprise): Allow local alias create RPCs to persist alias metadata [[GH-changelog:_2747](https://github.com/hashicorp/vault/pull/changelog:_2747)]
* core/managed-keys (enterprise): Allow PKCS#11 managed keys to use 0 as a slot number
* core/metrics: Fix incorrect table size metric for local mounts [[GH-14755](https://github.com/hashicorp/vault/pull/14755)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* raft: Ensure initialMmapSize is set to 0 on Windows [[GH-14977](https://github.com/hashicorp/vault/pull/14977)]
* replication (enterprise): fix panic due to missing entity during invalidation of local aliases. [[GH-14622](https://github.com/hashicorp/vault/pull/14622)]
* secrets/database: Ensure that a `connection_url` password is redacted in all cases. [[GH-14744](https://github.com/hashicorp/vault/pull/14744)]
* secrets/pki: Fix handling of "any" key type with default zero signature bits value. [[GH-14875](https://github.com/hashicorp/vault/pull/14875)]
* secrets/pki: Fixed bug where larger SHA-2 hashes were truncated with shorter ECDSA CA certificates [[GH-14943](https://github.com/hashicorp/vault/pull/14943)]
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not excepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging in with OIDC from a listed auth mounts tab [[GH-14916](https://github.com/hashicorp/vault/pull/14916)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]

## 1.10.0
### March 23, 2022

CHANGES:

* core (enterprise): requests with newly generated tokens to perf standbys which are lagging behind the active node return http 412 instead of 400/403/50x.
* core: Changes the unit of `default_lease_ttl` and `max_lease_ttl` values returned by
the `/sys/config/state/sanitized` endpoint from nanoseconds to seconds. [[GH-14206](https://github.com/hashicorp/vault/pull/14206)]
* core: Bump Go version to 1.17.7. [[GH-14232](https://github.com/hashicorp/vault/pull/14232)]
* plugin/database: The return value from `POST /database/config/:name` has been updated to "204 No Content" [[GH-14033](https://github.com/hashicorp/vault/pull/14033)]
* secrets/azure: Changes the configuration parameter `use_microsoft_graph_api` to use the Microsoft 
Graph API by default. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]
* storage/etcd: Remove support for v2. [[GH-14193](https://github.com/hashicorp/vault/pull/14193)]
* ui: Upgrade Ember to version 3.24 [[GH-13443](https://github.com/hashicorp/vault/pull/13443)]

FEATURES:

* **Database plugin multiplexing**: manage multiple database connections with a single plugin process [[GH-14033](https://github.com/hashicorp/vault/pull/14033)]
* **Login MFA**: Single and two phase MFA is now available when authenticating to Vault. [[GH-14025](https://github.com/hashicorp/vault/pull/14025)]
* **Mount Migration**: Vault supports moving secrets and auth mounts both within and across namespaces.
* **Postgres in the UI**: Postgres DB is now supported by the UI [[GH-12945](https://github.com/hashicorp/vault/pull/12945)]
* **Report in-flight requests**: Adding a trace capability to show in-flight requests, and a new gauge metric to show the total number of in-flight requests [[GH-13024](https://github.com/hashicorp/vault/pull/13024)]
* **Server Side Consistent Tokens**: Service tokens have been updated to be longer (a minimum of 95 bytes) and token prefixes for all token types are updated from s., b., and r. to hvs., hvb., and hvr. for service, batch, and recovery tokens respectively. Vault clusters with integrated storage will now have read-after-write consistency by default. [[GH-14109](https://github.com/hashicorp/vault/pull/14109)]
* **Transit SHA-3 Support**: Add support for SHA-3 in the Transit backend. [[GH-13367](https://github.com/hashicorp/vault/pull/13367)]
* **Transit Time-Based Key Autorotation**: Add support for automatic, time-based key rotation to transit secrets engine, including in the UI. [[GH-13691](https://github.com/hashicorp/vault/pull/13691)]
* **UI Client Count Improvements**: Restructures client count dashboard, making use of billing start date to improve accuracy. Adds mount-level distribution and filtering. [[GH-client-counts](https://github.com/hashicorp/vault/pull/client-counts)]
* **Agent Telemetry**: The Vault Agent can now collect and return telemetry information at the `/agent/v1/metrics` endpoint.

IMPROVEMENTS:

* agent: Adds ability to configure specific user-assigned managed identities for Azure auto-auth. [[GH-14214](https://github.com/hashicorp/vault/pull/14214)]
* agent: The `agent/v1/quit` endpoint can now be used to stop the Vault Agent remotely [[GH-14223](https://github.com/hashicorp/vault/pull/14223)]
* api: Allow cloning `api.Client` tokens via `api.Config.CloneToken` or `api.Client.SetCloneToken()`. [[GH-13515](https://github.com/hashicorp/vault/pull/13515)]
* api: Define constants for X-Vault-Forward and X-Vault-Inconsistent headers [[GH-14067](https://github.com/hashicorp/vault/pull/14067)]
* api: Implements Login method in Go client libraries for GCP and Azure auth methods [[GH-13022](https://github.com/hashicorp/vault/pull/13022)]
* api: Implements Login method in Go client libraries for LDAP auth methods [[GH-13841](https://github.com/hashicorp/vault/pull/13841)]
* api: Trim newline character from wrapping token in logical.Unwrap from the api package [[GH-13044](https://github.com/hashicorp/vault/pull/13044)]
* api: add api method for modifying raft autopilot configuration [[GH-12428](https://github.com/hashicorp/vault/pull/12428)]
* api: respect WithWrappingToken() option during AppRole login authentication when used with secret ID specified from environment or from string [[GH-13241](https://github.com/hashicorp/vault/pull/13241)]
* audit: The audit logs now contain the port used by the client [[GH-12790](https://github.com/hashicorp/vault/pull/12790)]
* auth/aws: Enable region detection in the CLI by specifying the region as `auto` [[GH-14051](https://github.com/hashicorp/vault/pull/14051)]
* auth/cert: Add certificate extensions as metadata [[GH-13348](https://github.com/hashicorp/vault/pull/13348)]
* auth/jwt: The Authorization Code flow makes use of the Proof Key for Code Exchange (PKCE) extension. [[GH-13365](https://github.com/hashicorp/vault/pull/13365)]
* auth/kubernetes: Added support for dynamically reloading short-lived tokens for better Kubernetes 1.21+ compatibility [[GH-13595](https://github.com/hashicorp/vault/pull/13595)]
* auth/ldap: Add a response warning and server log whenever the config is accessed
if `userfilter` doesn't consider `userattr` [[GH-14095](https://github.com/hashicorp/vault/pull/14095)]
* auth/ldap: Add username to alias metadata [[GH-13669](https://github.com/hashicorp/vault/pull/13669)]
* auth/ldap: Add username_as_alias configurable to change how aliases are named [[GH-14324](https://github.com/hashicorp/vault/pull/14324)]
* auth/okta: Update [okta-sdk-golang](https://github.com/okta/okta-sdk-golang) dependency to version v2.9.1 for improved request backoff handling [[GH-13439](https://github.com/hashicorp/vault/pull/13439)]
* auth/token: The `auth/token/revoke-accessor` endpoint is now idempotent and will
not error out if the token has already been revoked. [[GH-13661](https://github.com/hashicorp/vault/pull/13661)]
* auth: reading `sys/auth/:path` now returns the configuration for the auth engine mounted at the given path [[GH-12793](https://github.com/hashicorp/vault/pull/12793)]
* cli: interactive CLI for login mfa [[GH-14131](https://github.com/hashicorp/vault/pull/14131)]
* command (enterprise): "vault license get" now uses non-deprecated endpoint /sys/license/status
* core/ha: Add new mechanism for keeping track of peers talking to active node, and new 'operator members' command to view them. [[GH-13292](https://github.com/hashicorp/vault/pull/13292)]
* core/identity: Support updating an alias' `custom_metadata` to be empty. [[GH-13395](https://github.com/hashicorp/vault/pull/13395)]
* core/pki: Support Y10K value in notAfter field to be compliant with IEEE 802.1AR-2018 standard [[GH-12795](https://github.com/hashicorp/vault/pull/12795)]
* core/pki: Support Y10K value in notAfter field when signing non-CA certificates [[GH-13736](https://github.com/hashicorp/vault/pull/13736)]
* core: Add duration and start_time to completed requests log entries [[GH-13682](https://github.com/hashicorp/vault/pull/13682)]
* core: Add support to list password policies at `sys/policies/password` [[GH-12787](https://github.com/hashicorp/vault/pull/12787)]
* core: Add support to list version history via API at `sys/version-history` and via CLI with `vault version-history` [[GH-13766](https://github.com/hashicorp/vault/pull/13766)]
* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]
* core: Periodically test the health of connectivity to auto-seal backends [[GH-13078](https://github.com/hashicorp/vault/pull/13078)]
* core: Reading `sys/mounts/:path` now returns the configuration for the secret engine at the given path [[GH-12792](https://github.com/hashicorp/vault/pull/12792)]
* core: Replace "master key" terminology with "root key" [[GH-13324](https://github.com/hashicorp/vault/pull/13324)]
* core: Small changes to ensure goroutines terminate in tests [[GH-14197](https://github.com/hashicorp/vault/pull/14197)]
* core: Systemd unit file included with the Linux packages now sets the service type to notify. [[GH-14385](https://github.com/hashicorp/vault/pull/14385)]
* core: Update github.com/prometheus/client_golang to fix security vulnerability CVE-2022-21698. [[GH-14190](https://github.com/hashicorp/vault/pull/14190)]
* core: Vault now supports the PROXY protocol v2. Support for UNKNOWN connections
has also been added to the PROXY protocol v1. [[GH-13540](https://github.com/hashicorp/vault/pull/13540)]
* http (enterprise): Serve /sys/license/status endpoint within namespaces
* identity/oidc: Adds a default OIDC provider [[GH-14119](https://github.com/hashicorp/vault/pull/14119)]
* identity/oidc: Adds a default key for OIDC clients [[GH-14119](https://github.com/hashicorp/vault/pull/14119)]
* identity/oidc: Adds an `allow_all` assignment that permits all entities to authenticate via an OIDC client [[GH-14119](https://github.com/hashicorp/vault/pull/14119)]
* identity/oidc: Adds proof key for code exchange (PKCE) support to OIDC providers. [[GH-13917](https://github.com/hashicorp/vault/pull/13917)]
* sdk: Add helper for decoding root tokens [[GH-10505](https://github.com/hashicorp/vault/pull/10505)]
* secrets/azure: Adds support for rotate-root. [#70](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/70) [[GH-13034](https://github.com/hashicorp/vault/pull/13034)]
* secrets/consul: Add support for consul enterprise namespaces and admin partitions. [[GH-13850](https://github.com/hashicorp/vault/pull/13850)]
* secrets/consul: Add support for consul roles. [[GH-14014](https://github.com/hashicorp/vault/pull/14014)]
* secrets/database/influxdb: Switch/upgrade to the `influxdb1-client` module [[GH-12262](https://github.com/hashicorp/vault/pull/12262)]
* secrets/database: Add database configuration parameter 'disable_escaping' for username and password when connecting to a database. [[GH-13414](https://github.com/hashicorp/vault/pull/13414)]
* secrets/kv: add full secret path output to table-formatted responses [[GH-14301](https://github.com/hashicorp/vault/pull/14301)]
* secrets/kv: add patch support for KVv2 key metadata [[GH-13215](https://github.com/hashicorp/vault/pull/13215)]
* secrets/kv: add subkeys endpoint to retrieve a secret's stucture without its values [[GH-13893](https://github.com/hashicorp/vault/pull/13893)]
* secrets/pki: Add ability to fetch individual certificate as DER or PEM [[GH-10948](https://github.com/hashicorp/vault/pull/10948)]
* secrets/pki: Add count and duration metrics to PKI issue and revoke calls. [[GH-13889](https://github.com/hashicorp/vault/pull/13889)]
* secrets/pki: Add error handling for error types other than UserError or InternalError [[GH-14195](https://github.com/hashicorp/vault/pull/14195)]
* secrets/pki: Allow URI SAN templates in allowed_uri_sans when allowed_uri_sans_template is set to true. [[GH-10249](https://github.com/hashicorp/vault/pull/10249)]
* secrets/pki: Allow other_sans in sign-intermediate and sign-verbatim [[GH-13958](https://github.com/hashicorp/vault/pull/13958)]
* secrets/pki: Calculate the Subject Key Identifier as suggested in [RFC 5280, Section 4.2.1.2](https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.2). [[GH-11218](https://github.com/hashicorp/vault/pull/11218)]
* secrets/pki: Restrict issuance of wildcard certificates via role parameter (`allow_wildcard_certificates`) [[GH-14238](https://github.com/hashicorp/vault/pull/14238)]
* secrets/pki: Return complete chain (in `ca_chain` field) on calls to `pki/cert/ca_chain` [[GH-13935](https://github.com/hashicorp/vault/pull/13935)]
* secrets/pki: Use application/pem-certificate-chain for PEM certificates, application/x-pem-file for PEM CRLs [[GH-13927](https://github.com/hashicorp/vault/pull/13927)]
* secrets/pki: select appropriate signature algorithm for ECDSA signature on certificates. [[GH-11216](https://github.com/hashicorp/vault/pull/11216)]
* secrets/ssh: Add support for generating non-RSA SSH CAs [[GH-14008](https://github.com/hashicorp/vault/pull/14008)]
* secrets/ssh: Allow specifying multiple approved key lengths for a single algorithm [[GH-13991](https://github.com/hashicorp/vault/pull/13991)]
* secrets/ssh: Use secure default for algorithm signer (rsa-sha2-256) with RSA SSH CA keys on new roles [[GH-14006](https://github.com/hashicorp/vault/pull/14006)]
* secrets/transit: Don't abort transit encrypt or decrypt batches on single item failure. [[GH-13111](https://github.com/hashicorp/vault/pull/13111)]
* storage/aerospike: Upgrade `aerospike-client-go` to v5.6.0. [[GH-12165](https://github.com/hashicorp/vault/pull/12165)]
* storage/raft: Set InitialMmapSize to 100GB on 64bit architectures [[GH-13178](https://github.com/hashicorp/vault/pull/13178)]
* storage/raft: When using retry_join stanzas, join against all of them in parallel. [[GH-13606](https://github.com/hashicorp/vault/pull/13606)]
* sys/raw: Enhance sys/raw to read and write values that cannot be encoded in json. [[GH-13537](https://github.com/hashicorp/vault/pull/13537)]
* ui: Add support for ECDSA and Ed25519 certificate views [[GH-13894](https://github.com/hashicorp/vault/pull/13894)]
* ui: Add version diff view for KV V2 [[GH-13000](https://github.com/hashicorp/vault/pull/13000)]
* ui: Added client side paging for namespace list view [[GH-13195](https://github.com/hashicorp/vault/pull/13195)]
* ui: Adds flight icons to UI [[GH-12976](https://github.com/hashicorp/vault/pull/12976)]
* ui: Adds multi-factor authentication support [[GH-14049](https://github.com/hashicorp/vault/pull/14049)]
* ui: Allow static role credential rotation in Database secrets engines [[GH-14268](https://github.com/hashicorp/vault/pull/14268)]
* ui: Display badge for all versions in secrets engine header [[GH-13015](https://github.com/hashicorp/vault/pull/13015)]
* ui: Swap browser localStorage in favor of sessionStorage [[GH-14054](https://github.com/hashicorp/vault/pull/14054)]
* ui: The integrated web terminal now accepts both `-f` and `--force` as aliases
for `-force` for the `write` command. [[GH-13683](https://github.com/hashicorp/vault/pull/13683)]
* ui: Transform advanced templating with encode/decode format support [[GH-13908](https://github.com/hashicorp/vault/pull/13908)]
* ui: Updates ember blueprints to glimmer components [[GH-13149](https://github.com/hashicorp/vault/pull/13149)]
* ui: customizes empty state messages for transit and transform [[GH-13090](https://github.com/hashicorp/vault/pull/13090)]

BUG FIXES:

* Fixed bug where auth method only considers system-identity when multiple identities are available. [#50](https://github.com/hashicorp/vault-plugin-auth-azure/pull/50) [[GH-14138](https://github.com/hashicorp/vault/pull/14138)]
* activity log (enterprise): allow partial monthly client count to be accessed from namespaces [[GH-13086](https://github.com/hashicorp/vault/pull/13086)]
* agent: Fixes bug where vault agent is unaware of the namespace in the config when wrapping token
* api/client: Fixes an issue where the `replicateStateStore` was being set to `nil` upon consecutive calls to `client.SetReadYourWrites(true)`. [[GH-13486](https://github.com/hashicorp/vault/pull/13486)]
* auth/approle: Fix regression where unset cidrlist is returned as nil instead of zero-length array. [[GH-13235](https://github.com/hashicorp/vault/pull/13235)]
* auth/approle: Fix wrapping of nil errors in `login` endpoint [[GH-14107](https://github.com/hashicorp/vault/pull/14107)]
* auth/github: Use the Organization ID instead of the Organization name to verify the org membership. [[GH-13332](https://github.com/hashicorp/vault/pull/13332)]
* auth/kubernetes: Properly handle the migration of role storage entries containing an empty `alias_name_source` [[GH-13925](https://github.com/hashicorp/vault/pull/13925)]
* auth/kubernetes: ensure valid entity alias names created for projected volume tokens [[GH-14144](https://github.com/hashicorp/vault/pull/14144)]
* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13492](https://github.com/hashicorp/vault/pull/13492)]
* cli: Fix using kv patch with older server versions that don't support HTTP PATCH. [[GH-13615](https://github.com/hashicorp/vault/pull/13615)]
* core (enterprise): Fix a data race in logshipper.
* core (enterprise): Workaround AWS CloudHSM v5 SDK issue not allowing read-only sessions
* core/api: Fix overwriting of request headers when using JSONMergePatch. [[GH-14222](https://github.com/hashicorp/vault/pull/14222)]
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13093](https://github.com/hashicorp/vault/pull/13093)]
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13476](https://github.com/hashicorp/vault/pull/13476)]
* core/token: Fix null token panic from 'v1/auth/token/' endpoints and return proper error response. [[GH-13233](https://github.com/hashicorp/vault/pull/13233)]
* core/token: Fix null token_type panic resulting from 'v1/auth/token/roles/{role_name}' endpoint [[GH-13236](https://github.com/hashicorp/vault/pull/13236)]
* core: Fix warnings logged on perf standbys re stored versions [[GH-13042](https://github.com/hashicorp/vault/pull/13042)]
* core: `-output-curl-string` now properly sets cURL options for client and CA
certificates. [[GH-13660](https://github.com/hashicorp/vault/pull/13660)]
* core: add support for go-sockaddr templates in the top-level cluster_addr field [[GH-13678](https://github.com/hashicorp/vault/pull/13678)]
* core: authentication to "login" endpoint for non-existent mount path returns permission denied with status code 403 [[GH-13162](https://github.com/hashicorp/vault/pull/13162)]
* core: revert some unintentionally downgraded dependencies from 1.9.0-rc1 [[GH-13168](https://github.com/hashicorp/vault/pull/13168)]
* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* http:Fix /sys/monitor endpoint returning streaming not supported [[GH-13200](https://github.com/hashicorp/vault/pull/13200)]
* identity/oidc: Adds support for port-agnostic validation of loopback IP redirect URIs. [[GH-13871](https://github.com/hashicorp/vault/pull/13871)]
* identity/oidc: Check for a nil signing key on rotation to prevent panics. [[GH-13716](https://github.com/hashicorp/vault/pull/13716)]
* identity/oidc: Fixes inherited group membership when evaluating client assignments [[GH-14013](https://github.com/hashicorp/vault/pull/14013)]
* identity/oidc: Fixes potential write to readonly storage on performance secondary clusters during key rotation [[GH-14426](https://github.com/hashicorp/vault/pull/14426)]
* identity/oidc: Make the `nonce` parameter optional for the Authorization Endpoint of OIDC providers. [[GH-13231](https://github.com/hashicorp/vault/pull/13231)]
* identity/token: Fixes a bug where duplicate public keys could appear in the .well-known JWKS [[GH-14543](https://github.com/hashicorp/vault/pull/14543)]
* identity: Fix possible nil pointer dereference. [[GH-13318](https://github.com/hashicorp/vault/pull/13318)]
* identity: Fix regression preventing startup when aliases were created pre-1.9. [[GH-13169](https://github.com/hashicorp/vault/pull/13169)]
* identity: Fixes a panic in the OIDC key rotation due to a missing nil check. [[GH-13298](https://github.com/hashicorp/vault/pull/13298)]
* kmip (enterprise): Fix locate by name operations fail to find key after a rekey operation.
* licensing (enterprise): Revert accidental inclusion of the TDE feature from the `prem` build.
* metrics/autosnapshots (enterprise) : Fix bug that could cause
vault.autosnapshots.save.errors to not be incremented when there is an
autosnapshot save error.
* physical/mysql: Create table with wider `vault_key` column when initializing database tables. [[GH-14231](https://github.com/hashicorp/vault/pull/14231)]
* plugin/couchbase: Fix an issue in which the locking patterns did not allow parallel requests. [[GH-13033](https://github.com/hashicorp/vault/pull/13033)]
* replication (enterprise): When using encrypted secondary tokens, only clear the
private key after a successful connection to the primary cluster
* sdk/framework: Generate proper OpenAPI specs for path patterns that use an alternation as the root. [[GH-13487](https://github.com/hashicorp/vault/pull/13487)]
* sdk/helper/ldaputil: properly escape a trailing escape character to prevent panics. [[GH-13452](https://github.com/hashicorp/vault/pull/13452)]
* sdk/queue: move lock before length check to prevent panics. [[GH-13146](https://github.com/hashicorp/vault/pull/13146)]
* sdk: Fixes OpenAPI to distinguish between paths that can do only List, or both List and Read. [[GH-13643](https://github.com/hashicorp/vault/pull/13643)]
* secrets/azure: Fixed bug where Azure environment did not change Graph URL [[GH-13973](https://github.com/hashicorp/vault/pull/13973)]
* secrets/azure: Fixes service principal generation when assigning roles that have [DataActions](https://docs.microsoft.com/en-us/azure/role-based-access-control/role-definitions#dataactions). [[GH-13277](https://github.com/hashicorp/vault/pull/13277)]
* secrets/azure: Fixes the [rotate root](https://www.vaultproject.io/api-docs/secret/azure#rotate-root) 
operation for upgraded configurations with a `root_password_ttl` of zero. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]
* secrets/database/cassandra: change connect_timeout to 5s as documentation says [[GH-12443](https://github.com/hashicorp/vault/pull/12443)]
* secrets/database/mssql: Accept a boolean for `contained_db`, rather than just a string. [[GH-13469](https://github.com/hashicorp/vault/pull/13469)]
* secrets/gcp: Fixed bug where error was not reported for invalid bindings [[GH-13974](https://github.com/hashicorp/vault/pull/13974)]
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13548](https://github.com/hashicorp/vault/pull/13548)]
* secrets/openldap: Fix panic from nil logger in backend [[GH-14171](https://github.com/hashicorp/vault/pull/14171)]
* secrets/pki: Default value for key_bits changed to 0, enabling key_type=ec key generation with default value [[GH-13080](https://github.com/hashicorp/vault/pull/13080)]
* secrets/pki: Fix issuance of wildcard certificates matching glob patterns [[GH-14235](https://github.com/hashicorp/vault/pull/14235)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-13759](https://github.com/hashicorp/vault/pull/13759)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-2456](https://github.com/hashicorp/vault/pull/2456)]
* secrets/pki: Fixes around NIST P-curve signature hash length, default value for signature_bits changed to 0. [[GH-12872](https://github.com/hashicorp/vault/pull/12872)]
* secrets/pki: Recognize ed25519 when requesting a response in PKCS8 format [[GH-13257](https://github.com/hashicorp/vault/pull/13257)]
* secrets/pki: Skip signature bits validation for ed25519 curve key type [[GH-13254](https://github.com/hashicorp/vault/pull/13254)]
* secrets/transit: Ensure that Vault does not panic for invalid nonce size when we aren't in convergent encryption mode. [[GH-13690](https://github.com/hashicorp/vault/pull/13690)]
* secrets/transit: Return an error if any required parameter is missing. [[GH-14074](https://github.com/hashicorp/vault/pull/14074)]
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* storage/raft: Fix issues allowing invalid nodes to become leadership candidates. [[GH-13703](https://github.com/hashicorp/vault/pull/13703)]
* storage/raft: Fix regression in 1.9.0-rc1 that changed how time is represented in Raft logs; this prevented using a raft db created pre-1.9. [[GH-13165](https://github.com/hashicorp/vault/pull/13165)]
* storage/raft: On linux, use map_populate for bolt files to improve startup time. [[GH-13573](https://github.com/hashicorp/vault/pull/13573)]
* storage/raft: Units for bolt metrics now given in milliseconds instead of nanoseconds [[GH-13749](https://github.com/hashicorp/vault/pull/13749)]
* ui: Adds pagination to auth methods list view [[GH-13054](https://github.com/hashicorp/vault/pull/13054)]
* ui: Do not show verify connection value on database connection config page [[GH-13152](https://github.com/hashicorp/vault/pull/13152)]
* ui: Fix client count current month data not showing unless monthly history data exists [[GH-13396](https://github.com/hashicorp/vault/pull/13396)]
* ui: Fix default TTL display and set on database role [[GH-14224](https://github.com/hashicorp/vault/pull/14224)]
* ui: Fix incorrect validity message on transit secrets engine [[GH-14233](https://github.com/hashicorp/vault/pull/14233)]
* ui: Fix issue where UI incorrectly handled API errors when mounting backends [[GH-14551](https://github.com/hashicorp/vault/pull/14551)]
* ui: Fix kv engine access bug [[GH-13872](https://github.com/hashicorp/vault/pull/13872)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes caching issue on kv new version create [[GH-14489](https://github.com/hashicorp/vault/pull/14489)]
* ui: Fixes displaying empty masked values in PKI engine [[GH-14400](https://github.com/hashicorp/vault/pull/14400)]
* ui: Fixes horizontal bar chart hover issue when filtering namespaces and mounts [[GH-14493](https://github.com/hashicorp/vault/pull/14493)]
* ui: Fixes issue logging out with wrapped token query parameter [[GH-14329](https://github.com/hashicorp/vault/pull/14329)]
* ui: Fixes issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with SearchSelect component not holding focus [[GH-13590](https://github.com/hashicorp/vault/pull/13590)]
* ui: Fixes issue with automate secret deletion value not displaying initially if set in secret metadata edit view [[GH-13177](https://github.com/hashicorp/vault/pull/13177)]
* ui: Fixes issue with correct auth method not selected when logging out from OIDC or JWT methods [[GH-14545](https://github.com/hashicorp/vault/pull/14545)]
* ui: Fixes issue with placeholder not displaying for automatically deleted secrets when deletion time has passed [[GH-13166](https://github.com/hashicorp/vault/pull/13166)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]
* ui: Fixes long secret key names overlapping masked values [[GH-13032](https://github.com/hashicorp/vault/pull/13032)]
* ui: Fixes node-forge error when parsing EC (elliptical curve) certs [[GH-13238](https://github.com/hashicorp/vault/pull/13238)]
* ui: Redirects to managed namespace if incorrect namespace in URL param [[GH-14422](https://github.com/hashicorp/vault/pull/14422)]
* ui: Removes ability to tune token_type for token auth methods [[GH-12904](https://github.com/hashicorp/vault/pull/12904)]
* ui: trigger token renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.9.10
### September 30, 2022

BUG FIXES:

* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]
  
## 1.9.9
### August 31, 2022

CHANGES:

* core: Bump Go version to 1.17.13.

BUG FIXES:

* core (enterprise): Fix some races in merkle index flushing code found in testing
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]

SECURITY:

* identity/entity: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. [[HCSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

## 1.9.8
### July 21, 2022

CHANGES:

* core: Bump Go version to 1.17.12.

IMPROVEMENTS:

* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]

BUG FIXES:

* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* storage/raft (enterprise): Prevent unauthenticated voter status change with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* transform (enterprise): Fix a bug in the handling of nested or unmatched capture groups in FPE transformations.
* ui: Fix issue where metadata tab is hidden even though policy grants access [[GH-15824](https://github.com/hashicorp/vault/pull/15824)]
* ui: Updated `leasId` to `leaseId` in the "Copy Credentials" section of "Generate AWS Credentials" [[GH-15685](https://github.com/hashicorp/vault/pull/15685)]
  
## 1.9.7
### June 10, 2022

CHANGES:

* core: Bump Go version to 1.17.11. [[GH-go-ver-197](https://github.com/hashicorp/vault/pull/go-ver-197)]

IMPROVEMENTS:

* ui: Allow namespace param to be parsed from state queryParam [[GH-15378](https://github.com/hashicorp/vault/pull/15378)]

BUG FIXES:

* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* auth/ldap: The logic for setting the entity alias when `username_as_alias` is set
has been fixed. The previous behavior would make a request to the LDAP server to
get `user_attr` before discarding it and using the username instead. This would
make it impossible for a user to connect if this attribute was missing or had
multiple values, even though it would not be used anyway. This has been fixed
and the username is now used without making superfluous LDAP searches. [[GH-15525](https://github.com/hashicorp/vault/pull/15525)]
* core (enterprise): Fix overcounting of lease count quota usage at startup.
* core/config: Only ask the system about network interfaces when address configs contain a template having the format: {{ ... }} [[GH-15224](https://github.com/hashicorp/vault/pull/15224)]
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: Prevent metrics generation from causing deadlocks. [[GH-15693](https://github.com/hashicorp/vault/pull/15693)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* transform (enterprise): Fix non-overridable column default value causing tokenization tokens to expire prematurely when using the MySQL storage backend.
* ui: Fixes client count timezone bug [[GH-15743](https://github.com/hashicorp/vault/pull/15743)]
* ui: Fixes issue logging in with OIDC from a listed auth mounts tab [[GH-15666](https://github.com/hashicorp/vault/pull/15666)]

## 1.9.6
### April 29, 2022

BUG FIXES:

* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]


## 1.9.5
### April 22, 2022

CHANGES:

* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.17.9. [[GH-15045](https://github.com/hashicorp/vault/pull/15045)]

IMPROVEMENTS:

* auth/ldap: Add username_as_alias configurable to change how aliases are named [[GH-14324](https://github.com/hashicorp/vault/pull/14324)]
* core: Systemd unit file included with the Linux packages now sets the service type to notify. [[GH-14385](https://github.com/hashicorp/vault/pull/14385)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer
* website/docs: added a link to an Enigma secret plugin. [[GH-14389](https://github.com/hashicorp/vault/pull/14389)]

BUG FIXES:

* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* api: Respect increment value in grace period calculations in LifetimeWatcher [[GH-14836](https://github.com/hashicorp/vault/pull/14836)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* core (enterprise): Allow local alias create RPCs to persist alias metadata
* core/metrics: Fix incorrect table size metric for local mounts [[GH-14755](https://github.com/hashicorp/vault/pull/14755)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* identity/token: Fixes a bug where duplicate public keys could appear in the .well-known JWKS [[GH-14543](https://github.com/hashicorp/vault/pull/14543)]
* metrics/autosnapshots (enterprise) : Fix bug that could cause
vault.autosnapshots.save.errors to not be incremented when there is an
autosnapshot save error.
* replication (enterprise): fix panic due to missing entity during invalidation of local aliases. [[GH-14622](https://github.com/hashicorp/vault/pull/14622)]
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not excepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fix issue where UI incorrectly handled API errors when mounting backends [[GH-14551](https://github.com/hashicorp/vault/pull/14551)]
* ui: Fixes caching issue on kv new version create [[GH-14489](https://github.com/hashicorp/vault/pull/14489)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging out with wrapped token query parameter [[GH-14329](https://github.com/hashicorp/vault/pull/14329)]
* ui: Fixes issue with correct auth method not selected when logging out from OIDC or JWT methods [[GH-14545](https://github.com/hashicorp/vault/pull/14545)]
* ui: Redirects to managed namespace if incorrect namespace in URL param [[GH-14422](https://github.com/hashicorp/vault/pull/14422)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]


## 1.9.4
### March 3, 2022

SECURITY:
* secrets/pki: Vault and Vault Enterprise (“Vault”) allowed the PKI secrets engine under certain configurations to issue wildcard certificates to authorized users for a specified domain, even if the PKI role policy attribute allow_subdomains is set to false. This vulnerability, CVE-2022-25243, was fixed in Vault 1.8.9 and 1.9.4.
* transform (enterprise): Vault Enterprise (“Vault”) clusters using the tokenization transform feature can expose the tokenization key through the tokenization key configuration endpoint to authorized operators with read permissions on this endpoint. This vulnerability, CVE-2022-25244, was fixed in Vault Enterprise 1.7.10, 1.8.9, and 1.9.4.

CHANGES:

* secrets/azure: Changes the configuration parameter `use_microsoft_graph_api` to use the Microsoft
Graph API by default. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]

IMPROVEMENTS:

* core: Bump Go version to 1.17.7. [[GH-14232](https://github.com/hashicorp/vault/pull/14232)]
* secrets/pki: Restrict issuance of wildcard certificates via role parameter (`allow_wildcard_certificates`) [[GH-14238](https://github.com/hashicorp/vault/pull/14238)]

BUG FIXES:

* Fixed bug where auth method only considers system-identity when multiple identities are available. [#50](https://github.com/hashicorp/vault-plugin-auth-azure/pull/50) [[GH-14138](https://github.com/hashicorp/vault/pull/14138)]
* auth/kubernetes: Properly handle the migration of role storage entries containing an empty `alias_name_source` [[GH-13925](https://github.com/hashicorp/vault/pull/13925)]
* auth/kubernetes: ensure valid entity alias names created for projected volume tokens [[GH-14144](https://github.com/hashicorp/vault/pull/14144)]
* identity/oidc: Adds support for port-agnostic validation of loopback IP redirect URIs. [[GH-13871](https://github.com/hashicorp/vault/pull/13871)]
* identity/oidc: Fixes inherited group membership when evaluating client assignments [[GH-14013](https://github.com/hashicorp/vault/pull/14013)]
* secrets/azure: Fixed bug where Azure environment did not change Graph URL [[GH-13973](https://github.com/hashicorp/vault/pull/13973)]
* secrets/azure: Fixes the [rotate root](https://www.vaultproject.io/api-docs/secret/azure#rotate-root)
operation for upgraded configurations with a `root_password_ttl` of zero. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]
* secrets/gcp: Fixed bug where error was not reported for invalid bindings [[GH-13974](https://github.com/hashicorp/vault/pull/13974)]
* secrets/openldap: Fix panic from nil logger in backend [[GH-14171](https://github.com/hashicorp/vault/pull/14171)]
* secrets/pki: Fix issuance of wildcard certificates matching glob patterns [[GH-14235](https://github.com/hashicorp/vault/pull/14235)]
* storage/raft: Fix issues allowing invalid nodes to become leadership candidates. [[GH-13703](https://github.com/hashicorp/vault/pull/13703)]
* ui: Fix default TTL display and set on database role [[GH-14224](https://github.com/hashicorp/vault/pull/14224)]
* ui: Fix incorrect validity message on transit secrets engine [[GH-14233](https://github.com/hashicorp/vault/pull/14233)]
* ui: Fix kv engine access bug [[GH-13872](https://github.com/hashicorp/vault/pull/13872)]
* ui: Fix issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Trigger background token self-renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.9.3
### January 27, 2022

IMPROVEMENTS:

* auth/kubernetes: Added support for dynamically reloading short-lived tokens for better Kubernetes 1.21+ compatibility [[GH-13698](https://github.com/hashicorp/vault/pull/13698)]
* auth/ldap: Add username to alias metadata [[GH-13669](https://github.com/hashicorp/vault/pull/13669)]
* core/identity: Support updating an alias' `custom_metadata` to be empty. [[GH-13395](https://github.com/hashicorp/vault/pull/13395)]
* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]
* http (enterprise): Serve /sys/license/status endpoint within namespaces

BUG FIXES:

* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13492](https://github.com/hashicorp/vault/pull/13492)]
* cli: Fix using kv patch with older server versions that don't support HTTP PATCH. [[GH-13615](https://github.com/hashicorp/vault/pull/13615)]
* core (enterprise): Workaround AWS CloudHSM v5 SDK issue not allowing read-only sessions
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13476](https://github.com/hashicorp/vault/pull/13476)]
* core: add support for go-sockaddr templates in the top-level cluster_addr field [[GH-13678](https://github.com/hashicorp/vault/pull/13678)]
* identity/oidc: Check for a nil signing key on rotation to prevent panics. [[GH-13716](https://github.com/hashicorp/vault/pull/13716)]
* kmip (enterprise): Fix locate by name operations fail to find key after a rekey operation.
* secrets/database/mssql: Accept a boolean for `contained_db`, rather than just a string. [[GH-13469](https://github.com/hashicorp/vault/pull/13469)]
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13548](https://github.com/hashicorp/vault/pull/13548)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-13759](https://github.com/hashicorp/vault/pull/13759)]
* storage/raft: On linux, use map_populate for bolt files to improve startup time. [[GH-13573](https://github.com/hashicorp/vault/pull/13573)]
* storage/raft: Units for bolt metrics now given in milliseconds instead of nanoseconds [[GH-13749](https://github.com/hashicorp/vault/pull/13749)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]

## 1.9.2
### December 21, 2021

CHANGES:

* go: Update go version to 1.17.5 [[GH-13408](https://github.com/hashicorp/vault/pull/13408)]

IMPROVEMENTS:

* auth/jwt: The Authorization Code flow makes use of the Proof Key for Code Exchange (PKCE) extension. [[GH-13365](https://github.com/hashicorp/vault/pull/13365)]

BUG FIXES:

* ui: Fix client count current month data not showing unless monthly history data exists [[GH-13396](https://github.com/hashicorp/vault/pull/13396)]

## 1.9.1
### December 9, 2021

SECURITY:

* storage/raft: Integrated Storage backend could be caused to crash by an authenticated user with write permissions to the KV secrets engine. This vulnerability, CVE-2021-45042, was fixed in Vault 1.7.7, 1.8.6, and 1.9.1.

IMPROVEMENTS:

* storage/aerospike: Upgrade `aerospike-client-go` to v5.6.0. [[GH-12165](https://github.com/hashicorp/vault/pull/12165)]

BUG FIXES:

* auth/approle: Fix regression where unset cidrlist is returned as nil instead of zero-length array. [[GH-13235](https://github.com/hashicorp/vault/pull/13235)]
* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* http:Fix /sys/monitor endpoint returning streaming not supported [[GH-13200](https://github.com/hashicorp/vault/pull/13200)]
* identity/oidc: Make the `nonce` parameter optional for the Authorization Endpoint of OIDC providers. [[GH-13231](https://github.com/hashicorp/vault/pull/13231)]
* identity: Fixes a panic in the OIDC key rotation due to a missing nil check. [[GH-13298](https://github.com/hashicorp/vault/pull/13298)]
* sdk/queue: move lock before length check to prevent panics. [[GH-13146](https://github.com/hashicorp/vault/pull/13146)]
* secrets/azure: Fixes service principal generation when assigning roles that have [DataActions](https://docs.microsoft.com/en-us/azure/role-based-access-control/role-definitions#dataactions). [[GH-13277](https://github.com/hashicorp/vault/pull/13277)]
* secrets/pki: Recognize ed25519 when requesting a response in PKCS8 format [[GH-13257](https://github.com/hashicorp/vault/pull/13257)]
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* ui: Do not show verify connection value on database connection config page [[GH-13152](https://github.com/hashicorp/vault/pull/13152)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with automate secret deletion value not displaying initially if set in secret metadata edit view [[GH-13177](https://github.com/hashicorp/vault/pull/13177)]
* ui: Fixes issue with placeholder not displaying for automatically deleted secrets when deletion time has passed [[GH-13166](https://github.com/hashicorp/vault/pull/13166)]
* ui: Fixes node-forge error when parsing EC (elliptical curve) certs [[GH-13238](https://github.com/hashicorp/vault/pull/13238)]

## 1.9.0
### November 17, 2021

CHANGES:

* auth/kubernetes: `disable_iss_validation` defaults to true. [#127](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/127) [[GH-12975](https://github.com/hashicorp/vault/pull/12975)]
* expiration: VAULT_16_REVOKE_PERMITPOOL environment variable has been removed. [[GH-12888](https://github.com/hashicorp/vault/pull/12888)]
* expiration: VAULT_LEASE_USE_LEGACY_REVOCATION_STRATEGY environment variable has
been removed. [[GH-12888](https://github.com/hashicorp/vault/pull/12888)]
* go: Update go version to 1.17.2
* secrets/ssh: Roles with empty allowed_extensions will now forbid end-users
specifying extensions when requesting ssh key signing. Update roles setting
allowed_extensions to `*` to permit any extension to be specified by an end-user. [[GH-12847](https://github.com/hashicorp/vault/pull/12847)]

FEATURES:

* **Customizable HTTP Headers**: Add support to define custom HTTP headers for root path (`/`) and also on API endpoints (`/v1/*`) [[GH-12485](https://github.com/hashicorp/vault/pull/12485)]
* **Deduplicate Token With Entities in Activity Log**: Vault tokens without entities are now tracked with client IDs and deduplicated in the Activity Log [[GH-12820](https://github.com/hashicorp/vault/pull/12820)]
* **Elasticsearch Database UI**: The UI now supports adding and editing Elasticsearch connections in the database secret engine. [[GH-12672](https://github.com/hashicorp/vault/pull/12672)]
* **KV Custom Metadata**: Add ability in kv-v2 to specify version-agnostic custom key metadata via the
metadata endpoint. The data will be present in responses made to the data endpoint independent of the
calling token's `read` access to the metadata endpoint. [[GH-12907](https://github.com/hashicorp/vault/pull/12907)]
* **KV patch (Tech Preview)**: Add partial update support for the `/<mount>/data/:path` kv-v2
endpoint through HTTP `PATCH`.  A new `patch` ACL capability has been added and
is required to make such requests. [[GH-12687](https://github.com/hashicorp/vault/pull/12687)]
* **Key Management Secrets Engine (Enterprise)**: Adds support for distributing and managing keys in GCP Cloud KMS.
* **Local Auth Mount Entities (enterprise)**: Logins on `local` auth mounts will
generate identity entities for the tokens issued. The aliases of the entity
resulting from local auth mounts (local-aliases), will be scoped by the cluster.
This means that the local-aliases will never leave the geographical boundary of
the cluster where they were issued. This is something to be mindful about for
those who have implemented local auth mounts for complying with GDPR guidelines.
* **Namespaces (Enterprise)**: Adds support for locking Vault API for particular namespaces.
* **OIDC Identity Provider (Tech Preview)**: Adds support for Vault to be an OpenID Connect (OIDC) provider. [[GH-12932](https://github.com/hashicorp/vault/pull/12932)]
* **Oracle Database UI**: The UI now supports adding and editing Oracle connections in the database secret engine. [[GH-12752](https://github.com/hashicorp/vault/pull/12752)]
* **Postgres Database UI**: The UI now supports adding and editing Postgres connections in the database secret engine. [[GH-12945](https://github.com/hashicorp/vault/pull/12945)]

SECURITY:

* core/identity: A Vault user with write permission to an entity alias ID sharing a mount accessor with another user may acquire this other user’s policies by merging their identities. This vulnerability, CVE-2021-41802, was fixed in Vault and Vault Enterprise 1.7.5, 1.8.4, and 1.9.0.
* core/identity: Templated ACL policies would always match the first-created entity alias if multiple entity aliases existed for a specified entity and mount combination, potentially resulting in incorrect policy enforcement. This vulnerability, CVE-2021-43998, was fixed in Vault and Vault Enterprise 1.7.6, 1.8.5, and 1.9.0.

IMPROVEMENTS:

* agent/cache: Process persistent cache leases in dependency order during restore to ensure child leases are always correctly restored [[GH-12843](https://github.com/hashicorp/vault/pull/12843)]
* agent/cache: Use an in-process listener between consul-template and vault-agent when caching is enabled and either templates or a listener is defined [[GH-12762](https://github.com/hashicorp/vault/pull/12762)]
* agent/cache: tolerate partial restore failure from persistent cache [[GH-12718](https://github.com/hashicorp/vault/pull/12718)]
* agent/template: add support for new 'writeToFile' template function [[GH-12505](https://github.com/hashicorp/vault/pull/12505)]
* api: Add configuration option for ensuring isolated read-after-write semantics for all Client requests. [[GH-12814](https://github.com/hashicorp/vault/pull/12814)]
* api: adds native Login method to Go client module with different auth method interfaces to support easier authentication [[GH-12796](https://github.com/hashicorp/vault/pull/12796)]
* api: Move mergeStates and other required utils from agent to api module [[GH-12731](https://github.com/hashicorp/vault/pull/12731)]
* api: Support VAULT_HTTP_PROXY environment variable to allow overriding the Vault client's HTTP proxy [[GH-12582](https://github.com/hashicorp/vault/pull/12582)]
* auth/approle: The `role/:name/secret-id-accessor/lookup` endpoint now returns a 404 status code when the `secret_id_accessor` cannot be found [[GH-12788](https://github.com/hashicorp/vault/pull/12788)]
* auth/approle: expose secret_id_accessor as WrappedAccessor when creating wrapped secret-id. [[GH-12425](https://github.com/hashicorp/vault/pull/12425)]
* auth/aws: add profile support for AWS credentials when using the AWS auth method [[GH-12621](https://github.com/hashicorp/vault/pull/12621)]
* auth/kubernetes: validate JWT against the provided role on alias look ahead operations [[GH-12688](https://github.com/hashicorp/vault/pull/12688)]
* auth/kubernetes: Add ability to configure entity alias names based on the serviceaccount's namespace and name. [#110](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/110) [#112](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/112) [[GH-12633](https://github.com/hashicorp/vault/pull/12633)]
* auth/ldap: include support for an optional user filter field when searching for users [[GH-11000](https://github.com/hashicorp/vault/pull/11000)]
* auth/oidc: Adds the `skip_browser` CLI option to allow users to skip opening the default browser during the authentication flow. [[GH-12876](https://github.com/hashicorp/vault/pull/12876)]
* auth/okta: Send x-forwarded-for in Okta Push Factor request [[GH-12320](https://github.com/hashicorp/vault/pull/12320)]
* auth/token: Add `allowed_policies_glob` and `disallowed_policies_glob` fields to token roles to allow glob matching of policies [[GH-7277](https://github.com/hashicorp/vault/pull/7277)]
* cli: Operator diagnose now tests for missing or partial telemetry configurations. [[GH-12802](https://github.com/hashicorp/vault/pull/12802)]
* cli: add new http option : -header which enable sending arbitrary headers with the cli [[GH-12508](https://github.com/hashicorp/vault/pull/12508)]
* command: operator generate-root -decode: allow passing encoded token via stdin [[GH-12881](https://github.com/hashicorp/vault/pull/12881)]
* core/token: Return the token_no_default_policy config on token role read if set [[GH-12565](https://github.com/hashicorp/vault/pull/12565)]
* core: Add support for go-sockaddr templated addresses in config. [[GH-9109](https://github.com/hashicorp/vault/pull/9109)]
* core: adds custom_metadata field for aliases [[GH-12502](https://github.com/hashicorp/vault/pull/12502)]
* core: Update Oracle Cloud library to enable seal integration with the uk-gov-london-1 region [[GH-12724](https://github.com/hashicorp/vault/pull/12724)]
* core: Update github.com/ulikunitz/xz to fix security vulnerability GHSA-25xm-hr59-7c27. [[GH-12253](https://github.com/hashicorp/vault/pull/12253)]
* core: Upgrade github.com/gogo/protobuf [[GH-12255](https://github.com/hashicorp/vault/pull/12255)]
* core: build with Go 1.17, and mitigate a breaking change they made that could impact how approle and ssh interpret IPs/CIDRs [[GH-12868](https://github.com/hashicorp/vault/pull/12868)]
* core: observe the client counts broken down by namespace for partial month client count [[GH-12393](https://github.com/hashicorp/vault/pull/12393)]
* core: Artifact builds will now only run on merges to the release branches or to `main`
* core: The [dockerfile](https://github.com/hashicorp/vault/blob/main/Dockerfile) that is used to build the vault docker image available at [hashicorp/vault](https://hub.docker.com/repository/docker/hashicorp/vault) now lives in the root of this repo, and the entrypoint is available under [.release/docker/docker-entrypoint.sh](https://github.com/hashicorp/vault/blob/main/.release/docker/docker-entrypoint.sh)
* core: The vault linux packaging service configs and pre/post install scripts are now available under [.release/linux](https://github.com/hashicorp/vault/blob/main/.release/linux)
* core: Vault linux packages are now available for all supported linux architectures including arm, arm64, 386, and amd64
* db/cassandra: make the connect_timeout config option actually apply to connection timeouts, in addition to non-connection operations [[GH-12903](https://github.com/hashicorp/vault/pull/12903)]
* identity/token: Only return keys from the `.well-known/keys` endpoint that are being used by roles to sign/verify tokens. [[GH-12780](https://github.com/hashicorp/vault/pull/12780)]
* identity: fix issue where Cache-Control header causes stampede of requests for JWKS keys [[GH-12414](https://github.com/hashicorp/vault/pull/12414)]
* physical/etcd: Upgrade etcd3 client to v3.5.0 and etcd2 to v2.305.0. [[GH-11980](https://github.com/hashicorp/vault/pull/11980)]
* pki: adds signature_bits field to customize signature algorithm on CAs and certs signed by Vault [[GH-11245](https://github.com/hashicorp/vault/pull/11245)]
* plugin: update the couchbase gocb version in the couchbase plugin [[GH-12483](https://github.com/hashicorp/vault/pull/12483)]
* replication (enterprise): Add merkle.flushDirty.num_pages_outstanding metric which specifies number of
outstanding dirty pages that were not flushed. [[GH-2093](https://github.com/hashicorp/vault/pull/2093)]
* sdk/framework: The '+' wildcard is now supported for parameterizing unauthenticated paths. [[GH-12668](https://github.com/hashicorp/vault/pull/12668)]
* secrets/aws: Add conditional template that allows custom usernames for both STS and IAM cases [[GH-12185](https://github.com/hashicorp/vault/pull/12185)]
* secrets/azure: Adds support for rotate-root. [#70](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/70) [[GH-13034](https://github.com/hashicorp/vault/pull/13034)]
* secrets/azure: Adds support for using Microsoft Graph API since Azure Active Directory API is being removed in 2022. [#67](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/67) [[GH-12629](https://github.com/hashicorp/vault/pull/12629)]
* secrets/database: Update MSSQL dependency github.com/denisenkom/go-mssqldb to v0.11.0 and include support for contained databases in MSSQL plugin [[GH-12839](https://github.com/hashicorp/vault/pull/12839)]
* secrets/pki: Allow signing of self-issued certs with a different signature algorithm. [[GH-12514](https://github.com/hashicorp/vault/pull/12514)]
* secrets/pki: Use entropy augmentation when available when generating root and intermediate CA key material. [[GH-12559](https://github.com/hashicorp/vault/pull/12559)]
* secrets/pki: select appropriate signature algorithm for ECDSA signature on certificates. [[GH-11216](https://github.com/hashicorp/vault/pull/11216)]
* secrets/pki: Support ed25519 as a key for the pki backend [[GH-11780](https://github.com/hashicorp/vault/pull/11780)]
* secrets/rabbitmq: Update dependency github.com/michaelklishin/rabbit-hole to v2 and resolve UserInfo.tags regression from RabbitMQ v3.9 [[GH-12877](https://github.com/hashicorp/vault/pull/12877)]
* secrets/ssh: Let allowed_users template mix templated and non-templated parts. [[GH-10886](https://github.com/hashicorp/vault/pull/10886)]
* secrets/ssh: Use entropy augmentation when available for generation of the signing key. [[GH-12560](https://github.com/hashicorp/vault/pull/12560)]
* serviceregistration: add `external-source: "vault"` metadata value for Consul registration. [[GH-12163](https://github.com/hashicorp/vault/pull/12163)]
* storage/raft: Best-effort handling of cancelled contexts. [[GH-12162](https://github.com/hashicorp/vault/pull/12162)]
* transform (enterprise): Add advanced features for encoding and decoding for Transform FPE
* transform (enterprise): Add a `reference` field to batch items, and propogate it to the response
* ui: Add KV secret search box when no metadata list access. [[GH-12626](https://github.com/hashicorp/vault/pull/12626)]
* ui: Add custom metadata to KV secret engine and metadata to config [[GH-12169](https://github.com/hashicorp/vault/pull/12169)]
* ui: Creates new StatText component [[GH-12295](https://github.com/hashicorp/vault/pull/12295)]
* ui: client count monthly view [[GH-12554](https://github.com/hashicorp/vault/pull/12554)]
* ui: creates bar chart component for displaying client count data by namespace [[GH-12437](https://github.com/hashicorp/vault/pull/12437)]
* ui: Add creation time to KV 2 version history and version view [[GH-12663](https://github.com/hashicorp/vault/pull/12663)]
* ui: Added resize for JSON editor [[GH-12906](https://github.com/hashicorp/vault/pull/12906)] [[GH-12906](https://github.com/hashicorp/vault/pull/12906)]
* ui: Adds warning about white space in KV secret engine. [[GH-12921](https://github.com/hashicorp/vault/pull/12921)]
* ui: Click to copy database static role last rotation value in tooltip [[GH-12890](https://github.com/hashicorp/vault/pull/12890)]
* ui: Filter DB connection attributes so only relevant attrs POST to backend [[GH-12770](https://github.com/hashicorp/vault/pull/12770)]
* ui: Removes empty rows from DB config views [[GH-12819](https://github.com/hashicorp/vault/pull/12819)]
* ui: Standardizes toolbar presentation of destructive actions [[GH-12895](https://github.com/hashicorp/vault/pull/12895)]
* ui: Updates font for table row value fields [[GH-12908](https://github.com/hashicorp/vault/pull/12908)]
* ui: namespace search in client count views [[GH-12577](https://github.com/hashicorp/vault/pull/12577)]
* ui: parse and display pki cert metadata [[GH-12541](https://github.com/hashicorp/vault/pull/12541)]
* ui: replaces Vault's use of elazarl/go-bindata-assetfs in building the UI with Go's native Embed package [[GH-11208](https://github.com/hashicorp/vault/pull/11208)]
* ui: updated client tracking config view [[GH-12422](https://github.com/hashicorp/vault/pull/12422)]

DEPRECATIONS:

* auth/kubernetes: deprecate `disable_iss_validation` and `issuer` configuration fields [#127](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/127) [[GH-12975](https://github.com/hashicorp/vault/pull/12975)]

BUG FIXES:

* activity log (enterprise): allow partial monthly client count to be accessed from namespaces [[GH-13086](https://github.com/hashicorp/vault/pull/13086)]
* agent: Avoid possible `unexpected fault address` panic when using persistent cache. [[GH-12534](https://github.com/hashicorp/vault/pull/12534)]
* api: Fixes storage APIs returning incorrect error when parsing responses [[GH-12338](https://github.com/hashicorp/vault/pull/12338)]
* auth/aws: Fix ec2 auth on instances that have a cert in their PKCS7 signature [[GH-12519](https://github.com/hashicorp/vault/pull/12519)]
* auth/aws: Fixes ec2 login no longer supporting DSA signature verification [[GH-12340](https://github.com/hashicorp/vault/pull/12340)]
* auth/aws: fix config/rotate-root to store new key [[GH-12715](https://github.com/hashicorp/vault/pull/12715)]
* auth/jwt: Fixes OIDC auth from the Vault UI when using `form_post` as the `oidc_response_mode`. [[GH-12265](https://github.com/hashicorp/vault/pull/12265)]
* cli/api: Providing consistency for the use of comma separated parameters in auth/secret enable/tune [[GH-12126](https://github.com/hashicorp/vault/pull/12126)]
* cli: fixes CLI requests when namespace is both provided as argument and part of the path [[GH-12720](https://github.com/hashicorp/vault/pull/12720)]
* cli: fixes CLI requests when namespace is both provided as argument and part of the path [[GH-12911](https://github.com/hashicorp/vault/pull/12911)]
* cli: vault debug now puts newlines after every captured log line. [[GH-12175](https://github.com/hashicorp/vault/pull/12175)]
* core (enterprise): Allow deletion of stored licenses on DR secondary nodes
* core (enterprise): Disallow autogenerated licenses to be used in diagnose even when config is specified
* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Fix data race during perf standby sealing
* core (enterprise): Fixes reading raft auto-snapshot configuration from performance standby node [[GH-12317](https://github.com/hashicorp/vault/pull/12317)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* core (enterprise): namespace header included in responses, Go client uses it when displaying error messages [[GH-12196](https://github.com/hashicorp/vault/pull/12196)]
* core/api: Fix an arm64 bug converting a negative int to an unsigned int [[GH-12372](https://github.com/hashicorp/vault/pull/12372)]
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13093](https://github.com/hashicorp/vault/pull/13093)]
* core/identity: Cleanup alias in the in-memory entity after an alias deletion by ID [[GH-12834](https://github.com/hashicorp/vault/pull/12834)]
* core/identity: Disallow entity alias creation/update if a conflicting alias exists for the target entity and mount combination [[GH-12747](https://github.com/hashicorp/vault/pull/12747)]
* core: Fix a deadlock on HA leadership transfer [[GH-12691](https://github.com/hashicorp/vault/pull/12691)]
* core: Fix warnings logged on perf standbys re stored versions [[GH-13042](https://github.com/hashicorp/vault/pull/13042)]
* core: fix byte printing for diagnose disk checks [[GH-12229](https://github.com/hashicorp/vault/pull/12229)]
* core: revert some unintentionally downgraded dependencies from 1.9.0-rc1 [[GH-13168](https://github.com/hashicorp/vault/pull/13168)]
* database/couchbase: change default template to truncate username at 128 characters [[GH-12301](https://github.com/hashicorp/vault/pull/12301)]
* database/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* http: removed unpublished true from logical_system path, making openapi spec consistent with documentation [[GH-12713](https://github.com/hashicorp/vault/pull/12713)]
* identity/token: Adds missing call to unlock mutex in key deletion error handling [[GH-12916](https://github.com/hashicorp/vault/pull/12916)]
* identity: Fail alias rename if the resulting (name,accessor) exists already [[GH-12473](https://github.com/hashicorp/vault/pull/12473)]
* identity: Fix a panic on arm64 platform when doing identity I/O. [[GH-12371](https://github.com/hashicorp/vault/pull/12371)]
* identity: Fix regression preventing startup when aliases were created pre-1.9. [[GH-13169](https://github.com/hashicorp/vault/pull/13169)]
* identity: dedup from_entity_ids when merging two entities [[GH-10101](https://github.com/hashicorp/vault/pull/10101)]
* identity: disallow creation of role without a key parameter [[GH-12208](https://github.com/hashicorp/vault/pull/12208)]
* identity: do not allow a role's token_ttl to be longer than the signing key's verification_ttl [[GH-12151](https://github.com/hashicorp/vault/pull/12151)]
* identity: merge associated entity groups when merging entities [[GH-10085](https://github.com/hashicorp/vault/pull/10085)]
* identity: suppress duplicate policies on entities [[GH-12812](https://github.com/hashicorp/vault/pull/12812)]
* kmip (enterprise): Fix handling of custom attributes when servicing GetAttributes requests
* kmip (enterprise): Fix handling of invalid role parameters within various vault api calls
* kmip (enterprise): Forward KMIP register operations to the active node
* license: ignore stored terminated license while autoloading is enabled [[GH-2104](https://github.com/hashicorp/vault/pull/2104)]
* licensing (enterprise): Revert accidental inclusion of the TDE feature from the `prem` build.
* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* pki: Fix regression preventing email addresses being used as a common name within certificates [[GH-12716](https://github.com/hashicorp/vault/pull/12716)]
* plugin/couchbase: Fix an issue in which the locking patterns did not allow parallel requests. [[GH-13033](https://github.com/hashicorp/vault/pull/13033)]
* plugin/snowflake: Fixed bug where plugin would crash on 32 bit systems [[GH-12378](https://github.com/hashicorp/vault/pull/12378)]
* raft (enterprise): Fix panic when updating auto-snapshot config
* replication (enterprise): Fix issue where merkle.flushDirty.num_pages metric is not emitted if number
of dirty pages is 0. [[GH-2093](https://github.com/hashicorp/vault/pull/2093)]
* replication (enterprise): Fix merkle.saveCheckpoint.num_dirty metric to accurately specify the number 
of dirty pages in the merkle tree at time of checkpoint creation. [[GH-2093](https://github.com/hashicorp/vault/pull/2093)]
* sdk/database: Fix a DeleteUser error message on the gRPC client. [[GH-12351](https://github.com/hashicorp/vault/pull/12351)]
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/gcp: Fixes a potential panic in the service account policy rollback for rolesets. [[GH-12379](https://github.com/hashicorp/vault/pull/12379)]
* secrets/keymgmt (enterprise): Fix support for Azure Managed HSM Key Vault instances. [[GH-12934](https://github.com/hashicorp/vault/pull/12934)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12600](https://github.com/hashicorp/vault/pull/12600)]
* secrets/transit: Enforce minimum cache size for transit backend and init cache size on transit backend without restart. [[GH-12418](https://github.com/hashicorp/vault/pull/12418)]
* storage/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* storage/raft (enterprise): Ensure that raft autosnapshot backoff retry duration never hits 0s
* storage/raft: Detect incomplete raft snapshots in api.RaftSnapshot(), and thereby in `vault operator raft snapshot save`. [[GH-12388](https://github.com/hashicorp/vault/pull/12388)]
* storage/raft: Fix regression in 1.9.0-rc1 that changed how time is represented in Raft logs; this prevented using a raft db created pre-1.9. [[GH-13165](https://github.com/hashicorp/vault/pull/13165)]
* storage/raft: Support `addr_type=public_v6` in auto-join [[GH-12366](https://github.com/hashicorp/vault/pull/12366)]
* transform (enterprise): Enforce minimum cache size for Transform backend and reset cache size without a restart
* transform (enterprise): Fix an error where the decode response of an expired token is an empty result rather than an error.
* ui: Adds pagination to auth methods list view [[GH-13054](https://github.com/hashicorp/vault/pull/13054)]
* ui: Fix bug where capabilities check on secret-delete-menu was encoding the forward slashes. [[GH-12550](https://github.com/hashicorp/vault/pull/12550)]
* ui: Fix bug where edit role form on auth method is invalid by default [[GH-12646](https://github.com/hashicorp/vault/pull/12646)]
* ui: Fixed api explorer routing bug [[GH-12354](https://github.com/hashicorp/vault/pull/12354)]
* ui: Fixed text overflow in flash messages [[GH-12357](https://github.com/hashicorp/vault/pull/12357)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: Remove spinner after token renew [[GH-12887](https://github.com/hashicorp/vault/pull/12887)]
* ui: Removes ability to tune token_type for token auth methods [[GH-12904](https://github.com/hashicorp/vault/pull/12904)]
* ui: Show day of month instead of day of year in the expiration warning dialog [[GH-11984](https://github.com/hashicorp/vault/pull/11984)]
* ui: fix issue where on MaskedInput on auth methods if tab it would clear the value. [[GH-12409](https://github.com/hashicorp/vault/pull/12409)]
* ui: fix missing navbar items on login to namespace [[GH-12478](https://github.com/hashicorp/vault/pull/12478)]
* ui: update bar chart when model changes [[GH-12622](https://github.com/hashicorp/vault/pull/12622)]
* ui: updating database TTL picker help text. [[GH-12212](https://github.com/hashicorp/vault/pull/12212)]

## 1.8.12
### June 10, 2022

BUG FIXES:

* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* transform (enterprise): Fix non-overridable column default value causing tokenization tokens to expire prematurely when using the MySQL storage backend.

## 1.8.11
### April 29, 2022

BUG FIXES:

* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]

## 1.8.10
### April 22, 2022

CHANGES:

* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.16.15. [[GH-go-ver-1810](https://github.com/hashicorp/vault/pull/go-ver-1810)]

IMPROVEMENTS:

* auth/ldap: Add username_as_alias configurable to change how aliases are named [[GH-14324](https://github.com/hashicorp/vault/pull/14324)]
* core: Systemd unit file included with the Linux packages now sets the service type to notify. [[GH-14385](https://github.com/hashicorp/vault/pull/14385)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer

BUG FIXES:

* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* metrics/autosnapshots (enterprise) : Fix bug that could cause
vault.autosnapshots.save.errors to not be incremented when there is an
autosnapshot save error.
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not excepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging out with wrapped token query parameter [[GH-14329](https://github.com/hashicorp/vault/pull/14329)]
* ui: Fixes issue with correct auth method not selected when logging out from OIDC or JWT methods [[GH-14545](https://github.com/hashicorp/vault/pull/14545)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]


## 1.8.9
### March 3, 2022

* secrets/pki: Vault and Vault Enterprise (“Vault”) allowed the PKI secrets engine under certain configurations to issue wildcard certificates to authorized users for a specified domain, even if the PKI role policy attribute allow_subdomains is set to false. This vulnerability, CVE-2022-25243, was fixed in Vault 1.8.9 and 1.9.4.
* transform (enterprise): Vault Enterprise (“Vault”) clusters using the tokenization transform feature can expose the tokenization key through the tokenization key configuration endpoint to authorized operators with read permissions on this endpoint. This vulnerability, CVE-2022-25244, was fixed in Vault Enterprise 1.7.10, 1.8.9, and 1.9.4.

IMPROVEMENTS:

* secrets/pki: Restrict issuance of wildcard certificates via role parameter (`allow_wildcard_certificates`) [[GH-14238](https://github.com/hashicorp/vault/pull/14238)]

BUG FIXES:

* auth/aws: Fix ec2 auth on instances that have a cert in their PKCS7 signature [[GH-12519](https://github.com/hashicorp/vault/pull/12519)]
* database/mssql: Removed string interpolation on internal queries and replaced them with inline queries using named parameters. [[GH-13799](https://github.com/hashicorp/vault/pull/13799)]
* secrets/openldap: Fix panic from nil logger in backend [[GH-14170](https://github.com/hashicorp/vault/pull/14170)]
* secrets/pki: Fix issuance of wildcard certificates matching glob patterns [[GH-14235](https://github.com/hashicorp/vault/pull/14235)]
* ui: Fix issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Trigger background token self-renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.8.8
### January 27, 2022

IMPROVEMENTS:

* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]

BUG FIXES:

* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13494](https://github.com/hashicorp/vault/pull/13494)]
* core (enterprise): Workaround AWS CloudHSM v5 SDK issue not allowing read-only sessions
* kmip (enterprise): Fix locate by name operations fail to find key after a rekey operation.
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13549](https://github.com/hashicorp/vault/pull/13549)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-13759](https://github.com/hashicorp/vault/pull/13759)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-2456](https://github.com/hashicorp/vault/pull/2456)]
* storage/raft: Fix issues allowing invalid nodes to become leadership candidates. [[GH-13703](https://github.com/hashicorp/vault/pull/13703)]
* storage/raft: On linux, use map_populate for bolt files to improve startup time. [[GH-13573](https://github.com/hashicorp/vault/pull/13573)]
* storage/raft: Units for bolt metrics now given in milliseconds instead of nanoseconds [[GH-13749](https://github.com/hashicorp/vault/pull/13749)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]

## 1.8.7
### December 21, 2021

CHANGES:

* go: Update go version to 1.16.12 [[GH-13422](https://github.com/hashicorp/vault/pull/13422)]

## 1.8.6
### December 9, 2021

CHANGES:

* go: Update go version to 1.16.9 [[GH-13029](https://github.com/hashicorp/vault/pull/13029)]

SECURITY:

* storage/raft: Integrated Storage backend could be caused to crash by an authenticated user with write permissions to the KV secrets engine. This vulnerability, CVE-2021-45042, was fixed in Vault 1.7.7, 1.8.6, and 1.9.1.

BUG FIXES:

* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* ui: Adds pagination to auth methods list view [[GH-13054](https://github.com/hashicorp/vault/pull/13054)]
* ui: Do not show verify connection value on database connection config page [[GH-13152](https://github.com/hashicorp/vault/pull/13152)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]

## 1.8.5
### November 4, 2021

SECURITY:

* core/identity: Templated ACL policies would always match the first-created entity alias if multiple entity aliases existed for a specified entity and mount combination, potentially resulting in incorrect policy enforcement. This vulnerability, CVE-2021-43998, was fixed in Vault and Vault Enterprise 1.7.6, 1.8.5, and 1.9.0.

BUG FIXES:

* auth/aws: fix config/rotate-root to store new key [[GH-12715](https://github.com/hashicorp/vault/pull/12715)]
* core/identity: Cleanup alias in the in-memory entity after an alias deletion by ID [[GH-12834](https://github.com/hashicorp/vault/pull/12834)]
* core/identity: Disallow entity alias creation/update if a conflicting alias exists for the target entity and mount combination [[GH-12747](https://github.com/hashicorp/vault/pull/12747)]
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* identity/token: Adds missing call to unlock mutex in key deletion error handling [[GH-12916](https://github.com/hashicorp/vault/pull/12916)]
* kmip (enterprise): Fix handling of custom attributes when servicing GetAttributes requests
* kmip (enterprise): Fix handling of invalid role parameters within various vault api calls
* kmip (enterprise): Forward KMIP register operations to the active node
* secrets/keymgmt (enterprise): Fix support for Azure Managed HSM Key Vault instances. [[GH-12952](https://github.com/hashicorp/vault/pull/12952)]
* transform (enterprise): Fix an error where the decode response of an expired token is an empty result rather than an error.

## 1.8.4
### 6 October 2021

SECURITY:

* core/identity: A Vault user with write permission to an entity alias ID sharing a mount accessor with another user may acquire this other user’s policies by merging their identities. This vulnerability, CVE-2021-41802, was fixed in Vault and Vault Enterprise 1.7.5 and 1.8.4.

IMPROVEMENTS:

* core: Update Oracle Cloud library to enable seal integration with the uk-gov-london-1 region [[GH-12724](https://github.com/hashicorp/vault/pull/12724)]

BUG FIXES:

* core: Fix a deadlock on HA leadership transfer [[GH-12691](https://github.com/hashicorp/vault/pull/12691)]
* database/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* pki: Fix regression preventing email addresses being used as a common name within certificates [[GH-12716](https://github.com/hashicorp/vault/pull/12716)]
* storage/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* ui: Fix bug where edit role form on auth method is invalid by default [[GH-12646](https://github.com/hashicorp/vault/pull/12646)]

## 1.8.3
### 29 September 2021

IMPROVEMENTS:

* secrets/pki: Allow signing of self-issued certs with a different signature algorithm. [[GH-12514](https://github.com/hashicorp/vault/pull/12514)]

BUG FIXES:

* agent: Avoid possible `unexpected fault address` panic when using persistent cache. [[GH-12534](https://github.com/hashicorp/vault/pull/12534)]
* core (enterprise): Allow deletion of stored licenses on DR secondary nodes
* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* identity: Fail alias rename if the resulting (name,accessor) exists already [[GH-12473](https://github.com/hashicorp/vault/pull/12473)]
* raft (enterprise): Fix panic when updating auto-snapshot config
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12599](https://github.com/hashicorp/vault/pull/12599)]
* secrets/transit: Enforce minimum cache size for transit backend and init cache size on transit backend without restart. [[GH-12418](https://github.com/hashicorp/vault/pull/12418)]
* storage/raft: Detect incomplete raft snapshots in api.RaftSnapshot(), and thereby in `vault operator raft snapshot save`. [[GH-12388](https://github.com/hashicorp/vault/pull/12388)]
* ui: Fix bug where capabilities check on secret-delete-menu was encoding the forward slashes. [[GH-12550](https://github.com/hashicorp/vault/pull/12550)]
* ui: Show day of month instead of day of year in the expiration warning dialog [[GH-11984](https://github.com/hashicorp/vault/pull/11984)]

## 1.8.2
### 26 August 2021

CHANGES:

* Alpine: Docker images for Vault 1.6.6+, 1.7.4+, and 1.8.2+ are built with Alpine 3.14, due to CVE-2021-36159
* go: Update go version to 1.16.7 [[GH-12408](https://github.com/hashicorp/vault/pull/12408)]

BUG FIXES:

* auth/aws: Fixes ec2 login no longer supporting DSA signature verification [[GH-12340](https://github.com/hashicorp/vault/pull/12340)]
* cli: vault debug now puts newlines after every captured log line. [[GH-12175](https://github.com/hashicorp/vault/pull/12175)]
* database/couchbase: change default template to truncate username at 128 characters [[GH-12300](https://github.com/hashicorp/vault/pull/12300)]
* identity: Fix a panic on arm64 platform when doing identity I/O. [[GH-12371](https://github.com/hashicorp/vault/pull/12371)]
* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* plugin/snowflake: Fixed bug where plugin would crash on 32 bit systems [[GH-12378](https://github.com/hashicorp/vault/pull/12378)]
* sdk/database: Fix a DeleteUser error message on the gRPC client. [[GH-12351](https://github.com/hashicorp/vault/pull/12351)]
* secrets/gcp: Fixes a potential panic in the service account policy rollback for rolesets. [[GH-12379](https://github.com/hashicorp/vault/pull/12379)]
* ui: Fixed api explorer routing bug [[GH-12354](https://github.com/hashicorp/vault/pull/12354)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: fix issue where on MaskedInput on auth methods if tab it would clear the value. [[GH-12409](https://github.com/hashicorp/vault/pull/12409)]

## 1.8.1
### August 5th, 2021

CHANGES:

* go: Update go version to 1.16.6 [[GH-12245](https://github.com/hashicorp/vault/pull/12245)]

IMPROVEMENTS:

* serviceregistration: add `external-source: "vault"` metadata value for Consul registration. [[GH-12163](https://github.com/hashicorp/vault/pull/12163)]

BUG FIXES:

* auth/aws: Remove warning stating AWS Token TTL will be capped by the Default Lease TTL. [[GH-12026](https://github.com/hashicorp/vault/pull/12026)]
* auth/jwt: Fixes OIDC auth from the Vault UI when using `form_post` as the `oidc_response_mode`. [[GH-12258](https://github.com/hashicorp/vault/pull/12258)]
* core (enterprise): Disallow autogenerated licenses to be used in diagnose even when config is specified
* core: fix byte printing for diagnose disk checks [[GH-12229](https://github.com/hashicorp/vault/pull/12229)]
* identity: do not allow a role's token_ttl to be longer than the signing key's verification_ttl [[GH-12151](https://github.com/hashicorp/vault/pull/12151)]

## 1.8.0
### July 28th, 2021

CHANGES:

* agent: Errors in the template engine will no longer cause agent to exit unless
explicitly defined to do so. A new configuration parameter,
`exit_on_retry_failure`, within the new top-level stanza, `template_config`, can
be set to `true` in order to cause agent to exit. Note that for agent to exit if
`template.error_on_missing_key` is set to `true`, `exit_on_retry_failure` must
be also set to `true`. Otherwise, the template engine will log an error but then
restart its internal runner. [[GH-11775](https://github.com/hashicorp/vault/pull/11775)]
* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* core (enterprise): License/EULA changes that ensure the presence of a valid HashiCorp license to
start Vault. More information is available in the [Vault License FAQ](https://www.vaultproject.io/docs/enterprise/license/faqs)

FEATURES:

* **GCP Secrets Engine Static Accounts**: Adds ability to use existing service accounts for generation
  of service account keys and access tokens. [[GH-12023](https://github.com/hashicorp/vault/pull/12023)]
* **Key Management Secrets Engine (Enterprise)**: Adds general availability for distributing and managing keys in AWS KMS. [[GH-11958](https://github.com/hashicorp/vault/pull/11958)]
* **License Autoloading (Enterprise)**: Licenses may now be automatically loaded from the environment or disk.
* **MySQL Database UI**: The UI now supports adding and editing MySQL connections in the database secret engine [[GH-11532](https://github.com/hashicorp/vault/pull/11532)]
* **Vault Diagnose**: A new `vault operator` command to detect common issues with vault server setups.

SECURITY:

* storage/raft: When initializing Vault’s Integrated Storage backend, excessively broad filesystem permissions may be set for the underlying Bolt database used by Vault’s Raft implementation. This vulnerability, CVE-2021-38553, was fixed in Vault 1.8.0.
* ui: The Vault UI erroneously cached and exposed user-viewed secrets between authenticated sessions in a single shared browser, if the browser window / tab was not refreshed or closed between logout and a subsequent login. This vulnerability, CVE-2021-38554, was fixed in Vault 1.8.0 and will be addressed in pending 1.7.4 / 1.6.6 releases.

IMPROVEMENTS:

* agent/template: Added static_secret_render_interval to specify how often to fetch non-leased secrets [[GH-11934](https://github.com/hashicorp/vault/pull/11934)]
* agent: Allow Agent auto auth to read symlinked JWT files [[GH-11502](https://github.com/hashicorp/vault/pull/11502)]
* api: Allow a leveled logger to be provided to `api.Client` through `SetLogger`. [[GH-11696](https://github.com/hashicorp/vault/pull/11696)]
* auth/aws: Underlying error included in validation failure message. [[GH-11638](https://github.com/hashicorp/vault/pull/11638)]
* cli/api: Add lease lookup command [[GH-11129](https://github.com/hashicorp/vault/pull/11129)]
* core: Add `prefix_filter` to telemetry config [[GH-12025](https://github.com/hashicorp/vault/pull/12025)]
* core: Add a darwin/arm64 binary release supporting the Apple M1 CPU [[GH-12071](https://github.com/hashicorp/vault/pull/12071)]
* core: Add a small (<1s) exponential backoff to failed TCP listener Accept failures. [[GH-11588](https://github.com/hashicorp/vault/pull/11588)]
* core (enterprise): Add controlled capabilities to control group policy stanza
* core: Add metrics for standby node forwarding. [[GH-11366](https://github.com/hashicorp/vault/pull/11366)]
* core: Add metrics to report if a node is a perf standby, if a node is a dr secondary or primary, and if a node is a perf secondary or primary. [[GH-11472](https://github.com/hashicorp/vault/pull/11472)]
* core: Send notifications to systemd on start, stop, and configuration reload. [[GH-11517](https://github.com/hashicorp/vault/pull/11517)]
* core: add irrevocable lease list and count apis [[GH-11607](https://github.com/hashicorp/vault/pull/11607)]
* core: allow arbitrary length stack traces upon receiving SIGUSR2 (was 32MB) [[GH-11364](https://github.com/hashicorp/vault/pull/11364)]
* core: Improve renew/revoke performance using per-lease locks [[GH-11122](https://github.com/hashicorp/vault/pull/11122)]
* db/cassandra: Added tls_server_name to specify server name for TLS validation [[GH-11820](https://github.com/hashicorp/vault/pull/11820)]
* go: Update to Go 1.16.5 [[GH-11802](https://github.com/hashicorp/vault/pull/11802)]
* replication: Delay evaluation of X-Vault-Index headers until merkle sync completes.
* secrets/rabbitmq: Add ability to customize dynamic usernames [[GH-11899](https://github.com/hashicorp/vault/pull/11899)]
* secrets/ad: Add `rotate-role` endpoint to allow rotations of service accounts. [[GH-11942](https://github.com/hashicorp/vault/pull/11942)]
* secrets/aws: add IAM tagging support for iam_user roles [[GH-10953](https://github.com/hashicorp/vault/pull/10953)]
* secrets/aws: add ability to provide a role session name when generating STS credentials [[GH-11345](https://github.com/hashicorp/vault/pull/11345)]
* secrets/database/elasticsearch: Add ability to customize dynamic usernames [[GH-11957](https://github.com/hashicorp/vault/pull/11957)]
* secrets/database/influxdb: Add ability to customize dynamic usernames [[GH-11796](https://github.com/hashicorp/vault/pull/11796)]
* secrets/database/mongodb: Add ability to customize `SocketTimeout`, `ConnectTimeout`, and `ServerSelectionTimeout` [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database/mongodb: Increased throughput by allowing for multiple request threads to simultaneously update users in MongoDB [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database/mongodbatlas: Adds the ability to customize username generation for dynamic users in MongoDB Atlas. [[GH-11956](https://github.com/hashicorp/vault/pull/11956)]
* secrets/database/redshift: Add ability to customize dynamic usernames [[GH-12016](https://github.com/hashicorp/vault/pull/12016)]
* secrets/database/snowflake: Add ability to customize dynamic usernames [[GH-11997](https://github.com/hashicorp/vault/pull/11997)]
* ssh: add support for templated values in SSH CA DefaultExtensions [[GH-11495](https://github.com/hashicorp/vault/pull/11495)]
* storage/raft: Improve raft batch size selection [[GH-11907](https://github.com/hashicorp/vault/pull/11907)]
* storage/raft: change freelist type to map and set nofreelistsync to true [[GH-11895](https://github.com/hashicorp/vault/pull/11895)]
* storage/raft: Switch to shared raft-boltdb library and add boltdb metrics [[GH-11269](https://github.com/hashicorp/vault/pull/11269)]
* storage/raft: Support autopilot for HA only raft storage. [[GH-11260](https://github.com/hashicorp/vault/pull/11260)]
* storage/raft (enterprise): Enable Autopilot on DR secondary clusters
* ui: Add Validation to KV secret engine [[GH-11785](https://github.com/hashicorp/vault/pull/11785)]
* ui: Add database secret engine support for MSSQL [[GH-11231](https://github.com/hashicorp/vault/pull/11231)]
* ui: Add push notification message when selecting okta auth. [[GH-11442](https://github.com/hashicorp/vault/pull/11442)]
* ui: Add regex validation to Transform Template pattern input [[GH-11586](https://github.com/hashicorp/vault/pull/11586)]
* ui: Add specific error message if unseal fails due to license [[GH-11705](https://github.com/hashicorp/vault/pull/11705)]
* ui: Add validation support for open api form fields [[GH-11963](https://github.com/hashicorp/vault/pull/11963)]
* ui: Added auth method descriptions to UI login page [[GH-11795](https://github.com/hashicorp/vault/pull/11795)]
* ui: JSON fields on database can be cleared on edit [[GH-11708](https://github.com/hashicorp/vault/pull/11708)]
* ui: Obscure secret values on input and displayOnly fields like certificates. [[GH-11284](https://github.com/hashicorp/vault/pull/11284)]
* ui: Redesign of KV 2 Delete toolbar. [[GH-11530](https://github.com/hashicorp/vault/pull/11530)]
* ui: Replace tool partials with components. [[GH-11672](https://github.com/hashicorp/vault/pull/11672)]
* ui: Show description on secret engine list [[GH-11995](https://github.com/hashicorp/vault/pull/11995)]
* ui: Update ember to latest LTS and upgrade UI dependencies [[GH-11447](https://github.com/hashicorp/vault/pull/11447)]
* ui: Update partials to components [[GH-11680](https://github.com/hashicorp/vault/pull/11680)]
* ui: Updated ivy code mirror component for consistency [[GH-11500](https://github.com/hashicorp/vault/pull/11500)]
* ui: Updated node to v14, latest stable build [[GH-12049](https://github.com/hashicorp/vault/pull/12049)]
* ui: Updated search select component styling [[GH-11360](https://github.com/hashicorp/vault/pull/11360)]
* ui: add transform secrets engine to features list [[GH-12003](https://github.com/hashicorp/vault/pull/12003)]
* ui: add validations for duplicate path kv engine [[GH-11878](https://github.com/hashicorp/vault/pull/11878)]
* ui: show site-wide banners for license warnings if applicable [[GH-11759](https://github.com/hashicorp/vault/pull/11759)]
* ui: update license page with relevant autoload info [[GH-11778](https://github.com/hashicorp/vault/pull/11778)]

DEPRECATIONS:

* secrets/gcp: Deprecated the `/gcp/token/:roleset` and `/gcp/key/:roleset` paths for generating
  secrets for rolesets. Use `/gcp/roleset/:roleset/token` and `/gcp/roleset/:roleset/key` instead. [[GH-12023](https://github.com/hashicorp/vault/pull/12023)]

BUG FIXES:

* activity: Omit wrapping tokens and control groups from client counts [[GH-11826](https://github.com/hashicorp/vault/pull/11826)]
* agent/cert: Fix issue where the API client on agent was not honoring certificate
  information from the auto-auth config map on renewals or retries. [[GH-11576](https://github.com/hashicorp/vault/pull/11576)]
* agent/template: fix command shell quoting issue [[GH-11838](https://github.com/hashicorp/vault/pull/11838)]
* agent: Fixed agent templating to use configured tls servername values [[GH-11288](https://github.com/hashicorp/vault/pull/11288)]
* agent: fix timestamp format in log messages from the templating engine [[GH-11838](https://github.com/hashicorp/vault/pull/11838)]
* auth/approle: fixing dereference of nil pointer [[GH-11864](https://github.com/hashicorp/vault/pull/11864)]
* auth/jwt: Updates the [hashicorp/cap](https://github.com/hashicorp/cap) library to `v0.1.0` to
  bring in a verification key caching fix. [[GH-11784](https://github.com/hashicorp/vault/pull/11784)]
* auth/kubernetes: Fix AliasLookahead to correctly extract ServiceAccount UID when using ephemeral JWTs [[GH-12073](https://github.com/hashicorp/vault/pull/12073)]
* auth/ldap: Fix a bug where the LDAP auth method does not return the request_timeout configuration parameter on config read. [[GH-11975](https://github.com/hashicorp/vault/pull/11975)]
* cli: Add support for response wrapping in `vault list` and `vault kv list` with output format other than `table`. [[GH-12031](https://github.com/hashicorp/vault/pull/12031)]
* cli: vault delete and vault kv delete should support the same output options (e.g. -format) as vault write. [[GH-11992](https://github.com/hashicorp/vault/pull/11992)]
* core (enterprise): Fix orphan return value from auth methods executed on performance standby nodes.
* core (enterprise): Fix plugins mounted in namespaces being unable to use password policies [[GH-11596](https://github.com/hashicorp/vault/pull/11596)]
* core (enterprise): serialize access to HSM entropy generation to avoid errors in concurrent key generation.
* core/metrics: Add generic KV mount support for vault.kv.secret.count telemetry metric [[GH-12020](https://github.com/hashicorp/vault/pull/12020)]
* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* core: Fix edge cases in the configuration endpoint for barrier key autorotation. [[GH-11541](https://github.com/hashicorp/vault/pull/11541)]
* core: Fix goroutine leak when updating rate limit quota [[GH-11371](https://github.com/hashicorp/vault/pull/11371)]
* core (enterprise): Fix panic on DR secondary when there are lease count quotas [[GH-11742](https://github.com/hashicorp/vault/pull/11742)]
* core: Fix race that allowed remounting on path used by another mount [[GH-11453](https://github.com/hashicorp/vault/pull/11453)]
* core: Fix storage entry leak when revoking leases created with non-orphan batch tokens. [[GH-11377](https://github.com/hashicorp/vault/pull/11377)]
* core: Fixed double counting of http requests after operator stepdown [[GH-11970](https://github.com/hashicorp/vault/pull/11970)]
* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]
* identity: Use correct mount accessor when refreshing external group memberships. [[GH-11506](https://github.com/hashicorp/vault/pull/11506)]
* mongo-db: default username template now strips invalid '.' characters [[GH-11872](https://github.com/hashicorp/vault/pull/11872)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* replication: Fix panic trying to update walState during identity group invalidation.
* replication: Fix: mounts created within a namespace that was part of an Allow
  filtering rule would not appear on performance secondary if created after rule
  was defined.
* secret/pki: use case insensitive domain name comparison as per RFC1035 section 2.3.3
* secret: fix the bug where transit encrypt batch doesn't work with key_version [[GH-11628](https://github.com/hashicorp/vault/pull/11628)]
* secrets/ad: Forward all creds requests to active node [[GH-76](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/76)] [[GH-11836](https://github.com/hashicorp/vault/pull/11836)]
* secrets/database/cassandra: Fixed issue where hostnames were not being validated when using TLS [[GH-11365](https://github.com/hashicorp/vault/pull/11365)]
* secrets/database/cassandra: Fixed issue where the PEM parsing logic of `pem_bundle` and `pem_json` didn't work for CA-only configurations [[GH-11861](https://github.com/hashicorp/vault/pull/11861)]
* secrets/database/cassandra: Updated default statement for password rotation to allow for special characters. This applies to root and static credentials. [[GH-11262](https://github.com/hashicorp/vault/pull/11262)]
* secrets/database: Fix marshalling to allow providing numeric arguments to external database plugins. [[GH-11451](https://github.com/hashicorp/vault/pull/11451)]
* secrets/database: Fixed an issue that prevented external database plugin processes from restarting after a shutdown. [[GH-12087](https://github.com/hashicorp/vault/pull/12087)]
* secrets/database: Fixed minor race condition when rotate-root is called [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database: Fixes issue for V4 database interface where `SetCredentials` wasn't falling back to using `RotateRootCredentials` if `SetCredentials` is `Unimplemented` [[GH-11585](https://github.com/hashicorp/vault/pull/11585)]
* secrets/openldap: Fix bug where schema was not compatible with rotate-root [#24](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/24) [[GH-12019](https://github.com/hashicorp/vault/pull/12019)]
* storage/dynamodb: Handle throttled batch write requests by retrying, without which writes could be lost. [[GH-10181](https://github.com/hashicorp/vault/pull/10181)]
* storage/raft: Support cluster address change for nodes in a cluster managed by autopilot [[GH-11247](https://github.com/hashicorp/vault/pull/11247)]
* storage/raft: Tweak creation of vault.db file [[GH-12034](https://github.com/hashicorp/vault/pull/12034)]
* storage/raft: leader_tls_servername wasn't used unless leader_ca_cert_file and/or mTLS were configured. [[GH-11252](https://github.com/hashicorp/vault/pull/11252)]
* tokenutil: Perform the num uses check before token type. [[GH-11647](https://github.com/hashicorp/vault/pull/11647)]
* transform (enterprise): Fix an issue with malformed transform configuration
  storage when upgrading from 1.5 to 1.6.  See Upgrade Notes for 1.6.x.
* ui: Add role from database connection automatically populates the database for new role [[GH-11119](https://github.com/hashicorp/vault/pull/11119)]
* ui: Add root rotation statements support to appropriate database secret engine plugins [[GH-11404](https://github.com/hashicorp/vault/pull/11404)]
* ui: Automatically refresh the page when user logs out [[GH-12035](https://github.com/hashicorp/vault/pull/12035)]
* ui: Fix Version History queryParams on LinkedBlock [[GH-12079](https://github.com/hashicorp/vault/pull/12079)]
* ui: Fix bug where database secret engines with custom names cannot delete connections [[GH-11127](https://github.com/hashicorp/vault/pull/11127)]
* ui: Fix bug where the UI does not recognize version 2 KV until refresh, and fix [object Object] error message [[GH-11258](https://github.com/hashicorp/vault/pull/11258)]
* ui: Fix database role CG access [[GH-12111](https://github.com/hashicorp/vault/pull/12111)]
* ui: Fix date display on expired token notice [[GH-11142](https://github.com/hashicorp/vault/pull/11142)]
* ui: Fix entity group membership and metadata not showing [[GH-11641](https://github.com/hashicorp/vault/pull/11641)]
* ui: Fix error message caused by control group [[GH-11143](https://github.com/hashicorp/vault/pull/11143)]
* ui: Fix footer URL linking to the correct version changelog. [[GH-11283](https://github.com/hashicorp/vault/pull/11283)]
* ui: Fix issue where logging in without namespace input causes error [[GH-11094](https://github.com/hashicorp/vault/pull/11094)]
* ui: Fix namespace-bug on login [[GH-11182](https://github.com/hashicorp/vault/pull/11182)]
* ui: Fix status menu no showing on login [[GH-11213](https://github.com/hashicorp/vault/pull/11213)]
* ui: Fix text link URL on database roles list [[GH-11597](https://github.com/hashicorp/vault/pull/11597)]
* ui: Fixed and updated lease renewal picker [[GH-11256](https://github.com/hashicorp/vault/pull/11256)]
* ui: fix control group access for database credential [[GH-12024](https://github.com/hashicorp/vault/pull/12024)]
* ui: fix issue where select-one option was not showing in secrets database role creation [[GH-11294](https://github.com/hashicorp/vault/pull/11294)]
* ui: fix oidc login with Safari [[GH-11884](https://github.com/hashicorp/vault/pull/11884)]

## 1.7.10
### March 3, 2022

SECURITY:

* transform (enterprise): Vault Enterprise (“Vault”) clusters using the tokenization transform feature can expose the tokenization key through the tokenization key configuration endpoint to authorized operators with read permissions on this endpoint. This vulnerability, CVE-2022-25244, was fixed in Vault Enterprise 1.7.10, 1.8.9, and 1.9.4.

BUG FIXES:

* database/mssql: Removed string interpolation on internal queries and replaced them with inline queries using named parameters. [[GH-13799](https://github.com/hashicorp/vault/pull/13799)]
* ui: Fix issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Trigger background token self-renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.7.9
### January 27, 2022

IMPROVEMENTS:

* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]

BUG FIXES:

* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13493](https://github.com/hashicorp/vault/pull/13493)]
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13735](https://github.com/hashicorp/vault/pull/13735)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]

## 1.7.8
### December 21, 2021

CHANGES:

* go: Update go version to 1.16.12 [[GH-13422](https://github.com/hashicorp/vault/pull/13422)]

BUG FIXES:

* auth/aws: Fixes ec2 login no longer supporting DSA signature verification [[GH-12340](https://github.com/hashicorp/vault/pull/12340)]
* identity: Fix a panic on arm64 platform when doing identity I/O. [[GH-12371](https://github.com/hashicorp/vault/pull/12371)]

## 1.7.7
### December 9, 2021

SECURITY:

* storage/raft: Integrated Storage backend could be caused to crash by an authenticated user with write permissions to the KV secrets engine. This vulnerability, CVE-2021-45042, was fixed in Vault 1.7.7, 1.8.6, and 1.9.1.

BUG FIXES:

* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]

## 1.7.6
### November 4, 2021

SECURITY:

* core/identity: Templated ACL policies would always match the first-created entity alias if multiple entity aliases existed for a specified entity and mount combination, potentially resulting in incorrect policy enforcement. This vulnerability, CVE-2021-43998, was fixed in Vault and Vault Enterprise 1.7.6, 1.8.5, and 1.9.0.

BUG FIXES:

* auth/aws: fix config/rotate-root to store new key [[GH-12715](https://github.com/hashicorp/vault/pull/12715)]
* core/identity: Cleanup alias in the in-memory entity after an alias deletion by ID [[GH-12834](https://github.com/hashicorp/vault/pull/12834)]
* core/identity: Disallow entity alias creation/update if a conflicting alias exists for the target entity and mount combination [[GH-12747](https://github.com/hashicorp/vault/pull/12747)]
* core: Fix a deadlock on HA leadership transfer [[GH-12691](https://github.com/hashicorp/vault/pull/12691)]
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* kmip (enterprise): Fix handling of custom attributes when servicing GetAttributes requests
* kmip (enterprise): Fix handling of invalid role parameters within various vault api calls
* kmip (enterprise): Forward KMIP register operations to the active node
* secrets/keymgmt (enterprise): Fix support for Azure Managed HSM Key Vault instances. [[GH-12957](https://github.com/hashicorp/vault/pull/12957)]
* storage/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* database/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* transform (enterprise): Fix an error where the decode response of an expired token is an empty result rather than an error.

## 1.7.5
### 29 September 2021

SECURITY:

* core/identity: A Vault user with write permission to an entity alias ID sharing a mount accessor with another user may acquire this other user’s policies by merging their identities. This vulnerability, CVE-2021-41802, was fixed in Vault and Vault Enterprise 1.7.5 and 1.8.4.

IMPROVEMENTS:

* secrets/pki: Allow signing of self-issued certs with a different signature algorithm. [[GH-12514](https://github.com/hashicorp/vault/pull/12514)]

BUG FIXES:

* agent: Avoid possible `unexpected fault address` panic when using persistent cache. [[GH-12534](https://github.com/hashicorp/vault/pull/12534)]
* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* identity: Fail alias rename if the resulting (name,accessor) exists already [[GH-12473](https://github.com/hashicorp/vault/pull/12473)]
* raft (enterprise): Fix panic when updating auto-snapshot config
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12598](https://github.com/hashicorp/vault/pull/12598)]
* storage/raft: Detect incomplete raft snapshots in api.RaftSnapshot(), and thereby in `vault operator raft snapshot save`. [[GH-12388](https://github.com/hashicorp/vault/pull/12388)]
* ui: Fixed api explorer routing bug [[GH-12354](https://github.com/hashicorp/vault/pull/12354)]

## 1.7.4
### 26 August 2021

SECURITY:

* *UI Secret Caching*: The Vault UI erroneously cached and exposed user-viewed secrets between authenticated sessions in a single shared browser, if the browser window / tab was not refreshed or closed between logout and a subsequent login. This vulnerability, CVE-2021-38554, was fixed in Vault 1.8.0 and will be addressed in pending 1.7.4 / 1.6.6 releases.

CHANGES:

* Alpine: Docker images for Vault 1.6.6+, 1.7.4+, and 1.8.2+ are built with Alpine 3.14, due to CVE-2021-36159
* go: Update go version to 1.15.15 [[GH-12411](https://github.com/hashicorp/vault/pull/12411)]

IMPROVEMENTS:

* ui: Updated node to v14, latest stable build [[GH-12049](https://github.com/hashicorp/vault/pull/12049)]

BUG FIXES:

* replication (enterprise): Fix a panic that could occur when checking the last wal and the log shipper buffer is empty.
* cli: vault debug now puts newlines after every captured log line. [[GH-12175](https://github.com/hashicorp/vault/pull/12175)]
* database/couchbase: change default template to truncate username at 128 characters [[GH-12299](https://github.com/hashicorp/vault/pull/12299)]
* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* secrets/database/cassandra: Fixed issue where the PEM parsing logic of `pem_bundle` and `pem_json` didn't work for CA-only configurations [[GH-11861](https://github.com/hashicorp/vault/pull/11861)]
* secrets/database: Fixed an issue that prevented external database plugin processes from restarting after a shutdown. [[GH-12087](https://github.com/hashicorp/vault/pull/12087)]
* ui: Automatically refresh the page when user logs out [[GH-12035](https://github.com/hashicorp/vault/pull/12035)]
* ui: Fix database role CG access [[GH-12111](https://github.com/hashicorp/vault/pull/12111)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: fix control group access for database credential [[GH-12024](https://github.com/hashicorp/vault/pull/12024)]
* ui: fix oidc login with Safari [[GH-11884](https://github.com/hashicorp/vault/pull/11884)]

## 1.7.3
### June 16th, 2021

CHANGES:

* go: Update go version to 1.15.13 [[GH-11857](https://github.com/hashicorp/vault/pull/11857)]

IMPROVEMENTS:

* db/cassandra: Added tls_server_name to specify server name for TLS validation [[GH-11820](https://github.com/hashicorp/vault/pull/11820)]
* ui: Add specific error message if unseal fails due to license [[GH-11705](https://github.com/hashicorp/vault/pull/11705)]

BUG FIXES:

* auth/jwt: Updates the [hashicorp/cap](https://github.com/hashicorp/cap) library to `v0.1.0` to
bring in a verification key caching fix. [[GH-11784](https://github.com/hashicorp/vault/pull/11784)]
* core (enterprise): serialize access to HSM entropy generation to avoid errors in concurrent key generation.
* secret: fix the bug where transit encrypt batch doesn't work with key_version [[GH-11628](https://github.com/hashicorp/vault/pull/11628)]
* secrets/ad: Forward all creds requests to active node [[GH-76](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/76)] [[GH-11836](https://github.com/hashicorp/vault/pull/11836)]
* tokenutil: Perform the num uses check before token type. [[GH-11647](https://github.com/hashicorp/vault/pull/11647)]

## 1.7.2
### May 20th, 2021

SECURITY:

* Non-Expiring Leases: Vault and Vault Enterprise renewed nearly-expiring token
leases and dynamic secret leases with a zero-second TTL, causing them to be
treated as non-expiring, and never revoked. This issue affects Vault and Vault
Enterprise versions 0.10.0 through 1.7.1, and is fixed in 1.5.9, 1.6.5, and
1.7.2 (CVE-2021-32923).

CHANGES:

* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* auth/gcp: Update to v0.9.1 to use IAM Service Account Credentials API for
signing JWTs [[GH-11494](https://github.com/hashicorp/vault/pull/11494)]

IMPROVEMENTS:

* api, agent: LifetimeWatcher now does more retries when renewal failures occur.  This also impacts Agent auto-auth and leases managed via Agent caching. [[GH-11445](https://github.com/hashicorp/vault/pull/11445)]
* auth/aws: Underlying error included in validation failure message. [[GH-11638](https://github.com/hashicorp/vault/pull/11638)]
* http: Add optional HTTP response headers for hostname and raft node ID [[GH-11289](https://github.com/hashicorp/vault/pull/11289)]
* secrets/aws: add ability to provide a role session name when generating STS credentials [[GH-11345](https://github.com/hashicorp/vault/pull/11345)]
* secrets/database/mongodb: Add ability to customize `SocketTimeout`, `ConnectTimeout`, and `ServerSelectionTimeout` [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database/mongodb: Increased throughput by allowing for multiple request threads to simultaneously update users in MongoDB [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]

BUG FIXES:

* agent/cert: Fix issue where the API client on agent was not honoring certificate
information from the auto-auth config map on renewals or retries. [[GH-11576](https://github.com/hashicorp/vault/pull/11576)]
* agent: Fixed agent templating to use configured tls servername values [[GH-11288](https://github.com/hashicorp/vault/pull/11288)]
* core (enterprise): Fix plugins mounted in namespaces being unable to use password policies [[GH-11596](https://github.com/hashicorp/vault/pull/11596)]
* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]
* identity: Use correct mount accessor when refreshing external group memberships. [[GH-11506](https://github.com/hashicorp/vault/pull/11506)]
* replication: Fix panic trying to update walState during identity group invalidation. [[GH-1865](https://github.com/hashicorp/vault/pull/1865)]
* secrets/database: Fix marshalling to allow providing numeric arguments to external database plugins. [[GH-11451](https://github.com/hashicorp/vault/pull/11451)]
* secrets/database: Fixed minor race condition when rotate-root is called [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database: Fixes issue for V4 database interface where `SetCredentials` wasn't falling back to using `RotateRootCredentials` if `SetCredentials` is `Unimplemented` [[GH-11585](https://github.com/hashicorp/vault/pull/11585)]
* secrets/keymgmt (enterprise): Fixes audit logging for the read key response.
* storage/raft: Support cluster address change for nodes in a cluster managed by autopilot [[GH-11247](https://github.com/hashicorp/vault/pull/11247)]
* ui: Fix entity group membership and metadata not showing [[GH-11641](https://github.com/hashicorp/vault/pull/11641)]
* ui: Fix text link URL on database roles list [[GH-11597](https://github.com/hashicorp/vault/pull/11597)]

## 1.7.1
### 21 April 2021

SECURITY:

* The PKI Secrets Engine tidy functionality may cause Vault to exclude revoked-but-unexpired certificates from the
  Vault CRL. This vulnerability affects Vault and Vault Enterprise 1.5.1 and newer and was fixed in versions
  1.5.8, 1.6.4, and 1.7.1. (CVE-2021-27668)
* The Cassandra Database and Storage backends were not correctly verifying TLS certificates. This issue affects all
  versions of Vault and Vault Enterprise and was fixed in versions 1.6.4, and 1.7.1. (CVE-2021-27400)

CHANGES:

* go: Update to Go 1.15.11 [[GH-11395](https://github.com/hashicorp/vault/pull/11395)]

IMPROVEMENTS:

* auth/jwt: Adds ability to directly provide service account JSON in G Suite provider config. [[GH-11388](https://github.com/hashicorp/vault/pull/11388)]
* core: Add tls_max_version listener config option. [[GH-11226](https://github.com/hashicorp/vault/pull/11226)]
* core: Add metrics for standby node forwarding. [[GH-11366](https://github.com/hashicorp/vault/pull/11366)]
* core: allow arbitrary length stack traces upon receiving SIGUSR2 (was 32MB) [[GH-11364](https://github.com/hashicorp/vault/pull/11364)]
* storage/raft: Support autopilot for HA only raft storage. [[GH-11260](https://github.com/hashicorp/vault/pull/11260)]

BUG FIXES:

* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* core: Fix goroutine leak when updating rate limit quota [[GH-11371](https://github.com/hashicorp/vault/pull/11371)]
* core: Fix storage entry leak when revoking leases created with non-orphan batch tokens. [[GH-11377](https://github.com/hashicorp/vault/pull/11377)]
* core: requests forwarded by standby weren't always timed out. [[GH-11322](https://github.com/hashicorp/vault/pull/11322)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* replication: Fix: mounts created within a namespace that was part of an Allow
  filtering rule would not appear on performance secondary if created after rule
  was defined.
* replication: Perf standby nodes on newly enabled DR secondary sometimes couldn't connect to active node with TLS errors. [[GH-1823](https://github.com/hashicorp/vault/pull/1823)]
* secrets/database/cassandra: Fixed issue where hostnames were not being validated when using TLS [[GH-11365](https://github.com/hashicorp/vault/pull/11365)]
* secrets/database/cassandra: Updated default statement for password rotation to allow for special characters. This applies to root and static credentials. [[GH-11262](https://github.com/hashicorp/vault/pull/11262)]
* storage/dynamodb: Handle throttled batch write requests by retrying, without which writes could be lost. [[GH-10181](https://github.com/hashicorp/vault/pull/10181)]
* storage/raft: leader_tls_servername wasn't used unless leader_ca_cert_file and/or mTLS were configured. [[GH-11252](https://github.com/hashicorp/vault/pull/11252)]
* storage/raft: using raft for ha_storage with a different storage backend was broken in 1.7.0, now fixed. [[GH-11340](https://github.com/hashicorp/vault/pull/11340)]
* ui: Add root rotation statements support to appropriate database secret engine plugins [[GH-11404](https://github.com/hashicorp/vault/pull/11404)]
* ui: Fix bug where the UI does not recognize version 2 KV until refresh, and fix [object Object] error message [[GH-11258](https://github.com/hashicorp/vault/pull/11258)]
* ui: Fix OIDC bug seen when running on HCP [[GH-11283](https://github.com/hashicorp/vault/pull/11283)]
* ui: Fix namespace-bug on login [[GH-11182](https://github.com/hashicorp/vault/pull/11182)]
* ui: Fix status menu no showing on login [[GH-11213](https://github.com/hashicorp/vault/pull/11213)]
* ui: fix issue where select-one option was not showing in secrets database role creation [[GH-11294](https://github.com/hashicorp/vault/pull/11294)]

## 1.7.0
### 24 March 2021

CHANGES:

* agent: Failed auto-auth attempts are now throttled by an exponential backoff instead of the
~2 second retry delay. The maximum backoff may be configured with the new `max_backoff` parameter,
which defaults to 5 minutes. [[GH-10964](https://github.com/hashicorp/vault/pull/10964)]
* aws/auth: AWS Auth concepts and endpoints that use the "whitelist" and "blacklist" terms
have been updated to more inclusive language (e.g. `/auth/aws/identity-whitelist` has been
updated to`/auth/aws/identity-accesslist`). The old and new endpoints are aliases,
sharing the same underlying data. The legacy endpoint names are considered **deprecated**
and will be removed in a future release (not before Vault 1.9). The complete list of
endpoint changes is available in the [AWS Auth API docs](/api-docs/auth/aws#deprecations-effective-in-vault-1-7).
* go: Update Go version to 1.15.10 [[GH-11114](https://github.com/hashicorp/vault/pull/11114)] [[GH-11173](https://github.com/hashicorp/vault/pull/11173)]

FEATURES:

* **Aerospike Storage Backend**: Add support for using Aerospike as a storage backend [[GH-10131](https://github.com/hashicorp/vault/pull/10131)]
* **Autopilot for Integrated Storage**: A set of features has been added to allow for automatic operator-friendly management of Vault servers. This is only applicable when integrated storage is in use.
  * **Dead Server Cleanup**: Dead servers will periodically be cleaned up and removed from the Raft peer set, to prevent them from interfering with the quorum size and leader elections.
  * **Server Health Checking**: An API has been added to track the state of servers, including their health.
  * **New Server Stabilization**: When a new server is added to the cluster, there will be a waiting period where it must be healthy and stable for a certain amount of time before being promoted to a full, voting member.
* **Tokenization Secrets Engine (Enterprise)**: The Tokenization Secrets Engine is now generally available. We have added support for MySQL, key rotation, and snapshot/restore.
* replication (enterprise): The log shipper is now memory as well as length bound, and length and size can be separately configured.
* agent: Support for persisting the agent cache to disk [[GH-10938](https://github.com/hashicorp/vault/pull/10938)]
* auth/jwt: Adds `max_age` role parameter and `auth_time` claim validation. [[GH-10919](https://github.com/hashicorp/vault/pull/10919)]
* core (enterprise): X-Vault-Index and related headers can be used by clients to manage eventual consistency.
* kmip (enterprise): Use entropy augmentation to generate kmip certificates
* sdk: Private key generation in the certutil package now allows custom io.Readers to be used. [[GH-10653](https://github.com/hashicorp/vault/pull/10653)]
* secrets/aws: add IAM tagging support for iam_user roles [[GH-10953](https://github.com/hashicorp/vault/pull/10953)]
* secrets/database/cassandra: Add ability to customize dynamic usernames [[GH-10906](https://github.com/hashicorp/vault/pull/10906)]
* secrets/database/couchbase: Add ability to customize dynamic usernames [[GH-10995](https://github.com/hashicorp/vault/pull/10995)]
* secrets/database/mongodb: Add ability to customize dynamic usernames [[GH-10858](https://github.com/hashicorp/vault/pull/10858)]
* secrets/database/mssql: Add ability to customize dynamic usernames [[GH-10767](https://github.com/hashicorp/vault/pull/10767)]
* secrets/database/mysql: Add ability to customize dynamic usernames [[GH-10834](https://github.com/hashicorp/vault/pull/10834)]
* secrets/database/postgresql: Add ability to customize dynamic usernames [[GH-10766](https://github.com/hashicorp/vault/pull/10766)]
* secrets/db/snowflake: Added support for Snowflake to the Database Secret Engine [[GH-10603](https://github.com/hashicorp/vault/pull/10603)]
* secrets/keymgmt (enterprise): Adds beta support for distributing and managing keys in AWS KMS.
* secrets/keymgmt (enterprise): Adds general availability for distributing and managing keys in Azure Key Vault.
* secrets/openldap: Added dynamic roles to OpenLDAP similar to the combined database engine [[GH-10996](https://github.com/hashicorp/vault/pull/10996)]
* secrets/terraform: New secret engine for managing Terraform Cloud API tokens [[GH-10931](https://github.com/hashicorp/vault/pull/10931)]
* ui: Adds check for feature flag on application, and updates namespace toolbar on login if present [[GH-10588](https://github.com/hashicorp/vault/pull/10588)]
* ui: Adds the wizard to the Database Secret Engine [[GH-10982](https://github.com/hashicorp/vault/pull/10982)]
* ui: Database secrets engine, supporting MongoDB only [[GH-10655](https://github.com/hashicorp/vault/pull/10655)]

IMPROVEMENTS:

* agent: Add a `vault.retry` stanza that allows specifying number of retries on failure; this applies both to templating and proxied requests. [[GH-11113](https://github.com/hashicorp/vault/pull/11113)]
* agent: Agent can now run as a Windows service. [[GH-10231](https://github.com/hashicorp/vault/pull/10231)]
* agent: Better concurrent request handling on identical requests proxied through Agent. [[GH-10705](https://github.com/hashicorp/vault/pull/10705)]
* agent: Route templating server through cache when persistent cache is enabled. [[GH-10927](https://github.com/hashicorp/vault/pull/10927)]
* agent: change auto-auth to preload an existing token on start [[GH-10850](https://github.com/hashicorp/vault/pull/10850)]
* auth/approle: Secrets ID generation endpoint now returns `secret_id_ttl` as part of its response. [[GH-10826](https://github.com/hashicorp/vault/pull/10826)]
* auth/ldap: Improve consistency in error messages [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]
* auth/okta: Adds support for Okta Verify TOTP MFA. [[GH-10942](https://github.com/hashicorp/vault/pull/10942)]
* changelog: Add dependencies listed in dependencies/2-25-21 [[GH-11015](https://github.com/hashicorp/vault/pull/11015)]
* command/debug: Now collects logs (at level `trace`) as a periodic output. [[GH-10609](https://github.com/hashicorp/vault/pull/10609)]
* core (enterprise): "vault status" command works when a namespace is set. [[GH-10725](https://github.com/hashicorp/vault/pull/10725)]
* core (enterprise): Update Trial Enterprise license from 30 minutes to 6 hours
* core/metrics: Added "vault operator usage" command. [[GH-10365](https://github.com/hashicorp/vault/pull/10365)]
* core/metrics: New telemetry metrics reporting lease expirations by time interval and namespace [[GH-10375](https://github.com/hashicorp/vault/pull/10375)]
* core: Added active since timestamp to the status output of active nodes. [[GH-10489](https://github.com/hashicorp/vault/pull/10489)]
* core: Check audit device with a test message before adding it. [[GH-10520](https://github.com/hashicorp/vault/pull/10520)]
* core: Track barrier encryption count and automatically rotate after a large number of operations or on a schedule [[GH-10774](https://github.com/hashicorp/vault/pull/10774)]
* core: add metrics for active entity count [[GH-10514](https://github.com/hashicorp/vault/pull/10514)]
* core: add partial month client count api [[GH-11022](https://github.com/hashicorp/vault/pull/11022)]
* core: dev mode listener allows unauthenticated sys/metrics requests [[GH-10992](https://github.com/hashicorp/vault/pull/10992)]
* core: reduce memory used by leases [[GH-10726](https://github.com/hashicorp/vault/pull/10726)]
* secrets/gcp: Truncate ServiceAccount display names longer than 100 characters. [[GH-10558](https://github.com/hashicorp/vault/pull/10558)]
* storage/raft (enterprise): Listing of peers is now allowed on DR secondary
cluster nodes, as an update operation that takes in DR operation token for
authenticating the request.
* transform (enterprise): Improve FPE transformation performance
* transform (enterprise): Use transactions with batch tokenization operations for improved performance
* ui: Clarify language on usage metrics page empty state [[GH-10951](https://github.com/hashicorp/vault/pull/10951)]
* ui: Customize MongoDB input fields on Database Secrets Engine [[GH-10949](https://github.com/hashicorp/vault/pull/10949)]
* ui: Upgrade Ember-cli from 3.8 to 3.22. [[GH-9972](https://github.com/hashicorp/vault/pull/9972)]
* ui: Upgrade Storybook from 5.3.19 to 6.1.17. [[GH-10904](https://github.com/hashicorp/vault/pull/10904)]
* ui: Upgrade date-fns from 1.3.0 to 2.16.1. [[GH-10848](https://github.com/hashicorp/vault/pull/10848)]
* ui: Upgrade dependencies to resolve potential JS vulnerabilities [[GH-10677](https://github.com/hashicorp/vault/pull/10677)]
* ui: better errors on Database secrets engine role create [[GH-10980](https://github.com/hashicorp/vault/pull/10980)]

BUG FIXES:

* agent: Only set the namespace if the VAULT_NAMESPACE env var isn't present [[GH-10556](https://github.com/hashicorp/vault/pull/10556)]
* agent: Set TokenParent correctly in the Index to be cached. [[GH-10833](https://github.com/hashicorp/vault/pull/10833)]
* agent: Set namespace for template server in agent. [[GH-10757](https://github.com/hashicorp/vault/pull/10757)]
* api/sys/config/ui: Fixes issue where multiple UI custom header values are ignored and only the first given value is used [[GH-10490](https://github.com/hashicorp/vault/pull/10490)]
* api: Fixes CORS API methods that were outdated and invalid [[GH-10444](https://github.com/hashicorp/vault/pull/10444)]
* auth/jwt: Fixes `bound_claims` validation for provider-specific group and user info fetching. [[GH-10546](https://github.com/hashicorp/vault/pull/10546)]
* auth/jwt: Fixes an issue where JWT verification keys weren't updated after a `jwks_url` change. [[GH-10919](https://github.com/hashicorp/vault/pull/10919)]
* auth/jwt: Fixes an issue where `jwt_supported_algs` were not being validated for JWT auth using
`jwks_url` and `jwt_validation_pubkeys`. [[GH-10919](https://github.com/hashicorp/vault/pull/10919)]
* auth/oci: Fixes alias name to use the role name, and not the literal string `name` [[GH-10](https://github.com/hashicorp/vault-plugin-auth-oci/pull/10)] [[GH-10952](https://github.com/hashicorp/vault/pull/10952)]
* consul-template: Update consul-template vendor version and associated dependencies to master, 
pulling in https://github.com/hashicorp/consul-template/pull/1447 [[GH-10756](https://github.com/hashicorp/vault/pull/10756)]
* core (enterprise): Limit entropy augmentation during token generation to root tokens. [[GH-10487](https://github.com/hashicorp/vault/pull/10487)]
* core (enterprise): Vault EGP policies attached to path * were not correctly scoped to the namespace.
* core/identity: Fix deadlock in entity merge endpoint. [[GH-10877](https://github.com/hashicorp/vault/pull/10877)]
* core: Avoid disclosing IP addresses in the errors of unauthenticated requests [[GH-10579](https://github.com/hashicorp/vault/pull/10579)]
* core: Fix client.Clone() to include the address [[GH-10077](https://github.com/hashicorp/vault/pull/10077)]
* core: Fix duplicate quotas on performance standby nodes. [[GH-10855](https://github.com/hashicorp/vault/pull/10855)]
* core: Fix rate limit resource quota migration from 1.5.x to 1.6.x by ensuring `purgeInterval` and
`staleAge` are set appropriately. [[GH-10536](https://github.com/hashicorp/vault/pull/10536)]
* core: Make all APIs that report init status consistent, and make them report
initialized=true when a Raft join is in progress. [[GH-10498](https://github.com/hashicorp/vault/pull/10498)]
* core: Make the response to an unauthenticated request to sys/internal endpoints consistent regardless of mount existence. [[GH-10650](https://github.com/hashicorp/vault/pull/10650)]
* core: Turn off case sensitivity for allowed entity alias check during token create operation. [[GH-10743](https://github.com/hashicorp/vault/pull/10743)]
* http: change max_request_size to be unlimited when the config value is less than 0 [[GH-10072](https://github.com/hashicorp/vault/pull/10072)]
* license: Fix license caching issue that prevents new licenses to get picked up by the license manager [[GH-10424](https://github.com/hashicorp/vault/pull/10424)]
* metrics: Protect emitMetrics from panicking during post-seal [[GH-10708](https://github.com/hashicorp/vault/pull/10708)]
* quotas/rate-limit: Fix quotas enforcing old rate limit quota paths [[GH-10689](https://github.com/hashicorp/vault/pull/10689)]
* replication (enterprise): Fix bug with not starting merkle sync while requests are in progress
* secrets/database/influxdb: Fix issue where not all errors from InfluxDB were being handled [[GH-10384](https://github.com/hashicorp/vault/pull/10384)]
* secrets/database/mysql: Fixes issue where the DisplayName within generated usernames was the incorrect length [[GH-10433](https://github.com/hashicorp/vault/pull/10433)]
* secrets/database: Sanitize `private_key` field when reading database plugin config [[GH-10416](https://github.com/hashicorp/vault/pull/10416)]
* secrets/gcp: Fix issue with account and iam_policy roleset WALs not being removed after attempts when GCP project no longer exists [[GH-10759](https://github.com/hashicorp/vault/pull/10759)]
* secrets/transit: allow for null string to be used for optional parameters in encrypt and decrypt [[GH-10386](https://github.com/hashicorp/vault/pull/10386)]
* serviceregistration: Fix race during shutdown of Consul service registration. [[GH-10901](https://github.com/hashicorp/vault/pull/10901)]
* storage/raft (enterprise): Automated snapshots with Azure required specifying
`azure_blob_environment`, which should have had as a default `AZUREPUBLICCLOUD`.
* storage/raft (enterprise): Reading a non-existent auto snapshot config now returns 404.
* storage/raft (enterprise): The parameter aws_s3_server_kms_key was misnamed and
didn't work.  Renamed to aws_s3_kms_key, and make it work so that when provided
the given key will be used to encrypt the snapshot using AWS KMS.
* transform (enterprise): Fix bug tokenization handling metadata on exportable stores
* transform (enterprise): Fix bug where tokenization store changes are persisted but don't take effect
* transform (enterprise): Fix transform configuration not handling `stores` parameter on the legacy path
* transform (enterprise): Make expiration timestamps human readable
* transform (enterprise): Return false for invalid tokens on the validate endpoint rather than returning an HTTP error
* ui: Add role from database connection automatically populates the database for new role [[GH-11119](https://github.com/hashicorp/vault/pull/11119)]
* ui: Fix bug in Transform secret engine when a new role is added and then removed from a transformation [[GH-10417](https://github.com/hashicorp/vault/pull/10417)]
* ui: Fix bug that double encodes secret route when there are spaces in the path and makes you unable to view the version history. [[GH-10596](https://github.com/hashicorp/vault/pull/10596)]
* ui: Fix expected response from feature-flags endpoint [[GH-10684](https://github.com/hashicorp/vault/pull/10684)]
* ui: Fix footer URL linking to the correct version changelog. [[GH-10491](https://github.com/hashicorp/vault/pull/10491)]

DEPRECATIONS:
* aws/auth: AWS Auth endpoints that use the "whitelist" and "blacklist" terms have been deprecated.
Refer to the CHANGES section for additional details.

## 1.6.7
### 29 September 2021

BUG FIXES:

* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12597](https://github.com/hashicorp/vault/pull/12597)]

## 1.6.6
### 26 August 2021

SECURITY:

* *UI Secret Caching*: The Vault UI erroneously cached and exposed user-viewed secrets between authenticated sessions in a single shared browser, if the browser window / tab was not refreshed or closed between logout and a subsequent login. This vulnerability, CVE-2021-38554, was fixed in Vault 1.8.0 and will be addressed in pending 1.7.4 / 1.6.6 releases.

CHANGES:

* Alpine: Docker images for Vault 1.6.6+, 1.7.4+, and 1.8.2+ are built with Alpine 3.14, due to CVE-2021-36159
* go: Update go version to 1.15.15 [[GH-12423](https://github.com/hashicorp/vault/pull/12423)]

IMPROVEMENTS:

* db/cassandra: Added tls_server_name to specify server name for TLS validation [[GH-11820](https://github.com/hashicorp/vault/pull/11820)]

BUG FIXES:

* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* secret: fix the bug where transit encrypt batch doesn't work with key_version [[GH-11628](https://github.com/hashicorp/vault/pull/11628)]
* secrets/database: Fixed an issue that prevented external database plugin processes from restarting after a shutdown. [[GH-12087](https://github.com/hashicorp/vault/pull/12087)]
* ui: Automatically refresh the page when user logs out [[GH-12035](https://github.com/hashicorp/vault/pull/12035)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: fix oidc login with Safari [[GH-11884](https://github.com/hashicorp/vault/pull/11884)]

## 1.6.5
### May 20th, 2021

SECURITY:

* Non-Expiring Leases: Vault and Vault Enterprise renewed nearly-expiring token
leases and dynamic secret leases with a zero-second TTL, causing them to be
treated as non-expiring, and never revoked. This issue affects Vault and Vault
Enterprise versions 0.10.0 through 1.7.1, and is fixed in 1.5.9, 1.6.5, and
1.7.2 (CVE-2021-32923).

CHANGES:

* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* auth/gcp: Update to v0.8.1 to use IAM Service Account Credentials API for
signing JWTs [[GH-11498](https://github.com/hashicorp/vault/pull/11498)]

BUG FIXES:

* core (enterprise): Fix plugins mounted in namespaces being unable to use password policies [[GH-11596](https://github.com/hashicorp/vault/pull/11596)]
* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]
* secrets/database: Fix marshalling to allow providing numeric arguments to external database plugins. [[GH-11451](https://github.com/hashicorp/vault/pull/11451)]
* secrets/database: Fixes issue for V4 database interface where `SetCredentials` wasn't falling back to using `RotateRootCredentials` if `SetCredentials` is `Unimplemented` [[GH-11585](https://github.com/hashicorp/vault/pull/11585)]
* ui: Fix namespace-bug on login [[GH-11182](https://github.com/hashicorp/vault/pull/11182)]
  
## 1.6.4
### 21 April 2021

SECURITY:

* The PKI Secrets Engine tidy functionality may cause Vault to exclude revoked-but-unexpired certificates from the
  Vault CRL. This vulnerability affects Vault and Vault Enterprise 1.5.1 and newer and was fixed in versions
  1.5.8, 1.6.4, and 1.7.1. (CVE-2021-27668)
* The Cassandra Database and Storage backends were not correctly verifying TLS certificates. This issue affects all
  versions of Vault and Vault Enterprise and was fixed in versions 1.6.4, and 1.7.1. (CVE-2021-27400)

CHANGES:

* go: Update to Go 1.15.11 [[GH-11396](https://github.com/hashicorp/vault/pull/11396)]

IMPROVEMENTS:

* command/debug: Now collects logs (at level `trace`) as a periodic output. [[GH-10609](https://github.com/hashicorp/vault/pull/10609)]
* core: Add tls_max_version listener config option. [[GH-11226](https://github.com/hashicorp/vault/pull/11226)]
* core: allow arbitrary length stack traces upon receiving SIGUSR2 (was 32MB) [[GH-11364](https://github.com/hashicorp/vault/pull/11364)]

BUG FIXES:

* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* core: Fix goroutine leak when updating rate limit quota [[GH-11371](https://github.com/hashicorp/vault/pull/11371)]
* core: Fix storage entry leak when revoking leases created with non-orphan batch tokens. [[GH-11377](https://github.com/hashicorp/vault/pull/11377)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* pki: Preserve ordering of all DN attribute values when issuing certificates [[GH-11259](https://github.com/hashicorp/vault/pull/11259)]
* replication: Fix: mounts created within a namespace that was part of an Allow
  filtering rule would not appear on performance secondary if created after rule
  was defined.
* secrets/database/cassandra: Fixed issue where hostnames were not being validated when using TLS [[GH-11365](https://github.com/hashicorp/vault/pull/11365)]
* storage/raft: leader_tls_servername wasn't used unless leader_ca_cert_file and/or mTLS were configured. [[GH-11252](https://github.com/hashicorp/vault/pull/11252)]


## 1.6.3
### February 25, 2021

SECURITY:

* Limited Unauthenticated License Metadata Read: We addressed a security vulnerability that allowed for the unauthenticated
reading of Vault license metadata from DR Secondaries. This vulnerability affects Vault Enterprise and is
fixed in 1.6.3 (CVE-2021-27668).

CHANGES:

* secrets/mongodbatlas: Move from whitelist to access list API [[GH-10966](https://github.com/hashicorp/vault/pull/10966)]

IMPROVEMENTS:

* ui: Clarify language on usage metrics page empty state [[GH-10951](https://github.com/hashicorp/vault/pull/10951)]

BUG FIXES:

* auth/kubernetes: Cancel API calls to TokenReview endpoint when request context
is closed [[GH-10930](https://github.com/hashicorp/vault/pull/10930)]
* core/identity: Fix deadlock in entity merge endpoint. [[GH-10877](https://github.com/hashicorp/vault/pull/10877)]
* quotas: Fix duplicate quotas on performance standby nodes. [[GH-10855](https://github.com/hashicorp/vault/pull/10855)]
* quotas/rate-limit: Fix quotas enforcing old rate limit quota paths [[GH-10689](https://github.com/hashicorp/vault/pull/10689)]
* replication (enterprise): Don't write request count data on DR Secondaries.
Fixes DR Secondaries becoming out of sync approximately every 30s. [[GH-10970](https://github.com/hashicorp/vault/pull/10970)]
* secrets/azure (enterprise): Forward service principal credential creation to the
primary cluster if called on a performance standby or performance secondary. [[GH-10902](https://github.com/hashicorp/vault/pull/10902)]

## 1.6.2
### January 29, 2021

SECURITY:

* IP Address Disclosure: We fixed a vulnerability where, under some error
conditions, Vault would return an error message disclosing internal IP
addresses. This vulnerability affects Vault and Vault Enterprise and is fixed in
1.6.2 (CVE-2021-3024).
* Limited Unauthenticated Remove Peer: As of Vault 1.6, the remove-peer command
on DR secondaries did not require authentication. This issue impacts the
stability of HA architecture, as a bad actor could remove all standby
nodes from a DR
secondary. This issue affects Vault Enterprise 1.6.0 and 1.6.1, and is fixed in
1.6.2 (CVE-2021-3282).
* Mount Path Disclosure: Vault previously returned different HTTP status codes for
existent and non-existent mount paths. This behavior would allow unauthenticated
brute force attacks to reveal which paths had valid mounts. This issue affects
Vault and Vault Enterprise and is fixed in 1.6.2 (CVE-2020-25594).

CHANGES:

* go: Update go version to 1.15.7 [[GH-10730](https://github.com/hashicorp/vault/pull/10730)]

FEATURES:

* ui: Adds check for feature flag on application, and updates namespace toolbar on login if present [[GH-10588](https://github.com/hashicorp/vault/pull/10588)]

IMPROVEMENTS:

* core (enterprise): "vault status" command works when a namespace is set. [[GH-10725](https://github.com/hashicorp/vault/pull/10725)]
* core: reduce memory used by leases [[GH-10726](https://github.com/hashicorp/vault/pull/10726)]
* storage/raft (enterprise): Listing of peers is now allowed on DR secondary
cluster nodes, as an update operation that takes in DR operation token for
authenticating the request.
* core: allow setting tls_servername for raft retry/auto-join [[GH-10698](https://github.com/hashicorp/vault/pull/10698)]

BUG FIXES:

* agent: Set namespace for template server in agent. [[GH-10757](https://github.com/hashicorp/vault/pull/10757)]
* core: Make the response to an unauthenticated request to sys/internal endpoints consistent regardless of mount existence. [[GH-10650](https://github.com/hashicorp/vault/pull/10650)]
* metrics: Protect emitMetrics from panicking during post-seal [[GH-10708](https://github.com/hashicorp/vault/pull/10708)]
* secrets/gcp: Fix issue with account and iam_policy roleset WALs not being removed after attempts when GCP project no longer exists [[GH-10759](https://github.com/hashicorp/vault/pull/10759)]
* storage/raft (enterprise): Automated snapshots with Azure required specifying
`azure_blob_environment`, which should have had as a default `AZUREPUBLICCLOUD`.
* storage/raft (enterprise): Autosnapshots config and storage weren't excluded from
performance replication, causing conflicts and errors.
* ui: Fix bug that double encodes secret route when there are spaces in the path and makes you unable to view the version history. [[GH-10596](https://github.com/hashicorp/vault/pull/10596)]
* ui: Fix expected response from feature-flags endpoint [[GH-10684](https://github.com/hashicorp/vault/pull/10684)]

## 1.6.1
### December 16, 2020

SECURITY:

* LDAP Auth Method: We addressed an issue where error messages returned by the
  LDAP auth method allowed user enumeration [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]. This vulnerability affects Vault OSS and Vault 
  Enterprise and is fixed in 1.5.6 and 1.6.1 (CVE-2020-35177).
* Sentinel EGP: We've fixed incorrect handling of namespace paths to prevent
  users within namespaces from applying Sentinel EGP policies to paths above
  their namespace. This vulnerability affects Vault Enterprise and is fixed in
  1.5.6 and 1.6.1 (CVE-2020-35453).

IMPROVEMENTS:

* auth/ldap: Improve consistency in error messages [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]
* core/metrics: Added "vault operator usage" command. [[GH-10365](https://github.com/hashicorp/vault/pull/10365)]
* secrets/gcp: Truncate ServiceAccount display names longer than 100 characters. [[GH-10558](https://github.com/hashicorp/vault/pull/10558)]

BUG FIXES:

* agent: Only set the namespace if the VAULT_NAMESPACE env var isn't present [[GH-10556](https://github.com/hashicorp/vault/pull/10556)]
* auth/jwt: Fixes `bound_claims` validation for provider-specific group and user info fetching. [[GH-10546](https://github.com/hashicorp/vault/pull/10546)]
* core (enterprise): Vault EGP policies attached to path * were not correctly scoped to the namespace.
* core: Avoid deadlocks by ensuring that if grabLockOrStop returns stopped=true, the lock will not be held. [[GH-10456](https://github.com/hashicorp/vault/pull/10456)]
* core: Fix client.Clone() to include the address [[GH-10077](https://github.com/hashicorp/vault/pull/10077)]
* core: Fix rate limit resource quota migration from 1.5.x to 1.6.x by ensuring `purgeInterval` and
`staleAge` are set appropriately. [[GH-10536](https://github.com/hashicorp/vault/pull/10536)]
* core: Make all APIs that report init status consistent, and make them report
initialized=true when a Raft join is in progress. [[GH-10498](https://github.com/hashicorp/vault/pull/10498)]
* secrets/database/influxdb: Fix issue where not all errors from InfluxDB were being handled [[GH-10384](https://github.com/hashicorp/vault/pull/10384)]
* secrets/database/mysql: Fixes issue where the DisplayName within generated usernames was the incorrect length [[GH-10433](https://github.com/hashicorp/vault/pull/10433)]
* secrets/database: Sanitize `private_key` field when reading database plugin config [[GH-10416](https://github.com/hashicorp/vault/pull/10416)]
* secrets/transit: allow for null string to be used for optional parameters in encrypt and decrypt [[GH-10386](https://github.com/hashicorp/vault/pull/10386)]
* storage/raft (enterprise): The parameter aws_s3_server_kms_key was misnamed and didn't work.  Renamed to aws_s3_kms_key, and make it work so that when provided the given key will be used to encrypt the snapshot using AWS KMS.
* transform (enterprise): Fix bug tokenization handling metadata on exportable stores
* transform (enterprise): Fix transform configuration not handling `stores` parameter on the legacy path
* transform (enterprise): Make expiration timestamps human readable
* transform (enterprise): Return false for invalid tokens on the validate endpoint rather than returning an HTTP error
* transform (enterprise): Fix bug where tokenization store changes are persisted but don't take effect
* ui: Fix bug in Transform secret engine when a new role is added and then removed from a transformation [[GH-10417](https://github.com/hashicorp/vault/pull/10417)]
* ui: Fix footer URL linking to the correct version changelog. [[GH-10491](https://github.com/hashicorp/vault/pull/10491)]
* ui: Fox radio click on secrets and auth list pages. [[GH-10586](https://github.com/hashicorp/vault/pull/10586)]

## 1.6.0
### November 11th, 2020

NOTE:

Binaries for 32-bit macOS (i.e. the `darwin_386` build) will no longer be published. This target was dropped in the latest version of the Go compiler.

CHANGES:

* agent: Agent now properly returns a non-zero exit code on error, such as one due to template rendering failure. Using `error_on_missing_key` in the template config will cause agent to immediately exit on failure. In order to make agent properly exit due to continuous failure from template rendering errors, the old behavior of indefinitely restarting the template server is now changed to exit once the default retry attempt of 12 times (with exponential backoff) gets exhausted. [[GH-9670](https://github.com/hashicorp/vault/pull/9670)]  
* token: Periodic tokens generated by auth methods will have the period value stored in its token entry. [[GH-7885](https://github.com/hashicorp/vault/pull/7885)]
* core: New telemetry metrics reporting mount table size and number of entries [[GH-10201](https://github.com/hashicorp/vault/pull/10201)]
* go: Updated Go version to 1.15.4 [[GH-10366](https://github.com/hashicorp/vault/pull/10366)]

FEATURES:

* **Couchbase Secrets**: Vault can now manage static and dynamic credentials for Couchbase. [[GH-9664](https://github.com/hashicorp/vault/pull/9664)]
* **Expanded Password Policy Support**: Custom password policies are now supported for all database engines.
* **Integrated Storage Auto Snapshots (Enterprise)**: This feature enables an operator to schedule snapshots of the integrated storage backend and ensure those snapshots are persisted elsewhere.
* **Integrated Storage Cloud Auto Join**: This feature for integrated storage enables Vault nodes running in the cloud to automatically discover and join a Vault cluster via operator-supplied metadata.
* **Key Management Secrets Engine (Enterprise; Tech Preview)**: This new secret engine allows securely distributing and managing keys to Azure cloud KMS services.
* **Seal Migration**: With Vault 1.6, we will support migrating from an auto unseal mechanism to a different mechanism of the same type. For example, if you were using an AWS KMS key to automatically unseal, you can now migrate to a different AWS KMS key.
* **Tokenization (Enterprise; Tech Preview)**: Tokenization supports creating irreversible “tokens” from sensitive data. Tokens can be used in less secure environments, protecting the original data.
* **Vault Client Count**: Vault now counts the number of active entities (and non-entity tokens) per month and makes this information available via the "Metrics" section of the UI.

IMPROVEMENTS:

* auth/approle: Role names can now be referenced in templated policies through the `approle.metadata.role_name` property [[GH-9529](https://github.com/hashicorp/vault/pull/9529)]
* auth/aws: Improve logic check on wildcard `BoundIamPrincipalARNs` and include role name on error messages on check failure [[GH-10036](https://github.com/hashicorp/vault/pull/10036)]
* auth/jwt: Add support for fetching groups and user information from G Suite during authentication. [[GH-123](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/123)]
* auth/jwt: Adding EdDSA (ed25519) to supported algorithms [[GH-129](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/129)]
* auth/jwt: Improve cli authorization error [[GH-137](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/137)]
* auth/jwt: Add OIDC namespace_in_state option [[GH-140](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/140)]
* secrets/transit: fix missing plaintext in bulk decrypt response [[GH-9991](https://github.com/hashicorp/vault/pull/9991)]
* command/server: Delay informational messages in -dev mode until logs have settled. [[GH-9702](https://github.com/hashicorp/vault/pull/9702)]
* command/server: Add environment variable support for `disable_mlock`. [[GH-9931](https://github.com/hashicorp/vault/pull/9931)]
* core/metrics: Add metrics for storage cache [[GH_10079](https://github.com/hashicorp/vault/pull/10079)]
* core/metrics: Add metrics for leader status [[GH 10147](https://github.com/hashicorp/vault/pull/10147)]
* physical/azure: Add the ability to use Azure Instance Metadata Service to set the credentials for Azure Blob storage on the backend. [[GH-10189](https://github.com/hashicorp/vault/pull/10189)]
* sdk/framework: Add a time type for API fields. [[GH-9911](https://github.com/hashicorp/vault/pull/9911)]
* secrets/database: Added support for password policies to all databases [[GH-9641](https://github.com/hashicorp/vault/pull/9641),
  [and more](https://github.com/hashicorp/vault/pulls?q=is%3Apr+is%3Amerged+dbpw)]
* secrets/database/cassandra: Added support for static credential rotation [[GH-10051](https://github.com/hashicorp/vault/pull/10051)]
* secrets/database/elasticsearch: Added support for static credential rotation [[GH-19](https://github.com/hashicorp/vault-plugin-database-elasticsearch/pull/19)]
* secrets/database/hanadb: Added support for root credential & static credential rotation [[GH-10142](https://github.com/hashicorp/vault/pull/10142)]
* secrets/database/hanadb: Default password generation now includes dashes. Custom statements may need to be updated
  to include quotes around the password field [[GH-10142](https://github.com/hashicorp/vault/pull/10142)]
* secrets/database/influxdb: Added support for static credential rotation [[GH-10118](https://github.com/hashicorp/vault/pull/10118)]
* secrets/database/mongodbatlas: Added support for root credential rotation [[GH-14](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/14)]
* secrets/database/mongodbatlas: Support scopes field in creations statements for MongoDB Atlas database plugin [[GH-15](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/15)]
* seal/awskms: Add logging during awskms auto-unseal [[GH-9794](https://github.com/hashicorp/vault/pull/9794)]
* storage/azure: Update SDK library to use [azure-storage-blob-go](https://github.com/Azure/azure-storage-blob-go) since previous library has been deprecated. [[GH-9577](https://github.com/hashicorp/vault/pull/9577/)]
* secrets/ad: `rotate-root` now supports POST requests like other secret engines [[GH-70](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/70)]
* ui: Add ui functionality for the Transform Secret Engine [[GH-9665](https://github.com/hashicorp/vault/pull/9665)]
* ui: Pricing metrics dashboard [[GH-10049](https://github.com/hashicorp/vault/pull/10049)]

BUG FIXES:

* auth/jwt: Fix bug preventing config edit UI from rendering [[GH-141](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/141)]
* cli: Don't open or overwrite a raft snapshot file on an unsuccessful `vault operator raft snapshot` [[GH-9894](https://github.com/hashicorp/vault/pull/9894)]
* core: Implement constant time version of shamir GF(2^8) math [[GH-9932](https://github.com/hashicorp/vault/pull/9932)]
* core: Fix resource leak in plugin API (plugin-dependent, not all plugins impacted) [[GH-9557](https://github.com/hashicorp/vault/pull/9557)]
* core: Fix race involved in enabling certain features via a license change
* core: Fix error handling in HCL parsing of objects with invalid syntax [[GH-410](https://github.com/hashicorp/hcl/pull/410)]
* identity: Check for timeouts in entity API [[GH-9925](https://github.com/hashicorp/vault/pull/9925)]
* secrets/database: Fix handling of TLS options in mongodb connection strings [[GH-9519](https://github.com/hashicorp/vault/pull/9519)]
* secrets/gcp: Ensure that the IAM policy version is appropriately set after a roleset's bindings have changed. [[GH-93](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/93)]
* ui: Mask LDAP bindpass while typing [[GH-10087](https://github.com/hashicorp/vault/pull/10087)]
* ui: Update language in promote dr modal flow [[GH-10155](https://github.com/hashicorp/vault/pull/10155)]
* ui: Update language on replication primary dashboard for clarity [[GH-10205](https://github.com/hashicorp/vault/pull/10217)]
* core: Fix bug where updating an existing path quota could introduce a conflict. [[GH-10285](https://github.com/hashicorp/vault/pull/10285)]

## 1.5.9
### May 20th, 2021

SECURITY:

* Non-Expiring Leases: Vault and Vault Enterprise renewed nearly-expiring token
leases and dynamic secret leases with a zero-second TTL, causing them to be
treated as non-expiring, and never revoked. This issue affects Vault and Vault
Enterprise versions 0.10.0 through 1.7.1, and is fixed in 1.5.9, 1.6.5, and
1.7.2 (CVE-2021-32923).

CHANGES:

* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* auth/gcp: Update to v0.7.2 to use IAM Service Account Credentials API for
signing JWTs [[GH-11499](https://github.com/hashicorp/vault/pull/11499)]

BUG FIXES:

* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]

## 1.5.8
### 21 April 2021

SECURITY:

* The PKI Secrets Engine tidy functionality may cause Vault to exclude revoked-but-unexpired certificates from the
  Vault CRL. This vulnerability affects Vault and Vault Enterprise 1.5.1 and newer and was fixed in versions
  1.5.8, 1.6.4, and 1.7.1. (CVE-2021-27668)

CHANGES:

* go: Update to Go 1.14.15 [[GH-11397](https://github.com/hashicorp/vault/pull/11397)]

IMPROVEMENTS:

* core: Add tls_max_version listener config option. [[GH-11226](https://github.com/hashicorp/vault/pull/11226)]

BUG FIXES:

* core/identity: Fix deadlock in entity merge endpoint. [[GH-10877](https://github.com/hashicorp/vault/pull/10877)]
* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* core: Avoid deadlocks by ensuring that if grabLockOrStop returns stopped=true, the lock will not be held. [[GH-10456](https://github.com/hashicorp/vault/pull/10456)]

## 1.5.7
### January 29, 2021

SECURITY:

* IP Address Disclosure: We fixed a vulnerability where, under some error
conditions, Vault would return an error message disclosing internal IP
addresses. This vulnerability affects Vault and Vault Enterprise and is fixed in
1.6.2 and 1.5.7 (CVE-2021-3024).
* Mount Path Disclosure: Vault previously returned different HTTP status codes for
existent and non-existent mount paths. This behavior would allow unauthenticated
brute force attacks to reveal which paths had valid mounts. This issue affects
Vault and Vault Enterprise and is fixed in 1.6.2 and 1.5.7 (CVE-2020-25594).

IMPROVEMENTS:

* storage/raft (enterprise): Listing of peers is now allowed on DR secondary
cluster nodes, as an update operation that takes in DR operation token for
authenticating the request.

BUG FIXES:

* core: Avoid disclosing IP addresses in the errors of unauthenticated requests [[GH-10579](https://github.com/hashicorp/vault/pull/10579)]
* core: Make the response to an unauthenticated request to sys/internal endpoints consistent regardless of mount existence. [[GH-10650](https://github.com/hashicorp/vault/pull/10650)]

## 1.5.6
### December 16, 2020

SECURITY:

* LDAP Auth Method: We addressed an issue where error messages returned by the
  LDAP auth method allowed user enumeration [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]. This vulnerability affects Vault OSS and Vault 
  Enterprise and is fixed in 1.5.6 and 1.6.1 (CVE-2020-35177).
* Sentinel EGP: We've fixed incorrect handling of namespace paths to prevent
  users within namespaces from applying Sentinel EGP policies to paths above
  their namespace. This vulnerability affects Vault Enterprise and is fixed in
  1.5.6 and 1.6.1.

IMPROVEMENTS:

* auth/ldap: Improve consistency in error messages [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]

BUG FIXES:

* core (enterprise): Vault EGP policies attached to path * were not correctly scoped to the namespace.
* core: Fix bug where updating an existing path quota could introduce a conflict [[GH-10285](https://github.com/hashicorp/vault/pull/10285)]
* core: Fix client.Clone() to include the address [[GH-10077](https://github.com/hashicorp/vault/pull/10077)]
* quotas (enterprise): Reset cache before loading quotas in the db during startup
* secrets/transit: allow for null string to be used for optional parameters in encrypt and decrypt [[GH-10386](https://github.com/hashicorp/vault/pull/10386)]

## 1.5.5
### October 21, 2020

IMPROVEMENTS:

* auth/aws, core/seal, secret/aws: Set default IMDS timeouts to match AWS SDK [[GH-10133](https://github.com/hashicorp/vault/pull/10133)]

BUG FIXES:

* auth/aws: Restrict region selection when in the aws-us-gov partition to avoid IAM errors [[GH-9947](https://github.com/hashicorp/vault/pull/9947)]
* core (enterprise): Allow operators to add and remove (Raft) peers in a DR secondary cluster using Integrated Storage.
* core (enterprise): Add DR operation token to the remove peer API and CLI command (when DR secondary).
* core (enterprise): Fix deadlock in handling EGP policies
* core (enterprise): Fix extraneous error messages in DR Cluster
* secrets/mysql: Conditionally overwrite TLS parameters for MySQL secrets engine [[GH-9729](https://github.com/hashicorp/vault/pull/9729)]
* secrets/ad: Fix bug where `password_policy` setting was not using correct key when `ad/config` was read [[GH-71](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/71)]
* ui: Fix issue with listing roles and methods on the same auth methods with different names [[GH-10122](https://github.com/hashicorp/vault/pull/10122)]

## 1.5.4
### September 24th, 2020

SECURITY:

* Batch Token Expiry: We addressed an issue where batch token leases could outlive their TTL because we were not scheduling the expiration time correctly. This vulnerability affects Vault OSS and Vault Enterprise 1.0 and newer and is fixed in 1.4.7 and 1.5.4 (CVE-2020-25816).

IMPROVEMENTS:

* secrets/pki: Handle expiration of a cert not in storage as a success [[GH-9880](https://github.com/hashicorp/vault/pull/9880)]
* auth/kubernetes: Add an option to disable defaulting to the local CA cert and service account JWT when running in a Kubernetes pod [[GH-97]](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/97)
* secrets/gcp: Add check for 403 during rollback to prevent repeated deletion calls [[GH-97](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/97)]
* core: Disable usage metrics collection on performance standby nodes. [[GH-9966](https://github.com/hashicorp/vault/pull/9966)]
* credential/aws: Added X-Amz-Content-Sha256 as a default STS request header [[GH-10009](https://github.com/hashicorp/vault/pull/10009)]

BUG FIXES:

* agent: Fix `disable_fast_negotiation` not being set on the auth method when configured by user. [[GH-9892](https://github.com/hashicorp/vault/pull/9892)]
* core (enterprise): Fix hang when cluster-wide plugin reload cleanup is slow on unseal
* core (enterprise): Fix an error in cluster-wide plugin reload cleanup following such a reload
* core: Fix crash when metrics collection encounters zero-length keys in KV store [[GH-9811](https://github.com/hashicorp/vault/pull/9881)]
* mfa (enterprise): Fix incorrect handling of PingID responses that could result in auth requests failing
* replication (enterprise): Improve race condition when using a newly created token on a performance standby node
* replication (enterprise): Only write failover cluster addresses if they've changed
* ui: fix bug where dropdown for identity/entity management is not reflective of actual policy [[GH-9958](https://github.com/hashicorp/vault/pull/9958)]

## 1.5.3
### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust
* secrets/ssh: Fixed a bug with role option for SSH signing algorithm to allow more than RSA signing

## 1.5.2.1
### August 21st, 2020
### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.5.2
### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.5.2 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.5.2)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.


## 1.5.1
### August 20th, 2020

SECURITY:

* When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
* When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)
* When using Vault Agent with cert auto-auth and caching enabled, under certain circumstances, clients without permission to access agent's token may retrieve the token without login credentials. This vulnerability affects Vault Agent 1.1.0 and newer and is fixed in 1.5.1 (CVE-2020-17455)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.5.1 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.5.1)

CHANGES:

* pki: The tidy operation will now remove revoked certificates if the parameter `tidy_revoked_certs` is set to `true`. This will result in certificate entries being immediately removed, as opposed to awaiting until its NotAfter time. Note that this only affects certificates that have been already revoked. [[GH-9609](https://github.com/hashicorp/vault/pull/9609)]
* go: Updated Go version to 1.14.7

IMPROVEMENTS:

* auth/jwt: Add support for fetching groups and user information from G Suite during authentication. [[GH-9574](https://github.com/hashicorp/vault/pull/9574)]
* auth/jwt: Add EdDSA to supported algorithms. [[GH-129](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/129)]
* secrets/openldap: Add "ad" schema that allows the engine to correctly rotate AD passwords. [[GH-9740](https://github.com/hashicorp/vault/pull/9740)]
* pki: Add a `allowed_domains_template` parameter that enables the use of identity templating within the `allowed_domains` parameter. [[GH-8509](https://github.com/hashicorp/vault/pull/8509)]
* secret/azure: Use write-ahead-logs to cleanup any orphaned Service Principals [[GH-9773](https://github.com/hashicorp/vault/pull/9773)]
* ui: Wrap TTL option on transit engine export action is updated to a new component. [[GH-9632](https://github.com/hashicorp/vault/pull/9632)]
* ui: Wrap Tool uses newest version of TTL Picker component. [[GH-9691](https://github.com/hashicorp/vault/pull/9691)]

BUG FIXES:

* secrets/gcp: Ensure that the IAM policy version is appropriately set after a roleset's bindings have changed. [[GH-9603](https://github.com/hashicorp/vault/pull/9603)]
* replication (enterprise): Fix status API output incorrectly stating replication is in `idle` state.
* replication (enterprise): Use PrimaryClusterAddr if it's been set
* core: Fix panic when printing over-long info fields at startup [[GH-9681](https://github.com/hashicorp/vault/pull/9681)]
* core: Seal migration using the new minimal-downtime strategy didn't work properly with performance standbys. [[GH-9690](https://github.com/hashicorp/vault/pull/9690)]
* core: Vault failed to start when there were non-string values in seal configuration [[GH-9555](https://github.com/hashicorp/vault/pull/9555)]
* core: Handle a trailing slash in the API address used for enabling replication

## 1.5.0
### July 21st, 2020

CHANGES:

* audit: Token TTL and issue time are now provided in the auth portion of audit logs. [[GH-9091](https://github.com/hashicorp/vault/pull/9091)]
* auth/gcp: Changes the default name of the entity alias that gets created to be the role ID for both IAM and GCE authentication. [[GH-99](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/99)]
* core: Remove the addition of newlines to parsed configuration when using integer/boolean values [[GH-8928](https://github.com/hashicorp/vault/pull/8928)]
* cubbyhole: Reject reads and writes to an empty ("") path. [[GH-8971](https://github.com/hashicorp/vault/pull/8971)]
* secrets/azure: Default password generation changed from uuid to cryptographically secure randomized string [[GH-40](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/40)]
* storage/gcs: The `credentials_file` config option has been removed. The `GOOGLE_APPLICATION_CREDENTIALS` environment variable
  or default credentials may be used instead [[GH-9424](https://github.com/hashicorp/vault/pull/9424)]
* storage/raft: The storage configuration now accepts a new `max_entry_size` config that will limit
  the total size in bytes of any entry committed via raft. It defaults to `"1048576"` (1MiB). [[GH-9027](https://github.com/hashicorp/vault/pull/9027)]
* token: Token creation with custom token ID via `id` will no longer allow periods (`.`) as part of the input string. 
  The final generated token value may contain periods, such as the `s.` prefix for service token 
  indication. [[GH-8646](https://github.com/hashicorp/vault/pull/8646/files)]
* token: Token renewals will now return token policies within the `token_policies` , identity policies within `identity_policies`, and the full policy set within `policies`. [[GH-8535](https://github.com/hashicorp/vault/pull/8535)]
* go: Updated Go version to 1.14.4

FEATURES:

* **Monitoring**: We have released a Splunk App [9] for Enterprise customers. The app is accompanied by an updated monitoring guide and a few new metrics to enable OSS users to effectively monitor Vault.
* **Password Policies**: Allows operators to customize how passwords are generated for select secret engines (OpenLDAP, Active Directory, Azure, and RabbitMQ).
* **Replication UI Improvements**: We have redesigned the replication UI to highlight the state and relationship between primaries and secondaries and improved management workflows, enabling a more holistic understanding of multiple Vault clusters.
* **Resource Quotas**: As of 1.5, Vault supports specifying a quota to rate limit requests on OSS and Enterprise. Enterprise customers also have access to set quotas on the number of leases that can be generated on a path.
* **OpenShift Support**: We have updated the Helm charts to allow users to install Vault onto their OpenShift clusters.
* **Seal Migration**: We have made updates to allow migrations from auto unseal to Shamir unseal on Enterprise.
* **AWS Auth Web Identity Support**: We've added support for AWS Web Identities, which will be used in the credentials chain if present.
* **Vault Monitor**: Similar to the monitor command for Consul and Nomad, we have added the ability for Vault to stream logs from other Vault servers at varying log levels.
* **AWS Secrets Groups Support**: IAM users generated by Vault may now be added to IAM Groups.
* **Integrated Storage as HA Storage**: In Vault 1.5, it is possible to use Integrated Storage as HA Storage with a different storage backend as regular storage.
* **OIDC Auth Provider Extensions**: We've added support to OIDC Auth to incorporate IdP-specific extensions. Currently this includes expanded Azure AD groups support.
* **GCP Secrets**: Support BigQuery dataset ACLs in absence of IAM endpoints.
* **KMIP**: Add support for signing client certificates requests (CSRs) rather than having them be generated entirely within Vault.

IMPROVEMENTS:

* audit: Replication status requests are no longer audited. [[GH-8877](https://github.com/hashicorp/vault/pull/8877)]
* audit: Added mount_type field to requests and responses. [[GH-9167](https://github.com/hashicorp/vault/pull/9167)]
* auth/aws: Add support for Web Identity credentials [[GH-7738](https://github.com/hashicorp/vault/pull/7738)]
* auth/jwt: Support users that are members of more than 200 groups on Azure [[GH-120](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/120)]
* auth/kerberos: Support identities without userPrincipalName [[GH-44](https://github.com/hashicorp/vault-plugin-auth-kerberos/issues/44)]
* auth/kubernetes: Allow disabling `iss` validation [[GH-91](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/91)]
* auth/kubernetes: Try reading the ca.crt and TokenReviewer JWT from the default service account [[GH-83](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/83)]
* cli: Support reading TLS parameters from file for the `vault operator raft join` command. [[GH-9060](https://github.com/hashicorp/vault/pull/9060)]
* cli: Add a new subcommand, `vault monitor`, for tailing server logs in the console. [[GH-8477](https://github.com/hashicorp/vault/pull/8477)]
* core: Add the Go version used to build a Vault binary to the server message output. [[GH-9078](https://github.com/hashicorp/vault/pull/9078)]
* core: Added Password Policies for user-configurable password generation [[GH-8637](https://github.com/hashicorp/vault/pull/8637)]
* core: New telemetry metrics covering token counts, token creation, KV secret counts, lease creation. [[GH-9239](https://github.com/hashicorp/vault/pull/9239)] [[GH-9250](https://github.com/hashicorp/vault/pull/9250)] [[GH-9244](https://github.com/hashicorp/vault/pull/9244)] [[GH-9052](https://github.com/hashicorp/vault/pull/9052)]
* physical/gcs: The storage backend now uses a dedicated client for HA lock updates to prevent lock table update failures when flooded by other client requests. [[GH-9424](https://github.com/hashicorp/vault/pull/9424)]
* physical/spanner: The storage backend now uses a dedicated client for HA lock updates to prevent lock table update failures when flooded by other client requests. [[GH-9423](https://github.com/hashicorp/vault/pull/9423)]
* plugin: Add SDK method, `Sys.ReloadPlugin`, and CLI command, `vault plugin reload`, for reloading plugins. [[GH-8777](https://github.com/hashicorp/vault/pull/8777)]
* plugin (enterprise): Add a scope field to plugin reload, which when global, reloads the plugin anywhere in a cluster. [[GH-9347](https://github.com/hashicorp/vault/pull/9347)] 
* sdk/framework: Support accepting TypeFloat parameters over the API [[GH-8923](https://github.com/hashicorp/vault/pull/8923)]
* secrets/aws: Add iam_groups parameter to role create/update [[GH-8811](https://github.com/hashicorp/vault/pull/8811)]
* secrets/database: Add static role rotation for MongoDB Atlas database plugin [[GH-11](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/11)]
* secrets/database: Add static role rotation for MSSQL database plugin [[GH-9062](https://github.com/hashicorp/vault/pull/9062)]
* secrets/database: Allow InfluxDB to use insecure TLS without cert bundle [[GH-8778](https://github.com/hashicorp/vault/pull/8778)]
* secrets/gcp: Support BigQuery dataset ACLs in absence of IAM endpoints [[GH-78](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/78)]
* secrets/pki: Allow 3072-bit RSA keys [[GH-8343](https://github.com/hashicorp/vault/pull/8343)]
* secrets/ssh: Add a CA-mode role option to specify signing algorithm [[GH-9096](https://github.com/hashicorp/vault/pull/9096)]
* secrets/ssh: The [Vault SSH Helper](https://github.com/hashicorp/vault-ssh-helper) can now be configured to reference a mount in a namespace [[GH-44](https://github.com/hashicorp/vault-ssh-helper/pull/44)]
* secrets/transit: Transit requests that make use of keys now include a new field  `key_version` in their responses [[GH-9100](https://github.com/hashicorp/vault/pull/9100)]
* secrets/transit: Improving transit batch encrypt and decrypt latencies [[GH-8775](https://github.com/hashicorp/vault/pull/8775)]
* sentinel: Add a sentinel config section, and "additional_enabled_modules", a list of Sentinel modules that may be imported in addition to the defaults.
* ui: Update TTL picker styling on SSH secret engine [[GH-8891](https://github.com/hashicorp/vault/pull/8891)]
* ui: Only render the JWT input field of the Vault login form on mounts configured for JWT auth [[GH-8952](https://github.com/hashicorp/vault/pull/8952)]
* ui: Add replication dashboards.  Improve replication management workflows. [[GH-8705]](https://github.com/hashicorp/vault/pull/8705).
* ui: Update alert banners to match design systems black text. [[GH-9463]](https://github.com/hashicorp/vault/pull/9463).

BUG FIXES:

* auth/oci: Fix issue where users of the Oracle Cloud Infrastructure (OCI) auth method could not authenticate when the plugin backend was mounted at a non-default path. [[GH-7](https://github.com/hashicorp/vault-plugin-auth-oci/pull/7)]
* core: Extend replicated cubbyhole fix in 1.4.0 to cover case where a performance primary is also a DR primary [[GH-9148](https://github.com/hashicorp/vault/pull/9148)]
* replication (enterprise): Use the PrimaryClusterAddr if it's been set
* seal/awskms: fix AWS KMS auto-unseal when AWS_ROLE_SESSION_NAME not set [[GH-9416](https://github.com/hashicorp/vault/pull/9416)]
* sentinel: fix panic due to concurrent map access when rules iterate over metadata maps
* secrets/aws: Fix issue where performance standbys weren't able to generate STS credentials after an IAM access key rotation in AWS and root IAM credential update in Vault [[GH-9186](https://github.com/hashicorp/vault/pull/9186)]
* secrets/database: Fix issue where rotating root database credentials while Vault's storage backend is unavailable causes Vault to lose access to the database [[GH-8782](https://github.com/hashicorp/vault/pull/8782)]
* secrets/database: Fix issue that prevents performance standbys from connecting to databases after a root credential rotation [[GH-9129](https://github.com/hashicorp/vault/pull/9129)]
* secrets/database: Fix parsing of multi-line PostgreSQL statements [[GH-8512](https://github.com/hashicorp/vault/pull/8512)]
* secrets/gcp: Fix issue were updates were not being applied to the `token_scopes` of a roleset. [[GH-90](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/90)]
* secrets/kv: Return the value of delete_version_after when reading kv/config, even if it is set to the default. [[GH-42](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/42)]
* ui: Add Toggle component into core addon so it is available in KMIP and other Ember Engines.[[GH-8913]](https://github.com/hashicorp/vault/pull/8913)
* ui: Disallow max versions value of large than 9999999999999999 on kv2 secrets engine. [[GH-9242](https://github.com/hashicorp/vault/pull/9242)]
* ui: Add and upgrade missing dependencies to resolve a failure with `make static-dist`. [[GH-9277](https://github.com/hashicorp/vault/pull/9371)]

## 1.4.7.1
### October 15th, 2020
### Enterprise Only

BUG FIXES:
* replication (enterprise): Fix panic when old filter path evaluation fails

## 1.4.7
### September 24th, 2020

SECURITY:

* Batch Token Expiry: We addressed an issue where batch token leases could outlive their TTL because we were not scheduling the expiration time correctly. This vulnerability affects Vault OSS and Vault Enterprise 1.0 and newer and is fixed in 1.4.7 and 1.5.4 (CVE-2020-25816).

IMPROVEMENTS:

* secret/azure: Use write-ahead-logs to cleanup any orphaned Service Principals [[GH-9773](https://github.com/hashicorp/vault/pull/9773)]

BUG FIXES:
* replication (enterprise): Don't stop replication if old filter path evaluation fails

## 1.4.6
### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust
* secrets/ssh: Fixed a bug with role option for SSH signing algorithm to allow more than RSA signing [[GH-9824](https://github.com/hashicorp/vault/pull/9824)]

## 1.4.5.1
### August 21st, 2020
### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.4.5
### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.4.5 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.4.5)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.


## 1.4.4
### August 20th, 2020

SECURITY:

* When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
* When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.4.4 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.4.4)

BUG FIXES:

* auth/okta: fix bug introduced in 1.4.0: only 200 external groups were fetched even if user belonged to more [[GH-9580](https://github.com/hashicorp/vault/pull/9580)]
* seal/awskms: fix AWS KMS auto-unseal when AWS_ROLE_SESSION_NAME not set [[GH-9416](https://github.com/hashicorp/vault/pull/9416)]
* secrets/aws: Fix possible issue creating access keys when using Performance Standbys  [[GH-9606](https://github.com/hashicorp/vault/pull/9606)]

IMPROVEMENTS:
* auth/aws: Retry on transient failures during AWS IAM auth login attempts [[GH-8727](https://github.com/hashicorp/vault/pull/8727)]
* ui: Add transit key algorithms aes128-gcm96, ecdsa-p384, ecdsa-p521 to the UI. [[GH-9070](https://github.com/hashicorp/vault/pull/9070)] & [[GH-9520](https://github.com/hashicorp/vault/pull/9520)]

## 1.4.3
### July 2nd, 2020

IMPROVEMENTS:

* auth/aws: Add support for Web Identity credentials [[GH-9251](https://github.com/hashicorp/vault/pull/9251)]
* auth/kerberos: Support identities without userPrincipalName [[GH-44](https://github.com/hashicorp/vault-plugin-auth-kerberos/issues/44)]
* core: Add the Go version used to build a Vault binary to the server message output. [[GH-9078](https://github.com/hashicorp/vault/pull/9078)]
* secrets/database: Add static role rotation for MongoDB Atlas database plugin [[GH-9311](https://github.com/hashicorp/vault/pull/9311)]
* physical/mysql: Require TLS or plaintext flagging in MySQL configuration [[GH-9012](https://github.com/hashicorp/vault/pull/9012)]
* ui: Link to the Vault Changelog in the UI footer [[GH-9216](https://github.com/hashicorp/vault/pull/9216)]

BUG FIXES:

* agent: Restart template server when it shuts down [[GH-9200](https://github.com/hashicorp/vault/pull/9200)]
* auth/oci: Fix issue where users of the Oracle Cloud Infrastructure (OCI) auth method could not authenticate when the plugin backend was mounted at a non-default path. [[GH-9278](https://github.com/hashicorp/vault/pull/9278)]
* replication: The issue causing cubbyholes in namespaces on performance secondaries to not work, which was fixed in 1.4.0, was still an issue when the primary was both a performance primary and DR primary.
* seal: (enterprise) Fix issue causing stored seal and recovery keys to be mistaken as sealwrapped values
* secrets/aws: Fix issue where performance standbys weren't able to generate STS credentials after an IAM access key rotation in AWS and root IAM credential update in Vault [[GH-9207](https://github.com/hashicorp/vault/pull/9207)]
* secrets/database: Fix issue that prevents performance standbys from connecting to databases after a root credential rotation [[GH-9208](https://github.com/hashicorp/vault/pull/9208)]
* secrets/gcp: Fix issue were updates were not being applied to the `token_scopes` of a roleset. [[GH-9277](https://github.com/hashicorp/vault/pull/9277)]


## 1.4.2 (May 21st, 2020)

SECURITY:
* core: Proxy environment variables are now redacted before being logged, in case the URLs include a username:password. This vulnerability, CVE-2020-13223, is fixed in 1.3.6 and 1.4.2, but affects 1.4.0 and 1.4.1, as well as older versions of Vault [[GH-9022](https://github.com/hashicorp/vault/pull/9022)]
* secrets/gcp: Fix a regression in 1.4.0 where the system TTLs were being used instead of the configured backend TTLs for dynamic service accounts. This vulnerability is CVE-2020-12757. [[GH-85](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/85)]

IMPROVEMENTS:

* storage/raft: The storage stanza now accepts `leader_ca_cert_file`, `leader_client_cert_file`, and 
  `leader_client_key_file` parameters to read and parse TLS certificate information from paths on disk.
  Existing non-path based parameters will continue to work, but their values will need to be provided as a 
  single-line string with newlines delimited by `\n`.  [[GH-8894](https://github.com/hashicorp/vault/pull/8894)]
* storage/raft: The `vault status` CLI command and the `sys/leader` API now contain the committed and applied
  raft indexes. [[GH-9011](https://github.com/hashicorp/vault/pull/9011)]
  
BUG FIXES:

* auth/aws: Fix token renewal issues caused by the metadata changes in 1.4.1 [[GH-8991](https://github.com/hashicorp/vault/pull/8991)] 
* auth/ldap: Fix 1.4.0 regression that could result in auth failures when LDAP auth config includes upndomain. [[GH-9041](https://github.com/hashicorp/vault/pull/9041)]
* secrets/ad: Forward rotation requests from standbys to active clusters [[GH-66](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/66)]
* secrets/database: Prevent generation of usernames that are not allowed by the MongoDB Atlas API [[GH-9](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/9)]
* secrets/database: Return an error if a manual rotation of static account credentials fails [[GH-9035](https://github.com/hashicorp/vault/pull/9035)]
* secrets/openldap: Forward all rotation requests from standbys to active clusters [[GH-9028](https://github.com/hashicorp/vault/pull/9028)]
* secrets/transform (enterprise): Fix panic that could occur when accessing cached template entries, such as a requests
  that accessed templates directly or indirectly from a performance standby node.
* serviceregistration: Fix a regression for Consul service registration that ignored using the listener address as
  the redirect address unless api_addr was provided. It now properly uses the same redirect address as the one
  used by Vault's Core object. [[GH-8976](https://github.com/hashicorp/vault/pull/8976)]
* storage/raft: Advertise the configured cluster address to the rest of the nodes in the raft cluster. This fixes
  an issue where a node advertising 0.0.0.0 is not using a unique hostname. [[GH-9008](https://github.com/hashicorp/vault/pull/9008)]
* storage/raft: Fix panic when multiple nodes attempt to join the cluster at once. [[GH-9008](https://github.com/hashicorp/vault/pull/9008)]
* sys: The path provided in `sys/internal/ui/mounts/:path` is now namespace-aware. This fixes an issue
  with `vault kv` subcommands that had namespaces provided in the path returning permission denied all the time.
  [[GH-8962](https://github.com/hashicorp/vault/pull/8962)]
* ui: Fix snowman that appears when namespaces have more than one period [[GH-8910](https://github.com/hashicorp/vault/pull/8910)]

## 1.4.1 (April 30th, 2020)

CHANGES: 

* auth/aws: The default set of metadata fields added in 1.4.1 has been changed to `account_id` and `auth_type` [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]
* storage/raft: Disallow `ha_storage` to be specified if `raft` is set as the `storage` type. [[GH-8707](https://github.com/hashicorp/vault/pull/8707)]

IMPROVEMENTS:

* auth/aws: The set of metadata stored during login is now configurable [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]
* auth/aws: Improve region selection to avoid errors seen if the account hasn't enabled some newer AWS regions [[GH-8679](https://github.com/hashicorp/vault/pull/8679)]
* auth/azure: Enable login from Azure VMs with user-assigned identities [[GH-33](https://github.com/hashicorp/vault-plugin-auth-azure/pull/33)]
* auth/gcp: The set of metadata stored during login is now configurable [[GH-92](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/92)]
* auth/gcp: The type of alias name used during login is now configurable [[GH-95](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/95)]
* auth/ldap: Improve error messages during LDAP operation failures [[GH-8740](https://github.com/hashicorp/vault/pull/8740)]
* identity: Add a batch delete API for identity entities [[GH-8785]](https://github.com/hashicorp/vault/pull/8785)
* identity: Improve performance of logins when no group updates are needed [[GH-8795]](https://github.com/hashicorp/vault/pull/8795)
* metrics: Add `vault.identity.num_entities` metric [[GH-8816]](https://github.com/hashicorp/vault/pull/8816)
* secrets/kv: Allow `delete-version-after` to be reset to 0 via the CLI [[GH-8635](https://github.com/hashicorp/vault/pull/8635)]
* secrets/rabbitmq: Improve error handling and reporting [[GH-8619](https://github.com/hashicorp/vault/pull/8619)]
* ui: Provide One Time Password during Operation Token generation process [[GH-8630]](https://github.com/hashicorp/vault/pull/8630)

BUG FIXES:

* auth/okta: Fix MFA regression (introduced in [GH-8143](https://github.com/hashicorp/vault/pull/8143)) from 1.4.0 [[GH-8807](https://github.com/hashicorp/vault/pull/8807)]
* auth/userpass: Fix upgrade value for `token_bound_cidrs` being ignored due to incorrect key provided [[GH-8826](https://github.com/hashicorp/vault/pull/8826/files)]
* config/seal: Fix segfault when seal block is removed [[GH-8517](https://github.com/hashicorp/vault/pull/8517)]
* core: Fix an issue where users attempting to build Vault could receive Go module checksum errors [[GH-8770](https://github.com/hashicorp/vault/pull/8770)]
* core: Fix blocked requests if a SIGHUP is issued during a long-running request has the state lock held. 
  Also fixes deadlock that can happen if `vault debug` with the config target is ran during this time.
  [[GH-8755](https://github.com/hashicorp/vault/pull/8755)]
* core: Always rewrite the .vault-token file as part of a `vault login` to ensure permissions and ownership are set correctly [[GH-8867](https://github.com/hashicorp/vault/pull/8867)]
* database/mongodb: Fix context deadline error that may result due to retry attempts on failed commands
  [[GH-8863](https://github.com/hashicorp/vault/pull/8863)]
* http: Fix superflous call messages from the http package on logs caused by missing returns after
  `respondError` calls [[GH-8796](https://github.com/hashicorp/vault/pull/8796)]
* namespace (enterprise): Fix namespace listing to return `key_info` when a scoping namespace is also provided. 
* seal/gcpkms: Fix panic that could occur if all seal parameters were provided via environment
  variables [[GH-8840](https://github.com/hashicorp/vault/pull/8840)]
* storage/raft: Fix memory allocation and incorrect metadata tracking issues with snapshots [[GH-8793](https://github.com/hashicorp/vault/pull/8793)]
* storage/raft: Fix panic that could occur if `disable_clustering` was set to true on Raft storage cluster [[GH-8784](https://github.com/hashicorp/vault/pull/8784)]
* storage/raft: Handle errors returned from the API during snapshot operations [[GH-8861](https://github.com/hashicorp/vault/pull/8861)]
* sys/wrapping: Allow unwrapping of wrapping tokens which contain nil data [[GH-8714](https://github.com/hashicorp/vault/pull/8714)]

## 1.4.0 (April 7th, 2020)

CHANGES:

* cli: The raft configuration command has been renamed to list-peers to avoid
  confusion.

FEATURES:

* **Kerberos Authentication**: Vault now supports Kerberos authentication using a SPNEGO token. 
   Login can be performed using the Vault CLI, API, or agent.
* **Kubernetes Service Discovery**: A new Kubernetes service discovery feature where, if 
   configured, Vault will tag Vault pods with their current health status. For more, see [#8249](https://github.com/hashicorp/vault/pull/8249).
* **MongoDB Atlas Secrets**: Vault can now generate dynamic credentials for both MongoDB Atlas databases
  as well as the [Atlas programmatic interface](https://docs.atlas.mongodb.com/tutorial/manage-programmatic-access/).
* **OpenLDAP Secrets Engine**: We now support password management of existing OpenLDAP user entries. For more, see [#8360](https://github.com/hashicorp/vault/pull/8360/).
* **Redshift Database Secrets Engine**: The database secrets engine now supports static and dynamic secrets for the Amazon Web Services (AWS) Redshift service.
* **Service Registration Config**: A newly introduced `service_registration` configuration stanza, that allows for service registration to be configured separately from the storage backend. For more, see [#7887](https://github.com/hashicorp/vault/pull/7887/).
* **Transform Secrets Engine (Enterprise)**: A new secrets engine that handles secure data transformations against provided input values.
* **Integrated Storage**: Promoted out of beta and into general availability for both open-source and enterprise workloads.

IMPROVEMENTS:

* agent: add option to force the use of the auth-auth token, and ignore the Vault token in the request [[GH-8101](https://github.com/hashicorp/vault/pull/8101)]
* api: Restore and fix DNS SRV Lookup [[GH-8520](https://github.com/hashicorp/vault/pull/8520)]
* audit: HMAC http_raw_body in audit log; this ensures that large authenticated Prometheus metrics responses get
  replaced with short HMAC values [[GH-8130](https://github.com/hashicorp/vault/pull/8130)]
* audit: Generate-root, generate-recovery-token, and generate-dr-operation-token requests and responses are now audited. [[GH-8301](https://github.com/hashicorp/vault/pull/8301)]
* auth/aws: Reduce the number of simultaneous STS client credentials needed  [[GH-8161](https://github.com/hashicorp/vault/pull/8161)]
* auth/azure: subscription ID, resource group, vm and vmss names are now stored in alias metadata [[GH-30](https://github.com/hashicorp/vault-plugin-auth-azure/pull/30)]
* auth/jwt: Additional OIDC callback parameters available for CLI logins [[GH-80](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/80) & [GH-86](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/86)]
* auth/jwt: Bound claims may be optionally configured using globs [[GH-89](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/89)]
* auth/jwt: Timeout during OIDC CLI login if process doesn't complete within 2 minutes [[GH-97](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/97)]
* auth/jwt: Add support for the `form_post` response mode [[GH-98](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/98)]
* auth/jwt: add optional client_nonce to authorization flow [[GH-104](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/104)]
* auth/okta: Upgrade okta sdk lib, which should improve handling of groups [[GH-8143](https://github.com/hashicorp/vault/pull/8143)]
* aws: Add support for v2 of the instance metadata service (see [issue 7924](https://github.com/hashicorp/vault/issues/7924) for all linked PRs)
* core: Separate out service discovery interface from storage interface to allow
  new types of service discovery not coupled to storage [[GH-7887](https://github.com/hashicorp/vault/pull/7887)]
* core: Add support for telemetry option `metrics_prefix` [[GH-8340](https://github.com/hashicorp/vault/pull/8340)]
* core: Entropy Augmentation can now be used with AWS KMS and Vault Transit seals
* core: Allow tls_min_version to be set to TLS 1.3 [[GH-8305](https://github.com/hashicorp/vault/pull/8305)]
* cli: Incorrect TLS configuration will now correctly fail [[GH-8025](https://github.com/hashicorp/vault/pull/8025)]
* identity: Allow specifying a custom `client_id` for identity tokens [[GH-8165](https://github.com/hashicorp/vault/pull/8165)]
* metrics/prometheus: improve performance with high volume of metrics updates [[GH-8507](https://github.com/hashicorp/vault/pull/8507)]
* replication (enterprise): Fix race condition causing clusters with high throughput writes to sometimes
  fail to enter streaming-wal mode
* replication (enterprise): Secondary clusters can now perform an extra gRPC call to all nodes in a primary 
  cluster in an attempt to resolve the active node's address
* replication (enterprise): The replication status API now outputs `last_performance_wal`, `last_dr_wal`, 
  and `connection_state` values
* replication (enterprise): DR secondary clusters can now be recovered by the `replication/dr/secondary/recover`
  API
* replication (enterprise): We now allow for an alternate means to create a Disaster Recovery token, by using a batch
  token that is created with an ACL that allows for access to one or more of the DR endpoints.
* secrets/database/mongodb: Switched internal MongoDB driver to mongo-driver [[GH-8140](https://github.com/hashicorp/vault/pull/8140)]
* secrets/database/mongodb: Add support for x509 client authorization to MongoDB [[GH-8329](https://github.com/hashicorp/vault/pull/8329)]
* secrets/database/oracle: Add support for static credential rotation [[GH-26](https://github.com/hashicorp/vault-plugin-database-oracle/pull/26)]
* secrets/consul: Add support to specify TLS options per Consul backend [[GH-4800](https://github.com/hashicorp/vault/pull/4800)]
* secrets/gcp: Allow specifying the TTL for a service key [[GH-54](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/54)]
* secrets/gcp: Add support for rotating root keys [[GH-53](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/53)]
* secrets/gcp: Handle version 3 policies for Resource Manager IAM requests [[GH-77](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/77)]
* secrets/nomad: Add support to specify TLS options per Nomad backend [[GH-8083](https://github.com/hashicorp/vault/pull/8083)]
* secrets/ssh: Allowed users can now be templated with identity information [[GH-7548](https://github.com/hashicorp/vault/pull/7548)]
* secrets/transit: Adding RSA3072 key support [[GH-8151](https://github.com/hashicorp/vault/pull/8151)]
* storage/consul: Vault returns now a more descriptive error message when only a client cert or
  a client key has been provided [[GH-4930]](https://github.com/hashicorp/vault/pull/8084)
* storage/raft: Nodes in the raft cluster can all be given possible leader
  addresses for them to continuously try and join one of them, thus automating
  the process of join to a greater extent [[GH-7856](https://github.com/hashicorp/vault/pull/7856)]
* storage/raft: Fix a potential deadlock that could occur on leadership transition [[GH-8547](https://github.com/hashicorp/vault/pull/8547)]
* storage/raft: Refresh TLS keyring on snapshot restore [[GH-8546](https://github.com/hashicorp/vault/pull/8546)]
* storage/etcd: Bumped etcd client API SDK [[GH-7931](https://github.com/hashicorp/vault/pull/7931) & [GH-4961](https://github.com/hashicorp/vault/pull/4961) & [GH-4349](https://github.com/hashicorp/vault/pull/4349) & [GH-7582](https://github.com/hashicorp/vault/pull/7582)]
* ui: Make Transit Key actions more prominent [[GH-8304](https://github.com/hashicorp/vault/pull/8304)]
* ui: Add Core Usage Metrics [[GH-8347](https://github.com/hashicorp/vault/pull/8347)]
* ui: Add refresh Namespace list on the Namespace dropdown, and redesign of Namespace dropdown menu [[GH-8442](https://github.com/hashicorp/vault/pull/8442)]
* ui: Update transit actions to codeblocks & automatically encode plaintext unless indicated [[GH-8462](https://github.com/hashicorp/vault/pull/8462)]
* ui: Display the results of transit key actions in a modal window [[GH-8462](https://github.com/hashicorp/vault/pull/8575)]
* ui: Transit key version styling updates & ability to copy key from dropdown [[GH-8480](https://github.com/hashicorp/vault/pull/8480)]

BUG FIXES:

* agent: Fix issue where TLS options are ignored for agent template feature [[GH-7889](https://github.com/hashicorp/vault/pull/7889)]
* auth/jwt: Use lower case role names for `default_role` to match the `role` case convention [[GH-100](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/100)]
* auth/ldap: Fix a bug where the UPNDOMAIN parameter was wrongly used to lookup the group
  membership of the given user [[GH-6325]](https://github.com/hashicorp/vault/pull/8333)
* cli: Support autocompletion for nested mounts [[GH-8303](https://github.com/hashicorp/vault/pull/8303)]
* cli: Fix CLI namespace autocompletion [[GH-8315](https://github.com/hashicorp/vault/pull/8315)]
* identity: Fix incorrect caching of identity token JWKS responses [[GH-8412](https://github.com/hashicorp/vault/pull/8412)]
* metrics/stackdriver: Fix issue that prevents the stackdriver metrics library to create unnecessary stackdriver descriptors [[GH-8073](https://github.com/hashicorp/vault/pull/8073)]
* replication (enterprise): Fix issue causing cubbyholes in namespaces on performance secondaries to not work.
* replication (enterprise): Unmounting a dynamic secrets backend could sometimes lead to replication errors.  Change the order of operations to prevent that.
* seal (enterprise): Fix seal migration when transactional seal wrap backend is in use.
* secrets/database/influxdb: Fix potential panic if connection to the InfluxDB database cannot be established [[GH-8282](https://github.com/hashicorp/vault/pull/8282)]
* secrets/database/mysql: Ensures default static credential rotation statements are used [[GH-8240](https://github.com/hashicorp/vault/pull/8240)]
* secrets/database/mysql: Fix inconsistent query parameter names: {{name}} or {{username}} for
  different queries. Now it allows for either for backwards compatibility [[GH-8240](https://github.com/hashicorp/vault/pull/8240)]
* secrets/database/postgres: Fix inconsistent query parameter names: {{name}} or {{username}} for
  different queries. Now it allows for either for backwards compatibility [[GH-8240](https://github.com/hashicorp/vault/pull/8240)]
* secrets/pki: Support FQDNs in DNS Name [[GH-8288](https://github.com/hashicorp/vault/pull/8288)]
* storage/raft: Allow seal migration to be performed on Vault clusters using raft storage [[GH-8103](https://github.com/hashicorp/vault/pull/8103)]
* telemetry: Prometheus requests on standby nodes will now return an error instead of forwarding 
  the request to the active node [[GH-8280](https://github.com/hashicorp/vault/pull/8280)]
* ui: Fix broken popup menu on the transit secrets list page [[GH-8348](https://github.com/hashicorp/vault/pull/8348)]
* ui: Update headless Chrome flag to fix `yarn run test:oss` [[GH-8035](https://github.com/hashicorp/vault/pull/8035)]
* ui: Update CLI to accept empty strings as param value to reset previously-set values
* ui: Fix bug where error states don't clear when moving between action tabs on Transit [[GH-8354](https://github.com/hashicorp/vault/pull/8354)]

## 1.3.10
### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust

## 1.3.9.1
### August 21st, 2020
### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.3.9
### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.3.9 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.3.9)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.

## 1.3.8
### August 20th, 2020

SECURITY:

* When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
* When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.3.8 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.3.8)

## 1.3.7
### July 2nd, 2020

BUG FIXES:

* seal: (enterprise) Fix issue causing stored seal and recovery keys to be mistaken as sealwrapped values
* secrets/aws: Fix issue where performance standbys weren't able to generate STS credentials after an IAM access key rotation in AWS and root IAM credential update in Vault [[GH-9363](https://github.com/hashicorp/vault/pull/9363)]

## 1.3.6 (May 21st, 2020)

SECURITY:
* core: proxy environment variables are now redacted before being logged, in case the URLs include a username:password. This vulnerability, CVE-2020-13223, is fixed in 1.3.6 and 1.4.2, but affects 1.4 and 1.4.1, as well as older versions of Vault [[GH-9022](https://github.com/hashicorp/vault/pull/9022)]

BUG FIXES:

* auth/aws: Fix token renewal issues caused by the metadata changes in 1.3.5 [[GH-8991](https://github.com/hashicorp/vault/pull/8991)] 
* replication: Fix mount filter bug that allowed replication filters to hide local mounts on a performance secondary

## 1.3.5 (April 28th, 2020)

CHANGES: 

* auth/aws: The default set of metadata fields added in 1.3.2 has been changed to `account_id` and `auth_type` [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]

IMPROVEMENTS:

* auth/aws: The set of metadata stored during login is now configurable [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]

## 1.3.4 (March 19th, 2020)

SECURITY:

* A vulnerability was identified in Vault and Vault Enterprise such that, under certain circumstances,  an Entity's Group membership may inadvertently include Groups the Entity no longer has permissions to. This vulnerability, CVE-2020-10660, affects Vault and Vault Enterprise versions 0.9.0 and newer, and is fixed in 1.3.4. [[GH-8606](https://github.com/hashicorp/vault/pull/8606)]
* A vulnerability was identified in Vault Enterprise such that, under certain circumstances, existing nested-path policies may give access to Namespaces created after-the-fact. This vulnerability, CVE-2020-10661, affects Vault Enterprise versions 0.11 and newer, and is fixed in 1.3.4.

## 1.3.3 (March 5th, 2020)

BUG FIXES:

* approle: Fix excessive locking during tidy, which could potentially block new approle logins for long enough to cause an outage [[GH-8418](https://github.com/hashicorp/vault/pull/8418)]
* cli: Fix issue where Raft snapshots from standby nodes created an empty backup file [[GH-8097](https://github.com/hashicorp/vault/pull/8097)]
* identity: Fix incorrect caching of identity token JWKS responses [[GH-8412](https://github.com/hashicorp/vault/pull/8412)]
* kmip: role read now returns tls_client_ttl
* kmip: fix panic when templateattr not provided in rekey request
* secrets/database/influxdb: Fix potential panic if connection to the InfluxDB database cannot be established [[GH-8282](https://github.com/hashicorp/vault/pull/8282)]
* storage/mysql: Fix potential crash when using MySQL as coordination for high availability [[GH-8300](https://github.com/hashicorp/vault/pull/8300)]
* storage/raft: Fix potential crash when using Raft as coordination for high availability [[GH-8356](https://github.com/hashicorp/vault/pull/8356)]
* ui: Fix missing License menu item [[GH-8230](https://github.com/hashicorp/vault/pull/8230)]
* ui: Fix bug where default auth method on login is defaulted to auth method that is listing-visibility=unauth instead of “other” [[GH-8218](https://github.com/hashicorp/vault/pull/8218)]
* ui: Fix bug where KMIP details were not shown in the UI Wizard [[GH-8255](https://github.com/hashicorp/vault/pull/8255)]
* ui: Show Error messages on Auth Configuration page when you hit permission errors [[GH-8500](https://github.com/hashicorp/vault/pull/8500)]
* ui: Remove duplicate form inputs for the GitHub config [[GH-8519](https://github.com/hashicorp/vault/pull/8519)]
* ui: Correct HMAC capitalization [[GH-8528](https://github.com/hashicorp/vault/pull/8528)]
* ui: Fix danger message in DR [[GH-8555](https://github.com/hashicorp/vault/pull/8555)]
* ui: Fix certificate field for LDAP config [[GH-8573](https://github.com/hashicorp/vault/pull/8573)]

## 1.3.2 (January 22nd, 2020)

SECURITY:
 * When deleting a namespace on Vault Enterprise, in certain circumstances, the deletion
   process will fail to revoke dynamic secrets for a mount in that namespace. This will
   leave any dynamic secrets in remote systems alive and will fail to clean them up. This
   vulnerability, CVE-2020-7220, affects Vault Enterprise 0.11.0 and newer.

IMPROVEMENTS:
 * auth/aws: Add aws metadata to identity alias [[GH-7985](https://github.com/hashicorp/vault/pull/7985)]
 * auth/kubernetes: Allow both names and namespaces to be set to "*" [[GH-78](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/78)]

BUG FIXES:

* auth/azure: Fix Azure compute client to use correct base URL [[GH-8072](https://github.com/hashicorp/vault/pull/8072)]
* auth/ldap: Fix renewal of tokens without configured policies that are
  generated by an LDAP login [[GH-8072](https://github.com/hashicorp/vault/pull/8072)]
* auth/okta: Fix renewal of tokens without configured policies that are
  generated by an Okta login [[GH-8072](https://github.com/hashicorp/vault/pull/8072)]
* core: Fix seal migration error when attempting to migrate from auto unseal to shamir [[GH-8172](https://github.com/hashicorp/vault/pull/8172)]
* core: Fix seal migration config issue when migrating from auto unseal to auto unseal [[GH-8172](https://github.com/hashicorp/vault/pull/8172)]
* plugin: Fix issue where a plugin unwrap request potentially used an expired token [[GH-8058](https://github.com/hashicorp/vault/pull/8058)]
* replication: Fix issue where a forwarded request from a performance/standby node could run into
  a timeout
* secrets/database: Fix issue where a manual static role rotation could potentially panic [[GH-8098](https://github.com/hashicorp/vault/pull/8098)]
* secrets/database: Fix issue where a manual root credential rotation request is not forwarded
  to the primary node [[GH-8125](https://github.com/hashicorp/vault/pull/8125)]
* secrets/database: Fix issue where a manual static role rotation request is not forwarded
  to the primary node [[GH-8126](https://github.com/hashicorp/vault/pull/8126)]
* secrets/database/mysql: Fix issue where special characters for a MySQL password were encoded [[GH-8040](https://github.com/hashicorp/vault/pull/8040)]
* ui: Fix deleting namespaces [[GH-8132](https://github.com/hashicorp/vault/pull/8132)]
* ui: Fix Error handler on kv-secret edit and kv-secret view pages [[GH-8133](https://github.com/hashicorp/vault/pull/8133)]
* ui: Fix OIDC callback to check storage [[GH-7929](https://github.com/hashicorp/vault/pull/7929)].
* ui: Change `.box-radio` height to min-height to prevent overflow issues [[GH-8065](https://github.com/hashicorp/vault/pull/8065)]

## 1.3.1 (December 18th, 2019)

IMPROVEMENTS:

* agent: Add ability to set `exit-after-auth` via the CLI [[GH-7920](https://github.com/hashicorp/vault/pull/7920)]
* auth/ldap: Add a `request_timeout` configuration option to prevent connection
  requests from hanging [[GH-7909](https://github.com/hashicorp/vault/pull/7909)]
* auth/kubernetes: Add audience to tokenreview API request for Kube deployments where issuer
  is not Kube. [[GH-74](https://github.com/hashicorp/vault/pull/74)]
* secrets/ad: Add a `request_timeout` configuration option to prevent connection
  requests from hanging [[GH-59](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/59)]
* storage/postgresql: Add support for setting `connection_url` from enviornment
  variable `VAULT_PG_CONNECTION_URL` [[GH-7937](https://github.com/hashicorp/vault/pull/7937)]
* telemetry: Add `enable_hostname_label` option to telemetry stanza [[GH-7902](https://github.com/hashicorp/vault/pull/7902)]
* telemetry: Add accept header check for prometheus mime type [[GH-7958](https://github.com/hashicorp/vault/pull/7958)]

BUG FIXES:

* agent: Fix issue where Agent exits before all templates are rendered when
  using and `exit_after_auth` [[GH-7899](https://github.com/hashicorp/vault/pull/7899)]
* auth/aws: Fixes region-related issues when using a custom `sts_endpoint` by adding
  a `sts_region` parameter [[GH-7922](https://github.com/hashicorp/vault/pull/7922)]
* auth/token: Fix panic when getting batch tokens on a performance standby from a role
  that does not exist [[GH-8027](https://github.com/hashicorp/vault/pull/8027)]
* core: Improve warning message for lease TTLs [[GH-7901](https://github.com/hashicorp/vault/pull/7901)]
* identity: Fix identity token panic during invalidation [[GH-8043](https://github.com/hashicorp/vault/pull/8043)]
* plugin: Fix a panic that could occur if a mount/auth entry was unable to
  mount the plugin backend and a request that required the system view to be
  retrieved was made [[GH-7991](https://github.com/hashicorp/vault/pull/7991)]
* replication: Add `generate-public-key` endpoint to list of allowed endpoints
  for existing DR secondaries
* secrets/gcp: Fix panic if bindings aren't provided in roleset create/update. [[GH-56](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/56)]
* secrets/pki: Prevent generating certificate on performance standby when storing
  [[GH-7904](https://github.com/hashicorp/vault/pull/7904)]
* secrets/transit: Prevent restoring keys to new names that are sub paths [[GH-7998](https://github.com/hashicorp/vault/pull/7998)]
* storage/s3: Fix a bug in configurable S3 paths that was preventing use of S3 as
  a source during `operator migrate` operations [[GH-7966](https://github.com/hashicorp/vault/pull/7966)]
* ui: Ensure secrets with a period in their key can be viewed and copied [[GH-7926](https://github.com/hashicorp/vault/pull/7926)]
* ui: Fix status menu after demotion [[GH-7997](https://github.com/hashicorp/vault/pull/7997)]
* ui: Fix select dropdowns in Safari when running Mojave [[GH-8023](https://github.com/hashicorp/vault/pull/8023)]

## 1.3 (November 14th, 2019)

CHANGES:

 * Secondary cluster activation: There has been a change to the way that activating
   performance and DR secondary clusters works when using public keys for
   encryption of the parameters rather than a wrapping token. This flow was
   experimental and never documented. It is now officially supported and
   documented but is not backwards compatible with older Vault releases.
 * Cluster cipher suites: On its cluster port, Vault will no longer advertise
   the full TLS 1.2 cipher suite list by default. Although this port is only
   used for Vault-to-Vault communication and would always pick a strong cipher,
   it could cause false flags on port scanners and other security utilities
   that assumed insecure ciphers were being used. The previous behavior can be
   achieved by setting the value of the (undocumented) `cluster_cipher_suites`
   config flag to `tls12`.
 * API/Agent Renewal behavior: The API now allows multiple options for how it
   deals with renewals. The legacy behavior in the Agent/API is for the renewer
   (now called the lifetime watcher) to exit on a renew error, leading to a
   reauthentication. The new default behavior is for the lifetime watcher to
   ignore 5XX errors and simply retry as scheduled, using the existing lease
   duration. It is also possible, within custom code, to disable renewals
   entirely, which allows the lifetime watcher to simply return when it
   believes it is time for your code to renew or reauthenticate.

FEATURES:

 * **Vault Debug**: A new top-level subcommand, `debug`, is added that allows
   operators to retrieve debugging information related to a particular Vault
   node. Operators can use this simple workflow to capture triaging information,
   which can then be consumed programmatically or by support and engineering teams.
   It has the abilitity to probe for config, host, metrics, pprof, server status,
   and replication status.
 * **Recovery Mode**: Vault server can be brought up in recovery mode to resolve
   outages caused due to data store being in bad state. This is a privileged mode
   that allows `sys/raw` API calls to perform surgical corrections to the data
   store. Bad storage state can be caused by bugs. However, this is usually
   observed when known (and fixed) bugs are hit by older versions of Vault.
 * **Entropy Augmentation (Enterprise)**: Vault now supports sourcing entropy from
   external source for critical security parameters. Currently an HSM that
   supports PKCS#11 is the only supported source.
 * **Active Directory Secret Check-In/Check-Out**: In the Active Directory secrets
   engine, users or applications can check out a service account for use, and its
   password will be rotated when it's checked back in.
 * **Vault Agent Template**: Vault Agent now supports rendering templates containing
   Vault secrets to disk, similar to Consul Template [[GH-7652](https://github.com/hashicorp/vault/pull/7652)]
 * **Transit Key Type Support**: Signing and verification is now supported with the P-384
   (secp384r1) and P-521 (secp521r1) ECDSA curves [[GH-7551](https://github.com/hashicorp/vault/pull/7551)] and encryption and
   decryption is now supported via AES128-GCM96 [[GH-7555](https://github.com/hashicorp/vault/pull/7555)]
 * **SSRF Protection for Vault Agent**: Vault Agent has a configuration option to
   require a specific header before allowing requests [[GH-7627](https://github.com/hashicorp/vault/pull/7627)]
 * **AWS Auth Method Root Rotation**: The credential used by the AWS auth method can
   now be rotated, to ensure that only Vault knows the credentials it is using [[GH-7131](https://github.com/hashicorp/vault/pull/7131)]
 * **New UI Features**: The UI now supports managing users and groups for the
   Userpass, Cert, Okta, and Radius auth methods.
 * **Shamir with Stored Master Key**: The on disk format for Shamir seals has changed,
   allowing for a secondary cluster using Shamir downstream from a primary cluster
   using Auto Unseal. [[GH-7694](https://github.com/hashicorp/vault/pull/7694)]
 * **Stackdriver Metrics Sink**: Vault can now send metrics to
   [Stackdriver](https://cloud.google.com/stackdriver/). See the [configuration
   documentation](https://www.vaultproject.io/docs/config/index.html) for
   details. [[GH-6957](https://github.com/hashicorp/vault/pull/6957)]
 * **Filtered Paths Replication (Enterprise)**: Based on the predecessor Filtered Mount Replication,
   Filtered Paths Replication allows now filtering of namespaces in addition to mounts.
   With this feature, Filtered Mount Replication should be considered deprecated.
 * **Token Renewal via Accessor**: Tokens can now be renewed via the accessor value through
   the new `auth/token/renew-accessor` endpoint if the caller's token has
   permission to access that endpoint.
 * **Improved Integrated Storage (Beta)**: Improved raft write performance, added support for
   non-voter nodes, along with UI support for: using raft storage, joining a raft cluster,
   and downloading and restoring a snapshot.

IMPROVEMENTS:

 * agent: Add ability to set the TLS SNI name used by Agent [[GH-7519](https://github.com/hashicorp/vault/pull/7519)]
 * agent & api: Change default renewer behavior to ignore 5XX errors [[GH-7733](https://github.com/hashicorp/vault/pull/7733)]
 * auth/jwt: The redirect callback host may now be specified for CLI logins
   [[GH-71](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/71)]
 * auth/jwt: Bound claims may now contain boolean values [[GH-73](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/73)]
 * auth/jwt: CLI logins can now open the browser when running in WSL [[GH-77](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/77)]
 * core: Exit ScanView if context has been cancelled [[GH-7419](https://github.com/hashicorp/vault/pull/7419)]
 * core: re-encrypt barrier and recovery keys if the unseal key is updated
   [[GH-7493](https://github.com/hashicorp/vault/pull/7493)]
 * core: Don't advertise the full set of TLS 1.2 cipher suites on the cluster
   port, even though only strong ciphers were used [[GH-7487](https://github.com/hashicorp/vault/pull/7487)]
 * core (enterprise): Add background seal re-wrap
 * core/metrics: Add config parameter to allow unauthenticated sys/metrics
   access. [[GH-7550](https://github.com/hashicorp/vault/pull/7550)]
 * metrics: Upgrade DataDog library to improve performance [[GH-7794](https://github.com/hashicorp/vault/pull/7794)]
 * replication (enterprise): Write-Ahead-Log entries will not duplicate the
   data belonging to the encompassing physical entries of the transaction,
   thereby improving the performance and storage capacity.
 * replication (enterprise): Added more replication metrics
 * replication (enterprise): Reindex process now compares subpages for a more
   accurate indexing process.
 * replication (enterprise): Reindex API now accepts a new `skip_flush`
   parameter indicating all the changes should not be flushed while the tree is
   locked.
 * secrets/aws: The root config can now be read [[GH-7245](https://github.com/hashicorp/vault/pull/7245)]
 * secrets/aws: Role paths may now contain the '@' character [[GH-7553](https://github.com/hashicorp/vault/pull/7553)]
 * secrets/database/cassandra: Add ability to skip verfication of connection
   [[GH-7614](https://github.com/hashicorp/vault/pull/7614)]
 * secrets/gcp: Fix panic during rollback if the roleset has been deleted
   [[GH-52](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/52)]
 * storage/azure: Add config parameter to Azure storage backend to allow
   specifying the ARM endpoint [[GH-7567](https://github.com/hashicorp/vault/pull/7567)]
 * storage/cassandra: Improve storage efficiency by eliminating unnecessary
   copies of value data [[GH-7199](https://github.com/hashicorp/vault/pull/7199)]
 * storage/raft: Improve raft write performance by utilizing FSM Batching
   [[GH-7527](https://github.com/hashicorp/vault/pull/7527)]
 * storage/raft: Add support for non-voter nodes [[GH-7634](https://github.com/hashicorp/vault/pull/7634)]
 * sys: Add a new `sys/host-info` endpoint for querying information about
   the host [[GH-7330](https://github.com/hashicorp/vault/pull/7330)]
 * sys: Add a new set of endpoints under `sys/pprof/` that allows profiling
   information to be extracted [[GH-7473](https://github.com/hashicorp/vault/pull/7473)]
 * sys: Add endpoint that counts the total number of active identity entities
   [[GH-7541](https://github.com/hashicorp/vault/pull/7541)]
 * sys: `sys/seal-status` now has a `storage_type` field denoting what type of
   storage
   the cluster is configured to use
 * sys: Add a new `sys/internal/counters/tokens` endpoint, that counts the
   total number of active service token accessors in the shared token storage.
   [[GH-7541](https://github.com/hashicorp/vault/pull/7541)]
 * sys/config: Add  a new endpoint under `sys/config/state/sanitized` that
   returns the configuration state of the server. It excludes config values
   from `storage`, `ha_storage`, and `seal` stanzas and some values
   from `telemetry` due to potential sensitive entries in those fields.
 * ui: when using raft storage, you can now join a raft cluster, download a
   snapshot, and restore a snapshot from the UI [[GH-7410](https://github.com/hashicorp/vault/pull/7410)]
 * ui: clarify when secret version is deleted in the secret version history
   dropdown [[GH-7714](https://github.com/hashicorp/vault/pull/7714)]

BUG FIXES:

 * agent: Fix a data race on the token value for inmemsink [[GH-7707](https://github.com/hashicorp/vault/pull/7707)]
 * api: Fix Go API using lease revocation via URL instead of body [[GH-7777](https://github.com/hashicorp/vault/pull/7777)]
 * api: Allow setting a function to control retry behavior [[GH-7331](https://github.com/hashicorp/vault/pull/7331)]
 * auth/gcp: Fix a bug where region information in instance groups names could
   cause an authorization attempt to fail [[GH-74](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/74)]
 * cli: Fix a bug where a token of an unknown format (e.g. in ~/.vault-token)
   could cause confusing error messages during `vault login` [[GH-7508](https://github.com/hashicorp/vault/pull/7508)]
 * cli: Fix a bug where the `namespace list` command with JSON formatting
   always returned an empty object [[GH-7705](https://github.com/hashicorp/vault/pull/7705)]
 * cli: Command timeouts are now always specified solely by the
   `VAULT_CLIENT_TIMEOUT` value. [[GH-7469](https://github.com/hashicorp/vault/pull/7469)]
 * core: Don't allow registering a non-root zero TTL token lease. This is purely
   defense in depth as the lease would be revoked immediately anyways, but
   there's no real reason to allow registration. [[GH-7524](https://github.com/hashicorp/vault/pull/7524)]
 * core: Correctly revoke the token that's present in the response auth from a
   auth/token/ request if there's partial failure during the process. [[GH-7835](https://github.com/hashicorp/vault/pull/7835)]
 * identity (enterprise): Fixed identity case sensitive loading in secondary
   cluster [[GH-7327](https://github.com/hashicorp/vault/pull/7327)]
 * identity: Ensure only replication primary stores the identity case sensitivity state [[GH-7820](https://github.com/hashicorp/vault/pull/7820)]
 * raft: Fixed VAULT_CLUSTER_ADDR env being ignored at startup [[GH-7619](https://github.com/hashicorp/vault/pull/7619)]
 * secrets/pki: Don't allow duplicate SAN names in issued certs [[GH-7605](https://github.com/hashicorp/vault/pull/7605)]
 * sys/health: Pay attention to the values provided for `standbyok` and
   `perfstandbyok` rather than simply using their presence as a key to flip on
   that behavior [[GH-7323](https://github.com/hashicorp/vault/pull/7323)]
 * ui: using the `wrapped_token` query param will work with `redirect_to` and
   will automatically log in as intended [[GH-7398](https://github.com/hashicorp/vault/pull/7398)]
 * ui: fix an error when initializing from the UI using PGP keys [[GH-7542](https://github.com/hashicorp/vault/pull/7542)]
 * ui: show all active kv v2 secret versions even when `delete_version_after` is configured [[GH-7685](https://github.com/hashicorp/vault/pull/7685)]
 * ui: Ensure that items in the top navigation link to pages that users have access to [[GH-7590](https://github.com/hashicorp/vault/pull/7590)]

## 1.2.7
### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust

## 1.2.6.1
### August 21st, 2020
### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.2.6
### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.2.6 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.2.6)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.

## 1.2.5
### August 20th, 2020

SECURITY:

 * When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
 * When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.2.5 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.2.5)
  
BUG FIXES:
* seal: (enterprise) Fix issue causing stored seal and recovery keys to be mistaken as sealwrapped values

## 1.2.4 (November 7th, 2019)

SECURITY:

 * In a non-root namespace, revocation of a token scoped to a non-root
   namespace did not trigger the expected revocation of dynamic secret leases
   associated with that token. As a result, dynamic secret leases in non-root
   namespaces may outlive the token that created them.  This vulnerability,
   CVE-2019-18616, affects Vault Enterprise 0.11.0 and newer.
 * Disaster Recovery secondary clusters did not delete already-replicated data
   after a mount filter has been created on an upstream Performance secondary
   cluster. As a result, encrypted secrets may remain replicated on a Disaster
   Recovery secondary cluster after application of a mount filter excluding
   those secrets from replication. This vulnerability, CVE-2019-18617, affects
   Vault Enterprise 0.8 and newer.
 * Update version of Go to 1.12.12 to fix Go bug golang.org/issue/34960 which
   corresponds to CVE-2019-17596.

CHANGES:

 * auth/aws: If a custom `sts_endpoint` is configured, Vault Agent and the CLI
   should provide the corresponding region via the `region` parameter (which
   already existed as a CLI parameter, and has now been added to Agent). The
   automatic region detection added to the CLI and Agent in 1.2 has been removed.

IMPROVEMENTS:

  * cli: Ignore existing token during CLI login [[GH-7508](https://github.com/hashicorp/vault/pull/7508)]
  * core: Log proxy settings from environment on startup [[GH-7528](https://github.com/hashicorp/vault/pull/7528)]
  * core: Cache whether we've been initialized to reduce load on storage [[GH-7549](https://github.com/hashicorp/vault/pull/7549)]

BUG FIXES:

 * agent: Fix handling of gzipped responses [[GH-7470](https://github.com/hashicorp/vault/pull/7470)]
 * cli: Fix panic when pgp keys list is empty [[GH-7546](https://github.com/hashicorp/vault/pull/7546)]
 * cli: Command timeouts are now always specified solely by the
   `VAULT_CLIENT_TIMEOUT` value. [[GH-7469](https://github.com/hashicorp/vault/pull/7469)]
 * core: add hook for initializing seals for migration [[GH-7666](https://github.com/hashicorp/vault/pull/7666)]
 * core (enterprise): Migrating from one auto unseal method to another never
   worked on enterprise, now it does.
 * identity: Add required field `response_types_supported` to identity token
   `.well-known/openid-configuration` response [[GH-7533](https://github.com/hashicorp/vault/pull/7533)]
 * identity: Fixed nil pointer panic when merging entities [[GH-7712](https://github.com/hashicorp/vault/pull/7712)]
 * replication (Enterprise): Fix issue causing performance standbys nodes
   disconnecting when under high loads.
 * secrets/azure: Fix panic that could occur if client retries timeout [[GH-7793](https://github.com/hashicorp/vault/pull/7793)]
 * secrets/database: Fix bug in combined DB secrets engine that can result in
   writes to static-roles endpoints timing out [[GH-7518](https://github.com/hashicorp/vault/pull/7518)]
 * secrets/pki: Improve tidy to continue when value is nil [[GH-7589](https://github.com/hashicorp/vault/pull/7589)]
 * ui (Enterprise): Allow kv v2 secrets that are gated by Control Groups to be
   viewed in the UI [[GH-7504](https://github.com/hashicorp/vault/pull/7504)]

## 1.2.3 (September 12, 2019)

FEATURES:

* **Oracle Cloud (OCI) Integration**: Vault now support using Oracle Cloud for
  storage, auto unseal, and authentication.

IMPROVEMENTS:

 * auth/jwt: Groups claim matching now treats a string response as a single
   element list [[GH-63](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/63)]
 * auth/kubernetes: enable better support for projected tokens API by allowing
   user to specify issuer [[GH-65](https://github.com/hashicorp/vault/pull/65)]
 * auth/pcf: The PCF auth plugin was renamed to the CF auth plugin, maintaining
   full backwards compatibility [[GH-7346](https://github.com/hashicorp/vault/pull/7346)]
 * replication: Premium packages now come with unlimited performance standby
   nodes

BUG FIXES:

 * agent: Allow batch tokens and other non-renewable tokens to be used for
   agent operations [[GH-7441](https://github.com/hashicorp/vault/pull/7441)]
 * auth/jwt: Fix an error where newer (v1.2) token_* configuration parameters
   were not being applied to tokens generated using the OIDC login flow
   [[GH-67](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/67)]
 * raft: Fix an incorrect JSON tag on `leader_ca_cert` in the join request [[GH-7393](https://github.com/hashicorp/vault/pull/7393)]
 * seal/transit: Allow using Vault Agent for transit seal operations [[GH-7441](https://github.com/hashicorp/vault/pull/7441)]
 * storage/couchdb: Fix a file descriptor leak [[GH-7345](https://github.com/hashicorp/vault/pull/7345)]
 * ui: Fix a bug where the status menu would disappear when trying to revoke a
   token [[GH-7337](https://github.com/hashicorp/vault/pull/7337)]
 * ui: Fix a regression that prevented input of custom items in search-select
   [[GH-7338](https://github.com/hashicorp/vault/pull/7338)]
 * ui: Fix an issue with the namespace picker being unable to render nested
   namespaces named with numbers and sorting of namespaces in the picker
   [[GH-7333](https://github.com/hashicorp/vault/pull/7333)]

## 1.2.2 (August 15, 2019)

CHANGES:

 * auth/pcf: The signature format has been updated to use the standard Base64
   encoding instead of the URL-safe variant. Signatures created using the
   previous format will continue to be accepted [PCF-27]
 * core: The http response code returned when an identity token key is not found
   has been changed from 400 to 404

IMPROVEMENTS:

 * identity: Remove 512 entity limit for groups [[GH-7317](https://github.com/hashicorp/vault/pull/7317)]

BUG FIXES:

 * auth/approle: Fix an error where an empty `token_type` string was not being
   correctly handled as `TokenTypeDefault` [[GH-7273](https://github.com/hashicorp/vault/pull/7273)]
 * auth/radius: Fix panic when logging in [[GH-7286](https://github.com/hashicorp/vault/pull/7286)]
 * ui: the string-list widget will now honor multiline input [[GH-7254](https://github.com/hashicorp/vault/pull/7254)]
 * ui: various visual bugs in the KV interface were addressed [[GH-7307](https://github.com/hashicorp/vault/pull/7307)]
 * ui: fixed incorrect URL to access help in LDAP auth [[GH-7299](https://github.com/hashicorp/vault/pull/7299)]

## 1.2.1 (August 6th, 2019)

BUG FIXES:

 * agent: Fix a panic on creds pulling in some error conditions in `aws` and
   `alicloud` auth methods [[GH-7238](https://github.com/hashicorp/vault/pull/7238)]
 * auth/approle: Fix error reading role-id on a role created pre-1.2 [[GH-7231](https://github.com/hashicorp/vault/pull/7231)]
 * auth/token: Fix sudo check in non-root namespaces on create [[GH-7224](https://github.com/hashicorp/vault/pull/7224)]
 * core: Fix health checks with perfstandbyok=true returning the wrong status
   code [[GH-7240](https://github.com/hashicorp/vault/pull/7240)]
 * ui: The web CLI will now parse input as a shell string, with special
   characters escaped [[GH-7206](https://github.com/hashicorp/vault/pull/7206)]
 * ui: The UI will now redirect to a page after authentication [[GH-7088](https://github.com/hashicorp/vault/pull/7088)]
 * ui (Enterprise): The list of namespaces is now cleared when logging
   out [[GH-7186](https://github.com/hashicorp/vault/pull/7186)]

## 1.2.0 (July 30th, 2019)

CHANGES:

 * Token store roles use new, common token fields for the values
   that overlap with other auth backends. `period`, `explicit_max_ttl`, and
   `bound_cidrs` will continue to work, with priority being given to the
   `token_` prefixed versions of those parameters. They will also be returned
   when doing a read on the role if they were used to provide values initially;
   however, in Vault 1.4 if `period` or `explicit_max_ttl` is zero they will no
   longer be returned. (`explicit_max_ttl` was already not returned if empty.)
 * Due to underlying changes in Go version 1.12 and Go > 1.11.5, Vault is now
   stricter about what characters it will accept in path names. Whereas before
   it would filter out unprintable characters (and this could be turned off),
   control characters and other invalid characters are now rejected within Go's
   HTTP library before the request is passed to Vault, and this cannot be
   disabled. To continue using these (e.g. for already-written paths), they
   must be properly percent-encoded (e.g. `\r` becomes `%0D`, `\x00` becomes
   `%00`, and so on).
 * The user-configured regions on the AWSKMS seal stanza will now be preferred
   over regions set in the enclosing environment.  This is a _breaking_ change.
 * All values in audit logs now are omitted if they are empty.  This helps
   reduce the size of audit log entries by not reproducing keys in each entry
   that commonly don't contain any value, which can help in cases where audit
   log entries are above the maximum UDP packet size and others.
 * Both PeriodicFunc and WALRollback functions will be called if both are
   provided. Previously WALRollback would only be called if PeriodicFunc was
   not set. See [[GH-6717](https://github.com/hashicorp/vault/pull/6717)] for
   details.
 * Vault now uses Go's official dependency management system, Go Modules, to
   manage dependencies. As a result to both reduce transitive dependencies for
   API library users and plugin authors, and to work around various conflicts,
   we have moved various helpers around, mostly under an `sdk/` submodule. A
   couple of functions have also moved from plugin helper code to the `api/`
   submodule. If you are a plugin author, take a look at some of our official
   plugins and the paths they are importing for guidance.
 * AppRole uses new, common token fields for values that overlap
   with other auth backends. `period` and `policies` will continue to work,
   with priority being given to the `token_` prefixed versions of those
   parameters. They will also be returned when doing a read on the role if they
   were used to provide values initially.
 * In AppRole, `"default"` is no longer automatically added to the `policies`
   parameter. This was a no-op since it would always be added anyways by
   Vault's core; however, this can now be explicitly disabled with the new
   `token_no_default_policy` field.
 * In AppRole, `bound_cidr_list` is no longer returned when reading a role
 * rollback: Rollback will no longer display log messages when it runs; it will
   only display messages on error.
 * Database plugins will now default to 4 `max_open_connections`
   rather than 2.

FEATURES:

 * **Integrated Storage**: Vault 1.2 includes a _tech preview_ of a new way to
   manage storage directly within a Vault cluster. This new integrated storage
   solution is based on the Raft protocol which is also used to back HashiCorp
   Consul and HashiCorp Nomad.
 * **Combined DB credential rotation**: Alternative mode for the Combined DB
   Secret Engine to automatically rotate existing database account credentials
   and set Vault as the source of truth for credentials.
 * **Identity Tokens**: Vault's Identity system can now generate OIDC-compliant
   ID tokens. These customizable tokens allow encapsulating a signed, verifiable
   snapshot of identity information and metadata. They can be use by other
   applications—even those without Vault authorization—as a way of establishing
   identity based on a Vault entity.
 * **Pivotal Cloud Foundry plugin**: New auth method using Pivotal Cloud
   Foundry certificates for Vault authentication.
 * **ElasticSearch database plugin**: New ElasticSearch database plugin issues
   unique, short-lived ElasticSearch credentials.
 * **New UI Features**: An HTTP Request Volume Page and new UI for editing LDAP
   Users and Groups have been added.
 * **HA support for Postgres**: PostgreSQL versions >= 9.5 may now but used as
   and HA storage backend.
 * **KMIP secrets engine (Enterprise)**: Allows Vault to operate as a KMIP
   Server, seamlessly brokering cryptographic operations for traditional
   infrastructure.
 * Common Token Fields: Auth methods now use common fields for controlling
   token behavior, making it easier to understand configuration across methods.
 * **Vault API explorer**: The Vault UI now includes an embedded API explorer
   where you can browse the endpoints avaliable to you and make requests. To try
   it out, open the Web CLI and type `api`.

IMPROVEMENTS:

 * agent: Allow EC2 nonce to be passed in [[GH-6953](https://github.com/hashicorp/vault/pull/6953)]
 * agent: Add optional `namespace` parameter, which sets the default namespace
   for the auto-auth functionality [[GH-6988](https://github.com/hashicorp/vault/pull/6988)]
 * agent: Add cert auto-auth method [[GH-6652](https://github.com/hashicorp/vault/pull/6652)]
 * api: Add support for passing data to delete operations via `DeleteWithData`
   [[GH-7139](https://github.com/hashicorp/vault/pull/7139)]
 * audit/file: Dramatically speed up file operations by changing
   locking/marshaling order [[GH-7024](https://github.com/hashicorp/vault/pull/7024)]
 * auth/jwt: A JWKS endpoint may now be configured for signature verification [[GH-43](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/43)]
 * auth/jwt: A new `verbose_oidc_logging` role parameter has been added to help
   troubleshoot OIDC configuration [[GH-57](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/57)]
 * auth/jwt: `bound_claims` will now match received claims that are lists if any element
   of the list is one of the expected values [[GH-50](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/50)]
 * auth/jwt: Leeways for `nbf` and `exp` are now configurable, as is clock skew
   leeway [[GH-53](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/53)]
 * auth/kubernetes: Allow service names/namespaces to be configured as globs
   [[GH-58](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/58)]
 * auth/token: Allow the support of the identity system for the token backend
   via token roles [[GH-6267](https://github.com/hashicorp/vault/pull/6267)]
 * auth/token: Add a large set of token configuration options to token store
   roles [[GH-6662](https://github.com/hashicorp/vault/pull/6662)]
 * cli: `path-help` now allows `-format=json` to be specified, which will
   output OpenAPI [[GH-7006](https://github.com/hashicorp/vault/pull/7006)]
 * cli: Add support for passing parameters to `vault delete` operations
   [[GH-7139](https://github.com/hashicorp/vault/pull/7139)]
 * cli: Add a log-format CLI flag that can specify either "standard" or "json"
   for the log format for the `vault server`command. [[GH-6840](https://github.com/hashicorp/vault/pull/6840)]
 * cli: Add `-dev-no-store-token` to allow dev servers to not store the
   generated token at the tokenhelper location [[GH-7104](https://github.com/hashicorp/vault/pull/7104)]
 * identity: Allow a group alias' canonical ID to be modified
 * namespaces: Namespaces can now be created and deleted from performance
   replication secondaries
 * plugins: Change the default for `max_open_connections` for DB plugins to 4
   [[GH-7093](https://github.com/hashicorp/vault/pull/7093)]
 * replication: Client TLS authentication is now supported when enabling or
   updating a replication secondary
 * secrets/database: Cassandra operations will now cancel on client timeout
   [[GH-6954](https://github.com/hashicorp/vault/pull/6954)]
 * secrets/kv: Add optional `delete_version_after` parameter, which takes a
   duration and can be set on the mount and/or the metadata for a specific key
   [[GH-7005](https://github.com/hashicorp/vault/pull/7005)]
 * storage/postgres: LIST now performs better on large datasets [[GH-6546](https://github.com/hashicorp/vault/pull/6546)]
 * storage/s3: A new `path` parameter allows selecting the path within a bucket
   for Vault data [[GH-7157](https://github.com/hashicorp/vault/pull/7157)]
 * ui: KV v1 and v2 will now gracefully degrade allowing a write without read
   workflow in the UI [[GH-6570](https://github.com/hashicorp/vault/pull/6570)]
 * ui: Many visual improvements with the addition of Toolbars [[GH-6626](https://github.com/hashicorp/vault/pull/6626)], the restyling
   of the Confirm Action component [[GH-6741](https://github.com/hashicorp/vault/pull/6741)], and using a new set of glyphs for our
   Icon component [[GH-6736](https://github.com/hashicorp/vault/pull/6736)]
 * ui: Lazy loading parts of the application so that the total initial payload is
   smaller [[GH-6718](https://github.com/hashicorp/vault/pull/6718)]
 * ui: Tabbing to auto-complete in filters will first complete a common prefix if there
   is one [[GH-6759](https://github.com/hashicorp/vault/pull/6759)]
 * ui: Removing jQuery from the application makes the initial JS payload smaller [[GH-6768](https://github.com/hashicorp/vault/pull/6768)]

BUG FIXES:

 * audit: Log requests and responses due to invalid wrapping token provided
   [[GH-6541](https://github.com/hashicorp/vault/pull/6541)]
 * audit: Fix bug preventing request counter queries from working with auditing
   enabled [[GH-6767](https://github.com/hashicorp/vault/pull/6767)
 * auth/aws: AWS Roles are now upgraded and saved to the latest version just
   after the AWS credential plugin is mounted. [[GH-7025](https://github.com/hashicorp/vault/pull/7025)]
 * auth/aws: Fix a case where a panic could stem from a malformed assumed-role ARN
   when parsing this value [[GH-6917](https://github.com/hashicorp/vault/pull/6917)]
 * auth/aws: Fix an error complaining about a read-only view that could occur
   during updating of a role when on a performance replication secondary
   [[GH-6926](https://github.com/hashicorp/vault/pull/6926)]
 * auth/jwt: Fix a regression introduced in 1.1.1 that disabled checking of client_id
   for OIDC logins [[GH-54](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/54)]
 * auth/jwt: Fix a panic during OIDC CLI logins that could occur if the Vault server
   response is empty [[GH-55](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/55)]
 * auth/jwt: Fix issue where OIDC logins might intermittently fail when using
   performance standbys [[GH-61](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/61)]
 * identity: Fix a case where modifying aliases of an entity could end up
   moving the entity into the wrong namespace
 * namespaces: Fix a behavior (currently only known to be benign) where we
   wouldn't delete policies through the official functions before wiping the
   namespaces on deletion
 * secrets/database: Escape username/password before using in connection URL
   [[GH-7089](https://github.com/hashicorp/vault/pull/7089)]
 * secrets/pki: Forward revocation requests to active node when on a
   performance standby [[GH-7173](https://github.com/hashicorp/vault/pull/7173)]
 * ui: Fix timestamp on some transit keys [[GH-6827](https://github.com/hashicorp/vault/pull/6827)]
 * ui: Show Entities and Groups in Side Navigation [[GH-7138](https://github.com/hashicorp/vault/pull/7138)]
 * ui: Ensure dropdown updates selected item on HTTP Request Metrics page

## 1.1.4/1.1.5 (July 25th/30th, 2019)

NOTE:

Although 1.1.4 was tagged, we realized very soon after the tag was publicly
pushed that an intended fix was accidentally left out. As a result, 1.1.4 was
not officially announced and 1.1.5 should be used as the release after 1.1.3.

IMPROVEMENTS:

 * identity: Allow a group alias' canonical ID to be modified
 * namespaces: Improve namespace deletion performance [[GH-6939](https://github.com/hashicorp/vault/pull/6939)]
 * namespaces: Namespaces can now be created and deleted from performance
   replication secondaries

BUG FIXES:

 * api: Add backwards compat support for API env vars [[GH-7135](https://github.com/hashicorp/vault/pull/7135)]
 * auth/aws: Fix a case where a panic could stem from a malformed assumed-role
   ARN when parsing this value [[GH-6917](https://github.com/hashicorp/vault/pull/6917)]
 * auth/ldap: Add `use_pre111_group_cn_behavior` flag to allow recovering from
   a regression caused by a bug fix starting in 1.1.1 [[GH-7208](https://github.com/hashicorp/vault/pull/7208)]
 * auth/aws: Use a role cache to avoid separate locking paths [[GH-6926](https://github.com/hashicorp/vault/pull/6926)]
 * core: Fix a deadlock if a panic happens during request handling [[GH-6920](https://github.com/hashicorp/vault/pull/6920)]
 * core: Fix an issue that may cause key upgrades to not be cleaned up properly
   [[GH-6949](https://github.com/hashicorp/vault/pull/6949)]
 * core: Don't shutdown if key upgrades fail due to canceled context [[GH-7070](https://github.com/hashicorp/vault/pull/7070)]
 * core: Fix panic caused by handling requests while vault is inactive
 * identity: Fix reading entity and groups that have spaces in their names
   [[GH-7055](https://github.com/hashicorp/vault/pull/7055)]
 * identity: Ensure entity alias operations properly verify namespace [[GH-6886](https://github.com/hashicorp/vault/pull/6886)]
 * mfa: Fix a nil pointer panic that could occur if invalid Duo credentials
   were supplied
 * replication: Forward step-down on perf standbys to match HA behavior
 * replication: Fix various read only storage errors on performance standbys
 * replication: Stop forwarding before stopping replication to eliminate some
   possible bad states
 * secrets/database: Allow cassandra queries to be cancled [[GH-6954](https://github.com/hashicorp/vault/pull/6954)]
 * storage/consul: Fix a regression causing vault to not connect to consul over
   unix sockets [[GH-6859](https://github.com/hashicorp/vault/pull/6859)]
 * ui: Fix saving of TTL and string array fields generated by Open API [[GH-7094](https://github.com/hashicorp/vault/pull/7094)]

## 1.1.3 (June 5th, 2019)

IMPROVEMENTS:

 * agent: Now supports proxying request query parameters [[GH-6772](https://github.com/hashicorp/vault/pull/6772)]
 * core: Mount table output now includes a UUID indicating the storage path [[GH-6633](https://github.com/hashicorp/vault/pull/6633)]
 * core: HTTP server timeout values are now configurable [[GH-6666](https://github.com/hashicorp/vault/pull/6666)]
 * replication: Improve performance of the reindex operation on secondary clusters
   when mount filters are in use
 * replication: Replication status API now returns the state and progress of a reindex

BUG FIXES:

 * api: Return the Entity ID in the secret output [[GH-6819](https://github.com/hashicorp/vault/pull/6819)]
 * auth/jwt: Consider bound claims when considering if there is at least one
   bound constraint [[GH-49](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/49)]
 * auth/okta: Fix handling of group names containing slashes [[GH-6665](https://github.com/hashicorp/vault/pull/6665)]
 * cli: Add deprecated stored-shares flag back to the init command [[GH-6677](https://github.com/hashicorp/vault/pull/6677)]
 * cli: Fix a panic when the KV command would return no data [[GH-6675](https://github.com/hashicorp/vault/pull/6675)]
 * cli: Fix issue causing CLI list operations to not return proper format when
   there is an empty response [[GH-6776](https://github.com/hashicorp/vault/pull/6776)]
 * core: Correctly honor non-HMAC request keys when auditing requests [[GH-6653](https://github.com/hashicorp/vault/pull/6653)]
 * core: Fix the `x-vault-unauthenticated` value in OpenAPI for a number of
   endpoints [[GH-6654](https://github.com/hashicorp/vault/pull/6654)]
 * core: Fix issue where some OpenAPI parameters were incorrectly listed as
   being sent as a header [[GH-6679](https://github.com/hashicorp/vault/pull/6679)]
 * core: Fix issue that would allow duplicate mount names to be used [[GH-6771](https://github.com/hashicorp/vault/pull/6771)]
 * namespaces: Fix behavior when using `root` instead of `root/` as the
   namespace header value
 * pki: fix a panic when a client submits a null value [[GH-5679](https://github.com/hashicorp/vault/pull/5679)]
 * replication: Properly update mount entry cache on a secondary to apply all
   new values after a tune
 * replication: Properly close connection on bootstrap error
 * replication: Fix an issue causing startup problems if a namespace policy
   wasn't replicated properly
 * replication: Fix longer than necessary WAL replay during an initial reindex
 * replication: Fix error during mount filter invalidation on DR secondary clusters
 * secrets/ad: Make time buffer configurable [AD-35]
 * secrets/gcp: Check for nil config when getting credentials [[GH-35](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/35)]
 * secrets/gcp: Fix error checking in some cases where the returned value could
   be 403 instead of 404 [[GH-37](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/37)]
 * secrets/gcpkms: Disable key rotation when deleting a key [[GH-10](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/10)]
 * storage/consul: recognize `https://` address even if schema not specified
   [[GH-6602](https://github.com/hashicorp/vault/pull/6602)]
 * storage/dynamodb: Fix an issue where a deleted lock key in DynamoDB (HA)
   could cause constant switching of the active node [[GH-6637](https://github.com/hashicorp/vault/pull/6637)]
 * storage/dynamodb: Eliminate a high-CPU condition that could occur if an
   error was received from the DynamoDB API [[GH-6640](https://github.com/hashicorp/vault/pull/6640)]
 * storage/gcs: Correctly use configured chunk size values [[GH-6655](https://github.com/hashicorp/vault/pull/6655)]
 * storage/mssql: Use the correct database when pre-created schemas exist
   [[GH-6356](https://github.com/hashicorp/vault/pull/6356)]
 * ui: Fix issue with select arrows on drop down menus [[GH-6627](https://github.com/hashicorp/vault/pull/6627)]
 * ui: Fix an issue where sensitive input values weren't being saved to the
   server [[GH-6586](https://github.com/hashicorp/vault/pull/6586)]
 * ui: Fix web cli parsing when using quoted values [[GH-6755](https://github.com/hashicorp/vault/pull/6755)]
 * ui: Fix a namespace workflow mapping identities from external namespaces by
   allowing arbitrary input in search-select component [[GH-6728](https://github.com/hashicorp/vault/pull/6728)]

## 1.1.2 (April 18th, 2019)

This is a bug fix release containing the two items below. It is otherwise
unchanged from 1.1.1.

BUG FIXES:

 * auth/okta: Fix a potential dropped error [[GH-6592](https://github.com/hashicorp/vault/pull/6592)]
 * secrets/kv: Fix a regression on upgrade where a KVv2 mount could fail to be
   mounted on unseal if it had previously been mounted but not written to
   [[GH-31](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/31)]

## 1.1.1 (April 11th, 2019)

SECURITY:

 * Given: (a) performance replication is enabled; (b) performance standbys are
   in use on the performance replication secondary cluster; and (c) mount
   filters are in use, if a mount that was previously available to a secondary
   is updated to be filtered out, although the data would be removed from the
   secondary cluster, the in-memory cache of the data would not be purged on
   the performance standby nodes. As a result, the previously-available data
   could still be read from memory if it was ever read from disk, and if this
   included mount configuration data this could result in token or lease
   issuance. The issue is fixed in this release; in prior releases either an
   active node changeover (such as a step-down) or a restart of the standby
   nodes is sufficient to cause the performance standby nodes to clear their
   cache. A CVE is in the process of being issued; the number is
   CVE-2019-11075.
 * Roles in the JWT Auth backend using the OIDC login flow (i.e. role_type of
   “oidc”) were not enforcing bound_cidrs restrictions, if any were configured
   for the role. This issue did not affect roles of type “jwt”.

CHANGES:

 * auth/jwt: Disallow logins of role_type "oidc" via the `/login` path [[GH-38](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/38)]
 * core/acl:  New ordering defines which policy wins when there are multiple
   inexact matches and at least one path contains `+`. `+*` is now illegal in
   policy paths. The previous behavior simply selected any matching
   segment-wildcard path that matched. [[GH-6532](https://github.com/hashicorp/vault/pull/6532)]
 * replication: Due to technical limitations, mounting and unmounting was not
   previously possible from a performance secondary. These have been resolved,
   and these operations may now be run from a performance secondary.

IMPROVEMENTS:

 * agent: Allow AppRole auto-auth without a secret-id [[GH-6324](https://github.com/hashicorp/vault/pull/6324)]
 * auth/gcp: Cache clients to improve performance and reduce open file usage
 * auth/jwt: Bounds claims validiation will now allow matching the received
   claims against a list of expected values [[GH-41](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/41)]
 * secret/gcp: Cache clients to improve performance and reduce open file usage
 * replication: Mounting/unmounting/remounting/mount-tuning is now supported
   from a performance secondary cluster
 * ui: Suport for authentication via the RADIUS auth method [[GH-6488](https://github.com/hashicorp/vault/pull/6488)]
 * ui: Navigating away from secret list view will clear any page-specific
   filter that was applied [[GH-6511](https://github.com/hashicorp/vault/pull/6511)]
 * ui: Improved the display when OIDC auth errors [[GH-6553](https://github.com/hashicorp/vault/pull/6553)]

BUG FIXES:

 * agent: Allow auto-auth to be used with caching without having to define any
   sinks [[GH-6468](https://github.com/hashicorp/vault/pull/6468)]
 * agent: Disallow some nonsensical config file combinations [[GH-6471](https://github.com/hashicorp/vault/pull/6471)]
 * auth/ldap: Fix CN check not working if CN was not all in uppercase [[GH-6518](https://github.com/hashicorp/vault/pull/6518)]
 * auth/jwt: The CLI helper for OIDC logins will now open the browser to the correct
   URL when running on Windows [[GH-37](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/37)]
 * auth/jwt: Fix OIDC login issue where configured TLS certs weren't being used [[GH-40](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/40)]
 * auth/jwt: Fix an issue where the `oidc_scopes` parameter was not being included in
   the response to a role read request [[GH-35](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/35)]
 * core: Fix seal migration case when migrating to Shamir and a seal block
   wasn't explicitly specified [[GH-6455](https://github.com/hashicorp/vault/pull/6455)]
 * core: Fix unwrapping when using namespaced wrapping tokens [[GH-6536](https://github.com/hashicorp/vault/pull/6536)]
 * core: Fix incorrect representation of required properties in OpenAPI output
   [[GH-6490](https://github.com/hashicorp/vault/pull/6490)]
 * core: Fix deadlock that could happen when using the UI [[GH-6560](https://github.com/hashicorp/vault/pull/6560)]
 * identity: Fix updating groups removing existing members [[GH-6527](https://github.com/hashicorp/vault/pull/6527)]
 * identity: Properly invalidate group alias in performance secondary [[GH-6564](https://github.com/hashicorp/vault/pull/6564)]
 * identity: Use namespace context when loading entities and groups to ensure
   merging of duplicate entries works properly [[GH-6563](https://github.com/hashicorp/vault/pull/6563)]
 * replication: Fix performance standby election failure [[GH-6561](https://github.com/hashicorp/vault/pull/6561)]
 * replication: Fix mount filter invalidation on performance standby nodes
 * replication: Fix license reloading on performance standby nodes
 * replication: Fix handling of control groups on performance standby nodes
 * replication: Fix some forwarding scenarios with request bodies using
   performance standby nodes [[GH-6538](https://github.com/hashicorp/vault/pull/6538)]
 * secret/gcp: Fix roleset binding when using JSON [[GH-27](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/27)]
 * secret/pki: Use `uri_sans` param in when not using CSR parameters [[GH-6505](https://github.com/hashicorp/vault/pull/6505)]
 * storage/dynamodb: Fix a race condition possible in HA configurations that could
   leave the cluster without a leader [[GH-6512](https://github.com/hashicorp/vault/pull/6512)]
 * ui: Fix an issue where in production builds OpenAPI model generation was
   failing, causing any form using it to render labels with missing fields [[GH-6474](https://github.com/hashicorp/vault/pull/6474)]
 * ui: Fix issue nav-hiding when moving between namespaces [[GH-6473](https://github.com/hashicorp/vault/pull/6473)]
 * ui: Secrets will always show in the nav regardless of access to cubbyhole [[GH-6477](https://github.com/hashicorp/vault/pull/6477)]
 * ui: fix SSH OTP generation [[GH-6540](https://github.com/hashicorp/vault/pull/6540)]
 * ui: add polyfill to load UI in IE11 [[GH-6567](https://github.com/hashicorp/vault/pull/6567)]
 * ui: Fix issue where some elements would fail to work properly if using ACLs
   with segment-wildcard paths (`/+/` segments) [[GH-6525](https://github.com/hashicorp/vault/pull/6525)]

## 1.1.0 (March 18th, 2019)

CHANGES:

 * auth/jwt: The `groups_claim_delimiter_pattern` field has been removed. If the
   groups claim is not at the top level, it can now be specified as a
   [JSONPointer](https://tools.ietf.org/html/rfc6901).
 * auth/jwt: Roles now have a "role type" parameter with a default type of
   "oidc". To configure new JWT roles, a role type of "jwt" must be explicitly
   specified.
 * cli: CLI commands deprecated in 0.9.2 are now removed. Please see the CLI
   help/warning output in previous versions of Vault for updated commands.
 * core: Vault no longer automatically mounts a K/V backend at the "secret/"
   path when initializing Vault
 * core: Vault's cluster port will now be open at all times on HA standby nodes
 * plugins: Vault no longer supports running netRPC plugins. These were
   deprecated in favor of gRPC based plugins and any plugin built since 0.9.4
   defaults to gRPC. Older plugins may need to be recompiled against the latest
   Vault dependencies.

FEATURES:

 * **Vault Agent Caching**: Vault Agent can now be configured to act as a
   caching proxy to Vault. Clients can send requests to Vault Agent and the
   request will be proxied to the Vault server and cached locally in Agent.
   Currently Agent will cache generated leases and tokens and keep them
   renewed. The proxy can also use the Auto Auth feature so clients do not need
   to authenticate to Vault, but rather can make requests to Agent and have
   Agent fully manage token lifecycle.
 * **OIDC Redirect Flow Support**: The JWT auth backend now supports OIDC
   roles. These allow authentication via an OIDC-compliant provider via the
   user's browser. The login may be initiated from the Vault UI or through
   the `vault login` command.
 * **ACL Path Wildcard**: ACL paths can now use the `+` character to enable
   wild card matching for a single directory in the path definition.
 * **Transit Auto Unseal**: Vault can now be configured to use the Transit
   Secret Engine in another Vault cluster as an auto unseal provider.

IMPROVEMENTS:

 * auth/jwt: A default role can be set. It will be used during JWT/OIDC logins if
   a role is not specified.
 * auth/jwt: Arbitrary claims data can now be copied into token & alias metadata.
 * auth/jwt: An arbitrary set of bound claims can now be configured for a role.
 * auth/jwt: The name "oidc" has been added as an alias for the jwt backend. Either
   name may be specified in the `auth enable` command.
 * command/server: A warning will be printed when 'tls_cipher_suites' includes a
   blacklisted cipher suite or all cipher suites are blacklisted by the HTTP/2
   specification [[GH-6300](https://github.com/hashicorp/vault/pull/6300)]
 * core/metrics: Prometheus pull support using a new sys/metrics endpoint. [[GH-5308](https://github.com/hashicorp/vault/pull/5308)]
 * core: On non-windows platforms a SIGUSR2 will make the server log a dump of
   all running goroutines' stack traces for debugging purposes [[GH-6240](https://github.com/hashicorp/vault/pull/6240)]
 * replication: The initial replication indexing process on newly initialized or upgraded
   clusters now runs asynchronously
 * sentinel: Add token namespace id and path, available in rules as
   token.namespace.id and token.namespace.path
 * ui: The UI is now leveraging OpenAPI definitions to pull in fields for various forms.
   This means, it will not be necessary to add fields on the go and JS sides in the future.
   [[GH-6209](https://github.com/hashicorp/vault/pull/6209)]

BUG FIXES:

 * auth/jwt: Apply `bound_claims` validation across all login paths
 * auth/jwt: Update `bound_audiences` validation during non-OIDC logins to accept
   any matched audience, as documented and handled in OIDC logins [[GH-30](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/30)]
 * auth/token: Fix issue where empty values for token role update call were
   ignored [[GH-6314](https://github.com/hashicorp/vault/pull/6314)]
 * core: The `operator migrate` command will no longer hang on empty key names
   [[GH-6371](https://github.com/hashicorp/vault/pull/6371)]
 * identity: Fix a panic at login when external group has a nil alias [[GH-6230](https://github.com/hashicorp/vault/pull/6230)]
 * namespaces: Clear out identity store items upon namespace deletion
 * replication/perfstandby: Fixed a bug causing performance standbys to wait
   longer than necessary after forwarding a write to the active node
 * replication/mountfilter: Fix a deadlock that could occur when mount filters
   were updated [[GH-6426](https://github.com/hashicorp/vault/pull/6426)]
 * secret/kv: Fix issue where a v1→v2 upgrade could run on a performance
   standby when using a local mount
 * secret/ssh: Fix for a bug where attempting to delete the last ssh role
   in the zeroaddress configuration could fail [[GH-6390](https://github.com/hashicorp/vault/pull/6390)]
 * secret/totp: Uppercase provided keys so they don't fail base32 validation
   [[GH-6400](https://github.com/hashicorp/vault/pull/6400)]
 * secret/transit: Multiple HMAC, Sign or Verify operations can now be
   performed with one API call using the new `batch_input` parameter [[GH-5875](https://github.com/hashicorp/vault/pull/5875)]
 * sys: `sys/internal/ui/mounts` will no longer return secret or auth mounts
   that have been filtered. Similarly, `sys/internal/ui/mount/:path` will
   return a error response if a filtered mount path is requested. [[GH-6412](https://github.com/hashicorp/vault/pull/6412)]
 * ui: Fix for a bug where you couldn't access the data tab after clicking on
   wrap details on the unwrap page [[GH-6404](https://github.com/hashicorp/vault/pull/6404)]
 * ui: Fix an issue where the policies tab was erroneously hidden [[GH-6301](https://github.com/hashicorp/vault/pull/6301)]
 * ui: Fix encoding issues with kv interfaces [[GH-6294](https://github.com/hashicorp/vault/pull/6294)]

## 1.0.3.1 (March 14th, 2019) (Enterprise Only)

SECURITY:

 * A regression was fixed in replication mount filter code introduced in Vault
   1.0 that caused the underlying filtered data to be replicated to
   secondaries. This data was not accessible to users via Vault's API but via a
   combination of privileged configuration file changes/Vault commands it could
   be read.  Upgrading to this version or 1.1 will fix this issue and cause the
   replicated data to be deleted from filtered secondaries. More information
   was sent to customer contacts on file.

## 1.0.3 (February 12th, 2019)

CHANGES:

 * New AWS authentication plugin mounts will default to using the generated
   role ID as the Identity alias name. This applies to both EC2 and IAM auth.
   Existing mounts that explicitly set this value will not be affected but
   mounts that specified no preference will switch over on upgrade.
 * The default policy now allows a token to look up its associated identity
   entity either by name or by id [[GH-6105](https://github.com/hashicorp/vault/pull/6105)]
 * The Vault UI's navigation and onboarding wizard now only displays items that
   are permitted in a users' policy [[GH-5980](https://github.com/hashicorp/vault/pull/5980), [GH-6094](https://github.com/hashicorp/vault/pull/6094)]
 * An issue was fixed that caused recovery keys to not work on secondary
   clusters when using a different unseal mechanism/key than the primary. This
   would be hit if the cluster was rekeyed or initialized after 1.0. We recommend
   rekeying the recovery keys on the primary cluster if you meet the above
   requirements.

FEATURES:

 * **cURL Command Output**: CLI commands can now use the `-output-curl-string`
   flag to print out an equivalent cURL command.
 * **Response Headers From Plugins**: Plugins can now send back headers that
   will be included in the response to a client. The set of allowed headers can
   be managed by the operator.

IMPROVEMENTS:

 * auth/aws: AWS EC2 authentication can optionally create entity aliases by
   role ID [[GH-6133](https://github.com/hashicorp/vault/pull/6133)]
 * auth/jwt: The supported set of signing algorithms is now configurable [JWT
   plugin [GH-16](https://github.com/hashicorp/vault/pull/16)]
 * core: When starting from an uninitialized state, HA nodes will now attempt
   to auto-unseal using a configured auto-unseal mechanism after the active
   node initializes Vault [[GH-6039](https://github.com/hashicorp/vault/pull/6039)]
 * secret/database: Add socket keepalive option for Cassandra [[GH-6201](https://github.com/hashicorp/vault/pull/6201)]
 * secret/ssh: Add signed key constraints, allowing enforcement of key types
   and minimum key sizes [[GH-6030](https://github.com/hashicorp/vault/pull/6030)]
 * secret/transit: ECDSA signatures can now be marshaled in JWS-compatible
   fashion [[GH-6077](https://github.com/hashicorp/vault/pull/6077)]
 * storage/etcd: Support SRV service names [[GH-6087](https://github.com/hashicorp/vault/pull/6087)]
 * storage/aws: Support specifying a KMS key ID for server-side encryption
   [[GH-5996](https://github.com/hashicorp/vault/pull/5996)]

BUG FIXES:

 * core: Fix a rare case where a standby whose connection is entirely torn down
   to the active node, then reconnects to the same active node, may not
   successfully resume operation [[GH-6167](https://github.com/hashicorp/vault/pull/6167)]
 * cors: Don't duplicate headers when they're written [[GH-6207](https://github.com/hashicorp/vault/pull/6207)]
 * identity: Persist merged entities only on the primary [[GH-6075](https://github.com/hashicorp/vault/pull/6075)]
 * replication: Fix a potential race when a token is created and then used with
   a performance standby very quickly, before an associated entity has been
   replicated. If the entity is not found in this scenario, the request will
   forward to the active node.
 * replication: Fix issue where recovery keys would not work on secondary
   clusters if using a different unseal mechanism than the primary.
 * replication: Fix a "failed to register lease" error when using performance
   standbys
 * storage/postgresql: The `Get` method will now return an Entry object with
   the `Key` member correctly populated with the full path that was requested
   instead of just the last path element [[GH-6044](https://github.com/hashicorp/vault/pull/6044)]

## 1.0.2 (January 15th, 2019)

SECURITY:

 * When creating a child token from a parent with `bound_cidrs`, the list of
   CIDRs would not be propagated to the child token, allowing the child token
   to be used from any address.

CHANGES:

 * secret/aws: Role now returns `credential_type` instead of `credential_types`
   to match role input. If a legacy role that can supply more than one
   credential type, they will be concatenated with a `,`.
 * physical/dynamodb, autoseal/aws: Instead of Vault performing environment
   variable handling, and overriding static (config file) values if found, we
   use the default AWS SDK env handling behavior, which also looks for
   deprecated values. If you were previously providing both config values and
   environment values, please ensure the config values are unset if you want to
   use environment values.
 * Namespaces (Enterprise): Providing "root" as the header value for
   `X-Vault-Namespace` will perform the request on the root namespace. This is
   equivalent to providing an empty value. Creating a namespace called "root" in
   the root namespace is disallowed.

FEATURES:

 * **InfluxDB Database Plugin**: Use Vault to dynamically create and manage InfluxDB
   users

IMPROVEMENTS:

 * auth/aws: AWS EC2 authentication can optionally create entity aliases by
   image ID [[GH-5846](https://github.com/hashicorp/vault/pull/5846)]
 * autoseal/gcpckms: Reduce the required permissions for the GCPCKMS autounseal
   [[GH-5999](https://github.com/hashicorp/vault/pull/5999)]
 * physical/foundationdb: TLS support added. [[GH-5800](https://github.com/hashicorp/vault/pull/5800)]

BUG FIXES:

 * api: Fix a couple of places where we were using the `LIST` HTTP verb
   (necessary to get the right method into the wrapping lookup function) and
   not then modifying it to a `GET`; although this is officially the verb Vault
   uses for listing and it's fully legal to use custom verbs, since many WAFs
   and API gateways choke on anything outside of RFC-standardized verbs we fall
   back to `GET` [[GH-6026](https://github.com/hashicorp/vault/pull/6026)]
 * autoseal/aws: Fix reading session tokens when AWS access key/secret key are
   also provided [[GH-5965](https://github.com/hashicorp/vault/pull/5965)]
 * command/operator/rekey: Fix help output showing `-delete-backup` when it
   should show `-backup-delete` [[GH-5981](https://github.com/hashicorp/vault/pull/5981)]
 * core: Fix bound_cidrs not being propagated to child tokens
 * replication: Correctly forward identity entity creation that originates from
   performance standby nodes (Enterprise)
 * secret/aws: Make input `credential_type` match the output type (string, not
   array) [[GH-5972](https://github.com/hashicorp/vault/pull/5972)]
 * secret/cubbyhole: Properly cleanup cubbyhole after token revocation [[GH-6006](https://github.com/hashicorp/vault/pull/6006)]
 * secret/pki: Fix reading certificates on windows with the file storage backend [[GH-6013](https://github.com/hashicorp/vault/pull/6013)]
 * ui (enterprise): properly display perf-standby count on the license page [[GH-5971](https://github.com/hashicorp/vault/pull/5971)]
 * ui: fix disappearing nested secrets and go to the nearest parent when deleting
   a secret - [[GH-5976](https://github.com/hashicorp/vault/pull/5976)]
 * ui: fix error where deleting an item via the context menu would fail if the
   item name contained dots [[GH-6018](https://github.com/hashicorp/vault/pull/6018)]
 * ui: allow saving of kv secret after an errored save attempt [[GH-6022](https://github.com/hashicorp/vault/pull/6022)]
 * ui: fix display of kv-v1 secret containing a key named "keys" [[GH-6023](https://github.com/hashicorp/vault/pull/6023)]

## 1.0.1 (December 14th, 2018)

SECURITY:

 * Update version of Go to 1.11.3 to fix Go bug
   https://github.com/golang/go/issues/29233 which corresponds to
   CVE-2018-16875
 * Database user revocation: If a client has configured custom revocation
   statements for a role with a value of `""`, that statement would be executed
   verbatim, resulting in a lack of actual revocation but success for the
   operation. Vault will now strip empty statements from any provided; as a
   result if an empty statement is provided, it will behave as if no statement
   is provided, falling back to the default revocation statement.

CHANGES:

 * secret/database: On role read, empty statements will be returned as empty
   slices instead of potentially being returned as JSON null values. This makes
   it more in line with other parts of Vault and makes it easier for statically
   typed languages to interpret the values.

IMPROVEMENTS:

 * cli: Strip iTerm extra characters from password manager input [[GH-5837](https://github.com/hashicorp/vault/pull/5837)]
 * command/server: Setting default kv engine to v1 in -dev mode can now be
   specified via -dev-kv-v1 [[GH-5919](https://github.com/hashicorp/vault/pull/5919)]
 * core: Add operationId field to OpenAPI output [[GH-5876](https://github.com/hashicorp/vault/pull/5876)]
 * ui: Added ability to search for Group and Policy IDs when creating Groups
   and Entities instead of typing them in manually

BUG FIXES:

 * auth/azure: Cache azure authorizer [15]
 * auth/gcp: Remove explicit project for service account in GCE authorizer [[GH-58](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/58)]
 * cli: Show correct stored keys/threshold for autoseals [[GH-5910](https://github.com/hashicorp/vault/pull/5910)]
 * cli: Fix backwards compatibility fallback when listing plugins [[GH-5913](https://github.com/hashicorp/vault/pull/5913)]
 * core: Fix upgrades when the seal config had been created on early versions
   of vault [[GH-5956](https://github.com/hashicorp/vault/pull/5956)]
 * namespaces: Correctly reload the proper mount when tuning or reloading the
   mount [[GH-5937](https://github.com/hashicorp/vault/pull/5937)]
 * secret/azure: Cache azure authorizer [19]
 * secret/database: Strip empty statements on user input [[GH-5955](https://github.com/hashicorp/vault/pull/5955)]
 * secret/gcpkms: Add path for retrieving the public key [[GH-5](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/5)]
 * secret/pki: Fix panic that could occur during tidy operation when malformed
   data was found [[GH-5931](https://github.com/hashicorp/vault/pull/5931)]
 * secret/pki: Strip empty line in ca_chain output [[GH-5779](https://github.com/hashicorp/vault/pull/5779)]
 * ui: Fixed a bug where the web CLI was not usable via the `fullscreen`
   command - [[GH-5909](https://github.com/hashicorp/vault/pull/5909)]
 * ui: Fix a bug where you couldn't write a jwt auth method config [[GH-5936](https://github.com/hashicorp/vault/pull/5936)]

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

## 1.0.0 (December 3rd, 2018)

SECURITY:

 * When debugging a customer incident we discovered that in the case of
   malformed data from an autoseal mechanism, Vault's master key could be
   logged in Vault's server log. For this to happen, the data would need to be
   modified by the autoseal mechanism after being submitted to it by Vault but
   prior to encryption, or after decryption, prior to it being returned to
   Vault. To put it another way, it requires the data that Vault submits for
   encryption to not match the data returned after decryption. It is not
   sufficient for the autoseal mechanism to return an error, and it cannot be
   triggered by an outside attacker changing the on-disk ciphertext as all
   autoseal mechanisms use authenticated encryption. We do not believe that
   this is generally a cause for concern; since it involves the autoseal
   mechanism returning bad data to Vault but with no error, in a working Vault
   configuration this code path should never be hit, and if hitting this issue
   Vault will not be unsealing properly anyways so it will be obvious what is
   happening and an immediate rekey of the master key can be performed after
   service is restored. We have filed for a CVE (CVE-2018-19786) and a CVSS V3
   score of 5.2 has been assigned.

CHANGES:

 * Tokens are now prefixed by a designation to indicate what type of token they
   are. Service tokens start with `s.` and batch tokens start with `b.`.
   Existing tokens will still work (they are all of service type and will be
   considered as such). Prefixing allows us to be more efficient when consuming
   a token, which keeps the critical path of requests faster.
 * Paths within `auth/token` that allow specifying a token or accessor in the
   URL have been removed. These have been deprecated since March 2016 and
   undocumented, but were retained for backwards compatibility. They shouldn't
   be used due to the possibility of those paths being logged, so at this point
   they are simply being removed.
 * Vault will no longer accept updates when the storage key has invalid UTF-8
   character encoding [[GH-5819](https://github.com/hashicorp/vault/pull/5819)]
 * Mount/Auth tuning the `options` map on backends will now upsert any provided
   values, and keep any of the existing values in place if not provided. The
   options map itself cannot be unset once it's set, but the keypairs within the
   map can be unset if an empty value is provided, with the exception of the
   `version` keypair which is handled differently for KVv2 purposes.
 * Agent no longer automatically reauthenticates when new credentials are
   detected. It's not strictly necessary and in some cases was causing
   reauthentication much more often than intended.
 * HSM Regenerate Key Support Removed: Vault no longer supports destroying and
   regenerating encryption keys on an HSM; it only supports creating them.
   Although this has never been a source of a customer incident, it is simply a
   code path that is too trivial to activate, especially by mistyping
   `regenerate_key` instead of `generate_key`.
 * Barrier Config Upgrade (Enterprise): When upgrading from Vault 0.8.x, the
   seal type in the barrier config storage entry will be upgraded from
   "hsm-auto" to "awskms" or "pkcs11" upon unseal if using AWSKMS or HSM seals.
   If performing seal migration, the barrier config should first be upgraded
   prior to starting migration.
 * Go API client uses pooled HTTP client: The Go API client now uses a
   connection-pooling HTTP client by default. For CLI operations this makes no
   difference but it should provide significant performance benefits for those
   writing custom clients using the Go API library. As before, this can be
   changed to any custom HTTP client by the caller.
 * Builtin Secret Engines and Auth Methods are integrated deeper into the
   plugin system. The plugin catalog can now override builtin plugins with
   custom versions of the same name. Additionally the plugin system now
   requires a plugin `type` field when configuring plugins, this can be "auth",
   "database", or "secret".

FEATURES:

 * **Auto-Unseal in Open Source**: Cloud-based auto-unseal has been migrated
   from Enterprise to Open Source. We've created a migrator to allow migrating
   between Shamir seals and auto unseal methods.
 * **Batch Tokens**: Batch tokens trade off some features of service tokens for no
   storage overhead, and in most cases can be used across performance
   replication clusters.
 * **Replication Speed Improvements**: We've worked hard to speed up a lot of
   operations when using Vault Enterprise Replication.
 * **GCP KMS Secrets Engine**: This new secrets engine provides a Transit-like
   pattern to keys stored within GCP Cloud KMS.
 * **AppRole support in Vault Agent Auto-Auth**: You can now use AppRole
   credentials when having Agent automatically authenticate to Vault
 * **OpenAPI Support**: Descriptions of mounted backends can be served directly
   from Vault
 * **Kubernetes Projected Service Account Tokens**: Projected Service Account
   Tokens are now supported in Kubernetes auth
 * **Response Wrapping in UI**: Added ability to wrap secrets and easily copy
   the wrap token or secret JSON in the UI

IMPROVEMENTS:

 * agent: Support for configuring the location of the kubernetes service account
   [[GH-5725](https://github.com/hashicorp/vault/pull/5725)]
 * auth/token: New tokens are indexed in storage HMAC-SHA256 instead of SHA1
 * secret/totp: Allow @ character to be part of key name [[GH-5652](https://github.com/hashicorp/vault/pull/5652)]
 * secret/consul: Add support for new policy based tokens added in Consul 1.4
   [[GH-5586](https://github.com/hashicorp/vault/pull/5586)]
 * ui: Improve the token auto-renew warning, and automatically begin renewal
   when a user becomes active again [[GH-5662](https://github.com/hashicorp/vault/pull/5662)]
 * ui: The unbundled UI page now has some styling [[GH-5665](https://github.com/hashicorp/vault/pull/5665)]
 * ui: Improved banner and popup design [[GH-5672](https://github.com/hashicorp/vault/pull/5672)]
 * ui: Added token type to auth method mount config [[GH-5723](https://github.com/hashicorp/vault/pull/5723)]
 * ui: Display additonal wrap info when unwrapping. [[GH-5664](https://github.com/hashicorp/vault/pull/5664)]
 * ui: Empty states have updated styling and link to relevant actions and
   documentation [[GH-5758](https://github.com/hashicorp/vault/pull/5758)]
 * ui: Allow editing of KV V2 data when a token doesn't have capabilities to
   read secret metadata [[GH-5879](https://github.com/hashicorp/vault/pull/5879)]

BUG FIXES:

 * agent: Fix auth when multiple redirects [[GH-5814](https://github.com/hashicorp/vault/pull/5814)]
 * cli: Restore the `-policy-override` flag [[GH-5826](https://github.com/hashicorp/vault/pull/5826)]
 * core: Fix rekey progress reset which did not happen under certain
   circumstances. [[GH-5743](https://github.com/hashicorp/vault/pull/5743)]
 * core: Migration from autounseal to shamir will clean up old keys [[GH-5671](https://github.com/hashicorp/vault/pull/5671)]
 * identity: Update group memberships when entity is deleted [[GH-5786](https://github.com/hashicorp/vault/pull/5786)]
 * replication/perfstandby: Fix audit table upgrade on standbys [[GH-5811](https://github.com/hashicorp/vault/pull/5811)]
 * replication/perfstandby: Fix redirect on approle update [[GH-5820](https://github.com/hashicorp/vault/pull/5820)]
 * secrets/azure: Fix valid roles being rejected for duplicate ids despite
   having distinct scopes
   [[GH-16](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/16)]
 * storage/gcs: Send md5 of values to GCS to avoid potential corruption
   [[GH-5804](https://github.com/hashicorp/vault/pull/5804)]
 * secrets/kv: Fix issue where storage version would get incorrectly downgraded
   [[GH-5809](https://github.com/hashicorp/vault/pull/5809)]
 * secrets/kv: Disallow empty paths on a `kv put` while accepting empty paths
   for all other operations for backwards compatibility
   [[GH-19](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/19)]
 * ui: Allow for secret creation in kv v2 when cas_required=true [[GH-5823](https://github.com/hashicorp/vault/pull/5823)]
 * ui: Fix dr secondary operation token generation via the ui [[GH-5818](https://github.com/hashicorp/vault/pull/5818)]
 * ui: Fix the PKI context menu so that items load [[GH-5824](https://github.com/hashicorp/vault/pull/5824)]
 * ui: Update DR Secondary Token generation command [[GH-5857](https://github.com/hashicorp/vault/pull/5857)]
 * ui: Fix pagination bug where controls would be rendered once for each
   item when viewing policies [[GH-5866](https://github.com/hashicorp/vault/pull/5866)]
 * ui: Fix bug where `sys/leases/revoke` required 'sudo' capability to show
   the revoke button in the UI [[GH-5647](https://github.com/hashicorp/vault/pull/5647)]
 * ui: Fix issue where certain pages wouldn't render in a namespace [[GH-5692](https://github.com/hashicorp/vault/pull/5692)]

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
