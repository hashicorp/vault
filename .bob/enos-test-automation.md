# Enos Test Automation - Creation & Development Guide

**Activation Keywords:**
- **Enos Infrastructure**: "create enos test", "create enos scenario", "create enos module", "new enos", "enos automation"
- **Blackbox Tests**: "create sdk test", "create test sdk", "create blackbox test", "create test blackbox", "new blackbox test"

**Purpose:** Specialized guide for CREATING and DEVELOPING new Enos scenarios, modules, and test infrastructure.

**Special Focus for Blackbox Test Keywords:**
When you use keywords like "create sdk test", "create blackbox test", focus primarily on:
- `sdk/helper/testcluster/blackbox/` - Test SDK utilities and helpers
- `vault/external_tests/blackbox/` - Existing test examples and patterns

## Critical Directories

### 1. `enos/` - Test Infrastructure
```
enos/
├── enos-scenario-*.hcl          # Scenario definitions
├── enos-samples-*.hcl           # CI sample configurations
├── enos-modules.hcl             # Module registry
├── enos-globals.hcl             # Global variables
├── modules/                     # Terraform modules
│   ├── vault_cluster/           # Deploy Vault
│   ├── vault_run_blackbox_test/ # Execute tests
│   └── target_ec2_instances/    # Create infrastructure
```

### 2. `vault/external_tests/blackbox/` - Test Code
```
vault/external_tests/blackbox/
├── isolated/              # Isolated test categories
│   ├── auth/             # Auth method tests (userpass, token)
│   ├── secrets/          # Secrets engine tests (KV, PKI, Transit, SSH, KMIP)
│   ├── plugins/          # Plugin tests with external services
│   │   ├── ldap/         # LDAP secrets engine tests
│   │   ├── aws/          # AWS secrets engine tests
│   │   ├── mongodb/      # MongoDB database tests
│   │   └── postgresql/   # PostgreSQL database tests
│   └── verify/           # Version & state verification tests
├── scenario/             # Multi-step scenario tests
│   ├── ha/              # High availability scenarios
│   └── raft/            # Raft cluster operations
├── system/              # System endpoint tests
│   └── config/          # System config tests (billing, license)
└── verify/              # Legacy verify tests (being migrated)
```

### 3. `sdk/helper/testcluster/blackbox/` - Test SDK
```
sdk/helper/testcluster/blackbox/
├── session.go              # Main session & client
├── session_status.go       # Status assertions
├── session_raft.go         # Raft helpers
```

### 4. `.github/` - CI/CD Workflows
```
.github/
├── workflows/
│   ├── test-run-enos-scenario-matrix.yml  # Main test executor
│   ├── build.yml                          # Orchestration
└── actions/
    ├── run-enos-scenario/                 # Scenario execution
    └── create-dynamic-config/             # Config generation
```

**Execution Flow:**
```
GitHub Workflow → Enos Scenario → Blackbox Tests → SDK Helpers
(.github/)        (enos/)          (vault/...)      (sdk/...)
```

## Quick Start

```bash
# Create scenario
touch enos/enos-scenario-mytest.hcl
enos scenario validate mytest
enos scenario launch mytest

# Create test
touch vault/external_tests/blackbox/mypackage/mytest_test.go

# Debug
ENOS_DEBUG=1 enos scenario launch mytest
ENOS_VAR_blackbox_test_filter='["TestMyTest"]' enos scenario launch mytest
```

## Enos Scenario Template

```hcl
// File: enos/enos-scenario-mytest.hcl
scenario "mytest" {
  description = "Tests [feature]. Runs vault/external_tests/blackbox/[package]"

  matrix {
    arch            = global.archs
    artifact_source = global.artifact_sources
    backend         = global.backends
    distro          = global.distros
    edition         = global.editions
    
    exclude {
      edition = ["ce"]  // If ENT-only
    }
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [provider.aws.default, provider.enos.ubuntu]

  locals {
    enos_provider = {
      ubuntu = provider.enos.ubuntu
      rhel   = provider.enos.ec2_user
    }
  }

  step "build_vault" {
    module = "build_${matrix.artifact_source}"
    variables {
      goarch  = matrix.arch
      edition = matrix.edition
    }
  }

  step "create_vpc" {
    module = module.create_vpc
    variables { common_tags = global.tags }
  }

  step "create_vault_cluster_targets" {
    module     = module.target_ec2_instances
    depends_on = [step.create_vpc]
    providers { enos = local.enos_provider[matrix.distro] }
    variables {
      ami_id = step.ec2_info.ami_ids[matrix.arch][matrix.distro]
      vpc_id = step.create_vpc.id
    }
  }

  step "create_vault_cluster" {
    module     = module.vault_cluster
    depends_on = [step.create_vault_cluster_targets]
    providers { enos = local.enos_provider[matrix.distro] }
    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      storage_backend   = matrix.backend
    }
  }

  step "run_blackbox_tests" {
    module     = module.vault_run_blackbox_test
    depends_on = [step.create_vault_cluster]
    providers { enos = local.enos_provider[matrix.distro] }
    variables {
      leader_host       = step.get_cluster_ips.leader_host
      vault_root_token  = step.create_vault_cluster.root_token
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      // Choose appropriate path based on test type:
      // - isolated tests: "./vault/external_tests/blackbox/isolated/secrets/..."
      // - scenario tests: "./vault/external_tests/blackbox/scenario/raft/..."
      // - system tests: "./vault/external_tests/blackbox/system/..."
      test_package      = "./vault/external_tests/blackbox/isolated/secrets/..."
      test_names        = ["TestMyFeature"]  // Optional: specific test names
    }
  }

  output "test_results" {
    value = step.run_blackbox_tests.test_results_summary
  }
}
```

