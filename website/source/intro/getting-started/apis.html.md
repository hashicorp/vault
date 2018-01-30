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
[dev mode](/docs/concepts/dev-server.html), you could validate the
initialization status like this:

```text
$ curl http://127.0.0.1:8200/v1/sys/init
```

This will return a JSON response:

```json
{
  "initialized": true
}
```

## Accessing Secrets via the REST APIs

Machines that need access to information stored in Vault will most likely
access Vault via its REST API. For example, if a machine were using
[AppRole](/docs/auth/approle.html) for authentication, the application would
first authenticate to Vault which would return a Vault API token. The
application would use that token for future communication with Vault.

For the purpose of this guide, we will use the following configuration which
disables TLS and uses a file-based backend. TLS is disabled here only for
exemplary purposes; it should never be disabled in production.

```hcl
backend "file" {
  path = "vault"
}

listener "tcp" {
  tls_disable = 1
}
```

Save this file on disk as `config.hcl` and then start the Vault server:

```text
$ vault server -config=config.hcl
```

At this point, we can use Vault's API for all our interactions. For example, we
can initialize Vault like this:

```text
$ curl \
    --request POST \
    --data '{"secret_shares": 1, "secret_threshold": 1}' \
    http://127.0.0.1:8200/v1/sys/init
```

The response should be JSON and looks something like this:

```json
{
  "keys": [
    "373d500274dd8eb95271cb0f868e4ded27d9afa205d1741d60bb97cd7ce2fe41"
  ],
  "keys_base64": [
    "Nz1QAnTdjrlSccsPho5N7SfZr6IF0XQdYLuXzXzi/kE="
  ],
  "root_token": "6fa4128e-8bd2-fd02-0ea8-a5e020d9b766"
}
```

This response contains our initial root token. It also includes the unseal key.
You can use the unseal key to unseal the Vault and use the root token perform
other requests in Vault that require authentication.

To make this guide easy to copy-and-paste, we will be using the environment
variable `$VAULT_TOKEN` to store the root token. You can export this Vault
token in your current shell like this:

```sh
$ export VAULT_TOKEN=6fa4128e-8bd2-fd02-0ea8-a5e020d9b766
```

Using the unseal key (not the root token) from above, you can unseal the Vault
via the HTTP API:

```text
$ curl \
    --request POST \
    --data '{"key": "Nz1QAnTdjrlSccsPho5N7SfZr6IF0XQdYLuXzXzi/kE="}' \
    http://127.0.0.1:8200/v1/sys/unseal
```

Note that you should replace `Nz1QAnT...` with the generated key from your
output. This will return a JSON response:

```json
{
  "sealed": false,
  "t": 1,
  "n": 1,
  "progress": 0,
  "nonce": "",
  "version": "1.2.3",
  "cluster_name": "vault-cluster-9d524900",
  "cluster_id": "d69ab1b0-7e9a-2523-0d05-b0bfd09caeea"
}
```

Now any of the available auth methods can be enabled and configured.
For the purposes of this guide lets enable [AppRole](/docs/auth/approle.html)
authentication.

Start by enabling the AppRole authentication.

```text
$ curl \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data '{"type": "approle"}' \
    http://127.0.0.1:8200/v1/sys/auth/approle
```

Notice that the request to enable the AppRole endpoint needed an authentication
token. In this case we are passing the root token generated when we started
the Vault server. We could also generate tokens using any other authentication
mechanisms, but we will use the root token for simplicity.

Now create an AppRole with desired set of [ACL
policies](/docs/concepts/policies.html). In the following command, it is being
specified that the tokens issued under the AppRole `my-role`, should be
associated with `dev-policy` and the `my-policy`.

```text
$ curl \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data '{"policies": ["dev-policy", "my-policy"]}' \
    http://127.0.0.1:8200/v1/auth/approle/role/my-role
```

The AppRole backend, in its default configuration expects two hard to guess
credentials, a role ID and a secret ID. This command fetches the role ID of the
`my-role`.

```text
$ curl \
    --header "X-Vault-Token: $VAULT_TOKEN" \
     http://127.0.0.1:8200/v1/auth/approle/role/my-role/role-id
```

The response will include a `data` key with the `role_id`:

```json
{
  "data": {
    "role_id": "86a32a73-1f2b-05e0-113a-dfa930145d72"
  }
}
```

This command creates a new secret ID under the `my-role`.

```text
$ curl \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    http://127.0.0.1:8200/v1/auth/approle/role/my-role/secret-id
```

The response will include the `secret_id` in the `data` key:

```json
{
  "data": {
    "secret_id": "cd4b2002-3e3b-aceb-378d-5caa84dffd14",
    "secret_id_accessor": "6b9b58f6-d11a-c73c-ffa8-04a47d42716b"
  }
}
```

These two credentials can be supplied to the login endpoint to fetch a new
Vault token.

```text
$ curl \
    --request POST \
    --data '{"role_id": "86a32a73-1f2b-05e0-113a-dfa930145d72", "secret_id": "cd4b2002-3e3b-aceb-378d-5caa84dffd14"}' \
    http://127.0.0.1:8200/v1/auth/approle/login
```

The response will be JSON, under the key `auth`:

```json
{
  "auth": {
    "client_token": "50617721-dfb5-1916-7b13-4091e169d28c",
    "accessor": "ada8d354-47c0-5d9e-50f9-d74e6de2df9b",
    "policies": [
      "default",
      "dev-policy",
      "my-policy"
    ],
    "metadata": {
      "role_name": "my-role"
    },
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

The returned client token (`50617721...`) can be used to authenticate with
Vault. This token will be authorized with specific capabilities on all the
resources encompassed by the policies `default`, `dev-policy` and `my-policy`.

The newly acquired token can be exported as a new `VAULT_TOKEN` and use it to
authenticate Vault requests.

```sh
$ export VAULT_TOKEN="50617721-dfb5-1916-7b13-4091e169d28c"
```

```text
$ curl \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data '{"bar": "baz"}' \
    http://127.0.0.1:8200/v1/secret/foo
```

This will create a new secret named "foo" with the given JSON contents. We can
read this value back with the same token:

```text
$ curl \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    http://127.0.0.1:8200/v1/secret/foo
```

This should return a response like this:

```json
{
  "data": {
    "bar": "baz"
  },
  "lease_duration": 2764800,
  "renewable": false,
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
