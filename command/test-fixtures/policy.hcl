path "secret/foo" {
  policy = "write"
}

path "secret/bar/*" {
  capabilities = ["create", "read", "update"]
}
