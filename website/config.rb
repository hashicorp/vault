#-------------------------------------------------------------------------
# Configure Middleman
#-------------------------------------------------------------------------

set :base_url, "https://www.vault.io/"

activate :hashicorp do |h|
  h.version      = ENV['VAULT_VERSION']
  h.bintray_repo = 'mitchellh/vault'
  h.bintray_user = 'mitchellh'
  h.bintray_key  = ENV['BINTRAY_API_KEY']
end
