// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

sample "ce_default_linux_amd64_ent_docker" {
  subset "k8s" {
    matrix {
      repo    = ["docker", "ecr"]
      edition = ["ce"]
    }
  }
}

sample "ce_default_linux_arm64_ce_docker" {
  subset "k8s" {
    matrix {
      repo    = ["docker", "ecr"]
      edition = ["ce"]
    }
  }
}

sample "ce_ubi_linux_amd64_ce_redhat" {
  subset "k8s" {
    matrix {
      repo    = ["quay"]
      edition = ["ce"]
    }
  }
}

sample "ce_ubi_linux_arm64_ce_redhat" {
  subset "k8s" {
    matrix {
      repo    = ["quay"]
      edition = ["ce"]
    }
  }
}
