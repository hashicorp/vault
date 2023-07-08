# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

storage "inmem" {}
listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
  custom_response_headers  {
    "default" = {
       "Strict-Transport-Security" = ["max-age=1","domains"],
       "Content-Security-Policy" = ["default-src 'others'"],
       "X-Vault-Ignored" = ["ignored"],
       "X-Custom-Header" = ["Custom header value default"],
     }
     "307" = {
       "X-Custom-Header" = ["Custom header value 307"],
     }
     "3xx" = {
       "X-Vault-Ignored-3xx" = ["Ignored 3xx"],
       "X-Custom-Header" = ["Custom header value 3xx"]
     }
     "200" = {
       "someheader-200" = ["200"],
       "X-Custom-Header" = ["Custom header value 200"]
     }
     "2xx" = {
       "X-Custom-Header" = ["Custom header value 2xx"]
     }
     "400" = {
       "someheader-400" = ["400"]
     }
  }
}
disable_mlock = true
