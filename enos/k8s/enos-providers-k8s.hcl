# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

provider "enos" "default" {}

provider "helm" "default" {
  kubernetes {
    config_path = abspath(joinpath(path.root, "kubeconfig"))
  }
}
