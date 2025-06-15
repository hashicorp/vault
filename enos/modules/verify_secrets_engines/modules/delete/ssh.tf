# Delete SSH CA role
resource "enos_remote_exec" "ssh_delete_ca_role" {

  environment = {
    REQPATH           = "ssh/roles/${var.create_state.ssh.ca_role_name}"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/delete-payload.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Delete SSH OTP role
resource "enos_remote_exec" "ssh_delete_otp_role" {
  for_each = var.hosts

  environment = {
    REQPATH           = "ssh/roles/${var.create_state.ssh.otp_role_name}"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/delete-payload.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}