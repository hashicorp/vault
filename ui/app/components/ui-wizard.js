import { inject as service } from '@ember/service';
import { alias, or } from '@ember/object/computed';
import Component from '@ember/component';
import { matchesState } from 'xstate';

export default Component.extend({
  classNames: ['ui-wizard-container'],
  wizard: service(),
  auth: service(),

  shouldRender: or('auth.currentToken', 'wizard.showWhenUnauthenticated'),
  currentState: alias('wizard.currentState'),
  featureState: alias('wizard.featureState'),
  featureComponent: alias('wizard.featureComponent'),
  tutorialComponent: alias('wizard.tutorialComponent'),
  componentState: alias('wizard.componentState'),
  nextFeature: alias('wizard.nextFeature'),
  nextStep: alias('wizard.nextStep'),

  actions: {
    dismissWizard() {
      this.wizard.transitionTutorialMachine(this.currentState, 'DISMISS');
    },

    advanceWizard() {
      let inInit = matchesState('init', this.wizard.currentState);
      let event = inInit ? this.wizard.initEvent || 'CONTINUE' : 'CONTINUE';
      this.wizard.transitionTutorialMachine(this.currentState, event);
    },

    advanceFeature() {
      this.wizard.transitionFeatureMachine(this.featureState, 'CONTINUE');
    },

    finishFeature() {
      this.wizard.transitionFeatureMachine(this.featureState, 'DONE');
    },

    repeatStep() {
      this.wizard.transitionFeatureMachine(this.featureState, 'REPEAT', this.componentState);
    },

    resetFeature() {
      this.wizard.transitionFeatureMachine(this.featureState, 'RESET', this.componentState);
    },

    pauseWizard() {
      this.wizard.transitionTutorialMachine(this.currentState, 'PAUSE');
    },
  },
});
