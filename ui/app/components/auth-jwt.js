import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task, timeout } from 'ember-concurrency';

export default Component.extend({
  store: service(),
  onRoleChange() {},
  selectedAuthPath: null,
  roleName: null,
  fetchRole: task(function*(roleName) {
    // debounce
    this.set('roleName', roleName);
    yield timeout(500);
    let id = JSON.stringify([this.selectedAuthPath, roleName]);
    let role = yield this.store.findRecord('role-jwt', id);
    this.onRoleChange(role);
  }).restartable(),
});
