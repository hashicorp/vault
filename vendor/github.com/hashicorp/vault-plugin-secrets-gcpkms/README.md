# Vault Secrets Engine for Google Cloud KMS

[![Build Status](https://travis-ci.com/hashicorp/vault-plugin-secrets-gcpkms.svg?token=xjv5yxmcgdD1zvpeR4me&branch=master)](https://travis-ci.com/hashicorp/vault-plugin-secrets-gcpkms)

This is a plugin backend for [HashiCorp Vault][vault] that manages [Google Cloud
KMS][kms] keys and provides pass-through encryption/decryption of data through
KMS.

**Please note:** Security is taken seriously. If you believe you have found a
security issue, **do not open an issue**. Responsibly disclose by contacting
security@hashicorp.com.


## Usage

The Google Cloud KMS Vault secrets engine is automatically bundled and included
in [Vault][vault] distributions. To activate the plugin, run:

```text
$ vault secrets enable gcpkms
```

Optionally configure the backend with GCP credentials:

```text
$ vault write gcpkms/config credentials="..."
```

Ask Vault to generate a new Google Cloud KMS key:

```text
$ vault write gcpkms/keys/my-key \
    key_ring=projects/my-project/locations/global/keyRings/my-keyring \
    crypto_key=my-crypto-key
```

This will create a KMS key in Google Cloud and requests to Vault will be
encrypted/decrypted with that key.

Encrypt some data:

```text
$ vault write gcpkms/encrypt/my-key plaintext="hello world"
```

Decrypt the data:

```text
$ vault write gcpkms/decrypt/my-key ciphertext="..."
```


## Development

This plugin is automatically distributed and included with Vault. **These
instructions are only useful if you want to develop against the plugin.**

- Modern [Go](https://golang.org) (1.11+)
- Git

1. Clone the repo:

    ```text
    $ git clone https://github.com/hashicorp/vault-plugin-secrets-gcpkms
    $ cd vault-plugin-secrets-gcpkms
    ```

1. Build the binary:

    ```text
    $ make dev
    ```

    The plugin binary will be written to the `./bin` directory.

1. Run Vault plugins from that directory:

    ```text
    $ vault server -dev -dev-plugin-dir=./bin
    $ vault secrets enable -path=gcpkms -plugin=vault-plugin-secrets-gcpkms plugin
    ```

### Documentation

The documentation for the plugin lives in the [main Vault
repository](//github.com/hashicorp/vault) in the `website/` folder. Please make any
documentation updates as separate Pull Requests against that repo.

### Tests

This plugin has both unit tests and acceptance tests. To run the acceptance
tests, you must:

- Have a service account in the project with the roles "Cloud KMS Admin" and "Cloud KMS Crypto Operator"
- Set `GOOGLE_APPLICATION_CREDENTIALS` to the service account key credentials for the above account
- Set `GOOGLE_CLOUD_PROJECT` to the name of the project
- Request an increase to the Cloud Key Management Service (KMS) API Write-Requests quota to 600 per minute

We recommend running tests in a dedicated Google Cloud project. On a fresh
project, you will need to enable the Cloud KMS API. This operation only needs to
be completed once per project.

```text
$ gcloud services enable cloudkms.googleapis.com --project $GOOGLE_CLOUD_PROJECT
```

After the API is enabled, it may take a few minutes to propagate. Please wait
and try again.

To run the tests:

```text
$ make test
```

**Warning:** the acceptance tests change real resources which may incur real
costs. Please run acceptance tests at your own risk.

### Cleanup

If a test panics or fails to cleanup, you can be left with orphaned KMS keys.
While their monthly cost is minimal, this may be undesirable. As such, there a
cleanup script is included. To execute this script, run:

```text
$ export GOOGLE_CLOUD_PROJECT=my-test-project
$ go run test/cleanup/main.go
```

**WARNING!** This will delete all keys in most key rings, so do not run this
against a production project!

[kms]: https://cloud.google.com/kms
[vault]: https://www.vaultproject.io
