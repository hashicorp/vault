---
layout: "api"
page_title: "Related Tools"
sidebar_current: "docs-http-related"
description: |-
  Short list of third-party tools that work with or are related to Vault.
---

# Related Tools

## Hashicorp Tools

* The [Terraform Vault provider](https://www.terraform.io/docs/providers/vault/index.html) can read from, write to, and configure Vault from [HashiCorp Terraform](https://www.terraform.io/)
* [consul-template](https://github.com/hashicorp/consul-template) is a template renderer, notifier, and supervisor for HashiCorp Consul and Vault data
* [envconsul](https://github.com/hashicorp/envconsul) allows you to read and set environmental variables for processes from Consul and Vault data
* The [vault-ssh-helper](https://github.com/hashicorp/vault-ssh-helper) can be used to enable one-time passwords for SSH authentication via Vault

## Third-Party Tools

The following list of tools is maintained by the community of Vault users; HashiCorp has not tested or approved them and makes no claims as to their suitability or security.

* [HashiCorp Vault Jenkins plugin](https://plugins.jenkins.io/hashicorp-vault-plugin) - a Jenkins plugin for injecting Vault secrets into the build environment
* [Spring Vault](http://projects.spring.io/spring-vault/) - a Java Spring project for working with Vault secrets
* [vault-exec](https://github.com/kmanning/vault_exec) - a shell wrapper to execute arbitrary scripts using temporary AWS credentials managed by Vault
* [pouch](https://github.com/tuenti/pouch) - A set of tools to manage provisioning of secrets on hosts based on the AppRole authentication method of Vault
* [vault-aws-creds](https://github.com/jantman/vault-aws-creds) - Python helper to export Vault-provided temporary AWS creds into the environment

Want to add your own project, or one that you use? Additions are welcome via [pull requests](https://github.com/hashicorp/vault/blob/master/website/source/api/relatedtools.html.md).
