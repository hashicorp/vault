# Enos Test Automation — Creation & Development Guide

**Activation Keywords:**
- **Enos Scenarios/Modules**: "create enos test", "create enos scenario", "create enos module", "new enos", "enos automation"
- **Blackbox Go Tests**: "create sdk test", "create blackbox test", "new blackbox test", "convert to blackbox test"

**Purpose:** Guide for creating and developing Enos scenarios, modules, and blackbox Go tests.

---

## How to Understand the Codebase Before Writing Anything

### Writing an Enos scenario → read the ENTIRE `enos/` project first

Before creating or modifying any scenario, read enough of the project to understand the existing patterns:

| File | Purpose |
|------|---------|
| `enos/enos-globals.hcl` | Global values: archs, backends, distros, editions, ports, tags |
| `enos/enos-variables.hcl` | All input variables and their defaults |
| `enos/enos-modules.hcl` | Module registry — every available `module.*` reference |
| `enos/enos-modules-ent.hcl` | Enterprise-only module extensions |
| `enos/enos-qualities.hcl` | Quality assertions used in `verifies = [...]` blocks |
| `enos/enos-descriptions.hcl` | Reusable step descriptions (`global.description.*`) |
| `enos/enos-scenario-smoke-sdk.hcl` | **Primary reference scenario** for blackbox-based testing |
| `enos/enos-scenario-smoke.hcl` | Reference scenario for traditional enos-based testing |
| `enos/enos-scenario-*.hcl` | All other scenarios — check for reuse patterns |
| `enos/modules/vault_run_blackbox_test/` | The module that actually runs blackbox tests |

### Writing a blackbox Go test → read ALL three source trees

| Tree | What it tells you |
|------|-------------------|
| `sdk/helper/testcluster/blackbox/` | Every available SDK helper method — **read this first** |
| `vault/external_tests/blackbox/` | Existing test patterns, utilities, build tags |
| `enos/` | How tests are invoked, what env vars are injected, which scenarios run what |

---

## Directory Map

### Enos project: `enos/`
```
enos/
├── enos-globals.hcl                   # Global constants (archs, editions, backends, ports, tags)
├── enos-variables.hcl                 # All input variables
├── enos-modules.hcl                   # Module registry
├── enos-modules-ent.hcl               # Enterprise module extensions
├── enos-qualities.hcl                 # Quality assertions
├── enos-descriptions.hcl              # Step description strings
├── enos-scenario-smoke-sdk.hcl        # SDK/blackbox scenario (primary reference)
├── enos-scenario-smoke.hcl            # Traditional enos scenario
├── enos-scenario-*.hcl                # All other scenarios
└── modules/
    ├── vault_run_blackbox_test/       # Runs blackbox Go tests (main.tf, variables.tf, plugin.tf)
    ├── vault_cluster/                 # Deploy a Vault cluster
    ├── vault_get_cluster_ips/         # Resolve leader/follower IPs
    ├── vault_wait_for_leader/         # Wait for cluster to elect a leader
    ├── target_ec2_instances/          # Provision EC2 nodes
    ├── set_up_external_integration_target/  # LDAP, PostgreSQL, MongoDB containers
    └── ...
```

