import RoleEdit from '../role-edit';

export default RoleEdit.extend({
  actions: {
    delete() {
      this.model.save({ adapterOptions: { method: 'revoke' } });
    },
  },
});
