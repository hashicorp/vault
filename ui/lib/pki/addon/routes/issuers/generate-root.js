import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiIssuersGenerateRootRoute extends Route {
  @service secretMountPath;
  @service store;

  model() {
    return this.store.createRecord('pki/action');
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.secretMountPath.currentPath || 'pki';
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: backend, route: 'overview' },
      { label: 'generate root' },
    ];
  }
}
