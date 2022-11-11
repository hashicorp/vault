
terraform_cli "default" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null

  credentials "app.terraform.io" {
    token = var.tfc_api_token
  }
}
