---
layout: "docs"
page_title: "Agent Sidecar Injector Installation"
sidebar_current: "docs-platform-k8s-injector-installation"
sidebar_title: "Installation"
description: |-
  The Vault Agent Sidecar Injector can be installed in two ways: using Vault Helm or manually.
---

# Installing the Agent Injector

The following are the different methods of installing the Agent Injector in
Kubernetes.

## Using Vault Helm

To install the Vault Agent injector, enable the injection feature using
[Helm values](docs/platform/k8s/helm.html#configuration-values-) and
upgrade the installation using `helm upgrade` for existing installs or
`helm install` for a fresh install.

```bash
export CA_BUNDLE=$(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')

helm install --name=vault \
  --set="injector.enabled=true" \
  --set="injector.tls.caBundle=${CA_BUNDLE?}" \
  https://github.com/hashicorp/vault-helm/archive/v0.3.0tar.gz
``` 

Other values in the Helm chart can be used to limit the namespaces the injector
runs in, TLS options and more.

## Manual Installation

TODO

## TLS Options

TODO
