---
layout: "docs"
page_title: "Google Cloud - Auth Methods"
sidebar_current: "docs-auth-gcp"
description: |-
  The "gcp" auth method allows users and machines to authenticate to Vault using
  Google Cloud service accounts.
---

# Google Cloud Auth Method

The `gcp` auth method allows authentication against Vault using Google
credentials. It treats Google Cloud Platform (GCP) as a Trusted Third Party and
expects a [JSON Web Token][jwt] (JWT) signed by Google credentials from the
authenticating entity. This token can be generated through different GCP APIs
depending on the type of entity.

This plugin is developed in a separate GitHub repository at
[`hashicorp/vault-plugin-auth-gcp`](https://github.com/hashicorp/vault-plugin-auth-gcp),
but is automatically bundled in Vault releases. Please file all feature
requests, bugs, and pull requests specific to the GCP plugin under that
repository.

## Authentication

### Via the CLI

The default path is `/gcp`. If this auth method was enabled at a different
path, specify `-path=/my-path` in the CLI.

```text
$ vault login -method=gcp \
    role="my-role" \
    jwt="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

In this example, the "role" is the name of a configured role. The "jwt" is a
self-signed or Google-signed JWT token obtained using the
[`signJwt`][signjwt-method] API call.

Because the process to sign a service account JWT can be tedious, Vault includes
a CLI helper to generate the JWT token given the service account and parameters.
This process **only applies to `iam`-type roles!**

```text
$ vault login -method=gcp \
    role="my-role" \
    jwt_exp="15m" \
    credentials=@path/to/credentials.json \
    project="my-project" \
    service_account="service-account@my-project.iam.gserviceaccounts.com"
```

This signs a properly formatted service account JWT and authenticates to Vault
directly. For details on each field, please run `vault auth help gcp`.

### Via the API

```text
$ curl \
    --request POST \
    --data '{"role":"my-role", "jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}' \
    http://127.0.0.1:8200/v1/auth/gcp/login
```

The response will be in JSON. For example:

```javascript
{
  "auth": {
    "client_token": "f33f8c72-924e-11f8-cb43-ac59d697597c",
    "accessor": "0e9e354a-520f-df04-6867-ee81cae3d42d",
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "metadata": {
      "role": "my-role",
      "service_account_email": "dev1@project-123456.iam.gserviceaccount.com",
      "service_account_id": "111111111111111111111"
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

1. Enable the Google Cloud auth method:

    ```text
    $ vault auth enable gcp
    ```

1. Configure the auth method credentials:

    ```text
    $ vault write auth/gcp/config \
        credentials=@/path/to/credentials.json
    ```

    If you are using instance credentials or want to specify credentials via
    an environment variable, you can skip this step. To learn more, see the
    [Google Cloud Authentication](#google-cloud-authentication) section below.

1. Create a named role:

    For an `iam`-type role:

    ```text
    $ vault write auth/gcp/role/my-iam-role \
        type="iam" \
        project_id="my-project" \
        policies="dev,prod" \
        bound_service_accounts="my-service@my-project.iam.gserviceaccount.com"
    ```

    For a `gce`-type role:

    ```text
    $ vault write auth/gcp/role/my-gce-role \
        type="gce" \
        project_id="my-project" \
        policies="dev,prod" \
        bound_zones="us-east1-b" \
        bound_labels="foo:bar,zip:zap"
    ```

    For the complete list of configuration options for each type, please see the
    [API documentation][api-docs].


### Google Cloud Authentication

The Google Cloud Vault auth method uses the official Google Cloud Golang SDK.
This means it supports the common ways of [providing credentials to Google
Cloud][cloud-creds].

1. The environment variable `GOOGLE_APPLICATION_CREDENTIALS`. This is specified
as the **path** to a Google Cloud credentials file, typically for a service
account. If this environment variable is present, the resulting credentials are
used. If the credentials are invalid, an error is returned.

1. Default instance credentials. When no environment variable is present, the
default service account credentials are used.

For more information on service accounts, please see the [Google Cloud Service
Accounts documentation][service-accounts].

To use this storage backend, the service account must have the following
minimum scope(s):

```text
https://www.googleapis.com/auth/cloud-platform
```

## Workflow

This section describes the implementation details for how Vault communicates
with Google Cloud to authenticate and authorize JWT tokens. This information is
provided for those who are curious, but these implementation details are not
required knowledge for using the auth method.

### IAM Login

IAM login applies only to roles of type `iam`. The Vault authentication workflow
for IAM service accounts looks like this:

[![Vault Google Cloud IAM Login Workflow](/assets/images/vault-gcp-iam-auth-workflow.svg)](/assets/images/vault-gcp-iam-auth-workflow.svg)

  1. The client generates a signed JWT using the IAM
  [`projects.serviceAccounts.signJwt`][signjwt-method] method. For examples of
  how to do this, see the [Obtaining JWT Tokens](#obtaining-jwt-tokens) section.

  2. The client sends this signed JWT to Vault along with a role name.

  3. Vault extracts the `kid` header value, which contains the ID of the
  key-pair used to generate the JWT, and the `sub` ID/email to find the service
  account key. If the service account does not exist or the key is not linked to
  the service account, Vault denies authentication.

  4. Vault authorizes the confirmed service account against the given role. If
  that is successful, a Vault token with the proper policies is returned.

### GCE Login

GCE login only applies to roles of type `gce` and **must be completed on an
instance running in GCE**. These steps will not work from your local laptop or
another cloud provider.

[![Vault Google Cloud GCE Login Workflow](/assets/images/vault-gcp-gce-auth-workflow.svg)](/assets/images/vault-gcp-gce-auth-workflow.svg)

  1. The client obtains an [instance identity metadata token][instance-identity]
  on a GCE instance.

  2. The client sends this JWT to Vault along with a role name.

  3. Vault extracts the `kid` header value, which contains the ID of the
  key-pair used to generate the JWT, to find the OAuth2 public cert to verify
  this JWT.

  4. Vault authorizes the confirmed instance against the given role, ensuring
  the instance matches the bound zones, regions, or instance groups. If that is
  successful, a Vault token with the proper policies is returned.

## Obtaining JWT Tokens

Vault expects a signed JWT token to verify against. There are a few ways to
acquire a JWT token.

### Generating IAM Tokens

Vault includes a CLI helper for generating the signed JWT token and submitting
it to Vault for `iam`-type roles. If you want to generate the JWT token
yourself, follow this section.

#### Shell Example

The expected format of the JWT request payload is:

```javascript
{
  "sub": "$SERVICE_ACCOUNT",
  "aud": "vault/$ROLE",
  "exp": "$EXPIRATION" // optional
}
```

If specified, the expiration must be a
[NumericDate](https://tools.ietf.org/html/rfc7519#section-2) value (seconds from
Epoch). This value must be before the max JWT expiration allowed for a role.
This defaults to 15 minutes and cannot be more than 1 hour.

One you have all this information, the JWT token can be signed using curl and
[oauth2l](https://github.com/google/oauth2l):

```text
ROLE="my-role"
PROJECT="my-project"
SERVICE_ACCOUNT="service-account@my-project.iam.gserviceaccount.com"
OAUTH_TOKEN="$(oauth2l header cloud-platform)"

curl \
  --header "${OAUTH_TOKEN}" \
  --header "Content-Type: application/json" \
  --request POST \
  --data "{\"aud\":\"vault/${ROLE}\", \"sub\": \"${SERVICE_ACCOUNT}\"}" \
  "https://iam.googleapis.com/v1/projects/${PROJECT}/serviceAccounts/${SERVICE_ACCOUNT}:signJwt"
```

#### gcloud Example

```text
gcloud beta iam service-accounts sign-jwt credentials.json - \
  --iam-account=service-account@my-project.iam.gserviceaccount.com \
  --project=my-project
```

#### Golang Example

Read more on the
[Google Open Source blog](https://opensource.googleblog.com/2017/08/hashicorp-vault-and-google-cloud-iam.html).

### Generating GCE Tokens

GCE tokens can only be generated from a GCE instance. **You must run these
commands from the GCE instance.** The JWT token can be obtained from the
`service-accounts/default/identity` endpoint for a instance's metadata server.

```text
ROLE="my-gce-role"
SERVICE_ACCOUNT="service-account@my-project.iam.gserviceaccount.com"

curl \
  --header "Metadata-Flavor: Google" \
  --get \
  --data-urlencode "aud=http://vault/${ROLE}" \
  --data-urlencode "format=full" \
  "http://metadata/computeMetadata/v1/instance/service-accounts/${SERVICE_ACCOUNT}/identity"
```

## API

The GCP Auth Plugin has a full HTTP API. Please see the
[API docs][api-docs] for more details.

[jwt]: https://tools.ietf.org/html/rfc7519
[signjwt-method]: https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt
[cloud-creds]: https://cloud.google.com/docs/authentication/production#providing_credentials_to_your_application
[service-accounts]: https://cloud.google.com/compute/docs/access/service-accounts
[api-docs]: /api/auth/gcp/index.html
[instance-identity]: https://cloud.google.com/compute/docs/instances/verifying-instance-identity
