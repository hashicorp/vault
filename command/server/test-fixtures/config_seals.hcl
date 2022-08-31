listener "tcp" {
  address = "127.0.0.1:443"
}

backend "consul" {
}

seal "pkcs11" {
  purpose = "many,purposes"
  lib = "/usr/lib/libcklog2.so"
  slot = "0.0"
  pin = "XXXXXXXX"
  key_label = "HASHICORP"
  mechanism = "0x1082"
  hmac_mechanism = "0x0251"
  hmac_key_label = "vault-hsm-hmac-key"
  default_hmac_key_label = "vault-hsm-hmac-key"
  generate_key = "true"
}

seal "pkcs11" {
  purpose = "single"
  disabled = "true"
  lib = "/usr/lib/libcklog2.so"
  slot = "0.0"
  pin = "XXXXXXXX"
  key_label = "HASHICORP"
  mechanism = 0x1082
  hmac_mechanism = 0x0251
  hmac_key_label = "vault-hsm-hmac-key"
  default_hmac_key_label = "vault-hsm-hmac-key"
  generate_key = "true"
}

