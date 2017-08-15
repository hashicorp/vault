---
layout: "docs"
page_title: "Auth Plugin Backend: GCP"
sidebar_current: "docs-auth-gcp"
description: |-
  The gcp backend plugin allows automated authentication of AWS entities.
---



# Auth Plugin: gcp

The `gcp` backend plugin allows authentication against Vault using
Google credentials. It treats GCP as a Trusted Third Party and expects a
[JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519) signed by Google
credentials from the authenticating entity. This token can be generated through
different GCP APIs depending on the type of entity

Currently supports authentication for:

  * GCP IAM service accounts (`iam`)

We will update the documentation as we introduce more supported entities.

**Note**: The `gcp` backend is implemented as a
[Vault plugin](/docs/internals/plugins.html) backend. You must be using Vault
v0.8.0 to use plugins. The following documentation assumes that the backend has
 been [mounted](/docs/plugin/index.html) at `auth/gcp`.


## Quick start

The backend should initially be set up with:

* **Config**: This includes GCP credentials Vault will use these to make calls to
GCP APIs. If credentials are not configured or if the user explicitly sets the
config with no credentials, the Vault server will attempt to use
[Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)
as set on the Vault server.

```sh
$ vault write auth/gcp/config credentials=@path/to/creds.json
```

* **Roles**: Roles are associated with an authentication type/entity and a set of
Vault [policies](/docs/concepts/policies.html). Roles are configured with constraints
specific to the authentication type, as well as overall constraints and
configuration for the generated auth tokens.

Example IAM role:

```sh

$ vault write auth/gcp/role/myiamrole \
    type="iam" \
    project="project1234" \
    policies="prod,dev" \
    service_accounts="serviceaccount1@project1234.iam.gserviceaccount.com,uuid123,..."
    ...

```

Once the backend is setup and roles are registered with the backend,
the user can login against a specific role:

```sh

$ vault write auth/gcp/login role=myiamrole jwt=token

```

The format of the provided JWT differs depending on the authenticating entity.

