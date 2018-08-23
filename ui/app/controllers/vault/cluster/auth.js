import Ember from 'ember';
import { task, timeout } from 'ember-concurrency';

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

  updateNamespace: task(function*(value) {
    yield timeout(200);
    this.get('namespaceService').setNamespace(value, true);
    this.set('namespaceQueryParam', value);
  }).restartable(),
});
