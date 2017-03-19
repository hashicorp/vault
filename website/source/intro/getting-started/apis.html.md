---
layout: "intro"
page_title: "Using the HTTP APIs with Authentication"
sidebar_current: "gettingstarted-apis"
description: |-
  Using the HTTP APIs for authentication and secret access.
---

# Using the HTTP APIs with Authentication
All of Vault's capabilities are accessible via the HTTP API in addition to the
CLI. In fact, most calls from the CLI actually invoke the HTTP API. In some
cases, Vault features are not available via the CLI and can only be accessed
via the HTTP API.

Once you have started the Vault server, you can use `curl` or any other http
client to make API calls. For example, if you started the Vault server in
[development mode](/docs/concepts/dev-server.html), you could validate the
initialization status like this:

```javascript
$ curl http://127.0.0.1:8200/v1/sys/init
```

This will return a JSON response:

```javascript
{ "initialized": true }
```

## Accessing Secrets via the REST APIs
Machines that need access to information stored in Vault will most likely
access Vault via its REST API. For example, if a machine were using
[AppRole](/docs/auth/approle.html) for authentication, the application would
first authenticate to Vault which would return a Vault API token. The
application would use that token for future communication with Vault.

For the purpose of this guide, we will use the following configuration which
disables TLS and uses a file-based backend. TLS is disabled here only for
exemplary purposes and it should never be disabled in production.

```javascript
backend "file" {
  path = "vault"
}

listener "tcp" {
  tls_disable = 1
}
```

Save this file on disk and then start the Vault server with this command:

```javascript
$ vault server -config=/etc/vault.conf
```

At this point, we can use Vault's API for all our interactions. For example, we
can initialize Vault like this:

```javascript
$ curl \
  -X PUT \
  -d "{\"secret_shares\":1, \"secret_threshold\":1}" \
  http://127.0.0.1:8200/v1/sys/init
```

The response should be JSON and looks something like this:

```javascript
{
  "root_token": "4f66bdfa-f5e4-209f-096c-6e01d863c145",
  "keys_base64": [
    "FwwsSzMysLgYAvJFrs+q5UfLMKIxC+dDFbP6YzyjzvQ="
  ],
  "keys": [
    "170c2c4b3332b0b81802f245aecfaae547cb30a2310be74315b3fa633ca3cef4"
  ]
}
```

This response contains our initial root token. It also includes the unseal key.
You can use the unseal key to unseal the Vault and use the root token perform
other requests in Vault that require authentication.

To make this guide easy to copy-and-paste, we will be using the environment
variable `$VAULT_TOKEN` to store the root token. You can export this Vault
token in your current shell like this:

```javascript
$ export VAULT_TOKEN=4f66bdfa-f5e4-209f-096c-6e01d863c145
```

Using the unseal key (not the root token) from above, you can unseal the Vault
via the HTTP API:

```javascript
$ curl \
    -X PUT \
    -d '{"key": "FwwsSzMysLgYAvJFrs+q5UfLMKIxC+dDFbP6YzyjzvQ="}' \
    http://127.0.0.1:8200/v1/sys/unseal
```

Note that you should replace `FwwsSzM...` with the generated key from your
output. This will return a JSON response:

```javascript
{
  "cluster_id": "1c2523c9-adc2-7f3a-399f-7032da2b9faf",
  "cluster_name": "vault-cluster-9ac82317",
  "version": "0.6.2",
  "progress": 0,
  "n": 1,
  "t": 1,
  "sealed": false
}
```

Now any of the available authentication backends can be enabled and configured.
For the purposes of this guide lets enable [AppRole](/docs/auth/approle.html)
authentication.

Start by enabling the AppRole authentication.

```javascript
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"type":"approle"}' http://127.0.0.1:8200/v1/sys/auth/approle
```

Notice that the request to enable the AppRole endpoint needed an authentication
token. In this case we are passing the root token generated when we started
the Vault server. We could also generate tokens using any other authentication
mechanisms, but we will use the root token for simplicity.

Now create an AppRole with desired set of [ACL
policies](/docs/concepts/policies.html). In the following command, it is being
specified that the tokens issued under the AppRole `testrole`, should be
associated with `dev-policy` and the `test-policy`.

```javascript
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"policies":"dev-policy,test-policy"}' http://127.0.0.1:8200/v1/auth/approle/role/testrole
```

The AppRole backend, in its default configuration expects two hard to guess
credentials, a role ID and a secret ID. This command fetches the role ID of
the `testrole`.

```javascript
$ curl -X GET -H "X-Vault-Token:$VAULT_TOKEN" http://127.0.0.1:8200/v1/auth/approle/role/testrole/role-id | jq .
```

```javascript
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "role_id": "988a9dfd-ea69-4a53-6cb6-9d6b86474bba"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "ef5c9b3f-e15e-0527-5457-79b4ecfe7b60"
}
```

This command creates a new secret ID under the `testrole`.

```javascript
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" http://127.0.0.1:8200/v1/auth/approle/role/testrole/secret-id | jq .
```

```javascript
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "secret_id_accessor": "45946873-1d96-a9d4-678c-9229f74386a5",
    "secret_id": "37b74931-c4cd-d49a-9246-ccc62d682a25"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "c98fa1c2-7565-fd45-d9de-0b43c307f2aa"
}
```

These two credentials can be supplied to the login endpoint to fetch a new
Vault token.

```javascript
$ curl -X POST \
     -d '{"role_id":"988a9dfd-ea69-4a53-6cb6-9d6b86474bba","secret_id":"37b74931-c4cd-d49a-9246-ccc62d682a25"}' \
     http://127.0.0.1:8200/v1/auth/approle/login | jq .
```

```javascript
{
  "auth": {
    "renewable": true,
    "lease_duration": 2764800,
    "metadata": {},
    "policies": [
      "default",
      "dev-policy",
      "test-policy"
    ],
    "accessor": "5d7fb475-07cb-4060-c2de-1ca3fcbf0c56",
    "client_token": "98a4c7ab-b1fe-361b-ba0b-e307aacfd587"
  },
  "warnings": null,
  "wrap_info": null,
  "data": null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "988fb8db-ce3b-0167-0ac7-1a568b902d75"
}
```

The returned client token (`98a4c7ab-b1fe-361b-ba0b-e307aacfd587`) can now be
used to authenticate with Vault. This token will be authorized with specific
capabilities on all the resources encompassed by the policies `default`,
`dev-policy` and `test-policy`.

The newly acquired token can be exported as a new `VAULT_TOKEN` and use it to
authenticate Vault requests.

```javascript
$ export VAULT_TOKEN="98a4c7ab-b1fe-361b-ba0b-e307aacfd587"
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"bar":"baz"}' http://127.0.0.1:8200/v1/secret/foo
```

This will create a new secret named "foo" with the given JSON contents. We can
read this value back with the same token:

```javascript
$ curl -X GET -H "X-Vault-Token:$VAULT_TOKEN" http://127.0.0.1:8200/v1/secret/foo | jq .
```

This should return a response like this:

```javascript
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "bar": "baz"
  },
  "lease_duration": 2764800,
  "renewable": false,
  "lease_id": "",
  "request_id": "5e246671-ec05-6fc8-9f93-4fe4512f34ab"
}
```

You can see the documentation on the [HTTP APIs](/api/index.html) for
more details on other available endpoints.

Congratulations! You now know all the basics to get started with Vault.

## Next

Next, we have a page dedicated to [next
steps](/intro/getting-started/next-steps.html) depending on what you would like
to achieve.
