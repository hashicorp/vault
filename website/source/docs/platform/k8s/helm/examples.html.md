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
  image: "vault:1.3.0"

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

This example can be used to set up a single server Vault cluster using TLS.

1. Create key & certificate using Kubernetes CA
2. Store key & cert into [Kubernetes secrets store](https://kubernetes.io/docs/concepts/configuration/secret/)
3. Configure helm chart to use Kubernetes secret from step 2

### 1. Create key & certificate using Kubernetes CA

There are three variables that will be used in this example.

```bash
# SERVICE is the name of the Vault service in Kubernetes.
# It does not have to match the actual running service, though it may help for consistency.
SERVICE=vault-server-tls

# NAMESPACE where the Vault service is running.
NAMESPACE=vault-namespace

# SECRET_NAME to create in the Kubernetes secrets store.
SECRET_NAME=vault-server-tls

# TMPDIR is a temporary working directory.
TMPDIR=/tmp
```

1. Create a key for Kubernetes to sign.

    ```bash
    openssl genrsa -out ${TMPDIR}/vault.key 2048
    ```

2. Create a Certificate Signing Request (CSR).

    1. Create a file `${TMPDIR}/csr.conf` with the following contents:

        ```
        [req]
        req_extensions = v3_req
        distinguished_name = req_distinguished_name
        [req_distinguished_name]
        [ v3_req ]
        basicConstraints = CA:FALSE
        keyUsage = nonRepudiation, digitalSignature, keyEncipherment
        extendedKeyUsage = serverAuth
        subjectAltName = @alt_names
        [alt_names]
        DNS.1 = ${SERVICE}
        DNS.2 = ${SERVICE}.${NAMESPACE}
        DNS.3 = ${SERVICE}.${NAMESPACE}.svc
		DNS.4 = ${SERVICE}.${NAMESPACE}.svc.cluster.local
        IP.1 = 127.0.0.1
        EOF
        ```

    2. Create a CSR.

        ```bash
        openssl req -new -key ${TMPDIR}/vault.key -subj "/CN=${SERVICE}.${NAMESPACE}.svc" -out ${TMPDIR}/server.csr -config ${TMPDIR}/csr.conf
        ```

3. Create the certificate

    1. Create a file `${TMPDIR/csr.yaml` with the following contents:

        ```yaml
        apiVersion: certificates.k8s.io/v1beta1
        kind: CertificateSigningRequest
        metadata:
          name: ${CSR_NAME}
        spec:
          groups:
          - system:authenticated
          request: $(cat ${TMPDIR}/server.csr | base64 | tr -d '\n')
          usages:
          - digital signature
          - key encipherment
          - server auth
       ```
       -> `CSR_NAME` can be any name you want. It's the name of the CSR as seen by Kubernetes

    2. Send the CSR to Kubernetes.

        ```bash
        kubectl create -f ${TMPDIR}/csr.yaml
        ```
        -> If this process is automated, you may need to wait to ensure the CSR has been received and stored:
        `kubectl get csr ${CSR_NAME}`

    3. Approve the CSR in Kubernetes.

        ```bash
        kubectl certificate approve ${CSR_NAME}
        ```

### 2. Store key, cert, and Kubernetes CA into Kubernetes secrets store

1. Retrieve the certificate.

    ```bash
    serverCert=$(kubectl get csr ${csrName} -o jsonpath='{.status.certificate}')
    ```
   -> If this process is automated, you may need to wait to ensure the certificate has been created.
   If it hasn't, this will return an empty string.

2. Write the certificate out to a file.

    ```bash
    echo "${serverCert}" | openssl base64 -d -A -out ${TMPDIR}/vault.crt
    ```

3. Retrieve Kubernetes CA.

    ```bash
    kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}' | base64 -D > ${TMPDIR}/vault.ca
    ```

3. Store the key, cert, and Kubernetes CA into Kubernetes secrets.

    ```bash
    kubectl create secret generic ${SECRET_NAME} \
            --namespace ${NAMESPACE} \
            --from-file=vault.key=${TMPDIR}/vault.key \
            --from-file=vault.crt=${TMPDIR}/vault.crt \
            --from-file=vault.ca=${TMPDIR}/vault.ca
    ```



## Helm Configuration

The below `custom-values.yaml` can be used to set up a single server Vault cluster using TLS.
This assumes that a Kubernetes `secret` exists with the server certificate, key and
certificate authority:

```yaml
global:
  tlsDisable: false

server:
  extraEnvironmentVars:
    VAULT_CACERT: /vault/userconfig/vault-server-tls/vault.ca

  extraVolumes:
  - type: secret
    name: vault-server-tls # Matches the ${SECRET_NAME} from above

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
```

## Standalone Server with Audit Storage

The below `values.yaml` can be used to set up a single server Vault cluster with
auditing enabled.

```yaml
global:
  enabled: true
  image: "vault:1.3.0"

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
  image: "vault:1.3.0"

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
