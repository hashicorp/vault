import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KubernetesRoleDetailsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.get();
    const { name } = this.paramsFor('roles.role');
    return this.store.queryRecord('kubernetes/role', { backend, name });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'roles', route: 'roles' },
      { label: resolvedModel.name },
    ];
  }
}
