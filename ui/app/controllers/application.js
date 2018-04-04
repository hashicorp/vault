import Ember from 'ember';
import config from '../config/environment';

export default Ember.Controller.extend({
  env: config.environment,
  auth: Ember.inject.service(),
  vaultVersion: Ember.inject.service('version'),
  activeCluster: Ember.computed('auth.activeCluster', function() {
    return this.store.peekRecord('cluster', this.get('auth.activeCluster'));
  }),
  activeClusterName: Ember.computed('auth.activeCluster', function() {
    const activeCluster = this.store.peekRecord('cluster', this.get('auth.activeCluster'));
    return activeCluster ? activeCluster.get('name') : null;
  }),
  showNav: Ember.computed(
    'activeClusterName',
    'auth.currentToken',
    'activeCluster.dr.isSecondary',
    'activeCluster.{needsInit,sealed}',
    function() {
      if (
        this.get('activeCluster.dr.isSecondary') ||
        this.get('activeCluster.needsInit') ||
        this.get('activeCluster.sealed')
      ) {
        return false;
      }
      if (this.get('activeClusterName') && this.get('auth.currentToken')) {
        return true;
      }
    }
  ),
});
