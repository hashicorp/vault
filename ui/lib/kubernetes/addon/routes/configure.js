import FetchConfigRoute from './fetch-config';

export default class KubernetesConfigureRoute extends FetchConfigRoute {
  async model() {
    const backend = this.secretMountPath.get();
    return this.configModel || this.store.createRecord('kubernetes/config', { backend });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'configure' },
    ];
  }
}
