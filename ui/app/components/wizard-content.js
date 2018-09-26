import { inject as service } from '@ember/service';
import Component from '@ember/component';
export default Component.extend({
  wizard: service(),
  classNames: ['ui-wizard'],
  glyph: null,
  headerText: null,
  actions: {
    dismissWizard() {
      this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'DISMISS');
    },
  },
});