### SDK helpers: `sdk/helper/testcluster/blackbox/`
```
session.go              # New(), WithNoCleanup(), SkipIfVersion* — session lifecycle
session_auth.go         # Login, LoginUserpass, TryLoginUserpass, AssertWriteFails, AssertReadFails
session_autopilot.go    # AssertAutopilotHealthy
session_client.go       # Req(), WithClientRootNamespace/ParentNamespace/Timeout, GetParentNamespace
session_cluster.go      # AssertClusterHealthy, EventuallyClusterHealthyUnsealed
session_config.go       # MustSanitizedConfig, MustGetConfigStorageType
session_dynamic.go      # MustGenerateCreds, MustRevokeLease, AssertLeaseExists/Revoked,
                        #   MustConfigureDBConnection, MustCreateDBRole, MustCheckCreds
session_ha.go           # (internal) getClusterNodeCount, haNodes, haActiveNode
session_logical.go      # Read, MustRead, MustWrite, MustList, MustReadRequired,
                        #   MustDelete, MustWriteKV2, MustReadKV2
session_metrics.go      # AssertMetricGaugeValue
session_ops.go          # AssertUIAvailable, AssertFileDoesNotContainSecret
session_pki.go          # MustSetupPKIRoot
session_plugin.go       # MustRegisterPlugin, MustEnablePlugin, AssertPluginRegistered,
                        #   AssertPluginConfigured, PluginSession, ResponseValidator, SequenceOperation
session_raft.go         # AssertRaftStable, EventuallyRaftClusterHealthy, AssertRaftClusterHealthy,
                        #   MustRaftRemovePeer, AssertRaftPeerRemoved, MustGetCurrentLeader,
                        #   MustStepDownLeader, MustGetClusterNodeCount, MustGetNonLeaderNode
session_replication.go  # GetDRReplicationMode, AssertReplicationDisabled,
                        #   AssertDRReplicationStatus, AssertPerformanceReplicationStatus
session_seal.go         # (internal) sealed
session_status.go       # AssertUnsealed, AssertUnsealedAny, AssertVersion, AssertBuildDate,
                        #   AssertRevision, AssertCLIVersion, AssertServerVersion
session_sys.go          # MustEnableSecretsEngine, MustDisableSecretsEngine, MustEnableAuth,
                        #   MustWritePolicy
session_token.go        # MustCreateToken (TokenOptions: Policies, TTL, Renewable, NoParent)
                        #   AssertTokenIsValid
session_util.go         # Eventually, EventuallyWithTimeout, WithRootNamespace, WithParentNamespace
assertions.go           # AssertSecret → SecretAssertion → .Data() / .KV2()
                        #   MapAssertion: HasKey, GetKey, HasKeyCustom, HasKeyExists, GetMap, GetSlice
                        #   SliceAssertion: Length, FindMap, AllHaveKey, AllHaveKeyCustom, NoneHaveKeyVal
plugin_binary.go        # BuildPluginBinary, DeployPluginBinary, ValidatePluginBinary,
                        #   GetPluginBinaryInfo, ComparePluginVersions
plugin_config.go        # PluginConfig, NewPluginConfig, SetRequired/Optional, AddValidator,
                        #   SetConfig, SetEnvironment, ValidateConfig, LoadPluginConfig, ApplyPluginConfig
plugin_registry.go      # BuiltinPluginRegistry, DefaultBuiltinPluginRegistry, DefaultRegistry,
                        #   Register*/List*/GetBuiltin* helpers
plugin_utils.go         # ExtendedPluginRegistration, BatchPluginRegistration,
                        #   SetupBuiltin*/SetupExternal* convenience wrappers
```

### Go tests: `vault/external_tests/blackbox/`
```
doc.go                             # Package declaration
test_utils.go                      # Shared helpers: SetupKVEngine, WaitForKVEngineReady,
                                   #   SetupUserpassAuth, SetupStandardKVUserpass,
                                   #   AssertKVData, CreateTestToken
                                   #   StandardKVData, AltKVData, StandardOpsPolicy, ReadOnlyPolicy
isolated/
  auth/                            # //go:build isolated — auth method tests
    auth_engines_test.go
    auth_userpass_test.go
    token_test.go
  plugins/                         # //go:build isolated — plugin tests (require external services)
    aws/                           # Requires AWS credentials via env vars
    ldap/                          # Requires LDAP_URL_*, LDAP_BIND_*, LDAP_USERNAME env vars
    mongodb/                       # Requires MONGO_URL env var
    postgresql/                    # Requires PG_URL / PGHOST etc. env vars
  secrets/                         # //go:build isolated — secrets engine tests
    dynamic_test.go, engines_test.go, kvv2_test.go, pki_test.go
    secrets_identity_test.go, secrets_kmip_test.go, secrets_kv_test.go
    secrets_pki_test.go, secrets_ssh_test.go, secrets_transit_test.go
  verify/                          # //go:build isolated — state verification
    ibm_license_update_test.go     # Requires VAULT_IBM_LICENSE_EDITION env var
    performance_replication_test.go # Requires VAULT_SECONDARY_ADDR, VAULT_SECONDARY_TOKEN
    replication_test.go
    ui_test.go
    undo_logs_test.go
    unsealed_test.go
    version_core_test.go
    version_verification_test.go   # Requires VAULT_VERSION, VAULT_BUILD_DATE, VAULT_REVISION
scenario/
  ha/                              # //go:build scenario — HA tests
    stepdown_test.go
  raft/                            # //go:build scenario — Raft multi-node tests
    node_operations_test.go
    raft_node_removal_test.go
    voters_test.go
system/
  config/                          # //go:build system — sys/* endpoint tests
    billing_start_date_test.go
    billing_test.go
    license_test.go
verify/                            # //go:build testonly — legacy location; prefer isolated/verify
  replication_test.go
  ui_test.go
```

---

## Build Tags

Every test file must declare a build tag. The tag maps to the test type and is also used by Go to filter what's compiled and run.

