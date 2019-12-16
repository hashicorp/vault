---
layout: "docs"
page_title: "Examples"
sidebar_current: "docs-platform-k8s-agent-injector-examples"
sidebar_title: "Examples"
description: |-
  This section documents examples of using the Vault Agent Injector.
---

# Vault Agent Injector Examples

The following are different configuration examples to support a variety of
deployment models.

## Patching Existing Pods

To patch existing pods, a Kubernetes patch can be applied to add the required annoations 
to pods.  When applying a patch, the pods will be rescheduled.

The following example patches a deployment.  First, create the patch:

```bash
cat <<EOF >> ./patch.yaml
spec:
  template:
    metadata:
      annotations:
        vault.hashicorp.com/agent-inject: "true"
        vault.hashicorp.com/agent-inject-status: "update"
        vault.hashicorp.com/agent-inject-secret-db-creds: "database/creds/db-app"
        vault.hashicorp.com/agent-inject-template-db-creds: |
          {{- with secret "database/creds/db-app" -}}
          postgres://{{ .Data.username }}:{{ .Data.password }}@postgres:5432/appdb?sslmode=disable
          {{- end }}
        vault.hashicorp.com/role: "db-app"
        vault.hashicorp.com/ca-cert: "/vault/tls/ca.crt"
        vault.hashicorp.com/client-cert: "/vault/tls/client.crt"
        vault.hashicorp.com/client-key: "/vault/tls/client.key"
        vault.hashicorp.com/tls-secret: "vault-tls-client"
EOF
```

Next, apply the patch:

```bash
kubectl patch deployment <MY DEPLOYMENT> --patch "$(cat patch.yaml)"
```

## Deployment, Statefulsets, etc.

The annotations for configuring Vault Agent injection must be on the pod 
specification. Since higher level resources such as Deployments wrap pod 
specification templates, Vault Agent Injector can be used with all of these 
higher level constructs, too.

An example Deployment below shows how to enable Vault Agent injection:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-example-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: app-example
  template:
    metadata:
      labels:
        app: app-example
      annotations:
        vault.hashicorp.com/agent-inject: "true"
        vault.hashicorp.com/agent-inject-secret-db-creds: "database/creds/db-app"
        vault.hashicorp.com/agent-inject-template-db-creds: |
          {{- with secret "database/creds/db-app" -}}
          postgres://{{ .Data.username }}:{{ .Data.password }}@postgres:5432/appdb?sslmode=disable
          {{- end }}
        vault.hashicorp.com/role: "db-app"
        vault.hashicorp.com/ca-cert: "/vault/tls/ca.crt"
        vault.hashicorp.com/client-cert: "/vault/tls/client.crt"
        vault.hashicorp.com/client-key: "/vault/tls/client.key"
        vault.hashicorp.com/tls-secret: "vault-tls-client"
    spec:
      containers:
        - name:app 
          image: "app:1.0.0"
      serviceAccountName: app-example
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-example
```

~> A common mistake is to set the annotation on the Deployment or other resource. 
  Ensure that the injector annotations are specified on the pod specification 
  template as shown above.

## ConfigMap Example

The following example creates a deployment that mounts a Kubernetes ConfigMap 
containing Vault Agent configuration files.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-example-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: app-example
  template:
    metadata:
      labels:
        app: app-example
      annotations:
        vault.hashicorp.com/agent-inject: "true"
        vault.hashicorp.com/agent-configmap: "my-configmap"
        vault.hashicorp.com/tls-secret: "vault-tls-client"
    spec:
      containers:
        - name:app 
          image: "app:1.0.0"
      serviceAccountName: app-example
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-example
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap
agent-config
    app: app-example
data:
  config.hcl: |
    "auto_auth" = {
      "method" = {
        "config" = {
          "role" = "db-app"
        }
        "type" = "kubernetes"
      }

      "sink" = {
        "config" = {
          "path" = "/home/vault/.token"
        }

        "type" = "file"
      }
    }

    "exit_after_auth" = false
    "pid_file" = "/home/vault/.pid"

    "template" = {
      "contents" = "{{- with secret "database/creds/db-app" -}}postgres://{{ .Data.username }}:{{ .Data.password }}@postgres:5432/mydb?sslmode=disable{{- end }}"
      "destination" = "/vault/secrets/db-creds"
    }

    "vault" = {
      "address" = "https://vault.demo.svc.cluster.local:8200"
      "ca_cert" = "/vault/tls/ca.crt"
      "client_cert" = "/vault/tls/client.crt"
      "client_key" = "/vault/tls/client.key"
    }
  config-init.hcl: |
    "auto_auth" = {
      "method" = {
        "config" = {
          "role" = "db-app"
        }
        "type" = "kubernetes"
      }

      "sink" = {
        "config" = {
          "path" = "/home/vault/.token"
        }

        "type" = "file"
      }
    }

    "exit_after_auth" = true
    "pid_file" = "/home/vault/.pid"

    "template" = {
      "contents" = "{{- with secret "database/creds/db-app" -}}postgres://{{ .Data.username }}:{{ .Data.password }}@postgres:5432/mydb?sslmode=disable{{- end }}"
      "destination" = "/vault/secrets/db-creds"
    }

    "vault" = {
      "address" = "https://vault.demo.svc.cluster.local:8200"
      "ca_cert" = "/vault/tls/ca.crt"
      "client_cert" = "/vault/tls/client.crt"
      "client_key" = "/vault/tls/client.key"
    }
```
