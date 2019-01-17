import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task, timeout } from 'ember-concurrency';

export default Component.extend({
  store: service(),
  onRoleChange() {},
  selectedAuthPath: null,
  roleName: null,

  fetchRole: task(function*(roleName) {
    this.set('roleName', roleName);
    // debounce
    yield timeout(500);
    let path = this.selectedAuthPath || 'jwt';
    let id = JSON.stringify([path, roleName]);
    let role = yield this.store.findRecord('role-jwt', id);
    this.onRoleChange(role);
  })
    .restartable()
    .on('init'),
});
