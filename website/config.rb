#-------------------------------------------------------------------------
# Configure Middleman
#-------------------------------------------------------------------------

set :base_url, "https://www.vaultproject.io/"

activate :hashicorp do |h|
  h.version         = ENV["VAULT_VERSION"]
  h.bintray_enabled = ENV["BINTRAY_ENABLED"]
  h.bintray_repo    = "mitchellh/vault"
  h.bintray_user    = "mitchellh"
  h.bintray_key     = ENV["BINTRAY_API_KEY"]

  h.minify_javascript = false
end
