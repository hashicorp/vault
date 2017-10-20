---
layout: "docs"
page_title: "Auth Plugin Backend: Kubernetes"
sidebar_current: "docs-auth-kubernetes"
description: |-
  The Kubernetes auth backend allows automated authentication of Kubernetes
  Service Accounts.
---

# Auth Backend: Kubernetes

Name: `kubernetes`

The Kubernetes auth backend can be used to authenticate with Vault using a
Kubernetes Service Account Token. This method of authentication makes it easy to
introduce a Vault token into a Kubernetes Pod. 

## Authentication

#### Via the CLI

```
$ vault write auth/kubernetes/login role=demo jwt=...

Key                                   	Value
---                                   	-----
token                                 	1a445c6a-1ff5-7085-18f7-eca12210981d
token_accessor                        	fa82afb3-298b-41b0-6593-8b861bd3dc12
token_duration                        	768h0m0s
token_renewable                       	true
token_policies                        	[default]
token_meta_service_account_secret_name	"vault-auth-token-pd21c"
token_meta_service_account_uid        	"aa9aa8ff-98d0-11e7-9bb7-0800276d99bf"
token_meta_role                       	"demo"
token_meta_service_account_name       	"vault-auth"
token_meta_service_account_namespace  	"default"
```

#### Via the API

The endpoint for the kubernetes login is `auth/kubernetes/login`. 

The `kubernetes` mountpoint value in the url is the default mountpoint value.
If you have mounted the `kubernetes` backend with a different mountpoint, use that value.

```shell
$ curl $VAULT_ADDR/v1/auth/kubernetes/login \
    -d '{ "jwt": "your_service_account_jwt", "role": "demo" }'
```

The response will be in JSON. For example:

```javascript
{
	"request_id": "e344f8c2-fffc-c3e0-d118-e3a2e5de2d0d",
	"lease_id": "",
	"lease_duration": 0,
	"renewable": false,
	"data": null,
	"warnings": null,
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

First, you must enable the Kubernetes auth backend:

```
$ vault auth-enable kubernetes
Successfully enabled 'kubernetes' at 'kubernetes'!
```

Now when you run `vault auth -methods`, the Kubernetes backend is available:

```
Path         Type        Description
kubernetes/  kubernetes
token/       token       token based credentials
```

Prior to using the Kubernetes auth backend, it must be configured. To
configure it, use the `/config` endpoint.

```
$ vault write auth/kubernetes/config \
    kubernetes_host=https://192.168.99.100:8443 \
    kubernetes_ca_cert=@ca.crt
```

## Creating a Role

Authentication with this backend is role based. Before a token can be used to
login it first must be configured in a role.

```
vault write auth/kubernetes/role/demo \
    bound_service_account_names=vault-auth \ 
    bound_service_account_namespaces=default \
    policies=default \
    ttl=1h
```

This role Authorizes the vault-auth service account in the default namespace and
it gives it the default policy.

## Configuring Kubernetes

### Token Review Lookup
This backend accesses the [Kubernetes TokenReview
API](https://kubernetes.io/docs/api-reference/v1.7/#tokenreview-v1-authentication)
to validate the provided JWT is still valid. Kubernetes should be running with
`--service-account-lookup`. This is defaulted to true in Kubernetes 1.7, but any
versions prior should ensure the Kubernetes API server is started with with this
setting. Otherwise deleted tokens in Kubernetes will not be properly revoked and
will be able to authenticate to this backend. 

### RBAC Configuration

Service Accounts used in this backend will need to have access to the
TokenReview API. If Kubernetes is configured to use Role Based Access Control
the Service Account should be granted permissions to access this API. The
following example ClusterRoleBinding could be used to grant these permissions:

```
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

### GKE 

Currently the Token Review API endpoint is only available in alpha clusters on
Google Container Engine. This means on GKE this backend can only be used with an
alpha cluster.

## API

The Kubernetes Auth Plugin has a full HTTP API. Please see the
[API docs](/api/auth/kubernetes/index.html) for more details.


