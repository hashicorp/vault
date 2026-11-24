# PKI Ember Engine

This Ember Engine houses both **Internal** and **External** Vault PKI secrets engines.

### What is PKI?

**Public Key Infrastructure (PKI)** is a system of processes, technologies, and policies that allows you to encrypt and sign data.

**Certificate Authority (CA)** is a trusted organization that authenticates the digital identities of entities, from individuals to servers to devices. 

## Internal PKI vs External PKI

### Internal PKI (Original PKI Secrets Engine)

**Vault acts as the Certificate Authority**

- Vault generates and stores root or intermediate CA keys
- Vault is responsible for issuing, signing, and revoking certificates
- **Vault is the source of trust** and manages the entire certificate lifecycle
- Available since early Vault versions

**Use Cases:**
- Internal infrastructure certificates
- Development and testing environments
- Organizations that want full control over their PKI
- Private certificate issuance

### External/Public PKI (Added in Vault 2.0.0)

**Vault acts as a broker with external ACME-compatible CAs**

- Integrates with trusted third-party certificate providers (e.g., [GlobalSign](https://www.hashicorp.com/en/partners/tech/globalsign#all), [Sectigo](https://www.hashicorp.com/en/partners/tech/sectigo#all), [DigiCert](https://www.hashicorp.com/en/partners/tech/digicert#all))
- Automates certificate issuance from external CAs
- **External CA remains the source of trust**
- Vault serves as the secure automation and distribution layer

**Use Cases:**
- Publicly trusted certificates
- Compliance requirements for external CAs
- Organizations that need certificates trusted by browsers and operating systems
- Automated certificate lifecycle management with external providers

## Additional Resources

- [Internal PKI API Docs](https://developer.hashicorp.com/vault/api-docs/secret/pk)
- [External PKI API Docs](https://developer.hashicorp.com/vault/api-docs/secret/pki-external-ca)
- [ACME Protocol](https://datatracker.ietf.org/doc/html/rfc8555)
- [Public Key Infrastructure Overview](https://www.digicert.com/what-is-pki)
- [Certificate Authority Overview](https://www.cyberark.com/what-is/certificate-authority)
