# Vault PKI

Welcome to the Vault PKI (Ember) Engine! Below is an overview of PKI and resources for how to get started working within this engine.

## About PKI

> Public Key Infrastructure (PKI) is a system of processes, technologies, and policies that allows you to encrypt and sign data. (source: [digicert.com](https://www.digicert.com/what-is-pki))

The [Vault PKI Secrets Engine](https://developer.hashicorp.com/vault/api-docs/secret/pki) allows security engineers to [create a chain of PKI certificates](https://developer.hashicorp.com/vault/tutorials/secrets-management/pki-engine) much easier than they would with traditional workflows.

## About the UI engine

If you couldn't tell from the documentation above, PKI is _complex_. As such, the data doesn't map cleanly to a CRUD model and so the first thing you might notice is that the models and adapters for PKI (which [live in the main app](https://ember-engines.com/docs/addons#using-ember-data), not the engine) have some custom logic that differentiate it from most other secret engines. Below are the models used throughout PKI and how they are used to interact with the mount. Aside from `pki/action`, each model has a corresponding tab in the UI that takes you to its `LIST` view.

- ### [pki/action](../../app/models/pki/action.js)

  This model is used to perform different `POST` requests that receive similar parameters but don't create a single item (which would be a record in Ember data). These various actions may create multiple items that contain different attributes than those submitted in the `POST` request. For example:

  > - `POST pki/generate/root/:type` creates a new self-signed CA certificate (an issuer) and private key, which is only returned if `type = exported`
  > - `POST pki/issuer/:issuer_ref/sign-intermediate` creates a certificate, and returns issuing CA and CA chain data that is only available once

  The `pki/action`[adapter](../../app/adapters/pki/action.js) is used to map the desired action to the corresponding endpoint, and the `pki/action` [serializer](../../app/serializers/pki/action.js) includes logic to send the relevant attributes. The following PKI workflows use this model:

  - [Root generation and rotation](https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-root)
  - [Import CA cert and keys](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-ca-certificates-and-keys)
  - [Generate intermediate CSR](https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-intermediate-csr)
  - [Sign intermediate](https://developer.hashicorp.com/vault/api-docs/secret/pki#sign-intermediate)

- ### [pki/certificate/base](../../app/models/pki/certificate/base.js)

  This model is for specific interactions with certificate data. The base model contains attributes that make up a certificate's content. The other models that extend this model [certificate/generate](../../app/models/pki/certificate/generate.js) and [certificate/sign](../../app/models/pki/certificate/sign.js) include additional attributes to perform their relevant requests.

  The `parsedCertificate` attribute is an object that houses all of the parsed certificate data returned by the [parse-pki-cert.js](../../app/utils/parse-pki-cert.js) util.

- ### [pki/tidy](../../app/models/pki/tidy.js)

  This model is used to manage [tidy](https://developer.hashicorp.com/vault/api-docs/secret/pki#tidy) operations in a few different contexts. All of the following endpoints share the same parameters _except_ `enabled` and `interval_duration` which are reserved for auto-tidy operations only.

  > _`pki/tidy-status` does not use an Ember data model because it is read-only_

  - `POST pki/tidy` - perform a single, manual tidy operation
  - `POST pki/config/auto-tidy` - set configuration for automating the tidy process
  - `GET pki/config/auto-tidy` - read auto-tidy configuration settings

  The auto-tidy config is the only data that persists so `findRecord` and `updateRecord` in the `pki/tidy.js` [adapter](../../app/adapters/pki/tidy.js) only interact with the `/config/auto-tidy` endpoint. For each manual tidy operation, a new record is created so on `save()` the model uses the `createRecord` method which only ever uses the `/tidy` endpoint.

> _The following models more closely follow a CRUD pattern:_

- ### [pki/issuer](../../app/models/pki/issuer.js)

  > _Issuers are created by the `pki/action` model by either [importing a CA](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-ca-certificates-and-keys) or [generating a root](https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-root)_

  - [update](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-issuer-certificate)
  - [read](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-issuer-certificate)
  - [list](https://developer.hashicorp.com/vault/api-docs/secret/pki#list-issuers)

- ### [pki/role](../../app/models/pki/role.js)

  - [create/update](https://developer.hashicorp.com/vault/api-docs/secret/pki#create-update-role)
  - [read](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-role)
  - [list](https://developer.hashicorp.com/vault/api-docs/secret/pki#list-roles)

- ### [pki/key](../../app/models/pki/key.js)

  - `CREATE` has two options:
    - [generate](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-ca-certificates-and-keys)
    - [import](https://developer.hashicorp.com/vault/api-docs/secret/pki#import-key)
  - [read](https://developer.hashicorp.com/vault/api-docs/secret/pki#read-key)
  - [list](https://developer.hashicorp.com/vault/api-docs/secret/pki#list-keys)
