import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiConfigurationCreateRoute extends Route {
  @service secretMountPath;
  @service store;
  @service pathHelp;

  beforeModel() {
    return this.pathHelp.getNewModel('pki/urls', this.secretMountPath.currentPath);
  }

  model() {
    return {
      config: this.store.createRecord('pki/action'),
      urls: this.store.createRecord('pki/urls'),
    };
  }

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
