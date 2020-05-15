import RoleEdit from './role-edit';

export default RoleEdit.extend({
  init() {
    this._super(...arguments);
    this.set('backendType', 'ssh');
  },

  actions: {
    updateTtl(path, val) {
      const model = this.get('model');
      let valueToSet = val.enabled === true ? `${val.seconds}s` : undefined;
      model.set(path, valueToSet);
    },
  },
});
