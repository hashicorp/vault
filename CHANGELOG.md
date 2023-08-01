## 1.14.1
### July 25, 2023

CHANGES:

* auth/ldap: Normalize HTTP response codes when invalid credentials are provided [[GH-21282](https://github.com/hashicorp/vault/pull/21282)]
* core/namespace (enterprise): Introduce the concept of high-privilege namespace (administrative namespace),
which will have access to some system backend paths that were previously only accessible in the root namespace. [[GH-21215](https://github.com/hashicorp/vault/pull/21215)]
* secrets/transform (enterprise): Enforce a transformation role's max_ttl setting on encode requests, a warning will be returned if max_ttl was applied.
* storage/aerospike: Aerospike storage shouldn't be used on 32-bit architectures and is now unsupported on them. [[GH-20825](https://github.com/hashicorp/vault/pull/20825)]

IMPROVEMENTS:

* core/fips: Add RPM, DEB packages of FIPS 140-2 and HSM+FIPS 140-2 Vault Enterprise.
* eventbus: updated go-eventlogger library to allow removal of nodes referenced by pipelines (used for subscriptions) [[GH-21623](https://github.com/hashicorp/vault/pull/21623)]
* openapi: Better mount points for kv-v1 and kv-v2 in openapi.json [[GH-21563](https://github.com/hashicorp/vault/pull/21563)]
* replication (enterprise): Avoid logging warning if request is forwarded from a performance standby and not a performance secondary
* secrets/pki: Add a parameter to allow ExtKeyUsage field usage from a role within ACME. [[GH-21702](https://github.com/hashicorp/vault/pull/21702)]
* secrets/transform (enterprise): Switch to pgx PostgreSQL driver for better timeout handling
* sys/metrics (enterprise): Adds a gauge metric that tracks whether enterprise builtin secret plugins are enabled. [[GH-21681](https://github.com/hashicorp/vault/pull/21681)]

BUG FIXES:

* agent: Fix "generate-config" command documentation URL [[GH-21466](https://github.com/hashicorp/vault/pull/21466)]
* auth/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21800](https://github.com/hashicorp/vault/pull/21800)]
* auth/token, sys: Fix path-help being unavailable for some list-only endpoints [[GH-18571](https://github.com/hashicorp/vault/pull/18571)]
* auth/token: Fix parsing of `auth/token/create` fields to avoid incorrect warnings about ignored parameters [[GH-18556](https://github.com/hashicorp/vault/pull/18556)]
* awsutil: Update awsutil to v0.2.3 to fix a regression where Vault no longer
respects `AWS_ROLE_ARN`, `AWS_WEB_IDENTITY_TOKEN_FILE`, and `AWS_ROLE_SESSION_NAME`. [[GH-21951](https://github.com/hashicorp/vault/pull/21951)]
* core/managed-keys (enterprise): Allow certain symmetric PKCS#11 managed key mechanisms (AES CBC with and without padding) to operate without an HMAC.
* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-24170](https://github.com/hashicorp/vault/pull/24170)]
* core: Fixed issue with some durations not being properly parsed to include days. [[GH-21357](https://github.com/hashicorp/vault/pull/21357)]
* identity: Remove caseSensitivityKey to prevent errors while loading groups which could result in missing groups in memDB when duplicates are found. [[GH-20965](https://github.com/hashicorp/vault/pull/20965)]
* openapi: Fix response schema for PKI Issue requests [[GH-21449](https://github.com/hashicorp/vault/pull/21449)]
* openapi: Fix schema definitions for PKI EAB APIs [[GH-21458](https://github.com/hashicorp/vault/pull/21458)]
* replication (enterprise): update primary cluster address after DR failover
* secrets/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21631](https://github.com/hashicorp/vault/pull/21631)]
* secrets/pki: Fix bug with ACME tidy, 'unable to determine acme base folder path'. [[GH-21870](https://github.com/hashicorp/vault/pull/21870)]
* secrets/pki: Fix preserving acme_account_safety_buffer on config/auto-tidy. [[GH-21870](https://github.com/hashicorp/vault/pull/21870)]
* secrets/pki: Prevent deleted issuers from reappearing when migrating from a version 1 bundle to a version 2 bundle (versions including 1.13.0, 1.12.2, and 1.11.6); when managed keys were removed but referenced in the Vault 1.10 legacy CA bundle, this the error: `no managed key found with uuid`. [[GH-21316](https://github.com/hashicorp/vault/pull/21316)]
* secrets/transform (enterprise): Fix nil panic when deleting a template with tokenization transformations present
* secrets/transform (enterprise): Grab shared locks for various read operations, only escalating to write locks if work is required
* serviceregistration: Fix bug where multiple nodes in a secondary cluster could be labelled active after updating the cluster's primary [[GH-21642](https://github.com/hashicorp/vault/pull/21642)]
* ui: Adds missing values to details view after generating PKI certificate [[GH-21635](https://github.com/hashicorp/vault/pull/21635)]
* ui: Fixed an issue where editing an SSH role would clear `default_critical_options` and `default_extension` if left unchanged. [[GH-21739](https://github.com/hashicorp/vault/pull/21739)]
* ui: Fixed secrets, leases, and policies filter dropping focus after a single character [[GH-21767](https://github.com/hashicorp/vault/pull/21767)]
* ui: Fixes issue with certain navigational links incorrectly displaying in child namespaces [[GH-21562](https://github.com/hashicorp/vault/pull/21562)]
* ui: Fixes login screen display issue with Safari browser [[GH-21582](https://github.com/hashicorp/vault/pull/21582)]
* ui: Fixes problem displaying certificates issued with unsupported signature algorithms (i.e. ed25519) [[GH-21926](https://github.com/hashicorp/vault/pull/21926)]
* ui: Fixes styling of private key input when configuring an SSH key [[GH-21531](https://github.com/hashicorp/vault/pull/21531)]
* ui: Surface DOMException error when browser settings prevent localStorage. [[GH-21503](https://github.com/hashicorp/vault/pull/21503)]

## 1.14.0
### June 21, 2023

BREAKING CHANGES:

* secrets/pki: Maintaining running count of certificates will be turned off by default.
To re-enable keeping these metrics available on the tidy status endpoint, enable
maintain_stored_certificate_counts on tidy-config, to also publish them to the
metrics consumer, enable publish_stored_certificate_count_metrics . [[GH-18186](https://github.com/hashicorp/vault/pull/18186)]

CHANGES:

* auth/alicloud: Updated plugin from v0.14.0 to v0.15.0 [[GH-20758](https://github.com/hashicorp/vault/pull/20758)]
* auth/azure: Updated plugin from v0.13.0 to v0.15.0 [[GH-20816](https://github.com/hashicorp/vault/pull/20816)]
* auth/centrify: Updated plugin from v0.14.0 to v0.15.1 [[GH-20745](https://github.com/hashicorp/vault/pull/20745)]
* auth/gcp: Updated plugin from v0.15.0 to v0.16.0 [[GH-20725](https://github.com/hashicorp/vault/pull/20725)]
* auth/jwt: Updated plugin from v0.15.0 to v0.16.0 [[GH-20799](https://github.com/hashicorp/vault/pull/20799)]
* auth/kubernetes: Update plugin to v0.16.0 [[GH-20802](https://github.com/hashicorp/vault/pull/20802)]
* core: Bump Go version to 1.20.5.
* core: Remove feature toggle for SSCTs, i.e. the env var VAULT_DISABLE_SERVER_SIDE_CONSISTENT_TOKENS. [[GH-20834](https://github.com/hashicorp/vault/pull/20834)]
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]
* database/couchbase: Updated plugin from v0.9.0 to v0.9.2 [[GH-20764](https://github.com/hashicorp/vault/pull/20764)]
* database/redis-elasticache: Updated plugin from v0.2.0 to v0.2.1 [[GH-20751](https://github.com/hashicorp/vault/pull/20751)]
* replication (enterprise): Add a new parameter for the update-primary API call
that allows for setting of the primary cluster addresses directly, instead of
via a token.
* secrets/ad: Updated plugin from v0.10.1-0.20230329210417-0b2cdb26cf5d to v0.16.0 [[GH-20750](https://github.com/hashicorp/vault/pull/20750)]
* secrets/alicloud: Updated plugin from v0.5.4-beta1.0.20230330124709-3fcfc5914a22 to v0.15.0 [[GH-20787](https://github.com/hashicorp/vault/pull/20787)]
* secrets/aure: Updated plugin from v0.15.0 to v0.16.0 [[GH-20777](https://github.com/hashicorp/vault/pull/20777)]
* secrets/database/mongodbatlas: Updated plugin from v0.9.0 to v0.10.0 [[GH-20882](https://github.com/hashicorp/vault/pull/20882)]
* secrets/database/snowflake: Updated plugin from v0.7.0 to v0.8.0 [[GH-20807](https://github.com/hashicorp/vault/pull/20807)]
* secrets/gcp: Updated plugin from v0.15.0 to v0.16.0 [[GH-20818](https://github.com/hashicorp/vault/pull/20818)]
* secrets/keymgmt: Updated plugin to v0.9.1
* secrets/kubernetes: Update plugin to v0.5.0 [[GH-20802](https://github.com/hashicorp/vault/pull/20802)]
* secrets/mongodbatlas: Updated plugin from v0.9.1 to v0.10.0 [[GH-20742](https://github.com/hashicorp/vault/pull/20742)]
* secrets/pki: Allow issuance of root CAs without AIA, when templated AIA information includes issuer_id. [[GH-21209](https://github.com/hashicorp/vault/pull/21209)]
* secrets/pki: Warning when issuing leafs from CSRs with basic constraints. In the future, issuance of non-CA leaf certs from CSRs with asserted IsCA Basic Constraints will be prohibited. [[GH-20654](https://github.com/hashicorp/vault/pull/20654)]

FEATURES:

* **AWS Static Roles**: The AWS Secrets Engine can manage static roles configured by users. [[GH-20536](https://github.com/hashicorp/vault/pull/20536)]
* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* **Environment Variables through Vault Agent**: Introducing a new process-supervisor mode for Vault Agent which allows injecting secrets as environment variables into a child process using a new `env_template` configuration stanza. The process-supervisor configuration can be generated with a new `vault agent generate-config` helper tool. [[GH-20530](https://github.com/hashicorp/vault/pull/20530)]
* **MongoDB Atlas Database Secrets**: Adds support for client certificate credentials [[GH-20425](https://github.com/hashicorp/vault/pull/20425)]
* **MongoDB Atlas Database Secrets**: Adds support for generating X.509 certificates on dynamic roles for user authentication [[GH-20882](https://github.com/hashicorp/vault/pull/20882)]
* **NEW PKI Workflow in UI**: Completes generally available rollout of new PKI UI that provides smoother mount configuration and a more guided user experience [[GH-pki-ui-improvements](https://github.com/hashicorp/vault/pull/pki-ui-improvements)]
* **Secrets/Auth Plugin Multiplexing**: The plugin will be multiplexed when run
as an external plugin by vault versions that support secrets/auth plugin
multiplexing (> 1.12) [[GH-19215](https://github.com/hashicorp/vault/pull/19215)]
* **Sidebar Navigation in UI**: A new sidebar navigation panel has been added in the UI to replace the top navigation bar. [[GH-19296](https://github.com/hashicorp/vault/pull/19296)]
* **Vault PKI ACME Server**: Support for the ACME certificate lifecycle management protocol has been added to the Vault PKI Plugin. This allows standard ACME clients, such as the EFF's certbot and the CNCF's k8s cert-manager, to request certificates from a Vault server with no knowledge of Vault APIs or authentication mechanisms. For public-facing Vault instances, we recommend requiring External Account Bindings (EAB) to limit the ability to request certificates to only authenticated clients. [[GH-20752](https://github.com/hashicorp/vault/pull/20752)]
* **Vault Proxy**: Introduced Vault Proxy, a new subcommand of the Vault binary that can be invoked using `vault proxy -config=config.hcl`. It currently has the same feature set as Vault Agent's API proxy, but the two may diverge in the future. We plan to deprecate the API proxy functionality of Vault Agent in a future release. [[GH-20548](https://github.com/hashicorp/vault/pull/20548)]
* **OCI Auto-Auth**: Add OCI (Oracle Cloud Infrastructure) auto-auth method [[GH-19260](https://github.com/hashicorp/vault/pull/19260)]

IMPROVEMENTS:

* * api: Add Config.TLSConfig method to fetch the TLS configuration from a client config. [[GH-20265](https://github.com/hashicorp/vault/pull/20265)]
* * physical/etcd: Upgrade etcd3 client to v3.5.7 [[GH-20261](https://github.com/hashicorp/vault/pull/20261)]
* activitylog: EntityRecord protobufs now contain a ClientType field for
distinguishing client sources. [[GH-20626](https://github.com/hashicorp/vault/pull/20626)]
* agent: Add integration tests for agent running in process supervisor mode [[GH-20741](https://github.com/hashicorp/vault/pull/20741)]
* agent: Add logic to validate env_template entries in configuration [[GH-20569](https://github.com/hashicorp/vault/pull/20569)]
* agent: Added `reload` option to cert auth configuration in case of external renewals of local x509 key-pairs. [[GH-19002](https://github.com/hashicorp/vault/pull/19002)]
* agent: JWT auto-auth has a new config option, `remove_jwt_follows_symlinks` (default: false), that, if set to true will now remove the JWT, instead of the symlink to the JWT, if a symlink to a JWT has been provided in the `path` option, and the `remove_jwt_after_reading` config option is set to true (default). [[GH-18863](https://github.com/hashicorp/vault/pull/18863)]
* agent: Vault Agent now reports its name and version as part of the User-Agent header in all requests issued. [[GH-19776](https://github.com/hashicorp/vault/pull/19776)]
* agent: initial implementation of a process runner for injecting secrets via environment variables via vault agent [[GH-20628](https://github.com/hashicorp/vault/pull/20628)]
* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* api: property based testing for LifetimeWatcher sleep duration calculation [[GH-17919](https://github.com/hashicorp/vault/pull/17919)]
* audit: add plugin metadata, including plugin name, type, version, sha256, and whether plugin is external, to audit logging [[GH-19814](https://github.com/hashicorp/vault/pull/19814)]
* audit: forwarded requests can now contain host metadata on the node it was sent 'from' or a flag to indicate that it was forwarded.
* auth/cert: Better return OCSP validation errors during login to the caller. [[GH-20234](https://github.com/hashicorp/vault/pull/20234)]
* auth/kerberos: Enable plugin multiplexing
auth/kerberos: Upgrade plugin dependencies [[GH-20771](https://github.com/hashicorp/vault/pull/20771)]
* auth/ldap: allow configuration of alias dereferencing in LDAP search [[GH-18230](https://github.com/hashicorp/vault/pull/18230)]
* auth/ldap: allow providing the LDAP password via an env var when authenticating via the CLI [[GH-18225](https://github.com/hashicorp/vault/pull/18225)]
* auth/oidc: Adds support for group membership parsing when using IBM ISAM as an OIDC provider. [[GH-19247](https://github.com/hashicorp/vault/pull/19247)]
* build: Prefer GOBIN when set over GOPATH/bin when building the binary [[GH-19862](https://github.com/hashicorp/vault/pull/19862)]
* cli: Add walkSecretsTree helper function, which recursively walks secrets rooted at the given path [[GH-20464](https://github.com/hashicorp/vault/pull/20464)]
* cli: Improve addPrefixToKVPath helper [[GH-20488](https://github.com/hashicorp/vault/pull/20488)]
* command/server (enterprise): -dev-three-node now creates perf standbys instead of regular standbys. [[GH-20629](https://github.com/hashicorp/vault/pull/20629)]
* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* command/server: New -dev-cluster-json writes a file describing the dev cluster in -dev and -dev-three-node modes, plus -dev-three-node now enables unauthenticated metrics and pprof requests. [[GH-20224](https://github.com/hashicorp/vault/pull/20224)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core, secrets/pki, audit: Update dependency go-jose to v3 due to v2 deprecation. [[GH-20559](https://github.com/hashicorp/vault/pull/20559)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* core: Add possibility to decode a generated encoded root token via the rest API [[GH-20595](https://github.com/hashicorp/vault/pull/20595)]
* core: include namespace path in granting_policies block of audit log
* core: include reason for ErrReadOnly on PBPWF writing failures
* core: report intermediate error messages during request forwarding [[GH-20643](https://github.com/hashicorp/vault/pull/20643)]
* core:provide more descriptive error message when calling enterprise feature paths in open-source [[GH-18870](https://github.com/hashicorp/vault/pull/18870)]
* database/elasticsearch: Upgrade plugin dependencies [[GH-20767](https://github.com/hashicorp/vault/pull/20767)]
* database/mongodb: upgrade mongo driver to 1.11 [[GH-19954](https://github.com/hashicorp/vault/pull/19954)]
* database/redis: Upgrade plugin dependencies [[GH-20763](https://github.com/hashicorp/vault/pull/20763)]
* http: Support responding to HEAD operation from plugins [[GH-19520](https://github.com/hashicorp/vault/pull/19520)]
* openapi: Add openapi response definitions to /sys defined endpoints. [[GH-18633](https://github.com/hashicorp/vault/pull/18633)]
* openapi: Add openapi response definitions to pki/config_*.go [[GH-18376](https://github.com/hashicorp/vault/pull/18376)]
* openapi: Add openapi response definitions to vault/logical_system_paths.go defined endpoints. [[GH-18515](https://github.com/hashicorp/vault/pull/18515)]
* openapi: Consistently stop Vault server on exit in gen_openapi.sh [[GH-19252](https://github.com/hashicorp/vault/pull/19252)]
* openapi: Improve operationId/request/response naming strategy [[GH-19319](https://github.com/hashicorp/vault/pull/19319)]
* openapi: add openapi response definitions to /sys/internal endpoints [[GH-18542](https://github.com/hashicorp/vault/pull/18542)]
* openapi: add openapi response definitions to /sys/rotate endpoints [[GH-18624](https://github.com/hashicorp/vault/pull/18624)]
* openapi: add openapi response definitions to /sys/seal endpoints [[GH-18625](https://github.com/hashicorp/vault/pull/18625)]
* openapi: add openapi response definitions to /sys/tool endpoints [[GH-18626](https://github.com/hashicorp/vault/pull/18626)]
* openapi: add openapi response definitions to /sys/version-history, /sys/leader, /sys/ha-status, /sys/host-info, /sys/in-flight-req [[GH-18628](https://github.com/hashicorp/vault/pull/18628)]
* openapi: add openapi response definitions to /sys/wrapping endpoints [[GH-18627](https://github.com/hashicorp/vault/pull/18627)]
* openapi: add openapi response defintions to /sys/auth endpoints [[GH-18465](https://github.com/hashicorp/vault/pull/18465)]
* openapi: add openapi response defintions to /sys/capabilities endpoints [[GH-18468](https://github.com/hashicorp/vault/pull/18468)]
* openapi: add openapi response defintions to /sys/config and /sys/generate-root endpoints [[GH-18472](https://github.com/hashicorp/vault/pull/18472)]
* openapi: added ability to validate response structures against openapi schema for test clusters [[GH-19043](https://github.com/hashicorp/vault/pull/19043)]
* sdk/framework: Fix non-deterministic ordering of 'required' fields in OpenAPI spec [[GH-20881](https://github.com/hashicorp/vault/pull/20881)]
* sdk: Add new docker-based cluster testing framework to the sdk. [[GH-20247](https://github.com/hashicorp/vault/pull/20247)]
* secrets/ad: upgrades dependencies [[GH-19829](https://github.com/hashicorp/vault/pull/19829)]
* secrets/alicloud: upgrades dependencies [[GH-19846](https://github.com/hashicorp/vault/pull/19846)]
* secrets/consul: Improve error message when ACL bootstrapping fails. [[GH-20891](https://github.com/hashicorp/vault/pull/20891)]
* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]
* secrets/gcpkms: Enable plugin multiplexing
secrets/gcpkms: Upgrade plugin dependencies [[GH-20784](https://github.com/hashicorp/vault/pull/20784)]
* secrets/mongodbatlas: upgrades dependencies [[GH-19861](https://github.com/hashicorp/vault/pull/19861)]
* secrets/openldap: upgrades dependencies [[GH-19993](https://github.com/hashicorp/vault/pull/19993)]
* secrets/pki: Add missing fields to tidy-status, include new last_auto_tidy_finished field. [[GH-20442](https://github.com/hashicorp/vault/pull/20442)]
* secrets/pki: Add warning when issuer lacks KeyUsage during CRL rebuilds; expose in logs and on rotation. [[GH-20253](https://github.com/hashicorp/vault/pull/20253)]
* secrets/pki: Allow determining existing issuers and keys on import. [[GH-20441](https://github.com/hashicorp/vault/pull/20441)]
* secrets/pki: Include CA serial number, key UUID on issuers list endpoint. [[GH-20276](https://github.com/hashicorp/vault/pull/20276)]
* secrets/pki: Limit ACME issued certificates NotAfter TTL to a maximum of 90 days [[GH-20981](https://github.com/hashicorp/vault/pull/20981)]
* secrets/pki: Support TLS-ALPN-01 challenge type in ACME for DNS certificate identifiers. [[GH-20943](https://github.com/hashicorp/vault/pull/20943)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]
* secrets/postgresql: Add configuration to scram-sha-256 encrypt passwords on Vault before sending them to PostgreSQL [[GH-19616](https://github.com/hashicorp/vault/pull/19616)]
* secrets/terraform: upgrades dependencies [[GH-19798](https://github.com/hashicorp/vault/pull/19798)]
* secrets/transit: Add support to import public keys in transit engine and allow encryption and verification of signed data [[GH-17934](https://github.com/hashicorp/vault/pull/17934)]
* secrets/transit: Allow importing RSA-PSS OID (1.2.840.113549.1.1.10) private keys via BYOK. [[GH-19519](https://github.com/hashicorp/vault/pull/19519)]
* secrets/transit: Respond to writes with updated key policy, cache configuration. [[GH-20652](https://github.com/hashicorp/vault/pull/20652)]
* secrets/transit: Support BYOK-encrypted export of keys to securely allow synchronizing specific keys and version across clusters. [[GH-20736](https://github.com/hashicorp/vault/pull/20736)]
* ui: Add download button for each secret value in KV v2 [[GH-20431](https://github.com/hashicorp/vault/pull/20431)]
* ui: Add filtering by auth type and auth name to the Authentication Method list view. [[GH-20747](https://github.com/hashicorp/vault/pull/20747)]
* ui: Add filtering by engine type and engine name to the Secret Engine list view. [[GH-20481](https://github.com/hashicorp/vault/pull/20481)]
* ui: Adds whitespace warning to secrets engine and auth method path inputs [[GH-19913](https://github.com/hashicorp/vault/pull/19913)]
* ui: Remove the Bulma CSS framework. [[GH-19878](https://github.com/hashicorp/vault/pull/19878)]
* ui: Update Web CLI with examples and a new `kv-get` command for reading kv v2 data and metadata [[GH-20590](https://github.com/hashicorp/vault/pull/20590)]
* ui: Updates UI javascript dependencies [[GH-19901](https://github.com/hashicorp/vault/pull/19901)]
* ui: add allowed_managed_keys field to secret engine mount options [[GH-19791](https://github.com/hashicorp/vault/pull/19791)]
* ui: adds warning for commas in stringArray inputs and updates tooltip help text to remove references to comma separation [[GH-20163](https://github.com/hashicorp/vault/pull/20163)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]
* website/docs: Add rotate root documentation for azure secrets engine [[GH-19187](https://github.com/hashicorp/vault/pull/19187)]
* website/docs: fix database static-user sample payload [[GH-19170](https://github.com/hashicorp/vault/pull/19170)]

BUG FIXES:

* agent: Fix agent generate-config to accept -namespace, VAULT_NAMESPACE, and other client-modifying flags. [[GH-21297](https://github.com/hashicorp/vault/pull/21297)]
* agent: Fix bug with 'cache' stanza validation [[GH-20934](https://github.com/hashicorp/vault/pull/20934)]
* api: Addressed a couple of issues that arose as edge cases for the -output-policy flag. Specifically around properly handling list commands, distinguishing kv V1/V2, and correctly recognizing protected paths. [[GH-19160](https://github.com/hashicorp/vault/pull/19160)]
* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* auth/token: Fix cubbyhole and revocation for legacy service tokens [[GH-19416](https://github.com/hashicorp/vault/pull/19416)]
* cli/kv: add -mount flag to kv list [[GH-19378](https://github.com/hashicorp/vault/pull/19378)]
* core (enterprise): Don't delete backend stored data that appears to be filterable
on this secondary if we don't have a corresponding mount entry.
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry.  Fix wrong lock used in ListAuths link meta interface implementation. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix Forwarded Writer construction to correctly find active nodes, allowing PKI cross-cluster functionality to succeed on existing mounts.
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* core: Fix writes to readonly storage on performance standbys when user lockout feature is enabled. [[GH-20783](https://github.com/hashicorp/vault/pull/20783)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* license (enterprise): Fix bug where license would update even if the license didn't change.
* openapi: Small fixes for OpenAPI display attributes. Changed "log-in" to "login" [[GH-20285](https://github.com/hashicorp/vault/pull/20285)]
* plugin/reload:  Fix a possible data race with rollback manager and plugin reload [[GH-19468](https://github.com/hashicorp/vault/pull/19468)]
* replication (enterprise): Fix a caching issue when replicating filtered data to
a performance secondary. This resulted in the data being set to nil in the cache
and a "invalid value" error being returned from the API.
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs
* replication (enterprise): Fix regression causing token creation against a role
with a new entity alias to be incorrectly forwarded from perf standbys. [[GH-21100](https://github.com/hashicorp/vault/pull/21100)]
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* sdk/backend: prevent panic when computing the zero value for a `TypeInt64` schema field. [[GH-18729](https://github.com/hashicorp/vault/pull/18729)]
* secrets/pki: Support setting both maintain_stored_certificate_counts=false and publish_stored_certificate_count_metrics=false explicitly in tidy config. [[GH-20664](https://github.com/hashicorp/vault/pull/20664)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* secrets/transform: Added importing of keys and key versions into the Transform secrets engine using the command 'vault transform import' and 'vault transform import-version'. [[GH-20668](https://github.com/hashicorp/vault/pull/20668)]
* secrets/transit: Fix export of HMAC-only key, correctly exporting the key used for sign operations. For consumers of the previously incorrect key, use the plaintext export to retrieve these incorrect keys and import them as new versions.
* secrets/transit: Fix bug related to shorter dedicated HMAC key sizing.
* sdk/helper/keysutil: New HMAC type policies will have HMACKey equal to Key and be copied over on import. [[GH-20864](https://github.com/hashicorp/vault/pull/20864)]
* shamir: change mul and div implementations to be constant-time [[GH-19495](https://github.com/hashicorp/vault/pull/19495)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix secret render when path includes %. Resolves #11616. [[GH-20430](https://github.com/hashicorp/vault/pull/20430)]
* ui: Fixes issue unsealing cluster for seal types other than shamir [[GH-20897](https://github.com/hashicorp/vault/pull/20897)]
* ui: fixes auto_rotate_period ttl input for transit keys [[GH-20731](https://github.com/hashicorp/vault/pull/20731)]
* ui: fixes bug in kmip role form that caused `operation_all` to persist after deselecting all operation checkboxes [[GH-19139](https://github.com/hashicorp/vault/pull/19139)]
* ui: fixes key_bits and signature_bits reverting to default values when editing a pki role [[GH-20907](https://github.com/hashicorp/vault/pull/20907)]
* ui: wait for wanted message event during OIDC callback instead of using the first message event [[GH-18521](https://github.com/hashicorp/vault/pull/18521)]

## 1.13.5
### July 25, 2023

CHANGES:

* auth/ldap: Normalize HTTP response codes when invalid credentials are provided [[GH-21282](https://github.com/hashicorp/vault/pull/21282)]
* core/namespace (enterprise): Introduce the concept of high-privilege namespace (administrative namespace),
which will have access to some system backend paths that were previously only accessible in the root namespace. [[GH-21215](https://github.com/hashicorp/vault/pull/21215)]
* secrets/transform (enterprise): Enforce a transformation role's max_ttl setting on encode requests, a warning will be returned if max_ttl was applied.

IMPROVEMENTS:

* core/fips: Add RPM, DEB packages of FIPS 140-2 and HSM+FIPS 140-2 Vault Enterprise.
* core: Add a new periodic metric to track the number of available policies, `vault.policy.configured.count`. [[GH-21010](https://github.com/hashicorp/vault/pull/21010)]
* replication (enterprise): Avoid logging warning if request is forwarded from a performance standby and not a performance secondary
* secrets/transform (enterprise): Switch to pgx PostgreSQL driver for better timeout handling
* sys/metrics (enterprise): Adds a gauge metric that tracks whether enterprise builtin secret plugins are enabled. [[GH-21681](https://github.com/hashicorp/vault/pull/21681)]

BUG FIXES:

* auth/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21799](https://github.com/hashicorp/vault/pull/21799)]
* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-24170](https://github.com/hashicorp/vault/pull/24170)]
* identity: Remove caseSensitivityKey to prevent errors while loading groups which could result in missing groups in memDB when duplicates are found. [[GH-20965](https://github.com/hashicorp/vault/pull/20965)]
* replication (enterprise): update primary cluster address after DR failover
* secrets/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21632](https://github.com/hashicorp/vault/pull/21632)]
* secrets/pki: Prevent deleted issuers from reappearing when migrating from a version 1 bundle to a version 2 bundle (versions including 1.13.0, 1.12.2, and 1.11.6); when managed keys were removed but referenced in the Vault 1.10 legacy CA bundle, this the error: `no managed key found with uuid`. [[GH-21316](https://github.com/hashicorp/vault/pull/21316)]
* secrets/pki: Support setting both maintain_stored_certificate_counts=false and publish_stored_certificate_count_metrics=false explicitly in tidy config. [[GH-20664](https://github.com/hashicorp/vault/pull/20664)]
* secrets/transform (enterprise): Fix nil panic when deleting a template with tokenization transformations present
* secrets/transform (enterprise): Grab shared locks for various read operations, only escalating to write locks if work is required
* serviceregistration: Fix bug where multiple nodes in a secondary cluster could be labelled active after updating the cluster's primary [[GH-21642](https://github.com/hashicorp/vault/pull/21642)]
* ui: Fixed an issue where editing an SSH role would clear `default_critical_options` and `default_extension` if left unchanged. [[GH-21739](https://github.com/hashicorp/vault/pull/21739)]
* ui: Surface DOMException error when browser settings prevent localStorage. [[GH-21503](https://github.com/hashicorp/vault/pull/21503)]

## 1.13.4
### June 21, 2023
BREAKING CHANGES:

* secrets/pki: Maintaining running count of certificates will be turned off by default.
To re-enable keeping these metrics available on the tidy status endpoint, enable
maintain_stored_certificate_counts on tidy-config, to also publish them to the
metrics consumer, enable publish_stored_certificate_count_metrics . [[GH-18186](https://github.com/hashicorp/vault/pull/18186)]

CHANGES:

* core: Bump Go version to 1.20.5.

FEATURES:

* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* core (enterprise): Add background worker for automatic reporting of billing
information. [[GH-19625](https://github.com/hashicorp/vault/pull/19625)]

IMPROVEMENTS:

* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]

BUG FIXES:

* agent: Fix bug with 'cache' stanza validation [[GH-20934](https://github.com/hashicorp/vault/pull/20934)]
* core (enterprise): Don't delete backend stored data that appears to be filterable
on this secondary if we don't have a corresponding mount entry.
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry.  Fix wrong lock used in ListAuths link meta interface implementation. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs
* replication (enterprise): Fix regression causing token creation against a role
with a new entity alias to be incorrectly forwarded from perf standbys. [[GH-21100](https://github.com/hashicorp/vault/pull/21100)]
* storage/raft: Fix race where new follower joining can get pruned by dead server cleanup. [[GH-20986](https://github.com/hashicorp/vault/pull/20986)]

## 1.13.3
### June 08, 2023

CHANGES:

* core: Bump Go version to 1.20.4.
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]
* replication (enterprise): Add a new parameter for the update-primary API call
that allows for setting of the primary cluster addresses directly, instead of
via a token.
* storage/aerospike: Aerospike storage shouldn't be used on 32-bit architectures and is now unsupported on them. [[GH-20825](https://github.com/hashicorp/vault/pull/20825)]

IMPROVEMENTS:

* Add debug symbols back to builds to fix Dynatrace support [[GH-20519](https://github.com/hashicorp/vault/pull/20519)]
* audit: add a `mount_point` field to audit requests and response entries [[GH-20411](https://github.com/hashicorp/vault/pull/20411)]
* autopilot: Update version to v0.2.0 to add better support for respecting min quorum [[GH-19472](https://github.com/hashicorp/vault/pull/19472)]
* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* core: Add possibility to decode a generated encoded root token via the rest API [[GH-20595](https://github.com/hashicorp/vault/pull/20595)]
* core: include namespace path in granting_policies block of audit log
* core: report intermediate error messages during request forwarding [[GH-20643](https://github.com/hashicorp/vault/pull/20643)]
* openapi: Fix generated types for duration strings [[GH-20841](https://github.com/hashicorp/vault/pull/20841)]
* sdk/framework: Fix non-deterministic ordering of 'required' fields in OpenAPI spec [[GH-20881](https://github.com/hashicorp/vault/pull/20881)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]

BUG FIXES:

* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* cli: CLI should take days as a unit of time for ttl like flags [[GH-20477](https://github.com/hashicorp/vault/pull/20477)]
* cli: disable printing flags warnings messages for the ssh command [[GH-20502](https://github.com/hashicorp/vault/pull/20502)]
* command/server: fixes panic in Vault server command when running in recovery mode [[GH-20418](https://github.com/hashicorp/vault/pull/20418)]
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core/identity: Allow updates of only the custom-metadata for entity alias. [[GH-20368](https://github.com/hashicorp/vault/pull/20368)]
* core: Fix Forwarded Writer construction to correctly find active nodes, allowing PKI cross-cluster functionality to succeed on existing mounts.
* core: Fix writes to readonly storage on performance standbys when user lockout feature is enabled. [[GH-20783](https://github.com/hashicorp/vault/pull/20783)]
* core: prevent panic on login after namespace is deleted that had mfa enforcement [[GH-20375](https://github.com/hashicorp/vault/pull/20375)]
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* secrets/pki: Include per-issuer enable_aia_url_templating in issuer read endpoint. [[GH-20354](https://github.com/hashicorp/vault/pull/20354)]
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation
* secrets/transform: Added importing of keys and key versions into the Transform secrets engine using the command 'vault transform import' and 'vault transform import-version'. [[GH-20668](https://github.com/hashicorp/vault/pull/20668)]
* secrets/transit: Fix export of HMAC-only key, correctly exporting the key used for sign operations. For consumers of the previously incorrect key, use the plaintext export to retrieve these incorrect keys and import them as new versions.
secrets/transit: Fix bug related to shorter dedicated HMAC key sizing.
sdk/helper/keysutil: New HMAC type policies will have HMACKey equal to Key and be copied over on import. [[GH-20864](https://github.com/hashicorp/vault/pull/20864)]
* ui: Fixes issue unsealing cluster for seal types other than shamir [[GH-20897](https://github.com/hashicorp/vault/pull/20897)]
* ui: fixes issue creating mfa login enforcement from method enforcements tab [[GH-20603](https://github.com/hashicorp/vault/pull/20603)]
* ui: fixes key_bits and signature_bits reverting to default values when editing a pki role [[GH-20907](https://github.com/hashicorp/vault/pull/20907)]

## 1.13.2
### April 26, 2023

CHANGES:

* core: Bump Go version to 1.20.3.

SECURITY:

* core/seal: Fix handling of HMACing of seal-wrapped storage entries from HSMs using CKM_AES_CBC or CKM_AES_CBC_PAD which may have allowed an attacker to conduct a padding oracle attack. This vulnerability, CVE-2023-2197, affects Vault from 1.13.0 up to 1.13.1 and was fixed in 1.13.2. [[HCSEC-2023-14](https://discuss.hashicorp.com/t/hcsec-2023-14-vault-enterprise-vulnerable-to-padding-oracle-attacks-when-using-a-cbc-based-encryption-mechanism-with-a-hsm/53322)]

IMPROVEMENTS:

* Add debug symbols back to builds to fix Dynatrace support [[GH-20294](https://github.com/hashicorp/vault/pull/20294)]
* cli/namespace: Add detailed flag to output additional namespace information
such as namespace IDs and custom metadata. [[GH-20243](https://github.com/hashicorp/vault/pull/20243)]
* core/activity: add an endpoint to write test activity log data, guarded by a build flag [[GH-20019](https://github.com/hashicorp/vault/pull/20019)]
* core: Add a `raft` sub-field to the `storage` and `ha_storage` details provided by the
`/sys/config/state/sanitized` endpoint in order to include the `max_entry_size`. [[GH-20044](https://github.com/hashicorp/vault/pull/20044)]
* core: include reason for ErrReadOnly on PBPWF writing failures
* sdk/ldaputil: added `connection_timeout` to tune connection timeout duration
for all LDAP plugins. [[GH-20144](https://github.com/hashicorp/vault/pull/20144)]
* secrets/pki: Decrease size and improve compatibility of OCSP responses by removing issuer certificate. [[GH-20201](https://github.com/hashicorp/vault/pull/20201)]
* sys/wrapping: Add example how to unwrap without authentication in Vault [[GH-20109](https://github.com/hashicorp/vault/pull/20109)]
* ui: Allows license-banners to be dismissed. Saves preferences in localStorage. [[GH-19116](https://github.com/hashicorp/vault/pull/19116)]

BUG FIXES:

* auth/ldap: Add max_page_size configurable to LDAP configuration [[GH-19032](https://github.com/hashicorp/vault/pull/19032)]
* command/server: Fix incorrect paths in generated config for `-dev-tls` flag on Windows [[GH-20257](https://github.com/hashicorp/vault/pull/20257)]
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Fix regression breaking non-raft clusters whose nodes share the same cluster_addr/api_addr. [[GH-19721](https://github.com/hashicorp/vault/pull/19721)]
* helper/random: Fix race condition in string generator helper [[GH-19875](https://github.com/hashicorp/vault/pull/19875)]
* kmip (enterprise): Fix a problem decrypting with keys that have no Process Start Date attribute.
* pki: Fix automatically turning off CRL signing on upgrade to Vault >= 1.12, if CA Key Usage disallows it [[GH-20220](https://github.com/hashicorp/vault/pull/20220)]
* replication (enterprise): Fix a caching issue when replicating filtered data to
a performance secondary. This resulted in the data being set to nil in the cache
and a "invalid value" error being returned from the API.
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* sdk/helper/ocsp: Workaround bug in Go's ocsp.ParseResponse(...), causing validation to fail with embedded CA certificates.
auth/cert: Fix OCSP validation against Vault's PKI engine. [[GH-20181](https://github.com/hashicorp/vault/pull/20181)]
* secrets/aws: Revert changes that removed the lease on STS credentials, while leaving the new ttl field in place. [[GH-20034](https://github.com/hashicorp/vault/pull/20034)]
* secrets/pki: Ensure cross-cluster delta WAL write failure only logs to avoid unattended forwarding. [[GH-20057](https://github.com/hashicorp/vault/pull/20057)]
* secrets/pki: Fix building of unified delta CRLs and recovery during unified delta WAL write failures. [[GH-20058](https://github.com/hashicorp/vault/pull/20058)]
* secrets/pki: Fix patching of leaf_not_after_behavior on issuers. [[GH-20341](https://github.com/hashicorp/vault/pull/20341)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* ui: Fix OIDC provider logo showing when domain doesn't match [[GH-20263](https://github.com/hashicorp/vault/pull/20263)]
* ui: Fix bad link to namespace when namespace name includes `.` [[GH-19799](https://github.com/hashicorp/vault/pull/19799)]
* ui: fixes browser console formatting for help command output [[GH-20064](https://github.com/hashicorp/vault/pull/20064)]
* ui: fixes remaining doc links to include /vault in path [[GH-20070](https://github.com/hashicorp/vault/pull/20070)]
* ui: remove use of htmlSafe except when first sanitized [[GH-20235](https://github.com/hashicorp/vault/pull/20235)]
* website/docs: Fix Kubernetes Auth Code Example to use the correct whitespace in import. [[GH-20216](https://github.com/hashicorp/vault/pull/20216)]

## 1.13.1
### March 29, 2023

SECURITY:

* storage/mssql: When using Vault’s community-supported Microsoft SQL (MSSQL) database storage backend, a privileged attacker with the ability to write arbitrary data to Vault’s configuration may be able to perform arbitrary SQL commands on the underlying database server through Vault. This vulnerability, CVE-2023-0620, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-12](https://discuss.hashicorp.com/t/hcsec-2023-12-vault-s-microsoft-sql-database-storage-backend-vulnerable-to-sql-injection-via-configuration-file/52080)]
* secrets/pki: Vault’s PKI mount issuer endpoints did not correctly authorize access to remove an issuer or modify issuer metadata, potentially resulting in denial of service of the PKI mount. This bug did not affect public or private key material, trust chains or certificate issuance. This vulnerability, CVE-2023-0665, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-11](https://discuss.hashicorp.com/t/hcsec-2023-11-vault-s-pki-issuer-endpoint-did-not-correctly-authorize-access-to-issuer-metadata/52079)]
* core: HashiCorp Vault’s implementation of Shamir’s secret sharing used precomputed table lookups, and was vulnerable to cache-timing attacks. An attacker with access to, and the ability to observe a large number of unseal operations on the host through a side channel may reduce the search space of a brute force effort to recover the Shamir shares. This vulnerability, CVE-2023-25000, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-10](https://discuss.hashicorp.com/t/hcsec-2023-10-vault-vulnerable-to-cache-timing-attacks-during-seal-and-unseal-operations/52078)]

IMPROVEMENTS:

* auth/github: Allow for an optional Github auth token environment variable to make authenticated requests when fetching org id
website/docs: Add docs for `VAULT_AUTH_CONFIG_GITHUB_TOKEN` environment variable when writing Github config [[GH-19244](https://github.com/hashicorp/vault/pull/19244)]
* core: Allow overriding gRPC connect timeout via VAULT_GRPC_MIN_CONNECT_TIMEOUT. This is an env var rather than a config setting because we don't expect this to ever be needed.  It's being added as a last-ditch
option in case all else fails for some replication issues we may not have fully reproduced. [[GH-19676](https://github.com/hashicorp/vault/pull/19676)]
* core: validate name identifiers in mssql physical storage backend prior use [[GH-19591](https://github.com/hashicorp/vault/pull/19591)]
* database/elasticsearch: Update error messages resulting from Elasticsearch API errors [[GH-19545](https://github.com/hashicorp/vault/pull/19545)]
* events: Suppress log warnings triggered when events are sent but the events system is not enabled. [[GH-19593](https://github.com/hashicorp/vault/pull/19593)]

BUG FIXES:

* agent: Fix panic when SIGHUP is issued to Agent while it has a non-TLS listener. [[GH-19483](https://github.com/hashicorp/vault/pull/19483)]
* core (enterprise): Attempt to reconnect to a PKCS#11 HSM if we retrieve a CKR_FUNCTION_FAILED error.
* core: Fixed issue with remounting mounts that have a non-trailing space in the 'to' or 'from' paths. [[GH-19585](https://github.com/hashicorp/vault/pull/19585)]
* kmip (enterprise): Do not require attribute Cryptographic Usage Mask when registering Secret Data managed objects.
* kmip (enterprise): Fix a problem forwarding some requests to the active node.
* openapi: Fix logic for labeling unauthenticated/sudo paths. [[GH-19600](https://github.com/hashicorp/vault/pull/19600)]
* secrets/ldap: Invalidates WAL entry for static role if `password_policy` has changed. [[GH-19640](https://github.com/hashicorp/vault/pull/19640)]
* secrets/pki: Fix PKI revocation request forwarding from standby nodes due to an error wrapping bug [[GH-19624](https://github.com/hashicorp/vault/pull/19624)]
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* ui: Fixes crypto.randomUUID error in unsecure contexts from third party ember-data library [[GH-19428](https://github.com/hashicorp/vault/pull/19428)]
* ui: fixes SSH engine config deletion [[GH-19448](https://github.com/hashicorp/vault/pull/19448)]
* ui: fixes issue navigating back a level using the breadcrumb from secret metadata view [[GH-19703](https://github.com/hashicorp/vault/pull/19703)]
* ui: fixes oidc tabs in auth form submitting with the root's default_role value after a namespace has been inputted [[GH-19541](https://github.com/hashicorp/vault/pull/19541)]
* ui: pass encodeBase64 param to HMAC transit-key-actions. [[GH-19429](https://github.com/hashicorp/vault/pull/19429)]
* ui: use URLSearchParams interface to capture namespace param from SSOs (ex. ADFS) with decoded state param in callback url [[GH-19460](https://github.com/hashicorp/vault/pull/19460)]

## 1.13.0
### March 01, 2023

SECURITY:

* secrets/ssh: removal of the deprecated dynamic keys mode. **When any remaining dynamic key leases expire**, an error stating `secret is unsupported by this backend` will be thrown by the lease manager. [[GH-18874](https://github.com/hashicorp/vault/pull/18874)]
* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* auth/alicloud: require the `role` field on login [[GH-19005](https://github.com/hashicorp/vault/pull/19005)]
* auth/approle: Add maximum length of 4096 for approle role_names, as this value results in HMAC calculation [[GH-17768](https://github.com/hashicorp/vault/pull/17768)]
* auth: Returns invalid credentials for ldap, userpass and approle when wrong credentials are provided for existent users.
This will only be used internally for implementing user lockout. [[GH-17104](https://github.com/hashicorp/vault/pull/17104)]
* core: Bump Go version to 1.20.1.
* core: Vault version has been moved out of sdk and into main vault module.
Plugins using sdk/useragent.String must instead use sdk/useragent.PluginString. [[GH-14229](https://github.com/hashicorp/vault/pull/14229)]
* logging: Removed legacy environment variable for log format ('LOGXI_FORMAT'), should use 'VAULT_LOG_FORMAT' instead [[GH-17822](https://github.com/hashicorp/vault/pull/17822)]
* plugins: Mounts can no longer be pinned to a specific _builtin_ version. Mounts previously pinned to a specific builtin version will now automatically upgrade to the latest builtin version, and may now be overridden if an unversioned plugin of the same name and type is registered. Mounts using plugin versions without `builtin` in their metadata remain unaffected. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]
* plugins: `GET /database/config/:name` endpoint now returns an additional `plugin_version` field in the response data. [[GH-16982](https://github.com/hashicorp/vault/pull/16982)]
* plugins: `GET /sys/auth/:path/tune` and `GET /sys/mounts/:path/tune` endpoints may now return an additional `plugin_version` field in the response data if set. [[GH-17167](https://github.com/hashicorp/vault/pull/17167)]
* plugins: `GET` for `/sys/auth`, `/sys/auth/:path`, `/sys/mounts`, and `/sys/mounts/:path` paths now return additional `plugin_version`, `running_plugin_version` and `running_sha256` fields in the response data for each mount. [[GH-17167](https://github.com/hashicorp/vault/pull/17167)]
* sdk: Remove version package, make useragent.String versionless. [[GH-19068](https://github.com/hashicorp/vault/pull/19068)]
* secrets/aws: do not create leases for non-renewable/non-revocable STS credentials to reduce storage calls [[GH-15869](https://github.com/hashicorp/vault/pull/15869)]
* secrets/gcpkms: Updated plugin from v0.13.0 to v0.14.0 [[GH-19063](https://github.com/hashicorp/vault/pull/19063)]
* sys/internal/inspect: Turns of this endpoint by default. A SIGHUP can now be used to reload the configs and turns this endpoint on.
* ui: Upgrade Ember to version 4.4.0 [[GH-17086](https://github.com/hashicorp/vault/pull/17086)]

FEATURES:

* **User lockout**: Ignore repeated bad credentials from the same user for a configured period of time. Enabled by default.
* **Azure Auth Managed Identities**: Allow any Azure resource that supports managed identities to authenticate with Vault [[GH-19077](https://github.com/hashicorp/vault/pull/19077)]
* **Azure Auth Rotate Root**: Add support for rotate root in Azure Auth engine [[GH-19077](https://github.com/hashicorp/vault/pull/19077)]
* **Event System (Alpha)**: Vault has a new opt-in experimental event system. Not yet suitable for production use. Events are currently only generated on writes to the KV secrets engine, but external plugins can also be updated to start generating events. [[GH-19194](https://github.com/hashicorp/vault/pull/19194)]
* **GCP Secrets Impersonated Account Support**: Add support for GCP service account impersonation, allowing callers to generate a GCP access token without requiring Vault to store or retrieve a GCP service account key for each role. [[GH-19018](https://github.com/hashicorp/vault/pull/19018)]
* **Kubernetes Secrets Engine UI**: Kubernetes is now available in the UI as a supported secrets engine. [[GH-17893](https://github.com/hashicorp/vault/pull/17893)]
* **New PKI UI**: Add beta support for new and improved PKI UI [[GH-18842](https://github.com/hashicorp/vault/pull/18842)]
* **PKI Cross-Cluster Revocations**: Revocation information can now be
synchronized across primary and performance replica clusters offering
a unified CRL/OCSP view of revocations across cluster boundaries. [[GH-19196](https://github.com/hashicorp/vault/pull/19196)]
* **Server UDS Listener**: Adding listener to Vault server to serve http request via unix domain socket [[GH-18227](https://github.com/hashicorp/vault/pull/18227)]
* **Transit managed keys**: The transit secrets engine now supports configuring and using managed keys
* **User Lockout**: Adds support to configure the user-lockout behaviour for failed logins to prevent
brute force attacks for userpass, approle and ldap auth methods. [[GH-19230](https://github.com/hashicorp/vault/pull/19230)]
* **VMSS Flex Authentication**: Adds support for Virtual Machine Scale Set Flex Authentication [[GH-19077](https://github.com/hashicorp/vault/pull/19077)]
* **Namespaces (enterprise)**: Added the ability to allow access to secrets and more to be shared across namespaces that do not share a namespace hierarchy. Using the new `sys/config/group-policy-application` API, policies can be configured to apply outside of namespace hierarchy, allowing this kind of cross-namespace sharing.
* **OpenAPI-based Go & .NET Client Libraries (Beta)**: We have now made available two new [[OpenAPI-based Go](https://github.com/hashicorp/vault-client-go/)] & [[OpenAPI-based .NET](https://github.com/hashicorp/vault-client-dotnet/)] Client libraries (beta). You can use them to perform various secret management operations easily from your applications.

IMPROVEMENTS:

* **Redis ElastiCache DB Engine**: Renamed configuration parameters for disambiguation; old parameters still supported for compatibility. [[GH-18752](https://github.com/hashicorp/vault/pull/18752)]
* Bump github.com/hashicorp/go-plugin version from 1.4.5 to 1.4.8 [[GH-19100](https://github.com/hashicorp/vault/pull/19100)]
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
* auth/cf: Remove incorrect usage of CreateOperation from path_config [[GH-19098](https://github.com/hashicorp/vault/pull/19098)]
* auth/gcp: Upgrades dependencies [[GH-17858](https://github.com/hashicorp/vault/pull/17858)]
* auth/oidc: Adds `abort_on_error` parameter to CLI login command to help in non-interactive contexts [[GH-19076](https://github.com/hashicorp/vault/pull/19076)]
* auth/oidc: Adds ability to set Google Workspace domain for groups search [[GH-19076](https://github.com/hashicorp/vault/pull/19076)]
* auth/token (enterprise): Allow batch token creation in perfStandby nodes
* auth: Allow naming login MFA methods and using those names instead of IDs in satisfying MFA requirement for requests.
Make passcode arguments consistent across login MFA method types. [[GH-18610](https://github.com/hashicorp/vault/pull/18610)]
* auth: Provide an IP address of the requests from Vault to a Duo challenge after successful authentication. [[GH-18811](https://github.com/hashicorp/vault/pull/18811)]
* autopilot: Update version to v.0.2.0 to add better support for respecting min quorum
* cli/kv: improve kv CLI to remove data or custom metadata using kv patch [[GH-18067](https://github.com/hashicorp/vault/pull/18067)]
* cli/pki: Add List-Intermediates functionality to pki client. [[GH-18463](https://github.com/hashicorp/vault/pull/18463)]
* cli/pki: Add health-check subcommand to evaluate the health of a PKI instance. [[GH-17750](https://github.com/hashicorp/vault/pull/17750)]
* cli/pki: Add pki issue command, which creates a CSR, has a vault mount sign it, then reimports it. [[GH-18467](https://github.com/hashicorp/vault/pull/18467)]
* cli/pki: Added "Reissue" command which allows extracting fields from an existing certificate to create a new certificate. [[GH-18499](https://github.com/hashicorp/vault/pull/18499)]
* cli/pki: Change the pki health-check --list default config output to JSON so it's a usable configuration file [[GH-19269](https://github.com/hashicorp/vault/pull/19269)]
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
* database/redis-elasticache: changed config argument names for disambiguation [[GH-19044](https://github.com/hashicorp/vault/pull/19044)]
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
* secrets/azure: Adds ability to persist an application for the lifetime of a role. [[GH-19096](https://github.com/hashicorp/vault/pull/19096)]
* secrets/azure: upgrades dependencies [[GH-17964](https://github.com/hashicorp/vault/pull/17964)]
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* secrets/gcp: Upgrades dependencies [[GH-17871](https://github.com/hashicorp/vault/pull/17871)]
* secrets/kubernetes: Add /check endpoint to determine if environment variables are set [[GH-18](https://github.com/hashicorp/vault-plugin-secrets-kubernetes/pull/18)] [[GH-18587](https://github.com/hashicorp/vault/pull/18587)]
* secrets/kubernetes: add /check endpoint to determine if environment variables are set [[GH-19084](https://github.com/hashicorp/vault/pull/19084)]
* secrets/kv: Emit events on write if events system enabled [[GH-19145](https://github.com/hashicorp/vault/pull/19145)]
* secrets/kv: make upgrade synchronous when no keys to upgrade [[GH-19056](https://github.com/hashicorp/vault/pull/19056)]
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
* ui: adds allowed_response_headers as param for secret engine mount config [[GH-19216](https://github.com/hashicorp/vault/pull/19216)]
* ui: consolidate all <a> tag usage [[GH-17866](https://github.com/hashicorp/vault/pull/17866)]
* ui: mfa: use proper request id generation [[GH-17835](https://github.com/hashicorp/vault/pull/17835)]
* ui: remove wizard [[GH-19220](https://github.com/hashicorp/vault/pull/19220)]
* ui: update DocLink component to use new host url: developer.hashicorp.com [[GH-18374](https://github.com/hashicorp/vault/pull/18374)]
* ui: update TTL picker for consistency [[GH-18114](https://github.com/hashicorp/vault/pull/18114)]
* ui: use the combined activity log (partial + historic) API for client count dashboard and remove use of monthly endpoint [[GH-17575](https://github.com/hashicorp/vault/pull/17575)]
* vault/diagnose: Upgrade `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk`, `go.opentelemetry.io/otel/trace` to v1.11.2 [[GH-18589](https://github.com/hashicorp/vault/pull/18589)]

DEPRECATIONS:

* secrets/ad: Marks the Active Directory (AD) secrets engine as deprecated. [[GH-19334](https://github.com/hashicorp/vault/pull/19334)]

BUG FIXES:

* api: Remove timeout logic from ReadRaw functions and add ReadRawWithContext [[GH-18708](https://github.com/hashicorp/vault/pull/18708)]
* auth/alicloud: fix regression in vault login command that caused login to fail [[GH-19005](https://github.com/hashicorp/vault/pull/19005)]
* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* auth/cert: Address a race condition accessing the loaded crls without a lock [[GH-18945](https://github.com/hashicorp/vault/pull/18945)]
* auth/kubernetes: Ensure a consistent TLS configuration for all k8s API requests [[#173](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/173)] [[GH-18716](https://github.com/hashicorp/vault/pull/18716)]
* auth/kubernetes: fixes and dep updates for the auth-kubernetes plugin (see plugin changelog for details) [[GH-19094](https://github.com/hashicorp/vault/pull/19094)]
* auth/okta: fix a panic for AuthRenew in Okta [[GH-18011](https://github.com/hashicorp/vault/pull/18011)]
* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* cli/pki: Decode integer values properly in health-check configuration file [[GH-19265](https://github.com/hashicorp/vault/pull/19265)]
* cli/pki: Fix path for role health-check warning messages [[GH-19274](https://github.com/hashicorp/vault/pull/19274)]
* cli/pki: Properly report permission issues within health-check mount tune checks [[GH-19276](https://github.com/hashicorp/vault/pull/19276)]
* cli/transit: Fix import, import-version command invocation [[GH-19373](https://github.com/hashicorp/vault/pull/19373)]
* cli: Fix issue preventing kv commands from executing properly when the mount path provided by `-mount` flag and secret key path are the same. [[GH-17679](https://github.com/hashicorp/vault/pull/17679)]
* cli: Fix vault read handling to return raw data as secret.Data when there is no top-level data object from api response. [[GH-17913](https://github.com/hashicorp/vault/pull/17913)]
* cli: Remove empty table heading for `vault secrets list -detailed` output. [[GH-17577](https://github.com/hashicorp/vault/pull/17577)]
* command/namespace: Fix vault cli namespace patch examples in help text. [[GH-18143](https://github.com/hashicorp/vault/pull/18143)]
* core (enterprise): Fix missing quotation mark in error message
* core (enterprise): Fix panic that could occur with SSCT alongside invoking external plugins for revocation.
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
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
* core: Linux packages now have vendor label and set the default label to HashiCorp.
This fix is implemented for any future releases, but will not be updated for historical releases.
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
* kmip (enterprise): Fix a problem with some multi-part MAC Verify operations.
* kmip (enterprise): Only require data to be full blocks on encrypt/decrypt operations using CBC and ECB block cipher modes.
* license (enterprise): Fix bug where license would update even if the license didn't change.
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
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* sdk: Don't panic if system view or storage methods called during plugin setup. [[GH-18210](https://github.com/hashicorp/vault/pull/18210)]
* secret/pki: fix bug with initial legacy bundle migration (from < 1.11 into 1.11+) and missing issuers from ca_chain [[GH-17772](https://github.com/hashicorp/vault/pull/17772)]
* secrets/ad: Fix bug where updates to config would fail if password isn't provided [[GH-19061](https://github.com/hashicorp/vault/pull/19061)]
* secrets/gcp: fix issue where IAM bindings were not preserved during policy update [[GH-19018](https://github.com/hashicorp/vault/pull/19018)]
* secrets/mongodb-atlas: Fix a bug that did not allow WAL rollback to handle partial failures when creating API keys [[GH-19111](https://github.com/hashicorp/vault/pull/19111)]
* secrets/pki: Address nil panic when an empty POST request is sent to the OCSP handler [[GH-18184](https://github.com/hashicorp/vault/pull/18184)]
* secrets/pki: Allow patching issuer to set an empty issuer name. [[GH-18466](https://github.com/hashicorp/vault/pull/18466)]
* secrets/pki: Do not read revoked certificates from backend when CRL is disabled [[GH-17385](https://github.com/hashicorp/vault/pull/17385)]
* secrets/pki: Fix upgrade of missing expiry, delta_rebuild_interval by setting them to the default. [[GH-17693](https://github.com/hashicorp/vault/pull/17693)]
* secrets/pki: Fixes duplicate otherName in certificates created by the sign-verbatim endpoint. [[GH-16700](https://github.com/hashicorp/vault/pull/16700)]
* secrets/pki: OCSP GET request parameter was not being URL unescaped before processing. [[GH-18938](https://github.com/hashicorp/vault/pull/18938)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* secrets/pki: Revert fix for PR [18938](https://github.com/hashicorp/vault/pull/18938) [[GH-19037](https://github.com/hashicorp/vault/pull/19037)]
* secrets/pki: consistently use UTC for CA's notAfter exceeded error message [[GH-18984](https://github.com/hashicorp/vault/pull/18984)]
* secrets/pki: fix race between tidy's cert counting and tidy status reporting. [[GH-18899](https://github.com/hashicorp/vault/pull/18899)]
* secrets/transit: Do not warn about unrecognized parameter 'batch_input' [[GH-18299](https://github.com/hashicorp/vault/pull/18299)]
* secrets/transit: Honor `partial_success_response_code` on decryption failures. [[GH-18310](https://github.com/hashicorp/vault/pull/18310)]
* server/config:  Use file.Stat when checking file permissions when VAULT_ENABLE_FILE_PERMISSIONS_CHECK is enabled [[GH-19311](https://github.com/hashicorp/vault/pull/19311)]
* storage/raft (enterprise): An already joined node can rejoin by wiping storage
and re-issueing a join request, but in doing so could transiently become a
non-voter.  In some scenarios this resulted in loss of quorum. [[GH-18263](https://github.com/hashicorp/vault/pull/18263)]
* storage/raft: Don't panic on unknown raft ops [[GH-17732](https://github.com/hashicorp/vault/pull/17732)]
* storage/raft: Fix race with follower heartbeat tracker during teardown. [[GH-18704](https://github.com/hashicorp/vault/pull/18704)]
* ui/keymgmt: Sets the defaultValue for type when creating a key. [[GH-17407](https://github.com/hashicorp/vault/pull/17407)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]
* ui: Fixes issue with not being able to download raft snapshot via service worker [[GH-17769](https://github.com/hashicorp/vault/pull/17769)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]
* ui: Remove `default` and add `default-service` and `default-batch` to UI token_type for auth mount and tuning. [[GH-19290](https://github.com/hashicorp/vault/pull/19290)]
* ui: Remove default value of 30 to TtlPicker2 if no value is passed in. [[GH-17376](https://github.com/hashicorp/vault/pull/17376)]
* ui: allow selection of "default" for ssh algorithm_signer in web interface [[GH-17894](https://github.com/hashicorp/vault/pull/17894)]
* ui: cleanup unsaved auth method ember data record when navigating away from mount backend form [[GH-18651](https://github.com/hashicorp/vault/pull/18651)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]
* ui: fixes reliance on secure context (https) by removing methods using the Crypto interface [[GH-19403](https://github.com/hashicorp/vault/pull/19403)]
* ui: show Get credentials button for static roles detail page when a user has the proper permissions. [[GH-19190](https://github.com/hashicorp/vault/pull/19190)]

## 1.12.9
### July 25, 2023

CHANGES:

* secrets/transform (enterprise): Enforce a transformation role's max_ttl setting on encode requests, a warning will be returned if max_ttl was applied.

IMPROVEMENTS:

* core/fips: Add RPM, DEB packages of FIPS 140-2 and HSM+FIPS 140-2 Vault Enterprise.
* replication (enterprise): Avoid logging warning if request is forwarded from a performance standby and not a performance secondary
* secrets/transform (enterprise): Switch to pgx PostgreSQL driver for better timeout handling

BUG FIXES:

* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-24170](https://github.com/hashicorp/vault/pull/24170)]
* identity: Remove caseSensitivityKey to prevent errors while loading groups which could result in missing groups in memDB when duplicates are found. [[GH-20965](https://github.com/hashicorp/vault/pull/20965)]
* replication (enterprise): update primary cluster address after DR failover
* secrets/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21633](https://github.com/hashicorp/vault/pull/21633)]
* secrets/pki: Prevent deleted issuers from reappearing when migrating from a version 1 bundle to a version 2 bundle (versions including 1.13.0, 1.12.2, and 1.11.6); when managed keys were removed but referenced in the Vault 1.10 legacy CA bundle, this the error: `no managed key found with uuid`. [[GH-21316](https://github.com/hashicorp/vault/pull/21316)]
* secrets/pki: Support setting both maintain_stored_certificate_counts=false and publish_stored_certificate_count_metrics=false explicitly in tidy config. [[GH-20664](https://github.com/hashicorp/vault/pull/20664)]
* secrets/transform (enterprise): Fix nil panic when deleting a template with tokenization transformations present
* secrets/transform (enterprise): Grab shared locks for various read operations, only escalating to write locks if work is required
* serviceregistration: Fix bug where multiple nodes in a secondary cluster could be labelled active after updating the cluster's primary [[GH-21642](https://github.com/hashicorp/vault/pull/21642)]
* ui: Fixed an issue where editing an SSH role would clear `default_critical_options` and `default_extension` if left unchanged. [[GH-21739](https://github.com/hashicorp/vault/pull/21739)]

## 1.12.8
### June 21, 2023
BREAKING CHANGES:

* secrets/pki: Maintaining running count of certificates will be turned off by default.
To re-enable keeping these metrics available on the tidy status endpoint, enable
maintain_stored_certificate_counts on tidy-config, to also publish them to the
metrics consumer, enable publish_stored_certificate_count_metrics . [[GH-18186](https://github.com/hashicorp/vault/pull/18186)]

CHANGES:

* core: Bump Go version to 1.19.10.

FEATURES:

* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* core (enterprise): Add background worker for automatic reporting of billing
information. [[GH-19625](https://github.com/hashicorp/vault/pull/19625)]

IMPROVEMENTS:

* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]

BUG FIXES:

* core (enterprise): Don't delete backend stored data that appears to be filterable
on this secondary if we don't have a corresponding mount entry.
* core/activity: add namespace breakdown for new clients when date range spans multiple months, including the current month. [[GH-18766](https://github.com/hashicorp/vault/pull/18766)]
* core/activity: de-duplicate namespaces when historical and current month data are mixed [[GH-18452](https://github.com/hashicorp/vault/pull/18452)]
* core/activity: fix the end_date returned from the activity log endpoint when partial counts are computed [[GH-17856](https://github.com/hashicorp/vault/pull/17856)]
* core/activity: include mount counts when de-duplicating current and historical month data [[GH-18598](https://github.com/hashicorp/vault/pull/18598)]
* core/activity: report mount paths (rather than mount accessors) in current month activity log counts and include deleted mount paths in precomputed queries. [[GH-18916](https://github.com/hashicorp/vault/pull/18916)]
* core/activity: return partial month counts when querying a historical date range and no historical data exists. [[GH-17935](https://github.com/hashicorp/vault/pull/17935)]
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry.  Fix wrong lock used in ListAuths link meta interface implementation. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs
* storage/raft: Fix race where new follower joining can get pruned by dead server cleanup. [[GH-20986](https://github.com/hashicorp/vault/pull/20986)]

## 1.12.7
### June 08, 2023

CHANGES:

* core: Bump Go version to 1.19.9.
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]

IMPROVEMENTS:

* audit: add a `mount_point` field to audit requests and response entries [[GH-20411](https://github.com/hashicorp/vault/pull/20411)]
* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* core: include namespace path in granting_policies block of audit log
* openapi: Fix generated types for duration strings [[GH-20841](https://github.com/hashicorp/vault/pull/20841)]
* sdk/framework: Fix non-deterministic ordering of 'required' fields in OpenAPI spec [[GH-20881](https://github.com/hashicorp/vault/pull/20881)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]
* ui: update TTL picker for consistency [[GH-18114](https://github.com/hashicorp/vault/pull/18114)]

BUG FIXES:

* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* cli: CLI should take days as a unit of time for ttl like flags [[GH-20477](https://github.com/hashicorp/vault/pull/20477)]
* cli: disable printing flags warnings messages for the ssh command [[GH-20502](https://github.com/hashicorp/vault/pull/20502)]
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core: prevent panic on login after namespace is deleted that had mfa enforcement [[GH-20375](https://github.com/hashicorp/vault/pull/20375)]
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation
* secrets/transit: Fix export of HMAC-only key, correctly exporting the key used for sign operations. For consumers of the previously incorrect key, use the plaintext export to retrieve these incorrect keys and import them as new versions.
secrets/transit: Fix bug related to shorter dedicated HMAC key sizing.
sdk/helper/keysutil: New HMAC type policies will have HMACKey equal to Key and be copied over on import. [[GH-20864](https://github.com/hashicorp/vault/pull/20864)]
* ui: Fixes issue unsealing cluster for seal types other than shamir [[GH-20897](https://github.com/hashicorp/vault/pull/20897)]

## 1.12.6
### April 26, 2023

CHANGES:

* core: Bump Go version to 1.19.8.

IMPROVEMENTS:

* cli/namespace: Add detailed flag to output additional namespace information
such as namespace IDs and custom metadata. [[GH-20243](https://github.com/hashicorp/vault/pull/20243)]
* core/activity: add an endpoint to write test activity log data, guarded by a build flag [[GH-20019](https://github.com/hashicorp/vault/pull/20019)]
* core: Add a `raft` sub-field to the `storage` and `ha_storage` details provided by the
`/sys/config/state/sanitized` endpoint in order to include the `max_entry_size`. [[GH-20044](https://github.com/hashicorp/vault/pull/20044)]
* sdk/ldaputil: added `connection_timeout` to tune connection timeout duration
for all LDAP plugins. [[GH-20144](https://github.com/hashicorp/vault/pull/20144)]
* secrets/pki: Decrease size and improve compatibility of OCSP responses by removing issuer certificate. [[GH-20201](https://github.com/hashicorp/vault/pull/20201)]

BUG FIXES:

* auth/ldap: Add max_page_size configurable to LDAP configuration [[GH-19032](https://github.com/hashicorp/vault/pull/19032)]
* command/server: Fix incorrect paths in generated config for `-dev-tls` flag on Windows [[GH-20257](https://github.com/hashicorp/vault/pull/20257)]
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Fix regression breaking non-raft clusters whose nodes share the same cluster_addr/api_addr. [[GH-19721](https://github.com/hashicorp/vault/pull/19721)]
* helper/random: Fix race condition in string generator helper [[GH-19875](https://github.com/hashicorp/vault/pull/19875)]
* kmip (enterprise): Fix a problem decrypting with keys that have no Process Start Date attribute.
* openapi: Fix many incorrect details in generated API spec, by using better techniques to parse path regexps [[GH-18554](https://github.com/hashicorp/vault/pull/18554)]
* pki: Fix automatically turning off CRL signing on upgrade to Vault >= 1.12, if CA Key Usage disallows it [[GH-20220](https://github.com/hashicorp/vault/pull/20220)]
* replication (enterprise): Fix a caching issue when replicating filtered data to
a performance secondary. This resulted in the data being set to nil in the cache
and a "invalid value" error being returned from the API.
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* secrets/pki: Fix patching of leaf_not_after_behavior on issuers. [[GH-20341](https://github.com/hashicorp/vault/pull/20341)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* ui: Fix OIDC provider logo showing when domain doesn't match [[GH-20263](https://github.com/hashicorp/vault/pull/20263)]
* ui: Fix bad link to namespace when namespace name includes `.` [[GH-19799](https://github.com/hashicorp/vault/pull/19799)]
* ui: fixes browser console formatting for help command output [[GH-20064](https://github.com/hashicorp/vault/pull/20064)]
* ui: remove use of htmlSafe except when first sanitized [[GH-20235](https://github.com/hashicorp/vault/pull/20235)]

## 1.12.5
### March 29, 2023

SECURITY:

* storage/mssql: When using Vault’s community-supported Microsoft SQL (MSSQL) database storage backend, a privileged attacker with the ability to write arbitrary data to Vault’s configuration may be able to perform arbitrary SQL commands on the underlying database server through Vault. This vulnerability, CVE-2023-0620, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-12](https://discuss.hashicorp.com/t/hcsec-2023-12-vault-s-microsoft-sql-database-storage-backend-vulnerable-to-sql-injection-via-configuration-file/52080)]
* secrets/pki: Vault’s PKI mount issuer endpoints did not correctly authorize access to remove an issuer or modify issuer metadata, potentially resulting in denial of service of the PKI mount. This bug did not affect public or private key material, trust chains or certificate issuance. This vulnerability, CVE-2023-0665, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-11](https://discuss.hashicorp.com/t/hcsec-2023-11-vault-s-pki-issuer-endpoint-did-not-correctly-authorize-access-to-issuer-metadata/52079)]
* core: HashiCorp Vault’s implementation of Shamir’s secret sharing used precomputed table lookups, and was vulnerable to cache-timing attacks. An attacker with access to, and the ability to observe a large number of unseal operations on the host through a side channel may reduce the search space of a brute force effort to recover the Shamir shares. This vulnerability, CVE-2023-25000, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-10](https://discuss.hashicorp.com/t/hcsec-2023-10-vault-vulnerable-to-cache-timing-attacks-during-seal-and-unseal-operations/52078)]

IMPROVEMENTS:

* auth/github: Allow for an optional Github auth token environment variable to make authenticated requests when fetching org id
website/docs: Add docs for `VAULT_AUTH_CONFIG_GITHUB_TOKEN` environment variable when writing Github config [[GH-19244](https://github.com/hashicorp/vault/pull/19244)]
* core: Allow overriding gRPC connect timeout via VAULT_GRPC_MIN_CONNECT_TIMEOUT. This is an env var rather than a config setting because we don't expect this to ever be needed.  It's being added as a last-ditch
option in case all else fails for some replication issues we may not have fully reproduced. [[GH-19676](https://github.com/hashicorp/vault/pull/19676)]
* core: validate name identifiers in mssql physical storage backend prior use [[GH-19591](https://github.com/hashicorp/vault/pull/19591)]

BUG FIXES:

* cli: Fix vault read handling to return raw data as secret.Data when there is no top-level data object from api response. [[GH-17913](https://github.com/hashicorp/vault/pull/17913)]
* core (enterprise): Attempt to reconnect to a PKCS#11 HSM if we retrieve a CKR_FUNCTION_FAILED error.
* core: Fixed issue with remounting mounts that have a non-trailing space in the 'to' or 'from' paths. [[GH-19585](https://github.com/hashicorp/vault/pull/19585)]
* kmip (enterprise): Do not require attribute Cryptographic Usage Mask when registering Secret Data managed objects.
* kmip (enterprise): Fix a problem forwarding some requests to the active node.
* openapi: Fix logic for labeling unauthenticated/sudo paths. [[GH-19600](https://github.com/hashicorp/vault/pull/19600)]
* secrets/ldap: Invalidates WAL entry for static role if `password_policy` has changed. [[GH-19641](https://github.com/hashicorp/vault/pull/19641)]
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* ui: fixes issue navigating back a level using the breadcrumb from secret metadata view [[GH-19703](https://github.com/hashicorp/vault/pull/19703)]
* ui: pass encodeBase64 param to HMAC transit-key-actions. [[GH-19429](https://github.com/hashicorp/vault/pull/19429)]
* ui: use URLSearchParams interface to capture namespace param from SSOs (ex. ADFS) with decoded state param in callback url [[GH-19460](https://github.com/hashicorp/vault/pull/19460)]

## 1.12.4
### March 01, 2023

SECURITY:
* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* core: Bump Go version to 1.19.6.

IMPROVEMENTS:

* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]
* ui: remove wizard [[GH-19220](https://github.com/hashicorp/vault/pull/19220)]

BUG FIXES:

* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* license (enterprise): Fix bug where license would update even if the license didn't change.
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* secrets/ad: Fix bug where config couldn't be updated unless binddn/bindpass were included in the update. [[GH-18207](https://github.com/hashicorp/vault/pull/18207)]
* secrets/pki: Revert fix for PR [18938](https://github.com/hashicorp/vault/pull/18938) [[GH-19037](https://github.com/hashicorp/vault/pull/19037)]
* server/config:  Use file.Stat when checking file permissions when VAULT_ENABLE_FILE_PERMISSIONS_CHECK is enabled [[GH-19311](https://github.com/hashicorp/vault/pull/19311)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]
* ui: fixes reliance on secure context (https) by removing methods using the Crypto interface [[GH-19410](https://github.com/hashicorp/vault/pull/19410)]
* ui: show Get credentials button for static roles detail page when a user has the proper permissions. [[GH-19190](https://github.com/hashicorp/vault/pull/19190)]

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
* expiration: Prevent panics on perf standbys when an irrevocable lease gets deleted. [[GH-18401](https://github.com/hashicorp/vault/pull/18401)]
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

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

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

## 1.11.12
### June 21, 2023

CHANGES:

* core: Bump Go version to 1.19.10.
* licensing (enterprise): Terminated licenses will no longer result in shutdown. Instead, upgrades
will not be allowed if the license termination time is before the build date of the binary.

FEATURES:

* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* core (enterprise): Add background worker for automatic reporting of billing
information. [[GH-19625](https://github.com/hashicorp/vault/pull/19625)]

IMPROVEMENTS:

* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: generate hyperloglogs containing clientIds for each month during precomputation [[GH-16146](https://github.com/hashicorp/vault/pull/16146)]
* core/activity: refactor activity log api to reuse partial api functions in activity endpoint when current month is specified [[GH-16162](https://github.com/hashicorp/vault/pull/16162)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* core/activity: use monthly hyperloglogs to calculate new clients approximation for current month [[GH-16184](https://github.com/hashicorp/vault/pull/16184)]
* core: Activity log goroutine management improvements to allow tests to be more deterministic. [[GH-17028](https://github.com/hashicorp/vault/pull/17028)]
* core: Limit activity log client count usage by namespaces [[GH-16000](https://github.com/hashicorp/vault/pull/16000)]
* storage/raft: add additional raft metrics relating to applied index and heartbeating; also ensure OSS standbys emit periodic metrics. [[GH-12166](https://github.com/hashicorp/vault/pull/12166)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]

BUG FIXES:

* core/activity: add namespace breakdown for new clients when date range spans multiple months, including the current month. [[GH-18766](https://github.com/hashicorp/vault/pull/18766)]
* core/activity: de-duplicate namespaces when historical and current month data are mixed [[GH-18452](https://github.com/hashicorp/vault/pull/18452)]
* core/activity: fix the end_date returned from the activity log endpoint when partial counts are computed [[GH-17856](https://github.com/hashicorp/vault/pull/17856)]
* core/activity: include mount counts when de-duplicating current and historical month data [[GH-18598](https://github.com/hashicorp/vault/pull/18598)]
* core/activity: report mount paths (rather than mount accessors) in current month activity log counts and include deleted mount paths in precomputed queries. [[GH-18916](https://github.com/hashicorp/vault/pull/18916)]
* core/activity: return partial month counts when querying a historical date range and no historical data exists. [[GH-17935](https://github.com/hashicorp/vault/pull/17935)]
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs

## 1.11.11
### June 08, 2023

CHANGES:

* core: Bump Go version to 1.19.9.
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]

IMPROVEMENTS:

* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]
* ui: update TTL picker for consistency [[GH-18114](https://github.com/hashicorp/vault/pull/18114)]

BUG FIXES:

* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* cli: CLI should take days as a unit of time for ttl like flags [[GH-20477](https://github.com/hashicorp/vault/pull/20477)]
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core: prevent panic on login after namespace is deleted that had mfa enforcement [[GH-20375](https://github.com/hashicorp/vault/pull/20375)]
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation

## 1.11.10
### April 26, 2023

CHANGES:

* core: Bump Go version to 1.19.8.

IMPROVEMENTS:

* cli/namespace: Add detailed flag to output additional namespace information
such as namespace IDs and custom metadata. [[GH-20243](https://github.com/hashicorp/vault/pull/20243)]
* core/activity: add an endpoint to write test activity log data, guarded by a build flag [[GH-20019](https://github.com/hashicorp/vault/pull/20019)]
* core: Add a `raft` sub-field to the `storage` and `ha_storage` details provided by the
`/sys/config/state/sanitized` endpoint in order to include the `max_entry_size`. [[GH-20044](https://github.com/hashicorp/vault/pull/20044)]
* sdk/ldaputil: added `connection_timeout` to tune connection timeout duration
for all LDAP plugins. [[GH-20144](https://github.com/hashicorp/vault/pull/20144)]

BUG FIXES:

* auth/ldap: Add max_page_size configurable to LDAP configuration [[GH-19032](https://github.com/hashicorp/vault/pull/19032)]
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Fix regression breaking non-raft clusters whose nodes share the same cluster_addr/api_addr. [[GH-19721](https://github.com/hashicorp/vault/pull/19721)]
* helper/random: Fix race condition in string generator helper [[GH-19875](https://github.com/hashicorp/vault/pull/19875)]
* openapi: Fix many incorrect details in generated API spec, by using better techniques to parse path regexps [[GH-18554](https://github.com/hashicorp/vault/pull/18554)]
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* secrets/pki: Fix patching of leaf_not_after_behavior on issuers. [[GH-20341](https://github.com/hashicorp/vault/pull/20341)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* ui: Fix OIDC provider logo showing when domain doesn't match [[GH-20263](https://github.com/hashicorp/vault/pull/20263)]
* ui: Fix bad link to namespace when namespace name includes `.` [[GH-19799](https://github.com/hashicorp/vault/pull/19799)]
* ui: fixes browser console formatting for help command output [[GH-20064](https://github.com/hashicorp/vault/pull/20064)]
* ui: remove use of htmlSafe except when first sanitized [[GH-20235](https://github.com/hashicorp/vault/pull/20235)]

## 1.11.9
### March 29, 2023

SECURITY:

* storage/mssql: When using Vault’s community-supported Microsoft SQL (MSSQL) database storage backend, a privileged attacker with the ability to write arbitrary data to Vault’s configuration may be able to perform arbitrary SQL commands on the underlying database server through Vault. This vulnerability, CVE-2023-0620, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-12](https://discuss.hashicorp.com/t/hcsec-2023-12-vault-s-microsoft-sql-database-storage-backend-vulnerable-to-sql-injection-via-configuration-file/52080)]
* secrets/pki: Vault’s PKI mount issuer endpoints did not correctly authorize access to remove an issuer or modify issuer metadata, potentially resulting in denial of service of the PKI mount. This bug did not affect public or private key material, trust chains or certificate issuance. This vulnerability, CVE-2023-0665, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-11](https://discuss.hashicorp.com/t/hcsec-2023-11-vault-s-pki-issuer-endpoint-did-not-correctly-authorize-access-to-issuer-metadata/52079)]
* core: HashiCorp Vault’s implementation of Shamir’s secret sharing used precomputed table lookups, and was vulnerable to cache-timing attacks. An attacker with access to, and the ability to observe a large number of unseal operations on the host through a side channel may reduce the search space of a brute force effort to recover the Shamir shares. This vulnerability, CVE-2023-25000, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-10](https://discuss.hashicorp.com/t/hcsec-2023-10-vault-vulnerable-to-cache-timing-attacks-during-seal-and-unseal-operations/52078)]

IMPROVEMENTS:

* auth/github: Allow for an optional Github auth token environment variable to make authenticated requests when fetching org id
website/docs: Add docs for `VAULT_AUTH_CONFIG_GITHUB_TOKEN` environment variable when writing Github config [[GH-19244](https://github.com/hashicorp/vault/pull/19244)]
* core: Allow overriding gRPC connect timeout via VAULT_GRPC_MIN_CONNECT_TIMEOUT. This is an env var rather than a config setting because we don't expect this to ever be needed.  It's being added as a last-ditch
option in case all else fails for some replication issues we may not have fully reproduced. [[GH-19676](https://github.com/hashicorp/vault/pull/19676)]
* core: validate name identifiers in mssql physical storage backend prior use [[GH-19591](https://github.com/hashicorp/vault/pull/19591)]

BUG FIXES:

* auth/kubernetes: Ensure a consistent TLS configuration for all k8s API requests [[#190](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/190)] [[GH-19720](https://github.com/hashicorp/vault/pull/19720)]
* cli: Fix vault read handling to return raw data as secret.Data when there is no top-level data object from api response. [[GH-17913](https://github.com/hashicorp/vault/pull/17913)]
* core (enterprise): Attempt to reconnect to a PKCS#11 HSM if we retrieve a CKR_FUNCTION_FAILED error.
* core: Fixed issue with remounting mounts that have a non-trailing space in the 'to' or 'from' paths. [[GH-19585](https://github.com/hashicorp/vault/pull/19585)]
* openapi: Fix logic for labeling unauthenticated/sudo paths. [[GH-19600](https://github.com/hashicorp/vault/pull/19600)]
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* ui: fixes issue navigating back a level using the breadcrumb from secret metadata view [[GH-19703](https://github.com/hashicorp/vault/pull/19703)]
* ui: pass encodeBase64 param to HMAC transit-key-actions. [[GH-19429](https://github.com/hashicorp/vault/pull/19429)]
* ui: use URLSearchParams interface to capture namespace param from SSOs (ex. ADFS) with decoded state param in callback url [[GH-19460](https://github.com/hashicorp/vault/pull/19460)]

## 1.11.8
### March 01, 2023

SECURITY:

* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* core: Bump Go version to 1.19.6.

IMPROVEMENTS:

* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]

BUG FIXES:

* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* license (enterprise): Fix bug where license would update even if the license didn't change.
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* secrets/ad: Fix bug where config couldn't be updated unless binddn/bindpass were included in the update. [[GH-18208](https://github.com/hashicorp/vault/pull/18208)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]
* ui: show Get credentials button for static roles detail page when a user has the proper permissions. [[GH-19190](https://github.com/hashicorp/vault/pull/19190)]

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

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

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

SECURITY:

* core: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. This vulnerability, CVE-2022-40186, is fixed in 1.11.3, 1.10.6, and 1.9.9. [[HSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

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

SECURITY:

* storage/raft: Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]

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

## 1.10.11
### March 01, 2023

SECURITY:

* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* core: Bump Go version to 1.19.6.

IMPROVEMENTS:

* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]

BUG FIXES:

* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* secrets/ad: Fix bug where config couldn't be updated unless binddn/bindpass were included in the update. [[GH-18209](https://github.com/hashicorp/vault/pull/18209)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]

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

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

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

SECURITY:

* core: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. This vulnerability, CVE-2022-40186, is fixed in 1.11.3, 1.10.6, and 1.9.9. [[HSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

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

SECURITY:

* storage/raft: Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]

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

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

BUG FIXES:

* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]

## 1.9.9
### August 31, 2022

SECURITY:

* core: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. This vulnerability, CVE-2022-40186, is fixed in 1.11.3, 1.10.6, and 1.9.9. [[HSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

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

SECURITY:

* storage/raft: Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]

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
