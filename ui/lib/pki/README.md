# Vault PKI

Welcome to the Vault PKI (Ember) Engine! Below is an overview of PKI and resources for how to get started working within this engine.

## About PKI

> Public Key Infrastructure (PKI) is a system of processes, technologies, and policies that allows you to encrypt and sign data. (source: [digicert.com](https://www.digicert.com/what-is-pki))

The [Vault PKI Secrets Engine](https://developer.hashicorp.com/vault/api-docs/secret/pki) allows security engineers to [create a chain of PKI certificates](https://developer.hashicorp.com/vault/tutorials/secrets-management/pki-engine) much easier than they would with traditional workflows.

## About the UI engine

If you couldn't tell from the documentation above, PKI is _complex_. As such, the data doesn't map cleanly to a CRUD model and so the first thing you might notice is that the models and adapters for PKI (which [live in the main app](https://ember-engines.com/docs/addons#using-ember-data), not the engine) have some custom logic that differentiate it from most other secret engines. Below are the models used throughout PKI and how they are used to interact with the mount. Aside from `pki/action`, each model has a corresponding tab in the UI that takes you to its `LIST` view.

- ### [pki/action](../../app/models/pki/action.js)

  This model is used to perform different `POST` requests that receive similar parameters but don't create a single item (which would be a record in Ember data). These various actions may create multiple items, each with different parameters than those submitted in the `POST` request. For example:

  > - `pki/generate/root/:type` creates a new self-signed CA certificate (an issuer) and private key, which is only returned if `type = exported`.
  > - `pki/issuer/:issuer_ref/sign-intermediate` creates a certificate, and returns issuing CA and CA chain data that is only available once

  The `pki/action`[adapter](../../app/adapters/pki/action.js) is used to map the desired action to the corresponding endpoint, and the `pki/action` [serializer](../../app/serializers/pki/action.js) is leveraged to only send relevant attributes. The following PKI workflows use this model:

  - [Root generation and rotation](https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-root) generates a new self-signed CA certificate (an issuer) and private key (key data is only returned if type = `exported`)
  - [Import CA cert and keys](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-ca-certificates-and-keys)
  - [Generate intermediate CSR](https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-intermediate-csr)
  - [Sign intermediate](https://developer.hashicorp.com/vault/api-docs/secret/pki#sign-intermediate)

- ### [pki/certificate/base](../../app/models/pki/certificate/base.js)

  This model is for specific interactions with certificate data. The base model contains attributes that make up a certificate's content. The other models that extend this model [certificate/generate](../../app/models/pki/certificate/generate.js) and [certificate/sign](../../app/models/pki/certificate/sign.js) include additional attributes to perform their relevant requests.

  The `parsedCertificate` attribute is an object that houses all of the parsed certificate data returned by the [parse-pki-cert.js](../../app/utils/parse-pki-cert.js) util.

> The following models more closely follow a CRUD pattern:

- ### [pki/issuer](../../app/models/pki/issuer.js)

  > Issuers are created by the `pki/action` model, by [importing a CA](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-ca-certificates-and-keys) or [generating a root](https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-root)

  - [UPDATE](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-issuer-certificate)
  - [READ](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-issuer-certificate)
  - [LIST](https://developer.hashicorp.com/vault/api-docs/secret/pki#list-issuers)

- ### [pki/role](../../app/models/pki/role.js)

  - [CREATE/UPDATE](https://developer.hashicorp.com/vault/api-docs/secret/pki#create-update-role)
  - [READ](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-role)
  - [LIST](https://developer.hashicorp.com/vault/api-docs/secret/pki#list-roles)

- ### pki/key

  - `CREATE` has two options:
    - [GENERATE](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-ca-certificates-and-keys)
    - [IMPORT](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-key)
  - [READ](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-key)
  - [LIST](https://developer.hashicorp.com/vault/api-docs/secret/pki#list-keys)