| Tag | Location | When to use |
|-----|----------|-------------|
| `isolated` | `isolated/` | Tests that run independently without multi-node coordination |
| `scenario` | `scenario/` | Tests that need multi-node cluster operations (HA, Raft) |
| `system` | `system/` | Tests for `sys/` endpoints |
| `testonly` | `verify/` (legacy) | Old tag — prefer `isolated` for new tests |

```go
//go:build isolated
// +build isolated
```

---

## Environment Variables Injected by `vault_run_blackbox_test`

The module at `enos/modules/vault_run_blackbox_test/main.tf` injects these automatically:

| Variable | Source |
|----------|--------|
| `VAULT_TOKEN` | `var.vault_root_token` |
| `VAULT_ADDR` | `var.vault_addr` or derived from `var.leader_public_ip` |
| `VAULT_EDITION` | `var.vault_edition` (matrix.edition) |
| `VAULT_VERSION` | `var.vault_product_version` |
| `VAULT_REVISION` | `var.vault_revision` |
| `VAULT_BUILD_DATE` | `var.vault_build_date` |
| `VAULT_INSTALL_DIR` | `var.vault_install_dir` |
| `VAULT_IBM_LICENSE_EDITION` | `var.vault_ibm_license_edition` |
| `VAULT_NAMESPACE` | `var.vault_namespace` (optional) |
| `EXPECTED_STATE` | Auto-injected when `test_package` contains `isolated/verify` |
| `TIMEOUT_SECONDS` | Auto-injected when `test_package` contains `isolated/verify` |
| `RETRY_INTERVAL` | Auto-injected when `test_package` contains `isolated/verify` |
| `LDAP_URL_PRIVATE/PUBLIC`, `LDAP_BIND_DN/PASS`, `LDAP_USERNAME` | From `integration_host_state.ldap` |
| `PG_URL`, `POSTGRES_*`, `PGHOST/PORT/USER/PASSWORD/DATABASE` | From `integration_host_state.postgres` |
| `MONGO_URL`, `MONGO_URL_PRIVATE` | From `integration_host_state.mongodb` |

Override with `test_env_vars = { "KEY" = "value" }` in the scenario step.

---

## Writing a New Blackbox Go Test

### Step 1 — Choose location and build tag

| Test type | Location | Build tag |
|-----------|----------|-----------|
| Single-feature, no external deps | `isolated/secrets/` or `isolated/auth/` | `isolated` |
| Plugin with external service | `isolated/plugins/<service>/` | `isolated` |
| Multi-node / cluster operation | `scenario/raft/` or `scenario/ha/` | `scenario` |
| System endpoint (`sys/`) | `system/config/` | `system` |
| State verification | `isolated/verify/` | `isolated` |

### Step 2 — Check SDK helpers before writing anything

Open `sdk/helper/testcluster/blackbox/` and find:
- Is there already a method that does what you need? (`AssertRaftClusterHealthy`, `MustGenerateCreds`, etc.)
- Does `vault/external_tests/blackbox/test_utils.go` have a shared helper you should reuse?

### Step 3 — Write the test

```go
//go:build isolated
// +build isolated

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets  // matches the directory name

import (
    "testing"

    "github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
    "github.com/stretchr/testify/require"
    helpers "github.com/hashicorp/vault/vault/external_tests/blackbox" // shared utils
)

// TestFeatureName verifies [description].
// Replicates: enos/modules/MODULE_NAME/scripts/SCRIPT.sh  (if converting from shell)
func TestFeatureName(t *testing.T) {
    t.Parallel()
    v := blackbox.New(t)

    // Verify cluster is ready
    v.AssertUnsealedAny()

    // Test logic — prefer SDK helpers over raw v.Client.Logical() calls
    secret := v.MustRead("path/to/resource")
    require.NotNil(t, secret)

    v.AssertSecret(secret).Data().
        HasKey("field", "expected_value").
        HasKeyExists("other_field")
}
```

### Step 4 — Converting a shell script to Go

Map each shell construct:

| Shell | Go |
|-------|----|
| `vault read path` | `v.MustRead("path")` or `v.Client.Logical().Read("path")` |
| `vault write path k=v` | `v.MustWrite("path", map[string]any{"k": "v"})` |
| `vault list path` | `v.MustList("path")` |
| `jq -r ".data.field"` | `secret.Data["field"]` or `.AssertSecret(s).Data().HasKey(...)` |
| `if [ condition ]` | `require.True(t, condition, "msg")` |
| `retry N cmd` | `v.Eventually(fn)` or `v.EventuallyWithTimeout(fn, timeout)` |
| `echo "msg" 1>&2` | `t.Logf("msg")` |
| `exit 1` | `t.Fatal("msg")` |
| `vault auth enable userpass` | `v.MustEnableAuth("userpass", &api.EnableAuthOptions{Type: "userpass"})` |
| `vault secrets enable -path=p kv-v2` | `v.MustEnableSecretsEngine("p", &api.MountInput{Type: "kv-v2"})` |
| `vault write auth/userpass/users/u password=p policies=pol` | `v.MustWrite("auth/userpass/users/u", map[string]any{"password": "p", "policies": "pol"})` |

