# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "followers" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault follower hosts"
}

output "chosen_follower" {
  value = {
    0 : try(var.followers[0], null)
  }
}
