import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  currentState: computed.alias('wizard.currentState'),
  featureState: computed.alias('wizard.featureState'),
  selectedEngine: computed.alias('wizard.potentialSelection'),
  secretsEngines: [
    {
      key: 'secrets.ad',
      name: 'Active Directory',
      description:
        'The AD secrets engine rotates AD passwords dynamically, and is designed for a high-load environment where many instances may be accessing a shared password simultaneously.',
    },
    {
      key: 'secrets.aws',
      name: 'AWS',
      description:
        'The AWS secrets engine generates AWS access credentials dynamically based on IAM policies. This generally makes working with AWS IAM easier, since it does not involve clicking in the web UI. Additionally, the process is codified and mapped to internal auth methods (such as LDAP). The AWS IAM credentials are time-based and are automatically revoked when the Vault lease expires.',
    },
    {
      key: 'secrets.consul',
      name: 'Consul',
      description:
        'The Consul secrets engine generates Consul API tokens dynamically based on Consul ACL policies',
    },
    {
      key: 'secrets.ch',
      name: 'Cubbyhole',
      description:
        "The cubbyhole secrets engine is used to store arbitrary secrets within the configured physical storage for Vault namespaced to a token. In cubbyhole, paths are scoped per token. No token can access another token's cubbyhole. When the token expires, its cubbyhole is destroyed.",
    },
    {
      key: 'secrets.gcp',
      name: 'Google Cloud',
      description:
        'The Google Cloud Vault secrets engine dynamically generates Google Cloud service account keys and OAuth tokens based on IAM policies. This enables users to gain access to Google Cloud resources without needing to create or manage a dedicated service account.',
    },
    {
      key: 'secrets.kv',
      name: 'Key/Value',
      description:
        'The kv secrets engine is used to store arbitrary secrets within the configured physical storage for Vault. This backend can be run in one of two modes. It can be a generic Key-Value store that stores one value for a key. Versioning can be enabled and a configurable number of versions for each key will be stored.',
    },
    {
      key: 'secrets.nomad',
      name: 'Nomad',
      description:
        'The Nomad secret backend for Vault generates Nomad API tokens dynamically based on pre-existing Nomad ACL policies.',
    },
    {
      key: 'secrets.pki',
      name: 'PKI',
      description:
        'The PKI secrets engine generates dynamic X.509 certificates. With this secrets engine, services can get certificates without going through the usual manual process of generating a private key and CSR, submitting to a CA, and waiting for a verification and signing process to complete.',
    },
    {
      key: 'secrets.rabbitmq',
      name: 'RabbitMQ',
      description:
        'The RabbitMQ secrets engine generates user credentials dynamically based on configured permissions and virtual hosts. This means that services that need to access a virtual host no longer need to hardcode credentials.',
    },
    {
      key: 'secrets.ssh',
      name: 'SSH',
      description:
        'The Vault SSH secrets engine provides secure authentication and authorization for access to machines via the SSH protocol. The Vault SSH secrets engine helps manage access to machine infrastructure, providing several ways to issue SSH credentials.',
    },
    {
      key: 'secrets.totp',
      name: 'TOTP',
      description: 'The TOTP secrets engine generates time-based credentials according to the TOTP standard.',
    },
    {
      key: 'secrets.transit',
      name: 'Transit',
      description:
        'The transit secrets engine handles cryptographic functions on data in-transit. Vault does not store the data sent to the secrets engine. It can also be viewed as "cryptography as a service" or "encryption as a service".',
    },
  ],
  isIdle: computed('featureState', function() {
    return this.get('featureState') === 'idle';
  }),
  currentStepName: computed('featureState', function() {
    if (this.get('isIdle')) {
      return 'Enabling a secrets engine';
    }
  }),
  currentStepText: computed('featureState', function() {
    if (this.get('isIdle')) {
      return "Vault is all about managing secrets, so let's set up your first secrets engine. You can use a static engine to store your secrets locally in Vault, or connect to a cloud backend with one of the dynamic engines.";
    }
  }),
  engineDescription: computed('selectedEngine', function() {
    return this.get('secretsEngines')
      .filter(engine => engine.key === `secrets.${this.get('selectedEngine')}`)
      .objectAt(0).description;
  }),
  engineName: computed('selectedEngine', function() {
    return this.get('secretsEngines')
      .filter(engine => engine.key === `secrets.${this.get('selectedEngine')}`)
      .objectAt(0).name;
  }),

  dismissWizard() {
    this.get('wizard').transitionTutorialMachine(this.get('currentState'), 'DISMISS');
  },

  advanceWizard() {
    this.get('wizard').transitionFeatureMachine(this.get('featureState'), 'CONTINUE');
  },
});
