import PkiRolesIndexRoute from '../index';

export default class RolesRoleDetailsRoute extends PkiRolesIndexRoute {
  model() {
    const { role } = this.paramsFor('roles/role');
    return this.store.queryRecord('pki/role', {
      backend: this.secretMountPath.currentPath,
      id: role,
    });
  }
}