## Login Flows
### Overview: IAM
The Vault authentication workflow for IAM service accounts is as follows
(or see diagram below):

  1. A client with IAM service account credentials generates a signed JWT using the IAM [projects.serviceAccounts.signJwt](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt) method. See [usage](#iam-authentication-token) for the expected format and example code.
  2. The client sends this JWT to Vault in a login request with a role name. This role should have type `iam`
  3. Vault grabs the `kid` header value, which contains the ID of the key-pair used to generate the JWT, and the `sub` ID/email to find the service account key. If the service account does not exist or the key is not linked to the service account, Vault will deny authentication.
  4. Vault authorizes the confirmed service account against the given role. See [authorization section](#authorization) to see how each type of role handles authorization.
[![IAM Login Workflow](/assets/images/vault-gcp-iam-auth-workflow.svg)](/assets/images/vault-gcp-iam-auth-workflow.svg)

## Authorization
For `gcp`, login is per-role. Each role has a specific set of restrictions that
an authorized entity must fit in order to login. These restrictions are specific
to the role type.

Currently supported role types are: `iam`

Vault validates an authenticated entity against the role and uses the role to
determine information about the lease, including Vault policies assigned and
TTLs. For a full list of accepted arguments, see [role API docs](/api/auth/gcp/index.html#create-role)

### `iam` Roles

`iam` roles support the following role restrictions:

  * `service_accounts`: A list of service accounts that are allowed to login as
    this role. This list accepts either service account IDs or emails.

**Note**: In the future we hope to support restrictions related to
IAM policies/permissions. With current IAM APIs, querying all IAM roles or permissions
(set or inherited) on an entity is very difficult to do and overall ends with not a
great story. We plan to wait until more functionality is supported before
attempting to add permission-based constraints.

## Usage
### Config
```
vault write auth/gcp/config \
  credentials=...
```

See [API documentation](/api/auth/gcp/index.html#configure)
to learn more about parameters.

### Roles

```sh
vault write auth/gcp/role/dev-role \
  type='iam' \
  project_id='project-123456' \
  policies='dev,prod' \
  service_accounts='devAccount1@project-123456.iam.gserviceaccounts.com,...` \
  ...

vault read auth/gcp/role/dev-role

vault list auth/gcp/role
```

We also expose a helper path for updating the service accounts attached to
an existing `iam` role.

```sh
vault write auth/gcp/role/myrole/add-service-accounts \
  values='serviceAccountToAdd,...'

vault write auth/gcp/role/myrole/delete-service-accounts \
    values='serviceAccountToRemove,...' \
  ...
```

See [API docs](/api/auth/gcp/index.html#create-role) to view
parameters for role creation and updates.

### Generating Authentication Tokens

#### `iam` Authentication Token
The expected format of the JWT payload is as follows:

```json
{
  "sub" : "[SERVICE ACCOUNT IDENTIFIER]",
  "aud" : "vault/[ROLE NAME]",
  "exp" : "[EXPIRATION]",
}
```

Values:
  * `[SERVICE ACCOUNT ID OR EMAIL]`: Either the email or the unique ID of a service account.
  * `[ROLE NAME]`: Name of the role that this token will be used to login against. The full expected `aud` string should be "vault/$roleName".
  * `[EXPIRATION]` : A [NumericDate](https://tools.ietf.org/html/rfc7519#section-2) value (seconds from Epoch). This value must be before the max JWT expiration allowed for a role (see `max_jwt_exp` parameter for creating a role). This defaults to 15 minutes and cannot be more than a hour.

**Note:** By default, we enforce a shorter `exp` period than the default length
for a given token (1 hour) in order to make reuse of tokens difficult. You can
customize this value for a given role but it will be capped at an hour.

To generate this token, we use the IAM API method.

**HTTP Request Example**

This uses [Google API HTTP annotation](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto). Note the `$PAYLOAD` must be a marshaled JSON string with escaped double quotes.

```sh

curl -H "Authorization: Bearer $OAUTH_TOKEN" \
    -H "Content-Type: application/json" \
    -X POST -d '{"payload":"'$PAYLOAD'"}' https://iam.googleapis.com/v1/projects/$PROJECT/serviceAccounts/$SERVICE_ACCOUNT:signJwt
```

**Golang example**:
We use the Go OAuth2 libraries, GCP IAM API, and Vault API.

```go
# Abbreviated imports to show libraries.
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
		"aud": "auth/gcp/login",
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

### Login
```sh
vault write auth/gcp/login \
  role=$ROLE \
  jwt=$JWT
...
```
Parameters:
  * `role`: Required. Name of the role to login against.
  * `jwt`: Required. A signed JWT token for authenticating a role. Format
    depends on entity type.

### Other

#### Plugin Setup

Assuming you have saved the binary `vault-plugin-auth-gcp` to some folder
and configured the [plugin directory](/docs/internals/plugins.html#plugin-directory)
for your server at `path/to/plugins`:

```sh
# Write plugin to catalog
$ vault write sys/plugins/catalog/gcp-auth command='vault-plugin-auth-gcp' sha_256=$HASH
Success! Data written to: sys/plugins/catalog/gcp-auth

# Enable plugin backend for auth.
$ vault auth-enable -path=gcp -plugin-name=gcp-auth plugin
Successfully mounted plugin 'plugin' at 'gcp'!
```

#### API

The GCP Auth Plugin has a full HTTP API. Please see the
[API docs](/api/auth/gcp/index.html) for more details.

## Contributing
This plugin is developed in a separate Github repository: [`hashicorp/vault-plugin-auth-gcp`](https://github.com/hashicorp/vault-plugin-auth-gcp).
Please file all feature requests, bugs, and pull requests specific to the GCP plugin.
