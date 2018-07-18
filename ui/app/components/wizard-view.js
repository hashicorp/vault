import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),

  currentState: computed.alias('wizard.currentState'),
});
