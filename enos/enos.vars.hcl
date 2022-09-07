# aws_region is the AWS region where we'll create infrastructure
# for the smoke scenario
# aws_region = "us-west-1"

# aws_ssh_keypair_name is the AWS keypair to use for SSH
# aws_ssh_keypair_name = "enos-ci-ssh-key"

# aws_ssh_private_key_path is the path to the AWS keypair private key
# aws_ssh_private_key_path = "./support/private_key.pem"

# backend_instance_type is the instance type to use for the Vault backend
# backend_instance_type = "t3.small"

# tags are a map of tags that will be applied to infrastructure resources that
# support tagging.
# tags = { "Project Name" : "Vault", "Something Cool" : "Value" }

# terraform_plugin_cache_dir is the directory to cache Terraform modules and providers.
# It must exist.
# terraform_plugin_cache_dir = "/Users/<user>/.terraform/plugin-cache-dir

# tfc_api_token is the Terraform Cloud QTI Organization API token. We need this
# to download the enos Terraform provider and the enos Terraform modules.
# tfc_api_token = "XXXXX.atlasv1.XXXXX..."

# vault_bundle_path is the path to CRT generated or local vault.zip bundle. When
# using the "builder:local" variant a bundle will be built from the current branch.
# In CI it will use the output of the build workflow.
# vault_bundle_path = "./dist/vault.zip"

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

# vault_upgrade_initial_release is the Vault release to deploy before upgrading.