**Always notify** if any shell step cannot be converted, if an API endpoint differs, or if validation logic was simplified.

---

## Writing a New Enos Scenario

### Mandatory: read these files before starting

1. `enos/enos-globals.hcl` — available matrix values
2. `enos/enos-variables.hcl` — all available variables
3. `enos/enos-modules.hcl` (and `enos-modules-ent.hcl`) — every `module.*` you can reference
4. `enos/enos-qualities.hcl` — quality assertions for `verifies = [...]`
5. `enos/enos-descriptions.hcl` — reusable `global.description.*` strings
6. `enos/enos-scenario-smoke-sdk.hcl` — canonical reference for a blackbox-based scenario
7. The closest existing scenario to your use case

### Standard scenario step flow (for blackbox testing)

```
build_vault
  → ec2_info
  → create_vpc
    → read_backend_license (skip if raft or ce)
    → read_vault_license   (skip if ce)
    → create_seal_key
    → create_external_integration_target
    → create_vault_cluster_targets
    → create_vault_cluster_backend_targets
      → set_up_external_integration_target
      → create_backend_cluster
      → create_vault_cluster
        → get_local_metadata (skip if not local)
        → wait_for_leader
          → get_vault_cluster_ips
            → run_blackbox_tests
```

### Scenario template

```hcl
scenario "my_scenario" {
  description = <<-EOF
    Short description of what this scenario tests.
    Runs: vault/external_tests/blackbox/<type>/<package>
  EOF

  matrix {
    arch            = global.archs
    artifact_source = global.artifact_sources
    artifact_type   = global.artifact_types
    backend         = global.backends
    config_mode     = global.config_modes
    consul_edition  = global.consul_editions
    consul_version  = global.consul_versions
    distro          = global.distros_aws
    edition         = global.editions
    ip_version      = global.ip_versions
    seal            = global.seals

    exclude {
      artifact_source = ["local"]
      artifact_type   = ["package"]
    }
    exclude {
      seal    = ["pkcs11"]
      edition = [for e in matrix.edition : e if !strcontains(e, "hsm")]
    }
    exclude {
      seal   = ["pkcs11"]
      distro = ["leap", "sles"]
    }
    exclude {
      ip_version = ["6"]
      backend    = ["consul"]
    }
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ec2_user,
    provider.enos.ubuntu,
    provider.time.default,
  ]

  locals {
    artifact_path  = matrix.artifact_source != "artifactory" ? abspath(var.vault_artifact_path) : null
    enos_provider  = {
      amzn   = provider.enos.ec2_user
      leap   = provider.enos.ec2_user
      rhel   = provider.enos.ec2_user
      sles   = provider.enos.ec2_user
      ubuntu = provider.enos.ubuntu
    }
    manage_service = matrix.artifact_type == "bundle"

    // Test filter logic — accepts test function names ("Test*") or package paths
    is_test_name_filter   = length(var.blackbox_test_filter) > 0 &&
      length([for t in var.blackbox_test_filter : t if can(regex("^Test", t))]) > 0
    default_test_packages = ["isolated/secrets", "isolated/auth"]
    test_packages = length(var.blackbox_test_filter) > 0 && !local.is_test_name_filter ? [
      for pkg in var.blackbox_test_filter : "./vault/external_tests/blackbox/${pkg}/..."
      ] : [
      for pkg in local.default_test_packages : "./vault/external_tests/blackbox/${pkg}/..."
    ]
  }

  // --- infrastructure steps (copy from smoke_sdk reference) ---
  step "build_vault"     { ... }
  step "ec2_info"        { ... }
  step "create_vpc"      { ... }
  step "read_backend_license" { ... }
  step "read_vault_license"   { ... }
  step "create_seal_key" { ... }
  step "create_external_integration_target" { ... }
  step "create_vault_cluster_targets"         { ... }
  step "create_vault_cluster_backend_targets" { ... }
  step "set_up_external_integration_target"   { ... }
  step "create_backend_cluster" { ... }
  step "create_vault_cluster"   { ... }
  step "get_local_metadata"     { skip_step = matrix.artifact_source != "local" ... }
  step "wait_for_leader"        { ... }
  step "get_vault_cluster_ips"  { ... }

  step "run_blackbox_tests" {
    description = "Run blackbox SDK tests"
    module      = module.vault_run_blackbox_test
    depends_on  = [step.get_vault_cluster_ips, step.set_up_external_integration_target]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      leader_host            = step.get_vault_cluster_ips.leader_host
      leader_public_ip       = step.get_vault_cluster_ips.leader_public_ip
      vault_root_token       = step.create_vault_cluster.root_token
      test_names             = local.is_test_name_filter ? var.blackbox_test_filter : null
      test_package           = local.is_test_name_filter ? "./vault/external_tests/blackbox" : join(" ", local.test_packages)
      integration_host_state = step.set_up_external_integration_target.state
      vault_edition          = matrix.edition
      vault_product_version  = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_revision         = matrix.artifact_source == "local" ? step.get_local_metadata.revision : var.vault_revision
      vault_build_date       = matrix.artifact_source == "local" ? step.get_local_metadata.build_date : var.vault_build_date
    }
  }

  // --- standard outputs (copy from smoke_sdk reference) ---
  output "cluster_name"        { value = step.create_vault_cluster.cluster_name }
  output "root_token"          { value = step.create_vault_cluster.root_token }
  // ... etc
}
```

