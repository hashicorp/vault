import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  currentState: computed.alias('wizard.currentState'),
  featureState: computed.alias('wizard.featureState'),
  currentMachine: computed('wizard.featureList', function() {
    if (this.get('wizard.featureList') !== null) {
      let machine = this.get('wizard.featureList').toArray().objectAt(0);
      return machine.charAt(0).toUpperCase() + machine.slice(1);
    }
    return 'None';
  }),
  isActive: computed('currentState', function() {
    return this.get('currentState').indexOf('active') == 0;
  }),
  isDismissed: computed('currentState', function() {
    return this.get('currentState') === 'dismissed';
  }),
  isSelecting: computed('currentState', 'isActive', function() {
    return this.get('isActive') && this.get('currentState').indexOf('select') > 0;
  }),

  dismissWizard() {
    this.get('wizard').transitionTutorialMachine(this.get('currentState'), 'DISMISS');
  },

  advanceWizard() {
    this.get('wizard').transitionTutorialMachine(this.get('currentState'), 'CONTINUE');
  },
});