## Blackbox Test Organization

### Directory Structure Guidelines

**`isolated/`** - Independent tests that can run in parallel
- `auth/` - Authentication method tests
- `secrets/` - Secrets engine tests (KV, PKI, Transit, etc.)
- `plugins/` - Plugin tests requiring external services (LDAP, databases)
- `verify/` - Version and state verification tests

**`scenario/`** - Multi-step scenario tests
- `ha/` - High availability and failover scenarios
- `raft/` - Raft cluster operations and node management

**`system/`** - System-level endpoint tests
- `config/` - System configuration (billing, license, etc.)

**`verify/`** - Legacy location (being migrated to `isolated/verify/`)

### Choosing the Right Location

| Test Type | Location | Example |
|-----------|----------|---------|
| Single feature test | `isolated/secrets/` | KV write/read test |
| Plugin with external service | `isolated/plugins/ldap/` | LDAP dynamic roles |
| Multi-node operation | `scenario/raft/` | Node removal test |
| System endpoint | `system/config/` | Billing configuration |
| Version check | `isolated/verify/` | Version verification |

## Blackbox Test Template

```go
// File: vault/external_tests/blackbox/isolated/secrets/mytest_test.go
package secrets

import (
    "testing"
    "github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

func TestMyFeature(t *testing.T) {
    t.Parallel()
    v := blackbox.New(t)
    
    // SDK helpers from sdk/helper/testcluster/blackbox/
    v.AssertUnsealedAny()
    v.SkipIfVersionBelow("2.10.0")
    
    // Test logic
    v.MustWrite("secret/data/test", map[string]interface{}{
        "data": map[string]interface{}{"key": "value"},
    })
    
    secret := v.MustRead("secret/data/test")
    if secret.Data["data"] == nil {
        t.Fatal("Expected data")
    }
}
```

### Reference Existing Tests

**Examples in vault/external_tests/blackbox/:**
- `isolated/verify/version_verification_test.go` - Version checks
- `system/config/billing_test.go` - System endpoint testing
- `isolated/plugins/ldap/secrets_ldap_test.go` - Plugin with external services
- `scenario/raft/raft_node_removal_test.go` - Raft operations
- `isolated/secrets/secrets_kv_test.go` - KV secrets engine
- `isolated/auth/auth_userpass_test.go` - Auth method testing

## SDK Helpers Reference

### session.go
```go
v := blackbox.New(t)
v := blackbox.New(t, blackbox.WithNoCleanup())
secret := v.MustRead("path")
v.MustWrite("path", data)
v.MustDelete("path")
v.WithNamespace("admin", func() { /* ... */ })
```

### session_status.go
```go
v.AssertUnsealed("shamir")
v.AssertUnsealedAny()
v.AssertVersion("1.15.0")
v.SkipIfVersionBelow("2.10.0")
v.AssertBuildDate("1.15.0", "2024-01-01T00:00:00Z")
```

### session_raft.go
```go
v.EventuallyRaftClusterHealthy(timeout)
v.MustGetRaftConfiguration()
v.MustGetRaftAutopilotState()
```

## Common Patterns

### Dynamic Test Package Selection
```hcl
locals {
  // Detect if filter contains test names (starts with "Test") or package names
  is_test_name_filter = length(var.blackbox_test_filter) > 0 &&
    length([for t in var.blackbox_test_filter : t if can(regex("^Test", t))]) > 0
  
  // Default packages for smoke_sdk scenario
  default_test_packages = ["isolated/secrets", "isolated/auth", "scenario/raft", "system"]
  
  // Convert package names to full paths
  test_packages = length(var.blackbox_test_filter) > 0 && !local.is_test_name_filter ? [
    for pkg in var.blackbox_test_filter : "./vault/external_tests/blackbox/${pkg}/..."
  ] : [
    for pkg in local.default_test_packages : "./vault/external_tests/blackbox/${pkg}/..."
  ]
}

step "run_blackbox_tests" {
  variables {
    test_names   = local.is_test_name_filter ? var.blackbox_test_filter : null
    test_package = local.is_test_name_filter ?
      "./vault/external_tests/blackbox" :
      join(" ", local.test_packages)
  }
}

// Example filters:
// ENOS_VAR_blackbox_test_filter='["TestMyTest"]'           # Run specific test
// ENOS_VAR_blackbox_test_filter='["isolated/secrets"]'    # Run secrets tests
// ENOS_VAR_blackbox_test_filter='["scenario/raft"]'       # Run raft scenarios
// ENOS_VAR_blackbox_test_filter='["system/config"]'       # Run system tests
```

