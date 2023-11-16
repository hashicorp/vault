# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "enos" "default" {}

provider "helm" "default" {
  kubernetes {
    config_path = abspath(joinpath(path.root, "kubeconfig"))
  }
}
