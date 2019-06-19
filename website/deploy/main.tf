locals {
  github_parts = ["${split("/", var.github_repo)}"]
  github_full  = "${var.github_repo}"
  github_org   = "${local.github_parts[0]}"
  github_repo  = "${local.github_parts[1]}"
}

/*
-------------------------------------------------------------------
GitHub Resources
-------------------------------------------------------------------
*/

provider "github" {
  organization = "${local.github_org}"
  version      = "~>2.0.0"
}

// Configure the repository with the dynamically created Netlify key.
resource "github_repository_deploy_key" "key" {
  title      = "Netlify"
  repository = "${local.github_repo}"
  key        = "${netlify_deploy_key.key.public_key}"
  read_only  = false
}

// Create a webhook that triggers Netlify builds on push.
resource "github_repository_webhook" "main" {
  repository = "${local.github_repo}"
  events     = ["delete", "push", "pull_request"]

  configuration {
    content_type = "json"
    url          = "https://api.netlify.com/hooks/github"
    insecure_ssl = false
  }

  depends_on = ["netlify_site.main"]
}

/*
-------------------------------------------------------------------
Netlify Resources
-------------------------------------------------------------------
*/

provider "netlify" {
  version = "~> 0.1.0"
}

// A new, unique deploy key for this specific website
resource "netlify_deploy_key" "key" {}

resource "netlify_site" "main" {
  name          = "${var.name}"
  custom_domain = "${var.custom_site_domain}"

  repo {
    repo_branch   = "${var.github_branch}"
    command       = "cd website && bundle && cd assets && npm i && cd .. && middleman build --verbose"
    deploy_key_id = "${netlify_deploy_key.key.id}"
    dir           = "website/build"
    provider      = "github"
    repo_path     = "${local.github_full}"
  }
}
