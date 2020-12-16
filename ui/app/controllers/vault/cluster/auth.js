import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller, { inject as controller } from '@ember/controller';
import { task, timeout } from 'ember-concurrency';

export default Controller.extend({
  vaultController: controller('vault'),
  clusterController: controller('vault.cluster'),
  namespaceService: service('namespace'),
  configService: service('config'),
  namespaceQueryParam: alias('clusterController.namespaceQueryParam'),
  queryParams: [{ authMethod: 'with' }],
  wrappedToken: alias('vaultController.wrappedToken'),
  authMethod: '',
  redirectTo: alias('vaultController.redirectTo'),
  managedNamespaceRoot: alias('configService.managedNamespaceRoot'),

  get managedNamespaceChild() {
    let fullParam = this.namespaceQueryParam;
    let split = fullParam.split('/');
    if (split.length > 1) {
      split.shift();
      console.log(split);
      return `/${split.join('/')}`;
    }
    return '';
  },

  updateManagedNamespace: task(function*(value) {
    yield timeout(500);
    // TODO: Move this to shared fn
    const newNamespace = `${this.managedNamespaceRoot}${value}`;
    this.namespaceService.setNamespace(newNamespace, true);
    this.set('namespaceQueryParam', newNamespace);
  }).restartable(),

  updateNamespace: task(function*(value) {
    // debounce
    yield timeout(500);
    this.namespaceService.setNamespace(value, true);
    this.set('namespaceQueryParam', value);
  }).restartable(),
});
