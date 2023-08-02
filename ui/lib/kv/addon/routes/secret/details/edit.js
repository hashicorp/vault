import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KvSecretEditRoute extends Route {
  @service store;

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
      { label: resolvedModel.path, route: 'secret.details', model: resolvedModel.path },
      { label: 'edit' },
    ];
  }
}
