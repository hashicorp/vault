# Vault IBM Cloud Auth Method: Developer Guide

This guide is for developers who **build and test** the `ibmcloud` auth method from source. For using the
auth method (login flows, roles, trusted profiles, Vault Agent) see
[`IBMCLOUD-AUTH-QUICKSTART.md`](IBMCLOUD-AUTH-QUICKSTART.md).

It covers:

1. Where the feature lives in the source tree
2. Building the server, with and without the Linux binary
3. Starting a dev server
4. [`create-admin-apikey.sh`](create-admin-apikey.sh): Vault's admin identity
5. [`create-tester-apikey.sh`](create-tester-apikey.sh): identities to log in with
6. [`demo-setup.sh`](demo-setup.sh): enable and configure the auth method in Vault
7. An end-to-end walkthrough and troubleshooting

All commands run from the repository root.

---

## 1. Where the feature lives

| Path | Purpose |
|---|---|
| [`builtin/credential/ibmcloud/`](builtin/credential/ibmcloud) | The auth backend: config, roles, login, JWKS validation, `rotate-root`, IAM client, CLI handler |
| [`helper/builtinplugins/registry_full.go`](helper/builtinplugins/registry_full.go) | Registers the backend factory (`credIBMCloud.Factory`) |
| [`command/server/commands_full.go`](command/server/commands_full.go) | Registers `credIBMCloud.CLIHandler` so `vault login -method=ibmcloud` works |
| [`command/client/auth/ibmcloud/`](command/client/auth/ibmcloud) | Auto-auth method for Vault Agent |

Both registry files carry the `!minimal` build constraint, so **a `BUILD_MINIMAL` build does not contain the
`ibmcloud` auth method**. Never set `BUILD_MINIMAL` when building for this feature.

---

## 2. Building the server

### 2.1 Prerequisites

