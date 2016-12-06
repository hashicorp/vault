set :base_url, "https://www.vaultproject.io/"

activate :hashicorp do |h|
  h.name         = "vault"
  h.version      = "0.6.3"
  h.github_slug  = "hashicorp/vault"
  h.website_root = "website"
end
