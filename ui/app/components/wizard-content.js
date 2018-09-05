import Ember from 'ember';

const { Component, inject } = Ember;
export default Component.extend({
  wizard: inject.service(),
  classNames: ['ui-wizard'],
  glyph: null,
  headerText: null,
  actions: {
    dismissWizard() {
      this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'DISMISS');
    },
  },
});
