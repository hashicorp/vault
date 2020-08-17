import RoleEdit from './role-edit';
import transform from '../models/transform';

export default RoleEdit.extend({
  init() {
    this._super(...arguments);
    this.set('backendType', 'transform');
  },
});
