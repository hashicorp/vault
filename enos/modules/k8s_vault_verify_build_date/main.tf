
terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  // build date was introduced in Vault 1.11.0
  requires_build_date_check = length(regexall("^(?:1\\.1[^0]|[2-9]|\\d{2,})", var.vault_product_version)) > 0

  pods_to_check = toset([
    for idx in range(var.vault_instance_count) : tostring(idx)
    if local.requires_build_date_check
  ])
}

# Get the date from the vault status command      - status_date
# Format the original status output with ISO-8601 - formatted_date
# Format the original status output with awk      - awk_date
# Compare the formatted outputs                   - date_comparison
resource "enos_remote_exec" "status_date" {
  for_each = local.pods_to_check

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.vault_pods[each.key].name
      namespace         = var.vault_pods[each.key].namespace
    }
  }

  inline = ["${var.vault_bin_path} status -format=json | grep build_date | cut -d \\\" -f 4"]
}

resource "enos_remote_exec" "formatted_date" {
  for_each = local.pods_to_check

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.vault_pods[each.key].name
      namespace         = var.vault_pods[each.key].namespace
    }
  }

  inline = ["date -d \"${enos_remote_exec.status_date[each.key].stdout}\" -D '%Y-%m-%dT%H:%M:%SZ' -I"]
}

resource "enos_local_exec" "awk_date" {
  for_each = local.pods_to_check

  inline = ["echo ${enos_remote_exec.status_date[each.key].stdout} | awk -F\"T\" '{printf $1}'"]
}

resource "enos_local_exec" "date_comparison" {
  for_each = local.pods_to_check

  inline = ["[[ ${enos_local_exec.awk_date[each.key].stdout} == ${enos_remote_exec.formatted_date[each.key].stdout} ]] && echo \"Verification for build date format ${enos_remote_exec.status_date[each.key].stdout} succeeded\" || \"invalid build_date, must be formatted as RFC 3339: ${enos_remote_exec.status_date[each.key].stdout}\""]
}
