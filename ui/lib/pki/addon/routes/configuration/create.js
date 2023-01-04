import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiConfigurationCreateRoute extends Route {
  @service secretMountPath;

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.secretMountPath.currentPath || 'pki';
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: backend, route: 'overview' },
      { label: 'configure' },
    ];
  }
}
