terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

resource "enos_local_exec" "get_version" {
  scripts = ["${path.module}/scripts/version.sh"]
}

output "version" {
  value = trimspace(enos_local_exec.get_version.stdout)
}
