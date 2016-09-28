---
layout: "intro"
page_title: "Using the HTTP APIs with Authentication"
sidebar_current: "gettingstarted-apis"
description: |-
  Using the HTTP APIs for authentication and secret access.
---

# Using the HTTP APIs with Authentication
All of Vault's capabilities are accessible via the HTTP API in addition to the CLI. In fact, most calls from the CLI actually invoke the HTTP API. In some cases, Vault features are not available via the CLI and can only be accessed via the HTTP API.

Once you have started the Vault server, you can use `curl` or any other http client to make API calls. For example, if you started the Vault server in [development mode](/docs/concepts/dev-server.html), you could validate the initialization status like this:

```
$ curl http://127.0.0.1:8200/v1/sys/init
```

This will return a JSON response:

```javascript
{ "initialized": true }
```

## Accessing Secrets via the REST APIs
Machines that need access to information stored in Vault will most likely access Vault via its REST API. For example, if a machine were using [app-id](/docs/auth/app-id.html) for authentication, the application would first authenticate to Vault which would return a Vault API token. The application would use that token for future communication with Vault.

For the purpose of this guide, we will use the following configuration which disables TLS and uses a file-based backend. You should never disable TLS in production, but it is okay for the purposes of this tutorial.

```javascript
backend "file" {
  path = "vault"
}

listener "tcp" {
  tls_disable = 1
}
```

Save this file on disk and then start the Vault server with this command:

```
$ vault server -config=/etc/vault.conf
```

At this point, we can use Vault's API for all our interactions. For example, we can initialize Vault like this:

```
$ curl \
  -X PUT \
  -d "{\"secret_shares\":1, \"secret_threshold\":1}" \
  http://localhost:8200/v1/sys/init
```

The response should be JSON and look something like this:

```javascript
{
  "keys": ["69cf1c12a1f65dddd19472330b28cf4e95c657dfbe545877e5765d25d0592b16"],
  "root_token": "0e2ede5a-6664-a49e-ca33-8f204d1cdb95"
}
```

This response contains our initial root token. It also includes the unseal key. You can use the unseal key to unseal the Vault and use the root token perform other requests in Vault that require authentication.

To make this guide easy to copy-and-paste, we will be using the environment variable `$VAULT_TOKEN` to store the root token. You can export this Vault token in your current shell like this:

```
$ export VAULT_TOKEN=0e2ede5a-6664-a49e-ca33-8f204d1cdb95
```

Using the unseal key (not the root token) from above, you can unseal the Vault via the HTTP API:

```
$ curl \
    -X PUT \
    -d '{"key": "69cf1c12a1f65dddd19472330b28cf4e95c657dfbe545877e5765d25d0592b16"}' \
    http://127.0.0.1:8200/v1/sys/unseal
```

Note that you should replace `69cf1c1...` with the generated key from your output. This will return a JSON response:

```javascript
{
  "sealed": false,
  "t": 1,
  "n": 1,
  "progress": 0
}
```

Now we can enable an authentication backend such as [GitHub authentication](/docs/auth/github.html) or [App ID](/docs/auth/app-id.html). For the purposes of this guide, we will enable App ID authentication.

We can enable an authentication backend with the following `curl` command:

```
$ curl \
    -X POST \
    -H "X-Vault-Token:$VAULT_TOKEN" \
    -d '{"type":"app-id"}' \
    http://127.0.0.1:8200/v1/sys/auth/app-id
```

Notice that the request to the app-id endpoint needed an authentication token. In this case we are passing the root token generated when we started the Vault server. We could also generate tokens using any other authentication mechanisms, but we will use the root token for simplicity.

The last thing we need to do before using our App ID endpoint is writing the data to the store to associate an app id with a user id. For more information on this process, see the documentation on the [app-id auth backend](/docs/auth/app-id.html).

First, we need to associate the application with a particular [ACL policy](/docs/concepts/policies.html) in Vault. In the following command, we are going to associate the created tokens with the `root` policy. You would not want to do this in a real production scenario because the root policy allows complete read, write, and administrator access to Vault. For a production application, you should create an ACL policy (which is also possible via the HTTP API), but is not covered in this guide for simplicity.

```
$ curl \
    -X POST \
    -H "X-Vault-Token:$VAULT_TOKEN" \
    -d '{"value":"root", "display_name":"demo"}' \
    http://localhost:8200/v1/auth/app-id/map/app-id/152AEA38-85FB-47A8-9CBD-612D645BFACA
```

Note that `152AEA38-85FB-47A8-9CBD-612D645BFACA` is a randomly generated UUID. You can use any tool to generate a UUID, but make sure it is unique.

Next we need to map the application to a particular "user". In Vault, this is actually a particular application:

```
$ curl \
    -X POST \
    -H "X-Vault-Token:$VAULT_TOKEN" \
    -d '{"value":"152AEA38-85FB-47A8-9CBD-612D645BFACA"}' \
    http://localhost:8200/v1/auth/app-id/map/user-id/5ADF8218-D7FB-4089-9E38-287465DBF37E
```

Now your app can identify itself via the app-id and user-id and get access to Vault. The first step is to authenticate:

```
$ curl \
    -X POST \
    -d '{"app_id":"152AEA38-85FB-47A8-9CBD-612D645BFACA", "user_id": "5ADF8218-D7FB-4089-9E38-287465DBF37E"}' \
    "http://127.0.0.1:8200/v1/auth/app-id/login"
```

This will return a response that looks like the following:

```javascript
{
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "auth": {
    "client_token": "7a25c58b-9bad-5750-b579-edbb9f10a5ef",
    "policies": ["root"],
    "lease_duration": 0,
    "renewable": false,
    "metadata": {
      "app-id": "sha1:1c0401b419280b0771d006bcdae683989086a00e",
      "user-id": "sha1:4dbf74fce71648d54c42e28ad193253600853ca6"
    }
  }
}
```

The returned client token (`7a25c58b-9bad-5750-b579-edbb9f10a5ef`) can now be used to authenticate with Vault. As you can see from the returned payload, the App ID backend does not currently support lease expiration or renewal. If you authenticate with backend that does support leases, your application will have to track expiration and handle renewal, but that is a topic for another guide.

We can export this new token as our new `VAULT_TOKEN`:

```
$ export VAULT_TOKEN="7a25c58b-9bad-5750-b579-edbb9f10a5ef"
```

Be sure to replace this with the value returned from your API response. We can now use this token to authentication requests to Vault:

```
$ curl \
    -X POST \
    -H "X-Vault-Token:$VAULT_TOKEN" \
    -H 'Content-type: application/json' \
    -d '{"bar":"baz"}' \
    http://127.0.0.1:8200/v1/secret/foo
```

This will create a new secret named "foo" with the given JSON contents. We can read this value back with the same token:

```
$ curl \
    -H "X-Vault-Token:$VAULT_TOKEN" \
    http://127.0.0.1:8200/v1/secret/foo
```

This should return a response like this:

```javascript
{
  "lease_id": "secret/foo/cc529d06-36c8-be27-31f5-2390e1f6e2ae",
  "renewable": false,
  "lease_duration": 2764800,
  "data": {
    "bar": "baz"
  },
  "auth": null
}
```

You can see the documentation on the [HTTP APIs](/docs/http/index.html) for more details on other available endpoints.

Congratulations! You now know all the basics to get started with Vault.

## Next

Next, we have a page dedicated to
[next steps](/intro/getting-started/next-steps.html) depending on
what you would like to achieve.
