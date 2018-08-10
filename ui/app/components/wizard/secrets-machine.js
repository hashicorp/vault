import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  currentState: computed.alias('wizard.currentState'),
  featureState: computed.alias('wizard.featureState'),
  selectedEngine: computed.alias('wizard.potentialSelection'),
  secretsEngines: [
    {
      key: 'secrets.aws',
      name: 'AWS',
      description:
        'The AWS secrets engine generates AWS access credentials dynamically based on IAM policies. This generally makes working with AWS IAM easier, since it does not involve clicking in the web UI. Additionally, the process is codified and mapped to internal auth methods (such as LDAP). The AWS IAM credentials are time-based and are automatically revoked when the Vault lease expires.',
    },
    {
      key: 'secrets.ch',
      name: 'Cubbyhole',
      description:
        "The cubbyhole secrets engine is used to store arbitrary secrets within the configured physical storage for Vault namespaced to a token. In cubbyhole, paths are scoped per token. No token can access another token's cubbyhole. When the token expires, its cubbyhole is destroyed.",
    },
    {
      key: 'secrets.kv',
      name: 'Key/Value',
      description:
        'The kv secrets engine is used to store arbitrary secrets within the configured physical storage for Vault. This backend can be run in one of two modes. It can be a generic Key-Value store that stores one value for a key. Versioning can be enabled and a configurable number of versions for each key will be stored.',
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
