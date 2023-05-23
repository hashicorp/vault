auto_auth {
  method {
    type = "token_file"
    config {
      token_file_path = "/home/username/.vault-token"
    }
  }
}

env_template "MY_PASSWORD" {
  source = "/path/to/template/file"
}

exec {
  command = ["/path/to/my/app", "arg1", "arg2"]
}
