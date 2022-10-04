terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

resource "enos_local_exec" "get_build_date" {
  scripts = ["${path.module}/scripts/build_date.sh"]
}

output "build_date" {
  value = trimspace(enos_local_exec.get_build_date.stdout)
}

resource "enos_local_exec" "get_version" {
  scripts = ["${path.module}/scripts/version.sh"]
}

output "version" {
  value = trimspace(enos_local_exec.get_version.stdout)
}

resource "enos_local_exec" "get_revision" {
  inline = ["git rev-parse HEAD"]
}

output "revision" {
  value = trimspace(enos_local_exec.get_revision.stdout)
}
