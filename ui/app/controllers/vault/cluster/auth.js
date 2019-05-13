import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller, { inject as controller } from '@ember/controller';
import { task, timeout } from 'ember-concurrency';

export default Controller.extend({
  vaultController: controller('vault'),
  clusterController: controller('vault.cluster'),
  namespaceService: service('namespace'),
  namespaceQueryParam: alias('clusterController.namespaceQueryParam'),
  queryParams: [{ authMethod: 'with' }],
  wrappedToken: alias('vaultController.wrappedToken'),
  authMethod: '',
  redirectTo: null,

  updateNamespace: task(function*(value) {
    // debounce
    yield timeout(500);
    this.get('namespaceService').setNamespace(value, true);
    this.set('namespaceQueryParam', value);
  }).restartable(),
});
