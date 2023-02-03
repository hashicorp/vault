import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiRolesCreateRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store.createRecord('pki/role', {
      backend: this.secretMountPath.currentPath,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.secretMountPath.currentPath || 'pki';
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: backend, route: 'overview' },
      { label: 'roles', route: 'roles.index' },
      { label: 'create' },
    ];
  }
}
