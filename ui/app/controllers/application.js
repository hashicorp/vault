import Ember from 'ember';
import config from '../config/environment';

const { computed, inject } = Ember;
export default Ember.Controller.extend({
  env: config.environment,
  auth: inject.service(),
  vaultVersion: inject.service('version'),
  console: inject.service(),
  consoleOpen: computed.alias('console.isOpen'),
  activeCluster: computed('auth.activeCluster', function() {
    return this.store.peekRecord('cluster', this.get('auth.activeCluster'));
  }),
  activeClusterName: computed('auth.activeCluster', function() {
    const activeCluster = this.store.peekRecord('cluster', this.get('auth.activeCluster'));
    return activeCluster ? activeCluster.get('name') : null;
  }),
  showNav: computed(
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
  actions: {
    toggleConsole() {
      this.toggleProperty('consoleOpen');
    },
  },
});
