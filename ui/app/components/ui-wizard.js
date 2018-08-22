import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  auth: inject.service(),
  currentState: computed.alias('wizard.currentState'),
  featureState: computed.alias('wizard.featureState'),
  featureComponent: computed.alias('wizard.featureComponent'),
  tutorialComponent: computed.alias('wizard.tutorialComponent'),
  componentState: computed.alias('wizard.componentState'),

  actions: {
    dismissWizard() {
      this.get('wizard').transitionTutorialMachine(this.get('currentState'), 'DISMISS');
    },

    advanceWizard() {
      let event = this.get('auth.currentToken') ? 'AUTH' : 'CONTINUE';
      this.get('wizard').transitionTutorialMachine(this.get('currentState'), event);
    },

    advanceFeature() {
      this.get('wizard').transitionFeatureMachine(this.get('featureState'), 'CONTINUE');
    },

    pauseWizard() {
      this.get('wizard').transitionTutorialMachine(this.get('currentState'), 'PAUSE');
    },
  },
});
