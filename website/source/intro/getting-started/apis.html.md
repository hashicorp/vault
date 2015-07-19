---
layout: "intro"
page_title: "Using the HTTP APIs with Authentication"
sidebar_current: "gettingstarted-apis"
description: |-
  Using the HTTP APIs for authentication and secret access.
---

# Using the HTTP APIs with Authentication
All Vault capabilities can be accessed via HTTP, rather than the CLI. In fact some calls, for example [app-id](/docs/auth/app-id.html) authentication, cannot be called by the CLI at all. Once you have started your server, you can use curl, or any other http client to make API calls. For example, if you have started Vault in dev mode, you could validated initialization status like this:

```
$ curl http://127.0.0.1:8200/v1/sys/init

{"initialized":true}
```

# Access Secrets via the REST APIs
Machines will most likely access Vault via the REST APIs. Assuming a machine is using the [app-id](/docs/auth/app-id.html) backend for authentication, the flow would look like this:
![REST Sequence](/assets/images/app-id-api-sequence.png)

Dev mode doesn't require authentication, so we want a configuration something like this:

```
  backend "file" {
    path = "vault"
  }

  listener "tcp" {
    tls_disable = 1
  }
```
and start the server like this:

```
$ vault server -config=/etc/vault.conf

```
At this point, we can start using the APIs for all our interactions. For example, we can initialize the instance like this:

```
$ curl -f -XPUT --data "{\"secret_shares\":1, \"secret_threshold\":1}" http://localhost:8200/v1/sys/init

{"keys":["69cf1c12a1f65dddd19472330b28cf4e95c657dfbe545877e5765d25d0592b16"],"root_token":"0e2ede5a-6664-a49e-ca33-8f204d1cdb95"}
```
And now we have our root token, so we can unseal the vault, and enable app_id authentication through the REST APIs as well.

```
$ curl -XPUT --data '{"key": "69cf1c12a1f65dddd19472330b28cf4e95c657dfbe545877e5765d25d0592b16"}' http://127.0.0.1:8200/v1/sys/unseal

{"sealed":false,"t":1,"n":1,"progress":0}

$ curl -XPOST -H'X-Vault-Token:0e2ede5a-6664-a49e-ca33-8f204d1cdb95' --data '{"type":"app-id"}' http://127.0.0.1:8200/v1/sys/auth/app-id
```
Notice that the request to the app-id endpoint needed a token. In this case the only token we have is the root token so we can use it.

The last thing we need to do before using our app-id endpoint is writing the data to the store to associate an app id with a user id. For more information on this process, see the documentation on the [app-id auth backend](/docs/auth/app-id.html).

```
$ curl -XPOST -H'X-Vault-Token:0e2ede5a-6664-a49e-ca33-8f204d1cdb95' --data '{"value":"root", "display_name":"demo"}' http://localhost:8200/v1/auth/app-id/map/app-id/152AEA38-85FB-47A8-9CBD-612D645BFACA

$ curl -XPOST -H'X-Vault-Token:0e2ede5a-6664-a49e-ca33-8f204d1cdb95' --data '{"value":"152AEA38-85FB-47A8-9CBD-612D645BFACA"}' http://localhost:8200/v1/auth/app-id/map/user-id/5ADF8218-D7FB-4089-9E38-287465DBF37E
```
In the first request above, we associated the app with the ```root``` policy. You would not want to do this in a production scenario because the root policy allows complete read and write access. For a real deployed application you would [ create a policy](/docs/concepts/policies.html) with appropriate permissions, then associate the application id with it so you could control access to secrets.

Now your app can identify itself via the app-id and user-id and get access to the store. The first step is to authenticate:

```
$ curl -XPOST --data '{"app_id":"152AEA38-85FB-47A8-9CBD-612D645BFACA", "user_id": "5ADF8218-D7FB-4089-9E38-287465DBF37E"}' "http://127.0.0.1:8200/v1/auth/app-id/login"

{"lease_id":"","renewable":false,"lease_duration":0,"data":null,"auth":{"client_token":"7a25c58b-9bad-5750-b579-edbb9f10a5ef","policies":["root"],"metadata":{"app-id":"sha1:1c0401b419280b0771d006bcdae683989086a00e","user-id":"sha1:4dbf74fce71648d54c42e28ad193253600853ca6"},"lease_duration":0,"renewable":false}}
```

Now the token can be used to access the store. As you can see from the returned payload, the app-id backend does not currently support lease expiration or renewal. If you authenticate with backend that does support leases, your app will have to track expiration and handle renewal.

```
$ curl -XPOST -H 'X-Vault-Token: 7a25c58b-9bad-5750-b579-edbb9f10a5ef' -H 'Content-type: application/json' --data '{"bar":"baz"}' http://127.0.0.1:8200/v1/secret/foo

$ curl -H 'X-Vault-Token: 7a25c58b-9bad-5750-b579-edbb9f10a5ef' http://127.0.0.1:8200/v1/secret/foo

{"lease_id":"secret/foo/cc529d06-36c8-be27-31f5-2390e1f6e2ae","renewable":false,"lease_duration":2592000,"data":{"bar":"baz"},"auth":null}
```

See the documentation on the [HTTP APIs](/docs/http/index.html) for more details on other available endpoints.

Congratulations! You now know all the basics to get started with Vault.

## Next

Next, we have a page dedicated to
[next steps](/intro/getting-started/next-steps.html) depending on
what you'd like to achieve.
