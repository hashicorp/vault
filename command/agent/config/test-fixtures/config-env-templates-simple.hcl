env_template "MY_DATABASE_USER" {
  contents = "{{ with secret \"secret/db-secret\" }}{{ .Data.data.user }}{{ end }}"
  group    = "database_secret"
}

exec {
  command               = "/path/to/my/app"
  args                  = ["arg1", "arg2"]
}