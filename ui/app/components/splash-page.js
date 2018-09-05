import Ember from 'ember';

const { computed, inject } = Ember;

export default Ember.Component.extend({
  version: inject.service(),
  auth: inject.service(),
  store: inject.service(),
  tagName: '',

  activeCluster: computed('auth.activeCluster', function() {
    return this.get('store').peekRecord('cluster', this.get('auth.activeCluster'));
  }),
});
