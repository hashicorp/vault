import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  currentState: computed.alias('wizard.currentState'),
  isActive: computed('currentState', () => this.get('currentState') === 'active'),
  isDismissed: computed('currentState', () => this.get('currentState') === 'dismissed'),

  dismissWizard() {
    this.get('wizard').transitionMachine(this.get('currentState'), 'DISMISS');
  },
});
