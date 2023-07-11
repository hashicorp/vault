import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretEditRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // TODO add model. Needs to include a queryParam for version.
    const backend = this.secretMountPath.get();
    const { name } = this.paramsFor('secret');
    return hash({
      path: name,
      backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'list', linkExternal: true },
      { label: resolvedModel.backend, route: 'secret' },
      { label: resolvedModel.path, route: 'secret.details', model: resolvedModel.path },
      { label: 'edit' },
    ];
  }
}
