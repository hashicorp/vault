env_template "MY_DATABASE_USER" {
  contents = "{{ with secret \"secret/db-secret\" }}{{ .Data.data.user }}{{ end }}"
}

exec {
  command = ["/path/to/my/app", "arg1", "arg2"]
}
