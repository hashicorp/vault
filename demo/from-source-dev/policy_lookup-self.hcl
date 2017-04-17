path "auth/token/lookup-self" {
  policy = "read"
}

path "auth/userpass/users/*" {
  policy = "read"
}
