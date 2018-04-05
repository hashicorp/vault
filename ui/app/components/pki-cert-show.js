import RoleEdit from './role-edit';

export default RoleEdit.extend({
  actions: {
    delete() {
      this.get('model').save({ adapterOptions: { method: 'revoke' } });
    },
  },
});
