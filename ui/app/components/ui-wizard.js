import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { matchesState } from 'xstate';

export default Component.extend({
  classNames: ['ui-wizard-container'],
  wizard: service(),
  auth: service(),

  shouldRender: computed('wizard.showWhenUnauthenticated', 'auth.currentToken', function() {
    return this.get('auth.currentToken') || this.get('wizard.showWhenUnauthenticated');
  }),
  currentState: alias('wizard.currentState'),
  featureState: alias('wizard.featureState'),
  featureComponent: alias('wizard.featureComponent'),
  tutorialComponent: alias('wizard.tutorialComponent'),
  componentState: alias('wizard.componentState'),
  nextFeature: alias('wizard.nextFeature'),
  nextStep: alias('wizard.nextStep'),

  actions: {
    dismissWizard() {
      this.get('wizard').transitionTutorialMachine(this.get('currentState'), 'DISMISS');
    },

    advanceWizard() {
      let inInit = matchesState('init', this.get('wizard.currentState'));
      let event = inInit ? this.get('wizard.initEvent') || 'CONTINUE' : 'CONTINUE';
      this.get('wizard').transitionTutorialMachine(this.get('currentState'), event);
    },

    advanceFeature() {
      this.get('wizard').transitionFeatureMachine(this.get('featureState'), 'CONTINUE');
    },

    finishFeature() {
      this.get('wizard').transitionFeatureMachine(this.get('featureState'), 'DONE');
    },

    repeatStep() {
      this.get('wizard').transitionFeatureMachine(
        this.get('featureState'),
        'REPEAT',
        this.get('componentState')
      );
    },

    resetFeature() {
      this.get('wizard').transitionFeatureMachine(
        this.get('featureState'),
        'RESET',
        this.get('componentState')
      );
    },

    pauseWizard() {
      this.get('wizard').transitionTutorialMachine(this.get('currentState'), 'PAUSE');
    },
  },
});
