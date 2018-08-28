import RoleEdit from './role-edit';

export default RoleEdit.extend({
  init() {
    this._super(...arguments);
    this.set('backendType', 'pki');
  },
});