---

## `vault_run_blackbox_test` Module Variables Reference

| Variable | Type | Description |
|----------|------|-------------|
| `leader_host` | object | Host object with `private_ip` and `public_ip` |
| `leader_public_ip` | string | Public IP of the Vault leader |
| `vault_root_token` | string | Root token |
| `test_package` | string | Space-separated Go package path(s) |
| `test_names` | list(string) | Specific test function names; null = run all |
| `vault_addr` | string | Override Vault address (cloud environments) |
| `vault_namespace` | string | Vault namespace (HCP) |
| `integration_host_state` | any | Output of `set_up_external_integration_target` |
| `vault_edition` | string | e.g. `"ent"`, `"ce"`, `"ent.hsm"` |
| `vault_product_version` | string | e.g. `"2.0.0"` |
| `vault_revision` | string | Git commit SHA |
| `vault_build_date` | string | Build date string |
| `vault_ibm_license_edition` | string | IBM PAO license edition |
| `test_env_vars` | map(string) | Extra env vars to inject |
| `verify_expected_state` | string | Default `EXPECTED_STATE` (auto for `isolated/verify`) |
| `verify_timeout_seconds` | string | Default `TIMEOUT_SECONDS` (auto for `isolated/verify`) |
| `verify_retry_interval` | string | Default `RETRY_INTERVAL` (auto for `isolated/verify`) |

---

## Quick Reference

### Run a scenario

```bash
cd enos/
enos scenario validate smoke_sdk
enos scenario launch smoke_sdk
enos scenario destroy smoke_sdk

# Run with a specific test function
ENOS_VAR_blackbox_test_filter='["TestMyFeature"]' enos scenario launch smoke_sdk

# Run a specific package
ENOS_VAR_blackbox_test_filter='["isolated/secrets"]' enos scenario launch smoke_sdk

# Debug
ENOS_DEBUG=1 enos scenario launch smoke_sdk
```

### Key files at a glance

| What you need | Where it lives |
|---------------|----------------|
| Matrix values (archs, editions…) | `enos/enos-globals.hcl` |
| Input variables | `enos/enos-variables.hcl` |
| Available modules | `enos/enos-modules.hcl` |
| Scenario that runs blackbox tests | `enos/enos-scenario-smoke-sdk.hcl` |
| Module that executes Go tests | `enos/modules/vault_run_blackbox_test/` |
| All SDK helper methods | `sdk/helper/testcluster/blackbox/*.go` |
| Shared test utilities | `vault/external_tests/blackbox/test_utils.go` |
| Test files | `vault/external_tests/blackbox/<type>/<package>/*_test.go` |

### Reference tests

| What it shows | File |
|---------------|------|
| Version / build date checks | `isolated/verify/version_verification_test.go` |
| Performance replication (multi-cluster) | `isolated/verify/performance_replication_test.go` |
| System endpoint (`sys/`) | `system/config/billing_start_date_test.go` |
| LDAP plugin (external service) | `isolated/plugins/ldap/secrets_ldap_test.go` |
| Raft multi-node ops | `scenario/raft/raft_node_removal_test.go` |
| KV secrets engine | `isolated/secrets/secrets_kv_test.go` |
| Assertions API usage | `sdk/helper/testcluster/blackbox/assertions.go` |
