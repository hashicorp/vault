storage "inmem" {}
listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
  custom_response_headers  {
    "default" = {
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
listener "tcp" {
  address = "127.0.0.2:8200"
  tls_disable = true
  custom_response_headers  {
    "default" = {
       "Content-Security-Policy" = ["default-src 'others'"],
       "X-Vault-Ignored" = ["ignored"],
       "X-Custom-Header" = ["Custom header value default"],
     }
  }
}
listener "tcp" {
  address = "127.0.0.3:8200"
  tls_disable = true
  custom_response_headers  {
    "2xx" = {
       "X-Custom-Header" = ["Custom header value 2xx"]
    }
  }
}
listener "tcp" {
  address = "127.0.0.4:8200"
  tls_disable = true
}


disable_mlock = true
