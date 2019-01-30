---
layout: "docs"
page_title: "Vault Agent Auto-Auth Kubernetes Method"
sidebar_title: "Kubernetes"
sidebar_current: "docs-agent-autoauth-methods-kubernetes"
description: |-
  Kubernetes Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth Kubernetes Method 

The `kubernetes` method reads in a Kubernetes service account token from the
running pod (via `/var/run/secrets/kubernetes.io/serviceaccount/token`) and
sends it to the [Kubernetes Auth
method](https://www.vaultproject.io/docs/auth/kubernetes.html).

## Configuration

- `role` `(string: required)` - The role to authenticate against on Vault
- `token_path` `(string: optional)` - The file path to a custom JWT token to use
for authentication. If omitted, the default service account token path is used.
