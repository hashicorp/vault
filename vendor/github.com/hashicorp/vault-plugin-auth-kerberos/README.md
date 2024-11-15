# Vault Plugin: Kerberos Auth Backend

## Details

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin allows for users to authenticate with Vault via Kerberos/SPNEGO.

## Usage

### Authentication

You can authenticate by posting a valid SPNEGO Negotiate header to `/v1/auth/kerberos/login`.

```python
try:
    import kerberos
except:
    import winkerberos as kerberos
import requests

service = "HTTP@vault.domain"
rc, vc = kerberos.authGSSClientInit(service=service, mech_oid=kerberos.GSS_MECH_OID_SPNEGO)
kerberos.authGSSClientStep(vc, "")
kerberos_token = kerberos.authGSSClientResponse(vc)

r = requests.post("https://vault.domain:8200/v1/auth/kerberos/login",
                  headers={'authorization': 'Negotiate ' + kerberos_token})
print('Vault token:', r.json()['auth']['client_token'])
```

### Configuration

1. Install and register the plugin.

Put the plugin binary (`vault-plugin-auth-kerberos`) into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://developer.hashicorp.com/vault/docs/configuration#plugin_directory)
in the Vault config used to start the server.

```hcl
plugin_directory = "path/to/plugin/directory"
```

```sh
$ SHA256=$(shasum -a 256 'vault-plugin-auth-kerberos' | cut -d ' ' -f2)
$ vault plugin register \
        -sha256=$SHA256 \
        -command="vault-plugin-auth-kerberos" \
        -client-cert server.crt -client-key server.key \
        auth kerberos
```

2. Enable the Kerberos auth method:

```sh
$ vault auth enable -passthrough-request-headers=Authorization -allowed-response-headers=www-authenticate kerberos
Success! Enabled kerberos auth method at: kerberos/
```

3. Use the /config endpoint to configure Kerberos.

Create a keytab for the kerberos plugin:
```sh
$ ktutil
ktutil:  addent -password -p your_service_account@REALM.COM -e aes256-cts -k 1
Password for your_service_account@REALM.COM:
ktutil:  list -e
slot KVNO Principal
---- ---- ---------------------------------------------------------------------
   1    1            your_service_account@REALM.COM (aes256-cts-hmac-sha1-96)
ktutil:  wkt vault.keytab
```

The KVNO (`-k 1`) should match the KVNO of the service account. An error will show in the vault logs if this is incorrect.

Different encryption types can also be added to the keytab, for example `-e rc4-hmac` with additional `addent` commands.

Then base64 encode it:
```sh
base64 vault.keytab > vault.keytab.base64
```

```sh
vault write auth/kerberos/config keytab=@vault.keytab.base64 service_account="your_service_account"
```

4. Add a SPNs (Service Principal Names) to your KDC for your service and service account. This should map the vault service to the account it is running as:
```sh
# for Windows/Active Directory
setspn.exe -U -S HTTP/vault.domain:8200 your_service_account
setspn.exe -U -S HTTP/vault.domain your_service_account
```

5. Configure LDAP backend to look up Vault policies.
Configuration for LDAP is identical to the [LDAP](https://developer.hashicorp.com/vault/docs/auth/ldap)
auth method, but writing to to the Kerberos endpoint:

```sh
vault write auth/kerberos/config/ldap @vault-config/auth/ldap/config
vault write auth/kerberos/groups/example-role @vault-config/auth/ldap/groups/example-role
```

In non-kerberos mode, the LDAP bind and lookup works via the user that is currently trying to authenticate.
If you're running LDAP together with Kerberos you might want to set a binddn/bindpass in the ldap config.

## Developing

To run a development environment through Docker, use:

```sh
make dev-env
```

This will:
- Build the current local plugin code
- Start Vault in a Docker container
- Start a local Samba container to function as the domain server
- Start a local joined container that can be used for login testing
- Output a number of variables for you to export in your working terminal

Note: Press CTRL+C in your `make dev-env` window when you'd like to stop and tear down your 
dev environment.

To begin testing in a separate window, after exporting the variables given in `make dev-env`:

```sh
VAULT_PLUGIN_SHA=$(openssl dgst -sha256 pkg/linux_amd64/vault-plugin-auth-kerberos|cut -d ' ' -f2)
vault write sys/plugins/catalog/auth/kerberos sha_256=${VAULT_PLUGIN_SHA} command="vault-plugin-auth-kerberos"
vault auth enable \
    -path=kerberos \
    -passthrough-request-headers=Authorization \
    -allowed-response-headers=www-authenticate \
    vault-plugin-auth-kerberos
vault write auth/kerberos/config \
    keytab=@vault_svc.keytab.base64 \
    service_account="vault_svc"
vault write auth/kerberos/config/ldap \
    binddn=${DOMAIN_VAULT_ACCOUNT}@${REALM_NAME} \
    bindpass=${DOMAIN_VAULT_PASS} \
    groupattr=sAMAccountName \
    groupdn="${DOMAIN_DN}" \
    groupfilter="(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))" \
    insecure_tls=true \
    starttls=true \
    userdn="CN=Users,${DOMAIN_DN}" \
    userattr=sAMAccountName \
    upndomain=${REALM_NAME} \
    url=ldaps://${SAMBA_CONTAINER:0:12}.${DNS_NAME}
```

To authenticate, first drop into a Docker container, then its Python shell:

```sh
docker exec -it $DOMAIN_JOINED_CONTAINER /bin/bash
python
```

Revisit the VAULT_CONTAINER_PREFIX outputted earlier, as you'll need it below:

```
prefix = '<insert VAULT_CONTAINER_PREFIX here>'
import kerberos
import requests

host = prefix + ".matrix.lan:8200"
service = "HTTP@{}".format(host)
rc, vc = kerberos.authGSSClientInit(service=service, mech_oid=kerberos.GSS_MECH_OID_SPNEGO)
kerberos.authGSSClientStep(vc, "")
kerberos_token = kerberos.authGSSClientResponse(vc)

r = requests.post("http://{}/v1/auth/kerberos/login".format(host),
                  headers={'Authorization': 'Negotiate ' + kerberos_token})
print('Vault token:', r.json()['auth']['client_token'])
```

### Tests

If you are developing this plugin and want to verify it is still
functioning (and you haven't broken anything else), we recommend
running the tests.

To run the tests, invoke `make test`:

```sh
$ make test
```

You can also specify a `TESTARGS` variable to filter tests like so:

```sh
$ make test TESTARGS='--run=TestConfig'
```

### Acceptance Tests

Acceptance tests requires a Vault Enterprise license to be 
[provided](https://developer.hashicorp.com/vault/docs/commands#vault_license) through 
`VAULT_LICENSE` and the following tools to be installed:
- [Docker](https://docs.docker.com/get-docker/)
- [jq](https://stedolan.github.io/jq/)
- [bats](https://bats-core.readthedocs.io/en/stable)


Run the acceptance tests:

```sh
$ make test-acceptance VAULT_LICENSE=<vault-license>
```


## Contributors

  1. Clone the repo
  2. Make changes on a branch
  3. Test changes
  4. Submit a Pull Request to GitHub

Maintained with :heart: by Hashicorp.

With thanks to the original creators of this plugin:
  - [wintoncode](https://github.com/wintoncode)
  - @ah-
  - @sambott
  - @roederja2
  - @jcmturner
  - @kristian-lesko
