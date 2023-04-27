# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

storage "inmem" {}
listener "tcp" {
    address = "127.0.0.1:8200"
    tls_disable = true
}
cloud {
    resource_id = "organization/bc58b3d0-2eab-4ab8-abf4-f61d3c9975ff/project/1c78e888-2142-4000-8918-f933bbbc7690/hashicorp.example.resource/example"
    client_id = "J2TtcSYOyPUkPV2z0mSyDtvitxLVjJmu"
    client_secret = "N9JtHZyOnHrIvJZs82pqa54vd4jnkyU3xCcqhFXuQKJZZuxqxxbP1xCfBZVB82vY"
}
disable_mlock = true