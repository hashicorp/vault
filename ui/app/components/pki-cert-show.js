import RoleEdit from './role-edit';
// take whatever data model passed in
// node forge stuff here in getter

export default RoleEdit.extend({
  actions: {
    delete() {
      this.model.save({ adapterOptions: { method: 'revoke' } });
    },
  },
});
