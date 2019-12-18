---
layout: "docs"
page_title: "Agent Sidecar Injector Installation"
sidebar_current: "docs-platform-k8s-injector-installation"
sidebar_title: "Installation"
description: |-
  The Vault Agent Sidecar Injector can be installed using Vault Helm.
---

# Installing the Agent Injector

The following are the different methods of installing the Agent Injector in
Kubernetes.

~> The Vault Agent Injector requires Vault 1.3.1.

To install the Vault Agent injector, enable the injection feature using
[Helm values](/docs/platform/k8s/helm/configuration.html) and
upgrade the installation using `helm upgrade` for existing installs or
`helm install` for a fresh install.

To install a new instance of Vault and the Vault Agent Injector, run the following:

```bash
helm install --name=vault \
  --set="injector.enabled=true" \
  https://github.com/hashicorp/vault-helm/archive/v0.3.0tar.gz
``` 

Other values in the Helm chart can be used to limit the namespaces the injector
runs in, TLS options and more.

## TLS Options

Admission webhook controllers require TLS to run within Kubernetes.  At this time
the Vault Agent Injector supports two TLS options:

* Auto TLS generation (default)
* Manual TLS

### Auto TLS

By default, the Vault Agent Injector will bootstrap TLS by generating a certificate
authority and creating a certificate/key to be used by the controller.  If using
Vault Helm, the chart will automatically create the neccessary DNS entries for the
controller's service used to verify the certificate.

### Manual TLS

If desired, users can supply their own TLS certificates, key and certificate authority.
The following is required to configure TLS manually:

* Server certificate/key
* Base64 PEM encoded Certificate Authority bundle

For more information on configuring manual TLS, see the [Vault Helm cert values](/docs/platform/k8s/helm/configuration.html#certs).

## Namespace Selector

By default, the Vault Agent Injector will process all namespaces in Kubernetes except
the system namespaces `kube-system` and `kube-public`.  To limit what namespaces
the injector can work in a namespace selector can be defined to match labels attached
to namespaces.

For more information on configuring namespace selection, see the [Vault Helm namespaceSelector value](/docs/platform/k8s/helm/configuration.html#namespaceselector).
