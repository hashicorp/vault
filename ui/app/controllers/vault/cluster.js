import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller from '@ember/controller';
import { observer, computed } from '@ember/object';
export default Controller.extend({
  auth: service(),
  store: service(),
  media: service(),
  router: service(),
  permissions: service(),
  namespaceService: service('namespace'),

  vaultVersion: service('version'),
  console: service(),

  queryParams: [
    {
      namespaceQueryParam: {
        scope: 'controller',
        as: 'namespace',
      },
    },
  ],

  namespaceQueryParam: '',

  onQPChange: observer('namespaceQueryParam', function() {
    this.get('namespaceService').setNamespace(this.get('namespaceQueryParam'));
  }),

  consoleOpen: alias('console.isOpen'),

  activeCluster: computed('auth.activeCluster', function() {
    return this.get('store').peekRecord('cluster', this.get('auth.activeCluster'));
  }),

  activeClusterName: computed('activeCluster', function() {
    const activeCluster = this.get('activeCluster');
    return activeCluster ? activeCluster.get('name') : null;
  }),

  showNav: computed(
    'router.currentRouteName',
    'activeClusterName',
    'auth.currentToken',
    'activeCluster.{dr.isSecondary,needsInit,sealed}',
    function() {
      if (
        this.get('activeCluster.dr.isSecondary') ||
        this.get('activeCluster.needsInit') ||
        this.get('activeCluster.sealed')
      ) {
        return false;
      }
      if (
        this.activeClusterName &&
        this.auth.currentToken &&
        this.router.currentRouteName !== 'vault.cluster.auth'
      ) {
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
