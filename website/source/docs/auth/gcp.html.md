---
layout: "docs"
page_title: "GCP - Auth Methods"
sidebar_current: "docs-auth-gcp"
description: |-
  The gcp auth method plugin allows automated authentication of AWS entities.
---

# GCP Auth Method

The `gcp` auth method allows authentication against Vault using
Google credentials. It treats GCP as a Trusted Third Party and expects a
[JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519) signed by Google
credentials from the authenticating entity. This token can be generated through
different GCP APIs depending on the type of entity.

Currently supports authentication for:

  * GCP IAM service accounts (`iam`)
  * GCE IAM service accounts (`gce`)

## Quick Links
 * [API documentation](/api/auth/gcp/index.html)
 * [Usage: Setup](#setup): How to set up the auth backend on the server
 * [Usage: Authentication](#authentication): How to authenticate to Vault using this backend
 * [Workflows](#authentication-workflow): Authentication and authorization workflows in more detail

## Setup

The following documentation assumes that the method has been
[mounted](/docs/plugin/index.html) at `auth/gcp`.

```sh
$ vault auth enable 'gcp'
Success! Enabled gcp auth method at: gcp/
```

You must also [enable the following GCP APIs](https://support.google.com/cloud/answer/6158841?hl=en)
for your GCP project:

  * IAM API for both `iam` service accounts and `gce` instances
  * GCE API for just `gce` instances

There are generally two setup steps:

### 1. Config

Example:

```sh
$ vault write auth/gcp/config credentials=@/path/to/creds.json
```

This step allows you to explicitly set GCP credentials that Vault (this auth backend) uses.
If credentials are not configured or if the user explicitly sets the config with no credentials,
the Vault server will attempt to infer credentials. In order of preference, the auth method will use
the first credentials it finds from the following:

* Provided JSON as `credentials` arg
* Environment variables `GOOGLE_CREDENTIALS` or `GOOGLE_CLOUD_KEYFILE_JSON` set to the credentials JSON string
* JSON file  `~/.gcp/credentials`
* [Google Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)

For the complete list of configuration options, see the [API documentation](/api/auth/gcp/index.html#configure).

### 2. Roles

Example:

```sh
$ vault write auth/gcp/role/dev-iam-role \
    type="iam" \
    project_id="project-123456" \
    policies="dev" \
    bound_service_accounts="serviceaccount1@project1234.iam.gserviceaccount.com,uuid123,..."

$ vault write auth/gcp/role/prod-gce-role \
    type="gce" \
    project_id="project-123456" \
    policies="prod" \
    bound_zone="us-central1-a" \
    bound_instance_group="my-gce-group"
```

Users will login to Vault under roles (Vault-specific, unrelated to GCP IAM roles).
Roles are associated with an authentication type using the `type` parameter. Currently
`iam` or `gce` are acceptable values. The role also manages the [policies](/docs/concepts/policies.html)
the logged-in user will have and a set of constraints on authenticated users.
These are determined by the authenticating entity type and confirmed with GCP.

For the complete list of role options, please see the [API documentation](/api/auth/gcp/index.html#create-role).


## Authentication
These following examples assume the auth method is enabled at `/auth/gcp`
and use the CLI. For more complete documentation of parameters, see the
[HTTP API documentation](/api/auth/gcp/index.html#login).

There are two ways to login using the CLI:

### Client-provided JWT

In this case, you already have the GCE VM identity metadata token
or have obtained a IAM service account JWT, self-signed
or Google-signed using the GCP IAM API method
[projects.serviceAccounts.signJwt](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt).

```sh
$ vault write auth/gcp/login \
    role="dev-role" \
    jwt="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

* `role` is the role to login under
* `jwt` is the GCE metadata token or IAM service account token

See docs on the [IAM workflow](#iam) or [GCE workflow](#gce) to see how to obtain this JWT.

### (IAM ONLY) Client-side helper for service account JWT
Because the process to sign a service account JWT can be tedious,
we implemented a helper that can be run from client-side.

```sh
$ vault login -method=gcp \
    mount="gcp"
    role="dev-role" \
    jwt_exp
    credentials="@/path/to/creds.json" \
    project="project-for-authenticating-service-account" \
    service-account="authenticating-service-account@project.iam.gserviceaccounts.com"
```

This signs a properly formatted service account JWT for you and logs into Vault
directly. For more information, run `vault auth help gcp`.

Required Parameters:

- `role` `(string: <required>)` - The role with type=`iam` to login under.

Optional Parameters:

* `mount` (`string`: ``"gcp"`) - Name of the role.
* `jwt_exp` (`int`: `15`) - Number of minutes within which the JWT should expire.
        This defaults to the default value used when creating a role, and should only be set
        if set differently in the role (i.e. if the admin setting up your Vault server only accepts
        JWTs that expire within < 15 min).
* `credentials`: If not provided, attempts to use Application Default Credentials.
        These credentials are ONLY used to sign the service account JWT, i.e. call API
        method [projects.serviceAccounts.signJwt](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt),
        and thus must have IAM permission `iam.serviceAccounts.signJWT` or role `iam.serviceAccountTokenCreator`
        on the `service_account`, even if they are the same service account (**NOTE** this permission is not enabled by default).
                                          ```
* `service_account` `(string: `""`)`: Service account to sign JWT for i.e. login to Vault as
        and is listed under roles `bound_service_accounts`. If not provided, uses `client_email`
        from the credentials JSON. Only provide if the credentials are for a different account, acting as a "signer".
* `project`: Project of service account. If not provided, attempts to use `project_id` from the credentials JSON.
        Fails if not found and not provided by Application Default Credentials. Only provide if the credentials are for a
        different account, acting as a "signer". 

## Authentication Workflow

### IAM

The Vault authentication workflow for IAM service accounts is as follows:

  1. A client with IAM service account credentials generates a signed JWT using the IAM [projects.serviceAccounts.signJwt](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt) method. See [here](#the-iam-authentication-token) for the expected format and example code.
  2. The client sends this JWT to Vault in a login request with a role name. This role should have type `iam`.
  3. Vault grabs the `kid` header value, which contains the ID of the key-pair used to generate the JWT, and the `sub` ID/email to find the service account key. If the service account does not exist or the key is not linked to the service account, Vault will deny authentication.
  4. Vault authorizes the confirmed service account against the given role. See [authorization section](#authorization-workflow) to see how each type of role handles authorization.

[![IAM Login Workflow](/assets/images/vault-gcp-iam-auth-workflow.svg)](/assets/images/vault-gcp-iam-auth-workflow.svg)

#### The `iam` Authentication Token

The expected format of the JWT payload is as follows:

```json
{
  "sub" : "[SERVICE ACCOUNT IDENTIFIER]",
  "aud" : "vault/[ROLE NAME]",
  "exp" : "[EXPIRATION]"
}
```

- `[SERVICE ACCOUNT ID OR EMAIL]`: Either the email or the unique ID of a service account.

- `[ROLE NAME]`: Name of the role that this token will be used to login against. The full expected `aud` string should be "vault/$roleName".

- `[EXPIRATION]` : A [NumericDate](https://tools.ietf.org/html/rfc7519#section-2) value (seconds from Epoch). This value must be before the max JWT expiration allowed for a role (see `max_jwt_exp` parameter for creating a role). This defaults to 15 minutes and cannot be more than a hour.

**Note:** By default, we enforce a shorter `exp` period than the default length
for a given token (1 hour) in order to make reuse of tokens difficult. You can
customize this value for a given role but it will be capped at an hour.

To generate this token, we use the Google IAM API method [projects.serviceAccounts.signJwt](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt).
See an [example of how to generate this token](#generating-iam-token).

### GCE

The Vault authentication workflow for GCE instances is as follows:

  1. A client logins into a GCE instances and [obtains an instance identity metadata token](https://cloud.google.com/compute/docs/instances/verifying-instance-identity).
  2. The client request to login using this token (a JWT) and gives a role name to Vault.
  3. Vault uses the `kid` header value, which contains the ID of the key-pair used to generate the JWT, to find the OAuth2 public cert
  to verify this JWT.
  4. Vault authorizes the confirmed instance against the given role. See the [authorization section](#authorization-workflow) to see how each type of role handles authorization.

[![GCE Login Workflow](/assets/images/vault-gcp-gce-auth-workflow.svg)](/assets/images/vault-gcp-gce-auth-workflow.svg)

#### The `gce` Authentication Token

The token can be obtained from the `service-accounts/default/identity` endpoint for a instance's
[metadata server](https://cloud.google.com/compute/docs/storing-retrieving-metadata). You can use the
[example of how to obtain an instance metadata token](#generating-gce-token) to get started.

Learn more about the JWT format from the
[documentation](https://cloud.google.com/compute/docs/instances/verifying-instance-identity#token_format)
for the identity metadata token. The params the user provides are:

* `[AUD]`: The full expected `aud` string should end in "vault/$roleName". Note that Google requires the `aud`
    claim to contain a scheme or authority but Vault will only check for a suffix.
* `[FORMAT]`: MUST BE `full` for Vault. Format of the metadata token generated (`standard` or `full`).

### Examples for Obtaining Auth Tokens

#### Generating IAM Token

If you don't want to use the CLI helper, you can generate the JWT using the following methods.

**HTTP Request Example**

This uses [Google API HTTP annotation](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto).
Note the `$PAYLOAD` must be a marshaled JSON string with escaped double quotes.

```sh
#!/bin/sh
# [START PARAMS]
ROLE="test-role"
PROJECT="project-123456"
SERVICE_ACCOUNT="my-account@project-123456.iam.gserviceaccount.com"
OAUTH_TOKEN=$(oauth2l header cloud-platform)
# [END PARAMS]


PAYLOAD=$(echo "{ \"aud\": \"vault/$ROLE\", \"sub\": \"$SERVICE_ACCOUNT\"}" | sed -e 's/"/\\&/g')
curl -H "$OAUTH_TOKEN" \
    -H "Content-Type: application/json" \
    -X POST -d "{\"payload\":\"$PAYLOAD\"}" https://iam.googleapis.com/v1/projects/$PROJECT/serviceAccounts/$SERVICE_ACCOUNT:signJwt
```

**Golang Example**

We use the Go OAuth2 libraries, GCP IAM API, and Vault API. The example generates a token valid for the `dev-role` role (as indicated by the `aud` field of `jwtPayload`).

```go
// Abbreviated imports to show libraries.
import (
	vaultapi "github.com/hashicorp/vault/api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
	...
)

func main() {
	// Start [PARAMS]
	project := "project-123456"
	serviceAccount := "myserviceaccount@project-123456.iam.gserviceaccount.com"
	credsPath := "path/to/creds.json"

	os.Setenv("VAULT_ADDR", "https://vault.mycompany.com")
	defer os.Setenv("VAULT_ADDR", "")
	// End [PARAMS]

	// Start [GCP IAM Setup]
	jsonBytes, err := ioutil.ReadFile(credsPath)
	if err != nil {
		log.Fatal(err)
	}
	config, err := google.JWTConfigFromJSON(jsonBytes, iam.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := config.Client(oauth2.NoContext)
	iamClient, err := iam.New(httpClient)
	if err != nil {
		log.Fatal(err)
	}
	// End [GCP IAM Setup]

	// 1. Generate signed JWT using IAM.
	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", project, serviceAccount)
	jwtPayload := map[string]interface{}{
		"aud": "vault/dev-role",
		"sub": serviceAccount,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	payloadBytes, err := json.Marshal(jwtPayload)
	if err != nil {
		log.Fatal(err)
	}
	signJwtReq := &iam.SignJwtRequest{
		Payload: string(payloadBytes),
	}

	resp, err := iamClient.Projects.ServiceAccounts.SignJwt(resourceName, signJwtReq).Do()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Send signed JWT in login request to Vault.
	vaultClient, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	vaultResp, err := vaultClient.Logical().Write(
		"auth/gcp/login",
		map[string]interface{}{
			"role": "test",
			"jwt":  resp.SignedJwt,
	})

	if err != nil {
		log.Fatal(err)
	}

	// 3. Use auth token from response.
	log.Println("Access token: %s", vaultResp.Auth.ClientToken)
	vaultClient.SetToken(vaultResp.Auth.ClientToken)
	// ...
}
```

#### Generating GCE Token

**HTTP Request Example**

This uses [Google API HTTP annotation](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)
and must be run on a GCE instance.

```sh
# [START PARAMS]
VAULT_ADDR="https://127.0.0.1:8200/"
ROLE="my-gce-role"
SERVICE_ACCOUNT="default" # replace with an instance's service account if needed
# [END PARAMS]

curl -H "Metadata-Flavor: Google"\
     -G
     --data-urlencode "audience=$VAULT_ADDR/vault/$ROLE"\
     --data-urlencode "format=full" \
     "http://metadata/computeMetadata/v1/instance/service-accounts/$SERVICE_ACCOUNT/identity"
```

## Authorization Workflow

For `gcp`, login is per-role. Each role has a specific set of restrictions that
an authorized entity must fit in order to login. These restrictions are specific
to the role type.

Currently supported role types are:

* `iam` (Supports both IAM and inference for GCE tokens)
* `gce` (Only supports GCE tokens)

Vault validates an authenticated entity against the role and uses the role to
determine information about the lease, including Vault policies assigned and
TTLs. For a full list of accepted restrictions, see [role API docs](/api/auth/gcp/index.html#create-role).

If a GCE token is provided for login under an `iam` role, the service account associated with the token
(`sub` claim) is inferred and used to login.

## Usage

### Via the CLI.

#### Enable GCP authentication in Vault

```sh
$ vault auth enable gcp
```

#### Configure the GCP Auth Method

```sh
$ vault write auth/gcp/config credentials=@path/to/creds.json
```

**Configuration**: This includes GCP credentials Vault will use these to make calls to
GCP APIs. If credentials are not configured or if the user explicitly sets the
config with no credentials, the Vault server will attempt to use
[Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)
as set on the Vault server.

See [API documentation](/api/auth/gcp/index.html#configure)
to learn more about parameters.

#### Create a role

```sh
$ vault write auth/gcp/role/dev-role \
    type="iam" \
    project_id="project-123456" \
    policies="prod,dev" \
    bound_service_accounts="serviceaccount1@project1234.iam.gserviceaccount.com,uuid123,..."
    ...
```

**Roles**: Roles are associated with an authentication type/entity and a set of
Vault [policies](/docs/concepts/policies.html). Roles are configured with constraints
specific to the authentication type, as well as overall constraints and
configuration for the generated auth tokens.

We also expose a helper path for updating the service accounts attached to an existing `iam` role:

```sh
$ vault write auth/gcp/role/iam-role/service-accounts \
    add='serviceAccountToAdd,...' \
    remove='serviceAccountToRemove,...' \
```

and for updating the labels attached to an existing `gce` role:

```sh
$ vault write auth/gcp/role/gce-role/labels \
    add='label1:value1,foo:bar,...' \
    remove='key1,key2,...' \
```


See [API docs](/api/auth/gcp/index.html#create-role) to view
parameters for role creation and updates.

#### Login to get a Vault Token

Once the backend is setup and roles are registered with the backend,
the user can login against a specific role.

```sh
$ vault write auth/gcp/login role='dev-role' jwt='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
```

The `role` and `jwt` parameters are required. These map to the name of the
role to login against, and the signed JWT token for authenticating a role
respectively. The format of the provided JWT differs depending on the
authenticating entity.

### Via the API

#### Enable GCP authentication in Vault

```sh
$ curl $VAULT_ADDR/v1/sys/auth/gcp -d '{ "type": "gcp" }'
```

#### Configure the GCP Auth Method

```sh
$ curl $VAULT_ADDR/v1/auth/gcp/config \
-d '{  "credentials": "{...}" }'
```

#### Create a role

```sh
$ curl $VAULT_ADDR/v1/auth/gcp/role/dev-role \
-d '{ "type": "iam", "project_id": "project-123456", ...}'
```

#### Login to get a Vault Token

The endpoint for the GCP login is `auth/gcp/login`.

The `gcp` mountpoint value in the url is the default mountpoint value.
If you have mounted the `gcp` backend with a different mountpoint, use that value.

The `role` and `jwt` should be sent in the POST body encoded as JSON.

```sh
$ curl $VAULT_ADDR/v1/auth/gcp/login \
    -d '{ "role": "dev-role", "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }'
```

The response will be in JSON. For example:

```json
{
    "auth":{
        "client_token":"f33f8c72-924e-11f8-cb43-ac59d697597c",
        "accessor":"0e9e354a-520f-df04-6867-ee81cae3d42d",
        "policies":[
            "default",
            "dev",
            "prod"
        ],
        "metadata":{
            "role": "dev-role",
            "service_account_email": "dev1@project-123456.iam.gserviceaccount.com",
            "service_account_id": "111111111111111111111"
        },
        "lease_duration":2764800,
        "renewable":true
    },
    ...
}
```

### Contributing

This plugin is developed in a separate Github repository: [`hashicorp/vault-plugin-auth-gcp`](https://github.com/hashicorp/vault-plugin-auth-gcp). Please file all feature requests, bugs, and pull requests specific to the GCP plugin under that repository.

## API

The GCP Auth Plugin has a full HTTP API. Please see the
[API docs](/api/auth/gcp/index.html) for more details.
