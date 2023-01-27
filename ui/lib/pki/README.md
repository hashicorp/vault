# Vault PKI

Welcome to the Vault PKI (Ember) Engine! Below is an overview of PKI and resources for how to get started working within this engine.

## About PKI

> Public Key Infrastructure (PKI) is a system of processes, technologies, and policies that allows you to encrypt and sign data. (source: [digicert.com](https://www.digicert.com/what-is-pki))

The [Vault PKI Secrets Engine](https://developer.hashicorp.com/vault/api-docs/secret/pki) allows security engineers to [create a chain of PKI certificates](https://developer.hashicorp.com/vault/tutorials/secrets-management/pki-engine) much easier than they would with traditional workflows.

## About the UI engine

If you couldn't tell from the documentation above, PKI is _complex_. As such, the data doesn't map cleanly to a CRUD model and so the first thing you might notice is that the models and adapters for PKI (which [live in the main app](https://ember-engines.com/docs/addons#using-ember-data), not the engine) have some custom logic that differentiate it from most other secret engines. Below are the model

### pki/key

TBD

### pki/role

TBD

### pki/issuer

TBD

### pki/certificate/\*

TBD

### pki/action

TBD
