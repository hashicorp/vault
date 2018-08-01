import Ember from 'ember';

const { inject, computed, Controller } = Ember;
export default Controller.extend({
  vaultController: inject.controller('vault'),
  clusterController: inject.controller('vault.cluster'),
  namespaceService: inject.service('namespace'),
  namespaceQueryParam: computed.alias('clusterController.namespaceQueryParam'),
  queryParams: [{ authMethod: 'with' }],
  wrappedToken: computed.alias('vaultController.wrappedToken'),
  authMethod: '',
  redirectTo: null,

  actions: {
    updateNamespace(event) {
      let { value } = event.target;

      this.get('namespaceService').setNamespace(value, true);
      this.set('namespaceQueryParam', value);
    },
  },
});