### External Services Setup
```hcl
step "create_service_host" {
  module = module.target_ec2_instances
  variables {
    ami_id = step.ec2_info.ami_ids["arm64"]["ubuntu"]["26.04"]
    vpc_id = step.create_vpc.id
  }
}

step "setup_external_services" {
  module = module.set_up_external_integration_target
  variables {
    hosts    = step.create_service_host.hosts
    packages = ["podman", "podman-docker"]
    ports    = [389, 5432, 27017]  // LDAP, PostgreSQL, MongoDB
  }
}

step "run_blackbox_tests" {
  depends_on = [step.setup_external_services]
  variables {
    integration_host_state = step.setup_external_services.state
  }
}
```

## GitHub CI/CD Integration

### Main Workflow
**Location:** `.github/workflows/test-run-enos-scenario-matrix.yml`

**Flow:**
1. `build.yml` triggers on PR/push
2. Calls `test-run-enos-scenario-matrix.yml` with sample names
3. Reads samples from `enos/enos-samples-*.hcl`
4. Executes via `.github/actions/run-enos-scenario/`
5. Sets ENOS_VAR_*, VAULT_*, AWS_* environment variables
6. Runs `enos scenario launch`
7. Collects results from `/tmp/vault_test_results_*.json`

### Adding Tests to CI

1. Create scenario: `enos/enos-scenario-mytest.hcl`
2. Create sample: `enos/enos-samples-mytest.hcl`
```hcl
sample "mytest_linux_amd64" {
  attributes = global.sample_attributes
  subset "smoke" {
    scenario_name = "mytest"
    matrix {
      arch            = ["amd64"]
      artifact_source = ["local"]
      artifact_type   = ["bundle"]
      backend         = ["raft"]
      distro          = ["ubuntu"]
      edition         = ["ent"]
    }
  }
}
```
3. Sample auto-discovered by `build.yml`
4. Workflow executes automatically

### Debugging CI
- Check environment variables in workflow logs
- Download test results: `/tmp/vault_test_results_*.json`
- Download JUnit results: `/tmp/vault_test_results_*.xml`
- Common issues: missing licenses, timeouts, AWS cleanup failures

## Quick Reference

| Purpose | Directory | Example |
|---------|-----------|---------|
| Scenarios | `enos/enos-scenario-*.hcl` | `enos-scenario-plugin.hcl` |
| Modules | `enos/modules/` | `modules/vault_cluster/` |
| Isolated Tests | `vault/external_tests/blackbox/isolated/` | `isolated/plugins/ldap/` |
| Scenario Tests | `vault/external_tests/blackbox/scenario/` | `scenario/raft/` |
| System Tests | `vault/external_tests/blackbox/system/` | `system/config/` |
| SDK | `sdk/helper/testcluster/blackbox/` | `blackbox/session.go` |
| CI | `.github/workflows/` | `test-run-enos-scenario-matrix.yml` |
| Samples | `enos/enos-samples-*.hcl` | `enos-samples-ent-build.hcl` |

### Commands
```bash
# Scenario management
enos scenario validate mytest
enos scenario launch mytest
enos scenario launch mytest --no-destroy
enos scenario exec mytest run_blackbox_tests
enos scenario destroy mytest

# Debug mode
ENOS_DEBUG=1 enos scenario launch mytest

# Run specific test packages
ENOS_VAR_blackbox_test_filter='["isolated/secrets"]' enos scenario launch smoke_sdk
ENOS_VAR_blackbox_test_filter='["scenario/raft"]' enos scenario launch smoke_sdk
ENOS_VAR_blackbox_test_filter='["system/config"]' enos scenario launch smoke_sdk
ENOS_VAR_blackbox_test_filter='["isolated/plugins/ldap"]' enos scenario launch plugin

# Run specific test names
ENOS_VAR_blackbox_test_filter='["TestKVSecrets"]' enos scenario launch smoke_sdk
ENOS_VAR_blackbox_test_filter='["TestLDAPDynamicRoles"]' enos scenario launch plugin
```

### Key Files
- Scenarios: `enos/enos-scenario-*.hcl`
- Modules: `enos/modules/*/main.tf`
- Module registry: `enos/enos-modules.hcl`
- Globals: `enos/enos-globals.hcl`
- Variables: `enos/enos-variables.hcl`
- Local overrides: `enos/enos-local.vars.hcl` (gitignored)
