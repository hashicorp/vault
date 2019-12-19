---
layout: "docs"
page_title: "Examples"
sidebar_current: "docs-platform-k8s-examples"
sidebar_title: "Examples"
description: |-
  This section documents configuration options for the Vault Helm chart
---

# Helm Chart Examples

~> **Important Note:** This chart is not compatible with Helm 3. Please use Helm 2 with this chart.

The following are different configuration examples to support a variety of
deployment models.

## Standalone Server with Load Balanced UI

The below `values.yaml` can be used to set up a single server Vault cluster with a LoadBalancer to allow external access to the UI and API.

```yaml
global:
  enabled: true
  image: "vault:1.3.1"

server:
  standalone:
    enabled: true
    config: |
      ui = true

      listener "tcp" {
        tls_disable = 1
        address = "[::]:8200"
        cluster_address = "[::]:8201"
      }
      storage "file" {
        path = "/vault/data"
      }

  service:
    enabled: true

  dataStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce

ui:
  enabled: true
  serviceType: LoadBalancer
```

## Standalone Server with TLS

The below `values.yaml` can be used to set up a single server Vault cluster using TLS.
This assumes that a Kubernetes `secret` exists with the server certificate, key and
certificate authority:

```yaml
global:
  enabled: true
  image: "vault:1.3.1"
  tlsDisable: false

server:
  extraVolumes:
  - type: secret
    name: vault-server-tls

  standalone:
    enabled: true
    config: |
      listener "tcp" {
        address = "[::]:8200"
        cluster_address = "[::]:8201"
        tls_cert_file = "/vault/userconfig/vault-server-tls/vault.crt"
        tls_key_file  = "/vault/userconfig/vault-server-tls/vault.key"
        tls_client_ca_file = "/vault/userconfig/vault-server-tls/vault.ca"
      }

      storage "file" {
        path = "/vault/data"
      }

  service:
    enabled: true

  dataStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce
```

## Standalone Server with Audit Storage

The below `values.yaml` can be used to set up a single server Vault cluster with
auditing enabled.

```yaml
global:
  enabled: true
  image: "vault:1.3.1"

server:
  standalone:
    enabled: true
    config: |
      listener "tcp" {
        tls_disable = true
        address = "[::]:8200"
        cluster_address = "[::]:8201"
      }

      storage "file" {
        path = "/vault/data"
      }

  service:
    enabled: true

  dataStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce

  auditStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce
```

After Vault has been deployed, initialized and unsealed, auditing can be enabled
by running the following command against the Vault pod:

```bash
$ kubectl exec -ti <POD NAME> --  vault audit enable file file_path=/vault/audit/vault_audit.log
```

## Highly Available Vault Cluster with Consul

The below `values.yaml` can be used to set up a five server Vault cluster using
Consul as a highly available storage backend, Google Cloud KMS for Auto Unseal.

```yaml
global:
  enabled: true
  image: "vault:1.3.1"

server:
  extraEnvironmentVars:
    GOOGLE_REGION: global
    GOOGLE_PROJECT: myproject
    GOOGLE_APPLICATION_CREDENTIALS: /vault/userconfig/my-gcp-iam/myproject-creds.json

  extraVolumes: []
    - type: secret
      name: my-gcp-iam

  affinity: |
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchLabels:
              app: {{ template "vault.name" . }}
              release: "{{ .Release.Name }}"
              component: server
          topologyKey: kubernetes.io/hostname

  service:
    enabled: true

  ha:
    enabled: true
    replicas: 5

    config: |
      ui = true

      listener "tcp" {
        tls_disable = 1
        address = "[::]:8200"
        cluster_address = "[::]:8201"
      }

      storage "consul" {
        path = "vault"
        address = "HOST_IP:8500"
      }

      seal "gcpckms" {
         project     = "myproject"
         region      = "global"
         key_ring    = "vault-unseal-kr"
         crypto_key  = "vault-unseal-key"
      }
```
