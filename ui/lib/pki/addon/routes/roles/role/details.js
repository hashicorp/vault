import PkiRolesIndexRoute from '../index';

export default class RolesRoleDetailsRoute extends PkiRolesIndexRoute {
  model() {
    const { id } = this.paramsFor('roles/role');
    return this.store.queryRecord('pki/role', {
      backend: this.secretMountPath.currentPath,
      id,
    });
  }
}
