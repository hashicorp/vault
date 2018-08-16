import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  currentState: computed.alias('wizard.currentState'),
  featureState: computed.alias('wizard.featureState'),
  componentState: computed.alias('wizard.componentState'),
  initializationSteps: [
    {
      key: 'setup',
      title: 'Setting up your master keys',
      description:
        'This is the very first step of setting up a Vault server. Vault comes with an important security feature called "seal", which lets you shut down and secure your Vault installation if there is a security breach. This is the default state of Vault, so it is currently sealed since it was just installed. To unseal the vault, you will need to provide the key(s) that you generate here.',
      docs: {
        link: 'https://www.vaultproject.io/intro/getting-started/deploy.html#initializing-the-vault',
        text: 'Docs: Initialization',
      },
    },
    {
      key: 'save',
      title: 'Saving your keys',
      description:
        "Now that Vault is initialized, you'll want to save your root token and master key portions in a safe place. Distribute your keys to responsible people on your team. If these keys are lost, you may not be able to access your data again. Keep them safe!",
      docs: {
        link: 'https://www.vaultproject.io/intro/getting-started/deploy.html#initializing-the-vault',
        text: 'Docs: Initialization',
      },
    },
    {
      key: 'unseal',
      title: 'Unsealing your vault',
      description: '',
      docs: {
        link: 'https://www.vaultproject.io/docs/concepts/seal.html',
        text: 'Docs: Unseal',
      },
    },
    {
      key: 'login',
      title: 'Sign in to your vault',
      description: 'Congrats! Now that your vault is all set up you can sign in!',
    },
  ],
  inSetup: computed('currentState', function() {
    return this.get('currentState').indexOf('init.setup') > 0;
  }),
  isSaving: computed('currentState', function() {
    return this.get('currentState').indexOf('init.save') > 0;
  }),
  isUnsealing: computed('currentState', function() {
    return this.get('currentState').indexOf('init.unseal') > 0;
  }),
  inLogin: computed('currentState', function() {
    return this.get('currentState').indexOf('init.lgin') > 0;
  }),
  currentStep: computed('currentState', 'componentState', function() {
    const stateParts = this.get('currentState').split('.');
    let currentStep = this.get('initializationSteps')
      .filter(step => step.key === stateParts[stateParts.length - 1])
      .objectAt(0);
    if (this.get('isUnsealing')) {
      if (this.get('componentState')) {
        const keyWord = this.get('componentState.threshold') > 1 ? 'keys' : 'key';
        const providedWord = this.get('componentState.progress') > 1 ? 'have' : 'has';
        const keysLeft = this.get('componentState.threshold') - this.get('componentState.progress');
        const leftWord = keysLeft > 1 ? 'keys' : 'key';
        Ember.set(
          currentStep,
          'description',
          `Now we will provide the ${keyWord} that you copied or downloaded to unseal the vault so that we can get started using it. You'll need ${this.get(
            'componentState.threshold'
          )} ${keyWord} total, and ${this.get(
            'componentState.progress'
          )} ${providedWord} already been provided. Please provide ${keysLeft} more ${leftWord} to unseal.`
        );
      }
    }
    return currentStep;
  }),
});
