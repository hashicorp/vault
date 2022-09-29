terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

resource "enos_local_exec" "get_revision" {
  inline = ["git rev-parse HEAD"]
}

output "revision" {
  value = trimspace(enos_local_exec.get_revision.stdout)
}
