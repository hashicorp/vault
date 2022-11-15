# Enos

Enos is an quality testing framework that allows composing and executing quality
requirement scenarios as code. For Vault, it is currently used to perform
infrastructure integration testing using the artifacts that are created as part
of the `build` workflow. While intended to be executed via Github Actions using
the results of the `build` workflow, scenarios are also executable from a developer
machine that has the requisite dependencies and configuration.

Refer to the [Enos documentation](https://github.com/hashicorp/Enos-Docs)
for further information regarding installation, execution or composing Enos scenarios.

## When to use Enos
Determining whether to use `vault.NewTestCluster()` or Enos for testing a feature
or scenario is ultimately up to the author. Sometimes one, the other, or both
might be appropriate depending on the requirements. Generally, `vault.NewTestCluster()`
is going to give you faster feedback and execution time, whereas Enos is going
to give you a real-world execution and validation of the requirement. Consider
the following cases as examples of when one might opt for an Enos scenario:

* The feature require third-party integrations. Whether that be networked
  dependencies like a real Consul backend, a real KMS key to test awskms
  auto-unseal, auto-join discovery using AWS tags, or Cloud hardware KMS's.
* The feature might behave differently under multiple configuration variants
  and therefore should be tested with both combinations, e.g. auto-unseal and
  manual shamir unseal or replication in HA mode with integrated storage or
  Consul storage.
* The scenario requires coordination between multiple targets. For example,
  consider the complex lifecycle event of migrating the seal type or storage,
  or manually triggering a raft disaster scenario by partitioning the network
  between the leader and follower nodes. Or perhaps an auto-pilot upgrade between
  a stable version of Vault and our candidate version.
* The scenario has specific deployment strategy requirements. For example,
  if we want to add a regression test for an issue that only arises when the
  software is deployed in a certain manner.
* The scenario needs to use actual build artifacts that will be promoted
  through the pipeline.

## Requirements
* AWS access. HashiCorp Vault developers should use Doormat.
* Terraform >= 1.2
* Enos >= v0.0.10. You can [install it from a release channel](https://github.com/hashicorp/Enos-Docs/blob/main/installation.md).
* Access to the QTI org in Terraform Cloud. HashiCorp Vault developers can
  access a shared token in 1Password or request their own in #team-quality on
  Slack.
* An SSH keypair in the AWS region you wish to run the scenario. You can use
  Doormat to log in to the AWS console to create or upload an existing keypair.
* A Vault install bundle downloaded from releases.hashicorp.com or Artifactory
  when using the `builder:crt` variants. When using the `builder:local` variants
  Enos will build a Vault bundle from the current branch for you.

## Scenario Variables
In CI, each scenario is executed via Github Actions and has been configured using
environment variable inputs that follow the `ENOS_VAR_varname` pattern.

For local execution you can specify all the required variables using environment
variables, or you can update `enos.vars.hcl` with values and uncomment the lines.

Variables that are required:
* `aws_ssh_keypair_name`
* `aws_ssh_private_key_path`
* `tfc_api_token`
* `vault_bundle_path`
* `vault_license_path` (only required for non-OSS editions)

See [enos.vars.hcl](./enos.vars.hcl) or [enos-variables.hcl](./enos-variables.hcl)
for further descriptions of the variables.

## Executing Scenarios
From the `enos` directory:

```bash
# List all available scenarios
enos scenario list
# Run the smoke or upgrade scenario with an artifact that is built locally. Make sure
# the local machine has been configured as detailed in the requirements
# section. This will execute the scenario and clean up any resources if successful.
enos scenario run smoke builder:local
enos scenario run upgrade builder:local
# To run the same scenario variants that are run in CI, refer to the scenarios listed
# in .github/workflows/enos-run.yml under `jobs.enos.strategy.matrix.include`,
# adding `builder:local` to run locally.
enos scenario run smoke backend:consul consul_version:1.12.3 distro:ubuntu seal:awskms builder:local arch:amd64 edition:oss
# Launch an individual scenario but leave infrastructure up after execution
enos scenario launch smoke builder:local
# Check an individual scenario for validity. This is useful during scenario
# authoring and debugging.
enos scenario validate smoke builder:local
# If you've run the tests and desire to see the outputs, such as the URL or
# credentials, you can run the output command to see them. Please note that
# after "run" or destroy there will be no "outputs" as the infrastructure
# will have been destroyed and state cleared.
enos scenario output smoke builder:local
# Explicitly destroy all existing infrastructure
enos scenario destroy smoke builder:local
```

Refer to the [Enos documentation](https://github.com/hashicorp/Enos-Docs)
for further information regarding installation, execution or composing scenarios.

# Scenarios
There are current two scenarios: `smoke` and `upgrade`. Both begin by building Vault
as specified by the selected `builder` variant (see Variants section below for more
information).

## Smoke
The [`smoke` scenario](./enos-scenario-smoke.hcl) creates a Vault cluster using
the version from the current branch (either in CI or locally), with the backend
specified by the `backend` variant (`raft` or `consul`). Next, it unseals with the
appropriate method (`awskms` or `shamir`) and performs different verifications
depending on the backend and seal type.

## Upgrade
The [`upgrade` scenario](./enos-scenario-upgrade.hcl) creates a Vault cluster using
the version specified in `vault_upgrade_initial_release`, with the backend specified
by the `backend` variant (`raft` or `consul`). Next, it upgrades the Vault binary
that is determined by the `builder` variant. After the upgrade, it verifies that
cluster is at the desired version, along with additional verifications.


## Autopilot
The [`autopilot` scenario](./enos-scenario-autopilot.hcl) creates a Vault cluster using
the version specified in `vault_upgrade_initial_release`. Next, it creates additional
nodes with the candiate version of Vault as determined by the `builder` variant.
The module uses AWS auto-join to handle discovery and unseals with auto-unseal
or Shamir depending on the `seal` variant. After the new nodes have joined and been
unsealed, it waits for Autopilot to upgrade the new nodes and demote the old nodes.

# Variants
Both scenarios support a matrix of variants. In order to achieve broad coverage while
keeping test run time reasonable, the variants executed by the `enos-run` Github
Actions are tailored to maximize variant distribution per scenario.

## `builder:crt`
This variant is designed for use in Github Actions. The `enos-run.yml` workflow
downloads the artifact built by the `build.yml` workflow, unzips it, and sets the
`vault_bundle_path` to the zip file and the `vault_local_binary_path` to the binary.

## `builder:local`
This variant is for running the Enos scenario locally. It builds the Vault bundle
from the current branch, placing the bundle at the `vault_bundle_path` and the
unzipped Vault binary at the `vault_local_binary_path`.

# CI Bootstrap
In order to execute any of the scenarios in this repository, it is first necessary to bootstrap the 
CI AWS account with the required supporting AWS resources. At this time, the only resource that is 
required is an EC2 ssh key pair. The scenarios [bootstrap_ci](./ci/enos-scenario-bootstrap-ci.hcl) and 
[bootstrap_workspaces](./ci/enos-scenario-ci-bootstrap-workspaces.hcl) have been created to simplify
the process of bootstrapping the CI environment.

**Supported Regions** - enos scenarios are supported in the following regions: 
`"us-east-1", "us-east-2", "us-west-1", "us-west-2"`

## Bootstrap Process
These steps should be followed to bootstrap this repo for enos scenario execution:

1. **Setup the root Workspace** - In Terraform Cloud manually create a workspace in the `hashicorp-qti` 
   organization named `vault-ci-bootstrap`. This workspace will be used as the backend for the 
   `bootstrap_workspaces` scenario. When creating the workspace choose the execution mode `local`.


2. Get the enos ci ssh public key from a member of the QT team.


3. **Create the CI Bootstrap Workspaces** - Each region will have its own workspace and state, to 
   setup these workspaces the scenario `bootstrap_workspaces` should be executed:

```bash
> export ENOS_VAR_tfc_api_token=<tfc token>
> enos scenaio launch --no-reconfigure bootstrap_workspaces
```

4.**Bootstrap Vault CI** - Execute the following to boostrap the Vault CI account for all regions:

```bash
> export ENOS_VAR_aws_ssh_public_key_path=<path to the enos-ci-ssh-key file from 2>
> export ENOS_VAR_tfc_api_token=<tfc token>
> enos scneario launch --no-configure bootstrap_ci
```
