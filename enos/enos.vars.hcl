# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# artifactory_username is the username to use when testing an artifact stored in artfactory.
# artifactory_username = "yourname@hashicorp.com"

# artifactory_token is the token to use when authenticating to artifactory.
# artifactory_token = "yourtoken"

# artifactory_host is the artifactory host to search for vault artifacts.
# artifactory_host = "https://artifactory.hashicorp.engineering/artifactory"

# artifactory_repo is the artifactory repo to search for vault artifacts.
# artifactory_repo = "hashicorp-crt-stable-local*"

# aws_region is the AWS region where we'll create infrastructure
# for the smoke scenario
# aws_region = "us-east-1"

# aws_ssh_keypair_name is the AWS keypair to use for SSH
# aws_ssh_keypair_name = "enos-ci-ssh-key"

# aws_ssh_private_key_path is the path to the AWS keypair private key
# aws_ssh_private_key_path = "./support/private_key.pem"

# backend_license_path is the license for the backend if applicable (Consul Enterprise)".
# backend_license_path = "./support/consul.hclic"

# backend_log_level is the server log level for the backend. Supported values include 'trace',
# 'debug', 'info', 'warn', 'error'"
# backend_log_level = "trace"

# backend_instance_type is the instance type to use for the Vault backend. Must support arm64
# backend_instance_type = "t4g.small"

# project_name is the description of the project. It will often be used to tag infrastructure
# resources.
# project_name = "vault-enos-integration"

# rhel_distro_version is the version of RHEL to use for "distro:rhel" variants.
# rhel_distro_version = "9.1" // or "8.8"

# tags are a map of tags that will be applied to infrastructure resources that
# support tagging.
# tags = { "Project Name" : "Vault", "Something Cool" : "Value" }

# terraform_plugin_cache_dir is the directory to cache Terraform modules and providers.
# It must exist.
# terraform_plugin_cache_dir = "/Users/<user>/.terraform/plugin-cache-dir

# tfc_api_token is the Terraform Cloud QTI Organization API token. We need this
# to download the enos Terraform provider and the enos Terraform modules.
# tfc_api_token = "XXXXX.atlasv1.XXXXX..."

# ui_test_filter is the test filter to limit the ui tests to execute for the ui scenario. It will
# be appended to the ember test command as '-f=\"<filter>\"'.
# ui_test_filter = "sometest"

# ui_run_tests sets whether to run the UI tests or not for the ui scenario. If set to false a
# cluster will be created but no tests will be run.
# ui_run_tests = true

# ubuntu_distro_version is the version of ubuntu to use for "distro:ubuntu" variants
# ubuntu_distro_version = "22.04" // or "20.04", "18.04"

# vault_artifact_path is the path to CRT generated or local vault.zip bundle. When
# using the "builder:local" variant a bundle will be built from the current branch.
# In CI it will use the output of the build workflow.
# vault_artifact_path = "./dist/vault.zip"

# vault_artifact_type is the type of Vault artifact to use when installing Vault from artifactory.
# It should be 'package' for .deb or # .rpm package and 'bundle' for .zip bundles"
# vault_artifact_type = "bundle"

# vault_autopilot_initial_release is the version of Vault to deploy before doing an autopilot upgrade
# to the test artifact.
# vault_autopilot_initial_release = {
#     edition = "ent"
#     version = "1.11.0"
#   }
# }

# vault_build_date is the build date for Vault artifact. Some validations will require the binary build
# date to match"
# vault_build_date = "2023-07-07T14:06:37Z" // make ci-get-date for example

# vault_enable_file_audit_device sets whether or not to enable the 'file' audit device. It true it
# will be enabled at the path /var/log/vault_audit.log
# vault_enable_file_audit_device = true

# vault_install_dir is the directory where the vault binary will be installed on
# the remote machines.
# vault_install_dir = "/opt/vault/bin"

# vault_local_binary_path is the path of the local binary that we're upgrading to.
# vault_local_binary_path = "./support/vault"

# vault_instance_type is the instance type to use for the Vault backend
# vault_instance_type = "t3.small"

# vault_instance_count is how many instances to create for the Vault cluster.
# vault_instance_count = 3

# vault_license_path is the path to a valid Vault enterprise edition license.
# This is only required for non-oss editions"
# vault_license_path = "./support/vault.hclic"

# vault_local_build_tags override the build tags we pass to the Go compiler for builder:local variants.
# vault_local_build_tags = ["ui", "ent"]

# vault_log_level is the server log level for Vault logs. Supported values (in order of detail) are
# trace, debug, info, warn, and err."
# vault_log_level = "trace"

# vault_product_version is the version of Vault we are testing. Some validations will expect the vault
# binary and cluster to report this version.
# vault_product_version = "1.15.0"

# vault_upgrade_initial_release is the Vault release to deploy before upgrading.

# vault_revision is the git sha of Vault artifact we are testing. Some validations will expect the vault
# binary and cluster to report this revision.
# vault_revision = "df733361af26f8bb29b63704168bbc5ab8d083de"

# vault_upgrade_initial_release is the Vault release to deploy before doing an in-place upgrade.
# vault_upgrade_initial_release = {
#     edition = "oss"
#     // Vault 1.10.5 has a known issue with retry_join.
#     version = "1.10.4"
#   }
# }