- Go matching [`.go-version`](.go-version)
- `make`
- For the helper scripts: the [`ibmcloud` CLI](https://cloud.ibm.com/docs/cli) (logged in), `jq` and `curl`

### 2.2 Option A: build host binary only (no Linux binary)

Use this for local development and testing on your own machine.

```bash
make dev
```

Or, with the environment the project's build script uses (full build, no code generation):

```bash
unset BUILD_MINIMAL
CGO_ENABLED=0 BUILD_TAGS=testonly SKIP_GEN=1 make dev
```

Result:

- `./bin/vault`, for your current OS and architecture
- a copy in `$GOPATH/bin` (done by [`scripts/build.sh`](scripts/build.sh))

Notes:

- `make dev` adds the `testonly` build tag itself.
- `SKIP_GEN=1` skips `go generate` (`make prep`). Use it only when generated files are up to date. Drop it
  after changing `.proto` files or `go:generate` directives.
- `CGO_ENABLED` is `0` by default.

Verify that the auth method is compiled in:

```bash
./bin/vault version
./bin/vault auth help ibmcloud
```

### 2.3 Option B: build host binary and Linux binary

Use this when you want to run Vault on an IBM Cloud VSI or in a Linux container (see
[`deploy-vault-to-vsi.sh`](deploy-vault-to-vsi.sh)).

```bash
./build-vault-with-ibmcloud.sh
```

The script:

1. Checks it runs from the Vault repository root
2. Sets `CGO_ENABLED=0`, `BUILD_TAGS=testonly`, `SKIP_GEN=1` and unsets `BUILD_MINIMAL`
3. Runs `make dev` and produces `./bin/vault`
4. Checks `vault auth help ibmcloud` to prove the auth method is present
5. Cross-compiles **`./bin/vault-linux-amd64`** with `GOOS=linux GOARCH=amd64`

The Linux binary is not run by the script because the build host is not Linux, so it is not verified.

To build only the Linux binary, without the script:

```bash
GIT_COMMIT="$(scripts/ci-helper.sh revision)"
BUILD_DATE="$(scripts/ci-helper.sh date)"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X github.com/hashicorp/vault/version.GitCommit='${GIT_COMMIT}' -X github.com/hashicorp/vault/version.BuildDate=${BUILD_DATE}" \
    -o bin/vault-linux-amd64 \
    -tags testonly \
    ./command/server/cmd
```

For an ARM VSI use `GOARCH=arm64` and a different output name. Note that
[`deploy-vault-to-vsi.sh`](deploy-vault-to-vsi.sh) expects `bin/vault-linux-amd64`.

### 2.4 Which option?

| You want to | Use |
|---|---|
| Iterate on the backend and test on your laptop | `make dev` (Option A) |
| Check that the auth method is compiled in, plus produce a deployable Linux binary | `./build-vault-with-ibmcloud.sh` (Option B) |
| Deploy to a VSI to test trusted profile logins | Option B, then `./deploy-vault-to-vsi.sh` |

### 2.5 Unit tests

```bash
make test TEST=./builtin/credential/ibmcloud
make test TEST=./builtin/credential/ibmcloud TESTARGS="-run TestSpecificTest"
```

`make test` adds `testonly` automatically. See [`AGENTS.md`](AGENTS.md) for more build and test commands.

---

## 3. Starting a dev server

```bash
./restart-vault-dev.sh        # normal logging, log file /tmp/vault-dev.log
./restart-vault-debug.sh      # debug logging, log file /tmp/vault-dev-debug.log
```

Both scripts kill any running `vault server` process (`pkill -f "vault server"`), then start
`./bin/vault server -dev -dev-root-token-id=root` in the background. The debug variant logs IAM calls
(URLs, status codes, bodies, token claims), which is the first place to look when a login fails.

In each shell that talks to the server:

```bash
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=root
```

> The dev server is in-memory. Restarting it loses all auth config and roles. Re-run
> [`demo-setup.sh`](demo-setup.sh) after every restart.

---

## 4. `create-admin-apikey.sh`: Vault's admin identity

Vault needs its own **admin API key** to confirm that API-key identities belong to the account and to read
access group membership. This script creates a service ID with the least access needed.

### 4.1 Prerequisites

- `ibmcloud`, `jq`, `curl` installed
- `ibmcloud login` done against the target account and environment
- Permission to create service IDs, policies and API keys

### 4.2 Usage

```bash
./create-admin-apikey.sh [-n NAME] [-o FILE] [--skip-user-lookup] [--skip-verify]
```

| Option | Default | Meaning |
|---|---|---|
| `-n`, `--name` | `vault-ibmcloud-admin` | Service ID name. Must not exist yet |
| `-o`, `--output` | `$HOME/<NAME>.apikey.json` | File that receives the API key JSON (mode 600). Must not exist yet |
| `--skip-user-lookup` | off | Do not grant Viewer on `user-management`. Use when only service IDs will log in |
| `--skip-verify` | off | Do not test the new key against IAM |

| Environment | Default | Meaning |
|---|---|---|
| `ACCOUNT_ID` | account from `ibmcloud target` | Account to use |
| `IAM_ENDPOINT` | `https://iam.<host of 'ibmcloud api'>` | IAM endpoint to verify against |

### 4.3 What it does

1. Resolves the account and IAM endpoint, refuses to overwrite an existing service ID or key file
2. Creates the service ID
3. Grants three policies:

   | # | Role | Service | Scope | Used for |
   |---|---|---|---|---|
   | 1 | Viewer | `iam-identity` | all service IDs | Login: confirm a service ID belongs to the account |
   | 2 | Operator | `iam-identity` | resource type `serviceid`, only this service ID | `rotate-root`: create, read, delete this service ID's keys |
   | 3 | Viewer | `user-management` | account | Login: look up `User` identities (skipped by `--skip-user-lookup`) |

4. Creates the API key into the output file. The key is **never printed**
5. Verifies the key against IAM (token exchange, read own and another service ID, create/read/delete a key for
   itself, and confirm that creating a key for a **different** service ID is refused with 403)
6. On failure before the key exists, it rolls back by deleting the service ID. If only verification fails, the
   service ID and key are kept, because policies can take up to a minute to propagate. Re-test, or delete it
   with `ibmcloud iam service-id-delete <id> -f`

### 4.4 Example

```bash
ibmcloud login --sso
./create-admin-apikey.sh
# Done.
#   Service ID: ServiceId-xxxxxxxx-...
#   Account:    <account id>
#   Key file:   /Users/me/vault-ibmcloud-admin.apikey.json
```

The final output prints the exact `./demo-setup.sh` command to run next. Keep the printed **Service ID** and
**Account** values.

> If the `rotate-root` check prints `WARN ... worked only with iam_id=iam-<id>`, the IAM API accepted the
> `iam-` prefixed ID only, while Vault sends the bare service ID. `rotate-root` fails until
> `CreateAPIKey` in [`client.go`](builtin/credential/ibmcloud/client.go) sends the prefixed ID.

---

## 5. `create-tester-apikey.sh`: identities to log in with

Creates a **service ID that you log in to Vault with**. Run it several times with different names and groups to
build a test matrix.

### 5.1 Usage

```bash
./create-tester-apikey.sh [-n NAME] [-g ACCESS_GROUP] [-o FILE] [--skip-verify]
```

| Option | Default | Meaning |
|---|---|---|
| `-n`, `--name` | `vault-ibmcloud-tester` | Service ID name. Must not exist yet |
| `-g`, `--access-group` | none | Access group ID (`AccessGroupId-...`) or name. The service ID is added to it |
| `-o`, `--output` | `$HOME/<NAME>.apikey.json` | API key JSON file (mode 600). Must not exist yet |
| `--skip-verify` | off | Do not test the new key against IAM |

`ACCOUNT_ID` and `IAM_ENDPOINT` work as in the admin script.

### 5.2 With and without an access group

| Run | The service ID has | Role binding that matches it |
|---|---|---|
| no `-g` | no IAM policies, no group membership | `bound_service_ids=iam-<id>`, or only `bound_account_id` |
| with `-g GROUP` | membership in that group (and its access) | `bound_access_group_ids=<group id>` |

The access group is resolved by ID or name **before** anything is created, so a typo does not leave resources
behind.

### 5.3 What it verifies

- The key exchanges for an IAM token
- Token introspection shows: active, `sub_type` is `ServiceId`, the account matches, `iam_id` is `iam-<id>`
- With `-g`: the service ID's own group list contains the group, which is what Vault reads for
  `bound_access_group_ids` (retries for about 20 seconds)

### 5.4 Example

```bash
# Service ID with no access
./create-tester-apikey.sh -n tester-nogroup

# Service ID in an access group
./create-tester-apikey.sh -n tester-ingroup -g my-vault-group
```

The script prints the role and login commands that match the new identity. For the second example:

```bash
./bin/vault write auth/ibmcloud/role/tester-role \
    type=apikey \
    bound_account_id=<account id> \
    bound_service_ids=iam-<service id> \
    token_policies=default \
    bound_access_group_ids=<group id>

./bin/vault login -method=ibmcloud role=tester-role type=apikey \
    apikey="$(jq -r .apikey ~/tester-ingroup.apikey.json)" iam_endpoint=<iam endpoint>
```

Remove a tester when done:

```bash
ibmcloud iam service-id-delete <service id> -f && rm -f ~/tester-ingroup.apikey.json
```

> Role bindings use `iam-<service id>`, not the bare `ServiceId-...`.
> Roles created by [`demo-setup.sh`](demo-setup.sh) do not set `bound_service_ids`, so any identity in the
> account (and in the access group, if one was given) can log in to them.

---

## 6. `demo-setup.sh`: enable and configure the auth method

Sets up a running dev server: enables `auth/ibmcloud`, writes the config, creates two policies and two roles.

### 6.1 Usage

```bash
./demo-setup.sh <api_key> <account_id> [service_id] [access_group_id]
```

| Argument | Required | Meaning |
|---|---|---|
| `api_key` | yes | The admin API key (from `create-admin-apikey.sh`) |
| `account_id` | yes | Account every login must belong to |
| `service_id` | no | Service ID that owns `api_key`. Only needed for `rotate-root`. Pass `""` to skip it while still passing an access group |
| `access_group_id` | no | Binds both roles to this group via `bound_access_group_ids`. If omitted, login only requires account membership |

| Environment | Default | Meaning |
|---|---|---|
| `IAM_ENDPOINT` | `https://iam.test.cloud.ibm.com` | IAM endpoint written into the Vault config |
| `VAULT_ADDR`, `VAULT_TOKEN` | none | Must point to a running, unsealed Vault |

> **The default `IAM_ENDPOINT` is the IBM Cloud test environment.** For production IBM Cloud set
> `IAM_ENDPOINT=https://iam.cloud.ibm.com`. It must match the environment of the account, the admin key, and
> the `iam_endpoint` the client uses at login. A token from one environment is rejected by the other.

The script uses `./bin/vault` and reads `my-policy.hcl` and `my-writer-policy.hcl` from the current directory,
so run it from the repository root.

### 6.2 What it does

1. `vault auth enable ibmcloud` (continues if already enabled)
2. `vault write auth/ibmcloud/config` with `iam_endpoint`, `iam_client_timeout=30`, `api_key`, `account_id`
   and, if given, `service_id`
3. Writes policies from [`my-policy.hcl`](my-policy.hcl) (`my-policy`, read/list on `secret/data/myapp/*`) and
   [`my-writer-policy.hcl`](my-writer-policy.hcl) (`my-writer-policy`, full CRUD on `secret/data/myapp/*`)
4. Creates roles of `type=apikey`, with `bound_account_id`, TTL 1h, max TTL 24h:

   | Role | Policy |
   |---|---|
   | `my-service-role` | `my-policy` |
   | `my-service-writer-role` | `my-writer-policy` |

   `bound_access_group_ids` is added to both only if `access_group_id` was given.

The `api_key` is passed as a command-line argument, so it is visible in the process list while the script runs.
Use this only on development machines.

### 6.3 Examples

```bash
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=root
export IAM_ENDPOINT=https://iam.test.cloud.ibm.com   # or https://iam.cloud.ibm.com

KEY="$(jq -r .apikey ~/vault-ibmcloud-admin.apikey.json)"

# Account membership only, no rotate-root
./demo-setup.sh "$KEY" "$ACCOUNT_ID"

# Account membership + rotate-root
./demo-setup.sh "$KEY" "$ACCOUNT_ID" "$SID"

# Account membership + access group, no rotate-root
./demo-setup.sh "$KEY" "$ACCOUNT_ID" "" "$ACCESS_GROUP_ID"

# Everything
./demo-setup.sh "$KEY" "$ACCOUNT_ID" "$SID" "$ACCESS_GROUP_ID"
```

### 6.4 Re-running

Re-running is safe for the auth mount, policies and roles: the mount error is tolerated and the other writes
overwrite. To switch from "no access group" to "with access group", just run it again with the group ID.

---

## 7. End-to-end walkthrough

```bash
# 1. Build
./build-vault-with-ibmcloud.sh            # or: make dev

# 2. Start a dev server
./restart-vault-dev.sh
export VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root

# 3. IBM Cloud: admin identity for Vault
ibmcloud login --sso
export IAM_ENDPOINT=https://iam.test.cloud.ibm.com      # match your environment
./create-admin-apikey.sh
export ACCOUNT_ID=$(ibmcloud target --output json | jq -r '.account.guid')
export SID=$(ibmcloud iam service-id vault-ibmcloud-admin --output json | jq -r '.id')

# 4. Vault: enable and configure
./demo-setup.sh "$(jq -r .apikey ~/vault-ibmcloud-admin.apikey.json)" "$ACCOUNT_ID" "$SID"

# 5. IBM Cloud: an identity to log in with
./create-tester-apikey.sh -n my-tester

# 6. Log in using a role created by demo-setup.sh
./bin/vault login -method=ibmcloud role=my-service-writer-role type=apikey \
    apikey="$(jq -r .apikey ~/my-tester.apikey.json)" iam_endpoint="$IAM_ENDPOINT"

# 7. Check the policy: allowed under secret/data/myapp/*
VAULT_TOKEN=<token from step 6> ./bin/vault kv put -mount=secret myapp/test foo=bar

# 8. Optional: rotate the admin key so only Vault holds a valid copy, then delete the key file
./bin/vault write -f auth/ibmcloud/config/rotate-root
rm ~/vault-ibmcloud-admin.apikey.json
```

To exercise access groups, create a group, add the tester with `-g`, and run `demo-setup.sh` with the group ID
as the fourth argument. A tester created without `-g` must then be rejected on login.

---

## 8. Troubleshooting

| Symptom | Likely cause |
|---|---|
| `vault auth help ibmcloud` fails, or `auth enable ibmcloud` says unsupported | Binary built with `BUILD_MINIMAL`, or built from a branch without the registry entries (section 1). Rebuild |
| Login fails with a token or JWKS error | `IAM_ENDPOINT` mismatch between Vault config, admin key, and the client |
| Login rejected on account or access group | Admin key policies not yet propagated (wait about a minute), or `bound_account_id` / `bound_access_group_ids` do not match the identity |
| Access group binding seems ignored | Config has no `api_key` and `account_id`. Both are required for group and account checks |
| `rotate-root` fails | `service_id` was not given to `demo-setup.sh`, or the `iam-` prefix issue in section 4.4 |
| `create-*-apikey.sh`: "already exists" | Choose another `-n` or `-o`, or delete the old service ID and key file |
| `create-*-apikey.sh`: "could not determine the account" | Run `ibmcloud login` / `ibmcloud target -c <account>` or set `ACCOUNT_ID` |
| Config and roles disappeared | Dev server restarted, it is in-memory. Re-run `demo-setup.sh` |
| `demo-setup.sh`: policy file not found | Run it from the repository root |

For anything else, start the server with `./restart-vault-debug.sh` and read `/tmp/vault-dev-debug.log`.

---

## 9. Clean up

```bash
pkill -f "vault server"
ibmcloud iam service-id-delete <tester service id> -f
ibmcloud iam service-id-delete <admin service id> -f
rm -f ~/vault-ibmcloud-admin.apikey.json ~/my-tester.apikey.json
```
