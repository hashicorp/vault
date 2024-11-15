# Elasticsearch Database Secrets Engine
This plugin provides unique, short-lived credentials for Elasticsearch using native X-Pack Security.

## Getting Started

To take advantage of this plugin, you must first enable Elasticsearch's native realm of security by activating X-Pack. These
instructions will walk you through doing this using ElasticSearch 7.1.1.

### Version compatibility

If you would like to install a different version of this plugin to that which is bundled with Vault, versions v0.6.0 onwards of this plugin are incompatible with Vault versions before 1.6.0 due to an update of the database plugin interface.

### Enable X-Pack Security in Elasticsearch

Read [Securing the Elastic Stack](https://www.elastic.co/guide/en/elastic-stack-overview/7.1/elasticsearch-security.html) and
follow [its instructions for enabling X-Pack Security](https://www.elastic.co/guide/en/elasticsearch/reference/7.1/setup-xpack.html).

### Enable Encrypted Communications

This plugin communicates with Elasticsearch's security API. In ES 7.1.1, you must enable TLS to consume that API.

To set up TLS in Elasticsearch, first read [encrypted communications](https://www.elastic.co/guide/en/elastic-stack-overview/7.1/encrypting-communications.html)
and go through its instructions on [encrypting HTTP client communications](https://www.elastic.co/guide/en/elasticsearch/reference/7.1/configuring-tls.html#tls-http).

After enabling TLS on the Elasticsearch side, you'll need to convert the .p12 certificates you generated to other formats so they can be
used by Vault. [Here is an example using OpenSSL](https://stackoverflow.com/questions/15144046/converting-pkcs12-certificate-into-pem-using-openssl)
to convert our .p12 certs to the pem format.

Also, on the instance running Elasticsearch, we needed to install our newly generated CA certificate that was originally in the .p12 format.
We did this by converting the .p12 CA cert to a pem, and then further converting that
[pem to a crt](https://stackoverflow.com/questions/13732826/convert-pem-to-crt-and-key), adding that crt to `/usr/share/ca-certificates/extra`,
and using `sudo dpkg-reconfigure ca-certificates`.

The above instructions may vary if you are not using an Ubuntu machine. Please ensure you're using the methods specific to your operating
environment. Describing every operating environment is outside the scope of these instructions.

### Set Up Passwords

When done, verify that you've enabled X-Pack by running `$ $ES_HOME/bin/elasticsearch-setup-passwords interactive`. You'll
know it's been set up successfully if it takes you through a number of password-inputting steps.

### Create a Role for Vault

Next, in Elasticsearch, we recommend that you create a user just for Vault to use in managing secrets.

To do this, first create a role that will allow Vault the minimum privileges needed to administer users and passwords by performing a
POST to Elasticsearch. To do this, we used the `elastic` superuser whose password we created in the
`$ $ES_HOME/bin/elasticsearch-setup-passwords interactive` step.

```
$ curl \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"cluster": ["manage_security"]}' \
    http://elastic:$PASSWORD@localhost:9200/_security/role/vault
```

Next, create a user for Vault associated with that role.

```
$ curl \
    -X POST \
    -H "Content-Type: application/json" \
    -d @data.json \
    http://elastic:$PASSWORD@localhost:9200/_security/user/vault
```

The contents of `data.json` in this example are:
```
{
 "password" : "myPa55word",
 "roles" : [ "vault" ],
 "full_name" : "Hashicorp Vault",
 "metadata" : {
   "plugin_name": "Vault Plugin Secrets ElasticSearch",
   "plugin_url": "https://github.com/hashicorp/vault-plugin-secrets-elasticsearch"
 }
}
```

Now, Elasticsearch is configured and ready to be used with Vault.

## Example Walkthrough

Here is an example of how to successfully configure and use this secrets engine using the Vault CLI. Note that the
`plugin_name` may need to be `vault-plugin-database-elasticsearch` if you manually mounted it rather than using the
version of the plugin built in to Vault.
```
export ES_HOME=/home/somewhere/Applications/elasticsearch-7.1.1

vault secrets enable database

vault write database/config/my-elasticsearch-database \
    plugin_name="elasticsearch-database-plugin" \
    allowed_roles="internally-defined-role,externally-defined-role" \
    username=vault \
    password=myPa55word \
    url=http://localhost:9200 \
    ca_cert=/usr/share/ca-certificates/extra/elastic-stack-ca.crt.pem \
    client_cert=$ES_HOME/config/certs/elastic-certificates.crt.pem \
    client_key=$ES_HOME/config/certs/elastic-certificates.key.pem

# create and get creds with one type of role
vault write database/roles/internally-defined-role \
    db_name=my-elasticsearch-database \
    creation_statements='{"elasticsearch_role_definition": {"indices": [{"names":["*"], "privileges":["read"]}]}}' \
    default_ttl="1h" \
    max_ttl="24h"

vault read database/creds/internally-defined-role

# create and get creds with another type of role
vault write database/roles/externally-defined-role \
    db_name=my-elasticsearch-database \
    creation_statements='{"elasticsearch_roles": ["vault"]}' \
    default_ttl="1h" \
    max_ttl="24h"

vault read database/creds/externally-defined-role

# renew credentials
vault lease renew database/creds/internally-defined-role/nvJ6SveX9PN1E4BlxVWdKuX1

# revoke credentials
vault lease revoke database/creds/internally-defined-role/nvJ6SveX9PN1E4BlxVWdKuX1

# rotate root credentials
vault write -force database/rotate-root/my-elasticsearch-database
```

## Developing

The Vault plugin system is documented on the [Vault documentation site](https://developer.hashicorp.com/vault/docs/plugins).

The local_dev.sh script will build the plugin, start vault, mount the plugin,
and run any custom commands in `./scripts/custom.sh`:

```console
./scripts/local_dev.sh
```

## Testing

### Unit and integration tests

```console
make test
```

Requires docker.

### Acceptance tests

To also run acceptance tests against an elasticsearch cluster, follow the
instructions at the top of [acceptance_test.go](./acceptance_test.go), then

```console
make testacc
```

Or to run only the acceptance tests:

```console
./scripts/run_acceptance.sh
```

See the comments in [run_acceptance.sh](./scripts/run_acceptance.sh) for
configuration information.

**Note**: The acceptance test for 6.8.13 generally won't pass on M1 Macs due to the lack of an arm64 Elasticsearch 6.8.13 image.
