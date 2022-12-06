import { withConfirmLeave } from 'core/decorators/confirm-leave';
import PkiRolesIndexRoute from '../index';

@withConfirmLeave()
export default class PkiRoleEditRoute extends PkiRolesIndexRoute {
  model() {
    const { role } = this.paramsFor('roles/role');
    return this.store.queryRecord('pki/role', {
      backend: this.secretMountPath.currentPath,
      id: role,
    });
  }
}
