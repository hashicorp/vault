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

- The feature require third-party integrations. Whether that be networked
  dependencies like a real Consul backend, a real KMS key to test awskms
  auto-unseal, auto-join discovery using AWS tags, or Cloud hardware KMS's.
- The feature might behave differently under multiple configuration variants
  and therefore should be tested with both combinations, e.g. auto-unseal and
  manual shamir unseal or replication in HA mode with integrated storage or
  Consul storage.
- The scenario requires coordination between multiple targets. For example,
  consider the complex lifecycle event of migrating the seal type or storage,
  or manually triggering a raft disaster scenario by partitioning the network
  between the leader and follower nodes. Or perhaps an auto-pilot upgrade between
  a stable version of Vault and our candidate version.
- The scenario has specific deployment strategy requirements. For example,
  if we want to add a regression test for an issue that only arises when the
  software is deployed in a certain manner.
- The scenario needs to use actual build artifacts that will be promoted
  through the pipeline.

## Requirements
- AWS access. HashiCorp Vault developers should use Doormat.
- Terraform >= 1.7
- Enos >= v0.0.28. You can [download a release](https://github.com/hashicorp/enos/releases/) or
  install it with Homebrew:
  ```shell
  brew tap hashicorp/tap && brew update && brew install hashicorp/tap/enos
  ```
- An SSH keypair in the AWS region you wish to run the scenario. You can use
  Doormat to log in to the AWS console to create or upload an existing keypair.
- A Vault artifact is downloaded from the GHA artifacts when using the `artifact_source:crt` variants, from Artifactory when using `artifact_source:artifactory`, and is built locally from the current branch when using  `artifact_source:local` variant.

## Scenario Variables
In CI, each scenario is executed via Github Actions and has been configured using
environment variable inputs that follow the `ENOS_VAR_varname` pattern.

For local execution you can specify all the required variables using environment
variables, or you can update `enos.vars.hcl` with values and uncomment the lines.

Variables that are required:
* `aws_ssh_keypair_name`
* `aws_ssh_private_key_path`
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
enos scenario run smoke artifact_source:local
enos scenario run upgrade artifact_source:local
# To run the same scenario variants that are run in CI, refer to the scenarios listed
# in json files under .github/enos-run-matrices directory,
# adding `artifact_source:local` to run locally.
enos scenario run smoke backend:consul consul_version:1.12.3 distro:ubuntu seal:awskms artifact_source:local arch:amd64 edition:oss
# Launch an individual scenario but leave infrastructure up after execution
enos scenario launch smoke artifact_source:local
# Check an individual scenario for validity. This is useful during scenario
# authoring and debugging.
enos scenario validate smoke artifact_source:local
# If you've run the tests and desire to see the outputs, such as the URL or
# credentials, you can run the output command to see them. Please note that
# after "run" or destroy there will be no "outputs" as the infrastructure
# will have been destroyed and state cleared.
enos scenario output smoke artifact_source:local
# Explicitly destroy all existing infrastructure
enos scenario destroy smoke artifact_source:local
```

Refer to the [Enos documentation](https://github.com/hashicorp/Enos-Docs)
for further information regarding installation, execution or composing scenarios.

# Scenarios
There are current two scenarios: `smoke` and `upgrade`. Both begin by building Vault
as specified by the selected `artifact_source` variant (see Variants section below for more
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
that is determined by the `artifact_source` variant. After the upgrade, it verifies that
cluster is at the desired version, along with additional verifications.


## Autopilot
The [`autopilot` scenario](./enos-scenario-autopilot.hcl) creates a Vault cluster using
the version specified in `vault_upgrade_initial_release`. It writes test data to the Vault cluster. Next, it creates additional nodes with the candidate version of Vault as determined by the `vault_product_version` variable set.
The module uses AWS auto-join to handle discovery and unseals with auto-unseal
or Shamir depending on the `seal` variant. After the new nodes have joined and been
unsealed, it verifies reading stored data on the new nodes. Autopilot upgrade verification checks the upgrade status is "await-server-removal" and the target version is set to the version of upgraded nodes. This test also verifies the undo_logs status for Vault versions 1.13.x

## Replication
The [`replication` scenario](./enos-scenario-replication.hcl) creates two 3-node Vault clusters and runs following verification steps:

 1. Writes data on the primary cluster
 1. Enables performance replication
 1. Verifies reading stored data from secondary cluster
 1. Verifies initial replication status between both clusters
 1. Replaces the leader node and one standby node on the primary Vault cluster
 1. Verifies updated replication status between both clusters

 This scenario verifies the performance replication status on both clusters to have their connection_status as "connected" and that the secondary cluster has known_primaries cluster addresses updated to the active nodes IP addresses of the primary Vault cluster. This scenario currently works around issues VAULT-12311 and VAULT-12309.  The scenario fails when the primary storage backend is Consul due to issue VAULT-12332

## UI Tests
The [`ui` scenario](./enos-scenario-ui.hcl) creates a Vault cluster (deployed to AWS) using a version 
built from the current checkout of the project. Once the cluster is available the UI acceptance tests
are run in a headless browser.
### Variables
In addition to the required variables that must be set, as described in the [Scenario Variables](#Scenario Variables), 
the `ui` scenario has two optional variables:

**ui_test_filter** - An optional test filter to limit the tests that are run, i.e. `'!enterprise'`.
To set a filter export the variable as follows:
```shell
> export ENOS_VAR_ui_test_filter="some filter"
```
**ui_run_tests** - An optional boolean variable to run or not run the tests. The default value is true. 
Setting this value to false is useful in the case where you want to create a cluster, but run the tests 
manually. The section [Running the Tests](#Running the Tests) describes the different ways to run the
'UI' acceptance tests.

### Running the Tests
The UI tests can be run fully automated or manually.
#### Fully Automated
The following will deploy the cluster, run the tests, and subsequently tear down the cluster: 
```shell
> export ENOS_VAR_ui_test_filter="some filter" # <-- optional
> cd enos
> enos scenario ui run edition:oss
```
#### Manually
The UI tests can be run manually as follows:
```shell
> export ENOS_VAR_ui_test_filter="some filter" # <-- optional
> export ENOS_VAR_ui_run_tests=false
> cd enos
> enos scenario ui launch edition:oss
# once complete the scenario will output a set of environment variables that must be exported. The 
# output will look as follows:
export TEST_FILTER='some filter>' \
export VAULT_ADDR='http://<some ip address>:8200' \
export VAULT_TOKEN='<some token>' \
export VAULT_UNSEAL_KEYS='["<some key>","<some key>","<some key>"]'
# copy and paste the above into the terminal to export the values
> cd ../ui
> yarn test:enos # run headless
# or
> yarn test:enos -s # run manually in a web browser 
# once testing is complete
> cd ../enos
> enos scenario ui destroy edition:oss
```

# Variants
Both scenarios support a matrix of variants. In order to achieve broad coverage while
keeping test run time reasonable, the variants executed by the `enos-run` Github
Actions are tailored to maximize variant distribution per scenario.

## `artifact_source:crt`
This variant is designed for use in Github Actions. The `enos-run.yml` workflow
downloads the artifact built by the `build.yml` workflow, unzips it, and sets the
`vault_bundle_path` to the zip file and the `vault_local_binary_path` to the binary.

## `artifact_source:local`
This variant is for running the Enos scenario locally. It builds the Vault bundle
from the current branch, placing the bundle at the `vault_bundle_path` and the
unzipped Vault binary at the `vault_local_binary_path`.

## `artifact_source:artifactory`
This variant is for running the Enos scenario to test an artifact from Artifactory. It requires following Enos variables to be set:
* `artifactory_username`
* `artifactory_token`
* `aws_ssh_keypair_name`
* `aws_ssh_private_key_path`
* `vault_product_version`
* `vault_revision`

# CI Bootstrap
In order to execute any of the scenarios in this repository, it is first necessary to bootstrap the 
CI AWS account with the required permissions, service quotas and supporting AWS resources. There are 
two Terraform modules which are used for this purpose, [service-user-iam](./ci/service-user-iam) for 
the account permissions, and service quotas and [bootstrap](./ci/bootstrap) for the supporting resources.

**Supported Regions** - enos scenarios are supported in the following regions: 
`"us-east-1", "us-east-2", "us-west-1", "us-west-2"`

## Bootstrap Process
These steps should be followed to bootstrap this repo for enos scenario execution:

### Set up CI service user IAM role and Service Quotas
The service user that is used when executing enos scenarios from any GitHub Action workflow must have 
a properly configured IAM role granting the access required to create resources in AWS. Additionally,
service quotas need to be adjusted to ensure that normal use of the ci account does not cause any
service quotas to be exceeded. The [service-user-iam](./ci/service-user-iam) module contains the IAM 
Policy and Role for that grants this access as well as the service quota increase requests to adjust 
the service quotas. This module should be updated whenever a new AWS resource type is required for a 
scenario or a service quota limit needs to be increased. Since this is persistent and cannot be created 
and destroyed each time a scenario is run, the Terraform state will be managed by Terraform Cloud. 
Here are the steps to configure the GitHub Actions service user:

#### Pre-requisites
- Full access to the CI AWS account is required.

**Notes:**
- For help with access to Terraform Cloud and the CI Account, contact the QT team on Slack (#team-quality) 
  for an invite. After receiving an invite to Terraform Cloud, a personal access token can be created
  by clicking `User Settings` --> `Tokens` --> `Create an API token`.
- Access to the AWS account can be done via Doormat, at: https://doormat.hashicorp.services/.
  - For the vault repo the account is: `vault_ci` and for the vault-enterprise repo, the account is:
    `vault-enterprise_ci`.
  - Access can be requested by clicking: `Cloud Access` --> `AWS` --> `Request Account Access`.

1. **Create the Terraform Cloud Workspace** - The name of the workspace to be created depends on the 
   repository for which it is being created, but the pattern is: `<repository>-ci-service-user-iam`,
   e.g. `vault-ci-service-user-iam`. It is important that the execution mode for the workspace be set 
   to `local`. For help on setting up the workspace, contact the QT team on Slack (#team-quality)


2. **Execute the Terraform module**
```shell
> cd ./enos/ci/service-user-iam
> export TF_WORKSPACE=<repo name>-ci-service-user-iam
> export TF_TOKEN_app_terraform_io=<Terraform Cloud Token>
> export TF_VAR_repository=<repository name>
> terraform init
> terraform plan
> terraform apply -auto-approve
```

### Bootstrap the CI resources
Bootstrapping of the resources in the CI account is accomplished via the GitHub Actions workflow: 
[enos-bootstrap-ci](../.github/workflows/enos-bootstrap-ci.yml). Before this workflow can be run a 
workspace must be created as follows:

1. **Create the Terraform Cloud Workspace** - The name workspace to be created depends on the repository
   for which it is being created, but the pattern is: `<repository>-ci-bootstrap`, e.g.
   `vault-ci-bootstrap`. It is important that the execution mode for the workspace be set to
   `local`. For help on setting up the workspace, contact the QT team on Slack (#team-quality).

Once the workspace has been created, changes to the bootstrap module will automatically be applied via
the GitHub PR workflow. Each time a PR is created for changes to files within that module the module
will be planned via the workflow described above. If the plan is ok and the PR is merged, the module
will automatically be applied via the same workflow.
