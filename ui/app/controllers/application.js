import Ember from 'ember';
import config from '../config/environment';

const { Controller, computed, inject } = Ember;
export default Controller.extend({
  env: config.environment,
  auth: inject.service(),
  store: inject.service(),
  activeCluster: computed('auth.activeCluster', function() {
    return this.get('store').peekRecord('cluster', this.get('auth.activeCluster'));
  }),
  activeClusterName: computed('activeCluster', function() {
    const activeCluster = this.get('activeCluster');
    return activeCluster ? activeCluster.get('name') : null;
  }),
});
