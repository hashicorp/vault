import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from '../decorators/fetch-config';

@withConfig()
export default class KubernetesConfigureRoute extends Route {
  @service store;
  @service secretMountPath;

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
