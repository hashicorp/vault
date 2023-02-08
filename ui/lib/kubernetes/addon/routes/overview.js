import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'kubernetes/decorators/fetch-config';
import { hash } from 'rsvp';

@withConfig()
export default class KubernetesOverviewRoute extends Route {
  @service store;
  @service secretMountPath;

  async model() {
    const backend = this.secretMountPath.get();
    return hash({
      promptConfig: this.promptConfig,
      backend: this.modelFor('application'),
      roles: this.store.query('kubernetes/role', { backend }).catch(() => []),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend.id },
    ];
  }
}
