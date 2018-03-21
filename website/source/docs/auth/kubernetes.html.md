---
layout: "docs"
page_title: "Kubernetes - Auth Methods"
sidebar_current: "docs-auth-kubernetes"
description: |-
  The Kubernetes auth method allows automated authentication of Kubernetes
  Service Accounts.
---

# Kubernetes Auth Method

The `kubernetes` auth method can be used to authenticate with Vault using a
Kubernetes Service Account Token. This method of authentication makes it easy to
introduce a Vault token into a Kubernetes Pod.

## Authentication

### Via the CLI

The default path is `/kubernetes`. If this auth method was enabled at a
different path, specify `-path=/my-path` in the CLI.


```text
$ vault write auth/kubernetes/login role=demo jwt=...
```

### Via the API

The default endpoint is `auth/kubernetes/login`. If this auth method was enabled
at a different path, use that value instead of `kubernetes`.

```shell
$ curl \
    --request POST \
    --data '{"jwt": "your_service_account_jwt", "role": "demo"}' \
    https://vault.rocks/v1/auth/kubernetes/login
```

The response will contain a token at `auth.client_token`:

```json
{
  "auth": {
    "client_token": "38fe9691-e623-7238-f618-c94d4e7bc674",
    "accessor": "78e87a38-84ed-2692-538f-ca8b9f400ab3",
    "policies": [
      "default"
    ],
    "metadata": {
      "role": "test",
      "service_account_name": "vault-auth",
      "service_account_namespace": "default",
      "service_account_secret_name": "vault-auth-token-pd21c",
      "service_account_uid": "aa9aa8ff-98d0-11e7-9bb7-0800276d99bf"
    },
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

## Configuration

Auth methods must be configured in advance before users or machines can
authenticate. These steps are usually completed by an operator or configuration
management tool.


1. Enable the Kubernetes auth method:

    ```text
    $ vault auth enable kubernetes
    ```

1. Use the `/config` endpoint to configure Vault to talk to Kubernetes. For the
list of available configuration options, please see the API documentation.

    ```text
    $ vault write auth/kubernetes/config \
        token_reviewer_jwt="reviewer_service_account_jwt" \
        kubernetes_host=https://192.168.99.100:8443 \
        kubernetes_ca_cert=@ca.crt
    ```

1. Create a named role:

    ```text
    vault write auth/kubernetes/role/demo \
        bound_service_account_names=vault-auth \
        bound_service_account_namespaces=default \
        policies=default \
        ttl=1h
    ```

    This role authorizes the "vault-auth" service account in the default
    namespace and it gives it the default policy.

    For the complete list of configuration options, please see the API
    documentation.

## Configuring Kubernetes

This auth method accesses the [Kubernetes TokenReview API][k8s-tokenreview] to
validate the provided JWT is still valid. Kubernetes should be running with
`--service-account-lookup`. This is defaulted to true in Kubernetes 1.7, but any
versions prior should ensure the Kubernetes API server is started with with this
setting. Otherwise deleted tokens in Kubernetes will not be properly revoked and
will be able to authenticate to this auth method.

Service Accounts used in this auth method will need to have access to the
TokenReview API. If Kubernetes is configured to use RBAC roles the Service
Account should be granted permissions to access this API. The following
example ClusterRoleBinding could be used to grant these permissions:

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: role-tokenreview-binding
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: vault-auth
  namespace: default
```

## API

The Kubernetes Auth Plugin has a full HTTP API. Please see the
[API docs](/api/auth/kubernetes/index.html) for more details.

[k8s-tokenreview]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#tokenreview-v1-authentication
