#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

POD_NAME="${APP_POD_NAME:-vault-client-app}"

# kubectl is not installed on the runner. kind control-plane containers always
# have kubectl available, so proxy all kubectl calls through docker exec.
CONTROL_PLANE="${KIND_CLUSTER_NAME}-control-plane"
KUBECTL="docker exec ${CONTROL_PLANE} kubectl"

# Clean up existing resources if any
$KUBECTL delete pod "$POD_NAME" --ignore-not-found=true
$KUBECTL delete configmap vault-client-config --ignore-not-found=true
$KUBECTL delete secret vault-client-approle --ignore-not-found=true

# Create Secret with AppRole credentials
$KUBECTL create secret generic vault-client-approle \
  --from-literal=role-id="$ROLE_ID" \
  --from-literal=secret-id="$SECRET_ID"

# Create ConfigMap with vault agent / client configuration
docker exec -i "${CONTROL_PLANE}" kubectl apply -f - << 'EOF'
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault-client-config
data:
  agent-config.hcl: |
    pid_file = "/vault/vault-agent.pid"

    vault {
      address = "http://vault.default.svc.cluster.local:8200"
      retry {
        num_retries = 5
      }
    }

    auto_auth {
      method {
        type = "approle"
        config = {
          role_id_file_path   = "/vault/approle/role-id"
          secret_id_file_path = "/vault/approle/secret-id"
        }
      }
      sink {
        type = "file"
        config = {
          path = "/vault/token/vault-token"
        }
      }
    }

    # KV v1 path: secret/<key>. Update to secret/data/<key> if switching to KV v2.
    template {
      contents     = "{{ with secret \"secret/client-sidecar\" }}{{ .Data.test_key }}{{ end }}"
      destination  = "/vault/secrets/rendered-secret.txt"
    }

    cache {
      use_auto_auth_token = true
    }

    listener "tcp" {
      address     = "127.0.0.1:8200"
      tls_disable = true
    }
EOF

# Deploy Pod with application container and vault-client sidecar container
docker exec -i "${CONTROL_PLANE}" kubectl apply -f - << EOF
apiVersion: v1
kind: Pod
metadata:
  name: ${POD_NAME}
  labels:
    app: vault-client-test
spec:
  volumes:
    - name: vault-config
      configMap:
        name: vault-client-config
    - name: vault-approle
      secret:
        secretName: vault-client-approle
    - name: vault-token
      emptyDir: {}
    - name: vault-secrets
      emptyDir: {}

  containers:
    - name: app
      image: curlimages/curl:latest
      imagePullPolicy: IfNotPresent
      command: ["sh", "-c", "sleep 3600"]
      volumeMounts:
        - name: vault-secrets
          mountPath: /vault/data

    - name: vault-client-sidecar
      image: ${CLIENT_IMAGE}
      imagePullPolicy: Never
      command:
        - "vault"
        - "agent"
        - "-config=/vault/config/agent-config.hcl"
      volumeMounts:
        - name: vault-config
          mountPath: /vault/config
        - name: vault-approle
          mountPath: /vault/approle
        - name: vault-token
          mountPath: /vault/token
        - name: vault-secrets
          mountPath: /vault/secrets
EOF

# Wait for Pod to become ready
$KUBECTL wait --for=condition=Ready "pod/${POD_NAME}" --timeout=120s
